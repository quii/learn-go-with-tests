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

- We make a `Request` to send to our function. The `nil` argument refers to the request's body, which we don't need to set in this case.
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

We want to wire this up into an application. Why? All it does right now has little to do with the requirements given

- Have _actual working software_, we don't want to write tests for the sake of it, it's good to see the code in action.
- As we refactor our code, it's likely we will change the structure of the program. We want to make sure this is reflected in our application too as part of the incremental approach.

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

Next let's make our endpoint return a player's score. 

## Write the test first

We "know" that we need the concept of a `PlayerStore` at some point, but let's try and take smaller steps for now

```go
func TestGETPlayers(t *testing.T) {
	t.Run("returns the player's score", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		res := httptest.NewRecorder()

		PlayerServer(res, req)

		got := res.Body.String()
		want := "20"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})
}
```

We're now putting in something a bit more concrete in our test by saying if you ask for `/player/Pepper` you should get back 20.

We also know in the current state the resulting code to make it will be a little silly, but suspend your reservations for now.

## Try to run the test

```go
=== RUN   TestGETPlayers/returns_the_player's_score
    --- FAIL: TestGETPlayers/returns_the_player's_score (0.00s)
    	server_test.go:20: got 'Hello, world', want '20'
```

## Write enough code to make it pass

```go
// PlayerServer currently returns Hello, world given _any_ request
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
```

Yes, it's silly!

What we're going to do now is write _another_ test to force us into making a positive change

## Write the test first

```go
	t.Run("returns Floyd's score", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
		res := httptest.NewRecorder()

		PlayerServer(res, req)

		got := res.Body.String()
		want := "10"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})
```

## Try to run the test
```
=== RUN   TestGETPlayers/returns_the_Pepper's_score
    --- PASS: TestGETPlayers/returns_the_Pepper's_score (0.00s)
=== RUN   TestGETPlayers/returns_Floyd's_score
    --- FAIL: TestGETPlayers/returns_Floyd's_score (0.00s)
    	server_test.go:34: got '20', want '10'
```

## Write the minimal amount of code for the test to run and check the failing test output

By doing this the test has forced us to actually look at the request's URL and make some decision. So whilst in our heads we may have been worrying about stores and interfaces the next logical step actually seems to be about _routing_.

## Write enough code to make it pass

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/player/"):]

	if player == "Pepper" {
		fmt.Fprint(w, "20")
		return
	}

	if player == "Floyd" {
		fmt.Fprint(w, "10")
		return
	}
}
```

We're resisting the temptation to use any routing libraries right now, just the smallest step to get our test passing.

`r.URL.Path` returns the path of the request and then we are using slice syntax to slice it past the final slash after player. It's not very robust but will do the trick for now.

## Refactor

We can simplify the `PlayerServer` by separating out the score retrieval into a function

```go
// PlayerServer currently returns Hello, world given _any_ request
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/player/"):]

	fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(name string) string {
	if name == "Pepper" {
		return "20"
	}

	if name == "Floyd" {
		return "10"
	}

	return ""
}
```

And we can DRY up some of the code in the tests by making some helpers

```go
func TestGETPlayers(t *testing.T) {
	t.Run("returns the Pepper's score", func(t *testing.T) {
		req := newGetScoreRequest("Pepper")
		res := httptest.NewRecorder()

		PlayerServer(res, req)

		assertResponseBody(t, res.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		req := newGetScoreRequest("Floyd")
		res := httptest.NewRecorder()

		PlayerServer(res, req)

		assertResponseBody(t, res.Body.String(), "10")
	})
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got '%s' want '%s'", got, want)
	}
}
```

However we still shouldn't be happy. It doesn't feel right that our server knows the scores. 

Our refactoring has made it pretty clear what to do. 

We moved the score calculation out of the main body of our handler into a function `GetPlayerScore`. This feels like the right place to slice our functionality up using interfaces. 

Let's move our function we re-factored to be an interface instead

```go
type PlayerStore interface {
	GetPlayerScore(name string) string
}
```

For our `PlayerServer` to be able to use a `PlayerStore`, it will need a reference to one. Now feels like the right time to change our architecture so that our `PlayerServer` is now a `struct`

```go
type PlayerServer struct {
	store PlayerStore
}
```

Finally, we will now implement the `Handler` interface by adding a method to our new struct and putting in our existing handler code

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/player/"):]
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

The only other change is we now call our `store.GetPlayerStore` to get the score, rather than the local function we defined (which we can now delete).

Here is the full code listing of our server

```go
type PlayerStore interface {
	GetPlayerScore(name string) string
}

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/player/"):]
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

#### Fix the issues

This was quite a few changes and we know our tests and application will no longer compiler but just relax and let the compiler work through it

`./main.go:9:58: type PlayerServer is not an expression`

If the compiler was really nice it would say something like `PlayerServer is not a function, it's a type, dummy`. 

We need to change our tests to instead create a new instance of our `PlayerServer`

```go
func TestGETPlayers(t *testing.T) {
	t.Run("returns the Pepper's score", func(t *testing.T) {
		server := &PlayerServer{}
		req := newGetScoreRequest("Pepper")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertResponseBody(t, res.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		server := &PlayerServer{}
		req := newGetScoreRequest("Floyd")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertResponseBody(t, res.Body.String(), "10")
	})
}
```

Notice we're still not worrying about making stores _just yet_, we want to get into the state of the code at least compiling as soon as we can.

Now `main.go` won't compile for the same reason.

```go
func main() {
	server := &PlayerServer{}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

Finally everything is compiling but the tests are failing

```go
=== RUN   TestGETPlayers/returns_the_Pepper's_score
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
```

This is because we have not passed in a `PlayerStore` in our tests. We'll need to make a stub one up. 

```go
type StubPlayerStore struct {
	scores map[string]string
}

func (s *StubPlayerStore) GetPlayerScore(name string) string {
	score := s.scores[name]
	return score
}
```

A `map` is a quick and easy way of making a stub key/value store for our tests. Now let's create one of these stores for our tests and send it into our `PlayerServer`

```go
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]string{
			"Pepper": "20",
			"Floyd":  "10",
		},
	}
	server := &PlayerServer{&store}

	t.Run("returns the Pepper's score", func(t *testing.T) {
		req := newGetScoreRequest("Pepper")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertResponseBody(t, res.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		req := newGetScoreRequest("Floyd")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertResponseBody(t, res.Body.String(), "10")
	})
}
```

Our tests now pass and are looking better. The _intent_ behind our code is clearer now due to the introduction of the store. We're telling the reader that because we have _this data in a `PlayerStore`_ that when you use it with a `PlayerServer` you should get the following responses.
