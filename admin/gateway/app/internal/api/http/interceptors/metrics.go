package interceptors

import (
	"example/admin/gateway/internal/api/http/kernel"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
	"time"
)

func Metrics() func(http.Handler) http.Handler {
	promReqCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_server_handled_total",
		Help: "Общее число HTTP‑запросов",
	}, []string{"method", "path", "status"})

	promReqDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_server_handling_seconds",
		Help:    "Время обработки HTTP‑запроса (секунды)",
		Buckets: prometheus.DefBuckets,
	})

	prometheus.DefaultRegisterer.MustRegister(promReqCounter, promReqDuration)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			end := time.Now()

			status := strconv.Itoa(w.(*kernel.ResponseWriter).Status())
			promReqCounter.WithLabelValues(r.Method, r.URL.Path, status).Inc()

			durSec := end.Sub(start).Seconds()
			promReqDuration.Observe(durSec)
		})
	}
}
