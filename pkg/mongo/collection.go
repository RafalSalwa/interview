package mongo

type (
	Collection struct {
		Name string `mapstructure:"name"`
	}

	Collections struct {
		Collections []Collection `mapstructure:"name"`
	}
)
