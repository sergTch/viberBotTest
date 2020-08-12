package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	pub = "secure/pub.json"
)

var Cfg struct {
	MinAge         int    `json:"Age"`
	ApiUrl         string `json:"ApiUrl"`
	AgreementLink  string `json:"AgreementLink"`
	CountryID      string `json:"CountryID"`
	AcceptLanguage string `json:"Accept-Language"`
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
}
