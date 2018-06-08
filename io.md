# IO and sorting (WIP)

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/master/io)**

[In the previous chapter](json.md) we continued iterating on our application by adding a new endpoint `/league`. Along the way we learned about how to deal with JSON, embedding types and routing.

Our product-owner is somewhat perturbed by the software losing the scores when the server was restarted. This is because our implementation of our store is in-memory. She is also not pleased that we didn't interpret the `/league` endpoint should return the players ordered by number of wins!

## The code so far

```go
// server.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() []Player
}

// Player stores a name with a number of wins
type Player struct {
	Name string
	Wins int
}

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

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.Header().Set("content-type", jsonContentType)
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

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) GetLeague() []Player {
	var league []Player
	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return league
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

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
	server := NewPlayerServer(NewInMemoryPlayerStore())

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

You can find the corresponding tests in the link at the top of the chapter.

## Store the data

There are dozens of databases we could use for this but we're going to go for a very simple approach. We're going to store the data for this application in a file as JSON.

This keeps the data very portable and is relatively simple to implement.

It wont scale especially well but given this is a prototype it'll be fine for now. If our circumstances change and it's no longer appropriate it'll be simple to swap it out for something different because of the `PlayerStore` abstraction we have used.

We will keep the `InMemoryPlayerStore` for now so that the integration tests keep passing as we develop our new store. Once we are confident our new implementation is sufficient to make the integration test pass we will swap it in and then delete `InMemoryPlayerStore`.

## Write the test first

By now you should be familiar with the interfaces around the standard library for reading data (`io.Reader`), writing data (`io.Writer`) and how we can use the standard library to test these functions without having to use real files.

For this work to be complete we'll need to implement `PlayerStore` so we'll write tests for our store calling the methods we need to implement. We'll start with `GetLeague`.

```go
func TestFileSystemStore(t *testing.T) {

	t.Run("/league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)
	})
}
```

We're using `strings.NewReader` which will return us a `Reader`, which is what our `FileSystemStore` will use to read data. In `main` we will open a file, which is also a `Reader`.

## Try to run the test

```
# github.com/quii/learn-go-with-tests/json-and-io/v7
./FileSystemStore_test.go:15:12: undefined: FileSystemStore
```

## Write the minimal amount of code for the test to run and check the failing test output

Let's define `FileSystemStore` in a new file

```go
type FileSystemStore struct {}
```

Try again

```
# github.com/quii/learn-go-with-tests/json-and-io/v7
./FileSystemStore_test.go:15:28: too many values in struct initializer
./FileSystemStore_test.go:17:15: store.GetLeague undefined (type FileSystemStore has no field or method GetLeague)
```

It's complaining because we're passing in a `Reader` but not expecting one and it doesnt have `GetLeague` defined yet.

```go
type FileSystemStore struct {
	database io.Reader
}

func (f *FileSystemStore) GetLeague() []Player {
	return nil
}
```

One more try...

```
=== RUN   TestFileSystemStore//league_from_a_reader
    --- FAIL: TestFileSystemStore//league_from_a_reader (0.00s)
    	FileSystemStore_test.go:24: got [] want [{Cleo 10} {Chris 33}]
```

## Write enough code to make it pass

We've read JSON from a reader before

```go
func (f *FileSystemStore) GetLeague() []Player {
	var league []Player
	json.NewDecoder(f.database).Decode(&league)
	return league
}
```

The test should pass.

## Refactor

We _have_ done this before! Our test code for the server had to decode the JSON from the response.

Let's try DRYing this up into a function.

Create a new file called `league.go` and put this inside.

```go
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	return league, err
}
```

Call this in our implementation and in our test helper `getLeagueFromResponse` in `server_test.go`

```go
func (f *FileSystemStore) GetLeague() []Player {
	league, _ := NewLeague(f.database)
	return league
}
```

We haven't got a strategy yet for dealing with parsing errors but let's press on.

### Seeking problems

There is a flaw in our implementation. First of all let's remind ourselves how `io.Reader` is defined.

```go
type Reader interface {
        Read(p []byte) (n int, err error)
}
```

With our file you can imagine it reading through byte by byte until the end. What happens if you try and `Read` a second time?

Add the following to the end of our current test.

```go
// read again
got = store.GetLeague()
assertLeague(t, got, want)
```

We want this to pass, but if you run the test it doesn't.

The problem is our `Reader` has reached to the end so there is nothing more to read. We need a way to tell it to go back to the start.

[ReadSeeker](https://golang.org/pkg/io/#ReadSeeker) is another interface in the standard library that can help.

```go
type ReadSeeker interface {
        Reader
        Seeker
}
```

Remember embedding? This is an interface comprised of `Reader` and [`Seeker`](https://golang.org/pkg/io/#Seeker)

```go
type Seeker interface {
        Seek(offset int64, whence int) (int64, error)
}
```

This sounds good, can we change `FileSystemStore` to take this interface instead?

```go
type FileSystemStore struct {
	database io.ReadSeeker
}

