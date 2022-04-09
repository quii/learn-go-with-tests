package main

type Transaction struct {
	From string
	To   string
	Sum  float64
}

func BalanceFor(transactions []Transaction, name string) float64 {
	adjustBalance := func(acc float64, t Transaction) float64 {
		if t.From == name {
			return acc - t.Sum
		}
		if t.To == name {
			return acc + t.Sum
		}
		return acc
	}
	return Reduce(transactions, 0.0, adjustBalance)
}
