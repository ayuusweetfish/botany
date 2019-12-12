package models

import (
	"database/sql"
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
		"uid INTEGER NOT NULL",
		"contest INTEGER NOT NULL",
		"created_at BIGINT NOT NULL",
		"status SMALLINT NOT NULL DEFAULT "+strconv.Itoa(SubmissionStatusPending),
		"message TEXT NOT NULL DEFAULT ''",
		"contents TEXT NOT NULL",
		"ADD CONSTRAINT fk_users FOREIGN KEY (uid) REFERENCES users (id)",
		"ADD CONSTRAINT fk_contest FOREIGN KEY (contest) REFERENCES contest (id)",
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

func SubmissionHistory(uid int32, cid int32, limit, offset int) ([]map[string]interface{}, int, error) {
	var rows *sql.Rows
	var err error
	if limit != 0 {
		// XXX: DRY?
		// All submissions
		rows, err = db.Query("SELECT "+
			"submission.id, submission.created_at, submission.status, "+
			"users.id, users.handle, users.privilege, users.nickname "+
			"FROM submission "+
			"LEFT JOIN users ON submission.uid = users.id "+
			"WHERE contest = $1 "+
			"ORDER BY submission.created_at DESC LIMIT $2 OFFSET $3",
			cid, limit, offset)
	} else {
		// Specific user
		rows, err = db.Query("SELECT "+
			"submission.id, submission.created_at, submission.status, "+
			"users.id, users.handle, users.privilege, users.nickname "+
			"FROM submission "+
			"LEFT JOIN users ON submission.uid = users.id "+
			"WHERE uid = $1 AND contest = $2 "+
			"ORDER BY submission.created_at DESC",
			uid, cid)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var ss []map[string]interface{}
	for rows.Next() {
		s := Submission{Contest: cid}
		err := rows.Scan(&s.Id, &s.CreatedAt, &s.Status,
			&s.Rel.User.Id, &s.Rel.User.Handle,
			&s.Rel.User.Privilege, &s.Rel.User.Nickname)
		if err != nil {
			return nil, 0, err
		}
		s.User = s.Rel.User.Id
		ss = append(ss, s.ShortRepresentation())
	}
	var total int
	rows2 := db.QueryRow("SELECT COUNT(*) FROM submission WHERE contest = $1", cid)
	err = rows2.Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return ss, total, rows.Err()
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

// Call after LoadRel()
func (s *Submission) IsVisibleTo(u User) bool {
	return s.User == u.Id ||
		s.Rel.Contest.HasEnded() ||
		s.Rel.Contest.ParticipationOf(u) == ParticipationTypeModerator
}
