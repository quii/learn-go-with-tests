# Hello, world

It is traditional for your first program in a new language to be Hello, world. Create a file called `hello.go` and write this code. To run it type `go run hello.go`.

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
```

## How it works
When you write a program in Go you will have a `main` package defined with a `main` func inside it. The `func` keyword is how you define a function with a name and a body.

With `import "fmt"` we are importing a package which contains the `Println` function that we use to print

## How to test

How do you test this? It is good to separate your "domain" code from the outside world (side-effects). The `fmt.Println` is a side effect (printing to stdout) and the string we send in is our domain.

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

We have created a new function again with `func` but this time we've added another keyword `string` in the definition. This means this function returns a `string`. 

Now create a new file called `hello_test.go` where we are going to write a test for our `Hello` function

```go
package main

import "testing"

func TestHello(t *testing.T) {
	message := Hello()
	expected := "Hello, world"

	if message != expected {
		t.Errorf("expected '%s' but got '%s'", expected, message)
	}
}
```

Before explaining, let's just run the code. Type `go test`. It should've passed! Just to check, try deliberately breaking the test by changing the `expected` string.

Notice how you have not had to pick between multiple testing frameworks or decipher a testing DSL to write a test. Everything you need is built in to the language.

### Side-rant

Go famously "lacks" a number of programming language features and this is a constant point of discussion. I admittedly would be very pleased if Go supported generics. 

But it is important to understand that syntax is not the _only_ factor in how effective you can be as a programmer. We have just demonstrated one of the reasons Go is popular, it has first class support for testing out of the box and it is no different from writing the "real" code. 

### Back to the tests

Writing a test is just like writing a function, with a few rules

- It needs to be in a file with a name like `xxx_test.go`
- The test function must start with the word `Test`
- The test function takes one argument only `t *testing.T`

For now it's enough to know that your `t` of type `*testing.T` is your "hook" into the testing framework so you can do things like `t.Fail()` when you want to fail. 

#### New things

##### `if`

If statements in Go are very much like other programming languages. 

##### Declaring variables

We're declaring some variables with the syntax `varName := value`, which lets us re-use some values in our test for readability

##### `t.ErrorF`

We are calling the `ErrorF` _method_ on our `t` which will print out a message and fail the test. The `F` stands for format which allows us to build a string with values inserted into the placeholder values `%s`. When you made the test fail it should be clear how it works. 

We will later explore the difference between methods and functions.

### Go doc

Another quality of life feature of Go is the documenation. You can launch the docs locally by running `godoc -http :8000`. If you go to [localhost:8000/pkg](http://localhost:8000/pkg) you will see all the packages installed on your system.

The vast majority of the standard library has excellent documentation with examples. Navigating to [http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/) would be worthwhile to see what's available to you. 

### Hello, YOU

Now that we have a test we can iterate on our software safely. Our next requirement is to let us specify the recipient of the greeting. 

Let's start by capturing our requirements in the test. This is basic test driven development and allows us to make sure our test is _actually_ testing what we want. When you retrospectively write tests there is the risk that your test may continue to pass even if the code doesn't work as intended. 

```go
package main

import "testing"

func TestHello(t *testing.T) {
	message := Hello("Chris")
	expected := "Hello, Chris"

	if message != expected {
		t.Errorf("expected '%s' but got '%s'", expected, message)
	}
}
```

Now run `go test`, you should have a compilation error

```
./hello_test.go:6:18: too many arguments in call to Hello
	have (string)
	want ()
```

When using a statically typed language like Go it is important to _listen to the compiler_. The compiler understands how your code should snap together and work so you don't have to. 

In this case the compiler is telling you what you need to do to continue. We have to change our function `Hello` to accept an argument.

Edit the `Hello` function to accept an argument of type string 

```go
func Hello(name string) string {
	return "Hello, world"
}
```

If you try and run your tests again your `main.go` will fail to compile because you're not passing an argument. Send in "world" to make it pass.

```go
func main() {
	fmt.Println(Hello("world"))
}
```

Now when you run your tests you should see something like

```
hello_test.go:10: expected 'Hello, Chris' but got 'Hello, world'
```

We finally have a compiling program but it is not meeting our requirements according to the test. 

Let's make the test pass by using the name argument and concatenate it with `Hello, `

```go
func Hello(name string) string {
	return "Hello, " + name
}
```

When you run the tests they should now pass. Normally as part of the TDD cycle we should now *refactor*.

There's not a lot to refactor here, but we can introduce another language feature *constants*

### Constants

Constants are defined like so

```go
const helloPrefix = "Hello, "
```

We can now refactor our code like so

```go
const helloPrefix = "Hello, "

func Hello(name string) string {
	return helloPrefix + name
}
```

After refactoring, re-run your tests to make sure you haven't broken anything.

Constants should improve performance of your application as it saves you creating the `"Hello, "` string instance every time `Hello` is called. 

To be clear, this is incredibly negligble for this example! But it's worth thinking about creating constants to capture the meaning of values and sometimes to aid performance.

### Discipline

On the face of it, the cycle of writing a test, failing the compiler, making the code pass and then refactoring may seem tedious but sticking to the feedback loop is important. 

Not only does it ensure that you have *relevant tests* it helps ensure *you design good software* by refactoring with the safety of tests. 

By ensuring your tests are *fast* and setting up your tools so that running tests is simple you can get in to a state of flow when writing your code. 

By not writing tests you are committing to manually checking your code by running your software which breaks your state of flow and you wont be saving yourself any time, especially in the long run. 