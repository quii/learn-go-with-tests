# Command line (WIP)

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/master/command-line)**

Our product owner now wants to _pivot_ by introducing a second application. This will be a command line app which helps a group of people play Texas-Holdem Poker.

## Just enough information on poker

- N number of players sit in a circle.
- There is a dealer button, which gets passed to the left every round.
- To the left of the dealer there is the "small blind".
- To the left of the small blind there is the big blind (or just blind).
- These players *have* to contribute chips to the "pot". The big blind contributing twice as much as the small.
- This way, every player has to contribute money to the pot through the game, forcing people to play rather than "folding" all the time.
- The amount of chips the player has to contribute as a blind bet increases over time to ensure the game doesn't last too long.

Our application will help keep track of when the blind should go up, and how much it should be.

- Create a command line app.
- When it starts it asks how many players are playing. This determines the amount of time there is before the "blind" bet goes up.
  - There is a base amount of time of 5 minutes.
  - For every player, 1 minute is added.
  - e.g 6 players equals 11 minutes for the blind.
- After the blind time expires the game should alert the players the new amount the blind bet is.
- The blind starts at 100 chips, then 200, 400, 600, 1000, 2000 and continue to double until the game ends.
- When the game ends the user should be able to type "Chris wins" and that will record a win for the player in our existing database. This should then exit the program.

The product owner wants the database to be shared amongst the two applications so that the league updates according to wins recorded in the new application.

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
- Read everything from our `cli.in` to a string
- Extract out the winner by using `strings.Replace` which takes the string to replace, what substring to replace, its replacement and finally a flag to say how many instances to replace (`-1` means replace all).

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

Now that we have some passing tests, we should wire this up into `main`. Remember we should always strive to have fully-integrated working software as quickly as we can.

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

When you're writing a project with multiple packages I would strongly recommend that your test package name has `_test` at the end. When you do this you will only be able to have access to the public types in your package. This would help with this specific case but also helps enforce the discipline of only testing public APIs. If you still wish to test internals you can make a separate test with the package you want to test.

An adage with TDD is that if you cannot test your code then it is probably hard for users of your code to integrate with it. Using `package foo_test` will help with this by forcing you to test your code as if you are importing it like users of your package will.

Before fixing `main` let's change the package of our test inside `PokerCLI_test` to `poker_test`.

If you have a well configured IDE you will suddenly see a lot of red! If you run the compiler you'll get the following errors

```
./PokerCLI_test.go:12:19: undefined: StubPlayerStore
./PokerCLI_test.go:14:11: undefined: PokerCLI
./PokerCLI_test.go:17:3: undefined: assertPlayerWin
./PokerCLI_test.go:22:19: undefined: StubPlayerStore
./PokerCLI_test.go:24:11: undefined: PokerCLI
./PokerCLI_test.go:27:3: undefined: assertPlayerWin
```

We have now stumbled into more questions on package design. In order to test our software we made unexported stubs and helper functions which are no longer available for us to use in our `PokerCLI_test`.

#### `Do we want to have our stubs and helpers 'public' ?`

This is a subjective discussion. One could argue that you do not want to pollute your package's API with code to facilitate tests.

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

No matter what you type, you never see `2` logged. The reason is the `ReadAll`, we cant read "all" of `os.Stdin`, as you can just keep typing stuff in! The `os.Stdin` is attached to our process and is a stream that wont finish until the process finishes.

We want to test that if we read _more_ than beyond the first newline that we fail. This would mean a user could type `Bob wins` followed by a newline and it will be recorded, thus fixing the application.

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

### `time.AfterFunc`

We want to be able to schedule our program to print the blind bet values at certain durations dependant on the number of players.

To limit the scope of what we need to do, we'll forget about the number of players part and just assume there are 5 players so we'll test that _every 10 minutes the new value of the blind bet is printed_.

