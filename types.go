package metrics

// Metric contains name and value.
type Metric struct {
	Name  string
	Value string
}

// Collector return list of collected metrics.
type Collector interface {
	Metrics() []Metric
}

// Runnable represents any type with Run method.
type Runnable interface {
	Run()
}

// Stoppable represents any type with Stop method that can return an error.
type Stoppable interface {
	Stop() error
}

// Transport is an arbitrary object that can be started and stopped only once per object lifetime.
type Transport interface {
	Runnable
	Stoppable
	WithCollector(Collector) Transport
}

// ErrorLogger can be used to output errors.
type ErrorLogger interface {
	Errorf(format string, args ...interface{})
}
