package injector

import (
	"sync"
	"time"
)

// FailureType categorizes failures
type FailureType int

const (
	FailureCrash FailureType = iota
	FailurePartition
	FailureDelay
	FailureByzantine
)

func (f FailureType) String() string {
	switch f {
	case FailureCrash:
		return "crash"
	case FailurePartition:
		return "partition"
	case FailureDelay:
		return "delay"
	case FailureByzantine:
		return "byzantine"
	default:
		return "unknown"
	}
}

// Failure represents a scheduled or active failure
type Failure struct {
	ID        string
	Type      FailureType
	Target    string                 // Node ID or "partition:A:B"
	StartTime time.Duration          // Relative to simulation start
	Duration  time.Duration          // How long failure lasts (0 = permanent)
	Params    map[string]interface{}
	Active    bool
}

// NodeManager interface for controlling nodes
type NodeManager interface {
	CrashNode(nodeID string)
	RecoverNode(nodeID string)
	SetNodeDelay(nodeID string, delay time.Duration)
	ClearNodeDelay(nodeID string)
}

// NetworkManager interface for controlling network
type NetworkManager interface {
	CreatePartition(from, to string)
	HealPartition(from, to string)
	SetLatency(min, max time.Duration)
}

// EventEmitter interface for emitting events
type EventEmitter interface {
	Emit(eventType string, data map[string]interface{})
}

// Injector manages failure injection
type Injector struct {
	mu sync.RWMutex

	failures       map[string]*Failure
	scheduled      []*scheduledFailure
	nodeManager    NodeManager
	networkManager NetworkManager
	emitter        EventEmitter

	startTime time.Time
	running   bool
}

type scheduledFailure struct {
	failure   *Failure
	executeAt time.Time
	isRecover bool
}

// NewInjector creates a new failure injector
func NewInjector(nodeManager NodeManager, networkManager NetworkManager, emitter EventEmitter) *Injector {
	return &Injector{
		failures:       make(map[string]*Failure),
		scheduled:      make([]*scheduledFailure, 0),
		nodeManager:    nodeManager,
		networkManager: networkManager,
		emitter:        emitter,
	}
}

// InjectCrash immediately crashes a node
func (i *Injector) InjectCrash(nodeID string) *Failure {
	i.mu.Lock()
	defer i.mu.Unlock()

	failure := &Failure{
		ID:     generateID(),
		Type:   FailureCrash,
		Target: nodeID,
		Active: true,
	}

	i.failures[failure.ID] = failure

	if i.nodeManager != nil {
		i.nodeManager.CrashNode(nodeID)
	}

	if i.emitter != nil {
		i.emitter.Emit("node_crashed", map[string]interface{}{
			"nodeId":    nodeID,
			"failureId": failure.ID,
		})
	}

	return failure
}

// RecoverNode recovers a crashed node
func (i *Injector) RecoverNode(nodeID string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Find and deactivate the crash failure
	for id, f := range i.failures {
		if f.Target == nodeID && f.Type == FailureCrash && f.Active {
			f.Active = false
			delete(i.failures, id)
			break
		}
	}

	if i.nodeManager != nil {
		i.nodeManager.RecoverNode(nodeID)
	}

	if i.emitter != nil {
		i.emitter.Emit("node_recovered", map[string]interface{}{
			"nodeId": nodeID,
		})
	}
}

// InjectPartition creates a network partition between two nodes
func (i *Injector) InjectPartition(from, to string, bidirectional bool) *Failure {
	i.mu.Lock()
	defer i.mu.Unlock()

	target := from + ":" + to
	if bidirectional {
		target = target + ":bidirectional"
	}

	failure := &Failure{
		ID:     generateID(),
		Type:   FailurePartition,
		Target: target,
		Params: map[string]interface{}{
			"from":          from,
			"to":            to,
			"bidirectional": bidirectional,
		},
		Active: true,
	}

	i.failures[failure.ID] = failure

	if i.networkManager != nil {
		i.networkManager.CreatePartition(from, to)
		if bidirectional {
			i.networkManager.CreatePartition(to, from)
		}
	}

	if i.emitter != nil {
		i.emitter.Emit("partition_created", map[string]interface{}{
			"from":          from,
			"to":            to,
			"bidirectional": bidirectional,
			"failureId":     failure.ID,
		})
	}

	return failure
}

// HealPartition removes a network partition
func (i *Injector) HealPartition(from, to string, bidirectional bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Find and deactivate the partition failure
	for id, f := range i.failures {
		if f.Type == FailurePartition && f.Active {
			params := f.Params
			if params["from"] == from && params["to"] == to {
				f.Active = false
				delete(i.failures, id)
				break
			}
		}
	}

	if i.networkManager != nil {
		i.networkManager.HealPartition(from, to)
		if bidirectional {
			i.networkManager.HealPartition(to, from)
		}
	}

	if i.emitter != nil {
		i.emitter.Emit("partition_healed", map[string]interface{}{
			"from":          from,
			"to":            to,
			"bidirectional": bidirectional,
		})
	}
}

