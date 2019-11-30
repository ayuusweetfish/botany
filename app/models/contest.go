package models

type Contest struct {
	Id    int32
	Title string

	Owner     int32
	StartTime int64
	EndTime   int64

	IsVisible bool
	IsRegOpen bool

	Rel struct {
		Owner          User
		Participations []ContestParticipation
		MatchScripts   []ContestMatchScript
	}
}

const (
	ParticipationTypeModerator = iota
	ParticipationTypeContestant
)

type ContestParticipation struct {
	User    int32
	Contest int32
	Type    int8

	Rel struct {
		User    User
		Contest Contest
	}
}

type ContestMatchScript struct {
	Contest  int32
	Contents string
}

func init() {
	registerSchema("contest",
		"id SERIAL PRIMARY KEY",
		"title TEXT NOT NULL DEFAULT ''",
		"owner INTEGER NOT NULL REFERENCES users(id)",
		"start_time BIGINT NOT NULL",
		"end_time BIGINT NOT NULL",
		"is_visible BOOLEAN NOT NULL DEFAULT FALSE",
		"is_reg_open BOOLEAN NOT NULL DEFAULT FALSE",
	)
	registerSchema("contest_participation",
		"uid INTEGER NOT NULL REFERENCES users(id)",
		"contest INTEGER NOT NULL REFERENCES contest(id)",
		"type SMALLINT NOT NULL",
	)
	registerSchema("contest_match_script",
		"contest INTEGER NOT NULL REFERENCES contest(id)",
		"contents TEXT NOT NULL",
	)
}

func (u *Contest) Create() error {
	// TODO
	return nil
}

func (u *Contest) Read() error {
	// TODO
	return nil
}

func (u *Contest) LoadRel() error {
	// TODO
	return nil
}

func (u *Contest) Update() error {
	// TODO
	return nil
}
