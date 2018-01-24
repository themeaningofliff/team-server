package webservice

import (
	"crypto/rand"
	"encoding/base64"
)

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
