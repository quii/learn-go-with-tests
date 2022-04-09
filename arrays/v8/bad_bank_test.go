package main

import "testing"

func TestBadBank(t *testing.T) {
	transactions := []Transaction{
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
	balances := CalculateBalances(transactions)
	AssertEqual(t, balances["Riya"], 100)
	AssertEqual(t, balances["Chris"], -75)
	AssertEqual(t, balances["Adil"], -25)
}