func (f *FileSystemStore) GetLeague() []Player {
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
}
```

Try running the test, it now passes! Happily for us `string.NewReader` that we used in our test also implements `ReadSeeker` so we didn't have to make any other changes.

Next we'll implement `GetPlayerScore`

## Write the test first

```go
t.Run("get player score", func(t *testing.T) {
    database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)

    store := FileSystemPlayerStore{database}

    got := store.GetPlayerScore("Chris")

    want := 33

    if got != want {
        t.Errorf("got %d want %d", got, want)
    }
})
```

## Try to run the test

`./FileSystemStore_test.go:38:15: store.GetPlayerScore undefined (type FileSystemPlayerStore has no field or method GetPlayerScore)`

## Write the minimal amount of code for the test to run and check the failing test output

We need to add the method to our new type to get the test to compile.

```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	return 0
}
```

Now it compiles and the test fails

```
=== RUN   TestFileSystemStore/get_player_score
    --- FAIL: TestFileSystemStore//get_player_score (0.00s)
    	FileSystemStore_test.go:43: got 0 want 33
```

## Write enough code to make it pass

We can iterate over the league to find the player and return their score

```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	var wins int

	for _, player := range f.GetLeague() {
		if player.Name == name {
			wins = player.Wins
			break
		}
	}

	return wins
}
```

## Refactor

You will have seen dozens of test helper refactorings so I'll leave this to you to make it work

```go
t.Run("/get player score", func(t *testing.T) {
    database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)

    store := FileSystemPlayerStore{database}

    got := store.GetPlayerScore("Chris")
    want := 33
    assertScoreEquals(t, got, want)
})
```

Finally we need to start recording scores with `RecordWin`

## Write the test first

Our approach is fairly short-sighted for writes. We cant (easily) just update one "row" of JSON in a file. We'll need to store the _whole_ new representation of our database on every write.

How do we write? We'd normally use a `Writer` but we already have our `ReadSeeker`. Potentially we could have two dependencies but the standard library already has an interface for us `ReadWriteSeeker` which lets us do all the things we'll need to do with a file.

Let's update our type

```go
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}
```

See if it compiles

```go
./FileSystemStore_test.go:15:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
	*strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
./FileSystemStore_test.go:36:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
	*strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
```

It's not too surprising that `strings.Reader` does not implement `ReadWriteSeeker` so what do we do?

We have two choices

- We create a temporary file for each test. `*os.File` implements `ReadWriteSeeker`. The pro of this is it becomes more of an integration test, we're really reading and writing from the file system so it will give us a very high level of confidence. The cons are we prefer unit tests because they are faster and generally simpler. We will also need to do more work around creating temporary files and then making sure they're removed after the test.
- We could use a third party library. [github.com/mattetti](Mattetti) has written a library [filebuffer](https://github.com/mattetti/filebuffer) which implements the interface we need and doesn't touch the file system.

I don't think there's an especially wrong answer here, but by choosing to use a third party library I would have to explain dependency management! So we will use files instead.

Before adding our test we need to make our other tests compile by replacing the `strings.Reader` with an `os.File`.

Let's create a helper function which will create a temporary file with some data inside it

```go
func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
```

[TempFile](https://golang.org/pkg/io/ioutil/#TempDir) creates a temporary file for us to use. The `"db"` value we've passed in is a prefix put on a random file name it will create. This is to ensure it wont clash with other files by accident.

You'll notice we're not only returning our `ReadWriteSeeker` (the file) but also a function. We need to make sure that the file is removed once the test is finished. We don't want to leak details of the files into the test as it's prone to error and uninteresting for the reader. By returning a `removeFile` function, we can take care of the details in our helper and all the caller has to do is run `defer cleanDatabase()`.

```go
func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Chris")
		want := 33
		assertScoreEquals(t, got, want)
	})
}
```

Run the tests and they should be passing! That was a fair amount of changes but now it feels like we have our interface definition complete and it should be very easy to add new tests from now.

Let's get the first iteration of recording a win for an existing player

```go
t.Run("store wins for existing players", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    store.RecordWin("Chris")

    got := store.GetPlayerScore("Chris")
    want := 34
    assertScoreEquals(t, got, want)
})
```

## Try to run the test
`./FileSystemStore_test.go:67:8: store.RecordWin undefined (type FileSystemPlayerStore has no field or method RecordWin)`

## Write the minimal amount of code for the test to run and check the failing test output

Add the new method

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {

}
```

