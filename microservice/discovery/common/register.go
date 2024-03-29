package common

import (
	"fmt"
	"strings"
	"time"

	//etcdv3 "github.com/pandaychen/etcd_tools"
	"grpc-wrapper-framework/common/enums"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type ServiceBasicInfo struct {
	AddressInfo string
	Metadata    metadata.MD
}

// 注册到etcd中的key-value信息
type RegisterConfig struct {
	InstanceId     string //must be uniq random
	RegisterType   enums.RegType
	RootName       string //root-name
	ServiceName    string //service-name
	ServiceVersion string //version
	ServiceNodeID  string //node-name (IP:ADDR)
	RandomSuffix   string
	NodeData       ServiceBasicInfo
	Ttl            time.Duration

	Endpoint string
	//ETCD config
	//EtcdConfig *etcdv3.EtcdConfig
	Logger *zap.Logger
}

//
func (c *RegisterConfig) BuildEtcdKey() string {
	if strings.HasPrefix(c.RootName, "/") {
		return fmt.Sprintf("%s/%s/%s/%s%s", c.RootName, c.ServiceName, c.ServiceVersion, c.ServiceNodeID, c.RandomSuffix)
	} else {
		return fmt.Sprintf("/%s/%s/%s/%s%s", c.RootName, c.ServiceName, c.ServiceVersion, c.ServiceNodeID, c.RandomSuffix)
	}
}
