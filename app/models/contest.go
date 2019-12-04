package models

import (
	"database/sql"
	"strconv"
)

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
	Id       int32
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
		"descr TEXT NOT NULL DEFAULT ''",
		"details TEXT NOT NULL DEFAULT ''",
		"is_visible BOOLEAN NOT NULL DEFAULT FALSE",
		"is_reg_open BOOLEAN NOT NULL DEFAULT FALSE",
	)
	registerSchema("contest_participation",
		"uid INTEGER NOT NULL REFERENCES users(id)",
		"contest INTEGER NOT NULL REFERENCES contest(id)",
		"type SMALLINT NOT NULL",
		"ADD PRIMARY KEY (uid, contest)",
	)
	registerSchema("contest_match_script",
		"id SERIAL PRIMARY KEY",
		"contest INTEGER NOT NULL REFERENCES contest(id)",
		"hook SMALLINT NOT NULL DEFAULT "+strconv.Itoa(MatchScriptHookManual),
		"contents TEXT NOT NULL",
	)
}

func (c *Contest) Representation(uid int32) map[string]interface{} {
	myRole := int8(-1)
	if uid == c.Owner {
		myRole = ParticipationTypeModerator
	} else if uid != -1 {
		p := ContestParticipation{
			User:    uid,
			Contest: c.Id,
		}
		if err := p.Read(); err != nil {
			if err != sql.ErrNoRows {
				panic(err)
			}
			// No participation, keep myRole as -1
		} else {
			myRole = p.Type
		}
	}

	return map[string]interface{}{
		"id":          c.Id,
		"title":       c.Title,
		"banner":      c.Banner,
		"start_time":  c.StartTime,
		"end_time":    c.EndTime,
		"desc":        c.Desc,
		"details":     c.Details,
		"is_reg_open": c.IsRegOpen,
		"owner":       c.Rel.Owner.ShortRepresentation(),
		"my_role":     myRole,
	}
}

func (c *Contest) ShortRepresentation() map[string]interface{} {
	return map[string]interface{}{
		"id":          c.Id,
		"title":       c.Title,
		"banner":      c.Banner,
		"start_time":  c.StartTime,
		"end_time":    c.EndTime,
		"desc":        c.Desc,
		"is_reg_open": c.IsRegOpen,
	}
}

func (c *Contest) Create() error {
	err := db.QueryRow("INSERT INTO "+
		"contest(title, banner, owner, start_time, end_time, descr, details, is_visible, is_reg_open) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		c.Title,
		c.Banner,
		c.Owner,
		c.StartTime,
		c.EndTime,
		c.Desc,
		c.Details,
		c.IsVisible,
		c.IsRegOpen,
	).Scan(&c.Id)
	return err
}

func (c *Contest) Read() error {
	err := db.QueryRow("SELECT "+
		"title, banner, owner, start_time, end_time, descr, details, is_visible, is_reg_open "+
		"FROM contest WHERE id = $1",
		c.Id,
	).Scan(
		&c.Title,
		&c.Banner,
		&c.Owner,
		&c.StartTime,
		&c.EndTime,
		&c.Desc,
		&c.Details,
		&c.IsVisible,
		&c.IsRegOpen,
	)
	return err
}

func ContestReadAll() ([]Contest, error) {
	rows, err := db.Query("SELECT " +
		"id, title, banner, owner, start_time, end_time, descr, is_visible, is_reg_open " +
		"FROM contest",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cs := []Contest{}
	for rows.Next() {
		c := Contest{}
		err := rows.Scan(
			&c.Id,
			&c.Title,
			&c.Banner,
			&c.Owner,
			&c.StartTime,
			&c.EndTime,
			&c.Desc,
			&c.IsVisible,
			&c.IsRegOpen,
		)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, rows.Err()
}

func (c *Contest) LoadRel() error {
	c.Rel.Owner.Id = c.Owner
	return c.Rel.Owner.ReadById()
}

func (c *Contest) Update() error {
	// TODO
	return nil
}

func (p *ContestParticipation) Create() error {
	_, err := db.Exec("INSERT INTO "+
		"contest_participation(uid, contest, type) "+
		"VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		p.User,
		p.Contest,
		p.Type,
	)
	return err
}

func (p *ContestParticipation) Read() error {
	err := db.QueryRow("SELECT type "+
		"FROM contest_participation WHERE uid = $1 AND contest = $2",
		p.User,
		p.Contest,
	).Scan(&p.Type)
	return err
}

func (p *ContestParticipation) Update() error {
	_, err := db.Exec("UPDATE contest_participation SET "+
		"type = $1 WHERE uid = $1 AND contest = $2",
		p.Type,
		p.User,
		p.Contest,
	)
	return err
}
