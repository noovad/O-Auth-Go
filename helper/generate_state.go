package helper

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
