package bot

import (
	"github.com/orsenkucher/viber"
)

func BuildButton(v *viber.Viber, cols int, rows int, image string, text string, actionIDs ...string) *viber.Button {
	actBody := "#butt"
	for _, id := range actionIDs {
		actBody += "/" + id
	}
	return v.NewButton(cols, rows, viber.None, actBody, text, image)
}
