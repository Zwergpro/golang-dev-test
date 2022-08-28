package counters

type Counter interface {
	Inc()
	String() string
}
