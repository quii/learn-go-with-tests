package iteration

import (
	"fmt"
	"testing"
)

func ExampleRepeat() {
	got := Repeat("ab", 2)
	fmt.Println(got)
	// Output: abab
}

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 5)
	expected := "aaaaa"

	if repeated != expected {
		t.Errorf("Expected %q but got %q", expected, repeated)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}
