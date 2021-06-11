package blogposts_test

import (
	"errors"
	blogposts "github.com/quii/learn-go-with-tests/reading-files"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

type StubFailingFS struct {
}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
	)

	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	posts, err := blogposts.NewPostsFromFS(fs)

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

	t.Run("it parses the description", func(t *testing.T) {
		got := posts[0].Description
		want := "Description 1"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("it extracts the tags", func(t *testing.T) {
		got := posts[0].Tags
		want := []string{"tdd", "go"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it extracts the body", func(t *testing.T) {
		got := posts[0].Body
		want := `Hello
World`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
