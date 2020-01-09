module main

go 1.13

require (
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gorilla/sessions v1.2.0
	github.com/kawa-yoiko/botany/app/controllers v0.0.0-incompatible
	github.com/kawa-yoiko/botany/app/globals v0.0.0-incompatible
	github.com/kawa-yoiko/botany/app/models v0.0.0-incompatible
	github.com/lib/pq v1.2.0
	golang.org/x/crypto/blowfish v0.0.0 // indirect
)

replace github.com/kawa-yoiko/botany/app/controllers => ./controllers

replace github.com/kawa-yoiko/botany/app/models => ./models

replace github.com/mojocn/base64Captcha => /home/sakura/go/src/github.com/mojocn/base64Captcha

replace golang.org/x/crypto => /home/sakura/go/pkg/mod/golang.org/x/crypto@v0.0.0

replace golang.org/x/crypto/bcrypt => /home/sakura/go/pkg/mod/golang.org/x/crypto@v0.0.0/bcrypt

replace golang.org/x/crypto/blowfish => /home/sakura/go/pkg/mod/golang.org/x/crypto@v0.0.0/blowfish

replace github.com/kawa-yoiko/botany/app/globals => ./globals
