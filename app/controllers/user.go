package controllers

import (
	"github.com/kawa-yoiko/botany/app/globals"
	"github.com/kawa-yoiko/botany/app/models"

	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	id, err := r.Cookie("QAQ")
	if err == http.ErrNoCookie {
		//w.WriteHeader()
		//？遗留问题暂待处理
	}
	if r.Method == "POST" {
		r.ParseForm()
		u := models.User{
			Name:     strings.Join(r.Form["username"], ""),
			Password: strings.Join(r.Form["password"], ""),
		}
		captcha := strings.Join(r.Form["captcha"], "")
		captchaPass := globals.CaptchaVerfiy(id.Value, captcha)
		if captchaPass {
			err := u.Read()
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "该用户不存在"}`))
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal Server Error"}`))
				return
			} else {
				passwordPass := u.VerifyPassword()
				if passwordPass {
					w.Write([]byte(`{"success": "登录成功"}`))
				} else {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"error": "密码错误"}`))
				}
			}
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	middlewareProcessSession(w, r)
	id, err := r.Cookie("QAQ")
	if err == http.ErrNoCookie {
		//w.WriteHeader()
		//？遗留问题暂待处理
	}
	if r.Method == "POST" {
		r.ParseForm()
		u := models.User{
			Name:     strings.Join(r.Form["username"], ""),
			Password: strings.Join(r.Form["password"], ""),
			Email:    strings.Join(r.Form["email"], ""),
		}
		captcha := strings.Join(r.Form["captcha"], "")
		captchaPass := globals.CaptchaVerfiy(id.Value, captcha)
		if captchaPass {
			var uid int
			err := u.Read()
			if err == sql.ErrNoRows {
				u.Create()
				w.Write([]byte(fmt.Sprintf(`{"success: 注册成功", uid: %d}`, uid)))
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Internal Server Error"}`))
				return
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "用户名已存在"}`))
			}
		}
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
		_, captcha := globals.CaptchaCreate(id.Value)
		result := []byte(fmt.Sprintf(`{"pic": "%s"}`, captcha))
		w.Write(result)
		return
	}
}
