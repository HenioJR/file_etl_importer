package transform

import (
	"regexp"
	"strings"
)

func RemoveQuotesInsideString(s string) string {
	s = strings.Replace(s, "'", "", -1)
	re := regexp.MustCompile(`\b"\b`)
	return re.ReplaceAllString(s, " ")
}

func RemoveParenthesis(s string) string {
	s = strings.Replace(s, "(", " ", -1)
	s = strings.Replace(s, ")", " ", -1)
	return s
}

func RemoveDoubleQuote(s string) string {
	s = strings.Replace(s, "\"", "", -1)
	return s
}

func RemoveSpecialCharactersFromHeader(s string) string {
	s = strings.Replace(s, "/", " ", -1)
	s = strings.Replace(s, "\\", " ", -1)
	s = strings.Replace(s, "-", " ", -1)
	return s
}

func RemoveSeparatorInsideString(line string, separator string) string {
	//only work when have quotes on string
	re1 := regexp.MustCompile(`\b` + separator + `\b`)
	re2 := regexp.MustCompile(`\s,`)
	re3 := regexp.MustCompile(`,\s`)
	line = re1.ReplaceAllString(line, " ")
	line = re2.ReplaceAllString(line, " ")
	line = re3.ReplaceAllString(line, " ")

	return line
}
