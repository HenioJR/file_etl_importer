package transform

import (
	"testing"
)

func TestRemoveQuotesInsideString(t *testing.T) {
	var a = "\"test1\",\"test2\",\"tes\"t3\""
	var mockResult = "\"test1\",\"test2\",\"tes t3\""
	a = RemoveQuotesInsideString(a)

	if a != mockResult {
		t.Error("Error on TestRemoveQuotesInsideString: ", a, mockResult)
	}
}

func TestRemoveParenthesis(t *testing.T) {
	var a = "bus_zipcode (5)"
	var mockResult = "bus_zipcode  5 "
	a = RemoveParenthesis(a)

	if a != mockResult {
		t.Error("Error on TestRemoveParenthesis: ", a, mockResult)
	}
}

func TestRemoveSpecialCharactersFromHeader(t *testing.T) {
	var a = "cnpj/cpf-doc"
	var mockResult = "cnpj cpf doc"
	a = RemoveSpecialCharactersFromHeader(a)

	if a != mockResult {
		t.Error("Error on TestRemoveSpecialCharactersFromHeader: ", a, mockResult)
	}
}

func TestRemoveSeparatorInsideString(t *testing.T) {
	//c := config.GetConfig()
	//separator := c.File.Separator

	var a = "\"te,st1\",\"test2\",\"tes,t3\""
	var mockResult = "\"te st1\",\"test2\",\"tes t3\""
	a = RemoveSeparatorInsideString(a, ",")

	if a != mockResult {
		t.Error("Error on TestRemoveSeparatorInsideString: ", a, mockResult)
	}
}

func TestRemoveDoubleQuote(t *testing.T) {
	var a = "\"string with quotes\""
	var mockResult = "string with quotes"
	a = RemoveDoubleQuote(a)

	if a != mockResult {
		t.Error("Error on TestRemoveDoubleQuote: ", a, mockResult)
	}
}
