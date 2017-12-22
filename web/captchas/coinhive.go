package captchas

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

type ch struct {
	siteID string
	hashes string
}

func (self ch) Validate(token string) bool {
	v := url.Values{}
	v.Add("secret", self.siteID)
	v.Add("token", token)
	v.Add("hashes", self.hashes)

	r, err := http.PostForm("https://api.coinhive.com/token/verify", v)
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

func NewCoinhiveCaptcha(siteid string, hashes string) Captcha {
	return &ch{siteid, hashes}
}
