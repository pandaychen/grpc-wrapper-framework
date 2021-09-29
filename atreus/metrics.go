package atreus

import (
	"context"

	"grpc-wrapper-framework/atreus/codes"
	"grpc-wrapper-framework/metrics"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	metricRpcServerReqDur = metrics.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: metrics.DefaultNamespace,
		Subsystem: metrics.DefaultSubsystem,
		Name:      "duration_ms",
		Help:      "atreus rpc server requests duration(ms)",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
)

func (s *Server) Metrics2Prometheus() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//startTime := time.Now()
		resp, err := handler(ctx, req)
		if err != nil {
			s.Logger.Error("[Metrics2Prometheus]rpc error", zap.Any("errmsg", err))
		}
		code := codes.TransformError2Codes(err)
		metrics.ServerHandleCounter.Inc(metrics.LABLES_NAME_RPCTYPE_UNARY, info.FullMethod, code.GetMessage())
		return resp, err
	}
}