```
=== RUN   TestFileSystemStore/store_wins_for_existing_players
    --- FAIL: TestFileSystemStore/store_wins_for_existing_players (0.00s)
    	FileSystemStore_test.go:71: got 33 want 34
```

Our implementation is empty so the old score is getting returned

## Write enough code to make it pass

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()

	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
		}
	}

	f.database.Seek(0,0)
	json.NewEncoder(f.database).Encode(league)
}
```

You may be asking yourself why I am doing `league[i].Wins++` rather than `player.Wins++`.

When you `range` over a slice you are returned the current index of the loop (in our case `i`) and a _copy_ of the element at that index. Changing the `Wins` value of a copy wont have any effect on the `league` slice that we iterate on. For that reason we need to get the reference to the actual value by doing `league[i]` and then changing that value instead.

If you run the tests, they should now be passing.

## Refactor

In `GetPlayerScore` and `RecordWin` we are iterating over `[]Player` to find a player by name.

We could refactor this common code in the internals of `FileSystemStore` but to me it feels like this is maybe useful code we can lift into a new type. Working with a "League" so far has always been with `[]Player` but we can create a new type called `League`. This will be easier for other developers to understand and then we can attach useful methods onto that type for us to use.

Inside `league.go` add the following

```go
type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name==name {
			return &l[i]
		}
	}
	return nil
}
```

Now if anyone has a `League` they can easily find a given player.

And change our `PlayerStore` interface to return `League` rather than `[]Player`. Try and re-run the tests, you'll get a compilation problem because we've changed the interface but it's very easy to fix; just change the return type from `[]Player` to `League`.

This lets us simplify our methods in `FileSystemStore`.

```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	player := f.GetLeague().Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}
```

This is looking much better and we can see how we might be able to find other useful functionality around `League` can be refactored.

We now need to handle the scenario of recording wins of new players.

## Write the test first
```go
t.Run("store wins for existing players", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    store.RecordWin("Pepper")

    got := store.GetPlayerScore("Pepper")
    want := 1
    assertScoreEquals(t, got, want)
})
```
## Try to run the test

```
=== RUN   TestFileSystemStore/store_wins_for_existing_players#01
    --- FAIL: TestFileSystemStore/store_wins_for_existing_players#01 (0.00s)
    	FileSystemStore_test.go:86: got 0 want 1
```
## Write enough code to make it pass

We just need to handle the scenario where `Find` returns `nil` because it couldn't find the player.

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}
```

The happy path is looking ok so we can now try using our new `Store` in the integration test. This will give us more confidence that the software works and then we can delete the redundant `InMemoryPlayerStore`.

In `TestRecordingWinsAndRetrievingThem` replace the old store.

```go
database, cleanDatabase := createTempFile(t, "")
defer cleanDatabase()
store := &FileSystemPlayerStore{database}
```

In you run the test it should pass and now we can delete `InMemoryPlayerStore`. `main.go` will now have compilation problems which will motivate us to now use our new store in the "real" code.

