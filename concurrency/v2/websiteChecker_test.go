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
		"waat://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		false,
		true,
	}

	actualResults := WebsiteChecker(fakeIsWebsiteOK, websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	if !sameResults(expectedResults, actualResults) {
		t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
	}
}

func sameResults(as, bs []bool) bool {
	for index, a := range as {
		if a != bs[index] {
			return false
		}
	}
	return true
}
