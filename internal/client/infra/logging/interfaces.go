package logging

type ILogger interface {
	Info(format string, args ...interface{})
	Success(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Error(format string, args ...interface{})
}
