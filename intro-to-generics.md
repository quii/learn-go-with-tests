# Intro to generics

(At the time of writing) Go does not have support for user-defined generics, but it has been confirmed that it will be included in version 1.18.

However, there are ways to experiment with the upcoming implementation using the go2go playground _today_.

This chapter will give you a brief introduction to generics, hopefully dispel any reservations you may have about them and give you an idea of how you will be able to simplify some of your code in the future.

The code we write here will be the foundation for future chapters around generics.

## Setting up the playground

In the go2go playground we can't run `go test`, so how are we going to write tests to explore generic code?

The playground _does_ let us execute code, and because we're programmers that means we can work around the lack of a test runner by **making one of our own**.

## Our own test helpers (`AssertEqual`, `AssertNotEqual`)

//TODO: without generics

### Tradeoffs made without generics

- Can't use `==`
- Not typesafe

## Our own test helpers with generics

- Typesafe
- Probably marginally faster due to no reflection
- Simpler code because we've given the compiler more information vs just `interface{}`

## Wrapping up

There's a lot of FUD in the Go community about generics leading to nightmare abstractions and baffling code bases.

This is usually caveatted with "they must be used carefully". Whilst this is true, it's not especially useful advice because this is true of any language feature.

I know this because I have written extremely awful code _without_ generics.

### You're already using generics

The FUD becomes an even sillier statement when you consider that if you've used arrays, slices or maps; you've already been a consumer of generic code.

//TODO, some examples of how it's typesafe

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

It's easy to dunk on AbstractSingletonMethodFactory but let's not pretend a code base with no abstraction isn't just as bad.
