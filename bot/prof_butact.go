package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/data"
)

func ProfileChange(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	_, err := v.SendTextMessage(u.ID, data.Translate(user.Language, "Для использования бонусов, нужно заполнить обязательные поля (* - обязательно к заполнению)"))
	check(err)
	prof, err := abm.Client.Profile(user.Token)
	// if user == nil {
	// 	panic("panica: user was nil")
	// }
	checkServerError(err, v, u, m, token, t)
	if err != nil {
		return
	}
	text := data.Format(data.Translate(user.Language, "*Номер: {phone_number}\n"), "phone_number", user.PhoneNumber)
	text += prof.ToString()
	msg := v.NewTextMessage(text)
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.Menu, true, "mnu"))

	full := true
	for _, field := range prof.Fields {
		if field.Required && (field.Value == nil || field.Value == 0 || field.Value == "") {
			full = false
		}
	}

	if !full {
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FillRequired, true, "frq"))
	} else {
		for _, field := range prof.Fields {
			if field.Key == "id_region" {
				keyboard.AddButtons(*TxtBuildCfgButton(v, data.ButtCfg.ProfField, prof.Region.Name+"/"+prof.City.Name, true, "prof", field.Key))
			} else if field.Key != "id_city" && field.Key != "mobile" {
				keyboard.AddButtons(*TxtBuildCfgButton(v, data.ButtCfg.ProfField, field.Name, true, "prof", field.Key))
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
	checkServerError(err, v, u, m, token, t)
	if err != nil {
		return
	}
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
		_, err := v.SendTextMessage(u.ID, data.Translate(user.Language, "Все обязательные поля уже заполнены"))
		check(err)
		ProfileChange(v, u, m, token, t)
	}
}

func ChangeProfField(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time, fkey string) {
	fmt.Println("Change field: ", fkey)
	fmt.Println(abm.DataType)
	fmt.Println(abm.FieldType)
	user := UserIDMap[u.ID]
	prof, err := abm.Client.Profile(user.Token)
	checkServerError(err, v, u, m, token, t)
	if err != nil {
		return
	}
	fields := UserFields[u.ID]
	field, ok := prof.Main[fkey]
	if !ok {
		field = prof.Additional[fkey]
	}
	if fkey == "id_region" {
		regions, err := abm.Client.Regions()
		check(err)
		text := data.Format(data.Translate(user.Language, "Редактируем '{region_name}'. Выберите свой вариант"), "region_name", prof.Region.Name)
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
		for _, region := range regions {
			keyboard.AddButtons(*TxtCfgButton(v, viber.Reply, data.ButtCfg.Region, strconv.Itoa(region.RegionID), region.RegionName, true))
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
		text := data.Format(data.Translate(user.Language, "Редактируем '{field_name}'. Напишите новый вариант"), "field_name", field.Name)
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		if len(fields) == 0 || fields[0].Key != field.Key {
			UserFields[u.ID] = []*abm.Field{field}
		}
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Dropdown list" {
		text := data.Format(data.Translate(user.Language, "Редактируем '{field_name}'. Выберите свой вариант"), "field_name", field.Name)
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
		for _, sch := range field.Schema {
			keyboard.AddButtons(*TxtCfgButton(v, viber.Reply, data.ButtCfg.DropDown, sch.ID, sch.Value, true))
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
		text := data.Format(data.Translate(user.Language, "Редактируем '{field_name}'. Введите дату в формате ГГГГ-ММ-ДД"), "field_name", field.Name)
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		if len(fields) == 0 || fields[0].Key != field.Key {
			UserFields[u.ID] = []*abm.Field{field}
		}
		_, err := v.SendMessage(u.ID, msg)
		check(err)
	}
	if abm.DataType[field.DataType] == "Checkbox" {
		text := data.Format(data.Translate(user.Language, "Редактируем '{field_name}'. Выберите свой вариант"), "field_name", field.Name)
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
		keyboard.AddButtons(*CfgButton(v, viber.Reply, data.ButtCfg.Yes, "1", true))
		keyboard.AddButtons(*CfgButton(v, viber.Reply, data.ButtCfg.No, "0", true))
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
