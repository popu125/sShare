package web

import (
	"encoding/json"
)

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
