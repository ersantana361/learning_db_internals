package clock

import "sync"

// LamportClock implements a Lamport logical clock
type LamportClock struct {
	mu   sync.RWMutex
	time uint64
}

// NewLamportClock creates a new Lamport clock starting at 0
func NewLamportClock() *LamportClock {
	return &LamportClock{time: 0}
}

// Time returns the current clock value
func (c *LamportClock) Time() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.time
}

// Increment advances the clock by 1 and returns the new value
// Called before sending a message or on local events
func (c *LamportClock) Increment() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.time++
	return c.time
}

// Update updates the clock based on a received message timestamp
// Sets clock to max(local, received) + 1
func (c *LamportClock) Update(received uint64) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if received > c.time {
		c.time = received
	}
	c.time++
	return c.time
}

// Compare compares two Lamport timestamps
// Returns:
//   -1 if a happens-before b (a < b)
//    1 if b happens-before a (a > b)
//    0 if they are concurrent or equal
//
// Note: Lamport clocks can only determine happens-before if a < b,
// but a < b does NOT imply a happens-before b (they could be concurrent)
func Compare(a, b uint64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
