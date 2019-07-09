# Roman Numerals

Some companies will ask you to do the [Roman Numeral Kata](http://codingdojo.org/kata/RomanNumerals/) as part of the interview process. This chapter will show how you can tackle it with TDD.

We are going to write a function which converts an (Arabic number)[https://en.wikipedia.org/wiki/Arabic_numerals] (numbers 0 to 9) to a Roman Numeral.

If you haven't heard of [Roman Numerals](https://en.wikipedia.org/wiki/Roman_numerals) they are how the Romans wrote down numbers.

You build them by sticking symbols together and those symbols represent numbers

So `I` is "one". `III` is three. 

Seems easy but there's a few interesting rules. `V` means five, but `IV` is 4 (not `IIII`). 

`MCMLXXXIV` is 1984. That looks complicated and it's hard to imagine how we can write code to figure this out right from the start.

As this book stresses, a key skill for software developers is to try and identify "thin vertical slices" of _useful_ functionality and then **iterating**. The TDD workflow helps facilitate iterative development.

So rather than 1984, let's start with 1.

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

Create our function but don't make the test pass yet, always make sure the tests fails how you expect

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

We can use subtests to nicely group our tests

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

You may not have used [`strings.Builder`](https://golang.org/pkg/strings/#Builder) before

> A Builder is used to efficiently build a string using Write methods. It minimizes memory copying.

Normally I wouldn't bother with such optimisations until I have an actual performance problem but the amount of code is not much larger than a "manual" appending on a string so we may as well use the faster approach.

The code looks better to me and describes the domain _as we know it right now_.

### The Romans were into DRY too...

Things start getting more complicated now. The Romans in their wisdom thought repeating characters would become hard to read and count. So a rule with Roman Numerals is you cant have the same character repeated 3 times in a row. 

Instead you take the next highest symbol and then "subtract" by putting a symbol to the left of it. The symbol must be base 10.

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

In order for 4 to "fit" with my current thinking I now count down from the Arabic number, adding symbols to our string as we progress. Not sure if this will work in the long run but let's see!

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

We are looping over our Arabic number and if we hit certain symbols we are calling `break` but what we are _really_ doing is subtracting over `i` in a ham-fisted manner.

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

- Given the signals I'm reading from our code, driven from our tests of some very basic scenarios I can see that to build a Roman Numeral I need to subtract from `arabic` as I apply symbols
- The `for` loop no longer relies on an `i` and instead we will keep building our string until we have subtracted enough symbols away from `arabic`.

I'm pretty sure this approach will be valid for 6 (VI), 7 (VII) and 8 (VIII) too. Nonetheless add the cases in to our test suite and check (I wont include the code for brevity, check the github for samples if you're unsure).

9 follows the same rule as 4 in that we should subtract `I` from the representation of the following number. 10 is represented in Roman Numerals with `X`; so therefore 9 should be `IX`.

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

If you've ever done OO programming, you'll know that you should view `switch` statements with a bit of suspicion. Usually you are capturing a concept or data inside some imperative code when in fact it could be captured in a class structure instead. 

Go isn't strictly OO but that doesn't mean we ignore the lessons OO offers entirely (as much as some would like to tell you). 

Our switch statement is describing some truths about Roman Numerals along with behaviour. 

We can refactor this by decoupling the data from the behaviour.

```go
type RomanNumeral struct {
	Value  int
	Symbol string
}

var RomanNumerals = []RomanNumeral {
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

This feels much better. We've declared some rules around the numerals as data rather than hidden in an algorithm and we can see how we just work through the Arabic number, trying to add symbols to our result if they fit. 

Does this abstraction work for bigger numbers? Extend the test suite so it works for the Roman number for 50 which is `L`. 

Here are some test cases, try and make them pass.

```go
{"40 gets converted to XL", 40, "XL"},
{"47 gets converted to XLVII", 47, "XLVII"},
{"49 gets converted to XLIX", 49, "XLIX"},
{"50 gets converted to XLIX", 50, "L"},
``` 

If you're a cheater, all you needed to add to the `RomanNumerals` array is

```go
{50, "L"},
{40, "XL"},
```

## And the rest!

Here are the remaining symbols

| Arabic        | Roman           |
| ------------- |:-------------:|
| 100     | C      |
| 500 | D      |
| 1000 | M      |

Take the same approach for the remaining symbols, it should just be a matter of adding data to both the tests and our array of symbols.

Does your code work for `1984`: `MCMLXXXIV` ?

Here is my final test suite

```go
func TestRomanNumerals(t *testing.T) {
	cases := []struct {
		Arabic int
		Roman  string
	}{
		{Arabic: 1, Roman: "I"},
		{Arabic: 2, Roman: "II"},
		{Arabic: 3, Roman: "III"},
		{Arabic: 4, Roman: "IV"},
		{Arabic: 5, Roman: "V"},
		{Arabic: 6, Roman: "VI"},
		{Arabic: 7, Roman: "VII"},
		{Arabic: 8, Roman: "VIII"},
		{Arabic: 9, Roman: "IX"},
		{Arabic: 10, Roman: "X"},
		{Arabic: 14, Roman: "XIV"},
		{Arabic: 18, Roman: "XVIII"},
		{Arabic: 20, Roman: "XX"},
		{Arabic: 39, Roman: "XXXIX"},
		{Arabic: 40, Roman: "XL"},
		{Arabic: 47, Roman: "XLVII"},
		{Arabic: 49, Roman: "XLIX"},
		{Arabic: 50, Roman: "L"},
		{Arabic: 100, Roman: "C"},
		{Arabic: 90, Roman: "XC"},
		{Arabic: 400, Roman: "CD"},
		{Arabic: 500, Roman: "D"},
		{Arabic: 900, Roman: "CM"},
		{Arabic: 1000, Roman: "M"},
		{Arabic: 1984, Roman: "MCMLXXXIV"},
		{Arabic: 3999, Roman: "MMMCMXCIX"},
		{Arabic: 2014, Roman: "MMXIV"},
		{Arabic: 1006, Roman: "MVI"},
		{Arabic: 798, Roman: "DCCXCVIII"},
	}
	for _, test := range cases {
		t.Run(fmt.Sprintf("%d gets converted to '%s", test.Arabic, test.Roman), func(t *testing.T) {
			got := ConvertToRoman(test.Arabic)
			if got != test.Roman {
				t.Errorf("got '%s', want '%s'", got, test.Roman)
			}
		})
	}
}
```

- I removed `description` as I felt the _data_ described enough of the information.
- I added a few other edge cases I found just to give me a little more confidence. With table based tests this is very cheap to do.

I didn't change the algorithm, all I had to do was update the `RomanNumerals` array.

```go
var RomanNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}
```

## Parsing Roman Numerals

We're not done yet. Next we're going to write a function that converts _from_ a Roman Numeral to an `int`


## Write the test first

We can re-use our test cases here with a little refactoring

Move the `cases` variable outside of the test as a package variable in a `var` block.

```go
func TestConvertingToArabic(t *testing.T) {
	for _, test := range cases[:1] {
		t.Run(fmt.Sprintf("'%s' gets converted to %d", test.Roman, test.Arabic), func(t *testing.T) {
			got := ConvertToArabic(test.Roman)
			if got != test.Arabic {
				t.Errorf("got %d, want %d", got, test.Arabic)
			}
		})
	}
}
```

Notice I am using the slice functionality to just run one of the tests for now (`cases[:1]`) as trying to make all of those tests pass all at once is too big a leap

## Try to run the test

```
./numeral_test.go:60:11: undefined: ConvertToArabic
```

## Write the minimal amount of code for the test to run and check the failing test output

Add our new function definition

```go
func ConvertToArabic(roman string) int {
	return 0
}
```

The test should now run and fail

```
--- FAIL: TestConvertingToArabic (0.00s)
    --- FAIL: TestConvertingToArabic/'I'_gets_converted_to_1 (0.00s)
        numeral_test.go:62: got 0, want 1
```

## Write enough code to make it pass

You know what to do

```go
func ConvertToArabic(roman string) int {
	return 1
}
```

Next, change the slice index in our test to move to the next test case (e.g. `cases[:2]`). Make it pass yourself with the dumbest code you can think of, continue writing dumb code (best book ever right?) for the third case too. Here's my dumb code.

```go
func ConvertToArabic(roman string) int {
	if roman == "III" {
		return 3
	}
	if roman == "II" {
		return 2
	}
	return 1
}
```

Through the dumbness of _real code that works_ we can start to see a pattern like before. We need to iterate through the input and build _something_, in this case a total.

```go
func ConvertToArabic(roman string) int {
	total := 0
	for range roman {
		total++
	}
	return total
}
```

## Write the test first

Next we move to `cases[:4]` (`IV`) which now fails because it gets 2 back as that's the length of the string.

## Write enough code to make it pass

```go
// earlier..
type RomanNumerals []RomanNumeral

func (r RomanNumerals) ValueOf(symbol string) int {
	for _, s := range r {
		if s.Symbol == symbol {
			return s.Value
		}
	}

	return 0
}

// later..
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		// look ahead to next symbol if we can and, the current symbol is base 10 (only valid subtractors)
		if i+1 < len(roman) && symbol == 'I' {
			nextSymbol := roman[i+1]

			// build the two character string
			potentialNumber := string([]byte{symbol, nextSymbol})

			// get the value of the two character string
			value := romanNumerals.ValueOf(potentialNumber)

			if value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++
			}
		} else {
			total++
		}
	}
	return total
}
```

This is horrible but it does work. It's so bad I felt the need to add comments. 

- I wanted to be able to look up an integer value for a given roman numeral so I made a type from our array of `RomanNumeral`s and then added a method to it, `ValueOf`
- Next in our loop we need to look ahead _if_ the string is big enough _and the current symbol is a valid subtractor_. At the moment it's just `I` (1) but can also be `X` (10) or `C` (100).
    - If it satisfies both of these conditions we need to lookup the value and add it to the total _if_ it is one of the special subtractors, otherwise ignore it
    - Then we need to further increment `i` so we dont count this symbol twice

## Refactor

I'm not entirely convinced this will be the long-term approach and there's potentially some interesting refactors we could do, but I'll resist that in case our approach is totally wrong. I'd rather make a few more tests pass first and see. For the meantime I made the first `if` statement slightly less horrible.

```go
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		if couldBeSubtractive(i, symbol, roman) {
			nextSymbol := roman[i+1]

			// build the two character string
			potentialNumber := string([]byte{symbol, nextSymbol})

			// get the value of the two character string
			value := romanNumerals.ValueOf(potentialNumber)

			if value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++
			}
		} else {
			total++
		}
	}
	return total
}

