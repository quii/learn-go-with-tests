package main

import (
	"database/sql"
)

// PostgresPlayerStore collects data about players in PostgreSQL DB.
type PostgresPlayerStore struct {
	db *sql.DB
}

// GetPlayerScore retrieves scores for a given player.
func (p *PostgresPlayerStore) GetPlayerScore(_ string) int {
	panic("implement me")
}

// RecordWin will record a player's win.
func (p *PostgresPlayerStore) RecordWin(name string) {
	exec, _ := p.db.Exec(`
insert into playerWins (player, wins) values ($1, 1)
on conflict (player)
do update set
  wins = playerWins.wins + 1
where playerWins.player = $2 ;`, name, name)
	_, _ = exec.RowsAffected()
}
