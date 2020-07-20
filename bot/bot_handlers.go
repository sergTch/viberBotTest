package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/orsenkucher/viber"
)

func MyConversaionStarted(v *viber.Viber, u viber.User, conversationType, context string, subscribed bool, token uint64, t time.Time) viber.Message {
	fmt.Println("new subscriber", u.ID)

	//b := v.NewButton(2, 2, viber.Reply, "qwe", "1", "")
	b := v.NewButton(2, 2, viber.Reply, "qwe", "1", "https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Smiley.svg/1200px-Smiley.svg.png")
	k := v.NewKeyboard("", false)
	k.AddButtons(*b)
	msg := v.NewTextMessage("Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ")
	msg.SetKeyboard(k)
	return msg
}

// myMsgReceivedFunc will be called everytime when user send us a message.
func MyMsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {

	switch m := m.(type) {
	case *viber.TextMessage:
		fmt.Println(u.Mcc, u.Mnc, u.DeviceType, u.Name, u.PrimaryDeviceOs, u.Country)
		_, _ = v.SendTextMessage(u.ID, "Thank you for your message")
		txt := m.Text
		_, _ = v.SendTextMessage(u.ID, "This is the text you have sent to me "+txt)

		fmt.Printf("msg:%s\n", txt)
		fmt.Printf("eq:%v", txt == "button")
		if strings.Contains(txt, "button") {
			fmt.Println("button")

			b := v.NewButton(2, 2, viber.SharePhone, "qwe", "1", "")
			k := v.NewKeyboard("", false)
			k.AddButtons(*b)

			// b.Text = "2"
			// k.AddButtons(*b)

			// b = v.NewButton(2, 2, viber.SharePhone, "numberqwe", "Give num", "")
			// k.AddButtons(*b)

			msg := v.NewTextMessage("qwe")
			msg.MinAPIVersion = 3
			msg.Keyboard = k
			_, err := v.SendMessage(u.ID, msg)
			if err != nil {
				fmt.Println(err)
			}
		}

	case *viber.URLMessage:
		url := m.Media
		_, _ = v.SendTextMessage(u.ID, "You have sent me an interesting link "+url)

	case *viber.PictureMessage:
		_, _ = v.SendTextMessage(u.ID, "Nice pic!")

	case *viber.ContactMessage:
		fmt.Printf("%+v", m)
		_, _ = v.SendTextMessage(u.ID, fmt.Sprintf("%s %s", m.Contact.Name, m.Contact.PhoneNumber))
	}
}

func MyDeliveredFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "delivered to user ID", userID)
}

func MySeenFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "seen by user ID", userID)
}
