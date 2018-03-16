package concurrency

func WebsiteChecker(urls []string) (results []bool) {
	for _, url := range urls {
		results = append(results, IsWebsiteOK(url))
	}

	return
}
