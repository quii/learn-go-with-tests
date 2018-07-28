package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// PokerCLI helps players through a game of poker
type PokerCLI struct {
	playerStore PlayerStore
	in          *bufio.Reader
	alerter     BlindAlerter
}

func NewPokerCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *PokerCLI {
	return &PokerCLI{
		playerStore: store,
		in:          bufio.NewReader(in),
		alerter:     alerter,
	}
}

// PlayPoker starts the game
func (cli *PokerCLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput, _ := cli.in.ReadString('\n')
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *PokerCLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins\n", "", 1)
}
