package floatUtils

import (
	"Quote/utils/stringUtils"
	"regexp"
)

func ExtractFloat64(text string) (float64, error) {
	re := regexp.MustCompile(`\d+\.?\d+`)
	matches := re.FindStringSubmatch(text)

	if len(matches) == 0 {
		return 0, nil
	}

	return stringUtils.ToFloat64(matches[0])
}
