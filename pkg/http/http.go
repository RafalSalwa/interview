package http

import "time"

type Config struct {
	Addr                string        `mapstructure:"addr"`
	Development         bool          `mapstructure:"development"`
	BasePath            string        `mapstructure:"basePath"`
	DebugHeaders        bool          `mapstructure:"debugHeaders"`
	HTTPClientDebug     bool          `mapstructure:"httpClientDebug"`
	DebugErrorsResponse bool          `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string      `mapstructure:"ignoreLogUrls"`
	TimeoutRead         time.Duration `mapstructure:"SERVER_TIMEOUT_READ"`
	TimeoutWrite        time.Duration `mapstructure:"SERVER_TIMEOUT_WRITE"`
	TimeoutIdle         time.Duration `mapstructure:"SERVER_TIMEOUT_IDLE"`
}
