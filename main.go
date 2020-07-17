package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/mileusna/viber"
	"github.com/orsenkucher/nothing/encio"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var s = flag.String("s", "", "provide encio password")
	flag.Parse()
	if *s == "" {
		return errors.New("[-s] -> encio must be handled")
	}

	key := encio.NewEncIO(*s)

	cfg, err := key.GetConfig("secure/viber.json")
	if err != nil {
		return err
	}

	v := &viber.Viber{
		AppKey: cfg["token"].(string),
		Sender: viber.Sender{
			Name: "Loyalty bot",
			// Avatar: "https://mysite.com/img/avatar.jpg",
		},
		Message:   myMsgReceivedFunc, // your function for handling messages
		Delivered: myDeliveredFunc,   // your function for delivery report
	}
	v.Seen = mySeenFunc // or assign events after declaration

	// you really need this only once, remove after you set the webhook

	hook, err := v.SetWebhook("https://loyalty-vbot.abmloyalty.app/viber/webhook", nil)
	if err != nil {
		return err
	}
	log.Printf("%+v", hook)

	// userID := "Goxxuipn9xKKRqkFOOwKnw==" // fake user ID, use the real one
	// // send text message
	// token, err := v.SendTextMessage(userID, "Hello, World!")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("Message sent, message token:", token)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "viber://pa?chatURI=abmloyaltytest", http.StatusFound)
	})

	http.Handle("/viber/webhook", v)
	err = http.ListenAndServe(":9094", nil)
	if err != nil {
		return err
	}

	return nil
}

// myMsgReceivedFunc will be called everytime when user send us a message
func myMsgReceivedFunc(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	switch m.(type) {

	case *viber.TextMessage:
		v.SendTextMessage(u.ID, "Thank you for your message")
		txt := m.(*viber.TextMessage).Text
		v.SendTextMessage(u.ID, "This is the text you have sent to me "+txt)

	case *viber.URLMessage:
		url := m.(*viber.URLMessage).Media
		v.SendTextMessage(u.ID, "You have sent me an interesting link "+url)

	case *viber.PictureMessage:
		v.SendTextMessage(u.ID, "Nice pic!")

	}
}

func myDeliveredFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "delivered to user ID", userID)
}

func mySeenFunc(v *viber.Viber, userID string, token uint64, t time.Time) {
	log.Println("Message ID", token, "seen by user ID", userID)
}
