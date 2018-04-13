# Select (WIP)

You have been asked to make a function called `WebsiteRacer` which takes two URLs and "races" them by hitting them with a HTTP GET and returning the URL which returned first. If none of them return within 10 seconds then it should return an `error`

For this we will be using

- `net/http` to make the HTTP calls.
- `net/http/httptest` to help us test them.
- _Mocking_ to let us control our tests, keep them fast and test edge cases.
- `go` routines.
- `select` to synchronise processes. 

## Write the test first

Let's start with something naive to get us going.

```go
func TestRacer(t *testing.T) {
	slowURL := "http://www.facebook.com"
	fastURL := "http://www.quii.co.uk"

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want{
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
```

We know this isn't perfect and has problems but it will get us going. It's important not to get too hung-up on getting things perfect first time.

## Try to run the test

`./racer_test.go:14:9: undefined: Racer`


## Write the minimal amount of code for the test to run and check the failing test output

```go
func Racer(a, b string) (winner string) {
	return
}
```

`racer_test.go:25: got '', want 'http://127.0.0.1:54228'`

## Write enough code to make it pass

```go
func Racer(a, b string) (winner string) {
	startA := time.Now()
	http.Get(a)
	aDuration := time.Since(startA)

	startB := time.Now()
	http.Get(b)
	bDuration := time.Since(startB)

	if aDuration < bDuration {
		return a
	}

	return b
}
```

For each url: 

1. We use `time.Now()` to record just before we try and get the `URL`
2. Then we use `http.Get` to try and get the contents of the `URL`. This function returns a response and an `error` but so far we are not interested in these values
3. `time.Since` takes the start time and returns a `time.Duration` of the difference.

Once we have done this we simply compare the durations to see which is the quickest.

### Problems

This may or may not make the test pass for you. The problem is we're reaching out to real websites to test our own logic.

Testing code that uses HTTP is so common that Go has tools in the standard library to help you test it.

In the mocking and dependency injection chapters we covered how ideally we dont want to be relying on external services to test our code because they can be 

- Slow
- Flaky
- Cant test edge cases

In go there is a package `net/http/httptest` where you can easily create a mock http server that you can fully control.

Let's change our tests to be a little more reliable. 

```go
func TestRacer(t *testing.T) {

	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}

	slowServer.Close()
	fastServer.Close()
}
```

The syntax may look a bit busy but just take your time.

`httptest.NewServer` takes a `http.HandlerFunc` which we are sending in via an _anonymous function_. 

`http.HandlerFunc` is a type that looks like this `type HandlerFunc func(ResponseWriter, *Request)`

All it's really saying is it needs a function takes a `ResponseWriter` and a `Request`, which is not too surprising for a HTTP server

It turns out there's really no extra magic here, **this is also how you would write a _real_ HTTP server in Go**. The only difference is we are wrapping it in a `httptest.NewServer` which makes it easier to use with testing, as it finds an open port to listen on and then you can close it when you're done with your test.

Inside our two servers we make the slow one have a short `time.Sleep` when we get a request to make it slower than the other one. Both servers then write an `OK` response back to the caller.

If you re-run the test it will definitely pass now and should be faster. Try playing with sleeps to deliberately break the test.

## Refactor

We have some duplication in both our production code and test code.

```go
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
```

This DRY-ing up makes our `Racer` code a lot easier to read.

```go
func TestRacer(t *testing.T) {

	slowServer := makeDelayedServer(20 * time.Millisecond)
	fastServer := makeDelayedServer(0 * time.Millisecond)

	defer slowServer.Close()
	defer fastServer.Close()

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
```

We've refactored our server making into `makeDelayedServer` just to move some uninteresting code out of the test and reduce repetition.

There's a keyword that is maybe unfamiliar to you called `defer`. What this means it will run the function _at the end of the containing function_. Before we had the two `Close` calls at the end of our test. `defer` is useful when you want to keep the context of these important cleanup operations closer to where it's relevant. We're telling the reader (and the compiler) to remember to close our servers once the function is finished.

This is an improvement but there's more we can do. Why are we testing the speeds of the websites one after another when Go is great at concurrency. We should be able to test them at the same time and whatever comes back first wins. 

