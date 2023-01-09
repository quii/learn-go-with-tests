# Learn Go with Tests - Scaling Acceptance Tests (and light intro to gRPC)

This chapter is a follow-up to [Intro to acceptance tests](https://quii.gitbook.io/learn-go-with-tests/testing-fundamentals/intro-to-acceptance-tests). You can find [the finished code for this chapter on GitHub](https://github.com/quii/go-specs-greet).

Acceptance tests are essential, and they directly impact your ability to confidently evolve your system over time, with a reasonable cost of change.

They're also a fantastic tool to help you work with legacy code. When faced with a poor codebase without any tests, please resist the temptation to start refactoring. Instead, write some acceptance tests to give you a safety net to freely change the system's internals without affecting its functional external behaviour. ATs need not be concerned with internal quality, so they're a great fit in these situations.

After reading this, you'll appreciate that acceptance tests are useful for verification and can also be used in the development process by helping us change our system more deliberately and methodically, reducing wasted effort.

## Prerequisite material

The inspiration for this chapter is borne of many years of frustration with acceptance tests. Two videos I would recommend you watch are:

- Dave Farley - [How to write acceptance tests](https://www.youtube.com/watch?v=JDD5EEJgpHU)
- Nat Pryce - [E2E functional tests that can run in milliseconds](https://www.youtube.com/watch?v=Fk4rCn4YLLU)

"Growing Object Oriented Software" (GOOS) is such an important book for many software engineers, including myself. The approach it prescribes is the one I coach engineers I work with to follow.

- [GOOS](http://www.growing-object-oriented-software.com) - Nat Pryce & Steve Freeman

Finally, [Riya Dattani](https://twitter.com/dattaniriya) and I spoke about this topic in the context of BDD in our talk, [Acceptance tests, BDD and Go](https://www.youtube.com/watch?v=ZMWJCk_0WrY).

## Recap

We're talking about "black-box" tests that verify your system behaves as expected from the outside, from a "**business perspective**". The tests do not have access to the innards of the system it tests; they're only concerned with **what** your system does rather than **how**.

## Anatomy of bad acceptance tests

Over many years, I've worked for several companies and teams. Each of them recognised the need for acceptance tests; some way to test a system from a user's point of view and to verify it works how it's intended, but almost without exception, the cost of these tests became a real problem for the team.

- Slow to run
- Brittle
- Flaky
- Expensive to maintain, and seem to make changing the software harder than it ought to be
- Can only run in a particular environment, causing slow and poor feedback loops

Let's say you intend to write an acceptance test around a website you're building. You decide to use a headless web browser (like [Selenium](https://www.selenium.dev)) to simulate a user clicking buttons on your website to verify it does what it needs to do.

Over time, your website's markup has to change as new features are discovered, and engineers bike-shed over whether something should be an `<article>` or a `<section>` for the billionth time.

Even though your team are only making minor changes to the system, barely noticeable to the actual user, you find yourself wasting lots of time updating your ATs.

### Tight-coupling

Think about what prompts acceptance tests to change:

- An external behaviour change. If you want to change what the system does, changing the acceptance test suite seems reasonable, if not desirable.
- An implementation detail change / refactoring. Ideally, this shouldn't prompt a change, or if it does, a minor one.

Too often, though, the latter is the reason acceptance tests have to change. To the point where engineers even become reluctant to change their system because of the perceived effort of updating tests!

![Riya and myself talking about separating concerns in our tests](https://i.imgur.com/bbG6z57.png)

These problems stem from not applying well-established and practised engineering habits written by the authors mentioned above. **You can't write acceptance tests like unit tests**; they require more thought and different practices.

## Anatomy of good acceptance tests

If we want acceptance tests that only change when we change behaviour and not implementation detail, it stands to reason that we need to separate those concerns.

### On types of complexity

As software engineers, we have to deal with two kinds of complexity.

- **Accidental complexity** is the complexity we have to deal with because we're working with computers, stuff like networks, disks, APIs, etc.

- **Essential complexity** is sometimes referred to as "domain logic". It's the particular rules and truths within your domain.
  - For example, "if an account owner withdraws more money than is available, they are overdrawn". This statement says nothing about computers; this statement was true before computers were even used in banks!

Essential complexity should be expressible to a non-technical person, and it's valuable to have modelled it in our "domain" code, and in our acceptance tests.

### Separation of concerns

What Dave Farley proposed in the video earlier, and what Riya and I also discussed, is we should have the idea of **specifications**. Specifications describe the behaviour of the system we want without being coupled with accidental complexity or implementation detail.

This idea should feel reasonable to you. In production code, we frequently strive to separate concerns and decouple units of work. Would you not hesitate to introduce an `interface` to allow your `HTTP` handler to decouple it from non-HTTP concerns? Let's take this same line of thinking for our acceptance tests.

Dave Farley describes a specific structure.

![Dave Farley on Acceptance Tests](https://i.imgur.com/nPwpihG.png)

At GopherconUK, Riya and I put this in Go terms.

![Separation of concerns](https://i.imgur.com/qdY4RJe.png)

### Testing on steroids

Decoupling how the specification is executed allows us to reuse it in different scenarios. We can:

#### Make our drivers configurable

This means you can run your ATs locally, in your staging and (ideally) production environments.
- Too many teams engineer their systems such that acceptance tests are impossible to run locally. This introduces an intolerably slow feedback loop. Wouldn't you rather be confident your ATs will pass _before_ integrating your code? If the tests start breaking, is it acceptable that you'd be unable to reproduce the failure locally and instead, have to commit changes and cross your fingers that it'll pass 20 minutes later in a different environment?
- Remember, just because your tests pass in staging doesn't mean your system will work. Dev/Prod parity is, at best, a white lie. [I test in prod](https://increment.com/testing/i-test-in-production/).
- There are always differences between the environments that can affect the *behaviour* of your system. A CDN could have some cache headers incorrectly set; a downstream service you depend on may behave differently; a configuration value may be incorrect. But wouldn't it be nice if you could run your specifications in prod to catch these problems quickly?

#### Plug in _different_ drivers to test other parts of your system

This flexibility allows us to test behaviours at different abstraction and architectural layers, which allows us to have more focused tests beyond black-box tests.
- For instance, you may have a web page with an API behind it. Why not use the same specification to test both? You can use a headless web browser for the web page, and HTTP calls for the API.
- Taking this idea further, ideally, we want the **code to model essential complexity** (as "domain" code) so we should also be able to use our specifications for unit tests. This will give swift feedback that the essential complexity in our system is modelled and behaves correctly.


### Acceptance tests changing for the right reasons

With this approach, the only reason for your specifications to change is if the behaviour of the system changes, which is reasonable.

- If your HTTP API has to change, you have one obvious place to update it, the driver.
- If your markup changes, again, update the specific driver.

As your system grows, you'll find yourself reusing drivers for multiple tests, which again means if implementation detail changes, you only have to update one, usually obvious place.

When done right, this approach gives us flexibility in our implementation detail and stability in our specifications. Importantly, it provides a simple and obvious structure for managing change, which becomes essential as a system and its team grows.

### Acceptance tests as a method for software development

In our talk, Riya and I discussed acceptance tests and their relation to BDD. We talked about how starting your work by trying to _understand the problem you're trying to solve_ and expressing it as a specification helps focus your intent and is a great way to start your work.

I was first introduced to this way of working in GOOS. A while ago, I summarised the ideas on my blog. Here is an extract from my post [Why TDD](https://quii.dev/The_Why_of_TDD)

---

TDD is focused on letting you design for the behaviour you precisely need, iteratively. When starting a new area, you must identify a key, necessary behaviour and aggressively cut scope.

Follow a "top-down" approach, starting with an acceptance test (AT) that exercises the behaviour from the outside. This will act as a north-star for your efforts. All you should be focused on is making that test pass. This test will likely be failing for a while whilst you develop enough code to make it pass.

![](https://i.imgur.com/pxTaYu4.png)

Once your AT is set up, you can break into the TDD process to drive out enough units to make the AT pass. The trick is to not worry too much about design at this point; get enough code to make the AT pass because you're still learning and exploring the problem.

Taking this first step is often more extensive than you think, setting up web servers, routing, configuration, etc., which is why keeping the scope of the work small is essential. We want to make that first positive step on our blank canvas and have it backed by a passing AT so we can continue to iterate quickly and safely.

![](https://i.imgur.com/t5y5opw.png)

As you develop, listen to your tests, and they should give you signals to help you push your design in a better direction but, again, anchored to the behaviour rather than our imagination.

Typically, your first "unit" that does the hard work to make the AT pass will grow too big to be comfortable, even for this small amount of behaviour. This is when you can start thinking about how to break the problem down and introduce new collaborators.

![](https://i.imgur.com/UYqd7Cq.png)

This is where test doubles (e.g. fakes, mocks) are handy because most of the complexity that lives internally within software doesn't usually reside in implementation detail but "between" the units and how they interact.

#### The perils of bottom-up

This is a "top-down" approach rather than a "bottom-up". Bottom-up has its uses, but it carries an element of risk. By building "services" and code without it being integrated into your application quickly and without verifying a high-level test, **you risk wasting lots of effort on unvalidated ideas**.

This is a crucial property of the acceptance-test-driven approach, using tests to get real validation of our code.

Too many times, I've encountered engineers who have made a chunk of code, in isolation, bottom-up, they think is going to solve a job, but it:

- Doesn't work how we want to
- Does stuff we don't need
- Doesn't integrate easily
- Requires a ton of re-writing anyway

This is waste.

## Enough talk, time to code

Unlike other chapters, you'll need [Docker](https://www.docker.com) installed because we'll be running our applications in containers. It's assumed at this point in the book you're comfortable writing Go code, importing from different packages, etc.

Create a new project with `go mod init github.com/quii/go-specs-greet` (you can put whatever you like here but if you change the path you will need to change all internal imports to match)

Make a folder `specifications` to hold our specifications, and add a file `greet.go`

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet() (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet()
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, world")
}
```

My IDE (Goland) takes care of the fuss of adding dependencies for me, but if you need to do it manually, you'd do

`go get github.com/alecthomas/assert/v2`

Given Farley's acceptance test design (Specification->DSL->Driver->System), we now have a decoupled specification from implementation. It doesn't know or care about _how_ we `Greet`; it's just concerned with the essential complexity of our domain. Admittedly this complexity isn't much right now, but we'll expand upon the spec to add more functionality as we further iterate. It's always important to start small!

You could view the interface as our first step of a DSL; as the project grows, you may find the need to abstract differently, but for now, this is fine.

At this point, this level of ceremony to decouple our specification from implementation might make some people accuse us of "overly abstracting". **I promise you that acceptance tests that are too coupled to implementation become a real burden on engineering teams**. I am confident that most acceptance tests out in the wild are expensive to maintain due to this inappropriate coupling; rather than the reverse of being overly abstract.

We can use this specification to verify any "system" that can `Greet`.

### First system: HTTP API

We require to provide a "greeter service" over HTTP. So we'll need to create:

1. A **driver**. In this case, one works with an HTTP system by using an **HTTP client**. This code will know how to work with our API. Drivers translate DSLs into system-specific calls; in our case, the driver will implement the interface specifications define.
2. An **HTTP server** with a greet API
3. A **test**, which is responsible for managing the life-cycle of spinning up the server and then plugging the driver into the specification to run it as a test

## Write the test first

The initial process for creating a black-box test that compiles and runs your program, executes the test and then cleans everything up can be quite labour intensive. That's why it's preferable to do it at the start of your project with minimal functionality. I typically start all my projects with a "hello world" server implementation, with all of my tests set up and ready for me to build the actual functionality quickly.

The mental model of "specifications", "drivers", and "acceptance tests" can take a little time to get used to, so follow carefully. It can be helpful to "work backwards" by trying to call the specification first.

Create some structure to house the program we intend to ship.

`mkdir -p cmd/httpserver`

Inside the new folder, create a new file `greeter_server_test.go`, and add the following.

```go
package main_test

import (
	"testing"

	"github.com/quii/specifications"
)

func TestGreeterServer(t *testing.T) {
	specifications.GreetSpecification(t, nil)
}
```

We wish to run our specification in a Go test. We already have access to a `*testing.T`, so that's the first argument, but what about the second?

`specifications.Greeter` is an interface, which we will implement with a `Driver` by changing the new TestGreeterServer code to the following:

```go
func TestGreeterServer(t *testing.T) {
  	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080"}
	specifications.GreetSpecification(t, driver)
}
```

It would be favourable for our `Driver` to be configurable to run it against different environments, including locally, so we have added a `BaseURL` field.

## Try to run the test

```
./greeter_server_test.go:46:12: undefined: go_specs_greet.Driver
```

We're still practising TDD here! It's a big first step we have to make; we need to make a few files and write maybe more code than we're typically used to, but when you're first starting, this is often the case. It's so important we try to remember the red step's rules.

> Commit as many sins as necessary to get the test passing

## Write the minimal amount of code for the test to run and check the failing test output

Hold your nose; remember, we can refactor when the test has passed. Here's the code for the driver in `driver.go` which we will place in the project root:

```go
package go_specs_greet

import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
}

func (d Driver) Greet() (string, error) {
	res, err := http.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```


Notes:

- You could argue that I should be writing tests to drive out the various `if err != nil`, but in my experience, so long as you're not doing anything with the `err`, tests that say "you return the error you get" are relatively low value.
- **You shouldn't use the default HTTP client**. Later we'll pass in an HTTP client to configure it with timeouts etc., but for now, we're just trying to get ourselves to a passing test.

Try and rerun the tests; they should now compile but not pass.

```
Get "http://localhost:8080/greet": dial tcp [::1]:8080: connect: connection refused
```

We have a `Driver`, but we have not started our application yet, so it cannot do an HTTP request. We need our acceptance test to coordinate building, running and finally killing our system for the test to run.

### Running our application

It's common for teams to build Docker images of their systems to deploy, so for our test we'll do the same

To help us use Docker in our tests, we will use [Testcontainers](https://golang.testcontainers.org). Testcontainers gives us a programmatic way to build Docker images and manage container life-cycles.

`go get github.com/testcontainers/testcontainers-go`

Now you can edit `cmd/httpserver/greeter_server_test.go` to read as follows:

```go
package main_test

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGreeterServer(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "./cmd/httpserver/Dockerfile",
	  // set to false if you want less spam, but this is helpful if you're having troubles
	  PrintBuildLog: true,
		},
	ExposedPorts: []string{"8080:8080"},
		WaitingFor:   wait.ForHTTP("/").WithPort("8080"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})

	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080"}
	specifications.GreetSpecification(t, driver)
}
```

Try and run the test.

```
=== RUN   TestGreeterHandler
2022/09/10 18:49:44 Starting container id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Waiting for container id 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Container is ready id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
    greeter_server_test.go:32: Did not expect an error but got:
        Error response from daemon: Cannot locate specified Dockerfile: ./cmd/httpserver/Dockerfile: failed to create container
--- FAIL: TestGreeterHandler (0.59s)
```

We need to create a Dockerfile for our program. Inside our `httpserver` folder, create a `Dockerfile` and add the following.

```dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/httpserver/*.go

EXPOSE 8080
CMD [ "./svr" ]
```

Don't worry too much about the details here; it can be refined and optimised, but for this example, it'll suffice. The advantage of our approach here is we can later improve our Dockerfile and have a test to prove it works as we intend it to. This is a real strength of having black-box tests!

Try and rerun the test; it should complain about not being able to build the image. Of course, that's because we haven't written a program to build yet!

For the test to fully execute, we'll need to create a program that listens on `8080`, but **that's all**. Stick to the TDD discipline, don't write the production code that would make the test pass until we've verified the test fails as we'd expect.

Create a `main.go` inside our `httpserver` folder with the following

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

Try to run the test again, and it should fail with the following.

```
    greet.go:16: Expected values to be equal:
        +Hello, World
        \ No newline at end of file
--- FAIL: TestGreeterHandler (2.09s)
```

## Write enough code to make it pass

Update the handler to behave how our specification wants it to

```go
func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

## Refactor

Whilst this technically isn't a refactor, we shouldn't rely on the default HTTP client, so let's change our Driver, so we can supply one, which our test will give.

```go
type Driver struct {
	BaseURL string
	Client *http.Client
}

func (d Driver) Greet() (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

In our test in `cmd/httpserver/greeter_server_test.go`, update the creation of the driver to pass in a client.

```go
client := http.Client{
  Timeout: 1 * time.Second,
}

driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080", Client: &client}
specifications.GreetSpecification(t, driver)
```

It's good practice to keep `main.go` as simple as possible; it should only be concerned with piecing together the building blocks you make into an application.

Create a file in the project root called `handler.go` and move our code into there.

```go
package go_specs_greet

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}
```

Update `main.go` to import and use the handler instead.

```go
package main

import (
	"net/http"

	go_specs_greet "github.com/quii/go-specs-greet"
)

func main() {
	handler := http.HandlerFunc(go_specs_greet.Handler)
	http.ListenAndServe(":8080", handler)
}
```

## Reflect

The first step felt like an effort. We've made several `go` files to create and test an HTTP handler that returns a hard-coded string. This "iteration 0" ceremony and setup will serve us well for further iterations.

Changing functionality should be simple and controlled by driving it through the specification and dealing with whatever changes it forces us to make. Now the `DockerFile` and `testcontainers` are set up for our acceptance test; we shouldn't have to change these files unless the way we construct our application changes.

We'll see this with our following requirement, greet a particular person.

## Write the test first

Edit our specification

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet(name string) (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet("Mike")
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, Mike")
}
```

To allow us to greet specific people, we need to change the interface to our system to accept a `name` parameter.

## Try to run the test

```
./greeter_server_test.go:48:39: cannot use driver (variable of type go_specs_greet.Driver) as type specifications.Greeter in argument to specifications.GreetSpecification:
	go_specs_greet.Driver does not implement specifications.Greeter (wrong type for Greet method)
		have Greet() (string, error)
		want Greet(name string) (string, error)
```

The change in the specification has meant our driver needs to be updated.

## Write the minimal amount of code for the test to run and check the failing test output

Update the driver so that it specifies a `name` query value in the request to ask for a particular `name` to be greeted.

```go
func (d Driver) Greet(name string) (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet?name=" + name)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

The test should now run, and fail.

```
    greet.go:16: Expected values to be equal:
        -Hello, world
        \ No newline at end of file
        +Hello, Mike
        \ No newline at end of file
--- FAIL: TestGreeterHandler (1.92s)
```

## Write enough code to make it pass

Extract the `name` from the request and greet.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.URL.Query().Get("name"))
}
```

The test should now pass.

## Refactor

In [HTTP Handlers Revisited,](https://github.com/quii/learn-go-with-tests/blob/main/http-handlers-revisited.md) we discussed how important it is for HTTP handlers should only be responsible for handling HTTP concerns; any "domain logic" should live outside of the handler. This allows us to develop domain logic in isolation from HTTP, making it simpler to test and understand.

Let's pull apart these concerns.

Update our handler in `./handler.go` as follows:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, Greet(name))
}
```

Create new file `./greet.go`:

```go
package go_specs_greet

import "fmt"

func Greet(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}
```

## A slight diversion in to the "adapter" design pattern

Now that we've separated our domain logic of greeting people into a separate function, we are now free to write unit tests for our greet function. This is undoubtedly a lot simpler than testing it through a specification that goes through a driver that hits a web server, to get a string!

Wouldn't it be nice if we could reuse our specification here too? After all, the specification's point is decoupled from implementation details. If the specification captures our **essential complexity** and our "domain" code is supposed to model it, we should be able to use them together.

Let's give it a go by creating  `./greet_test.go` as follows:

```go
package go_specs_greet_test

import (
	"testing"

	"github.com/quii/go-specs-greet/specifications"
	go_specs_greet "github.com/quii/go-specs-greet"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecification(t, go_specs_greet.Greet)
}

```

This would be nice, but it doesn't work

```
./greet_test.go:11:39: cannot use go_specs_greet.Greet (value of type func(name string) string) as type specifications.Greeter in argument to specifications.GreetSpecification:
	func(name string) string does not implement specifications.Greeter (missing Greet method)
```

Our specification wants something that has a method `Greet()` not a function.

The compilation error is frustrating; we have a thing that we "know" is a `Greeter`, but it's not quite in the right **shape** for the compiler to let us use it. This is what the **adapter** pattern caters for.

> In [software engineering](https://en.wikipedia.org/wiki/Software_engineering), the **adapter pattern** is a [software design pattern](https://en.wikipedia.org/wiki/Software_design_pattern) (also known as [wrapper](https://en.wikipedia.org/wiki/Wrapper_function), an alternative naming shared with the [decorator pattern](https://en.wikipedia.org/wiki/Decorator_pattern)) that allows the [interface](https://en.wikipedia.org/wiki/Interface_(computer_science)) of an existing [class](https://en.wikipedia.org/wiki/Class_(computer_science)) to be used as another interface.[[1\]](https://en.wikipedia.org/wiki/Adapter_pattern#cite_note-HeadFirst-1) It is often used to make existing classes work with others without modifying their [source code](https://en.wikipedia.org/wiki/Source_code).

A lot of fancy words for something relatively simple, which is often the case with design patterns, which is why people tend to roll their eyes at them. The value of design patterns is not specific implementations but a language to describe specific solutions to common problems engineers face. If you have a team that has a shared vocabulary, it reduces the friction in communication.

Add this code in `./specifications/adapters.go`

```go
type GreetAdapter func(name string) string

func (g GreetAdapter) Greet(name string) (string, error) {
	return g(name), nil
}
```

We can now use our adapter in our test to plug our `Greet` function into the specification.

```go
package go_specs_greet_test

import (
	"testing"

	"github.com/quii/go-specs-greet/specifications"
	gospecsgreet "github.com/quii/go-specs-greet"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		specifications.GreetAdapter(gospecsgreet.Greet),
	)
}
```

The adapter pattern is handy when you have a type that exhibits the behaviour that an interface wants, but isn't in the right shape.

## Reflect

The behaviour change felt simple, right? OK, maybe it was simply due to the nature of the problem, but this method of work gives you discipline and a simple, repeatable way of changing your system from top to bottom:

- Analyse your problem and identify a slight improvement to your system that pushes you in the right direction
- Capture the new essential complexity in a specification
- Follow the compilation errors until the AT runs
- Update your implementation to make the system behave according to the specification
- Refactor

After the pain of the first iteration, we didn't have to edit our acceptance test code because we have the separation of specifications, drivers and implementation. Changing our specification required us to update our driver and finally our implementation, but the boilerplate code around _how_ to spin up the system as a container was unaffected.

Even with the overhead of building a docker image for our application and spinning up the container, the feedback loop for testing our **entire** application is very tight:

```
quii@Chriss-MacBook-Pro go-specs-greet % go test ./...
ok  	github.com/quii/go-specs-greet	0.181s
ok  	github.com/quii/go-specs-greet/cmd/httpserver	2.221s
?   	github.com/quii/go-specs-greet/specifications	[no test files]
```

Now, imagine your CTO has now decided that gRPC is _the future_. She wants you to expose this same functionality over a gRPC server whilst maintaining the existing HTTP server.

This is an example of **accidental complexity**. Remember, accidental complexity is the complexity we have to deal with because we're working with computers, stuff like networks, disks, APIs, etc. **The essential complexity has not changed**, so we shouldn't have to change our specifications.

Many repository structures and design patterns are mainly dealing with separating types of complexity. For instance, "ports and adapters" ask that you separate your domain code from anything to do with accidental complexity; that code lives in an "adapters" folder.

### Making the change easy

Sometimes, it makes sense to do some refactoring _before_ making a change.

> First make the change easy, then make the easy change

~Kent Beck

For that reason, let's move our `http` code - `driver.go` and `handler.go` - into a package called `httpserver` within an `adapters` folder and change their package names to `httpserver`.

You'll now need to import the root package into `handler.go` to refer to the Greet method...

```go
package httpserver

import (
    "fmt"
    "net/http"

    go_specs_greet "github.com/quii/go-specs-greet/domain/interactions"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, go_specs_greet.Greet(name))
}

```

import your httpserver adapter into main.go:

```go
package main

import (
	"net/http"

	 "github.com/quii/go-specs-greet/adapters/httpserver"
)

func main() {
	handler := http.HandlerFunc(httpserver.Handler)
	http.ListenAndServe(":8080", handler)
}

```

and update the import and reference to `Driver` in greeter_server_test.go:

```go
import (
	...
	"github.com/quii/go-specs-greet/adapters/httpserver"
	...
)

	  ...
	  driver := httpserver.Driver{BaseURL: "http://localhost:8080", Client: &client}
	  ...

```

Finally, it's helpful to gather our domain level code in to its own folder too. Don't be lazy and have a `domain` folder in your projects with hundreds of unrelated types and functions. Make an effort to think about your domain and group ideas that belong together, together. This will make your project easier to understand and will improve the quality of your imports.

Rather than seeing

```go
domain.Greet
```

Which is just a bit weird, instead favour

```go
interactions.Greet
```

Create a `domain` folder to house all your domain code, and within it, an `interactions` folder. Depending on your tooling, you may have to update some imports and code.

Our project tree should now look like this:

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── Dockerfile
├── Makefile
├── README.md
├── adapters
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   └── httpserver
│       ├── greeter_server_test.go
│       └── main.go
├── domain
│   └── interactions
│       ├── greet.go
│       └── greet_test.go
├── go.mod
├── go.sum
└── specifications
    └── adapters.go
    └── greet.go

```

Our domain code, **essential complexity**, lives at the root of our go module, and code that will allow us to use them in "the real world" are organised into **adapters**. The `cmd` folder is where we can compose these logical groupings into practical applications, which have black-box tests to verify it all works. Nice!

Finally, we can do a _tiny_ bit of tidying up our acceptance test. If you consider the high-level steps of our acceptance test:

- Build a docker image
- Wait for it to be listening on _some_ port
- Create a driver that understands how to translate the DSL into system specific calls
- Plug in the driver into the specification

... you'll realise we have the same requirements for an acceptance test for the gRPC server!

The `adapters` folder seems a good place as any, so inside a file called `docker.go`, encapsulate the first two steps in a function that we'll reuse next.

```go
package adapters

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartDockerServer(
	t testing.TB,
	port string,
	dockerFilePath string,
) {
	ctx := context.Background()
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: dockerFilePath,
			PrintBuildLog: true,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

This gives us an opportunity to clean up our acceptance test a little

```go
func TestGreeterServer(t *testing.T) {
	var (
		port           = "8080"
		dockerFilePath = "./cmd/httpserver/Dockerfile"
		baseURL        = fmt.Sprintf("http://localhost:%s", port)
		driver         = httpserver.Driver{BaseURL: baseURL, Client: &http.Client{
			Timeout: 1 * time.Second,
		}}
	)

	adapters.StartDockerServer(t, port, dockerFilePath)
	specifications.GreetSpecification(t, driver)
}
```

This should make writing the _next_ test simpler.

## Write the test first

This new functionality can be accomplished by creating a new `adapter` to interact with our domain code. For that reason we:

- Shouldn't have to change the specification;
- Should be able to reuse the specification;
- Should be able to reuse the domain code.

Create a new folder `grpcserver` inside `cmd` to house our new program and the corresponding acceptance test. Inside `cmd/grpc_server/greeter_server_test.go`, add an acceptance test, which looks very similar to our HTTP server test, not by coincidence but by design.

```go
package main_test

import (
	"fmt"
	"testing"

	"github.com/quii/go-specs-greet/adapters"
	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	var (
		port           = "50051"
		dockerFilePath = "./cmd/grpcserver/Dockerfile"
		driver         = grpcserver.Driver{Addr: fmt.Sprintf("localhost:%s", port)}
	)

	adapters.StartDockerServer(t, port, dockerFilePath)
	specifications.GreetSpecification(t, &driver)
}
```

The only differences are:

- We use a different docker file, because we're building a different program
- This means we'll need a new `Driver`, that'll use `gRPC` to interact with our new program

## Try to run the test

```
./greeter_server_test.go:26:12: undefined: grpcserver
```

We haven't created a `Driver` yet, so it won't compile.

## Write the minimal amount of code for the test to run and check the failing test output

Create a `grpcserver` folder inside `adapters` and inside it create `driver.go`

```go
package grpcserver

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	return "", nil
}
```

If you run again, it should now _compile_ but not pass because we haven't created a Dockerfile and corresponding program to run.

Create a new `Dockerfile` inside `cmd/grpcserver`.

```dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/grpcserver/*.go

EXPOSE 8080
CMD [ "./svr" ]
```

And a `main.go`

```go
package main

import "fmt"

func main() {
	fmt.Println("implement me")
}
```

You should find now that the test fails because our server is not listening on the port. Now is the time to start building our client and server with gRPC.

## Write enough code to make it pass

### gRPC

If you're unfamiliar with gRPC, I'd start by looking at the [gRPC website](https://grpc.io). Still, for this chapter, it's just another kind of adapter into our system, a way for other systems to call (**r**emote **p**rocedure **c**all) our excellent domain code.

The twist is you define a "service definition" using Protocol Buffers. You then generate server and client code from the definition. This not only works for Go but for most mainstream languages too. This means you can share a definition with other teams in your company who may not even write Go and can still do service-to-service communication smoothly.

If you haven't used gRPC before, you'll need to install a **Protocol buffer compiler** and some **Go plugins**. [The gRPC website has clear instructions on how to do this](https://grpc.io/docs/languages/go/quickstart/).

Inside the same folder as our new driver, add a `greet.proto` file with the following

```protobuf
syntax = "proto3";

option go_package = "github.com/quii/adapters/grpcserver";

package grpcserver;

service Greeter {
  rpc Greet (GreetRequest) returns (GreetReply) {}
}

message GreetRequest {
  string name = 1;
}

message GreetReply {
  string message = 1;
}
```

To understand this definition, you don't need to be an expert in Protocol Buffers. We define a service with a Greet method and then describe the incoming and outgoing message types.

Inside `adapters/grpcserver` run the following to generate the client and server code

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet.proto
```

If it worked, we would have some code generated for us to use. Let's start by using the generated client code inside our `Driver`.

```go
package grpcserver

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	//todo: we shouldn't redial every time we call greet, refactor out when we're green
	conn, err := grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := NewGreeterClient(conn)
	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}
```

Now that we have a client, we need to update our `main.go` to create a server. Remember, at this point; we're just trying to get our test to pass and not worrying about code quality.

```go
package main

import (
	"context"
	"log"
	"net"

	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	grpcserver.RegisterGreeterServer(s, &GreetServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: "fixme"}, nil
}
```

To create our gRPC server, we have to implement the interface it generated for us

```go
// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	Greet(context.Context, *GreetRequest) (*GreetReply, error)
	mustEmbedUnimplementedGreeterServer()
}
```

Our `main` function:

- Listens on a port
- Creates a `GreetServer` that implements the interface, and then registers it with `grpcServer.RegisterGreeterServer`, along with a `grpc.Server`.
- Uses the server with the listener

It wouldn't be a massive extra effort to call our domain code inside `greetServer.Greet` rather than hard-coding `fix-me` in the message, but I'd like to run our acceptance test first to see if everything is working on a transport level and verify the failing test output.

```
greet.go:16: Expected values to be equal:
-fixme
\ No newline at end of file
+Hello, Mike
\ No newline at end of file
```

Nice! We can see our driver is able to connect to our gRPC server in the test.

Now, call our domain code inside our `GreetServer`

```go
type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: interactions.Greet(request.Name)}, nil
}
```

Finally, it passes! We have an acceptance test that proves our gRPC greet server behaves how we'd like.

## Refactor

We committed several sins to get the test passing, but now they're passing, we have the safety net to refactor.

### Simplify main

As before, we don't want `main` to have too much code inside it. We can move our new `GreetServer` into `adapters/grpcserver` as that's where it should live. In terms of cohesion, if we change the service definition, we want the "blast-radius" of change to be confined to that area of our code.

### Don't redial in our driver every time

We only have one test, but if we expand our specification (we will), it doesn't make sense for the Driver to redial for every RPC call.

```go
package grpcserver

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string

	connectionOnce sync.Once
	conn           *grpc.ClientConn
	client         GreeterClient
}

func (d *Driver) Greet(name string) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}

	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}

func (d *Driver) getClient() (GreeterClient, error) {
	var err error
	d.connectionOnce.Do(func() {
		d.conn, err = grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		d.client = NewGreeterClient(d.conn)
	})
	return d.client, err
}
```

Here we're showing how we can use [`sync.Once`](https://pkg.go.dev/sync#Once) to ensure our `Driver` only attempts to create a connection to our server once.

Let's take a look at the current state of our project structure before moving on.

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── Makefile
├── README.md
├── adapters
│   ├── docker.go
│   ├── grpcserver
│   │   ├── driver.go
│   │   ├── greet.pb.go
│   │   ├── greet.proto
│   │   ├── greet_grpc.pb.go
│   │   └── server.go
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   ├── grpcserver
│   │   ├── Dockerfile
│   │   ├── greeter_server_test.go
│   │   └── main.go
│   └── httpserver
│       ├── Dockerfile
│       ├── greeter_server_test.go
│       └── main.go
├── domain
│   └── interactions
│       ├── greet.go
│       └── greet_test.go
├── go.mod
├── go.sum
└── specifications
    └── greet.go
```

- `adapters` have cohesive units of functionality grouped together
- `cmd` holds our applications and corresponding acceptance tests
- Our code is totally decoupled from any accidental complexity

### Consolidating `Dockerfile`

You've probably noticed the two `Dockerfiles` are almost identical beyond the path to the binary we wish to build.

`Dockerfiles` can accept arguments to let us reuse them in different contexts, which sounds perfect. We can delete our 2 Dockerfiles and instead have one at the root of the project with the following

```dockerfile
FROM golang:1.18-alpine

WORKDIR /app

ARG bin_to_build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/${bin_to_build}/main.go

CMD [ "./svr" ]
```

We'll have to update our `StartDockerServer` function to pass in the argument when we build the images

```go
func StartDockerServer(
	t testing.TB,
	port string,
	binToBuild string,
) {
	ctx := context.Background()
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "Dockerfile",
			BuildArgs: map[string]*string{
				"bin_to_build": &binToBuild,
			},
			PrintBuildLog: true,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

And finally, update our tests to pass in the image to build (do this for the other test and change `grpcserver` to `httpserver`).

```go
func TestGreeterServer(t *testing.T) {
	var (
		port   = "50051"
		driver = grpcserver.Driver{Addr: fmt.Sprintf("localhost:%s", port)}
	)

	adapters.StartDockerServer(t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
}
```

### Separating different kinds of tests

Acceptance tests are great in that they test the whole system works from a pure user-facing, behavioural POV, but they do have their downsides compared to unit tests:

- Slower
- Quality of feedback is often not as focused as a unit test
- Doesn't help you with internal quality, or design

[The Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html) guides us on the kind of mix we want for our test suite, you should read Fowler's post for more detail, but the very simplistic summary for this post is "lots of unit tests and a few acceptance tests".

For that reason, as a project grows you often may be in situations where the acceptance tests can take a few minutes to run. To offer a friendly developer experience for people checking out your project, you can enable developers to run the different kinds of tests separately.

It's preferable that running `go test ./...` should be runnable with no further set up from an engineer, beyond say a few key dependencies such as the Go compiler (obviously) and perhaps Docker.

Go provides a mechanism for engineers to run only "short" tests with the [short flag](https://pkg.go.dev/testing#Short)

`go test -short ./...`

We can add to our acceptance tests to see if the user wants to run our acceptance tests by inspecting the value of the flag

```go
if testing.Short() {
  t.Skip()
}
```

I made a `Makefile` to show this usage

```makefile
build:
	golangci-lint run
	go test ./...

unit-tests:
	go test -short ./...
```

### When should I write acceptance tests?

The best practice is to favour having lots of fast running unit tests and a few acceptance tests, but how do you decide when you should write an acceptance test, vs unit tests?

It's difficult to give a concrete rule, but the questions I typically ask myself are:

- Is this an edge case? I'd prefer to unit test those
- Is this something that the non-computer people talk about a lot? I would prefer to have a lot of confidence the key thing "really" works, so I'd add an acceptance test
- Am I describing a user journey, rather than a specific function? Acceptance test
- Would unit tests give me enough confidence? Sometimes you're taking an existing journey that already has an acceptance test, but you're adding other functionality to deal with different scenarios due to different inputs. In this case, adding another acceptance test adds a cost but brings little value, so I'd prefer some unit tests.

## Iterating on our work

With all this effort, you'd hope extending our system will now be simple. Making a system that is simple to work on, is not necessarily easy, but it's worth the time, and is substantially easier to do when you start a project.

Let's extend our API to include a "curse" functionality.

## Write the test first

This is brand-new behaviour, so we should start with an acceptance test. In our specification file, add the following

```go
type MeanGreeter interface {
	Curse(name string) (string, error)
}

func CurseSpecification(t *testing.T, meany MeanGreeter) {
	got, err := meany.Curse("Chris")
	assert.NoError(t, err)
	assert.Equal(t, got, "Go to hell, Chris!")
}
```

Pick one of our acceptance tests and try to use the specification

```go
func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port   = "50051"
		driver = grpcserver.Driver{Addr: fmt.Sprintf("localhost:%s", port)}
	)

	t.Cleanup(driver.Close)
	adapters.StartDockerServer(t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
	specifications.CurseSpecification(t, &driver)
}
```

## Try to run the test

```
# github.com/quii/go-specs-greet/cmd/grpcserver_test [github.com/quii/go-specs-greet/cmd/grpcserver.test]
./greeter_server_test.go:27:39: cannot use &driver (value of type *grpcserver.Driver) as type specifications.MeanGreeter in argument to specifications.CurseSpecification:
	*grpcserver.Driver does not implement specifications.MeanGreeter (missing Curse method)
```

Our `Driver` doesn't support `Curse` yet.

## Write the minimal amount of code for the test to run and check the failing test output

Remember we're just trying to get the test to run, so add the method to `Driver`

```go
func (d *Driver) Curse(name string) (string, error) {
	return "", nil
}
```

If you try again, the test should compile, run, and fail

```
greet.go:26: Expected values to be equal:
+Go to hell, Chris!
\ No newline at end of file
```

## Write enough code to make it pass

We'll need to update our protocol buffer specification have a `Curse` method on it, and then regenerate our code.

```protobuf
service Greeter {
  rpc Greet (GreetRequest) returns (GreetReply) {}
  rpc Curse (GreetRequest) returns (GreetReply) {}
}
```

You could argue that reusing the types `GreetRequest` and `GreetReply` is inappropriate coupling, but we can deal with that in the refactoring stage. As I keep stressing, we're just trying to get the test passing, so we verify the software works, _then_ we can make it nice.

Re-generate our code with (inside `adapters/grpcserver`).

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet.proto
```

### Update driver

Now the client code has been updated, we can now call `Curse` in our `Driver`

```go
func (d *Driver) Curse(name string) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}

	greeting, err := client.Curse(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}
```

### Update server

Finally, we need to add the `Curse` method to our `Server`

```go
package grpcserver

import (
	"context"
	"fmt"

	"github.com/quii/go-specs-greet/domain/interactions"
)

type GreetServer struct {
	UnimplementedGreeterServer
}

func (g GreetServer) Curse(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: fmt.Sprintf("Go to hell, %s!", request.Name)}, nil
}

func (g GreetServer) Greet(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: interactions.Greet(request.Name)}, nil
}
```

The tests should now pass.

## Refactor

Try doing this yourself.

- Extract the `Curse` "domain logic", away from the grpc server, as we did for `Greet`. Use the specification as a unit test against your domain logic
- Have different types in the protobuf to ensure the message types for `Greet` and `Curse` are decoupled.

## Implementing `Curse` for the HTTP server

Again, an exercise for you, the reader. We have our domain-level specification and our domain-level logic neatly separated. If you've followed this chapter, this should be very straightforward.

- Add the specification to the existing acceptance test for the HTTP server
- Update your `Driver`
- Add the new endpoint to the server, and reuse the domain code to implement the functionality. You may wish to use `http.NewServeMux` to handle the routeing to the separate endpoints.

Remember to work in small steps, commit and run your tests frequently. If you get really stuck [you can find my implementation on GitHub](https://github.com/quii/go-specs-greet).

## Enhance both systems by updating the domain logic with a unit test

As mentioned, not every change to a system should be driven via an acceptance test. Permutations of business rules and edge cases should be simple to drive via a unit test if you have separated concerns well.

Add a unit test to our `Greet` function to default the `name` to `World` if it is empty. You should see how simple this is, and then the business rules are reflected in both applications for "free".

## Wrapping up

Building systems with a reasonable cost of change requires you to have ATs engineered to help you, not become a maintenance burden. They can be used as a means of guiding, or as a GOOS says, "growing" your software methodically.

Hopefully, with this example, you can see our application's predictable, structured workflow for driving change and how you could use it for your work.

You can imagine talking to a stakeholder who wants to extend the system you work on in some way. Capture it in a domain-centric, implementation-agnostic way in a specification, and use it as a north star towards your efforts. Riya and I describe leveraging BDD techniques like "Example Mapping" [in our GopherconUK talk](https://www.youtube.com/watch?v=ZMWJCk_0WrY) to help you understand the essential complexity more deeply and allow you to write more detailed and meaningful specifications.

Separating essential and accidental complexity concerns will make your work less ad-hoc and more structured and deliberate; this ensures the resiliency of your acceptance tests and helps them become less of a maintenance burden.

Dave Farley gives an excellent tip:

> Imagine the least technical person that you can think of, who understands the problem-domain, reading your Acceptance Tests. The tests should make sense to that person.

Specifications should then double up as documentation. They should specify clearly how a system should behave. This idea is the principle around tools like [Cucumber](https://cucumber.io), which offers you a DSL for capturing behaviours as code, and then you convert that DSL into system calls, just like we did here.

### What has been covered

- Writing abstract specifications allows you to express the essential complexity of the problem you're solving and remove accidental complexity. This will enable you to reuse the specifications in different contexts.
- How to use [Testcontainers](https://golang.testcontainers.org) to manage the life-cycle of your system for ATs. This allows you to thoroughly test the image you intend to ship on your computer, giving you fast feedback and confidence.
- A brief intro into containerising your application with Docker
- gRPC
- Rather than chasing canned folder structures, you can use your development approach to naturally drive out the structure of your application, based on your own needs

### Further material

- In this example, our "DSL" is not much of a DSL; we just used interfaces to decouple our specification from the real world and allow us to express domain logic cleanly. As your system grows, this level of abstraction might become clumsy and unclear. [Read into the "Screenplay Pattern"](https://cucumber.io/blog/bdd/understanding-screenplay-(part-1)/) if you want to find more ideas as to how to structure your specifications.
- For emphasis, [Growing Object-Oriented Software, Guided by Tests,](http://www.growing-object-oriented-software.com) is a classic. It demonstrates applying this "London style", "top-down" approach to writing software. Anyone who has enjoyed Learn Go with Tests should get much value from reading GOOS.
- [In the example code repository](https://github.com/quii/go-specs-greet), there's more code and ideas I haven't written about here, such as multi-stage docker build, you may wish to check this out.
  - In particular, *for fun*, I made a **third program**, a website with some HTML forms to `Greet` and `Curse`. The `Driver` leverages the excellent-looking [https://github.com/go-rod/rod](https://github.com/go-rod/rod) module, which allows it to work with the website with a browser, just like a user would. Looking at the git history, you can see how I started not using any templating tools "just to make it work" Then, once I passed my acceptance test, I had the freedom to do so without fear of breaking things.


