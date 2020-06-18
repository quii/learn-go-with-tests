package poker

// Game manages the state of a game.
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}
