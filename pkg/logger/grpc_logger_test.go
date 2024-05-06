package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewGRPCLogger(t *testing.T) {
	// Call the NewGRPCLogger function
	entry := NewGRPCLogger()

	// Assert that the returned entry is not nil
	assert.NotNil(t, entry, "NewGRPCLogger should return a non-nil logrus entry")

	// Assert that the logger instance in the entry is of type *logrus.Logger
	assert.IsType(t, logrus.New(), entry.Logger, "Logger instance should be of type *logrus.Logger")
}
