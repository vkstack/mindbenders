package logging

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (dlogger *dlogger) initExporter() {
	dlogger.collector = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "log_count",
		Subsystem: "dotpe",
		Namespace: "logs",
		Help:      "log counts from application",
	}, []string{"level", "msg", "caller"})
	prometheus.Register(dlogger.collector)
}

func (dlogger *dlogger) addMetrics(level, msg, caller string) {
	dlogger.collector.WithLabelValues(level, msg, caller).Add(1)
}
