# Intro to generics

(At the time of writing) Go does not have support for user-defined generics, but [the proposal](https://blog.golang.org/generics-proposal) [has been accepted](https://github.com/golang/go/issues/43651#issuecomment-776944155) and will be included in version 1.18.

However, there are ways to experiment with the upcoming implementation using the [go2go playground](https://go2goplay.golang.org/) _today_. So to work through this chapter you'll have to leave your precious editor of choice and instead do the work within the playground.

This chapter will give you an introduction to generics, dispel reservations you may have about them and, give you an idea how to simplify some of your code in the future. After reading this you'll know how to write:

- A function that takes generic aguments
- A generic data-structure

## Setting up the playground

In the _go2go playground_ we can't run `go test`. How are we going to write tests to explore generic code?

The playground _does_ let us execute code, and because we're programmers that means we can work around the lack of a test runner by **making one of our own**.

## Our own test helpers (`AssertEqual`, `AssertNotEqual`)

To explore generics we'll write some test helpers that'll kill the program and print something useful if a test fails.

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

	AssertEqual(50, 100) // this should fail

	AssertNotEqual(2, 2) // so you wont see this print
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

[This program prints](https://go2goplay.golang.org/p/WywgJnAp34v):

```
2009/11/10 23:00:00 PASSED: 1 did equal 1
2009/11/10 23:00:00 PASSED: 1 did not equal 2
2009/11/10 23:00:00 FAILED: got 50, want 100
```

### Assert on strings

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

If you take your time to read the error, you'll see the compiler is complaining that we're trying to pass a `string` to a function that expects an `integer`.

#### Recap on type-safety

If you've read the previous chapters of this book, or have experience with statically typed languages, this should not surprise you. The Go compiler expects you to write your functions, structs etc by describing what types you wish to work with.

You can't pass a `string` to a function that expects an `integer`.

Whilst this can feel like ceremony, it can be extremely helpful. By describing these constraints you,

- Make function implementation simpler. By describing to the compiler what types you work with, you **constrain the number of possible valid implementations**. You can't "add" a `Person` and a `BankAccount`. You can't capitalise an `integer`. In software, constraints are often extremely helpful.
- Are prevented from accidentally passing data to a function you didn't mean to.

Go currently offers you a way to be more abstract with your types with interfaces, so that you can design functions that do not take concrete types but instead, types that offer the behaviour you need. This gives you some flexibility whilst maintaining type-safety.

### A function that takes a string or an integer? (or indeed, other things)

The other option that Go _currently_ gives is declaring the type of your argument as `interface{}` which means "anything".

Try changing the signatures to use this type instead.

```go
func AssertEqual(got, want interface{}) {

func AssertNotEqual(got, want interface{}) {

```

The tests should now compile and pass. The output will be a bit ropey because we're using the integer `%d` format string to print our messages, so change them to the general `%+v` format for a better output of any kind of value.

### Tradeoffs made without generics

Our `AssertX` functions are quite naive but conceptually aren't too different to how other [popular libraries offer this functionality](https://github.com/matryer/is/blob/master/is.go#L150)

```go
func (is *I) Equal(a, b interface{}) {
```

So what's the problem?

By using `interface{}` the compiler can't help us when writing our code, because we're not telling it anything useful about the types of things passed to the function. Go back to the _go2go playground_ and try comparing two different types,

```go
AssertNotEqual(1, "1")
```

In this case, we get away with it; the test compiles, and it fails as we'd hope, although the error message `got 1, want 1` is unclear; but do we want to be able to compare strings with integers? What about comparing a `Person` with an `Airport`?

Writing functions that take `interface{}` can be extremely challenging and bug-prone because we've _lost_ our constraints, and we have no information at compile time as to what kinds of data we're dealing with.

Often developers have to use reflection to implement these *ahem* generic functions, which is usually painful and can hurt the performance of your program.

## Our own test helpers with generics

Ideally, we don't want to have to make specific `AssertX` functions for every type we ever deal with. We'd like to be able to have _one_ `AssertEqual` function that works with _any_ type but does not let you compare [apples and oranges](https://en.wikipedia.org/wiki/Apples_and_oranges).

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

We're using `comparable` because we want to describe to the compiler that we wish to use the `==` and `!=` operators - we want to compare! If you try changing the type to `any`,

```go
func AssertNotEqual[T any](got, want T) {
```

You'll get the following error:

```
prog.go2:15:5: cannot compare got != want (operator != not defined for T)
```

Which makes a lot of sense, because you can't use those operators on every (or `any`) type.

## Next: Generic data types

We're going to create a [stack](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)) data type. Stacks should be fairly straightforward to understand from a requirements point of view. They're a collection of items where you can `Push` items to the "top" and to get items back again you `Pop` items from the top (LIFO - last in, first out).

For the sake of brevity I've omitted the TDD process that arrived me at the [following code](https://go2goplay.golang.org/p/HghXymv1OKm) for a stack of `int`s, and a stack of `string`s.

```go
package main

import (
	"log"
)

type StackOfInts struct {
	values []int
}

func (s *StackOfInts) Push(value int) {
	s.values = append(s.values, value)
}

func (s *StackOfInts) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *StackOfInts) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}

type StackOfStrings struct {
	values []string
}

func (s *StackOfStrings) Push(value string) {
	s.values = append(s.values, value)
}

func (s *StackOfStrings) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *StackOfStrings) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}

func main() {
	// INT STACK

	myStackOfInts := new(StackOfInts)

	// check stack is empty
	AssertTrue(myStackOfInts.IsEmpty())

	// add a thing, then check it's not empty
	myStackOfInts.Push(123)
	AssertFalse(myStackOfInts.IsEmpty())

	// add another thing, pop it back again
	myStackOfInts.Push(456)
	value, _ := myStackOfInts.Pop()
	AssertEqual(value, 456)
	value, _ = myStackOfInts.Pop()
	AssertEqual(value, 123)
	AssertTrue(myStackOfInts.IsEmpty())

	// STRING STACK

	myStackOfStrings := new(StackOfStrings)

	// check stack is empty
	AssertTrue(myStackOfStrings.IsEmpty())

	// add a thing, then check it's not empty
	myStackOfStrings.Push("one two three")
	AssertFalse(myStackOfStrings.IsEmpty())

	// add another thing, pop it back again
	myStackOfStrings.Push("four five six")
	strValue, _ := myStackOfStrings.Pop()
	AssertEqual(strValue, "four five six")
	strValue, _ = myStackOfStrings.Pop()
	AssertEqual(strValue, "one two three")
	AssertTrue(myStackOfStrings.IsEmpty())
}

func AssertTrue(thing bool) {
    if thing {
        log.Printf("PASSED: Expected thing to be true and it was\n")
    } else {
        log.Fatalf("FAILED: expected true but got false")
    }
}

func AssertFalse(thing bool) {
    if !thing {
        log.Printf("PASSED: Expected thing to be false and it was\n")
    } else {
        log.Fatalf("FAILED: expected false but got true")
    }
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

### Problems

- The code for both `StackOfStrings` and `StackOfInts` is almost identical. Whilst duplication isn't always the end of the world, this doesn't feel great and does add an increased maintenance cost.
- As we're duplicating the logic across two types, we've had to duplicate the tests too.

We really want to capture the _idea_ of a stack in one type, and have one set of tests for them. We should be wearing our refactoring hat right now which means we should not be changing the tests because we want to maintain the same behaviour.

Pre-generics, this is what we _could_ do

```go
type StackOfInts = Stack
type StackOfStrings = Stack

type Stack struct {
	values []interface{}
}

func (s *Stack) Push(value interface{}) {
	s.values = append(s.values, value)
}

func (s *Stack) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		var zero interface{}
		return zero, false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}
```

- We're aliasing our previous implementations of `StackOfInts` and `StackOfStrings` to a new unified type `Stack`
- We've removed the type safety from the `Stack` by making it so `values` is a [slice](https://github.com/quii/learn-go-with-tests/blob/main/arrays-and-slices.md) of `interface{}`

... And our tests still pass. Who needs generics?

### The problem with throwing out type safety

The first problem is the same as we saw with our `AssertEquals` - we've lost type safety. I can now `Push` apples onto a stack of oranges.

Even if we have the discipline not to do this, the code is still unpleasant to work with because when methods **return `interface{}` they are horrible to work with**.

Add the following test,

```go
myStackOfInts.Push(1)
myStackOfInts.Push(2)
firstNum, _ := myStackOfInts.Pop()
secondNum, _ := myStackOfInts.Pop()
AssertEqual(firstNum+secondNum, 3)
```

You get a compiler error, showing the weakness of losing type-safety:

```go
prog.go2:59:14: invalid operation: operator + not defined for firstNum (variable of type interface{})
```

When `Pop` returns `interface{}` it means the compiler has no information about what the data is and therefore severely limits what we can do. It can't know that it should be an integer, so it does not let us use the `+` operator.

To get around this, the caller has to do a [type assertion](https://golang.org/ref/spec#Type_assertions) for each value.

```go
myStackOfInts.Push(1)
myStackOfInts.Push(2)
firstNum, _ := myStackOfInts.Pop()
secondNum, _ := myStackOfInts.Pop()

// get our ints from out interface{}
reallyFirstNum, ok := firstNum.(int)
AssertTrue(ok) // need to check we definitely got an int out of the interface{}

reallySecondNum, ok := secondNum.(int)
AssertTrue(ok) // and again!

AssertEqual(reallyFirstNum+reallySecondNum, 3)
```

The unpleasantness radiating from this test would be repeated for every potential user of our `Stack` implementation, yuck.

### Generic data structures to the rescue

Just like you can define generic arguments to functions, you can define generic data structures.

Here's our new `Stack` implementation, featuring a generic data type and the tests, showing them working how we'd like them to work, with full type-safety. ([Full code listing here](https://go2goplay.golang.org/p/xAWcaMelgQV))

```go
package main

import (
    "log"
)

type Stack[T any] struct {
    values []T
}

func (s *Stack[T]) Push(value T) {
    s.values = append(s.values, value)
}

func (s *Stack[T]) IsEmpty() bool {
    return len(s.values)==0
}

func (s *Stack[T]) Pop() (T, bool) {
    if s.IsEmpty() {
        var zero T
        return zero, false
    }

    index := len(s.values) -1
    el := s.values[index]
    s.values = s.values[:index]
    return el, true
}

func main() {
    myStackOfInts := new(Stack[int])

    // check stack is empty
    AssertTrue(myStackOfInts.IsEmpty())

    // add a thing, then check it's not empty
    myStackOfInts.Push(123)
    AssertFalse(myStackOfInts.IsEmpty())

    // add another thing, pop it back again
    myStackOfInts.Push(456)
    value, _ := myStackOfInts.Pop()
    AssertEqual(value, 456)
    value, _ = myStackOfInts.Pop()
    AssertEqual(value, 123)
    AssertTrue(myStackOfInts.IsEmpty())

    // can get the numbers we put in as numbers, not untyped interface{}
    myStackOfInts.Push(1)
    myStackOfInts.Push(2)
    firstNum, _ := myStackOfInts.Pop()
    secondNum, _ := myStackOfInts.Pop()
    AssertEqual(firstNum+secondNum, 3)
}
```

You'll notice the syntax for defining generic data structures is consistent with defining generic arguments to functions.

```go
type Stack[T any] struct {
    values []T
}
```

It's _almost_ the same as before, it's just that what we're saying is the **type of the stack constrains what type of values you can work with**.

Once you create a `Stack[Orange]` or a `Stack[Apple]` the methods defined on our stack will only let you pass in and will only return the particular type of the stack you're working with:

```go
func (s *Stack[T]) Pop() (T, bool) {
```

You can imagine the types of implementation being somehow generated for you, depending on what type of stack you create:

```go
func (s *Stack[Orange]) Pop() (Orange, bool) {
```

```go
func (s *Stack[Apple]) Pop() (Apple, bool) {
```

Now that we have done this refactoring, we can safely remove the string stack test because we don't need to prove the same logic over and over.

Using a generic data type we have:

- Reduced duplication of important logic.
- Made `Pop` return `T` so that if we create a `Stack[int]` we in practice get back `int` from `Pop`; we can now use `+` without the need for type assertion gymnastics.
- Prevented misuse at compile time. You cannot `Push` oranges to an apple stack.

## Wrapping up

This chapter should have given you a taste of generics syntax, and some ideas as to why generics might be helpful. We've written our own `Assert` functions which we can safely re-use to experiment with other ideas around generics, and we've implemented a simple data structure to store any type of data we wish, in a type-safe manner.

### Generics are simpler than using `interface{}` in most cases

If you're inexperienced with statically-typed languages, the point of generics may not be immediately obvious, but I hope the examples in this chapter have illustrated where the Go language isn't as expressive as we'd like. In particular using `interface{}` makes your code:

- Less safe (mix apples and oranges), requires more error handling
- Less expressive, `interface{}` tells you nothing about the data
- More likely to rely on [reflection](https://github.com/quii/learn-go-with-tests/blob/main/reflection.md), type-assertions etc which makes your code more difficult to work with and more error prone as it pushes checks from compile-time to runtime

Using statically typed languages is an act of describing constraints. If you do it well you create code that is not only safe and simple to use but also simpler to write because the possible solution space is smaller.

Generics gives us a new way to express constraints in our code, which as demonstrated will allow us to consolidate and simplify code that is not possible to do today.

### Will generics turn Go into Java?

- No.

There's a lot of [FUD (fear, uncertainty and doubt)](https://en.wikipedia.org/wiki/Fear,_uncertainty,_and_doubt) in the Go community about generics leading to nightmare abstractions and baffling code bases. This is usually caveatted with "they must be used carefully".

Whilst this is true, it's not especially useful advice because this is true of any language feature.

When you define your own interfaces you are describing constraints (to use this function, give me something with this method) just like you do with generics. It is possible to make poor design decisions when you do this, generics are not unique in this respect.

### You're already using generics

When you consider that if you've used arrays, slices or maps; you've _already been a consumer of generic code_.

```go
var myApples []Apples
// You cant do this!
append(myApples, Orange{})
```

### Abstraction is not a dirty word

It's easy to dunk on [AbstractSingletonProxyFactoryBean](https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/aop/framework/AbstractSingletonProxyFactoryBean.html) but let's not pretend a code base with no abstraction at all isn't also bad. It's your job to _gather_ related concepts when appropriate, so your system is easier to understand and change; rather than being a collection of disparate functions and types with a lack of clarity.

### [Make it work, make it right, make it fast](https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast#:~:text=%22Make%20it%20work%2C%20make%20it,to%20DesignForPerformance%20ahead%20of%20time.)

People run in to problems with generics when they're abstracting too quickly without enough information to make good design decisions.

The TDD cycle of red, green, refactor means that you have more guidance as to what code you _actually need_ to deliver your behaviour, **rather than imagining abstractions up front**; but you still need to be careful.

There's no hard and fast rules here but resist making things generic until you can see that you have a useful generalisation. When we created the various `Stack` implementations we importantly started with _concrete_ behaviour like `StackOfStrings` and `StackOfInts` backed by tests. From our _real_ code we could start to see real patterns, and backed by our tests, we could explore refactoring toward a more general-purpose solution.

People often advise you to only generalise when you see the same code three times, which seems like a good starting rule of thumb.

A common path I've taken in other programming languages has been:

- One TDD cycle to drive some behaviour
- Another TDD cycle to exercise some other related scenarios

> Hmm, these things look similar - but a little duplication is better than coupling to a bad abstraction

- Sleep on it
- Another TDD cycle

> OK, I'd like to try to see if I can generalise this thing. Thank goodness I am so smart and good-looking because I use TDD, so I can refactor whenever I wish, and the process has helped me understand what behaviour I actually need before designing too much.

- This abstraction feels nice! The tests are still passing, and the code is simpler
- I can now delete a number of tests, I've captured the _essence_ of the behaviour and removed unnecessary detail
