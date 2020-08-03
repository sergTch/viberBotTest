package abm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

type Field struct {
	Value     interface{} `json:"-"`
	Name      string      `json:"name"`
	Key       string      `json:"key"`
	Required  bool        `json:"required"`
	FieldType int         `json:"field_type"`
	DataType  int         `json:"data_type"`
	Schema    schema      `json:"items"`
}

type Profile struct {
	Fields     []*Field
	Region     *Field
	City       *Field
	schemas    map[string]interface{}
	Additional map[string]*Field
	Main       map[string]*Field
	DataType   map[int]string
	FieldType  map[int]string
	Required   map[int]string
}

func NewProfile() *Profile {
	return &Profile{
		schemas:    map[string]interface{}{},
		Additional: map[string]*Field{},
		Main:       map[string]*Field{},
		DataType:   map[int]string{},
		FieldType:  map[int]string{},
		Required:   map[int]string{},
	}
}

func (c *client) Profile(token string) (*Profile, error) {
	r, err := c.profileParams()
	if err != nil {
		return nil, err
	}

	p := NewProfile()
	err = p.readParams(r)
	if err != nil {
		return nil, err
	}

	r, err = c.profileFields(token)
	if err != nil {
		return nil, err
	}

	err = p.readFields(r)
	if err != nil {
		return nil, err
	}

	p.fillMainParams()

	r, err = c.profileLoad(token)
	if err != nil {
		return nil, err
	}

	err = p.readProfile(r)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Profile) readParams(r io.ReadCloser) error {
	defer r.Close()
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

	p.schemas = resp2.Data.Schema
	fmt.Printf("%+v\n", p.schemas)

	for k, v := range resp.Data.Params.Required {
		s, _ := p.Schema(k)
		field := Field{Name: k, Key: k, Required: v, Schema: s}
		p.Main[k] = &field
		p.Fields = append(p.Fields, &field)
	}

	return nil
}

func (p *Profile) readFields(r io.ReadCloser) error {
	defer r.Close()
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
	r = ioutil.NopCloser(buf)

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

	for _, v := range resp.Data.DataType {
		id, _ := strconv.Atoi(v.ID)
		p.DataType[id] = v.Value
	}
	for _, v := range resp.Data.FieldType {
		id, _ := strconv.Atoi(v.ID)
		p.FieldType[id] = v.Value
	}
	for _, v := range resp.Data.Required {
		id, _ := strconv.Atoi(v.ID)
		p.Required[id] = v.Value
	}
	for _, v := range resp.Data.Fields {
		field := v
		p.Fields = append(p.Fields, &field)
		p.Additional[v.Key] = &field
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
		ID    interface{} `json:"id"`
		Name  string      `json:"name"`
		Value string      `json:"value"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.Value = v.Name
	if e.Value == "" {
		e.Value = v.Value
	}

	if id, ok := v.ID.(string); ok {
		e.ID = id
	} else if id, ok := v.ID.(float64); ok {
		e.ID = fmt.Sprintf("%v", int(id))
	}

	return nil
}

func (p *Profile) Schema(param string) (s schema, ok bool) {
	val, ok := p.schemas[param+"_params"]
	if !ok {
		return
	}

	fmt.Println(val)
	bs, err := json.Marshal(val)
	if err != nil {
		return
	}

	fmt.Println(string(bs))

	err = json.Unmarshal(bs, &s)
	if err != nil {
		return
	}
	fmt.Println(s)

	ok = true
	return
}

func (p *Profile) ToString() string {
	text := ""
	for _, field := range p.Main {
		if field.Required {
			text += "*"
		}
		text += field.Name + ": " + "\n"
	}
	for _, field := range p.Additional {
		if field.Required {
			text += "*"
		}
		text += field.Name + ": " + "\n"
	}
	if p.City.Required {
		text += "*"
	}
	text += p.City.Name + ": " + "\n"
	if p.Region.Required {
		text += "*"
	}
	text += p.Region.Name + ": " + "\n"
	return text
}

func (p *Profile) readProfile(r io.ReadCloser) error {
	defer r.Close()
	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		return err
	}

	for _, v := range p.Fields {
		if val, ok := resp.Data[v.Key]; ok {
			v.Value = val
		}
	}
	return nil
}
