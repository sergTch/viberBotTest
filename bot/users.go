package bot

import (
	"github.com/orsenkucher/viber"
)

var UserPhoneMap map[string]*User
var UserIDMap map[string]*User

type User struct {
	Account viber.Account
	Contact viber.Contact
}
