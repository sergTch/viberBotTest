package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/data"
)

var ButtActions map[string]*ButtAction

type ButtAction struct {
	Act func(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time)
	ID  string
}

func init() {
	ButtActions = map[string]*ButtAction{}
	ButtActions["agr"] = &ButtAction{Act: AgreementMsg, ID: "agr"}
	ButtActions["str"] = &ButtAction{Act: StartMsg, ID: "str"}
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
	linkB := v.NewButton(6, 1, viber.OpenURL, data.AgreementLink, "Условия", "")
	phoneB := v.NewButton(3, 1, viber.SharePhone, "", "Принять", "")
	cancelB := BuildButton(v, 3, 1, "", "Отмена", "str")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*linkB, *phoneB, *cancelB)
	msg := v.NewTextMessage(fmt.Sprint("Вам уже исполнилось ", data.Age, " лет и Вы принимаете Условия программы лояльности?"))
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	if err != nil {
		fmt.Println(err)
	}
}
