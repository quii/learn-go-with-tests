package concurrency

func CheckWebsites(urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = CheckWebsite(url)
	}

	return results
}
