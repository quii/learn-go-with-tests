# Arrays - WIP

Arrays allow you to store multiple elements of the same type in a variable

When you have an array, it is very common to have to iterate over them so let's use our new-found knowledge of `for` to make a `Sum` function. `Sum` will take an array of numbers and return the total.

Let's use our TDD skills

## Write the test first

In `sum_test.go`
```go
func TestSum(t *testing.T) {

	numbers := [5]int{1, 2, 3, 4, 5}

	got := Sum(numbers)
	want := 15

	if want != got {
		t.Errorf("got %d want %d given, %v", got, want, numbers)
	}
}
```

Arrays have a _fixed capacity_ which you define when you declare the variable. It is sometimes useful to also print the inputs to the function in the error message and we are using the `%v` placeholder which is the "default" format, which works well for arrays.

[Read more about the format strings](https://golang.org/pkg/fmt/)

## Try and run the test

By running `go test` the compiler will fail with `./sum_test.go:10:15: undefined: Sum`

## Write the minimal amount of code for the test to run and check the failing test output

In `sum.go`

```go
func Sum(numbers [5]int) (sum int) {
	return
}
```

Your test should now fail with _a clear error message_

`sum_test.go:13: got 0 want 15 given, [1 2 3 4 5]`

## Write enough code to make it pass

```go
func Sum(numbers [5]int) (sum int) {
	for i := 0; i < 5; i++ {
		sum += numbers[i]
	}
	return
}
```

To get the value out of an array at a particular index, just use `array[index]` syntax. In this case we are using `for` to iterate 5 times to work through the array and add each item onto `sum`

#### A note on source control

At this point if you are using source control (which you should!) I would `commit` the code as it is. We have working software backed by a test. 

I _wouldnt_ push to master though, because I plan to refactor next. It is nice to commit at this point in case you somehow get in to a mess with refactoring - you can always go back to the working version.

## Refactor

Let's introduce `range` to help clean up our code

```go
func Sum(numbers [5]int) (sum int) {
	for _, number := range numbers {
		sum += number
	}
	return
}
```

`range` lets you iterate over an array. Every time it is called it returns two values, the index and the value. We are choosing to ignore the index value by using `_`

An interesting property of arrays though is the size is encoded in its type. If you try and pass an `[4]int` into a function that expects `[5]int`, it wont compile.

You may be thinking it's quite cumbersome that arrays are fixed length and most of the time you probably wont be using them! Go has _slices_ which are dynamic in size and most of the time you will probably be using them instead.

The next requirement will be to sum collections of varying sizes

## Write the test first

We will now use the slice type which allows us to have collections of any size. The syntax is very similar to arrays, you just omit the size when declaring them

```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := [5]int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {

		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```
## Try and run the test

This does not compile

`./sum_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Sum`

## Write the minimal amount of code for the test to run and check the failing test output

The problem here is we can either

- Break the existing API by changing the argument to `Sum` to be a slice rather than an array. When we do this we will know we have potentially ruined someone's day because our _other_ test will not compile! 
- Create a new function

In our case, no-one else is using our function so rather than having two functions to maintain lets just have one. 

```go
func Sum(numbers []int) (sum int) {
	for _, number := range numbers {
		sum += number
	}
	return
}
```

If you try and run the tests they will still not compile, you will have to change the first test to pass in a slice rather than an array. 

## Write enough code to make it pass

It turns out that fixing the compiler problems were all we need to do here and the tests pass!

## Refactor

We had already refactored `Sum` and all we've done is change from arrays to slices so there's not a lot to do here but our tests could do with some love. 

```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {

		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```

It is important to question the value of your tests. It should not be a goal to have as many tests as possible, but rather to have as much *confidence* as possible in your code base. Having too many tests can turn in to a real problem and it just adds more overhead in maintenance. 

In our case, you can see that having two tests for this function is redundant. If it works for a slice of one size it's very likely it'll work for any size (within reason).

Go's built-in testing toolkit features a coverage tool, which can help identify areas of your code you have not covered. I do want to stress that having 100% coverage should not be your goal, it's just a tool to give you an idea of your coverage. If you have been strict with TDD it's quite likely you'll have close to 100% coverage anyway.

Try running 

`go test -cover`

You should see 

```
PASS
coverage: 100.0% of statements
```

Now delete one of the tests and check the coverage again.

Now that we are happy we have a well tested function you should commit your great work before taking on the next challenge.

We need a few function called `SumAll` which will take a varying number of slices, returning a new slice containing the totals for each slice pass in. 

For example

`SumAll([]int{1,2}, []int{0,9})` would return `[]int{3, 9}`

or 

`SumAll([]int{1,1,1})` would return `[]int{3}`

## Write the test first

```go
func TestSumAll(t *testing.T)  {

	got := SumAll([]int{1,2}, []int{0,9})
	want := []int{3, 9}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## Try and run the test

`./sum_test.go:23:9: undefined: SumAll`

## Write the minimal amount of code for the test to run and check the failing test output

We need to define SumAll according to what our test wants.

Go can let you write _variadic functions_ that can take a variable number of arguments.

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
	return
}
```

Try and compile but our tests still don't compile! 

`./sum_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

Go does not let you use equality operators with slices. You _could_ write a function to iterate over each `got` and `want` slice and check their values but for convenience sake we can use `reflect.DeepEqual` which is useful for seeing if _any_ two variables are the same.

```go
func TestSumAll(t *testing.T)  {

	got := SumAll([]int{1,2}, []int{0,9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

It's important to note that `reflect.DeepEqual` is not "type safe", the code will compile even if you did something a bit silly. To see this in action, temporarily change the test to:

```go
func TestSumAll(t *testing.T)  {

	got := SumAll([]int{1,2}, []int{0,9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, "bob") {
		t.Errorf("got %v want %v", got, want)
	}
}
```

What we have done here is try to compare a `slice` with a `string`. Which makes no sense, but the test compiles! So while using `reflect.DeepEqual` is a convenient way of comparing slices (and other things) you must be careful when using it.

Change the test back again and run it, you should have test output looking like this

`sum_test.go:30: got [] want [3 9]`

## Write enough code to make it pass

What we need to do is iterate over the varargs, calculate the sum using our `Sum` function from before and then add it to the slice we will return

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
	sums = make([]int, len(numbersToSum))

	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}

	return
}
```

Lots of new things to learn! 

There's a new way to create a slice. `make` allows you to create a slice with a starting capacity of the `len` of the `numbersToSum` we need to work through.

You can index slices like arrays with `mySlice[N]` to get the value out or assign it a new value with `=`

The tests should now pass

## Refactor

As mentioned, slices have a capacity. If you have a slice with a capacity of 2 and try and do `mySlice[10] = 1` you will get a _runtime_ error. 

However you can use the `append` function which takes a slice and a new value, returning a new slice with all the items in it.

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return
}
```

In this implementation we are worrying less about capacity. We start with an empty slice (defined in the function signature) and append to it the result of `Sum` as we work through the varargs. 

`TODO: Something around slicing slices e.g mySlice[:3]`

## `TODO` Wrapping up