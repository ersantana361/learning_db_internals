package engine

import (
	"errors"
	"sync"
	"time"

	"github.com/ersantana/db-internals/packages/protocol"
)

var (
	ErrNotInitialized   = errors.New("simulation not initialized")
	ErrAlreadyRunning   = errors.New("simulation already running")
	ErrNoMoreSteps      = errors.New("no more steps available")
	ErrNoPreviousSteps  = errors.New("no previous steps available")
	ErrInvalidStepIndex = errors.New("invalid step index")
)

// Step represents a single visualization step
type Step struct {
	Index       int
	Title       string
	Description string
	Execute     func() StepResult
	Highlights  []protocol.Highlight
}

// StepResult contains the result of executing a step
type StepResult struct {
	Success     bool
	Error       error
	Highlights  []protocol.Highlight
	Data        map[string]interface{}
	Description string
}

// Simulation is the interface all simulations must implement
type Simulation interface {
	// Identity
	Name() string
	Description() string

	// Lifecycle
	Initialize(config map[string]interface{}) error
	Reset() error

	// Step execution
	GenerateSteps() []Step
	CurrentStep() int
	ExecuteStep(index int) StepResult
	CanStepForward() bool
	CanStepBackward() bool

	// State
	GetState() interface{}
	GetVisualizationData() map[string]interface{}
}

// EventEmitter is called when simulation state changes
type EventEmitter func(event string, data interface{})

// Engine orchestrates simulation execution
type Engine struct {
	simulation  Simulation
	mode        protocol.SimulationMode
	currentStep int
	steps       []Step
	history     []StepResult
	speed       float64
	emitter     EventEmitter
	ticker      *time.Ticker
	stopChan    chan struct{}
	mu          sync.RWMutex
}

// NewEngine creates a new simulation engine
func NewEngine(sim Simulation, emitter EventEmitter) *Engine {
	return &Engine{
		simulation:  sim,
		mode:        protocol.ModeIdle,
		currentStep: -1,
		steps:       []Step{},
		history:     []StepResult{},
		speed:       1.0,
		emitter:     emitter,
		stopChan:    make(chan struct{}),
	}
}

// Initialize initializes the simulation with the given config
func (e *Engine) Initialize(config map[string]interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.simulation.Initialize(config); err != nil {
		return err
	}

	e.steps = e.simulation.GenerateSteps()
	e.currentStep = -1
	e.history = make([]StepResult, 0, len(e.steps))
	e.mode = protocol.ModeIdle

	e.emit("initialized", e.GetState())
	return nil
}

// Reset resets the simulation to its initial state
func (e *Engine) Reset() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Stop()

	if err := e.simulation.Reset(); err != nil {
		return err
	}

	e.steps = e.simulation.GenerateSteps()
	e.currentStep = -1
	e.history = make([]StepResult, 0, len(e.steps))
	e.mode = protocol.ModeIdle

	e.emit("reset", e.GetState())
	return nil
}

// StepForward executes the next step
func (e *Engine) StepForward() (*StepResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.steps) == 0 {
		return nil, ErrNotInitialized
	}

	nextStep := e.currentStep + 1
	if nextStep >= len(e.steps) {
		return nil, ErrNoMoreSteps
	}

	result := e.simulation.ExecuteStep(nextStep)
	e.currentStep = nextStep
	e.history = append(e.history, result)
	e.mode = protocol.ModeStep

	e.emit("step_forward", e.buildStepUpdate(result))
	return &result, nil
}

// StepBackward reverts to the previous step
func (e *Engine) StepBackward() (*StepResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.currentStep < 0 {
		return nil, ErrNoPreviousSteps
	}

	e.currentStep--
	if e.currentStep >= 0 && len(e.history) > e.currentStep {
		result := e.history[e.currentStep]
		e.emit("step_backward", e.buildStepUpdate(result))
		return &result, nil
	}

	// Reset to beginning
	e.simulation.Reset()
	e.history = e.history[:0]
	e.mode = protocol.ModeIdle

	e.emit("step_backward", e.GetState())
	return nil, nil
}

