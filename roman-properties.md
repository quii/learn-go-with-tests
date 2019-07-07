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

### The Romans were into DRY too...

Things start getting more complicated now. The Romans in their wisdom thought repeating characters would become hard to read and count. So a rule with Roman numerals is you cant have the same character repeated 3 times. Instead you take the next symbol and then "subtract" by putting symbols to the left of it. 

For example `5` in Roman Numerals is `V`. To create 4 you do not do `IIII`, instead you do `IV`. 

## Write the test first

```go
{"4 gets converted to IV (cant repeat more than 3 times)", 4, "IV"},
```

## Try to run the test

```
=== RUN   TestRomanNumerals/4_gets_converted_to_IV_(cant_repeat_more_than_3_times)
    --- FAIL: TestRomanNumerals/4_gets_converted_to_IV_(cant_repeat_more_than_3_times) (0.00s)
        numeral_test.go:24: got 'IIII', want 'IV'
```

## Write enough code to make it pass

```go
func ConvertToRoman(arabic int) string {
	
	if arabic == 4 {
		return "IV"
	}

	var result strings.Builder

	for i:=0; i<arabic; i++ {
		result.WriteString("I")
	}

	return result.String()
}
```

## Refactor  

I dont "like" that we have broken our string building pattern and I want to carry on with it.

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for i := arabic; i > 0; i-- {
		if i == 4 {
			result.WriteString("IV")
			break
		}
		result.WriteString("I")
	}

	return result.String()
}
```

In order for 4 to "fit" with my current thinking I now count down from the arabic number, adding symbols to our string as we progress. Not sure if this will work but let's see!

Let's make 5 work

## Write the test first

```go
{"5 gets converted to V", 5, "V"},
```

## Try to run the test

```
=== RUN   TestRomanNumerals/5_gets_converted_to_V
    --- FAIL: TestRomanNumerals/5_gets_converted_to_V (0.00s)
        numeral_test.go:25: got 'IIV', want 'V'
```

## Write enough code to make it pass

Just copy the approach we did for 4

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for i := arabic; i > 0; i-- {
		if i == 5 {
			result.WriteString("V")
			break
		}
		if i == 4 {
			result.WriteString("IV")
			break
		}
		result.WriteString("I")
	}

	return result.String()
}
```

## Refactor

Repetition in loops like this are usually a sign of an abstraction waiting to be called out. Short-circuiting loops can be an effective tool for reabability but it could also be telling you something else.

We are looping over our arabic number and if we hit certain symbols we are calling `break` but what we are _really_ doing is subtracting over `i` in a ham-fisted manner.

```go
func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for arabic > 0 {
		switch {
		case arabic > 4:
			result.WriteString("V")
			arabic -= 5
		case arabic > 3:
			result.WriteString("IV")
			arabic -= 4
		default:
			result.WriteString("I")
			arabic--
		}
	}

	return result.String()
}

```

- Given the signals I'm reading from our code, driven from our tests of some very basic scenarios I can see that to build a Roman numeral I need to subtract from `arabic` as I apply symbols
- The `for` loop no longer relies on an `i` and instead we will keep building our string until we have subtracted enough symbols away from `arabic`.

I'm pretty sure this approach will be valid for 6 (VI), 7 (VII) and 8 (VIII) too. Nonetheless add the cases in to our test suite and check (I wont include the code for brevity, check the github for samples if you're unsure).

9 follows the same rule as 4 in that we should subtract `I` from the representation of the following number. 10 is represented in Roman numerals with `X`; so therefore 9 should be `IX`.

## Write the test first

```go
{"9 gets converted to IX", 9, "IX"}
```
## Try to run the test

```
=== RUN   TestRomanNumerals/9_gets_converted_to_IX
    --- FAIL: TestRomanNumerals/9_gets_converted_to_IX (0.00s)
        numeral_test.go:29: got 'VIV', want 'IX'
```

## Write enough code to make it pass

We should be able to adopt the same approach as before

```go
case arabic > 8:
    result.WriteString("IX")
    arabic -= 9
```

