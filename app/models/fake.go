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
	numbers := []string{"zero", "one", "two", "three", "four", "five"}
	for i := 1; i <= 5; i++ {
		s := "This is the description for contest number " + numbers[i] + "!\n"
		s += "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst."
		script := `
local count = 0
function on_timer(all)
    count = count + 1
    -- if count < 5 then return end
    count = 0
    print('Creating matches for contest #` + strconv.Itoa(i) + `')
    for i = 1, #all do
        print(string.format('Contest %s (%d)', get_handle(all[i]), all[i]))
        if i > 1 then create_match(all[i], all[i - 1]) end
    end
end
`
		c := Contest{
			Title:     "Grand Contest " + strconv.Itoa(i),
			Banner:    "banner.png",
			Owner:     int32(1 + i),
			StartTime: t + 3600*int64(-3+i),
			EndTime:   t + 3600*int64(-1+i),
			Desc:      "Really big contest, number " + numbers[i],
			Details:   s,
			IsVisible: i != 1,
			IsRegOpen: i != 5,
			Script:    script,
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
