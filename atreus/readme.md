## 简介

atreus 是一个基于 gRPC 封装的脚手架

## go.validator 接入 proto 验证方法

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
