package concurrency

import "testing"

func TestWebsiteChecker(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"http://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		true,
		false,
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
