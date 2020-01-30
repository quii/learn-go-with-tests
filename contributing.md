# Contributing

Contributions are very welcome. I hope for this to become a great home for guides of how to learn Go by writing tests. Consider submitting a PR or creating an issue which you can do [here](https://github.com/quii/learn-go-with-tests).

## What we're looking for

* Teaching Go features \(e.g things like `if`, `select`, structs, methods, etc\).
* Showcase interesting functionality within the standard library. Show off how easy it is to TDD a HTTP server for instance.
* Show how Go's tooling, like benchmarking, race detectors, etc can help you arrive at great software.

If you don't feel confident to submit your own guide, submitting an issue for something you want to learn is still a valuable contribution.

### ⚠️ Get feedback quickly for new content ⚠️

- TDD teaches us to work iteratively and get feedback and I strongly suggest you do the same if you wish to contribute
    - Open a PR with your first test and implementation, discuss your approach so I can offer feedback and course correct
- This is of course open-source but I do have strong opinions on the content. The sooner you talk to me the better.

## Style guide

* Always be reinforcing the TDD cycle. Take a look at the [Chapter Template](template.md).
* Emphasis on iterating over functionality driven by tests. The Hello, world example works well because we gradually make it more sophisticated and learning new techniques _driven_ by the tests. For example:
  * `Hello()` &lt;- how to write functions, return types.
  * `Hello(name string)` &lt;- arguments, constants.
  * `Hello(name string)` &lt;- default to "world" using `if`.
  * `Hello(name, language string)` &lt;- `switch`.
* Try and minimise the surface area of required knowledge.
  * Thinking of examples that showcase what you're trying to teach without confusing the reader with other features is important.
  * For example you can learn about `struct`s without understanding pointers.
  * Brevity is king.
* Follow the [Code Review Comments style guide](https://github.com/golang/go/wiki/CodeReviewComments). It's important for a consistent style across all the sections.
* Your section should have a runnable application at the end \(e.g `package main` with a `main` func\) so users can see it in action and play with it.
* All tests should pass.
* Run `./build.sh` before raising PR.
