package atreus

import (
	"github.com/renstrom/shortuuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type globalReqIDKey struct{}

var DefaultAtreusReqIDKey = "atreus-requestid"
var DefaultAtreusReqIDName = "atreus-reqid-name"

type AtreusReqId struct {
	Name string
	Id   string
}

func NewAtreusReqId(name string) *AtreusReqId {
	return &AtreusReqId{
		Name: name,
		Id:   shortuuid.New(),
	}
}

//校验reqid
func (d *AtreusReqId) Validate() bool {
	//todo: verify logic
	return true
}

func generateReqIdFromCtx(ctx context.Context) *AtreusReqId {
	mdict, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}

	header, ok := mdict[DefaultAtreusReqIDKey]
	if !ok || len(header) == 0 {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}

	reqID := header[0].(AtreusReqId)
	if reqID.Id == "" {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}

	if reqID.Validate() {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}

	return &reqID
}
