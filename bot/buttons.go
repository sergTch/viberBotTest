package bot

import "github.com/mileusna/viber"

var Actions map[string]*Action

type Action struct {
	Act func()
	ID  *string
}

func BuildButton(v *viber.Viber, actionIDs []string, image string, text string) *viber.Button {
	actBody := "#butt"
	for _, id := range actionIDs {
		actBody += "/" + id
	}
	return v.NewButton(1, 1, viber.Reply, actBody, "0", image)
}
