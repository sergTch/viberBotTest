package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func AddOpperation(v *viber.Viber, msg *viber.RichMediaMessage, item abm.HistoryItem, n int) {
	rows := 0
	msg.AddButton(v.NewButton(6, 1, viber.None, "", item.Type, "", true))
	rows++

	if item.Type == "check" || item.Type == "check_return" {
		sum, ok := item.Data["details_sum"]
		if ok && sum != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Сумма: ", sum), "", true))
			rows++
		}
		written, ok := item.Data["written_off"]
		if ok && written != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Списано: ", written), "", true))
			rows++
		}
		accrued, ok := item.Data["accrued"]
		if ok && accrued != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Начислено: ", accrued), "", true))
			rows++
		}
		name, ok := item.Data["shop_name"]
		if ok {
			adress, ok := item.Data["shop_address"]
			if ok {
				msg.AddButton(v.NewButton(3, 1, viber.None, "", fmt.Sprint("Магазин: ", name), "", true))
				msg.AddButton(v.NewButton(3, 1, viber.None, "", fmt.Sprint("Адресс: ", adress), "", true))
				rows++
			}
		}
		date, ok := item.Data["date"]
		if ok {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Дата: ", parseDate(date)), "", true))
			rows++
		}
		msg.AddButton(BuildButton(v, 6, 1, "", "Детали", "histdet", strconv.Itoa(n)))
		rows++
	}

	if item.Type == "pending" {
		category, ok := item.Data["category"]
		if ok {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Категория: ", category), "", true))
			rows++
		}
		bonus, ok := item.Data["bonus"]
		if ok && bonus != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Бонус: ", bonus), "", true))
			rows++
		}
		date, ok := item.Data["date"]
		if ok {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Действ. до ", parseDate(date)), "", true))
			rows++
		}
	}

	if item.Type == "gift" {
		bonus, ok := item.Data["bonus"]
		if ok && bonus != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Бонус: ", bonus), "", true))
			rows++
		}
		date, ok := item.Data["date"]
		if ok {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Действ. до ", parseDate(date)), "", true))
			rows++
		}
	}

	if item.Type == "burn" {
		bonus, ok := item.Data["bonus"]
		if ok && bonus != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Бонус: ", bonus), "", true))
			rows++
		}
	}

	if item.Type == "withdraw" {
		bonus, ok := item.Data["bonus"]
		if ok && bonus != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Бонус: ", bonus), "", true))
			rows++
		}
	}

	if rows < 7 {
		msg.AddButton(v.NewButton(6, 7-rows, viber.None, "", " ", "", true))
	}
}

func AddNews(v *viber.Viber, msg *viber.RichMediaMessage, news *abm.News, n int) {
	rows := 0
	msg.AddButton(v.NewButton(6, 1, viber.None, "", news.Name, "", true))
	rows++

	if news.Image != "" {
		msg.AddButton(v.NewButton(6, 3, viber.None, "", "", news.Image, true))
	} else {
		msg.AddButton(v.NewButton(6, 3, viber.None, "", " ", "", true))
	}
	rows += 3
	msg.AddButton(v.NewButton(6, 2, viber.None, "", news.Short, "", true))
	rows += 2

	msg.AddButton(BuildButton(v, 6, 1, "", "Детали", "newsdet", strconv.Itoa(n)))
	rows++

	if rows < 7 {
		msg.AddButton(v.NewButton(6, 7-rows, viber.None, "", " ", "", true))
	}
}

func AddAction(v *viber.Viber, msg *viber.RichMediaMessage, action *abm.Actions, n int) {
	rows := 0
	msg.AddButton(v.NewButton(6, 1, viber.None, "", action.Title, "", true))
	rows++

	if action.Image != "" {
		msg.AddButton(v.NewButton(6, 3, viber.None, "", "", action.Image, true))
	} else {
		msg.AddButton(v.NewButton(6, 3, viber.None, "", " ", "", true))
	}
	rows += 3
	msg.AddButton(v.NewButton(6, 2, viber.None, "", "C "+action.From+" по "+action.To+"\n"+action.Content, "", true))
	rows += 2

	msg.AddButton(BuildButton(v, 6, 1, "", "Детали", "actsdet", strconv.Itoa(n)))
	rows++

	if rows < 7 {
		msg.AddButton(v.NewButton(6, 7-rows, viber.None, "", " ", "", true))
	}
}

func parseDate(sec interface{}) string {
	seconds, err := strconv.ParseInt(fmt.Sprintf("%v", sec), 10, 64)
	check(err)
	fmt.Println(seconds)
	t := time.Unix(seconds, seconds)
	return fmt.Sprint(t.Day(), "-", int(t.Month()), "-", t.Year())
}
