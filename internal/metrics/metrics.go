package metrics

import (
	"expvar"
	"homework-1/internal/metrics/counters"
)

type Metrics struct {
	IncomingRequestCounter   counters.Counter
	OutgoingRequestCounter   counters.Counter
	SuccessfulRequestCounter counters.Counter
	FailedRequestCounter     counters.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		IncomingRequestCounter:   &counters.IncomingRequestCounter{},
		OutgoingRequestCounter:   &counters.OutgoingRequestCounter{},
		SuccessfulRequestCounter: &counters.SuccessfulRequestCounter{},
		FailedRequestCounter:     &counters.FailedRequestCounter{},
	}
}

func (m *Metrics) Publish() {
	expvar.Publish("IncomingRequestCounter", m.IncomingRequestCounter)
	expvar.Publish("OutgoingRequestCounter", m.OutgoingRequestCounter)
	expvar.Publish("SuccessfulRequestCounter", m.SuccessfulRequestCounter)
	expvar.Publish("FailedRequestCounter", m.FailedRequestCounter)
}
