package concurrency

import "net/http"

// CheckWebsite returns true if the URL returns a 200 status code, false otherwise.
func CheckWebsite(url string) bool {
	response, err := http.Head(url)
	if err != nil {
		return false
	}

	return response.StatusCode == http.StatusOK
}
