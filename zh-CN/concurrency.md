# 并发

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/concurrency)**

这是我们的计划：同事已经写了一个 `CheckWebsites` 的函数检查 URL 列表的状态。

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        results[url] = wc(url)
    }

    return results
}
```

它返回一个 map，由每个 url 检查后的得到的布尔值组成，成功响应的值为 `true`，错误响应的值为 `false`。

你还必须传入一个 `WebsiteChecker` 处理单个 URL 并返回一个布尔值。它会被函数调用以检查所有的网站。

使用 [依赖注入](zh-CN/dependency-injection.md)，允许在不发起真实 HTTP 请求的情况下测试函数，这使测试变得可靠和快速。

这是他们写的测试：

```go
package concurrency

import (
    "reflect"
    "testing"
)

func mockWebsiteChecker(url string) bool {
    if url == "waat://furhurterwe.geds" {
        return false
    }
    return true
}

func TestCheckWebsites(t *testing.T) {
    websites := []string{
        "http://google.com",
        "http://blog.gypsydave5.com",
        "waat://furhurterwe.geds",
    }

    actualResults := CheckWebsites(mockWebsiteChecker, websites)

    want := len(websites)
    got := len(actualResults)
    if want != got {
        t.Fatalf("Wanted %v, got %v", want, got)
    }

    expectedResults := map[string]bool{
        "http://google.com":          true,
        "http://blog.gypsydave5.com": true,
        "waat://furhurterwe.geds":    false,
    }

    if !reflect.DeepEqual(expectedResults, actualResults) {
        t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
    }
}
```
该功能在生产环境中被用于检查数百个网站。但是你的同事开始抱怨它速度很慢，所以他们请你帮忙为程序提速。

## 写一个测试

首先我们对 `CheckWebsites` 做一个基准测试，这样就能看到我们修改的影响。

```go
package concurrency

import (
    "testing"
    "time"
)

func slowStubWebsiteChecker(_ string) bool {
    time.Sleep(20 * time.Millisecond)
    return true
}

func BenchmarkCheckWebsites(b *testing.B) {
    urls := make([]string, 100)
    for i := 0; i < len(urls); i++ {
        urls[i] = "a url"
    }

    for i := 0; i < b.N; i++ {
        CheckWebsites(slowStubWebsiteChecker, urls)
    }
}
```

基准测试使用一百个网址的 slice 对 `CheckWebsites` 进行测试，并使用 `WebsiteChecker` 的伪造实现。`slowStubWebsiteChecker` 故意放慢速度。它使用 `time.Sleep` 明确等待 20 毫秒，然后返回 true。

当我们运行基准测试时使用 `go test -bench=.` 命令 (如果在 Windows Powershell 环境下使用 `go test -bench="."`)：

```
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v0
BenchmarkCheckWebsites-4               1        2249228637 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v0        2.268s
```
`CheckWebsite` 经过基准测试的时间为 2249228637 纳秒，大约 2.25 秒。

让我们尝试去让它运行得更快。

## 编写足够的代码让它通过

现在我们终于可以谈论并发了，以下内容是为了说明「不止一件事情正在进行中」。这是我们每天很自然在做的事情。

比如，今天早上我泡了一杯茶。我放上水壶，然后在等待它煮沸时，从冰箱里取出了牛奶，把茶从柜子里拿出来，找到我最喜欢的杯子，把茶袋放进杯子里，然后等水壶沸了，把水倒进杯子里。

我 *没有* 做的事情是放上水壶，然后呆呆地盯着水壶等水煮沸，然后在煮沸后再做其他事情。

如果你能理解为什么第一种方式泡茶更快，那你就可以理解我们如何让 `CheckWebsites` 变得更快。与其等待网站响应之后再发送下一个网站的请求，不如告诉计算机在等待时就发起下一个请求。

通常在 Go 中，当调用函数 `doSomething()` 时，我们等待它返回（即使它没有值返回，我们仍然等待它完成）。我们说这个操作是 *阻塞* 的 —— 它让我们等待它完成。Go 中不会阻塞的操作将在称为 *goroutine* 的单独 *进程* 中运行。将程序想象成从上到下读 Go 的 代码，当函数被调用执行读取操作时，进入每个函数「内部」。当一个单独的进程开始时，就像开启另一个 reader（阅读程序）在函数内部执行读取操作，原来的 reader 继续向下读取 Go 代码。

要告诉 Go 开始一个新的 goroutine，我们把一个函数调用变成 `go` 声明，通过把关键字 `go` 放在它前面：`go doSomething()`。

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        go func() {
            results[url] = wc(url)
        }()
    }

    return results
}
```

因为开启 goroutine 的唯一方法就是将 `go` 放在函数调用前面，所以当我们想要启动 goroutine 时，我们经常使用 *匿名函数（anonymous functions）*。一个匿名函数文字看起来和正常函数声明一样，但没有名字（意料之中）。你可以在 上面的 `for` 循环体中看到一个。