func couldBeSubtractive(index int, currentSymbol uint8, roman string) bool {
	return index+1 < len(roman) && currentSymbol == 'I'
}
```

## Write the test first

Let's move on to `cases[:5]`

```
=== RUN   TestConvertingToArabic/'V'_gets_converted_to_5
    --- FAIL: TestConvertingToArabic/'V'_gets_converted_to_5 (0.00s)
        numeral_test.go:62: got 1, want 5
```

## Write enough code to make it pass

Apart from when it is subtractive our code assumes that every character is a `I` which is why the value is 1. We should be able to re-use our `ValueOf` method to fix this.

```go
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		// look ahead to next symbol if we can and, the current symbol is base 10 (only valid subtractors)
		if couldBeSubtractive(i, symbol, roman) {
			nextSymbol := roman[i+1]

			// build the two character string
			potentialNumber := string([]byte{symbol, nextSymbol})

			if value := romanNumerals.ValueOf(potentialNumber); value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++ // this is fishy...
			}
		} else {
			total+=romanNumerals.ValueOf(string([]byte{symbol}))
		}
	}
	return total
}
```

## Refactor

When you index strings in Go, you get a `byte`. This is why when we build up the string again we have to do stuff like `string([]byte{symbol})`. It's repeated a couple of times, let's just move that functionality so that `ValueOf` takes some bytes instead.

```go
func (r RomanNumerals) ValueOf(symbols ...byte) int {
	symbol := string(symbols)
	for _, s := range r {
		if s.Symbol == symbol {
			return s.Value
		}
	}

	return 0
}
```

Then we can just pass in the bytes as is, to our function

```go
func ConvertToArabic(roman string) int {
	total := 0

	for i := 0; i < len(roman); i++ {
		symbol := roman[i]

		if couldBeSubtractive(i, symbol, roman) {
			if value := romanNumerals.ValueOf(symbol, roman[i+1]); value != 0 {
				total += value
				i++ // move past this character too for the next loop
			} else {
				total++ // this is fishy...
			}
		} else {
			total+=romanNumerals.ValueOf(symbol)
		}
	}
	return total
}
```

It's still pretty nasty, but it's getting there.

If you start moving our `cases[:xx]` number through you'll see that quite a few are passing now. Remove the slice operator entirely and see which ones fail, here's some examples from my suite

```
=== RUN   TestConvertingToArabic/'XL'_gets_converted_to_40
    --- FAIL: TestConvertingToArabic/'XL'_gets_converted_to_40 (0.00s)
        numeral_test.go:62: got 60, want 40
