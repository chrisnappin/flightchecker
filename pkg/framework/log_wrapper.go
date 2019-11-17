package framework

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LogWrapper wraps the logrus logger in methods compatible with domain.Logger.
type LogWrapper struct {
	logger *logrus.Entry
}

// NewLogWrapper creates a logger with the specified name.
func NewLogWrapper(name string, debug bool) *LogWrapper {
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

	return &LogWrapper{logger.WithFields(logrus.Fields{"name": name})}
}

// Debug writes a static message at debug level.
func (l *LogWrapper) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf writes a formatted message at debug level.
func (l *LogWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info writes a static message at info level.
func (l *LogWrapper) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof writes a formatted message at debug level.
func (l *LogWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn writes a static message at warn level.
func (l *LogWrapper) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf writes a formatted message at warn level.
func (l *LogWrapper) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error writes a static message at error level.
func (l *LogWrapper) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf writes a formatted message at error level.
func (l *LogWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal writes a static message at fatal level, then quits.
func (l *LogWrapper) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf writes a formatted message at fatal level, then quits.
func (l *LogWrapper) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
