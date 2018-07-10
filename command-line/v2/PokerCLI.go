package poker

// PokerCLI helps players through a game of poker
type PokerCLI struct {
	playerStore PlayerStore
}

// PlayPoker starts the game
func (cli *PokerCLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
