package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/data"
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
		smsID, err := abm.Client.Register(user.PhoneNumber, m.Text)
		fmt.Println(smsID)
		if err != nil {
			fmt.Println(err)
			_, err = v.SendTextMessage(u.ID, data.Translate(user.Language, "Плохой пароль, попробуйте другой"))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		user.Password = m.Text
		_, err = v.SendTextMessage(u.ID, data.Translate(user.Language, "Введите код из SMS. Будет отправлен вам в течении 2-х минут"))
		check(err)
		UserSMS[u.ID] = SMS{ID: smsID, ConfirmType: "registration-confirm"}
		UserTxtAct[u.ID] = []*TextAction{{Act: SMSConfirm}}
	} else {
		_, err := v.SendTextMessage(u.ID, data.Translate("", "Введите другой пароль, он должен не содержать пробелов и состоять минимум из 6-ти символов"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func CheckPassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	uToken, err := abm.Client.AuthPhone(user.PhoneNumber, m.Text, func() { ReenterPassword(v, u.ID) })
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Println(err)
		msg := v.NewTextMessage(data.Translate(user.Language, "Неправильный пароль, попробуйте другой"))
		keyboard := v.NewKeyboard("", false)
		//keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Забыл Пароль", "chp"))
		keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.ForgotPass, true, "chp"))

		msg.Keyboard = keyboard
		_, err := v.SendMessage(u.ID, msg)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	user.Password = m.Text
	user.Token = uToken
	DB.Save(&user)
	Menu(v, u, m, token, t)
}

func ReadNewPassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	if !strings.Contains(m.Text, " ") && len(m.Text) > 5 {
		user := UserIDMap[u.ID]
		smsID, err := abm.Client.ChangePassword(user.PhoneNumber, m.Text)
		fmt.Println(smsID)
		if err != nil {
			fmt.Println(err)
			_, err = v.SendTextMessage(u.ID, data.Translate(user.Language, "Плохой пароль, попробуйте другой"))
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		user.Password = m.Text
		_, err = v.SendTextMessage(u.ID, data.Translate(user.Language, "Введите код из SMS. Будет отправлен вам в течении 2-х минут"))
		check(err)
		UserSMS[u.ID] = SMS{ID: smsID, ConfirmType: "change-password-confirm"}
		UserTxtAct[u.ID] = []*TextAction{{Act: SMSConfirm}}
	} else {
		_, err := v.SendTextMessage(u.ID, data.Translate("", "Введите другой пароль, он должен не содержать пробелов и состоять минимум из 6-ти символов"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SMSConfirm(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	sms, ok := UserSMS[u.ID]
	if !ok {
		text := data.Translate(user.Language, "Извините, что-то пошло не так, давайе попробуем ещё раз:(\nДля регистрации в программе лояльности придумайте и отправьте мне пароль. Пароль должен состоять минимум из 6-ти символов")
		_, err := v.SendTextMessage(u.ID, text)
		check(err)
		UserTxtAct[u.ID] = []*TextAction{{Act: SetPassword}}
		return
	}
	regToken, resp, err := abm.Client.Confirm(m.Text, sms.ID, sms.ConfirmType)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.Ok {
		user.Token = abm.NewSmartToken(abm.Client, regToken, user.PhoneNumber, user.Password, func() { ReenterPassword(v, u.ID) })
		DB.Save(&user)
		act := NextAction[user.ViberUser.ID]
		act.Act(v, u, m, token, t)
	} else {
		fmt.Println(resp.Err)
		_, err := v.SendTextMessage(u.ID, resp.Err)
		check(err)
	}
}

func SetCard(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	_, resp, err := abm.Client.SetCard(user.Token, m.Text)
	check(err)
	if resp.Ok {
		_, barcode, err := abm.Client.BarCode(user.Token)
		check(err)
		msg := v.NewPictureMessage("bar-code", barcode, "")
		msg.Text = ""
		_, err = v.SendMessage(u.ID, msg)
		check(err)
		FillInfQuestion(v, u, m, token, t)
	} else {
		_, err := v.SendTextMessage(u.ID, resp.Err)
		check(err)
		CardExistQuestion(v, u, m, token, t)
	}
}
