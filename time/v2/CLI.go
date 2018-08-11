package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Game manages the state of a game
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

// CLI helps players through a game of poker
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	out         io.Writer
	game        Game
}

// NewCLI creates a CLI for playing poker
func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

// PlayerPrompt is the text asking the user for the number of players
const PlayerPrompt = "Please enter the number of players: "

// ErrorPlayerNumberPrompt tells the user they entered in the value wrong
const ErrorPlayerNumberPrompt = "ERROR: Please enter the number of players as a number: "

// PlayPoker starts the game
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, _ := strconv.Atoi(cli.readLine())

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
