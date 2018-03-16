package concurrency

import "testing"

func fakeIsWebsiteOK(url string) bool {
	if url == "http://blog.gypsydave5.com" {
		return false
	}
	return true
}

func TestWebsiteChecker(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"http://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		false,
		true,
	}

	actualResults := websiteChecker(fakeIsWebsiteOK, websites)

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
