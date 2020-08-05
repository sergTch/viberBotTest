package bot

import (
	"fmt"
	"strconv"
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

	full := true
	for _, field := range prof.Fields {
		if field.Required && (field.Value == nil || field.Value == 0 || field.Value == "") {
			full = false
		}
	}

	if !full {
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Заполнить анкету", "frq"))
	} else {
		for _, field := range prof.Fields {
			if field.Key == "id_region" {
				keyboard.AddButtons(*BuildButton(v, 6, 1, "", prof.Region.Name+"/"+prof.City.Name, "prof", field.Key))
			} else if field.Key != "id_city" && field.Key != "mobile" {
				keyboard.AddButtons(*BuildButton(v, 6, 1, "", field.Name, "prof", field.Key))
			}
		}
	}

	keyboard.InputFieldState = viber.HiddenInputField
	msg.Keyboard = keyboard
	_, err = v.SendMessage(u.ID, msg)
	check(err)
}

func FillRequired(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	prof, err := abm.Client.Profile(user.Token)
	check(err)
	fields := []*abm.Field{}
	for _, field := range prof.Fields {
		if field.Key == "id_region" || field.Key == "id_city" || field.Key == "has_smartphone" {
			fmt.Println(field.Key, " ", field.Value)
			fmt.Println(field.Value == nil || fmt.Sprint(field.Value) == "0" || field.Value == "")
			fmt.Println(field.Required)
		}
		if field.Required && (field.Value == nil || fmt.Sprint(field.Value) == "0" || field.Value == "") {
			fields = append(fields, field)
		}
	}
	for _, field := range prof.Fields {
		if !field.Required && (field.Value == nil || field.Value == 0 || field.Value == "") {
			fields = append(fields, field)
		}
	}
	UserFields[u.ID] = fields
	if len(fields) != 0 {
		ChangeProfField(v, u, m, token, t, fields[0].Key)
	} else {
		_, err := v.SendTextMessage(u.ID, "Все обязательные поля уже заполнены")
		check(err)
		ProfileChange(v, u, m, token, t)
	}
}

func ChangeProfField(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, fkey string) {
	user := UserIDMap[u.ID]
	prof, err := abm.Client.Profile(user.Token)
	fields := UserFields[u.ID]
	check(err)
	field, ok := prof.Main[fkey]
	if !ok {
		field = prof.Additional[fkey]
	}
	if fkey == "id_region" {
		regions, err := abm.Client.Regions()
		check(err)
		msg := v.NewTextMessage("Редактируем '" + prof.Region.Name + "'" + ". Выберите свой вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Закончить позже", "prf"))
		for _, region := range regions {
			keyboard.AddButtons(*v.NewButton(3, 1, viber.Reply, strconv.Itoa(region.RegionID), region.RegionName, "", true))
		}
		keyboard.InputFieldState = viber.HiddenInputField
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		if len(fields) == 0 || fields[0].Key != prof.Region.Key {
			UserFields[u.ID] = []*abm.Field{prof.Region, prof.City}
		}
		_, err = v.SendMessage(u.ID, msg)
		check(err)
		return
	}
	if field == nil {
		return
	}
	if abm.DataType[field.DataType] == "Text" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Напишите новый вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Закончить позже", "prf"))
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		if len(fields) == 0 || fields[0].Key != field.Key {
			UserFields[u.ID] = []*abm.Field{field}
		}
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Dropdown list" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Выберите свой вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Закончить позже", "prf"))
		for _, sch := range field.Schema {
			keyboard.AddButtons(*v.NewButton(6, 1, viber.Reply, sch.ID, sch.Value, "", true))
		}
		keyboard.InputFieldState = viber.HiddenInputField
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		if len(fields) == 0 || fields[0].Key != field.Key {
			UserFields[u.ID] = []*abm.Field{field}
		}
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Date" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Введите дату в формате ГГГГ-ММ-ДД")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Закончить позже", "prf"))
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		if len(fields) == 0 || fields[0].Key != field.Key {
			UserFields[u.ID] = []*abm.Field{field}
		}
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Checkbox" {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Выберите свой вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Закончить позже", "prf"))
		keyboard.AddButtons(*v.NewButton(6, 1, viber.Reply, "1", "да", "", true))
		keyboard.AddButtons(*v.NewButton(6, 1, viber.Reply, "0", "нет", "", true))
		keyboard.InputFieldState = viber.HiddenInputField
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		if len(fields) == 0 || fields[0].Key != field.Key {
			UserFields[u.ID] = []*abm.Field{field}
		}
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
}
