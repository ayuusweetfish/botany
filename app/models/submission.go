package models

type Submission struct {
	Id            int32
	Participation int32

	Rel struct {
		Participation ContestParticipation
	}
}

func init() {
	registerSchema("submission",
		"id SERIAL PRIMARY KEY",
		"participation INTEGER NOT NULL REFERENCES contest_participation(id)",
	)
}

func (u *Submission) Create() error {
	// TODO
	return nil
}

func (u *Submission) Read() error {
	// TODO
	return nil
}

func (u *Submission) LoadRel() error {
	// TODO
	return nil
}

func (u *Submission) Update() error {
	// TODO
	return nil
}
