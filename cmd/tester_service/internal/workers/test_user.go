package workers

const (
	emailDomain = "@interview.com"
	password    = "VeryG00dPass!"
	numUsers    = 50
)

type testUser struct {
	ValidationCode string
	Username       string
	Email          string
	Password       string
}
