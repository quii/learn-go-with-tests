## Outline / notes

- Create a trivial, but paralizable, example - i.e. fetching multiple webpages
- Introduce goroutines to improve the performance
- Note that the test now fails as we can't be sure when the goroutines return
- Make tests pass using a sleep, but demonstrate that this doesn't scale
- introduce Channels for synchronization between processes

---

change of approach!!!

- start with a result object coming back from the original function
- avoids this mess with the ordering of the results
- is this copping out?

- use map!!!

---

# Concurrency

Here is are a pair of functions that somebody else has written:

```go
package concurrency

import "net/http"

func CheckWebsite(url string) bool {
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

```go
package concurrency

func CheckWebsites(urls []string) []bool {
	results := make([]bool, len(urls))

	for index, url := range urls {
		results[index] = isOK(url)
	}

	return results
}
```

`CheckWebsite` takes a URL as an argument and returns `true` if that URL returns
a 200 status code to an HTTP HEAD request, `false` for any other response or
error.

`CheckWebsites` takes a slice of URLs and returns a `map[string]bool`, mapping
together the original URLs to the result of them being checked by
`CheckWebsite`.

Here is the test for `CheckWebsites`:

```go
package concurrency

import (
	"reflect"
	"testing"
)

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	actualResults := CheckWebsites(websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	expectedResults := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	if !reflect.DeepEqual(expectedResults, actualResults) {
		t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
	}
}
```

This test sends a slice of URLs to `CheckWebsites`, and asserts on the length of
the map returned, as well as the contents of that map (using `reflect.DeepEqual`).

`CheckWebsites` is currently being used with very large slices of URLs, and
people are complaining about the speed. We've been asked to make it faster.

If the above isn't familiar to you, don't worry about it. The key thing for this
exercise is that the above function _may_ take some time to return a result.

### Write the test first

```go
func BenchmarkCheckWebsites(b *testing.B) {
	websites := make([]string, 100)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.com"
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(websites)
	}
}
```

A simple benchmark that executes `CheckWebsites` using a slice of 100 URLs. When
we run this with `go test -bench=.`:

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v1
BenchmarkCheckWebsites-4               1        9940439437 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        10.723s
```

9940439437 nanoseconds, around 10 seconds.

This test is slow to run, which is a bad thing. If a test is slow we will
be less willing to run it, and if a test isn't being run then that test is
useless.

Secondly, the test is making real network requests to google.com. This will make
the benchmark inconsistent - what if google.com is down, or we have no Internet
connection? The speed that the benchmark reports will fluctuate, and we will
find it harder to be sure that any changes we make have sped up or slowed down
`CheckWebsites`.

```go
func fakeWebsiteChecker(url string) bool {
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	websites := make([]string, 100)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.com"
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(fakeWebsiteChecker, websites)
	}
}
```

We have written a function called `stubWebsiteChecker` which performs in the
same way as the real `CheckWebsite` - it takes a `string` and returns a `bool` -
only without making any network calls. We are using `stubWebsiteChecker` as the
first argument to `CheckWebsites`.

### Try and run the test

When we run `go build`:

```sh
# github.com/gypsydave5/learn-go-with-tests/concurrency/v2
./CheckWebsites_test.go:22:32: too many arguments in call to CheckWebsites
        have (func(string) bool, []string)
        want ([]string)
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v2 [build failed]
```

This is telling us that we need to update `CheckWebsites` to take an extra
argument of the type `func(string) bool`.

Just by reading the output of the compiler carefully, we now know exactly what
we have to do next. When performing TDD it is _vital_ that you read the output of
the compiler and your tests carefully. These outputs, especially those from the
compiler, will be telling you exactly what we have to do next. While there
should only be one way in which your tests can pass, there are many ways in
which they can fail. If you can read and understand why they are failing, you
will be 90% of the way to making them pass.


## Dependency Injection

### Write the test first

So far, so good. But there are already two problems with what we've got so far.

1. If `google.co.uk` goes down, (or someone puts a website at `waat://furhurterwe.geds`), our expectations will be wrong for our tests.
2. If we turn off the Internet, our tests will always fail.

And in true TDD style, we can demonstrate this with a failing test. So, turn off
the computer's WiFi / unplug the network cable and run the tests again:

```sh
--- FAIL: TestWebsiteChecker (0.00s)
        websiteChecker_test.go:39: expected value of key 'http://google.com' in actual results to be 'true', but it was 'false'
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v2        0.018s
```

