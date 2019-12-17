package controllers

import (
	"github.com/kawa-yoiko/botany/app/models"

	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/blake2b"
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

	if digest != sig {
		w.WriteHeader(403)
		return
	}

	sid, _ := strconv.Atoi(mux.Vars(r)["sid"])
	s := models.Submission{Id: int32(sid)}
	if err := s.Read(); err != nil {
		panic(err)
	}
	w.Write([]byte(strconv.Itoa(len(s.Contents))))
	w.Write([]byte(" "))
	w.Write([]byte(s.Language))
	w.Write([]byte("\n"))
	w.Write([]byte(s.Contents))
}

func init() {
	registerRouterFunc("/judge/{sid:[0-9]+}", judgeHandler, "GET")
}
