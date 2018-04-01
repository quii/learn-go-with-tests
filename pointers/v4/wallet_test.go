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
			wantedErr        error
		}{
			{
				name:             "sufficient funds",
				wallet:           Wallet{Bitcoin(20)},
				amountToWithdraw: Bitcoin(10),
				wantedBalance:    Bitcoin(10),
				wantedErr:        nil,
			},
			{
				name:             "insufficient funds",
				wallet:           Wallet{Bitcoin(20)},
				amountToWithdraw: Bitcoin(100),
				wantedBalance:    Bitcoin(20),
				wantedErr:        InsufficientFundsError,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.wallet.Withdraw(tt.amountToWithdraw)

				assertBalance(t, tt.wallet, tt.wantedBalance)

				gotAnError := err != nil
				wantAnError := tt.wantedErr != nil

				if gotAnError != wantAnError {
					t.Fatalf("got error '%s' want '%s'", err, tt.wantedErr)
				}

				if wantAnError && err != tt.wantedErr {
					t.Errorf("got err '%s' want '%s'", err.Error(), tt.wantedErr)
				}
			})
		}
	})
}
