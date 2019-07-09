package v1

import "testing"

func TestRomanNumerals(t *testing.T) {
	got := ConvertToRoman(1)
	want := "I"

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func ConvertToRoman(arabic int) string {
	return "I"
}
