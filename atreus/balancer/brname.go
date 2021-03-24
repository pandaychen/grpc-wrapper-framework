package balancer

import (
	"strconv"

	"github.com/pandaychen/grpc-wrapper-framework/common/vars"
	"google.golang.org/grpc/metadata"
)

//注册到 grpc中负载均衡器的全局名字

type GRPC_BALANCER_NAME string

const (
	BALANCER_RandomWeight_NAME   GRPC_BALANCER_NAME = "RandomWeight"
	BALANCER_SimpleWeightRR_NAME GRPC_BALANCER_NAME = "SimpleWeightRR"
)

func GetServerWeightValue(mdata interface{}) int {
	md, ok := mdata.(metadata.MD)
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
