# JSON and IO (WIP)

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/master/json-and-io)**

[In the previous chapter](http-server.md) we created a web server to store how many games a players have won. 

Our product-owner was mostly delighted but was somewhat perturbed by the software losing the scores when the server was restarted. This was because our implementation of our store was in-memory. 

She also has a new requirement; to have a new endpoint called `/league` which returns a list of all players stored, ordered by wins. She would like this to be returned as JSON.

She says she doesn't have a preference for how we persist the scores, she trusts us - just make it work!

## Here is the code we have so far

```go
// server.go
package main

import (
	"fmt"
	"net/http"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
``` 

```go
// InMemoryPlayerStore.go
package main

// NewInMemoryPlayerStore initialises an empty player store
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

// InMemoryPlayerStore collects data about players in memory
type InMemoryPlayerStore struct {
	store map[string]int
}

// RecordWin will record a player's win
func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

// GetPlayerScore retrieves scores for a given player
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

```

```go
// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

You can find the corresponding tests in the link at the top of the chapter.

We'll start by making the league table endpoint.

## Write the test first

We'll extend the existing suite as we have some useful test functions and a fake `PlayerServer` to use.

```go
func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := &PlayerServer{&store}

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}
```

Before worrying about actual scores and JSON we will try and keep the changes small with the plan to iterate toward our goal. The simplest start is to check we can hit `/league` and get an `OK` back. 

## Try to run the test

```
=== RUN   TestLeague/it_returns_200_on_/league
panic: runtime error: slice bounds out of range [recovered]
	panic: runtime error: slice bounds out of range

goroutine 6 [running]:
testing.tRunner.func1(0xc42010c3c0)
	/usr/local/Cellar/go/1.10/libexec/src/testing/testing.go:742 +0x29d
panic(0x1274d60, 0x1438240)
	/usr/local/Cellar/go/1.10/libexec/src/runtime/panic.go:505 +0x229
github.com/quii/learn-go-with-tests/json-and-io/v2.(*PlayerServer).ServeHTTP(0xc420048d30, 0x12fc1c0, 0xc420010940, 0xc420116000)
	/Users/quii/go/src/github.com/quii/learn-go-with-tests/json-and-io/v2/server.go:20 +0xec
```

Your `PlayerServer` should be panicking like this. Go to the line of code in the stack trace which is pointing to `server.go`.

```go
player := r.URL.Path[len("/players/"):]
```

In the previous chapter we mentioned this was a fairly naive way of doing our routing. What is happening is it's trying to split the string of the path starting at an index beyond `/league` so it is `slice bounds out of range`

## Write enough code to make it pass

Go does have a built in routing mechanism called `ServeMux` (server multiplexer) which lets you attach `http.Handler`s to particular paths.

Let's commit some sins and get the tests passing in the quickest way we can, knowing we can refactor it with safety once we know the tests are passing

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	router.Handle("/players/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := r.URL.Path[len("/players/"):]

		switch r.Method {
		case http.MethodPost:
			p.processWin(w, player)
		case http.MethodGet:
			p.showScore(w, player)
		}
	}))

	router.ServeHTTP(w, r)
}
```

- When the request starts we create a router and then we tell it for `x` path using `y` handler.
- So for our new endpoint, we use `http.HandlerFunc` and an _anonymous function_ to `w.WriteHeader(http.StatusOK)` when `/league` is requested to make our new test pass
- For the `/players/` route we just cut and paste our code into another `http.HandlerFunc`
- Finally we handle the request that came in by calling our new router's `ServeHTTP` (notice how `ServeMux` is _also_ a `http.Handler`?)

If you run all the tests, it should all be passing.

## Refactor

There's a few improvements we can make.

`ServeHTTP` is looking quite big, we can separate things out a bit by refactoring our handlers into separate methods. 

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	router.ServeHTTP(w, r)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}
```

Looking better! 

Next, it's quite odd (and inefficient) to be setting up a router as a request comes in and then calling it. What we ideally want to do is have some kind of `NewPlayerServer` function which will take our dependencies and do the one time setup of creating the router. Each request can then just use that one instance of the router.

Here are the relevant changes

```go
type PlayerServer struct {
	store  PlayerStore
	router *http.ServeMux
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := &PlayerServer{
		store,
		http.NewServeMux(),
	}

	p.router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	p.router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	return p
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
```

- `PlayerServer` now needs to store a router
- We have moved the routing creation out of `ServeHTTP` and into our `NewPlayerServer` so this only has to be done once, not per request
- You will need to update all the test and production code where we used to do `PlayerServer{&store}` with `NewPlayerServer(&store)`

### One final refactor

Try changing the code to the following

```go
type PlayerServer struct {
	store  PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}
```

Finally make sure you **delete** `func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request)` as it is no longer needed!

### Embedding 

The first change is the second property of `PlayerServer`, we have removed the name of the field (it was `router`); this is called _embedding_. 

> Go does not provide the typical, type-driven notion of subclassing, but it does have the ability to “borrow” pieces of an implementation by embedding types within a struct or interface. 

[Effective Go - Embedding](https://golang.org/doc/effective_go.html#embedding)

What this means is that our `PlayerServer` now has all the methods that `http.ServeMux` has.

By doing this our type now implements `http.Handler` by virtue of having `http.ServeMux` embedded in it (because it has the method `ServeHTTP`). 

This lets us remove our own `ServeHTTP` method, as we are already exposing one. 

Embedding is a very interesting language feature. You can use it with interfaces too to compose new interfaces

```go
type Animal interface{
	Eater()
	Sleeper()
}
```

#### Any downsides?

You must be careful with embedding types because you will expose all public methods and properties of the type you embed. In our case it is ok because we embedded just the _interface_ that we wanted to expose (`http.Handler`) 

If we had been lazy and embedded `http.ServeMux` instead (the concrete type) it would still work _but_ users of `PlayerServer` would be able to add new routes to our server because `Handle(path, handler)` would be public.

**When embedding types, really think about what impact that has on your public API**

# Wrapping up

What we've covered:

- `struct` embedding
- JSON deserializing and serializing
- Routing