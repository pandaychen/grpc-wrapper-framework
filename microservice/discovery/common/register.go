package common

import (
	"time"

	etcdv3 "github.com/pandaychen/etcd_tools"
	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	"github.com/pandaychen/grpc-wrapper-framework/discovery/etcdv3"
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
	ServiceNodeID  string //node-name
	RandomSuffix   string
	Info           ServiceBasicInfo
	Ttl            time.Duration
	Logger         *zap.Logger

	//ETCD config
	EtcdConfig *etcdv3.EtcdConfig
}