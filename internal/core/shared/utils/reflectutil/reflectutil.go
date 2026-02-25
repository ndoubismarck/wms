package reflectutil

import (
	"fmt"
	"reflect"
	"strings"
)

func CopyStruct(source, destination any) error {
	sourceValue := reflect.ValueOf(source)
	destinationValue := reflect.ValueOf(destination).Elem()

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldType := sourceValue.Type().Field(i).Type
		destinationFieldType := destinationValue.Type().Field(i).Type
		if sourceFieldType != destinationFieldType {
			return fmt.Errorf("field %s value of type %s is not assignable to type %s", sourceValue.Type().Field(i).Name, sourceFieldType, destinationFieldType)
		}
		destinationValue.Field(i).Set(sourceValue.Field(i))
	}
	return nil
}

func JsonStructToMap(input any) (map[string]any, error) {
	result := make(map[string]any)
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}
	valType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := valType.Field(i).Tag.Get("json")
		if len(fieldName) == 0 {
			fieldName = valType.Field(i).Name
		}
		if strings.Contains(fieldName, "omitempty") {
			if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				continue
			}
		}
		filedData := strings.Split(fieldName, ",")
		if len(filedData) == 0 {
			result[fieldName] = field.Interface()
		} else {
			result[filedData[0]] = field.Interface()
		}
	}
	return result, nil
}
