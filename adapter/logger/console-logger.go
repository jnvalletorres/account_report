package logger

import (
	"account_report/domain/port"
	"account_report/env"
	"fmt"
	"time"
)

type consoleLogger struct {
	localEnv env.App
}

func NewConsoleLogger(environment env.Env) port.LoggerPort {
	return &consoleLogger{
		localEnv: environment.App,
	}
}

func (l *consoleLogger) Debug(message string, args ...interface{}) {
	l.log("DEBUG", message, args...)
}

func (l *consoleLogger) Info(message string, args ...interface{}) {
	l.log("INFO ", message, args...)
}

func (l *consoleLogger) Error(message string, args ...interface{}) {
	l.log("ERROR", message, args...)
}

func (l *consoleLogger) log(level string, message string, args ...interface{}) {
	base := fmt.Sprintf("%s %s [%s]: %s ", time.Now().UTC().Format(time.RFC3339), level, l.localEnv.AppId, message)
	fmt.Println(fmt.Sprintf(base, args...))
}
