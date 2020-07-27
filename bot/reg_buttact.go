package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/data"
)

var ButtActIDs map[string]*ButtAction

type ButtAction struct {
	Act func(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time)
	ID  string
}

func init() {
	ButtActIDs = map[string]*ButtAction{}
	ButtActIDs["agr"] = &ButtAction{Act: AgreementMsg, ID: "agr"}
	ButtActIDs["str"] = &ButtAction{Act: StartMsg, ID: "str"}
	ButtActIDs["ceq"] = &ButtAction{Act: CardExistQuestion, ID: "ceq"}
	ButtActIDs["ccr"] = &ButtAction{Act: CardCreate, ID: "ccr"}
	ButtActIDs["cin"] = &ButtAction{Act: CardInput, ID: "cin"}
}

//id: str
func StartMsg(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	startB := BuildButton(v, 6, 1, "", "СТАРТ", "agr", "qwe")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*startB)
	keyboard.InputFieldState = viber.HiddenInputField
	UserTxtAct[u.ID] = []*TextAction{}
	msg := v.NewTextMessage("Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ")
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	check(err)
}

//id: agr
func AgreementMsg(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	linkB := v.NewButton(6, 1, viber.OpenURL, data.AgreementLink, "Условия", "", true)
	phoneB := v.NewButton(3, 1, viber.SharePhone, "qwe", "Принять", "", true)
	cancelB := BuildButton(v, 3, 1, "", "Отмена", "str")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*linkB, *phoneB, *cancelB)
	keyboard.InputFieldState = viber.HiddenInputField
	UserTxtAct[u.ID] = []*TextAction{}
	msg := v.NewTextMessage(fmt.Sprint("Вам уже исполнилось ", data.MinAge, " лет и Вы принимаете Условия программы лояльности?"))
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	check(err)
}

func CardExistQuestion(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage("У вас уже есть бонусная карта?")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 3, 1, "", "Да", "cin"), *BuildButton(v, 3, 1, "", "Да", "ccr"))
	keyboard.InputFieldState = viber.HiddenInputField
	_, err := v.SendMessage(u.ID, msg)
	check(err)
	UserTxtAct[u.ID] = []*TextAction{}
}

func CardInput(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	_, err := v.SendTextMessage(u.ID, "Введите номер вашей карты")
	UserTxtAct[u.ID] = []*TextAction{{Act: BarCode}}
	check(err)
}

func CardCreate(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	_, err := v.SendTextMessage(u.ID, "Введите номер вашей карты")
	check(err)
}
