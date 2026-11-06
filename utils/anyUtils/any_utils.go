package anyUtils

import (
	"Quote/utils/timeUtils"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToString(value any) (string, error) {
	return fmt.Sprintf("%v", value), nil
}

func ToFloat32(value any) (float32, error) {
	switch value := value.(type) {
	case uint64:
		return float32(value), nil
	case int64:
		return float32(value), nil
	case float64:
		return float32(value), nil
	case float32:
		return value, nil
	case int32:
		return float32(value), nil
	case int16:
		return float32(value), nil
	case int8:
		return float32(value), nil
	case string:
		result, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0, err
		}
		return float32(result), nil
	}

	return 0, fmt.Errorf("failed to convert %v to float32", value)
}

func ToFloat64(value any) (float64, error) {
	switch value := value.(type) {
	case uint64:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case float64:
		return value, nil
	case float32:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int8:
		return float64(value), nil
	case string:
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err
		}
		return float64(result), nil
	}

	return 0, fmt.Errorf("failed to convert %v to float64", value)
}

func ToInt(value any) (int, error) {
	switch value := value.(type) {
	case int:
		return value, nil
	case *int:
		return *value, nil
	case uint64:
		return int(value), nil
	case int64:
		return int(value), nil
	case float64:
		return int(value), nil
	case float32:
		return int(value), nil
	case int32:
		return int(value), nil
	case int16:
		return int(value), nil
	case int8:
		return int(value), nil
	case string:
		result, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return result, nil
	}

	return 0, fmt.Errorf("failed to convert %v to int", value)
}

func ToInt64(value any) (int64, error) {
	switch value := value.(type) {
	case uint64:
		return int64(value), nil
	case int64:
		return value, nil
	case float64:
		return int64(value), nil
	case float32:
		return int64(value), nil
	case int32:
		return int64(value), nil
	case int16:
		return int64(value), nil
	case int8:
		return int64(value), nil
	case string:
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
		return result, nil
	}

	return 0, fmt.Errorf("failed to convert %v to uint64", value)
}

func ToUint8(value any) (uint8, error) {
	switch value := value.(type) {
	case uint64:
		return uint8(value), nil
	case int64:
		return uint8(value), nil
	case float64:
		return uint8(value), nil
	case float32:
		return uint8(value), nil
	case int32:
		return uint8(value), nil
	case int16:
		return uint8(value), nil
	case int8:
		return uint8(value), nil
	case uint8:
		return value, nil
	case string:
		result, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return 0, err
		}
		return uint8(result), nil
	}

	return 0, fmt.Errorf("failed to convert %v to uint64", value)
}

func ToUint64(value any) (uint64, error) {
	switch value := value.(type) {
	case uint64:
		return value, nil
	case int64:
		return uint64(value), nil
	case float64:
		return uint64(value), nil
	case float32:
		return uint64(value), nil
	case int32:
		return uint64(value), nil
	case int16:
		return uint64(value), nil
	case int8:
		return uint64(value), nil
	case *uint64:
		return *value, nil
	case string:
		result, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, err
		}
		return result, nil
	}

	return 0, fmt.Errorf("failed to convert %v to uint64", value)
}

func ToTime(value any) (time.Time, error) {
	if value == nil {
		return time.Now(), nil
	}

	var convertedValue time.Time
	var err error

	switch value := value.(type) {
	case string:
		convertedValue, err = time.Parse(time.RFC3339, value)
		if err != nil {
			tz := timeUtils.EmptyLocalisedDateTime().Format("Z07:00")
			extendedValue := fmt.Sprintf("%sT00:00:00%s", value, tz)
			convertedValue, err = time.Parse(time.RFC3339, extendedValue)
		}

		if err == nil {
			convertedValue = convertedValue.In(timeUtils.GetLocation())
		}
	case time.Time:
		convertedValue = value.In(timeUtils.GetLocation())
	case *time.Time:
		_time := *value
		convertedValue = _time.In(timeUtils.GetLocation())
	}

	return convertedValue, err
}

func ToStringArray(value any) ([]string, error) {
	if value == nil {
		return []string{}, nil
	}

	convertedArray := []string{}

	switch value := value.(type) {
	case []any:
		for _, item := range value {
			convertedArray = append(convertedArray, item.(string))
		}
	case []string:
		convertedArray = value
	}

	return convertedArray, nil
}

func ToIntArray(value any) ([]int, error) {
	if value == nil {
		return []int{}, nil
	}

	convertedArray := []int{}

	switch value := value.(type) {
	case []any:
		for _, item := range value {
			num, err := ToInt(item)
			if err != nil {
				return convertedArray, err
			}
			convertedArray = append(convertedArray, num)
		}
	case []int:
		convertedArray = value
	}

	return convertedArray, nil
}

func ToUInt8Array(value any) ([]uint8, error) {
	if value == nil {
		return []uint8{}, nil
	}

	convertedArray := []uint8{}

	switch value := value.(type) {
	case []any:
		for _, item := range value {
			num, err := ToUint8(item)
			if err != nil {
				return convertedArray, err
			}
			convertedArray = append(convertedArray, num)
		}
	case string:
		strValue, err := ToString(value)
		if err != nil {
			return convertedArray, err
		}
		return []uint8(strValue), nil
	case []uint8:
		convertedArray = value
	}

	return convertedArray, nil
}

func ToTimePtr(value any) (*time.Time, error) {
	if isNilInterface(value) {
		return nil, nil
	}

	convertedValue, err := ToTime(value)
	return &convertedValue, err
}

func ToStringPtr(value any) (*string, error) {
	if isNilInterface(value) {
		return nil, nil
	}

	convertedValue, err := ToString(value)
	return &convertedValue, err
}

func ToIntPtr(value any) (*int, error) {
	if isNilInterface(value) {
		return nil, nil
	}

	convertedValue, err := ToInt(value)
	return &convertedValue, err
}

func ToUInt64Ptr(value any) (*uint64, error) {
	if isNilInterface(value) {
		return nil, nil
	}

	convertedValue, err := ToUint64(value)
	return &convertedValue, err
}

func ToJsonMap(value any) ([]map[string]any, error) {
	if value == nil {
		return []map[string]any{}, nil
	}

	convertedMap := []map[string]any{}

	switch value := value.(type) {
	case []any:
		for _, item := range value {
			convertedItem := item.(map[string]any)
			convertedMap = append(convertedMap, convertedItem)
		}
	case []map[string]any:
		convertedMap = value
	}

	return convertedMap, nil
}

func isNilInterface(value any) bool {
	iv := reflect.ValueOf(value)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
