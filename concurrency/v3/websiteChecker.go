package concurrency

type URLchecker func(string) bool
type Result struct {
	index   int
	success bool
	url     string
}

func WebsiteChecker(isOK URLchecker, urls []string) []bool {
	results := make([]bool, len(urls))
	urlChan := ofStrings(urls)
	resultsChan := checkEach(urlChan, isOK)

	resultsCount := 0
	for resultsCount != len(urls) {
		r := <-resultsChan
		resultsCount += 1
		results[r.index] = r.success
	}

	return results
}

func checkEach(urls <-chan Result, f URLchecker) <-chan Result {
	out := make(chan Result)

	for r := range urls {
		go func(r Result) {
			out <- Result{r.index, f(r.url), r.url}
		}(r)
	}

	return out
}

func ofStrings(urls []string) <-chan Result {
	c := make(chan Result)

	go func() {
		for index, url := range urls {
			c <- Result{
				index,
				false,
				url,
			}
		}
		close(c)
	}()

	return c
}
