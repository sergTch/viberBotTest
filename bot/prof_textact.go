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
			_, err := v.SendTextMessage(u.ID, data.Translate(user.Language, "Дата должна быть выписана в формате ГГГГ-ММ-ДД. Повторите ещё раз"))
			check(err)
			text := data.Format(data.Translate(user.Language, "Редактируем '{field_name}'. Введите дату в формате ГГГГ-ММ-ДД"), "field_name", field[0].Name)
			msg := v.NewTextMessage(text)
			keyboard := v.NewKeyboard("", false)
			keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
			msg.Keyboard = keyboard
			_, err = v.SendMessage(u.ID, msg)
			check(err)
			return
		}
		if time.Since(date).Hours() < float64(data.Cfg.MinAge)*24*365.25 && field[0].Key == "birth_day" {
			text := data.Format(data.Translate(user.Language, "Минимальный возраст для регистрации {min_age} лет."), "min_age", data.Cfg.MinAge)
			_, err := v.SendTextMessage(u.ID, text)
			check(err)
			text = data.Format(data.Translate(user.Language, "Редактируем '{field_name}'. Введите дату в формате ГГГГ-ММ-ДД"), "field_name", field[0].Name)
			msg := v.NewTextMessage(text)
			keyboard := v.NewKeyboard("", false)
			keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
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
		text := data.Format(data.Translate(user.Language, "Редактируем '{field_name}'. Введите несколько первых букв вашего города"), "field_name", field[1].Name)
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
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
		text := data.Format(data.Translate(user.Language, "Редактируем '{city_name}'. Выберите свой вариант"), "city_name", prof.City.Name)
		msg := v.NewTextMessage(text)
		keyboard := v.NewKeyboard("", false)
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.FinishLater, true, "prf"))
		for _, city := range cities {
			if city.RegionID == region.RegionID {
				keyboard.AddButtons(*TxtCfgButton(v, viber.Reply, data.ButtCfg.Region, strconv.Itoa(city.CityID), city.CityName, true))
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
