package balancer

import (
	"strconv"

	"grpc-wrapper-framework/common/vars"

	"google.golang.org/grpc/metadata"
)

//注册到 grpc中负载均衡器的全局名字

type GRPC_BALANCER_NAME string

const (
	BALANCER_DEFAULT_RR_NAME                        = "round_robin" //默认grpc实现
	BALANCER_RandomWeight_NAME   GRPC_BALANCER_NAME = "RandomWeight"
	BALANCER_SimpleWeightRR_NAME GRPC_BALANCER_NAME = "SimpleWeightRR"
	BALANCER_LeastConn_NAME      GRPC_BALANCER_NAME = "LeastConn"
)

func GetServerWeightValue(mdata interface{}) int {
	//md, ok := mdata.(metadata.MD)
	md, ok := mdata.(*metadata.MD)
	if ok {
		values := md.Get(vars.SERVICE_WEIGHT_KEY)
		if len(values) > 0 {
			weight, err := strconv.Atoi(values[0])
			if err == nil {
				return weight
			}
		}
	}
	return vars.DEFAULT_WEIGHT
}

//新版本Address结构：https://pkg.go.dev/google.golang.org/grpc/resolver#Address
