package abm

func (p *Profile) fillMainParams() {
	strField := -1
	numberField := -1
	dayField := -1
	dropData := -1
	textData := -1
	dayData := -1

	for id, dataType := range p.DataType {
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
	for id, fieldType := range p.FieldType {
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
}
