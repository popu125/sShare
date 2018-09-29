package safe

import (
	"encoding/hex"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru"
)

var (
	anticc_ret      = []byte(`<html><head><meta http-equiv="refresh" content="3;url=/"></head></html>`)
	anticc_api      = []byte(`{"Status":"ERR_ANTICC"}`)
	anticc_redirect = []byte(`<html><head><meta http-equiv="refresh" content="3;url=/get_cctoken"></head></html>`)
)

type anticc struct {
	tokens *lru.Cache
}

func NewAntiCC(size int) *anticc {
	cache, err := lru.New(size)
	if err != nil {
		panic(err)
	}
	return &anticc{tokens: cache}
}

func (acc *anticc) Redirect(w http.ResponseWriter, r *http.Request) {
	tmp := make([]byte, 64)
	rand.Read(tmp)
	token := hex.EncodeToString(tmp)

	expires := time.Now().Add(time.Minute * 20)
	acc.tokens.Add(token, &expires)

	http.SetCookie(w, &http.Cookie{
		Name:    "cctoken",
		Value:   token,
		Expires: time.Now().Add(time.Minute * 20),
	})

	w.Write(anticc_ret)
}

func (acc *anticc) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/get_cctoken" {
			next.ServeHTTP(w, r)
			return
		}
		tokenC, err := r.Cookie("cctoken")
		if err != nil {
			w.Write(anticc_redirect)
			return
		}
		token := tokenC.Value

		if value, found := acc.tokens.Get(token); found && value.(*time.Time).After(time.Now()) {
			next.ServeHTTP(w, r)
			return
		} else {
			if strings.HasPrefix(r.URL.Path, "/api") {
				w.Write(anticc_api)
			} else {
				w.Write(anticc_redirect)
			}
		}
	})
}
