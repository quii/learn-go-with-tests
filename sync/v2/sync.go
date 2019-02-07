package v1

import "sync"

type Counter struct {
	value int
	lock sync.Mutex
}

func (c *Counter) Inc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}
