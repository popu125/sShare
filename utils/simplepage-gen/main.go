package main

import (
	"flag"
	"html/template"
	"os"
	"io/ioutil"
	"encoding/json"
)

var (
	tpl                  *template.Template
	name, desc, location string
)

func main() {
	tpl = template.Must(template.ParseFiles("index.tpl"))

	flag.StringVar(&name, "name", "sShare", "Name of site.")
	flag.StringVar(&desc, "desc", "生活不止眼前的苟，还有身后的苟。", "Description of site.")
	flag.StringVar(&location, "l", "server-a.bobo.moe", "Server's location provided to user.")
	helpMsg := flag.Bool("h", false, "Show this message.")
	flag.Parse()
	if *helpMsg {
		flag.Usage()
		return
	}

	f, err := os.Create("index.out.html")
	check(err)
	defer f.Close()
	conf := &config{}
	conf.load("config.json")
	check(tpl.Execute(f, conf))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type config struct {
	SiteName        string
	SiteDescription string
	Location        string

	Captcha   captchaConf `json:"captcha"`
	TTLString string      `json:"ttl"`
	Limit     uint32      `json:"limit"`
}

type captchaConf struct {
	Name   string `json:"name"`
	SiteID string `json:"site_id"`
	Extra  string `json:"extra"`
}

func (self *config) load(fn string) {
	data, err := ioutil.ReadFile(fn)
	check(err)
	check(json.Unmarshal(data, self))
	self.SiteName = name
	self.SiteDescription = desc
	self.Location = location
}
