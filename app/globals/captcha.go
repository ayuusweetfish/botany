package globals

import (
	"github.com/mojocn/base64Captcha"

	"time"
)

var store = base64Captcha.NewMemoryStore(65536, 5*time.Minute)

var captchaDriver = base64Captcha.DriverDigit{
	Height:   80,
	Width:    240,
	Length:   4,
	MaxSkew:  0.7,
	DotCount: 80,
}

func CaptchaCreate() (string, string) {
	captcha := base64Captcha.NewCaptcha(&(captchaDriver), store)
	idKey, base64string, err := captcha.Generate()
	if err != nil {
		panic(err)
	}
	return idKey, base64string
}

func CaptchaVerfiy(idKey, verifyValue string) bool {
	return store.Verify(idKey, verifyValue, true)
}
