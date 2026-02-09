package mysolution

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Repeat(t *testing.T) {
	repeated := Repeat("a")
	expected := "aaaaa"

	assert.Equal(t, expected, repeated)
}

func Test_Repeat2(t *testing.T) {
	repeated := Repeat("b")
	expected := "bbbbb"

	assert.Equal(t, expected, repeated)
}

func Test_RepeatRandomCharacters(t *testing.T) {

	testCases := map[string]struct {
		input          string
		expectedString string
	}{
		"a": {
			input:          "a",
			expectedString: "aaaaa",
		},
		"b": {
			input:          "b",
			expectedString: "bbbbb",
		},
		"cc": {
			input:          "cc",
			expectedString: "cccccccccc",
		},
		"hej": {
			input:          "hej",
			expectedString: "hejhejhejhejhej",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			value := Repeat(tc.input)

			assert.Equal(t, tc.expectedString, value)
		})
	}

}
