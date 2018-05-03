# JSON and IO (WIP)

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/master/json-and-io)**

[In the previous chapter](http-server.md) we created a web server to store how many games a player has won.

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
