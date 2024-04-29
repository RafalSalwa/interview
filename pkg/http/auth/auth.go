package auth

import (
	"time"
)

type (
	Auth struct {
		AuthMethod  string    `mapstructure:"method"`
		APIKey      string    `mapstructure:"apiKey"`
		BearerToken string    `mapstructure:"bearerToken"`
		BasicAuth   BasicAuth `mapstructure:"basic"`
		JWTToken    JWTConfig `mapstructure:"jwt"`
	}
	BasicAuth struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
	JWTConfig struct {
		Access  Token `mapstructure:"accessToken"`
		Refresh Token `mapstructure:"refreshToken"`
	}
	Token struct {
		PrivateKey string        `mapstructure:"privateKey"`
		PublicKey  string        `mapstructure:"publicKey"`
		ExpiresIn  time.Duration `mapstructure:"expiresIn"`
		MaxAge     int           `mapstructure:"maxAge"`
	}
)

func NewAuthorizer(cfg *Auth) (IAuthType, error) {
	at, err := NewAuthMethod(cfg)
	if err != nil {
		return nil, err
	}
	return at, nil
}
