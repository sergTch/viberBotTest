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

func SetPassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	if !strings.Contains(m.Text, " ") && len(m.Text) > 5 {
		user := UserIDMap[u.ID]
		smsID, err := abm.Client.Register(user.PhoneNumber, m.Text, "Ydfsdf464s")
		fmt.Println(smsID)
		if err != nil {
			fmt.Println(err)
			_, err = v.SendTextMessage(u.ID, "Плохой пароль, попробуйте другой")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		user.Password = m.Text
		_, err = v.SendTextMessage(u.ID, "Введите код из SMS. Будет отправлен вам в течении 2-х минут")
		check(err)
		UserSMS[u.ID] = SMS{ID: smsID, ConfirmType: "registration-confirm"}
		UserTxtAct[u.ID] = []*TextAction{{Act: SMSConfirm}}
	} else {
		_, err := v.SendTextMessage(u.ID, "Введите другой пароль, он должен не содержать пробелов и состоять минимум из 6-ти символов")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func CheckPassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	uToken, err := abm.Client.AuthPhone(user.PhoneNumber, m.Text, "Ydfsdf464s")
	if err != nil {
		fmt.Println(err)
		_, err = v.SendTextMessage(u.ID, "Неправильный пароль, попробуйте другой")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	user.Password = m.Text
	user.Token = uToken
	Menu(v, u, m, token, t)
}

func ReadNewPassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	if !strings.Contains(m.Text, " ") && len(m.Text) > 5 {
		user := UserIDMap[u.ID]
		smsID, err := abm.Client.ChangePassword(user.PhoneNumber, m.Text, "Ydfsdf464s")
		fmt.Println(smsID)
		if err != nil {
			fmt.Println(err)
			_, err = v.SendTextMessage(u.ID, "Плохой пароль, попробуйте другой")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		user.Password = m.Text
		_, err = v.SendTextMessage(u.ID, "Введите код из SMS. Будет отправлен вам в течении 2-х минут")
		check(err)
		UserSMS[u.ID] = SMS{ID: smsID, ConfirmType: "change-password-confirm"}
		UserTxtAct[u.ID] = []*TextAction{{Act: SMSConfirm}}
	} else {
		_, err := v.SendTextMessage(u.ID, "Введите другой пароль, он должен не содержать пробелов и состоять минимум из 6-ти символов")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SMSConfirm(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	sms, ok := UserSMS[u.ID]
	if !ok {
		_, err := v.SendTextMessage(u.ID, "Извините, что-то пошло не так, давайе попробуем ещё раз:(\nДля регистрации в программе лояльности придумайте и отправьте мне пароль. Пароль должен состоять минимум из 6-ти символов")
		check(err)
		UserTxtAct[u.ID] = []*TextAction{{Act: SetPassword}}
		return
	}
	regToken, ok, err := abm.Client.Confirm(m.Text, sms.ID, sms.ConfirmType)
	if err != nil {
		fmt.Println(err)
		return
	}
	if ok {
		user.Token = regToken
		act := NextAction[user.ViberUser.ID]
		act.Act(v, u, m, token, t)
	} else {
		fmt.Println("")
	}
}

func SetCard(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	_, ok, err := abm.Client.SetCard(user.Token, m.Text)
	check(err)
	if ok {
		_, barcode, err := abm.Client.BarCode(user.Token)
		check(err)
		msg := v.NewPictureMessage("bar-code", barcode, "")
		msg.Text = ""
		_, err = v.SendMessage(u.ID, msg)
		check(err)
		FillInfQuestion(v, u, m, token, t)
	} else {
		_, err := v.SendTextMessage(u.ID, "That card is invalid, sorry, you can answer no for following question to get newone")
		check(err)
		CardExistQuestion(v, u, m, token, t)
	}
}
