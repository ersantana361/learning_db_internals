package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/ersantana/db-internals/packages/protocol"
	"github.com/ersantana/db-internals/packages/simulation/engine"
)

// SimulationFactory creates a new simulation instance
type SimulationFactory func() engine.Simulation

// Session represents a client's active simulation session
type Session struct {
	ClientID string
	Project  string
	Engine   *engine.Engine
	Created  time.Time
}

// SimulationManager manages simulation sessions for clients
type SimulationManager struct {
	hub       *Hub
	factories map[string]SimulationFactory
	sessions  map[string]*Session
	mu        sync.RWMutex
}

// NewSimulationManager creates a new simulation manager
func NewSimulationManager(hub *Hub) *SimulationManager {
	return &SimulationManager{
		hub:       hub,
		factories: make(map[string]SimulationFactory),
		sessions:  make(map[string]*Session),
	}
}

// RegisterProject registers a simulation factory for a project
func (m *SimulationManager) RegisterProject(name string, factory SimulationFactory) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.factories[name] = factory
}

// HandleMessage processes incoming WebSocket messages
func (m *SimulationManager) HandleMessage(clientID string, data []byte) {
	msg, err := protocol.ParseMessage(data)
	if err != nil {
		m.sendError(clientID, "parse_error", "Failed to parse message")
		return
	}

	switch msg.Type {
	case protocol.MsgStartSimulation:
		m.handleStartSimulation(clientID, msg)
	case protocol.MsgStepForward:
		m.handleStepForward(clientID)
	case protocol.MsgStepBackward:
		m.handleStepBackward(clientID)
	case protocol.MsgResumeSimulation:
		m.handlePlay(clientID)
	case protocol.MsgPauseSimulation:
		m.handlePause(clientID)
	case protocol.MsgReset:
		m.handleReset(clientID)
	case protocol.MsgSetSpeed:
		m.handleSetSpeed(clientID, msg)
	case protocol.MsgExecuteOperation:
		m.handleExecuteOperation(clientID, msg)
	case protocol.MsgGetState:
		m.handleGetState(clientID)
	default:
		m.sendError(clientID, "unknown_message", "Unknown message type")
	}
}

