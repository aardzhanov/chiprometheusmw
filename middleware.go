package prometheusmw

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	CounterName   = "http_request_count"
	HistogramName = "http_request_duration_seconds"
)

type Middleware struct {
	counter   *prometheus.CounterVec
	histogram *prometheus.HistogramVec
}

func NewMiddleware(name string, buckets ...float64) Middleware {

	var m Middleware
	m.counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        CounterName,
			Help:        "Count HTTP requests partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		}, []string{"code", "method", "path"})

	prometheus.MustRegister(m.counter)

	if len(buckets) == 0 {
		buckets = []float64{.025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	}

	m.histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        HistogramName,
		Help:        "Measure exec time by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	}, []string{"code", "method", "path"})

	prometheus.MustRegister(m.histogram)

	return m
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		if routeCtx := chi.RouteContext(r.Context()); routeCtx != nil {
			routePattern := routeCtx.RoutePattern()
			m.counter.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, routePattern).Inc()
			m.histogram.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, routePattern).Observe(time.Since(start).Seconds())
		}
	}
	return http.HandlerFunc(fn)
}
