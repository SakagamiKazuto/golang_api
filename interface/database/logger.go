package database

import "github.com/sirupsen/logrus"

// map[string]interface{}として宣言されたlogrusのFieldsをwrapする
type Fields = logrus.Fields

type Logger interface {
	Info(args ...interface{})
	InfoWithFields(inf interface{}, f Fields)

	Warn(args ...interface{})
	WarnWithFields(inf interface{}, f Fields)

	Fatal(args ...interface{})
	FatalWithFields(inf interface{}, f Fields)

	WithFields(Fields) *logrus.Entry

	SetFormat(formatter logrus.Formatter)
	SetLogLevel(level logrus.Level)
}
