package captchas

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type rv struct {
	siteKey   string
	secretKey string
}

func (self rv) Validate(token string) bool {
	v := url.Values{}
	v.Add("response", token)
	v.Add("secret", self.secretKey)

	r, err := http.PostForm("https://recaptcha.net/recaptcha/api/siteverify", v)
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

func NewReCaptcha(siteKey string, secretKey string) Captcha {
	return rv{siteKey: siteKey, secretKey: secretKey}
}
