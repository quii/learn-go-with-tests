# Hello, World

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/hello-world)**

It is traditional for your first program in a new language to be [Hello, World](https://en.m.wikipedia.org/wiki/%22Hello,_World!%22_program).

- Create a folder wherever you like
- Put a new file in it called `hello.go` and put the following code inside it

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
```

 `go run hello.go` 명령어로 실행합니다.

## How it works

Go로 프로그램을 작성할 때 보통 `main` 함수로 정의된 `main` 패키지를 가지게 됩니다. 패키지란 관련된 Go 코드를 함께 그루핑하기 위한 방식입니다.

`func` 키워드는 함수를 이름과 본문으로 정의합니다.

`import "fmt"`로 출력할 때 사용하는 `PrintIn` 함수를 담고 있는 패키지를 임포트할 수 있습니다.


## How to test

테스트를 하기전에, **도메인 코드**를 외부로 부터(사이드 이펙트)로 분리하는 것이 좋은 방법입니다. 여기서 fmt.PrintIn이 사이드 이펙트이고 문자열을 전달하는 곳이 도메인입니다.

So let's separate these concerns so it's easier to test

```go
package main

import "fmt"

func Hello() string {
	return "Hello, world"
}

func main() {
	fmt.Println(Hello())
}
```
위와 같이 `func` 키워드로 새로운 함수를 생성하였지만, 이번에는 함수를 정의할 때 `string`이라는 다른 키워드가 추가되었습니다. 이는 이 함수가 `string`을 반환한다는 뜻입니다.



Now create a new file called `hello_test.go` where we are going to write a test for our `Hello` function

```go
package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	want := "Hello, world"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

## Go modules?

The next step is to run the tests. Enter `go test` in your terminal. If the tests pass, then you are probably using an earlier version of Go. However, if you are using Go 1.16 or later, then the tests will likely not run at all. Instead, you will see an error message like this in the terminal:

```shell
$ go test
go: cannot find main module; see 'go help modules'
```

What's the problem? In a word, [modules](https://blog.golang.org/go116-module-changes). Luckily, the problem is easy to fix. Enter `go mod init hello` in your terminal. That will create a new file with the following contents:

```
module hello

go 1.16
```

This file tells the `go` tools essential information about your code. If you planned to distribute your application, you would include where the code was available for download as well as information about dependencies. For now, your module file is minimal, and you can leave it that way. To read more about modules, [you can check out the reference in the Golang documentation](https://golang.org/doc/modules/gomod-ref). We can get back to testing and learning Go now since the tests should run, even on Go 1.16.

In future chapters you will need to run `go mod init SOMENAME` in each new folder before running commands like `go test` or `go build`.

## Back to Testing

Run `go test` in your terminal. It should've passed! Just to check, try deliberately breaking the test by changing the `want` string.

Notice how you have not had to pick between multiple testing frameworks and then figure out how to install. Everything you need is built in to the language and the syntax is the same as the rest of the code you will write.

### Writing tests

테스트를 작성은 함수를 작성하는 것과 매우 비슷하며 다음과 같은 룰을 가집니다.

Writing a test is just like writing a function, with a few rules

* `xxx_test.go`와 같은 파일명으로 작성합니다.
* 테스트 함수는 반드시 `Test`라는 단어로 시작해야합니다.
* 테스트 함수는 `t *testing.T`라는 단 하나의 인자만 가집니다.
* `*testing.T` 타입을 사용하기 위해 `fmt`를 임포트한 것처럼 `testing` 패키지를 임포트해야합니다.

For now, it's enough to know that your `t` of type `*testing.T` is your "hook" into the testing framework so you can do things like `t.Fail()` when you want to fail.

We've covered some new topics:

#### `if`
If statements in Go are very much like other programming languages.

#### Declaring variables

We're declaring some variables with the syntax `varName := value`, which lets us re-use some values in our test for readability.

#### `t.Errorf`

We are calling the `Errorf` _method_ on our `t` which will print out a message and fail the test. The `f` stands for format which allows us to build a string with values inserted into the placeholder values `%q`. When you made the test fail it should be clear how it works.

You can read more about the placeholder strings in the [fmt go doc](https://golang.org/pkg/fmt/#hdr-Printing). For tests `%q` is very useful as it wraps your values in double quotes.

We will later explore the difference between methods and functions.

### Go doc

Another quality of life feature of Go is the documentation. You can launch the docs locally by running `godoc -http :8000`. If you go to [localhost:8000/pkg](http://localhost:8000/pkg) you will see all the packages installed on your system.

The vast majority of the standard library has excellent documentation with examples. Navigating to [http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/) would be worthwhile to see what's available to you.

If you don't have `godoc` command, then maybe you are using the newer version of Go (1.14 or later) which is [no longer including `godoc`](https://golang.org/doc/go1.14#godoc). You can manually install it with `go install golang.org/x/tools/cmd/godoc@latest`.

### Hello, YOU

Now that we have a test we can iterate on our software safely.

In the last example we wrote the test _after_ the code had been written just so you could get an example of how to write a test and declare a function. From this point on we will be _writing tests first_.

Our next requirement is to let us specify the recipient of the greeting.

Let's start by capturing these requirements in a test. This is basic test driven development and allows us to make sure our test is _actually_ testing what we want. When you retrospectively write tests there is the risk that your test may continue to pass even if the code doesn't work as intended.

```go
package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Chris")
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

`go test`로 테스트를 실행하면 컴파일 에러를 겪게 됩니다.

```text
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

Go와 같은 정적 타이핑 언어를 사용할 때는 컴파일러에 귀기울이는 것이 중요합니다. 왜냐면 컴파일러는 코드가 어떻게 엮여서 작동해야하는지 이해하기 때문입니다.

위와 같은 케이스에서는 컴파일러가 다음을 진행하기 위해 무엇이 필요한지 말해줍니다. 그러므로 `Hello` 함수가 인자를 허용하도록 코드를 수정해야합니다.

Edit the `Hello` function to accept an argument of type string

```go
func Hello(name string) string {
	return "Hello, world"
}
```

If you try and run your tests again your `hello.go` will fail to compile because you're not passing an argument. Send in "world" to make it compile.

```go
func main() {
	fmt.Println(Hello("world"))
}
```

Now when you run your tests you should see something like

```text
hello_test.go:10: got 'Hello, world' want 'Hello, Chris''
```

We finally have a compiling program but it is not meeting our requirements according to the test.

Let's make the test pass by using the name argument and concatenate it with `Hello,`

```go
func Hello(name string) string {
	return "Hello, " + name
}
```

When you run the tests they should now pass. Normally as part of the TDD cycle we should now _refactor_.

### A note on source control

At this point, if you are using source control \(which you should!\) I would
`commit` the code as it is. We have working software backed by a test.

I _wouldn't_ push to master though, because I plan to refactor next. It is nice
to commit at this point in case you somehow get into a mess with refactoring - you can always go back to the working version.

There's not a lot to refactor here, but we can introduce another language feature, _constants_.

### Constants

상수(Constants)는 다음과 같이 선언할 수 있습니다.

```go
const englishHelloPrefix = "Hello, "
```

그리고 코드를 다음과 같이 리팩토링할 수 있습니다.

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
	return englishHelloPrefix + name
}
```

After refactoring, re-run your tests to make sure you haven't broken anything.

상수를 사용함으로써 어플리케이션의 퍼포먼스를 향상 시킬 수 있습니다. 왜냐면 `Hello` 함수가 매번 호출 될 때 마다 `"Hello, "` 문자열 인스턴스를 생성했기 때문입니다.


To be clear, the performance boost is incredibly negligible for this example! But it's worth thinking about creating constants to capture the meaning of values and sometimes to aid performance.

## Hello, world... again

The next requirement is when our function is called with an empty string it defaults to printing "Hello, World", rather than "Hello, ".

Start by writing a new failing test

```go
func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris")
		want := "Hello, Chris"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
```

여기서는 이전의 테스트와 다른 형태의 서브 테스트가 사용되었습니다. 가끔은 여러 테스트를 하나의 그룹으로 묶고 서브 테스트들로 여러 다른 시나리오를 표현하는 것이 유용합니다.

이런 접근의 베네핏은 공통으로 사용하는 코드를 작성하여 다른 테스트에서도 재사용할 수 있다는 점입니다.

보이는 바와 같이 테스트에서 메시지가 예상한 바와 같은지를 체크하는 부분의 코드가 반복되고 있습니다.

리팩토링은 프로덕션 코드만을 위한 것이 아닙니다. 왜냐면 실제 코드가 어떻게 작동해야 하는지의 대한 명확한 정의가 테스트에서 설명 되어야하기 때문입니다.

그러므로 다음과 같이 테스트 코드를 리팩토링할 수 있습니다.

```go
func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("empty string defaults to 'world'", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})

}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

위와 같이 테스트 검증 코드를 하나의 함수로 만듬으로써 테스트의 가독성을 향상시키고 중복을 줄였습니다. 여기서 `t *testing.T`을 전달하므로 필요할 때 테스트 코드를 실패하게 할 수 있습니다.

`*testing.T`와 `*testing.B`를 모두 만족하는 인터페이스인 `testing.TB`를 사용하는 것은 매우 좋은 아이디어입니다. 헬퍼 함수를 테스트로부터 혹은 벤치마크로부터 호출할 수 있기 때문입니다.


`t.Helper()`는 test suite(test case들의 묶음)에게 이 함수가 헬퍼 함수라는 것을 알려주기 위해 필요합니다. 이렇게 함으로써 테스트가 실패할 경우 라인 넘버가 테스트 헬퍼 안에서 리포트 되지 않고 호출한 함수에서 리포트 되게 됩니다. 또한 다른 개발자들이 문제를 더 쉽게 파악하도록 해줍니다. If you still don't understand, comment it out, make a test fail and observe the test output. Comments in Go are a great way to add additional information to your code, or in this case, a quick way to tell the compiler to ignore a line. You can comment out the `t.Helper()` code by adding two forward slashes `//` at the beginning of the line. You should see that line turn grey or change to another color than the rest of your code to indicate it's now commented out.

Now that we have a well-written failing test, let's fix the code, using an `if`.

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return englishHelloPrefix + name
}
```

If we run our tests we should see it satisfies the new requirement and we haven't accidentally broken the other functionality.

### Back to source control

Now we are happy with the code I would amend the previous commit so we only
check in the lovely version of our code with its test.

### Discipline

개발 사이클을 다시 복습하면 다음과 같습니다.

* 테스트 작성를 작성합니다
* 컴파일러가 통과하도록 합니다
* 테스트를 실행하고, 실패하는 것을 확인한 후에 에러 메시지가 유의미한지 체크합니다.
* 테스트를 통과하기 위해 충분한 코드를 작성합니다
* 리팩토링합니다.

On the face of it this may seem tedious but sticking to the feedback loop is important.

Not only does it ensure that you have _relevant tests_, it helps ensure _you design good software_ by refactoring with the safety of tests.

Seeing the test fail is an important check because it also lets you see what the error message looks like. As a developer it can be very hard to work with a codebase when failing tests do not give a clear idea as to what the problem is.

By ensuring your tests are _fast_ and setting up your tools so that running tests is simple you can get in to a state of flow when writing your code.

By not writing tests you are committing to manually checking your code by running your software which breaks your state of flow and you won't be saving yourself any time, especially in the long run.

## Keep going! More requirements

Goodness me, we have more requirements. We now need to support a second parameter, specifying the language of the greeting. If a language is passed in that we do not recognise, just default to English.

We should be confident that we can use TDD to flesh out this functionality easily!

Write a test for a user passing in Spanish. Add it to the existing suite.

```go
	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})
