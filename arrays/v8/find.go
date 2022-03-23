package main

func Find[A any](items []A, predicate func(A) bool) (A, bool) {
	var item A
	for _, v := range items {
		if predicate(v) {
			return v, true
		}
	}
	return item, false
}
