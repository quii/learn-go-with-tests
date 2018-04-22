# HTTP Server (WIP)

You have been asked to create a web server where users can track how many games players have won

- `GET /players/{name}` should return a number indicating total number of wins
- `POST /players/{name}/win` should increment the number of wins 

We will follow the TDD approach, getting to "working software" as quickly as we can and then incrementing until we have the solution

### Chicken and egg

How can we incrementally build this? We cant `GET` a player without having stored something and it seems hard to know if `POST` has worked without the `GET` endpoint already existing. 

This is where _mocking_ shines. 

- `GET` will need a `PlayerStore` to reach out to. This should be an interface so when we test we can create a simple stub to work with.
- For `POST` we can _spy_ on its calls to `PlayerStore` to make sure it stores players correctly. 
- For having some working software quickly we can make a very simple in-memory implementation and then later we can create an implementation backed by whatever store we want. 

## Write the test first

Before worrying about any kind of domain-level logic, we should get the overall application scaffolding sorted. We'll start by testing that if we hit an endpoint we get back "Hello, world" and then wire it up into a real application.

To create a web server in Go you will typically call [https://golang.org/pkg/net/http/#ListenAndServe](ListenAndServe).

```go
func ListenAndServe(addr string, handler Handler) error
```

The [`Handler`](https://golang.org/pkg/net/http/#Handler) part is the bit we'll be working with.

```go
type Handler interface {
        ServeHTTP(ResponseWriter, *Request)
}
```

It's an interface which expects two arguments, the first being where we _write our response_ and the second being the HTTP request that was sent to us.

```go
func TestGETPlayers(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	PlayerServer(res, req)

	t.Run("hello world in response body", func(t *testing.T) {
		got := res.Body.String()
		want := "Hello, world"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})
}
```

In order to test our server, we will need a `Request` to send in and we'll want to _spy_ on what our handler writes to the `ResponseWriter`. 

- We make a _real_ `Request` as all it really is a collection of data sent into our function. The `nil` argument refers to the request's body, which we don't need to set in this case.
- `net/http/httptest` has a spy already made for us called `ResponseRecorder` so we can use that. 

## Try to run the test

`./server_test.go:13:2: undefined: PlayerServer`

## Write the minimal amount of code for the test to run and check the failing test output

The compiler is here to help, just listen to it.

Define `PlayerServer`

```go
func PlayerServer() {}
```

Try again

```
./server_test.go:13:14: too many arguments in call to PlayerServer
	have (*httptest.ResponseRecorder, *http.Request)
	want ()
```

Add the arguments to our function

```go
import "net/http"

func PlayerServer(w http.ResponseWriter, r *http.Request) {

}
```

The code now compiles and the test fails

```
=== RUN   TestGETPlayers/hello_world_in_response_body
    --- FAIL: TestGETPlayers/hello_world_in_response_body (0.00s)
    	server_test.go:20: got '', want 'Hello, world'
```

## Write enough code to make it pass

From the DI chapter we touched on HTTP servers with a `Greet` function. We learned that net/http `ResponseWriter` also implements io `Writer` so we can use `fmt.Fprint` to send strings as HTTP responses

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}
```

The test should now pass

## Complete the scaffolding

We want to wire this up into an application. Why?

- Have _actual working software_, we don't want to write tests for the sake of it, it's good to see the code in action.
- As we refactor our code, it's likely we will change the structure of the program. We want to make sure this is reflected in our application too.

Create a new file for our application and put this code in.

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":5000", http.HandlerFunc(PlayerServer)); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

So far all of our "mains" have been in one file, however this isn't best practice for larger projects where you'll want to separate things into different files. 

To run this, do `go build` which will take all the `.go` files in the directory and build you a program. You can then execute it with `./myprogram`

### `http.HandlerFunc`

Earlier we explored that the `Handler` interface is what we need to implement in order to make a server. _Typically_ we do that by creating a `struct` and make it implement the interface. However the use-case for structs is for holding data but _currently_ we have no state, so it doesn't feel right to be creating one.

[HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc) lets us avoid this.

> The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers. If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f. 

```go
type HandlerFunc func(ResponseWriter, *Request)
```

So we use this to wrap our `PlayerServer` function so that it now conforms to `Handler`.


### `http.ListenAndServe(":5000"...`

ListenAndServer takes a port to listen on and a `Handler` indefinitely. If the port is already being listened to it will return an `error` so we are using an `if` statement to capture that scenario and log the problem to the user.

