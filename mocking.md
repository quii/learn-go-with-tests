# Mocking

We'll next cover _mocking_ and it's relation to DI with a case-study. 

You have been asked to write a program which will count 5 seconds, printing each number on a new line and when it reaches zero it will print "Go!" and exit. 

```
5
4
3
2
1
Go!
```

We'll tackle this by writing a function called `Countdown` which we will then put inside a `main` program so it looks something like this:

```go
package main

func main() {
    Countdown()
}
```

While this is a pretty trivial program, to test it fully we will need as always to take an _iterative_, _test-driven_ approach. 

What do I mean by iterative? We make sure we take the smallest steps we can to have _useful software_. We dont want to spend a long time with code that will theoretically work after some hacking. 

- Print 5 to 1, followed by "Go!" to stdout
- Wait a second between each

Let's just work on the first one, writing a program which prints 5 to 1 to stdout. 

## Write the test first

Our software needs to print to stdout and we saw how we could use DI to facilitate testing this in the DI section.

```go
func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	Countdown(buffer)

	got := buffer.String()
	want := "5"

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
```

If anything like `buffer` is unfamiliar to you, re-read the previous section.

We know we want our `Countdown` function to write data somewhere and `io.Writer` is the de-facto way of capturing that as an interface in Go. 

- In `main` we will send to `os.Stdout` so our users see the countdown printed to the terminal
- In test we will send to `bytes.Buffer` so our tests can capture what data is being generated

Notice that we're just capturing the first thing that should be printed first, it's an important skill to be able to slice up requirements as small as you can so you can have _working software_.

## Try and run the test

`./countdown_test.go:11:2: undefined: Countdown`


## Write the minimal amount of code for the test to run and check the failing test output

Define `Countdown`

```go
func Countdown() {}
```

Try again

```go
./countdown_test.go:11:11: too many arguments in call to Countdown
	have (*bytes.Buffer)
    want ()
```

The compiler is telling you what your function signature could be, so update it.

```go
func Countdown(out *bytes.Buffer) {}
```

`countdown_test.go:17: got '' want '5'`

Perfect!

## Write enough code to make it pass

```go
func Countdown(out *bytes.Buffer) {
	fmt.Fprint(out, "5")
}
```

We're using `fmt.Fprint` which takes an `io.Writer` (like `*bytes.Buffer`) and sends a `string` to it. The test should pass. 

## Refactor

We know that while `*bytes.Buffer` works, it would be better to use a general purpose interface instead.

```go
func Countdown(out io.Writer) {
	fmt.Fprint(out, "5")
}
```

Re-run the tests and they should be passing. 

To complete matters, let's now wire up our function into a `main` so we have some working software to reassure ourselves we're making progress.

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func Countdown(out io.Writer) {
	fmt.Fprint(out, "5")
}

func main() {
	Countdown(os.Stdout)
}
```

Try and run the program and be amazed at your handywork. 

Yes this seems trivial but this approach is what I would recommend for any project. **Take a thin slice of functionality and make it work end-to-end, backed by tests.**

Next we can make it print 4,3,2,1 and then "Go!".

## Write the test first

```go
func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	Countdown(buffer)

	got := buffer.String()
	want := `5
4
3
2
1
Go!`

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
```

The backtick syntax is another way of creating a `string` but lets you put things like newlines which is perfect for our test.

Isn't it nice that by focusing on getting our software working end-to-end that iterating on the requirements is very straightforward. By investing in getting the overall plumbing working right, we can focus on the next requirements easily.

## Try and run the test

```
countdown_test.go:21: got '5' want '5
		4
		3
		2
    1
		Go!'
```
## Write enough code to make it pass

```go
func Countdown(out io.Writer) {
	for i := 5; i > 0; i-- {
		fmt.Fprintln(out, i)
    }
    fmt.Fprint(w, "Go!")
}
```

Use a `for` loop counting backwards with `i--` and use `fmt.Fprintln` to print to `out` with our number followed by a newline character. Finally use `fmt.Fprint` to send "Go!" aftward

## Refactor

There's not much to refactor other than removing some magic values.

```go
const finalWord = "Go!"
const countdownStart = 5

