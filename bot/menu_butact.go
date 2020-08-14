package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func Menu(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Покупки", "lop"))
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

func LastOperations(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	if user, ok := UserIDMap[u.ID]; ok {
		fmt.Println(user.Token.Token())
		history, err := abm.Client.ClientHistory(user.Token, 1)
		check(err)
		msg := v.NewRichMediaMessage(2, 2, "")
		msg.AddButton(v.NewButton(6, 6, viber.None, "", "first", "", true))
		msg.AddButton(v.NewButton(6, 1, viber.None, "", "last", "", true))
		msg.AddButton(v.NewButton(6, 3, viber.None, "", "1-3", "", true))
		msg.AddButton(v.NewButton(6, 2, viber.None, "", "4-5", "", true))
		msg.AddButton(v.NewButton(6, 1, viber.None, "", "6", "", true))
		msg.AddButton(v.NewButton(6, 1, viber.None, "", "7", "", true))
		_, err = v.SendMessage(u.ID, msg)
		check(err)
		Menu(v, u, m, token, t)
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
