package atreus

//grpc 服务的 graceful 重启实现

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	xsys "grpc-wrapper-framework/pkg/system"

	"go.uber.org/zap"
)

const (
	defaultListenerFilename = "ATERUS_LISTENER"
)

var (
	ppid = os.Getppid()
)

func checkInheritSign() bool {
	if os.Getenv(defaultListenerFilename) != "" {
		//TODO: check env valid
		return true
	} else {
		return false
	}
}

type Graceful struct {
	Logger *zap.Logger
}

// 继承或创建 listener
func (g *Graceful) RenewListener(new_bindaddr string) (net.Listener, error) {
	if checkInheritSign() {
		return g.inheritListener()
	}

	return g.createListener(new_bindaddr)
}

func (g *Graceful) createListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		g.Logger.Error("createListener error", zap.Any("errmsg", err))
		return nil, err
	}

	return ln, nil
}

func (g *Graceful) inheritListener() (net.Listener, error) {
	// get formmer listener from current process env
	lnFilename := os.Getenv(defaultListenerFilename)
	if lnFilename == "" {
		return nil, errors.New("get LISTENER environment variable error")
	}

	// 根据当前进程 env 中的文件名和描述符，创建一个新的文件
	newListenFd := os.NewFile(uintptr(3), lnFilename)
	if newListenFd == nil {
		return nil, fmt.Errorf("create listener file error: %s", lnFilename)
	}
	defer newListenFd.Close()

	// 创建新的 listener
	newlistener, err := net.FileListener(newListenFd)
	if err != nil {
		g.Logger.Error("inheritListener error", zap.Any("errmsg", err))
		return nil, err
	}

	return newlistener, nil
}

type GracefulGrpcAppserver struct {
	AtreusServer *Server
	Addr         string
	Listener     net.Listener // 当前 GracefulGrpcAppserver 对应的 listener
	Logger       *zap.Logger
}

// 以 GracefulGrpcAppserver 启动并创建 Listener
func NewGracefulGrpcAppserver(srv *Server, new_bindaddr string) (*GracefulGrpcAppserver, error) {
	gsrv := &GracefulGrpcAppserver{
		AtreusServer: srv,
		Addr:         new_bindaddr,
		Logger:       srv.Logger,
	}

	g := &Graceful{}

	// create or inherit listener
	new_listener, err := g.RenewListener(new_bindaddr)
	if err != nil {
		gsrv.Logger.Error("NewGracefulGrpcAppserver RenewListener error", zap.Any("errmsg", err))
		return nil, err
	}

	//update listener
	gsrv.Listener = new_listener

	return gsrv, nil
}

// 使用 GracefulGrpcAppserver 的 run 方法启动服务，代替 Server 启动
func (g *GracefulGrpcAppserver) RunServer() error {
	go func() {
		err := g.AtreusServer.Serve(g.Listener)
		if err != nil {
			gsrv.Logger.Error("[NewGracefulGrpcAppserver]gRPC server start failed", zap.Any("errmsg", err))
			panic(err)
		}
	}()

	ppid := os.Getppid()
	// Close the parent if we inherited and it wasn't init that started us.
	if checkInheritSign() && ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("[NewGracefulGrpcAppserver]failed to close parent: %v", err)
		}
	}

	//monitor and restart process
	g.waitForSignals()
	return nil
}

//
func (g *GracefulGrpcAppserver) waitForSignals() {
	signalCh := make(chan os.Signal, 1024)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		sig := <-signalCh
		g.Logger.Info("[GracefulGrpcAppserver]waitForSignals recv signal", zap.String("signal", sig.String()))
		switch sig {
		case syscall.SIGUSR2:
			child, err := g.forkChild()
			if err != nil {
				g.Logger.Error("[GracefulGrpcAppserver]forkChild error", zap.Any("errmsg", err))
				continue
			}
			g.Logger.Info("[GracefulGrpcAppserver]Forked child succ", zap.Int("newprocess", child.Pid))
		case syscall.SIGINT, syscall.SIGTERM:
			//quit
			g.Logger.Info("[GracefulGrpcAppserver]Receive quit signal..", zap.Any("pid", os.Getpid()))
			signal.Stop(signalCh)
			g.AtreusServer.Shutdown()
			return
		}
	}
}

// forkChild 子进程，替换被 kill 掉的父进程
func (g *GracefulGrpcAppserver) forkChild() (*os.Process, error) {
	// 获取当前进程的 listener 的文件描述符
	lnFile, err := xsys.ExtractListenerFile(g.Listener)
	if err != nil {
		g.Logger.Error("[GracefulGrpcAppserver] ExtractListenerFile error", zap.Any("errmsg", err))
		return nil, err
	}
	defer lnFile.Close()

	// 当前进程的 listener 的文件描述符名字通过环境变量传递给子进程
	environment := append(os.Environ(), fmt.Sprintf("%s=%s", defaultListenerFilename, lnFile.Name()))

	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		g.Logger.Error("[GracefulGrpcAppserver] LookPath error", zap.Any("errmsg", err))
		return nil, err
	}

	// 将标准输入、标准输出、标准错误输出、当前进程的 listener 的文件描述符
	// 4 个 fd 传递给子进程
	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
		lnFile}

	workDir, err := os.Getwd()
	if err != nil {
		g.Logger.Error("[GracefulGrpcAppserver] Getwd error", zap.Any("errmsg", err))
		return nil, err
	}

	// 通过 StartProcess 方式启动子进程
	child_process, err := os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   workDir,
		Env:   environment,
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})

	return child_process, err
}