As usual the standard library has us covered with [`func AfterFunc(d Duration, f func()) *Timer`](https://golang.org/pkg/time/#AfterFunc)

> `AfterFunc` waits for the duration to elapse and then calls f in its own goroutine. It returns a `Timer` that can be used to cancel the call using its Stop method.

When we call `PlayPoker` we'll schedule all of our blind alerts.

Testing this may be a little tricky though. We'll want to verify that each time period is scheduled with the correct blind amount but if you look at the signature of `time.AfterFunc` its second argument is the function it will run. You cannot compare functions in Go so we'd be unable to test what function has been sent in. So we'll need to write some kind of wrapper around `time.AfterFunc` which will take the time to run and the amount to print so we can spy on that.

## Write the test first

Add a new test to our suite

```go
t.Run("it schedules printing of blind values", func(t *testing.T) {
    in := strings.NewReader("Chris wins\n")
    playerStore := &poker.StubPlayerStore{}
    blindAlerter := &SpyBlindAlerter{}

    cli := poker.NewPokerCLI(playerStore, in, blindAlerter)
    cli.PlayPoker()
    
    if len(blindAlerter.alerts) != 1 {
        t.Fatal("expected a blind alert to be scheduled")
    }
})
```

You'll notice we've made a `SpyBlindAlerter` which we are trying to inject into our `PokerCLI` and then checking that after we call `PlayerPoker` that an alert is scheduled.

(Remember we are just going for the simplest scenario first and then we'll iterate.)

Here's the definition of `SpyBlindAlerter`

```go
type SpyBlindAlerter struct {
	alerts []struct{
		scheduledAt time.Duration
		amount int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct {
		scheduledAt time.Duration
		amount int
	}{duration,  amount})
}

```


## Try to run the test

```
./PokerCLI_test.go:32:27: too many arguments in call to poker.NewPokerCLI
	have (*poker.StubPlayerStore, *strings.Reader, *SpyBlindAlerter)
	want (poker.PlayerStore, io.Reader)
```

## Write the minimal amount of code for the test to run and check the failing test output

We have added a new argument and the compiler is complaining. _Strictly speaking_ the minimal amount of code is to make `NewPokerCLI` accept a `*SpyBlindAlerter` but let's cheat a little and just define the dependency as an interface.

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}
```

And then add it to the constructor

```go
func NewPokerCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *PokerCLI
```

Your other tests will now fail as they dont have a `BlindAlerter` passed in to `NewPokerCLI`. 

Spying on BlindAlerter is not relevant for the other tests so in the test file so add

```go
var dummySpyAlerter = &SpyBlindAlerter{}
```

Then passed that into the other tests.

The tests should now compile and our new test fails

```
=== RUN   TestCLI
=== RUN   TestCLI/it_schedules_printing_of_blind_values
--- FAIL: TestCLI (0.00s)
    --- FAIL: TestCLI/it_schedules_printing_of_blind_values (0.00s)
    	PokerCLI_test.go:38: expected a blind alert to be scheduled
```

## Write enough code to make it pass

We'll need to add the `BlindAlerter` as a field on our `PokerCLI` so we can reference it in our `PlayPoker` method.

```go
type PokerCLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	alerter     BlindAlerter
}

func NewPokerCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *PokerCLI {
	return &PokerCLI{
		playerStore: store,
		in:          bufio.NewReader(in),
		alerter:     alerter,
	}
}
```

To make the test pass, we can call our `BlindAlerter` with anything we like

```go
func (cli *PokerCLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5 * time.Second, 100)
	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

Next we'll want to check it schedules all the alerts we'd hope for, for 5 players

## Write the test first

```go
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewPokerCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		cases := []struct{
			expectedScheduleTime time.Duration
			expectedAmount       int
		} {
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduleTime), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				amountGot := alert.amount
				if amountGot != c.expectedAmount {
					t.Errorf("got amount %d, want %d", amountGot, c.expectedAmount)
				}

				gotScheduledTime := alert.scheduledAt
				if gotScheduledTime != c.expectedScheduleTime {
					t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, c.expectedScheduleTime)
				}
			})
		}
	})
```

Table-based test works nicely here and clearly illustrate what our requirements are. We run through the table and check the `SpyBlindAlerter` to see if the alert has been scheduled with the correct values.

## Try to run the test

You should have a lot of failures looking like this

```go
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/it_schedules_printing_of_blind_values
    --- FAIL: TestCLI/it_schedules_printing_of_blind_values (0.00s)
=== RUN   TestCLI/it_schedules_printing_of_blind_values/100_scheduled_for_0s
        --- FAIL: TestCLI/it_schedules_printing_of_blind_values/100_scheduled_for_0s (0.00s)
        	PokerCLI_test.go:71: got scheduled time of 5s, want 0s
=== RUN   TestCLI/it_schedules_printing_of_blind_values/200_scheduled_for_10m0s
        --- FAIL: TestCLI/it_schedules_printing_of_blind_values/200_scheduled_for_10m0s (0.00s)
        	PokerCLI_test.go:59: alert 1 was not scheduled [{5000000000 100}]
```

## Write enough code to make it pass

