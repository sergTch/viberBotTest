package bot

import (
	"fmt"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func AddOpperation(v *viber.Viber, msg *viber.RichMediaMessage, item abm.HistoryItem) {
	rows := 0
	if item.Type == "other" {
		return
	}
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
		adress, ok := item.Data["shop_address"]
		if ok && adress != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Начислено: ", adress), "", true))
			rows++
		}
		name, ok := item.Data["shop_name"]
		if ok && name != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Начислено: ", name), "", true))
			rows++
		}
	}

	if item.Type == "pending" {
		category, ok := item.Data["category"]
		if ok && category != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Категория: ", category), "", true))
			rows++
		}
		bonus, ok := item.Data["bonus"]
		if ok && bonus != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Бонус: ", bonus), "", true))
			rows++
		}
	}

	if item.Type == "gift" {
		bonus, ok := item.Data["bonus"]
		if ok && bonus != "0" {
			msg.AddButton(v.NewButton(6, 1, viber.None, "", fmt.Sprint("Бонус: ", bonus), "", true))
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
		msg.AddButton(v.NewButton(6, rows-7, viber.None, "", "", "", true))
	}
}
