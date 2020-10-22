package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/quii/learn-go-with-tests/time/v3"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	StartCalled     bool
	StartCalledWith int

	FinishedCalled   bool
	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalled = true
	g.StartCalledWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedCalled = true
	g.FinishCalledWith = winner
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("it prints an error when the winner is declared incorrectly", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("8", "Lloyd is a killer")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)
	})
}

func assertGameStartedWith(t testing.TB, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()
	if game.StartCalledWith != numberOfPlayersWanted {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.StartCalledWith)
	}
}

func assertGameNotFinished(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.FinishedCalled {
		t.Errorf("game should not have finished")
	}
}

func assertGameNotStarted(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()
	if game.FinishCalledWith != winner {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertScheduledAlert(t testing.TB, got, want poker.ScheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
