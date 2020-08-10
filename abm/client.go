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
		apiURL: data.ApiUrl,
	}
}

func (c *client) url(endpoint string) string {
	return c.apiURL + endpoint
}

func (c *client) CheckPhone(number string) (ok bool, err error) {
	values := url.Values{}
	values.Set("phone", number)

	r, err := c.client.PostForm(c.url("/v2/client/check-phone"), values)
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

func (c *client) Register(phone, password, signature string) (smsID int, err error) {
	values := url.Values{}
	values.Set("phone", phone)
	values.Set("password", password)
	values.Set("signature", signature)

	r, err := c.client.PostForm(c.url("/v2.1/client/registration"), values)
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

	token, phone, password, signature string
}

func NewSmartToken(client *client, token, phone, password, signature string) *SmartToken {
	return &SmartToken{
		client:    client,
		token:     token,
		phone:     phone,
		password:  password,
		signature: signature,
	}
}

func (s *SmartToken) Token() string {
	return s.token
}

func (s *SmartToken) Renew() (token *SmartToken, err error) {
	fmt.Println("*** RENEWING TOKEN ***")
	fmt.Printf("params: %v, %v, %v\n", s.phone, s.password, s.signature)
	token, err = s.client.AuthPhone(s.phone, s.password, s.signature)
	fmt.Printf("New token, err: %v, %v\n", token, err)
	if err != nil {
		return
	}

	*s = *token
	fmt.Printf("struct token, local token: %v, %v\n", *s, *token)
	return
}

func (c *client) AuthPhone(phone, password, signature string) (token *SmartToken, err error) {
	values := url.Values{}
	values.Set("phone", phone)
	values.Set("password", password)
	values.Set("signature", signature)

	r, err := c.client.PostForm(c.url("/v2.1/client/auth-phone"), values)
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

	token = NewSmartToken(c, resp.Data.Token, phone, password, signature)
	return
}

func (c *client) ChangePassword(phone, password, signature string) (smsID int, err error) {
	values := url.Values{}
	values.Set("phone", phone)
	values.Set("password", password)
	values.Set("signature", signature)

	r, err := c.client.PostForm(c.url("/v2.1/client/change-password"), values)
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

func (c *client) Confirm(code string, smsID int, confirmType string) (token string, ok bool, err error) {
	values := url.Values{}
	values.Set("code", code)
	values.Set("sms_id", strconv.Itoa(smsID))

	r, err := c.client.PostForm(c.url("/v2.1/client/"+confirmType), values)
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
			Token string `json:"token"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	//ok = resp.Success
	ok = true
	token = resp.Data.Token
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

func (c *client) SetCard(token *SmartToken, cardNumber string) (card *Card, ok bool, err error) {
	card, ok, err = c.setCard(token, cardNumber)
	if err == nil {
		return
	}

	token, err = token.Renew()
	if err != nil {
		return
	}

	return c.setCard(token, cardNumber)
}

func (c *client) setCard(token *SmartToken, cardNumber string) (card *Card, ok bool, err error) {
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
	r, err := c.client.Do(req)
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
			Card *Card `json:"card"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	ok = resp.Success
	card = resp.Data.Card
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
	r, err := c.client.Do(req)
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

	r, err := c.client.Do(req)
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

	r, err := c.client.Do(req)
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

	r, err := c.client.Do(req)
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
	fmt.Println(string(bytes))
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
	r, err := c.client.Do(req)
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
	r, err := c.client.Do(req)
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
	req, err := http.NewRequest("", c.url(fmt.Sprintf("/v2/client/geo/%s/regions", data.CountryID)), nil)
	if err != nil {
		return
	}

	r, err := c.client.Do(req)
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
	req, err := http.NewRequest("", c.url(fmt.Sprintf("/v2/client/geo/%s/%s/search-city", data.CountryID, city)), nil)
	if err != nil {
		return
	}

	r, err := c.client.Do(req)
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

	r, err := c.client.Do(req)
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

	r, err := c.client.Do(req)
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
