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
example clock, so let's think about what the important parameters are going to be.

```
<line x1="150" y1="150" x2="114.150000" y2="132.260000"
        style="fill:none;stroke:#000;stroke-width:7px;"/>
```

The centre of the clock (the attributes `x1` and `y1` for this line) is the same
for each hand of the clock. The numbers that need to change for each hand of the
clock - the parameters to whatever builds the SVG - are the `x2` and `y2`
attributes. We'll need an X and a Y for each of the hands of the clock.

Another thing to note about SVGs: the origin - point (0,0) - is at the _top left_ hand corner, not the _bottom left_ as we might expect.

I _could_ think about more parameters - the radius of the clockface circle, the size of the SVG, the colours of the hands, their shape, etc... but it's better to start off by solving a simple, concrete problem with a simple, concrete solution, and then to start adding parameters to make it generalised.

So we'll say that every clock has a centre of (150,150), and that the hour hand is 50 long, the minute hand is 80 long and the second hand is 90 long.

Finally, I'm not deciding _how_ to construct the SVG - we could use a template
from the [`text/template`][texttemplate] package, or we could just send bytes into
a `bytes.Buffer` or a writer. But we know we'll need those numbers, so let's
focus on testing something that creates them.

## Write the test first

So my first test looks like this:

```go
package clockface_test

import (
	"testing"
	"time"

	"github.com/gypsydave5/learn-go-with-tests/math/v1/clockface"
)

func TestSecondHandAtMidnight(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := clockface.Point{X: 150, Y: 150 - 90}
	got := clockface.SecondHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
```

Remember how SVGs start with 0 from the top of the Y axis? To place the second hand at midnight we say that it hasn't moved from the centre of the clockface on the X axis - still 150 - and the Y axis is the length of the hand 'up' from the centre; 150 minus 90.

## Try to run the test

This drives out the expected failures around the missing functions and types:

```
--- FAIL: TestSecondHandAtMidnight (0.00s)
# github.com/gypsydave5/learn-go-with-tests/math/v1/clockface_test [github.com/gypsydave5/learn-go-with-tests/math/v1/clockface.test]
./clockface_test.go:13:10: undefined: clockface.Point
./clockface_test.go:14:9: undefined: clockface.SecondHand
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v1/clockface [build failed]
```

So a `Point` where the tip of the hand should go, and a function to get it.

## Write the minimal amount of code for the test to run and check the failing test output

Let's implement those types to get the code to compile

```go
package clockface

import "time"

type Point struct {
	X float64
	Y float64
}

func SecondHand(t time.Time) Point {
	return Point{}
}
```

and now we get

```
--- FAIL: TestSecondHandAtMidnight (0.00s)
    clockface_test.go:17: Got {0 0}, wanted {150 60}
FAIL
exit status 1
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v1/clockface	0.006s
```

## Write enough code to make it pass

When we get the expected failure, we can fill in the return value of `HandsAt`:

```go
func SecondHand(t time.Time) Point {
	return Point{150, 60}
}
```

Behold, a passing test.

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v1/clockface	0.006s
```

## Refactor

No need to refactor yet - there's barely enough code!

## Repeat for new requirements

We probably need to do some work here that doesn't just involve returning
a clock that shows midnight for every time...

## Write the test first

```go
func TestSecondHandAt30Seconds(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)

	want := clockface.Point{X: 150, Y: 150 + 90}
	got := clockface.SecondHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
