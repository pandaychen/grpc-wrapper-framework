package atreus

import (
	"google.golang.org/grpc"
)

type Client struct {
	Conf *AtreusClientConfig //指向 客户端配置

	DialOpts []grpc.DialOption             //grpc-客户端option
	handlers []grpc.UnaryClientInterceptor //GRPC拦截器数组
}

func (c *Client) AddCliOpt(opts ...grpc.DialOption) *Client {
	c.DialOpts = append(c.DialOpts, opts...)
	return c
}
