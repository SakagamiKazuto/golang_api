package logger

import (
	"github.com/sirupsen/logrus"
)

func init() {
	// JSONフォーマット
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
