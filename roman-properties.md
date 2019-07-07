# Roman Properties

In this chapter we are going to cover two subjects:

- A number of companies will ask you to do the "Roman Numeral Kata" as part of the interview process. This chapter will show how you can tackle it with TDD.
- All of the tests we have shown so far could be described as "example based tests", where we provide example scenarios and expectations. _Property based testing_ allows us to _describe the properties of our domain_ and then the tests will exercise those properties against our code. Sounds abstract, but all will become clear!

We are going to write a function which converts an Arabic number to a Roman Numeral. 

If you haven't heard of Roman Numerals, it's how the Romans wrote down numbers. They rely on broadly straight lines which are easier to carve into tablets and stuff!

`MCMLXXXIV` is 1984 for instance.

That's a bit complicated and it's hard to imagine how we can write code to figure this out right from the start but as this book stresses a key skill for software developers is to try and identify "thin vertical slices" of _useful_ functionality and then **iterating**. The TDD workflow shows us a clear path as to how to do this. 

So rather than 1984, let's start with `1`.

## Write the test first

```go
func TestRomanNumerals(t *testing.T) {
	got := ConvertToRoman(1)
	want := "I"
	
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
```

If you've got this far in the book this is hopefully feeling very boring and routine to you. That's a good thing.

## Try to run the test

`./numeral_test.go:6:9: undefined: ConvertToRoman`

Let the compiler guide the way

## Write the minimal amount of code for the test to run and check the failing test output

Create our function but not make it pass yet, always make sure the tests fails how you expect

```go
func ConvertToRoman(arabic int) string {
	return ""
}
```

It should run now

```go
=== RUN   TestRomanNumerals
--- FAIL: TestRomanNumerals (0.00s)
    numeral_test.go:10: got '', want 'I'
FAIL
```

## Write enough code to make it pass

```go
func ConvertToRoman(arabic int) string {
	return "I"
}
```

## Refactor

Not much to refactor yet. 

_I know_ it feels weird just to hard-code the result but with TDD we want to stay out of "red" for as long as possible. It may _feel_ like we haven't accomplished much but we've defined our API and got a test capturing one of our rules; even if the "real" code is pretty dumb. 

Now use that uneasy feeling to write a new test to force us to write slightly less dumb code.

## Write the test first

We can use subtests to group our tests nice

```go
func TestRomanNumerals(t *testing.T) {
	t.Run("1 gets converted to I", func(t *testing.T) {
		got := ConvertToRoman(1)
		want := "I"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}	
	})

	t.Run("2 gets converted to II", func(t *testing.T) {
		got := ConvertToRoman(2)
		want := "II"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})
}
```

## Try to run the test

```
=== RUN   TestRomanNumerals/2_gets_converted_to_II
    --- FAIL: TestRomanNumerals/2_gets_converted_to_II (0.00s)
        numeral_test.go:20: got 'I', want 'II'
```

Not much surprise there

## Write enough code to make it pass

```go
func ConvertToRoman(arabic int) string {
	if arabic == 2 {
		return "II"
	}
	return "I"
}
```

Yup, it still feels like we're not actually tackling the problem. So we need to write more tests to drive us forward.

## Refactor

We have some repetition in our tests. When you're testing something which feels like it's a matter of "given input X, we expect Y" you should probably use table based tests.

```go
func TestRomanNumerals(t *testing.T) {
	cases := []struct {
		Description string
		Arabic      int
		Want        string
	}{
		{"1 gets converted to I", 1, "I"},
		{"2 gets converted to II", 2, "II"},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got := ConvertToRoman(test.Arabic)
			if got != test.Want {
				t.Errorf("got '%s', want '%s'", got, test.Want)
			}
		})
	}
}
```

We can now easily add more cases without having to write any more test boilerplate.

Let's push on and go for 3

## Write the test first

Add the following to our cases

```go
{"3 gets converted to III", 3, "III"},
```

## Try to run the test

```
=== RUN   TestRomanNumerals/3_gets_converted_to_III
    --- FAIL: TestRomanNumerals/3_gets_converted_to_III (0.00s)
        numeral_test.go:20: got 'I', want 'III'
```

## Write enough code to make it pass

```go
func ConvertToRoman(arabic int) string {
	if arabic == 3 {
		return "III"
	}
	if arabic == 2 {
		return "II"
	}
	return "I"
}
```

## Refactor

OK so I'm starting to not enjoy these if statements and if you look at the code hard enough you can see that we're building a string of `I` based on the size of `arabic`. 

We "know" that for more complicated numbers we will be doing some kind of arithmetic and string concatenation.  

Let's try a refactor with these thoughts in mind, it _might not_ be suitable for the end solution but that's OK. We can always throw our code away and start afresh with the tests we have to guide us.

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for i:=0; i<arabic; i++ {
		result.WriteString("I")
	}

	return result.String()
}
```

The code looks better to me and describes the domain _as we know it right now_.
 



