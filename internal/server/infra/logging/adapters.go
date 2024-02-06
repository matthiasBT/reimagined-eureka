package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLogger() ILogger {
	return setupLogrusLogger()
}

func setupLogrusLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}
