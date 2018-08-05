package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CLI helps players through a game of poker
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	out         io.Writer
	game        *Game
}

// NewCLI creates a CLI for playing poker
func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		in:  bufio.NewReader(in),
		out: out,
		game: &Game{
			alerter: alerter,
			store:   store,
		},
	}
}

// PlayerPrompt is the text asking the user for the number of players
const PlayerPrompt = "Please enter the number of players: "

// PlayPoker starts the game
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput, _ := cli.in.ReadString('\n')
	numberOfPlayers, _ := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	cli.game.Start(numberOfPlayers)

	winnerInput, _ := cli.in.ReadString('\n')
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins\n", "", 1)
}
