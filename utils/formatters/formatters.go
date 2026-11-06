package formatters

import (
	"Quote/utils/stringUtils"
	"fmt"
	"strings"
)

func FormatPhone(phone string) string {
	result := phone

	if len(result) >= 8 {
		if result[0] == '+' || result[0] == ' ' {
			result = fmt.Sprintf("00%s", phone[1:])
		}

		result = stringUtils.ExtractDigits(result)

		if len(result) == 8 {
			result = fmt.Sprintf("00356%s", result)
		}
	}

	return result
}

func ToReadableWholeNumber(number int) string {
	// Format to 2 decimal places
	formatted := fmt.Sprintf("%d", number)

	// Split into integer and decimal parts
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]

	// Add commas to the integer part
	var result strings.Builder
	for i, digit := range integerPart {
		result.WriteRune(digit)
		if (len(integerPart)-i-1)%3 == 0 && i != len(integerPart)-1 {
			result.WriteString(",")
		}
	}
	return result.String()
}

func ToReadableNumber_2dp(number float64) string {
	// Format to 2 decimal places
	formatted := fmt.Sprintf("%.2f", number)

	// Split into integer and decimal parts
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = "." + parts[1]
	}

	// Add commas to the integer part
	var result strings.Builder
	for i, digit := range integerPart {
		result.WriteRune(digit)
		if (len(integerPart)-i-1)%3 == 0 && i != len(integerPart)-1 {
			result.WriteString(",")
		}
	}
	return result.String() + decimalPart
}
