# Structs - WIP

The first requirement we have is to write a `Perimeter(width, height float64)` function, which will calculate the perimeter of a square given a width and height. `float64` is a type like `int` but allows you to add precision like `123.45`

The TDD cycle should be pretty familiar to you by now.

## Write the test first

```go
func TestPerimeter(t *testing.T) {
	got := Perimeter(10.0, 10.0)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

Notice the new format string? The `f` is for our `float64` and the `.2` means print 2 decimal places (todo: eugh word this better)

## Try and run the test

`./shapes_test.go:6:9: undefined: Perimeter`

## Write the minimal amount of code for the test to run and check the failing test output

```go
func Perimeter(width float64, height float64) (perimeter float64) {
	return
}
```

Results in `shapes_test.go:10: got 0 want 40`

## Write enough code to make it pass

```go
func Perimeter(width float64, height float64) (perimeter float64) {
	return 2 * (width + height)
}
```

So far, so easy. Now you need to create a function called `Area(width, height float64)` which returns the area of a rectangle.

Try and do it yourself, following the TDD cycle.

You should end up with tests like this

```go
func TestPerimeter(t *testing.T) {
	got := Perimeter(10.0, 10.0)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	got := Area(12.0, 6.0)
	want := 72.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

And code like this

```go
func Perimeter(width float64, height float64) (perimeter float64) {
	return 2 * (width + height)
}

func Area(width float64, height float64) (area float64) {
	return width * height
}
```

## Refactor

So far we have been talking about rectangles a lot but it's not reflected much in our code. We pass width and height float64o our functions, but they could be the width and height of a circle. 

We of course could possible name our functions more specifically but instead we could define our own _type_ called `Rectangle` which encapsulates this concept for us. We can then use that type as an argument to our functions instead.

A struct is just a named collection of fields where you can store data.

Declare a struct like this

```go
type Rectangle struct {
	width float64
	height float64
}
```

Now let's refactor the tests to use `Rectangle` instead of plain `float64`s.

```go
func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	rectangle := Rectangle{12.0, 6.0}
	got := Area(rectangle)
	want := 72.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

Remember to run your tests before attempting to fix, you should get a helpful error like

```
./shapes_test.go:7:18: not enough arguments in call to Perimeter
	have (Rectangle)
	want (float64, float64)
```

You can access the fields of a struct with the syntax of `myStruct.field`. 

Change the two functions to fix the test.

```go
func Perimeter(rectangle Rectangle) (perimeter float64) {
	return 2 * (rectangle.width + rectangle.height)
}

func Area(rectangle Rectangle) (area float64) {
	return rectangle.width * rectangle.height
}
```

I hope you'll agree that passing a `Rectangle` to a function conveys our float64ent more clearly but there are more benefits of using structs that we will get on to.

Our next requirement is to write an `Area` function for circles.

## Write the test first

```go
func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		got := Area(rectangle)
		want := 72.0

		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := Area(circle)
		want := 314.16

		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	})

}
```

## Try and run the test

`./shapes_test.go:28:13: undefined: Circle`

## Write the minimal amount of code for the test to run and check the failing test output

We need to define our `Circle` type.

```go
type Circle struct {
	radius float64
}
```

Now try and run the tests again

`./shapes_test.go:29:14: cannot use circle (type Circle) as type Rectangle in argument to Area`

Some programming languages allow you to do something like this:

```go
func Area(circle Circle) float64 { ... }
func Area(rectangle Rectangle) float64 { ... }
```

But you cannot in Go

`./shapes.go:20:32: Area redeclared in this block`

We have two choices

- You can have functions with the same name declared in different _packages_. So we could create our `Area(Circle)` in a new package, but that feels overkill here
- We can define _methods_ on our newly defined types instead. 

### What are methods?

So far we have only been writing *functions* but we have been using some methods. When we call `t.Errof` we are calling the method `ErrorF` on the instance of our `t` (`testing.T`). 

Methods are very similar to functions but they are called by invoking them on an instance of a particular type. Where you can just call functions wherever you like, such as `Area(rectangle)` you can only call methods on "things".

As always, an example will help so let's change our tests first to call methods instead and then fix the code

```go
func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		got := rectangle.Area()
		want := 72.0

		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := circle.Area()
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	})

}
```

If we try to run the tests we get

```
./shapes_test.go:19:19: rectangle.Area undefined (type Rectangle has no field or method Area)
./shapes_test.go:29:16: circle.Area undefined (type Circle has no field or method Area)
```

## Write the minimal amount of code for the test to run and check the failing test output

Let's add some methods to our types

```go
type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() (area float64)  {
	return
}

