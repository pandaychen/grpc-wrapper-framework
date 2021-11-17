package atreus

//封装grpc的metadata，使之接口更易于用户调用，主要关注这几种场景
//1. 服务端：接收时，从header中获取metadata；（向header中注入metadata）
//2. tracing中的metadata传递
//3. 客户端：发送时，向header注入metadata；（从header中获取metadata）
//为了便于操作，规定所有的key均为小写

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

//type MD map[string][]string
//MD其实就是一个map结构，必须使用md := MD{}初始化

type WMetadata struct {
	metadata.MD
	//context.Context
}

//从服务端ctx获取metadata
func (m *WMetadata) FromIncoming(ctx context.Context) bool {
	if m == nil {
		m = &WMetadata{metadata.MD{}}
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}
	//m.MD["aaa"] = nil
	m.MD = md
	return true
}

//从客户端ctx获取metadata
func (m *WMetadata) FromOutgoing(ctx context.Context) bool {
	if m == nil {
		m = &WMetadata{metadata.MD{}}
		/*
			m = &WMetadata{
				MD: metadata.Pairs(),
			}
		*/
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return false
	}
	m.MD = md
	return true
}

//客户端注入数据
func (m *WMetadata) ToOutgoing(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD(m.MD))
}

//服务端注入数据
func (m *WMetadata) ToIncoming(ctx context.Context) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD(m.MD))
}

func (m *WMetadata) GetArr(key string) []string {
	lk := strings.ToLower(key)
	vv, ok := m.MD[lk]
	if !ok {
		return nil
	}
	return vv
}

func (m *WMetadata) Get(key string) string {
	lk := strings.ToLower(key)
	vv, ok := m.MD[lk]
	if !ok {
		return ""
	}
	return vv[0]
}

func (m *WMetadata) Del(key string) {
	lk := strings.ToLower(key)
	delete(m.MD, lk)
}

func (m *WMetadata) Set(key string, value string) {
	lk := strings.ToLower(key)
	m.MD[lk] = []string{value}
	return
}

func (m *WMetadata) Add(key string, value string) {
	lk := strings.ToLower(key)
	m.MD[lk] = append(m.MD[lk], value)
	return
}

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
