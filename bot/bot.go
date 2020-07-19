package bot

import (
	"github.com/mileusna/viber"
	"github.com/orsenkucher/nothing/encio"
)

func NewBot(cfg encio.Config) *viber.Viber {
	return &viber.Viber{
		AppKey: cfg["token"].(string),
		Sender: viber.Sender{
			Name: "Loyalty bot",
			// Avatar: "https://mysite.com/img/avatar.jpg",
		},
		Subscribed: MySubscribed,      // your function for handling new subscribers
		Message:    MyMsgReceivedFunc, // your function for handling messages
		Delivered:  MyDeliveredFunc,   // your function for delivery report
		Seen:       MySeenFunc,        // or assign events after declaration
	}

}
