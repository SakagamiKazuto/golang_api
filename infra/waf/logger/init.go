package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = NewLogger()

func init() {
	defer Log.Info("Logger Settints is set")

	Log.SetFormat(&logrus.JSONFormatter{})

	switch os.Getenv("APP_MODE") {
	case "production", "staging":
		Log.SetLogLevel(logrus.InfoLevel)
	default:
		Log.SetLogLevel(logrus.DebugLevel)
	}
}
