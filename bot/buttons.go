package bot

import (
	"github.com/orsenkucher/viber"
	"github.com/sergTch/viberBotTest/data"
)

func BuildButton(v *viber.Viber, cols int, rows int, image string, text string, actionIDs ...string) *viber.Button {
	actBody := "#butt"
	for _, id := range actionIDs {
		actBody += "/" + id
	}
	return v.NewButton(cols, rows, viber.Reply, actBody, text, image, true)
}

func BuildCfgButton(v *viber.Viber, butt data.Butt, silent bool, actionIDs ...string) *viber.Button {
	actBody := "#butt"
	for _, id := range actionIDs {
		actBody += "/" + id
	}
	return v.NewButton(butt.Col, butt.Row, viber.Reply, actBody, butt.Text, butt.Image, silent)
}

func CfgButton(v *viber.Viber, actType viber.ActionType, butt data.Butt, action string, silent bool) *viber.Button {
	return v.NewButton(butt.Col, butt.Row, actType, action, butt.Text, butt.Image, silent)
}

func TxtCfgButton(v *viber.Viber, actType viber.ActionType, butt data.Butt, text string, action string, silent bool) *viber.Button {
	return v.NewButton(butt.Col, butt.Row, actType, action, text, butt.Image, silent)
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
	ButtActions["prf"] = &ButtAction{Act: ProfileChange, ID: "prf"}
	ButtActions["frq"] = &ButtAction{Act: FillRequired, ID: "frq"}
	ButtActions["sbq"] = &ButtAction{Act: ShowBarcode, ID: "sbq"}
	ButtActions["sbl"] = &ButtAction{Act: ShowBalance, ID: "sbl"}
}
