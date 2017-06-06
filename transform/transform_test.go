package transform

import (
	"testing"
)

func TestRemoveQuotes(t *testing.T) {
	var a = "\"test1\",\"test2\",\"tes\"t3\""
	var mockResult = "\"test1\",\"test2\",\"tes t3\""
	a = RemoveQuotesInsideString(a)

	if a != mockResult {
		t.Error("Error or TestRemoveQuotes: ", a)
	}
}

func TestTransformRemoveParenthesis(t *testing.T) {
	var a = "bus_zipcode (5)"
	var mockResult = "bus_zipcode  5 "
	a = RemoveParenthesis(a)

	if a != mockResult {
		t.Error("TestTransformRemoveParenthesis: ", a)
	}
}

func TestTransformRemoveSeparator(t *testing.T) {
	var a = "\"test1\",\"test2\",\"tes, t3\""
	var mockResult = "\"test1\",\"test2\",\"tes t3\""
	a = RemoveSeparatorInsideString(a)

	if a != mockResult {
		t.Error("TestTransformRemoveSeparator: ", a)
	}
}
