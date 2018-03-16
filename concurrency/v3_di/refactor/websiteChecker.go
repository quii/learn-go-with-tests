package concurrency

type testURL func(string) bool

func websiteChecker(isOK testURL, urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = isOK(url)
	}

	return results
}
