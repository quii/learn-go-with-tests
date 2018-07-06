package poker

import (
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// PokerCLI helps players through a game of poker
type PokerCLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func NewPokerCLI(store PlayerStore, in io.Reader) *PokerCLI {
	return &PokerCLI{
		playerStore: store,
		in:          in,
	}
}

// PlayPoker starts the game
func (cli *PokerCLI) PlayPoker() {
	log.Println("1")
	userInput, _ := ioutil.ReadAll(cli.in)
	log.Println("2")
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput []byte) string {
	return strings.Replace(string(userInput), " wins\n", "", 1)
}
