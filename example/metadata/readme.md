##  metadata 测试

客户端发送数据，可见 gRPC 操作的数据，仅仅是 `key=metadata.mdOutgoingKey` 的部分：
```javascript
context.Background.WithValue(type metadata.mdOutgoingKey, val <not Stringer>)
context.Background.WithValue(type metadata.mdOutgoingKey, val <not Stringer>).WithValue(type string, val value3)
context.Background.WithValue(type metadata.mdOutgoingKey, val <not Stringer>).WithValue(type string, val value3)
context.Background.WithValue(type metadata.mdOutgoingKey, val <not Stringer>).WithValue(type string, val value3).WithValue(type metadata.mdOutgoingKey, val <not Stringer>)
metadata.FromOutgoingContext(ctx)= map[app:[test] atreus-requestid:[cvalue] key2:[value2] key4:[value4] method:[normal] token:[test]]
```