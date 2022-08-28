package counters

import (
	"strconv"
	"sync/atomic"
)

type FailedRequestCounter struct {
	cnt uint64
}

func (c *FailedRequestCounter) Inc() {
	atomic.AddUint64(&c.cnt, 1)
}

func (c *FailedRequestCounter) String() string {
	return strconv.FormatUint(c.cnt, 10)
}
