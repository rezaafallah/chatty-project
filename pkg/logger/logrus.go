// pkg/logger/logrus.go
package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

// Setup creates a standardized logger
func Setup() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	return l
}