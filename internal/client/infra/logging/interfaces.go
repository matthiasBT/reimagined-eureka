package logging

type ILogger interface {
	Infoln(format string, args ...interface{})
	Info(format string, args ...interface{})
	Successln(format string, args ...interface{})
	Warningln(format string, args ...interface{})
	Debugln(format string, args ...interface{})
	Failureln(format string, args ...interface{})
}
