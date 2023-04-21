package logging

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

func WithMetric(level logrus.Level) Option {
	return func(dlogger *dlogger) {
		dlogger.collector = promauto.NewCounterVec(prometheus.CounterOpts{
			Name:      "counter",
			Subsystem: "logs",
			Namespace: "dotpe",
			Help:      "log counts from application",
		}, []string{"level", "msg", "caller"})
		prometheus.Register(dlogger.collector)
		dlogger.metricCollectionLevel = level
	}
}

func (dlogger *dlogger) addMetrics(level logrus.Level, msg, caller string) {
	if dlogger.collector != nil && level <= dlogger.metricCollectionLevel {
		parts := strings.Split(caller, "/")
		parts = parts[len(parts)-3:]
		caller = strings.Join(parts, "/")
		dlogger.collector.WithLabelValues(level.String(), msg, caller).Add(1)
	}
}
