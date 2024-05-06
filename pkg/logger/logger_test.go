package logger

import (
	"testing"
)

func TestNewConsole(t *testing.T) {
	_ = NewConsole()
}

func TestLogger_Print(t *testing.T) {
	logger := NewConsole()
	// Call the Print method
	logger.Print("Test message")
	logger.Error().Msg("Error message")
	logger.Warn().Msg("Warn message")
	logger.Info().Msg("Info message")
	logger.Debug().Msg("Debug message")
	logger.Print("Test message")
	logger.Printf("Test %s", "message")
	logger.Log().Msg("Test message")
	// logger.Fatal().Msg("Fatal message")
}
