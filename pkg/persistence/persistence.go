package persistence

type Config struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
	PoolSize int    `mapstructure:"poolSize"`
}
