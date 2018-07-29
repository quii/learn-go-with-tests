package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// PokerCLI helps players through a game of poker
type PokerCLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	out         io.Writer
	alerter     BlindAlerter
}

func NewPokerCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *PokerCLI {
	return &PokerCLI{
		playerStore: store,
		in:          bufio.NewReader(in),
		out:         out,
		alerter:     alerter,
	}
}

const PlayerPrompt = "Please enter the number of players: "

// PlayPoker starts the game
func (cli *PokerCLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput, _ := cli.in.ReadString('\n')
	numberOfPlayers, _ := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	cli.scheduleBlindAlerts(numberOfPlayers)

	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *PokerCLI) scheduleBlindAlerts(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins\n", "", 1)
}