```

Same idea, but now the second hand is pointing _downwards_ so we _add_ the
length to the Y axis.

This will compile... but how do we make it pass?

## Thinking time

How are we going to solve this problem?

Every minute the second hand goes through the same 60 states, pointing in 60
different directions. When it's 0 seconds it points to the top of the clockface,
when it's 30 seconds it points to the bottom of the clockface. Easy enough.

So if I wanted to think about in what direction the second hand was pointing at,
say, 37 seconds, I'd want the angle between 12 o'clock and 37/60ths around the
circle. In degrees this is `(360 / 60 ) * 37 = 222`, but it's easier just to
remember that it's `37/60` of a complete rotation.

But the angle is only half the story; we need to know the X and Y coordinate
that the tip of the second hand is pointing at. How can we work that out?

## Math

Imagine a circle with a radius of 1 drawn around the origin - the coordinate `0,
0`.

![picture of the unit circle](math/images/unit_circle.png)

This is called the 'unit circle' because... well, the radius is 1 unit!

The circumference of the circle is made of points on the grid - more
coordinates. The x and y components of each of these coordinates form
a triangle, the hypotenuse of which is always 1 - the radius of the circle

![picture of the unit circle with a point defined on the circumference](math/images/unit_circle_coords.png)

Now, trigonometry will let us work out the lengths of X and Y for each triangle
if we know the angle they make with the origin. The X coordinate will be cos(a),
and the Y coordinate will be sin(a), where a is the angle made between the line
and the (positive) x axis.

![picture of the unit circle with the x and y elements of a ray defined as cos(a) and sin(a) respectively, where a is the angle made by the ray with the x axis](math/images/unit_circle_params.png)

(If you don't believe this, [go and look at Wikipedia...][circle])

One final twist - because we want to measure the angle from 12 o'clock rather
than from the X axis (3 o'clock), we need to swap the axis around; now
x = sin(a) and y = cos(a).

![unit circle ray defined from by angle from y axis](math/images/unit_circle_12_oclock.png)

So now we know how to get the angle of the second hand (1/60th of a circle for
each second) and the X and Y coordinates. We'll need functions for both `sin`
and `cos`.

## `math`

Happily the Go `math` package has both, with one small snag we'll need to get
our heads around; if we look at the description of [`math.Cos`][mathcos]:

> Cos returns the cosine of the radian argument x.

It wants the angle to be in radians. So what's a radian? Instead of defining the full turn of a circle to be made up of 360 degrees, we define a full turn as being 2π radians. There are good reasons to do this that we won't go in to.[^2]

Now that we've done some reading, some learning and some thinking, we can write
our next test.

## Write the test first

All this maths is hard and confusing. I'm not confident I understand what's
going on - so let's write a test! We don't need to solve the whole problem in
one go - let's start off with working out the correct angle, in radians, for the
second hand at a particular time.

I'm going to write these tests _within_ the `clockface` package; they may never
get exported, and they may get deleted (or moved) once I have a better grip on
what's going on.

I'm also going to _comment out_ the acceptance test that I was working on while
I'm working on these tests - I don't want to get distracted by that test while
I'm getting this one to pass.

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

## Try to run the test

```
# github.com/gypsydave5/learn-go-with-tests/math/v2/clockface [github.com/gypsydave5/learn-go-with-tests/math/v2/clockface.test]
./clockface_test.go:12:9: undefined: secondsInRadians
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v2/clockface [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
func secondsInRadians(t time.Time) float64 {
	return 0
}
```

```
--- FAIL: TestSecondsInRadians (0.00s)
    clockface_test.go:15: Wanted 3.141592653589793 radians, but got 0
FAIL
exit status 1
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v2/clockface	0.007s
```

## Write enough code to make it pass

```go
func secondsInRadians(t time.Time) float64 {
	return math.Pi
}
```

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v2/clockface	0.011s
```

## Refactor

Nothing needs refactoring yet

## Repeat for new requirements

Now we can extend the test to cover a few more scenarios. I'm going to skip
forward a bit and show some already refactored test code - it should be clear
enough how I got where I want to.

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
FAIL
exit status 1
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface	0.007s
```

Time to implement all of that maths stuff we were talking about above:

```go
func secondsInRadians(t time.Time) float64 {
	return float64(t.Second()) * (math.Pi / 30)
}
```

One second is (2π / 60) radians... cancel out the 2 and we get π/30 radians.
Multiply that by the number of seconds (as a `float64`) and we should now have
all the tests passing...


```
--- FAIL: TestSecondsInRadians (0.00s)
    --- FAIL: TestSecondsInRadians/00:00:30 (0.00s)
        clockface_test.go:24: Wanted 3.141592653589793 radians, but got 3.1415926535897936
FAIL
exit status 1
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface	0.006s
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

<!---
v3 here
-->

## Repeat for new requirements

So we've got the first part covered here - we know what angle the second hand
will be pointing at in radians. Now we need to work out the coordinates.

Again, let's keep this as simple as possible and only work with the _unit
circle_; the circle with a radius of 1. This means that our hands will all have
a length of one but, on the bright side, it means that the maths will be easy
for us to swallow.

## Write the test first

```go
func TestSecondHandVector(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 30), Point{0, -1}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if got != c.point {
				t.Fatalf("Wanted %v Point, but got %v", c.point, got)
			}
		})
	}
}
```

## Try to run the test

```
# github.com/gypsydave5/learn-go-with-tests/math/v4/clockface [github.com/gypsydave5/learn-go-with-tests/math/v4/clockface.test]
./clockface_test.go:40:11: undefined: secondHandPoint
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v4/clockface [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
func secondHandPoint(t time.Time) Point {
	return Point{}
}
```

```
--- FAIL: TestSecondHandPoint (0.00s)
    --- FAIL: TestSecondHandPoint/00:00:30 (0.00s)
        clockface_test.go:42: Wanted {0 -1} Point, but got {0 0}
FAIL
exit status 1
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v4/clockface	0.010s
```

## Write enough code to make it pass

```go
func secondHandPoint(t time.Time) Point {
	return Point{0, -1}
}
```

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v4/clockface	0.007s
```

## Repeat for new requirements

```go
func TestSecondHandPoint(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 30), Point{0, -1}},
		{simpleTime(0, 0, 45), Point{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if got != c.point {
				t.Fatalf("Wanted %v Point, but got %v", c.point, got)
			}
		})
	}
}
```

## Try to run the test

```
--- FAIL: TestSecondHandPoint (0.00s)
    --- FAIL: TestSecondHandPoint/00:00:45 (0.00s)
        clockface_test.go:43: Wanted {-1 0} Point, but got {0 -1}
