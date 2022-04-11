// Code generated by protoc-gen-go.
// source: service.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	service.proto
	service_stream.proto

It has these top-level messages:
	HelloRequest
	HelloReply
	SimpleRequest
	SimpleResponse
	StreamRequest
	StreamResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/envoyproxy/protoc-gen-validate/validate"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type HelloRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *HelloRequest) Reset()                    { *m = HelloRequest{} }
func (m *HelloRequest) String() string            { return proto1.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()               {}
func (*HelloRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type HelloReply struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *HelloReply) Reset()                    { *m = HelloReply{} }
func (m *HelloReply) String() string            { return proto1.CompactTextString(m) }
func (*HelloReply) ProtoMessage()               {}
func (*HelloReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto1.RegisterType((*HelloRequest)(nil), "proto.HelloRequest")
	proto1.RegisterType((*HelloReply)(nil), "proto.HelloReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GreeterService service

type GreeterServiceClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterServiceClient struct {
	cc *grpc.ClientConn
}

func NewGreeterServiceClient(cc *grpc.ClientConn) GreeterServiceClient {
	return &greeterServiceClient{cc}
}

func (c *greeterServiceClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := grpc.Invoke(ctx, "/proto.GreeterService/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GreeterService service

type GreeterServiceServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServiceServer(s *grpc.Server, srv GreeterServiceServer) {
	s.RegisterService(&_GreeterService_serviceDesc, srv)
}

func _GreeterService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GreeterService/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServiceServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GreeterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GreeterService",
	HandlerType: (*GreeterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _GreeterService_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

func init() { proto1.RegisterFile("service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 209 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xe2, 0x65,
	0x89, 0x39, 0x99, 0x29, 0x89, 0x25, 0xa9, 0xfa, 0x30, 0x06, 0x44, 0x5e, 0x49, 0x9b, 0x8b, 0xc7,
	0x23, 0x35, 0x27, 0x27, 0x3f, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x48, 0x9a, 0x8b, 0x25,
	0x2f, 0x31, 0x37, 0x55, 0x82, 0x51, 0x81, 0x51, 0x83, 0xd3, 0x89, 0xfd, 0x97, 0x13, 0x4b, 0x11,
	0x93, 0x80, 0x48, 0x10, 0x58, 0x50, 0x49, 0x8d, 0x8b, 0x0b, 0xaa, 0xb8, 0x20, 0xa7, 0x52, 0x48,
	0x82, 0x8b, 0x3d, 0x37, 0xb5, 0xb8, 0x38, 0x31, 0x1d, 0xaa, 0x3a, 0x08, 0xc6, 0x35, 0x72, 0xe3,
	0xe2, 0x73, 0x2f, 0x4a, 0x4d, 0x2d, 0x49, 0x2d, 0x0a, 0x86, 0x38, 0x46, 0xc8, 0x84, 0x8b, 0x23,
	0x38, 0xb1, 0x12, 0xac, 0x59, 0x48, 0x18, 0x62, 0xb5, 0x1e, 0xb2, 0xbd, 0x52, 0x82, 0xa8, 0x82,
	0x05, 0x39, 0x95, 0x4a, 0x0c, 0x4e, 0x06, 0x5c, 0xd2, 0x99, 0xf9, 0x7a, 0xe9, 0x45, 0x05, 0xc9,
	0x7a, 0xa9, 0x15, 0x89, 0xb9, 0x05, 0x39, 0xa9, 0xc5, 0x7a, 0x19, 0x20, 0x05, 0xe5, 0xf9, 0x45,
	0x39, 0x29, 0x4e, 0xfc, 0x60, 0xc5, 0xe1, 0x20, 0x76, 0x00, 0x48, 0x73, 0x00, 0x63, 0x12, 0x1b,
	0xd8, 0x14, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x05, 0xa5, 0x1e, 0x3b, 0x06, 0x01, 0x00,
	0x00,
}
