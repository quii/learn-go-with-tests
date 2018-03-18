package concurrency

import "time"

type TestURL func(string) bool

func WebsiteChecker(isOK TestURL, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func(u string) {
			results[u] = isOK(u)
		}(url)
	}

	time.Sleep(2 * time.Second)

	return results
}
