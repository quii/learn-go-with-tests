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
	dict := map[string]string{"test": "this is just a test"}

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

The key is special. It can only be a comparable type. Comparable types are explained in depth in the [language spec](https://golang.org/ref/spec#Comparison_operators). But the simple version is:
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
	dict := map[string]string{"test": "this is just a test"}

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
