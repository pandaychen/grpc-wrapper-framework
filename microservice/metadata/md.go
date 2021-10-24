package metadata

import (
	"context"
	"errors"
	"strconv"
)

//from kratos，这里用于传输 grpc 头部的各类字段

// metadata const key
const (

	// Network
	RemoteIPKey   = "remote_ip"
	RemotePortKey = "remote_port"
	ServerAddrKey = "server_addr"
	ClientAddrKey = "client_addr"

	CPUloadKey = "cpu_usage"

	CallerKey = "caller"
)

//for client md
var outgoingKeyMap = map[string]struct{}{
	RemoteIPKey:   struct{}{},
	RemotePortKey: struct{}{},
}

//for server md
var incomingKeyMap = map[string]struct{}{
	CallerKey: struct{}{},
}

func GetIncomingMedataMap(t bool) map[string]struct{} {
	if t {
		return incomingKeyMap
	}
	return outgoingKeyMap
}

// IsOutgoingKey represent this key should propagate by rpc.
func IsOutgoingKey(key string) bool {
	_, ok := outgoingKeyMap[key]
	return ok
}

// IsIncomingKey represent this key should extract from rpc metadata.
func IsIncomingKey(key string) (ok bool) {
	_, ok = outgoingKeyMap[key]
	if ok {
		return
	}
	_, ok = incomingKeyMap[key]
	return
}

//////////////////////

// MD is a mapping from metadata keys to values.
type XMetaData map[string]interface{}

type DefaultMdKey struct{}

//pointer
func NewXMetaData() XMetaData {
	return XMetaData{}
}

func (m XMetaData) Len() int {
	return len(m)
}

// Copy returns a copy of md.
func (m XMetaData) Copy() XMetaData {
	return New(m)
}

// New creates an MD from a given key-value map.
func New(m map[string]interface{}) XMetaData {
	newm := XMetaData{}
	for k, v := range m {
		newm[k] = v
	}
	return newm
}

func Combine(mds ...XMetaData) XMetaData {
	out := XMetaData{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = v
		}
	}
	return out
}

//create mds from kv list
func Pairs(kvlist ...interface{}) XMetaData {
	if len(kvlist)%2 != 0 {
		panic(errors.New("Len must be Even"))
	}
	nm := XMetaData{}
	var (
		savekey string
		ok      bool
	)
	for index, k := range kvlist {
		if index%2 == 0 {
			savekey, ok = k.(string)
			if !ok {
				continue
			}
			continue
		}
		if savekey != "" {
			nm[savekey] = k
		}

	}
	return nm
}

// 给 ctx 加入 md 的 kv 传递，返回一个子 ctx
func NewContext(ctx context.Context, md XMetaData) context.Context {
	// 以 DefaultMdKey{} 为 key
	// func WithValue(parent Context, key, val interface{}) Context
	return context.WithValue(ctx, DefaultMdKey{}, md)
}

// FromContext returns the incoming metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromContext(ctx context.Context) (XMetaData, bool) {
	md, ok := ctx.Value(DefaultMdKey{}).(XMetaData)
	return md, ok
}

// get value from metadata in context return nil if not found
func Value(ctx context.Context, key string) interface{} {
	md, ok := ctx.Value(DefaultMdKey{}).(XMetaData)
	if !ok {
		return nil
	}
	return md[key] //is not found,return nil
}

// get string value from metadata in context
func StringValctx(ctx context.Context, key string) string {
	md, ok := ctx.Value(DefaultMdKey{}).(XMetaData)
	if !ok {
		return ""
	}
	s, ok := md[key].(string)
	if !ok {
		return ""
	}
	return s
}

// get int64 value from metadata in context
func Int64Valctx(ctx context.Context, key string) int64 {
	md, ok := ctx.Value(DefaultMdKey{}).(XMetaData)
	if !ok {
		return -1
	}
	i, ok := md[key].(int64)
	if !ok {
		return -1
	}
	return i
}

// get boolean from metadata in context use strconv.Parse.
func BoolValctx(ctx context.Context, key string) bool {
	md, ok := ctx.Value(DefaultMdKey{}).(XMetaData)
	if !ok {
		return false
	}

	switch md[key].(type) {
	case bool:
		return md[key].(bool)
	case string:
		ok, err := strconv.ParseBool(md[key].(string))
		if err != nil {
			return false
		}
		return ok
	default:
		return false
	}
}
