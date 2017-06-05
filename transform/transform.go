package transform

import (
	"strings"
)

func RemoveQuotes(s string) string {
	s = strings.Replace(s, "\"", "", -1)
	return s
}

func RemoveParenthesis(s string) string {
	s = strings.Replace(s, "(", " ", -1)
	s = strings.Replace(s, ")", " ", -1)
	return s
}
