package framework

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger handles writing to a log.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// logWrapper handles using the logger.
type logWrapper struct {
	logger *logrus.Entry
}

// NewLogger creates a logger with the specified name.
func NewLogger(name string, debug bool) Logger {
	logger := logrus.New()

	if debug {
		logger.SetLevel(logrus.DebugLevel)

		// adds func and file fields, has small runtime overhead
		// logger.SetReportCaller(true)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		// DisableColors: true, // sets logfmt format
		FullTimestamp: true,
	})

	return &logWrapper{logger.WithFields(logrus.Fields{"name": name})}
}

// Debug writes a static message at debug level.
func (l *logWrapper) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf writes a formatted message at debug level.
func (l *logWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info writes a static message at info level.
func (l *logWrapper) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof writes a formatted message at debug level.
func (l *logWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn writes a static message at warn level.
func (l *logWrapper) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf writes a formatted message at warn level.
func (l *logWrapper) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error writes a static message at error level.
func (l *logWrapper) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf writes a formatted message at error level.
func (l *logWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal writes a static message at fatal level, then quits.
func (l *logWrapper) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf writes a formatted message at fatal level, then quits.
func (l *logWrapper) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
