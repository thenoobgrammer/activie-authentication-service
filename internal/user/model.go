package user

import (
	"github.com/lib/pq"
)

type User struct {
	ID          string
	Email       string
	SystemRoles *pq.StringArray
}
