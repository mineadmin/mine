package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os/exec"
)

// GenerateJwtSecret generates a random JWT secret
func GenerateJwtSecret() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// GenerateAppKey generates a random application key
func GenerateAppKey() (string, error) {
	bytes := make([]byte, 16) // 128 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// CheckCommandExists checks if a command exists in the system
func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
