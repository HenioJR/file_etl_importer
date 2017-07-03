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

func RemoveSpecialCharactersFromHeader(s string) string {
	s = strings.Replace(s, "/", " ", -1)
	s = strings.Replace(s, "\\", " ", -1)
	return s
}

func RemoveSeparatorInsideString(s string) string {
	// change to use separator of config file
	//c := config.GetConfig()
	//separator := c.File.Separator
	re1 := regexp.MustCompile(`\b,\b`)
	//fmt.Println("separator: ", separator)
	//regex1 := `\b` + separator + `\b`
	//fmt.Println(regex1)
	//re1 := regexp.MustCompile(regex1)
	re2 := regexp.MustCompile(`\s,`)
	re3 := regexp.MustCompile(`,\s`)
	s = re1.ReplaceAllString(s, " ")
	s = re2.ReplaceAllString(s, " ")
	s = re3.ReplaceAllString(s, " ")

	return s
}
