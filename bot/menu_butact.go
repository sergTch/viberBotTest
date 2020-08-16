package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func Menu(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Покупки", "hist", "0"))
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Карточка", "sbq"))
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Go to start", "str"))
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Go to profile", "prf"))
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
			Menu(v, u, m, token, t)
			return
		}
		check(err)
		msg := v.NewRichMediaMessage(6, 7, "#FFFFFF")
		for i := n % 20; i < len(history.Items) && i < 5; i++ {
			AddOpperation(v, msg, history.Items[i])
		}
		// keyboard := v.NewKeyboard("", false)
		// if n > 0 {
		// 	if n < 5 {
		// 		n = 5
		// 	}
		// 	keyboard.AddButtons(*BuildButton(v, 2, 1, "", "<-", "hist", strconv.Itoa(n-5)))
		// } else {
		// 	keyboard.AddButtons(*v.NewButton(2, 1, viber.None, "", "--", "", false))
		// }
		// keyboard.AddButtons(*BuildButton(v, 2, 1, "", "Меню", "mnu"))
		// if history.Meta.CurrentPage < history.Meta.PageCount || n+5 < history.Meta.TotalCount {
		// 	keyboard.AddButtons(*BuildButton(v, 2, 1, "", "->", "hist", strconv.Itoa(n+5)))
		// } else {
		// 	keyboard.AddButtons(*v.NewButton(2, 1, viber.None, "", "--", "", false))
		// }
		// msg.SetKeyboard(keyboard)
		_, err = v.SendMessage(user.ViberUser.ID, msg)
		check(err)
		fmt.Printf("%+v\n", history)
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
