package controllers

import (
	_ "github.com/kawa-yoiko/botany/app/models"

	"golang.org/x/crypto/blake2b"

	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

const judgeTimestampThreshold = 30000

func judgeHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["ts"]) < 1 || len(r.URL.Query()["sig"]) < 1 {
		w.WriteHeader(400)
		return
	}

	ts := r.URL.Query()["ts"][0]
	sig := r.URL.Query()["sig"][0]
	tsVal, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	cur := time.Now().Unix()
	if cur < tsVal-judgeTimestampThreshold ||
		cur > tsVal+judgeTimestampThreshold {
		w.WriteHeader(403)
		return
	}

	h, err := blake2b.New256([]byte("aha"))
	if err != nil {
		panic(err)
	}
	h.Write([]byte(ts))
	digest := hex.EncodeToString(h.Sum(nil))
	println(digest) // For debug use

	if digest != sig {
		w.WriteHeader(403)
		return
	}
}

func init() {
	registerRouterFunc("/judge/{sid:[0-9]+}", judgeHandler, "GET")
}
