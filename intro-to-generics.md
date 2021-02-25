# Intro to generics

(At the time of writing) Go does not have support for user-defined generics, but [the proposal](https://blog.golang.org/generics-proposal) [has been accepted](https://github.com/golang/go/issues/43651#issuecomment-776944155) and will be included in 1.18

However, there are ways to experiment with the upcoming implementation using the [go2go playground](https://go2goplay.golang.org/) _today_.

This chapter will give you a brief introduction to generics, hopefully dispel any reservations you may have about them and give you an idea of how you will be able to simplify some of your code in the future.

The code we write here will be the foundation for future chapters around generics.

## Setting up the playground

In the go2go playground we can't run `go test`, so how are we going to write tests to explore generic code?

The playground _does_ let us execute code, and because we're programmers that means we can work around the lack of a test runner by **making one of our own**.

## Our own test helpers (`AssertEqual`, `AssertNotEqual`)

To explore generics in future chapters we need to write some test helpers that'll kill the program and print something useful if a test fails.

### Assert on integers

Let's start with something basic and iterate toward our goal

```go
package main

import (
	"log"
)

func main() {
	AssertEqual(1, 1)
	AssertNotEqual(1, 2)

	AssertEqual(50, 100)

	AssertNotEqual(2, 2) // wont see this
}

func AssertEqual(got, want int) {
	if got != want {
		log.Fatalf("FAILED: got %d, want %d", got, want)
	} else {
		log.Printf("PASSED: %d did equal %d\n", got, want)
	}
}

func AssertNotEqual(got, want int) {
	if got == want {
		log.Fatalf("FAILED: got %d, want %d", got, want)
	} else {
		log.Printf("PASSED: %d did not equal %d\n", got, want)
	}
}
```

[This program prints](https://go2goplay.golang.org/p/WywgJnAp34v)

```
2009/11/10 23:00:00 PASSED: 1 did equal 1
2009/11/10 23:00:00 PASSED: 1 did not equal 2
2009/11/10 23:00:00 FAILED: got 50, want 100
```

### Iteration 2

Being able to assert on the equality of integers is great but what if we want to assert on `string` ?

```go
func main() {
	AssertEqual("CJ", "CJ")
}
```

You'll get an error

```
type checking failed for main
prog.go2:8:14: cannot use "CJ" (untyped string constant) as int value in argument to AssertEqual
```

If you take your time to read the error, you'll see the compiler is complaining we're trying to pass a `string` to a function that expects an `integer`.

#### Recap on type-safety

If you've read the previous chapters of this book, or have experience with statically typed languages this should not surprise you. The Go compiler expects you to write your functions, structs etc by describing what types you wish to work with.

You can't pass a `string` to a function that expects an `integer`.

Whilst this can feel like ceremony, it can be extremely helpful. By describing these constraints you

- Make function implementation simpler. By describing to the compiler what types you work with you **constrain the number of possible valid implementations**. You can't "add" a `Person` and a `BankAccount`. You can't capitalise an `integer`. In software constraints are often extremely helpful.
- Prevents you accidentally passing data to a function that you didn't mean to

Go currently offers you a way to be more abstract with your types with interfaces so that you can design functions that do not take concrete types but instead types that offer the behaviour you need. This gives you some flexibility whilst maintaining type-safety.

### A function that takes a string or an integer? (or indeed, other things)

The other option that Go _currently_ gives is declaring the type of your argument as `interface{}` which means "anything".

Try changing the signatures to use this type instead.

```go
func AssertEqual(got, want interface{}) {

func AssertNotEqual(got, want interface{}) {

```

The tests should now compile and pass. The output will be a bit ropey because we're using the `%d` format string to print our messages so change them to `%+v` for a better output.

### Tradeoffs made without generics

Our `AssertX` functions are quite naive but conceptually aren't too different to how other [popular libraries offer this functionality](https://github.com/matryer/is/blob/master/is.go#L150)

```go
func (is *I) Equal(a, b interface{}) {
```

So what's the problem?

By using `interface{}` the compiler can't help us when writing our code, because we're not telling it anything useful about the types of things passed to the function. Go back to the go2go playground and try comparing two different types

```go
AssertNotEqual(1, "1")
```

Now to be fair in this case, we get away with it; the test compiles, and it fails as we'd hope; but in a cosy type-safe world do we want to be able to compare strings and integers?

In our case we get away with it but writing functions that take `interface{}` can be extremely challenging and bug-prone because we've _lost_ our constraints, and we have no information at compile time as to what kind of data we're dealing with.

Developers often have to use reflection to implement these *ahem* generic functions, which is usually painful and can hurt the performance of your program.

## Our own test helpers with generics

Ideally, we don't want to have to make specific `AssertX` functions for every type we ever deal with. We'd like to be able to have _one_ `AssertEqual` function that works with _any_ type but does not let you compare apples with oranges.

Generics offer us a new way to make abstractions (like interfaces) by letting us **describe our constraints** in ways we cannot currently do.

```go
package main

import (
    "log"
)

func main() {
    AssertEqual(1, 1)
    AssertEqual("1", "1")
    AssertNotEqual(1, 2)
    //AssertEqual(1, "1") - uncomment me to see compilation error
}

func AssertEqual[T comparable](got, want T) {
    if got != want {
        log.Fatalf("FAILED: got %+v, want %+v", got, want)
    } else {
        log.Printf("PASSED: %+v did equal %+v\n", got, want)
    }
}

func AssertNotEqual[T comparable](got, want T) {
    if got == want {
        log.Fatalf("FAILED: got %+v, want %+v", got, want)
    } else {
        log.Printf("PASSED: %+v did not equal %+v\n", got, want)
    }
}
```

[go2go playground link](https://go2goplay.golang.org/p/a-6MzWrjeAx)

To write generic functions in Go, you need to provide "type parameters" which is just a fancy way of saying "describe your generic type and give it a label".

In our case the type of our type parameter is [`comparable`](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md#comparable-types-in-constraints) and we've given it the label of `T`. This label then lets us describe the types for the arguments to our function.

We're using `comparable` because we want to describe to the compiler that we wish to use `==` and `!=` - we want to compare! If you try changing the type to `any`

```go
func AssertNotEqual[T any](got, want T) {
```

You'll get the following error

```
prog.go2:15:5: cannot compare got != want (operator != not defined for T)
```

Which makes a lot of sense, because you can't use those operators on every (or `any`) type.

## Next: Generic data types

//todo: Implement a generic stack (start with stack of ints, stack of strings, refactor into generic version)
// https://go2goplay.golang.org/p/cbtP3zCNh7v

## Wrapping up

Hopefully this chapter has given you a taste of generics syntax and give you some ideas as to why they might be helpful. We've written our own `Assert` functions which we can safely re-use to experiment with other ideas around generics, and we implemented a simple data structure where it can store any type of data we wish in a type-safe manner.

The next chapters will explore:

- Defining our own type parameters
- Multiple type parameters

### Will generics turn Go into Java?

- No.
- Stop being rude about Java, it's not nice. It's nice to be nice.

There's a lot of [FUD](https://en.wikipedia.org/wiki/Fear,_uncertainty,_and_doubt) in the Go community about generics leading to nightmare abstractions and baffling code bases.

This is usually caveatted with "they must be used carefully". Whilst this is true, it's not especially useful advice because this is true of any language feature.

I know this because I have written extremely awful code _without_ generics.

### You're already using generics

The FUD becomes even sillier when you consider that if you've used arrays, slices or maps; you've already been a consumer of generic code.

```go
var myApples []Apples
// You cant do this!
append(myApples, Orange{})
```

### Make it work, make it right, make it fast

People run in to problems with generics when they're abstracting too quickly without enough information.

The TDD cycle of red, green, refactor means that you have more guidance as to what code you _actually need_ to deliver your behaviour, **rather than imagining abstractions up front**; but you still need to be careful.

There's no hard and fast rules here but resist making things generic until you can see that you have a useful generalisation. This may take the form of writing a number of tests and _then_ noticing a pattern when you're refactoring.

People often advise you to only generalise when you see the same code 3 times, which seems like a good starting rule of thumb.

A common path I've taken in other programming languages has been:

- One TDD cycle to drive some behaviour
- Another TDD cycle to exercise some other related scenarios

> Hmm, these things look similar - but a little duplication is better than coupling to a bad abstraction

- Sleep on it
- Another TDD cycle

> OK, I'd like to try to see if I can generalise this thing. Thank goodness I am so smart and handsome because I use TDD, so I can refactor whenever I wish, and the process has helped me understand what behaviour I actually need before designing too much.

- The abstraction feels nice! The tests are still passing, and the code is simpler
- I can now delete a number of tests, I've captured the _essence_ of the behaviour and removed unnecessary detail


### Abstraction is not a dirty word

It's easy to dunk on [AbstractSingletonProxyFactoryBean](https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/aop/framework/AbstractSingletonProxyFactoryBean.html) but let's not pretend a code base with no abstraction at all isn't also bad. It's your job to _gather_ related concepts when appropriate, so your system is easier to understand and change; rather than being a collection of disparate functions and types with a lack of clarity.
