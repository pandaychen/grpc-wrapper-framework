package atreus

//原则：grpc最终返回的错误类型是status.Status

import (
	"context"
	"strconv"

	//"github.com/pandaychen/grpc-wrapper-framework/errcode"	//这是最通用的引用方式
	"grpc-wrapper-framework/errcode"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 转换ctx错误为grpc的标注错误
func TransContextErr2GrpcErr(err error) error {
	switch err {
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return status.Error(codes.Canceled, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

// 判断err是否为ctx错误（DeadlineExceeded || Canceled）
func IsContextError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded || code == codes.Canceled
}

// ToErrEcode convert grpc.status to ecode.Codes
// 外部接口
func ToErrEcode(gst *status.Status) errcode.Codes {
	// 获取status.Status
	details := gst.Details() //返回值是[]interface{}类型
	/*
		func (s *Status) Details() []interface{}
		Details returns a slice of details messages attached to the status. If a detail cannot be decoded, the error is returned in place of the detail.

	*/
	for _, detail := range details {
		// convert detail to status only use first detail
		if pb, ok := detail.(proto.Message); ok {
			return errcode.FromProto(pb)
		}
	}

	//否则，根据gst的错误码再次构造
	return toErrCode(gst)
}

// 将grpc的标准错误码转换为项目定义的错误码
func toErrCode(grcpStauts *status.Status) errcode.Code {
	gcode := grcpStauts.Code()
	switch gcode {
	case codes.OK:
		return errcode.OK
	case codes.InvalidArgument:
		return errcode.RequestErr
	case codes.NotFound:
		return errcode.NothingFound
	case codes.PermissionDenied:
		return errcode.AccessDenied
	case codes.Unauthenticated:
		return errcode.Unauthorized
	case codes.ResourceExhausted:
		return errcode.LimitExceed
	case codes.Unimplemented:
		return errcode.MethodNotAllowed
	case codes.DeadlineExceeded:
		return errcode.Deadline
	case codes.Unavailable:
		return errcode.ServiceUnavailable
	case codes.Unknown:
		return errcode.String(grcpStauts.Message())
	}

	//默认
	return errcode.ServerErr
}

// 将errcode转换为grpc的标准错误码
func togRPCCode(code errcode.Codes) codes.Code {
	switch code.Code() {
	case errcode.OK.Code():
		return codes.OK
	case errcode.RequestErr.Code():
		return codes.InvalidArgument
	case errcode.NothingFound.Code():
		return codes.NotFound
	case errcode.Unauthorized.Code():
		return codes.Unauthenticated
	case errcode.AccessDenied.Code():
		return codes.PermissionDenied
	case errcode.LimitExceed.Code():
		return codes.ResourceExhausted
	case errcode.MethodNotAllowed.Code():
		return codes.Unimplemented
	case errcode.Deadline.Code():
		return codes.DeadlineExceeded
	case errcode.ServiceUnavailable.Code():
		return codes.Unavailable
	}
	return codes.Unknown
}

//将errcode转换为grpcStatus（errcode.Codes是interface{}接口类型）
func gRPCStatusFromEcode(pcode errcode.Codes) (*status.Status, error) {
	var (
		st *errcode.Status
	)

	switch v := pcode.(type) {
	case *errcode.Status:
		st = v
	case errcode.Code:
		st = errcode.FromCode(v)
	default:
		//重新构造status.Status
		st = errcode.Error(errcode.Code(pcode.Code()), pcode.Message())
		for _, detail := range pcode.Details() {
			if msg, ok := detail.(proto.Message); ok {
				st.WithDetails(msg)
			}
		}
	}
	gst := status.New(codes.Unknown, strconv.Itoa(st.Code()))
	//func (s *Status) WithDetails(details ...proto.Message) (*Status, error)
	return gst.WithDetails(st.Proto())
}

// ConvertNormalError convert error for service reply and try to convert it to grpc.Status.
//
func ConvertNormalError(svrErr error) (gst *status.Status) {
	var (
		err error
	)
	//剥离，获取最原始的错误
	svrErr = errors.Cause(svrErr)
	if code, ok := svrErr.(errcode.Codes); ok {
		// TODO: deal with err
		if gst, err = gRPCStatusFromEcode(code); err == nil {
			return
		}
	}
	// for some special error convert context.Canceled to ecode.Canceled,
	// context.DeadlineExceeded to ecode.DeadlineExceeded only for raw error
	// if err be wrapped will not effect.

	//context的错误
	switch svrErr {
	case context.Canceled:
		gst, _ = gRPCStatusFromEcode(errcode.Canceled)
	case context.DeadlineExceeded:
		gst, _ = gRPCStatusFromEcode(errcode.Deadline)
	default:
		//调用默认grpc的错误封装方法
		gst, _ = status.FromError(svrErr)
	}
	return
}

// 服务端错误统一化处理
func (s *Server) TransError() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		//统一转换错误
		return resp, ConvertNormalError(err).Err()
		return
	}
}

// 客户端错误统一处理，将服务端返回的err类型（status.Status）统一转换为errcode.Codes类型
// 因为熔断器需要errcode.Codes类型
func (c *Client) TransError() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		err = invoker(ctx, method, req, reply, cc, opts...)
		//call grpc.Status package
		gst, _ := status.FromError(err)
		ec := ToErrEcode(gst)
		//是想把服务端的错误返回给被调用方
		err = errors.WithMessage(ec, gst.Message()) //将status.Status通过pkg/errors包发送给调用方
		return
	}
}
