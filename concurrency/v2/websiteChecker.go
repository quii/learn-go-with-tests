package concurrency

type TestURL func(string) bool

func WebsiteChecker(isOK TestURL, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = isOK(url)
	}

	return results
}
