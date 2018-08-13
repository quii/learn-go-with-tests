package poker

import (
	"bufio"
	"io"
	"strings"
)

// CLI helps players through a game of poker
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
}

// NewCLI creates a CLI for playing poker
func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewReader(in),
	}
}

// PlayPoker starts the game
func (cli *CLI) PlayPoker() {
	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins\n", "", 1)
}
