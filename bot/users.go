package bot

import (
	"github.com/jinzhu/gorm"
	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

//var UserPhoneMap map[string]*User
var UserIDMap map[string]*User
var UserSMS map[string]SMS
var UserFields map[string][]*abm.Field
var NextAction map[string]*ButtAction

var DB *gorm.DB

func init() {
	//UserPhoneMap = map[string]*User{}
	UserIDMap = map[string]*User{}
	UserFields = map[string][]*abm.Field{}
	UserSMS = map[string]SMS{}
	NextAction = map[string]*ButtAction{}
}

type User struct {
	gorm.Model
	ViberUser   viber.User
	PhoneNumber string
	Token       *abm.SmartToken
	Password    string
	Language    string
}

func LoadUsers() {
	DB.AutoMigrate(&User{})

	var users []User
	DB.Find(&users)

	for i := range users {
		UserIDMap[users[i].ViberUser.ID] = &users[i]
	}
}

func CleanBase() {
	DB.DropTableIfExists("users")
}
