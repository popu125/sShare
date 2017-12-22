package web

import (
	"github.com/popu125/sShare/web/captchas"
	"github.com/popu125/sShare/config"
)

type Captcha interface {
	Validate(token string) bool
}

func NewCaptcha(c config.CaptchaConf) Captcha {
	siteid, extra := c.SiteID, c.Extra
	switch c.Name {
	case "ppoi":
		return captchas.NewPpoiCaptcha(siteid, extra)
	case "coinhive":
		return captchas.NewCoinhiveCaptcha(siteid, extra)
	case "recaptcha":
		return captchas.NewReCaptcha(siteid, extra)
	default:
		return captchas.NewTestCaptcha(siteid)
	}
}
