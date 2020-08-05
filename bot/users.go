package bot

import (
	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

//var UserPhoneMap map[string]*User
var UserIDMap map[string]*User
var UserSMS map[string]SMS
var UserFields map[string][]*abm.Field
var NextAction map[string]*ButtAction

func init() {
	//UserPhoneMap = map[string]*User{}
	UserIDMap = map[string]*User{}
	UserFields = map[string][]*abm.Field{}
	UserSMS = map[string]SMS{}
	NextAction = map[string]*ButtAction{}
}

type User struct {
	ViberUser   viber.User
	PhoneNumber string
	Token       string
	Password    string
}
