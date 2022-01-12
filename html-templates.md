# HTML Templates

**[You can find all the code here](https://github.com/quii/learn-go-with-tests/tree/main/htmlrenderer)**

We live in a world where everyone wants to build web applications with the latest flavour of the month frontend framework built upon gigabytes of transpiled JavaScript, working with a Byzantine build system; [but maybe that's not always neccessary](https://quii.dev/The_Web_I_Want).  

I'd say most Go developers value a simple, stable & fast toolchain but the frontend world frequently fails to deliver on this front.

Many websites do not need React. **HTML and CSS are fantastic ways of delivering content** and you can use Go to make a website to deliver HTML. 

If you wish your site to still have some dynamic elements, you can still sprinkle in some client side JavaScript or you may even want to try experimenting with [Hotwire](https://hotwired.dev) which allows you to deliver a dynamic experience with a server-side approach. 

Of course you can generate your HTML with elaborate use of `fmt.FPrintf`, but the standard library has tools to help us generate HTML in a more maintainable and flexible way. 

## What we're going to build

In the [Reading Files](/reading-files.md) chapter we wrote some code that would take an `fs.FS` and return a slice of `Post` for each markdown file it encountered.

```go
type Post struct {
    Title, Description, Body string
    Tags []string
}
```

If we continue our journey of writing some blog software, we'd take this data and generate HTML.

For our blog software, we want to generate two kinds of page:

1. **View post**. Renders a page. The `Body` field in `Post` is a string containing markdown so that should be coverted to HTML. 
2. **Index**. Lists all of the posts, most recent first, with hyperlinks to view the specific post.

We'll also want a consistent look and feel across our site, so we'll have the usual HTML furniture like `<html>` and a `<head>` containing links to CSS stylesheets and whatever else we may want.

When you're building blog software you have a few options in terms of approach of how you build and send HTML to the user's browser. 

We're assuming the preferred approach is that the HTML files will be generated as some kind of build step and a web server will serve the files. For the sake of brevity, the scope of this chapter is the writing files part, not the web server part. 

## Write the test first

As always, it's important to think about requirements before diving in too fast. How can we take this large-ish set of requirements and break it down in to a small, achievable step that we can focus on?

In my view, actually viewing content is higher priority than an index page. We could launch this product and share direct links to our wonderful content. An index page which cant link to the actual content isn't useful.

Still, rendering a post as described earlier still feels big. All the HTML furniture, converting the body markdown into HTML, listing tags, e.t.c. 

At this stage I'm not overly concerned with the specific markup, and an easy first step would be just to check we can render the post's title as an `<h1>`. This *feels* like the smallest first step that can move us forward a bit.

```go
package blogrenderer_test

import (
	"bytes"
	"github.com/quii/learn-go-with-tests/blogrenderer"
	"testing"
)

func TestRender(t *testing.T) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := blogrenderer.Render(&buf, aPost)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
                want := `<h1>hello world</h1>`
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}
```

Notice also that we're taking an `io.Writer`. This makes the code simple to test but also means we can write the contents to file, or if we change our mind we can send them to a `http.ResponseWriter` and just write the posts directly to a HTTP response.

## Try to run the test

If you've read the previous chapters of this book you should be well-practiced at this now. You won't be able to run the test because we don't have the package defined or the `Render` function. Try and follow the compiler messages yourself and get to a state where you can run the test and see that it fails with a clear message. 

It's really important that you exercise your tests failing, you'll thank yourself when you accidentally make a test fail 6 months later and you put in the effort now to check it fails with a clear message.

## Write the minimal amount of code for the test to run and check the failing test output

This is the minimal code to get the test running

```go
package blogrenderer

type Post struct {
	Title, Description, Body string
	Tags                     []string
}

func Render(w io.Writer, p Post) error {
	return nil
}
```

The test should complain that an empty string doesn't equal what we want.

## Write enough code to make it pass

```go
func Render(w io.Writer, p Post) error {
	_, err := fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
	return err
}
```

Remember, software development is primarily a learning activity. In order to discover and learn as we work, we need to work in a way that gives us frequent, high-quality feedback loops, and the easiest way to do that is work in small steps. 

So we're not worrying about using any templating libraries right now. You can make HTML just with "normal" string templating just fine, and by skipping the template part we can validate a small bit of useful behaviour and we've done a small bit of design work for our package's API.

## Refactor

Not much to refactor yet, so let's move to the next iteration

## Write the test first

Now we have a very basic version working, we can now iterate on the test to expand on the functionality. In this case, adding more information from the `Post`.

```go
	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := blogrenderer.Render(&buf, aPost)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>hello world</h1>
<p>This is a description</p>
Tags: <ul><li>go</li><li>tdd</li></ul>`

		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
```

Notice that writing this, *feels* a bit rubbish. Seeing all that markup in the test feels naff, and we haven't even put the body in, or the actual HTML we'd want with all of the `<head>` content and whatever page furniture we need.

Nonetheless, let's put up with the pain *for now*.

## Try to run the test

It should fail, complaining it doesn't have the string we expect, as we're not rendering the description and tags. 

## Write enough code to make it pass

Try and do this yourself rather than copying the code. What you should find is that making this test pass _is a bit annoying_!. When I tried, my first attempt got this error

```
=== RUN   TestRender
=== RUN   TestRender/it_converts_a_single_post_into_HTML
    renderer_test.go:32: got '<h1>hello world</h1><p>This is a description</p><ul><li>go</li><li>tdd</li></ul>' want '<h1>hello world</h1>
        <p>This is a description</p>
        Tags: <ul><li>go</li><li></li></ul>'
```

New lines! Who cares? Well, our test does. Should it? I removed the newlines for now just to get the test passing.

```go
func Render(w io.Writer, p Post) error {
	_, err := fmt.Fprintf(w, "<h1>%s</h1><p>%s</p>", p.Title, p.Description)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "Tags: <ul>")
	if err != nil {
		return err
	}

	for _, tag := range p.Tags {
		_, err = fmt.Fprintf(w, "<li>%s</li>", tag)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "</ul>")
	if err != nil {
		return err
	}

	return nil
}
```

**Yikes**. Not the nicest code i've written, and we're still only at a very early implementation of our markup. We'll need so much more content and things on our page, we're quickly seeing that this approach is not appropriate. 

Crucially though, we have a passing test; we have working software

## Refactor

With the safety-net of a passing test for working code, we can now think about changing our implementation approach at the refactoring stage. 

### Introducing templates

Go has two templating packages [text/template](https://pkg.go.dev/text/template) and [html/template](https://pkg.go.dev/html/template) and they share the same interface. The HTML one is a more specialised version of the text one. What they both do is allow you to is combine a template and some data to produce a string. 

The templating language is very similar to [Mustache](https://mustache.github.io) and allows you to dynamically generate content in a very clean fashion with a nice separation of concerns. Compared to other templating languages you may have used, it is very constrained or "logic-less" as Mustache likes to say. This is an important design decision.

Whilst we're focusing on generating HTML here, if your project is doing complex string concatenations and incantations, you might want to reach for `text/template` to clean up your code.

### Back to the code

Here is a template for our blog: 

`<h1>{{.Title}}</h1><p>{{.Description}}</p>Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>`

Where do we define this string? Well, we have a few options, but to keep the steps small, let's just start with a plain old string

```go
package blogrenderer

import (
	"html/template"
	"io"
)

const (
	postTemplate = `<h1>{{.Title}}</h1><p>{{.Description}}</p>Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>`
)

func Render(w io.Writer, p Post) error {
	templ, err := template.New("blog").Parse(postTemplate)
	if err != nil {
		return err
	}
	
	if err := templ.Execute(w, p); err != nil {
		return err
	}

	return nil
}
```

We create a new template with a name, and then parse our template string. We can then use the `Execute` method on it, passing in our data, in this case the `Post`. The template will substitute things like `{{.Description}}` with the content of `p.Description`.

This should be a pure refactor. We shouldn't need to change our tests and they should continue to pass. Importantly, our code is far easier to read and has far less annoying error handling to contend with. 

Frequently people complain about the verbosity of error handling in Go, but you might find you can find better ways to write your code so it's less error-prone in the first place, like here.

### More refactoring

Using the `html/template` has definitely been an improvement, but having it as a string constant in our code isn't great:

- Still quite difficult to read
- Not IDE/editor friendly. No syntax highlighting, ability to reformat, refactor, e.t.c.
- It looks like HTML, but you can't really work with it like you could a "normal" HTML file

What we'd really like to do is have our templates live in separate files so we can better organise them, and work with them as if they're HTML files.

Create a folder called "templates" and inside it make a file called `blog.gohtml`, paste our template into the file.

Now change our code to embed the file systems using the [embedding functionality included in go 1.16](https://pkg.go.dev/embed).

```go
package blogrenderer

import (
	"embed"
	"html/template"
	"io"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

func Render(w io.Writer, p Post) error {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return err
	}

	if err := templ.Execute(w, p); err != nil {
		return err
	}

	return nil
}
```

TODO: explain embed! Relate to the fs.FS stuff we saw in the reading files chapter.

By embedding a "file system" into our code, we can load multiple templates and combine them freely. This will become useful when we want to share rendering logic across different templates, such as a header for the top of the HTML page and a footer.

## Next: Make the template "nice"

We don't really want our template to be defined as a one line string. We want to be able to space it out to make it easier to read and work with, something like this:

```handlebars
<h1>{{.Title}}</h1>

<p>{{.Description}}</p>

Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
```

But if we do this, our test fails. This is because our test is expecting a very specific string to be returned. 

But really, we don't actually care about whitespace. Maintaining this test will become a nightmare if we have to keep painstakingly updating the assertion string every time we make minor changes to the markup. As the template grows, these kind of edits become harder to manage and the costs of work will spiral out of control.

## Introducing Approval Tests

[Go Approval Tests](https://github.com/approvals/go-approval-tests)

> ApprovalTests allows for easy testing of larger objects, strings and anything else that can be saved to a file (images, sounds, csv, etc...)

The idea is similar to "golden" files, or snapshot testing. Rather than awkwardly maintaining strings within a test file, the approval tool can compare the output for you against an "approved" file you created. You then simply copy over the new version if you approve it. Re-run the test and you're back to green.

Add a dependency to `"github.com/approvals/go-approval-tests"` to your project and edit the test to the following

```go
func TestRender(t *testing.T) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}

		if err := blogrenderer.Render(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}
```

The first time you run it, it will fail because we haven't approved anything yet

```
=== RUN   TestRender
=== RUN   TestRender/it_converts_a_single_post_into_HTML
    renderer_test.go:29: Failed Approval: received does not match approved.
```

It will have created two files, that look like the following

- `renderer_test.TestRender.it_converts_a_single_post_into_HTML.received.txt`
- `renderer_test.TestRender.it_converts_a_single_post_into_HTML.approved.txt`

The received file has the new, unapproved version of the output. Copy that into the empty approved file and re-run the test.

By copying the new version you have "approved" the change, and the test now passes.

To see the workflow in action, edit the template to how we discussed to make it easier to read (but semantically, it's the same).

```handlebars
<h1>{{.Title}}</h1>

<p>{{.Description}}</p>

Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
```

Re-run the test. A new "received" file will be generated because the output of our code differs to the approved version. Give them a look, and if you're happy with the changes, simply copy over the new version and re-run the test. Be sure to commit the approved files to source control.

The advantage of this approach is you can more easily use a diff tool to view and manage the differences, and it keeps your test code cleaner.

![Use diff tool to manage changes](https://i.imgur.com/0MoNdva.png)

This is actually a fairly minor usage of approval tests, which are an extremely useful tool in your testing arsenal. [Emily Bache](https://twitter.com/emilybache) has an [incredible video where she uses approval tests to add an incredibly extensive set of tests to a complicated codebase that has zero tests](https://www.youtube.com/watch?v=zyM2Ep28ED8). "Combinatorial Testing" is definitely something worth looking into.

Now that we have made this change, we still benefit from having our code well-tested, but the tests won't get in the way too much when we're tinkering with the markup.

### Are we still doing TDD?

An interesting side-effect of this approach is it takes us away from TDD. Of course you _could_ manually edit the approved files to the state you want, run your tests and then fix the templates so they output what you defined. 

But that's just silly! TDD is a method for doing work, specifically designing; but that doesn't mean we have to dogmatically use it for **everything**. 

The important thing is, we've done the right thing and used TDD as a **design tool** to design our package's API. For templates changes our process can be:

- Make a small change to the template
- Run the approval test
- Eyeball the output to check it looks correct
- Make the approval
- Repeat

We still shouldn't give up the value of working in small achievable steps. Try to find ways to make the changes small and keep re-running the tests to get real feedback on what you're doing.

If we start doing things like changing the code _around_ the templates, then of course that may warrant going back to our TDD method of work. 

## Expand the markup

Most websites richer HTML than we have right now. For starters, a `html` element, along with a `head`, perhaps some `nav` too. Usually there's an idea of a footer too.

If our site is going to have different pages, we'd want to define these things in one place to keep our site looking consistent. Go templates support us defining sections which we can then import in to other templates.

Edit our existing template to import a top and bottom template

```handlebars
{{template "top" .}}
<h1>{{.Title}}</h1>

<p>{{.Description}}</p>

Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
{{template "bottom" .}}
```

Then create `top.gohtml` with the following

```handlebars
{{define "top"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <title>My amazing blog!</title>
    <meta charset="UTF-8"/>
    <meta name="description" content="Wow, like and subscribe, it really helps the channel guys" lang="en"/>
</head>
<body>
<nav role="navigation">
    <div>
        <h1>Budding Gopher's blog</h1>
        <ul>
            <li><a href="/">home</a></li>
            <li><a href="about">about</a></li>
            <li><a href="archive">archive</a></li>
        </ul>
    </div>
</nav>
<main>
{{end}}
```

And `bottom.gohtml`

```handlebars
{{define "bottom"}}
<footer>
    <ul>
        <li><a href="https://twitter.com/quii">Twitter</a></li>
        <li><a href="https://github.com/quii">GitHub</a></li>
    </ul>
</footer>
</main>
</body>
</html>
{{end}}
```

Re-run your test. A new "received" file should be made and the test will fail. Check it over and if you're happy, approve it by copying it over the old version. Re-run the test again and it should pass.

## An excuse to mess around with Benchmarking

Before pressing on, let's consider what our code does.

```go
func Render(w io.Writer, p Post) error {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return err
	}

	if err := templ.Execute(w, p); err != nil {
		return err
	}

	return nil
}
```

- Parse the templates
- Use the template to render a post to an `io.Writer`

Whilst the performance impact of re-parsing the templates for each post in most cases will be fairly negligible, the effort to *not* do this is also pretty negligible and should tidy the code up a bit too.

To see the impact of not doing this parsing over and over, we can use the benchmarking tool to see how fast our function is.

```go
func BenchmarkRender(b *testing.B) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blogrenderer.Render(io.Discard, aPost)
	}
}
```

On my computer, here are the results

```
BenchmarkRender-8 22124 53812 ns/op
```

To stop us having to re-parse the templates over and over, we'll create a type that'll hold the parsed template, and that'll have a method to do the rendering

```go
type PostRenderer struct {
	templ *template.Template
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}
	
	return &PostRenderer{templ: templ}, nil
}

