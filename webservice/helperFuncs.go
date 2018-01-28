package webservice

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
)

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func pprint(v interface{}, label string) {
	res1B, _ := json.MarshalIndent(v, "", "    ")
	log.Println(label, "\n", string(res1B))
}
