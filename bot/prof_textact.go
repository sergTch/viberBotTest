package bot

import (
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
		return
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
	check(err)
	ProfileChange(v, u, m, token, t)
}