=== RUN   TestConvertingToArabic/'XLVII'_gets_converted_to_47
    --- FAIL: TestConvertingToArabic/'XLVII'_gets_converted_to_47 (0.00s)
        numeral_test.go:62: got 67, want 47
=== RUN   TestConvertingToArabic/'XLIX'_gets_converted_to_49
    --- FAIL: TestConvertingToArabic/'XLIX'_gets_converted_to_49 (0.00s)
        numeral_test.go:62: got 69, want 49
```

I think all we're missing is an update to `couldBeSubtractive` so that it accounts for the other kinds of subtractive symbols

```go
func couldBeSubtractive(index int, currentSymbol uint8, roman string) bool {
	isSubtractiveSymbol := currentSymbol == 'I' || currentSymbol == 'X' || currentSymbol =='C'
	return index+1 < len(roman) && isSubtractiveSymbol
}
```

Try again, they still fail. However we left a comment earlier...

```go
total++ // this is fishy...
```

We should never be just increment total as that implies every symbol is a `I`. Replace it with

```go
total += romanNumerals.ValueOf(symbol)
```

And all the tests pass! Now that we have fully working software we can indulge ourselves in some refactoring, with confidence

## Refactor

Here is all the code I finished up with. I had a few failed attempts but as I keep stressing, that's fine and the tests help me play around with the code freely.

```go
import "strings"

