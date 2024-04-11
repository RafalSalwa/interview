package encdec

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const (
	key string = "RHCaDQMfmrBhuaKLDmBNGOEgERuGzwSi"
)

var bytes = []byte{104, 32, 18, 239, 16, 250, 111, 197, 34, 150, 248, 7, 222, 146, 58, 151}

func Encrypt(plaintext string) string {
	block, _ := aes.NewCipher([]byte(key))
	mode := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plaintext))
	mode.XORKeyStream(cipherText, []byte(plaintext))

	return base64.StdEncoding.EncodeToString(cipherText)
}
