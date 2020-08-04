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
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "В меню", "mnu"))

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

func ChangeProfField(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, fkey string) {
	user := UserIDMap[u.ID]
	prof, err := abm.Client.Profile(user.Token)
	check(err)
	field, ok := prof.Main[fkey]
	if !ok {
		field = prof.Additional[fkey]
	}
	if fkey == "city" {
		field = prof.City
	}
	if fkey == "region" {
		field = prof.Region
	}
	if field == nil {
		return
	}
	if abm.DataType[field.DataType] == "Text" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Напишите новый вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		UserField[u.ID] = field
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Dropdown list" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Выберите свой вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
		for _, sch := range field.Schema {
			keyboard.AddButtons(*v.NewButton(6, 1, viber.Reply, sch.ID, sch.Value, "", true))
		}
		keyboard.InputFieldState = viber.HiddenInputField
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		UserField[u.ID] = field
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Date" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Введите дату в формате ГГГГ-ММ-ДД")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		UserField[u.ID] = field
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Checkbox" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Выберите свой вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
		keyboard.AddButtons(*v.NewButton(6, 1, viber.Reply, "1", "да", "", true))
		keyboard.AddButtons(*v.NewButton(6, 1, viber.Reply, "0", "нет", "", true))
		keyboard.InputFieldState = viber.HiddenInputField
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		UserField[u.ID] = field
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
}
