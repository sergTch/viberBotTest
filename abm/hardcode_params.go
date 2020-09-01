package abm

import "strings"

var requiredParams []string
var htmlTags map[string]string

func init() {
	requiredParams = []string{
		"email",
		"first_name",
		"middle_name",
		"last_name",
		"address",
		"gender",
		"birth_day",
		"region",
		"city"}
	DataType = map[int]string{}
	FieldType = map[int]string{}
	DataType[1] = "Text"
	DataType[2] = "Dropdown list"
	DataType[3] = "Date"
	DataType[4] = "Checkbox"
	FieldType[0] = "String"
	FieldType[1] = "Integer"
	FieldType[2] = "Birthday"

	htmlTags = map[string]string{}
	htmlTags["<span style=\"color: "] = "<font color=\""
	htmlTags["</span>"] = "</font>"
	htmlTags[";\">"] = "\">"
	htmlTags["<p>"] = ""
	htmlTags["</p>"] = ""
	htmlTags["<b>"] = ""
	htmlTags["</b>"] = ""
	htmlTags["<strong>"] = "<b>"
	htmlTags["</strong>"] = "</b>"
	htmlTags["<em>"] = "<i>"
	htmlTags["</em>"] = "</i>"
}

func remakeHtml(s *string) {
	for k, v := range htmlTags {
		*s = strings.ReplaceAll(*s, k, v)
	}
}

func (p *Profile) fillMainParams() {
	strField := -1
	numberField := -1
	dayField := -1
	dropData := -1
	textData := -1
	dayData := -1

	for id, dataType := range DataType {
		if dataType == "Dropdown list" {
			dropData = id
		}
		if dataType == "Date" {
			dayData = id
		}
		if dataType == "Text" {
			textData = id
		}
	}
	for id, fieldType := range FieldType {
		if fieldType == "Integer" {
			numberField = id
		}
		if fieldType == "String" {
			strField = id
		}
		if fieldType == "Birthday" {
			dayField = id
		}
	}

	if f, ok := p.Main["gender"]; ok {
		f.DataType = dropData
		f.FieldType = numberField
	}

	if f, ok := p.Main["email"]; ok {
		f.DataType = textData
		f.FieldType = strField
	}

	if f, ok := p.Main["first_name"]; ok {
		f.DataType = textData
		f.FieldType = strField
	}

	if f, ok := p.Main["middle_name"]; ok {
		f.DataType = textData
		f.FieldType = strField
	}

	if f, ok := p.Main["last_name"]; ok {
		f.DataType = textData
		f.FieldType = strField
	}

	if f, ok := p.Main["address"]; ok {
		f.DataType = textData
		f.FieldType = strField
	}

	if f, ok := p.Main["birth_day"]; ok {
		f.DataType = dayData
		f.FieldType = dayField
	}

	p.City = p.Main["city"]
	p.City.FieldType = numberField
	p.City.Key = "id_city"
	p.Region = p.Main["region"]
	p.Region.FieldType = numberField
	p.Region.Key = "id_region"
	delete(p.Main, "city")
	delete(p.Main, "region")
}
