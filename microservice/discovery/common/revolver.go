package common

import "go.uber.org/zap"

type ResolverConfig struct {
	WatchKey string
	Endpoint string
	Logger   *zap.Logger
}
