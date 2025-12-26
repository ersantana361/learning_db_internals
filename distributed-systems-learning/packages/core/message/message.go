package message

import (
	"encoding/json"
	"time"
)

// MessageType categorizes messages
type MessageType string

// Message is the base interface for all messages
type Message interface {
	Type() MessageType
	Payload() interface{}
}

// BaseMessage provides a basic message implementation
type BaseMessage struct {
	MsgType    MessageType `json:"type"`
	MsgPayload interface{} `json:"payload"`
}

func (m *BaseMessage) Type() MessageType {
	return m.MsgType
}

func (m *BaseMessage) Payload() interface{} {
	return m.MsgPayload
}

// NewMessage creates a new base message
func NewMessage(msgType MessageType, payload interface{}) *BaseMessage {
	return &BaseMessage{
		MsgType:    msgType,
		MsgPayload: payload,
	}
}

// Envelope wraps a message with routing and timing metadata
type Envelope struct {
	ID          string            `json:"id"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	Message     Message           `json:"message"`
	SentAt      time.Time         `json:"sentAt"`
	ReceivedAt  time.Time         `json:"receivedAt,omitempty"`
	LamportTime uint64            `json:"lamportTime,omitempty"`
	VectorClock map[string]uint64 `json:"vectorClock,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewEnvelope creates a new message envelope
func NewEnvelope(id, from, to string, msg Message) *Envelope {
	return &Envelope{
		ID:       id,
		From:     from,
		To:       to,
		Message:  msg,
		SentAt:   time.Now(),
		Metadata: make(map[string]interface{}),
	}
}

// Clone creates a copy of the envelope
func (e *Envelope) Clone() *Envelope {
	clone := &Envelope{
		ID:          e.ID,
		From:        e.From,
		To:          e.To,
		Message:     e.Message,
		SentAt:      e.SentAt,
		ReceivedAt:  e.ReceivedAt,
		LamportTime: e.LamportTime,
		Metadata:    make(map[string]interface{}),
	}

	if e.VectorClock != nil {
		clone.VectorClock = make(map[string]uint64)
		for k, v := range e.VectorClock {
			clone.VectorClock[k] = v
		}
	}

	for k, v := range e.Metadata {
		clone.Metadata[k] = v
	}

	return clone
}

// ToJSON serializes the envelope to JSON
func (e *Envelope) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// Queue is a thread-safe message queue
type Queue struct {
	messages chan *Envelope
	capacity int
}

// NewQueue creates a new message queue
func NewQueue(capacity int) *Queue {
	return &Queue{
		messages: make(chan *Envelope, capacity),
		capacity: capacity,
	}
}

// Enqueue adds a message to the queue
func (q *Queue) Enqueue(env *Envelope) bool {
	select {
	case q.messages <- env:
		return true
	default:
		return false // Queue full
	}
}

// Dequeue removes and returns a message from the queue
func (q *Queue) Dequeue() *Envelope {
	select {
	case env := <-q.messages:
		return env
	default:
		return nil
	}
}

// DequeueBlocking blocks until a message is available
func (q *Queue) DequeueBlocking() *Envelope {
	return <-q.messages
}

// Channel returns the underlying channel for select statements
func (q *Queue) Channel() <-chan *Envelope {
	return q.messages
}

// Len returns the current number of messages in the queue
func (q *Queue) Len() int {
	return len(q.messages)
}

// Close closes the queue
func (q *Queue) Close() {
	close(q.messages)
}
