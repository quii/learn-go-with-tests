# Reflection (WIP)

> #golang challenge: write a function `walk(x interface{}, fn func(string))` which takes a struct `x` and calls `fn` for all strings fields found inside.difficulty level: recursively.

Let's do it! 

## What is `interface` ?

We have enjoyed the type-safety that Go has offered us in terms of functions that work with known types, such as `strings`, `int32` and our own types like `BankAccount`.

This means that we get some documentation for free and the compiler will complain if you try and pass the wrong type to a function.

You may come across scenarios though where you want to write a function where you dont know the type at compile time.

Go lets us get around this with the type `interface{}` which you can think of as just _any_ type.

So `walk(x interface{}, fn func(string))` will accept any value for `x`.

### So why not use `interface` for everything and have really flexible functions?

- As a user of a function that takes `interface` you lose type safety.What if you meant to pass `Foo.bar` of type `string` into a function but instead did `Foo.baz` which is an `int`? The compiler wont be able to inform you of your mistake
- As a write of such a function you have to be able to inspect _anything_ that has been passed to you and try and figure out what the type is and what you can do with it.This is done using _reflection_.This can be quite clumsy and difficult to read for some and is generally less performant (as you have to do checks at runtime).

In short only use reflection if you really need to.

If you want polymorphic functions, consider if you could design it around an interface (not `interface`, confusingly) so that users can use your function with multiple types if they implement whatever methods you need for your function to work.

Our function will need to be able to work with lots of different things.As always we'll take an iterative approach, writing tests for each new thing we want to support and refactoring along the way until we're done.

## Write the test first

We'll want to call our function with a struct that has a string field in it (`x`).Then we can spy on the function (`fn`) passed in to see if it is called.

```go
type CallSpy []string

func (c *CallSpy) Fn(input string) {
    *c = append(*c, input)
}

func TestWalk(t *testing.T) {

    expected := "Chris"
	
    x := struct {
        Name string
    }{expected}

    var fnSpy CallSpy
	
    walk(x, fnSpy.Fn)
	
    if len(fnSpy) != 1 {
        t.Errorf("wrong number of calls to CallSpy, got %d want %d", len(fnSpy), 1)
    }
}
``` 

- We want to store a slice of strings which represent the function calls we got so we've made a new type `CallSpy` based on `[]string`.This lets us add a method which we pass in to `walk`.Every time the method is called we add the `input` to the underlying `[]string` so we can spy on the calls that `walk` makes.
- We use an anonymous `struct` with a `Name` field of type string to go for the simplest "happy" path.
- Finally call `walk` with `x` and the spy and for now just check the number of calls, we'll be more specific with our assertions once we've got something very basic working.


## Try to run the test

```
./reflection_test.go:21:2: undefined: walk
```

## Write the minimal amount of code for the test to run and check the failing test output

We need to define `walk`

```go
func walk(x interface{}, fn func(input string)) {

}
```

Try and run the test again

```
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
	reflection_test.go:24: wrong number of calls to CallSpy, got 0 want 1
FAIL
```

## Write enough code to make it pass

We can call the spy with any string to make this pass.

```go
func walk(x interface{}, fn func(input string)) {
    fn("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}
```

The test should now be passing.The next thing we'll need to do is make a more specific assertion on what our `fn` is being called with.
