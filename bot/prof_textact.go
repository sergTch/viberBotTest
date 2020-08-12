package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/data"
)

func ChangeField(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	var err error
	user := UserIDMap[u.ID]
	field, ok := UserFields[u.ID]
	if !ok {
		Menu(v, u, m, token, t)
	}
	if abm.FieldType[field[0].FieldType] == "String" {
		field[0].Value = m.Text
	} else if abm.FieldType[field[0].FieldType] == "Integer" {
		field[0].Value, err = strconv.Atoi(m.Text)
		check(err)
	} else if abm.FieldType[field[0].FieldType] == "Birthday" {
		date, err := time.Parse("2006-01-02", m.Text)
		if err != nil {
			_, err := v.SendTextMessage(u.ID, "Дата должна быть выписана в формате ГГГГ-ММ-ДД. Повторите ещё раз")
			check(err)
			msg := v.NewTextMessage("Редактируем '" + field[0].Name + "'" + ". Введите дату в формате ГГГГ-ММ-ДД")
			keyboard := v.NewKeyboard("", false)
			keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
			msg.Keyboard = keyboard
			_, err = v.SendMessage(u.ID, msg)
			check(err)
			return
		}
		if time.Since(date).Hours() < float64(data.Cfg.MinAge)*24*365.25 && field[0].Key == "birth_day" {
			_, err := v.SendTextMessage(u.ID, fmt.Sprintf("Минимальный возраст для регистрации %v лет.", data.Cfg.MinAge))
			check(err)
			msg := v.NewTextMessage("Редактируем '" + field[0].Name + "'" + ". Введите дату в формате ГГГГ-ММ-ДД")
			keyboard := v.NewKeyboard("", false)
			keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
			msg.Keyboard = keyboard
			_, err = v.SendMessage(u.ID, msg)
			check(err)
			return
		}
		field[0].Value = m.Text
	}
	fmt.Println(field[0].Name, field[0].Value)
	err = abm.Client.FieldSave(user.Token, field[0])
	check(err)
	if field[0].Key == "id_region" {
		if len(field) <= 1 {
			ProfileChange(v, u, m, token, t)
			return
		}
		check(err)
		msg := v.NewTextMessage("Редактируем '" + field[1].Name + "' Введите несколько первых букв вашего города")
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Отмена", "prf"))
		msg.Keyboard = keyboard
		UserFields[u.ID] = field[1:]
		UserTxtAct[u.ID] = []*TextAction{{Act: SearchCity}}
		_, err = v.SendMessage(u.ID, msg)
		check(err)
		return
	}
	UserFields[u.ID] = field[1:]
	if len(field) == 1 {
		ProfileChange(v, u, m, token, t)
	} else {
		ChangeProfField(v, u, m, token, t, field[1].Key)
	}
}

func SearchCity(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	prof, err := abm.Client.Profile(user.Token)
	checkServerError(err, v, u, m, token, t)
	if err != nil {
		return
	}
	region, err := abm.Client.GetRegion(fmt.Sprint(prof.Region.Value))
	checkServerError(err, v, u, m, token, t)
	if err != nil {
		return
	}
	check(err)
	fmt.Println(region)

	cities, err := abm.Client.SearchCity(m.Text)
	if len(cities) > 0 {
		field, ok := UserFields[u.ID]
		fmt.Println("red city ", ok, field[0])
		msg := v.NewTextMessage("Редактируем '" + prof.City.Name + "'" + ". Выберите свой вариант")
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
		_, err = v.SendMessage(u.ID, msg)
		check(err)
	}
	check(err)
}
