package errortypes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BadStatusError struct {
	URL    string
	Status int
}

func (b BadStatusError) Error() string {
	return fmt.Sprintf("did not get 200 from %s, got %d", b.URL, b.Status)
}

// DumbGetter will get the string body of url if it gets a 200
func DumbGetter(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("problem fetching from %s, %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return "", BadStatusError{URL: url, Status: res.StatusCode}
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body) // ignoring err for brevity

	return string(body), nil
}

func TestDumbGetter(t *testing.T) {

	t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

		svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusTeapot)
		}))
		defer svr.Close()

		_, err := DumbGetter(svr.URL)

		if err == nil {
			t.Fatal("expected an error")
		}

		got, isStatusErr := err.(BadStatusError)

		if !isStatusErr {
			t.Fatalf("was not a BadStatusError, got %T", err)
		}

		want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
