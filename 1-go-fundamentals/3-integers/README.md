# Integers

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/integers)**

Integers work as you would expect. Let's write an `Add` function to try things out. Create a test file called `adder_test.go` and write this code.

**Note:** Go source files can only have one `package` per directory, make sure that your files are organised separately. [Here is a good explanation on this.](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)

## Write the test first

```go
package integers

import "testing"

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}
```

전에 작성한 테스트 코드와는 달리 포맷 문자열로 %q 대신 %d가 사용되었습니다. 이는 문자열을 프린트하지 않고 정수를 프린트하기 위함입니다.

또한 더 이상 메인 패키지도 사용하지 않고 `integers` 패키지를 사용하는데, 이름에서도 알 수 있듯이 `Add`와 같이 정수를 다루는 함수들을 그루핑하기 위함입니다.


## Try and run the test

`go test`로 테스트를 실행하면 다음과 같은 컴파일 에러를 마주하게 됩니다.


`./adder_test.go:6:9: undefined: Add`

## Write the minimal amount of code for the test to run and check the failing test output

컴파일러를 만족시키기 위해 충분한 코드를 작성합니다 - 테스트가 정확한 이유로 실패하는 것을 확인하기 위함임을 기억해야합니다.

```go
package integers

func Add(x, y int) int {
	return 0
}
```

만약 하나 이상의 같은 타입을 가진 인자가 존재할 때는 `(x int, y int)`가 아닌 `(x, y int)`와 같은 형식으로 줄여 작성할 수 있습니다.

Now run the tests and we should be happy that the test is correctly reporting what is wrong.

`adder_test.go:10: expected '4' but got '0'`

If you have noticed we learnt about _named return value_ in the [last](hello-world.md#one...last...refactor?) section but aren't using the same here. It should generally be used when the meaning of the result isn't clear from context, in our case it's pretty much clear that `Add` function will add the parameters. You can refer [this](https://github.com/golang/go/wiki/CodeReviewComments#named-result-parameters) wiki for more details.

## Write enough code to make it pass

In the strictest sense of TDD we should now write the _minimal amount of code to make the test pass_. A pedantic programmer may do this

```go
func Add(x, y int) int {
	return 4
}
```

Ah hah! Foiled again, TDD is a sham right?

We could write another test, with some different numbers to force that test to fail but that feels like [a game of cat and mouse](https://en.m.wikipedia.org/wiki/Cat_and_mouse).

Once we're more familiar with Go's syntax I will introduce a technique called *"Property Based Testing"*, which would stop annoying developers and help you find bugs.

For now, let's fix it properly

```go
func Add(x, y int) int {
	return x + y
}
```

If you re-run the tests they should pass.

## Refactor

There's not a lot in the _actual_ code we can really improve on here.

We explored earlier how by naming the return argument it appears in the documentation but also in most developer's text editors.

This is great because it aids the usability of code you are writing. It is preferable that a user can understand the usage of your code by just looking at the type signature and documentation.

주석을 닮으로써 함수를 문서화 할 수 있고, 이는 일반 라이브러리 문서처럼 Go Doc에 노출됩니다. 

```go
// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
	return x + y
}
```

### Examples

만약 정말로 테스트를 더 보강하고 싶으면 [examples](https://blog.golang.org/examples) 을 만들 수 있습니다. 실제 일반 라이브러리에서 `example`이 많이 사용되는 것을 확인할 수 있습니다.

보통 코드 예제는 코드 베이스 밖에 리드미와 같은 파일에서 찾을 수 있고, 체크 되지 않기 때문에 실제 코드와 비교해서 오래되었거나 잘못된 경우가 많습니다.

`Go example`은 테스트와 같이 실행되므로 코드가 실제로 수행하는 작업의 결과를 확인하고, 그렇기에 작동을 신뢰할 수 있게 도와줍니다.

`Example`은 패키지의 테스트 수이트의 부분으로 컴파일 됩니다(혹은 옵셔널하게 실행됩니다).

As with typical tests, examples are functions that reside in a package's `_test.go` files. Add the following `ExampleAdd` function to the `adder_test.go` file.

```go
func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
```

(If your editor doesn't automatically import packages for you, the compilation step will fail because you will be missing `import "fmt"` in `adder_test.go`. It is strongly recommended you research how to have these kind of errors fixed for you automatically in whatever editor you are using.)

만약 코드가 변경되면 `example`은 더 이상 유효하지 않기 때문에 테스트는 실패하게 됩니다.

Running the package's test suite, we can see the example function is executed with no further arrangement from us:

```bash
$ go test -v
=== RUN   TestAdder
--- PASS: TestAdder (0.00s)
=== RUN   ExampleAdd
--- PASS: ExampleAdd (0.00s)
```

여기서 중요한 점은 `// Output: 6`이라는 주석을 삭제하면 `example` 함수는 실행되지 않는다는 점입니다. 더 정확히 함수는 컴파일 되지만 실행되지 않는다고 할 수 있습니다.

By adding this code the example will appear in the documentation inside `godoc`, making your code even more accessible.

To try this out, run `godoc -http=:6060` and navigate to `http://localhost:6060/pkg/`

Inside here you'll see a list of all the packages and you'll be able to find your example documentation.

If you publish your code with examples to a public URL, you can share the documentation of your code at [pkg.go.dev](https://pkg.go.dev/). For example, [here](https://pkg.go.dev/github.com/quii/learn-go-with-tests/integers/v2) is the finalised API for this chapter. This web interface allows you to search for documentation of standard library packages and third-party packages.

## Wrapping up

What we have covered:

* More practice of the TDD workflow
* Integers, addition
* Writing better documentation so users of our code can understand its usage quickly
* Examples of how to use our code, which are checked as part of our tests
