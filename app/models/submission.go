package models

type Submission struct {
	Id            int32
	Participation int32
	Contents      string

	Rel struct {
		Participation ContestParticipation
	}
}

func init() {
	registerSchema("submission",
		"id SERIAL PRIMARY KEY",
		"participation INTEGER NOT NULL REFERENCES contest_participation(id)",
		"contents TEXT NOT NULL",
	)
}

func (s *Submission) Create() error {
	// TODO
	return nil
}

func (s *Submission) Read() error {
	// TODO
	return nil
}

func (s *Submission) LoadRel() error {
	// TODO
	return nil
}

func (s *Submission) Update() error {
	// TODO
	return nil
}
