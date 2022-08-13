# Iteration

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/for)**

Go에서 반복시키고 싶은 것이 있다면 `for`을 사용하면 된다. Go는 `while`, `do`, `until`과 같은 키워드는 가지지 않고 오직 `for`만 사용한다!

먼저 문자를 다섯번 반복하는 함수의 테스트 코드를 작성해보자.

There's nothing new so far, so try and write it yourself for practice.

## Write the test first

```go
package iteration

import "testing"

func TestRepeat(t *testing.T) {
	repeated := Repeat("a")
	expected := "aaaaa"

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}
```

## Try and run the test

`./repeat_test.go:6:14: undefined: Repeat`

## Write the minimal amount of code for the test to run and check the failing test output

_Keep the discipline!_ You don't need to know anything new right now to make the test fail properly.

All you need to do right now is enough to make it compile so you can check your test is written well.

```go
package iteration

func Repeat(character string) string {
	return ""
}
```

Isn't it nice to know you already know enough Go to write tests for some basic problems? This means you can now play with the production code as much as you like and know it's behaving as you'd hope.

`repeat_test.go:10: expected 'aaaaa' but got ''`

## Write enough code to make it pass

Go의 `for` 문법은 C언어로부터 파생된 언어들과 매우 비슷합니다.

```go
func Repeat(character string) string {
	var repeated string
	for i := 0; i < 5; i++ {
		repeated = repeated + character
	}
	return repeated
}
```

C, Java, 혹은 JavaScript와 같은 언어들과  달리 Go에서는 for문의 각 요소를 감싸는 소괄호(`()`)가 필요하지 않고, 중괄호(`{}`)가 항상 요구됩니다. 


```go
var repeated string
```

지금까지 변수를 선언하고 초기화하기위해 `:=`를 사용해왔습니다. 하지만 이 방식은 단순히 빨리 적기([속기](https://gobyexample.com/variables)) 위함입니다. 여기서는 `string` 변수만을 선언해줬습니다. 더 명백히는 ₩를 함수를 선언할 때도 사용히 가능합니다.

Run the test and it should pass.

Additional variants of the for loop are described [here](https://gobyexample.com/for).

## Refactor

Now it's time to refactor and introduce another construct `+=` assignment operator.

```go
const repeatCount = 5

func Repeat(character string) string {
	var repeated string
	for i := 0; i < repeatCount; i++ {
		repeated += character
	}
	return repeated
}
```
`+=`는 특정 값을 동시에 **더하고 할당**해주는 할당 연산자입니다. 오른쪽 피연산자를 왼쪽 피연산자와 더해주고 그 결과값을 왼쪽 피연산자로 할당해줍니다. 문자열뿐만 아니라 숫자형 타입과도 사용이 가능합니다.

### Benchmarking

Writing [benchmarks](https://golang.org/pkg/testing/#hdr-Benchmarks) in Go is another first-class feature of the language and it is very similar to writing tests.

```go
func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a")
	}
}
```

You'll see the code is very similar to a test.

`testing.B`를 통해 숨겨진 `b.N`이라는 이름에 접근할 수 있게 됩니다.

벤치마크 코드가 실행되면 b,N 횟수만큼 코드가 실행되고 얼마나 코드의 실행속도가 측정됩니다.

The amount of times the code is run shouldn't matter to you, the framework will determine what is a "good" value for that to let you have some decent results.

To run the benchmarks do `go test -bench=.` (or if you're in Windows Powershell `go test -bench="."`)

```text
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000           136 ns/op
PASS
```

What `136 ns/op` means is our function takes on average 136 nanoseconds to run \(on my computer\). Which is pretty ok! To test this it ran it 10000000 times.

_NOTE_ by default Benchmarks are run sequentially.

## Practice exercises

* Change the test so a caller can specify how many times the character is repeated and then fix the code
* Write `ExampleRepeat` to document your function
* Have a look through the [strings](https://golang.org/pkg/strings) package. Find functions you think could be useful and experiment with them by writing tests like we have here. Investing time learning the standard library will really pay off over time.

## Wrapping up

* More TDD practice
* Learned `for`
* Learned how to write benchmarks
