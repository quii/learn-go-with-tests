package concurrency

import (
	"testing"
	"time"
)

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

func BenchmarkWebsiteCheckerWithManyURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		websites := make([]string, 100)
		for index, _ := range websites {
			websites[index] = "http://google.co.uk"
		}

		WebsiteChecker(slowIsWebsiteOK, websites)
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

func slowIsWebsiteOK(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func fakeIsWebsiteOK(url string) bool {
	if url == "http://blog.gypsydave5.com" {
		return false
	}
	return true
}
