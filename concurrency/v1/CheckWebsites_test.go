package concurrency

import (
	"reflect"
	"testing"
)

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	actualResults := CheckWebsites(websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	expectedResults := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	if !reflect.DeepEqual(expectedResults, actualResults) {
		t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
	}
}

func BenchmarkCheckWebsites(b *testing.B) {
	websites := make([]string, 100)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.com"
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(websites)
	}
}
