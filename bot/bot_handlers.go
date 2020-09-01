package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/data"
)

func generatePhoneNumber(number string) string {
	str := ""
	for _, ch := range number {
		if ch >= '0' && ch <= '9' {
			str += string(ch)
		}
	}
	return "+" + str[0:3] + "(" + str[3:5] + ")" + str[5:8] + "-" + str[8:10] + "-" + str[10:12]
	//return str
}

func MyConversaionStarted(v *viber.Viber, u viber.User, conversationType, context string, subscribed bool, token uint64, t time.Time) viber.Message {
	fmt.Println("new subscriber", u.ID)

	startB := BuildCfgButton(v, data.ButtCfg.Start, true, "agr")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*startB)
	keyboard.InputFieldState = viber.HiddenInputField
	UserTxtAct[u.ID] = []*TextAction{}
	msg := v.NewTextMessage(data.Translate("", data.Translate("", "Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ")))
	msg.SetKeyboard(keyboard)
	return msg
}

// myMsgReceivedFunc will be called everytime when user send us a message.
func MyMsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	fmt.Println(u.ID, " response")
	msg, txtMsg := m.(*viber.TextMessage)
	if _, ok := UserIDMap[u.ID]; !ok && txtMsg && msg.Text != "#butt/str" && msg.Text != "#butt/agr" {
		StartMsg(v, u, *v.NewTextMessage(""), token, t)
		return
	}
	switch m := m.(type) {
	case *viber.TextMessage:
		txt := m.Text
		fmt.Println(txt)
		parts := strings.Split(txt, "/")
		if parts[0] == "#butt" {
			if parts[1] == "prof" {
				ChangeProfField(v, u, *m, token, t, parts[2])
			} else if parts[1] == "hist" {
				n, err := strconv.Atoi(parts[2])
				if err != nil {
					return
				}
				LastOperations(v, u, *m, token, t, n)
			} else if parts[1] == "news" {
				n, err := strconv.Atoi(parts[2])
				if err != nil {
					return
				}
				News(v, u, *m, token, t, n)
			} else if parts[1] == "acts" {
				n, err := strconv.Atoi(parts[2])
				if err != nil {
					return
				}
				Actions(v, u, *m, token, t, n)
			} else if parts[1] == "histdet" {
				n, err := strconv.Atoi(parts[2])
				if err != nil {
					return
				}
				LastOperationsDet(v, u, *m, token, t, n)
			} else if parts[1] == "newsdet" {
				n, err := strconv.Atoi(parts[2])
				if err != nil {
					return
				}
				NewsDet(v, u, *m, token, t, n)
			} else if parts[1] == "actsdet" {
				n, err := strconv.Atoi(parts[2])
				if err != nil {
					return
				}
				ActionsDet(v, u, *m, token, t, n)
			} else {
				for _, actionID := range parts {
					if action, ok := ButtActions[actionID]; ok {
						action.Act(v, u, *m, token, t)
					}
				}
			}
		} else {
			for _, actions := range UserTxtAct {
				for _, action := range actions {
					action.Act(v, u, *m, token, t)
				}
			}
		}

	case *viber.URLMessage:
		url := m.Media
		_, _ = v.SendTextMessage(u.ID, "You have sent me an interesting link "+url)

	case *viber.PictureMessage:
		_, _ = v.SendTextMessage(u.ID, "Nice pic!")

	case *viber.ContactMessage:
		user := User{ViberUser: u, PhoneNumber: m.Contact.PhoneNumber, Language: data.Cfg.AcceptLanguage}
		UserIDMap[user.ViberUser.ID] = &user
		DB.Save(&user)
		//UserPhoneMap[user.Contact.PhoneNumber] = &user
		fmt.Println(m.Contact.PhoneNumber)
		fmt.Println(generatePhoneNumber(m.Contact.PhoneNumber))
		ok, err := abm.Client.CheckPhone(generatePhoneNumber(m.Contact.PhoneNumber))
		if err != nil {
			fmt.Println(err)
			return
		}
		if !ok {
			msg := v.NewTextMessage(data.Translate("", "Для регистрации в программе лояльности придумайте и отправьте мне пароль. Пароль должен состоять минимум из 6-ти символов"))
			keyboard := v.NewKeyboard("", false)
			//keyboard.AddButtons(*v.NewButton(6, 1, viber.None, "", "Вводите пароль", "", true))
			keyboard.AddButtons(*CfgButton(v, viber.None, data.ButtCfg.EnterPass, "", true))
			msg.Keyboard = keyboard
			_, err := v.SendMessage(u.ID, msg)
			check(err)
			NextAction[user.ViberUser.ID] = ButtActions["ceq"]
			UserTxtAct[u.ID] = []*TextAction{{Act: SetPassword}}
		} else {
			EnterPassword(v, u, *v.NewTextMessage(""), token, t)
			NextAction[u.ID] = ButtActions["mnu"]
		}
		//_, _ = v.SendTextMessage(u.ID, fmt.Sprintf("%s %s", m.Contact.Name, m.Contact.PhoneNumber))
	}
}

func MyDeliveredFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "delivered to user ID", userID)
}

func MySeenFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "seen by user ID", userID)
}
