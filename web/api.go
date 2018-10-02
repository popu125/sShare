package web

import (
	"fmt"
	"github.com/popu125/sShare/config"
	"github.com/popu125/sShare/pool"
	"log"
	"net/http"
)

type ApiServe struct {
	pool    *pool.Pool
	captcha Captcha
	conf    *config.Config

	l *log.Logger
}

func NewApiServe(conf *config.Config, l log.Logger) (*ApiServe, *pool.Pool) {
	p := pool.NewPool(conf, l)
	captcha := NewCaptcha(conf.Captcha)
	l.SetPrefix("[WEB] ")
	api := &ApiServe{
		pool:    p,
		captcha: captcha,
		l:       &l,
		conf:    conf,
	}

	return api, p
}

func (apis *ApiServe) serveCount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, apis.pool.Count())
}

func (apis *ApiServe) newProc(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := r.PostForm.Get("token")
	if len(token) == 0 || !apis.captcha.Validate(token) {
		apis.l.Println("ERR_NO_CAPTCHA", r.RemoteAddr)
		w.Write(transInfo{Status: "ERR_NO_CAPTCHA"}.Json())
		return
	}

	port, info, err := apis.pool.NewProc()
	if err != nil {
		apis.l.Println("ERR_FULL", r.RemoteAddr)
		w.Write(transInfo{Status: err.Error()}.Json())
		return
	}

	apis.l.Println("ACCEPT", r.RemoteAddr, port)
	w.WriteHeader(http.StatusCreated)
	w.Write(transInfo{Status: "ACCEPT", Port: port, Pass: info}.Json())
	return
}
