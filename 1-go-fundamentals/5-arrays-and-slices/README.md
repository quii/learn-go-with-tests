# Arrays and slices

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/arrays)**

**배열**(array)을 사용하면 특정한 순서로 같은 타입의 여러 요소들을 한 변수에 저장할 수 있습니다.

특히 배열을 다룰 때는 배열을 순환(iterate)하는 것이 매우 흔한 방법입니다. So let's
use [our new-found knowledge of `for`](iteration.md) to make a `Sum` function. `Sum` will
take an array of numbers and return the total.

Let's use our TDD skills

## Write the test first

Create a new folder to work in. Create a new file called `sum_test.go` and insert the following:

```go
package main

import "testing"

func TestSum(t *testing.T) {

	numbers := [5]int{1, 2, 3, 4, 5}

	got := Sum(numbers)
	want := 15

	if got != want {
		t.Errorf("got %d want %d given, %v", got, want, numbers)
	}
}
```
배열은 변수를 선언할 때 정의한 고정된 수용량을 가집니다. 

다음과 같은 두 가지 방법으로 배열을 초기화할 수 있습니다.

* \[N\]type{value1, value2, ..., valueN} e.g. `numbers := [5]int{1, 2, 3, 4, 5}`
* \[...\]type{value1, value2, ..., valueN} e.g. `numbers := [...]int{1, 2, 3, 4, 5}`

때때로 함수에 전달된 인풋을 에러메시지에서 프린트 해주는 것이 유용한데, 여기서 `%v`를 사용해 배열을 **"default"** 포맷으로 프린트해주었습니다.

