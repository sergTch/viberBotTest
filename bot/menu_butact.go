package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func Menu(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Новости", "news", "0"))
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Покупки", "hist", "0"))
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Баланс", "sbl"))
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Карточка", "sbq"))
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Профиль", "prf"))
	keyboard.InputFieldState = viber.HiddenInputField
	UserTxtAct[u.ID] = []*TextAction{}
	msg := v.NewTextMessage("Вы попали в меню...")
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	check(err)
}

func LastOperations(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, n int) {
	if user, ok := UserIDMap[u.ID]; ok {
		fmt.Println(user.Token.Token())
		history, err := abm.Client.ClientHistory(user.Token, n/20+1)
		if (err != nil) || history.Meta.TotalCount < n%20 {
			_, err := v.SendTextMessage(u.ID, "Нету предидущих транзакций;(")
			check(err)
			Menu(v, u, m, token, t)
			return
		}
		check(err)
		msg := v.NewRichMediaMessage(6, 7, "#FFFFFF")
		for i := n % 20; i < len(history.Items) && i < n%20+5; i++ {
			AddOpperation(v, msg, history.Items[i])
		}
		keyboard := v.NewKeyboard("", false)
		if n > 0 {
			if n < 5 {
				n = 5
			}
			keyboard.AddButtons(*BuildButton(v, 2, 1, "", "<-", "hist", strconv.Itoa(n-5)))
		} else {
			keyboard.AddButtons(*v.NewButton(2, 1, viber.None, "", "--", "", false))
		}
		keyboard.AddButtons(*BuildButton(v, 2, 1, "", "Меню", "mnu"))
		if n+5 < history.Meta.TotalCount {
			keyboard.AddButtons(*BuildButton(v, 2, 1, "", "->", "hist", strconv.Itoa(n+5)))
		} else {
			keyboard.AddButtons(*v.NewButton(2, 1, viber.None, "", "--", "", false))
		}
		msg.SetKeyboard(keyboard)
		_, err = v.SendMessage(user.ViberUser.ID, msg)
		check(err)
		fmt.Printf("%+v\n", history)
	}
}

func News(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, n int) {
	if user, ok := UserIDMap[u.ID]; ok {
		fmt.Println(user.Token.Token())
		news, meta, err := abm.Client.News(user.Token, n/20+1)
		if (err != nil) || meta.TotalCount < n%20 {
			_, err := v.SendTextMessage(u.ID, "Нету новостей;(")
			check(err)
			Menu(v, u, m, token, t)
			return
		}
		check(err)
		msg := v.NewRichMediaMessage(6, 7, "#FFFFFF")
		for i := n % 20; i < len(news) && i < n%20+5; i++ {
			AddNews(v, msg, &news[i])
		}
		keyboard := v.NewKeyboard("", false)
		if n > 0 {
			if n < 5 {
				n = 5
			}
			keyboard.AddButtons(*BuildButton(v, 2, 1, "", "<-", "news", strconv.Itoa(n-5)))
		} else {
			keyboard.AddButtons(*v.NewButton(2, 1, viber.None, "", "--", "", false))
		}
		keyboard.AddButtons(*BuildButton(v, 2, 1, "", "Меню", "mnu"))
		if n+5 < meta.TotalCount {
			keyboard.AddButtons(*BuildButton(v, 2, 1, "", "->", "news", strconv.Itoa(n+5)))
		} else {
			keyboard.AddButtons(*v.NewButton(2, 1, viber.None, "", "--", "", false))
		}
		msg.SetKeyboard(keyboard)
		_, err = v.SendMessage(user.ViberUser.ID, msg)
		check(err)
		fmt.Printf("%+v\n", news)
	}
}

func ShowBarcode(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	_, barcode, err := abm.Client.BarCode(user.Token)
	check(err)
	msg := v.NewPictureMessage("bar-code", barcode, "")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Меню", "mnu"))
	msg.SetKeyboard(keyboard)
	_, err = v.SendMessage(u.ID, msg)
	check(err)
}

func ShowBalance(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	balance, err := abm.Client.Balance(user.Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(balance)
	msg := v.NewTextMessage("Баланс: " + balance.Balance + balance.Currency + "\n" + "Доступно к списанию: " + balance.Avialable + balance.Currency)
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Меню", "mnu"))
	msg.SetKeyboard(keyboard)
	_, err = v.SendMessage(u.ID, msg)
	check(err)
}
