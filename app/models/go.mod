module models

go 1.13

require (
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/yuin/gopher-lua v0.0.0-20191213034115-f46add6fdb5c
	golang.org/x/crypto v0.0.0 // indirect
	golang.org/x/crypto/bcrypt v0.0.0
)

replace golang.org/x/crypto => /home/sakura/go/pkg/mod/golang.org/x/crypto@v0.0.0

replace golang.org/x/crypto/bcrypt => /home/sakura/go/pkg/mod/golang.org/x/crypto@v0.0.0/bcrypt
