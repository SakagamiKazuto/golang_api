package logger

import (
	"github.com/sirupsen/logrus"
)

func init() {
	// JSONフォーマット
	logrus.SetFormatter(&logrus.JSONFormatter{})

	//// 標準エラー出力でなく標準出力とする
	//log.SetOutput(os.Stdout)

	//// Warningレベル以上を出力
	//log.SetLevel(log.WarnLevel)
}
