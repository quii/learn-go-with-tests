# Maps

**[You can find all the code for this chapter
here](https://github.com/quii/learn-go-with-tests/tree/master/maps)**

In the previous chapter, you saw how to store values in order. Now, we will
look at a way to store items by a `key` and look them up quickly.

Maps allow you to store items in a manner similar to a dictionary. You can
think of the `key` as the word and the `value` as the definition. And what better
way is there to learn about Maps than to build our own dictionary?

## Write the test first

In `dict_test.go`

```go
package main

import "testing"

func TestSearch(t *testing.T) {
	dict := Dict{"test": "this is just a test"}

	got := Search(dict, "test")
	want := "this is just a test"

	if got != want {
		t.Errorf("got %s want %s given, %s", got, want, "test")
	}
}
```
Declaring a Map is somewhat similar to an array. Except, it starts with the
`map` keyword and requires two types. The first is the key, which is written
inside the `[]`. The second is the value, which goes right after the `[]`.

The key is special. It can only be a comparable type. Comparable types are
explained in depth in the [language
spec](https://golang.org/ref/spec#Comparison_operators). But the simple version
is:
* boolean
* numeric
* string
* pointer
* channel
* interface types
* structs that contain comparable types
* arrays that contain comparable types

*if you don't know what some of these are yet, don't worry. We will get to
them later in the book.*

The value, on the other hand, can be any type you want. It can even be
another Map.

Everything else in this test should be familiar.

## Try to run the test

By running `go test` the compiler will fail with `./dict_test.go:8:9:
undefined: Search`.

## Write the minimal amount of code for the test to run and check the output

In `dict.go`

```go
package main

func Search(dict map[string]string, word string) string {
	return ""
}
```

Your test should now fail with a *clear error message*

`dict_test.go:12: got  want this is just a test given, test`

## Write enough code to make it pass

```go
func Search(dict map[string]string, word string) string {
	return dict[word]
}
```

Getting a value our of a Map is the same as getting a value out of Array
`map[key]`.

## Refactor

Our test output wasn't very clear. Let's make a small change to increase
readability and extract our assertion.

```go
func TestSearch(t *testing.T) {
	dict := Dict{"test": "this is just a test"}

	got := Search(dict, "test")
	want := "this is just a test"

	assertStrings(t, got, want)
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
```

With this in place our failing test looks a lot clearer
`dict_test.go:12: got '' want 'this is just a test'`.

I also decided to get rid of the given piece. That way this assertion is
more generally useful.

### Using a type alias

We can greatly improve our module's usage by aliasing the Map and making
`Search` a method.

In `dict_test.go`:

```go
func TestSearch(t *testing.T) {
	dict := Dict{"test": "this is just a test"}

	got := dict.Search("test")
	want := "this is just a test"

	assertStrings(t, got, want)
}
```

We switched to using a `Dict` struct which we have not defined yet, and call
`Search` on the newly created `Dict` instance.

We do not need to change the `assertStrings`.

In `dict.go`:

```go
type Dict map[string]string

func (d Dict) Search(word string) string {
	return d[word]
}
```

Here we created a type alias which acts as a thin wrapper around
the actual type. The advantage of using a type alias is that we can now
create our own methods on our Map type.

## Write the test first

The basic search was very easy to implement, but what will happen if we
supply a word that's not in our dictionary?

We actually get nothing back. This is good because the program can continue to
run, but there is a better approach. The function can report that the word is
not in the dictionary. This way, the user isn't left wondering if the word
doesn't exist or if there is just no definition (this might not seem very useful
for a dictionary.However, it's a scenario that could be key in other usecases).

```go
func TestSearch(t *testing.T) {
	dict := Dict{"test": "this is just a test"}

	t.Run("known word", func(t *testing.T) {
		got, _ := dict.Search("test")
		want := "this is just a test"

		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, got := dict.Search("test")
		want := "could not find the word you were looking for"

		if got == nil {
			t.Error("expected to receive and error.")
		} else {
			assertStrings(t, got.Error(), want)
		}
	})
}
```

The way to handle this scenario in Go is to return a second argument which is
an `Error` type.

`Error`s can be converted to a string with the `.Error()`
method, which we do when passing it to the assertion. We are also wrapping
`assertString` in an if to ensure we don't call `.Error()` on `nil`.

## Try and run the test

This does not compile

```
./dict_test.go:18:10: assignment mismatch: 2 variables but 1 values
```

## Write the minimal amount of code for the test to run and check the output

```go
func (d Dict) Search(word string) (string, error) {
	return d[word], nil
}
```

Your test should now fails with a much clearer error message.

`dict_test.go:22: expected to receive and error.`

## Write enough code to make it pass

```go
func (d Dict) Search(word string) (string, error) {
	def, ok := d[word]
	if !ok {
		return "", errors.New("could not find the word you were looking for")
	}

	return def, nil
}
```

In order to make this pass we are using an interesting property of the Map
lookup. It can return 2 values. The second value being a boolean which
indicates if the key was found successfully.

This property allows us to differentiate between a word that doesn't exist
and a word that just doesn't have a definition.

## Refactor

```go
var NotFoundError = errors.New("could not find the word you were looking for")

func (d Dict) Search(word string) (string, error) {
	def, ok := d[word]
	if !ok {
		return "", NotFoundError
	}

	return def, nil
}
```

We can get rid of the magic error in our `Search` function by bringing it up
into a constant. This will also allow us to have a better test.

```go
t.Run("unknown word", func(t *testing.T) {
    _, got := dict.Search("unknown")

    assertError(t, got, NotFoundError)
})

func assertError(t *testing.T, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error '%s' want '%s'", got, want)
	}
}
```

By creating a new helper we were able to simplify our test, and start using
our `NotFoundError` variable so our test doesn't fail if we change the error
text in the future.

## Write the test first

We have a great way to search the dictionary. However, we have no way to add
new words to our dictionary.

```go
func TestAdd(t *testing.T) {
	dict := Dict{}
	dict.Add("test", "this is just a test")

	want := "this is just a test"
	got, err := dict.Search("test")
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if want != got {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
```

In this test, we are utilizing our `Search` function to make the
validation of the dictionary a little easier.

## Write the minimal amount of code for the test to run and check output

In `dict.go`

```go
func (d Dict) Add(word, def string) {
}
```

Your test should now fail

```
dict_test.go:31: should find added word: could not find the word you were
looking for
```

## Write enough code to make it pass

```go
func (d Dict) Add(word, def string) {
	dict[word] = def
}
```

Adding to a Map is also similar to an Array. You just need to specify
key and set it equal to a value.

Another interesting property of Maps is that you can modify them without passing
them as a pointer. This is because maps are a reference types. They don't
actually hold any values. Instead, they point to the underlying data structure
which houses the data.

## Refactor

There isn't much to refactor in our implementation but the test could use
a little simplification.

```go
func TestAdd(t *testing.T) {
	dict := Dict{}
	word := "test"
	def := "this is just a test"

	dict.Add(word, def)

	assertDef(t, dict, word, def)
}

func assertDef(t *testing.T, dict Dict, word, def string) {
	t.Helper()

	got, err := dict.Search(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}

	if def != got {
		t.Errorf("got '%s' want '%s'", got, def)
	}
}
```

We made variables for word and definition, and moved the definition
assertion into it's own helper function.

Our `Add` is looking good. Except, we didn't consider what happens when the
value we are trying to add already exists!

Map will not throw an error if the value already exists. Instead, they will go
ahead and overwrite the value with the newly provided value. This can
be convenient in practice, but makes our function name less than
accurate. `Add` should not modify existing values. It should only add new
words to our dictionary.

## Write the test first

```go
func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		dict := Dict{}
		word := "test"
		def := "this is just a test"

		err := Add(dict, word, def)

		assertError(t, err, nil)
		assertDef(t, dict, word, def)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		def := "this is just a test"
		dict := Dict{word: def}
		err := Add(dict, word, "new test")

		assertError(t, err, WordExistsError)
		assertDef(t, dict, word, def)
	})
}
```

For this test we modified `Add` to return an error, which we are
validating against a new error variable, `WordExistsError`. We also modified
the previous test to check for a `nil` error.

## Try to run test

The compiler will fail because we are not return a value for `Add`.

```bash
./dict_test.go:30:13: Add(dict, word, def) used as value
./dict_test.go:41:13: Add(dict, word, "new test") used as value
```

## Write the minimal amount of code for the test to run and check the output

In `dict.go`

```go
var (
	NotFoundError   = errors.New("could not find the word you were looking for")
	WordExistsError = errors.New("cannot add word because it already exists")
)

func Add(dict map[string]string, word, def string) error {
	dict[word] = def
	return nil
}
```

Now we get two more errors. We are still modifying the value, and
returning a `nil` error.

```bash
        dict_test.go:43: got error '%!s(<nil>)' want 'cannot add word because
        it already exists'
        dict_test.go:44: got 'new test' want 'this is just a test'
```

## Write enough code to make it pass

```go
func Add(dict map[string]string, word, def string) error {
	_, err := Search(dict, word)
	switch err {
	case NotFoundError:
		dict[word] = def
	case nil:
		return WordExistsError
	default:
		return err

	}

	return nil
}
```

Here we are using a `switch` statement to match on the error. Having a `switch`
like this provides an extra safety net, in case `Search`returns an error other
than `NotFoundError`.

## Refactor

We don't have too much to refactor, but as our error usage grow we can make
a few modifications.

```go
const (
	ErrNotFound   = DictErr("could not find the word you were looking for")
	ErrWordExists = DictErr("cannot add word because it already exists")
)

type DictErr string

func (e DictErr) Error() string {
	return string(e)
}
```

The first thing you will notice, is we made the errors constant. This required
us to create our own `DictErr` type which implements the `error`
interface. You can read more about
the details in [this excellent article by Dave
Cheney](https://dave.cheney.net/2016/04/07/constant-errors). Simply put, it
makes the errors more reusable and immutable.

We also changed the names of the errors to make them IDE friendly.
If you have auto completion enables it's nice to be able to see all your
errors by typing `Err`. You can perform this change manually or try out
[gorename](https://godoc.org/golang.org/x/tools/refactor/rename), which is
a great refactoring tool!

## Write the test first

```go
func TestUpdate(t *testing.T) {
	word := "test"
	def := "this is just a test"
	dict := Dict{word: def}
	newDef := "new def"

	Update(dict, word, newDef)

	assertDef(t, dict, word, newDef)
}
```

`Update` is very closely related to `Create` and will be our next
implementation.

## Try and run the test
```
./dict_test.go:53:2: undefined: Update
```

## Write minimal amount of code for the test to run and check the failing test output

We already know how to deal with an error like this. We need to define our
function.

```go
func Update(dict map[string]string, word, def string) {}
```

With that in place we are able to see that we need to change the definition of
the word.

```
	dict_test.go:55: got 'this is just a test' want 'new def'
```

## Write enough code to make it pass

We already saw how to do this when we fixed the issue with create. So let's
implement something really similar to create.

```go
func Update(dict map[string]string, word, def string) {
	dict[word] = def
}
```

There is no refactoring we need to do on this since it was a simple change.
However, we now have the same issue as with create. If we pass in a new word,
`Update` will add it to the dictionary.

## Write the test first

```go
t.Run("existing word", func(t *testing.T) {
    word := "test"
    def := "this is just a test"
    newDef := "new def"
    dict := Dict{word: def}

    err := Update(dict, word, newDef)

    assertError(t, err, nil)
    assertDef(t, dict, word, newDef)
})

t.Run("new word", func(t *testing.T) {
    word := "test"
    def := "this is just a test"
    dict := Dict{}

    err := Update(dict, word, def)

    assertError(t, err, ErrWordDoesNotExist)
})
```

We added yet another error type for when the word does not exist. We also
modified `Update` to return an `error` value.

## Try and run the test

```
./dict_test.go:53:16: Update(dict, word, "new test") used as value
./dict_test.go:64:16: Update(dict, word, def) used as value
./dict_test.go:66:23: undefined: ErrWordDoesNotExists
```

We get 3 errors this time, but we know how to deal with these.

## Write the minimal amount of code for the test to run and check the failing test output

```go
const (
	ErrNotFound         = DictErr("could not find the word you were looking for")
	ErrWordExists       = DictErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictErr("cannot update word because it does not exist")
)

func Update(dict map[string]string, word, def string) error {
	dict[word] = def
	return nil
}
```

We added our own error type and are returning a `nil` error.

With these changes, we now get a very clear error:

```
dict_test.go:66: got error '%!s(<nil>)' want 'cannot update word because it does not exist'
```

## Write enough code to make it pass

```go
func Update(dict map[string]string, word, def string) error {
	_, err := Search(dict, word)
	switch err {
	case ErrNotFound:
		return ErrWordDoesNotExist
	case nil:
		dict[word] = def
	default:
		return err

	}

	dict[word] = def
	return nil
}
```

This algorithm looks almost identical to `Add` except we switched when we update
the `dict` and when we return an error.

### Note on declaring a new error for Update

We could reused `ErrNotFound` and not added a new error. However, it is often
better to have a precise error for when an update fails.

Having specific errors allows your application to know more about what went
wrong. For example, if you are running a website. You might not want the user to
see `ErrNotFound`, but instead redirect them to a add page. While
`ErrWordDoesNotExist` would be displayed when they are trying to update a word.

## Write the test first
