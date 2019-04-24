# Mathematics

For all the power of modern computers to perform huge sums at lightning speed,
the average developer rarely uses any maths to do their job.  But not in this
example! Today we'll use mathematics to solve a _real_ problem. And not boring
mathematics - we're going to use trigonometry and vectors and all sorts of stuff
that you always said you'd never have to use after highschool.

## The Problem

You want to make an SVG of a clock. Not a digital clock - no, that
would be easy - an _analogue_ clock, with hands. You're not looking for anything
fancy, just a nice function that takes a `Time` from the `time` package and
spits out an SVG of a clock with all the hands - hour, minute and second
- pointing in the right direction. How hard can that be?

First we're going to need an SVG of a clock for us to play with. SVGs are
a fantastic image format to manipulate programmatically because they're written
as a series of shapes, described in XML. So this clock:

![an svg of a clock](math/example_clock.svg)

Is described like this:

```xml
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="100%"
     height="100%"
     viewBox="0 0 300 300"
     version="2.0">

  <!-- bezel -->
  <circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>

  <!-- hour hand -->
  <line x1="150" y1="150" x2="114.150000" y2="132.260000"
        style="fill:none;stroke:#000;stroke-width:7px;"/>

  <!-- minute hand -->
  <line x1="150" y1="150" x2="101.290000" y2="99.730000"
        style="fill:none;stroke:#000;stroke-width:7px;"/>

  <!-- second hand -->
  <line x1="150" y1="150" x2="77.190000" y2="202.900000"
        style="fill:none;stroke:#f00;stroke-width:3px;"/>
</svg>
```

It's a circle with three lines, each of the lines starting in the middle of the
circle (x=150, y=150), and ending some distance away.

So what we're going to do is reconstruct the above somehow, but change the lines
so they point in the right directions for some time.

## An Acceptance Test

Before we get too stuck in, lets think about an acceptance test. We've got an
example clock, let's turn it into a template using the
[`text/template`][texttemplate] package:

```svg
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="100%"
     height="100%"
     viewBox="0 0 300 300"
     version="2.0">

  <!-- bezel -->
  <circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>

  <!-- hour hand -->
  <line x1="150" y1="150" x2="{{.Hour.X}}" y2="{{.Hour.Y}}"
        style="fill:none;stroke:#000;stroke-width:7px;"/>
  <!-- minute hand -->
  <line x1="150" y1="150" x2="{{.Minute.X}}" y2="{{.Minute.Y}}"
        style="fill:none;stroke:#000;stroke-width:7px;"/>
  <!-- second hand -->
  <line x1="150" y1="150" x2="{{.Second.X}}" y2="{{.Second.Y}}"
        style="fill:none;stroke:#f00;stroke-width:3px;"/>
</svg>
```

Although I'm ignorant of how I'm going to achieve it, I know that I'm going to
have to give a coordinate each of where the hour, minute and second hands are
going to point to. So I'm imagining a data structure with each hand in it -
`Hour`, `Minute` and `Second` - and each of those hands having a property of `X`
and `Y`. So if I hand this template the right data structure I'll get a clock
showing whatever time I like.

Now I could write a test that builds an SVG and compares it to another SVG, but
this would be (a) boring, (b) time consuming and (c) fragile. What would be much
better would be to test this data structure that I want to pass to the
template - it's the thing that's doing all of the work after all.

So my first test looks like this:

```go
package clockface_test

import (
	"testing"
	"time"

	"github.com/gypsydave5/learn-go-with-tests/math/v1/clockface"
)

func HandsAtMidnigthTest(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := Hands{
		Hour:   Vector{X: 0, Y: 150},
		Minute: Vector{X: 0, Y: 150},
		Second: Vector{X: 0, Y: 150},
	}

	got := clockface.HandsAt(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
```

Which drives out the expected failures

```
# github.com/gypsydave5/learn-go-with-tests/math/v1/clockface_test [github.com/gypsydave5/learn-go-with-tests/math/v1/clockface.test]
./acceptance_test.go:13:10: undefined: Hands
./acceptance_test.go:19:9: undefined: clockface.HandsAt
FAIL    github.com/gypsydave5/learn-go-with-tests/math/v1/clockface [build failed]
```

We want a type called `Hands` which describes where the end of each hand is
sitting as a coordinate. I'm thinking of this coordinate as a `Vector` - when we
think of the centre of the clock as the origin of the hand, each hand becomes a
vector from that point - a pair of numbers indicating a direction and a
magnitude - i.e. how long the hand is and where it's pointing.

This test might need some refining as we discover more about what we're trying
to achieve, but it's a good start.

Let's implement those types to get the code to compile

```go
package clockface

import "time"

type Hands struct {
	Hour   Vector
	Minute Vector
	Second Vector
}

type Vector struct {
	X int
	Y int
}

func HandsAt(t time.Time) (hands Hands) {
	return
}
```

When we get the expected failure, we can fill in the return value of `HandsAt`:

