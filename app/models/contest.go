package models

import "strconv"

type Contest struct {
	Id     int32
	Title  string
	Banner string

	Owner     int32
	StartTime int64
	EndTime   int64

	Desc    string
	Details string

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

const (
	MatchScriptHookManual = iota
	MatchScriptHookSubmission
	MatchScriptHookTimed
)

type ContestMatchScript struct {
	Contest  int32
	Hook     int8
	Contents string
}

func init() {
	registerSchema("contest",
		"id SERIAL PRIMARY KEY",
		"title TEXT NOT NULL DEFAULT ''",
		"banner TEXT NOT NULL DEFAULT ''",
		"owner INTEGER NOT NULL REFERENCES users(id)",
		"start_time BIGINT NOT NULL",
		"end_time BIGINT NOT NULL",
		"desc TEXT NOT NULL DEFAULT ''",
		"details TEXT NOT NULL DEFAULT ''",
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
		"hook SMALLINT NOT NULL DEFAULT "+strconv.Itoa(MatchScriptHookManual),
		"contents TEXT NOT NULL",
	)
}

func (c *Contest) Create() error {
	// TODO
	return nil
}

func (c *Contest) Read() error {
	// TODO
	return nil
}

func (c *Contest) LoadRel() error {
	// TODO
	return nil
}

func (c *Contest) Update() error {
	// TODO
	return nil
}
