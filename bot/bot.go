package bot

import (
	"fmt"

	"github.com/orsenkucher/nothing/encio"
	"github.com/orsenkucher/viber"
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

func Send(v *viber.Viber, id string, m viber.Message) {
	_, err := v.SendMessage(id, m)
	if err != nil {
		fmt.Println(err)
	}
}
