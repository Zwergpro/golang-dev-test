package counters

import (
	"strconv"
	"sync/atomic"
)

type IncomingRequestCounter struct {
	cnt uint64
}

func (c *IncomingRequestCounter) Inc() {
	atomic.AddUint64(&c.cnt, 1)
}

func (c *IncomingRequestCounter) String() string {
	return strconv.FormatUint(c.cnt, 10)
}
