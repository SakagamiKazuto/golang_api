package logger

import (
	"github.com/sirupsen/logrus"
)

var Log = NewLogger()

func init() {
	defer Log.Info("Logger Settints is set")
	// JSONフォーマット
	Log.SetFormat(&logrus.JSONFormatter{})
}
