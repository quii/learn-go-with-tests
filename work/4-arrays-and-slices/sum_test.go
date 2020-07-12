package arrays

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {

	// t.Run("collection of 5 numbers", func(t *testing.T) {
	// 	numbers := []int{1, 2, 3, 4, 5}

	// 	got := Sum(numbers)
	// 	want := 15

	// 	if got != want {
	// 		t.Errorf("got %d want %d, given %v", got, want, numbers)
	// 	}
	// })

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4}

		got := Sum(numbers)
		want := 10

		if got != want {
			t.Errorf("got %d want %d, given %v", got, want, numbers)
		}
	})
}

func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSumAllTails(t *testing.T) {

	checkSums := func(t *testing.T, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("sum some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 3}, []int{4, 7})
		want := []int{5, 7}

		checkSums(t, got, want)
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{4, 7})
		want := []int{0, 7}

		checkSums(t, got, want)

	})

}
