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
- As a writer of such a function you have to be able to inspect _anything_ that has been passed to you and try and figure out what the type is and what you can do with it. This is done using _reflection_. This can be quite clumsy and difficult to read and is generally less performant (as you have to do checks at runtime).

In short only use reflection if you really need to.

If you want polymorphic functions, consider if you could design it around an interface (not `interface`, confusingly) so that users can use your function with multiple types if they implement whatever methods you need for your function to work.

Our function will need to be able to work with lots of different things. As always we'll take an iterative approach, writing tests for each new thing we want to support and refactoring along the way until we're done.

## Write the test first

We'll want to call our function with a struct that has a string field in it (`x`).Then we can spy on the function (`fn`) passed in to see if it is called.

```go
func TestWalk(t *testing.T) {

	expected := "Chris"
	var got []string

	x := struct {
		Name string
	}{expected}

	walk(x, func(input string) {
		got = append(got, input)
	})

	if len(got) != 1 {
		t.Errorf("wrong number of function calls, got %d want %d", len(got), 1)
	}
}
``` 

- We want to store a slice of strings (`got`) which stores which strings were passed into `fn` by `walk`. Often in previous chapters we have made dedicated types for this to spy on function/method invocations but in this case we can just pass in an anonymous function for `fn` that closes over `got`
- We use an anonymous `struct` with a `Name` field of type string to go for the simplest "happy" path.
- Finally call `walk` with `x` and the spy and for now just check the length of `got`, we'll be more specific with our assertions once we've got something very basic working.


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
	reflection_test.go:19: wrong number of function calls, got 0 want 1
FAIL
```

## Write enough code to make it pass

We can call the spy with any string to make this pass.

```go
func walk(x interface{}, fn func(input string)) {
    fn("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}
```

The test should now be passing. The next thing we'll need to do is make a more specific assertion on what our `fn` is being called with.

## Write the test first

Add the following to the existing test to check the string passed to `fn` is correct

```go
if got[0] != expected {
    t.Errorf("got '%s', want '%s'", got[0], expected)
}
```

## Try to run the test

```
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
	reflection_test.go:23: got 'I still can't believe South Korea beat Germany 2-0 to put them last in their group', want 'Chris'
FAIL
```

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)
	field := val.Field(0)
	fn(field.String())
}
```

This code is _very unsafe and very naive_ but remember our goal when we are in "red" (the tests failing) is to write the smallest amount of code possible. We then write more tests to address our concerns.

We need to use reflection to have a look at `x` and try and look at its properties.

The [reflect package](https://godoc.org/reflect) has a function `ValueOf` which returns us a `Value` of a given variable. This has ways for us to inspect a value, including its fields which we use on the next line. 

We then make some very silly assumptions about the the value passed in
- We look at the first and only field, there may be no fields at all which would cause a panic
- We then call `String()` which returns the underlying value as a string but we know it would be wrong if the field was something other than a string.

## Refactor

Our code is passing for the simple case but we know there's a lot of shortcomings in our code. 

We're going to be writing a number of tests where we pass in different values and checking the array of strings that `fn` was called with. 

We should refactor our test into a table based test to make this easier to continue testing new scenarios.

```go
func TestWalk(t *testing.T) {

	cases := []struct{
		Name string
		Input interface{}
		ExpectedCalls []string
	} {
		{
			"Struct with one string field",
			struct {
				Name string
			}{ "Chris"},
			[]string{"Chris"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}
}
```

Now we can easily add a scenario to see what happens if we have more than one string field

## Write the test first

Add the following scenario to the `cases`.

```go
{
    "Struct with two string fields",
    struct {
        Name string
        City string
    }{"Chris", "London"},
    []string{"Chris", "London"},
}
```

## Try to run the test

```
=== RUN   TestWalk/Struct_with_two_string_fields
    --- FAIL: TestWalk/Struct_with_two_string_fields (0.00s)
    	reflection_test.go:40: got [Chris], want [Chris London]
```

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	for i:=0; i<val.NumField(); i++ {
		field := val.Field(i)
		fn(field.String())
	}
}
```

`value` has a method `NumField` which returns the number of fields in the value. This lets us iterate over the fields and call `fn` which passes our test.

## Refactor

