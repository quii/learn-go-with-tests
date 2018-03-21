package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(websiteChecker WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = websiteChecker(url)
	}

	return results
}
