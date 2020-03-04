package config

//使用lumberjack.v2对zap日志进行压缩

import (
	"fmt"
	"os"

	zlog "github.com/pandaychen/grpc-wrapper-framework/atreus/zlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	ServiceName string //服务名
}

func InitLumberjack() *lumberjack.Logger {
	hook := lumberjack.Logger{
		Filename:   "../log/default.log", // 日志文件路径
		MaxSize:    512,                  // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                   // 日志文件最多保存多少个备份
		MaxAge:     7,                    // 文件最多保存多少天
		Compress:   true,                 // 是否压缩
	}
	return &hook
}

func (c *LogConfig) CreateNewLogger(svcname string) *zlog.ZapLogger {
	lumberhooker := InitLumberjack()
	if lumberhooker == nil {
		fmt.Println("[ERROR]Lumberjack Init Error")
		return nil
	}

	c.ServiceName = svcname
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberhooker)),
		atomicLevel,
	)

	caller := zap.AddCaller()
	development := zap.Development()
	//add default fileds
	filed := zap.Fields(zap.String("ServiceName", string(c.ServiceName)))
	logger := zap.New(core, caller, development, filed)

	zlogger := zlog.NewZapLogger(logger)

	return zlogger
}
