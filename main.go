package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/orsenkucher/nothing/encio"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/bot"
	"github.com/sergTch/viberBotTest/data"
)

func main() {
	cfg := flag.String("cfg", "secure", "config directory")
	s := flag.String("s", "", "provide encio password")
	h := flag.Bool("h", false, "set webhook if true")

	flag.Parse()

	data.Init(*cfg)
	abm.Init()

	if err := run(*s, *h); err != nil {
		log.Fatalln(err)
	}
}

func run(s string, h bool) error {
	if s == "" {
		return errors.New("[-s] -> encio must be handled")
	}

	key := encio.NewEncIO(s)

	db, err := NewDB(key)
	if err != nil {
		return err
	}
	defer db.Close()
	bot.DB = db
	bot.LoadUsers()

	cfg, err := key.GetConfig(data.Viber)
	if err != nil {
		return err
	}

	v := bot.NewBot(cfg)

	// you really need this only once, remove after you set the webhook
	if h {
		hook, err := v.SetWebhook(data.Cfg.Webhook, nil)
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
		http.Redirect(w, r, fmt.Sprintf("viber://pa?chatURI=%s", data.Cfg.ChatURI), http.StatusFound)
	})

	http.Handle("/viber/webhook", v)

	err = http.ListenAndServe(fmt.Sprintf(":%v", data.Cfg.Port), nil)
	if err != nil {
		return err
	}

	return nil
}

func NewDB(key encio.EncIO) (*gorm.DB, error) {
	var cfg struct {
		Driver string `json:"driver"`
		Host   string `json:"host"`
		Port   int    `json:"port"`
		User   string `json:"user"`
		Pass   string `json:"pass"`
		Name   string `json:"name"`
	}

	bytes, err := key.ReadFile(data.Gorm)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name,
	)

	db, err := gorm.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w\n", err)
	}

	return db, nil
}
