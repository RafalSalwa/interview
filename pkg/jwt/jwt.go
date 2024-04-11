package jwt

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenPair struct {
	AccessToken  string `mapstructure:"accessToken"`
	RefreshToken string `mapstructure:"refreshToken"`
}

type JWTConfig struct {
	Access  Token `mapstructure:"accessToken"`
	Refresh Token `mapstructure:"refreshToken"`
}

type Token struct {
	PrivateKey string        `mapstructure:"privateKey"`
	PublicKey  string        `mapstructure:"publicKey"`
	ExpiresIn  time.Duration `mapstructure:"expiresIn"`
	MaxAge     int           `mapstructure:"maxAge"`
}

func CreateToken(ttl time.Duration, payload int64, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = strconv.FormatInt(payload, 10)
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func DecodeToken(token, publicKey string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	_, err = jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return t, nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

func ValidateToken(token string, publicKey string) (string, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	cc := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, cc, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return "", fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return "", fmt.Errorf("validate: %w", err)
	}
	err = cc.Valid()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", cc["sub"]), nil
}
