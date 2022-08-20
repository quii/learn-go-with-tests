# Structs, methods & interfaces

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/structs)**

만약 높이와 너비가 주어지고 직사각형의 둘레를 계산하는 기하학 코드가 필요하다고 가정해봅시다. 
그렇다면 `Perimeter(width float64, height float64)`와 같은 코드를 작성할 수 있는데, 여기서 `float64`은 `123.45`와 같은 부동 수소점 숫자입니다

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

Notice the new format string? The `f` is for our `float64` and the `.2` means print 2 decimal places.

## Try to run the test

`./shapes_test.go:6:9: undefined: Perimeter`

## Write the minimal amount of code for the test to run and check the failing test output

```go
func Perimeter(width float64, height float64) float64 {
	return 0
}
```

Results in `shapes_test.go:10: got 0.00 want 40.00`.

## Write enough code to make it pass
정사각형 계산 = `4 * a`

직사각형 계산 = `2 * (a + b)`
```go
func Perimeter(width float64, height float64) float64 {
	return 2 * (width + height)
}
```

So far, so easy. Now let's create a function called `Area(width, height float64)` which returns the area of a rectangle.

Try to do it yourself, following the TDD cycle.

You should end up with tests like this

넓이 계산 = `width * height`

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
func Perimeter(width float64, height float64) float64 {
	return 2 * (width + height)
}

