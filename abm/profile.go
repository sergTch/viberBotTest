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
				Name     string `json:"name"`
				Key      string `json:"key"`
				Required bool   `json:"required"`
				Schema   []struct {
					ID    string `json:"id"`
					Value string `json:"value"`
				} `json:"items"`
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
		schema := schema{}
		for _, v := range f.Schema {
			schema = append(schema, struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{ID: v.ID, Name: v.Value})
		}
		p.schemas[f.Key] = schema
		fmt.Println("!!! !!!")
		fmt.Println(f.Schema)
		fmt.Println(schema)
	}

	p.Fields = fields
	return nil
}

type schema []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (p *Profile) Schema(param string) (s schema, ok bool) {
	val, ok := p.schemas[param+"_params"]
	if !ok {
		return
	}

	fmt.Println(val)
	b, err := json.Marshal(val)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}
	fmt.Println(s)

	ok = true
	return
}
