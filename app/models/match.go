package models

import (
	"strconv"
)

const (
	MatchStatusPending     = 0
	MatchStatusRunning     = 1
	MatchStatusDone        = 9
	MatchStatusSystemError = -9
)

type Match struct {
	Id      int32
	Contest int32
	Status  int8
	Report  string

	Rel struct {
		Contest Contest
		Parties []Submission
	}
}

type MatchParty struct {
	Match      int32
	Index      int32
	Submission int32
	Log        string

	Rel struct {
		Match      Match
		Submission Submission
	}
}

func init() {
	registerSchema("match",
		"id SERIAL PRIMARY KEY",
		"contest INTEGER NOT NULL",
		"status SMALLINT NOT NULL DEFAULT "+strconv.Itoa(MatchStatusPending),
		"report TEXT NOT NULL DEFAULT ''",
		"ADD CONSTRAINT fk_contest FOREIGN KEY (contest) REFERENCES contest (id)",
	)
	registerSchema("match_party",
		"match INTEGER NOT NULL",
		"index INTEGER NOT NULL",
		"submission INTEGER NOT NULL",
		"log TEXT NOT NULL DEFAULT ''",
		"ADD CONSTRAINT fk_match FOREIGN KEY (match) REFERENCES match (id)",
		"ADD CONSTRAINT fk_submission FOREIGN KEY (submission) REFERENCES submission (id)",
	)
}

func (m *Match) Create() error {
	// TODO: Combine into an transaction
	err := db.QueryRow("INSERT INTO "+
		"match(contest, report) "+
		"VALUES ($1, $2) RETURNING id",
		m.Contest,
		m.Report,
	).Scan(&m.Id)
	if err != nil {
		return err
	}

	// Create MatchParty records
	for i, s := range m.Rel.Parties {
		_, err := db.Exec("INSERT INTO "+
			"match_party(match, index, submission) VALUES ($1, $2, $3)",
			m.Id, i, s.Id)
		if err != nil {
			return err
		}
		err = s.Read()
		if err != nil {
			return err
		}
		err = s.LoadRel()
		if err != nil {
			return err
		}
		m.Rel.Parties[i] = s
	}

	return nil
}

func (m *Match) ShortRepresentation() map[string]interface{} {
	parties := []map[string]interface{}{}
	for _, p := range m.Rel.Parties {
		parties = append(parties, p.ShortRepresentation())
	}
	return map[string]interface{}{
		"id":      m.Id,
		"parties": parties,
		"status":  m.Status,
		"contest": map[string]interface{}{
			"id":    m.Rel.Contest.Id,
			"title": m.Rel.Contest.Title,
		},
	}
}

func (m *Match) Representation() map[string]interface{} {
	r := m.ShortRepresentation()
	r["report"] = m.Report
	return r
}

func (m *Match) Read() error {
	err := db.QueryRow("SELECT "+
		"contest, status, report "+
		"FROM match WHERE id = $1", m.Id,
	).Scan(
		&m.Contest,
		&m.Status,
		&m.Report,
	)
	return err
}

func ReadByContest(cid int32) ([]Match, error) {
	rows, err := db.Query("SELECT DISTINCT match.id, match.status, match.report, "+
		"c.id, c.title, c.banner, c.owner, c.start_time, c.end_time, c.descr, c.details, c.is_visible, c.is_reg_open, c.script, "+
		"submission.id, submission.uid, submission.contest, submission.created_at, submission.status, submission.message, submission.language, submission.contents, "+
		"u.id, u.handle, u.email, u.password, u.privilege, u.joined_at, u.nickname, u.bio "+
		"FROM match "+
		"JOIN contest c on match.contest = c.id "+
		"JOIN match_party ON match.id = match_party.match "+
		"JOIN submission ON match_party.submission = submission.id "+
		"JOIN users u on submission.uid = u.id "+
		"WHERE match.contest = $1 ORDER BY match.id DESC", cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ms := []Match{}
	for rows.Next() {
		m := Match{Contest: cid}
		s := Submission{}
		u := User{}
		err := rows.Scan(&m.Id, &m.Status, &m.Report, &m.Rel.Contest.Id, &m.Rel.Contest.Title, &m.Rel.Contest.Banner, &m.Rel.Contest.Owner, &m.Rel.Contest.StartTime,
			&m.Rel.Contest.EndTime, &m.Rel.Contest.Desc, &m.Rel.Contest.Details, &m.Rel.Contest.IsVisible, &m.Rel.Contest.IsRegOpen, &m.Rel.Contest.Script,
			&s.Id, &s.User, &s.Contest, &s.CreatedAt, &s.Status, &s.Message, &s.Language, &s.Contents,
			&u.Id, &u.Handle, &u.Email, &u.Password, &u.Privilege, &u.JoinedAt, &u.Nickname, &u.Bio)
		if err != nil {
			return nil, err
		}
		if len(ms) == 0 {
			s.Rel.User = u
			s.Rel.Contest = m.Rel.Contest
			m.Rel.Parties = append(m.Rel.Parties, s)
			ms = append(ms, m)
		} else {
			if m.Id == ms[len(ms)-1].Id {
				s.Rel.User = u
				s.Rel.Contest = m.Rel.Contest
				ms[len(ms)-1].Rel.Parties = append(ms[len(ms)-1].Rel.Parties, s)
			} else {
				s.Rel.User = u
				s.Rel.Contest = m.Rel.Contest
				m.Rel.Parties = append(m.Rel.Parties, s)
				ms = append(ms, m)
			}
		}
	}
	return ms, rows.Err()
}

func (m *Match) PartiesCount() (int, error) {
	var num int
	err := db.QueryRow("SELECT COUNT(*) FROM match_party WHERE match = $1",
		m.Id).Scan(&num)
	return num, err
}

func (m *Match) LoadRel() error {
	m.Rel.Contest.Id = m.Contest
	if err := m.Rel.Contest.Read(); err != nil {
		return err
	}

	// Find out all parties
	rows, err := db.Query("SELECT submission FROM match_party "+
		"WHERE match = $1 ORDER BY index ASC", m.Id)
	if err != nil {
		return err
	}
	defer rows.Close()
	m.Rel.Parties = []Submission{}
	for rows.Next() {
		// TODO: Optimize
		s := Submission{}
		if err := rows.Scan(&s.Id); err != nil {
			return err
		}
		if err := s.Read(); err != nil {
			return err
		}
		if err := s.LoadRel(); err != nil {
			return err
		}
		m.Rel.Parties = append(m.Rel.Parties, s)
	}
	return rows.Err()
}

func (m *Match) Update() error {
	_, err := db.Exec("UPDATE match SET "+
		"status = $1, report = $2 WHERE id = $3",
		m.Status, m.Report, m.Id)
	return err
}

func (p *MatchParty) LoadLog() error {
	return db.QueryRow("SELECT log FROM match_party WHERE match = $1 AND index = $2",
		p.Match, p.Index).Scan(&p.Log)
}

func (p *MatchParty) UpdateLog() error {
	_, err := db.Exec("UPDATE match_party SET log = $1 WHERE match = $2 AND index = $3",
		p.Log, p.Match, p.Index)
	return err
}
