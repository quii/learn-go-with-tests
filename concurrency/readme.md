## Outline / notes

- Create a trivial, but paralizable, example - i.e. fetching multiple webpages
- Introduce goroutines to improve the performance
- Note that the test now fails as we can't be sure when the goroutines return
- Make tests pass using a sleep, but demonstrate that this doesn't scale
- introduce Channels for synchronization between processes

# Concurrency

Your colleague, Jo, has written a function that checks whether a webpage is
working or not. It's called `IsWebsiteOK`

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

### Write the test first

We'd like a function that takes a slice of strings and returns a slice of
bools. For each URL that returns a `true`, the returned slice will have
a `true` at that index - and vice versa for `false`s

We'll use a list of three URLs for now; the first two we know _should_ work; the
last one shouldn't work.

The first test checks that `WebsiteChecker` returns the same number of results
as websites; the second test checks that the results are what we expect for each
website.

```go
package concurrency

import "testing"

func TestWebsiteChecker(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		true,
		false,
	}

	actualResults := WebsiteChecker(websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
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

### Try and run the test

When we run the tests we see

```sh
# github.com/gypsydave5/learn-go-with-tests/concurrency/v1
./websiteChecker_test.go:18:10: undefined: websiteChecker
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1 [build failed]
```

### Write the minimal amount of code for the test to run and check the failing test output

To get this to pass we need to implement `websiteChecker` with the correct type
signature - a function that takes a single argument of a slice of `string`s
(`[]string`) and returns a slice of `bool`s (`[]bool`).

The simplest implementation of this is:
```go
func WebsiteChecker(_ []string) (result []bool) {
	return
}
```

Now when we run the tests we get
```sh
--- FAIL: TestWebsiteChecker (0.00s)
        websiteChecker_test.go:23: Wanted 3, got 0
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s
```

### Write enough code to make it pass

The code now builds (hooray), but the first test fails because the length of the
slice returned is too short. This is easy enough to fix by defining a length for
the slice:

```go
func WebsiteChecker(_ []string) (result [3]bool) {
	return
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

To make the second test pass we will iterate through the slice of URLs using
a `for...  range` loop. For each URL we will call `IsWebsiteOK` with the URL and
then store the answer in the `results` slice.  We grow the `results` slice by
using the `append` function, as seen in [the Arrays Tutorial][Arrays], so the
response from `IsWebsiteOK` will be at the same index as the `url` we're
checking.

When we've checked all of the URLs we'll finally return the `results` slice.

```go
func WebsiteChecker(urls []string) (results []bool) {
	for _, url := range urls {
		results = append(results, IsWebsiteOK(url))
	}

	return
}
```

Now when we run our tests:
```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.269s
```

#### Refactor

Could we refactor this? Well, the actual code isn't too bad at all, but we might
want to try and make our tests a little more readable by changing the test that
the two `bool` slices are equal. We could use `reflect.DeepEqual` as seen
previously, but instead this time we'll a small helper function to iterate
through both slices and flag up the difference (if any).

```go
package concurrency

import "testing"

func TestWebsiteChecker(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		true,
		false,
	}

	actualResults := WebsiteChecker(websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	if !sameResults(expectedResults, actualResults) {
		t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
	}
}

func sameResults(as, bs []bool) bool {
	for index, a := range as {
		if a != bs[index] {
			return false
		}
	}
	return true
}
```

## Dependency Injection

### Write the test first

So far, so good. But there are already two problems with what we've got so far.

1. If `google.co.uk` goes down, (or someone puts a website at `waat://furhurterwe.geds`), our expectations will be wrong for our tests.
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
can say that our dependency is on Jo's function `IsWebsiteOK`. If that function
stops working for any reason at all - whether the network cuts out or Jo creates
a bug in her code - our code will stop working and our tests will fail.

To mitigate this problem we can make `IsWebsiteOK` an extra argument to our
`WebsiteChecker` function. Then, in the tests, we can use a different function
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
		"waat://furhurterwe.geds",
	}

	expectedResults := []bool{
		true,
		false,
		true,
	}

	actualResults := WebsiteChecker(fakeIsWebsiteOK, websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	if !sameResults(expectedResults, actualResults) {
		t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
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

The way we want this to work is for `WebsiteChecker` to take our
`fakeIsWebsiteOK` function as it's first argument and to use it to 'check' the
websites. So that's what we've written in the test

```go
  actualResults := websiteChecker(fakeIsWebsiteOK, websites)
```

### Try and run the test

If we run this we get

```sh
./websiteChecker_test.go:25:33: too many arguments in call to websiteChecker
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

We should now be able to change `websiteChecker` to at least get the compilation
to work, just by adding an extra argument of type `func(string) bool` to its
argument list.

```go
package concurrency

func WebsiteChecker(_ func(string) bool, urls []string) (results []bool) {
	for _, url := range urls {
		results = append(results, IsWebsiteOK(url))
	}

	return
}
```

```sh
WebsiteChecker_test.go:36: Wanted false, got true
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v2     0.241s
```

Which is because we're still not using the function we're passing in.

### Write enough code to make it pass

This is easily fixed by naming and using the function:

```go
package concurrency

func WebsiteChecker(isOK func(string) bool, urls []string) (results []bool) {
	for _, url := range urls {
		results = append(results, isOK(url))
	}

	return
}
```

And now...

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2     0.013s
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

func WebsiteChecker(isOK testURL, urls []string) []bool {
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


##### Note on Dependency Injection and Test Doubles

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

## Concurrency

This is all great, but what happens when we try and check more websites. A _lot_
more websites. Let's check `http://google.co.uk` fifty times. To begin with,
let's use the real version of `IsWebsiteOK`; we'll want to change it pretty soon
but this test will give us a good idea of the pain point we're going to hit
(hint: it involves concurrency. Finally.)

```go
package concurrency

import "testing"

func TestWebsiteCheckerWithManyURLs(t *testing.T) {
	websites := make([]string, 50)
	for i := 0; i < len(websites); i++ {
		websites[i] = "http://google.co.uk"
	}

	expectedResults := make([]bool, len(websites))
	for i := 0; i < len(websites); i++ {
		expectedResults[i] = true
	}

	actualResults := websiteChecker(IsWebsiteOK, websites)

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

Well at least it's consistent... consistently slower, that is!

We're looking for a way of testing the _speed_ of our code now. Happily Go's
testing library supports benchmarking so that we can show that our code is
speeding up.

```go
func BenchmarkWebsiteCheckerWithManyURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		websites := make([]string, 100)
		for index, _ := range websites {
			websites[index] = "http://google.co.uk"
		}

		WebsiteChecker(IsWebsiteOK, websites)
	}
}
```

Here's your first benchmark. Benchmarks in Go are characterized by the `for...`
loop on the outside of the code that you want to benchmark. What it does in
effect is repeat the loop a number of times until it doesn't differ
significantly from the previous runs - until it is 'stable'.

To run it you need to add a flag to your `go test` command: `go test -benchmark=.`

```sh
goos: darwin
goarch: amd64
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v3
BenchmarkWebsiteCheckerWithManyURLs-4                  1        11352126640 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3        25.216s
```

The key number we want to read here is the one before `ns/op` - this is the
number of nanoseconds that it took, on average, to perform the operation in the
benchmark loop. 11352126640 nanoseconds is about 10 seconds, so the benchmark
confirms what our ad hoc testing has shown us.

Finally, let's stop annoying the good people at Google with hundreds of requests
everytime we run our tests. We can use another fake version of
`fakeIsWebsiteOK`, but this time we'll make it slow - say abut 20ms.

```go
func slowIsWebsiteOK(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}
```

The `Sleep()` function in from the `time` package is fairly self explanitory.
When we plug _that_ into the code in our benchmark, things get a lot faster.

```sh
goos: darwin
goarch: amd64
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v3
BenchmarkWebsiteCheckerWithManyURLs-4                  1        2267018950 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3        2.281s
```

Our goal now should be to make that 2 seconds duration much closer to 2 miliseconds.

### Write enough code to make it pass

For the purposes of this test, 'passing' should be thought of as being
synonymous with 'making it a lot faster'.

Which means that we finally get to do something with concurrency in Go!

```go
func WebsiteChecker(isOK URLchecker, urls []string) (results []bool) {
	for _, url := range urls {
		go func(url string) {
			results = append(results, isOK(url))
		}()
	}

	return
}
```

Concurrency in Go is built up from the snappily-named 'goroutines'. In any place
where you can call a function, you can place the keyword `go` in front of it and
the function will execute as a separate process to the parent process.

Here we are executing an anonymous function as a goroutine inside the `for` loop
we had before. The body of the function is just the same as the loop body was
before. The only difference here is that each iteration of the loop will spin
off a new process, concurrent with the current process (the `WebsiteChecker`
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

None of the goroutines that span off in our for loop had enough time to append
their result to the `results` slice; `WebsiteChecker` is too fast for them, and
it returns the still empty slice.

To fix this we can just wait while all the goroutines do their work, and then
return. Two seconds ought to do it

```go
func WebsiteChecker(isOK URLchecker, urls []string) (results []bool) {
	for _, url := range urls {
		go func(url string) {
			results = append(results, isOK(url))
		}(url)
	}

	time.Sleep(2 * time.Second)

	return
}
```

Now when we run the tests

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v3        2.022s
```

