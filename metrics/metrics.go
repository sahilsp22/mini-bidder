package metrics

import(
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Counter prometheus.Counter
type Opts prometheus.CounterOpts
// type CounterVec *prometheus.CounterVec

func NewCounter(opts Opts) Counter {
	return promauto.NewCounter(prometheus.CounterOpts(opts))
}

func NewCounterVec(opts Opts, labels []string) *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts(opts),labels)
}