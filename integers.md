# Integers

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/integers)**

Integers work as you would expect. Let's write an `Add` function to try things out. Create a test file called `adder_test.go` and write this code.

**Note:** Go source files can only have one `package` per directory. Make sure that your files are organised into their own packages. [Here is a good explanation on this.](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)

Your project directory might look something like this:

```
learnGoWithTests
    |
    |-> helloworld
    |    |- hello.go
    |    |- hello_test.go
    |
    |-> integers
    |    |- adder_test.go
    |
    |- go.mod
    |- README.md
```

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

You will notice that we're using `%d` as our format strings rather than `%q`. That's because we want it to print an integer rather than a string.

Also note that we are no longer using the main package, instead we've defined a package named `integers`, as the name suggests this will group functions for working with integers such as `Add`.

## Try and run the test

Run the test `go test`

Inspect the compilation error

`./adder_test.go:6:9: undefined: Add`

## Write the minimal amount of code for the test to run and check the failing test output

Write enough code to satisfy the compiler _and that's all_ - remember we want to check that our tests fail for the correct reason.

```go
package integers

func Add(x, y int) int {
	return 0
}
```

Remember, when you have more than one argument of the same type \(in our case two integers\) rather than having `(x int, y int)` you can shorten it to `(x, y int)`.

Now run the tests, and we should be happy that the test is correctly reporting what is wrong.

`adder_test.go:10: expected '4' but got '0'`

If you have noticed we learnt about _named return value_ in the [last](hello-world.md#one...last...refactor?) section but aren't using the same here. It should generally be used when the meaning of the result isn't clear from context, in our case it's pretty much clear that `Add` function will add the parameters. You can refer [this](https://go.dev/wiki/CodeReviewComments#named-result-parameters) wiki for more details.

## Write enough code to make it pass

In the strictest sense of TDD we should now write the _minimal amount of code to make the test pass_. A pedantic programmer may do this

```go
func Add(x, y int) int {
	return 4
}
```

Ah hah! Foiled again, TDD is a sham right?

We could write another test, with some different numbers to force that test to fail but that feels like [a game of cat and mouse](https://en.m.wikipedia.org/wiki/Cat_and_mouse).

Once we're more familiar with Go's syntax I will introduce a technique called _"Property Based Testing"_, which would stop annoying developers and help you find bugs.

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

You can add documentation to functions with comments, and these will appear in Go Doc just like when you look at the standard library's documentation.

```go
// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
	return x + y
}
```

### Testable Examples

If you really want to go the extra mile you can make [Testable Examples](https://blog.golang.org/examples). You will find many examples in the standard library documentation.

Often code examples that can be found outside the codebase, such as a readme file, become out of date and incorrect compared to the actual code because they don't get checked.

Example functions are compiled whenever tests are executed. Because such examples are validated by the Go compiler, you can be confident your documentation's examples always reflect current code behavior.

Example functions begin with `Example` (much like test functions begin with `Test`), and reside in a package's `_test.go` files. Add the following `ExampleAdd` function to the `adder_test.go` file.

```go
func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
```

(If your editor doesn't automatically import packages for you, the compilation step will fail because you will be missing `import "fmt"` in `adder_test.go`. It is strongly recommended you research how to have these kind of errors fixed for you automatically in whatever editor you are using.)

Adding this code will cause the example to appear in your documentation, making your code even more accessible. If ever your code changes so that the example is no longer valid, your build will fail.

Running the package's test suite, we can see the example `ExampleAdd` function is executed with no further arrangement from us:

```bash
$ go test -v
=== RUN   TestAdder
--- PASS: TestAdder (0.00s)
=== RUN   ExampleAdd
--- PASS: ExampleAdd (0.00s)
```

Notice the special format of the comment, `// Output: 6`. While the example will always be compiled, adding this comment means the example will also be executed. Go ahead and temporarily remove the comment `// Output: 6`, then run `go test`, and you will see `ExampleAdd` is no longer executed.

Examples without output comments are useful for demonstrating code that cannot run as unit tests, such as that which accesses the network, while guaranteeing the example at least compiles.

To view example documentation, let's take a quick look at `pkgsite`. Navigate to your project's directory, then run `pkgsite -open .`, which should open a web browser for you, pointing to `http://localhost:8080`. Inside here you'll see a list of all of Go's Standard Library packages, plus Third Party packages you have installed, under which you should see your example documentation for `github.com/quii/learn-go-with-tests`. Follow that link, and then look under `Integers`, then under `func Add`, then expand `Example` and you should see the example you added for `sum := Add(1, 5)`.

If you publish your code with examples to a public URL, you can share the documentation of your code at [pkg.go.dev](https://pkg.go.dev/). For example, [here](https://pkg.go.dev/github.com/quii/learn-go-with-tests/integers/v2) is the finalised API for this chapter. This web interface allows you to search for documentation of standard library packages and third-party packages.

## Wrapping up

What we have covered:

*   More practice of the TDD workflow
*   Integers, addition
*   Writing better documentation so users of our code can understand its usage quickly
*   Examples of how to use our code, which are checked as part of our tests
