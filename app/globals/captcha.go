package globals

import (
	"github.com/mojocn/base64Captcha"
	captchaStore "github.com/mojocn/base64Captcha/store"

	"time"
)

var captchaConfig = base64Captcha.ConfigDigit{
	Height:     80,
	Width:      240,
	MaxSkew:    0.7,
	DotCount:   80,
	CaptchaLen: 4,
}

func CaptchaCreate() (string, string) {
	idKey, captcha := base64Captcha.GenerateCaptcha("", captchaConfig)
	base64string := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
	return idKey, base64string
}

func CaptchaVerfiy(idKey, verifyValue string) bool {
	return base64Captcha.VerifyCaptcha(idKey, verifyValue)
}

func init() {
	base64Captcha.SetCustomStore(captchaStore.NewMemoryStore(65536, 5*time.Minute))
}
