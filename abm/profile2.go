package abm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Item struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Field struct {
	//Value     interface{} `json:"value"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	FieldType int    `json:"field_type"`
	DataType  int    `json:"data_type"`
	Required  bool   `json:"required"`
	Items     []Item `json:"items"`
}

type Profile2 struct {
	Gender     Field
	Region     Field
	City       Field
	Main       map[string]Field
	Additional map[string]Field
}

func (c *client) NewProfile2(token string) *Profile2 {
	prof := Profile2{}
	err := c.LoadFields(&prof, token)
	if err != nil {
		fmt.Println(err)
	}
	return &prof
}

func (c *client) LoadFields(prof *Profile2, token string) (err error) {
	req, err := http.NewRequest("", c.url("/v2/client/profile-params"), nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(token, "")
	r, err := c.client.Do(req)
	if err != nil {
		return
	}
	if r.StatusCode != 200 {
		fmt.Println("status: ", r.StatusCode)
		err = errors.New("Not 200 status")
		return
	}
	//buf := &bytes.Buffer{}
	//tee := io.TeeReader(r.Body, buf)
	//bytes, _ := ioutil.ReadAll(tee)
	var resp struct {
		Data struct {
			Fields []Field `json:"fields"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}
	for _, field := range resp.Data.Fields {
		prof.Additional[field.Key] = field
	}
	return
}

func (c *client) ProfileLoadTest(token string) (err error) {
	req, err := http.NewRequest("", c.url("/v2/client/profile-params"), nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(token, "")

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

	return nil
}
