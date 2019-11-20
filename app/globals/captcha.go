package globals

import (
	"github.com/mojocn/base64Captcha"
)

var config = base64Captcha.ConfigDigit{
	Height:     80,
	Width:      240,
	MaxSkew:    0.7,
	DotCount:   80,
	CaptchaLen: 4,
}

func CaptchaCreate(idKey string) (id string, base64 string) {
	idKey, captcha := base64Captcha.GenerateCaptcha(idKey, config)
	base64string := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
	return idKey, base64string
}

func CaptchaVerfiy(idkey, verifyValue string) bool {
	verifyResult := base64Captcha.VerifyCaptcha(idkey, verifyValue)
	if verifyResult {
		return true
	} else {
		return false
	}
}