This dependency on the Internet is a bad thing because these failures will have
nothing to do with any changes to the behaviour of our code. More precisely, we
can say that our dependency is on Jo's function `IsWebsiteOK`. If that function
stops working for any reason at all - whether the network cuts out or the code
has a bug - our tests will fail.

To mitigate this problem we can make the function `IsWebsiteOK` an
argument to our `WebsiteChecker` function. Then, in the tests, we can use
a different function with the same interface as that behaves in a way that we
can control.

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
		"waat://furhurterwe.geds",
	}

	actualResults := WebsiteChecker(fakeIsWebsiteOK, websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	expectedResults := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": false,
		"waat://furhurterwe.geds":    true,
	}

	assertSameResults(t, expectedResults, actualResults)
}
```

We've added a new function, `fakeIsWebsiteOK`, which has the same behaviour as
`IsWebsiteOK`. From the outside you couldn't tell the difference between them -
they take a `string` and return a `bool`. But on the inside `fakeIsWebsiteOK`
is just an `if` statement that always returns `true` unless the `string`
argument is `"http://blog.gypsydave5.com"`. It's a function we have complete
control over - because we wrote it.

The expectations have also been updated; we now expect the middle one to fail.

The way we want this to work is for `WebsiteChecker` to take our
`fakeIsWebsiteOK` function as it's first argument and to use it to 'check' the
websites. So that's what we've written in the test

```go
  actualResults := websiteChecker(fakeIsWebsiteOK, websites)
```

### Try and run the test

```sh
./websiteChecker_test.go:21:33: too many arguments in call to WebsiteChecker
        have (func(string) bool, []string)
        want ([]string)
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v2 [build failed]
```

A faiure to compile. If we read the output of the compiler we can see that
`websiteChecker` wants `([]string)` as it's arguments, but we gave it
`(func(string) bool, []string)`. So we've learnt that `func(string) bool` is the
type of our `fakeWebsiteOK` function in the same way as `[]string` is the type
of the slice of strings we're passing in as the second argument.

Just by reading the output of the compiler carefully, we now know exactly what
we have to do next. When performing TDD it is _vital_ that you read the output of
the compiler and your tests carefully. These outputs, especially those from the
compiler, will be telling you exactly what we have to do next. While there
should only be one way in which your tests can pass, there are many ways in
which they can fail. If you can read and understand why they are failing, you
will be 90% of the way to making them pass.

### Write the minimal amount of code for the test to run and check the failing test output

```go
package concurrency

func WebsiteChecker(_ func(string) bool, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = IsWebsiteOK(url)
	}

	return results
}
```

```sh
# github.com/gypsydave5/learn-go-with-tests/concurrency/v2
./CheckWebsites_test.go:21:32: not enough arguments in call to CheckWebsites
        have ([]string)
        want (func(string) bool, []string)
./CheckWebsites_test.go:35:6: undefined: reflect
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v2 [build failed]
```

Just need to update the original test...

```go
package concurrency

import (
	"reflect"
	"testing"
)

