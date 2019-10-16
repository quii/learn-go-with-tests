package main

import (
	"fmt"
)

func main() {
	store, removeStore := NewStore()
	fmt.Printf("Connected to the store!, %+v\n", store)
	defer removeStore()
}
