package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

var Key string

func Encrypt(data []byte) ([]byte, error) {
	secKeyBytes := sha256.Sum256([]byte(Key))
	secKey := secKeyBytes[0:32]

	block, err := aes.NewCipher(secKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	cipheredRaw := gcm.Seal(nonce, nonce, data, nil)
	return cipheredRaw, nil
}

func Decrypt(ciphertext []byte) ([]byte, error) {
	secKeyBytes := sha256.Sum256([]byte(Key))
	secKey := secKeyBytes[0:32]

	block, err := aes.NewCipher(secKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("too short message")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	cipherdata := ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, cipherdata, nil)
}
