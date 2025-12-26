package events

import (
	"sync"
)

// Listener is a function that handles events
type Listener func(event Event)

// EventBus manages event distribution
type EventBus struct {
	mu        sync.RWMutex
	listeners []Listener
	channels  []chan Event
	buffer    []Event
	recording bool
	closed    bool
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make([]Listener, 0),
		channels:  make([]chan Event, 0),
		buffer:    make([]Event, 0),
		recording: false,
	}
}

// Subscribe registers a listener function
func (eb *EventBus) Subscribe(listener Listener) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.listeners = append(eb.listeners, listener)
}

// SubscribeChannel returns a channel for receiving events
func (eb *EventBus) SubscribeChannel(bufferSize int) <-chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(chan Event, bufferSize)
	eb.channels = append(eb.channels, ch)
	return ch
}

// Emit broadcasts an event to all subscribers
func (eb *EventBus) Emit(event Event) {
	eb.mu.RLock()
	if eb.closed {
		eb.mu.RUnlock()
		return
	}

	listeners := eb.listeners
	channels := eb.channels
	recording := eb.recording
	eb.mu.RUnlock()

	// Record if enabled
	if recording {
		eb.mu.Lock()
		eb.buffer = append(eb.buffer, event)
		eb.mu.Unlock()
	}

	// Notify function listeners
	for _, listener := range listeners {
		go listener(event)
	}

	// Send to channel subscribers
	for _, ch := range channels {
		select {
		case ch <- event:
		default:
			// Channel full, skip (non-blocking)
		}
	}
}

// StartRecording starts recording events for replay
func (eb *EventBus) StartRecording() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.recording = true
	eb.buffer = make([]Event, 0)
}

// StopRecording stops recording events
func (eb *EventBus) StopRecording() []Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.recording = false
	events := eb.buffer
	eb.buffer = make([]Event, 0)
	return events
}

// GetRecordedEvents returns all recorded events
func (eb *EventBus) GetRecordedEvents() []Event {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	events := make([]Event, len(eb.buffer))
	copy(events, eb.buffer)
	return events
}

// ClearRecording clears the recorded events
func (eb *EventBus) ClearRecording() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.buffer = make([]Event, 0)
}

// Close closes all channels and stops the event bus
func (eb *EventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.closed = true
	for _, ch := range eb.channels {
		close(ch)
	}
	eb.channels = nil
	eb.listeners = nil
}

// Replay replays recorded events with optional delay
type Replay struct {
	events []Event
	index  int
}

// NewReplay creates a new replay from recorded events
func NewReplay(events []Event) *Replay {
	return &Replay{
		events: events,
		index:  0,
	}
}

// Next returns the next event, or nil if done
func (r *Replay) Next() Event {
	if r.index >= len(r.events) {
		return nil
	}
	event := r.events[r.index]
	r.index++
	return event
}

// HasNext returns true if there are more events
func (r *Replay) HasNext() bool {
	return r.index < len(r.events)
}

// Reset resets the replay to the beginning
func (r *Replay) Reset() {
	r.index = 0
}

// Len returns the total number of events
func (r *Replay) Len() int {
	return len(r.events)
}

// Current returns the current position
func (r *Replay) Current() int {
	return r.index
}
