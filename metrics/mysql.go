package metrics

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const defaultDuration = 5 * time.Second

var dbcollector, dbWaitdurationCollector *prometheus.CounterVec

func dbpoolcollector() {
	dbcollector = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "conn_wait_count",
		Subsystem: "dotpe",
		Namespace: "mysql",
		Help:      "Mysql Pool Connection request wait count",
	}, []string{"db", "app"})

	dbWaitdurationCollector = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "conn_wait_duration",
		Subsystem: "dotpe",
		Namespace: "mysql",
		Help:      "Mysql Pool Connection request wait duration",
	}, []string{"db", "app"})
}

type DBCollectorLables struct {
	//db type... ie. payment-master | payment-slave
	Type string
	App  string
}

var oncedbpoolcollectorInit sync.Once

func WithDBPoolMetrics(db *sql.DB, d time.Duration, lablevalue DBCollectorLables) {
	oncedbpoolcollectorInit.Do(dbpoolcollector)
	if d.Seconds() < 10 {
		d = defaultDuration
	}
	var waitCount int64
	var waitDurTillnow time.Duration
	go func() {
		for {
			time.Sleep(d)
			stats := db.Stats()
			dbcollector.WithLabelValues(lablevalue.Type, lablevalue.App).Add(float64(stats.WaitCount - waitCount))
			dbWaitdurationCollector.WithLabelValues(lablevalue.Type, lablevalue.App).Add((stats.WaitDuration - waitDurTillnow).Seconds())
			waitDurTillnow = stats.WaitDuration
			waitCount = stats.WaitCount
		}
	}()
}