```

Remember not to cheat! _Test first_. When you try and run the test, the compiler _should_ complain because you are calling `Hello` with two arguments rather than one.

```text
./hello_test.go:27:19: too many arguments in call to Hello
    have (string, string)
    want (string)
```

Fix the compilation problems by adding another string argument to `Hello`

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}
	return englishHelloPrefix + name
}
```

When you try and run the test again it will complain about not passing through enough arguments to `Hello` in your other tests and in `hello.go`

```text
./hello.go:15:19: not enough arguments in call to Hello
    have (string)
    want (string, string)
```

Fix them by passing through empty strings. Now all your tests should compile _and_ pass, apart from our new scenario

```text
hello_test.go:29: got 'Hello, Elodie' want 'Hola, Elodie'
```

We can use `if` here to check the language is equal to "Spanish" and if so change the message

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language == "Spanish" {
		return "Hola, " + name
	}
	return englishHelloPrefix + name
}
```

The tests should now pass.

Now it is time to _refactor_. You should see some problems in the code, "magic" strings, some of which are repeated. Try and refactor it yourself, with every change make sure you re-run the tests to make sure your refactoring isn't breaking anything.

```go
const spanish = "Spanish"
const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language == spanish {
		return spanishHelloPrefix + name
	}
	return englishHelloPrefix + name
}
```

### French

* Write a test asserting that if you pass in `"French"` you get `"Bonjour, "`
* See it fail, check the error message is easy to read
* Do the smallest reasonable change in the code

You may have written something that looks roughly like this

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language == spanish {
		return spanishHelloPrefix + name
	}
	if language == french {
		return frenchHelloPrefix + name
	}
	return englishHelloPrefix + name
}
```

