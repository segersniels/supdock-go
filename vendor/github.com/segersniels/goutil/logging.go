package util

import (
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

// Log : custom TTY logging
func Log(a ...interface{}) {
	logrus.Info(strings.TrimSpace(spew.Sdump(a...)))
}

// Warn : custom TTY warning logging
func Warn(a ...interface{}) {
	logrus.Warn(a...)
}

// Error : custom TTY error logging and throw exit code 1
func Error(a ...interface{}) {
	logrus.Error(a...)
	os.Exit(1)
}

// Fatal : custom TTY fatal logging
func Fatal(a ...interface{}) {
	logrus.Fatal(a...)
}

// Panic : custom TTY panic logging
func Panic(a ...interface{}) {
	logrus.Panic(a...)
}

// Fields : fields type needed to pass to LogWithFields
type Fields map[string]interface{}

// LogWithFields : custom logging with fields
func LogWithFields(fields map[string]interface{}, value string) {
	logrus.WithFields(fields).Info(value)
}

// WarnWithFields : custom warning with fields
func WarnWithFields(fields map[string]interface{}, value string) {
	logrus.WithFields(fields).Warn(value)
}

// ErrorWithFields : custom error with fields and exit code 0
func ErrorWithFields(fields map[string]interface{}, value string) {
	logrus.WithFields(fields).Error(value)
	os.Exit(0)
}

// FatalWithFields : custom panic with fields
func FatalWithFields(fields map[string]interface{}, value string) {
	logrus.WithFields(fields).Fatal(value)
}

// PanicWithFields : custom panic with fields
func PanicWithFields(fields map[string]interface{}, value string) {
	logrus.WithFields(fields).Panic(value)
}
