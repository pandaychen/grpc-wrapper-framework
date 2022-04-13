package atreus

//客户端熔断拦截器

import (
	"context"
	"grpc-wrapper-framework/errcode"
	"path"

	"github.com/sony/gobreaker"
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
		//属于熔断错误判断范围，必须返回false
		//c.Logger.Error("IsBreakerNeedError need error ok", zap.Any("errmsg", err))
		return false
	default:
		return true
	}
}

//如果框架使用grpc的原生错误，那么必须使用status.Code(err)方法对errors进行转换
func (c *Client) IsBreakerNeedErrorV2(err error) bool {
	//fmt.Println(err, status.Code(err))
	if err != nil {
		if errcode.EqualError(errcode.ServerErr, err) || errcode.EqualError(errcode.ServiceUnavailable, err) || errcode.EqualError(errcode.Deadline, err) || errcode.EqualError(errcode.LimitExceed, err) {
			//触发熔断的错误
			return false
		}
	} else {
		//其他类型的错误，纳入熔断成功计数范围
		return true
	}
}
