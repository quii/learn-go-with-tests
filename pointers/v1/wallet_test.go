package main

import (
	"testing"
)

func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()

	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}