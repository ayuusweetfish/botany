package models

import (
	"database/sql"
	"strconv"
	"time"
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
	User        int32
	Contest     int32
	Type        int8
	Rating      int64
	Performance string

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
		"rating BIGINT NOT NULL DEFAULT 1500",
		"performance TEXT NOT NULL DEFAULT ''",
		"ADD PRIMARY KEY (uid, contest)",
	)
	registerSchema("contest_match_script",
		"id SERIAL PRIMARY KEY",
		"contest INTEGER NOT NULL REFERENCES contest(id)",
		"hook SMALLINT NOT NULL DEFAULT "+strconv.Itoa(MatchScriptHookManual),
		"contents TEXT NOT NULL",
	)
}

func (c *Contest) Representation(u User) map[string]interface{} {
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
		"my_role":     c.ParticipationOf(u),
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

func (c *Contest) AllParticipations() ([]ContestParticipation, error) {
	rows, err := db.Query("SELECT "+
		"contest_participation.type, "+
		"contest_participation.rating, "+
		"contest_participation.performance, "+
		"users.id, users.handle, users.privilege, users.nickname "+
		"FROM contest_participation "+
		"LEFT JOIN users ON contest_participation.uid = users.id "+
		"WHERE contest = $1 "+
		"ORDER BY contest_participation.rating DESC",
		c.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ps := []ContestParticipation{}
	for rows.Next() {
		p := ContestParticipation{Contest: c.Id}
		err := rows.Scan(&p.Type, &p.Rating, &p.Performance,
			&p.Rel.User.Id, &p.Rel.User.Handle,
			&p.Rel.User.Privilege, &p.Rel.User.Nickname)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, rows.Err()
}

func (c *Contest) Update() error {
	_, err := db.Exec("UPDATE contest SET "+
		"title = $1, banner = $2, owner = $3, "+
		"start_time = $4, end_time = $5, descr = $6, "+
		"details = $7, is_visible = $8, is_reg_open = $9 "+
		"WHERE id = $10",
		c.Title,
		c.Banner,
		c.Owner,
		c.StartTime,
		c.EndTime,
		c.Desc,
		c.Details,
		c.IsVisible,
		c.IsRegOpen,
		c.Id,
	)
	return err
}

func (c *Contest) HasStarted() bool {
	return time.Now().Unix() >= c.StartTime
}

func (c *Contest) HasEnded() bool {
	return time.Now().Unix() >= c.EndTime
}

func (c *Contest) IsRunning() bool {
	return c.HasStarted() && !c.HasEnded()
}

func (c *Contest) ParticipationOf(u User) int8 {
	if c.Owner == u.Id || u.Privilege == UserPrivilegeSuperuser {
		return ParticipationTypeModerator
	}

	// Look for the participation record
	p := ContestParticipation{
		User:    u.Id,
		Contest: c.Id,
	}
	if err := p.Read(); err != nil {
		if err == sql.ErrNoRows {
			// Did not participate
			return -1
		} else {
			panic(err)
		}
	}
	return p.Type
}

func (c *Contest) IsVisibleTo(u User) bool {
	return c.IsVisible || c.ParticipationOf(u) != -1
}

func (p *ContestParticipation) Representation() map[string]interface{} {
	return map[string]interface{}{
		"participant": p.Rel.User.ShortRepresentation(),
		"rating":      p.Rating,
		"performance": p.Performance,
	}
}

func (p *ContestParticipation) Create() error {
	_, err := db.Exec("INSERT INTO "+
		"contest_participation(uid, contest, type, rating, performance) "+
		"VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
		p.User,
		p.Contest,
		p.Type,
		p.Rating,
		p.Performance,
	)
	return err
}

func (p *ContestParticipation) Read() error {
	err := db.QueryRow("SELECT type, rating, performance "+
		"FROM contest_participation WHERE uid = $1 AND contest = $2",
		p.User,
		p.Contest,
	).Scan(&p.Type, &p.Rating, &p.Performance)
	return err
}

func (p *ContestParticipation) Update() error {
	_, err := db.Exec("UPDATE contest_participation SET "+
		"type = $1, rating = $2, performance = $3 "+
		"WHERE uid = $4 AND contest = $5",
		p.Type,
		p.Rating,
		p.Performance,
		p.User,
		p.Contest,
	)
	return err
}
