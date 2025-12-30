package protocol

import "encoding/json"

// MessageType represents WebSocket message types
type MessageType string

// Client -> Server message types
const (
	// Simulation control
	MsgStartSimulation  MessageType = "start_simulation"
	MsgPauseSimulation  MessageType = "pause_simulation"
	MsgResumeSimulation MessageType = "resume_simulation"
	MsgStopSimulation   MessageType = "stop_simulation"
	MsgStepForward      MessageType = "step_forward"
	MsgStepBackward     MessageType = "step_backward"
	MsgSetSpeed         MessageType = "set_speed"
	MsgReset            MessageType = "reset"

	// User interactions
	MsgExecuteOperation MessageType = "execute_operation"
	MsgSelectScenario   MessageType = "select_scenario"
	MsgGetState         MessageType = "get_state"
)

// Server -> Client message types
const (
	// State updates
	MsgSimulationState MessageType = "simulation_state"
	MsgStepUpdate      MessageType = "step_update"

	// Visualization events
	MsgNodeHighlight  MessageType = "node_highlight"
	MsgNodeUpdate     MessageType = "node_update"
	MsgEdgeHighlight  MessageType = "edge_highlight"
	MsgAnimationStart MessageType = "animation_start"
	MsgAnimationEnd   MessageType = "animation_end"

	// Errors
	MsgError MessageType = "error"
)

// Message is the base WebSocket message structure
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// SimulationMode represents the current simulation state
type SimulationMode string

const (
	ModeIdle    SimulationMode = "idle"
	ModePlaying SimulationMode = "playing"
	ModePaused  SimulationMode = "paused"
	ModeStep    SimulationMode = "step"
)

// SimulationConfig holds configuration for starting a simulation
type SimulationConfig struct {
	Project    string                 `json:"project"`    // "btree", "mvcc", "query-parser"
	Scenario   string                 `json:"scenario"`   // Predefined scenario name
	StepMode   bool                   `json:"stepMode"`   // Start in step mode
	Speed      float64                `json:"speed"`      // Animation speed multiplier
	Parameters map[string]interface{} `json:"parameters"` // Topic-specific parameters
}

// StepInfo describes the current step in visualization
type StepInfo struct {
	Index       int                    `json:"index"`
	Total       int                    `json:"total"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Highlights  []Highlight            `json:"highlights"`
	Data        map[string]interface{} `json:"data"`
}

// Highlight represents a visual highlight on the visualization
type Highlight struct {
	Type      string `json:"type"`      // "node", "edge", "cell", "row", "token"
	ID        string `json:"id"`        // Element identifier
	Color     string `json:"color"`     // Highlight color
	Animation string `json:"animation"` // "pulse", "flash", "fade", "none"
	Duration  int    `json:"duration"`  // Duration in milliseconds
}

// SimulationState represents the full state of a simulation
type SimulationState struct {
	Project     string                 `json:"project"`
	Mode        SimulationMode         `json:"mode"`
	Speed       float64                `json:"speed"`
	CurrentStep *StepInfo              `json:"currentStep"`
	TotalSteps  int                    `json:"totalSteps"`
	Data        map[string]interface{} `json:"data"`
}

// --- Request Payloads ---

// StartSimulationRequest is the payload for start_simulation
type StartSimulationRequest struct {
	Config SimulationConfig `json:"config"`
}

// SetSpeedRequest is the payload for set_speed
type SetSpeedRequest struct {
	Speed float64 `json:"speed"`
}

// ExecuteOperationRequest is the payload for execute_operation
type ExecuteOperationRequest struct {
	Operation string                 `json:"operation"` // "insert", "delete", "search", etc.
	Params    map[string]interface{} `json:"params"`
}

// SelectScenarioRequest is the payload for select_scenario
type SelectScenarioRequest struct {
	ScenarioID string `json:"scenarioId"`
}

// --- Response Payloads ---

// SimulationStateResponse is the payload for simulation_state
type SimulationStateResponse struct {
	State SimulationState `json:"state"`
	Steps []StepInfo      `json:"steps"`
}

// StepUpdateResponse is the payload for step_update
type StepUpdateResponse struct {
	Project    string                 `json:"project"`
	Step       StepInfo               `json:"step"`
	Highlights []Highlight            `json:"highlights"`
	Data       map[string]interface{} `json:"data"`
}

// ErrorResponse is the payload for error messages
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// --- Helper Functions ---

// NewMessage creates a new message with the given type and payload
func NewMessage(msgType MessageType, payload interface{}) (*Message, error) {
	var rawPayload json.RawMessage
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		rawPayload = data
	}
	return &Message{Type: msgType, Payload: rawPayload}, nil
}

// ParseMessage parses a raw JSON message
func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ParsePayload parses the payload into the given type
func (m *Message) ParsePayload(v interface{}) error {
	if m.Payload == nil {
		return nil
	}
	return json.Unmarshal(m.Payload, v)
}

// ToJSON serializes the message to JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
