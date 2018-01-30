package web

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"sync"
	"strconv"
	"log"

	"github.com/popu125/sShare/config"
	"github.com/popu125/sShare/pool"
)

type ApiServe struct {
	pool    *pool.Pool
	captcha Captcha

	cmd  string
	args []string

	ports     []int
	plock     sync.Mutex
	portStart int
	portLimit int

	l *log.Logger

	passGen func() string
}

func NewApiServe(conf config.Config, l log.Logger) *ApiServe {
	p := pool.NewPool(conf, l)
	captcha := NewCaptcha(conf.Captcha)
	l.SetPrefix("[WEB] ")
	api := &ApiServe{
		pool:      p,
		captcha:   captcha,
		portStart: int(conf.PortStart),
		portLimit: int(conf.PortRange),
		l:         &l,
	}
	if !conf.GenUUID {
		api.passGen = newPass
	} else {
		api.passGen = newUUID
	}
	return api
}

func (self *ApiServe) serveCount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, self.pool.Count())
}

func (self *ApiServe) newProc(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := r.PostForm.Get("token")
	if len(token) == 0 || !self.captcha.Validate(token) {
		self.l.Println("ERR_NO_CAPTCHA", r.RemoteAddr)
		w.Write(transInfo{Status: "ERR_NO_CAPTCHA"}.Json())
		return
	}

	var port int
	self.plock.Lock()
OUT:
	for {
		port = self.portStart + rand.Intn(self.portLimit)
		for _, p := range self.ports {
			if p == port {
				continue OUT
			}
		}
		break
	}

	pass := self.passGen()

	self.plock.Unlock()
	status := self.pool.NewProc(strconv.Itoa(port), pass)
	if status == pool.ERR_FULL {
		self.l.Println("ERR_FULL", r.RemoteAddr)
		w.Write(transInfo{Status: "ERR_FULL"}.Json())
		return
	}

	self.l.Println("ACCEPT", r.RemoteAddr, port, pass)
	cookie := &http.Cookie{
		Name:  "checkid",
		Value: strconv.FormatUint(uint64(status), 10),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
	w.Write(transInfo{Status: "ACCEPT", Port: port, Pass: pass}.Json())
	return
}

func (self ApiServe) procCheck(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	pid, err := strconv.Atoi(v["id"])
	if err != nil {
		w.WriteHeader(400)
	} else if self.pool.Check(uint(pid)) {
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("nope"))
	}
	return
}