// ScheduleFailure schedules a failure for future execution
func (i *Injector) ScheduleFailure(failure *Failure) {
	i.mu.Lock()
	defer i.mu.Unlock()

	executeAt := i.startTime.Add(failure.StartTime)

	i.scheduled = append(i.scheduled, &scheduledFailure{
		failure:   failure,
		executeAt: executeAt,
		isRecover: false,
	})

	// Schedule recovery if duration is set
	if failure.Duration > 0 {
		i.scheduled = append(i.scheduled, &scheduledFailure{
			failure:   failure,
			executeAt: executeAt.Add(failure.Duration),
			isRecover: true,
		})
	}
}

// Start starts the failure injection scheduler
func (i *Injector) Start() {
	i.mu.Lock()
	i.startTime = time.Now()
	i.running = true
	i.mu.Unlock()

	go i.runScheduler()
}

// Stop stops the failure injection scheduler
func (i *Injector) Stop() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.running = false
}

// runScheduler runs the failure scheduler
func (i *Injector) runScheduler() {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		i.mu.RLock()
		running := i.running
		i.mu.RUnlock()

		if !running {
			return
		}

		now := time.Now()
		i.mu.Lock()
		toExecute := make([]*scheduledFailure, 0)
		remaining := make([]*scheduledFailure, 0)

		for _, sf := range i.scheduled {
			if now.After(sf.executeAt) || now.Equal(sf.executeAt) {
				toExecute = append(toExecute, sf)
			} else {
				remaining = append(remaining, sf)
			}
		}
		i.scheduled = remaining
		i.mu.Unlock()

		for _, sf := range toExecute {
			i.executeScheduled(sf)
		}

		<-ticker.C
	}
}

// executeScheduled executes a scheduled failure or recovery
func (i *Injector) executeScheduled(sf *scheduledFailure) {
	if sf.isRecover {
		i.executeRecovery(sf.failure)
	} else {
		i.executeFailure(sf.failure)
	}
}

// executeFailure executes a failure
func (i *Injector) executeFailure(f *Failure) {
	switch f.Type {
	case FailureCrash:
		i.InjectCrash(f.Target)
	case FailurePartition:
		from := f.Params["from"].(string)
		to := f.Params["to"].(string)
		bidir := false
		if b, ok := f.Params["bidirectional"].(bool); ok {
			bidir = b
		}
		i.InjectPartition(from, to, bidir)
	case FailureDelay:
		if i.nodeManager != nil {
			delay := f.Params["delay"].(time.Duration)
			i.nodeManager.SetNodeDelay(f.Target, delay)
		}
	}
}

// executeRecovery executes a recovery
func (i *Injector) executeRecovery(f *Failure) {
	switch f.Type {
	case FailureCrash:
		i.RecoverNode(f.Target)
	case FailurePartition:
		from := f.Params["from"].(string)
		to := f.Params["to"].(string)
		bidir := false
		if b, ok := f.Params["bidirectional"].(bool); ok {
			bidir = b
		}
		i.HealPartition(from, to, bidir)
	case FailureDelay:
		if i.nodeManager != nil {
			i.nodeManager.ClearNodeDelay(f.Target)
		}
	}
}

// GetActiveFailures returns all active failures
func (i *Injector) GetActiveFailures() []*Failure {
	i.mu.RLock()
	defer i.mu.RUnlock()

	failures := make([]*Failure, 0, len(i.failures))
	for _, f := range i.failures {
		if f.Active {
			failures = append(failures, f)
		}
	}
	return failures
}

// ClearAll clears all active failures
func (i *Injector) ClearAll() {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, f := range i.failures {
		if f.Active {
			switch f.Type {
			case FailureCrash:
				if i.nodeManager != nil {
					i.nodeManager.RecoverNode(f.Target)
				}
			case FailurePartition:
				if i.networkManager != nil {
					from := f.Params["from"].(string)
					to := f.Params["to"].(string)
					i.networkManager.HealPartition(from, to)
					if bidir, ok := f.Params["bidirectional"].(bool); ok && bidir {
						i.networkManager.HealPartition(to, from)
					}
				}
			}
		}
	}

	i.failures = make(map[string]*Failure)
	i.scheduled = make([]*scheduledFailure, 0)
}

var idCounter int
var idMu sync.Mutex

func generateID() string {
	idMu.Lock()
	defer idMu.Unlock()
	idCounter++
	return "failure-" + string(rune('a'+idCounter%26)) + time.Now().Format("150405")
}
