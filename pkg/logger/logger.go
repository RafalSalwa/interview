package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger      *zerolog.Logger
	nrLoggerApp *newrelic.Application
}

type Config struct {
	LogLevel string `mapstructure:"level"`
	DevMode  bool   `mapstructure:"devMode"`
}

func NewConsole() *Logger {
	logLevel := zerolog.InfoLevel

	zerolog.SetGlobalLevel(logLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("interview-srv"),
		newrelic.ConfigLicense("asd"),
		newrelic.ConfigAppLogDecoratingEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(false),
		func(config *newrelic.Config) {
			config.Enabled = true
		},
	)

	return &Logger{logger: &logger, nrLoggerApp: app}
}

func (l *Logger) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

func (l *Logger) Level(level zerolog.Level) zerolog.Logger {
	return l.logger.Level(level)
}

func (l *Logger) Sample(s zerolog.Sampler) zerolog.Logger {
	return l.logger.Sample(s)
}

func (l *Logger) Hook(h zerolog.Hook) zerolog.Logger {
	return l.logger.Hook(h)
}

func (l *Logger) Err(err error) *zerolog.Event {
	return l.logger.Err(err)
}

func (l *Logger) Trace() *zerolog.Event {
	return l.logger.Trace()
}

func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}
func (l *Logger) WithTracing(span trace.Span) *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

func (l *Logger) Log() *zerolog.Event {
	return l.logger.Log()
}

func (l *Logger) Print(v ...interface{}) {
	l.logger.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}

func (l *Logger) Ctx(ctx context.Context) *Logger {
	return &Logger{logger: zerolog.Ctx(ctx)}
}

func (l *Logger) Println(v ...interface{}) {
	l.Printf("%+v\n", v...)
}
