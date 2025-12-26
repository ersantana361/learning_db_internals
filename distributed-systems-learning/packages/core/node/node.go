package node

import (
	"context"
	"sync"

	"github.com/ersantana/distributed-systems-learning/packages/core/message"
)

// ID uniquely identifies a node
type ID string

// State represents the current state of a node
type State int

const (
	StateRunning State = iota
	StateCrashed
	StatePartitioned
	StateByzantine
)

func (s State) String() string {
	switch s {
	case StateRunning:
		return "running"
	case StateCrashed:
		return "crashed"
	case StatePartitioned:
		return "partitioned"
	case StateByzantine:
		return "byzantine"
	default:
		return "unknown"
	}
}

// Node is the base interface for all distributed nodes
type Node interface {
	// Identity
	ID() ID

	// State management
	State() State
	SetState(state State)

	// Lifecycle
	Start(ctx context.Context) error
	Stop() error

	// Messaging
	Send(to ID, msg message.Message) error
	Receive(env *message.Envelope)
	Inbox() *message.Queue

	// Failure injection hooks
	Crash()
	Recover()

	// Visualization support
	GetVisualizationState() map[string]interface{}
}

// SendFunc is a function type for sending messages
type SendFunc func(from, to ID, msg message.Message) error

// EventEmitter is a function type for emitting events
type EventEmitter func(eventType string, data interface{})

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	mu       sync.RWMutex
	id       ID
	state    State
	inbox    *message.Queue
	sendFunc SendFunc
	emitter  EventEmitter
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewBaseNode creates a new base node
func NewBaseNode(id ID, sendFunc SendFunc, emitter EventEmitter) *BaseNode {
	return &BaseNode{
		id:       id,
		state:    StateRunning,
		inbox:    message.NewQueue(1000),
		sendFunc: sendFunc,
		emitter:  emitter,
	}
}

// ID returns the node's unique identifier
func (n *BaseNode) ID() ID {
	return n.id
}

// State returns the current state of the node
func (n *BaseNode) State() State {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.state
}

// SetState sets the node's state
func (n *BaseNode) SetState(state State) {
	n.mu.Lock()
	defer n.mu.Unlock()
	oldState := n.state
	n.state = state
	if n.emitter != nil {
		n.emitter("node_state_changed", map[string]interface{}{
			"nodeID":   string(n.id),
			"oldState": oldState.String(),
			"newState": state.String(),
		})
	}
}

// Start starts the node
func (n *BaseNode) Start(ctx context.Context) error {
	n.mu.Lock()
	n.ctx, n.cancel = context.WithCancel(ctx)
	n.state = StateRunning
	n.mu.Unlock()
	return nil
}

// Stop stops the node
func (n *BaseNode) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.cancel != nil {
		n.cancel()
	}
	n.inbox.Close()
	return nil
}

// Send sends a message to another node
func (n *BaseNode) Send(to ID, msg message.Message) error {
	n.mu.RLock()
	state := n.state
	sendFunc := n.sendFunc
	n.mu.RUnlock()

	if state != StateRunning {
		return nil // Silently drop if crashed/partitioned
	}

	if sendFunc != nil {
		return sendFunc(n.id, to, msg)
	}
	return nil
}

// Receive receives a message into the node's inbox
func (n *BaseNode) Receive(env *message.Envelope) {
	n.mu.RLock()
	state := n.state
	n.mu.RUnlock()

	if state != StateRunning {
		return // Silently drop if crashed
	}

	n.inbox.Enqueue(env)
}

// Inbox returns the node's message queue
func (n *BaseNode) Inbox() *message.Queue {
	return n.inbox
}

// Crash simulates a node crash
func (n *BaseNode) Crash() {
	n.SetState(StateCrashed)
}

// Recover recovers a crashed node
func (n *BaseNode) Recover() {
	n.SetState(StateRunning)
}

// Emit emits an event for visualization
func (n *BaseNode) Emit(eventType string, data interface{}) {
	if n.emitter != nil {
		n.emitter(eventType, data)
	}
}

// Context returns the node's context
func (n *BaseNode) Context() context.Context {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.ctx
}

// GetVisualizationState returns state for UI rendering
func (n *BaseNode) GetVisualizationState() map[string]interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return map[string]interface{}{
		"id":         string(n.id),
		"state":      n.state.String(),
		"inboxSize":  n.inbox.Len(),
	}
}

// IsRunning returns true if the node is in running state
func (n *BaseNode) IsRunning() bool {
	return n.State() == StateRunning
}
