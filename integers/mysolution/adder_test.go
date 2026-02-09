package mysolution

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		x, y        int
		expectedSum int
	}{
		"OnePlusOne": {
			1,
			1,
			2,
		},
		"OnePlusTwo": {
			1,
			2,
			3,
		},
		"ThreePlusTwo": {
			3,
			2,
			5,
		},
		"123Plus234": {
			123,
			234,
			357,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expectedSum, Add(tc.x, tc.y))
		})
	}
}
