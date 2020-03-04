package log

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger: logger,
	}
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

func (l *ZapLogger) Infoln(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}
func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

func (l *ZapLogger) Warning(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

func (l *ZapLogger) Warningln(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

func (l *ZapLogger) Warningf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *ZapLogger) Errorln(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
}

func (l *ZapLogger) Fatalln(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

func (l *ZapLogger) V(v int) bool {
	return false
}
