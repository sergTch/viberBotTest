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
	Act func(v *viber.Viber, u viber.User, m *viber.TextMessage, token uint64, t time.Time)
}

func init() {

}

func Registration(v *viber.Viber, u viber.User, m *viber.TextMessage, token uint64, t time.Time) {
	if !strings.Contains(m.Text, " ") && len(m.Text) > 5 {
		user := UserIDMap[u.ID]
		_, err := abm.Client.Register(user.Contact.PhoneNumber, m.Text, "Ydfsdf464s")
		if err != nil {
			fmt.Println(err)
			_, err = v.SendTextMessage(u.ID, "Плохой пароль, попробуйте другой")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		_, err = v.SendTextMessage(u.ID, "Введите код из SMS. Будет отправлен вам в течении 2-х минут")
		if err != nil {
			fmt.Println(err)
		}
		UserTxtAct[u.ID] = []*TextAction{}
	} else {
		_, err := v.SendTextMessage(u.ID, "Введите другой пароль, он должен не содержать пробелов и состоять минимум из 6-ти символов")
		if err != nil {
			fmt.Println(err)
		}
	}
}
