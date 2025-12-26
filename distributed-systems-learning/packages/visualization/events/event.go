package events

import (
	"encoding/json"
	"time"
)

// EventType categorizes events
type EventType string

const (
	// Message events
	EventMessageSent     EventType = "message_sent"
	EventMessageReceived EventType = "message_received"
	EventMessageDropped  EventType = "message_dropped"

	// Node events
	EventNodeStateChanged EventType = "node_state_changed"
	EventNodeCrashed      EventType = "node_crashed"
	EventNodeRecovered    EventType = "node_recovered"

	// Network events
	EventPartitionCreated EventType = "partition_created"
	EventPartitionHealed  EventType = "partition_healed"

	// Algorithm-specific events
	EventLeaderElected    EventType = "leader_elected"
	EventVoteRequested    EventType = "vote_requested"
	EventVoteCast         EventType = "vote_cast"
	EventConsensusReached EventType = "consensus_reached"
	EventLogAppended      EventType = "log_appended"
	EventLogCommitted     EventType = "log_committed"

	// Transaction events
	EventTransactionStarted   EventType = "transaction_started"
	EventTransactionPrepared  EventType = "transaction_prepared"
	EventTransactionCommitted EventType = "transaction_committed"
	EventTransactionAborted   EventType = "transaction_aborted"

	// Clock events
	EventClockTick   EventType = "clock_tick"
	EventClockMerge  EventType = "clock_merge"
	EventClockUpdate EventType = "clock_update"
)

// Event is the base interface for all visualization events
type Event interface {
	EventType() EventType
	Timestamp() time.Time
	Data() map[string]interface{}
	ToJSON() ([]byte, error)
}

// BaseEvent provides common event functionality
type BaseEvent struct {
	Type      EventType              `json:"type"`
	Time      time.Time              `json:"timestamp"`
	EventData map[string]interface{} `json:"data"`
}

func (e *BaseEvent) EventType() EventType {
	return e.Type
}

func (e *BaseEvent) Timestamp() time.Time {
	return e.Time
}

func (e *BaseEvent) Data() map[string]interface{} {
	return e.EventData
}

func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// NewEvent creates a new base event
func NewEvent(eventType EventType, data map[string]interface{}) *BaseEvent {
	return &BaseEvent{
		Type:      eventType,
		Time:      time.Now(),
		EventData: data,
	}
}

// MessageSentEvent represents a message being sent
type MessageSentEvent struct {
	BaseEvent
	From        string            `json:"from"`
	To          string            `json:"to"`
	MessageID   string            `json:"messageId"`
	MessageType string            `json:"messageType"`
	Payload     interface{}       `json:"payload,omitempty"`
	Clock       map[string]uint64 `json:"clock,omitempty"`
}

// NewMessageSentEvent creates a message sent event
func NewMessageSentEvent(from, to, msgID, msgType string, payload interface{}, clock map[string]uint64) *MessageSentEvent {
	return &MessageSentEvent{
		BaseEvent: BaseEvent{
			Type: EventMessageSent,
			Time: time.Now(),
			EventData: map[string]interface{}{
				"from":        from,
				"to":          to,
				"messageId":   msgID,
				"messageType": msgType,
			},
		},
		From:        from,
		To:          to,
		MessageID:   msgID,
		MessageType: msgType,
		Payload:     payload,
		Clock:       clock,
	}
}

// MessageReceivedEvent represents a message being received
type MessageReceivedEvent struct {
	BaseEvent
	At        string        `json:"at"`
	From      string        `json:"from"`
	MessageID string        `json:"messageId"`
	Latency   time.Duration `json:"latency"`
}

// NewMessageReceivedEvent creates a message received event
func NewMessageReceivedEvent(at, from, msgID string, latency time.Duration) *MessageReceivedEvent {
	return &MessageReceivedEvent{
		BaseEvent: BaseEvent{
			Type: EventMessageReceived,
			Time: time.Now(),
			EventData: map[string]interface{}{
				"at":        at,
				"from":      from,
				"messageId": msgID,
				"latency":   latency.Milliseconds(),
			},
		},
		At:        at,
		From:      from,
		MessageID: msgID,
		Latency:   latency,
	}
}

// MessageDroppedEvent represents a message being dropped
type MessageDroppedEvent struct {
	BaseEvent
	MessageID string `json:"messageId"`
	From      string `json:"from"`
	To        string `json:"to"`
	Reason    string `json:"reason"`
}

// NewMessageDroppedEvent creates a message dropped event
func NewMessageDroppedEvent(msgID, from, to, reason string) *MessageDroppedEvent {
	return &MessageDroppedEvent{
		BaseEvent: BaseEvent{
			Type: EventMessageDropped,
			Time: time.Now(),
			EventData: map[string]interface{}{
				"messageId": msgID,
				"from":      from,
				"to":        to,
				"reason":    reason,
			},
		},
		MessageID: msgID,
		From:      from,
		To:        to,
		Reason:    reason,
	}
}

// NodeStateChangedEvent represents a node state change
type NodeStateChangedEvent struct {
	BaseEvent
	NodeID   string                 `json:"nodeId"`
	OldState string                 `json:"oldState"`
	NewState string                 `json:"newState"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// NewNodeStateChangedEvent creates a node state changed event
func NewNodeStateChangedEvent(nodeID, oldState, newState string, details map[string]interface{}) *NodeStateChangedEvent {
	return &NodeStateChangedEvent{
		BaseEvent: BaseEvent{
			Type: EventNodeStateChanged,
			Time: time.Now(),
			EventData: map[string]interface{}{
				"nodeId":   nodeID,
				"oldState": oldState,
				"newState": newState,
			},
		},
		NodeID:   nodeID,
		OldState: oldState,
		NewState: newState,
		Details:  details,
	}
}

// LeaderElectedEvent represents a leader being elected
type LeaderElectedEvent struct {
	BaseEvent
	LeaderID string `json:"leaderId"`
	Term     int    `json:"term"`
}

// NewLeaderElectedEvent creates a leader elected event
func NewLeaderElectedEvent(leaderID string, term int) *LeaderElectedEvent {
	return &LeaderElectedEvent{
		BaseEvent: BaseEvent{
			Type: EventLeaderElected,
			Time: time.Now(),
			EventData: map[string]interface{}{
				"leaderId": leaderID,
				"term":     term,
			},
		},
		LeaderID: leaderID,
		Term:     term,
	}
}

// ConsensusReachedEvent represents consensus being reached
type ConsensusReachedEvent struct {
	BaseEvent
	Value        interface{} `json:"value"`
	Term         int         `json:"term,omitempty"`
	Participants []string    `json:"participants"`
}

// NewConsensusReachedEvent creates a consensus reached event
func NewConsensusReachedEvent(value interface{}, term int, participants []string) *ConsensusReachedEvent {
	return &ConsensusReachedEvent{
		BaseEvent: BaseEvent{
			Type: EventConsensusReached,
			Time: time.Now(),
			EventData: map[string]interface{}{
				"value":        value,
				"term":         term,
				"participants": participants,
			},
		},
		Value:        value,
		Term:         term,
		Participants: participants,
	}
}
