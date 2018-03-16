package concurrency

import "time"

type URLchecker func(string) bool

func WebsiteChecker(isOK URLchecker, urls []string) (results []bool) {
	for _, url := range urls {
		go func(url string) {
			results = append(results, isOK(url))
		}(url)
	}

	time.Sleep(2 * time.Second)

	return
}
