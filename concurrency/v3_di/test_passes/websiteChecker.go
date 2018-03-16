package concurrency

func websiteChecker(isOK func(string) bool, urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = isOK(url)
	}

	return results
}
