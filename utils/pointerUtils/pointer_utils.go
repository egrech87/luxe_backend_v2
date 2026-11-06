package pointerUtils

import "reflect"

func ToPointer[T any](obj T) *T {
	return &obj
}

func IsNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
