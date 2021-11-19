package metadata

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

//metadata API for client
func CloneClientOutgoingData(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	//return a copy
	return md.Copy()
}

//metadata API for server
func CloneServerIncomingData(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	//return a copy
	return md.Copy()
}

type gRPCMD metadata.MD

func ExtractIncoming(ctx context.Context) gRPCMD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		//create a new one
		return gRPCMD(metadata.Pairs())
	}
	return gRPCMD(md)
}

func ExtractOutgoing(ctx context.Context) gRPCMD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return gRPCMD(metadata.Pairs())
	}
	return gRPCMD(md)
}

func (m gRPCMD) Clone(copiedKeys ...string) gRPCMD {
	newMd := gRPCMD(metadata.Pairs())
	for k, vv := range m {
		found := false
		if len(copiedKeys) == 0 {
			found = true
		} else {
			for _, allowedKey := range copiedKeys {
				if strings.EqualFold(allowedKey, k) {
					found = true
					break
				}
			}
		}
		if !found {
			continue
		}
		newMd[k] = make([]string, len(vv))
		copy(newMd[k], vv)
	}
	return gRPCMD(newMd)
}

func (m gRPCMD) ToOutgoing(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD(m))
}

func (m gRPCMD) ToIncoming(ctx context.Context) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD(m))
}

func (m gRPCMD) Get(key string) string {
	k := strings.ToLower(key)
	vv, ok := m[k]
	if !ok {
		return ""
	}
	return vv[0]
}

func (m gRPCMD) Del(key string) gRPCMD {
	k := strings.ToLower(key)
	delete(m, k)
	return m
}

func (m gRPCMD) Set(key string, value string) gRPCMD {
	k := strings.ToLower(key)
	m[k] = []string{value}
	return m
}

func (m gRPCMD) Add(key string, value string) gRPCMD {
	k := strings.ToLower(key)
	m[k] = append(m[k], value)
	return m
}
