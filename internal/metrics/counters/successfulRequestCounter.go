package counters

import (
	"strconv"
	"sync/atomic"
)

type SuccessfulRequestCounter struct {
	cnt uint64
}

func (c *SuccessfulRequestCounter) Inc() {
	atomic.AddUint64(&c.cnt, 1)
}

func (c *SuccessfulRequestCounter) String() string {
	return strconv.FormatUint(c.cnt, 10)
}
