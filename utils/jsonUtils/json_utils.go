package jsonUtils

import "Quote/utils/arrayUtils"

// Returns a new json having only the attributes from json parameter that are within the attributeList parameter
func Json_ExtractOnly(json map[string]any, attributeList []string) map[string]any {
	result := map[string]any{}

	for key, value := range json {
		if arrayUtils.ContainsValue(attributeList, key) {
			result[key] = value
		}
	}

	return result
}