func ConvertToArabic(roman string) (total int) {
	for _, symbols := range windowedRoman(roman).Symbols() {
		total += allRomanNumerals.ValueOf(symbols...)
	}
	return
}

func ConvertToRoman(arabic int) string {
	var result strings.Builder

	for _, numeral := range allRomanNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}

	return result.String()
}

type romanNumeral struct {
	Value  int
	Symbol string
}

type romanNumerals []romanNumeral

func (r romanNumerals) ValueOf(symbols ...byte) int {
	symbol := string(symbols)
	for _, s := range r {
		if s.Symbol == symbol {
			return s.Value
		}
	}

	return 0
}

func (r romanNumerals) Exists(symbols ...byte) bool {
	symbol := string(symbols)
	for _, s := range r {
		if s.Symbol == symbol {
			return true
		}
	}
	return false
}

var allRomanNumerals = romanNumerals{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

type windowedRoman string

func (w windowedRoman) Symbols() (symbols [][]byte) {
	for i := 0; i < len(w); i++ {
		symbol := w[i]
		notAtEnd := i+1 < len(w)

		if notAtEnd && isSubtractive(symbol) && allRomanNumerals.Exists(symbol, w[i+1]) {
			symbols = append(symbols, []byte{byte(symbol), byte(w[i+1])})
			i++
		} else {
			symbols = append(symbols, []byte{byte(symbol)})
		}
	}
	return
}

func isSubtractive(symbol uint8) bool {
	return symbol == 'I' || symbol == 'X' || symbol == 'C'
}
```

My main problem with the previous code is similar to our refactor from earlier. We had too many concerns coupled together. We wrote an algorithm which was trying to extract Roman Numerals from a string _and_ then find their values. 

So I created a new type `windowedRoman` which took care of extracting the numerals, offering a `Symbols` method to retrieve them as a slice. This meant out `ConvertToArabic` function could simply iterate over the symbols and total them.

I broke the code down a bit by extracting some functions, especially around the wonky if statement to figure out if the symbol we are currently dealing with is a two character subtractive symbol.

There's probably a more elegant way but I'm not going to sweat it. The code is there and it works and it is tested. If I (or anyone else) finds a better way they can safely change it - the hard work is done.


## Wrapping up

Nothing new in this chapter, just more TDD practice! 

Did the thought of writing code that converts 1984 into MCMLXXXIV feel intimidating to you at first? It did to me and I've been writing software for quite a long time. 

The trick, as always, is to **get started with something simple** and take **small steps**. 

At no point in this process did we make any large leaps, do any huge refactorings, or get in a mess.

I can hear someone cynically saying "this is just a kata". I cant argue with that, but I still take this same approach for every project I work on. I never ship a big distributed system in my first step, I find the simplest thing the team could ship (usually a "Hello world" website) and then iterate on small bits of functionality in manageable chunks, just like how we did here.

The skill is knowing _how_ to split work up, and that comes with practice and with some lovely TDD to help you on your way.
