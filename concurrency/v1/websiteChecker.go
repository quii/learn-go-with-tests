package concurrency

func WebsiteChecker(urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = IsWebsiteOK(url)
	}

	return results
}
