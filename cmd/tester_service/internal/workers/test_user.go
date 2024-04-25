package workers

type testUser struct {
	ValidationCode string
	Username       string
	Email          string
	Password       string
}

const (
	emailDomain = "@interview.com"
)
