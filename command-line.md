# Command line (WIP)

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/master/command-line)**

Our product owner now wants to _pivot_ by introducing a second application. This will be a command line app which helps a group of people play Texas-Holdem Poker.

## Just enough information on poker

- N number of players sit in a circle.
- There is a dealer button, which gets passed to the left every round.
- To the left of the dealer there is the "small blind".
- To the left of the small blind there is the big blind (or just blind).
- These players *have* to contribute chips to the "pot". The big blind contributing twice as much as the small.
- This way, every player has to contribute money to the pot, forcing people to play the game rather than "folding" all the time.
- The amount of chips the player has to contribute as a blind bet increases over time to ensure the game doesn't last too long.

Our application will help keep track of when the blind should go up, and how much it should be.

- Create a command line app.
- When it starts it asks how many players are playing. This determines the amount of time there is before the "blind" bet goes up.
  - There is a base amount of time of 15 minutes.
  - For every player, 1 minute is added.
  - e.g 6 players equals 21 minutes for the blind.
- After the blind time expires the game should alert the players the new amount the blind bet is.
- The blind starts at 100 chips, then 200, 400, 600, 1000, 2000 and continue to double until the game ends.
- When the game ends the user should be able to type "Chris wins" and that will record a win for the player in our existing database. This should then exit the program.

The product owner wants the database to be shared amongst the two applications so that the league update according to wins recorded in the new application.

## A reminder of the code

We have an application with a `main.go` file that launches a HTTP server. The HTTP server wont be interesting to us for this exercise but the abstraction it uses will. It depends on a `PlayerStore`.

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
    GetLeague() League
}
```

In the previous chapter we made a `FileSystemPlayerStore` which implements that interface. We should be able to re-use some of this for our new application

## Some project refactoring first

Our project now needs to create two binaries, our existing web server and the command line app.

Before we get stuck in to our new work we should structure our project to accommodate this.

So far all the work has lived in one folder, and we'll assume the code on your computer is living somewhere like

`$GOPATH/src/github.com/your-name/my-app`

It is good practice not to go over-the-top with package structure and thankfully it's pretty straightforward to add structure _when you need it_.

Inside the existing project create a `cmd` directory with a `webserver` directory inside that (e.g `mkdir -p cmd/webserver`).

Move the `main.go` inside there.

If you have `tree` installed you should run it and your structure should look like this

```
.
├── FileSystemStore.go
├── FileSystemStore_test.go
├── cmd
│   └── webserver
│       └── main.go
├── league.go
├── server.go
├── server_integration_test.go
├── server_test.go
├── tape.go
└── tape_test.go
```

We now effectively have a separation between our application and the library code but we now need to change some package names. Remember when you build a Go application it's package _must_ be `main`.

Change all the other code to have a package called `poker`.

Finally we need to import this package into `main.go` so we can use it to create our web server. Then we can use our library code by using `poker.FunctionName`

The paths will be different on your computer, but it should be similar to this:

```go
package main

import (
    "log"
    "net/http"
    "os"
    "github.com/quii/learn-go-with-tests/command-line/v1"
)

const dbFileName = "game.db.json"

func main() {
    db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("problem opening %s %v", dbFileName, err)
    }

    store, err := poker.NewFileSystemPlayerStore(db)

    if err != nil {
        log.Fatalf("problem creating file system player store, %v ", err)
    }

    server := poker.NewPlayerServer(store)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

### Final checks

- Inside the root run `go test` and check they're still passing
- Go inside our `cmd/webserver` and do `go run main.go`
  - Visit `http://localhost:5000/league` and you should see it's still working

### Walking skeleton

Before we get stuck in to writing tests, let's add a new application that our project will build. Create another directory inside `cmd` called `cli` (command line interface) and add a `main.go` with the following

```go
package main

import "fmt"

func main() {
	fmt.Println("Let's play poker")
}
```

The first requirement we'll tackle is recording a win when the user types `{PlayerName} wins`

## Write the test first

Inside `PokerCLI_test.go` (in the root of the project, not inside `cmd`)

We know we need to make something called `PokerCLI` which will allow us to `Play` poker. It'll need to read user input and then record wins to a `PlayerStore`.

Before we jump too far ahead though, let's just write a test to check it integrates with the `PlayerStore` how we'd like

```go
func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	cli := &PokerCLI{playerStore}
	cli.PlayPoker()

	if len(playerStore.winCalls) !=1 {
		t.Fatal("expected a win call but didnt get any")
	}
}
```

- We can use our `StubPlayerStore` from other tests. 
- We pass in our dependency into our not yet existing `PokerCLI` type
- Trigger the game by an unwritten `PlayPoker` method
- Check that a win is recorded

## Try to run the test

```
# github.com/quii/learn-go-with-tests/command-line/v2
./cli_test.go:25:10: undefined: CLI
```

## Write the minimal amount of code for the test to run and check the failing test output

At this point you should be comfortable enough to create our new `CLI` struct with the respective field for our dependency and add a method. 

You should end up with code like this

```go
type PokerCLI struct {
	playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {}
```

