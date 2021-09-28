package atreus

import (
	"context"
	"path"

	"github.com/sony/gobreaker"
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
			return nil, err
		})
		if err != nil {
			// error：circuit breaker is open
			return err
		}
		return
	}
}
