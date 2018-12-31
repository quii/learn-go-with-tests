# Q&A: OS Exec

[https://www.reddit.com/user/keith6014](https://www.reddit.com/user/keith6014) asks on [reddit](https://www.reddit.com/r/golang/comments/aaz8ji/testdata_and_function_setup_help/)

> I am executing a command using os/exec.Command() which generated XML data. The command will be executed in a function called GetData().
> In order to test GetData(), I have some testdata which I created.
> In my _test.go I have a TestGetData which calls GetData() but that will use os.exec, instead I would like for it to use my testdata.
> What is a good way to acheive this? When calling GetData should I have a "test" flag mode so it will read a file ie GetData(mode string)?

A few things

- When something is difficult to test, it's probably due to the separation of concerns not being quite right
- Dont add "test modes" into your code, instead use [Dependency Injection](/dependency-injection.md) so that you can model your dependencies and separate concerns. 

I have taken the liberty of guessing what the code might look like

```go
type Payload struct {
	Message string `xml:"message"`
}

func KeithCommand() string {
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
- The rest of the code is more or less copy and pasted from the [excellent documentation](https://golang.org/pkg/os/exec/#example_Cmd_StdoutPipe) 

Here is what is contained inside `msg.xml`

```xml
<payload>
    <message>Happy New Year!</message>
</payload>
```

I wrote a simple test to show it in action

```go
func TestKeithCommand(t *testing.T) {
	got := KeithCommand()
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
```

## Testable code

Testable code is decoupled and single purpose. To me it feels like there are two main concerns for this code

1. Retrieving the raw XML data
2. Decoding the XML data and applying our business logic (in this case `strings.ToUpper` on the `<message>`)

The first part is _more or less_ boilerplate, it's just copying the example from the standard lib. 

The second part is where we have our business logic and by looking at the code we can see where the "seam" in our logic starts; it's where we get our `io.ReadCloser`. We can use this existing abstraction to separate concerns and make our code testable.

Our `TestKeithCommand` can act as our integration test between our two concerns so we'll keep hold of that and make sure it keeps working.

Here is what the newly separated code looks like

```go
type Payload struct {
	Message string `xml:"message"`
}

func KeithCommand(data io.Reader) string {
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

func TestKeithCommandIntegration(t *testing.T) {
	got := KeithCommand(getXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
```

Now that `KeithCommand` takes it's input from just an `io.Reader` we have made it testable; and people can re-use the function with anything that returns an `io.Reader` (which is extremely common).

```go
func TestKeithCommand(t *testing.T) {
	
	t.Run("an example", func(t *testing.T) {
		
		input := strings.NewReader(`
<payload>
    <message>Cats are the best animal</message>
</payload>`)

		got := KeithCommand(input)
		want := "CATS ARE THE BEST ANIMAL"

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})
}
```

Here is an example of a unit test for `KeithCommand`. 

By separating the concerns and using existing abstractions within Go testing our important business logic is a breeze.
