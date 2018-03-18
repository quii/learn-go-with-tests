package concurrency

type TestURL func(string) bool
type result struct {
	string
	bool
}

func WebsiteChecker(isOK TestURL, urls []string) map[string]bool {
	results := make(map[string]bool)
	urlChannel := make(chan string)
	resultChannel := make(chan result)

	go func() {
		for {
			url := <-urlChannel

			good := isOK(url)
			result := result{url, good}
			resultChannel <- result
		}
	}()

	go func() {
		for _, url := range urls {
			urlChannel <- url
		}
	}()

	for i := 0; i < len(urls); i++ {
		result := <-resultChannel
		results[result.string] = result.bool
	}

	return results
}
