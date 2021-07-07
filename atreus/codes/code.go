package codes

import (
	"sync"

	"github.com/golang/protobuf/ptypes/any"
	pbstatus "google.golang.org/genproto/googleapis/rpc/status"
	rpcodes "google.golang.org/grpc/codes"
)

var (
	GlobalAtreusCodesStore sync.Map //全局错误码数组
	//
	PB_SUCC = &pbstatus.Status{
		Code:    int32(rpcodes.OK),
		Message: "SUCC",
		Details: make([]*any.Any, 0),
	}

	PB_STAT_SUCC = &PbStatusWrapper{
		Status: PB_SUCC,
	}
)
