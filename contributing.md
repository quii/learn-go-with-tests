# Contributing

Contributions are very welcome, I hope for this to become a great home for guides of how to learn go by writing tests

## What we're looking for

- Teaching Go features (e.g things like `if`, `select`, structs, methods, etc)
- Showcase interesting functionality within the standard library. Show off how easy it is to TDD a HTTP server for instance.

If you don't feel confident to submit your own guide, submitting an issue for something you want to learn is still a valuable contribution. 

## Style guide

- Always be reinforcing the TDD cycle. Take a look at [template.md](template.md)
- Emphasis on iterating over functionality driven by tests. The Hello, world example works well because we gradually make it more sophisticated and learning new techniques *driven* by the tests. For example: 
    - `Hello()` <- how to write functions, return types.
    - `Hello(name string)` <- arguments, constants
    - `Hello(name string)` <- default to "world" using `if`
    - `Hello(name, language string)` <- `switch`
    
- Try and minimise the surface area of required knowledge. Thinking of examples that showcase what you're trying to teach without confusing the reader with lots of other features is important.
- Follow the [Code Review Comments style guide](https://github.com/golang/go/wiki/CodeReviewComments). It's important for a consistent style across all the sections.
- Your section should have a runnable application at the end (e.g `package main` with a `main` func) so users can see it in action and play with it
- All tests should pass
- Run `./build.sh` before raising PR