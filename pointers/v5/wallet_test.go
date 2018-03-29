package main

import (
	"testing"
)

func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {

		cases := []struct {
			description      string
			wallet           Wallet
			amountToWithdraw Bitcoin
			wantedBalance    Bitcoin
			wantedErr        *WithdrawError
		}{
			{
				description:      "happy withdraw",
				wallet:           Wallet{balance: Bitcoin(10)},
				amountToWithdraw: Bitcoin(5),
				wantedBalance:    Bitcoin(5),
				wantedErr:        nil,
			},
			{
				description:      "not enough funds",
				wallet:           Wallet{balance: Bitcoin(10)},
				amountToWithdraw: Bitcoin(20),
				wantedBalance:    Bitcoin(10),
				wantedErr:        &WithdrawError{AmountToWithdraw: Bitcoin(20), CurrentBalance: Bitcoin(10)},
			},
		}

		for _, tt := range cases {
			t.Run(tt.description, func(t *testing.T) {
				err := tt.wallet.Withdraw(tt.amountToWithdraw)

				assertBalance(t, tt.wallet, tt.wantedBalance)

				if tt.wantedErr != nil {
					assertWithdrawError(t, err, *tt.wantedErr)
				}
			})
		}
	})

}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertWithdrawError(t *testing.T, err error, want WithdrawError) {
	got, isWithdrawErr := err.(WithdrawError)

	if !isWithdrawErr {
		t.Fatalf("did not get a withdraw error %#v", err)
	}

	if want != got {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
