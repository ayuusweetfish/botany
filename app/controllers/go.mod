module controllers

go 1.13

require (
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/kawa-yoiko/botany/app/globals v0.0.0-incompatible
	github.com/kawa-yoiko/botany/app/models v0.0.0-incompatible
)

replace github.com/kawa-yoiko/botany/app/models => ../models

replace github.com/kawa-yoiko/botany/app/globals => ../globals
