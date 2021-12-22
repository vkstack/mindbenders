package metrics

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const defaultDuration = 5 * time.Second

func WithDBPoolMetrics(db interface{ Stats() sql.DBStats }, d time.Duration, lablevalue map[string]string) {
	if d.Seconds() > 5 {
		d = defaultDuration
	}
	var waitCount int64
	v1 := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mysql_conn_wait_count",
		Help: "",
	}, nil)
	go func() {
		for {
			time.Sleep(d)
			stats := db.Stats()
			v1.With(lablevalue).Add(float64(stats.WaitCount - waitCount))
			waitCount = stats.WaitCount
		}
	}()
}
