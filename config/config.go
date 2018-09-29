package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func LoadConfig(fn string) *Config {
	data, err := ioutil.ReadFile(fn)
	check(err)
	conf := Config{}
	err2 := json.Unmarshal(data, &conf)
	check(err2)
	conf.TTL, err = time.ParseDuration(conf.TTLString)
	check(err)
	if !conf.RunCmd.Enabled {
		fmt.Println("Run Command should be enabled.")
		os.Exit(1)
	}
	if conf.PortRange < conf.Limit || conf.PortStart+conf.PortRange > 65536 {
		fmt.Println("Port range settings illegal, please check your config.")
		os.Exit(1)
	}
	return &conf
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Config struct {
	RunCmd  CmdConf `json:"run_command"`
	ExitCmd CmdConf `json:"exit_command"`

	Captcha   CaptchaConf `json:"captcha"`
	TTLString string      `json:"ttl"`
	Limit     uint32      `json:"limit"`
	Addr      string      `json:"web_addr"`
	RandSeed  int64       `json:"rand_seed"`
	PortStart uint32      `json:"port_start"`
	PortRange uint32      `json:"port_range"`

	Safe         SafeConf `json:"safe"`
	NoCheckAlive bool     `json:"no_check_alive"`

	TTL time.Duration
}

type CaptchaConf struct {
	Name   string `json:"name"`
	SiteID string `json:"site_id"`
	Extra  string `json:"extra"`
}

type CmdConf struct {
	Cmd     string `json:"cmd"`
	Arg     string `json:"arg"`
	Enabled bool   `json:"enabled"`
}

type SafeConf struct {
	AntiCC     bool   `json:"anti_cc"`
	CityCheck  bool   `json:"city_check"`
	CityName   string `json:"city_name"`
	CityFile   string `json:"city_file"`
	CDNEnabled bool   `json:"cdn_enabled"`
}
