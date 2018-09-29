package web

import (
	"encoding/json"
)

type transInfo struct {
	Status string
	Port   int
	Pass   interface{}
}

func (self transInfo) Json() []byte {
	j, err := json.Marshal(self)
	if err != nil {
		return []byte("ERR_UNKNOWN")
	}
	return j
}
