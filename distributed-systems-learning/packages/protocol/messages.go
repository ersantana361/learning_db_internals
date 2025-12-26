package protocol

import (
	"encoding/json"
)

// MessageType defines WebSocket message types
type MessageType string

// Client -> Server message types
const (
	// Simulation control
	MsgStartSimulation   MessageType = "start_simulation"
	MsgPauseSimulation   MessageType = "pause_simulation"
	MsgResumeSimulation  MessageType = "resume_simulation"
	MsgStopSimulation    MessageType = "stop_simulation"
	MsgStepForward       MessageType = "step_forward"
	MsgSetSpeed          MessageType = "set_speed"

	// Failure injection
	MsgInjectCrash     MessageType = "inject_crash"
	MsgRecoverNode     MessageType = "recover_node"
	MsgInjectPartition MessageType = "inject_partition"
	MsgHealPartition   MessageType = "heal_partition"

	// User interactions
	MsgSendClientRequest MessageType = "send_client_request"
	MsgSelectScenario    MessageType = "select_scenario"

	// Query state
	MsgGetState MessageType = "get_state"
)

// Server -> Client message types
const (
	// State updates
	MsgSimulationState  MessageType = "simulation_state"
	MsgNodeStateUpdate  MessageType = "node_state_update"

	// Events
	MsgMessageSent     MessageType = "message_sent"
	MsgMessageReceived MessageType = "message_received"
	MsgMessageDropped  MessageType = "message_dropped"
	MsgLeaderElected   MessageType = "leader_elected"
	MsgConsensusReached MessageType = "consensus_reached"
	MsgTransactionState MessageType = "transaction_state"

	// Visualization
	MsgTimelineEvent MessageType = "timeline_event"
	MsgClockUpdate   MessageType = "clock_update"

	// Errors
	MsgError MessageType = "error"
)

// BaseMessage is the base structure for all messages
type BaseMessage struct {
	Type MessageType `json:"type"`
}

// StartSimulationRequest starts a simulation
type StartSimulationRequest struct {
	Type     MessageType `json:"type"`
	Project  string      `json:"project"`
	Scenario string      `json:"scenario,omitempty"`
	Config   struct {
		NodeCount int     `json:"nodeCount,omitempty"`
		Speed     float64 `json:"speed,omitempty"`
		StepMode  bool    `json:"stepMode,omitempty"`
	} `json:"config,omitempty"`
}

// SetSpeedRequest sets simulation speed
type SetSpeedRequest struct {
	Type  MessageType `json:"type"`
	Speed float64     `json:"speed"`
}

// InjectCrashRequest crashes a node
type InjectCrashRequest struct {
	Type   MessageType `json:"type"`
	NodeID string      `json:"nodeId"`
}

// RecoverNodeRequest recovers a crashed node
type RecoverNodeRequest struct {
	Type   MessageType `json:"type"`
	NodeID string      `json:"nodeId"`
}

// InjectPartitionRequest creates a network partition
type InjectPartitionRequest struct {
	Type          MessageType `json:"type"`
	From          string      `json:"from"`
	To            string      `json:"to"`
	Bidirectional bool        `json:"bidirectional,omitempty"`
}

// HealPartitionRequest heals a network partition
type HealPartitionRequest struct {
	Type          MessageType `json:"type"`
	From          string      `json:"from"`
	To            string      `json:"to"`
	Bidirectional bool        `json:"bidirectional,omitempty"`
}

// ClientRequest sends a client request to the simulation
type ClientRequest struct {
	Type    MessageType            `json:"type"`
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// SimulationStateResponse contains the full simulation state
type SimulationStateResponse struct {
	Type        MessageType              `json:"type"`
	VirtualTime int64                    `json:"virtualTime"`
	Mode        string                   `json:"mode"`
	Speed       float64                  `json:"speed"`
	Running     bool                     `json:"running"`
	Nodes       map[string]NodeState     `json:"nodes"`
	Messages    []MessageState           `json:"messages,omitempty"`
	Partitions  []PartitionState         `json:"partitions,omitempty"`
	Timeline    []TimelineEvent          `json:"timeline,omitempty"`
}

// NodeState represents a node's state
type NodeState struct {
	ID          string                 `json:"id"`
	Status      string                 `json:"status"`
	Role        string                 `json:"role,omitempty"`
	Term        int                    `json:"term,omitempty"`
	VotedFor    string                 `json:"votedFor,omitempty"`
	Log         []LogEntry             `json:"log,omitempty"`
	CommitIndex int                    `json:"commitIndex,omitempty"`
	Clock       map[string]uint64      `json:"clock,omitempty"`
	CustomState map[string]interface{} `json:"customState,omitempty"`
}

// LogEntry represents a log entry
type LogEntry struct {
	Index   int         `json:"index"`
	Term    int         `json:"term"`
	Command interface{} `json:"command"`
}

// MessageState represents an in-flight message
type MessageState struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Type    string `json:"type"`
	Status  string `json:"status"` // "pending", "delivered", "dropped"
}

// PartitionState represents a network partition
type PartitionState struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// TimelineEvent represents an event in the timeline
type TimelineEvent struct {
	Time int64                  `json:"time"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// NodeStateUpdateResponse updates a single node's state
type NodeStateUpdateResponse struct {
	Type     MessageType            `json:"type"`
	NodeID   string                 `json:"nodeId"`
	OldState string                 `json:"oldState,omitempty"`
	NewState string                 `json:"newState"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// MessageEventResponse represents a message event
type MessageEventResponse struct {
	Type        MessageType       `json:"type"`
	MessageID   string            `json:"messageId"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	MessageType string            `json:"messageType"`
	Payload     interface{}       `json:"payload,omitempty"`
	Clock       map[string]uint64 `json:"clock,omitempty"`
	Reason      string            `json:"reason,omitempty"` // For dropped messages
	Latency     int64             `json:"latency,omitempty"` // For received messages
}

// ErrorResponse represents an error
type ErrorResponse struct {
	Type    MessageType `json:"type"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
}

// ParseMessage parses a JSON message and returns its type
func ParseMessage(data []byte) (MessageType, error) {
	var base BaseMessage
	if err := json.Unmarshal(data, &base); err != nil {
		return "", err
	}
	return base.Type, nil
}

// ParseStartSimulation parses a start simulation message
func ParseStartSimulation(data []byte) (*StartSimulationRequest, error) {
	var msg StartSimulationRequest
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ParseSetSpeed parses a set speed message
func ParseSetSpeed(data []byte) (*SetSpeedRequest, error) {
	var msg SetSpeedRequest
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ParseInjectCrash parses an inject crash message
func ParseInjectCrash(data []byte) (*InjectCrashRequest, error) {
	var msg InjectCrashRequest
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// NewSimulationState creates a new simulation state response
func NewSimulationState(virtualTime int64, mode string, speed float64, running bool, nodes map[string]NodeState) *SimulationStateResponse {
	return &SimulationStateResponse{
		Type:        MsgSimulationState,
		VirtualTime: virtualTime,
		Mode:        mode,
		Speed:       speed,
		Running:     running,
		Nodes:       nodes,
	}
}

// NewError creates a new error response
func NewError(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Type:    MsgError,
		Code:    code,
		Message: message,
	}
}

// ToJSON serializes a message to JSON
func ToJSON(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}
