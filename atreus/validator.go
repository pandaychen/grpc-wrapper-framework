package atreus

//pb协议validator中间件

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate(all bool) error
}

type validatorLegacy interface {
	Validate() error
}

func do_validate(req interface{}) error {
	switch v := req.(type) {
	case validatorLegacy:
		//from proto/service.validator.pb.go
		//func (this *HelloRequest) Validate() error {...}
		if err := v.Validate(); err != nil {
			//if param invalid, return codes.InvalidArgument
			return status.Error(codes.InvalidArgument, err.Error())
		}
	case validator:
		//from proto/service.pb.validate.go
		/*func (m *HelloReply) Validate() error {
			return m.validate(false)
		}*/
		if err := v.Validate(false); err != nil {
			//rpc error: code = InvalidArgument desc = invalid HelloRequest.Name: value length must be at least 20 runes
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return nil
}

func (s *Server) SrvValidator() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := do_validate(req); err != nil {
			s.Logger.Error("[SrvValidator]check params error", zap.Any("param", req))
			return nil, err
		}

		//go to next interceptor
		return handler(ctx, req)
	}
}

//客户端参数校验
func (c *Client) ClientValidator() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := do_validate(req); err != nil {
			c.Logger.Error("[ClientValidator]check params error", zap.Any("param", req))
			return status.Error(codes.InvalidArgument, err.Error())
		}

		//go to next interceptor
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
