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

	var resp struct {
		Data map[string]interface{} `json:"data"`
	}

	buf := &bytes.Buffer{}
	tee := io.TeeReader(r.Body, buf)
	bytes, _ := ioutil.ReadAll(tee)
	fmt.Println(string(bytes))
	r.Body = ioutil.NopCloser(buf)

	// var v struct {
	// 	Data struct {
	// 		Params struct {
	// 			Vals []struct {
	// 				ID   string `json:"id"`
	// 				Name string `json:"value"`
	// 			} `json:"work_status_params"`
	// 		} `json:"params"`
	// 	} `json:"data"`
	// }

	// var v struct {
	// 	Data struct {
	// 		Params map[string]interface{} `json:"params"`
	// 	} `json:"data"`
	// }

	//s := v.Data.Params["work_status_params"]

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return
	}

	fmt.Printf("%+v", resp)

	return nil
}
