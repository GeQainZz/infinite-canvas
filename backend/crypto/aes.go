package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func Encrypt(key, plaintext string) (string, error) {
	if key == "" {
		return plaintext, nil
	}
	block, err := aes.NewCipher([]byte(padKey(key)))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key, ciphertext string) (string, error) {
	if key == "" {
		return ciphertext, nil
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errors.New("invalid encrypted data")
	}
	block, err := aes.NewCipher([]byte(padKey(key)))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func padKey(key string) string {
	const keyLen = 32 // AES-256
	if len(key) >= keyLen {
		return key[:keyLen]
	}
	padded := make([]byte, keyLen)
	copy(padded, key)
	return string(padded)
}
