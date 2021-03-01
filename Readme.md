## grpc-wrapper-framework

> 一个 2021 全年的计划

在业余时间，封装一个 gRPC 封装的微服务框架，希望支持如下特性：

- 支持 Kubernetes 部署
  - 考虑容器 GOMAXPROCS 的影响
  - 支持（默认）DNS 为 Kubernetes 服务发现
  - 支持（API）方式获取 Service 下对应的 Pods 列表
  - 支持 gRPC 的健康检查协议（考虑 CPU、内存及多种因素的 healthy-check）
- 支持 Etcd/Consul 为服务发现的方式
- 支持多种负载均衡算法
  - Nginx 的 WRR 算法
  - Nginx 的 P2C 算法
  - 一致性 hash 算法
  - 普通的 WRR 算法
- 实现多种实用的拦截器实现接口
  - 拦截器链（chain）实现
    - 客户端拦截器
    - 服务端拦截器
  - recovery panic
  - global request id
    - 贯穿一个 RPC 生命周期的 requestid
  - 通用的 zap-Logger
    - 按照请求记录日志，关联到 ctx
    - zap 包 + ctx 关联的日志输出
  - 限流算法
    - 令牌桶
    - 漏桶
  - 熔断器
- 支持 jaeger/zipkin 链路追踪
- 支持动态配置更新（远程和本地）
- 支持内置健康检查服务
- gRPC 封装
  - 客户端的封装，包括拦截器 chain、tls、keepalive、服务发现、熔断器、负载均衡算法等
  - 服务端的封装，包括拦截器 chain、tls、keepalive、链路追踪、服务注册等

grpc 开启的参数：

- keepalives 配置（双向，长连接）
- backoff
