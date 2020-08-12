package abm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

var DataType map[int]string
var FieldType map[int]string
var Required map[int]string

type Field struct {
	Value     interface{} `json:"-"`
	Name      string      `json:"name"`
	Key       string      `json:"key"`
	Required  bool        `json:"required"`
	FieldType int         `json:"field_type"`
	DataType  int         `json:"data_type"`
	Schema    Schema      `json:"items"`
}

type Profile struct {
	Fields     []*Field
	Region     *Field
	City       *Field
	schemas    map[string]interface{}
	Additional map[string]*Field
	Main       map[string]*Field
}

func NewProfile() *Profile {
	return &Profile{
		schemas:    map[string]interface{}{},
		Additional: map[string]*Field{},
		Main:       map[string]*Field{},
		//DataType:   map[int]string{},
		//FieldType:  map[int]string{},
		//Required:   map[int]string{},
	}
}

func (c *client) Profile(token *SmartToken) (*Profile, error) {
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

	for _, k := range requiredParams {
		s, _ := p.Schema(k)
		field := Field{Name: k, Key: k, Required: resp.Data.Params.Required[k], Schema: s}
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
			DataType  Schema  `json:"data_type_param"`
			FieldType Schema  `json:"field_type_param"`
			Required  Schema  `json:"required_param"`
		} `json:"data"`
	}

	buf := &bytes.Buffer{}
	tee := io.TeeReader(r, buf)
	fmt.Println("$$$ $$$")
	bs, _ := ioutil.ReadAll(tee)
	_ = fmt.Sprint(string(bs))
	r = ioutil.NopCloser(buf)

	err := json.NewDecoder(r).Decode(&resp)
	if err != nil {
		return err
	}

	fields := map[string]bool{}
	for _, f := range resp.Data.Fields {
		fields[f.Key] = f.Required
		p.schemas[f.Key] = f.Schema
	}

	Required = map[int]string{}
	for _, v := range resp.Data.Required {
		id, _ := strconv.Atoi(v.ID)
		Required[id] = v.Value
	}
	for _, v := range resp.Data.Fields {
		field := v
		p.Fields = append(p.Fields, &field)
		p.Additional[v.Key] = &field
	}

	return nil
}

type Schema []Entry
type Entry struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (e *Entry) UnmarshalJSON(data []byte) error {
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

func (p *Profile) Schema(param string) (s Schema, ok bool) {
	val, ok := p.schemas[param+"_params"]
	if !ok {
		return
	}

	fmt.Println(val)
	bs, err := json.Marshal(val)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &s)
	if err != nil {
		return
	}

	ok = true
	return
}

func valueTosString(value interface{}) string {
	switch value.(type) {
	case string:
		return fmt.Sprint(value)
	case float64, float32:
		return strings.Split(fmt.Sprintf("%f\n", value), ".")[0]
	default:
		return fmt.Sprintf("%v\n", value)
	}
}

func (f *Field) ToString() string {
	text := ""
	if f.Required {
		text += "*"
	}
	text += f.Name + ": "
	if f.Key == "id_region" {
		if valueTosString(f.Value) == "0" {
			return text
		} else {
			region, err := Client.GetRegion(valueTosString(f.Value))
			if err != nil {
				fmt.Println(err)
				return text
			}
			return text + region.RegionName
		}
	}
	if f.Key == "id_city" {
		fmt.Println(valueTosString(f.Value))
		if valueTosString(f.Value) == "0" {
			return text
		} else {
			city, err := Client.GetCity(valueTosString(f.Value))
			if err != nil {
				fmt.Println(err)
				return text
			}
			return text + city.CityName
		}
	}
	if f.Value != nil {
		if FieldType[f.FieldType] == "Integer" {
			str := valueTosString(f.Value)
			for _, ent := range f.Schema {
				if ent.ID == str {
					text += ent.Value
				}
			}
			if len(f.Schema) == 0 {
				if str == "0" {
					text += "нет"
				} else if str == "1" {
					text += "да"
				}
			}
		} else {
			text += fmt.Sprint(f.Value)
		}
	}
	return text
}

func (p *Profile) ToString() string {
	text := ""
	for _, field := range p.Fields {
		text += field.ToString() + "\n"
	}
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
