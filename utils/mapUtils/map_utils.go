package mapUtils

import (
	"Quote/utils/anyUtils"
	"reflect"

	"github.com/pkg/errors"
)

type fieldName func(reflect.StructField) string

func GetKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

func GetValues[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func toStruct(json map[string]any, s any, getFieldName fieldName) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		if value, ok := json[getFieldName(field)]; ok && value != nil {
			stValue.Field(i).Set(reflect.ValueOf(value))
		}
	}
	return nil
}

func ToStruct_usingFieldNames(m map[string]any, s any) error {
	return toStruct(m, s,
		func(field reflect.StructField) string {
			return field.Name
		})
}

func ToStruct_usingJsonTags(json map[string]any, s any) error {
	return toStruct(json, s,
		func(field reflect.StructField) string {
			return field.Tag.Get("json")
		})
}

// The returned "bool" value indicates whether the given fieldname exists as a key in the given map
func FieldAnyToType(m map[string]any, fieldName string, basicDataType string) (bool, error) {
	originalValue, ok := m[fieldName]
	if !ok {
		return false, nil
	}

	var convertedValue any
	var err error
	switch basicDataType {
	case "int":
		convertedValue, err = anyUtils.ToInt(originalValue)
	case "int64":
		convertedValue, err = anyUtils.ToInt64(originalValue)
	case "uint64":
		convertedValue, err = anyUtils.ToUint64(originalValue)
	case "float32":
		convertedValue, err = anyUtils.ToFloat32(originalValue)
	case "float64":
		convertedValue, err = anyUtils.ToFloat64(originalValue)
	case "string":
		convertedValue, err = anyUtils.ToString(originalValue)
	case "time.Time":
		convertedValue, err = anyUtils.ToTime(originalValue)
	case "[]string":
		convertedValue, err = anyUtils.ToStringArray(originalValue)
	case "[]int":
		convertedValue, err = anyUtils.ToIntArray(originalValue)
	case "[]uint8":
		convertedValue, err = anyUtils.ToUInt8Array(originalValue)
	case "*string":
		convertedValue, err = anyUtils.ToStringPtr(originalValue)
	case "*int":
		convertedValue, err = anyUtils.ToIntPtr(originalValue)
	case "*uint64":
		convertedValue, err = anyUtils.ToUInt64Ptr(originalValue)
	case "*time.Time":
		convertedValue, err = anyUtils.ToTimePtr(originalValue)
	default:
		return true, errors.Errorf("Data type '%v' not supported for conversion from type any", basicDataType)
	}

	if err != nil {
		return false, err
	}
	m[fieldName] = convertedValue

	return true, nil
}

func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// Returns true if the given map has at least one of the given keys
func HasAnyKey[K comparable, V any](m map[K]V, keyList []K) bool {
	for _, key := range keyList {
		hasKey := HasKey(m, key)
		if hasKey == true {
			return true
		}
	}

	return false
}

func RemoveKey[K comparable, V any](m map[K]V, key K) {
	delete(m, key)
}

// Removes all the given keys from the given map
func RemoveKeyList[K comparable, V any](m map[K]V, keyList []K) {
	for _, key := range keyList {
		RemoveKey(m, key)
	}
}

func Merge[K comparable, V any](m1 map[K]V, m2 map[K]V) map[K]V {
	merged := make(map[K]V)
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}
