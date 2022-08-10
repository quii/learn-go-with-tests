# Introduction to acceptance testing

At `$WORK`, we've been running into the need to have "graceful shutdown" for our services. Graceful shutdown is making sure your system finishes its work properly before it is terminated. A real-world analogy would be someone trying to wrap up a phone call properly before moving on to the next meeting, rather than just hanging up mid-sentence.

This chapter will give an intro to graceful shutdown in the context of an HTTP server, and how to write "acceptance tests" to give yourself confidence in the behaviour of your code.

After reading this you'll know how to share packages with excellent tests, reduce maintenance efforts and, increase confidence in the quality of your work.

## Just enough info about Kubernetes

We run our software on [Kubernetes](https://kubernetes.io/) (K8s). K8s will terminate "pods" (in practice, our software) for various reasons, and a common one is when we push new code that we want to deploy.

We are setting ourselves high standards regarding [DORA metrics](https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance), so we work in a way where we deploy small, incremental improvements and features to production multiple times per day.

When k8s wishes to terminate a pod, it initiates a ["termination lifecycle"](https://cloud.google.com/blog/products/containers-kubernetes/kubernetes-best-practices-terminating-with-grace), and a part of that is sending a SIGTERM signal to our software. This is k8s telling our code:

> You need to shut yourself down, finish whatever work you're doing because after a certain "grace period", I will send SIGKILL, and it's lights out for you.

On `SIGKILL` any work your program might've been doing will be immediately stopped.

## If you do not have grace

Depending on the nature of your software, if you ignore `SIGTERM`, you can run into problems.

Our specific problem was with in-flight HTTP requests. When an automated test was exercising our API, if k8s decided to stop the pod, the server would die, the test would not get a response from the server, and the test will fail.

This would trigger an alert in our incidents channel which requires a dev to stop what they're doing and address the problem. These intermittent failures are an annoying distraction for our team.

These problems are not unique to our tests. If a user sends a request to your system and the process gets terminated mid-flight, they'll likely be greeted with a 5xx HTTP error, not the kind of user experience you want to deliver.

## When you have grace

What we want to do is listen for `SIGTERM`, and rather than instantly killing the server, we want to:

- Stop listening to any more requests;
- Allow any in-flight requests to finish;
- *Then* terminate the process.

## How to have grace

Thankfully, Go already has a mechanism for gracefully shutting down a server with [net/http/Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown).

> Shutdown gracefully shuts down the server without interrupting any active connections. Shutdown works by first closing all open listeners, then closing all idle connections, and then waiting indefinitely for connections to return to idle and then shut down. If the provided context expires before the shutdown is complete, Shutdown returns the context's error, otherwise it returns any error returned from closing the Server's underlying Listener(s).

To handle `SIGTERM` we can use [os/signal.Notify](https://pkg.go.dev/os/signal#Notify), which will send any incoming signals to a channel we provide.

By using these two features from the standard library, you can listen for `SIGTERM` and shutdown gracefully.

## Graceful shutdown package

To that end, I wrote [https://pkg.go.dev/github.com/quii/go-graceful-shutdown](https://pkg.go.dev/github.com/quii/go-graceful-shutdown). It provides a decorator function for a `*http.Server` to call its `Shutdown` method when a `SIGTERM` signal is detected

```go
func main() {
	httpServer := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(acceptancetests.SlowHandler)}

	server := gracefulshutdown.NewServer(httpServer)

	if err := server.ListenAndServe(); err != nil {
		// this will typically happen if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
```

The specifics around the code are not too important for this read, but it is worth having a quick look over the code before carrying on.

## Tests and feedback loops

When the package was written, there were unit tests to prove the `gracefulshutdown` package behaves correctly, and they gave us the confidence to aggressively refactor, but we still didn't feel "confident" that it **really** worked.

We added a `cmd` package and made a real program to use the package we were writing, and we'd manually fire it up, fire off an HTTP request to it, and then send a `SIGTERM` to see what would happen.

**The engineer in you should be feeling uncomfortable with manual testing**, it's boring, doesn't scale, it's inaccurate, and it's wasteful. If you're writing a package you intend to share, but also keep it simple and cheap to change, manual testing is not going to cut it.

## Acceptance tests

If you’ve read the rest of this book, you will have mostly written "unit tests". Unit tests are a fantastic tool for enabling fearless refactoring, driving good modular design, preventing regressions, and facilitating fast feedback.

By their nature, they only test small parts of your system. Usually, unit tests alone are *not enough* for an effective testing strategy. Remember, we want our systems to **always be shippable**. We can't rely on manual testing, so we need another kind of testing; **acceptance tests**.

### What are they?

They are a kind of "black-box test". They are sometimes referred to as "functional tests". They should exercise the system as a user of the system would.

The term "black-box" refers to the idea that the test code has no access to the internals of the system, it can only use its public interface and make assertions on the behaviours it observes. This means they can only test the system as a whole.

This is an advantageous trait because it means the tests exercise the system the same as a user would, it can't use any special workarounds that could make a test pass, but not actually prove what you need to prove. This is similar to the principle of preferring your unit test files to live inside a separate test package, for example, `package mypkg_test` rather than `package mypg`.

### Benefits of acceptance tests

- When acceptance tests pass, you know your entire system behaves how you want it to;
- More accurate, less effort and quicker than manual testing;
- When written well, they act as accurate, verified documentation of your system. It doesn't fall into the trap of documentation that diverges from the real behaviour of the system;
- No mocking! It's all real.

### Potential drawbacks vs unit tests

- Expensive to write;
- Slower;
- Very dependant on the design of the system;
- When they fail, typically don't give you a root cause, and can be difficult to debug;
- Don't give you feedback on the internal quality of your system. You could write total garbage and still make an acceptance test pass;
- Not all scenarios, are practical to exercise due to the black-box nature.

For this reason, it is foolish to only rely on acceptance tests. They do not have many of the qualities unit tests have, and a system with a large number of acceptance tests will tend to suffer in terms of maintenance costs and poor lead time.

#### Lead time?

Lead time refers to how long it takes from a commit being merged into your main branch, to it being deployed in production. This number can vary from weeks and even months for some teams to a matter of minutes. Again at `$WORK`, we value DORA's findings and want to keep our lead time to under 10 minutes.

A balanced testing approach is required for a reliable system with excellent lead time, and this is usually described in terms of [The Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html).

## How to write basic acceptance tests

How does this relate to the original problem? We've just written a package here, and it is entirely unit-testable.

As I mentioned, the unit tests weren't quite giving us the confidence we needed. We want to be *really* sure the package works when integrated with a real, running program. We should be able to automate the manual checks we were making.

Let's take a look at the test program

```
func main() {
	httpServer := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(acceptancetests.SlowHandler)}

	server := gracefulshutdown.NewServer(httpServer)

	if err := server.ListenAndServe(); err != nil {
		// this will typically happen if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didnt shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
```

You may have guessed that `SlowHandler` has a `time.Sleep` to delay responding so I had time to `SIGTERM` and see what happens. The rest is fairly boilerplate:

- Make a `net/http/Server`;
- Wrap it in the library (see: [Decorator pattern](https://en.wikipedia.org/wiki/Decorator_pattern));
- Use the wrapped version to `ListenAndServe`.

### High-level steps for an acceptance test

- Build the program;
- Run it (and wait for it listen on `8080`);
- Send an HTTP request to the server;
- Before the server has a chance to respond to the rest, send SIGTERM;
- See if we still get a response.

### Building and running the program

```
package acceptancetests

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	baseBinName = "temp-testbinary"
)

func LaunchTestProgram(port string) (cleanup func(), sendInterrupt func() error, err error) {
	binName, err := buildBinary()
	if err != nil {
		return nil, nil, err
	}

	sendInterrupt, kill, err := runServer(binName, port)

	cleanup = func() {
		if kill != nil {
			kill()
		}
		os.Remove(binName)
	}

	if err != nil {
		cleanup() // even though it's not listening correctly, the program could still be running
		return nil, nil, err
	}

	return cleanup, sendInterrupt, nil
}

func buildBinary() (string, error) {
	binName := randomString(10) + "-" + baseBinName

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		return "", fmt.Errorf("cannot build tool %s: %s", binName, err)
	}
	return binName, nil
}

func runServer(binName string, port string) (sendInterrupt func() error, kill func(), err error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	cmdPath := filepath.Join(dir, binName)

	cmd := exec.Command(cmdPath)

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("cannot run temp converter: %s", err)
	}

	kill = func() {
		_ = cmd.Process.Kill()
	}

	sendInterrupt = func() error {
		return cmd.Process.Signal(syscall.SIGTERM)
	}

	err = waitForServerListening(port)

	return
}

func waitForServerListening(port string) error {
	for i := 0; i < 30; i++ {
		conn, _ := net.Dial("tcp", net.JoinHostPort("localhost", port))
		if conn != nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("nothing seems to be listening on localhost:%s", port)
}
```

`LaunchTestProgram` is responsible for building the program, launching it, waiting for it to listen on port and providing:

- a `cleanup` function to kill the program and delete it, to ensure when our tests finish that we're left in a clean state
- an interrupt function to send the program a `SIGTERM` to let us test the behaviour

Admittedly, this is not the nicest code in the world, but just focus on the exported function LaunchTestProgram, the un-exported functions it calls are uninteresting boilerplate.

As discussed, acceptance testing tends to be trickier to set up. This code does make the *testing* code substantially simpler to read, and often with acceptance tests once you've written the ceremonious code once, it's done and you can forget about it.

### The acceptance test(s)

We wanted to have two acceptance tests for two programs, one with graceful shutdown and one without, so we, and the readers can see the difference in behaviour. With `LaunchTestProgram` to build and run the programs, it's quite simple to write acceptance tests for both, and we benefit from re-use with some helper functions.

Here is the test for the server *with* a graceful shutdown, [you can find the test without on GitHub](https://github.com/quii/go-graceful-shutdown/blob/main/acceptancetests/withoutgracefulshutdown/main_test.go)

```
package main

import (
	"testing"
	"time"

	"github.com/quii/go-graceful-shutdown/acceptancetests"
	"github.com/quii/go-graceful-shutdown/assert"
)

const (
	port = "8080"
	url  = "<http://localhost:"> + port
)

func TestGracefulShutdown(t *testing.T) {
	cleanup, sendInterrupt, err := acceptancetests.LaunchTestProgram(port)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)

	// just check the server works before we shut things down
	assert.CanGet(t, url)

	// fire off a request, and before it has a chance to respond send SIGTERM.
	time.AfterFunc(50*time.Millisecond, func() {
		assert.NoError(t, sendInterrupt())
	})
	// Without graceful shutdown, this would fail
	assert.CanGet(t, url)

	// after interrupt, the server should be shutdown, and no more requests will work
	assert.CantGet(t, url)
}
```

With the setup encapsulated away, the tests are comprehensive, describe the behaviour and are relatively easy to follow.

`assert.CanGet/CantGet` are helper functions I made to DRY up this common assertion for this suite.

```
func CanGet(t testing.TB, url string) {
	errChan := make(chan error)

	go func() {
		res, err := http.Get(url)
		if err != nil {
			errChan <- err
			return
		}
		res.Body.Close()
		errChan <- nil
	}()

	select {
	case err := <-errChan:
		NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for request to %q", url)
	}
}
```

This will fire off a `GET` to `URL` on a goroutine, and if it responds without error before 3 seconds, then it will not fail. `CantGet` is omitted for brevity, [but you can view it on GitHub here](https://github.com/quii/go-graceful-shutdown/blob/main/assert/assert.go#L61).

It's important to note again, Go has all the tools you need to write acceptance tests out of the box. You don't *need* a special framework to build acceptance tests.

### Small investment with a big pay-off

With these tests, readers can look at the example programs and be confident that the example *actually* works, so they can be confident in the package's claims.

Importantly, as the author, we get **fast feedback** and **massive confidence** that the package works in a real-world setting.

```
go test -count=1 ./...
ok  	github.com/quii/go-graceful-shutdown	0.196s
?   	github.com/quii/go-graceful-shutdown/acceptancetests	[no test files]
ok  	github.com/quii/go-graceful-shutdown/acceptancetests/withgracefulshutdown	4.785s
ok  	github.com/quii/go-graceful-shutdown/acceptancetests/withoutgracefulshutdown	2.914s
?   	github.com/quii/go-graceful-shutdown/assert	[no test files]
```

## Wrapping up

In this blog post, we introduced acceptance tests into your testing tool belt. They are invaluable when you start to build real systems and are an important complement to your unit tests.

The nature of *how* to write acceptance tests depends on the system you're building, but the principles stay the same. Treat your system like a "black box". If you're making a website, your tests should act like a user and you'll want to use a headless web browser like [Selenium](https://www.selenium.dev/), to click on links, fill in forms, etc. For a RESTful API, you'll send HTTP requests using a client.

### Taking it further for more complicated systems

Non-trivial systems don't tend to be single-process applications like the one we've discussed. Typically you'll depend on other systems such as a database. For these scenarios, you'll need to automate a local environment to test with. Tools like [docker-compose](https://docs.docker.com/compose/) are useful for spinning up containers of the environment you need to run your system locally.

### The next chapter

In this post the acceptance test was written retrospectively. However, in [Growing Object-Oriented Software](http://www.growing-object-oriented-software.com) the authors show that we can use acceptance tests in a test-driven approach to act as a "north-star" to guide our efforts.

As systems get more complex, the costs of writing and maintaining acceptance tests can quickly spiral out of control. There are countless stories of development teams being hamstrung by expensive acceptance test suites.

The next chapter will introduce using acceptance test to guide our design and principles and techniques for managing the costs of acceptance tests.

### Improving the quality of open-source

If you're writing packages you intend to share, I'd encourage you to create simple example programs demonstrating what your package does and invest time in having simple-to-follow acceptance tests to give yourself, and potential users of your work confidence.

Like [Testable Examples](https://go.dev/blog/examples), seeing this little extra effort in developer experience goes a long way toward building trust in your work, and will reduce your own maintenance costs.

## Recruitment plug for `$WORK`

If you fancy working in an environment with other engineers solving interesting problems, live near or around London or Porto, and enjoy the contents of this chapter and book -  please [reach out to me on Twitter](https://twitter.com/quii) and maybe we can work together soon!
