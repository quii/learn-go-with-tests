# Integers

Integers work as you would expect. Let's write an add function to try things out

## Test first

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

Run the test `go test`

Inspect the compilation error

`./adder_test.go:6:9: undefined: Add`

Write enough code to satisfy the compiler *and that's all* - remember we want to check that our tests fail for the correct reason.

```go
func Add(x, y int) (sum int) {
	return
}
```

We've introduced a few new concepts when writing functions here. 

- When you have more than one argument of the same type (in our case two integers) rather than having `(x int, y int)` you can shorten it to `(x, y int)`
- You can assign a name to the return value `(sum int)` 
    - This will create a variable called `sum` in your function
    - It will be assigned the "zero" value. This depends on the type, for example `int`s are 0 and for strings it is `""` 
     - You can return whatever it's set to by just calling `return` rather than `return sum`. 
    - This will display in the Go Doc for your function so it can make the intent of your code clearer.
    
Now run the tests and we should be happy that the test is correctly reporting what is wrong.

`adder_test.go:10: expected '4' but got '0'`

In the strictest sense of TDD we should now write the _minimal amount of code to make the test pass_. A pedantic programmer may do this

```go
func Add(x, y int) (sum int) {
	return 4
}
```

Ah hah! Foiled again, TDD is a sham right?

We could write another test, with some different numbers to force that test to fail but that feels like a game of cat and mouse. 

What I am going to introduce is an interesting approach to testing called "Property Based Tests"

