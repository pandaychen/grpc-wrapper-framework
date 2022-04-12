package errcode

//对外部提供生成Status对象的方法！

import (
	"fmt"
	"strconv"

	"grpc-wrapper-framework/errcode/types"

	"github.com/golang/protobuf/proto"
)

// 根据int和message构造 Status对象
func Error(code Code, message string) *Status {
	return &Status{s: &types.Status{Code: int32(code.Code()), Message: message}}
}

// Errorf new status with code and message
func Errorf(code Code, format string, args ...interface{}) *Status {
	return Error(code, fmt.Sprintf(format, args...))
}

// FromCode create status from ecode
func FromCode(code Code) *Status {
	return &Status{s: &types.Status{Code: int32(code), Message: code.Message()}}
}

// FromProto new status from grpc detail
func FromProto(pbMsg proto.Message) Codes {
	//尝试解析pbMsg，pbMsg的生成方式如下面所示
	//	err, _ := errcode.Error(errcode.AccessDenied, "AccessDenied").WithDetails(&pb.HelloReply{Success: true, Message: "this is test detail"})
	if msg, ok := pbMsg.(*types.Status); ok {
		if msg.Message == "" || msg.Message == strconv.FormatInt(int64(msg.Code), 10) {
			// NOTE: if message is empty convert to pure Code, will get message from config center.
			// 当msg.Message的字符串为空，或者是纯数字（错误码）的时候，重新构造Codes类型返回
			return Code(msg.Code)
		}
		return &Status{s: msg}
	}
	return Errorf(ServerErr, "invalid proto message get %v", pbMsg)
}