```go
package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store := &FileSystemPlayerStore{db}
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

- We create a file for our database.
- The 2nd argument to `os.OpenFile` lets you define the permissions for opening the file, in our case `O_RDWR` means we want to read and write _and_ `os.O_CREATE` means create the file if it doesn't exist.
- The 3rd argument means sets permissions for the file, in our case all users can read and write the file. [(See superuser.com for a more detailed explanation)](https://superuser.com/questions/295591/what-is-the-meaning-of-chmod-666)

Running the program now persists the data in a file in between restarts, hooray!

## Error handling

Before we start working on sorting we should make sure we're happy with our current code and remove any technical debt we may have. It's an important principle to get to working software as quickly as possible (stay out of the red state) but that doesn't mean we should ignore error cases!

If we go back to `FileSystemStore.go` we have

`league, _ := NewLeague(f.database)`

`NewLeague` can return an error if it is unable to parse the league from the `io.Reader` that we provide.

It was pragmatic to ignore that at the time as we already had failing tests. If we had tried to tackle it at the same time we would be juggling two things at once.

If we get an error we'll want to inform the user there was some kind of problem by returning a `500` status code and some kind of message. We'll also want to log it, but we'll get onto that in a later chapter.

Let's try and return the error in our function

```go
func (f *FileSystemPlayerStore) GetLeague() (League, error) {
	f.database.Seek(0, 0)
	return NewLeague(f.database)
}
```

### Try and compile

```
./FileSystemStore.go:22:23: multiple-value f.GetLeague() in single-value context
./FileSystemStore.go:33:23: multiple-value f.GetLeague() in single-value context
./main.go:19:27: cannot use store (type *FileSystemPlayerStore) as type PlayerStore in argument to NewPlayerServer:
	*FileSystemPlayerStore does not implement PlayerStore (wrong type for GetLeague method)
		have GetLeague() (League, error)
		want GetLeague() League
./FileSystemStore_test.go:38:25: multiple-value store.GetLeague() in single-value context
./FileSystemStore_test.go:48:24: multiple-value store.GetLeague() in single-value context
./server_integration_test.go:13:27: cannot use store (type *FileSystemPlayerStore) as type PlayerStore in argument to NewPlayerServer:
	*FileSystemPlayerStore does not implement PlayerStore (wrong type for GetLeague method)
		have GetLeague() (League, error)
		want GetLeague() League
```

This looks bad, but again this is actually a good thing. In a dynamic language you would not get the computer telling you exactly what is wrong.

We just need to work through the errors from the top until we get green again. Make sure after every change to try recompiling to ensure that the change you do has fixed the problem.

Once we've done that we can take stock of the changes and add whatever tests we feel we need.

> `./FileSystemStore.go:22:23: multiple-value f.GetLeague() in single-value context`

```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	league, _ := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}
```

We are going to ignore the error here for now as we're just trying to get green again as quickly as possible.

> `./FileSystemStore.go:33:23: multiple-value f.GetLeague() in single-value context`

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league, _ := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
}
```

Same deal again, just ignore it for now

> `./main.go:19:27: cannot use store (type *FileSystemPlayerStore) as type PlayerStore in argument to NewPlayerServer:
	*FileSystemPlayerStore does not implement PlayerStore (wrong type for GetLeague method)
		have GetLeague() (League, error)
		want GetLeague() League`

Our `FileSystemStore` no longer implements `PlayerStore`. We will have to update the interface to work with our new reality that `GetLeague` could fail.

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() (League, error)
}
```

Try and compile again

```
./server.go:47:27: too many arguments in call to json.NewEncoder(w).Encode
	have (League, error)
	want (interface {})
./FileSystemStore_test.go:38:25: multiple-value store.GetLeague() in single-value context
./FileSystemStore_test.go:48:24: multiple-value store.GetLeague() in single-value context
./server_test.go:40:28: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
	*StubPlayerStore does not implement PlayerStore (wrong type for GetLeague method)
		have GetLeague() League
		want GetLeague() (League, error)
./server_test.go:78:28: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
	*StubPlayerStore does not implement PlayerStore (wrong type for GetLeague method)
		have GetLeague() League
		want GetLeague() (League, error)
./server_test.go:110:29: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
	*StubPlayerStore does not implement PlayerStore (wrong type for GetLeague method)
		have GetLeague() League
		want GetLeague() (League, error)
```

Keep going! The compiler is helping us make robust software, just keep ticking off the errors.

> `./server.go:47:27: too many arguments in call to json.NewEncoder(w).Encode
  	have (League, error)
  	want (interface {})`

This is good, this is where we'll actually have to handle the error but let's resist the temptation for now. We'll want to write a test to exercise this scenario but we musn't add any more code than necessary while we are in a state of the code not compiling

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	//todo: handle the error by responding differently
	league, _ := p.store.GetLeague()
	json.NewEncoder(w).Encode(league)
}
```

