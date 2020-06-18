# OS Exec

**[You can find all the code here](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/os-exec)**

[keith6014](https://www.reddit.com/user/keith6014) asks on [reddit](https://www.reddit.com/r/golang/comments/aaz8ji/testdata_and_function_setup_help/)

> I am executing a command using os/exec.Command() which generated XML data. The command will be executed in a function called GetData().

> In order to test GetData(), I have some testdata which I created.

> In my _test.go I have a TestGetData which calls GetData() but that will use os.exec, instead I would like for it to use my testdata.

> What is a good way to achieve this? When calling GetData should I have a "test" flag mode so it will read a file ie GetData(mode string)?

A few things

- When something is difficult to test, it's often due to the separation of concerns not being quite right
- Don't add "test modes" into your code, instead use [Dependency Injection](/dependency-injection.md) so that you can model your dependencies and separate concerns.

I have taken the liberty of guessing what the code might look like

```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData() string {
	cmd := exec.Command("cat", "msg.xml")

	out, _ := cmd.StdoutPipe()
	var payload Payload
	decoder := xml.NewDecoder(out)

	// these 3 can return errors but I'm ignoring for brevity
	cmd.Start()
	decoder.Decode(&payload)
	cmd.Wait()

	return strings.ToUpper(payload.Message)
}
```

- It uses `exec.Command` which allows you to execute an external command to the process
- We capture the output in `cmd.StdoutPipe` which returns us a `io.ReadCloser` (this will become important)
- The rest of the code is more or less copy and pasted from the [excellent documentation](https://golang.org/pkg/os/exec/#example_Cmd_StdoutPipe).
    - We capture any output from stdout into an `io.ReadCloser` and then we `Start` the command and then wait for all the data to be read by calling `Wait`. In between those two calls we decode the data into our `Payload` struct.

Here is what is contained inside `msg.xml`

```xml
<payload>
    <message>Happy New Year!</message>
</payload>
```

I wrote a simple test to show it in action

```go
func TestGetData(t *testing.T) {
	got := GetData()
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

## Testable code

Testable code is decoupled and single purpose. To me it feels like there are two main concerns for this code

1. Retrieving the raw XML data
2. Decoding the XML data and applying our business logic (in this case `strings.ToUpper` on the `<message>`)

The first part is just copying the example from the standard lib.

The second part is where we have our business logic and by looking at the code we can see where the "seam" in our logic starts; it's where we get our `io.ReadCloser`. We can use this existing abstraction to separate concerns and make our code testable.

**The problem with GetData is the business logic is coupled with the means of getting the XML. To make our design better we need to decouple them**

Our `TestGetData` can act as our integration test between our two concerns so we'll keep hold of that to make sure it keeps working.

Here is what the newly separated code looks like

```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := ioutil.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

func TestGetDataIntegration(t *testing.T) {
	got := GetData(getXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

Now that `GetData` takes its input from just an `io.Reader` we have made it testable and it is no longer concerned how the data is retrieved; people can re-use the function with anything that returns an `io.Reader` (which is extremely common). For example we could start fetching the XML from a URL instead of the command line.

```go
func TestGetData(t *testing.T) {
	input := strings.NewReader(`
<payload>
    <message>Cats are the best animal</message>
</payload>`)

	got := GetData(input)
	want := "CATS ARE THE BEST ANIMAL"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

```

Here is an example of a unit test for `GetData`.

By separating the concerns and using existing abstractions within Go testing our important business logic is a breeze.
