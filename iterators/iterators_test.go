package iterators

import (
	"iter"
	"slices"
	"testing"
)

func Concatenate(seq iter.Seq[string]) string {
	var result string
	for s := range seq {
		result += s
	}
	return result
}

// annoyingly, there is no builtin way to go from seq2, to seq (e.g just get the values)
func Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// WIP!
func TestConcatenate(t *testing.T) {
	t.Run("values of a slice", func(t *testing.T) {
		got := Concatenate(slices.Values([]string{"a", "b", "c"}))
		want := "abc"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("values of a slice backwards", func(t *testing.T) {
		backward := slices.Backward([]string{"a", "b", "c"})

		got := Concatenate(Values(backward))
		want := "cba"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("values of a slice sorted", func(t *testing.T) {
		got := Concatenate(slices.Values(slices.Sorted(slices.Values([]string{"c", "a", "b"}))))
		want := "abc"
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
