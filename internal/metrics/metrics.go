package metrics

import (
	"expvar"
	"homework-1/internal/metrics/counters"
)

type Metrics struct {
	IncomingRequestCounter     counters.Counter
	OutgoingRequestCounter     counters.Counter
	SuccessfulRequestCounter   counters.Counter
	UnsuccessfulRequestCounter counters.Counter
	FailedRequestCounter       counters.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		IncomingRequestCounter:     counters.NewIntCounter(),
		OutgoingRequestCounter:     counters.NewIntCounter(),
		SuccessfulRequestCounter:   counters.NewIntCounter(),
		UnsuccessfulRequestCounter: counters.NewIntCounter(),
		FailedRequestCounter:       counters.NewIntCounter(),
	}
}

func (m *Metrics) Publish() {
	expvar.Publish("IncomingRequestCounter", m.IncomingRequestCounter)
	expvar.Publish("OutgoingRequestCounter", m.OutgoingRequestCounter)
	expvar.Publish("SuccessfulRequestCounter", m.SuccessfulRequestCounter)
	expvar.Publish("UnsuccessfulRequestCounter", m.UnsuccessfulRequestCounter)
	expvar.Publish("FailedRequestCounter", m.FailedRequestCounter)
}
