package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
)

var UserTxtAct map[string][]*TextAction

type TextAction struct {
	Act func(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time)
}

func init() {
	UserTxtAct = map[string][]*TextAction{}
}

func Registration(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	if !strings.Contains(m.Text, " ") && len(m.Text) > 5 {
		user := UserIDMap[u.ID]
		smsID, err := abm.Client.Register(user.Contact.PhoneNumber, m.Text, "Ydfsdf464s")
		if err != nil {
			fmt.Println(err)
			_, err = v.SendTextMessage(u.ID, "Плохой пароль, попробуйте другой")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		_, err = v.SendTextMessage(u.ID, "Введите код из SMS. Будет отправлен вам в течении 2-х минут")
		check(err)
		UserData[u.ID] = SMS{ID: smsID}
		UserTxtAct[u.ID] = []*TextAction{{Act: RegistrationConfirm}}
	} else {
		_, err := v.SendTextMessage(u.ID, "Введите другой пароль, он должен не содержать пробелов и состоять минимум из 6-ти символов")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func RegistrationConfirm(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	//user := UserIDMap[u.ID]
	data, ok := UserData[u.ID].(SMS)
	if !ok {
		_, err := v.SendTextMessage(u.ID, "Извините, что-то пошло не так, давайе попробуем ещё раз:(\nДля регистрации в программе лояльности придумайте и отправьте мне пароль. Пароль должен состоять минимум из 6-ти символов")
		check(err)
		UserTxtAct[u.ID] = []*TextAction{{Act: Registration}}
		return
	}
	_, ok, err := abm.Client.Confirm(m.Text, data.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if ok {
		_, err := v.SendTextMessage(u.ID, "У вас уже есть бонусная карта?")
		check(err)
		UserTxtAct[u.ID] = []*TextAction{}
	} else {
		fmt.Println("")
	}
}
