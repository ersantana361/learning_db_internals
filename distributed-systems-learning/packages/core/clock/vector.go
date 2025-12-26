package clock

import "sync"

// CausalRelation represents the causal relationship between two events
type CausalRelation int

const (
	// HappensBefore indicates a causally precedes b
	HappensBefore CausalRelation = iota
	// HappensAfter indicates b causally precedes a
	HappensAfter
	// Concurrent indicates neither event causally precedes the other
	Concurrent
	// Equal indicates the clocks are identical
	Equal
)

// VectorClock implements a vector clock for detecting causality
type VectorClock struct {
	mu     sync.RWMutex
	nodeID string
	clock  map[string]uint64
}

// NewVectorClock creates a new vector clock for the given node
func NewVectorClock(nodeID string, allNodes []string) *VectorClock {
	vc := &VectorClock{
		nodeID: nodeID,
		clock:  make(map[string]uint64),
	}
	for _, n := range allNodes {
		vc.clock[n] = 0
	}
	return vc
}

// NodeID returns the ID of the node this clock belongs to
func (vc *VectorClock) NodeID() string {
	return vc.nodeID
}

// Time returns a copy of the current clock values
func (vc *VectorClock) Time() map[string]uint64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.copy()
}

// Get returns the clock value for a specific node
func (vc *VectorClock) Get(nodeID string) uint64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.clock[nodeID]
}

// Increment increments this node's clock component
// Called before sending a message or on local events
func (vc *VectorClock) Increment() map[string]uint64 {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.clock[vc.nodeID]++
	return vc.copy()
}

// Merge merges a received vector clock with the local clock
// Sets each component to max(local, received), then increments own component
func (vc *VectorClock) Merge(received map[string]uint64) map[string]uint64 {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	for k, v := range received {
		if v > vc.clock[k] {
			vc.clock[k] = v
		}
	}
	vc.clock[vc.nodeID]++
	return vc.copy()
}

// Compare determines the causal relationship between this clock and another
func (vc *VectorClock) Compare(other map[string]uint64) CausalRelation {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return CompareVectorClocks(vc.clock, other)
}

// CompareVectorClocks compares two vector clocks
func CompareVectorClocks(a, b map[string]uint64) CausalRelation {
	aLessOrEqual := true
	bLessOrEqual := true
	equal := true

	// Collect all keys from both maps
	allKeys := make(map[string]bool)
	for k := range a {
		allKeys[k] = true
	}
	for k := range b {
		allKeys[k] = true
	}

	for k := range allKeys {
		aVal := a[k]
		bVal := b[k]

		if aVal != bVal {
			equal = false
		}
		if aVal > bVal {
			bLessOrEqual = false
		}
		if bVal > aVal {
			aLessOrEqual = false
		}
	}

	if equal {
		return Equal
	}
	if aLessOrEqual && !bLessOrEqual {
		return HappensBefore // a -> b
	}
	if bLessOrEqual && !aLessOrEqual {
		return HappensAfter // b -> a
	}
	return Concurrent
}

// HappensBefore returns true if this clock causally precedes other
func (vc *VectorClock) HappensBefore(other map[string]uint64) bool {
	return vc.Compare(other) == HappensBefore
}

// IsConcurrent returns true if this clock is concurrent with other
func (vc *VectorClock) IsConcurrent(other map[string]uint64) bool {
	return vc.Compare(other) == Concurrent
}

// copy returns a copy of the internal clock map (must be called with lock held)
func (vc *VectorClock) copy() map[string]uint64 {
	result := make(map[string]uint64, len(vc.clock))
	for k, v := range vc.clock {
		result[k] = v
	}
	return result
}

// Clone creates an independent copy of the vector clock
func (vc *VectorClock) Clone() *VectorClock {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	clone := &VectorClock{
		nodeID: vc.nodeID,
		clock:  make(map[string]uint64, len(vc.clock)),
	}
	for k, v := range vc.clock {
		clone.clock[k] = v
	}
	return clone
}
