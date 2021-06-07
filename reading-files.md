# Reading files

In this chapter we're going to learn how to read some files, get some data out of them and do something useful.

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

This requires us to break up our work into iterative steps, but we should be careful not to fall in to the trap of taking a "bottom up" approach. That might involve us making some kind of abstraction that is only validated once we stick everything together.

This is _not_ iterative! This is missing out on the tight feedback loops that TDD is supposed to bring us.

Kent Beck says:

> Optimism is an occupational hazard of programming. Feedback is the treatment.

Iterative means we work in "thin" vertical slices, as end-to-end as possible, keeping scope small but useful and validated with tight feedback loops.

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

Create a new project to work through this chapter

- `mkdir blogposts`
- `cd blogposts`
- `go mod init github.com/{your-name}/blogposts`
- `touch blogposts_test.go`.

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

This forces us to make the skeleton of our function to make the test run. Remember not to overthink the code at this point, we're just trying to get a running test and make sure it fails as we'd expect. If we skip this step we might be skipping over assumptions and not write a useful test.

```go
package blogposts

import "testing/fstest"

type Post struct {

}

func New(fileSystem fstest.MapFS) []Post {
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

This will make the test pass. The TDD practitioner in you might be annoyed we didn't see a failing test before writing the code to propagate the error from `fs.ReadDir`. To do this "properly", we'd need a new test where we inject a failing `fs.FS` to make `fs.ReadDir` return an `error`.

In some cases that might be the pragmatic thing to do but in our case we're not doing anything _interesting_ with the error, we're just propagating it; so it's probably not worth the hassle of writing a new test.

Logically our next iterations will be around expanding our `Post` type, so it has some useful data. As we are iterating on the same functionality it might be simpler to re-use some of the setup in our existing test rather than creating new test data for each step toward our goal. Re-work the test to allow us to make further assertions a bit easier

```go
func TestNewBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte("hi")},
		"hello-world2.md": {Data: []byte("hola")},
	}

	posts, err := blogposts.New(fs)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("it creates a post for each file in the file system", func(t *testing.T) {
		got := len(posts)
		want := len(fs)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
```

## Write the test first

We'll start with the first line in the proposed blog post schema, the title field.

We need to change the contents of the test files, so they match what was specified, and then we can make an assertion that we parse it correctly.
```go
func TestNewBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte("Title: Post 1")},
		"hello-world2.md": {Data: []byte("Title: Post 2")},
	}

	posts, err := blogposts.New(fs)

	if err != nil {
		t.Fatal(err)
	}

	t.Run("it creates a post for each file in the file system", func(t *testing.T) {
		got := len(posts)
		want := len(fs)

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it parses the title", func(t *testing.T) {
		got := posts[0].Title
		want := "Post 1"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
```

## Try to run the test
```
./blogpost_test.go:31:18: posts[0].Title undefined (type blogposts.Post has no field or method Title)
```

## Write the minimal amount of code for the test to run and check the failing test output

Add the new field to our `Post` type so that the test will run

```go
type Post struct {
	Title string
}
```

Re-run the test and you should get a clear, failing test

```
=== RUN   TestNewBlogPosts
=== RUN   TestNewBlogPosts/it_parses_the_title
    blogpost_test.go:35: got "", want "Post 1"
```

## Write enough code to make it pass

We'll need to open each file and then extract the title

```go
func New(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(fileSystem, f)
		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
	postFile, err := fileSystem.Open(f.Name())
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()

	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

Remember our focus at this point is not to write elegant code, it's just to get to a point where we have working software.

Even though this feels like a small increment forward it still required us to write a fair amount of code and make some assumptions in respect to error handling. This would be a point where you should talk to your colleagues and decide the best approach. Fast feedback loops and iterative development give us a better chance of raising these issues sooner so that we can make better decisions.

`fs.FS` gives us a way of opening a file within it by name with its `Open` method. From there we read the data from the file and for now we do not need any sophisticated parsing, just cut out the `Title: ` text by slicing the string.

## Refactor

Separating out the opening file part code from parsing its contents will make the code simpler to understand and work with.

```go
func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
	postFile, err := fileSystem.Open(f.Name())
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile fs.File) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

When you refactor out new functions or methods, take care and think about the arguments. You're designing here and because we have passing tests we are free to think deeply about what is appropriate. Think about coupling and cohesion. In this case you should ask yourself:

> Does `newPost` have to be coupled to an `fs.File` ? Do we use all the methods and data from this type? What do we _really_ need?

In our case we only use it as an argument to `io.ReadAll` which needs an `io.Reader`. So we should just loosen our coupling in our function and also ask for an `io.Reader`.

```go
func newPost(postFile io.Reader) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

You can make a similar argument for our `getPost` function which takes an `fs.DirEntry` argument but simply calls `Name()` to get the file name. We don't need all that, decouple ourselves from that type and just pass the file name through as a string. Here's the fully refactored code:

```go
func New(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(fileSystem, f.Name())
		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile io.Reader) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

From now on, most of our efforts can be neatly contained within `newPost`. The concerns of opening and iterating over files are done, and now we can focus on extracting the data for our `Post` type. Whilst not technically necessary, files are a nice way to logically group related things together so I also moved the `Post` type and `newPost` into a new `post.go` file.

## Write the test first

Let's extend our test further to extract the next line from the file, the description. Up until making it pass should now feel comfortable and familiar.

```go
func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1`
		secondBody = `Title: Post 2
Description: Description 2`
	)

	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

    // SNIP: all the previous test stuff

	t.Run("it parses the description", func(t *testing.T) {
		got := posts[0].Description
		want := "Description 1"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
```

## Try to run the test

```
./blogpost_test.go:47:18: posts[0].Description undefined (type blogposts.Post has no field or method Description)
```

## Write the minimal amount of code for the test to run and check the failing test output

Add the new field to `Post`.

```go
type Post struct {
	Title       string
	Description string
}
```

The tests should now compile, and fail.

```
=== RUN   TestNewBlogPosts
=== RUN   TestNewBlogPosts/it_creates_a_post_for_each_file_in_the_file_system
=== RUN   TestNewBlogPosts/it_parses_the_title
    blogpost_test.go:42: got "Post 1\nDescription: Description 1", want "Post 1"
=== RUN   TestNewBlogPosts/it_parses_the_description
    blogpost_test.go:51: got "", want "Description 1"
```

You'll notice that not only does our new test fail, but the title test fails too. This is because both tests are coupled to the test data and implementation. There are probably things you could do to prevent this but at some level you have to acknowledge that these things are _just coupled_ and you may as well live with it, at least for the short term.

## Write enough code to make it pass

The standard library has a handy library for helping you scan through data, line by line; [`bufio.Scanner`](https://golang.org/pkg/bufio/#Scanner)

> Scanner provides a convenient interface for reading data such as a file of newline-delimited lines of text.

```go
func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	scanner.Scan()
	titleLine := scanner.Text()

	scanner.Scan()
	descriptionLine := scanner.Text()

	return Post{Title: titleLine[7:], Description: descriptionLine[13:]}, nil
}
```

Handily, it also takes an `io.Reader` to read through (thank you again, loose-coupling) so we don't need to change our function arguments at all.

We then just `Scan` to read a line, and then extract the data using `Text`.

You'll notice this function as it stands could never return an `error`. It would be tempting at this point to remove the argument, but we do know we'll have to handle invalid file structures at some point so we may as well leave it. This has the benefit of us not having to edit the calling code back and forth.

## Refactor

We have some repitition around scanning a line and then reading the text. We know we're going to do this operation at least one more time, it's a simple refactor to DRY up so let's start with that.

```go
func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	title := readLine()[7:]
	description := readLine()[13:]

	return Post{Title: title, Description: description}, nil
}
```

This has barely saved any lines of code but that's rarely the point of refactoring. What I'm trying to do here is just separating the _what_ from the _how_ of reading lines to make the code a little more declarative to the reader.

Whilst the magic numbers of 7 and 13 get the job done, they're not awfully descriptive.

```go
const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	title := readLine()[len(titleSeparator):]
	description := readLine()[len(descriptionSeparator):]

	return Post{Title: title, Description: description}, nil
}
```

Now that I'm staring at the code with my creative refactoring mind, I'd like to try making our readLine function take care of removing the tag.

```go
func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return scanner.Text()[len(tagName):]
	}

	return Post{
		Title: readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
	}, nil
}
```

You may or may not like this, I do though. The point is in the refactoring state we are free to play with the internal details, and you can keep running your tests to check things still behave correctly. We can always go back to previous states if we're not happy.

## Wrapping up

# Notes
- Make a recording on twitch
