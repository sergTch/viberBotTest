package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/orsenkucher/nothing/encio"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/bot"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var s = flag.String("s", "", "provide encio password")
	var h = flag.Bool("h", false, "set webhook if true")

	flag.Parse()

	profile, err := abm.Client.Profile("")
	fmt.Printf("%+v %+v\n", profile, err)
	println("------")
	fmt.Println(profile.Schema("birth_day"))
	fmt.Println(profile.Schema("gender"))

	if *s == "" {
		return errors.New("[-s] -> encio must be handled")
	}

	key := encio.NewEncIO(*s)

	cfg, err := key.GetConfig("secure/viber.json")
	if err != nil {
		return err
	}

	v := bot.NewBot(cfg)

	// you really need this only once, remove after you set the webhook

	if *h {
		hook, err := v.SetWebhook("https://loyalty-vbot.abmloyalty.app/viber/webhook", nil)
		if err != nil {
			return err
		}

		log.Printf("%+v", hook)
	}

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
