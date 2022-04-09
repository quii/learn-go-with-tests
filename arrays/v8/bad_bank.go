package main

type Transaction struct {
	From string
	To   string
	Sum  float64
}

type Balances map[string]float64

func CalculateBalances(transactions []Transaction) Balances {
	adjustAccountsByTransaction := func(b Balances, t Transaction) Balances {
		b[t.From] -= t.Sum
		b[t.To] += t.Sum
		return b
	}
	return Reduce(transactions, make(Balances), adjustAccountsByTransaction)
}
