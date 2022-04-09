package main

func Find[A any](items []A, predicate func(A) bool) (value A, found bool) {
	for _, v := range items {
		if predicate(v) {
			return v, true
		}
	}
	return
}

func Reduce[A, B any](collection []A, accumulator B, f func(B, A) B) B {
	for _, x := range collection {
		accumulator = f(accumulator, x)
	}
	return accumulator
}