```go
func HandsAt(t time.Time) Hands {
	return Hands{
		Hour:   Vector{X: 0, Y: 150},
		Minute: Vector{X: 0, Y: 150},
		Second: Vector{X: 0, Y: 150},
	}
}
```

To make it pass, and then supply another test to force us to actually do some
work:

```go
func TestHandsAtSixOclock(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 6, 0, 0, 0, time.UTC)

	want := clockface.Hands{
		Hour:   clockface.Vector{X: 0, Y: -150},
		Minute: clockface.Vector{X: 0, Y: 150},
		Second: clockface.Vector{X: 0, Y: 150},
	}

	got := clockface.HandsAt(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
```

## Thinking without tests

Let's not rush in and write another test yet. Let's do some thinking first. How
are we going to solve this problem?

Let's start with the second hand - it's the easiest hand to reason about. Every
minute it goes through the same 60 states, pointing in 60 different
directions. When it's 0 seconds it points to the top of the clockface, when it's
30 seconds it points to the bottom of the clockface. Easy enough.

So if I wanted to think about in what direction the second hand was pointing at,
say, 37 seconds, I'd want the angle between 12 o'clock and 37/60ths around the
circle. In degrees this is `(360 / 60 ) * 37 = 222`, but it's easier just to
remember that it's `37/60` of a complete rotation.

But the angle is only half the story; we need to know the X and Y coordinate
that the tip of the second hand is pointing at. How can we work that out?

## Math

Imagine a circle with a radius of 1 drawn around the origin - the coordinate `0,
0`.

