# Context (WIP)

The [Go blog describes the motivation for working with `context` excellently](https://blog.golang.org/context)

> In Go servers, each incoming request is handled in its own goroutine. Request handlers often start additional goroutines to access backends such as databases and RPC services. The set of goroutines working on a request typically needs access to request-specific values such as the identity of the end user, authorization tokens, and the request's deadline. When a request is canceled or times out, all the goroutines working on that request should exit quickly so the system can reclaim any resources they are using.

In this chapter we'll cover some usage with some simple examples of how to manage long running processes. 

We're going to start with a classic example of a web server that when hit kicks off a potentially long-running process to fetch some data for it to return in the response. 

We will exercise a scenario where a user cancels the request before the data can be retrieved and we'll make sure the process is told to give up. 

I've set up some code on the happy path to get us started. Here is our server code

```go
func NewHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, store.Fetch())
	}
}
```

The function `NewHandler` takes a `Store` and returns us a `http.HandlerFunc`. Store is defined as:

```go
type Store interface {
	Fetch() string
}
```

The returned function calls the `store`'s `Fetch` method to get the data and writes it to the response.

We have a corresponding stub for `Store` which we use in a test.

```go
type StubStore struct {
	response string
}

func (s *StubStore) Fetch() string {
	return s.response
}

func TestHandler(t *testing.T) {
	data := "hello, world"
	svr := NewHandler(&StubStore{data})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	svr.ServeHTTP(response, request)

	if response.Body.String() != data {
		t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
	}
}
```

Now that we have a happy path, we want to make a more realistic scenario where the `Store` cant finish a`Fetch` before the user cancels the request.

## Write the test first

Our handler will need a way of telling the `Store` to cancel the work so update the interface.

```go
type Store interface {
	Fetch() string
	Cancel()
}
```

We will need to adjust our spy so it takes some time to return `data` and a way of knowing it has been told to cancel. We'll also rename it to `SpyStore` as we are now observing the way it is called. It'll have to add `Cancel` as a method to implement the `Store` interface.

```go
type SpyStore struct {
	response string
	cancelled bool
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}
```

Let's add a new test where we cancel the request before 100 milliseconds and check the store to see if it gets cancelled.

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
      store := &SpyStore{response: data}
      svr := Server(store)
  
      request := httptest.NewRequest(http.MethodGet, "/", nil)
      
      cancellingCtx, cancel := context.WithCancel(request.Context())
      time.AfterFunc(5 * time.Millisecond, cancel)
      request = request.WithContext(cancellingCtx)
      
      response := httptest.NewRecorder()
  
      svr.ServeHTTP(response, request)
  
      if !store.cancelled {
          t.Errorf("store was not told to cancel")
      }
  })
```

From the google blog again

> The context package provides functions to derive new Context values from existing ones. These values form a tree: when a Context is canceled, all Contexts derived from it are also canceled.

It's important that you derive your contexts so that cancellations are propegated throughout the call stack for a given request. 

What we do is derive a new `cancellingCtx` from our `request` which gives us access to a `cancel` function. We then schedule that function to be called in 5 milliseconds by using `time.AfterFunc`. Finally we use this new context in our request by calling `request.WithContext`

## Try to run the test

The test fails as we'd expect

```go
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.00s)
    	context_test.go:62: store was not told to cancel
```

## Write enough code to make it pass

Remember to be disciplined with TDD. Write the _minimal_ amount of code to make our test pass. 

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.Cancel()
		fmt.Fprint(w, store.Fetch())
	}
}
```

This makes this test pass but it doesn't feel good does it! We surely shouldn't be cancelling `Store` before we fetch on _every request_. 

By being disciplined it highlighted a flaw in our tests, this is a good thing! 

We'll need to update our happy path test to assert that it does not get cancelled. 

```go
t.Run("returns data from store", func(t *testing.T) {
    store := SpyStore{response: data}
    svr := Server(&store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }
    
    if store.cancelled {
        t.Error("it should not have cancelled the store")
    }
})
```

Run both tests and the happy path test should now be failing and now we're forced to do a more sensible implementation.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}
	}
}
```

What have we done here?

`context` has a method `Done()` which returns a channel which gets sent a signal when the context is "done" or "cancelled". We want to listen to that signal and call `store.Cancel` if we get it but we want to ignore it if our `Store` manages to `Fetch` before it.

To manage this we run `Fetch` in a goroutine and it will write the result into a new channel `data`. We then use `select` to effectively race to the two asynchronous processes and then we either write a response or `Cancel`

## Refactor

We can refactor our test code a bit by making assertion methods on our spy

```go
func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Errorf("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Errorf("store was told to cancel")
	}
}
```

Remember to pass in the `*testing.T` when creating the spy. 

```go
func TestServer(t *testing.T) {
	data := "hello, world"

	t.Run("returns data from store", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}

		store.assertWasNotCancelled()
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		request = request.WithContext(cancellingCtx)

		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		store.assertWasCancelled()
	})
}
```

This approach is ok, but is it idiomatic? 

Does it make sense for our web server to be concerned with manually cancelling `Store`? What if `Store` also happens to depend on other slow-running processes? We'll have to make sure that `Store.Cancel` correctly propagates the cancellation to all of its dependants. 

One of the main points of `context` is that it is a consistent way of offering cancellation. 

From the Google blog again

> At Google, we require that Go programmers pass a Context parameter as the first argument to every function on the call path between incoming and outgoing requests. This allows Go code developed by many different teams to interoperate well. It provides simple control over timeouts and cancelation and ensures that critical values like security credentials transit Go programs properly.

Maybe it would be better for us to follow that approach and instead pass through the `context` to our `Store` and let it be responsible.


## notes for later...

- Cover context.value, but warn against putting stupid stuff like databases in it, keep it request scoped. Obscures the inputs and outputs of functions and is not typesafe. 
