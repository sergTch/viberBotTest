package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func MyConversaionStarted(v *viber.Viber, u viber.User, conversationType, context string, subscribed bool, token uint64, t time.Time) viber.Message {
	fmt.Println("new subscriber", u.ID)

	startB := BuildButton(v, 6, 1, "", "СТАРТ", "agr", "qwe")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*startB)
	msg := v.NewTextMessage("Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ")
	msg.SetKeyboard(keyboard)
	UserTxtAct[u.ID] = []*TextAction{}
	return msg
}

// myMsgReceivedFunc will be called everytime when user send us a message.
func MyMsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	fmt.Println(u.ID, " response")
	switch m := m.(type) {
	case *viber.TextMessage:
		txt := m.Text
		fmt.Println(txt)
		parts := strings.Split(txt, "/")
		if parts[0] == "#butt" {
			for _, actionID := range parts {
				if action, ok := ButtActions[actionID]; ok {
					action.Act(v, u, m, token, t)
				}
			}
		} else {
			for _, actions := range UserTxtAct {
				for _, action := range actions {
					action.Act(v, u, m, token, t)
				}
			}
		}

	case *viber.URLMessage:
		url := m.Media
		_, _ = v.SendTextMessage(u.ID, "You have sent me an interesting link "+url)

	case *viber.PictureMessage:
		_, _ = v.SendTextMessage(u.ID, "Nice pic!")

	case *viber.ContactMessage:
		fmt.Printf("%+v", m)
		ok, err := abm.Client.CheckPhone(m.Contact.PhoneNumber)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ok)
		if !ok {
			_, err := v.SendTextMessage(u.ID, "Для регистрации в программе лояльности придумайте и отправьте мне пароль. Пароль должен состоять минимум из 6-ти символов")
			if err != nil {
				fmt.Println(err)
			}
			UserTxtAct[u.ID] = []*TextAction{{Act: Registration}}
		} else {
			v.SendTextMessage(u.ID, "Дратути")
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
