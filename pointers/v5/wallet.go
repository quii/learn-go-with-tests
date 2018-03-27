package main

import (
	"fmt"
)

type Bitcoin int

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

type WithdrawError struct {
	CurrentBalance   Bitcoin
	AmountToWithdraw Bitcoin
}

func (w WithdrawError) Error() string {
	return fmt.Sprintf("cannot withdraw %s, insufficient funds (%s)", w.AmountToWithdraw, w.CurrentBalance)
}

type Wallet struct {
	balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return WithdrawError{
			CurrentBalance:   w.balance,
			AmountToWithdraw: amount,
		}
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}
