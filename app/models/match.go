package models

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
	Score      int32

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
	)
}

func (m *Match) Create() error {
	// TODO
	return nil
}

func (m *Match) Read() error {
	// TODO
	return nil
}

func (m *Match) LoadRel() error {
	// TODO
	return nil
}

func (m *Match) Update() error {
	// TODO
	return nil
}
