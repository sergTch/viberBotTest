package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/data"
)

var Actions map[string]*Action

type Action struct {
	Act func(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time)
	ID  string
}

func init() {
	Actions = map[string]*Action{}
	Actions["arg"] = &Action{Act: AgreementMsg, ID: "arg"}
	Actions["str"] = &Action{Act: AgreementMsg, ID: "str"}
}

//id: str
func StartMsg(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	startB := BuildButton(v, 6, 1, "", "СТАРТ", "agr", "qwe")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*startB)
	msg := v.NewTextMessage("Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ")
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	if err != nil {
		fmt.Println(err)
	}
}

//id: agr
func AgreementMsg(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time) {
	phoneB := v.NewButton(6, 1, viber.SharePhone, "", "Принять", "")
	linkB := v.NewButton(3, 1, viber.OpenURL, "", "Условия", data.AgreementLink)
	cancelB := BuildButton(v, 3, 1, "", "Отмена", "str")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*phoneB, *linkB, *cancelB)
	msg := v.NewTextMessage("Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ")
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	if err != nil {
		fmt.Println(err)
	}
}