func stubWebsiteChecker(url string) bool {
	if url == "waat://furhurterwe.geds" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	actualResults := CheckWebsites(stubWebsiteChecker, websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	expectedResults := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	if !reflect.DeepEqual(expectedResults, actualResults) {
		t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
	}
}
```

And now it runs the test OK

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v2
BenchmarkCheckWebsites-8               1        8772112216 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        8.934s
```

But it's still really slow.

### Write enough code to make it pass

It won't speed up until we actually use the function we're passing in to check
the websites:

```go
package concurrency

func CheckWebsites(websiteChecker func(string) bool, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = websiteChecker(url)
	}

	return results
}
```

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v2
BenchmarkCheckWebsites-8         1000000              1942 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        1.975s
```

#### Refactor

`func(string) bool` doesn't exactly trip off the tongue when trying to describe
what the function is doing - you can tell the behaviour, but it's hard to say
what the intention of it is.

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(websiteChecker WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = websiteChecker(url)
	}

	return results
}
```

We've used the `type` keyword here to say that we'd like `func(string) bool` to
also be known as `TestURL`. This is a useful technique to help your code read
nicely, especially with long function types.

This technique of handling the dependencies of your software is called *Dependency
Injection*. The thing our code depends on to work, the `IsWebsiteOK` function,
is injected, in this case as an argument, into our code.

TDD will inspire you to perform Dependency Injection in order to make testing
easier, but the real benefits come when you are able to understand your code in
discrete, individual parts.

Finally, the technique we've used here of sending in a fake version of our
dependency in our tests is called "Mocking" or "Stubbing out" the dependency.
It's an excellent technique that allows us to control the behaviour of things in
our tests that we either don't own or want to test elsewhere.

>>>CARRY ON HERE<<<

### Concurrency

This is all great, but what happens when we try and check more websites. A _lot_
more websites. Let's check `http://google.co.uk` one hundred times.

```go
package concurrency

import "testing"

func TestWebsiteCheckerWithManyURLs(t *testing.T) {
	websites := make([]string, 100)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.co.uk"
	}

	expectedResults := make(map[string]bool)

	for i := 0; i < len(websites); i++ {
		expectedResults["http://google.co.uk"] = true
	}

	actualResults := WebsiteChecker(IsWebsiteOK, websites)

	want := len(expectedResults)
	got := len(actualResults)
	if len(actualResults) != len(websites) {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	assertSameResults(t, expectedResults, actualResults)
}
```
We've written a new test that uses the real version of `IsWebsiteOK` for now.

Run this test and, after a bit of thumb-twiddling, we finally get:

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3        11.122s
```

Ten seconds or so. So if we kick the number of checks up to 500...?

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        51.523s
```

A proportional increase in latency (i.e. it was five times as slow)

We're looking for a way of testing the _speed_ of our code now. We can do this
by using a benchmark again as we saw in [the `for` tutorial][For].

```go
func BenchmarkWebsiteChecker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		websites := make([]string, 100)
		for index, _ := range websites {
			websites[index] = "http://google.co.uk"
		}

		WebsiteChecker(IsWebsiteOK, websites)
	}
}
```

When we run `go test -benchmark=.`

```sh
goos: darwin
goarch: amd64
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v3
BenchmarkWebsiteChecker-4              1        11100348034 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3        21.936s
```

The key number we want to read here is the one before `ns/op` - this is the
number of nanoseconds that it took, on average, to perform the operation in the
benchmark loop. 11100348034 nanoseconds is about 10 seconds, so the benchmark
confirms what our ad hoc testing has shown us.

Finally, let's stop annoying Google with hundreds of requests everytime we run
our tests. We can use another fake version of `fakeIsWebsiteOK`, but this time
we'll make it slow - say abut 20ms.

```go
func slowIsWebsiteOK(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}
```

The `Sleep()` function in from the `time` package is fairly self explanitory.

```go
func BenchmarkWebsiteChecker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		websites := make([]string, 100)
		for index, _ := range websites {
			websites[index] = "http://google.co.uk"
		}

		WebsiteChecker(slowIsWebsiteOK, websites)
	}
}
```

```sh
goos: darwin
goarch: amd64
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v3
BenchmarkWebsiteCheckerWithManyURLs-4                  1        2267018950 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3        2.281s
```

Our goal now should be to make that 2 seconds duration much closer to 2 milliseconds.

### Write enough code to make it pass

```go
func WebsiteChecker(isOK URLchecker, urls []string) (results []bool) {
	results := make(map[string]bool)

	for _, url := range urls {
		go func(u string) {
			results[u] = isOK(u)
		}(url)
	}

	return results
}
```

Concurrency in Go is built up from the _goroutines_. In any place where you can
call a function, you can place the keyword `go` in front of it and the function
will execute as a separate process to the parent process.

Here we are executing an anonymous function as a goroutine inside the `for` loop
we had before. The body of the function is just the same as the loop body was
before. The only difference is that each iteration of the loop will start
a new process, in parallel to with the current process (the `WebsiteChecker`
function) each of which will append its result to the `results` slice.

But when we give this a go:

```sh
--- FAIL: TestWebsiteChecker (0.00s)
        websiteChecker_test.go:26: Wanted 3, got 0
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v3        0.015s
```

We are caught by the first test we wrote; `WebsiteChecker` is now returning an
empty slice. What went wrong?

None of the goroutines that our `for` loop started had enough time to add
their result to the `results` map; the `WebsiteChecker` function is too fast for
them, and it returns the still empty map.

To fix this we can just wait while all the goroutines do their work, and then
return. Two seconds ought to do it

```go
package concurrency

import "time"

type TestURL func(string) bool

func WebsiteChecker(isOK TestURL, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func(u string) {
			results[u] = isOK(u)
		}(url)
	}

	time.Sleep(2 * time.Second)

	return results
}
```

Now when we run the tests you might get

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3        2.022s
```

But if you're unlucky (this is more likely if you run them with `go test -bench=.`)

```sh
fatal error: concurrent map writes

goroutine 8 [running]:
runtime.throw(0x12c5895, 0x15)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/panic.go:605 +0x95 fp=0xc420037700 sp=0xc4200376e0 pc=0x102d395
runtime.mapassign_faststr(0x1271d80, 0xc42007acf0, 0x12c6634, 0x17, 0x0)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:783 +0x4f5 fp=0xc420037780 sp=0xc420037700 pc=0x100eb65
github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1(0xc42007acf0, 0x12d3938, 0x12c6634, 0x17)
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x71 fp=0xc4200377c0 sp=0xc420037780 pc=0x12308f1
runtime.goexit()
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/asm_amd64.s:2337 +0x1 fp=0xc4200377c8 sp=0xc4200377c0 pc=0x105cf01
created by github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xa1

        ... many more scary lines of text ...
```

Errors this long tend to freak me out, but the headline is what we should be
paying attention to. `fatal error: concurrent map writes`. Maps in Go don't like
it when more than one thing tries to write to them at once.

What we have here is a classic _race condition_, an bug that occurs when the
output of our software is dependent on the timing and sequence of events that we
have no control over. Because we cannot control exactly when each goroutine
writes to the results, we are vulnerable to two goroutines writing to it at the
same time.

Go can help us to spot race conditions with its built in [_race detetector_][godoc_race_detector].
To enable this feature, run the tests with the `race` flag: `go test -race`.

You should get some pretty verbose output that looks a bit like this:

```sh
==================
WARNING: DATA RACE
Write at 0x00c420084d20 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Previous write at 0x00c420084d20 by goroutine 7:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Goroutine 8 (running) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c

Goroutine 7 (finished) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c
==================
```

Again, the details are hard to read, but the headline isn't: `WARNING: DATA
RACE` is pretty unambiguous. Reading into the body of the error we can see two
different goroutines performing writes on a map - pretty much as we suspected.

### Channels

We can solve this problem by coordinating our goroutines using _channels_.
Channels are a Go data structure that can both receive and send values. These
operations, along with their details, allow communication between processes.

```go
package concurrency

type TestURL func(string) bool
type result struct {
	string
	bool
}

func WebsiteChecker(isOK TestURL, urls []string) map[string]bool {
	results := make(map[string]bool)
	urlChannel := make(chan string)
	resultChannel := make(chan result)

	go func() {
		for {
			url := <-urlChannel

			good := isOK(url)
			result := result{url, good}
			resultChannel <- result
		}
	}()

	for _, url := range urls {
		urlChannel <- url
	}

	for i := 0; i < len(urls); i++ {
		result := <-resultChannel
		results[result.string] = result.bool
	}

	return results
}
```

This version of `WebsiteChecker` uses two channels to organise the work of one
anonymous goroutine which is where the bulk of the work gets done.

The first change to note is that we've introduced a new type, `result`, which is
a struct of a `string` and a `bool`. We will use this to keep the information
about a URL and its test result together.

First, we `make` two channels, one of which communicates using strings (`chan
string`), and the other which communicates using results (`chan result`). We
will use the first channel to manage the URLs and the second to manage the
results of using the `TestURL` function.

Next comes the main goroutine. It's made of an infinite loop - a `for` that has
no terminating condition. First it takes a value from the `urlChannel` with the
receive operation: `url := <-urlChannel`. We then test the url using our `isOK`
function, package the response up into a `result` along with the original URL,
and then uses the send operation to put the result onto the `resultChannel`:
`resultChannel <- result`.

Because this loop is in a goroutine it won't block the main process of
`WebsiteChecker`, but will keep on receiving and sending values to the two
channels as long as there are values to receive and send.

[^1]: For further reading on Test Doubles, Stubs, Mocks and the like, see https://martinfowler.com/articles/mocksArentStubs.html

[Arrays]: ../arrays/
[For]: ../for/
[godoc_maps]: https://blog.golang.org/go-maps-in-action
[godoc_zero_values]: https://golang.org/ref/spec#The_zero_value
[godoc_race_detector]: https://blog.golang.org/race-detector

## An observation

Did you notice that the time it took for the `websiteCheckerTest` to run
increased dramatically when we were really checking websites? It added around
a quarter of a second to the total time. Although the Internet is fast, and the
response we're getting from the websites is coming back pretty quickly, it still
takes time for our functions to make those real requests.
