package atreus

//实现Unary interceptor-chain
//参考：https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/chain.go

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*
	opt = append(opt, grpc.UnaryInterceptor(BuildUnaryInterceptorChain(Interceptor1, Interceptor2, Interceptor3, Interceptor4)))
*/
func (s *Server) BuildUnaryInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	// here returns a function grpc.UnaryServerInterceptor
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		chain := handler // 原始的RPC方法
		for i := len(interceptors) - 1; i >= 0; i-- {
			//从数组最后一个interceptor开始，依次和前一个建立chain关系
			if s.IsDebug {
				s.Logger.Info("[BuildUnaryInterceptorChain]createSubchain", zap.String("method", info.FullMethod), zap.Any("req", req), zap.Int("intercepor index", i))
			}
			chain = createSubchain(interceptors[i], chain, info)
		}

		//返回的是第0号位置上的chain
		return chain(ctx, req)
	}
}

func createSubchain(us_interceptor grpc.UnaryServerInterceptor, us_handler grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return us_interceptor(ctx, req, info, us_handler)
	}
}

// Use method attachs a global inteceptor to the server
func (s *Server) Use(handlers ...grpc.UnaryServerInterceptor) *Server {
	new_size := len(s.InnerHandlers) + len(handlers)
	if new_size >= int(vars.ATREUS_MAX_INTERCEPTOR_NUM) {
		//限制拦截器的使用上限
		panic("too many interceptors")
	}
	mergedHandlers := make([]grpc.UnaryServerInterceptor, new_size)

	//warning: Should keep interceptors order
	copy(mergedHandlers, s.InnerHandlers)

	//copy new handles
	copy(mergedHandlers[len(s.InnerHandlers):], handlers)
	s.InnerHandlers = mergedHandlers
	return s
}

func (s *Server) BuildUnaryInterceptorChain2(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		i     int
		chain grpc.UnaryHandler
	)

	if len(s.InnerHandlers) == 0 {
		// no middleware, return rpc method
		return handler(ctx, req)
	}

	chain = func(ic context.Context, ir interface{}) (interface{}, error) {
		if i == len(s.InnerHandlers)-1 {
			//拦截器数组中最后一个位置
			return handler(ic, ir)
		}
		i++
		if s.IsDebug {
			s.Logger.Info("[BuildUnaryInterceptorChain2]createSubchain", zap.String("method", args.FullMethod), zap.Any("req", ir), zap.Int("intercepor index", i))
		}
		return s.InnerHandlers[i](ic, ir, args, chain)
	}

	//返回第0号位置上的拦截器
	return s.InnerHandlers[0](ctx, req, args, chain)
}
