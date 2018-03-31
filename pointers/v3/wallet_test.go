package main

import (
	"testing"
)

func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		cases := []struct {
			name             string
			wallet           Wallet
			amountToWithdraw Bitcoin
			wantedBalance    Bitcoin
			wantedErr        bool
		}{
			{
				name:             "sufficient funds",
				wallet:           Wallet{Bitcoin(20)},
				amountToWithdraw: Bitcoin(10),
				wantedBalance:    Bitcoin(10),
				wantedErr:        false,
			},
			{
				name:             "insufficient funds",
				wallet:           Wallet{Bitcoin(20)},
				amountToWithdraw: Bitcoin(100),
				wantedBalance:    Bitcoin(20),
				wantedErr:        true,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.wallet.Withdraw(tt.amountToWithdraw)

				if tt.wallet.Balance() != tt.wantedBalance {
					t.Errorf("got balance %s want %s", tt.wallet.Balance(), tt.wantedBalance)
				}

				if tt.wantedErr && err == nil {
					t.Error("wanted an error but didnt get one")
				}

				if !tt.wantedErr && err != nil {
					t.Errorf("didnt want an error but got one %s", err)
				}
			})
		}
	})
}
