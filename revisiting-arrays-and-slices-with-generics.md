# Revisiting arrays and slices with generics (DRAFT)

**[The code for this chapter is a continuation from Arrays and Slices, found here](https://github.com/quii/learn-go-with-tests/tree/main/arrays)**

Take a look at both `SumAll` and `SumAllTails` that we wrote in [arrays and slices](arrays-and-slices.md). If you don't have your version please copy the code from the [arrays and slices](arrays-and-slices.md) chapter along with the tests.

```go
// Sum calculates the total from a slice of numbers.
func Sum(numbers []int) int {
	var sum int
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// SumAllTails calculates the sums of all but the first number given a collection of slices.
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
}
```

Do you see a recurring pattern?

- Create some kind of "initial" value, or accumulator.
- Iterate over the collection, applying some kind of operation (or function) to the accumulator and the next item in the slice.
- Return the accumulator.

This idea is commonly talked about in functional programming circles, often times called 'reduce' or [fold](https://en.wikipedia.org/wiki/Fold_(higher-order_function)).

> In functional programming, fold (also termed reduce, accumulate, aggregate, compress, or inject) refers to a family of higher-order functions that analyze a recursive data structure and through use of a given combining operation, recombine the results of recursively processing its constituent parts, building up a return value. Typically, a fold is presented with a combining function, a top node of a data structure, and possibly some default values to be used under certain conditions. The fold then proceeds to combine elements of the data structure's hierarchy, using the function in a systematic way.

Go has always had higher-order functions, and as of version 1.18 it also has [generics](./generics.md), so it is now possible to define some of these functions discussed in our wider field. There's no point burying your head in the sand, this is a very common abstraction outside the Go ecosystem and it'll be beneficial to understand it.

Now, I know some of you are probably cringing at this.

> Go is supposed to be simple

**Don't conflate easiness, with simplicity**. Doing loops and copy-pasting code is easy, but it's not necessarily simple. For more on simple vs easy, watch [Rich Hickey's masterpiece of a talk - Simple Made Easy](https://www.youtube.com/watch?v=SxdOUGdseq4).

**Don't conflate unfamiliarity, with complexity**. Fold/reduce may initially sound scary and computer-sciencey but all it really is, is an abstraction over a very common operation. Taking a collection, and combining it into one item. When you step back, you'll realise you probably do this _a lot_.

## A generic refactor

A mistake people often make with shiny new language features is they start by using them without having a concrete use-case. They rely on conjecture and guesswork to guide their efforts.

Thankfully we've written our "useful" functions and have tests around them, so now we are free to experiment with ideas in the refactoring stage of TDD and know that whatever we're trying, has a verification of its value via our unit tests.

Using generics as a tool for simplifying code via the refactoring step is far more likely to guide you to useful improvements, rather than premature abstractions.

We are safe to try things out, re-run our tests, if we like the change we can commit. If not, just revert the change. This freedom to experiment is one of the truly huge values of TDD.

You should be familiar with the generics syntax [from the previous chapter](generics.md), try and write your own `Reduce` function and use it inside `Sum` and `SumAllTails`.

### Hints

- You only need to work with one type, but one will return `int`, the other `[]int`.
- If you think about the arguments to your function first, it'll give you a very small set of valid solutions
  - The array you want to reduce
  - Some kind of combining function

"Reduce" is an incredibly well documented pattern, there's no need to re-invent the wheel. [Read the wiki, in particular the lists section](https://en.wikipedia.org/wiki/Fold_(higher-order_function)), it should prompt you for another argument you'll need.

> In practice, it is convenient and natural to have an initial value

### My reduce function

```go
func Reduce[A any](collection []A, accumulator A, f func(A, A) A) A {
	for _, x := range collection {
		accumulator = f(accumulator, x)
	}
	return accumulator
}
```

Reduce captures the _essence_ of the pattern, it's a function that takes a collection, an initial value and a combining function, and returns a single value. There's no messy distractions around concrete types.

If you understand generics syntax, you should have no problem understanding what this function does. By using the recognised term `Reduce`, programmers from other languages understand the intent too.

### The usage

```go
// Sum calculates the total from a slice of numbers.
func Sum(numbers []int) int {
	add := func(acc, x int) int { return acc + x }
	return Reduce(numbers, 0, add)
}

// SumAllTails calculates the sums of all but the first number given a collection of slices.
func SumAllTails(numbers ...[]int) []int {
	sumTail := func(acc, x []int) []int {
		if len(x) == 0 {
			return append(acc, 0)
		} else {
			tail := x[1:]
			return append(acc, Sum(tail))
		}
	}

	return Reduce(numbers, []int{}, sumTail)
}
```

`Sum` and `SumAllTails` now describe the behaviour of their computations as the functions declared on their first lines respectively. The act of running the computation on the collection is abstracted away in `Reduce`.

## Further applications of reduce

Using tests we can play around with our reduce function to see how re-usable it is. I have copied over our generic assertion functions from the previous chapter.

```go
func TestReduce(t *testing.T) {
	t.Run("multiplication of all elements", func(t *testing.T) {
		multiply := func(x, y int) int {
			return x * y
		}

		AssertEqual(t, Reduce([]int{1, 2, 3}, 1, multiply), 6)
	})

	t.Run("concatenate strings", func(t *testing.T) {
		concatenate := func(x, y string) string {
			return x + y
		}

		AssertEqual(t, Reduce([]string{"a", "b", "c"}, "", concatenate), "abc")
	})
}
```

### The zero value

In the multiplication example, we show the reason for having a default value as an argument to `Reduce`. If we relied on Go's default value of 0, we'd multiply our initial value by 0, and then the following ones, so you'd only ever get 0. By setting it to 1, the first element in the slice will stay the same, and the rest will multiply by the next elements.

If you wish to sound clever with your nerd friends, you'd call this [The Identity Element](https://en.wikipedia.org/wiki/Identity_element).

> In mathematics, an identity element, or neutral element, of a binary operation operating on a set is an element of the set which leaves unchanged every element of the set when the operation is applied.

In addition, the identity element is 0.

`1 + 0 = 1`

With multiplication, it is 1.

`1 * 1 = 1`

## What if we wish to reduce into a different type from `A`?

Suppose we had a list of transactions `Transaction` and we wanted a function that would take them plus a name to figure out their bank balance.

Let's follow the TDD process.

## Write the test first

```go
func TestBadBank(t *testing.T) {
	transactions := []Transaction{
		{
			From: "Chris",
			To:   "Riya",
			Sum:  100,
		},
		{
			From: "Adil",
			To:   "Chris",
			Sum:  25,
		},
	}

	AssertEqual(t, BalanceFor(transactions, "Riya"), 100)
	AssertEqual(t, BalanceFor(transactions, "Chris"), -75)
	AssertEqual(t, BalanceFor(transactions, "Adil"), -25)
}
```

## Try to run the test
```
# github.com/quii/learn-go-with-tests/arrays/v8 [github.com/quii/learn-go-with-tests/arrays/v8.test]
./bad_bank_test.go:6:20: undefined: Transaction
./bad_bank_test.go:18:14: undefined: BalanceFor
```

## Write the minimal amount of code for the test to run and check the failing test output

We don't have our types or functions yet, add them to make the test run.

```go
type Transaction struct {
	From string
	To   string
	Sum  float64
}

func BalanceFor(transactions []Transaction, name string) float64 {
	return 0.0
}
```

When you run the test you should see the following:

```
=== RUN   TestBadBank
    bad_bank_test.go:19: got 0, want 100
    bad_bank_test.go:20: got 0, want -75
    bad_bank_test.go:21: got 0, want -25
--- FAIL: TestBadBank (0.00s)
```

## Write enough code to make it pass

Let's write the code as if we didn't have a `Reduce` function first.

```go
func BalanceFor(transactions []Transaction, name string) float64 {
	var balance float64
	for _, t := range transactions {
		if t.From == name {
			balance -= t.Sum
		}
		if t.To == name {
			balance += t.Sum
		}
	}
	return balance
}
```

## Refactor

At this point, have some source control discipline and commit your work. We have working software, ready to challenge Monzo, Barclays, et al.

Now our work is committed, we are free to play around with it, and try some different ideas out in the refactoring phase. To be fair, the code we have isn't exactly bad, but for the sake of this exercise, I want to demonstrate the same code using `Reduce`.

```go
func BalanceFor(transactions []Transaction, name string) float64 {
	adjustBalance := func(acc float64, t Transaction) float64 {
		if t.From == name {
			return acc - t.Sum
		}
		if t.To == name {
			return acc + t.Sum
		}
		return acc
	}
	return Reduce(transactions, 0.0, adjustBalance)
}
```

But this won't compile.

```
./bad_bank.go:19:35: type func(acc float64, t Transaction) float64 of adjustBalance does not match inferred type func(Transaction, Transaction) Transaction for func(A, A) A
```

The reason is we're trying to reduce to a _different_ type than the type of the collection. This sounds scary, but actually just requires us to adjust the type signature of `Reduce` to make it work. We won't have to change the function body, and we won't have to change any of our existing callers.

```go
func Reduce[A, B any](collection []A, accumulator B, f func(B, A) B) B {
	for _, x := range collection {
		accumulator = f(accumulator, x)
	}
	return accumulator
}
```

We've added a second type constraint which has allowed us to loosen the constraints on `Reduce`, whilst keeping it type-safe. This makes it more general-purpose and reusable.

### Fold/reduce are pretty universal

The possibilities are endless™️ with `Reduce` (or `Fold`). It's a common pattern for a reason, it's not just for arithmetic or string concatenation. Try a few other applications.

- Why not mix some `color.RGBA` into a single colour?
- Total up the number of votes in a poll, or items in a shopping basket.
- More or less anything involving processing a list.

## Find

Now that Go has generics, combining them with higher-order-functions, we can reduce a lot of boilerplate code within our projects.

No longer do you need to write specific `Find` functions for each type of collection you want to search, instead re-use or write a `Find` function. If you understood the `Reduce` function above, writing a `Find` function will be trivial.

Here's a test

```go
func TestFind(t *testing.T) {
	t.Run("find first even number", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		firstEvenNumber, found := Find(numbers, func(x int) bool {
			return x%2 == 0
		})
		AssertTrue(t, found)
		AssertEqual(t, firstEvenNumber, 2)
	})
}
```

And here's the implementation

```go
func Find[A any](items []A, predicate func(A) bool) (value A, found bool) {
	for _, v := range items {
		if predicate(v) {
			return v, true
		}
	}
	return
}
```

Again, because it takes a generic type, we can re-use it in many ways

```go
type Person struct {
    Name string
}

t.Run("Find the best programmer", func(t *testing.T) {
    people := []Person{
        Person{Name: "Kent Beck"},
        Person{Name: "Martin Fowler"},
        Person{Name: "Chris James"},
    }

    king, found := Find(people, func(p Person) bool {
        return strings.Contains(p.Name, "Chris")
    })

    AssertTrue(t, found)
    AssertEqual(t, king, Person{Name: "Chris James"})
})
```

As you can see, this code is flawless.

## Wrapping up

When done tastefully, higher-order functions like these will make your code simpler to read and maintain, but remember the rule of thumb:

Use the TDD process to drive out real, specific behaviour that you actually need, in the refactoring stage you then _might_ discover some useful abstractions to help tidy the code up.

Practice combining TDD with good source control habits. Commit your work when your test is passing, _before_ trying to refactor. This way if you make a mess, you can easily get yourself back to your working state.

### Names matter

Make an effort to do some research outside of Go, so you don't re-invent patterns that already exist with an already established name.

Writing a function takes a collection of `A` and converts them to `B`? Don't call it `Convert`, that's [`Map`](https://en.wikipedia.org/wiki/Map_(higher-order_function)). Using the "proper" name for these items will reduce the cognitive burden for others and make it more search engine friendly to learn more.

### This doesn't feel idiomatic?

Try to have an open-mind.

Whilst the idioms of Go won't, and shouldn't _radically_ change due to generics being released, they _will_ change - due to the language changing!

Discuss with your colleagues patterns and style of code based on their merits rather than dogma. So long as you have well-designed tests, you'll always be able to refactor and shift things as you understand what works well for you, and your team.

### Resources

Fold is a real fundamental in computer science. Here's some interesting resources if you wish to dig more into it
- [Wikipedia: Fold](https://en.wikipedia.org/wiki/Fold)
- [A tutorial on the universality and expressiveness of fold](http://www.cs.nott.ac.uk/~pszgmh/fold.pdf)
