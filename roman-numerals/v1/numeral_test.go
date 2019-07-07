package v1

import "testing"

func TranslateToRoman(number int) string {
	return "I"
}

func TestRomanNumberals(t *testing.T) {

	cases := []struct {
		Name   string
		Number int
		Want   string
	}{
		{"One", 1, "I"},
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := TranslateToRoman(test.Number)

			if got != test.Want {
				t.Errorf("got '%s' want '%s'", got, test.Want)
			}
		})
	}
}
