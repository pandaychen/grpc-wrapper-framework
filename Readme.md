## grpc-wrapper-framework

> 一个基于 gRPC 封装的脚手架

已支持如下特性：

- gRPC 库封装
  - gRPC 服务端：
    - 拦截器链（interceptor chain）


##  测试
客户端测试需要加入参数，否则会报错 `transport: authentication handshake failed: x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0`
```bash
GODEBUG=x509ignoreCN=0 ./client
```