package encryption

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"log"
)

// Encode a string to Base64
func encodeToString(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

// GenerateCryptoID generate random string of 6 bytes and encode to base64 format
func GenerateCryptoID() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("ERROR_WHILE_GNERATING_UNIQ_ID")
		return "", err
	}
	return encodeToString(hex.EncodeToString(bytes)), nil
}
