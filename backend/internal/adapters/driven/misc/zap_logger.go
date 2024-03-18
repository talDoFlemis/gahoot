package misc

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(logger *zap.SugaredLogger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (z *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		z.logger.Infow(msg)
	} else {
		z.logger.Infow(msg, keysAndValues...)
	}
}

func (z *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		z.logger.Errorf(msg)
	} else {
		z.logger.Errorf(msg, keysAndValues...)
	}
}

func (z *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		z.logger.Warnf(msg)
	} else {
		z.logger.Warnf(msg, keysAndValues...)
	}
}

func (z *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	if keysAndValues == nil {
		z.logger.Debugf(msg)
	} else {
		z.logger.Debugf(msg, keysAndValues...)
	}
}

func (z *ZapLogger) Warnf(msg string, args ...interface{}) {
	z.logger.Warnf(msg, args...)
}

func (z *ZapLogger) Errorf(msg string, args ...interface{}) {
	z.logger.Errorf(msg, args...)
}

func (z *ZapLogger) Infof(msg string, args ...interface{}) {
	z.logger.Infof(msg, args...)
}

func (z *ZapLogger) Debugf(msg string, args ...interface{}) {
	z.logger.Debugf(msg, args...)
}
