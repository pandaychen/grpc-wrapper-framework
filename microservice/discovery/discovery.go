package discovery

import (
	"errors"

	etcdv3 "github.com/pandaychen/etcd_tools"
	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	_ "github.com/pandaychen/grpc-wrapper-framework/discovery/common"
	"github.com/pandaychen/grpc-wrapper-framework/discovery/etcdv3"
)

type ServiceRegisterWrapper interface {
	ServiceRegister() error
	ServiceUnRegister() error
	Close()
}

func NewDiscoveryRegister(conf *RegisterConfig) (*ServiceRegisterWrapper, error) {
	switch conf.RegisterType {
	case enums.REG_TYPE_ETCD:
		return etcdv3.NewRegister(conf)
	default:
		return nil, errors.New("not support register method")
	}
}
