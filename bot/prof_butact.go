package bot

import (
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func ProfileChange(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	_, err := v.SendTextMessage(u.ID, "Для использования бонусов, нужно заполнить обязательные поля (* - обязательно к заполнению)")
	check(err)
	user := UserIDMap[u.ID]
	prof, err := abm.Client.Profile(user.Token)
	check(err)
	text := "*Номер: " + user.PhoneNumber + "\n"
	text += prof.ToString()
	msg := v.NewTextMessage(text)
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*v.NewButton(6, 1, viber.None, "", "Вводите номер карты", "", true))

	for _, field := range prof.Main {
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", field.Name, "prof", field.Key))
	}
	for _, field := range prof.Additional {
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", field.Name, "prof", field.Key))
	}

	keyboard.InputFieldState = viber.HiddenInputField
	msg.Keyboard = keyboard
	_, err = v.SendMessage(u.ID, msg)
	check(err)
}
