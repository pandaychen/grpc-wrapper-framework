package rdb

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	lock sync.RWMutex
	// 声明一个全局的rdb Map
	redisClientMap = make(map[string]*redis.Client)
)

func GetRdbClient(name string) *redis.Client {
	lock.RLock()
	defer lock.RUnlock()
	if _, exists := redisClientMap[name]; exists {
		return redisClientMap[name]
	}

	return nil
}

type RdbOption struct {
	Name         string
	Password     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	IdleTimeout  time.Duration
	Addr         string
	Db           int
}

func NewClient(opts *RdbOption) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		DB:           opts.Db,
		Addr:         opts.Addr,
		Password:     opts.Password,
		PoolSize:     opts.PoolSize,
		ReadTimeout:  opts.ReadTimeout * time.Second,
		WriteTimeout: opts.WriteTimeout * time.Second,
		MinIdleConns: opts.MinIdleConns,
		IdleTimeout:  opts.IdleTimeout * time.Second,
	})
	//嵌入tracing的钩子
	redisClient.AddHook(rdbTracingHook{})

	_, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	lock.Lock()
	redisClientMap[opts.Name] = redisClient
	lock.Unlock()

	return redisClient, nil
}

func combineCommand(cmds ...redis.Cmder) string {
	var commands []string
	for _, cmd := range cmds {
		commands = append(commands, fmt.Sprintf("%v", cmd.String()))
	}
	return strings.Join(commands, ", ")
}
