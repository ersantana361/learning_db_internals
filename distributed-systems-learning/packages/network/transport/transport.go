package transport

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MessageType identifies the type of message
type MessageType string

// Envelope wraps a message with routing metadata
type Envelope struct {
	ID          string                 `json:"id"`
	From        string                 `json:"from"`
	To          string                 `json:"to"`
	Type        MessageType            `json:"type"`
	Payload     interface{}            `json:"payload"`
	SentAt      time.Time              `json:"sentAt"`
	ReceivedAt  time.Time              `json:"receivedAt,omitempty"`
	LamportTime uint64                 `json:"lamportTime,omitempty"`
	VectorClock map[string]uint64      `json:"vectorClock,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewEnvelope creates a new message envelope
func NewEnvelope(from, to string, msgType MessageType, payload interface{}) *Envelope {
	return &Envelope{
		ID:       uuid.New().String(),
		From:     from,
		To:       to,
		Type:     msgType,
		Payload:  payload,
		SentAt:   time.Now(),
		Metadata: make(map[string]interface{}),
	}
}

// DeliveryHandler is called when a message is delivered
type DeliveryHandler func(env *Envelope)

// DropHandler is called when a message is dropped
type DropHandler func(env *Envelope, reason string)

// Transport defines the network transport interface
type Transport interface {
	// Send sends a message (may fail depending on implementation)
	Send(ctx context.Context, env *Envelope) error

	// RegisterHandler registers a handler for incoming messages
	RegisterHandler(nodeID string, handler DeliveryHandler)

	// Configure failure characteristics
	SetLatency(min, max time.Duration)
	SetPacketLoss(probability float64)
	SetPartition(from, to string, enabled bool)
	ClearPartition(from, to string)
	ClearAllPartitions()

	// Event handlers
	OnDrop(handler DropHandler)

	// Close shuts down the transport
	Close()
}

// NetworkTransport implements Transport with configurable reliability
type NetworkTransport struct {
	mu sync.RWMutex

	handlers   map[string]DeliveryHandler
	dropHandler DropHandler

	// Network characteristics
	minLatency   time.Duration
	maxLatency   time.Duration
	packetLoss   float64 // 0.0 to 1.0

	// Partitions: partitions[from][to] = true means messages from->to are blocked
	partitions map[string]map[string]bool

	// Pending messages (for step mode)
	pending []*pendingMessage

	closed bool
}

type pendingMessage struct {
	env       *Envelope
	deliverAt time.Time
}

// NewNetworkTransport creates a new network transport
func NewNetworkTransport() *NetworkTransport {
	return &NetworkTransport{
		handlers:   make(map[string]DeliveryHandler),
		partitions: make(map[string]map[string]bool),
		minLatency: 0,
		maxLatency: 0,
		packetLoss: 0,
	}
}

// RegisterHandler registers a delivery handler for a node
func (t *NetworkTransport) RegisterHandler(nodeID string, handler DeliveryHandler) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.handlers[nodeID] = handler
}

// OnDrop sets the drop handler
func (t *NetworkTransport) OnDrop(handler DropHandler) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.dropHandler = handler
}

// Send sends a message through the network
func (t *NetworkTransport) Send(ctx context.Context, env *Envelope) error {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return nil
	}

	// Check for partition
	if t.isPartitioned(env.From, env.To) {
		dropHandler := t.dropHandler
		t.mu.RUnlock()
		if dropHandler != nil {
			dropHandler(env, "network_partition")
		}
		return nil
	}

	// Check for packet loss
	if t.packetLoss > 0 && rand.Float64() < t.packetLoss {
		dropHandler := t.dropHandler
		t.mu.RUnlock()
		if dropHandler != nil {
			dropHandler(env, "packet_loss")
		}
		return nil
	}

	handler := t.handlers[env.To]
	minLat := t.minLatency
	maxLat := t.maxLatency
	t.mu.RUnlock()

	if handler == nil {
		return nil // No handler registered
	}

	// Calculate latency
	latency := minLat
	if maxLat > minLat {
		latency = minLat + time.Duration(rand.Int63n(int64(maxLat-minLat)))
	}

	// Deliver with latency
	if latency > 0 {
		go func() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(latency):
				envCopy := *env
				envCopy.ReceivedAt = time.Now()
				handler(&envCopy)
			}
		}()
	} else {
		envCopy := *env
		envCopy.ReceivedAt = time.Now()
		go handler(&envCopy)
	}

	return nil
}

// SetLatency sets the min and max latency for message delivery
func (t *NetworkTransport) SetLatency(min, max time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.minLatency = min
	t.maxLatency = max
}

// SetPacketLoss sets the probability of packet loss (0.0 to 1.0)
func (t *NetworkTransport) SetPacketLoss(probability float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if probability < 0 {
		probability = 0
	}
	if probability > 1 {
		probability = 1
	}
	t.packetLoss = probability
}

// SetPartition creates a network partition between two nodes
func (t *NetworkTransport) SetPartition(from, to string, enabled bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if enabled {
		if t.partitions[from] == nil {
			t.partitions[from] = make(map[string]bool)
		}
		t.partitions[from][to] = true
	} else {
		if t.partitions[from] != nil {
			delete(t.partitions[from], to)
		}
	}
}

// ClearPartition removes a partition between two nodes
func (t *NetworkTransport) ClearPartition(from, to string) {
	t.SetPartition(from, to, false)
}

// ClearAllPartitions removes all network partitions
func (t *NetworkTransport) ClearAllPartitions() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.partitions = make(map[string]map[string]bool)
}

// isPartitioned checks if there's a partition between from and to
func (t *NetworkTransport) isPartitioned(from, to string) bool {
	if t.partitions[from] != nil && t.partitions[from][to] {
		return true
	}
	return false
}

// CreateBidirectionalPartition creates a partition in both directions
func (t *NetworkTransport) CreateBidirectionalPartition(a, b string) {
	t.SetPartition(a, b, true)
	t.SetPartition(b, a, true)
}

// ClearBidirectionalPartition clears a partition in both directions
func (t *NetworkTransport) ClearBidirectionalPartition(a, b string) {
	t.SetPartition(a, b, false)
	t.SetPartition(b, a, false)
}

// Close shuts down the transport
func (t *NetworkTransport) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
}

// GetNetworkStats returns current network configuration
func (t *NetworkTransport) GetNetworkStats() map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	partitionList := make([]map[string]string, 0)
	for from, tos := range t.partitions {
		for to := range tos {
			partitionList = append(partitionList, map[string]string{
				"from": from,
				"to":   to,
			})
		}
	}

	return map[string]interface{}{
		"minLatency":  t.minLatency.String(),
		"maxLatency":  t.maxLatency.String(),
		"packetLoss":  t.packetLoss,
		"partitions":  partitionList,
	}
}