func (r *PostRenderer) Render(w io.Writer, p Post) error {
	
	if err := r.templ.Execute(w, p); err != nil {
		return err
	}

	return nil
}
```

This does change the interface of our code, so we'll need to update our test

```go
func TestRender(t *testing.T) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	postRenderer, err := blogrenderer.NewPostRenderer()

	if err != nil {
		t.Fatal(err)
	}

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}

		if err := postRenderer.Render(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}
```

And our benchmark

```go
func BenchmarkRender(b *testing.B) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	postRenderer, err := blogrenderer.NewPostRenderer()

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postRenderer.Render(io.Discard, aPost)
	}
}
```

The test should continue to pass. How about our benchmark?

`BenchmarkRender-8 362124 3131 ns/op`. The old NS per op were `53812 ns/op`, so this is a decent improvement! As we add other methods to render, say an Index page, it should simplify the code as we don't need to duplicate the template parsing.

## Back to the real work

In terms of rendering posts, the important part left is actually rendering the `Body`. If you recall, that should be markdown that the author has written, so it'll need converting to HTML. 

We'll leave this as an exercise for you, the reader. You should be able to find a Go library to do this for you. Use the approval test to validate what you're doing. 

The next bit of functionality we're going to do is rendering an Index, listing the posts as a HTML ordered list. 

We're expanding upon our API, so we'll put our TDD hat back on. 

## Write the test first

On the face of it an index page seems simple, but writing the test still prompts us to make some design choices

```go
t.Run("it renders an index of posts", func(t *testing.T) {
   buf := bytes.Buffer{}
   posts := []blogrenderer.Post{{Title: "Hello World"}, {Title: "Hello World 2"}}

   if err := postRenderer.RenderIndex(&buf, posts); err != nil {
      t.Fatal(err)
   }

   got := buf.String()
   want := `<ol><li><a href="/post/hello-world">Hello World</a></li><li><a href="/post/hello-world-2">Hello World 2</a></li></ol>`

   if got != want {
      t.Errorf("got %q want %q", got, want)
   }
})
```

1. We're using the `Post`'s title field as a part of the path of the URL, but we don't really want spaces in the URL so we're replacing them with hyphens.
2. We've added a `RenderIndex` method to our `PostRenderer` that again takes an `io.Writer` and a slice of `Post`.

If we had stuck with a test-after, approval tests approach here we would not be answering these questions in a controlled environment. **Tests give us space to think**. 

## Try to run the test

```
./renderer_test.go:41:13: undefined: blogrenderer.RenderIndex
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	return nil
}
```

The above should get the following test failure

```
=== RUN   TestRender
=== RUN   TestRender/it_renders_an_index_of_posts
    renderer_test.go:49: got "" want "<ol><li><a href=\"/post/hello-world\">Hello World</a></li><li><a href=\"/post/hello-world-2\">Hello World 2</a></li></ol>"
