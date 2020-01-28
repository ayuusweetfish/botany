package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Contest struct {
	Id     int32
	Title  string
	Banner int32

	Owner     int32
	StartTime int64
	EndTime   int64

	Desc    string
	Details string

	IsVisible bool
	IsRegOpen bool

	Judge    int32
	Script   string
	Playback string

	Rel struct {
		Owner          User
		Participations []ContestParticipation
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
	Delegate    int32
	Rating      int64
	Performance string

	Rel struct {
		User    User
		Contest Contest
	}
}

func init() {
	registerSchema("contest",
		"id SERIAL PRIMARY KEY",
		"title TEXT NOT NULL DEFAULT ''",
		"banner INTEGER", // Nullable
		"owner INTEGER NOT NULL",
		"start_time BIGINT NOT NULL",
		"end_time BIGINT NOT NULL",
		"descr TEXT NOT NULL DEFAULT ''",
		"details TEXT NOT NULL DEFAULT ''",
		"is_visible BOOLEAN NOT NULL DEFAULT FALSE",
		"is_reg_open BOOLEAN NOT NULL DEFAULT FALSE",
		"judge INTEGER", // Nullable
		"script TEXT NOT NULL DEFAULT ''",
		"script_log TEXT NOT NULL DEFAULT ''",
		"playback TEXT NOT NULL DEFAULT ''",
		"ADD CONSTRAINT fk_banner FOREIGN KEY (banner) REFERENCES file (id)",
		"ADD CONSTRAINT fk_users FOREIGN KEY (owner) REFERENCES users (id)",
		"ADD CONSTRAINT fk_judge FOREIGN KEY (judge) REFERENCES submission (id)",
	)
	registerSchema("contest_participation",
		"uid INTEGER NOT NULL",
		"contest INTEGER NOT NULL",
		"type SMALLINT NOT NULL",
		"delegate INTEGER", // Nullable
		"rating BIGINT NOT NULL DEFAULT 0",
		"performance TEXT NOT NULL DEFAULT ''",
		"ADD PRIMARY KEY (uid, contest)",
		"ADD CONSTRAINT fk_users FOREIGN KEY (uid) REFERENCES users (id)",
		"ADD CONSTRAINT fk_contest FOREIGN KEY (contest) REFERENCES contest (id)",
		"ADD CONSTRAINT fk_delegate FOREIGN KEY (delegate) REFERENCES submission (id)",
	)
}

func (c *Contest) Representation(u User) map[string]interface{} {
	mods := []map[string]interface{}{}
	rows, err := db.Query("SELECT uid FROM contest_participation where contest = $1 AND type = $2", c.Id, ParticipationTypeModerator)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	for rows.Next() {
		u := User{}
		_ = rows.Scan(&u.Id)
		u.ReadById()
		mods = append(mods, u.ShortRepresentation())
	}
	return map[string]interface{}{
		"id":          c.Id,
		"title":       c.Title,
		"banner":      c.Banner,
		"start_time":  c.StartTime,
		"end_time":    c.EndTime,
		"desc":        c.Desc,
		"details":     c.Details,
		"is_visible":  c.IsVisible,
		"is_reg_open": c.IsRegOpen,
		"judge":       c.Judge,
		"script":      c.Script,
		"owner":       c.Rel.Owner.ShortRepresentation(),
		"moderators":  mods,
		"my_role":     c.ParticipationOf(u),
	}
}

func (c *Contest) ShortRepresentation(u User) map[string]interface{} {
	return map[string]interface{}{
		"id":          c.Id,
		"title":       c.Title,
		"banner":      c.Banner,
		"start_time":  c.StartTime,
		"end_time":    c.EndTime,
		"desc":        c.Desc,
		"is_reg_open": c.IsRegOpen,
		"my_role":     c.ParticipationOf(u),
	}
}

func (c *Contest) Create() error {
	err := db.QueryRow("INSERT INTO "+
		"contest(title, owner, start_time, end_time, descr, details, is_visible, is_reg_open, script, playback) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		c.Title,
		c.Owner,
		c.StartTime,
		c.EndTime,
		c.Desc,
		c.Details,
		c.IsVisible,
		c.IsRegOpen,
		c.Script,
		c.Playback,
	).Scan(&c.Id)
	return err
}

func (c *Contest) Read() error {
	err := db.QueryRow("SELECT "+
		"title, COALESCE(banner, -1), owner, start_time, end_time, descr, details, "+
		"is_visible, is_reg_open, COALESCE(judge, -1) AS judge, script "+
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
		&c.Judge,
		&c.Script,
	)
	return err
}

func (c *Contest) ReadModerators() ([]int32, error) {
	mods := []int32{}
	var uid int32
	rows, err := db.Query("SELECT uid FROM contest_participation where contest = $1 AND type = $2", c.Id, ParticipationTypeModerator)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	for rows.Next() {
		_ = rows.Scan(&uid)
		mods = append(mods, uid)
	}
	err = db.QueryRow("SELECT owner FROM contest WHERE id = $1", c.Id).Scan(&c.Owner)
	if err != nil {
		return nil, err
	}
	mods = append(mods, c.Owner)
	return mods, nil
}

func (c *Contest) ReadScriptLog() (error, string) {
	var s string
	err := db.QueryRow("SELECT script_log FROM contest WHERE id = $1", c.Id).Scan(&s)
	return err, s
}

func (c *Contest) AppendScriptLog(s string) error {
	_, err := db.Exec("UPDATE contest SET script_log = script_log || $1 WHERE id = $2", s, c.Id)
	return err
}

func ContestReadAll() ([]Contest, error) {
	rows, err := db.Query("SELECT " +
		"id, title, COALESCE(banner, -1), owner, start_time, end_time, descr, " +
		"is_visible, is_reg_open, COALESCE(judge, -1) AS judge, script " +
		"FROM contest ORDER BY id ASC",
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
			&c.Judge,
			&c.Script,
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

func (c *Contest) LoadBanner() (File, error) {
	err := db.QueryRow("SELECT COALESCE(banner, -1) FROM contest WHERE id = $1",
		c.Id).Scan(&c.Banner)
	if err != nil {
		return File{}, err
	}
	if c.Banner == -1 {
		return File{Id: -1, Content: nil}, nil
	}
	f := File{Id: c.Banner}
	if err := f.Read(); err != nil {
		return File{}, err
	}
	return f, nil
}

func (c *Contest) UpdateBanner() error {
	_, err := db.Exec("UPDATE contest SET "+
		"banner = $1 WHERE id = $2", c.Banner, c.Id)
	return err
}

func (c *Contest) AllParticipationsRequiresDelegate(d bool) ([]ContestParticipation, error) {
	delegateCond := ""
	if d {
		delegateCond = " AND delegate != -1"
	}
	rows, err := db.Query("SELECT "+
		"contest_participation.type, "+
		"COALESCE(contest_participation.delegate, -1) AS delegate, "+
		"contest_participation.rating, "+
		"contest_participation.performance, "+
		"users.id, users.handle, users.privilege, users.nickname "+
		"FROM contest_participation "+
		"LEFT JOIN users ON contest_participation.uid = users.id "+
		"WHERE contest = $1 AND type = $2"+delegateCond+
		"ORDER BY contest_participation.rating DESC",
		c.Id, ParticipationTypeContestant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ps := []ContestParticipation{}
	for rows.Next() {
		p := ContestParticipation{Contest: c.Id}
		err := rows.Scan(&p.Type, &p.Delegate, &p.Rating, &p.Performance,
			&p.Rel.User.Id, &p.Rel.User.Handle,
			&p.Rel.User.Privilege, &p.Rel.User.Nickname)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, rows.Err()
}

func (c *Contest) AllParticipations() ([]ContestParticipation, error) {
	return c.AllParticipationsRequiresDelegate(false)
}

func (c *Contest) AllParticipationsWithDelegate() ([]ContestParticipation, error) {
	return c.AllParticipationsRequiresDelegate(true)
}

func (c *Contest) PartParticipation(limit, offset int) ([]ContestParticipation, int, error) {
	rows, err := db.Query("SELECT "+
		"contest_participation.type, "+
		"contest_participation.rating, "+
		"contest_participation.performance, "+
		"COALESCE(contest_participation.delegate, -1), "+
		"users.id, users.handle, users.privilege, users.nickname "+
		"FROM contest_participation "+
		"LEFT JOIN users ON contest_participation.uid = users.id "+
		"WHERE contest = $1 AND type = $2"+
		"ORDER BY contest_participation.rating DESC, users.handle ASC LIMIT $3 OFFSET $4",
		c.Id, ParticipationTypeContestant, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	ps := []ContestParticipation{}
	for rows.Next() {
		p := ContestParticipation{Contest: c.Id}
		err := rows.Scan(&p.Type, &p.Rating, &p.Performance, &p.Delegate,
			&p.Rel.User.Id, &p.Rel.User.Handle,
			&p.Rel.User.Privilege, &p.Rel.User.Nickname)
		if err != nil {
			return nil, 0, err
		}
		ps = append(ps, p)
	}
	var total int
	rows2 := db.QueryRow("SELECT COUNT(*) from contest_participation "+
		"where contest = $1 AND type = $2", c.Id, ParticipationTypeContestant)
	err = rows2.Scan(&total)
	return ps, total, rows.Err()
}

func (c *Contest) Update() error {
	_, err := db.Exec("UPDATE contest SET "+
		"title = $1, owner = $2, "+
		"start_time = $3, end_time = $4, descr = $5, details = $6, "+
		"is_visible = $7, is_reg_open = $8, judge = NULLIF($9, -1), script = $10, playback = $11 "+
		"WHERE id = $12",
		c.Title,
		c.Owner,
		c.StartTime,
		c.EndTime,
		c.Desc,
		c.Details,
		c.IsVisible,
		c.IsRegOpen,
		c.Judge,
		c.Script,
		c.Playback,
		c.Id,
	)
	return err
}

func (c *Contest) UpdateModerators(uids []int64) error {
	_, err := db.Exec("DELETE FROM contest_participation "+
		"WHERE contest = $1 AND type = $2",
		c.Id,
		ParticipationTypeModerator)
	if err != nil {
		return err
	}

	if len(uids) > 0 {
		stmt := "INSERT INTO contest_participation" +
			"(uid, contest, type, rating, performance) VALUES "
		vals := []interface{}{}
		for i, uid := range uids {
			if i != 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("($%d, $%d, $%d, 0, '')", i*3+1, i*3+2, i*3+3)
			vals = append(vals, uid, c.Id, ParticipationTypeModerator)
		}
		stmt += " ON CONFLICT (uid, contest) DO UPDATE SET type = " + strconv.Itoa(ParticipationTypeModerator)
		if _, err := db.Exec(stmt, vals...); err != nil {
			return err
		}
	}

	return nil
}

func (c *Contest) LoadPlayback() error {
	return db.QueryRow("SELECT playback FROM contest WHERE id = $1",
		c.Id).Scan(&c.Playback)
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
		"delegate":    p.Delegate,
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
	p.Delegate = -1
	return err
}

func (p *ContestParticipation) Read() error {
	err := db.QueryRow("SELECT type, COALESCE(delegate, -1), rating, performance "+
		"FROM contest_participation WHERE uid = $1 AND contest = $2",
		p.User,
		p.Contest,
	).Scan(&p.Type, &p.Delegate, &p.Rating, &p.Performance)
	return err
}

func (p *ContestParticipation) Update() error {
	_, err := db.Exec("UPDATE contest_participation SET "+
		"type = $1, rating = $2, delegate = NULLIF($3, -1), "+
		"performance = $4 "+
		"WHERE uid = $5 AND contest = $6",
		p.Type,
		p.Rating,
		p.Delegate,
		p.Performance,
		p.User,
		p.Contest,
	)
	return err
}

func (p *ContestParticipation) UpdateStats() error {
	_, err := db.Exec("UPDATE contest_participation SET "+
		"rating = $1, performance = $2 "+
		"WHERE uid = $3 AND contest = $4",
		p.Rating, p.Performance,
		p.User, p.Contest,
	)
	return err
}
