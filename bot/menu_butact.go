package bot

import (
	"time"

	"github.com/orsenkucher/viber"
)

func Menu(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Покупки", "lop"))
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

}
