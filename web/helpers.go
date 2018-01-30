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

func newUUID() string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	uuid := []byte("")
	for i := 36; i > 0; i-- {
		switch i {
		case 27, 22, 17, 12:
			uuid = append(uuid, 45) // 45 is "-"
		default:
			uuid = append(uuid, chars[rand.Intn(36)])
		}
	}
	return string(uuid)
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
