package templates

import (
	"github.com/kwinso/kubsu-web-api/db"
)

type IndexContextSubmission struct {
	ID        int32
	Name      string
	Phone     string
	Email     string
	BirthDate string
	Bio       string
	Sex       int
	Username  string
	Password  string
	Languages []int32
}

type IndexContext struct {
	Success    bool
	Submission IndexContextSubmission
	Errors     map[string]string
	Languages  []db.Language
}