FAIL
exit status 1
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v4/clockface	0.006s
```

## Write enough code to make it pass

Remember our unit circle picture?

![picture of the unit circle with the x and y elements of a ray defined as cos(a) and sin(a) respectively, where a is the angle made by the ray with the x axis](math/images/unit_circle_params.png)

We now want the equation that produces X and Y. Let's write it into seconds:

```go
func secondHandPoint(t time.Time) Point {
    angle := secondsInRadians(t)
	x := math.Sin(angle)
	y := math.Cos(angle)

	return Point{x, y}
}
```
Now we get


```
--- FAIL: TestSecondHandPoint (0.00s)
    --- FAIL: TestSecondHandPoint/00:00:30 (0.00s)
        clockface_test.go:43: Wanted {0 -1} Point, but got {1.2246467991473515e-16 -1}
    --- FAIL: TestSecondHandPoint/00:00:45 (0.00s)
        clockface_test.go:43: Wanted {-1 0} Point, but got {-1 -1.8369701987210272e-16}
FAIL
exit status 1
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v4/clockface	0.007s
```

Wait, what (again)? Looks like we've been cursed by the floats once again - both
of those unexpected numbers are _infinitessimal_ - way down at the 16th decimal
place. So again we can either choose to try to increase precision, or to just
say that they're roughly equal and get on with our lives.

One option to increase the accuracy of these angles would be to use the rational
type `Rat` from the `math/big` package. But given the objective is to draw an
SVG and not brain surgery I think we can live with a bit of fuzziness.

```go
func TestSecondHandPoint(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 30), Point{0, -1}},
		{simpleTime(0, 0, 45), Point{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if !roughlyEqualPoint(got, c.point) {
				t.Fatalf("Wanted %v Point, but got %v", c.point, got)
			}
		})
	}
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b Point) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}
```

We've defined two functions to define approximate equality between two `Points`
- they'll work if the X and Y elements are within 0.0000001 of each other.
  That's still pretty accurate.

and now we get

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v4/clockface	0.007s
```

## Refactor

I'm still pretty happy with this.

<!--
v4 ends
--->



## Write the test first
## Try to run the test
## Write the minimal amount of code for the test to run and check the failing test output
## Write enough code to make it pass
## Refactor

## Repeat for new requirements
## Wrapping up

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

OK - time for the back up plan from above: let's define angle equality to a precision that's 'good enough' for a clock face, because what's one one-hundred-quadrillionth between friends? A good way of expressing this is to say that the difference between the two numbers is less than some small number.

```go
func roughlyEqual(a, b float64) bool {
	float64EqualityThreshold := 1e-7
	return math.Abs(a-b) < float64EqualityThreshold
}
```

`1e-7` is a floating point literal, meaning `1 x 10⁻⁷` or `0.0000001`. We're saying that the two numbers shouldn't differ by more than this.

