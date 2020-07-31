package bot

import (
	"github.com/orsenkucher/viber"
)

//var UserPhoneMap map[string]*User
var UserIDMap map[string]*User
var UserSMS map[string]SMS
var NextAction map[string]*ButtAction

func init() {
	//UserPhoneMap = map[string]*User{}
	UserIDMap = map[string]*User{}
	UserSMS = map[string]SMS{}
	NextAction = map[string]*ButtAction{}
}

type User struct {
	ViberUser   viber.User
	PhoneNumber string
	Token       string
	Password    string
}
