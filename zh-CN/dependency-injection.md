# 依赖注入

**[您可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/di)**

因为我们需要懂得一些接口的知识，所以我们假设你已经阅读了前面的结构体篇。

在编程社区里，对于**依赖注入**（dependency injection）存在诸多误解。我们希望本篇会向你展示为什么：

* 你不需要一个框架
* 它不会过度复杂化你的设计
* 它易于测试
* 它能让你编写优秀和通用的函数

就像我们在 hello-world 篇做的那样，我们想要编写一个欢迎某人的函数，只不过这次我们希望测试实际的打印（actual printing）。

回顾一下，这个函数应该长这个样子：

```go
func Greet(name string) {
	fmt.Printf("Hello, %s", name)
}
```

那么我们该如何测试它呢？调用 `fmt.Printf` 会打印到标准输出，用测试框架来捕获它会非常困难。

我们所需要做的就是**注入**（这只是一个等同于“传入”的好听的词）打印的依赖。

**我们的函数不需要关心在哪里打印，以及如何打印，所以我们应该接收一个接口，而非一个具体的类型**。

如果我们这样做的话，就可以通过改变接口的实现，控制打印的内容，于是就能测试它了。在实际情况中，你可以注入一些写入标准输出的内容。

如果你看看 `fmt.Printf` 的源码，你可以找到一种挂钩的方式：

```go
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}
```

有意思！在 `Printf` 内部，只是传入 `os.Stdout`，并调用了 `Fprintf`。

`os.Stdout` 究竟是什么？`Fprintf` 期望第一个参数传递过来什么？

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	p := newPrinter()
	p.doPrintf(format, a)
	n, err = w.Write(p.buf)
	p.free()
	return
}
```

`io.Writer` 是：

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

如果你写过很多 Go 代码的话，你会发现这个接口出现的频率很高，因为 `io.Writer` 是一个很好的通用接口，用于“将数据放在某个地方”。

所以我们知道了，在幕后我们其实是用 `Writer` 来把问候发送到某处。我们现在来使用这个抽象，让我们的代码可以测试，并且重用性更好。

## 测试优先

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

`bytes` 包中的 `buffer` 类型实现了 `Writer` 接口。

因此，我们可以在测试中，用它来作为我们的 `Writer`，接着调用了 `Greet` 后，我们可以用它来检查写入了什么。

## 尝试运行测试

这个测试编译会报错：

```bash
./di_test.go:10:7: too many arguments in call to Greet
	have (*bytes.Buffer, string)
	want (string)
```

## 编写最小化代码供测试运行，并检查失败的测试输出

根据编译器，修复问题。

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Printf("Hello, %s", name)
}
```

```bash
Hello, Chris di_test.go:16: got '' want 'Hello, Chris'
```

测试失败了。注意到可以打印出 `name`，不过它传入到了标准输出。

## 编写代码使其通过

用 `writer` 把问候发送到我们测试中的缓冲区。记住 `fmt.Fprintf` 和 `fmt.Printf` 一样，只不过 `fmt.Fprintf` 会接收一个 `Writer` 参数，用于把字符串传递过去，而 `fmt.Printf` 默认是标准输出。

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}
```

现在测试就可以通过了。

## 重构

早些时候，编译器会告诉我们需要传入一个指向 `bytes.Buffer` 的指针。这在技术上是正确的，但却不是很有用。

为了展示这一点，我们把 `Greet` 函数接入到一个 Go 应用里面，其中我们会打印到标准输出。

```go
func main() {
	Greet(os.Stdout, "Elodie")
}
```

```
./di.go:14:7: cannot use os.Stdout (type *os.File) as type *bytes.Buffer in argument to Greet
```

我们前面讨论过，`fmt.Fprintf` 允许传入一个 `io.Writer` 接口，我们知道 `os.Stdout` 和 `bytes.Buffer` 都实现了它。

我们可以修改一下代码，使用更为通用的接口，于是我们现在可以在测试和应用中都使用这个函数了。

```go
package main

import (
	"fmt"
	"os"
	"io"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func main() {
	Greet(os.Stdout, "Elodie")
}
```

## 关于 io.Writer 的更多内容

通过使用 `io.Writer`，我们还可以将数据写入哪些地方？我们的 `Greet` 函数的通用性怎么样了？

### 互联网

运行下面代码：

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

前往 [http://localhost:5000](http://localhost:5000)。你会看到使用到了你的 `greeting` 函数。

在下一章会介绍 HTTP 服务器，所以不要太担心这些细节。

当你编写一个 HTTP 处理器（handler）时，你需要给出 `http.ResponseWriter` 和用于创建请求的 `http.Request`。在你实现服务器时，你使用 `writer` **写入**了请求。

你可能已经猜到，`http.ResponseWriter` 也实现了 `io.Writer`，所以我们可以重用处理器中的 `Greet` 函数。

## 圆满完成

我们第一轮迭代的代码不易测试，因为它把数据写到了我们无法控制的地方。

通过测试的启发，我们重构了代码。因为有了注入依赖，我们可以控制数据向哪儿写入，它允许我们：

* **测试代码**。如果你不能很轻松地测试函数，这通常是因为有依赖硬链接到了函数或全局状态。例如，如果某个服务层使用了全局的数据库连接池，这通常难以测试，并且运行速度会很慢。DI 提倡你注入一个数据库依赖（通过接口），然后就可以在测试中控制你的模拟数据了。
* **关注点分离**，解耦了**数据到达的地方**和**如何产生数据**。如果你感觉一个方法/函数负责太多功能了（产生数据**并且**写入一个数据库？处理 HTTP 请求**并且**处理业务级别的逻辑），那么你可能就需要 DI 这项工具了。
* **使得代码可以在不同环境下重用**。我们的代码所处的第一个“新”环境就是在运行的测试。但是随后，如果其他人想要用你的代码尝试点新东西，他们只要注入他们自己的依赖就可以了。

### 什么是模拟？我听说 DI 要用到模拟，它可讨厌了

模拟（mocking）会在后面详细讨论（它并不坏）。你会使用模拟来代替真实事物，用一个模拟版本来注入，于是可以控制和检查你的测试。在我们的例子中，标准库已经有工具供我们使用了。

### Go 标准库真的很棒，花时间好好研究它吧

通过熟悉 `io.Writer` 接口，我们可以用测试中的 `bytes.Buffer` 来作为 `Writer`，然后我们可以使用标准库中的其他的 `Writer`，在命令行应用或 web 服务器中使用这个函数。

随着你越来越熟悉标准库，你会越来越了解，这些在代码中重用的通用接口，会使得你的软件在许多场景都可以重用。

本例深深受到 [The Go Programming language](https://www.amazon.co.uk/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440) 中一个章节的启发，如果你喜欢的话，去它买吧！

----------------

作者：[Chris James](https://dev.to/quii)
译者：[Noluye](https://github.com/Noluye)
校对：[rxcai](https://github.com/rxcai)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
