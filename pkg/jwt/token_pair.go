package jwt

func GenerateTokenPair(c JWTConfig, uid int64) (*TokenPair, error) {
	t, err := CreateToken(c.Access.ExpiresIn, uid, c.Access.PrivateKey)
	if err != nil {
		return nil, err
	}

	rt, err := CreateToken(c.Refresh.ExpiresIn, uid, c.Refresh.PrivateKey)
	if err != nil {
		return nil, err
	}

	tp := &TokenPair{AccessToken: t, RefreshToken: rt}
	return tp, nil
}
