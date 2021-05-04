package logger

import (
	"fmt"
	"github.com/SakagamiKazuto/golang_api/interface/database"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type logger struct {
	logger *logrus.Logger
}

// NewLogger returns a new Logger logging to out.
func NewLogger() database.Logger {
	l := logrus.New()
	return logger{l}
}

func (l logger) Info(args ...interface{}) {
	if l.logger.Level >= logrus.InfoLevel {
		entry := l.WithFields(database.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

func (l logger) InfoWithFields(inf interface{}, f database.Fields) {
	if l.logger.Level >= logrus.InfoLevel {
		entry := l.WithFields(f)
		entry.Data["file"] = fileInfo(2)
		entry.Info(inf)
	}
}

func (l logger) Warn(args ...interface{}) {
	if l.logger.Level >= logrus.WarnLevel {
		entry := l.WithFields(database.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

func (l logger) WarnWithFields(inf interface{}, f database.Fields) {
	if l.logger.Level >= logrus.WarnLevel {
		entry := l.WithFields(f)
		entry.Data["file"] = fileInfo(2)
		entry.Warn(inf)
	}
}

func (l logger) Fatal(args ...interface{}) {
	if l.logger.Level >= logrus.FatalLevel {
		entry := l.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

func (l logger) FatalWithFields(inf interface{}, f database.Fields) {
	if l.logger.Level >= logrus.FatalLevel {
		entry := l.logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(inf)
	}
}

func (l logger) SetLogLevel(level logrus.Level) {
	l.logger.Level = level
}

func (l logger) SetFormat(formatter logrus.Formatter) {
	l.logger.Formatter = formatter
}

func (l logger) WithFields(f database.Fields) *logrus.Entry {
	return l.logger.WithFields(f)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
