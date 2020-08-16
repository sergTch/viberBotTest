package bot

import (
	"github.com/orsenkucher/viber"
)

func BuildButton(v *viber.Viber, cols int, rows int, image string, text string, actionIDs ...string) *viber.Button {
	actBody := "#butt"
	for _, id := range actionIDs {
		actBody += "/" + id
	}
	return v.NewButton(cols, rows, viber.Reply, actBody, text, image, true)
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
	//ButtActions["lop"] = &ButtAction{Act: LastOperations, ID: "lop"}
	ButtActions["sbq"] = &ButtAction{Act: ShowBarcode, ID: "sbq"}
}