--- FAIL: TestRender (0.00s)
```

## Write enough code to make it pass

Even though this _feels_ like it should be easy, it is a bit awkward. I did it in multiple steps

```go
func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	indexTemplate := `<ol>{{range .}}<li><a href="/post/{{.Title}}">{{.Title}}</a></li>{{end}}</ol>`

	templ, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return err
	}

	if err := templ.Execute(w, posts); err != nil {
		return err
	}

	return nil
}
```

I didn't want to bother with separate filles at first, I just wanted to get it working. I view the upfront template parsing as refactoring I can do later. 

```sequence
client->server : hello
server->client: sick 
```



This doesn't pass, but it's close.

```
=== RUN   TestRender
=== RUN   TestRender/it_renders_an_index_of_posts
    renderer_test.go:49: got "<ol><li><a href=\"/post/Hello%20World\">Hello World</a></li><li><a href=\"/post/Hello%20World%202\">Hello World 2</a></li></ol>" want "<ol><li><a href=\"/post/hello-world\">Hello World</a></li><li><a href=\"/post/hello-world-2\">Hello World 2</a></li></ol>"
--- FAIL: TestRender (0.00s)
    --- FAIL: TestRender/it_renders_an_index_of_posts (0.00s)
```

You can see that the templating code is escaping the spaces in the `href` attributes. We need a way to do a string replace of spaces with hyphens. We can't just loop through the `[]Post` and replace them in-memory because we still want the spaces displayed to the user in the anchors. 

We have a few options. The first one we'll explore is passing a function in to our template. 

### Passing functions into templates 

```go
func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	indexTemplate := `<ol>{{range .}}<li><a href="/post/{{sanitiseTitle .Title}}">{{.Title}}</a></li>{{end}}</ol>`

	templ, err := template.New("index").Funcs(template.FuncMap{
		"sanitiseTitle": func(title string) string {
			return strings.ToLower(strings.Replace(title, " ", "-", -1))
		},
	}).Parse(indexTemplate)
	if err != nil {
		return err
	}

	if err := templ.Execute(w, posts); err != nil {
		return err
	}

	return nil
}
```

_Before you parse a template_ you can add a `template.FuncMap` into your template, which allows you to define functions that can be called within your template. In this case we've made a `sanitiseTitle` function which we then call inside our template with `{{sanitiseTitle .Title}}`.

This is a powerful feature, being able to send functions in to your template will allow you to do some very cool things, but, should you? Going back to the principles of Mustache and logic-less templates, why did they advocate for logic-less? **What is wrong with logic in templates?** 

As we've shown, in order to test our templates, *we've had to introduce a whole different kind of testing*. 

Imagine you introduce a function into a template which has a few different permutations of behaviour and edge cases, **how will you test it**? With this current design, your only means of testing this logic is by _rendering HTML and comparing strings_. This is not an easy or sane way of testing logic, and definitely not what you'd want for _important_ business logic. 

What Mustache-influenced templating engines give you is a useful constraint, don't try to circumvent it too often. Instead, embrace the idea of ViewModels, where you construct the bag of data you need to render in a way that's convienient for the templating language. 

This way, whatever important business logic you use to generate that bag of data can be unit tested separately, away from the messy world of HTML and templating. 

## Refactor



## Wrapping up

### On logic-less templates

As always, this is all about **separation of concerns**. It's important we consider what the responsibilities are of the various parts of our system. Too often people leak important business logic into templates, mixing up concerns and making systems difficult to understand, maintain and test.

### References and further material 

- [John Calhoun's 'Learn Web Development with Go](https://www.calhoun.io/intro-to-templates-p1-contextual-encoding/) has a number of excellent articles on templating.