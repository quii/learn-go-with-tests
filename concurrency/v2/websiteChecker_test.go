package concurrency

import "testing"

func TestWebsiteChecker(t *testing.T) {
	websites := make([]string, 500)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.co.uk"
	}

	expectedResults := make([]bool, len(websites))
	for i := 0; i < len(websites); i++ {
		expectedResults[i] = true
	}

	actualResults := websiteChecker(websites)

	want := len(websites)
	got := len(actualResults)
	if len(actualResults) != len(websites) {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	for index, want := range expectedResults {
		got := actualResults[index]
		if want != got {
			t.Fatalf("Wanted %v, got %v", want, got)
		}
	}
}
