package concurrency

import (
	"fmt"
	"testing"
)

func TestMultipleWebsiteChecker(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"http://furhurterwe.geds",
	}

	wants := []bool{
		true,
		true,
		false,
	}

	gots := multipleWebsiteChecker(websites)
	fmt.Println(gots)

	for index, want := range wants {
		got := gots[index]

		if want != got {
			t.Errorf("Wanted %v, got %v", want, got)
		}
	}
}