## `switch`

만약 특정한 값을 체크하기위해 많은 `if`문을 사용해야한다면 `if`문 대신 `switch`문을 사용하는 것이 보편적입니다. `switch`문을 사용해 코드를 리팩토링함으로써 가독성을 향상시킬 수 있고 다양한 언어를 추후에 더 추가한다고 했을 때 확장성있는 코드를 작성할 수 있습니다.


```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	prefix := englishHelloPrefix

	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	}

	return prefix + name
}
```

Write a test to now include a greeting in the language of your choice and you should see how simple it is to extend our _amazing_ function.

### one...last...refactor?

이쯤되면 `Hello` 함수의 코드 사이즈가 크다는 논란의 여지가 생길 수 있습니다. 가장 간단하게 해결할 수 있는 방법은 특정 기능을 다른 함수로 분리하는 것입니다.

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}
```

A few new concepts:

* 함수 시그니처에서 `(prefix string)`이라는 이름을 가진 반환값을 만들었습니다.
* 이렇게 함으로써 함수 안에 `prefix`라는 변수가 생성됩니다.
  * 기본으로 **"제로"** 값이 할당됩니다. 이는 타입에 기반하는데, 예를 들어 정수면 `0`, 문자열이면 `""`이 할당됩니다.
    * This will display in the Go Doc for your function so it can make the intent of your code clearer.
* 스위치문 안에 `default`는 매칭되는 케이스문이 없는 경우 `if none`으로 분기됩니다.
* Go `public` 함수는 대문자로, `private` 함수는 소문자로 시작합니다. We don't want the internals of our algorithm to be exposed to the world, so we made this function private.

## Wrapping up

Who knew you could get so much out of `Hello, world`?

By now you should have some understanding of:

### Some of Go's syntax around

* Writing tests
* Declaring functions, with arguments and return types
* `if`, `const` and `switch`
* Declaring variables and constants

### The TDD process and _why_ the steps are important

* _Write a failing test and see it fail_ so we know we have written a _relevant_ test for our requirements and seen that it produces an _easy to understand description of the failure_
* Writing the smallest amount of code to make it pass so we know we have working software
* _Then_ refactor, backed with the safety of our tests to ensure we have well-crafted code that is easy to work with

In our case we've gone from `Hello()` to `Hello("name")`, to `Hello("name", "French")` in small, easy to understand steps.

This is of course trivial compared to "real world" software but the principles still stand. TDD is a skill that needs practice to develop, but by breaking problems down into smaller components that you can test, you will have a much easier time writing software.
