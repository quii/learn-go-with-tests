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

It doesn't look like there's any obvious refactors here that would improve the code so let's press on

The next shortcoming in `walk` is that it assumes every field is a `string`. Let's write a test for this scenario

## Write the test first

Add the following case

```go
{
    "Struct with non string field",
    struct {
        Name string
        Age  int
    }{"Chris", 33},
    []string{"Chris"},
},
```

## Try to run the test

```go
=== RUN   TestWalk/Struct_with_non_string_field
    --- FAIL: TestWalk/Struct_with_non_string_field (0.00s)
    	reflection_test.go:46: got [Chris <int Value>], want [Chris]
```

## Write enough code to make it pass

We need to check that the type of the field is a `string`.

```go
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String {
			fn(field.String())
		}
	}
}
```

We can do that by checking its [`Kind`](https://godoc.org/reflect#Kind)

## Refactor

Again it looks like the code is reasonable enough for now. 

The next scenario is what if it isn't a "flat" `struct`? In other words what happens if we have a `struct` with some nested fields?

## Write the test first

We have been using the anonymous struct syntax to declare types ad-hocly for our tests so we could continue to do that like so

```go
{
    "Nested fields",
    struct {
        Name string
        Profile struct {
            Age  int
            City string
        }
    }{"Chris", struct {
        Age  int
        City string
    }{33, "London"}},
    []string{"Chris", "London"},
},
```

But we can see that when you get inner anonymous structs the syntax gets a little messy. [There is a proposal to make it so the syntax would be nicer](https://github.com/golang/go/issues/12854).

Let's just refactor this by making a known type for this scenario and reference it in the test. There is a little indirection in that some of the code for our test is outside the test but readers should be able to infer the structure of the `struct` by looking at the initialisation.

Add the following type declarations somewhere in your test file

```go
type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}
```

Now we can add this to our cases which reads a lot clearer than before
```go
{
    "Nested fields",
    Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## Try to run the test

```
=== RUN   TestWalk/Nested_fields
    --- FAIL: TestWalk/Nested_fields (0.00s)
    	reflection_test.go:54: got [Chris], want [Chris London]
```

The problem is we're only iterating on the fields on the first level of the type's hierarchy 

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String {
			fn(field.String())
		}

		if field.Kind() == reflect.Struct {
			walk(field.Interface(), fn)
		}
	}
}
```

The solution is quite simple, we again inspect its `Kind` and if it happens to be a `struct` we just call `walk` again on that inner `struct`

## Refactor

```go
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.String:
			fn(field.String())
		case reflect.Struct:
			walk(field.Interface(), fn)
		}
	}
}
```

When you're doing a comparison on the same value more than once _generally_ refactoring into a `switch` will improve readability and make your code easier to extend.

What if the value of the struct passed in is a pointer?

## Write the test first

Add this case

```go
{
    "Pointers to things",
    &Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## Try to run the test

```
=== RUN   TestWalk/Pointers_to_things
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
	panic: reflect: call of reflect.Value.NumField on ptr Value
```

## Write enough code to make it pass

```go
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.String:
			fn(field.String())
		case reflect.Struct:
			walk(field.Interface(), fn)
		}
	}
}
```

You cant use `NumField` on a pointer `Value`, we need to extract the underlying value before we can do that by using `Elem()`

## Refactor

Let's encapsulate the responsibility of extracting the `reflect.Value` from a given `interface{}` into a function

```go
func walk(x interface{}, fn func(input string)) {
	val := getValue(x)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.String:
			fn(field.String())
		case reflect.Struct:
			walk(field.Interface(), fn)
		}
	}
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val
}
```

This actually adds _more_ code but I feel the abstraction level is right
- Get the `reflect.Value` of `x` so i can inspect it, I dont care how.
- Iterate over the fields, doing whatever needs to be done depending on its type
