package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	pub = "secure/pub.json"
)

var (
	MinAge        int
	ApiUrl        string
	AgreementLink string
)

func init() {
	var cfg map[string]interface{}
	bytes, err := ioutil.ReadFile(pub)
	if err != nil {
		panic(fmt.Errorf("Failed to read '%s': %w", pub, err))
	}

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(fmt.Errorf("Failed to parse '%s': %w", pub, err))
	}

	AgreementLink = cfg["AgreementLink"].(string)
	MinAge = int(cfg["Age"].(float64))
	ApiUrl = cfg["ApiUrl"].(string)
}