匿名函数有许多有用的特性，其中两个上面正在使用。首先，它们可以在声明的同时执行 —— 这就是匿名函数末尾的 `()` 实现的。其次，它们维护对其所定义的词汇作用域的访问权 —— 在声明匿名函数时所有可用的变量也可在函数体内使用。

上面匿名函数的主体和之前循环体中的完全一样。唯一的区别是循环的每次迭代都会启动一个新的 goroutine，与当前进程（`WebsiteChecker` 函数）同时发生，每个循环都会将结果添加到 `results` map 中。

但是当我们执行 `go test`：

```
-------- FAIL: TestCheckWebsites (0.00s)
        CheckWebsites_test.go:31: Wanted map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], got map[]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s
```

## 快速进入平行宇宙......

你可能不会得到这个结果。你可能会得到一个 panic 信息，这个稍后再谈。如果你得到的是那些结果，不要担心，只要继续运行测试，直到你得到上述结果。或假装你得到了，这取决于你。欢迎来到并发编程的世界：如果处理不正确，很难预测会发生什么。别担心 —— 这就是我们编写测试的原因，当处理并发时，测试帮助我们预测可能发生的情况。

## ... 重新回到这些问题。

让我们困惑的是，原来的测试 `WebsiteChecker` 现在返回空的 map。哪里出问题了？

我们 `for` 循环开始的 `goroutines` 没有足够的时间将结果添加结果到 `results` map 中；`WebsiteChecker` 函数对于它们来说太快了，以至于它返回时仍为空的 map。

为了解决这个问题，我们可以等待所有的 goroutine 完成他们的工作，然后返回。两秒钟应该能完成了，对吧？

```go
package concurrency

import "time"

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        go func() {
            results[url] = wc(url)
        }()
    }

    time.Sleep(2 * time.Second)

    return results
}
```

现在当我们运行测试时获得的结果（如果没有得到 —— 参考上面的做法）：

```
-------- FAIL: TestCheckWebsites (0.00s)
        CheckWebsites_test.go:31: Wanted map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], got map[waat://furhurterwe.geds:false]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s
```
这不是很好 - 为什么只有一个结果？我们可以尝试通过增加等待的时间来解决这个问题 —— 如果你愿意，可以试试。但没什么作用。这里的问题是变量 `url` 被重复用于 `for` 循环的每次迭代 —— 每次都会从 `urls` 获取新值。但是我们的每个 goroutine 都是 `url` 变量的引用 —— 它们没有自己的独立副本。所以他们 *都* 会写入在迭代结束时的 `url` —— 最后一个 url。这就是为什么我们得到的结果是最后一个 url。

解决这个问题：

```go
import (
    "time"
)

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)

    for _, url := range urls {
        go func(u string) {
            results[u] = wc(u)
        }(url)
    }

    time.Sleep(2 * time.Second)

    return results
}
```
通过给每个匿名函数一个参数 url(`u`)，然后用 `url` 作为参数调用匿名函数，我们确保 `u` 的值固定为循环迭代的 `url` 值，重新启动 `goroutine`。`u` 是 `url` 值的副本，因此无法更改。

现在，如果你幸运的话，你会得到：

```
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        2.012s
```

但是，如果你不走运（如果你运行基准测试，这很可能会发生，因为你将发起多次的尝试）。

```
fatal error: concurrent map writes

goroutine 8 [running]:
runtime.throw(0x12c5895, 0x15)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/panic.go:605 +0x95 fp=0xc420037700 sp=0xc4200376e0 pc=0x102d395
runtime.mapassign_faststr(0x1271d80, 0xc42007acf0, 0x12c6634, 0x17, 0x0)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:783 +0x4f5 fp=0xc420037780 sp=0xc420037700 pc=0x100eb65
github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1(0xc42007acf0, 0x12d3938, 0x12c6634, 0x17)
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x71 fp=0xc4200377c0 sp=0xc420037780 pc=0x12308f1
runtime.goexit()
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/asm_amd64.s:2337 +0x1 fp=0xc4200377c8 sp=0xc4200377c0 pc=0x105cf01
created by github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xa1

        ... many more scary lines of text ...
```

这看上去冗长、可怕，我们需要深呼吸并阅读错误：`fatal error: concurrent map writes`。有时候，当我们运行我们的测试时，两个 goroutines 完全同时写入 `results` map。Go 的 Maps 不喜欢多个事物试图一次性写入，所以就导致了 `fatal error`。

这是一种 *race condition（竞争条件）*，当软件的输出取决于事件发生的时间和顺序时，因为我们无法控制，bug 就会出现。因为我们无法准确控制每个 goroutine 写入结果 map 的时间，两个 goroutines 同一时间写入时程序将非常脆弱。

