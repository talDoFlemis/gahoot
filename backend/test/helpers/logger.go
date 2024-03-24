package testshelpers

import (
	"io"
	"log"
)

type DummyLogger struct {
	w io.Writer
}

func NewDummyLogger(w io.Writer) *DummyLogger {
	log.SetOutput(w)
	return &DummyLogger{
		w: w,
	}
}

func (d *DummyLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf(msg, keysAndValues...)
}
func (DummyLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf(msg, keysAndValues...)
}
func (DummyLogger) Warn(msg string, keysAndValues ...interface{}) {
	log.Printf(msg, keysAndValues...)
}
func (DummyLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.Printf(msg, keysAndValues...)
}
func (DummyLogger) Warnf(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}
func (DummyLogger) Errorf(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}
func (DummyLogger) Infof(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}
func (DummyLogger) Debugf(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}
