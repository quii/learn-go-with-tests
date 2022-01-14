package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRecordWin(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	playerName := "Ron"
	mock.ExpectExec(`
insert into playerWins (player, wins) values ($1, 1)
on conflict (player)
do update set wins = playerWins.wins + 1 where playerWins.player = $2 ;`,
	).WithArgs(playerName, playerName).WillReturnResult(sqlmock.NewResult(1, 1))

	storage := PostgresPlayerStore{db}
	storage.RecordWin(playerName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectation: %s", err)
	}
}
