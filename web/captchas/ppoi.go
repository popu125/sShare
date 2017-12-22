package captchas

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

type p struct {
	siteID string
	hashes string
}

func (self p) Validate(token string) bool {
	v := url.Values{}
	v.Add("secret", self.siteID)
	v.Add("token", token)
	v.Add("hashes", self.hashes)

	r, err := http.PostForm("https://api.ppoi.org/token/verify", v)
	if err != nil {
		return false
	}

	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false
	}
	chR := struct {
		Success bool `json:"success"`
	}{}
	json.Unmarshal(resp, &chR)
	if chR.Success {
		return true
	} else {
		return false
	}
}

func NewPpoiCaptcha(siteid string, hashes string) Captcha {
	return &p{siteid, hashes}
}