func Countdown(out io.Writer) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
	}
	fmt.Fprint(out, finalWord)
}
```

If you run the program now, you should get the desired output but we dont have it as a dramatic countdown with the 1 second pauses. 

Go let's you achieve this with `time.Sleep`. Try adding it in to our code.

```go
func Countdown(out io.Writer) {
	for i := countdownStart; i > 0; i-- {
		time.Sleep(1 * time.Second)
		fmt.Fprintln(out, i)
	}
	
	time.Sleep(1 * time.Second)
	fmt.Fprint(out, finalWord)
}
```

If you run the program it works as we want it to. The tests still pass, but they now take 6 seconds. 

Not only that, but this seems like an important property of the function that we have not tested. 

## Mocking

We have a dependency on `Sleep`ing which we need to extract so we can then control it in our tests.

We want to assert that after every count we `Sleep` for a second.

If we can _mock_ `time.Sleep` we can use _dependency injection_ to use it instead of a "real" `time.Sleep` and then we can **spy on the calls** to make assertions on them. 

## Write the test first

Let's define our dependency

```go
type Sleeper interface {
	Sleep()
}
```

I made a design decision that our `Countdown` function would not be responsible
for how long the sleep is. This simplifies our code a little for now at least
and means a user of our function can configure that sleepiness however they
like.

Now we need to make a _mock_ of it for our tests to use. 

```go
type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}
```

There's nothing new to learn here, we just need to be able to _spy_ on the calls to our mock so we can check it has been called `N` times.

Update the tests to inject a dependency on our Spy and assert that the sleep has been called 6 times.

```go
func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}
	spySleeper := &SpySleeper{}

	Countdown(buffer, spySleeper)

	got := buffer.String()
	want := `5
4
3
2
1
Go!`

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}

	if spySleeper.Calls != 6 {
		t.Errorf("not enough calls to sleeper, want 6 got %d", spySleeper.Calls)
	}
}
```

## Try and run the test

```
too many arguments in call to Countdown
	have (*bytes.Buffer, Sleeper)
	want (io.Writer)
```

## Write the minimal amount of code for the test to run and check the failing test output

We need to update `Countdown` to accept our `Sleeper`

```go
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		time.Sleep(1 * time.Second)
		fmt.Fprintln(out, i)
	}

	time.Sleep(1 * time.Second)
	fmt.Fprint(out, finalWord)
}
```

If you try again, your `main` will no longer compile for the same reason

```
./main.go:26:11: not enough arguments in call to Countdown
	have (*os.File)
	want (io.Writer, Sleeper)
```

Let's create a _real_ sleeper which implements the interface we need

```go
type ConfigurableSleeper struct {
	duration time.Duration
}

func (o *ConfigurableSleeper) Sleep() {
	time.Sleep(o.duration)
}
```

I decided to make a little extra effort and make it so our real sleeper is
configurable but you could just as easily not bother and hard-code it for
1 second (YAGNI right?).

We can then use it in our real application like so

```go
func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second}
	Countdown(os.Stdout, sleeper)
}
```

## Write enough code to make it pass

The test is now compiling but not passing because we're still calling the `time.Sleep` rather than the injected in dependency. Let's fix that.

```go
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.sleep()
		fmt.Fprintln(out, i)
	}

	sleeper.sleep()
	fmt.Fprint(out, finalWord)
}
```

Now the test should be passing (and no longer taking 6 seconds!).

## Refactor

todo


## But isn't mocking evil?

You may have heard mocking is evil. Just like anything in software development it can be used for evil, just like DRY. 

People normally get in to a bad state when they dont _listen to their tests_ and _not respecting the refactoring stage_. 

If your mocking code is becoming complicated or you are having to mock out lots of things to test something, you should _listen_ to that bad feeling and think about your code. Usually it is a sign of

- The thing you are testing is having to do too many things. 
- Its dependencies are too fine-grained

Normally a lot of mocking points to _bad abstraction_ in your code. 

**What people see here is a weakness in TDD but it is actually a strength**, more often than not poor test code is a result of bad design or put more nicely, well-designed code is easy to test. 