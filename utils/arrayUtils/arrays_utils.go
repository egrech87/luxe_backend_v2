package arrayUtils

import (
	"Quote/utils/anyUtils"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
)

const (
	ALL_ITEMS int8 = iota
	ANY_ITEMS
)

func StringArrayToUint64Array(strArray []string, restriction int8) (uint64Array []uint64, err error) {
	uint64Array = []uint64{}

	for _, valueS := range strArray {
		value, err := strconv.ParseUint(strings.Trim(valueS, " "), 10, 64)
		if err != nil {
			if restriction == ALL_ITEMS {
				return []uint64{}, errors.Errorf(fmt.Sprintf("not all items are positive whole numbers: %v", err.Error()))
			}
		} else {
			uint64Array = append(uint64Array, value)
		}
	}

	return uint64Array, nil
}

func StringArrayToIntArray(strArray []string, restriction int8) (intArray []int, err error) {
	intArray = []int{}

	for _, valueS := range strArray {
		value, err := strconv.Atoi(strings.Trim(valueS, " "))
		if err != nil {
			if restriction == ALL_ITEMS {
				return []int{}, errors.Errorf(fmt.Sprintf("not all items are whole numbers: %v", err.Error()))
			}
		} else {
			intArray = append(intArray, value)
		}
	}

	return intArray, nil
}

func ContainsValues[E comparable](array []E, values []E) bool {
	for _, value := range values {
		valueFound := ContainsValue(array, value)

		if !valueFound {
			return false
		}
	}

	return true
}

func ContainsValue[E comparable](array []E, value E) bool {
	for _, a := range array {
		if value == a {
			return true
		}
	}

	return false
}

// returns true if the arrays have the same items in any order
func Equal[E comparable](array1 []E, array2 []E) bool {
	if len(array1) != len(array2) || !ContainsValues(array1, array2) {
		return false
	}

	return true
}

// returns true if the arrays have the same items in same order
func SequencialEqual[E comparable](array1 []E, array2 []E) bool {
	if len(array1) != len(array2) {
		return false
	}

	for i, v := range array1 {
		if v != array2[i] {
			return false
		}
	}

	return true
}

func ToString[T constraints.Integer | constraints.Float | string | bool](a []T, delim string) string {
	result := ""

	totalItems := len(a)
	for i, item := range a {
		itemStr, err := anyUtils.ToString(any(item))
		if err != nil {
			itemStr = ""
		}
		result += itemStr

		if i < totalItems-1 {
			result += delim
		}
	}

	return result
}

// trims the items and removes any empty items
func FromString(text string, delim rune) []string {
	splitByCommaFn := func(c rune) bool {
		return c == delim
	}

	array := strings.FieldsFunc(strings.TrimSpace(text), splitByCommaFn)

	for i, s := range array {
		array[i] = strings.Trim(s, " ")
	}

	return array
}

func AppendUnique[E comparable](array []E, value E) ([]E, bool) {
	if ContainsValue(array, value) {
		return array, false
	}

	return append(array, value), true
}

func Clear[E any](array []E, target *[]E) {
	*target = array[:0]
}

// Removes an element from the array at the given index. The original array is not changed. Result is returned by the parameter "target".
func RemoveItemAt[E any](array []E, target *[]E, index int) error {
	if index < 0 || index > len(array) {
		return errors.Errorf("Invalid array index")
	}

	*target = append(array[:index], array[index+1:]...)
	return nil
}

// Removes multiple elements from the array at the given indexes. The original array is not changed. Result is returned by the parameter "target".
func RemoveItemsAt[E any](array []E, target *[]E, indexes []int) error {
	size := len(indexes)
	if size == 0 {
		return errors.Errorf("Failed to remove from the array since no index is given")
	}

	slices.Sort(indexes)
	if indexes[0] < 0 || indexes[size-1] >= len(array) {
		return errors.Errorf("Invalid array indexes")
	}

	*target = array[:indexes[0]]
	for i := 1; i < size; i++ {
		*target = append(*target, array[indexes[i-1]+1:indexes[i]]...)
	}
	*target = append(*target, array[indexes[size-1]+1:]...)

	return nil
}

// Removes the given element "item" from the array. The original array is not changed. Result is returned by the parameter "target".
func RemoveItem[E comparable](array []E, target *[]E, item E) (index int, err error) {
	for i, itemInArray := range array {
		if itemInArray == item {
			return i, RemoveItemAt(array, target, i)
		}
	}

	return -1, nil
}

// Removes the given element "items" from the array. The original array is not changed. Result is returned by the parameter "target".
func RemoveItems[E comparable](array []E, target *[]E, items []E) ([]int, error) {
	indexes := []int{}
	for _, item := range items {
		for i, itemInArray := range array {
			if itemInArray == item {
				indexes = append(indexes, i)
				break
			}
		}
	}

	if err := RemoveItemsAt(array, target, indexes); err != nil {
		return indexes, err
	}

	return indexes, nil
}

// Fills the target array with items that exist in both array1 and array2
func SetIntersection[E comparable](array1, array2 []E, target *[]E) {
	set := make(map[E]bool)

	// Create a set from the first array
	for _, item := range array1 {
		set[item] = true // setting the initial value to true
	}

	// Check elements in the second array against the set
	for _, item := range array2 {
		if set[item] {
			*target = append(*target, item)
		}
	}
}

// Removes all items from the array that end with the given postfix. By default, caseSensitive is false
func RemoveItemsHavingSuffix(array *[]string, suffix string, caseSensitive ...bool) error {
	if len(caseSensitive) > 1 {
		return errors.Errorf("Case sensitive parameter can hold one value only")
	}

	isCaseSensitive := false
	if len(caseSensitive) == 1 {
		isCaseSensitive = caseSensitive[0]
	}

	indexes := []int{}

	for i, item := range *array {
		if isCaseSensitive {
			item = strings.ToLower(item)
		}

		if strings.HasSuffix(item, suffix) {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) > 0 {
		if err := RemoveItemsAt(*array, array, indexes); err != nil {
			return err
		}
	}

	return nil
}

func Reverse[E comparable](array []E) []E {
	for i := 0; i < len(array)/2; i++ {
		j := len(array) - i - 1
		array[i], array[j] = array[j], array[i]
	}
	return array
}
