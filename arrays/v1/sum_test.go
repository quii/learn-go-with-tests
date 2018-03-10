package main

import "testing"

func TestSum(t *testing.T) {

	numbers := [5]int{1, 2, 3, 4, 5}

	expectedSum := 15
	actualSum := Sum(numbers)

	if expectedSum != actualSum {
		t.Errorf("expected the sum to be %d but was %d given, %v", expectedSum, actualSum, numbers)
	}
}