// handleStartSimulation starts a new simulation
func (m *SimulationManager) handleStartSimulation(clientID string, msg *protocol.Message) {
	var req protocol.StartSimulationRequest
	if err := msg.ParsePayload(&req); err != nil {
		m.sendError(clientID, "invalid_payload", "Invalid start simulation request")
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Get factory for project
	factory, ok := m.factories[req.Config.Project]
	if !ok {
		m.sendError(clientID, "unknown_project", "Unknown project: "+req.Config.Project)
		return
	}

	// Create simulation instance
	sim := factory()

	// Create event emitter that sends to client
	emitter := func(event string, data interface{}) {
		m.sendEvent(clientID, event, data)
	}

	// Create engine
	eng := engine.NewEngine(sim, emitter)

	// Initialize with config
	if err := eng.Initialize(req.Config.Parameters); err != nil {
		m.sendError(clientID, "init_error", "Failed to initialize: "+err.Error())
		return
	}

	// Store session
	m.sessions[clientID] = &Session{
		ClientID: clientID,
		Project:  req.Config.Project,
		Engine:   eng,
		Created:  time.Now(),
	}

	// Send initial state
	m.sendState(clientID, eng)
}

// handleStepForward executes the next step
func (m *SimulationManager) handleStepForward(clientID string) {
	session := m.getSession(clientID)
	if session == nil {
		m.sendError(clientID, "no_session", "No active simulation")
		return
	}

	result, err := session.Engine.StepForward()
	if err != nil {
		if err == engine.ErrNoMoreSteps {
			m.sendError(clientID, "no_more_steps", "No more steps available")
		} else {
			m.sendError(clientID, "step_error", err.Error())
		}
		return
	}

	// State is sent via event emitter
	_ = result
}

// handleStepBackward goes to previous step
func (m *SimulationManager) handleStepBackward(clientID string) {
	session := m.getSession(clientID)
	if session == nil {
		m.sendError(clientID, "no_session", "No active simulation")
		return
	}

	_, err := session.Engine.StepBackward()
	if err != nil {
		m.sendError(clientID, "step_error", err.Error())
	}
}

// handlePlay starts automatic playback
func (m *SimulationManager) handlePlay(clientID string) {
	session := m.getSession(clientID)
	if session == nil {
		m.sendError(clientID, "no_session", "No active simulation")
		return
	}

	if err := session.Engine.Play(); err != nil {
		m.sendError(clientID, "play_error", err.Error())
	}
}

// handlePause pauses playback
func (m *SimulationManager) handlePause(clientID string) {
	session := m.getSession(clientID)
	if session == nil {
		return
	}

	session.Engine.Pause()
}

// handleReset resets the simulation
func (m *SimulationManager) handleReset(clientID string) {
	session := m.getSession(clientID)
	if session == nil {
		m.sendError(clientID, "no_session", "No active simulation")
		return
	}

	if err := session.Engine.Reset(); err != nil {
		m.sendError(clientID, "reset_error", err.Error())
		return
	}

	m.sendState(clientID, session.Engine)
}

// handleSetSpeed sets playback speed
func (m *SimulationManager) handleSetSpeed(clientID string, msg *protocol.Message) {
	session := m.getSession(clientID)
	if session == nil {
		m.sendError(clientID, "no_session", "No active simulation")
		return
	}

	var req protocol.SetSpeedRequest
	if err := msg.ParsePayload(&req); err != nil {
		m.sendError(clientID, "invalid_payload", "Invalid speed request")
		return
	}

	session.Engine.SetSpeed(req.Speed)
}

// handleExecuteOperation executes a custom operation
func (m *SimulationManager) handleExecuteOperation(clientID string, msg *protocol.Message) {
	// This will be implemented by individual simulations
	// For now, just acknowledge
	log.Printf("Execute operation for client %s", clientID)
}

// handleGetState returns current simulation state
func (m *SimulationManager) handleGetState(clientID string) {
	session := m.getSession(clientID)
	if session == nil {
		m.sendError(clientID, "no_session", "No active simulation")
		return
	}

	m.sendState(clientID, session.Engine)
}

// getSession returns the session for a client
func (m *SimulationManager) getSession(clientID string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[clientID]
}

// RemoveSession removes a client's session
func (m *SimulationManager) RemoveSession(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, ok := m.sessions[clientID]; ok {
		session.Engine.Stop()
		delete(m.sessions, clientID)
	}
}

// sendState sends the current simulation state
func (m *SimulationManager) sendState(clientID string, eng *engine.Engine) {
	state := eng.GetState()
	steps := eng.GetSteps()

	resp := protocol.SimulationStateResponse{
		State: *state,
		Steps: steps,
	}

	m.sendMessage(clientID, protocol.MsgSimulationState, resp)
}

// sendEvent sends an event to the client
func (m *SimulationManager) sendEvent(clientID string, event string, data interface{}) {
	// Convert engine events to protocol messages
	switch event {
	case "step_forward", "step_backward":
		if update, ok := data.(*protocol.StepUpdateResponse); ok {
			m.sendMessage(clientID, protocol.MsgStepUpdate, update)
		}
	case "initialized", "reset", "play", "pause":
		if state, ok := data.(*protocol.SimulationState); ok {
			m.sendMessage(clientID, protocol.MsgSimulationState, protocol.SimulationStateResponse{
				State: *state,
			})
		}
	}
}

// sendMessage sends a message to a specific client
func (m *SimulationManager) sendMessage(clientID string, msgType protocol.MessageType, payload interface{}) {
	msg, err := protocol.NewMessage(msgType, payload)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		return
	}

	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("Error serializing message: %v", err)
		return
	}

	m.hub.mu.RLock()
	client, ok := m.hub.clients[clientID]
	m.hub.mu.RUnlock()

	if ok {
		select {
		case client.Send <- data:
		default:
			log.Printf("Client %s send buffer full", clientID)
		}
	}
}

// sendError sends an error message to a client
func (m *SimulationManager) sendError(clientID string, code, message string) {
	resp := protocol.ErrorResponse{
		Code:    code,
		Message: message,
	}
	m.sendMessage(clientID, protocol.MsgError, resp)
}

// GetRegisteredProjects returns the list of registered projects
func (m *SimulationManager) GetRegisteredProjects() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	projects := make([]string, 0, len(m.factories))
	for name := range m.factories {
		projects = append(projects, name)
	}
	return projects
}
