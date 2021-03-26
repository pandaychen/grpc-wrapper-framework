package codes

//GRPC的错误码封装

import (
	//"google.golang.org/grpc/status"
	//https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/status
	"github.com/golang/protobuf/ptypes/any"
	pbstatus "google.golang.org/genproto/googleapis/rpc/status"
)

type PbStatusWrapper struct {
	*pbstatus.Status
}

func NewPbStatusWrapper(code int, errmsg string) *PbStatusWrapper {
	return &PbStatusWrapper{
		&pbstatus.Status{
			Code:    int32(code),
			Message: errmsg,
			Details: make([]*any.Any, 0),
		},
	}
}

func (p *PbStatusWrapper) GetCode2Int() int {
	return int(p.Code)
}

func (p *PbStatusWrapper) GetCode2Uint32() uint32 {
	return uint32(p.Code)
}
