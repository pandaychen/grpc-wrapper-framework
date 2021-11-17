package atreus

//重试retry拦截器

import (
	"context"
	"strconv"

	//"grpc-wrapper-framework/microservice/retrys"
	vctx "grpc-wrapper-framework/common/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 服务端重试检测：检测ctx中的重传次数是否满足服务端限制
func (s *Server) RetryChecking() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var wmd WMetadata
		ok := wmd.FromIncoming(ctx) //get ctx
		if ok {
			attemptArr := wmd.Get(vctx.CtxAttemptKey)
			if attemptArr != "" {
				if attempt_int, err := strconv.Atoi(attemptArr); err == nil {
					if attempt_int > s.MaxRetry {
						s.Logger.Error("[Server]RetryChecking limit")
						return nil, status.Error(codes.FailedPrecondition, "Max Retries exceeded")
					}
				}
			}
		}

		//not found ctxkey
		return handler(ctx, req)
	}
}
