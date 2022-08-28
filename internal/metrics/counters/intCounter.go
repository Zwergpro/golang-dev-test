package counters

import (
	"strconv"
	"sync/atomic"
)

type IntCounter struct {
	cnt uint64
}

func (c *IntCounter) Inc() {
	atomic.AddUint64(&c.cnt, 1)
}

func (c *IntCounter) String() string {
	return strconv.FormatUint(c.cnt, 10)
}

func NewIntCounter() *IntCounter {
	return &IntCounter{}
}
