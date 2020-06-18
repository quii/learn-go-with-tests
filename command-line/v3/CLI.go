package poker

import (
	"bufio"
	"io"
	"strings"
)

// CLI helps players through a game of poker.
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

// NewCLI creates a CLI for playing poker.
func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}

// PlayPoker starts the game.
func (cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