>./FileSystemStore_test.go:38:25: multiple-value store.GetLeague() in single-value context
 ./FileSystemStore_test.go:48:24: multiple-value store.GetLeague() in single-value context

These two scenarios are the same, just fix the tests by using the underscore syntax to ignore the error.

Change `got := store.GetLeague()` to `got, _ := store.GetLeague()` and leave a `todo` to remind ourselves to assert there are no errors later.

> `./server_test.go:40:28: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
   	*StubPlayerStore does not implement PlayerStore (wrong type for GetLeague method)
   		have GetLeague() League
   		want GetLeague() (League, error)`

We changed the interface of `PlayerStore` so `StubPlayerStore` needs updating.

```go
func (s *StubPlayerStore) GetLeague() (League, error) {
	return s.league, nil
}
```

When you try and run the tests it should now all be passing. We have updated our `PlayerStore` interface to reflect the new reality of stores that can fail which will enable us to handle errors better.

This may have felt arduous but once you become familiar with compiler errors and are handy with your tooling fixing this kind of error only really takes a few minutes.

Now we can write a test for our `Server` to log and respond with a `500` when we cannot load the league.

## Write the test first

```go
t.Run("it returns a 500 when the league cannot be loaded", func(t *testing.T) {

    store := FailingPlayerStore{}
    server := NewPlayerServer(&store)

    request := newLeagueRequest()
    response := httptest.NewRecorder()

    server.ServeHTTP(response, request)

    assertStatus(t, response.Code, http.StatusInternalServerError)
})
```

What's a `FailingPlayerStore?`. We could've added some flexibility to our `StubPlayerStore` to somehow make it so it fails given some kind of `fail` flag but I felt it would be simpler and perhaps clearer to make a new stub that explicitly fails.

```go
type FailingPlayerStore struct {
	PlayerStore
}

func (f * FailingPlayerStore) GetLeague() (League, error) {
	return League{}, errors.New("cannot load league")
}
```

I did not want to have to implement the _whole_ interface (e.g also have methods for `GetPlayerScore` and `RecordWin`) for this test so i _embedded the interface_ into our new type. By doing this our new stub implements `PlayerStore` and then we add our specific implementation for `GetLeague` to make it fail. If you try and call any of the other methods that we have not implemented it will panic. This technique is useful when you want to mock an interface with multiple methods but you're not concerned with every method.

## Try to run the test

```
=== RUN   TestLeague/it_returns_a_500_when_the_league_cannot_be_loaded
    --- FAIL: TestLeague/it_returns_a_500_when_the_league_cannot_be_loaded (0.00s)
    	server_test.go:144: did not get correct status, got 200, want 500
```

## Write enough code to make it pass

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	league, err := p.store.GetLeague()

	if err != nil {
		http.Error(w, fmt.Sprintf("could not load league, %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(league)
}
```

`http.Error` is a convenient method for when you want to return an error. We are using `fmt.Sprintf` to ensure we add some context to the error message.

## More cleaning up

As we changed the `GetLeague` interface we left some `TODO`s around for us to come back to. `TODO` is often a dangerous tool which could be renamed to `SOMEDAY I WILL TACKLE THIS, BUT WHATEVER`. We are better than this!

In our `FileSystemStoreTest` we need to update the tests to check we don't get an error.

Make a helper and then fix the `TODO`s

```go
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.FatalF("unexpected error %v", err)
	}
}
```

We should probably check that if we cannot read the JSON into a league that it actually fails so let's add another test to this suite.

```go
t.Run("return an error when league cannot be read", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `not very good JSON`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    _, err := store.GetLeague()

    if err == nil {
        t.Error("expected an error but didn't get one")
    }
})
```

If we run this test it actually passes. To check it works how we'd hope, change `GetLeague` to return `nil` for the error in all scenarios and check the test output is what you expect. It's very important you check tests fail how you expect them if you didn't follow the strict TDD cycle.

## Remaining technical debt

As we changed the `PlayerStore` interface for `GetLeague` we had to take on more technical debt in `FileSystemStore`.

In both `GetPlayerScore` and `RecordWin` we have this line.

```go
league, _ := f.GetLeague()
```

Getting leagues can fail so therefore these other two methods can fail too.

We need to go through the same exercise again of changing the interface of these two methods to be able to return `error`.

We're not going to document the process again, this is an exercise for you to do.

Commit your code to source control first in-case you get stuck.

Just follow these steps for each method carefully, trying to re-run the compiler after every change.

1. Try and return the error
2. Compile
3. The compiler will complain the method should only return one thing. Change the method signature, compile again
4. The compiler will now complain that `FileSystemStore` no longer implements the interface, so change it
5. The compiler will complain about `StubFileStore` no longer implements the interface, so fix it
6. The compiler will complain about `multiple-value p.store.XXX() in single-value context`, so fix them and handle the errors. For the server return a `500` (don't forget to write a test) and in the tests ensure an error isn't returned.

tl;dr - Make the change you want and use the compiler to help you get back to working code.

If you get stuck, start over. If you get really stuck, [have a look at the current state of the code here](https://github.com/quii/learn-go-with-tests/tree/master/io/v7)

### Integration test woes

By no longer ignoring the errors when doing `GetLeague` we introduced a bug which is highlighted by our integration test.

```
=== RUN   TestRecordingWinsAndRetrievingThem/get_score
    --- FAIL: TestRecordingWinsAndRetrievingThem/get_score (0.00s)
    	server_integration_test.go:24: did not get correct status, got 500, want 200
    	server_integration_test.go:26: response body is wrong, got 'could not show score, EOF
    		' want '3'