And look - another useful `math` function: `math.Abs` returns the absolute value of a number - or in other words is gets rid of the minus sign if it's present. This is a good way of not having to worry about whether `a` is bigger or smallerthan `b`.

Now the tests pass.

### Hours

```go
func TestHoursInRadians(t *testing.T) {
	cases := []struct {
		time  time.Time
		angle float64
	}{
		{simpleTime(6, 0, 0), math.Pi},
		{simpleTime(18, 0, 0), math.Pi},
		{simpleTime(0, 0, 0), 0},
		{simpleTime(12, 0, 0), 0},
		{simpleTime(3, 0, 0), math.Pi / 2},
		{simpleTime(7, 0, 0), (math.Pi / 6) * 7},
		{simpleTime(0, 30, 0), (math.Pi / 60)},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := hoursInRadians(c.time)
			if !roughlyEqual(got, c.angle) {
				t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
			}
		})
	}
}
```

```
# github.com/gypsydave5/learn-go-with-tests/math/v3/clockface [github.com/gypsydave5/learn-go-with-tests/math/v3/clockface.test]
./clockface_test.go:68:11: undefined: hoursInRadians
```


```go
func hoursInRadians(t time.Time) float64 {
	return 0
}
```

```
--- FAIL: TestHoursInRadians (0.00s)
    --- FAIL: TestHoursInRadians/06:00:00 (0.00s)
        clockface_test.go:70: Wanted 3.141592653589793 radians, but got 0
    --- FAIL: TestHoursInRadians/18:00:00 (0.00s)
        clockface_test.go:70: Wanted 3.141592653589793 radians, but got 0
    --- FAIL: TestHoursInRadians/03:00:00 (0.00s)
        clockface_test.go:70: Wanted 4.71238898038469 radians, but got 0
    --- FAIL: TestHoursInRadians/07:00:00 (0.00s)
        clockface_test.go:70: Wanted 0.7330382858376184 radians, but got 0
    --- FAIL: TestHoursInRadians/00:30:00 (0.00s)
        clockface_test.go:70: Wanted 0.05235987755982989 radians, but got 0
```

```go
func hoursInRadians(t time.Time) float64 {
	seconds := secondsInRadians(t) / (60 * 60)
	minutes := minutesInRadians(t) / 60
	hours := math.Pi / (6 / (float64(t.Hour() % 12)))
	return seconds + minutes + hours
}
```

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface	0.005s
```

### Refactor

I prefer the way the `hoursInRadians` function looks - it's more readable. So I think I'll refactor `minutesInRadians` to look the same.

```go
func minutesInRadians(t time.Time) float64 {
	seconds := secondsInRadians(t) / 60
	minutes := math.Pi / (30 / float64(t.Minute()))
	return seconds + minutes
}
```

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface	0.006
```

## Next Test

So we've got the first part covered here - we know what angle the hands will be pointing at in radians. Now we need to work out their coordinates.

Again, let's keep this as simple as possible and only work with the _unit circle_; the circle with a radius of 1. This means that our hands will all have a length of one but, on the bright side, it means that the maths will be easy for us to swallow.

```go
func TestSecondHandVector(t *testing.T) {
	cases := []struct {
		time   time.Time
		vector Vector
	}{
		{simpleTime(0, 0, 30), Vector{0, -1}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandVector(c.time)
			if got != c.vector {
				t.Fatalf("Wanted %v Vector, but got %v", c.vector, got)
			}
		})
	}
}
```

```
# github.com/gypsydave5/learn-go-with-tests/math/v3/clockface [github.com/gypsydave5/learn-go-with-tests/math/v3/clockface.test]
./clockface_test.go:86:11: undefined: secondHandVector
FAIL	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface [build failed]
go: exit 2
```

```go
func secondHandVector(t time.Time) Vector {
	return Vector{}
}
```


```
--- FAIL: TestSecondHandVector (0.00s)
    --- FAIL: TestSecondHandVector/00:00:30 (0.00s)
        clockface_test.go:88: Wanted {0 -1} Vector, but got {0 0}
FAIL
```

```go
func secondHandVector(t time.Time) Vector {
	return Vector{0, -1}
}
```

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface	0.005s
```

```go
	cases := []struct {
		time   time.Time
		vector Vector
	}{
		{simpleTime(0, 0, 30), Vector{0, -1}},
		{simpleTime(0, 0, 45), Vector{-1, 0}},
	}