type Circle struct {
	radius float64
}

func (c Circle) Area() (area float64)  {
	return
}
```

The syntax for declaring methods is almost the same as functions and that's because they're so similar. The only difference is the syntax of the method receiver `func (receiverName RecieverType) MethodName(args)`.

When your method is called on an variable of that type, you get your reference to it's data via the receiverName variable. In many other programming languages this is done implicitly and you access the reciever via `this`.

It is a convention in Go to have the receiver variable be the first letter of the type.

If you try and re-run the tests they should now compile and give you some failing output

## Write enough code to make it pass

Now let's make our rectangle tests pass by fixing our new method

```go
func (r Rectangle) Area() (area float64)  {
	return r.width * r.height
}
```

If you re-run the tests the rectangle tests should be passing but circle should still be failing.

To make circle's `Area` function pass we will borrow a constant from the `math` package

```go
func (c Circle) Area() (area float64)  {
	return math.Pi * c.radius * c.radius
}
```

## Refactor

There is some duplication in our tests. If you zoom out a bit all we want to do is take a collection of _shapes_, call the `Area()` method on them and then check the result. 

Our shapes share a common _interface_: `Area() float64`. 

In Go, if you want to write functions which can be called with different types, like `Rectangle` and `Circle`s that share *the same interface*, you can define your function to say

> I only accept arguments that have methods called `Area` which return `float64`

You do this with `interface`. Interfaces are a very powerful concept in statically typed languages like Go, they allow you to make re-useable functions that can be used with different types and create highly-decoupled code. 

Let's introduce this by refactoring our tests.


```go
func TestArea(t *testing.T) {

	checkArea := func(t *testing.T, shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	}

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		checkArea(t, rectangle, 72.0)
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		checkArea(t, circle, 314.1592653589793)
	})

}
```

We are creating a helper function like we have in other exercises but this time we are asking for a `Shape` to be passed in. If we try and call this with something that isn't a shape then it will not compile.

How does something become a shape? We just tell Go what a `Shape` is using an interface declaration

```go
type Shape interface {
	Area() float64
}
```

Once you add this to the code, the tests will pass. 

Notice how our helper does not need to concern itself with whether the shape is a rectangle or a square or a triangle. By declaring an interface the helper is _decoupled_ from the concrete types and just has the method it needs to do it's job. 

## Further refactoring

Now that you have some understanding of structs we can now introduce "table based tests"

Table based tests are useful when you want to build a list of test cases that can be tested in the same manner.

```go
func TestArea(t *testing.T) {

	areaTests := []struct {
		shape Shape
		want  float64
	}{
		{Rectangle{12, 6}, 72.0},
		{Circle{10}, 314.1592653589793},
	}

	for _, tt := range areaTests {
		got := tt.shape.Area()
		if got != tt.want {
			t.Errorf("got %.2f want %.2f", got, tt.want)
		}
	}

}

```

The only new syntax here is creating an "anonymous struct". We are declaring a slice of structs with two fields, the `shape` and the `want`. Then we fill the array with cases. 

We then iterate over them just like we do any other slice, using the struct fields to run our tests.

You can see how it would be very easy for a developer to introduce a new shape, implement `Area` and then add it to the test cases. In addition if a bug is found with `Area` it is very easy to add a new test case to exercise it before fixing it.

Table based tests can be a great item in your toolbox but be sure that you have a need for the extra noise in the tests. If you wish to test various implementations of an interface, or if the data being passed in to a function has lots of different requirements that need testing then they are a great fit.

## Wrapping up

What we have covered

- Declaring structs
- Adding methods
- Interfaces
- Table based tests 