package pool

import (
	"encoding/json"
	"math/rand"
)

var plist = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type runArgData struct {
	PassMap map[string]string `json:"pass_map"`
	UUIDMap map[string]string `json:"uuid_map"`

	Port int `json:"port"`
}

func newArgData(port int) *runArgData {
	return &runArgData{
		PassMap: make(map[string]string, 64),
		UUIDMap: make(map[string]string, 64),
		Port:    port,
	}
}

func (g *runArgData) Pass(name string) string {
	if pass, ok := g.PassMap[name]; ok {
		return pass
	}

	pass := make([]byte, 10)
	for i := 0; i < 10; i++ {
		pass[i] = plist[rand.Intn(62)]
	}
	g.PassMap[name] = string(pass)
	return string(pass)
}

func (g *runArgData) UUID(name string) string {
	if uuid, ok := g.UUIDMap[name]; ok {
		return uuid
	}

	uuid := []byte("")
	for i := 36; i > 0; i-- {
		switch i {
		case 27, 22, 17, 12:
			uuid = append(uuid, 45) // 45 is "-"
		default:
			uuid = append(uuid, chars[rand.Intn(36)])
		}
	}
	g.UUIDMap[name] = string(uuid)
	return string(uuid)
}

func (g runArgData) Data() string {
	if tmp, err := json.Marshal(g); err != nil {
		return ""
	} else {
		return string(tmp)
	}
}
