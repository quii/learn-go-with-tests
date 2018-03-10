# Integers

Integers work as you would expect. Let's write an add function to try things out

## Write the test first

```go
package main

import "testing"

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4
	
	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}
```

You will notice that we're using `%d` as our format strings rather than `%s`. That's because we want it to print an integer rather than a string.

## 2. Try and run the test

Run the test `go test`

Inspect the compilation error

`./adder_test.go:6:9: undefined: Add`

## 3. Write the minimal amount of code for the test to run and check the failing test output

Write enough code to satisfy the compiler *and that's all* - remember we want to check that our tests fail for the correct reason.

```go
func Add(x, y int) (sum int) {
	return
}
```

When you have more than one argument of the same type (in our case two integers) rather than having `(x int, y int)` you can shorten it to `(x, y int)`    
    
Now run the tests and we should be happy that the test is correctly reporting what is wrong.

`adder_test.go:10: expected '4' but got '0'`

## 4. Write enough code to make it pass

In the strictest sense of TDD we should now write the _minimal amount of code to make the test pass_. A pedantic programmer may do this

```go
func Add(x, y int) (sum int) {
	return 4
}
```

Ah hah! Foiled again, TDD is a sham right?

We could write another test, with some different numbers to force that test to fail but that feels like a game of cat and mouse. 

Once we're more familiar with Go's syntax I will introduce a technique called Property Based Testing, which would stop annoying developers and help find you bugs

For now, let's fix it properly

```go
func Add(x, y int) (sum int) {
	return x + y
}
```

If you re-run the tests they should pass. 

## 5. Refactor

There's not a lot in the _actual_ code we can really improve on here.

We explored earlier how by naming the return argument it appears in the documentation but also in most developer's text editors. 

This is great because it aids the usability of code you are writing. It is preferable that a user can understand the usage of your code by just looking at the type signature and documentation. 

You can add documentation to functions with comments, and these will appear in Go Doc just like when you look at the standard library's documentation.

```go
// Add takes two integers and returns the sum of them
func Add(x, y int) (sum int) {
	return x + y
}
```

#### Examples

If you really want to go the extra mile you can make examples. You will find a lot of examples in the documentation of the standard library

Often code examples go out of date with what the actual code does because they live outside of the real code and don't get checked. 

Go examples are executed just like tests so you can be confident examples reflect what the code actually does. 

```go
func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
```

If your code changes so that the example is no longer valid, your build will fail. 

By adding this code in the example will appear in the documentation inside `godoc` making your code even more accessible.

## Wrapping up

What we have covered:

- More practice of the TDD workflow
- Integers, addition
- Writing better documentation so users of our code can understand its usage quickly
- Examples of how to use our code, which are checked as part of our tests