[Read more about the format strings](https://golang.org/pkg/fmt/)

## Try to run the test

If you had initialized go mod with `go mod init main` you will be presented with an error
`_testmain.go:13:2: cannot import "main"`. This is because according to common practice,
package main will only contain integration of other packages and not unit-testable code and
hence Go will not allow you to import a package with name `main`.

To fix this, you can rename the main module in `go.mod` to any other name.

Once the above error is fixed, if you run `go test` the compiler will fail with the familiar
`./sum_test.go:10:15: undefined: Sum` error. Now we can proceed with writing the actual method
to be tested.

## Write the minimal amount of code for the test to run and check the failing test output

In `sum.go`

```go
package main

func Sum(numbers [5]int) int {
	return 0
}
```

Your test should now fail with _a clear error message_

`sum_test.go:13: got 0 want 15 given, [1 2 3 4 5]`

## Write enough code to make it pass

```go
func Sum(numbers [5]int) int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += numbers[i]
	}
	return sum
}
```

만약 배열의 특정 인덱스 값을 얻고 싶다면 `array[index]` 문법을 사용하면 됩니다. 위 예제에서는 `for`문을 사용해 5번 배열을 순환하고 각 아이템을 `sum` 변수에 더해줍니다.

## Refactor

Let's introduce [`range`](https://gobyexample.com/range) to help clean up our code

```go
func Sum(numbers [5]int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```

`range`를 사용하면 배열을 순환할 수 있다. 매 순환마다 `range`는 두 값을 반환합니다. 그 중 하나는 값의 인덱스를 반환하는데, 이 코드에서는 인덱스를 사용하지 않으므로 [blank identifier](https://golang.org/doc/effective_go.html#blank)인 `_`을 사용합니다.


### Arrays and their type

배열의 재밌는 점 중 하나는 바로 사이즈가 사이즈의 타입에 인코딩되어 있다는 점입니다. 만약 `[4]int`를 [`5]int`를 예상하는 함수에 전달하면 컴파일 되지 않습니다. 왜냐면 둘은 서로 다른 타입이고 마치 `int`를 원하는 함수에 `string` 타입을 전달하는 것과 같기 때문입니다.

You may be thinking it's quite cumbersome that arrays have a fixed length, and most
of the time you probably won't be using them!

Go는 컬렉션의 사이즈를 인코딩하지 않고 아무 사이즈나 가질 수 있게 하는 `slice` 기능을 제공합니다.

The next requirement will be to sum collections of varying sizes.

## Write the test first

We will now use the [slice type][slice] which allows us to have collections of
any size. The syntax is very similar to arrays, you just omit the size when
declaring them

`mySlice := []int{1,2,3}` rather than `myArray := [3]int{1,2,3}`

```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := [5]int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```

## Try and run the test

This does not compile

`./sum_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Sum`

## Write the minimal amount of code for the test to run and check the failing test output

The problem here is we can either

* Break the existing API by changing the argument to `Sum` to be a slice rather
  than an array. When we do this, we will potentially ruin
  someone's day because our _other_ test will no longer compile!
* Create a new function

In our case, no one else is using our function, so rather than having two functions to maintain, let's have just one.

```go
func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```

If you try to run the tests they will still not compile, you will have to change the first test to pass in a slice rather than an array.

## Write enough code to make it pass

It turns out that fixing the compiler problems were all we need to do here and the tests pass!

## Refactor

We already refactored `Sum` - all we did was replace arrays with slices, so no extra changes are required.
Remember that we must not neglect our test code in the refactoring stage - we can further improve our `Sum` tests.

```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```

항상 테스트가 유의미한지 질문해보는 것이 중요합니다. 왜냐면 테스트를 만들 수 있는 만큼 최대한 많이 만드는 것이 목표가 아니라, 코드 베이스를 최대한 신뢰할 수 있도록 하는 것이 목표이기 때문입니다.
너무 많은 테스트를 가지는 것은 그저 관리해야하는 것들을 더해줄 뿐입니다. 모든 테스트는 코스트가 있습니다.

지금 케이스에서도 함수를 테스트하는 테스트가 두개 있는 것은 반복적이라고 할 수 있습니다. 만약 다른 사이즈의 슬라이스가 작동한다면 다른 사이즈와도 작동한다고 추정할 수 있을 것입니다.

Go에는 [coverage tool](https://blog.golang.org/cover)이라는 내장된 테스팅 툴킷 기능이 있습니다.
100% 테스트 커버리지를 만드는 것이 최종 목표가 되면 안되겠지만, 커버리지 툴은 테스트로 커버되지 않는 코드를 식별할 수 있습니다.
만약 TDD에 매우 엄격한 사람이라면, 100% 커버리지를 유지하는 것이 좋습니다.

Try running

`go test -cover`

You should see

```bash
PASS
coverage: 100.0% of statements
```

Now delete one of the tests and check the coverage again.

Now that we are happy we have a well-tested function you should commit your
great work before taking on the next challenge.

We need a new function called `SumAll` which will take a varying number of
slices, returning a new slice containing the totals for each slice passed in.

For example

`SumAll([]int{1,2}, []int{0,9})` would return `[]int{3, 9}`

or

`SumAll([]int{1,1,1})` would return `[]int{3}`

## Write the test first

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## Try and run the test

`./sum_test.go:23:9: undefined: SumAll`

## Write the minimal amount of code for the test to run and check the failing test output

We need to define `SumAll` according to what our test wants.

Go 언어에서는 [_variadic functions_(가변 함수)](https://gobyexample.com/variadic-functions)라는 것을 통해여러 가변 인자를 작성할 수 있습니다.

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
	return
}
```

This is valid, but our tests still won't compile!

`./sum_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

Go에서는 슬라이스와 동등연산자를 함께 사용할 수 없습니다. 
이런 경우 각 `got`과 `want`를 순환하는 함수를 작성해서 값을 확인할 수 있지만, [`reflect.DeepEqual`][deepEqual]란 것을 사용해서 아무 두 변수가 동일한 값을 가지는지 확인할 수 있습니다.

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

\(make sure you `import reflect` in the top of your file to have access to `DeepEqual`\)

여기서 한가지 주의해야할 점은 `reflect.DeepEqual`은 타입을 보장하지 않는다는 점입니다. 즉, 이 코드는 어떤 타입이라도 컴파일하게 됩니다.

To see this in action, temporarily change the test to:

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := "bob"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

위 예제에서 `string`과 `slice`를 비교했습니다. 상식적으로는 말이 되지 않지만, 테스트는 위 코드를 성공적으로 컴파일합니다. 그러므로 `reflect.DeepEqual`은 매우 편리한 방법이지만 조심해서 사용해야합니다.

Change the test back again and run it. You should have test output like the following

`sum_test.go:30: got [] want [3 9]`

## Write enough code to make it pass

What we need to do is iterate over the varargs, calculate the sum using our
existing `Sum` function, then add it to the slice we will return

```go
func SumAll(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)
	sums := make([]int, lengthOfNumbers)

	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}

	return sums
}
```

Lots of new things to learn!

`make`는 `numbersToSu`m의 `len`만큼 슬라이스를 생성해줍니다.

You can index slices like arrays with `mySlice[N]` to get the value out or
assign it a new value with `=`

The tests should now pass.

## Refactor

위에서 언급했듯이, 슬라이스는 수용량을 가집니다. 예를 들어 둘의 수용량을 가진 슬라이가 있을때 `mySlice[10] = 1` 실행하면 *런타임 에러*가 발생하게 됩니다.

그러나 `append`를 사용하면 슬라이스를 취해 새로운 값을 만들고 새로운 슬라이스와 슬라이스 안에 모든 아이템을 반환할 수 있습니다.

```go
func SumAll(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}
```

위와 같이 구현하면 더 이상 수용량에 대해 걱정하지 않아도 됩니다. 위 코드는 먼저 빈 슬라인스인 `sums`로 시작하여 `Sum` 함수의 결과 값을 더해줄 수 있습니다.

Our next requirement is to change `SumAll` to `SumAllTails`, where it will
calculate the totals of the "tails" of each slice. The tail of a collection is
all items in the collection except the first one \(the "head"\).

## Write the test first

```go
func TestSumAllTails(t *testing.T) {
	got := SumAllTails([]int{1, 2}, []int{0, 9})
	want := []int{2, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## Try and run the test

`./sum_test.go:26:9: undefined: SumAllTails`

## Write the minimal amount of code for the test to run and check the failing test output

Rename the function to `SumAllTails` and re-run the test

`sum_test.go:30: got [3 9] want [2 9]`

## Write enough code to make it pass

```go
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		tail := numbers[1:]
		sums = append(sums, Sum(tail))
	}

	return sums
}
```

위와 같이 슬라이스는 슬라이싱 될 수 있습니다. `slice[low:high]`와 같은 문법을 사용합니다. If you omit the value on
one of the sides of the `:` it captures everything to that side of it. In our
case, we are saying "take from 1 to the end" with `numbers[1:]`. You may wish to
spend some time writing other tests around slices and experiment with the
slice operator to get more familiar with it.

## Refactor

Not a lot to refactor this time.

What do you think would happen if you passed in an empty slice into our
function? What is the "tail" of an empty slice? What happens when you tell Go to
capture all elements from `myEmptySlice[1:]`?

## Write the test first

```go
func TestSumAllTails(t *testing.T) {

	t.Run("make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

}
```

## Try and run the test

```text
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range
```

런타임 에러가 발생했지만 테스트 코드가 컴파일 되었다는 점을 인지하는 것은 중요합니다. 컴파일 타임 에러는 개발자가 작동하는 코드를 작성하도록 도와주지만, 런타임 에러는 유저에게 영향을 미치므로 피해야하는 부분입니다.

## Write enough code to make it pass

```go
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
}
```

## Refactor

Our tests have some repeated code around the assertions again, so let's extract those into a function.

```go
func TestSumAllTails(t *testing.T) {

	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("make the sums of tails of", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})

}
```

기존에 리팩토링 방식처럼 checkSums라는 함수를 생성했지만, 이 케이스에서는 함수를 변수에 할당하는 새로운 기술이 소개되었습니다. 약간 이상하게 보일 순 있지만
이 방식은 정수나 문자열을 변수에 할당하는 것과 다르지 않습니다. 왜냐면 함수도 결국 값이기 때문입니다.

여기선 더 자세히 설명이 안되었지만, 이런 방식은 함수를 다른 로컬 변수와 스코프 안에서 묶어줄 때 매우 유용한 방식이라고 할 수 있습니다(e.g between some `{}`). 또한 API의 면적을 줄여주기도 합니다.

이렇게 테스트 안에서 함수를 정의함으로써, 이 패키지 안에 다른 함수에서 해당 함수는 사용할 수 없게 됩니다. 외부에서 사용될 필요가 없는 변수나 함수를 숨겨주는 것은 소프트웨어 디자인적으로 매우 중요한 부분입니다.

A handy side-effect of this is this adds a little type-safety to our code. If
a developer mistakenly adds a new test with `checkSums(t, got, "dave")` the compiler
will stop them in their tracks.

```bash
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## Wrapping up

We have covered

* Arrays
* Slices
    * The various ways to make them
    * How they have a _fixed_ capacity but you can create new slices from old ones
      using `append`
    * How to slice, slices!
* `len` to get the length of an array or slice
* Test coverage tool
* `reflect.DeepEqual` and why it's useful but can reduce the type-safety of your code

We've used slices and arrays with integers but they work with any other type
too, including arrays/slices themselves. So you can declare a variable of
`[][]string` if you need to.

[Check out the Go blog post on slices][blog-slice] for an in-depth look into
slices. Try writing more tests to solidify what you learn from reading it.

Another handy way to experiment with Go other than writing tests is the Go
playground. You can try most things out and you can easily share your code if
you need to ask questions. [I have made a go playground with a slice in it for you to experiment with.](https://play.golang.org/p/ICCWcRGIO68)

[Here is an example](https://play.golang.org/p/bTrRmYfNYCp) of slicing an array
and how changing the slice affects the original array; but a "copy" of the slice
will not affect the original array.
[Another example](https://play.golang.org/p/Poth8JS28sc) of why it's a good idea
to make a copy of a slice after slicing a very large slice.

[for]: ../iteration.md#
[blog-slice]: https://blog.golang.org/go-slices-usage-and-internals
[deepEqual]: https://golang.org/pkg/reflect/#DeepEqual
[slice]: https://golang.org/doc/effective_go.html#slices
