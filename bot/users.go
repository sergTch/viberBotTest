package bot

import (
	"github.com/orsenkucher/viber"
)

var UserPhoneMap map[string]*User
var UserIDMap map[string]*User

func init() {
	UserPhoneMap = map[string]*User{}
	UserIDMap = map[string]*User{}
}

type User struct {
	ViberUser viber.User
	Contact   viber.Contact
}
