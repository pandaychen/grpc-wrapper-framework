## 简介

atreus 是一个基于 gRPC 封装的脚手架

## 0x01 错误处理（规范）

#### 服务端错误返回

服务端错误返回需要使用如下代码完成，其中第一个参数来源于[官方](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)，第二个参数为自定义，实现[代码](https://github.com/grpc/grpc-go/blob/v1.42.0/status/status.go#L57)：

```golang
return status.Error(codes.Internal, pyerrors.InternalError)
```

#### 客户端错误返回

需要纳入熔断错误计算的类型（超时类、服务器错误等）：

- `codes.Unknown`：异常错误（recover 拦截器）
- `codes.DeadlineExceeded`：服务端 ctx 超时（timeout 拦截器）
- `codes.ResourceExhausted`：服务端限速丢弃（limiter 拦截器）

不纳入的（逻辑错误等）：

- `codes.InvalidArgument`：非法参数（ACL 拦截器）
- `codes.Unauthenticated`：未认证（auth 拦截器）
- `codes.InvalidArgument`：非法参数（validator 拦截器）

## 0x02 参数校验：go.validator 接入 proto 的步骤

#### 使用 github.com/mwitkow/go-proto-validators/protoc-gen-govalidators 包

1. 设置 `GOPATH`，如本机的 `GOPATH` 地址为 `/root/go/`

2. 下载 `https://github.com/protocolbuffers/protobuf` 项目

3. 将 `protobuf/src/*` 目录复制到 `GOPATH` 中的如下路径：

```bash
cp src/ ${GOPATH}/src/github.com/google/protobuf/src -r
cp src/ /root/go/src/github.com/google/protobuf/ -r
```

4. 下载 `protoc-gen-govalidators` 包：

```bash
go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
```

5. 编写 proto 文件，注意添加 `validator.proto` 包及协议字段的 validato 规则

```protobuf
syntax = "proto3";

// protoc -I=. *.proto --go_out=plugins=grpc:.

option java_multiple_files = true;
option java_package = "io.grpc.examples.helloworld";
option java_outer_classname = "HelloWorldProto";

import "github.com/mwitkow/go-proto-validators/validator.proto";

package proto;

//a common RPC names with Serivce suffix
service GreeterService {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    //rpc ErrorSayHello(HelloRequest) returns  (HelloReply) {}
}

message HelloRequest {
    string name = 1 [(validator.field) = {regex: "^[a-z]{2,5}$"}];
}

message HelloReply {
    string message = 1;
}
```

5. 生成 `pb.go` 及 `validator.pb.go` 文件，完成：

```bash
protoc    --proto_path=${GOPATH}/src   --proto_path=${GOPATH}/src/github.com/google/protobuf/src   --proto_path=.   --go_out=.   --govalidators_out=. --go_out=plugins=grpc:./   *.proto
```

#### 使用 github.com/envoyproxy/protoc-gen-validate

1.  安装 `protoc-gen-validate` 工具

```bash
go get -d github.com/envoyproxy/protoc-gen-validate
cd ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate
make build
```

安装成功：

```bash
GOBIN=/root/go/src/github.com/envoyproxy/protoc-gen-validate/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
protoc -I . \
        --plugin=protoc-gen-go=/root/go//bin/protoc-gen-go \
        --go_opt=paths=source_relative \
        --go_out="Mvalidate/validate.proto=github.com/envoyproxy/protoc-gen-validate/validate,Mgoogle/protobuf/any.proto=google.golang.org/protobuf/types/known/anypb,Mgoogle/protobuf/duration.proto=google.golang.org/protobuf/types/known/durationpb,Mgoogle/protobuf/struct.proto=google.golang.org/protobuf/types/known/structpb,Mgoogle/protobuf/timestamp.proto=google.golang.org/protobuf/types/known/timestamppb,Mgoogle/protobuf/wrappers.proto=google.golang.org/protobuf/types/known/wrapperspb,Mgoogle/protobuf/descriptor.proto=google.golang.org/protobuf/types/descriptorpb:." validate/validate.proto
go install .
go: downloading github.com/lyft/protoc-gen-star v0.5.3
go: downloading github.com/iancoleman/strcase v0.2.0
```

2. 编写 proto 文件，如下：

```protobuf
syntax = "proto3";

// protoc -I=. *.proto --go_out=plugins=grpc:.

option java_multiple_files = true;
option java_package = "io.grpc.examples.helloworld";
option java_outer_classname = "HelloWorldProto";


import "validate/validate.proto";

package proto;

//a common RPC names with Serivce suffix
service GreeterService {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
    //rpc ErrorSayHello(HelloRequest) returns  (HelloReply) {}
}

message HelloRequest {
    string name = 1  [(validate.rules).string = {
                      pattern:   "^[a-z]{2,5}$",
                      max_bytes: 256,
                   }];
}

message HelloReply {
    string message = 1;
}
```

3. 编译 proto 文件，生成 `pb.go` 及 `pb.validator.go` ，完成。

```bash
protoc   -I .   -I ${GOPATH}/src   -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate     --validate_out="lang=go:." --go_out=plugins=grpc:./   *.proto
```

4. 规则可见：https://github.com/envoyproxy/protoc-gen-validate#constraint-rules
