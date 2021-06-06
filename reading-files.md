# Reading files

In this chapter we're going to learn how to read some files, get some data out of them and do something useful. Naturally it'll be in a test-driven manner so we have some modular, nicely tested and simple to maintain code.

We've been asked to create a package that converts a given folder of blog posts and return a collection of `Post` which represents each file parsed with information about its contents.

### Example data

hello world.md
```
Title: Hello, TDD world!
Description: First post on our wonderful blog
Tags: tdd, go
---

Hello world!

The body of posts starts after the `---`
```

### Expected data

```go
type Post struct {
	Title, Description, Body string
	Tags []string
}
```

## Iterative, test-driven development

As always we'll take an iterative approach where we're always taking simple, safe steps toward our goal.

This requires us to break up our work into iterative steps, but we should be careful not to fall in to the trap of taking a "bottom up" approach. That might involve us making some kind of abstraction that is only validated once we stick everything together. This is _not_ iterative! This is missing out on the tight feedback loops that TDD is supposed to bring us.

Iterative means we work in "thin" vertical slices, as end-to-end as possible, keeping scope small but useful and validated.

Let's remind ourselves of our mindset and goals when starting:

- Write the test we want to see
- Focused on the what, rather than the how
- Consumer focused

## Thinking about the kind of test we want to see

Our package needs to offer a function that can be pointed at some kind of folder, do the hard work and return us some posts.

```go
var posts blogposts.Post
posts = blogposts.New("some-folder")
```

To write a test around this, we'd need some kind of test folder with some example posts in it. _There's nothing wrong with this_, but it's not strictly a unit test because it goes to the file system. This means it'll be slightly slower, and a little more difficult to maintain although admittedly, not by a huge amount. We're also at risk of our code coupling itself to a specific file-system implementation that we don't really need to do.

### File system abstractions introduced in Go 1.16

