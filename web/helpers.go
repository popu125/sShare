package web

import (
	"math/rand"
	"encoding/json"
)

var plist = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func newPass() string {
	pass := make([]byte, 10)
	for i := 0; i < 10; i++ {
		pass[i] = plist[rand.Intn(62)]
	}
	return string(pass)
}

type transInfo struct {
	Status string
	Port   int
	Pass   string
}

func (self transInfo) Json() []byte {
	j, err := json.Marshal(self)
	if err != nil {
		return []byte("ERR_UNKNOWN")
	}
	return j
}
