package atreus

import (
	//auth "grpc-wrapper-framework/microservice/authentication"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Authorize ：Server端认证的一元拦截器
func (s *Server) Authorize() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//auth
		if err := s.Auther.Authenticate(ctx); err != nil {
			// TODO：需要区分逻辑错误不应该成为客户端熔断机制触发的错误
			s.Logger.Error("[Server] Authorize error", zap.Any("errmsg", err))
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(ctx, req)
	}
}
