package common

import (
	"fmt"
	"strings"

	"github.com/pandaychen/grpc-wrapper-framework/common/enums"
	"go.uber.org/zap"
)

type ResolverConfig struct {
	RegisterType   enums.RegType
	Schemename     string
	RootName       string //root-name
	ServiceName    string //service-name
	ServiceVersion string //version
	ServiceNodeID  string //node-name (IP:ADDR)

	Endpoint string
	Logger   *zap.Logger
}

func (c *ResolverConfig) BuildEtcdPrefix() string {
	if strings.HasPrefix(c.RootName, "/") {
		return fmt.Sprintf("%s/%s/%s/", c.RootName, c.ServiceName, c.ServiceVersion)
	} else {
		return fmt.Sprintf("/%s/%s/%s/", c.RootName, c.ServiceName, c.ServiceVersion)
	}
}
