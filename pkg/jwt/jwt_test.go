//go:build unit

package jwt

import (
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initConfig() JWTConfig {
	return JWTConfig{
		Access: Token{
			PrivateKey: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlCUEFJQkFBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VXNjbGFFKzlaUUg5Q2VpOGIxcUVmCnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUUpCQUw4ZjRBMUlDSWEvQ2ZmdWR3TGMKNzRCdCtwOXg0TEZaZXMwdHdtV3Vha3hub3NaV0w4eVpSTUJpRmI4a25VL0hwb3piTnNxMmN1ZU9wKzVWdGRXNApiTlVDSVFENm9JdWxqcHdrZTFGY1VPaldnaXRQSjNnbFBma3NHVFBhdFYwYnJJVVI5d0loQVBOanJ1enB4ckhsCkUxRmJxeGtUNFZ5bWhCOU1HazU0Wk1jWnVjSmZOcjBUQWlFQWhML3UxOVZPdlVBWVd6Wjc3Y3JxMTdWSFBTcXoKUlhsZjd2TnJpdEg1ZGdjQ0lRRHR5QmFPdUxuNDlIOFIvZ2ZEZ1V1cjg3YWl5UHZ1YStxeEpXMzQrb0tFNXdJZwpQbG1KYXZsbW9jUG4rTkVRdGhLcTZuZFVYRGpXTTlTbktQQTVlUDZSUEs0PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQ==",
			PublicKey:  "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBTzVIKytVM0xrWC91SlRvRHhWN01CUURXSTdGU0l0VQpzY2xhRSs5WlFIOUNlaThiMXFFZnJxR0hSVDVWUis4c3UxVWtCUVpZTER3MnN3RTVWbjg5c0ZVQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==",
			ExpiresIn:  5,
			MaxAge:     5,
		},
		Refresh: Token{
			PrivateKey: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlCT1FJQkFBSkJBSWFJcXZXeldCSndnYjR1SEhFQ01RdHFZMTI5b2F5RzVZMGlGcG51a0J1VHpRZVlQWkE4Cmx4OC9lTUh3Rys1MlJGR3VxMmE2N084d2s3TDR5dnY5dVY4Q0F3RUFBUUpBRUZ6aEJqOUk3LzAxR285N01CZUgKSlk5TUJLUEMzVHdQQVdwcSswL3p3UmE2ZkZtbXQ5NXNrN21qT3czRzNEZ3M5T2RTeWdsbTlVdndNWXh6SXFERAplUUloQVA5UStrMTBQbGxNd2ZJbDZtdjdTMFRYOGJDUlRaZVI1ZFZZb3FTeW40YmpBaUVBaHVUa2JtZ1NobFlZCnRyclNWZjN0QWZJcWNVUjZ3aDdMOXR5MVlvalZVRlVDSUhzOENlVHkwOWxrbkVTV0dvV09ZUEZVemhyc3Q2Z08KU3dKa2F2VFdKdndEQWlBdWhnVU8yeEFBaXZNdEdwUHVtb3hDam8zNjBMNXg4d012bWdGcEFYNW9uUUlnQzEvSwpNWG1heWtsaFRDeWtXRnpHMHBMWVdkNGRGdTI5M1M2ZUxJUlNIS009Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0t",
			PublicKey:  "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBSWFJcXZXeldCSndnYjR1SEhFQ01RdHFZMTI5b2F5Rwo1WTBpRnBudWtCdVR6UWVZUFpBOGx4OC9lTUh3Rys1MlJGR3VxMmE2N084d2s3TDR5dnY5dVY4Q0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==",
			ExpiresIn:  10,
			MaxAge:     10,
		},
	}
}

func TestGenerateTokenPair(t *testing.T) {
	c := initConfig()

	issuedAt := time.Now()

	tp, err := GenerateTokenPair(c, 1)
	assert.NoError(t, err)

	at := tp.AccessToken
	token, err := ValidateToken(at, c.Access.PublicKey)
	assert.NoError(t, err)
	expectedMap := map[string]interface{}{
		"user_id": 1.0,
	}
	exp := &UserClaims{
		ID: int64(1),
	}
	err = mapstructure.Decode(expectedMap, &exp)
	assert.NoError(t, err)
	assert.Equal(t, exp, token)
	expiredAt := time.Now()
	require.WithinDuration(t, issuedAt, expiredAt, time.Second)
}

func TestCreateToken(t *testing.T) {
	c := initConfig()
	at, err := CreateToken(c.Access.ExpiresIn, "1", c.Access.PrivateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, at)
	if len(at) < 180 {
		t.Error("token not valid")
	}
}

func TestDecodeToken(t *testing.T) {
	c := initConfig()
	badToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	_, err := DecodeToken(badToken, c.Access.PublicKey)
	assert.Error(t, err)
}
