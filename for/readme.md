# Iteration

To do stuff repeatedly in Go, you'll need `for`. In go there are no `while`, `do`, `until` kewords, you can only use `for`. Which is a good thing!

## Repeat a character

Let's write a test for a function that repeats a character 5 times.

There's nothing new so far, so try and write it yourself for practice.

```go
package main

import "testing"

func TestRepeat(t *testing.T) {
	repeated := Repeat("a")
	expected := "aaaaa"
	
	if repeated != expected {
		t.Errorf("expected '%s' but got '%s'", expected, repeated)
	}
}
```

_Keep the discipline!_ You don't need to know anything new right now

- Run the test
- Listen to the compiler
- Make it compile

```go
package main

func Repeat(character string) (repeated string)  {
	return
}
```

Isn't it nice to know you already know enough Go to write tests for some basic problems? This means you can now play with the production code as much as you like and know it's behaving as you'd hope.

The `for` syntax is very unremarkable and follows most C-like languages.

```go
func Repeat(character string) (repeated string) {
	for i := 0; i < 5; i++ {
		repeated = repeated + character
	}
	return
}
```

Run the test and it should pass. 

Now it's time to refactor and introduce another construct `+=`

```go
const repeatCount = 5

func Repeat(character string) (repeated string) {
	for i := 0; i < repeatCount; i++ {
		repeated += character
	}
	return
}
```

`+=` adds a value to another. It works other types like integers.

## Next iteration

todo: Think of a way to include `continue` in exercise.

## Benchmarking

Writing benchmarks in Go is another first-class feature of the language and it is very similar to writing tests. 

```go
func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a")
	}
}
```

You'll see the code is very similar to a test.

The `testing.B` gives you access to the cryptically named `b.N`. 

When the benchmark is run the code is ran `b.N` times, and measures how long it takes. 

The amount of times the code is ran shouldnt matter to you, the framework will determine what is a "good" value for that to let you have some decent results.

To run the benchmarks do `go test -bench=.`

```
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000	       136 ns/op
PASS
```

What that means is our function takes 136 nanoseconds to run (on my computer). Which is pretty ok!

### Practice exercises

- Change the test so a caller can specify how many times the character is repeated and then fix the code
- Write `ExampleRepeat` to document your function
- Have a look through the [the strings package](https://golang.org/pkg/strings)  package. Find functions you think could be useful and experiment with them by writing tests like we have here. Investing time learning the standard library will really pay off over time.