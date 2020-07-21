package bot

import (
	"fmt"
	"time"

	"github.com/orsenkucher/viber"
)

var TextActions map[string]*TextAction

type TextAction struct {
	Act func(v *viber.Viber, u viber.User, m viber.Message, token uint64, t time.Time)
	ID  string
}

func init() {
	fmt.Println("init2")
}
