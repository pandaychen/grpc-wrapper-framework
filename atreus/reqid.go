package atreus

import (
	"fmt"

	"github.com/renstrom/shortuuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	DefaultAtreusReqIDKey  = "atreus-requestid"
	DefaultAtreusReqIDVal  = "atreus-requestid-value"
	DefaultAtreusReqIDName = "atreus-reqid-name"
)

type globalReqIDKey struct{}

var DefaultAtreusReqIDSKey = globalReqIDKey{}

func GetGlobalReqIDFromContext(ctx context.Context) string {
	id := ctx.Value(DefaultAtreusReqIDKey)
	return id.(string)
}

type AtreusReqId string

func NewAtreusReqId(name string) AtreusReqId {
	id := shortuuid.New()
	return AtreusReqId(id)

}

//校验reqid
func (d *AtreusReqId) Validate() bool {
	//todo: verify logic
	return true
}

func composeReqIdFromContext(ctx context.Context) AtreusReqId {
	//warning：FromOutgoingContext 一般用于客户端取自己的context中存储的数据
	//FromIncomingContext --用于服务端接收客户端的context
	mdict, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}

	header, ok := mdict[DefaultAtreusReqIDKey]
	if !ok || len(header) == 0 {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}

	reqID := AtreusReqId(header[0])
	if reqID == "" {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}

	if !reqID.Validate() {
		return NewAtreusReqId(DefaultAtreusReqIDName)
	}
	newid := string(NewAtreusReqId(DefaultAtreusReqIDName))
	return AtreusReqId(fmt.Sprintf("%s,%s", reqID, newid))
}

func (s *Server) AtreusXRequestId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newReqID string
		newReqID = string(composeReqIdFromContext(ctx))

		//passing reqid with ctx
		ctx = context.WithValue(ctx, DefaultAtreusReqIDKey, newReqID)
		resp, err := handler(ctx, req)
		return resp, err
	}
}

/*
func main() {
	ctxbase := context.Background()
	md := metadata.New(map[string]string{DefaultAtreusReqIDKey: DefaultAtreusReqIDVal})
	//md := metadata.Pairs(DefaultAtreusReqIDKey: DefaultAtreusReqIDVal)

	ctx := metadata.NewIncomingContext(ctxbase, md)
	mdict, ok := metadata.FromIncomingContext(ctx)
	fmt.Println(mdict, ok) //map[atreus-requestid:[atreus-requestid-value]] true

	requestID := ComposeReqIdFromContext(ctx)

	ctx = context.WithValue(ctx, globalReqIDKey{}, "12345678")
	fmt.Println(requestID, GetGlobalReqIDFromContext(ctx)) //atreus-requestid-value 12345678
}
*/
