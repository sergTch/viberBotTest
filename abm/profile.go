package abm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Profile struct {
	params  map[string]bool
	fields  map[string]bool
	schemas map[string]interface{}
}

func NewProfile() *Profile {
	return &Profile{
		params:  map[string]bool{},
		fields:  map[string]bool{},
		schemas: map[string]interface{}{},
	}
}

func (p *Profile) readParams(r io.Reader) error {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)

	var resp struct {
		Data struct {
			Params struct {
				Required map[string]bool `json:"required"`
			} `json:"params"`
		} `json:"data"`
	}

	err := json.NewDecoder(tee).Decode(&resp)
	if err != nil {
		return err
	}

	var resp2 struct {
		Data struct {
			Schema map[string]interface{} `json:"params"`
		} `json:"data"`
	}

	err = json.NewDecoder(buf).Decode(&resp2)
	if err != nil {
		return err
	}

	p.params = resp.Data.Params.Required
	fmt.Printf("%+v\n", p.params)
	p.schemas = resp2.Data.Schema
	fmt.Printf("%+v\n", p.schemas)
	return nil
}

func (p *Profile) readFileds(r io.Reader) {}

type schema []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (p *Profile) Schema(param string) (s schema, ok bool) {
	val, ok := p.schemas[param+"_params"]
	if !ok {
		return
	}

	b, err := json.Marshal(val)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}

	ok = true
	return
}
