package atreus

import (
	//auth "grpc-wrapper-framework/microservice/authentication"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*
// AuthorizeStream：Server端Stream认证拦截器
func (s *Server) AuthorizeStream(authenticator *auth.Authenticator) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := s.Auther.Authenticate(ctx); err != nil {
			// TODO：需要区分逻辑错误不应该成为客户端熔断机制触发的错误
			s.Logger.Error("[Server] AuthorizeStream error", zap.Any("errmsg", err))
			return nil, err
		}

		return handler(srv, stream)
	}
}
*/

// Authorize ：Server端认证的一元拦截器
func (s *Server) Authorize() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//auth
		if err := s.Auther.Authenticate(ctx); err != nil {
			// TODO：需要区分逻辑错误不应该成为客户端熔断机制触发的错误
			s.Logger.Error("[Server] Authorize error", zap.Any("errmsg", err))
			return nil, err
		}

		return handler(ctx, req)
	}
}
