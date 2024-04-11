package grpc

type Config struct {
	Addr        string `mapstructure:"addr"`
	Development bool   `mapstructure:"development"`
}
