package middlewares

import (
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func Metrics() grpc.UnaryServerInterceptor {
	provider := grpcPrometheus.NewServerMetrics(
		grpcPrometheus.WithServerHandlingTimeHistogram(),
	)

	prometheus.DefaultRegisterer.MustRegister(provider)

	return provider.UnaryServerInterceptor()
}