```go
func (cli *PokerCLI) PlayPoker() {

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}

	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

It's not a lot more complicated than what we already had. We're just now iterating over an array of `blinds` and calling the scheduler on an increasing `blindTime`

## Refactor

We can encapsulate our scheduled alerts into a method just to make `PlayPoker` read a little clearer.

```go
func (cli *PokerCLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *PokerCLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}
```

Finally our tests are looking a little clunky. We have two anonymous structs representing the same thing, a `ScheduledAlert`. Let's refactor that into a new type and then make some helps to compare them.

```go
type scheduledAlert struct {
	at time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}
```

We've added a `String()` method to our type so it prints nicely if the test fails

Update our test to use our new type

```go
t.Run("it schedules printing of blind values", func(t *testing.T) {
    in := strings.NewReader("Chris wins\n")
    playerStore := &poker.StubPlayerStore{}
    blindAlerter := &SpyBlindAlerter{}

    cli := poker.NewPokerCLI(playerStore, in, blindAlerter)
    cli.PlayPoker()

    cases := []scheduledAlert {
        {0 * time.Second, 100},
        {10 * time.Minute, 200},
        {20 * time.Minute, 300},
        {30 * time.Minute, 400},
        {40 * time.Minute, 500},
        {50 * time.Minute, 600},
        {60 * time.Minute, 800},
        {70 * time.Minute, 1000},
        {80 * time.Minute, 2000},
        {90 * time.Minute, 4000},
        {100 * time.Minute, 8000},
    }

    for i, want := range cases {
        t.Run(fmt.Sprint(want), func(t *testing.T) {

            if len(blindAlerter.alerts) <= i {
                t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
            }

            got := blindAlerter.alerts[i]
            assertScheduledAlert(t, got, want)
        })
    }
})
```

Implement `assertScheduledAlert` yourself.

We've spent a fair amount of time here writing tests and have been somewhat naughty not integrating with our application. Let's address that before we pile on any more requirements.

Try running the app and it wont compile, complaining about not enough args to `NewPokerCLI`.

Let's create an implementation of `BlindAlerter` that we can use in our application.

Create `BlindAlerter.go` and move our `BlindAlerter` interface and add the new things below

```go
package poker

import (
	"time"
	"fmt"
	"os"
)

// BlindAlerter schedules alerts for blind amounts
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
```

Remember that any _type_ can implement an interface, not just `structs`. If you are making a library that exposes an interface with one function defined it is a common idiom to also expose a `MyInterfaceFunc` type. This type will be a `func` which will also implement your interface. That way users of your interface can convieniently implement your interface with just a function rather than having to create an empty `struct` type.

We then create the function `StdOutAlerter` which has the same signature as the function and just use `time.AfterFunc` to schedule it to print to `os.Stdout`.

Update `main` where we create `NewPokerCLI` to see this in action

```go
game := poker.NewPokerCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter))
```

Before running you might want to change the `blindTime` increment in `PokerCLI` to be 10 seconds rather than 10 minutes just so you can see it in action.

You should see it print the blind values as we'd expect every 10 seconds. Notice how you can still type "Shaun wins" into the CLI and it will stop the program how we'd expect.

The game wont always be played with 5 people so we need to prompt the user to enter a number of players before the game starts. 

## Write the test first

We'll want to record what is written to StdOut. We've done this a few times now, we know that `os.Stdout` is an `io.Writer` so we can check what is written if we use dependency injection to pass in a `bytes.Buffer` in our test and see what our code will write.

We don't care about our other collaborators in this test just yet so we've made some dummies in our test file. We should be a little wary that we now have 4 dependencies for `PokerCLI`, that feels like maybe it is starting to have too many responsiblities. Let's live with it for now and see if a refactoring emerges as we add this new functionality.

```go
var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}
```

And here is our new test

```go
t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
    stdout := &bytes.Buffer{}
    cli := poker.NewPokerCLI(dummyPlayerStore, dummyStdIn, stdout, dummyBlindAlerter)
    cli.PlayPoker()

    got :=stdout.String()
    want := "Please enter the number of players: "

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
})
```

We pass in what will be `os.Stdout` in `main` and see what is written.

## Try to run the test

```
./PokerCLI_test.go:38:27: too many arguments in call to poker.NewPokerCLI
	have (*poker.StubPlayerStore, *bytes.Buffer, *bytes.Buffer, *SpyBlindAlerter)
	want (poker.PlayerStore, io.Reader, poker.BlindAlerter)
