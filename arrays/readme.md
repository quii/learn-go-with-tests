# Arrays

Arrays allow you to store multiple elements of the same type in a variable

When you have an array, it is very common to have to iterate over them so let's use our new-found knowledge of `for` to make a `Sum` function. `Sum` will take an array of numbers and return the total.

Let's use our TDD skills

## 1. Write the test first

In `sum_test.go`
```go
func TestSum(t *testing.T) {

	numbers := [5]int{1, 2, 3, 4, 5}

	expectedSum := 15
	actualSum := Sum(numbers)

	if expectedSum != actualSum {
		t.Errorf("expected the sum to be %d but was %d, given %v", expectedSum, actualSum, numbers)
	}
}
```

Arrays have a _fixed capacity_ which you define when you declare the variable. It is sometimes useful to also print the inputs to the function in the error message and we are using the `%v` placeholder which is the "default" format, which works well for arrays.

[Read more about the format strings](https://golang.org/pkg/fmt/)

## 2. Try and run the test

By running `go test` the compiler will fail with `./sum_test.go:10:15: undefined: Sum`

## 3. Write the minimal amount of code for the test to run and check the failing test output

In `sum.go`

```go
func Sum(numbers [5]int) (sum int) {
	return
}
```

Your test should now fail with _a clear error message_

`sum_test.go:13: expected the sum to be 15 but was 0 given, [1 2 3 4 5]`

## 4. Write enough code to make it pass

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

## 5. Refactor

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