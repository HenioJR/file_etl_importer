package transform

import (
	"testing"
)

func TestRemoveQuotes(t *testing.T) {
	var a = "bus_\"zipcode\""
	var mockResult = "bus_zipcode"
	a = RemoveQuotes(a)

	if a != mockResult {
		t.Error("Error or transform.RemoveQuotes")
	}
}

func TestTransformRemoveParenthesis(t *testing.T) {
	var a = "bus_zipcode (5)"
	var mockResult = "bus_zipcode  5 "
	a = RemoveParenthesis(a)

	if a != mockResult {
		t.Error("Error or transform.RemoveParenthesis")
	}
}