```

## Write the minimal amount of code for the test to run and check the failing test output

We have a new dependency so we'll have to update `NewPokerCLI`

```go
func NewPokerCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *PokerCLI
```

Now the _other_ tests will fail to compile because they dont have an `io.Writer` being passed into `NewPokerCLI`. Add `dummyStdout` for the other tests.

The new test should fail like so

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players
    --- FAIL: TestCLI/it_prompts_the_user_to_enter_the_number_of_players (0.00s)
    	PokerCLI_test.go:46: got '', want 'Please enter the number of players: '
FAIL
```

## Write enough code to make it pass

We need to add our new dependency to our `PokerCLI` so we can reference it in `PlayPoker`

```go
type PokerCLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	out         io.Writer
	alerter     BlindAlerter
}

func NewPokerCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *PokerCLI {
	return &PokerCLI{
		playerStore: store,
		in:          bufio.NewReader(in),
		out:         out,
		alerter:     alerter,
	}
}
```

Then finally we can write our prompt at the start of the game

```go
func (cli *PokerCLI) PlayPoker() {
	fmt.Fprint(cli.out, "Please enter the number of players: ")
	cli.scheduleBlindAlerts()
	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

## Refactor

We have a duplicate string for the prompt which we should extract into a constant

```go
const PlayerPrompt = "Please enter the number of players: "
```

Use this in both the test code and `PokerCLI`.

Now we need to send in a number and extract it out. The only way we'll know if it has had the desired effect is by seeing what blind alerts were scheduled.

## Write the test first

```go
t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
    stdout := &bytes.Buffer{}
    in := strings.NewReader("7\n")
    blindAlerter := &SpyBlindAlerter{}

    cli := poker.NewPokerCLI(dummyPlayerStore, in, stdout, blindAlerter)
    cli.PlayPoker()

    got :=stdout.String()
    want := poker.PlayerPrompt

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }

    cases := []scheduledAlert{
        {0 * time.Second, 100},
        {12 * time.Minute, 200},
        {24 * time.Minute, 300},
        {36 * time.Minute, 400},
    }

    for i, want := range cases {
        t.Run(fmt.Sprint(want), func(t *testing.T) {

            if len(blindAlerter.alerts) <= i {
                t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
            }

            got := blindAlerter.alerts[i]
            assertScheduledAlert(t, got, want)
        })
    }
})
```

Ouch! A lot of changes. 

- We remove our dummy for StdIn and instead send in a mocked version representing our user entering 7
- We also remove our dummy on the blind alerter so we can see that the number of players has had an effect on the scheduling
- We test what alerts are scheduled

## Try to run the test

The test should still compile and fail reporting that the scheduled times are wrong because we've hard-coded for the game to be based on having 5 players

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players
    --- FAIL: TestCLI/it_prompts_the_user_to_enter_the_number_of_players (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players/100_chips_at_0s
        --- PASS: TestCLI/it_prompts_the_user_to_enter_the_number_of_players/100_chips_at_0s (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players/200_chips_at_12m0s
```

## Write enough code to make it pass

Remember, we are free to commit whatever sins we need to make this work. Once we have working software we can then work on refactoring the mess we're about to make!


```go
func (cli *PokerCLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	
	numberOfPlayersInput, _ := cli.in.ReadString('\n')
	numberOfPlayers, _ := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	cli.scheduleBlindAlerts(numberOfPlayers)

	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *PokerCLI) scheduleBlindAlerts(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}
```

- We read in the `numberOfPlayersInput` into a string
- We `Trim` the string of the newline entered and use `Atoi` to convert it into an integer - ignoring any error scenarios. We'll need to write a test for that scenario later.
- From here we change `scheduleBlindAlerts` to accept a number of players. We then calculate a `blindIncrement` time to use to add to `blindTime` as we iterate over the blind amounts

While our new test has been fixed, a lot of others have failed because now our system only works if the game starts with a user entering a number. You'll need to fix the tests by changing the user inputs so that a number followed by a newline is added (this is highlighting yet more flaws in our approach right now).

## Refactor

Let's _listen to our tests_. 

- In order to test that we are scheduling some alerts we set up 4 different dependencies. Whenever you have a lot of dependencies for a _thing_ in your system, it implies it's doing too much. Visually we can see it in how cluttered our test is.
- To me it feels like **we need to make a cleaner abstraction between reading user input and the business logic we want to do** 
- A better test would be _given this user input, do we call the `BlindAlerter` with the correct number of players_. 
- We would then extract the testing of the scheduling into the tests for our new `BlindAlerter`.

We can refactor our `BlindAlerter` first and our test should continue to pass. Once we've made the structural changes we want we can think about how we can refactor the tests to reflect our new separation of concerns

