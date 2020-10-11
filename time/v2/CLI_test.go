package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "github.com/quii/learn-go-with-tests/time/v2"
)

var dummyBlindAlerter = &poker.SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	StartCalledWith  int
	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalledWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishCalledWith = winner
}

func TestCLI(t *testing.T) {

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartCalledWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartCalledWith)
		}
	})

	t.Run("finish game with 'Chris' as winner", func(t *testing.T) {
		in := strings.NewReader("1\nChris wins\n")
		game := &GameSpy{}
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		if game.FinishCalledWith != "Chris" {
			t.Errorf("expected finish called with 'Chris' but got %q", game.FinishCalledWith)
		}
	})

	t.Run("record 'Cleo' win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nCleo wins\n")
		game := &GameSpy{}
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		if game.FinishCalledWith != "Cleo" {
			t.Errorf("expected finish called with 'Cleo' but got %q", game.FinishCalledWith)
		}
	})

	t.Run("do not read beyond the first newline", func(t *testing.T) {
		in := failOnEndReader{
			t,
			strings.NewReader("1\nChris wins\n hello there"),
		}

		game := poker.NewTexasHoldem(dummyBlindAlerter, dummyPlayerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()
	})

}

type failOnEndReader struct {
	t   *testing.T
	rdr io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {

	n, err = m.rdr.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Read to the end when you shouldn't have")
	}

	return n, err
}

func assertScheduledAlert(t testing.TB, got, want poker.ScheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
