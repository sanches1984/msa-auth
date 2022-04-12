package metrics

import (
	"context"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type Service struct {
	httpServer  *http.Server
	grpcMetrics *grpc_prometheus.ServerMetrics
}

func NewService(addr string) *Service {
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcMetrics)
	reg.MustRegister(requestTimeHist)

	return &Service{
		httpServer:  &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: addr},
		grpcMetrics: grpcMetrics,
	}
}

func (l *Service) Initialize(server *grpc.Server) {
	l.grpcMetrics.InitializeMetrics(server)
	requestTimeHist.GetMetricWithLabelValues(fieldMethodName)

}

func (l Service) GRPCMetricsInterceptor() grpc.UnaryServerInterceptor {
	return l.grpcMetrics.UnaryServerInterceptor()
}

func (l Service) AppMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ts := time.Now()
		resp, err := handler(ctx, req)
		requestTimeHist.WithLabelValues(info.FullMethod).Observe(time.Now().Sub(ts).Seconds())
		return resp, err
	}
}

func (l *Service) Listen() error {
	return l.httpServer.ListenAndServe()
}
