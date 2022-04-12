package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "auth_service"

const (
	fieldMethodName = "grpc_method"
)

var requestTimeHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Subsystem: namespace,
	Name:      "request_duration_seconds",
	Help:      "Request duration per grpc method.",
}, []string{fieldMethodName})
