package main

import (
	"github.com/mojocn/base64Captcha"
)

var config = base64Captcha.ConfigDigit{
	Height:     80,
	Width:      240,
	MaxSkew:    0.7,
	DotCount:   80,
	CaptchaLen: 5,
}

func captchaCreate() (id string, base64 string) {
	idKey, captcha := base64Captcha.GenerateCaptcha("", config)
	base64string := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
	return idKey, base64string
}

func captchaUpdate(id string) (basee64 string) {
	_, captcha := base64Captcha.GenerateCaptcha(id, config)
	base64string := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
	return base64string
}

func verfiyCaptcha(idkey, verifyValue string) bool {
	verifyResult := base64Captcha.VerifyCaptcha(idkey, verifyValue)
	if verifyResult {
		return true
	} else {
		return false
	}
}
