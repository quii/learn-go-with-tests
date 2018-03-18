package concurrency

import (
	"testing"
)

func TestWebsiteChecker(t *testing.T) {

	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	actualResults := WebsiteChecker(websites)

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

	assertSameResults(t, expectedResults, actualResults)
}

func assertSameResults(t *testing.T, expectedResults, actualResults map[string]bool) {
	for expectedKey, expectedValue := range expectedResults {
		actualValue, ok := actualResults[expectedKey]
		if !ok {
			t.Fatalf("actual results did not contain expected key: '%s'", expectedKey)
		}
		if actualValue != expectedValue {
			t.Fatalf("expected value of key '%s' in actual results to be '%v', but it was '%v'", expectedKey, expectedValue, actualValue)
		}
	}

	for actualKey, _ := range actualResults {
		if _, ok := expectedResults[actualKey]; !ok {
			t.Fatalf("found unexpected key in actual results: '%s'", actualKey)
		}
	}
}
