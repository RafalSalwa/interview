package generator

import (
	crand "crypto/rand"
	"errors"
	"io"
	"math/big"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	minLength   = 32
	maxLength   = 64
)

var ErrLengthInvalid = errors.New("code should be between 6 and 10 letters")

func RandomString(length int) (string, error) {
	if length < minLength || length > maxLength {
		return "", ErrLengthInvalid
	}

	b := make([]byte, length)
	_, _ = crand.Read(b)
	_, _ = io.ReadFull(crand.Reader, b)

	for i := 0; i < length; i++ {
		num, _ := crand.Int(crand.Reader, big.NewInt(int64(len(letterBytes))))
		b[i] = letterBytes[num.Int64()]
	}

	return string(b), nil
}
