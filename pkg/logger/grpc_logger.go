package logger

import (
	"github.com/sirupsen/logrus"
)

func NewGRPCLogger() *logrus.Entry {
	logger := logrus.New()
	entry := logrus.NewEntry(logger)
	return entry
}
