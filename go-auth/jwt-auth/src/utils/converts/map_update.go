package converts

import (
	"reflect"
	"strings"
)

func MapToKeyAndValueUptade(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		// Pega a chave da tag JSON
		jsonTag := fieldType.Tag.Get("json")
		var key string
		if jsonTag == "" || jsonTag == "-" {
			// Usa o nome do campo, tudo em minúsculo
			key = strings.ToLower(fieldType.Name)
		} else {
			// Usa a parte antes da vírgula na tag json (caso tenha algo tipo `json:"title,omitempty"`)
			key = strings.Split(jsonTag, ",")[0]
			if key == "" || key == "-" {
				key = strings.ToLower(fieldType.Name)
			}
		}
		// Ignora campos com valor zero
		if !fieldValue.IsZero() {
			result[key] = fieldValue.Interface()
		}
	}
	return result
}
