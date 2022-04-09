package main

import "testing"

func TestBadBank(t *testing.T) {
	var (
		transactions = []Transaction{
			{
				From: "Chris",
				To:   "Riya",
				Sum:  100,
			},
			{
				From: "Adil",
				To:   "Chris",
				Sum:  25,
			},
		}
		riya  = Account{Name: "Riya", Balance: 100}
		chris = Account{Name: "Chris", Balance: 75}
		adil  = Account{Name: "Adil", Balance: 200}
	)

	newBalanceFor := func(account Account) float64 {
		return NewBalanceFor(account, transactions).Balance
	}

	AssertEqual(t, newBalanceFor(riya), 200)
	AssertEqual(t, newBalanceFor(chris), 0)
	AssertEqual(t, newBalanceFor(adil), 175)
}
