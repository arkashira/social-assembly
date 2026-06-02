package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type EncryptedUserData struct {
	Email           string `json:"email"`
	IP              string `json:"ip"`
	PrivateMessages string `json:"private_messages"`
}

// EncryptUserData encrypts user data using AES-GCM (authenticated encryption).
func EncryptUserData(userData *User) (*EncryptedUserData, error) {
	if userData == nil {
		return nil, errors.New("user data is nil")
	}

	// Generate a random nonce (12 bytes for AES-GCM)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher([]byte(AESKey))
	if err != nil {
		return nil, err
	}

	// Create GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Encrypt each field separately (or combine into a single struct)
	encryptField := func(data string) (string, error) {
		if data == "" {
			return "", nil // Skip empty fields
		}
		encrypted := gcm.Seal(nil, nonce, []byte(data), nil)
		return base64.StdEncoding.EncodeToString(encrypted), nil
	}

	emailEnc, err := encryptField(userData.Email)
	if err != nil {
		return nil, err
	}

	ipEnc, err := encryptField(userData.IP)
	if err != nil {
		return nil, err
	}

	pmEnc, err := encryptField(userData.PrivateMessages)
	if err != nil {
		return nil, err
	}

	return &EncryptedUserData{
		Email:           emailEnc,
		IP:              ipEnc,
		PrivateMessages: pmEnc,
	}, nil
}