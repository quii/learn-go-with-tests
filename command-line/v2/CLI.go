package poker

// CLI helps players through a game of poker
type CLI struct {
	playerStore PlayerStore
}

// PlayPoker starts the game
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