To do this, we're going to introduce a new construct called `select` which helps us synchronise processes really easily and clearly.  

```go
func Racer(a, b string) (winner string) {
	select {
	case <-measureResponseTime(a):
		return a
	case <-measureResponseTime(b):
		return b
	}
}

func measureResponseTime(url string) chan interface{} {
	ch := make(chan interface{})
	go func() {
		fmt.Println("getting", url)
		http.Get(url)
		ch <- true
	}()
	return ch
}
```

If you recall from the concurrency chapter, you can wait for values to be sent to a channel with `myVar := <-ch`. This is a _blocking_ call, as you're waiting for a value. 

What `select` lets you do is wait on _multiple_ channels. The first one to send a value "wins" and the code underneath the `case` is executed. 

In our case we have defined a function `measureResponseTime` which creates a `chan interface` and returns it. `inteface` is a type in Go which means "i dont know what the type is". In our case, we don't really _care_ what the type is returned, we just want to send a signal back in the channel to say we're finished. 

Inside the same function we start a go routine which will send a signal into that channel once we have completed `http.Get(url)`

We use this function in our `select` to set up two channels for each of our `URL`s. Whichever one writes to its channel first will have its code executed in the `select`, which results in its `URL` being returned (and being the winner). 

### Timeouts

Our final requirement was to return an error if the test took longer than 10 seconds. 

## Write the test first

```go
t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
    serverA := makeDelayedServer(11 * time.Second)
    serverB := makeDelayedServer(12 * time.Second)

    defer serverA.Close()
    defer serverB.Close()
    
    _, err := Racer(serverA.URL, serverB.URL)

    if err == nil {
        t.Error("expected an error but didn't get one")
    }
})
```

We've made our test servers take longer than 10s to return to exercise this scenario and we are expecting `Racer` to return two values now, the winning URL (which we ignore in this test with `_`) and an `error`.

## Try to run the test

`./racer_test.go:37:10: assignment mismatch: 2 variables but 1 values`

## Write the minimal amount of code for the test to run and check the failing test output

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-measureResponseTime(a):
		return a, nil
	case <-measureResponseTime(b):
		return b, nil
	}
}
```

Change the signature of `Racer` to return the winner and an `error`. Return `nil` for our happy cases.

The compiler will complain about your _first test_ only looking for one value so change that line to

`got, _ := Racer(slowURL, fastURL)`, knowing that we should check we _don't_ get an error in our happy scenario. 

If you run it now after 11 seconds it will fail

```
--- FAIL: TestRacer (12.00s)
    --- FAIL: TestRacer/returns_an_error_if_a_server_doesn't_respond_within_10s (12.00s)
    	racer_test.go:40: expected an error but didn't get one
```

## Write enough code to make it pass

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-measureResponseTime(a):
		return a, nil
	case <-measureResponseTime(b):
		return b, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

`time.After` is a very handy function when using `select`. Although it didn't happen in our case you can potentially write code that blocks forever if the channels you're listening on never return a value. `time.After` returns a `chan` like `measureResponseTime` and will send a signal down it after the amount of time you define. 

For us this is perfect, if `a` or `b` manage to return they win, but if we get to 10 seconds then our `time.After` will send a signal and we'll return an `error`

### Slow tests

The problem we have is that this test takes 10 seconds to run. For such a simple bit of logic this doesn't feel great. 

What we can do is make the timeout configurable so in our test we can have a very short timeout and then when the code is used in the real world it can be set to 10 seconds. 

```go
func Racer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-measureResponseTime(a):
		return a, nil
	case <-measureResponseTime(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

Our tests now wont compile because we're not supplying a timeout

Before rushing in to add this default value to both our tests let's just listen to them.

- Do we care about the timeout in the "happy" test
- The requirements were explicit about the timeout

Given this knowledge, let's do a little refactoring to be sympathetic to both our tests and the users of our code

```go
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-measureResponseTime(a):
		return a, nil
	case <-measureResponseTime(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

Our users and our first test can just use `Racer` (which uses `ConfigurableRacer` under the hood) and our sad path test can use `ConfigurableRacer`.

```go
func TestRacer(t *testing.T) {

	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, _ := Racer(slowURL, fastURL)

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}
```

## Wrapping up

### `select`

### `httptest`