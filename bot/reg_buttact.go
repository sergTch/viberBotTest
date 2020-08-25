package bot

import (
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

//id: str
func StartMsg(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	//startB := BuildButton(v, 6, 1, "", "СТАРТ", "agr")
	startB := BuildCfgButton(v, data.ButtCfg.Start, true, "agr")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*startB)
	keyboard.InputFieldState = viber.HiddenInputField
	UserTxtAct[u.ID] = []*TextAction{}
	msg := v.NewTextMessage(data.Translate("", "Приветствуем в програме лояльности ABMLoyalty! Для начала работы нажмите СТАРТ"))
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	check(err)
}

//id: agr
func AgreementMsg(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	//linkB := v.NewButton(6, 1, viber.OpenURL, data.Cfg.AgreementLink, "Условия", "", true)
	linkB := CfgButton(v, viber.OpenURL, data.ButtCfg.Agreement, data.Cfg.AgreementLink, true)
	//phoneB := v.NewButton(3, 1, viber.SharePhone, "qwe", "Принять", "", true)
	phoneB := CfgButton(v, viber.SharePhone, data.ButtCfg.Agree, "qwe", true)
	// cancelB := BuildButton(v, 3, 1, "", "Отмена", "str")
	cancelB := BuildCfgButton(v, data.ButtCfg.Back, true, "str")
	keyboard := v.NewKeyboard("", false)
	keyboard.AddButtons(*linkB, *phoneB, *cancelB)
	keyboard.InputFieldState = viber.HiddenInputField
	UserTxtAct[u.ID] = []*TextAction{}
	text := data.Format(data.Translate("", "Вам уже исполнилось {min_age} лет и Вы принимаете Условия программы лояльности?"), "min_age", data.Cfg.MinAge)
	msg := v.NewTextMessage(text)
	msg.SetKeyboard(keyboard)
	_, err := v.SendMessage(u.ID, msg)
	check(err)
}

func CardExistQuestion(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage(data.Translate("", "У вас уже есть бонусная карта?"))
	keyboard := v.NewKeyboard("", false)
	//keyboard.AddButtons(*BuildButton(v, 3, 1, "", "Да", "cin"), *BuildButton(v, 3, 1, "", "Нет", "ccr"))
	keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.Yes, true, "cin"), *BuildCfgButton(v, data.ButtCfg.No, true, "ccr"))
	keyboard.InputFieldState = viber.HiddenInputField
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	check(err)
	UserTxtAct[u.ID] = []*TextAction{}
}

func CardInput(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage(data.Translate("", "Введите номер вашей карты"))
	keyboard := v.NewKeyboard("", false)
	// keyboard.AddButtons(*v.NewButton(6, 1, viber.None, "", "Вводите номер карты", "", true))
	keyboard.AddButtons(*CfgButton(v, viber.None, data.ButtCfg.EnterCard, "", true), *BuildCfgButton(v, data.ButtCfg.NoCard, true, "ccr"))
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
	text := data.Translate("", "Для возможности использовать накопленные бонусы необходимо заполнить обязательные поля анкеты в разделе меню Мой профиль. Заполнить анкету сейчас?")
	msg := v.NewTextMessage(text)
	keyboard := v.NewKeyboard("", false)
	//keyboard.AddButtons(*BuildButton(v, 3, 1, "", "да", "prf"), *BuildButton(v, 3, 1, "", "Заполнить позже", "mnu"))
	keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.Yes, true, "prf"), *BuildCfgButton(v, data.ButtCfg.No, true, "mnu"))
	keyboard.InputFieldState = viber.HiddenInputField
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	UserTxtAct[u.ID] = []*TextAction{}
	check(err)
}

func EnterPassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage(data.Translate("", "Ваш номер уже зарегистрирован в программе лояльности. Для авторизации отправьте свой пароль"))
	keyboard := v.NewKeyboard("", false)
	// keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Забыл Пароль", "chp"))
	keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.ForgotPass, true, "chp"))
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	UserTxtAct[u.ID] = []*TextAction{{Act: CheckPassword}}
	check(err)
}

func ReenterPassword(v *viber.Viber, uid string) {
	msg := v.NewTextMessage(data.Translate("", "Введите пароль, он возможно был изменён"))
	keyboard := v.NewKeyboard("", false)
	// keyboard.AddButtons(*BuildButton(v, 6, 1, "", "Забыл Пароль", "chp"))
	keyboard.AddButtons(*BuildCfgButton(v, data.ButtCfg.ForgotPass, true, "chp"))
	msg.Keyboard = keyboard
	_, err := v.SendMessage(uid, msg)
	UserTxtAct[uid] = []*TextAction{{Act: CheckPassword}}
	check(err)
}

func ChangePassword(v *viber.Viber, u viber.User, m viber.TextMessage, token uint64, t time.Time) {
	msg := v.NewTextMessage(data.Translate("", "Придумайте и отправьте мне новый пароль. Пароль должен состоять минимум из 6-ти символов"))
	keyboard := v.NewKeyboard("", false)
	// keyboard.AddButtons(*v.NewButton(6, 1, viber.None, "", "Вводите новый пароль", "", true))
	keyboard.AddButtons(*CfgButton(v, viber.None, data.ButtCfg.EnterNewPass, "", true))
	msg.Keyboard = keyboard
	_, err := v.SendMessage(u.ID, msg)
	UserTxtAct[u.ID] = []*TextAction{{Act: ReadNewPassword}}
	check(err)
}
