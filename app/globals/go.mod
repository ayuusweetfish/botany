module globals

go 1.13

require (
	github.com/gorilla/sessions v1.2.0
	github.com/mojocn/base64Captcha v0.0.0
	github.com/kawa-yoiko/botany/app/models v0.0.0-incompatible
)

replace github.com/mojocn/base64Captcha => /home/sakura/go/src/github.com/mojocn/base64Captcha

replace github.com/kawa-yoiko/botany/app/models => ../models
