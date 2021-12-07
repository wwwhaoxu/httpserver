package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	MetricsNamespace = "httpserver"
	functionLatency  = CreateExecutionTimeMetric(MetricsNamespace,
		"Time spent.")
)

func Register() {
	err := prometheus.Register(functionLatency)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}

func NewExecutionTimer(histogram *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histogram: histogram,
		start:     now,
		last:      now,
	}
}

type ExecutionTimer struct {
	histogram *prometheus.HistogramVec
	start     time.Time
	last      time.Time
}

func (t ExecutionTimer) ObserveTotal() {
	(*t.histogram).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}
