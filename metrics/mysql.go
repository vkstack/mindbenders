package metrics

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const defaultDuration = 5 * time.Second

var dbcollector *prometheus.CounterVec

func dbpoolcollector() {
	dbcollector = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "conn_wait_count",
		Subsystem: "dotpe",
		Namespace: "mysql",
		Help:      "Mysql Pool Connection request wait count",
	}, []string{"db"})
}

type DBCollectorLables struct {
	//db type... ie. payment-master | payment-slave
	Type string
}

var oncedbpoolcollectorInit sync.Once

func WithDBPoolMetrics(db interface{ Stats() sql.DBStats }, d time.Duration, lablevalue DBCollectorLables) {
	oncedbpoolcollectorInit.Do(dbpoolcollector)
	if d.Seconds() < 10 {
		d = defaultDuration
	}
	var waitCount int64
	go func() {
		for {
			time.Sleep(d)
			stats := db.Stats()
			dbcollector.WithLabelValues(lablevalue.Type).Add(float64(stats.WaitCount - waitCount))
			waitCount = stats.WaitCount
		}
	}()
}
