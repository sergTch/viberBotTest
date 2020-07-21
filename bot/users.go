package bot

import (
	"time"

	"github.com/orsenkucher/viber"
)

type State interface {
	Act(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time)
}

type User struct {
	State   State
	Account viber.Account
	Contact viber.Contact
}