Go 可以帮助我们通过其内置的 [race detector](https://blog.golang.org/race-detector) 来发现竞争条件。要启用此功能，请使用 `race` 标志运行测试：`go test -race`。

你应该得到一些如下所示的输出：

```
==================
WARNING: DATA RACE
Write at 0x00c420084d20 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Previous write at 0x00c420084d20 by goroutine 7:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Goroutine 8 (running) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c

Goroutine 7 (finished) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c
==================
```

细节还是难以阅读 - 但 `WARNING: DATA RACE` 相当明确。阅读错误的内容，我们可以看到两个不同的 goroutines 在 map 上执行写入操作：

`Write at 0x00c420084d20 by goroutine 8:`

正在写入相同的内存块

`Previous write at 0x00c420084d20 by goroutine 7:`

最重要的是，我们可以看到发生写入的代码行：

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12`

和 goroutines 7 和 8 开始的代码行号：

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11`

你需要知道的所有内容都会打印到你的终端上 - 你只需耐心阅读就可以了。

## Channels

我们可以通过使用 *channels* 协调我们的 goroutines 来解决这个数据竞争。channels 是一个 Go 数据结构，可以同时接收和发送值。这些操作以及细节允许不同进程之间的通信。

在这种情况下，我们想要考虑父进程和每个 goroutine 之间的通信，goroutine 使用 url 来执行 `WebsiteChecker` 函数。

```go
package concurrency

type WebsiteChecker func(string) bool
type result struct {
    string
    bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
    results := make(map[string]bool)
    resultChannel := make(chan result)

    for _, url := range urls {
        go func(u string) {
            resultChannel <- result{u, wc(u)}
        }(url)
    }

    for i := 0; i < len(urls); i++ {
        result := <-resultChannel
        results[result.string] = result.bool
    }

    return results
}
```

除了 `results` map 之外，我们现在还有一个 `resultChannel` 的变量，同样使用 `make` 方法创建。`chan result` 是 channel 类型的 —— `result` 的 channel。新类型的 `result` 是将 `WebsiteChecker` 的返回值与正在检查的 url 相关联 —— 它是一个 `string` 和 `bool` 的结构。因为我们不需要任何一个要命名的值，它们中的每一个在结构中都是匿名的；这在很难知道用什么命名值的时候可能很有用。

现在，当我们迭代 urls 时，不是直接写入 `map`，而是使用 *send statement* 将每个调用 `wc` 的 `result` 结构体发送到 `resultChannel`。这使用 `<-` 操作符，channel 放在左边，值放在右边：

```
// send statement
resultChannel <- result{u, wc(u)}
```

下一个 `for` 循环为每个 url 迭代一次。 我们在内部使用 *receive expression*，它将从通道接收到的值分配给变量。这也使用 `<-` 操作符，但现在两个操作数颠倒过来：现在 channel 在右边，我们指定的变量在左边：

```
// receive expression
result := <-resultChannel
```

然后我们使用接收到的 `result` 更新 map。

通过将结果发送到通道，我们可以控制每次写入 `results` map 的时间，确保每次写入一个结果。虽然 `wc` 的每个调用都发送给结果通道，但是它们在其自己的进程内并行发生，因为我们将结果通道中的值与接收表达式一起逐个处理一个结果。

我们已经将想要加快速度的那部分代码并行化，同时确保不能并发的部分仍然是线性处理。我们使用 channel 在多个进程间通信。

当我们运行基准时：

```
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v2
BenchmarkCheckWebsites-8             100          23406615 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        2.377s
```

23406615 纳秒 —— 0.023 秒，速度大约是最初函数的一百倍，这是非常成功的。

## 总结

这个比在 TDD 上的寻常练习轻松一些。某种程度说，我们已经参与了 `CheckWebsites` 函数的一个长期重构；输入和输出从未改变，它只是变得更快了。但是我们所做的测试以及我们编写的基准测试允许我们重构 `CheckWebsites`，让我们有信心保证软件仍然可以工作，同时也证明它确实变得更快了。

在使它更快的过程中，我们明白了

- *goroutines* 是 Go 的基本并发单元，它让我们可以同时检查多个网站。
- *anonymous functions（匿名函数）*，我们用它来启动每个检查网站的并发进程。
- *channels*，用来组织和控制不同进程之间的交流，使我们能够避免 *race condition（竞争条件）* 的问题。
- *the race detector（竞争探测器）* 帮助我们调试并发代码的问题。

### 使程序加快

一种构建软件的敏捷方法，常常被错误地归属于 Kent Beck，即：

> [让它运作，使它正确，使它快速](http://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast)

「运作」是通过测试，「正确」是重构代码，而「快速」是优化代码以使其快速运行。一旦我们使程序可以正确运行，我们能做的就只有使它快速。很幸运，我们得到的代码已经被证明是可以运作的，并且不需要重构。在另外两个步骤执行之前，我们绝不应该试图「使它快速」，因为

> [过早的优化是万恶之源](http://wiki.c2.com/?PrematureOptimization)
> —— Donald Knuth

---

作者：[David Wickes](https://dev.to/gypsydave5)
译者：[Donng](https://github.com/Donng)
校对：[polaris1119](https://github.com/polaris1119)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