Remember we're just trying to get the test running so we can check the test fails how we'd hope

```
--- FAIL: TestCLI (0.00s)
	cli_test.go:30: expected a win call but didnt get any
FAIL
```

## Write enough code to make it pass

```go
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
```

That should make it pass. 

Next we need to simulate reading from `Stdin` (the input from the user) so that we can record wins for specific players.

Let's extend our test to exercise this

## Write the test first

```go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &PokerCLI{playerStore, in}
	cli.PlayPoker()

	if len(playerStore.winCalls) < 1 {
		t.Fatal("expected a win call but didnt get any")
	}

	got := playerStore.winCalls[0]
	want := "Chris"

	if got != want {
		t.Errorf("didnt record correct winner, got '%s', want '%s'", got, want)
	}
}
```

`os.Stdin` is what we'll use in `main` to capture the user's input. It is a `*File` under the hood which means it implements `io.Reader` which as we know by now is a handy way of capturing text.

We create a `io.Reader` in our test using the handy `strings.NewReader`, filling it with what we expect the user to type. 

## Try to run the test

`./PokerCLI_test.go:12:32: too many values in struct initializer`

## Write the minimal amount of code for the test to run and check the failing test output

We need to add our new dependency into `PokerCLI`.

```go
type PokerCLI struct {
	playerStore PlayerStore
	in io.Reader
}
```

## Write enough code to make it pass

```
--- FAIL: TestCLI (0.00s)
	PokerCLI_test.go:23: didnt record correct winner, got 'Cleo', want 'Chris'
FAIL
```

Remember to do the strictly easiest thing first

```go
func (cli *PokerCLI) PlayPoker() {
	cli.playerStore.RecordWin("Chris")
}
```

The test passes. We'll add another test to force us to write some real code next, but first let's refactor

## Refactor

In `server_test` we earlier did checks to see if wins are recorded like we have here. Let's DRY that assertion up into a helper

```go
func assertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], winner)
	}
}
```

Now replace the assertions in both `server_test.go` and `PokerCLI_test.go`

The test should now read like so

```go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &PokerCLI{playerStore, in}
	cli.PlayPoker()

	assertPlayerWin(t, playerStore, "Chris")
}
```

Now let's write _another_ test with different user input to force us into actually reading it.

## Write the test first

```go
func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}

		cli := &PokerCLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}

		cli := &PokerCLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Cleo")
	})

}
```

## Try to run the test

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/record_chris_win_from_user_input
    --- PASS: TestCLI/record_chris_win_from_user_input (0.00s)
=== RUN   TestCLI/record_cleo_win_from_user_input
    --- FAIL: TestCLI/record_cleo_win_from_user_input (0.00s)
    	PokerCLI_test.go:27: did not store correct winner got 'Chris' want 'Cleo'
FAIL
```

## Write enough code to make it pass

```go
func (cli *PokerCLI) PlayPoker() {
	userInput, _ := ioutil.ReadAll(cli.in)

	winner := strings.Replace(string(userInput), " wins\n", "", -1)
	cli.playerStore.RecordWin(winner)
}
```

The easiest way to make this test pass is: 
- Read all the contents of the string
- Extract out the winner by using `strings.Replace` which takes the string to replace, what string to replace, its replacement and finally a flag to say how many instances to replace

## Refactor

We can extract getting the winner's name into a meaningful function

```go
func (cli *PokerCLI) PlayPoker() {
	userInput, _ := ioutil.ReadAll(cli.in)

	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput []byte) string {
	return strings.Replace(string(userInput), " wins\n", "", 1)
}
```

Now that we have some working(ish) software, we should wire this up into `main`. Remember we should always strive to have fully-integrated working software as quickly as we can.

In `main.go` add the following and run it.

```go
package main

