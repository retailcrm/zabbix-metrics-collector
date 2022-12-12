package metrics

import "log"

var (
	DefaultLogger = &defaultLogger{}
	NilLogger     = &nilLogger{}
)

type defaultLogger struct{}

func (l *defaultLogger) Errorf(fmt string, args ...interface{}) {
	log.Printf("error: "+fmt, args...)
}

type nilLogger struct{}

func (l *nilLogger) Errorf(fmt string, args ...interface{}) {}