```

The problem is when the file is empty (when we first start) it cannot be read into JSON.

When we were ignoring the error, we would then carry on and start afresh with a new database which is why it was passing before.

We need a way of making our `FileSystemStore` initialising itself if the file is empty. Thankfully we already have a failing test for this so we can just work with that.

```go
func NewFileSystemPlayerStore(database io.ReadWriteSeeker) (*FileSystemPlayerStore, error) {
	buf := &bytes.Buffer{}
	length, err := io.Copy(buf, database)

	if err != nil {
		return nil, err
	}

	if length == 0 {
		json.NewEncoder(database).Encode(League{})
	}

	return &FileSystemPlayerStore{
		database: database,
	}, nil
}
```

What we need to do is to make a "constructor" which
- Tries to read the `database`, `io.Copy` will return the number of bytes it has read.
- If there is an error then return it.
- If the length of the bytes read is zero then we need to encode an empty `League` into the database to initialise ourselves.

When we use this function in the integration test it now passes. Make sure to update `main` to use it too.

## Sorting

Our product owner wants `/league` to return the players sorted by their scores.

The main decision to make here is where in the software should this happen. If we were using a "real" database we would use things like `ORDER BY` so the sorting is super fast so for that reason it feels like implementations of `PlayerStore` should be responsible.

## Write the test first

We can update the assertion on our first test in `TestFileSystemStore`

```go
	t.Run("league from a reader, sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got, err := store.GetLeague()
		assertNoError(t, err)

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		// read again
		got, err = store.GetLeague()
		assertNoError(t, err)
		assertLeague(t, got, want)
	})

```

The order of the JSON coming in is in the wrong order and our `want` will check that it is returned to the caller in the correct order.

## Try to run the test

```
=== RUN   TestFileSystemStore/league_from_a_reader,_sorted
    --- FAIL: TestFileSystemStore/league_from_a_reader,_sorted (0.00s)
    	FileSystemStore_test.go:46: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
    	FileSystemStore_test.go:51: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
```

## Write enough code to make it pass

```go
func (f *FileSystemPlayerStore) GetLeague() (League, error) {
	f.database.Seek(0, 0)
	league, err := NewLeague(f.database)

	sort.Slice(league, func(i, j int) bool {
		return league[i].Wins > league[j].Wins
	})

	return league, err
}
```

[`sort.Slice`](https://golang.org/pkg/sort/#Slice)

>  Slice sorts the provided slice given the provided less function.

Easy!

## Wrapping up

What we've covered:

- The `Seeker` interface and its relation with `Reader` and `Writer`.
- Working with files.
- Returning errors as HTTP responses.
- Creating an easy to use helper for testing with files that hides all the messy stuff.
- Using embedding when you want to be lazy about mocking just a part of an interface.
- `sort.Slice` for sorting slices.
- Using the compiler to help us make structural changes to the application safely.

Where our software is at:

- We have a HTTP API where you can create players and increment their score.
- We can return a league of everyone's scores as JSON.
- The data is persisted as a JSON file.