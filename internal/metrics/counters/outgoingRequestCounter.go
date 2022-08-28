package counters

import (
	"strconv"
	"sync/atomic"
)

type OutgoingRequestCounter struct {
	cnt uint64
}

func (c *OutgoingRequestCounter) Inc() {
	atomic.AddUint64(&c.cnt, 1)
}

func (c *OutgoingRequestCounter) String() string {
	return strconv.FormatUint(c.cnt, 10)
}
