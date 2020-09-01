package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/data"
)

func LastOperationsDet(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, n int) {
	if user, ok := UserIDMap[u.ID]; ok {
		fmt.Println(user.Token.Token())
		history, err := abm.Client.ClientHistory(user.Token, n/20+1)
		if (err != nil) || history.Meta.TotalCount < n%20 {
			_, err := v.SendTextMessage(u.ID, data.Translate(user.Language, "Нету предидущих транзакций;("))
			check(err)
			Menu(v, u, m, token, t)
			return
		}
		item := history.Items[n%20]
		text := item.Type + "\n"
		sum, ok := item.Data["details_sum"]
		if ok && sum != "0" {
			text += fmt.Sprintln("Сумма: ", sum)
		}
		written, ok := item.Data["written_off"]
		if ok && written != "0" {
			text += fmt.Sprintln("Списано: ", written)
		}
		text += "________________\n"
		for _, det := range item.Details {
			text += fmt.Sprintln(det["product_name"], "\n", det["product_price"], "x", det["product_amount"], "=", det["product_sum"])
		}
		text += "________________\n"
		name, ok := item.Data["shop_name"]
		if ok {
			adress, ok := item.Data["shop_address"]
			if ok {
				text += fmt.Sprintln("Магазин: ", name)
				text += fmt.Sprintln("Адресс: ", adress)
			}
		}
		date, ok := item.Data["date"]
		if ok {
			text += fmt.Sprint("Дата: ", parseDate(date))
		}
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 3, 1, "", "Назад", "hist", strconv.Itoa(n-n%5)))
		keyboard.AddButtons(*BuildButton(v, 3, 1, "", data.Translate(user.Language, "Меню"), "mnu"))
		keyboard.InputFieldState = viber.HiddenInputField
		msg.SetKeyboard(keyboard)

		_, err = v.SendMessage(u.ID, msg)
		check(err)
	}
}

func ActionsDet(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, n int) {
	if user, ok := UserIDMap[u.ID]; ok {
		fmt.Println(user.Token.Token())
		actions, meta, err := abm.Client.Actions(user.Token, n/20+1)
		if (err != nil) || meta.TotalCount < n%20 {
			_, err := v.SendTextMessage(u.ID, data.Translate(user.Language, "Нету новостей;("))
			check(err)
			Menu(v, u, m, token, t)
			return
		}
		msg := v.NewTextMessage(actions[n%20].Title + "\n" + actions[n%20].Content + "\n" + parseDate(actions[n%20].From) + "\n" + parseDate(actions[n%20].To))
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 3, 1, "", "Назад", "acts", strconv.Itoa(n-n%5)))
		keyboard.AddButtons(*BuildButton(v, 3, 1, "", data.Translate(user.Language, "Меню"), "mnu"))
		keyboard.InputFieldState = viber.HiddenInputField
		msg.SetKeyboard(keyboard)

		_, err = v.SendMessage(u.ID, msg)
		check(err)
	}
}

func NewsDet(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, n int) {
	if user, ok := UserIDMap[u.ID]; ok {
		fmt.Println(user.Token.Token())
		news, meta, err := abm.Client.News(user.Token, n/20+1)
		if (err != nil) || meta.TotalCount < n%20 {
			_, err := v.SendTextMessage(u.ID, data.Translate(user.Language, "Нету новостей;("))
			check(err)
			Menu(v, u, m, token, t)
			return
		}
		msg := v.NewTextMessage(news[n%20].Name + "\n" + news[n%20].Descr)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 3, 1, "", "Назад", "news", strconv.Itoa(n-n%5)))
		keyboard.AddButtons(*BuildButton(v, 3, 1, "", data.Translate(user.Language, "Меню"), "mnu"))
		keyboard.InputFieldState = viber.HiddenInputField
		msg.SetKeyboard(keyboard)

		_, err = v.SendMessage(u.ID, msg)
		check(err)
	}
}
