package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/popu125/sShare/config"
	"github.com/popu125/sShare/web"
)

func main() {
	confFn := flag.String("c", "config.json", "Config file to load")
	logFn := flag.String("log", "", "File to save log")
	flag.Parse()
	var l *log.Logger
	if *logFn == "" {
		l = log.New(os.Stderr, "[MAIN] ", log.LstdFlags)
	} else {
		f, err := os.Create(*logFn)
		if err != nil {
			fmt.Println("Connot open log file for writing.")
			os.Exit(1)
		}
		l = log.New(f, "[MAIN] ", log.LstdFlags)
		defer f.Close()
	}

	conf := config.LoadConfig(*confFn)
	api, pool := web.NewApiServe(conf, *l)
	defer pool.Purge()
	route := web.GetRouter(api)
	rand.Seed(time.Now().Unix() + conf.RandSeed)

	http.Handle("/", route)
	l.Println("sShare Server running on:", conf.Addr)
	l.Fatal(http.ListenAndServe(conf.Addr, nil))
}
