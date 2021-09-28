package discovery

import (
	"errors"

	//etcdv3 "github.com/pandaychen/etcd_tools"
	"grpc-wrapper-framework/common/enums"
	com "grpc-wrapper-framework/microservice/discovery/common"
	"grpc-wrapper-framework/microservice/discovery/etcdv3"

	"google.golang.org/grpc/resolver"
)

type ServiceRegisterWrapper interface {
	ServiceRegister() error
	ServiceUnRegister() error
}

type ServiceResolverWrapper interface {
	Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error)
	Scheme() string
	ResolveNow(o resolver.ResolveNowOptions)
	Close()
}

func NewDiscoveryRegister(conf *com.RegisterConfig) (ServiceRegisterWrapper, error) {
	switch conf.RegisterType {
	case enums.REG_TYPE_ETCD:
		return etcdv3.NewRegister(conf)
	default:
		return nil, errors.New("not support register method")
	}
}

//Create grpc resolver
func NewDiscoveryResolver(conf *com.ResolverConfig) (ServiceResolverWrapper, error) {
	switch conf.RegisterType {
	case enums.REG_TYPE_ETCD:
		return etcdv3.NewResolverRegister(conf)
	default:
		return nil, errors.New("not support resolve method")
	}
}
