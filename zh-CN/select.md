# Select

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/select)**

你被要求编写一个叫做 `WebsiteRacer` 的函数，用来对比请求两个 URL 来「比赛」，并返回先响应的 URL。如果两个 URL 在 10 秒内都未返回结果，那么应该返回一个 `error`。

实现这个功能我们需要用到

- `net/http` 用来调用 HTTP 请求
- `net/http/httptest` 用来测试这些请求
- Go 程（goroutines）
- `select` 用来同步进程

## 先写测试

我们从最幼稚的做法开头把事情开展起来。

```go
func TestRacer(t *testing.T) {
    slowURL := "http://www.facebook.com"
    fastURL := "http://www.quii.co.uk"

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}
```

我们知道这样不完美并且有问题，但这样可以把事情开展起来。重要的是，不要徘徊在第一次就想把事情做到完美。

## 尝试运行测试

`./racer_test.go:14:9: undefined: Racer`

## 为测试的运行编写最少量的代码，并检查失败测试的输出

```go
func Racer(a, b string) (winner string) {
    return
}
```

`racer_test.go:25: got '', want 'http://www.quii.co.uk'`

## 编写足够的代码使程序通过

```go
func Racer(a, b string) (winner string) {
    startA := time.Now()
    http.Get(a)
    aDuration := time.Since(startA)

    startB := time.Now()
    http.Get(b)
    bDuration := time.Since(startB)

    if aDuration < bDuration {
        return a
    }

    return b
}
```

对每个 URL：

