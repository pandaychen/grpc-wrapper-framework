package etcdv3

import (
	"encoding/json"
	"sync"
	"time"

	com "grpc-wrapper-framework/microservice/discovery/common"

	"go.etcd.io/etcd/api/v3/mvccpb"
	etcd3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/resolver"
)

const (
	DEFAULT_CHANNEL_SIZE = 64
)

//独立封装watcher
type EtcdWatcher struct {
	WatchKey  string
	Client    *etcd3.Client
	Ctx       *context.Context
	Cancel    context.CancelFunc
	AddrsList []resolver.Address
	WatchCh   etcd3.WatchChan // watch() RETURN channel
	Logger    *zap.Logger
	Wg        *sync.WaitGroup
}

func (w *EtcdWatcher) Close() {
	w.Cancel()
}

//create a etcd watcher,which belongs to etcd resolver
func NewEtcdWatcher(ctx *context.Context, key string, etcdclient *etcd3.Client, zaploger *zap.Logger, wg *sync.WaitGroup) *EtcdWatcher {
	new_ctx, new_cancel := context.WithCancel(*ctx)
	watcher := &EtcdWatcher{
		WatchKey: key,
		Client:   etcdclient,
		Logger:   zaploger,
		Ctx:      &new_ctx,
		Cancel:   new_cancel,
		Wg:       wg,
	}
	return watcher
}

// sync full addrs
func (w *EtcdWatcher) GetAllAddresses() []resolver.Address {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	total_addrlist := []resolver.Address{}

	//get all prefix keys
	getResp, err := w.Client.Get(ctx, w.WatchKey, etcd3.WithPrefix())

	if err == nil {
		addrs := w.ExtractAddrs(getResp)
		if len(addrs) > 0 {
			for _, saddr := range addrs {
				total_addrlist = append(total_addrlist, resolver.Address{
					Addr:     saddr.AddressInfo, // Addr 和grpc的resolver中的结构体格式保持一致
					Metadata: &saddr.Metadata,   // Metadata is the information associated with Addr, which may be used
				})
			}
		}
	} else {
		w.Logger.Error("Watcher: get all keys withprefix error", zap.String("errmsg", err.Error()))
	}
	return total_addrlist
}

//返回range channel
func (w *EtcdWatcher) start() chan []resolver.Address {
	retchannel := make(chan []resolver.Address, DEFAULT_CHANNEL_SIZE)
	w.Wg.Add(1)
	go func() {
		defer func() {
			close(retchannel)
			w.Wg.Done()
		}()

		//init
		w.AddrsList = w.GetAllAddresses()
		retchannel <- w.cloneAddresses(w.AddrsList)

		//starting a watching channel
		w.WatchCh = w.Client.Watch(*w.Ctx, w.WatchKey, etcd3.WithPrefix(), etcd3.WithPrevKV())
		for wresp := range w.WatchCh {
			//block and go range,watching etcd events change
			for _, ev := range wresp.Events {
				//range  wresp.Events slice
				switch ev.Type {
				case mvccpb.PUT:
					jsonobj := com.ServiceBasicInfo{}
					err := json.Unmarshal([]byte(ev.Kv.Value), &jsonobj)
					if err != nil {
						w.Logger.Error("Parse node data error", zap.String("errmsg", err.Error()))
						continue
					}
					//generate grpc Address struct
					addr := resolver.Address{Addr: jsonobj.AddressInfo, Metadata: &jsonobj.Metadata}
					if w.addAddr(addr) {
						//if-add-new,return new
						retchannel <- w.cloneAddresses(w.AddrsList)
					}
				case mvccpb.DELETE:
					jsonobj := com.ServiceBasicInfo{}
					err := json.Unmarshal([]byte(ev.PrevKv.Value), &jsonobj)
					w.Logger.Info("key", zap.String("prevalue", string(ev.PrevKv.Value)), zap.String("value", string(ev.Kv.Value)))
					w.Logger.Info("value", zap.String("preky", string(ev.PrevKv.Key)), zap.String("Key", string(ev.Kv.Key)))

					w.Logger.Info("value", zap.String("value", string(ev.PrevKv.Value)))
					w.Logger.Info("key", zap.String("key", string(ev.Kv.Key)))
					if err != nil {
						w.Logger.Error("Parse node data error", zap.String("errmsg", err.Error()))
						continue
					}
					addr := resolver.Address{Addr: jsonobj.AddressInfo, Metadata: &jsonobj.Metadata}
					if w.removeAddr(addr) {
						retchannel <- w.cloneAddresses(w.AddrsList)
					}
				}
			}
		}
	}()

	//直接返回一个可以range的带缓冲channel
	return retchannel
}

//get keys from etcdctl response
func (w *EtcdWatcher) ExtractAddrs(etcdresponse *etcd3.GetResponse) []com.ServiceBasicInfo {
	addrs := []com.ServiceBasicInfo{}

	//KVS is slice
	if etcdresponse == nil || etcdresponse.Kvs == nil {
		return addrs
	}

	for i := range etcdresponse.Kvs {
		if v := etcdresponse.Kvs[i].Value; v != nil {
			//parse string to node-data
			jsonobj := com.ServiceBasicInfo{}
			err := json.Unmarshal(v, &jsonobj)
			if err != nil {
				w.Logger.Error("Parse node data error", zap.String("errmsg", err.Error()))
				continue
			}
			addrs = append(addrs, jsonobj)
		}
	}
	return addrs
}

//
func (w *EtcdWatcher) cloneAddresses(in []resolver.Address) []resolver.Address {
	out := make([]resolver.Address, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[i]
	}
	return out
}

// 检查addr是否已存在,如果没存在就增加
func (w *EtcdWatcher) addAddr(addr resolver.Address) bool {
	for _, v := range w.AddrsList {
		if addr.Addr == v.Addr {
			return false
		}
	}
	w.AddrsList = append(w.AddrsList, addr)
	return true
}

// 检查addr是否已存在,如果没存在就删除
func (w *EtcdWatcher) removeAddr(addr resolver.Address) bool {
	for i, v := range w.AddrsList {
		if addr.Addr == v.Addr {
			w.AddrsList = append(w.AddrsList[:i], w.AddrsList[i+1:]...)
			return true
		}
	}
	return false
}
