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

## A brief encounter with Dependency Injection...

So far, so good. But there are already two problems with what we've got so far.

1. If `google.co.uk` goes down, (or someone puts a website at `http://furhurterwe.geds`), our expectations will be wrong for our tests.
2. If we turn off the Internet, our tests will always fail.

And in true TDD style, we should demonstrate this with a failing test. So, let's
turn off the computer's WiFi / unplug the network cable and run the tests again:

```sh
--- FAIL: TestWebsiteChecker (0.20s)
        websiteChecker_test.go:27: Wanted true, got false
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v2        0.214s
```

This dependency on the Internet is a bad thing because these failures will have
nothing to do with any changes to the behaviour of our code. More precisely, we
can say that the dependency is on Jo's function `IsWebsiteOK`. If that function
stops working for any reason at all - whether the network cuts out or Jo creates
a bug in her code - our code will stop working and our tests will fail.

To mitigate this problem we can make `IsWebsiteOK` an extra argument to our
`websiteChecker` function. Then, in the tests, we can use a different function
with the same interface as `IsWebsiteOK` that behaves in a way that we can
control in our tests.

Let's try it out. So, leaving the internet off, let's make some changes to
`websiteChecker_test.go`. We're going to change our original test like so

```go
package concurrency

import "testing"

func fakeIsWebsiteOK(url string) bool {
    if url == "http://blog.gypsydave5.com" {
        return false
    }
    return true
}

func TestWebsiteChecker(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"http://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		false,
		true,
	}

	actualResults := websiteChecker(fakeIsWebsiteOK, websites)

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

We've added a new function, `fakeIsWebsiteOK`, which has the same behaviour as
`IsWebsiteOK`. From the outside you couldn't tell the difference between them
they take a `string` and return a `bool`. But on the inside `fakeIsWebsiteOK`
is just an `if` statement that always returns `true` unless the `string`
argument is `"http://blog.gypsydave5.com"`. It's a function we have complete
control over - because we wrote it!

The expectations have also been updated; we now expect the middle one to fail.

The way we want this to work is for `websiteChecker` to take our
`fakeIsWebsiteOK` function as it's first argument and to use it to 'check' the
websites. So that's what we've written in the test

```go
  actualResults := websiteChecker(fakeIsWebsiteOK, websites)
```

If we run this we get

```sh
./websiteChecker_test.go:25:33: too many arguments in call to websiteChecker
        have (func(string) bool, []string)
        want ([]string)
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v3_di [build failed]
```

A faiure to compile. `websiteChecker` wants `([]string)` as it's arguments, but we gave it
`(func(string) bool, []string)`. So we've learnt that `func(string) bool` is the
type of our `fakeWebsiteOK` function in the same way as `[]string` is the type
of the slice of strings we're passing in as the second argument. We should now
be able to change `websiteChecker` to at least get the compilation to work.

```go
package concurrency

func websiteChecker(isOK func(string) bool, urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = IsWebsiteOK(url)
	}

	return results
}
```

I've named the function we've passed in `isOK`. Now if we run the tests:

```sh
websiteChecker_test.go:36: Wanted false, got true
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v3_di     0.241s
```

Which is because we're still not using the `isOK` function. This is easily fixed
by changing `IsWebsiteOK` to `isOK`:

```go
package concurrency

func websiteChecker(isOK func(string) bool, urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = isOK(url)
	}

	return results
}
```

And now...

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3_di     0.013s
```

#### Refactor

`func(string) bool` doesn't exactly trip off the tongue when trying to describe
what the function is doing - you can tell the behaviour, but it's hard to say
what the intention of it is. Happily in Go we can give a type an alias - like
a nickname we can use for a type. This will help us to remember that the
function we pass in is for checking a website's status.

```go
package concurrency

type testURL func(string) bool

func websiteChecker(isOK testURL, urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = isOK(url)
	}

	return results
}
```

We've used the `type` keyword to say that we'd like `func(string) bool` to also
be known as `testURL` in the rest of this package. This is a useful technique to
help your code read nicely, especially when some of the function types get
really, really long.


#### Note on Dependency Injection and Test Doubles

This technique of handling the dependencies of your software is called *Dependency
Injection*. The thing our code depends on to work, the `IsWebsiteOK` function,
is injected, in this case as an argument, into our code.

TDD draws you to perform Dependency Injection in order to make testing easier,
but the real benefits come when you are able to understand your code in
discrete, individual parts.

Finally, the technique we've used here of sending in a fake version of our
dependency in our tests is called "Mocking" or "Stubbing out" the dependency.
It's an excellent technique that allows us to control the behaviour of things in
our tests that we either don't own or want to test elsewhere.

## A New Requirement...

This is all great, but what happens when we try and check more websites. A _lot_
more websites. Let's check `http://google.co.uk` fifty times.

```go
package concurrency

import "testing"

func TestWebsiteChecker(t *testing.T) {
	websites := make([]string, 50)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.co.uk"
	}

	expectedResults := make([]bool, len(websites))
	for i := 0; i < len(websites); i++ {
		expectedResults[i] = true
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

Run this test and, after a bit of thumb-twiddling, we finally get:

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        10.320s
```

Ten seconds. So if we kick the number of checks up to 500...?

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        51.523s
```

Well at least it's consistent.

[^1]: For further reading on Test Doubles, Stubs, Mocks and the like, see https://martinfowler.com/articles/mocksArentStubs.html
