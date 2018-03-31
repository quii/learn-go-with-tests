package main

import (
	"fmt"
)

type WithdrawError struct {
	AmountToWithdraw Bitcoin
	CurrentBalance   Bitcoin
}

func (w WithdrawError) Error() string {
	return fmt.Sprintf("cannot withdraw %s, insufficient funds - current balance is %s", w.AmountToWithdraw, w.CurrentBalance)
}

type Bitcoin int

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
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
			AmountToWithdraw: amount,
			CurrentBalance:   w.balance,
		}
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}
