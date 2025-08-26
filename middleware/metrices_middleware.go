package middleware
import (
	"github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"

)

var (
	// Define a histogram metric to track request durations
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
		},
		[]string{"path", "method"},
	)
	statusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "Total number of HTTP responses by status code",
		},
		[]string{"path", "method", "status_code"},
	)
)
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func init() {
	// Register the metrics with Prometheus's default registry
	prometheus.MustRegister(requestCounter, requestDuration,statusCounter)
}

func MetricMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Increment the request counter
		ww := &responseWriter{ResponseWriter: w} 
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()
		requestCounter.WithLabelValues(r.URL.Path, r.Method).Inc()
		requestDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
		statusCounter.WithLabelValues(r.URL.Path, r.Method, http.StatusText(ww.statusCode)).Inc()

	})
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}