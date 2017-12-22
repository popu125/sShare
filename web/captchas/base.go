package captchas

type Captcha interface {
	Validate(token string) bool
}

type TestCaptcha struct {
	siteID string
}

func (self *TestCaptcha) Validate(code string) bool {
	return true
}

func NewTestCaptcha(siteid string) Captcha {
	return &TestCaptcha{siteid}
}
