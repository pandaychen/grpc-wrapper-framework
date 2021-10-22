package logger

//使用lumberjack.v2对zap日志进行压缩

import (
	"errors"
	"fmt"
	"grpc-wrapper-framework/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	ServiceName string //服务名
	LogPath     string
}

func InitLumberjack(config *config.LogConfig) *lumberjack.Logger {
	if config.FileName == "" {
		config.FileName = "../log/atreus_service.log"
	}
	hook := lumberjack.Logger{
		Filename:   config.FileName,   // 日志文件路径
		MaxSize:    config.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     config.MaxAge,     // 文件最多保存多少天
		Compress:   config.Compress,   // 是否压缩
	}
	return &hook
}

func (c *LogConfig) CreateNewLogger(config *config.LogConfig) (*zap.Logger, error) {
	lumberhooker := InitLumberjack(config)
	if lumberhooker == nil {
		fmt.Println("[ERROR]CreateNewLogger Lumberjack Init Error")
		return nil, errors.New("create lumberhooker error")
	}

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

	//zlogger := NewZapLogger(logger)

	return logger, nil
}
