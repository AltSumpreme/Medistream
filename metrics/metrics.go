package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
	DBLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_latency_seconds",
			Help:    "Latency of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_seconds",
			Help: "Total number of cache hits.",
		},
		[]string{"cache_key"},
	)

	CacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_seconds",
			Help: "Total number of cache misses",
		},
		[]string{"cache_key"},
	)
)

func MetricsInit() {
	prometheus.MustRegister(HTTPRequestDuration, DBLatency, CacheHits, CacheMisses)
}
