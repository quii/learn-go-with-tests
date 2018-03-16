# Dependency injection and interfaces - WIP

It is assumed that you have read the structs section before reading this.

There is _a lot_ of misunderstandings around dependency injection around the programming community. Hopefully this guide will show you how

- You dont need a framework
- It does not overcomplicate your design
- It facilitates testing
- It allows you to write great, general-purpose functions. 


We want to write a function that greets someone, just like we did in the hello-world chapter but this time we are going to be testing the _actual printing_. 

Just to recap, here is what that function could look like

```go
func Greet(name string) {
	fmt.Printf("Hello, %s", name)
}
```

But how can we test this? Calling `fmt.Printf` prints to stdout, which is pretty hard for us to capture using the testing framework. 

What we need to do is to be able to **inject** (which is just a fancy word for pass in) the dependency of printing. If we do that, we can then change the implementation to print to something we control so that we can test it. In "real life" you would inject in something that writes to stdout.  

If you look at the source code of `fmt.Printf` you can see a way for us to hook in

```go
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}
```

Interesting! Under the hood `Printf` just calls `Fprintf` passing in `os.Stdout`.

What exactly _is_ an `os.Stdout` ? What does `Fprintf` expect to get passed to it in 1st argument?

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	p := newPrinter()
	p.doPrintf(format, a)
	n, err = w.Write(p.buf)
	p.free()
	return
}
```

An `io.Writer`

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

As you write more Go code you will find this interface popping up a lot because it's a great general purpose interface for "put this data somewhere".

Now that we know this, we can write a test!

## Write the test first

```go
func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer,"Chris")

	got := buffer.String()
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}
```

The `buffer` type from the `bytes` package implements the `Writer` interface so it is perfect for us to try and record what is being written.

We call the `Greet` function and afterwards read the buffer into a `string` so we can assert on it.

## Try and run the test

The test will not compile

```
./di_test.go:10:7: too many arguments in call to Greet
	have (*bytes.Buffer, string)
	want (string)
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Printf("Hello, %s", name)
}
```

`Hello, Chris	di_test.go:16: got '' want 'Hello, Chris'`

The test fails. Notice that the name is getting printed out, but it's going to stdout.

## Write enough code to make it pass

Use the writer to send the greeting to the buffer in our test

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Fprintf(writer,"Hello, %s", name)
}
```

The test now pass

## Refactor

Earlier the compiler told us to pass in a pointer to a `bytes.Buffer`. This is technically correct but not very useful. 

To demonstrate this, try wiring up the `Greet` function into a Go application where we want it to print to stdout.

```go
func main() {
	Greet(os.Stdout, "Elodie")
}
```

`./di.go:14:7: cannot use os.Stdout (type *os.File) as type *bytes.Buffer in argument to Greet`

As discussed earlier `fmt.Fprintf` allows you to pass in an `io.Writer` which we know both `os.Stdout` and `bytes.Buffer` both implement.

If we change our code to use the more general purpose interface we can then use it in both tests and in our application

```go
package main

import (
	"fmt"
	"os"
	"io"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer,"Hello, %s", name)
}

func main() {
	Greet(os.Stdout, "Elodie")
}
```

## More on io.Writer

What other places can we write data to using `io.Writer` ? Just how general purpose is our `Greet` function?

### The internet

Run the following

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "world")
}

func main() {
	http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))
}
```

Go to [http://localhost:5000](http://localhost:5000). You'll see your greeting function being used. 

HTTP servers will be covered in a later chapter so dont worry too much about the details. 

When you write a HTTP handler, you are given a `http.ResponseWriter` and the `http.Request` that was used to make the request. When you implement your server you _write_ your response using the writer. 

You can probably guess that `http.ResponseWriter` also implements `io.Writer` so this is why we could re-use our `Greet` function inside our handler.

The point I am trying to make is that by injecting the dependency of "where to write the data" we have made a function that can be used to write to files, stdout, the internet and lots more. 

## Wrapping up 

wip wip

