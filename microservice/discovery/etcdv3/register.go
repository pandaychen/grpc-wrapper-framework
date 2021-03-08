package etcdv3

//etcdv3的服务注册封装

import (
	"encoding/json"
	"fmt"
	"time"

	com "github.com/pandaychen/grpc-wrapper-framework/microservice/discovery/common"
	etcd3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type EtcdRegister struct {
	Etcd3Client *etcd3.Client
	Logger      *zap.Logger
	Key         string //service uniq-key
	Value       string //micro-service ip+port+weight
	Ttl         time.Duration
	Ctx         context.Context
	Cancel      context.CancelFunc
	ApiOn       bool
	Leaseid     clientv3.LeaseID
}

func NewRegister(config *com.RegisterConfig) (*EtcdRegister, error) {
	//TODO：use https://github.com/pandaychen/etcd_tools/blob/master/clientv3.go instead
	etcdConfg := etcd3.Config{
		Endpoints: config.Endpoint,
	}

	client, err := etcd3.New(etcdConfg)
	if err != nil {
		config.Logger.Error("[NewRegister]Create etcdv3 client error", zap.String("errmsg", err.Error()))
		return nil, err
	}

	//check format
	val, err := json.Marshal(config.NodeData)
	if err != nil {
		config.Logger.Error("[NewRegister]Create etcdv3 value error", zap.String("errmsg", err.Error()))
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	registry := &EtcdRegister{
		Etcd3Client: client,
		Logger:      config.Logger,
		Ttl:         config.Ttl / time.Second,
		Ctx:         ctx,
		Cancel:      cancel,
		Key:         config.BuildEtcdKey(),
		Value:       string(val),
		ApiOn:       true,
	}
	return registry, nil
}

func (r *EtcdRegister) ServiceRegister() error {
	if r.ApiOn {
		resp, err := r.Etcd3Client.Grant(r.Ctx, int64(r.Ttl))
		if err != nil {
			r.Logger.Error("Register Grant error", zap.String("errmsg", err.Error()))
			return fmt.Errorf("create etcd3 lease failed: %v", err)
		}
		r.Leaseid = resp.ID
		if _, err := r.Etcd3Client.Put(context.TODO(), r.Key, r.Value, etcd3.WithLease(resp.ID)); err != nil {
			r.Logger.Error("Set key with ttl error", zap.String("key", r.Key), zap.String("leaseid", fmt.Sprintf("%x", resp.ID)), zap.String("errmsg", err.Error()))
			return fmt.Errorf("set service '%s' with ttl to etcd3 failed: %s", r.Key, err.Error())
		}

		//in keepalive,start with a new groutine for loop
		if leaseRespChan, err := r.Etcd3Client.KeepAlive(context.TODO(), resp.ID); err != nil {
			r.Logger.Error("Set key keepalive error", zap.String("key", r.Key), zap.String("leaseid", fmt.Sprintf("%x", resp.ID)), zap.String("errmsg", err.Error()))
			return fmt.Errorf("refresh service '%s' with ttl to etcd3 failed: %s", r.Key, err.Error())
		}
		go ListenLeaseChan(leaseRespChan)
	} else {

	}

	return nil
}

func (r *EtcdRegister) ServiceUnRegister() error {

	return nil
}
func (r *EtcdRegister) Close() {
	return
}

func ListenLeaseChan(leaseRespChan <-chan *etcd3.LeaseKeepAliveResponse) {
	var (
		leaseKeepResp *etcd3.LeaseKeepAliveResponse
	)
	for {
		select {
		case leaseKeepResp = <-leaseRespChan:
			if leaseKeepResp == nil {
				fmt.Println("Etcd leaseid not effectiveness,quit")
				//TODO:PANDAYCHEN,这里需要更多处理
				//一个思路是发送channel通知应用层，重新续期(告警)
				goto END
			} else {
				//fmt.Println("Etcd leaseid effectiveness", leaseKeepResp.ID)
			}
		}
	}
END:
}
