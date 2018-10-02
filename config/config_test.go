package config

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	testConfigFilename = "config_test.json"
	trueConfig         = []byte(`{
  "run_command": {
    "cmd": "client",
    "arg": "-p {{.Port}} -k {{.Pass \"ss_pass\"}}",
    "enabled": true
  },
  "exit_command": {
    "cmd": "client",
    "arg": "kill -p {{.Port}} -k {{.Pass \"ss_pass\"}}",
    "enabled": false
  },
  "captcha": {
    "name": "base",
    "site_id": "23333",
    "extra": "66666"
  },
  "ttl": "20m",
  "limit": 20,
  "web_addr": ":9527",
  "port_start": 2000,
  "port_range": 200,
  "rand_seed": 23343,
  "no_check_alive": false,
  "safe": {
    "anti_cc": false,
    "city_check": false,
    "city_name": "Beijing",
    "city_file":"/path/to/your/file",
    "cdn_enabled": false
  }
}`)
)

func TestLoadConfig(t *testing.T) {
	ioutil.WriteFile(testConfigFilename, trueConfig, os.ModePerm)
	defer os.Remove(testConfigFilename)
	LoadConfig(testConfigFilename)
}
