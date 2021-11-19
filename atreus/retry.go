package atreus

//重试retry拦截器

import (
	"context"
	"strconv"

	actx "grpc-wrapper-framework/common/context"
	"grpc-wrapper-framework/logger"
	"grpc-wrapper-framework/microservice/retrys"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 服务端重试检测：检测ctx中的重传次数是否满足服务端限制
func (s *Server) RetryChecking() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var (
			wmd      WMetadata
			wmd_copy *WMetadata
		)
		ok := wmd.FromIncoming(ctx) //get ctx
		if ok {
			wmd_copy = wmd.Copy() //Fix：we use a copy of wmd
			attemptArr := wmd_copy.Get(actx.CtxAttemptKey)
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

//客户端重试
func (c *Client) DoClientRetry(optFuncs ...retrys.CallOption) grpc.UnaryClientInterceptor {
	if c.MaxRetry == 0 {
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	//初始化应用retrys配置
	intOpts := retrys.ReuseOrNewWithCallOptions(retrys.DefaultOptions, optFuncs)
	return func(parentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var (
			lastErr error
		)
		//注意第一个参数为上一级context
		grpcOpts, retryOpts := retrys.SplitCallOptions(opts)

		//嵌入logger
		parentCtx = logger.NewContext(parentCtx, c.Logger)

		//更新retrys配置
		callOpts := retrys.ReuseOrNewWithCallOptions(intOpts, retryOpts)

		// 前面已判断
		if callOpts.GetMaxretries() == 0 {
			return invoker(parentCtx, method, req, reply, cc, grpcOpts...)
		}
		for attempt := 0; attempt < callOpts.GetMaxretries(); attempt++ {
			//Sleeping....
			if err := retrys.WaitRetryBackoff(parentCtx, attempt, callOpts); err != nil {
				c.Logger.Error("[Client]DoClientRetry WaitRetryBackoff error", zap.Int("attempt", attempt), zap.Any("errmsg", err))
				return err
			}

			//检查是否设置了单次请求的超时时间
			newCallerCtx := retrys.PerCallContext(parentCtx, attempt, callOpts)

			//do real rpc...
			lastErr = invoker(newCallerCtx, method, req, reply, cc, grpcOpts...)

			//rpc请求完成之后的处理

			if lastErr == nil {
				return nil
			}

			if retrys.IsContextError(lastErr) {
				if parentCtx.Err() != nil {
					//如果是上一层报错，那么直接退出
					c.Logger.Error("[Client]DoClientRetry IsContextError error", zap.Any("parent error", parentCtx.Err()), zap.Any("errmsg", lastErr))
					return lastErr
				} else if callOpts.GetperCallTimeout() != 0 {
					// We have set a perCallTimeout in the retry middleware, which would result in a context error if
					// the deadline was exceeded, in which case try again.
					continue
				}
			}
			if retrys.CheckIsRetriable(lastErr, callOpts) == false {
				//非重传类错误，直接返回
				//在server超时拦截器中返回的codes.DeadlineExceeded错误（服务端处理超时），不在此错误集合中，会跳过重传逻辑
				return lastErr
			}
			continue
		}
		return lastErr
	}
}
