package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/mileusna/viber"
)

func MyConversaionStarted(v *viber.Viber, u viber.User, conversationType, context string, subscribed bool, token uint64, t time.Time) viber.Message {
	fmt.Println("new subscriber", u.ID)

	//b := v.NewButton(2, 2, viber.Reply, "qwe", "1", "")
	b := v.NewButton(2, 2, viber.Reply, "qwe", "1", "https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Smiley.svg/1200px-Smiley.svg.png")
	k := v.NewKeyboard("", false)
	k.AddButton(b)
	msg := v.NewTextMessage("Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ")
	msg.SetKeyboard(k)
	return msg
}

// myMsgReceivedFunc will be called everytime when user send us a message.
func MyMsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {

	switch tm := m.(type) {
	case *viber.TextMessage:
		fmt.Println(u.Mcc, u.Mnc, u.DeviceType, u.Name, u.PrimaryDeviceOs, u.Country)
		_, _ = v.SendTextMessage(u.ID, "Thank you for your message")
		txt := tm.Text
		_, _ = v.SendTextMessage(u.ID, "This is the text you have sent to me "+txt)

		if txt == "button" {
			fmt.Println("button")

			b := v.NewButton(2, 2, viber.Reply, "qwe", "1", "")
			k := v.NewKeyboard("", false)
			k.AddButton(b)

			b.Text = "2"
			k.AddButton(b)

			msg := v.NewTextMessage("qwe")
			msg.Keyboars = k
			_, _ = v.SendMessage(u.ID, msg)
		}

	case *viber.URLMessage:
		url := m.(*viber.URLMessage).Media
		_, _ = v.SendTextMessage(u.ID, "You have sent me an interesting link "+url)

	case *viber.PictureMessage:
		_, _ = v.SendTextMessage(u.ID, "Nice pic!")
	}
}

func MyDeliveredFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "delivered to user ID", userID)
}

func MySeenFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "seen by user ID", userID)
}
