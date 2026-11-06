package stringUtils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	STRINGS_NOT_FOUND = -1
)

func ToTitleCase(userString string) string {
	return cases.Title(language.English, cases.NoLower).String(strings.ToLower(userString))
}

// Example of usage.. I want to extract all the text within the string "abc<de>fg<hi>jkl<m>" that is within
// the < and >. And I want also that the extracted text remains within the < and >. Finally, I want that
// each extracted text will be separated by the comma. Input shall be:
// - text 						-> "abc<de>fg<hi>jkl"  	... the original string
// - preSep 					-> "<"                 	... since I want to extract text that starts with <
// - postSep 					-> ">"				 	... since I want to extract text that ends with >
// - resultSep 					-> ","			 		... since I want that have each extracted text separated by the comma
// - includePrePostSeps 		-> true			 		... since I want the < and > to be extracted as well
// Result will be <de>,<hi>,<m>
func ExtractFrom(text string, preSep string, postSep string, resultSep string, includePrePostSeps bool) string {
	extractedText := ""
	lenPreSep := len(preSep)
	lenPostSep := len(postSep)
	preIndex := 0

	for preIndex != STRINGS_NOT_FOUND {
		preIndex = IndexOf(text, preSep, preIndex)

		if preIndex != STRINGS_NOT_FOUND {
			postIndex := IndexOf(text, postSep, preIndex+lenPreSep)

			if postIndex != STRINGS_NOT_FOUND {
				if extractedText != "" {
					extractedText += resultSep
				}
				if includePrePostSeps {
					extractedText += text[preIndex : postIndex+lenPostSep]
				} else {
					extractedText += text[preIndex+lenPreSep : postIndex]
				}
			}

			preIndex++
		}
	}

	return extractedText
}

func IndexOf(text string, substring string, fromIndex int) int {
	if substring == "" {
		return 0
	}

	if fromIndex < 0 {
		return STRINGS_NOT_FOUND
	}

	subLen := len(substring)
	charMatched := 0
	startingIndex := STRINGS_NOT_FOUND

	for i := fromIndex; i < len(text); i++ {
		if substring[charMatched] == text[i] {
			if startingIndex == STRINGS_NOT_FOUND {
				startingIndex = i
			}

			charMatched++
			if charMatched == subLen {
				return startingIndex
			}
		} else {
			startingIndex = STRINGS_NOT_FOUND
			charMatched = 0
		}
	}

	return STRINGS_NOT_FOUND
}

// Example of usage.. I want to get the text before the first slash (/). Input shall be:
// - text 	-> "photos/BC230008/EE230003/resized/Back Panels View.jpg"  	... the original string
// - search -> "/"	... since I want to get all the text before the first slash (/)
// Result will be photos
func BeforeFirstSubstring(text, search string) string {
	firstIndex := strings.Index(text, search)

	if firstIndex < 0 {
		return text
	}

	return text[:firstIndex]
}

// Example of usage.. I want to get the text after the last slash (/). Input shall be:
// - text 	-> "photos/BC230008/EE230003/resized/Back Panels View.jpg"  	... the original string
// - search -> "/"	... since I want to get all the text after the last slash (/)
// Result will be Back Panels View.jpg
// Note: Returns empty string if search not found
func AfterLastSubstring(text, search string) string {
	lastIndex := strings.LastIndex(text, search)
	if lastIndex < 0 {
		return ""
	}

	return text[lastIndex+len(search):]
}

// Assumption: the after & before values exist only once in the string
func SubstringWithin(text, after, before string) string {
	return BeforeFirstSubstring(AfterLastSubstring(text, after), before)
}

func ExtractDigits(text string) string {
	regex := regexp.MustCompile(`\d[\d,]*[\.]?[\d{2}]*`)

	submatchall := regex.FindAllString(text, -1)
	digits := ""
	for _, element := range submatchall {
		digits = fmt.Sprintf("%s%s", digits, element)
	}

	return digits
}

func ReplaceFirstRune(text, replacement string) string {
	return string([]rune(text)[:0]) + replacement + string([]rune(text)[1:])
}

func ToFloat64(text string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(text), 64)
}
