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






[circle]: https://en.wikipedia.org/wiki/Sine#Unit_circle_definition
[mathcos]: https://golang.org/pkg/math/#Cos









## Our next test

Let's write something that lets us work out the angle of the second hand from
a time.

```go

```



[texttemplate]: https://golang.org/pkg/text/template/
