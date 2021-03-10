package discovery

import (
	"errors"

	//etcdv3 "github.com/pandaychen/etcd_tools"
	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	com "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/common"
	"github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/etcdv3"
)

type ServiceRegisterWrapper interface {
	ServiceRegister() error
	ServiceUnRegister() error
}

func NewDiscoveryRegister(conf *com.RegisterConfig) (ServiceRegisterWrapper, error) {
	switch conf.RegisterType {
	case enums.REG_TYPE_ETCD:
		return etcdv3.NewRegister(conf)
	default:
		return nil, errors.New("not support register method")
	}
}

func NewDiscoveryResolver() interface{} {
	switch conf.RegisterType {
	case enums.REG_TYPE_ETCD:
		return etcdv3.NewRegister(conf)
	default:
		return nil, errors.New("not support register method")
	}
}
