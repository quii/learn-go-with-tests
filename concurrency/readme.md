## Outline / notes

- Create a trivial, but paralizable, example - i.e. fetching multiple webpages
- Introduce goroutines to improve the performance
- Note that the test now fails as we can't be sure when the goroutines return
- Make tests pass using a sleep, but demonstrate that this doesn't scale
- introduce Channels for synchronization between processes

# Concurrency

Your colleague, Jo, has written a function in Go that checks whether a webpage
is working or not. It's called `IsWebsiteOK`

```go
package concurrency

import (
	"net/http"
)

// IsWebsiteOK returns true if the URL returns a 200 status code, false otherwise
func IsWebsiteOK(url string) bool {
	response, err := http.Head(url)
	if err != nil {
		return false
	}

	if response.StatusCode != http.StatusOK {
		return false
	}

	return true
}
```

If the above isn't familiar to you, don't worry about it. The key thing for this
exercise is that the above function _may_ take some time to return a result.

Jo's function is great, but you've been asked to make a version that takes
_multiple_ URLs and returns a list with the result that `IsWebsiteOK` gives for
each one.

## Test One

We'd like a function that takes a `slice` of `string`s and returns a `slice` of
`bool`s. For each URL that returns a `true`, the returned slice will have
a `true` at that index - and vice versa for `false`s

We'll use a list of three URLs for now; the first two we know _should_ work; the
last one shouldn't work.

The first test checks that `websiteChecker` returns the same number of results
as websites; the second test checks that the results are what we expect for each
website.

```
package concurrency

import "testing"

func TestWebsiteChecker(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"http://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		true,
		false,
	}

	actualResults := websiteChecker(websites)

	want := len(websites)
	got := len(actualResults)
	if len(actualResults) != len(websites) {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	for index, want := range expectedResults {
		got := actualResults[index]
		if want != got {
			t.Fatalf("Wanted %v, got %v", want, got)
		}
	}
}
```

When we run the tests we see

```sh
# github.com/gypsydave5/learn-go-with-tests/concurrency/v1
./websiteChecker_test.go:18:10: undefined: websiteChecker
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1 [build failed]
```

### Making the build pass

To get this to pass we need to implement `websiteChecker` with the correct type
signature - a function that takes a single argument of a slice of `string`s
(`[]string`) and returns a slice of `bool`s (`[]bool`).

The simplest implementation of this is:
```go
func websiteChecker(_ []string) []bool {
	return []bool{}
}
```

Now when we run the tests we get
```sh
--- FAIL: TestWebsiteChecker (0.00s)
        websiteChecker_test.go:23: Wanted 3, got 0
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s
#... etc
```

The code now builds (hooray), but the first test fails because the length of the
slice returned is too short. This is easy enough to fix by defining a length for
the slice:

```go
func websiteChecker(_ []string) []bool {
	return make([]bool, 3)
}
```

Now when we run the tests we get:
```sh
--- FAIL: TestWebsiteChecker (0.00s)
        websiteChecker_test.go:29: Wanted true, got false
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.013s
```

Which is OK. So now we have to do some work in getting the right results!

### Making the test pass

To make the test pass we will iterate through the slice of URLs using a `for...
range` loop. For each URL we will call `IsWebsiteOK` with the URL and then store
the answer in a `results` slice which we'll create at the top of our function.
The `results` slice will be the same length as the `urls` slice, so we'll save
the response from `IsWebsiteOK` at the same index as the `url` we're checking.

When we've checked all of the URLs we'll finally return the `results` slice.

```go
func websiteChecker(urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = IsWebsiteOK(url)
	}

	return results
}
```

Now when we run our tests:
```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.269s
```

## An observation

Did you notice that the time it took for the `websiteCheckerTest` to run
increased dramatically when we were really checking websites? It added around
a quarter of a second to the total time. Although the Internet is fast, and the
response we're getting from the websites is coming back pretty quickly, it still
takes time for our functions to make those real requests.