import (
	"fmt"
	"github.com/quii/learn-go-with-tests/command-line/v3"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	if err != nil{
		log.Fatalf("problem creating ")
	}
	
	game := poker.PokerCLI{store, os.Stdin}
	game.PlayPoker()
}
```

You should get an error

```
command-line/v3/cmd/cli/main.go:32:25: implicit assignment of unexported field 'playerStore' in poker.PokerCLI literal
command-line/v3/cmd/cli/main.go:32:34: implicit assignment of unexported field 'in' in poker.PokerCLI literal
```

This highlights the importance of _integrating your work_. We rightfully made the dependencies of our `PokerCLI` private but haven't made a way for users to construct it.

Is there a way to have caught this problem earlier?

### `package mypackage_test`

In all other examples so far when we make a test file we declare it as being in the same package that we are testing. 

This is fine and it means on the odd occasion where we want to test something internal to the package we have access to the unexported types. 

But given we have advocated for _not_ testing internal things _generally_, can Go help enforce that? What if we could test our code where we only have access to the exported types (like our `main` does)?

When you're writing a project with different packages I would strongly recommend that your test package name has `_test` at the end. When you do this you will only be able to have access to the public types in your package. This would help with this specific case but also helps enforce the discipline of only testing public APIs. 

An adage with TDD is that if you cannot test your code then it is probably poor code to work with. Using `package foo_test` will help with this.

Before fixing `main` let's change the package of our test inside `PokerCLI_test` to `poker_test`

If you have a well configured IDE you will suddenly see a lot of red! If you run the compiler you'll get the following errors

```
./PokerCLI_test.go:12:19: undefined: StubPlayerStore
./PokerCLI_test.go:14:11: undefined: PokerCLI
./PokerCLI_test.go:17:3: undefined: assertPlayerWin
./PokerCLI_test.go:22:19: undefined: StubPlayerStore
./PokerCLI_test.go:24:11: undefined: PokerCLI
./PokerCLI_test.go:27:3: undefined: assertPlayerWin
```

We have now stumbled into more questions on package design. 

#### `Do we want to have our stubs and helpers 'public' ?`

This is a subjective discussion. One could definitely argue that you do not want to pollute your package's API with code to facilitate tests.

In the presentation ["Advanced Testing with Go"](https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=53) by Mitchell Hashimoto it is described how at HashiCorp they advocate doing this so that users of the package can write tests without having to re-invent the wheel writing stubs. In our case this would mean anyone using our poker package wont have to create their own stub `PlayerStore` if they wish to work with our code.

Anecdotally I have used this technique in other shared packages and it has proved extremely useful in terms of users saving time when integrating with our packages.

So let's create a file called `testing.go` and add our stub and our helpers.

```go
package poker

import "testing"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], winner)
	}
}

// todo for you - the rest of the helpers
```

You'll need to make the helpers public (remember exporting is done with a capital letter at the start) if you want them to be exposed to importers of our package. 

In our CLI test you'll need to call the code as if you were using it within a different package.

```go
func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.PokerCLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.PokerCLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

}
```

You'll now see we have the same problems as we had in `main`

```
./PokerCLI_test.go:15:26: implicit assignment of unexported field 'playerStore' in poker.PokerCLI literal
./PokerCLI_test.go:15:39: implicit assignment of unexported field 'in' in poker.PokerCLI literal
./PokerCLI_test.go:25:26: implicit assignment of unexported field 'playerStore' in poker.PokerCLI literal
./PokerCLI_test.go:25:39: implicit assignment of unexported field 'in' in poker.PokerCLI literal
```

The easiest way to get around this is to make a constructor as we have for other types

```go
func NewPokerCLI(store PlayerStore, in io.Reader) *PokerCLI {
	return &PokerCLI{
		playerStore:store,
		in:in,
	}
}
```

Change the test to use the constructor instead and we should be back to the tests passing

Finally, we can go back to our new `main.go` and use the constructor we just made

```go
game := poker.NewPokerCLI(store, os.Stdin)
```
 
Try and run it, type "Bob wins".

### You cannot read "all" of os.Stdin

Nothing happens! You'll have to force the process to quit. What's going on? 

As an experiment change the code to the following

```go
func (cli *PokerCLI) PlayPoker() {
	log.Println("1")
	userInput, _ := ioutil.ReadAll(cli.in)
	log.Println("2")
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

No matter what you type, you never see `2` logged. The reason is the `ReadAll`, we cant read "all" of `os.Stdin`, as you can just keep typing stuff in! The `os.Stdin` is attached to our process and is a stream that wont finish.

We want to test that if we read _more_ than beyond the first newline that we fail.

```go
type failOnEndReader struct {
	t *testing.T
	rdr io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {

	n, err = m.rdr.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Read to the end when you shouldn't have")
	}

	return n, err
}
```

We've created a custom `io.Reader` wrapping around another and if we get to the end of the reader then we fail the test

We can now create a new test to try it out

```go
t.Run("do not read beyond the first newline", func(t *testing.T) {
    in := failOnEndReader{
        t,
        strings.NewReader("Chris wins\n hello there"),
    }

    playerStore := &poker.StubPlayerStore{}

    cli := poker.NewPokerCLI(playerStore, in)
    cli.PlayPoker()
})
```

It fails with

```
=== RUN   TestCLI/do_not_read_beyond_the_first_newline
    --- FAIL: TestCLI/do_not_read_beyond_the_first_newline (0.00s)
    	PokerCLI_test.go:56: Read to the end when you shouldn't have
```

To fix it, we cant use `io.ReadAll`. Instead we'll use a [`bufio.Reader`](https://golang.org/pkg/bufio/).

> Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer object, creating another object (Reader or Writer) that also implements the interface but provides buffering and some help for textual I/O. 

Update the code to the following

```go
type PokerCLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
}

func NewPokerCLI(store PlayerStore, in io.Reader) *PokerCLI {
	return &PokerCLI{
		playerStore: store,
		in:          bufio.NewReader(in),
	}
}

func (cli *PokerCLI) PlayPoker() {
	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins\n", "", 1)
}
```

The tests will now pass.

Now try to run the application in `main.go` again and it should work how we expect.

We will probably end up deleting this test in time as definitely _will_ want to read beyond the first line as we evaluate multiple commands from the user; but it was helpful to drive out a better solution and we didn't want to add new features while our application was not working properly.
