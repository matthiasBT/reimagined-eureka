package logging

import (
	"fmt"

	"github.com/fatih/color"
)

func SetupLogger() ILogger {
	return setupColorLogger()
}

func setupColorLogger() *ColorLogger {
	return &ColorLogger{}
}

type ColorLogger struct {
}

func (l *ColorLogger) print(printer func(string, ...interface{}) string, format string, args ...interface{}) {
	if args == nil {
		fmt.Println(printer(format))
	} else {
		fmt.Println(printer(format, args))
	}
}

func (l *ColorLogger) Info(format string, args ...interface{}) {
	l.print(color.CyanString, format, args...)
}

func (l *ColorLogger) Success(format string, args ...interface{}) {
	l.print(color.HiGreenString, format, args...)
}

func (l *ColorLogger) Warning(format string, args ...interface{}) {
	l.print(color.HiYellowString, format, args...)
}

func (l *ColorLogger) Debug(format string, args ...interface{}) {
	l.print(color.HiWhiteString, format, args...)
}

func (l *ColorLogger) Error(format string, args ...interface{}) {
	l.print(color.HiRedString, format, args...)
}
