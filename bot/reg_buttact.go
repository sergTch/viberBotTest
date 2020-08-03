package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/abm"
	"github.com/sergTch/viberBotTest/data"
)

var ButtActions map[string]*ButtAction

type ButtAction struct {
	Act func(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time)
	ID  string
}

func init() {
	ButtActions = map[string]*ButtAction{}
	ButtActions["agr"] = &ButtAction{Act: AgreementMsg, ID: "agr"}
	ButtActions["str"] = &ButtAction{Act: StartMsg, ID: "str"}
	ButtActions["ceq"] = &ButtAction{Act: CardExistQuestion, ID: "ceq"}
	ButtActions["ccr"] = &ButtAction{Act: CardCreate, ID: "ccr"}
	ButtActions["cin"] = &ButtAction{Act: CardInput, ID: "cin"}
	ButtActions["mnu"] = &ButtAction{Act: Menu, ID: "mnu"}
	ButtActions["chp"] = &ButtAction{Act: ChangePassword, ID: "chp"}
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
	keyboard.AddButtons(*BuildButton(v, 3, 1, "", "Да", "cin"), *BuildButton(v, 3, 1, "", "Нет", "ccr"))
	keyboard.InputFieldState = viber.HiddenInputField
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	check(err)
	UserTxtAct[u.ID] = []*TextAction{}
}

func CardInput(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage("Введите номер вашей карты")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*v.NewButton(6, 1, viber.None, "", "Вводите номер карты", "", true))
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	UserTxtAct[u.ID] = []*TextAction{{Act: SetCard}}
	check(err)
}

func CardCreate(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	user := UserIDMap[u.ID]
	_, barcode, err := abm.Client.BarCode(user.Token)
	check(err)
	msg := v.NewPictureMessage("bar-code", barcode, "")
	_, err = v.SendMessage(u.ID, msg)
	check(err)
	FillInfQuestion(v, u, m, token, t)
}

func FillInfQuestion(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage("Для возможности использовать накопленные бонусы необходимо заполнить обязательные поля анкеты в разделе меню Мой профиль. Заполнить анкету сейчас?")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 3, 1, "", "да"), *BuildButton(v, 3, 1, "", "Заполнить позже", "mnu"))
	keyboard.InputFieldState = viber.HiddenInputField
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	UserTxtAct[u.ID] = []*TextAction{}
	check(err)
}

func EnterPassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage("Ваш номер уже зарегистрирован в программе лояльности. Для авторизации отправьте свой пароль")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Забыл Пароль", "chp"))
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	UserTxtAct[u.ID] = []*TextAction{{Act: CheckPassword}}
	check(err)
}

func ChangePassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage("Придумайте и отправьте мне новый пароль. Пароль должен состоять минимум из 6-ти символов")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*v.NewButton(6, 1, viber.None, "", "Вводите новый пароль", "", true))
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	UserTxtAct[u.ID] = []*TextAction{{Act: ReadNewPassword}}
	check(err)
}

func Menu(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Go to start", "str"))
	keyboard.InputFieldState = viber.HiddenInputField
	UserTxtAct[u.ID] = []*TextAction{}
	msg := v.NewTextMessage("Вы попали в меню...")
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	check(err)

	user := UserIDMap[u.ID]
	check(err)

	profile := abm.Client.NewProfile2(user.Token)
	for key, field := range profile.Additional {
		fmt.Println(key, field)
	}
}
