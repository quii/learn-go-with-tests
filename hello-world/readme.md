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

Notice how you have not had to pick between multiple testing frameworks or decipher a testing DSL to write a test. Everything you need is built in to the language and the syntax is the same as the rest of the code you will write. 

### Writing tests

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

To be clear, the performance boost is incredibly negligible for this example! But it's worth thinking about creating constants to capture the meaning of values and sometimes to aid performance.

## Hello, world... again

The next requirement is when our function is called with an empty string it defaults to printing "Hello, World", rather than "Hello, "

Start by writing a new failing test

```go
func TestHello(t *testing.T) {

	t.Run("saying hello to people", func(t *testing.T) {
		message := Hello("Chris")
		expected := "Hello, Chris"

		if message != expected {
			t.Errorf("expected '%s' but got '%s'", expected, message)
		}
	})
	
	t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
		message := Hello("")
		expected := "Hello, World"

		if message != expected {
			t.Errorf("expected '%s' but got '%s'", expected, message)
		}
	})

}
```

Here we are introducing another tool in our testing arsenal, subtests. Sometimes it is useful to group tests around a "thing" and then have subtests describing different scenarios. 

A benefit of this approach is you can set up shared code that can be used in the other tests.

There is repeated code when we check if the message is what we expect. 

Refactoring is not _just_ for the production code! We can and should refactor our tests.

```go
func TestHello(t *testing.T) {

	assertCorrectMessage := func(expected, actual string) {
		t.Helper()
		if expected != actual {
			t.Errorf("expected '%s' but got '%s'", expected, actual)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		message := Hello("Chris")
		expected := "Hello, Chris"
		assertCorrectMessage(expected, message)
	})

	t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
		message := Hello("")
		expected := "Hello, World"
		assertCorrectMessage(expected, message)
	})

}
```

What have we done here? In Go you can declare functions inside other functions and then they can _close_ over other variables - in this case our `*testing.T`.

We've written a function to do our assertion. This reduces duplication and improves readability of our tests.

Now that we have a well-written failing test, let's fix the code, using the `else` keyword.

`TODO:// explain t.Helper()`

```go
const helloPrefix = "Hello, "

func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return helloPrefix + name
}
```

If we run our tests we should see it satisfies the new requirement and we haven't accidentally broken the other functionality

### Discipline

On the face of it, the cycle of writing a test, failing the compiler, making the code pass and then refactoring may seem tedious but sticking to the feedback loop is important. 

Not only does it ensure that you have *relevant tests* it helps ensure *you design good software* by refactoring with the safety of tests. 

Seeing the test fail is an important check because it also lets you see what the error message looks like. As a developer it can be very hard to work with a codebase when failing tests do not give a clear idea as to what the problem is. 

By ensuring your tests are *fast* and setting up your tools so that running tests is simple you can get in to a state of flow when writing your code. 

By not writing tests you are committing to manually checking your code by running your software which breaks your state of flow and you wont be saving yourself any time, especially in the long run. 

## Keep going! More requirements

Goodness me, we have more requirements. We now need to support a second parameter, specifying the language of the greeting. If a language is passed in that we do not recognise, just default to English.

We should be confident that we can use TDD to flesh out this functionality easily!

Write a test for a user passing in Spanish. Add it to the existing suite.

```go
	t.Run("say hello in Spanish", func(t *testing.T) {
		message := Hello("Elodie", "Spanish")
		expected := "Hola, Elodie"
		assertCorrectMessage(expected, message)
	})
```

Remember not to cheat! *Test first*. When you try and run the test, the compiler _should_ complain because you are calling `Hello` with two arguments rather than one.

```
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
	return helloPrefix + name
}
```

When you try and run the test again it will complain about not passing through enough arguments to `Hello` in your other tests and in `main.go`

```
./hello.go:15:19: not enough arguments in call to Hello
	have (string)
	want (string, string)
```

Fix them by passing through empty strings. Now all your tests should compile _and_ pass, apart from our new scenario

```
hello_test.go:29: expected 'Hola, Elodie' but got 'Hello, Elodie'
```

We can use `if` here to check the language is equal to "Spanish" and if so change the message

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language=="Spanish" {
		return "Hola, " + name
	}
	
	return helloPrefix + name
}
```

The tests should now pass. 

Now it is time to *refactor*. You should see some problems in the code, "magic" strings, some of which are repeated. Try and refactor it yourself, with every change make sure you re-run the tests to make sure your refactoring isn't breaking anything.

```go
const spanish = "Spanish"
const helloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language == spanish {
		return spanishHelloPrefix + name
	}

	return helloPrefix + name
}
```

### French

- Write a test asserting that if you pass in `"french"` you get `"Bonjour, "`
- See it fail, check the error message is easy to read
- Do the smallest reasonable change in the code

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

	return helloPrefix + name
}
```

## `switch`

When you have lots of `if` statements checking a particular value it is common to use a `switch` statement instead. We can use `switch` to refactor the code to make it easier to read and more extensible if we wish to add more language support later

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	prefix := helloPrefix
	
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

## Wrapping up

Who knew you could get so much out of `Hello, world` ?

By now you should have some understanding of

### Some of Go's syntax around

- Writing tests
- Declaring functions, with arguments and return types
- `if`, `else`, `switch`
- Declaring variables and constants

### An understanding of the TDD process and _why_ the steps are important

- *Write a failing test and see it fail* so we know we have written a _relevant_ test for our requirements and seen that it produces an _easy to understand description of the failure_
- Writing the smallest amount of code to make it pass so we know we have working software
- _Then_ refactor, backed with the safety of our tests to ensure we have well-crafted code that is easy to work with

 In our case we've gone from `Hello()` to `Hello("name")`, to `Hello("name", "french")` in small, easy to understand steps. 
 
 This is of course trivial compared to "real world" software but the principles still stand. TDD is a skill that needs practice to develop but by being able to break problems down into smaller components that you can test you will have a much easier time writing software.