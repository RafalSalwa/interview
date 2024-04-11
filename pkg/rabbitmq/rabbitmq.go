package rabbitmq

type Config struct {
	Addr     string    `mapstructure:"addr"`
	Username string    `mapstructure:"username"`
	Password string    `mapstructure:"password"`
	VHost    string    `mapstructure:"vhost"`
	Exchange *Exchange `mapstructure:"exchange"`
}
type Exchange struct {
	Name       string `mapstructure:"name"`
	Type       string `mapstructure:"type"`
	Durable    bool   `mapstructure:"durable"`
	Queue      string `mapstructure:"queue"`
	RoutingKey string `mapstructure:"routing_key"`
}
