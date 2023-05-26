package logging

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func WithMetric(level Level) Option {
	return func(dlogger *dlogger) {
		dlogger.collector = promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "counter",
			Subsystem: "logs",
			Namespace: "dotpe",
			Help:      "log counts from application",
		}, []string{"level", "caller"})
		prometheus.Register(dlogger.collector)
		dlogger.metricCollectionLevel = level
	}
}

func (dlogger *dlogger) addMetrics(level Level, caller string) {
	if dlogger.collector != nil && level < dlogger.metricCollectionLevel {
		parts := strings.Split(caller, "/")
		parts = parts[len(parts)-3:]
		caller = strings.Join(parts, "/")
		dlogger.collector.WithLabelValues(level.String(), caller).Add(1)
	}
}

// func AddMetrics(ctx context.Context, f Fields) {
// 	if c, ok := f["caller"]; ok {
// 		if c1, ok := c.(string); ok {
// 			parts := strings.Split(c1, "\n")
// 			caller := parts[len(parts)-1]
// 			dlogger.collector.WithLabelValues(level.String(), msg, caller).Add(1)
// 		}
// 	}
// }
