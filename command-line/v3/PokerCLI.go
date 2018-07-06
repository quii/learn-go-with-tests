package poker

import (
	"bufio"
	"io"
	"strings"
)

// PokerCLI helps players through a game of poker
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

// PlayPoker starts the game
func (cli *PokerCLI) PlayPoker() {
	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins\n", "", 1)
}
