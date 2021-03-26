package atreus

import (
	"context"

	"github.com/pandaychen/grpc-wrapper-framework/atreus/codes"
	"github.com/pandaychen/grpc-wrapper-framework/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
