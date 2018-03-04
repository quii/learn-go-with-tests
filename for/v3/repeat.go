package main

const repeatCount = 5

func Repeat(character string) (repeated string) {
	for i := 0; i < repeatCount; i++ {
		repeated += character
	}
	return
}
