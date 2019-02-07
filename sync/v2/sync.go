package v1

import "sync"

// Counter will increment a number
type Counter struct {
	value int
	lock  sync.Mutex
}

// NewCounter returns a new Counter
func NewCounter() *Counter {
	return &Counter{}
}

// Inc the count
func (c *Counter) Inc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value++
}

// Value returns the current count
func (c *Counter) Value() int {
	return c.value
}
