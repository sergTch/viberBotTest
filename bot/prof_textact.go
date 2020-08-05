package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

func ChangeField(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	var err error
	user := UserIDMap[u.ID]
	field, ok := UserField[u.ID]
	if !ok {
		Menu(v, u, m, token, t)
	}
	if abm.FieldType[field.FieldType] == "String" {
		field.Value = m.Text
	} else if abm.FieldType[field.FieldType] == "Integer" {
		field.Value, err = strconv.Atoi(m.Text)
		check(err)
	} else if abm.FieldType[field.FieldType] == "Birthday" {
		_, err := time.Parse("2006-01-02", m.Text)
		if err != nil {
			_, err := v.SendTextMessage(u.ID, "Дата должна быть выписана в формате ГГГГ-ММ-ДД. Повторите ещё раз")
			check(err)
			msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Введите дату в формате ГГГГ-ММ-ДД")
			keyboard := v.NewKeyboard("", false)
			keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
			msg.Keyboard = keyboard
			_, err = v.SendMessage(u.ID, msg)
			check(err)
			return
		}
		field.Value = m.Text
	}
	err = abm.Client.FieldSave(user.Token, field)
	if field.Key == "id_region" {
		check(err)
		msg := v.NewTextMessage("Редактируем '" + field.Name + "' Введите несколько первых букв вашего города")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: SearchCity}}
		UserField[u.ID] = field
		_, err = v.SendMessage(u.ID, msg)
		check(err)
		return
	}
	check(err)
	ProfileChange(v, u, m, token, t)
}

func SearchCity(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	field, ok := UserField[u.ID]
	if !ok {
		Menu(v, u, m, token, t)
	}

	user := UserIDMap[u.ID]
	prof, err := abm.Client.Profile(user.Token)
	check(err)
	region, err := abm.Client.GetRegion(fmt.Sprint(prof.Region.Value))
	check(err)
	fmt.Println(region)

	cities, err := abm.Client.SearchCity(m.Text)
	if len(cities) > 0 {
		msg := v.NewTextMessage("Редактируем '" + field.Name + "'" + ". Выберите свой вариант")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 3, 1, "", "Отмена", "prf"))
		for _, city := range cities {
			if city.RegionID == region.RegionID {
				keyboard.AddButtons(*v.NewButton(3, 1, viber.Reply, strconv.Itoa(city.CityID), city.CityName, "", true))
			}
		}
		keyboard.InputFieldState = viber.HiddenInputField
		msg.Keyboard = keyboard
		UserTxtAct[u.ID] = []*TextAction{{Act: ChangeField}}
		UserField[u.ID] = prof.City
		_, err = v.SendMessage(u.ID, msg)
		check(err)
	}
	check(err)
}
