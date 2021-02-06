# TDD Anti-patterns

This book tries to teach you some good habits in respect to TDD but from time to time it's still necessary to review your own technique and in particular remind yourself the kind of things you should be avoiding. This chapter lists a number of TDD anti-patterns and how to remedy them.

## Not doing it at all

Of course, it is possible to write great software without TDD, but a lot of problems I've seen with the design of code, and the quality of tests would be very difficult to arrive at if a disciplined approach to TDD had been used.

One of the strengths of TDD is it gives you a formal process for you to break down problems, understand what you're trying to achieve (red), get it done (green) and then have a good think about how to make it right (blue/refactor). Without this the process is often quite ad-hoc and loose which _can_ result in not the best engineering.

## Misunderstanding the constraints of the refactoring step

I have been in a number of workshops, mobbing or pairing sessions where someone has made the test pass and is now doing the refactoring stage. After some thought she thinks it would be good to abstract away some code into a new struct, a budding pedant yells:

> You're not allowed to do this! You should write a test for this first, we're doing TDD!

This seems to be a common misunderstanding. **You can do whatever you like to the code when the tests are green**, the only thing you're not allowed to do is **add or change behaviour**.

The point of these tests are to give you the _freedom to refactor_, find the right abstractions and make the code easier to change and understand.

## Having tests that won't fail

It's astonishing how often this comes up. You start debugging or changing some tests and then realise that there are no scenarios where a test can fail, or at least it won't fail the way the test is supposed to be protecting against.

This is _next to impossible_ with TDD if you're following **the first step**

> Write a test, see it fail

This is almost always a function of developers writing tests _after_ code is written and are probably chasing test coverage rather than a useful test suite.

## Useless assertions

Ever worked on a system and you've broken a test and you see this?

> `false was not equal to true`

I already know false is not equal to true. This is not a helpful message, it doesn't tell me what I've broken.

This is another symptom of not following the TDD process and writing tests after the fact.

Going back to the drawing board

> Write a test, see it fail (and don't be ashamed of the error message)

## Not listening to your tests

[Dave Farley in his video "When TDD goes wrong"](https://www.youtube.com/watch?v=UWtEVKVPBQ0&feature=youtu.be) points out

> TDD gives you the fastest feedback possible on your design

From my own experience, a lot of developers are trying to practice TDD but frequently ignore the signals coming back to them from the TDD process. So they're still stuck with fragile, annoying systems with a poor test suite.

Simply put, if testing your code is difficult, then _using_ your code is difficult too. Treat your tests as the first user of your code and then you'll see if your code is pleasant to work with or not.

I've emphasised this a lot in the book, and I'll say it again **listen to your tests**.

### Excessive setup

Ever looked at a test with 20, 50, 100, 200 lines of setup code before anything interesting in the test happens? Do you then have to change the code and then have to revisit this mess and wish you had a different career?

What are the signals here?

Complicated tests == complicated code. Why is your code complicated? Does it have to be?

### Too many test doubles with lots of setup

A specialisation of excessive setup.

- If you have lots of test doubles in your tests that means the code you're testing has lots of dependencies which means your design needs work.
- If your test is reliant on setting up various interactions with mocks that means your code has to do lots of interactions with its dependencies, which you may be able to simplify and consolidate

#### Think about the types of test doubles you use

- Mocks are sometimes helpful, but they're extremely powerful and therefore easy to misuse. Try giving yourself the constraint of using stubs instead.
- Verifying implementation detail with spies is sometimes helpful, but try to avoid it. Remember your implementation detail is usually not important, and you don't want your tests coupled to them. Look to couple your tests to **useful behaviour rather than incidental details**.

#### Consolidate dependencies

Here is some code for a `http.HandlerFunc` to handle new user registrations for a website.

```go
type User struct {
	// Some user fields
}

type UserStore interface {
	CheckEmailExists(email string) (bool, error)
	StoreUser(newUser User) error
}

type Emailer interface {
	SendEmail(to User, body string, subject string) error
}

func NewRegistrationHandler(userStore UserStore, emailer Emailer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// extract out the user from the request body (handle error)
		// check user exists (handle duplicates, errors)
		// store user (handle errors)
		// compose and send confirmation email (handle error)
		// if we got this far, return 2xx response
	}
}
```

At first pass it's reasonable to say the design isn't so bad. It only has 2 dependencies!

Let's re-evaluate this by considering the handler's responsibilities

- Parse the request body into a `User` ✅
- Go to a storage abstraction and check if the user exists ❓
- Go to the storage abstraction and store the user ❓
- Compose an email ❓
- Use an emailer abstraction to send the email ❓
- Return an appropriate http response, depending on success, errors, etc ✅

To exercise this code, you're going to have to write many tests with varying degrees of test double setups, spies, etc

- What if the requirements expand? Translations for the emails? Sending an SMS confirmation too? Does it make sense to you that you have to change a HTTP handler to accommodate this change?
- Does it feel right that the important rule of "we should send an email" resides within a HTTP handler?
- Why do you have to go through the ceremony of creating HTTP requests and reading responses to verify that rule?

**Listen to your tests**. Writing tests for this code in a TDD fashion should quickly make you feel uncomfortable. If it feels painful, stop and think.

What if the design was like this instead?

```go
type UserService interface {
	Register(newUser User) error
}

func NewRegistrationHandler(userService UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// parse user
		// register user
		// check error, send response
	}
}
```

- Simple to test the handler ✅
- Changes to the rules around registration are isolated away from HTTP, so they are also simpler to test ✅

### Lots of assertions

Many assertions can make tests difficult to read and challenging to debug when they fail.

A helpful rule of thumb is to aim to do make one assertion per test.


### Violating encapsulation







- `package foo_test`

## Summary

- Not actually following the TDD process
- Poor design

So learn about good software design!

The good news is TDD can help you learn them because as stated in the beginning:

**TDD's main purpose is to provide feedback on your design.**
