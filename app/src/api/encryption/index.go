package encryption

import (
	"cloudview/app/src/api/middleware/logger"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

/*
To encrypt data correctly it is recommended to use a
16 byte secret key.

The encryption will fail otherwise.
*/
func Encrypt(text string, secret string) (string, error) {
	logger.Logger.Log("Encrypting string")
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		logger.Logger.Error("Unable to encrypt data", err)
		return "", err
	}
	plainText := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logger.Logger.Error("Unable to encrypt data", err)
		return "", err
	}
	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	logger.Logger.Log("encryption success")
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(cipherstring string, keystring string) (string, error) {
	logger.Logger.Log("decrypting cipher string...")
	// Byte array of the string
	ciphertext, err := hex.DecodeString(cipherstring)
	if err != nil {
		logger.Logger.Error("Unable to decode hex string", err)
		return "", err
	}
	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Logger.Error("Unable to decrypt string", err)
		return "", err
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		logger.Logger.Error("Unable to decrypt string: Text is too short")
		return "", errors.New("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	logger.Logger.Log("decrypt success")
	return string(ciphertext), nil
}

func GenerateRandomSecretKey(bytes int) (string, error) {
	key := make([]byte, bytes)
	_, err := rand.Read(key)
	if err != nil {
		logger.Logger.Error("Unable to generate random key", err)
		return "", err
	}
	return hex.EncodeToString(key), nil
}
