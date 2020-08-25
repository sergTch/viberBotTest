package bot

import (
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/data"
)

func F(v *viber.Viber, user User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage(data.Translate("", "Для возможности использовать накопленные бонусы необходимо заполнить обязательные поля анкеты в разделе меню Мой профиль. Заполнить анкету сейчас?"))
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 3, 1, "", data.Translate("", "Заполнить позже"), "mnu"))
	keyboard.InputFieldState = viber.HiddenInputField
	msg.Keyboard = keyboard
}
