package concurrency

import (
	"testing"
	"time"
)

func slowWebsiteChecker(url string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	websites := make([]string, 100)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.com"
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(slowWebsiteChecker, websites)
	}
}
