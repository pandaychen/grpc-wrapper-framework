package atreus

//客户端熔断拦截器

import (
	"context"
	"fmt"
	"grpc-wrapper-framework/errcode"
	"path"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 客户端熔断拦截器
func (c *Client) CircuitBreaker() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		breakerName := path.Join(cc.Target(), method)
		if _, exist := c.CbBreakerMap[breakerName]; !exist {
			c.CbBreakerMap[breakerName] = gobreaker.NewCircuitBreaker(c.CbBreakerConfig)
		}

		_, err = c.CbBreakerMap[breakerName].Execute(func() (interface{}, error) {
			err = invoker(ctx, method, req, reply, cc, opts...)
			//c.Logger.Error("[Client]CircuitBreaker call error", zap.Any("errmsg", err))
			return nil, err
		})
		if err != nil {
			// error：circuit breaker is open
			return errcode.ServiceUnavailable
		}
		return
	}
}

//check  whether or not error is acceptable，根据服务端错误的返回，来判断哪些错误才进入熔断计算逻辑
//https://grpc.github.io/grpc/core/md_doc_statuscodes.html
//https://github.com/sony/gobreaker/blob/master/gobreaker.go#L113
func (c *Client) IsBreakerNeedError(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		//属于熔断错误判断范围
		c.Logger.Error("IsBreakerNeedError need error ok", zap.Any("errmsg", err))
		return true
	default:
		return false
	}
}

func (c *Client) IsBreakerNeedErrorV2(err error) bool {
	switch status.Code(err) {
	case codes.DeadlineExceeded, codes.Internal, codes.Unavailable, codes.DataLoss:
		//属于熔断错误判断范围
		c.Logger.Error("IsBreakerNeedError need error ok", zap.Any("errmsg", err))
		return true
	default:
		return false
	}
}
