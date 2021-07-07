package atreus

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 提供统一的limiter接口
type Limiter interface {
	Allow(method string) bool
}

// UnaryServerInterceptor returns a new unary server interceptors that performs request rate limiting.
func (s *Server) Limit(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Allow(info.FullMethod) {
			if s.IsDebug {
				//TODO: ADD sample log
				s.Logger.Error("Limit exceed", zap.String("method", info.FullMethod))
			}
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by ratelimit middleware", info.FullMethod)
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func (s *Server) LimitStream(limiter Limiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if limiter.Allow(info.FullMethod) {
			if s.IsDebug {
				//TODO: ADD sample log
				s.Logger.Error("Limit exceed", zap.String("method", info.FullMethod))
			}
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by ratelimit middleware.", info.FullMethod)
		}
		return handler(srv, stream)
	}
}

// 提供基础xrate的限速实现
type XRateLimiter struct {
	RateStore  map[string]*rate.Limiter
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

func (x *XRateLimiter) Allow(method string) bool {
	if _, exists := x.RateStore[method]; exists {
		//return !x.RateStore[method].Allow()
	} else {
		x.RateStore[method] = rate.NewLimiter(x.Rate, x.BucketSize)
		//return !x.RateStore[method].Allow()
	}

	return !x.RateStore[method].Allow()
}
