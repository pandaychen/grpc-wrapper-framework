package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	LABLES_NAME_RPCTYPE_UNARY  = "unary"
	LABLES_NAME_RPCTYPE_STREAM = "stream"
)

// 默认的采集器（定义），也可在代码中调用NewCounterVec等进行创建
var (
	ServerHandleCounter = CounterVecOption{
		Namespace: DefaultNamespace,
		Name:      "server_counter_total",
		Labels:    []string{"type", "method", "code"}, //默认的标签
	}.Build()

	ClientHandleCounter = CounterVecOption{
		Namespace: DefaultNamespace,
		Name:      "client_counter_total",
		Labels:    []string{"type", "name", "method", "code"},
	}.Build()
)

func init() {
	//TODO：改为注册路由的方式，默认以http方式启动
	go func() {
		addr := fmt.Sprintf(":%d", DEFAULT_PROMHTTP_PORT)
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(addr, nil)
	}()
}
