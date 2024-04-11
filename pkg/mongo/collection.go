package mongo

type Collection struct {
	Name string `mapstructure:"name"`
}

type Collections struct {
	Collections []Collection `mapstructure:"name"`
}
