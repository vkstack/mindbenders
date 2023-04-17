package logging

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (dlogger *dlogger) initExporter() {
	dlogger.collector = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "counter",
		Subsystem: "logs",
		Namespace: "dotpe",
		Help:      "log counts from application",
	}, []string{"level", "msg", "caller"})
	prometheus.Register(dlogger.collector)
}

func (dlogger *dlogger) addMetrics(level, msg, caller string) {
	parts := strings.Split(caller, "/")
	parts = parts[len(parts)-3:]
	caller = strings.Join(parts, "/")
	dlogger.collector.WithLabelValues(level, msg, caller).Add(1)
}