1. 我们用 `time.Now()` 来记录请求 `URL` 前的时间。
1. 然后用 [`http.Get`](https://golang.org/pkg/net/http/#Client.Get) 来请求 `URL` 的内容。这个函数返回一个 [`http.Response`](https://golang.org/pkg/net/http/#Response) 和一个 `error`，但目前我们不关心它们的值。
1. `time.Since` 获取开始时间并返回一个 `time.Duration` 时间差。

我们完成这些后就可以通过对比请求耗时来找出最快的了。

### 问题

这可能会让你的测试通过，也可能不会。问题是我们通过访问真实网站来测试我们的逻辑。

使用 HTTP 测试代码非常常见，Go 标准库有这类工具可以帮助测试。

在模拟和依赖注入章节中，我们讲到了理想情况下如何不依赖外部服务来进行测试，因为它们可能

- 速度慢
- 不可靠
- 无法进行边界条件测试

在标准库中有一个 [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/) 包，它可以让你轻易建立一个 HTTP 模拟服务器（mock HTTP server）。

我们改为使用模拟测试，这样我们就可以控制可靠的服务器来测试了。

```go
func TestRacer(t *testing.T) {

    slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(20 * time.Millisecond)
        w.WriteHeader(http.StatusOK)
    }))

    fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    slowURL := slowServer.URL
    fastURL := fastServer.URL

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }

    slowServer.Close()
    fastServer.Close()
}
```

语法看着有点儿复杂，没关系，慢慢来。

`httptest.NewServer` 接受一个我们传入的 _匿名函数_ `http.HandlerFunc`。

`http.HandlerFunc` 是一个看起来类似这样的类型：`type HandlerFunc func(ResponseWriter, *Request)`。

这些只是说它是一个需要接受一个 `ResponseWriter` 和 `Request` 参数的函数，这对于 HTTP 服务器来说并不奇怪。

结果呢，这里并没有什么彩蛋，**这也是如何在 Go 语言写一个 _真实的_ HTTP 服务器的方法**。唯一的区别就是我们把它封装成一个易于测试的 `httptest.NewServer`，它会找一个可监听的端口，然后测试完你就可以关闭它了。

我们让两个服务器中慢的那一个短暂地 `time.Sleep` 一段时间，当我们请求时让它比另一个慢一些。然后两个服务器都会通过 `w.WriteHeader(http.StatusOK)` 返回一个 `OK` 给调用者。

如果你重新运行测试，它现在肯定会通过并且会更快完成。调整 sleep 时间故意破坏测试。

## 重构

我们在主程序代码和测试代码里都有一些重复。

```go
func Racer(a, b string) (winner string) {
    aDuration := measureResponseTime(a)
    bDuration := measureResponseTime(b)

    if aDuration < bDuration {
        return a
    }

    return b
}

func measureResponseTime(url string) time.Duration {
    start := time.Now()
    http.Get(url)
    return time.Since(start)
}
```

这样简化代码后可以让 `Racer` 函数更加易读。

```go
func TestRacer(t *testing.T) {

    slowServer := makeDelayedServer(20 * time.Millisecond)
    fastServer := makeDelayedServer(0 * time.Millisecond)

    defer slowServer.Close()
    defer fastServer.Close()

    slowURL := slowServer.URL
    fastURL := fastServer.URL

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(delay)
        w.WriteHeader(http.StatusOK)
    }))
}
```

我们通过一个名为 `makeDelayedServer` 的函数重构了模拟服务器，以将一些不感兴趣的代码移出测试并减少了重复代码。

### `defer`

在某个函数调用前加上 `defer` 前缀会在 _包含它的函数结束时_ 调用它。

有时你需要清理资源，例如关闭一个文件，在我们的案例中是关闭一个服务器，使它不再监听一个端口。

你想让它在函数结束时执行（关闭服务器），但要把它放在你创建服务器语句附近，以便函数内后面的代码仍可以使用这个服务器。

我们的重构是一次改进，并且目前是涵盖 Go 语言特性提供的合理解决方案，但我们可以让它更简单。

### 进程同步

- Go 在并发方面很在行，为什么我们要一个接一个地测试哪个网站更快呢？我们应该能够同时测试两个。
- 我们并不关心请求的 _准确响应时间_，我们只是需要知道哪个更快返回而已。

想实现这个，我们要介绍一个叫 `select` 的新构造（construct），它可以帮我们轻易清晰地实现进程同步。

```go
func Racer(a, b string) (winner string) {
    select {
    case <-ping(a):
        return a
    case <-ping(b):
        return b
    }
}

func ping(url string) chan bool {
    ch := make(chan bool)
    go func() {
        http.Get(url)
        ch <- true
    }()
    return ch
}
```

#### `ping`

我们定义了一个可以创建 `chan bool` 类型并返回它的 `ping` 函数。

在这个案例中，我们并不 _关心_ channel 中发送的类型， _我们只是想发送一个信号_ 来说明已经发送完了，所以返回 bool 就可以了。

同样在这个函数中，当我们完成 `http.Get(url)` 时启动了一个用来给 channel 发送信号的 Go 程（goroutine）。

#### `select`

如果你记得并发那一章的内容，你可以通过 `myVar := <-ch` 来等待值发送给 channel。这是一个 _阻塞_ 的调用，因为你需要等待值返回。

`select` 则允许你同时在 _多个_ channel 等待。第一个发送值的 channel「胜出」，`case` 中的代码会被执行。

我们在 `select` 中使用 `ping` 为两个 `URL` 设置两个 channel。无论哪个先写入其 channel 都会使 `select` 里的代码先被执行，这会导致那个 `URL` 先被返回（胜出）。

做了这些修改后，我们的代码背后的意图就很明确了，实现起来也更简单。

### 超时

最后的需求是当 `Racer` 耗时超过 10 秒时返回一个 error。

## 先写测试

```go
t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
    serverA := makeDelayedServer(11 * time.Second)
    serverB := makeDelayedServer(12 * time.Second)

    defer serverA.Close()
    defer serverB.Close()

    _, err := Racer(serverA.URL, serverB.URL)

    if err == nil {
        t.Error("expected an error but didn't get one")
    }
})
```

为了练习这个场景，现在我们要使模拟服务器超过 10 秒后返回两个值，胜出的 URL（这个测试中我们用 `_` 忽略掉了）和一个 `error`。

## 尝试运行测试

`./racer_test.go:37:10: assignment mismatch: 2 variables but 1 values`

## 为测试的运行编写最少量的代码，并检查失败测试的输出

```go
func Racer(a, b string) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    }
}
```

修改 `Racer` 的函数签名来返回胜出者和一个 `error`。返回 `nil` 仅用于模拟顺利的场景（happy cases）。

编译器会报怨你的 _第一个测试_ 只期望一个值，所以把这行改为 `got, _ := Racer(slowURL, fastURL)`，要知道顺利的场景中我们不应得到一个 `error`。

现在运行测试会在超过 11 秒后失败。

```
--- FAIL: TestRacer (12.00s)
    --- FAIL: TestRacer/returns_an_error_if_a_server_doesn't_respond_within_10s (12.00s)
        racer_test.go:40: expected an error but didn't get one
```

## 编写足够的代码使程序通过

```go
func Racer(a, b string) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(10 * time.Second):
        return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
    }
}
```

使用 `select` 时，`time.After` 是一个很好用的函数。当你监听的 channel 永远不会返回一个值时你可以潜在地编写永远阻塞的代码，尽管在我们的案例中它没有发生。`time.After` 会在你定义的时间过后发送一个信号给 channel 并返回一个 `chan` 类型（就像 `ping` 那样）。

对我们来说这完美了；如果 `a` 或 `b` 谁胜出就返回谁，但如果测试达到 10 秒，那么 `time.After` 会发送一个信号并返回一个 `error`。

### 慢速测试

现在的问题是这个测试要耗时 10 秒以上。对这么简单的逻辑来说可不好。

我们可以做的就是让超时时间（timeout）可配置，这样测试就可以设置一个非常短的时间，并且代码在真实环境中可以被设置成 10 秒。

```go
func Racer(a, b string, timeout time.Duration) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
    }
}
```

现在代码不能编译了，因为我们没提供超时时间。

在急于将这个默认值添加到测试前，先让我们 _聆听他们_。

- 在顺利的情况「happy test」下我们是否关心超时时间？
- 需求对超时时间很明确

鉴于以上信息，我们再做一次小的重构来让我们的测试和代码的用户合意。

```go
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
    return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
    select {
    case <-ping(a):
        return a, nil
    case <-ping(b):
        return b, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
    }
}
```

我们的用户和第一个测试可以使用 `Racer`（使用 `ConfigurableRacer`），不顺利的场景测试可以使用 `ConfigurableRacer`。

```go
func TestRacer(t *testing.T) {

    t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
        slowServer := makeDelayedServer(20 * time.Millisecond)
        fastServer := makeDelayedServer(0 * time.Millisecond)

        defer slowServer.Close()
        defer fastServer.Close()

        slowURL := slowServer.URL
        fastURL := fastServer.URL

        want := fastURL
        got, err := Racer(slowURL, fastURL)

        if err != nil {
            t.Fatalf("did not expect an error but got one %v", err)
        }

        if got != want {
            t.Errorf("got '%s', want '%s'", got, want)
        }
    })

    t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
        server := makeDelayedServer(25 * time.Millisecond)

        defer server.Close()

        _, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

        if err == nil {
            t.Error("expected an error but didn't get one")
        }
    })
}
```

我在第一个测试最后加了一个检查来验证我们没得到一个 `error`。

## 总结

### `select`

- 可帮助你同时在多个 channel 上等待。
- 有时你想在你的某个「案例」中使用 `time.After` 来防止你的系统被永久阻塞。

### `httptest`

- 一种方便地创建测试服务器的方法，这样你就可以进行可靠且可控的测试。
- 使用和 `net/http` 相同的接口作为「真实的」服务器会和真实环境保持一致，并且只需更少的学习。

---

作者：[Chris James](https://dev.to/quii)
译者：[pityonline](https://github.com/pityonline)
校对：[Donng](https://github.com/Donng)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
