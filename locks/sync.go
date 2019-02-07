package locks

type Counter struct {
	value int
}

func (c *Counter) Inc() {
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}
