package main

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var config = base64Captcha.ConfigDigit{
	Height:     80,
	Width:      240,
	MaxSkew:    0.7,
	DotCount:   80,
	CaptchaLen: 4,
}

func captchaCreate(idKey string) (id string, base64 string) {
	idKey, captcha := base64Captcha.GenerateCaptcha(idKey, config)
	base64string := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
	return idKey, base64string
}

func captchaVerfiy(idkey, verifyValue string) bool {
	verifyResult := base64Captcha.VerifyCaptcha(idkey, verifyValue)
	if verifyResult {
		return true
	} else {
		return false
	}
}

func captchaHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	if r.Method == "GET" {
		id, err := r.Cookie("QAQ")

		if err == http.ErrNoCookie {
			//w.WriteHeader()
			//？遗留问题暂待处理
		}
		_, captcha := captchaCreate(id.Value)
		result := []byte(fmt.Sprintf(`{"pic": "%s"}`, captcha))
		w.Write(result)
		return
	}
}