Go 1.16 introduced an abstraction for file-systems; the [io/fs](https://golang.org/pkg/io/fs/) package.

> Package fs defines basic interfaces to a file system. A file system can be provided by the host operating system but also by other packages.

This lets us loosen our coupling to a specific file system, which will then let us inject different implementations according to our needs.

> [On the producer side of the interface, the new embed.FS type implements fs.FS, as does zip.Reader. The new os.DirFS function provides an implementation of fs.FS backed by a tree of operating system files.](https://golang.org/doc/go1.16#fs)

If we use this interface, users of our package have a number of options baked in to the standard library to use. Learning to leverage interfaces defined in Go's standard library (like this, but also `io.Reader` and `io.Writer`) is vital to writing loosely coupled packages that can be re-used in contexts different to what you imagined with minimal fuss from your consumers.

In our case maybe our consumer wants the posts to be embedded into the Go binary rather than files in a "real" filesystem, either way _our code doesn't need to care_.

For our tests, the package [testing/fstest](https://golang.org/pkg/testing/fstest/), offers us an implementation of [io/FS](https://golang.org/pkg/io/fs/#FS) to use, similar to the tools we're familiar with in [net/http/httptest](https://golang.org/pkg/net/http/httptest/).

Given this information, the following feels like a better approach

```go
var posts blogposts.Post
posts = blogposts.New(someFS)
```


## Write the test first

As discussed, we should try to keep scope as small and end-to-end as possible. A good first start to give us confidence is to prove we can read all the files in a directory and check the count of posts is the same as the number of files inside.

Create a new project (`mkdir blogposts`, cd into it, `go mod init github.com/{your-name}/blogposts`) and then create a new file `blogposts_test.go`.

```go
package blogposts_test

import (
	"testing"
	"testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {
	t.Run("it creates a post for each file in the file system", func(t *testing.T) {
		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte("hi")},
			"hello-world2.md": {Data: []byte("hola")},
		}

		posts := blogposts.New(fs)

		got := len(posts)
		want := len(fs)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

```

Notice that the package of our test is `blogposts_test`. Remember when TDD is practiced well we take a consumer-driven approach, we don't want to test internal details because consumers don't care about that. By appending `_test` to our intended package name it means we can only access exported members from our package, just like a real user of our package.

We've imported [`testing/fstest`](https://golang.org/pkg/testing/fstest/) which gives us access to the [`fstest.MapFS`](https://golang.org/pkg/testing/fstest/#MapFS) type for our fake file-system to pass to our package.

> A MapFS is a simple in-memory file system for use in tests, represented as a map from path names (arguments to Open) to information about the files or directories they represent.

This feels simpler than maintaining a folder of test files and will execute quicker too.

Finally, we've written down the usage of our API from a consumer's point of view and then check if it at least creates the correct number of posts

## Try to run the test

```
./blogpost_test.go:15:12: undefined: blogposts
```

## Write the minimal amount of code for the test to run and check the failing test output

The package doesn't exist, create a new file `blogposts.go` and put `package blogposts` inside it. You'll need to then import that package into your tests. For me the imports now look like:

```go
import (
	blogposts "github.com/quii/learn-go-with-tests/reading-files"
	"testing"
	"testing/fstest"
)
```

The tests still wont compile because our new package does not have a `New` function that returns some kind of collection

```
./blogpost_test.go:16:12: undefined: blogposts.New
```

This forces us to make the skeleton of our function to make the test run. Remember not to overthink the code at this point, we're just trying to get a running test.

```go
package blogposts

import "testing/fstest"

type Post struct {

}

func New(fs fstest.MapFS) []Post {
	return nil
}
```

The test should now correctly fail

```
=== RUN   TestNewBlogPosts/it_creates_a_post_for_each_file_in_the_file_system
    blogpost_test.go:22: got 0, want 2
    --- FAIL: TestNewBlogPosts/it_creates_a_post_for_each_file_in_the_file_system (0.00s)
```

## Write enough code to make it pass

We _could_ ["slime"](https://deniseyu.github.io/leveling-up-tdd/) this to make it pass:

```go
func New(fileSystem fstest.MapFS) []Post {
	return []Post{{},{}}
}
```

But, as Denise wrote:

>Sliming is useful for giving a “skeleton” to your object. Designing an interface and executing logic are two concerns, and sliming tests strategically lets you focus on one at a time.

We already have our structure. So what do we do instead?

As we've cut scope, all we need to do is read the directory and create a post for each file we encounter. We don't have to worry about opening files and parsing them just yet.

```go
func New(fileSystem fstest.MapFS) []Post {
	dir, _ := fs.ReadDir(fileSystem, ".")
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts
}
```

[`fs.ReadDir`](https://golang.org/pkg/io/fs/#ReadDir) reads a directory inside a given `fs.FS` returning [`[]DirEntry`](https://golang.org/pkg/io/fs/#DirEntry).

Already our idealised view of the world has been foiled because errors can happen but remember right now our focus is to make the test pass, not changing design, so we'll ignore the error for now.

The rest of the code is straightforward, iterate over the entries and create a `Post` for each one and return the slice.

## Refactor

Whilst our tests are passing, in practice we would not be able to use our new package outside this context because it is coupled to a concrete implementation, `fstest.MapFS` but as discussed it doesn't have to be. Change the argument to our `New` function to accept the interface from the standard library.

```go
func New(fileSystem fs.FS) []Post {
	dir, _ := fs.ReadDir(fileSystem, ".")
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts
}
```

Re-run the tests and everything should still be working.

### Error handling

We parked error handling before as we were focused on making the happy-path work. Before continuing to iterate on the functionality we should acknowledge that errors can happen when working with files. Beyond reading the directory, when we open the individual files we can also get problems so let's change our API (via our tests first, naturally) so that it can return an `error`.

```go
func TestNewBlogPosts(t *testing.T) {
	t.Run("it creates a post for each file in the file system", func(t *testing.T) {
		fs := fstest.MapFS{
			"hello world.md":  {Data: []byte("hi")},
			"hello-world2.md": {Data: []byte("hola")},
		}

		posts, err := blogposts.New(fs)

		if err != nil {
			t.Fatal(err)
		}

		got := len(posts)
		want := len(fs)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
```

Run the test, and it should complain about the wrong number of return values, fixing the code is straightforward.

```go
func New(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts, nil
}
```

This will make the test pass. The TDD practitioner in you might be annoyed we didn't see a failing test before writing the code to propagate the error from `fs.ReadDir`. To do this "properly", we'd need a new test where we inject a failing `fs.FS` to make `fs.ReadDir` return an `error`. In some cases that might be the pragmatic thing to do but in our case we're not doing anything _interesting_ with the error, we're just propagating it; so it's probably not worth the hassle of writing a test.


## Wrapping up

# Notes
- Make a recording on twitch
