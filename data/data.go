package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	pub        = "secure/pub.json"
	buttonspub = "secure/butt.json"
)

var Cfg struct {
	MinAge         int    `json:"Age"`
	ApiUrl         string `json:"ApiUrl"`
	AgreementLink  string `json:"AgreementLink"`
	CountryID      string `json:"CountryID"`
	AcceptLanguage string `json:"Accept-Language"`
	Currency       string `json:"Currency"`
}

type Butt struct {
	Col   int    `json:"col"`
	Row   int    `json:"row"`
	Image string `json:"image"`
	Text  string `json:"text"`
}

var ButtCfg struct {
	Start        Butt `json:"start"`
	Agree        Butt `json:"agree"`
	Back         Butt `json:"back"`
	Cancel       Butt `json:"cancel"`
	Agreement    Butt `json:"agreement"`
	ForgotPass   Butt `json:"forgot_pass"`
	Yes          Butt `json:"yes"`
	No           Butt `json:"no"`
	EnterCard    Butt `json:"enter_card"`
	EnterPass    Butt `json:"enter_pass"`
	EnterNewPass Butt `json:"enter_new_pass"`
	NoCard       Butt `json:"no_card"`
	FillRequired Butt `json:"fill_required"`
	ProfField    Butt `json:"prof_field"`
	DropDown     Butt `json:"drop_down"`
	Menu         Butt `json:"menu"`
	FinishLater  Butt `json:"finish_later"`
	Region       Butt `json:"region"`
	City         Butt `json:"city"`
}

func Init() {
	bytes, err := ioutil.ReadFile(pub)
	if err != nil {
		panic(fmt.Errorf("Failed to read '%s': %w", pub, err))
	}

	err = json.Unmarshal(bytes, &Cfg)
	if err != nil {
		panic(fmt.Errorf("Failed to parse '%s': %w", pub, err))
	}

	bytes, err = ioutil.ReadFile(buttonspub)
	if err != nil {
		panic(fmt.Errorf("Failed to read '%s': %w", buttonspub, err))
	}

	err = json.Unmarshal(bytes, &ButtCfg)
	if err != nil {
		panic(fmt.Errorf("Failed to parse '%s': %w", buttonspub, err))
	}
}
