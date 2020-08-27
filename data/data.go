package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	pub          = "%s/pub.json"
	buttonspub   = "%s/butt.json"
	translations = "%s/lang.json"
	Viber        = "%s/viber.json"
	Gorm         = "%s/gorm.json"
)

var Cfg struct {
	MinAge         int    `json:"Age"`
	ApiUrl         string `json:"ApiUrl"`
	Webhook        string `json:"Webhook"`
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

func Init(path string) {
	pub = fmt.Sprintf(pub, path)
	buttonspub = fmt.Sprintf(buttonspub, path)
	translations = fmt.Sprintf(translations, path)
	Viber = fmt.Sprintf(Viber, path)
	Gorm = fmt.Sprintf(Gorm, path)

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

	bytes, err = ioutil.ReadFile(translations)
	if err != nil {
		panic(fmt.Errorf("Failed to read '%s': %w", translations, err))
	}

	err = json.Unmarshal(bytes, &Translations)
	if err != nil {
		panic(fmt.Errorf("Failed to parse '%s': %w", translations, err))
	}
}

var Translations map[string]map[string]string

func Translate(lang, text string) string {
	if lang == "" {
		lang = Cfg.AcceptLanguage
	}

	if t, ok := translate(lang, text); ok {
		return t
	}

	return text
}

func translate(lang, text string) (string, bool) {
	if ts, ok := Translations[text]; ok {
		if t, ok := ts[lang]; ok {
			return t, true
		}
	}

	return "", false
}

func Format(format string, args ...interface{}) string {
	args2 := make([]string, len(args))
	for i, v := range args {
		if i%2 == 0 {
			args2[i] = fmt.Sprintf("{%v}", v)
		} else {
			args2[i] = fmt.Sprint(v)
		}
	}
	r := strings.NewReplacer(args2...)
	return r.Replace(format)
}
