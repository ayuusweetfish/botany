package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var rcli *redis.Client

func InitializeRedis(client *redis.Client) {
	rcli = client

	_, err := rcli.XGroupCreateMkStream("compile", "compile_group", "0").Result()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		panic(err)
	}

	go redisPollStatus()
}

func RedisSendForCompilation(s *Submission) error {
	_, err := rcli.XAdd(&redis.XAddArgs{
		Stream: "compile",
		ID:     "*",
		Values: map[string]interface{}{
			"sid":      s.Id,
			"contents": s.Contents,
		},
	}).Result()
	return err
}

func redisPollStatus() {
	for {
		println("Polling")
		r, err := rcli.BLPop(1*time.Second, "compile_result").Result()
		if err != nil && err.Error() != "redis: nil" {
			fmt.Println(err.Error())
			continue
		}
		// Assumes all data are well-formatted
		if r != nil {
			if r[0] == "compile_result" {
				sid, _ := strconv.ParseInt(r[1], 10, 32)
				r2, _ := rcli.LPop("compile_result").Result()
				status, _ := strconv.ParseInt(r2, 10, 8)
				r3, _ := rcli.LPop("compile_result").Result()
				if err := redisUpdateSubmissionStatus(int32(sid), int8(status), r3); err != nil {
					fmt.Println(err.Error())
					continue
				}
			}
		}
	}
}

func redisUpdateSubmissionStatus(sid int32, status int8, msg string) error {
	s := Submission{Id: sid}
	if err := s.Read(); err != nil {
		return err
	}

	s.Status = status
	s.Message = msg
	return s.Update()
}
