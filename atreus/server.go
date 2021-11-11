package atreus

//A wrapper grpc-server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"grpc-wrapper-framework/common/enums"
	"grpc-wrapper-framework/common/vars"
	"grpc-wrapper-framework/config"
	"grpc-wrapper-framework/logger"
	auth "grpc-wrapper-framework/microservice/authentication"
	dis "grpc-wrapper-framework/microservice/discovery"
	discom "grpc-wrapper-framework/microservice/discovery/common"
	"grpc-wrapper-framework/pkg/xrand"
)

const (
	DEFAULT_ATREUS_SERVICE_NAME = "atreus_svc"
	DEFAULT_TIME_TO_QUIT        = 5 * time.Second
)

//grpc-server核心结构（封装）
type Server struct {
	Logger   *zap.Logger
	Conf     *config.AtreusSvcConfig
	ConfLock *sync.RWMutex

	EtcdClient          *etcdv3.Client
	InnerHandlers       []grpc.UnaryServerInterceptor  //拦截器数组
	InnerStreamHandlers []grpc.StreamServerInterceptor //stream拦截器数组
	ServiceReg          dis.ServiceRegisterWrapper
	//auth
	Auther *auth.Authenticator //通用的验证接口
	//limiter
	Limiters *XRateLimiter
	//context
	Ctx     context.Context
	IsDebug bool

	//wrapper Server
	RpcServer *grpc.Server //原生Server
}

func NewServer(conf *config.AtreusSvcConfig, opt ...grpc.ServerOption) *Server {
	var err error
	if conf == nil {
		panic("atreus server config null")
	}
	/*
		var opt []grpc.ServerOption
		opt = append(opt, grpc.UnaryInterceptor(AtreusUnaryInterceptorChain(Recovery, Middle, Timing, Middle)))
		//return grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorChain(Recovery, Logging)))
	*/

	logconf := logger.LogConfig{
		ServiceName: DEFAULT_ATREUS_SERVICE_NAME,
	}

	logger, err := logconf.CreateNewLogger(conf.LogConf)
	if err != nil {
		panic(err)
	}
	srv := &Server{
		Logger:              logger,
		ConfLock:            new(sync.RWMutex),
		InnerHandlers:       make([]grpc.UnaryServerInterceptor, 0),
		InnerStreamHandlers: make([]grpc.StreamServerInterceptor, 0),
		Conf:                conf,
		Ctx:                 context.Background(),
	}

	if conf.SrvConf.Keepalive {
		//初始化gRPC-Server的keepalive参数
		keepaliveopts := grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(conf.SrvConf.IdleTimeout), //如果一个client空闲超过MaxConnectionIdle-s,发送一个GOAWAY,为了防止同一时间发送大量GOAWAY
			//假设MaxConnectionIdle=15s，那么会在15s时间间隔上下浮动MaxConnectionIdle*10%,即15+1.5或者15-1.5
			MaxConnectionAgeGrace: time.Duration(conf.SrvConf.ForceCloseWait),    //在强制关闭连接之间,允许有ForceCloseWait-s的时间完成pending的rpc请求
			Time:                  time.Duration(conf.SrvConf.KeepAliveInterval), //如果一个clinet空闲超过KeepAliveInterval-s,则发送一个ping请求
			Timeout:               time.Duration(conf.SrvConf.KeepAliveTimeout),  //如果ping请求KeepAliveTimeout-s内未收到回复,则认为该连接已断开
			MaxConnectionAge:      time.Duration(conf.SrvConf.MaxLifeTime),       //如果任意连接存活时间超过MaxLifeTime-s,发送一个GOAWAY
		})

		opt = append(opt, keepaliveopts, grpc.UnaryInterceptor(srv.BuildUnaryInterceptorChain2))
	} else {
		opt = append(opt, grpc.UnaryInterceptor(srv.BuildUnaryInterceptorChain2))
	}

	srv.RpcServer = grpc.NewServer(opt...)

	//Fill the interceptors

	//注意：Metrics2Prometheus必须放在Limiters的前面，否则，捕获不到Limiters返回的错误
	srv.Use(srv.Recovery(), srv.Timing(), srv.XRequestId(), srv.Metrics2Prometheus())

	if conf.LimiterConf.On {
		srv.Limiters = NewXRateLimiter(rate.Limit(conf.LimiterConf.LimiterRate), conf.LimiterConf.LimiterSize)
		srv.Use(srv.Limit(srv.Limiters))
	}

	//开启auth
	if conf.AuthConf.On {
		srv.Auther, err = auth.NewAuthenticator(&srv.Ctx)
		if err != nil {
			panic(err)
		}
		srv.Use(srv.Authorize())
	}

	srv.Use(srv.SrvValidator())

	if conf.RegistryConf.RegOn {
		nodeinfo := discom.ServiceBasicInfo{
			AddressInfo: conf.SrvConf.Addr,
			Metadata:    metadata.Pairs(vars.SERVICE_WEIGHT_KEY, conf.WeightConf.Weight),
		}

		srv.ServiceReg, err = dis.NewDiscoveryRegister(&discom.RegisterConfig{
			RegisterType:   enums.RegType(conf.RegistryConf.RegisterType),
			RootName:       conf.RegistryConf.RegisterRootPath,
			ServiceName:    conf.RegistryConf.RegisterService,
			ServiceVersion: conf.RegistryConf.RegisterServiceVer,
			ServiceNodeID:  fmt.Sprintf("addr-%s", conf.SrvConf.Addr),
			RandomSuffix:   string(xrand.RandomString(8)),
			Ttl:            conf.RegistryConf.RegisterTTL,
			Endpoint:       conf.RegistryConf.RegisterEndpoints,
			Logger:         logger,
			NodeData:       nodeinfo,
		})

		if err == nil {
			err = srv.ServiceReg.ServiceRegister()
			if err != nil {
				logger.Error("[NewServer]ServiceRegister error", zap.String("errmsg", err.Error()))
				return nil
			}
		} else {
			logger.Error("[NewServer]NewDiscoveryRegister error", zap.String("errmsg", err.Error()))
			panic(err)
		}
	}

	return srv
}

// Server return the grpc server for registering service.
func (s *Server) GetServer() *grpc.Server {
	return s.RpcServer
}

func (s *Server) Serve(lis net.Listener) error {
	return s.RpcServer.Serve(lis)
}

func (s *Server) ReloadConfig() error {
	//not necessary
	s.ConfLock.Lock()
	s.Conf = config.GetAtreusSvcConfig()
	s.ConfLock.Unlock()
	return nil
}

// 启动服务
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.Conf.SrvConf.Addr)
	if err != nil {
		s.Logger.Error("[Run]failed to listen", zap.Any("errmsg", err))
		return err
	}
	reflection.Register(s.RpcServer)
	return s.Serve(listener)
}

// 优雅退出
func (s *Server) Shutdown(ctx context.Context) error {
	var (
		err error
		ch  = make(chan struct{})
	)
	go func() {
		// 调用grpc的GracefulStop()
		s.RpcServer.GracefulStop()
		close(ch)
	}()
	select {
	//force to stop
	case <-ctx.Done():
		s.RpcServer.Stop()
		err = ctx.Err()
	case <-ch:
		return nil
	}
	return err
}

func (s *Server) ExitWithSignalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Logger.Info("Recv Signal to Quit", zap.String("signal", sig.String()))
			ctx, cancel := context.WithTimeout(s.Ctx, DEFAULT_TIME_TO_QUIT)
			defer cancel()
			//gracefully shutdown with timeout
			s.Shutdown(ctx)
			return
		default:
			return
		}
	}
}
