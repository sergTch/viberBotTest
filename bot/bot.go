package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/nothing/encio"
	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/data"
)

func NewBot(cfg encio.Config) *viber.Viber {
	return &viber.Viber{
		AppKey: cfg["token"].(string),
		MinAPI: 4,
		Sender: viber.Sender{
			Name: "Loyalty bot",
			// Avatar: "https://mysite.com/img/avatar.jpg",
		},
		ConversationStarted: MyConversaionStarted,
		Message:             MyMsgReceivedFunc, // your function for handling messages
		Delivered:           MyDeliveredFunc,   // your function for delivery report
		Seen:                MySeenFunc,        // or assign events after declaration
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func checkServerError(err error, v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	check(err)
	if err != nil {
		_, err = v.SendTextMessage(u.ID, data.Translate("", "Извините, что-то пошло не так. Попробуйте ещё раз"))
		check(err)
		Menu(v, u, m, token, t)
		return
	}
}
