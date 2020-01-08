package models

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var rcli *redis.Client = nil
var redisPolling = false

func InitializeRedis(client *redis.Client) {
	rcli = client

	_, err := rcli.XGroupCreateMkStream("compile", "judge_group", "0").Result()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		panic(err)
	}

	_, err = rcli.XGroupCreateMkStream("match", "judge_group", "0").Result()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		panic(err)
	}

	// Prevent multiple polling goroutines
	if !redisPolling {
		redisPolling = true
		go redisPollStatus()
	}
}

func (s *Submission) SendToQueue() error {
	if rcli == nil {
		return nil
	}
	_, err := rcli.XAdd(&redis.XAddArgs{
		Stream: "compile",
		ID:     "*",
		Values: map[string]interface{}{"sid": s.Id},
	}).Result()
	return err
}

func (m *Match) SendToQueue(judge int32) error {
	if rcli == nil {
		return nil
	}
	values := map[string]interface{}{
		"mid":         m.Id,
		"judge":       judge,
		"num_parties": len(m.Rel.Parties),
	}
	for i, p := range m.Rel.Parties {
		values["party_"+strconv.Itoa(i)] = p.Id
	}
	_, err := rcli.XAdd(&redis.XAddArgs{
		Stream: "match",
		ID:     "*",
		Values: values,
	}).Result()
	return err
}

func redisPollStatus() {
	for {
		// println("Polling")
		r, err := rcli.BLPop(1*time.Second, "compile_result", "match_result").Result()
		if err != nil && err.Error() != "redis: nil" {
			log.Println(err.Error())
			continue
		}
		// Assumes all data are well-formatted
		if r != nil {
			if r[0] == "compile_result" || r[0] == "match_result" {
				id, _ := strconv.ParseInt(r[1], 10, 32)
				r2, _ := rcli.LPop(r[0]).Result()
				status, _ := strconv.ParseInt(r2, 10, 8)
				r3, _ := rcli.LPop(r[0]).Result()
				var err error
				if r[0] == "compile_result" {
					err = redisUpdateSubmissionStatus(int32(id), int8(status), r3)
				} else {
					err = redisUpdateMatchStatus(int32(id), int8(status), r3)
				}
				if err != nil {
					for i := 0; i < len(r); i++ {
						println(r[i])
					}
					log.Println(err.Error())
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

func redisUpdateMatchStatus(mid int32, status int8, msg string) error {
	m := Match{Id: mid}
	if err := m.Read(); err != nil {
		return err
	}

	m.Status = status
	m.Report = msg
	if err := m.Update(); err != nil {
		return err
	}

	if status == MatchStatusDone {
		// Save participant logs
		num, err := m.PartiesCount()
		if err != nil {
			return err
		}
		println("Number of parties: ", num)
		for i := 0; i < num; i++ {
			log, _ := rcli.LPop("match_result").Result()
			println(i, log)
			p := MatchParty{Match: mid, Index: int32(i), Log: log}
			if err := p.UpdateLog(); err != nil {
				return err
			}
		}
	}

	return nil
}
