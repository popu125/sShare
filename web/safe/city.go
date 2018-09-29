package safe

import (
	"github.com/hashicorp/golang-lru"
	"github.com/ipipdotnet/datx-go"
	"net/http"
	"regexp"
	"strings"
)

const (
	cacheLength = 1024
)

var xffRegex, _ = regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")

type cityLimit struct {
	targetCity string

	db        *datx.City
	loggedIps *lru.Cache

	loadXFF bool
}

func NewCityLimit(path string) *cityLimit {
	db, err := datx.NewCity(path)
	if err != nil {
		panic(err)
	}

	cache, _ := lru.New(cacheLength)
	return &cityLimit{db: db, loggedIps: cache}
}

func (cl *cityLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string
		if cl.loadXFF {
			ip = xffRegex.FindString(r.Header.Get("X-Forwarded-For"))
		} else {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		}

		if authed, ok := cl.loggedIps.Get(ip); ok {
			if authed.(bool) {
				next.ServeHTTP(w, r)
				return
			}
		}

		location, err := cl.db.FindLocation(ip)
		if err != nil {
			return
		}
		if location.City == cl.targetCity {
			cl.loggedIps.Add(ip, true)
			next.ServeHTTP(w, r)
			return
		} else {
			cl.loggedIps.Add(ip, false)
		}

		http.Error(w, "I'm a teapot", 418)
	})
}
