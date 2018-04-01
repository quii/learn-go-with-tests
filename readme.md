# Learn Go with Tests

_Or learn test-driven development with Go_

[![Build Status](https://travis-ci.org/quii/learn-go-with-tests.svg?branch=master)](https://travis-ci.org/quii/learn-go-with-tests)

## Why

- Explore the Go language by writing tests
- **Get a grounding with TDD**. Go is a good language for learning TDD because it is a simple language to learn and testing is built in
- Be confident that you'll be able to start writing robust, well tested systems in Go

## Table of contents

If there's no link, it's not done yet! [why not contribute?](contributing.md)

### Go fundamentals

0. todo: Install Go, set up environment for productivity.
1. [Hello, world](/hello-world) - Declaring variables, constants, if/else statements, switch, write your first go program and write your first test. Sub-test syntax and closures.
2. [Integers](/integers) - Further Explore function declaration syntax and learn new ways to improve the documentation of your code.
3. [Iteration](/for) - Learn about `for` and benchmarking.
4. [Arrays and slices](/arrays) - Learn about arrays, slices, `len`, varargs, `range` and test coverage.
5. [Structs, methods & interfaces](/structs) - Learn about `struct`, methods, `interface` and table driven tests.
6. [Pointers & errors (WIP)](/pointers) - Learn about pointers, errors and type aliasing.
7. [Dependency Injection & interfaces (WIP)](/di-and-interfaces) - Learn about dependency injection, how it relates to using interfaces and a primer on io
8. Concurrency (WIP)
9. Errors

Property based tests (todo)

### Standard lib essentials

- `io`
- HTTP client
- HTTP server
- JSON marhshalling
- `time`
- Sorting

## Contributing

- *This project is work in progress* If you would like to contribute, please do get in touch.
- Read [contributing.md](contributing.md) for guidelines
- Any ideas? Create an issue

## Background

I have some experience introducing Go to development teams and have tried different approaches as to how to grow a team from some people curious about Go into highly effective writers of Go systems.

### What didn't work

#### Read _the_ book

An approach we tried was to take [the blue book](https://www.amazon.co.uk/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440) and every week discuss the next chapter along with the exercises.

I love this book but it requires a high level of commitment. The book is very detailed in explaining concepts, which is obviously great but it means that the progress is slow and steady - this is not for everyone.

I found that whilst a small number of people would read chapter X and do the exercises, many people didn't.

#### Solve some problems
Katas are fun but they are usually limited in their scope for learning a language; you're unlikely to use go routines to solve a kata.

Another problem is when you have varying levels of enthusiasm. Some people just learn way more of the language than others and when demonstrating what they have done end up confusing people with featues the others are not familiar with.

This ends up making the learning feel quite _unstructured_ and _ad hoc_.

### What did work

By far the most effective way was by slowly introducing the fundamentals of the language by reading through [go by example](https://gobyexample.com/), exploring them with examples and discussing them as a group. This was a more interactive approach than "read chapter x for homework".

Over time the team gained a solid foundation of the _grammar_ of the language so we could then start to build systems.

This to me seems analogous to practicing scales when trying to learn guitar.

It doesn't matter how artistic you think you are, you are unlikely to write good music without understanding the fundamentals and practicing the mechanics.

### What works for me
When *I* learn a new programming language I usually start by messing around in a REPL but eventually I need more structure.

What I like to do is explore concepts and then solidify the ideas with tests. Tests verify the code I write is correct and documents the feature I have learned.

Taking my experience of learning with a group and my own personal way I am going to try and create something that hopefully proves useful to other teams. Learning the fundamentals by writing small tests so that you can then take your existing software design skills and ship some great systems.

## Who this is for

- People who are interested in picking up Go
- People who already know some Go, but want to explore testing more

## What you'll need

- A computer!
- [Installed Go](https://golang.org/)
- A text editor
- Some experience with programming. Understanding of concepts like `if`, variables, functions etc.
- Comfortable with using the terminal

## Feedback

- Add issues or [tweet me @quii](https://twitter.com/quii)