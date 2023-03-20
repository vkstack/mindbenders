package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type profOpt func(*ginprometheus.Prometheus)

func WithRouter(router *gin.Engine) profOpt {
	return func(p *ginprometheus.Prometheus) {
		p.ReqCntURLLabelMappingFn = func(c *gin.Context) string { return c.FullPath() }
		router.Use(customCollector(p))
		p.Use(router)
	}
}

func AttachPrometheusExporter(opts ...profOpt) {
	p := ginprometheus.NewPrometheus("gin")
	for _, opt := range opts {
		opt(p)
	}
}

func SetPrometheusMetricsOnGin(router *gin.Engine) {
	p := ginprometheus.NewPrometheus("gin")
	WithRouter(router)(p)
}

func customCollector(p *ginprometheus.Prometheus) gin.HandlerFunc {
	CustomCollectors := initializeCollector()
	return func(c *gin.Context) {
		if c.Request.URL.String() == p.MetricsPath {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		for _, col := range CustomCollectors {
			var lables []string
			for _, arg := range col.Args {
				if cb, ok := argFuncs[arg]; ok {
					lables = append(lables, cb(c))
				} else {
					lables = append(lables, "-")
				}
			}
			if summaryCollector, ok := col.MetricCollector.(prometheus.SummaryVec); ok {
				summaryCollector.WithLabelValues(lables...).Observe(time.Since(start).Seconds())
			}

			switch t := col.MetricCollector.(type) {
			case *prometheus.SummaryVec:
				t.WithLabelValues(lables...).Observe(time.Since(start).Seconds())
			case *prometheus.HistogramVec:
				t.WithLabelValues(lables...).Observe(time.Since(start).Seconds())
			default:
				fmt.Println("unknown collector", t)
			}
		}
	}
}

func initializeCollector() []*ginprometheus.Metric {
	CustomCollectors := []*ginprometheus.Metric{
		{
			Name:        "custom_request_latency_summary",
			Description: "The HTTP request latencies in seconds per api. Summary",
			Type:        "summary_vec",
			Args:        []string{"code", "method", "handler", "host", "url"},
		},
		{
			Name:        "custom_request_latency_histo",
			Description: "The HTTP request latencies in seconds per api. Histograms",
			Type:        "histogram_vec",
			Args:        []string{"code", "method", "handler", "host", "url"},
		},
	}
	for _, m := range CustomCollectors {
		m.MetricCollector = ginprometheus.NewMetric(m, "dotpe")
		if err := prometheus.Register(m.MetricCollector); err != nil {
			fmt.Println(err, "\n", m)
		}
	}
	return CustomCollectors
}

var argFuncs = map[string]func(c *gin.Context) string{
	"code":    func(c *gin.Context) string { return strconv.Itoa(c.Writer.Status()) },
	"method":  func(c *gin.Context) string { return c.Request.Method },
	"handler": func(c *gin.Context) string { return c.HandlerName() },
	"host":    func(c *gin.Context) string { return c.Request.Host },
	"url":     func(c *gin.Context) string { return c.FullPath() },
}
