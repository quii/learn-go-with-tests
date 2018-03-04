package main

func Repeat(character string) (repeated string) {
	for i := 0; i < 5; i++ {
		repeated = repeated + character
	}
	return
}
