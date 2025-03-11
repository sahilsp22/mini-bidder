package metrics

import(
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Counter prometheus.Counter
type Opts prometheus.CounterOpts

func NewCounter(opts Opts) Counter {
	return promauto.NewCounter(prometheus.CounterOpts(opts))
}
