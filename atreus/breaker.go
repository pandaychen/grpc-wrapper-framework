package atreus

//客户端熔断拦截器

import (
	"context"
	"path"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (c *Client) CircuitBreaker() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		breakerName := path.Join(cc.Target(), method)
		if _, exist := c.CbBreakerMap[breakerName]; !exist {
			c.CbBreakerMap[breakerName] = gobreaker.NewCircuitBreaker(c.CbBreakerConfig)
		}

		_, err = c.CbBreakerMap[breakerName].Execute(func() (interface{}, error) {
			err = invoker(ctx, method, req, reply, cc, opts...)
			c.Logger.Error("[Client]CircuitBreaker call error", zap.Any("errmsg", err))
			return nil, err
		})
		if err != nil {
			// error：circuit breaker is open
			//TODO: 根据服务端错误的返回，来判断哪些错误才进入熔断计算逻辑
			return err
		}
		return
	}
}
