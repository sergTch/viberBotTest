package abm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Profile struct {
	Params  map[string]bool
	Fields  map[string]bool
	schemas map[string]interface{}
}

func NewProfile() *Profile {
	return &Profile{
		Params:  map[string]bool{},
		Fields:  map[string]bool{},
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

	p.Params = resp.Data.Params.Required
	fmt.Printf("%+v\n", p.Params)
	p.schemas = resp2.Data.Schema
	fmt.Printf("%+v\n", p.schemas)
	return nil
}

func (p *Profile) readFields(r io.Reader) error {
	var resp struct {
		Data struct {
			Fields []struct {
				Name string `json:"name"`
				Key  string `json:"key"`
				// FieldType string   `json:"field_type"`
				// DataType  string   `json:"data_type"`
				Required bool     `json:"required"`
				Values   []string `json:"value"`
			} `json:"fields"`
		} `json:"data"`
	}

	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)
	fmt.Println("$$$ $$$")
	bs, _ := ioutil.ReadAll(tee)
	fmt.Println(string(bs))
	r = buf

	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		return err
	}

	fields := map[string]bool{}
	for _, f := range resp.Data.Fields {
		fields[f.Key] = f.Required

		type entry struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		schema := []entry{}
		for i, v := range f.Values {
			schema = append(schema, entry{ID: i, Name: v})
		}

		p.schemas[f.Key] = schema
		fmt.Println("!!! !!!")
		fmt.Println(schema)
	}

	p.Fields = fields
	return nil
}

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
