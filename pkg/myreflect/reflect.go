package myreflect

import (
	"reflect"
)

func FillMapNotZeros(pairs map[string]any) map[string]any {
	fields := make(map[string]any)
	for key, value := range pairs {
		v := reflect.ValueOf(value)

		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			v = v.Elem()
		}

		if !v.IsZero() {
			fields[key] = value
		}
	}
	return fields
}

func IsEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
	}

	// Проверяем, что это структура
	if v.Kind() != reflect.Struct {
		return false
	}
	// Проверяем, что все поля структуры имеют нулевые значения
	for idx := range v.NumField() {
		if !v.Field(idx).IsZero() {
			return false
		}
	}

	return true
}
