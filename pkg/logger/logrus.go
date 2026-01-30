package logger

import (
	"os"
	"sync"
	"github.com/sirupsen/logrus"
)
// Logger Interface:
type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Warn(args ...interface{})
	WithError(err error) Logger
}

type logrusLogger struct {
	*logrus.Entry 
}

var (
	instance Logger
	once     sync.Once
)

// Setup
func Setup() Logger {
	once.Do(func() {
		l := logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{})
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.DebugLevel)

		instance = &logrusLogger{
			Entry: logrus.NewEntry(l),
		}
	})
	return instance
}

// (Wrapper Methods)
func (l *logrusLogger) Info(args ...interface{}) { l.Entry.Info(args...) }

func (l *logrusLogger) Infof(format string, args ...interface{}) { l.Entry.Infof(format, args...) }

func (l *logrusLogger) Error(args ...interface{}) { l.Entry.Error(args...) }

func (l *logrusLogger) Errorf(format string, args ...interface{}) { l.Entry.Errorf(format, args...) }

func (l *logrusLogger) Debug(args ...interface{}) { l.Entry.Debug(args...) }

func (l *logrusLogger) Debugf(format string, args ...interface{}) { l.Entry.Debugf(format, args...) }

func (l *logrusLogger) Warn(args ...interface{}) { l.Entry.Warn(args...) }

func (l *logrusLogger) WithError(err error) Logger {
	return &logrusLogger{
		Entry: l.Entry.WithError(err),
	}
}