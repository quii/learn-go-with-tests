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

	newBalanceFor := func(account Account) Account {
		return NewBalanceFor(transactions, account)
	}

	AssertEqual(t, newBalanceFor(riya).Balance, 200)
	AssertEqual(t, newBalanceFor(chris).Balance, 0)
	AssertEqual(t, newBalanceFor(adil).Balance, 175)
}
