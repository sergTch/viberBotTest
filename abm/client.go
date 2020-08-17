package abm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sergTch/viberBotTest/data"
)

var Client *client

func Init() {
	Client = New()
}

type client struct {
	client *http.Client
	apiURL string
}

func New() *client {
	return &client{
		client: &http.Client{},
		apiURL: data.Cfg.ApiUrl,
	}
}

func (c *client) url(endpoint string) string {
	return c.apiURL + endpoint
}

func (c *client) PostForm(url string, vals url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(vals.Encode()),
	)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return c.Do(req)
}

func (c *client) Do(req *http.Request) (*http.Response, error) {
	setAcceptLanguage(req)
	return c.client.Do(req)
}

func setAcceptLanguage(req *http.Request) {
	req.Header.Add("Accept-Language", data.Cfg.AcceptLanguage)
}

func (c *client) CheckPhone(number string) (ok bool, err error) {
	values := url.Values{}
	values.Set("phone", number)

	r, err := c.PostForm(c.url("/v2/client/check-phone"), values)
	if err != nil {
		return
	}

	defer r.Body.Close()

	if r.StatusCode != 201 {
		err = errors.New("Not 201 status")
		return
	}

	var resp struct {
		Data struct {
			Is_exist bool `json:"is_exist"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	ok = resp.Data.Is_exist
	return
}

func (c *client) Register(phone, password string) (smsID int, err error) {
	values := url.Values{}
	values.Set("phone", phone)
	values.Set("password", password)

	r, err := c.PostForm(c.url("/v2.1/client/registration"), values)
	if err != nil {
		return
	}

	defer r.Body.Close()

	if r.StatusCode != 201 {
		err = errors.New("Not 201 status")
		return
	}

	var resp struct {
		Data struct {
			Phone   string `json:"phone"`
			SMSID   int    `json:"sms_id"`
			Timeout int    `json:"timeout"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	smsID = resp.Data.SMSID
	return
}

type SmartToken struct {
	client *client

	token, phone, password string
}

func NewSmartToken(client *client, token, phone, password string) *SmartToken {
	return &SmartToken{
		client:   client,
		token:    token,
		phone:    phone,
		password: password,
	}
}

func (s *SmartToken) Token() string {
	return s.token
}

func (s *SmartToken) Renew() (token *SmartToken, err error) {
	fmt.Println("*** RENEWING TOKEN ***")
	fmt.Printf("params: %v, %v\n", s.phone, s.password)
	token, err = s.client.AuthPhone(s.phone, s.password)
	fmt.Printf("New token, err: %v, %v\n", token, err)
	if err != nil {
		return
	}

	*s = *token
	fmt.Printf("struct token, local token: %v, %v\n", *s, *token)
	return
}

func (c *client) AuthPhone(phone, password string) (token *SmartToken, err error) {
	values := url.Values{}
	values.Set("phone", phone)
	values.Set("password", password)

	r, err := c.PostForm(c.url("/v2.1/client/auth-phone"), values)
	if err != nil {
		return
	}

	defer r.Body.Close()

	if r.StatusCode != 201 {
		err = errors.New("Not 201 status" + "\nstatus: " + strconv.Itoa(r.StatusCode))
		return
	}

	var resp struct {
		Data struct {
			Token   string `json:"token"`
			Phone   string `json:"phone"`
			SMSID   int    `json:"sms_id"`
			Timeout int    `json:"timeout"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	token = NewSmartToken(c, resp.Data.Token, phone, password)
	fmt.Printf("token from api: %v", resp.Data.Token)
	fmt.Printf("token in token: %v", token.Token())
	return
}

func (c *client) ChangePassword(phone, password string) (smsID int, err error) {
	values := url.Values{}
	values.Set("phone", phone)
	values.Set("password", password)

	r, err := c.PostForm(c.url("/v2.1/client/change-password"), values)
	if err != nil {
		return
	}

	defer r.Body.Close()

	if r.StatusCode != 201 {
		err = errors.New("Not 201 status")
		return
	}

	var resp struct {
		Data struct {
			Phone   string `json:"phone"`
			SMSID   int    `json:"sms_id"`
			Timeout int    `json:"timeout"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	smsID = resp.Data.SMSID
	return
}

func (c *client) Confirm(code string, smsID int, confirmType string) (token string, resp struct {
	Ok  bool
	Err string
}, err error) {
	values := url.Values{}
	values.Set("code", code)
	values.Set("sms_id", strconv.Itoa(smsID))

	r, err := c.PostForm(c.url("/v2.1/client/"+confirmType), values)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var _resp struct {
		Data struct {
			Token   string `json:"token"`
			Message string `json:"message"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&_resp)
	if err != nil {
		return
	}

	resp.Ok = _resp.Success
	resp.Err = _resp.Data.Message
	token = _resp.Data.Token
	return
}

type Card struct {
	Created         int        `json:"created"`
	DateActivated   int        `json:"date_activated"`
	DateBlocked     int        `json:"date_blocked"`
	Number          string     `json:"number"`
	UserGuid        string     `json:"user_guid"`
	MainCardNumber  string     `json:"main_card_number"`
	SlaveCardNumber string     `json:"slave_card_number"`
	Type            CardType   // `json:"type"`
	Status          CardStatus // `json:"status"`
}

func (c *Card) UnmarshalJSON(data []byte) error {
	v := &struct {
		*Card
		Type   int `json:"type"`
		Status int `json:"status"`
	}{Card: c}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	c.Type = CardType(v.Type)
	c.Status = CardStatus(v.Status)
	return nil
}

type CardStatus int

const (
	NewCard CardStatus = iota
	ActiveCard
	BlockedCard
	PaymentCard
)

type CardType int

const (
	MainCard CardType = iota + 1
	SlaveCard
)

func (c *client) SetCard(token *SmartToken, cardNumber string) (card *Card, resp struct {
	Ok  bool
	Err string
}, err error) {
	card, resp, err = c.setCard(token, cardNumber)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.setCard(token, cardNumber)
}

func (c *client) setCard(token *SmartToken, cardNumber string) (card *Card, resp struct {
	Ok  bool
	Err string
}, err error) {
	values := url.Values{}
	values.Set("number", cardNumber)

	req, err := http.NewRequest(
		"POST",
		c.url("/v2/client/card/set-card"),
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return
	}

	req.SetBasicAuth(token.Token(), "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r, err := c.Do(req)
	if err != nil {
		return
	}

	defer r.Body.Close()

	var _resp struct {
		Data struct {
			Card    *Card  `json:"card"`
			Message string `json:"message"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&_resp)
	if err != nil {
		return
	}

	resp.Ok = _resp.Success
	resp.Err = _resp.Data.Message
	card = _resp.Data.Card
	return
}

func (c *client) BarCode(token *SmartToken) (userID int, barCode string, err error) {
	userID, barCode, err = c.barCode(token)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.barCode(token)
}

func (c *client) barCode(token *SmartToken) (userID int, barCode string, err error) {
	req, err := http.NewRequest("", c.url("/v2/client/bar-code"), nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(token.Token(), "")
	r, err := c.Do(req)
	if err != nil {
		return
	}

	defer r.Body.Close()

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	var resp struct {
		Data struct {
			UserID  int    `json:"user_id"`
			BarCode string `json:"bar_code"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	userID = resp.Data.UserID
	barCode = resp.Data.BarCode
	return
}

func (c *client) profileParams() (reader io.ReadCloser, err error) {
	req, err := http.NewRequest("", c.url("/v2/client/profile-params"), nil)
	if err != nil {
		return
	}

	r, err := c.Do(req)
	if err != nil {
		return
	}

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	return r.Body, nil
}

func (c *client) profileFields(token *SmartToken) (reader io.ReadCloser, err error) {
	reader, err = c._profileFields(token)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c._profileFields(token)
}

func (c *client) _profileFields(token *SmartToken) (reader io.ReadCloser, err error) {
	req, err := http.NewRequest("", c.url("/v2/system/profile-fields"), nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(token.Token(), "")

	r, err := c.Do(req)
	if err != nil {
		return
	}

	fmt.Println("@@@ @@@")
	fmt.Println(r.StatusCode)
	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	return r.Body, nil
}

func (c *client) profileLoad(token *SmartToken) (reader io.ReadCloser, err error) {
	reader, err = c._profileLoad(token)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c._profileLoad(token)
}

func (c *client) _profileLoad(token *SmartToken) (reader io.ReadCloser, err error) {
	req, err := http.NewRequest("", c.url("/v2/client/profile"), nil)
	if err != nil {
		return
	}

	fmt.Println(token.Token())
	req.SetBasicAuth(token.Token(), "")

	r, err := c.Do(req)
	if err != nil {
		return
	}

	fmt.Println("*** ***")
	fmt.Println(r.StatusCode)
	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	buf := &bytes.Buffer{}
	tee := io.TeeReader(r.Body, buf)
	bytes, _ := ioutil.ReadAll(tee)
	_ = fmt.Sprint(string(bytes))
	r.Body = ioutil.NopCloser(buf)

	return r.Body, nil
}

func (c *client) ProfileSave(token *SmartToken, profile *Profile) (err error) {
	err = c.profileSave(token, profile)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.profileSave(token, profile)
}

func (c *client) profileSave(token *SmartToken, profile *Profile) error {
	values := url.Values{}
	values.Set("channel_reg", "22")
	for _, f := range profile.Fields {
		values.Set(f.Key, fmt.Sprintf("%v", f.Value))
	}

	req, err := http.NewRequest(
		"PUT",
		c.url("/v2/client/profile"),
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return err
	}

	req.SetBasicAuth(token.Token(), "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r, err := c.Do(req)
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return err
	}

	return nil
}

func (c *client) FieldSave(token *SmartToken, field *Field) (err error) {
	err = c.fieldSave(token, field)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.fieldSave(token, field)
}

func (c *client) fieldSave(token *SmartToken, field *Field) error {
	values := url.Values{}
	values.Set(field.Key, fmt.Sprintf("%v", field.Value))

	req, err := http.NewRequest(
		"PUT",
		c.url("/v2/client/profile"),
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return err
	}

	req.SetBasicAuth(token.Token(), "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r, err := c.Do(req)
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return err
	}

	return nil
}

type Region struct {
	CountryID   int    `json:"country_id"`
	CountryName string `json:"country_name"`
	RegionID    int    `json:"region_id"`
	RegionName  string `json:"region_name"`
}

func (c *client) Regions() (regs []Region, err error) {
	req, err := http.NewRequest("", c.url(fmt.Sprintf("/v2/client/geo/%s/regions", data.Cfg.CountryID)), nil)
	if err != nil {
		return
	}

	r, err := c.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	var resp struct {
		Data struct {
			Country string   `json:"country"`
			Target  []Region `json:"target"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	return resp.Data.Target, nil
}

type City struct {
	Region
	CityID   int    `json:"city_id"`
	CityName string `json:"city_name"`
}

func (c *client) SearchCity(city string) (cs []City, err error) {
	req, err := http.NewRequest("", c.url(fmt.Sprintf("/v2/client/geo/%s/%s/search-city", data.Cfg.CountryID, city)), nil)
	if err != nil {
		return
	}

	r, err := c.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	var resp struct {
		Data struct {
			Country string `json:"country"`
			Target  []City `json:"target"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	return resp.Data.Target, nil
}

func (c *client) GetRegion(regID string) (reg Region, err error) {
	req, err := http.NewRequest("", c.url(fmt.Sprintf("/v2/client/geo/%s/get-region", regID)), nil)
	if err != nil {
		return
	}

	r, err := c.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	var resp struct {
		Data struct {
			ID     string `json:"id"`
			Target Region `json:"target"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	return resp.Data.Target, nil
}

func (c *client) GetCity(cityID string) (ct City, err error) {
	req, err := http.NewRequest("", c.url(fmt.Sprintf("/v2/client/geo/%s/get-city", cityID)), nil)
	if err != nil {
		return
	}

	r, err := c.Do(req)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		err = errors.New("Not 200 status")
		return
	}

	var resp struct {
		Data struct {
			ID     string `json:"id"`
			Target City   `json:"target"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	return resp.Data.Target, nil
}

type HistoryItem struct {
	Type    string
	Data    map[string]interface{}
	Details []map[string]interface{}
}

func (i *HistoryItem) UnmarshalJSON(data []byte) error {
	var v struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if det, ok := v.Data["details"]; ok {
		bytes, err := json.Marshal(det)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &i.Details)
		if err != nil {
			return err
		}
	}

	i.Data = v.Data
	i.Type = v.Type
	return nil
}

type ClientHistory struct {
	DateFrom string        `json:"dateFrom"`
	DateTo   string        `json:"dateTo"`
	Meta     PageMeta      `json:"_meta"`
	Items    []HistoryItem `json:"items"`
	Error    string        `json:"message"`
}

func (c *client) ClientHistory(token *SmartToken, page int) (history ClientHistory, err error) {
	history, err = c.clientHistory(token, page)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.clientHistory(token, page)
}

// page=1&dateFrom=2020-02-17&dateTo=2020-02-24
func (c *client) clientHistory(token *SmartToken, page int) (history ClientHistory, err error) {
	values := url.Values{}
	values.Set("page", fmt.Sprintf("%v", page))
	values.Set("dateFrom", "2015-01-01")
	year, mon, day := time.Now().Date()
	month := ""
	if int(mon) < 10 {
		month = "0" + strconv.Itoa(int(mon))
	} else {
		month = strconv.Itoa(int(mon))
	}
	values.Set("dateTo", fmt.Sprintf("%v-%v-%v", year, month, day))
	fmt.Println(fmt.Sprintf("%v-%v-%v", year, month, day))

	req, err := http.NewRequest(
		"",
		c.url("/v2/client/client-history"),
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return
	}

	req.SetBasicAuth(token.Token(), "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r, err := c.Do(req)
	if err != nil {
		return
	}

	var resp struct {
		Data    ClientHistory `json:"data"`
		Success bool          `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil || !resp.Success {
		return history, fmt.Errorf("%s\n%w\n", resp.Data.Error, err)
	}

	history = resp.Data
	return
}

type Balance struct {
	Currency  string `json:"currency"`
	Account   int    `json:"account"`
	Balance   string `json:"balance"`
	Avialable string `json:"avialable"`
	Error     string `json:"message"`
}

func (c *client) Balance(token *SmartToken) (bal Balance, err error) {
	currency := data.Cfg.Currency

	bal, err = c.balance(token, currency)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.balance(token, currency)
}

func (c *client) balance(token *SmartToken, currency string) (bal Balance, err error) {
	req, err := http.NewRequest(
		"",
		c.url(fmt.Sprintf("/v2/client/account/%s/balance", currency)),
		nil,
	)

	if err != nil {
		return
	}

	req.SetBasicAuth(token.Token(), "")
	r, err := c.Do(req)
	if err != nil {
		return
	}

	var resp struct {
		Data    Balance `json:"data"`
		Success bool    `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil || !resp.Success {
		return bal, fmt.Errorf("%s\n%w\n", resp.Data.Error, err)
	}

	bal = resp.Data
	return
}

type PageMeta struct {
	TotalCount  int `json:"totalCount"`
	PageCount   int `json:"pageCount"`
	CurrentPage int `json:"currentPage"`
	PerPage     int `json:"perPage"`
}

type Actions struct {
	Title   string `json:"title"`
	Image   string `json:"img_path"`
	From    string `json:"act_from"`
	To      string `json:"act_to"`
	Content string `json:"content"`
}

func (c *client) Actions(token *SmartToken, page int) (actions []Actions, meta PageMeta, err error) {
	actions, meta, err = c.actions(token, page)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.actions(token, page)
}

func (c *client) actions(token *SmartToken, page int) (actions []Actions, meta PageMeta, err error) {
	values := url.Values{}
	values.Set("page", fmt.Sprintf("%v", page))

	req, err := http.NewRequest(
		"",
		c.url("/v2/client/partner/actions"),
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return
	}

	req.SetBasicAuth(token.Token(), "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r, err := c.Do(req)
	if err != nil {
		return
	}

	var resp struct {
		Data struct {
			Items []Actions `json:"items"`
			Meta  PageMeta  `json:"_meta"`
			Error string    `json:"message"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil || !resp.Success {
		return actions, meta, fmt.Errorf("%s\n%w\n", resp.Data.Error, err)
	}

	meta = resp.Data.Meta
	actions = resp.Data.Items
	return
}

type News struct {
	Name  string `json:"name"`
	Descr string `json:"description"`
	Short string `json:"description_short"`
	Image string `json:"image_path"`
}

func (c *client) News(token *SmartToken, page int) (news []News, meta PageMeta, err error) {
	news, meta, err = c.news(token, page)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.news(token, page)
}

func (c *client) news(token *SmartToken, page int) (news []News, meta PageMeta, err error) {
	values := url.Values{}
	values.Set("page", fmt.Sprintf("%v", page))

	req, err := http.NewRequest(
		"",
		c.url("/v2/client/partner/news-all"),
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return
	}

	req.SetBasicAuth(token.Token(), "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r, err := c.Do(req)
	if err != nil {
		return
	}

	var resp struct {
		Data struct {
			Items []News   `json:"items"`
			Meta  PageMeta `json:"_meta"`
			Error string   `json:"message"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil || !resp.Success {
		return news, meta, fmt.Errorf("%s\n%w\n", resp.Data.Error, err)
	}

	meta = resp.Data.Meta
	news = resp.Data.Items
	return
}
