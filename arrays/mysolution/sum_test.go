package mysolution

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func Test_Sum(t *testing.T) {
	t.Parallel()

	numbers := [5]int{1, 2, 3, 4, 5} //specifically testing a slice of size 5

	got := Sum(numbers)

	assert.Equal(t, 15, got)
}
*/

/*
func Test_Sum(t *testing.T) {
	t.Parallel()

	numbers := []int{1, 2, 3, 4, 5} //an arbitrary sized slice

	got := Sum(numbers)

	assert.Equal(t, 15, got)
}
*/

func Test_Sum(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		numbers  []int
		expected int
	}{
		"slice of 5 numbers": {
			numbers: []int{
				1, 2, 3, 4, 5,
			},
			expected: 15,
		},

		"slice of 10 numbers": {
			numbers: []int{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			},
			expected: 55,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			result := Sum(testCase.numbers)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func Test_SumAll(t *testing.T) {
	/* function looks like this
	func SumAll(numbersOne []int, numbersTwo []int) []int {
		return []int{}
	}
	==> Test fails
	*/

	/* after refactoring our function by using minimal amount of code for test to run
	func SumAll(numbersOne []int, numbersTwo []int) []int {

		sumOne := Sum(numbersOne)
		sumTwo := Sum(numbersTwo)

		return []int{
			sumOne,
			sumTwo,
		}
	}
	==> Test passes
	*/

	t.Parallel()

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	assert.Equal(t, want, got)
}

/*
This should test for any arbitrary amount of numbers in the argument
For this we need a variadic function!
*/
func Test_SumAll2(t *testing.T) {
	/* function looks like this
	func SumAll(numbersOne []int, numbersTwo []int) []int {
		return []int{}
	}
	==> Test does not even compile
	*/

	/* after refactoring our function by using minimal amount of code for test to run
	func SumAll(numbersOne []int, numbersTwo []int) []int {
		return []int{}
	}
	==> Test fails
	*/

	t.Parallel()

	testCases := map[string]struct {
		slices   [][]int
		expected []int
	}{
		"one slice": {
			slices: [][]int{
				{1, 2, 3},
			},
			expected: []int{6},
		},
		"two slices": {
			slices: [][]int{
				{1, 2, 3, 4, 5},
				{0, 9},
			},
			expected: []int{15, 9},
		},
		"three slices": {
			slices: [][]int{
				{1, 2, 3, 4, 5},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			expected: []int{15, 55},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			sum := SumAll(testCase.slices...)

			assert.Equal(t, testCase.expected, sum)
		})
	}

}
