package paths

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNotAStruct = errors.New("CheckPointersNotNil expects a struct or pointer to struct")

func CheckPointersNotNil(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := t.Field(i)

		tagVal := fieldType.Tag.Get("svmkit")
		if tagVal == "optional" {
			continue
		}

		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			return fmt.Errorf("%s is nil", fieldType.Name)
		}
	}

	return nil
}
