package xorm

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var (
	lock       sync.RWMutex
	XormSqlMap = make(map[string]*XormClient)
)

type XormClient struct {
	name      string
	opts      atomic.Value
	incr      int64
	currCount int64
	lastUsed  time.Time
	*xorm.Engine
}

type XormOption struct {
	Name   string
	Driver string
	Dsn    string
}

func NewXormClient(option *XormOption) (*XormClient, error) {
	var (
		err     error
		xormCli XormClient
	)

	xormCli.Engine, err = xorm.NewEngine(option.Driver, option.Dsn)
	if err != nil {
		return nil, err
	}

	xormCli.SetDefaultContext(context.WithValue(context.Background(), clientInstance, &xormCli))

	xormCli.opts.Store(option)

	// 注入钩子实现
	xormCli.AddHook(NewXormHook(option.Name))

	lock.Lock()
	XormSqlMap[option.Name] = &xormCli
	lock.Unlock()
	return &xormCli, nil
}

func GetXormClient(name string) *XormClient {
	lock.RLock()
	defer lock.RUnlock()
	if _, exists := XormSqlMap[name]; exists {
		return XormSqlMap[name]
	}

	return nil
}
