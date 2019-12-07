package models

import (
	"fmt"
	"strconv"
	"time"
)

func fakeCreateUser(handle string, privilege int8, bio string) {
	u := User{
		Handle:    handle,
		Email:     handle + "@example.com",
		Password:  "qwq",
		Privilege: privilege,
		Nickname:  "~ " + handle + " ~",
		Bio:       bio,
	}
	if err := u.Create(); err != nil {
		panic(err)
	}
	println("User " + handle + " created")
}

func FakeDatabase() {
	// Clear database
	for _, schema := range schemata {
		_, err := db.Exec("DROP TABLE IF EXISTS " + schema.table + " CASCADE")
		if err != nil {
			panic(err)
		}
	}

	InitializeSchemata(db)

	// Users
	// - Superuser
	fakeCreateUser("su", UserPrivilegeSuperuser, "I have been notified")
	// - Organizers
	for i := 1; i <= 5; i++ {
		fakeCreateUser("o"+strconv.Itoa(i), UserPrivilegeOrganizer, "Enjoy the contests")
	}
	// - Participants
	for i := 1; i <= 20; i++ {
		fakeCreateUser("p"+strconv.Itoa(i), UserPrivilegeNormal, "I'm a teapot")
	}

	// Contests
	t := time.Now().Unix()
	for i := 1; i <= 5; i++ {
		c := Contest{
			Title:     "Grand Contest" + strconv.Itoa(i),
			Banner:    "banner.png",
			Owner:     int32(1 + i),
			StartTime: t + 3600*int64(-3+i),
			EndTime:   t + 3600*int64(-1+i),
			Desc:      "Really big contest #" + strconv.Itoa(i),
			Details:   "Lorem ipsum dolor sit amet",
			IsVisible: true,
			IsRegOpen: true,
		}
		if err := c.Create(); err != nil {
			panic(err)
		}

		sidFirst, sidLast := int32(-1), int32(-1)

		// Participants
		for j := 1 + i/2; j <= 20; j += i {
			fmt.Printf("User %d joins contest %d\n", j, i)
			p := ContestParticipation{
				User:    int32(6 + j),
				Contest: int32(i),
				Type:    ParticipationTypeContestant,
			}
			if err := p.Create(); err != nil {
				panic(err)
			}

			// Submissions
			for k := 1; k <= 2+(i+j)%3; k++ {
				s := Submission{
					User:     int32(6 + j),
					Contest:  int32(i),
					Contents: "print(" + strconv.Itoa(i+j+k) + ")",
				}
				if err := s.Create(); err != nil {
					panic(err)
				}
				if sidFirst == -1 {
					sidFirst = s.Id
				}
				sidLast = s.Id
			}
		}

		// Matches
		count := sidLast - sidFirst + 1
		seed := int32(129)
		for j := 1; j <= 30; j++ {
			m := Match{
				Contest: int32(i),
				Report:  "{\"winner\": \"In queue\"}",
			}
			u := seed % count
			seed = ((seed * 1103515245) + 12345) & 0x7fffffff
			v := seed % (count - 1)
			seed = ((seed * 1103515245) + 12345) & 0x7fffffff
			if u == v {
				v = count - 1
			}
			m.Rel.Parties = []Submission{
				Submission{Id: sidFirst + u},
				Submission{Id: sidFirst + v},
			}
			fmt.Printf("Match takes place between submissions %d and %d\n",
				sidFirst+u, sidFirst+v)
			if err := m.Create(); err != nil {
				panic(err)
			}
		}
	}
}
