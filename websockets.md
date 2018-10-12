# Websockets (WIP)

In this chapter we'll learn how to use websockets to improve our application. 

## Project recap

We have two applications in our poker codebase

- *Command line app*. Prompts the user to enter the number of players in a game. From then on informs the players of what the "blind bet" value is, which increases over time. At any point a user can enter `"{Playername} wins"` to finish the game and record the victor in a store.
- *Web app*. Allows users to record winners of games and displays a league table. Shares the same store as the command line app. 

## Next steps

The product owner is thrilled with the command line application but would prefer it if we could bring that functionality to the browser. She imagines a web page with a text box that allows the user to enter the number of players and when they submit the form the page displays the blind value and automatically updates it when appropriate. Like the command line application the user can declare the winner and it'll get saved in the database.

On the face of it, it sounds quite simple but as always we must emphasise taking an _iterative_ approach to writing software.

First of all we will need to serve HTML. So far all of our HTTP endpoints have returned either plaintext or JSON. We _could_ use the same techniques we know (as they're all ultimately strings) but we can also use the [/html/template](https://golang.org/pkg/html/template/) for a cleaner solution.

We also need to be able to asynchronously send messages to the user saying `The blind is now *y*` without having to refresh the browser. We can use [websockets](https://en.wikipedia.org/wiki/WebSocket) to facilitate this. 

> WebSocket is a computer communications protocol, providing full-duplex communication channels over a single TCP connection

Given we are taking on a number of techniques it's even more important we do the smallest amount of useful work possible first and then iterate. 

For that reason the first thing we'll do is create a web page with a form for the user to record a winner. Rather than using a plain form, we will use websockets to send that data to our server for it to record. 

After that we'll work on the blind alerts by which point we will have a bit of infrastructure code set up. 

### What about tests for the JavaScript ?

There will be some JavaScript written to do this but I wont go in to writing tests. 

It is of course possible but for the sake of brevity I wont be including any explanations for it. 

Sorry folks. Lobby O'Reilly to pay me to make a "Learn JavaScript with tests".

## Write the test first

First thing we need to do is serve up some HTML to the user when they hit `/game`. 

Here's a reminder of the pertinent code in our web server

```go
// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

// NewPlayerServer creates a PlayerServer with routing configured
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

The _easiest_ thing we can do for now is check when we `GET /game` that we get a `200`.

```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}
```

## Try to run the test
```
--- FAIL: TestGame (0.00s)
=== RUN   TestGame/GET_/game_returns_200
    --- FAIL: TestGame/GET_/game_returns_200 (0.00s)
    	server_test.go:109: did not get correct status, got 404, want 200
```

## Write enough code to make it pass

Our server has a router setup so it's relatively easy to fix. 

To our router add

```go
router.Handle("/game", http.HandlerFunc(p.game))
```

And then write the `game` method

```go
func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
```

## Refactor

The server code is already fine due to us slotting in more code into the existing well-factored code very easily. 

We can tidy up the test a little by adding a helper (write this yourself) to make the request to `/game`

```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request :=  newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})
}
```

You'll also notice I changed `assertStatus` to accept `response` rather than `response.Code` as I feel it reads better.


note to self: The key is to refactor `game` so that start takes a destination as to where to write the blind things
