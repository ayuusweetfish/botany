package models

import (
	"strconv"
)

const (
	SubmissionStatusPending           = 0
	SubmissionStatusCompiling         = 1
	SubmissionStatusAccepted          = 9
	SubmissionStatusCompilationFailed = -1
	SubmissionStatusSystemError       = -9
)

type Submission struct {
	Id            int32
	Participation int32
	CreatedAt     int64
	Status        int8
	Message       string
	Contents      string

	Rel struct {
		Participation ContestParticipation
	}
}

func init() {
	registerSchema("submission",
		"id SERIAL PRIMARY KEY",
		"participation INTEGER NOT NULL REFERENCES contest_participation(id)",
		"created_at BIGINT NOT NULL",
		"status SMALLINT NOT NULL DEFAULT "+strconv.Itoa(SubmissionStatusPending),
		"message TEXT NOT NULL",
		"contents TEXT NOT NULL",
	)
}

func (s *Submission) Create() error {
	// TODO
	return nil
}

func (s *Submission) Read() error {
	// TODO
	return nil
}

func (s *Submission) LoadRel() error {
	// TODO
	return nil
}

func (s *Submission) Update() error {
	// TODO
	return nil
}
