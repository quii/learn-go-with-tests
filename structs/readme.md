# Structs - WIP

The first requirement we have is to write a `Perimeter(width, height int)` function, which will calculate the perimeter of a square given a width and height.

The TDD cycle should be pretty familiar to you by now.

## Write the test first

```go
func TestPerimeter(t *testing.T) {
	got := Perimeter(10, 10)
	want := 40

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

## Try and run the test

`./shapes_test.go:6:9: undefined: Perimeter`

## Write the minimal amount of code for the test to run and check the failing test output

```go
func Perimeter(width int, height int) (perimeter int) {
	return
}
```

Results in `shapes_test.go:10: got 0 want 40`

## Write enough code to make it pass

```go
func Perimeter(width int, height int) (perimeter int) {
	return 2 * (width + height)
}
```

So far, so easy. Now you need to create a function called `Area(width, height int)` which returns the area of a rectangle.

Try and do it yourself, following the TDD cycle.

You should end up with tests like this

```go
func TestPerimeter(t *testing.T) {
	got := Perimeter(10, 10)
	want := 40

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestArea(t *testing.T) {
	got := Area(12, 6)
	want := 72

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

And code like this

```go
func Perimeter(width int, height int) (perimeter int) {
	return 2 * (width + height)
}

func Area(width int, height int) (area int) {
	return width * height
}
```

## Refactor

So far we have been talking about rectangles a lot but it's not reflected much in our code. We pass width and height into our functions, but they could be the width and height of a circle. 

We of course could possible name our functions more specifically but instead we could define our own _type_ called `Rectangle` which encapsulates this concept for us. We can then use that type as an argument to our functions instead.

A struct is just a named collection of fields where you can store data.

Declare a struct like this

```go
type Rectangle struct {
	width int
	height int
}
```

Now let's refactor the tests to use `Rectangle` instead of plain `int`s.

```go
func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10, 10}
	got := Perimeter(rectangle)
	want := 40

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestArea(t *testing.T) {
	rectangle := Rectangle{12, 6}
	got := Area(rectangle)
	want := 72

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

Remember to run your tests before attempting to fix, you should get a helpful error like

```
./shapes_test.go:7:18: not enough arguments in call to Perimeter
	have (Rectangle)
	want (int, int)
```

You can access the fields of a struct with the syntax of `myStruct.field`. 

Change the two functions to fix the test.

```go
func Perimeter(rectangle Rectangle) (perimeter int) {
	return 2 * (rectangle.width + rectangle.height)
}

func Area(rectangle Rectangle) (area int) {
	return rectangle.width * rectangle.height
}
```

I hope you'll agree that passing a `Rectangle` to a function conveys our intent more clearly but there are more benefits of using structs that we will get on to.