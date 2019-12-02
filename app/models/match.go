package models

type Match struct {
	Id      int32
	Contest int32
	Report  string

	Rel struct {
		Contest Contest
		Sides   []Submission
	}
}

type MatchSide struct {
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
	registerSchema("match_side",
		"match INTEGER NOT NULL REFERENCES match(id)",
		"submission INTEGER NOT NULL REFERENCES submission(id)",
		"score INTEGER NOT NULL DEFAULT 0",
	)
}

func (u *Match) Create() error {
	// TODO
	return nil
}

func (u *Match) Read() error {
	// TODO
	return nil
}

func (u *Match) LoadRel() error {
	// TODO
	return nil
}

func (u *Match) Update() error {
	// TODO
	return nil
}
