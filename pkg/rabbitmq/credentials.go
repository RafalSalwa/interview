package rabbitmq

type Credentials struct {
	Protocol string
	Username string
	Password string
	Addr     string
	VHost    string
	Exchange *Exchange
}

func (c *Credentials) GetURL() string {
	result := "amqp://"
	if c.Username != "" && c.Password != "" {
		result = result + c.Username + ":" + c.Password + "@"
	}
	result = result + c.Addr + "/" + c.VHost
	return result
}
