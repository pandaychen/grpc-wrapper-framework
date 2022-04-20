## 组件
Tracing

## 概览

* 基于 opentracing 语义
* 使用 protobuf 协议描述 trace 结构
* 参考 Kratos 代码改造，对比[opentracing-go](https://github.com/opentracing/opentracing-go)的实现更为轻量


##  结构定义
根据opentracing的语义：<br>
![base](https://raw.githubusercontent.com/pandaychen/pandaychen.github.io/master/blog_img/microservice/tracing_base_structure.png)


####    Tag
kv结构

####    Log
kv结构



##  封装
-   `Inject`：注入的过程就是把 context 的所有信息写入到一个叫 Carrier 的 map 中，然后把 map 中的所有 KV 对写入 HTTP Header 或者 grpc Metadata
-   `Extract`：抽取过程是注入的逆过程，从 carrier，也就是 HTTP Headers（grpc Metadata），构建 SpanContext

整个过程类似客户端和服务器传递数据的序列化和反序列化的过程。这里的 Carrier （Map）支持 Key 为 string 类型，value 为 string 或者 Binary 格式（Bytes）


##  参考
-   [Tracers](https://opentracing.io/docs/overview/tracers/)
-   [grpc-tracing例子](https://github.com/grpc-ecosystem/grpc-opentracing/tree/master/go/otgrpc)