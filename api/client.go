package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/sergTch/viberBotTest/data"
)

type Client struct {
	client *http.Client
	api    string
}

func New() *Client {
	return &Client{
		client: &http.Client{},
		api:    data.API,
	}
}

// client := api.New()
// ok, err := client.CheckPhone("380671810640")
// fmt.Printf("phone: %v %v", ok, err)

func (c *Client) CheckPhone(number string) (bool, error) {
	values := url.Values{}
	values.Set("phone", number)
	url := c.api + "/v2/client/check-phone"
	r, err := c.client.PostForm(url, values)
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	if r.StatusCode != 201 {
		return false, nil
	}

	var resp struct {
		Data struct {
			Is_exist bool `json:"is_exist"`
		} `json:"data"`
	}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, nil
	}

	return true, nil
}
