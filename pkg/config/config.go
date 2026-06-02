package config

import (
	"errors"
	"os"
)

var AESKey string

// LoadConfig loads the AES key from environment variables (secure) or a fallback.
func LoadConfig() error {
	// Prefer environment variable (most secure)
	if key := os.Getenv("AES_KEY"); key != "" {
		AESKey = key
		return nil
	}

	// Fallback: Load from file (less secure, but works in dev)
	// (Implementation omitted for brevity; use `ioutil.ReadFile` + `json.Unmarshal`)
	return errors.New("AES_KEY environment variable not set")
}