func Area(width float64, height float64) float64 {
	return width * height
}
```

## Refactor

위 코드가 정상적으로 작동은 하지만, 직사각형에 대한 명시적인 설명이 없습니다. 이 상황에서 부주의하고 게으른 개발자라면 잘못된 값을 반환할수도 있다는 생각을 못한채 삼각형과 같은 형태의 넓이와 높이를 반환할 수도 있습니다.

해결법으로 함수 이름을 `RectangleArea`와 같이 조금 더 특정 지을 수 있습니다. 조금 더 나은 방법은 이 컨셉을 캡슐화하는 `Rectangle`이라는 커스텀 타입을 정의할 수 있습니다. 




Go에서는 **구조체(struct)** 이라는 것을 이용해서 단순한 타입을 지정할 수 있다.. [구조체](https://golang.org/ref/spec#Struct_types) 는 다수의 필드를 가지고 있는 이름이 지정된 컬렉션이며 데이터를 저장할 수 있는 공간이다.

구조체는 다음과 같이 선언할 수 있다.

```go
type Rectangle struct {
	Width  float64
	Height float64
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

Remember to run your tests before attempting to fix. The tests should show a helpful error like

```text
./shapes_test.go:7:18: not enough arguments in call to Perimeter
    have (Rectangle)
    want (float64, float64)
```

구조체의 필드에 접근하기 위해서는 `myStruct.field`와 같은 문법을 사용하면 됩니다.

Change the two functions to fix the test.

```go
func Perimeter(rectangle Rectangle) float64 {
	return 2 * (rectangle.Width + rectangle.Height)
}

func Area(rectangle Rectangle) float64 {
	return rectangle.Width * rectangle.Height
}
```

I hope you'll agree that passing a `Rectangle` to a function conveys our intent more clearly, but there are more benefits of using structs that we will cover later.

Our next requirement is to write an `Area` function for circles.

## Write the test first

```go
func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		got := Area(rectangle)
		want := 72.0

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := Area(circle)
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

}
```

위에서 보는 것과 같이, `f`가 `g`로 대체되었음을 볼 수 있습니다. 이는 `g`를 사용하면 더 정확한 소숫점 숫자를 에러메시지에서 확인할 수 있기 때문입니다. ([fmt options](https://golang.org/pkg/fmt/))
. 예를 들어 원 면적 계산에 반지름 `1.5`를 사용한다면, `f`는 `7.068583`을 보여주지만, g는 `7.0685834705770345`를 보여주게 됩니다.

## Try to run the test

`./shapes_test.go:28:13: undefined: Circle`

## Write the minimal amount of code for the test to run and check the failing test output

We need to define our `Circle` type.

```go
type Circle struct {
	Radius float64
}
```

Now try to run the tests again

`./shapes_test.go:29:14: cannot use circle (type Circle) as type Rectangle in argument to Area`

Some programming languages allow you to do something like this:

```go
func Area(circle Circle) float64       {}
func Area(rectangle Rectangle) float64 {}
```

하지만 Go에서는 다음과 같이 재선언이 되었다는 에러가 발생합니다.

`./shapes.go:20:32: Area redeclared in this block`

위 문제를 두 가지 방법으로 해결할 수 있다:

* 동일한 이름의 두 함수가 서로 다른 패키지 안에 선언될 수 있다. 그러므로 `Area(Ciecle)`과 같은 함수를 새로운 패키지에 선언할 수 있지만 위와 같은 상황에서는 너무 지나친 방법입니다.

* [_메서드(methods)_](https://golang.org/ref/spec#Method_declarations) 를 새로 선언된 타입에 사용할 수도 있습니다.

### What are methods?

지금까지는 함수만 작성했지만, 메서드도 사용해왔습니다. 예를 들어 `t(testing.T)`의 인스턴스로써 `t.Errorf`를 호출할 때 메서드 `Errorf`를 호출한 것입니다.

메서드는 `receiver`를 가지는 함수를 뜻합니다.
메서드 선언은 식별자, 메서드 이름을 메서드에 묶고 메서드를 `receiver`의 기본 유형과 연관시킵니다.

메서드는 함수와 매우 비슷하지만 특정 타입의 인스턴스에서 불러져서 호출됩니다. `Area(rectangle)`과 같이 함수는 어디에서나 호출할 수 있지만, 메서드는 `things`에서만 호출이 됩니다.


An example will help so let's change our tests first to call methods instead and then fix the code.

```go
func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		got := rectangle.Area()
		want := 72.0

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := circle.Area()
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

}
```

If we try to run the tests, we get

```text
./shapes_test.go:19:19: rectangle.Area undefined (type Rectangle has no field or method Area)
./shapes_test.go:29:16: circle.Area undefined (type Circle has no field or method Area)
```

> type Circle has no field or method Area

I would like to reiterate how great the compiler is here. It is so important to take the time to slowly read the error messages you get, it will help you in the long run.

## Write the minimal amount of code for the test to run and check the failing test output

Let's add some methods to our types

```go
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return 0
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 0
}
```

메서드와 함수는 매우 비슷하기 때문에 작성하는 문법 또한 비슷합니다. 
유일하게 다른 점은 메서드 리시버는 다음과 같은 문법을 가진다는 점입니다. `func (receiverName ReceiverType) MethodName(args)`

메서드가 타입의 변수와 함께 호출되면, `receiverName` 변수를 통해서 데이터에 대한 참조를 얻게 됩니다. 다른 프로그래밍 언어에서는 해당 방식이 함축적으로 작동하고 `this`라는 리시버를 통해 접근합니다.

Go에서는 리시버 변수의 첫 글자를 타입의 첫 글자로부터 따오는 것이 컨벤션입니다.


```
r Rectangle
```

If you try to re-run the tests they should now compile and give you some failing output.

## Write enough code to make it pass

Now let's make our rectangle tests pass by fixing our new method

```go
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
```

If you re-run the tests the rectangle tests should be passing but circle should still be failing.

To make circle's `Area` function pass we will borrow the `Pi` constant from the `math` package \(remember to import it\).

```go
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
```

## Refactor

현재 테스트에는 다수의 중복이 존재합니다.

이 부분을 shape의 컬렉션을 취하고, `Area()` 메서드를 호출하고 결과를 확인하도록 수정할 수 있습니다.

`Rectangle`과 `Circle`을 전달 받지만 허용되는 모양(shape)가 아닌 다른 것을 전달 받았을 때에서는 컴파일 에러를 반환하는 `checkArea`과 같은 함수를 작성할 수 있습니다.

Go에서 이런 의도가 있을 때 사용할 수 있는 것이 바로 **인터페이스(interfaces)** 입니다.

[인터페이스 (Interfaces)](https://golang.org/ref/spec#Interface_types) 란 여전히 타입적으로 안전하게 유지하면서도 코드로부터 매우 분리되어있고 다른 타입과 사용이 가능한 함수를 작성할 수도 있도록 도와주는 매우 강력한 개념입니다(특히나 Go와 같은 정적 타이핑 언어에서 빛을 발합니다.)

Let's introduce this by refactoring our tests.

```go
func TestArea(t *testing.T) {

	checkArea := func(t testing.TB, shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		if got != want {
			t.Errorf("got %g want %g", got, want)
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

We are creating a helper function like we have in other exercises but this time we are asking for a `Shape` to be passed in. If we try to call this with something that isn't a shape, then it will not compile.

How does something become a shape? We just tell Go what a `Shape` is using an interface declaration

```go
type Shape interface {
	Area() float64
}
```

We're creating a new `type` just like we did with `Rectangle` and `Circle` but this time it is an `interface` rather than a `struct`.

Once you add this to the code, the tests will pass.

### Wait, what?

이런식의 인터페이스 사용은 일반적인 다른 프로그래밍 언어와 매우 다릅니다. 보통은 `My type Foo implements interface Bar`과 같은 코드를 사용하는 것이 일반적인데 말입니다.

하지만 위와 같은 코드에서는
* `Rectangle`은 `float64`를 반환하는 `Area` 메서드를 호출하기 때문에 `Shape` 인터페이스를 만족시킵니다.
* `Circle`은 `float64`를 반환하는 `Area` 메서드를 호출하기 때문에 `Shape` 인터페이스를 만족시킵니다.
* `string`은 메서드를 가지지 않으므로 인터페이스를 만족시키지 않습니다.


In Go **interface resolution is implicit**. If the type you pass in matches what the interface is asking for, it will compile.

### Decoupling

Notice how our helper does not need to concern itself with whether the shape is a `Rectangle` or a `Circle` or a `Triangle`. By declaring an interface, the helper is _decoupled_ from the concrete types and only has the method it needs to do its job.

위 코드에서 shape가 `Rectangle`, `Circle`, 혹은 `Triangle`인지 테스트 헬퍼가 고려하지 않아도 되는 부분을 주목해야합니다. 인터페이스를 선언함으로써, 헬퍼는 콘크리트 타입으로부터 분리되고 작업을 수행하기 위해 필요한 메서드만 가집니다.

이런 방식처럼 인터페이스를 사용하여 필요한 것만 선언하는 방식은 소프트웨어 디자인에서 매우 중요한 부분입니다.

## Further refactoring

Now that you have some understanding of structs we can introduce "table driven tests".

[테이블 주도 테스트 (Table driven tests)](https://github.com/golang/go/wiki/TableDrivenTests) 란 같은 방식으로 테스트해야하는 테스트 케이스의 리스트를 생성해야하는 경우 매우 우용하다.

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
			t.Errorf("got %g want %g", got, tt.want)
		}
	}

}
```

여기에서 추가된 새로운 문법은 단순히 `areaTests`라는 익명 구조체(anonymous struct)를 생성한 것 뿐입니다. `shape`와 `want`라는 두 개의 필드로 구성된 `[]struct`를 사용해 구조체의 슬라이스를 선언하였고 케이스로 슬라이스를 채워주었습니다.

이렇게 구조체 필드를 사용하여 다른 슬라이스와 동일하게 이터레이팅하여 테스틀르 실행할 수 있습니다.

위와 같은 방법을 통해서 새로운 `shape`을 사용하고 `Area`를 구현, 및 테스트 케이스에 추가하는 것이 얼마나 쉬운 방법인지 알  수 있습니다.

You can see how it would be very easy for a developer to introduce a new shape, implement `Area` and then add it to the test cases. In addition, if a bug is found with `Area` it is very easy to add a new test case to exercise it before fixing it.

테이블 주도 테스트는 이렇게 도구로써 좋은 역할을 합니다. 하지만 테스트에 추가적인 잡음이 필요해진다는 것도 인지해야합니다. 인터페이스의 다양한 구현체를 테스트하거나 다양한 요구사항을 가진 함수에 데이터를 전달할 때 테스트하기 좋은 방법입니다.

Let's demonstrate all this by adding another shape and testing it; a triangle.

## Write the test first

Adding a new test for our new shape is very easy. Just add `{Triangle{12, 6}, 36.0},` to our list.

```go
func TestArea(t *testing.T) {

	areaTests := []struct {
		shape Shape
		want  float64
	}{
		{Rectangle{12, 6}, 72.0},
		{Circle{10}, 314.1592653589793},
		{Triangle{12, 6}, 36.0},
	}

	for _, tt := range areaTests {
		got := tt.shape.Area()
		if got != tt.want {
			t.Errorf("got %g want %g", got, tt.want)
		}
	}

}
```

## Try to run the test

Remember, keep trying to run the test and let the compiler guide you toward a solution.

## Write the minimal amount of code for the test to run and check the failing test output

`./shapes_test.go:25:4: undefined: Triangle`

We have not defined `Triangle` yet

```go
type Triangle struct {
	Base   float64
	Height float64
}
```

Try again

```text
./shapes_test.go:25:8: cannot use Triangle literal (type Triangle) as type Shape in field value:
    Triangle does not implement Shape (missing Area method)
```

It's telling us we cannot use a `Triangle` as a shape because it does not have an `Area()` method, so add an empty implementation to get the test working

```go
func (t Triangle) Area() float64 {
	return 0
}
```

Finally the code compiles and we get our error

`shapes_test.go:31: got 0.00 want 36.00`

## Write enough code to make it pass

```go
func (t Triangle) Area() float64 {
	return (t.Base * t.Height) * 0.5
}
```

And our tests pass!

## Refactor

Again, the implementation is fine but our tests could do with some improvement.

When you scan this

```
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

It's not immediately clear what all the numbers represent and you should be aiming for your tests to be easily understood.

So far you've only been shown syntax for creating instances of structs `MyStruct{val1, val2}` but you can optionally name the fields.

Let's see what it looks like

```
        {shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
        {shape: Circle{Radius: 10}, want: 314.1592653589793},
        {shape: Triangle{Base: 12, Height: 6}, want: 36.0},
```

In [Test-Driven Development by Example](https://g.co/kgs/yCzDLF) Kent Beck refactors some tests to a point and asserts:

> The test speaks to us more clearly, as if it were an assertion of truth, **not a sequence of operations**

\(emphasis in the quote is mine\)

Now our tests - rather, the list of test cases - make assertions of truth about shapes and their areas.

## Make sure your test output is helpful

Remember earlier when we were implementing `Triangle` and we had the failing test? It printed `shapes_test.go:31: got 0.00 want 36.00`.

We knew this was in relation to `Triangle` because we were just working with it.
But what if a bug slipped in to the system in one of 20 cases in the table?
How would a developer know which case failed?
This is not a great experience for the developer, they will have to manually look through the cases to find out which case actually failed.

We can change our error message into `%#v got %g want %g`. The `%#v` format string will print out our struct with the values in its field, so the developer can see at a glance the properties that are being tested.

To increase the readability of our test cases further, we can rename the `want` field into something more descriptive like `hasArea`.

One final tip with table driven tests is to use `t.Run` and to name the test cases.

By wrapping each case in a `t.Run` you will have clearer test output on failures as it will print the name of the case

```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
```

And you can run specific tests within your table with `go test -run TestArea/Rectangle`.

Here is our final test code which captures this

```go
func TestArea(t *testing.T) {

	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
			}
		})

	}

}
```

## Wrapping up

This was more TDD practice, iterating over our solutions to basic mathematic problems and learning new language features motivated by our tests.

* Declaring structs to create your own data types which lets you bundle related data together and make the intent of your code clearer
* Declaring interfaces so you can define functions that can be used by different types \([parametric polymorphism](https://en.wikipedia.org/wiki/Parametric_polymorphism)\)
* Adding methods so you can add functionality to your data types and so you can implement interfaces
* Table driven tests to make your assertions clearer and your test suites easier to extend & maintain

This was an important chapter because we are now starting to define our own types. In statically typed languages like Go, being able to design your own types is essential for building software that is easy to understand, to piece together and to test.

Interfaces are a great tool for hiding complexity away from other parts of the system. In our case our test helper _code_ did not need to know the exact shape it was asserting on, only how to "ask" for its area.

As you become more familiar with Go you will start to see the real strength of interfaces and the standard library. You'll learn about interfaces defined in the standard library that are used _everywhere_ and by implementing them against your own types, you can very quickly re-use a lot of great functionality.
