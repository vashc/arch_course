package hw3

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dflBuckets = []float64{50, 500, 1000}
)

const (
	requestsTotalTag = "http_requests_total"
	errorRateTag     = "http_error_rate"
	latencyTag       = "http_request_duration_milliseconds"
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	requestsTotal *prometheus.CounterVec
	errorRate     *prometheus.CounterVec
	latency       *prometheus.HistogramVec
}

// NewPatternMiddleware returns a new prometheus Middleware handler that groups requests by the chi routing pattern.
// EX: /users/{firstName} instead of /users/bob
func NewPatternMiddleware(name string, buckets ...float64) func(next http.Handler) http.Handler {
	var m Middleware
	m.requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        requestsTotalTag,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.requestsTotal)

	m.errorRate = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: errorRateTag,
			Help: "Rate of responses with HTTP 400+ status.",
		},
		[]string{"code", "method", "path"},
	)

	if len(buckets) == 0 {
		buckets = dflBuckets
	}
	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyTag,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)
	return m.patternHandler
}

func (c Middleware) patternHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.ReplaceAll(routePattern, "/*/", "/")

		c.requestsTotal.WithLabelValues(
			http.StatusText(ww.Status()),
			r.Method,
			routePattern,
		).Inc()
		c.latency.WithLabelValues(
			http.StatusText(ww.Status()),
			r.Method,
			routePattern,
		).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)

		if ww.Status() >= 400 {
			c.errorRate.WithLabelValues(
				http.StatusText(ww.Status()),
				r.Method,
				routePattern,
			).Inc()
		}
	}
	return http.HandlerFunc(fn)
}
