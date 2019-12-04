package models

import (
	"strconv"
	"time"
)

const (
	SubmissionStatusPending           = 0
	SubmissionStatusCompiling         = 1
	SubmissionStatusAccepted          = 9
	SubmissionStatusCompilationFailed = -1
	SubmissionStatusSystemError       = -9
)

type Submission struct {
	Id        int32
	User      int32
	Contest   int32
	CreatedAt int64
	Status    int8
	Message   string
	Contents  string

	Rel struct {
		User    User
		Contest Contest
	}
}

func init() {
	registerSchema("submission",
		"id SERIAL PRIMARY KEY",
		"uid INTEGER NOT NULL REFERENCES users(id)",
		"contest INTEGER NOT NULL REFERENCES contest(id)",
		"created_at BIGINT NOT NULL",
		"status SMALLINT NOT NULL DEFAULT "+strconv.Itoa(SubmissionStatusPending),
		"message TEXT NOT NULL DEFAULT ''",
		"contents TEXT NOT NULL",
	)
}

func (s *Submission) Representation() map[string]interface{} {
	return map[string]interface{}{
		"id":          s.Id,
		"participant": s.Rel.User.ShortRepresentation(),
		"created_at":  s.CreatedAt,
		"status":      s.Status,
		"msg":         s.Message,
		"contents":    s.Contents,
	}
}

func (s *Submission) ShortRepresentation() map[string]interface{} {
	return map[string]interface{}{
		"id":          s.Id,
		"participant": s.Rel.User.ShortRepresentation(),
		"created_at":  s.CreatedAt,
		"status":      s.Status,
	}
}

func (s *Submission) Create() error {
	s.CreatedAt = time.Now().Unix()
	s.Status = SubmissionStatusPending
	s.Message = ""
	err := db.QueryRow("INSERT INTO "+
		"submission(uid, contest, created_at, contents) "+
		"VALUES ($1, $2, $3, $4) RETURNING id",
		s.User,
		s.Contest,
		s.CreatedAt,
		s.Contents,
	).Scan(&s.Id)
	return err
}

func (s *Submission) Read() error {
	err := db.QueryRow("SELECT "+
		"uid, contest, created_at, status, message, contents "+
		"FROM submission WHERE id = $1",
		s.Id,
	).Scan(&s.User, &s.Contest, &s.CreatedAt,
		&s.Status, &s.Message, &s.Contents)
	return err
}

func SubmissionHistory(uid int32, cid int32, limit int) ([]Submission, error) {
	rows, err := db.Query("SELECT "+
		"submission.id, submission.created_at, submission.status, "+
		"users.id, users.handle, users.privilege, users.nickname "+
		"FROM submission "+
		"LEFT JOIN users ON submission.uid = users.id "+
		"WHERE uid = $1 AND contest = $2 "+
		"ORDER BY submission.created_at DESC LIMIT $3",
		uid, cid, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ss := []Submission{}
	for rows.Next() {
		s := Submission{User: uid, Contest: cid}
		err := rows.Scan(&s.Id, &s.CreatedAt, &s.Status,
			&s.Rel.User.Id, &s.Rel.User.Handle,
			&s.Rel.User.Privilege, &s.Rel.User.Nickname)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return ss, rows.Err()
}

func (s *Submission) LoadRel() error {
	s.Rel.User.Id = s.User
	if err := s.Rel.User.ReadById(); err != nil {
		return err
	}
	s.Rel.Contest.Id = s.Contest
	return s.Rel.Contest.Read()
}

func (s *Submission) Update() error {
	// TODO
	return nil
}
