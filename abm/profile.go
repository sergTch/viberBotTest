package abm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Profile struct {
	Gender     Field
	Region     Field
	City       Field
	Params     map[string]bool
	Fields     map[string]bool
	schemas    map[string]interface{}
	Additional map[string]Field
	DataType   schema
	FieldType  schema
	Required   schema
}

func NewProfile() *Profile {
	return &Profile{
		Params:     map[string]bool{},
		Fields:     map[string]bool{},
		schemas:    map[string]interface{}{},
		Additional: map[string]Field{},
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
			Fields    []Field `json:"fields"`
			DataType  schema  `json:"data_type_param"`
			FieldType schema  `json:"field_type_param"`
			Required  schema  `json:"required_param"`
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
		p.schemas[f.Key] = f.Schema
		fmt.Println("!!! !!!")
		fmt.Println(f.Schema)
	}

	p.Fields = fields
	p.DataType = resp.Data.DataType
	p.FieldType = resp.Data.FieldType
	p.Required = resp.Data.Required
	for _, v := range resp.Data.Fields {
		p.Additional[v.Key] = v
	}

	return nil
}

type schema []entry
type entry struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (e *entry) UnmarshalJSON(data []byte) error {
	var v struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.Value = v.Name
	if e.Value == "" {
		e.Value = v.Value
	}

	e.ID = v.ID
	if e.ID == "" {
		var v struct {
			ID int `json:"id"`
		}

		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}

		e.ID = fmt.Sprintf("%v", v.ID)
	}
	return nil
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
