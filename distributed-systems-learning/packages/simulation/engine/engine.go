package engine

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// SimulationMode defines execution mode
type SimulationMode int

const (
	ModeRealtime SimulationMode = iota
	ModeStepByStep
	ModePaused
)

func (m SimulationMode) String() string {
	switch m {
	case ModeRealtime:
		return "realtime"
	case ModeStepByStep:
		return "step"
	case ModePaused:
		return "paused"
	default:
		return "unknown"
	}
}

// NodeController interface for simulation nodes
type NodeController interface {
	ID() string
	Start(ctx context.Context) error
	Stop() error
	Tick() // Process one step
	GetState() map[string]interface{}
}

// EventEmitter interface for emitting events
type EventEmitter interface {
	Emit(eventType string, data map[string]interface{})
}

// Config holds simulation configuration
type Config struct {
	Speed       float64 // Speed multiplier (1.0 = realtime)
	TickRate    time.Duration
	StepMode    bool
	ProjectName string
	Scenario    string
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Speed:    1.0,
		TickRate: 100 * time.Millisecond,
		StepMode: false,
	}
}

// Engine orchestrates distributed simulations
type Engine struct {
	mu sync.RWMutex

	nodes   map[string]NodeController
	emitter EventEmitter
	config  Config

	mode        SimulationMode
	stepCh      chan struct{}
	speed       float64
	virtualTime time.Time
	startTime   time.Time

	ctx    context.Context
	cancel context.CancelFunc

	running bool
}

// NewEngine creates a new simulation engine
func NewEngine(emitter EventEmitter, config Config) *Engine {
	return &Engine{
		nodes:   make(map[string]NodeController),
		emitter: emitter,
		config:  config,
		stepCh:  make(chan struct{}, 100),
		speed:   config.Speed,
		mode:    ModePaused,
	}
}

// AddNode registers a node with the simulation
func (e *Engine) AddNode(node NodeController) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.nodes[node.ID()] = node
}

// RemoveNode removes a node from the simulation
func (e *Engine) RemoveNode(nodeID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.nodes, nodeID)
}

// GetNode returns a node by ID
func (e *Engine) GetNode(nodeID string) NodeController {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.nodes[nodeID]
}

// Start starts the simulation
func (e *Engine) Start(ctx context.Context) error {
	e.mu.Lock()
	e.ctx, e.cancel = context.WithCancel(ctx)
	e.startTime = time.Now()
	e.virtualTime = e.startTime
	e.running = true

	if e.config.StepMode {
		e.mode = ModeStepByStep
	} else {
		e.mode = ModeRealtime
	}
	e.mu.Unlock()

	// Start all nodes
	for _, node := range e.nodes {
		if err := node.Start(e.ctx); err != nil {
			return err
		}
	}

	if e.emitter != nil {
		e.emitter.Emit("simulation_started", map[string]interface{}{
			"mode":   e.mode.String(),
			"speed":  e.speed,
			"config": e.config,
		})
	}

	// Start main loop
	go e.run()

	return nil
}

// Stop stops the simulation
func (e *Engine) Stop() error {
	e.mu.Lock()
	e.running = false
	if e.cancel != nil {
		e.cancel()
	}
	e.mu.Unlock()

	// Stop all nodes
	for _, node := range e.nodes {
		node.Stop()
	}

	if e.emitter != nil {
		e.emitter.Emit("simulation_stopped", map[string]interface{}{})
	}

	return nil
}

// run is the main simulation loop
func (e *Engine) run() {
	tickDuration := e.config.TickRate

	for {
		e.mu.RLock()
		running := e.running
		mode := e.mode
		speed := e.speed
		e.mu.RUnlock()

		if !running {
			return
		}

		switch mode {
		case ModeRealtime:
			e.tick()
			adjustedDuration := time.Duration(float64(tickDuration) / speed)
			time.Sleep(adjustedDuration)

		case ModeStepByStep:
			select {
			case <-e.stepCh:
				e.tick()
			case <-e.ctx.Done():
				return
			}

		case ModePaused:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// tick performs one simulation step
func (e *Engine) tick() {
	e.mu.Lock()
	e.virtualTime = e.virtualTime.Add(e.config.TickRate)
	e.mu.Unlock()

	// Process each node
	e.mu.RLock()
	nodes := make([]NodeController, 0, len(e.nodes))
	for _, node := range e.nodes {
		nodes = append(nodes, node)
	}
	e.mu.RUnlock()

	for _, node := range nodes {
		node.Tick()
	}

	if e.emitter != nil {
		e.emitter.Emit("simulation_tick", map[string]interface{}{
			"virtualTime": e.virtualTime.UnixMilli(),
		})
	}
}

// Step advances simulation by one step (for step-by-step mode)
func (e *Engine) Step() {
	e.stepCh <- struct{}{}
}

// StepN advances simulation by n steps
func (e *Engine) StepN(n int) {
	for i := 0; i < n; i++ {
		e.Step()
	}
}

// Pause pauses the simulation
func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.mode = ModePaused
	if e.emitter != nil {
		e.emitter.Emit("simulation_paused", map[string]interface{}{})
	}
}

// Resume resumes the simulation
func (e *Engine) Resume() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.config.StepMode {
		e.mode = ModeStepByStep
	} else {
		e.mode = ModeRealtime
	}
	if e.emitter != nil {
		e.emitter.Emit("simulation_resumed", map[string]interface{}{
			"mode": e.mode.String(),
		})
	}
}

// SetSpeed sets simulation speed multiplier
func (e *Engine) SetSpeed(speed float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if speed < 0.1 {
		speed = 0.1
	}
	if speed > 10.0 {
		speed = 10.0
	}
	e.speed = speed
}

// SetMode sets the simulation mode
func (e *Engine) SetMode(mode SimulationMode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.mode = mode
}

// GetState returns current simulation state for visualization
func (e *Engine) GetState() SimulationState {
	e.mu.RLock()
	defer e.mu.RUnlock()

	nodeStates := make(map[string]interface{})
	for id, node := range e.nodes {
		nodeStates[id] = node.GetState()
	}

	return SimulationState{
		Mode:        e.mode.String(),
		Speed:       e.speed,
		VirtualTime: e.virtualTime.UnixMilli(),
		Running:     e.running,
		Nodes:       nodeStates,
	}
}

// SimulationState represents the current state of the simulation
type SimulationState struct {
	Mode        string                 `json:"mode"`
	Speed       float64                `json:"speed"`
	VirtualTime int64                  `json:"virtualTime"`
	Running     bool                   `json:"running"`
	Nodes       map[string]interface{} `json:"nodes"`
}

// ToJSON serializes the state to JSON
func (s *SimulationState) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// GetVirtualTime returns the current virtual time
func (e *Engine) GetVirtualTime() time.Time {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.virtualTime
}

// IsRunning returns true if the simulation is running
func (e *Engine) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

// GetMode returns the current simulation mode
func (e *Engine) GetMode() SimulationMode {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mode
}

// NodeCount returns the number of nodes
func (e *Engine) NodeCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.nodes)
}

// GetNodeIDs returns all node IDs
func (e *Engine) GetNodeIDs() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ids := make([]string, 0, len(e.nodes))
	for id := range e.nodes {
		ids = append(ids, id)
	}
	return ids
}
