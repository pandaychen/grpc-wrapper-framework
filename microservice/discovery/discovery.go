package discovery

import (
	"errors"

	//etcdv3 "github.com/pandaychen/etcd_tools"
	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	com "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/common"
	"github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/etcdv3"
	"google.golang.org/grpc/resolver"
)

type ServiceRegisterWrapper interface {
	ServiceRegister() error
	ServiceUnRegister() error
}

type ServiceResolverWrapper interface {
	Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error)
	Scheme() string
	ResolveNow(o resolver.ResolveNowOption)
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
