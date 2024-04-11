package encdec

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func Decrypt(cipherText string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherTextByte, _ := base64.StdEncoding.DecodeString(cipherText)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherTextByte))
	cfb.XORKeyStream(plainText, cipherTextByte)
	return string(plainText), nil
}
