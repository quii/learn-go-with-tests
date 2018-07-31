package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// PokerGame manages a game of poker
type PokerGame struct {
	alerter BlindAlerter
	store   PlayerStore
}

func (p *PokerGame) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (p *PokerGame) Finish(winner string) {
	p.store.RecordWin(winner)
}

// PokerCLI helps players through a game of poker
type PokerCLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	out         io.Writer
	game        *PokerGame
}

func NewPokerCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *PokerCLI {
	return &PokerCLI{
		in:  bufio.NewReader(in),
		out: out,
		game: &PokerGame{
			alerter: alerter,
			store:   store,
		},
	}
}

const PlayerPrompt = "Please enter the number of players: "

// PlayPoker starts the game
func (cli *PokerCLI) PlayPoker() {
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
