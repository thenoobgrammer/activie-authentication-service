package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests made.",
		},
		[]string{"method", "endpoint", "status"},
	)
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

var (
	DBQueryFailuresTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_failures_total",
			Help: "Total number of failed database queries.",
		},
		[]string{"query_type", "error_type"},
	)
	DBQuerySuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_success_total",
			Help: "Total number of successful database queries.",
		},
		[]string{"query_type", "success_type"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(HTTPRequestsTotal, HTTPRequestDuration)
	prometheus.MustRegister(DBQueryFailuresTotal, DBQuerySuccessTotal)
}