## Refactor

It _feels_ like the code is still telling us there's a refactor somewhere but it's not totally obvious to me, so let's keep going. 

I'll skip the code for this too, but add to your test cases a test for `10` which should be `X` and make it pass before reading on.

Here are a few tests I added as I'm confident up to 39 our code should work

```go
{"10 gets converted to X", 10, "X"},
{"14 gets converted to XIV", 14, "XIV"},
{"18 gets converted to XVIII", 18, "XVIII"},
{"20 gets converted to XX", 20, "XX"},
{"39 gets converted to XXXIX", 39, "XXXIX"},
```

If you've ever done OO programming, you'll know that you should view `switch` statements with a bit of suspicion. Usually you are capturing a concept or data inside some imperative code when in fact it could be captured in a class structure instead. Go isn't strictly OO but that doesn't mean we ignore the lessons it offers entirely. Our switch statement is describing some truths about Roman Numerals along with behaviour. 

We can refactor this to simplify the code

```go
type RomanNumeral struct {
	Value  int
	Symbol string
}

var RomanNumerals = []RomanNumeral{
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func ConvertToRoman(arabic int) string {

	var result strings.Builder

	for _, numeral := range RomanNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}

	return result.String()
}
```




### todooooo
Here are the remaining symbols

| Arabic        | Roman           |
| ------------- |:-------------:|
| 50      | L |
| 100     | C      |
| 500 | D      |
| M | 1000      |

Take the same TDD approach for the remaining symbols
 
## TODO Property-based tests

8 wont be valid because again we will end up with repeating `I`. This is the second time we have come across this "rule" or **property** in our domain and it wont be the last as we work through the numbers. 

What we'll do is write a special kind of test called a "property based tests" where instead of providing example data, we express our domain rules in code and the test suite will throw random data at our code and see if our code respects these properties.

```go
func TestRomanNumeralProperties(t *testing.T) {
	assertion := func(arabic uint16) bool {
		return !strings.Contains(ConvertToRoman(arabic), "IIII")
	}

	if err := quick.Check(assertion, nil); err != nil {
		chError := err.(*quick.CheckError)
		in := chError.In[0].(uint16)
		t.Errorf("failed property on %d, %s", in, ConvertToRoman(in))
	}
}

```

This looks like a lot of code but it's all just Go code, dont fret.

- _Reading from the bottom_ the important call is `quick.Check`. This function will take an "assertion" and call it many times to see if it passes.
- Go lacks generics so the code around this is a bit wonky. We get ` quick.CheckError` if a check fails which contains

If you run this, your computer will hang for a while. Throw in some `log.Println` in the assertion and see what `arabic` is coming in as and you may see why. 

The library is throwing random `int`s at us and some of them are extremely large and some of them are very very negative. `int` is not actually a great datatype for us. `M` is the largest symbol with numerals which is 1000. Given our rule of no more than 3 consecutive digits that means we cant represent anything higher than 3999. 

**Just by writing the property based test it has challenged our knowledge of the domain**.

There's a few things we could do

- Create our own datatype that wraps around `int` and cant be created if it is too large or too small
- Return an error from our function if the user provides a bad input
- Return a canned string representing an error

The third option is _probably_ the worst but it is the one I will take. Remember when we have failing tests it's important to try and get out of that state as quickly as possible. Refactoring our API will lead to _more_ compilation problems. So let's just take this on the chin for now and we can maybe make it better later.

```go
const CantConvertToRoman = "Cannot convert numbers greater than 3999"

func ConvertToRoman(arabic uint16) string {
	
	if arabic > 3999 {
		return CantConvertToRoman
	}

	var result strings.Builder

	for arabic > 0 {
		switch {
		case arabic > 4:
			result.WriteString("V")
			arabic -= 5
		case arabic > 3:
			result.WriteString("IV")
			arabic -= 4
		default:
			result.WriteString("I")
			arabic--
		}
	}

	return result.String()
}
```

I _did_ change the API slightly by changing from `int` to `uint16` which means an _unsigned integer_, so we cant be passed negative numbers.

