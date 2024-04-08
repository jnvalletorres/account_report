package port

type LoggerPort interface {
	Debug(message string, args ...interface{})

	Info(message string, args ...interface{})

	Error(message string, args ...interface{})
}
