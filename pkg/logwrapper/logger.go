package logwrapper

import (
	"os"
	"github.com/sirupsen/logrus"
)

func NewLogger(name string, debug bool) *logrus.Entry {
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
		DisableColors: true, // sets logfmt format
		FullTimestamp: true,
	})

	return logger.WithFields(logrus.Fields{ "name": name })
}