```


```
--- FAIL: TestSecondHandVector (0.00s)
    --- FAIL: TestSecondHandVector/00:00:45 (0.00s)
        clockface_test.go:89: Wanted {-1 0} Vector, but got {0 -1}
FAIL
```

Phew... now we can do some work!

## Make it pass

Remember our unit circle picture?

![picture of the unit circle with the x and y elements of a ray defined as cos(a) and sin(a) respectively, where a is the angle made by the ray with the x axis](math/images/unit_circle_params.png)

We now want the equation that produces X and Y. Let's write it into seconds:

```go
func secondHandVector(t time.Time) Vector {
    angle := secondsInRadians(t)
	x := math.Sin(angle)
	y := math.Cos(angle)

	return Vector{x, y}
}
```

We run this and we get

```
./clockface.go:45:16: cannot use x (type float64) as type int in field value
./clockface.go:45:19: cannot use y (type float64) as type int in field value
FAIL
```

Whoops! Looks like I'll need to update the `Vector` type...

```go
type Vector struct {
	X float64
	Y float64
}
```


Now we get

```
--- FAIL: TestSecondHandVector (0.00s)
    --- FAIL: TestSecondHandVector/00:00:30 (0.00s)
        clockface_test.go:89: Wanted {0 -1} Vector, but got {1.2246467991473515e-16 -1}
    --- FAIL: TestSecondHandVector/00:00:45 (0.00s)
        clockface_test.go:89: Wanted {-1 0} Vector, but got {-1 -1.8369701987210272e-16}
```

Wait, what (again)? Looks like we've been cursed by the floats once again - both of those unexpected numbers are _infinitessimal_ - way down at the 16th decimal place. So again we can either choose to try to increase precision, or to just say that they're roughly equal and get on with our lives.

One option to increase the accuracy of these angles would be to use the rational type `Rat` from the `math/big` package. But given the objective is to draw an SVG and not brain surgery I think we can live with a bit of fuzziness.

```go
func TestSecondHandVector(t *testing.T) {
	cases := []struct {
		time   time.Time
		vector Vector
	}{
		{simpleTime(0, 0, 30), Vector{0, -1}},
		{simpleTime(0, 0, 45), Vector{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandVector(c.time)
			if !roughlyEqualVector(got, c.vector) {
				t.Fatalf("Wanted %v Vector, but got %v", c.vector, got)
			}
		})
	}
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualVector(a, b Vector) bool {
	return roughlyEqualFloat64(a.X, b.X) &&
		roughlyEqualFloat64(a.Y, b.Y)
}
```

and so

```
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface	0.006s
```

Let's throw a few more test cases into the test

```go
func TestSecondHandVector(t *testing.T) {
	cases := []struct {
		time   time.Time
		vector Vector
	}{
		{simpleTime(0, 0, 30), Vector{0, -1}},
		{simpleTime(0, 0, 45), Vector{-1, 0}},
		{simpleTime(0, 0, 15), Vector{1, 0}},
		{simpleTime(0, 0, 0), Vector{0, 1}},
		{simpleTime(0, 0, 5), Vector{0.5, 0.5 * math.Sqrt(3)}},
		{simpleTime(0, 0, 10), Vector{0.5 * math.Sqrt(3), 0.5}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandVector(c.time)
			if !roughlyEqualVector(got, c.vector) {
				t.Fatalf("Wanted %v Vector, but got %v", c.vector, got)
			}
		})
	}
}
```

```go
PASS
ok  	github.com/gypsydave5/learn-go-with-tests/math/v3/clockface	0.006s
```



[^1]: This is a lot easier than writing a name out by hand as a string and then having to keep it in sync with the actual time. Believe me you don't want to do that...

[^2]: In short it makes it easier to do calculus with circles as π just keeps coming up as an angle if you use normal degrees, so if you count your angles in πs it makes all the equations simpler.

[texttemplate]: https://golang.org/pkg/text/template/
[circle]: https://en.wikipedia.org/wiki/Sine#Unit_circle_definition
[mathcos]: https://golang.org/pkg/math/#Cos
[floatingpoint]: https://0.30000000000000004.com/
