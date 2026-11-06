// Package int64Array provides a Postgres bigint[] compatible slice type.
package int64Array

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Int64Array maps to a Postgres BIGINT[] column. It serialises to JSON as a plain array of numbers.
type Int64Array []int64

// GormDataType tells gorm which column type to use.
func (Int64Array) GormDataType() string {
	return "bigint[]"
}

// Value implements the driver Valuer interface. A nil array is stored as an empty array.
func (a Int64Array) Value() (driver.Value, error) {
	var sb strings.Builder
	sb.WriteByte('{')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('}')
	return sb.String(), nil
}

// Scan implements the sql Scanner interface. It accepts the Postgres text form, e.g. "{900,900}".
func (a *Int64Array) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		*a = Int64Array{}
		return nil
	case string:
		return a.parse(src)
	case []byte:
		return a.parse(string(src))
	case []int64:
		*a = append(Int64Array{}, src...)
		return nil
	case []any:
		result := make(Int64Array, 0, len(src))
		for _, item := range src {
			switch v := item.(type) {
			case int64:
				result = append(result, v)
			case int32:
				result = append(result, int64(v))
			case int:
				result = append(result, int64(v))
			case float64:
				result = append(result, int64(v))
			default:
				return fmt.Errorf("cannot convert %T to int64", item)
			}
		}
		*a = result
		return nil
	}
	return fmt.Errorf("cannot scan %T into Int64Array", src)
}

func (a *Int64Array) parse(text string) error {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "{")
	text = strings.TrimSuffix(text, "}")

	result := Int64Array{}
	if strings.TrimSpace(text) != "" {
		for _, part := range strings.Split(text, ",") {
			v, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
			if err != nil {
				return err
			}
			result = append(result, v)
		}
	}
	*a = result
	return nil
}

// Collapsed joins the values with sep, e.g. "900, 905". A uniform array collapses to its single value, e.g. "900".
func (a Int64Array) Collapsed(sep string) string {
	if len(a) == 0 {
		return "0"
	}
	uniform := true
	for _, v := range a {
		if v != a[0] {
			uniform = false
			break
		}
	}
	if uniform {
		return strconv.FormatInt(a[0], 10)
	}
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(parts, sep)
}

// AreaSqm returns the summed area in square metres of the units described by the given widths and heights (mm).
// Widths and heights are paired by index; any unpaired trailing values are ignored.
func AreaSqm(widths, heights Int64Array) float64 {
	n := min(len(widths), len(heights))
	total := 0.0
	for i := 0; i < n; i++ {
		total += float64(widths[i]) * float64(heights[i])
	}
	return total / 1000000.0
}