![](#todo-circle-picture)

This is called the 'unit circle' because... well, the radius is 1 unit!

The circumference of the circle is made of points on the grid - more
coordinates. The x and y components of each of these coordinates form
a triangle, the hypotenuse of which is always 1 - the radius of the circle

![](#todo-circle-picture-triangle)

Now, trigonometry will let us work out the lengths of X and Y for each triangle
if we know the angle they make with the origin. The X coordinate will be cos(a),
and the Y coordinate will be sin(a), where a is the angle made between the line
and the (positive) x axis.

![](#todo-circle-with-maths)

(If you don't believe this, [go and look at Wikipedia...][circle])

One final twist - because we want to measure the angle from 12 o'clock rather
than from the X axis (3 o'clock), we need to swap the axis around; now
x = sin(a) and y = cos(a).

So now we know how to get the angle of the second hand (1/60th of a circle for
each second) and the X and Y coordinates. We'll need functions for both `sin`
and `cos`.

## `math`

Happily the Go `math` package has both, with one small snag we'll need to get
our heads around; if we look at the description of [`math.Cos`][mathcos]:

> Cos returns the cosine of the radian argument x.

It wants the angle to be in radians. Instead of breaking a circle up into 360
degrees as we might be more used to, we break the full turn of the circle into
2π radians. There are good reasons to do this that we won't go in to.

Now that we've done some reading, some learning and some thinking, we can write
our next test.

## Second Test

All this maths is hard and confusing. I'm not confident I understand what's
going on - so let's write a test! We don't need to solve the whole problem in
one go - let's start off with working out the correct angle, in radians, for the
second hand at a particular time.

```go
package clockface

import (
	"math"
	"testing"
	"time"
)

func TestSecondsInRadians(t *testing.T) {
	thirtySeconds := time.Date(312, time.October, 28, 0, 0, 30, 0, time.UTC)
	want := math.Pi
	got := secondsInRadians(thirtySeconds)

	if want != got {
		t.Fatalf("Wanted %v radians, but got %v", want, got)
	}
}
```

Here we're testing that 30 seconds past the minute should put the second hand at
halfway around the clock - and it's Our first use of the `math` package;
f a full turn of a circle is 2π radians, we know that halfway round should just
be π radians. `math.Pi` provides us with a value for π.

A dumb implementation:


```go
func secondsInRadians(t time.Time) float64 {
	return math.Pi
}
```

Now we can extend the test to cover a few more scenarios.

```go
func TestSecondsInRadians(t *testing.T) {
	cases := []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(0, 0, 30), math.Pi},
		{simpleTime(0, 0, 0), 0},
		{simpleTime(0, 0, 45), (math.Pi / 2) * 3},
		{simpleTime(0, 0, 7), (math.Pi / 30) * 7},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondsInRadians(c.time)
			if got != c.angle {
				t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
			}
		})
	}
}
```

I added a couple of helper functions to make writing this table based test
a little less tedious. `testName` converts a time into a digital watch
format (HH:MM:SS), and `simpleTime` constructs a `time.Time` using only the
parts we actually care about (again, hours, minutes and seconds).[^1]

```go
func simpleTime(hours, minutes, seconds int) time.Time {
	return time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}
```

These two functions should help make these tests (and future tests) a little
easier to write and maintain.

This gives us some nice test output:

```
--- FAIL: TestSecondsInRadians (0.00s)
    --- FAIL: TestSecondsInRadians/00:00:00 (0.00s)
        clockface_test.go:24: Wanted 0 radians, but got 3.141592653589793
    --- FAIL: TestSecondsInRadians/00:00:45 (0.00s)
        clockface_test.go:24: Wanted 4.71238898038469 radians, but got 3.141592653589793
    --- FAIL: TestSecondsInRadians/00:00:07 (0.00s)
        clockface_test.go:24: Wanted 0.7330382858376184 radians, but got 3.141592653589793
```

Time to implement all of that maths stuff we were talking about above:

```go
func secondsinradians(t time.time) float64 {
	return float64(t.second()) * (math.pi / 30)
}
```

One second is (2π / 60) radians... cancel out the 2 and we get π/30 radians.
Multiply that by the number of seconds (as a `float64`) and we should now have
all the tests passing...

```
--- FAIL: TestSecondsInRadians (0.00s)
    --- FAIL: TestSecondsInRadians/00:00:30 (0.00s)
        clockface_test.go:24: Wanted 3.141592653589793 radians, but got 3.1415926535897936
```

Wait, what?

### Floats are horrible

Floating point arithmetic is [notoriously inaccurate][floatingpoint]. Computers
can only really handle integers, and rational numbers to some extent. Decimal
numbers start to become inaccurate, especially when we factor them up and down
as we are in the `secondsInRadians` function. By dividing `math.Pi` by 30 and
then by multiplying it by 30 we've ended up with a number that's no longer the
same as `math.Pi`.

There are two ways around this:

1. Live with the inaccuracy
2. Refactor our function by refactoring our equation

Now (1) may not seem all that appealing, but it's often the only way to make
floating point equality work. B.eing inaccurate by some infinitessimal fraction
is frankly not going to matter for the purposes of drawing a clockface, so we
could write a function that defines a 'close enough' equality for our angles.
But there's a simple way we can get the accuracy back: we rearrange the equation
so that we're no longer dividing down and then multiplying up. We can do it all
by just dividing.

So instead of `numberOfSeconds * π / 30`, we can write `π / (30
/ numberOfSeconds)`, which is equivalent.

In Go:

```go
func secondsinradians(t time.time) float64 {
	return (math.Pi / (30 / (float64(t.Second()))))
}
```

And we get a pass.

```
PASS
ok      github.com/gypsydave5/learn-go-with-tests/math/v2/clockface     0.005s
```

## More tests

This test has definitely boosted my confidence in the maths I'm doing. I hope we
can see where we're headed now - similar tests for functions called
`minutesInRadians` and `hoursInRadians`

### Minutes

```go
func TestMinutesInRadians(t *testing.T) {
	cases := []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(0, 30, 0), math.Pi},
		{simpleTime(0, 0, 0), 0},
		{simpleTime(0, 45, 0), (math.Pi / 2) * 3},
		{simpleTime(0, 7, 0), (math.Pi / 30) * 7},
		{simpleTime(0, 0, 30), (math.Pi / 60)},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := minutesInRadians(c.time)
			if got != c.angle {
				t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
			}
		})
	}
}
```

```go
func minutesInRadians(t time.Time) float64 {
	seconds := secondsInRadians(t) / 60
	minutes := math.Pi / (30 / (float64(t.Minute())))
	return seconds + minutes
}
```

```sh
--- FAIL: TestMinutesInRadians (0.00s)
    --- FAIL: TestMinutesInRadians/00:00:30 (0.00s)
        clockface_test.go:46: Wanted 0.05235987755982989 radians, but got 0.05235987755982988
```

Again some inaccuracy... one one-hundred-quadrillionth out. Can we fix it the
same way?

```go
func minutesInRadians(t time.Time) float64 {
	return math.Pi / ((30 * 60) / (float64(t.Second()) + (60 * float64(t.Minute()))))
}
```

```sh
--- FAIL: TestMinutesInRadians (0.00s)
    --- FAIL: TestMinutesInRadians/00:00:30 (0.00s)
        clockface_test.go:46: Wanted 0.05235987755982989 radians, but got 0.05235987755982988
```

No. :(

OK - time for the back up plan from above: let's define angle equality to
a precision that's 'good enough' for a clock face, because what's one
one-hundred-quadrillionth between friends? A good way of expressing this
is to say that the difference between the two numbers is less than some small
number.

```go
func roughlyEqual(a, b float64) bool {
	float64EqualityThreshold := 1e-7
	return math.Abs(a-b) < float64EqualityThreshold
}
```

`1e-7` is a floating point literal, meaning `1 x 10⁻⁷` or `0.0000001`. We're
saying that the two numbers shouldn't differ by more than this.

And look - another useful `math` function: `math.Abs` returns the absolute value
of a number - or in other words is gets rid of the minus sign if it's present.
This is a good way of not having to worry about whether `a` is bigger or smaller
than `b`.

[^1]: This is a lot easier than writing a name out by hand as a string and then
  having to keep it in sync with the actual time. Believe me you don't want to
  do that...


[texttemplate]: https://golang.org/pkg/text/template/
[circle]: https://en.wikipedia.org/wiki/Sine#Unit_circle_definition
[mathcos]: https://golang.org/pkg/math/#Cos
[floatingpoint]: https://0.30000000000000004.com/
