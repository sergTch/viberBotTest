package abm

import (
	"bytes"
	"encoding/json"
	"io"
)

type Profile struct {
	params map[string]bool
	fields map[string]bool
	data   map[string]interface{}
}

func NewProfile() *Profile {
	return &Profile{
		params: map[string]bool{},
		fields: map[string]bool{},
		data:   map[string]interface{}{},
	}
}

func (p *Profile) readParams(r io.Reader) error {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)
	var params struct {
		Data struct {
			Params struct {
				Required map[string]bool `json:"required"`
			} `json:"params"`
		} `json:"data"`
	}

	err := json.NewDecoder(tee).Decode(&params)
	if err != nil {
		return err
	}

	var data struct {
		Data struct {
			Params map[string]interface{} `json:"params"`
		} `json:"data"`
	}

	err = json.NewDecoder(buf).Decode(&data)
	if err != nil {
		return err
	}

	p.params = params.Data.Params.Required
	p.data = data.Data.Params
	return nil
}

func (p *Profile) readFileds(r io.Reader) {}

type schema []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (p *Profile) Schema(param string) schema {
	val, ok := p.data[param+"_param"]
	if !ok {
		return nil
	}
	return val.(schema)
}
