package bot

import (
	"github.com/orsenkucher/viber"
)

//var UserPhoneMap map[string]*User
var UserIDMap map[string]*User
var UserData map[string]interface{}

func init() {
	//UserPhoneMap = map[string]*User{}
	UserIDMap = map[string]*User{}
	UserData = map[string]interface{}{}
}

type User struct {
	ViberUser   viber.User
	PhoneNumber string
	Token       string
	Password    string
}
