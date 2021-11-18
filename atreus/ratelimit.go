package atreus

//提供官方golang.org/x/time/rate的令牌桶限流中间件

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc-wrapper-framework/common/pyerrors"
)

// 提供统一的limiter接口
type Limiter interface {
	Allow(method string) bool
}

// UnaryServerInterceptor returns a new unary server interceptors that performs request rate limiting.
func (s *Server) Limit(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Allow(info.FullMethod) {
			if s.Proba.TrueOrNot() {
				s.Logger.Error("Limit exceed", zap.String("method", info.FullMethod))
			}
			//在触发RPC调用前就return了，所以其他需要捕获错误的中间件需要设置在limiter之前
			//return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by ratelimit middleware", info.FullMethod)
			//for short metrics：atreusns_atreusss_server_counter_total{code="ErrRatelimit",method="/proto.GreeterService/SayHello",type="unary"} 2
			return nil, status.Error(codes.ResourceExhausted, pyerrors.RatelimiterServiceReject)
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func (s *Server) LimitStream(limiter Limiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if limiter.Allow(info.FullMethod) {
			if s.Proba.TrueOrNot() {
				s.Logger.Error("Limit exceed", zap.String("method", info.FullMethod))
			}
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by ratelimit middleware.", info.FullMethod)
		}
		return handler(srv, stream)
	}
}

// 提供基础xrate的限速实现
type XRateLimiter struct {
	RateStore  map[string]*rate.Limiter //按照RPC-method限流
	LogTime    int64
	Rate       rate.Limit
	BucketSize int
}

func NewXRateLimiter(rates rate.Limit, size int) *XRateLimiter {
	return &XRateLimiter{
		RateStore:  make(map[string]*rate.Limiter),
		LogTime:    time.Now().UnixNano(),
		Rate:       rates,
		BucketSize: size,
	}
}

//true：限速，请求丢弃
//false：请求放过
func (x *XRateLimiter) Allow(method string) bool {
	if _, exists := x.RateStore[method]; exists {
		//return !x.RateStore[method].Allow()
	} else {
		x.RateStore[method] = rate.NewLimiter(x.Rate, x.BucketSize)
		//return !x.RateStore[method].Allow()
	}

	return !x.RateStore[method].Allow()
}
