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

func (l *ColorLogger) println(printer func(string, ...interface{}) string, format string, args ...interface{}) {
	if args == nil {
		text := printer(format)
		fmt.Println(text)
	} else {
		text := printer(format, args...)
		fmt.Println(text)
	}
}

func (l *ColorLogger) print(printer func(string, ...interface{}) string, format string, args ...interface{}) {
	if args == nil {
		fmt.Print(printer(format))
	} else {
		fmt.Print(printer(format, args))
	}
}

func (l *ColorLogger) Infoln(format string, args ...interface{}) {
	l.println(color.CyanString, format, args...)
}

func (l *ColorLogger) Info(format string, args ...interface{}) {
	l.print(color.CyanString, format, args...)
}

func (l *ColorLogger) Successln(format string, args ...interface{}) {
	l.println(color.HiGreenString, format, args...)
}

func (l *ColorLogger) Warningln(format string, args ...interface{}) {
	l.println(color.HiYellowString, format, args...)
}

func (l *ColorLogger) Warning(format string, args ...interface{}) {
	l.print(color.HiYellowString, format, args...)
}

func (l *ColorLogger) Debugln(format string, args ...interface{}) {
	l.println(color.HiWhiteString, format, args...)
}

func (l *ColorLogger) Failureln(format string, args ...interface{}) {
	l.println(color.HiRedString, format, args...)
}