But if we run them again

```sh
--- FAIL: TestWebsiteChecker (2.00s)
        websiteChecker_test.go:26: Wanted 3, got 1
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v3        2.015s
```

and again

```sh
--- FAIL: TestWebsiteChecker (2.00s)
        websiteChecker_test.go:30: Wanted [true false true], got [true true false]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v3        2.020s
```

Your tests may give slightly different results, but we should still expect to
see one of the above three outputs. So now what's going wrong?

There are two issues. First, we're not waiting long enough - this is why we get
a set of results that isn't the right length. We could fix this by just bumping
up the time slightly - or just waiting for the results slice to be the right
length - if it wasn't for the other problem.

The other problem is that the goroutines are able to append to the results slice in
a different order to that which they were called in, which has the effect of the
results coming back in a different order.

This is not a problem that can be solved by sleeping for a few extra seconds; we
will need a completely approach to handling concurrency that allows coordination
between different processes.

#### Channels

[^1]: For further reading on Test Doubles, Stubs, Mocks and the like, see https://martinfowler.com/articles/mocksArentStubs.html


[Arrays]: ../arrays/

## An observation

Did you notice that the time it took for the `websiteCheckerTest` to run
increased dramatically when we were really checking websites? It added around
a quarter of a second to the total time. Although the Internet is fast, and the
response we're getting from the websites is coming back pretty quickly, it still
takes time for our functions to make those real requests.
