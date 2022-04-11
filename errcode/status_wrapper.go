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
	if msg, ok := pbMsg.(*types.Status); ok {
		if msg.Message == "" || msg.Message == strconv.FormatInt(int64(msg.Code), 10) {
			// NOTE: if message is empty convert to pure Code, will get message from config center.
			return Code(msg.Code)
		}
		return &Status{s: msg}
	}
	return Errorf(ServerErr, "invalid proto message get %v", pbMsg)
}
