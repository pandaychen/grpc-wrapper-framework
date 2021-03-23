package balancer

//注册到 grpc中负载均衡器的全局名字

type GRPC_BALANCER_NAME string

const (
	BALANCER_RandomWeight_NAME GRPC_BALANCER_NAME = "RandomWeight"
)
