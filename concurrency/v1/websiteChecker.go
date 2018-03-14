package concurrency

func websiteChecker(urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = IsWebsiteOK(url)
	}

	return results
}
