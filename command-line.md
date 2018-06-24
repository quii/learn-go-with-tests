# Command line (WIP)

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/master/command-line)**

Our product owner now wants to _pivot_ by introducing a second application. This will be a command line app which helps a group of people play Texas-Holdem Poker.

## Just enough information on poker. 

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
│   └── webserver
│       └── main.go
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
