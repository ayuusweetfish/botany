package models

import "fmt"

type Match struct {
	Id      int32
	Contest int32
	Report  string

	Rel struct {
		Contest Contest
		Parties []Submission
	}
}

type MatchParty struct {
	Match      int32
	Submission int32

	Rel struct {
		Match      Match
		Submission Submission
	}
}

func init() {
	registerSchema("match",
		"id SERIAL PRIMARY KEY",
		"contest INTEGER NOT NULL REFERENCES contest(id)",
		"report TEXT NOT NULL DEFAULT ''",
	)
	registerSchema("match_party",
		"match INTEGER NOT NULL REFERENCES match(id)",
		"submission INTEGER NOT NULL REFERENCES submission(id)",
		"score INTEGER NOT NULL DEFAULT 0",
		"is_winner BOOLEAN NOT NULL",
	)
}

func (m *Match) Create() error {
	err := db.QueryRow("INSERT INTO "+
		"match(contest, report) "+
		"VALUES ($1, $2) RETURNING id",
		m.Contest,
		m.Report,
	).Scan(&m.Id)
	return err
}

func (m *Match) ShortRepresentation() map[string]interface{} {
	return map[string]interface{}{
		"id":      m.Id,
		"contest": m.Contest,
		"parties": m.Rel,
	}
}

func (m *Match) Read() error {
	err := db.QueryRow("SELECT "+
		"contest, report "+
		"FROM match WHERE id = $1", m.Id,
	).Scan(
		&m.Contest,
		&m.Report,
	)
	return err
}

func ReadByContest(contest int32) ([]Match, error) {
	rows, err := db.Query("SELECT "+
		"match.id, s.id, s.uid, "+
		"s.created_at, s.status "+
		"FROM match "+
		"LEFT JOIN match_party ON match.id = match_party.match "+
		"LEFT JOIN submission s ON match_party.submission = s.id "+
		"WHERE match.contest = $1 ",
		contest)
	if err != nil {
		fmt.Println("bad", err)
		return nil, err
	}
	defer rows.Close()
	ms := []Match{}
	s := Submission{}
	m := Match{Contest: contest}
	var matchId int32 = -1
	matchEnd := 0 // true means one match finished.
	for rows.Next() {
		err = rows.Scan(&m.Id, &s.Id, &s.User, &s.CreatedAt, &s.Status)
		if err != nil {
			return nil, err
		}
		if matchId == -1 {
			matchId = m.Id
		} else if matchId != m.Id {
			m.Id, matchId = matchId, m.Id
			ms = append(ms, m)
			m.Id = matchId
			matchEnd = 1
		}
		m.Rel.Parties = append(m.Rel.Parties, s)
		matchEnd = 0
	}
	if matchEnd == 0 {
		ms = append(ms, m)
	}
	return ms, nil
}

func (m *Match) LoadRel() error {
	m.Rel.Contest.Id = m.Contest
	return m.Rel.Contest.Read()
}

func (m *Match) Update() error {
	// TODO
	return nil
}