// Play starts automatic step execution
func (e *Engine) Play() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.mode == protocol.ModePlaying {
		return ErrAlreadyRunning
	}

	if len(e.steps) == 0 {
		return ErrNotInitialized
	}

	e.mode = protocol.ModePlaying
	e.stopChan = make(chan struct{})

	go e.playLoop()

	e.emit("play", e.GetState())
	return nil
}

// Pause pauses automatic execution
func (e *Engine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.mode == protocol.ModePlaying {
		e.mode = protocol.ModePaused
		close(e.stopChan)
		if e.ticker != nil {
			e.ticker.Stop()
		}
	}

	e.emit("pause", e.GetState())
}

// Stop stops the simulation completely
func (e *Engine) Stop() {
	if e.mode == protocol.ModePlaying {
		close(e.stopChan)
		if e.ticker != nil {
			e.ticker.Stop()
		}
	}
	e.mode = protocol.ModeIdle
}

// SetSpeed sets the playback speed
func (e *Engine) SetSpeed(speed float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if speed < 0.25 {
		speed = 0.25
	}
	if speed > 4.0 {
		speed = 4.0
	}
	e.speed = speed

	e.emit("speed_changed", map[string]interface{}{"speed": speed})
}

// GetState returns the current simulation state
func (e *Engine) GetState() *protocol.SimulationState {
	var currentStep *protocol.StepInfo
	if e.currentStep >= 0 && e.currentStep < len(e.steps) {
		step := e.steps[e.currentStep]
		highlights := step.Highlights
		if e.currentStep < len(e.history) {
			highlights = e.history[e.currentStep].Highlights
		}
		currentStep = &protocol.StepInfo{
			Index:       step.Index,
			Total:       len(e.steps),
			Title:       step.Title,
			Description: step.Description,
			Highlights:  highlights,
		}
	}

	return &protocol.SimulationState{
		Project:     e.simulation.Name(),
		Mode:        e.mode,
		Speed:       e.speed,
		CurrentStep: currentStep,
		TotalSteps:  len(e.steps),
		Data:        e.simulation.GetVisualizationData(),
	}
}

// GetSteps returns all step info
func (e *Engine) GetSteps() []protocol.StepInfo {
	e.mu.RLock()
	defer e.mu.RUnlock()

	steps := make([]protocol.StepInfo, len(e.steps))
	for i, s := range e.steps {
		steps[i] = protocol.StepInfo{
			Index:       i,
			Total:       len(e.steps),
			Title:       s.Title,
			Description: s.Description,
			Highlights:  s.Highlights,
		}
	}
	return steps
}

// CanStepForward returns true if there are more steps
func (e *Engine) CanStepForward() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentStep+1 < len(e.steps)
}

// CanStepBackward returns true if we can go back
func (e *Engine) CanStepBackward() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentStep >= 0
}

// playLoop runs the automatic playback
func (e *Engine) playLoop() {
	interval := time.Duration(float64(time.Second) / e.speed)
	e.ticker = time.NewTicker(interval)
	defer e.ticker.Stop()

	for {
		select {
		case <-e.stopChan:
			return
		case <-e.ticker.C:
			_, err := e.StepForward()
			if err == ErrNoMoreSteps {
				e.Pause()
				return
			}
		}
	}
}

// emit sends an event through the emitter
func (e *Engine) emit(event string, data interface{}) {
	if e.emitter != nil {
		e.emitter(event, data)
	}
}

// buildStepUpdate creates a step update response
func (e *Engine) buildStepUpdate(result StepResult) *protocol.StepUpdateResponse {
	var step protocol.StepInfo
	if e.currentStep >= 0 && e.currentStep < len(e.steps) {
		s := e.steps[e.currentStep]
		step = protocol.StepInfo{
			Index:       e.currentStep,
			Total:       len(e.steps),
			Title:       s.Title,
			Description: result.Description,
			Highlights:  result.Highlights,
			Data:        result.Data,
		}
	}

	return &protocol.StepUpdateResponse{
		Project:    e.simulation.Name(),
		Step:       step,
		Highlights: result.Highlights,
		Data:       e.simulation.GetVisualizationData(),
	}
}
