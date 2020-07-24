package abm

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sergTch/viberBotTest/data"
)

type Client struct {
	client *http.Client
	apiURL string
}

func New() *Client {
	return &Client{
		client: &http.Client{},
		apiURL: data.ApiUrl,
	}
}

func (c *Client) url(endpoint string) string {
	return c.apiURL + endpoint
}

// client := abm.New()
// ok, err := client.CheckPhone("380671810640")
// fmt.Printf("phone: %v %v", ok, err)

func (c *Client) CheckPhone(number string) (ok bool, err error) {
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

	ok = true
	return
}

func (c *Client) Register(phone, password, signature string) (smsID int, err error) {
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

func (c *Client) Confirm(code string, smsID int) (token string, ok bool, err error) {
	values := url.Values{}
	values.Set("code", code)
	values.Set("sms_id", strconv.Itoa(smsID))

	r, err := c.client.PostForm(c.url("/v2.1/client/registration-confirm"), values)
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

	ok = resp.Success
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

func (c *Client) SetCard(number string) (card *Card, ok bool, err error) {
	values := url.Values{}
	values.Set("number", number)

	r, err := c.client.PostForm(c.url("/v2/client/card/set-card"), values)
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

func (c *Client) BarCode(token string) (userID int, barCode string, err error) {
	req, err := http.NewRequest("", c.url("/v2/client/bar-code"), nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(token, "")
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