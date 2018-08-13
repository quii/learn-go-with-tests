# HTTP 服务器

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/http-server)**

你被要求创建一个 Web 服务器，用户可以在其中跟踪玩家赢了多少场游戏。

- `GET /players/{name}` 应该返回一个表示获胜总数的数字
- `POST /players/{name}` 应该为玩家赢得游戏记录一次得分，并随着每次 `POST` 递增

我们将遵循 TDD 方法，尽可能快地让程序先可用，然后进行小步迭代改进，直到我们找到解决方案。通过采取这种方法我们

- 在任何给定时间保持问题都是小问题
- 不要陷入陷阱（rabbit holes）
- 如果我们卡住或迷失了方向，回退不会前功尽弃。

## 持续迭代，重构

在本书中，我们强调了编写测试并观察失败（红色），编写 _最少量_ 代码跑通测试（绿色）然后重构的 TDD 过程。

就 TDD 的安全性而言，编写最少量代码的这一规则非常重要。你应该尽快摆脱测试失败的状态（红色）。

Kent Beck 这样描述它：

> 快速跑通测试，为满足必要条件暂时犯错亦可。

你可以先写一些存在已知问题的代码，因为后面你会基于 TDD 安全地进行重构。

### 如果不这样做？

测试未通过前你写得越多，也就引入越来越多测试不能覆盖的问题。

这样做是为了小步快速迭代，测试驱使你不会掉入陷阱。

### 先有鸡还是先有蛋

我们如何逐步建立这个？我们不能在没有数据的前提下 `GET` 一个玩家，而且似乎很难知道 `POST` 在没有 `GET` 的情况下是否工作。

这就是 _模拟_ 测试的亮点。

- `GET` 需要一个类似 `PlayerStore` 的东西来获得玩家的分数。这应该是一个接口，所以测试时我们可以创建一个简单的存根来测试代码而无需实现任何真实的存储机制。
- 对于 `POST`，我们可以 _监听_ `PlayerStore` 的调用以确保它能正确存储玩家。我们的存储实现不会与检索相关联。
- 为了尽快让代码可运行，我们可以先在内存中写一个非常简单的实现，然后我们可以实现一个任何喜欢的存储机制。

## 先写测试

我们可以先写一个硬编码的值让测试通过来启动我们的工作。Kent Beck 把这称为“造数据”。一旦有了一个可以通过的测试用例，我们可以接着写更多测试来帮我们删除之前的硬编码代码。

这些微小的步骤对启动工作来说很重要，它令我们的项目先有一个正确的结构，而不用过于担心程序逻辑。

一般来说，你可以在 Go 中调用 [ListenAndServe](https://golang.org/pkg/net/http/#ListenAndServe) 来创建一个 web 服务器

```go
func ListenAndServe(addr string, handler Handler) error
```

这可以启动一个 web 服务器监听在一个端口上，为每个请求创建一个 Go 例程，并用一个 handler 处理这些请求。

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

这个函数期望有两个参数输入，第一个是 _响应_ 请求的，第二个是发送给服务器的 HTTP 请求。

我们为 `PlayerServer` 写一个测试函数，让它接受上面提到的两个参数。发送的请求将得到一个期望为 20 的玩家得分。

```go
t.Run("returns Pepper's score", func(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
    response := httptest.NewRecorder()

    PlayerServer(response, request)

    got := response.Body.String()
    want := "20"

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
})
```

为了测试服务器，我们需要通过 `Request` 来发送请求，并期望监听到 handler 向 `ResponseWriter` 写入了什么。

- 我们用 `http.NewRequest` 来创建一个请求。第一个参数是请求方法，第二个是请求路径。`nil` 是请求实体，不过在这个场景中不用发送请求实体。
- `net/http/httptest` 自带一个名为 `ResponseRecorder` 的监听器，所以我们可以用这个。它有很多有用的方法可以检查应答被写入了什么。

## 尝试运行测试

`./server_test.go:13:2: undefined: PlayerServer`

## 编写最少量的代码让测试运行起来，然后检查错误输出

编译器给出的信息非常有用，照它说的办。

声明 `PlayerServer` ：

```go
func PlayerServer() {}
```

再运行一次

```
./server_test.go:13:14: too many arguments in call to PlayerServer
    have (*httptest.ResponseRecorder, *http.Request)
    want ()
```

给函数添加参数：

```go
import "net/http"

func PlayerServer(w http.ResponseWriter, r *http.Request) {

}
```

代码可以编译了，只是测试还是失败的。

```
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- FAIL: TestGETPlayers/returns_Pepper's_score (0.00s)
        server_test.go:20: got '', want '20'
```

## 编写足够的代码让它通过

在依赖注入章节中，我们通过 `Greet` 函数接触到了 HTTP 服务器。我们知道了 net/http 的 `ResponseWriter` 也实现了 io `Writer`，所以我们可以用 `fmt.Fprint` 发送字符串来作为 HTTP 应答。

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "20")
}
```

现在测试应该通过了。

## 完成框架（scaffolding）

我们要把它联系起来，这很重要。因为

- 我们要写真正能用的代码，而不是为了测试而写测试，代码能用才是王道。
- 当我们重构时，可能会修改程序的结构。我们希望确保这也作为增量方法的一部分反映在程序中。

创建一个新文件，写入以下代码。

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    handler := http.HandlerFunc(PlayerServer)
    if err := http.ListenAndServe(":5000", handler); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

目前我们所有程序代码都在一个文件里，然而这不是那种把代码拆分为多个文件的大型项目的最佳实践。

通过 `go build` 把目录中所有 `.go` 文件编译成一个可运行的程序，然后你可以用 `./myprogram` 来运行它。

### `http.HandlerFunc`

之前我们探讨过 `Handler` 接口是为创建服务器而需要实现的。一般来说，我们通过创建 `struct` 来实现接口。然而，struct 的用途是用于存储数据，但是目前没有状态可存储，因此创建一个 struct 感觉不太对。

[HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc) 可以让我们避免这样。

> HandlerFunc 类型是一个允许将普通函数用作 HTTP handler 的适配器。如果 f 是具有适当签名的函数，则 HandlerFunc(f) 是一个调用 f 的 Handler。

```go
type HandlerFunc func(ResponseWriter, *Request)
```

所以我们用它来封装 `PlayerServer` 函数，使它现在符合 `Handler`。

### `http.ListenAndServe(":5000"...)`

`ListenAndServe` 会在 `Handler` 上监听一个端口。如果端口已被占用，它会返回一个 `error`，所以我们在一个 `if` 语句中捕获出错的场景并记录下来。

我们现在要做的是编写另一个测试来迫使我们尝试摆脱现有的硬编码。

## 先写测试

我们添加另一个子测试来尝试为不同的玩家获取得分，来破坏之前硬编码的实现。

```go
t.Run("returns Floyd's score", func(t *testing.T) {
    request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
    response := httptest.NewRecorder()

    PlayerServer(response, request)

    got := response.Body.String()
    want := "10"

    if got != want {
        t.Errorf("got '%s', want '%s'", got, want)
    }
})
```

你或许在想：

> 当然，我们需要一种存储机制来控制不同玩家的得分。在这个测试中那些值看起来很武断，这有点儿怪。

注意，我们只是尽量合理地小步前进，所以现在只改善硬编码的问题。

## 尝试运行测试

```
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- PASS: TestGETPlayers/returns_Pepper's_score (0.00s)
=== RUN   TestGETPlayers/returns_Floyd's_score
    --- FAIL: TestGETPlayers/returns_Floyd's_score (0.00s)
        server_test.go:34: got '20', want '10'
```

## 编写足够的代码让它通过

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    if player == "Pepper" {
        fmt.Fprint(w, "20")
        return
    }

    if player == "Floyd" {
        fmt.Fprint(w, "10")
        return
    }
}
```

这样迫使我们必须决定请求的 URL 怎么写。因此，在我们的脑海里可能一直担心玩家数据存储和接口，下一个逻辑步骤实际上应该是路由（routing）了。

如果我们从存储机制开始写，那么与此相比，我们必须做的工作量将非常大。 **这是朝着我们的最终目标迈出的一小步，并且是由测试驱动的**。

我们现在正在抵制使用任何路由库的诱惑，这是让测试通过的最小步骤。

`r.URL.Path` 返回请求的路径，然后我们用切片语法得到 `/players/` 最后的斜杠后的路径。这不太靠谱，但现在起码可行。

## 重构

我们可以通过将分数检索分离为函数来简化 `PlayerServer`

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(name string) string {
    if name == "Pepper" {
        return "20"
    }

    if name == "Floyd" {
        return "10"
    }

    return ""
}
```

我们可以创建一些辅助函数来避免测试中的重复代码

```go
func TestGETPlayers(t *testing.T) {
    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        PlayerServer(response, request)

        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        PlayerServer(response, request)

        assertResponseBody(t, response.Body.String(), "10")
    })
}

func newGetScoreRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
    return req
}

func assertResponseBody(t *testing.T, got, want string) {
    t.Helper()
    if got != want {
        t.Errorf("response body is wrong, got '%s' want '%s'", got, want)
    }
}
```

然而事情还没完。好像服务器并不知道那些得分。

不过，重构还是令代码更清晰可读了。

我们把得分计算从 handler 移到函数 `GetPlayerScore` 中，这就是使用接口重构的正确方法。

让我们把重构的函数改成接口。

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
}
```

为了让 `PlayerServer` 能够使用 `PlayerStore`，它需要一个引用。现在是改变架构的时候了，将 `PlayerServer` 改成一个 `struct`。

```go
type PlayerServer struct {
    store PlayerStore
}
```

最后，我们通过给这个 struct 添加一个方法来实现 `Handler` 接口，并把它放到已有的 handler 中。

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]
    fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

最后一个修改是现在调用 `store.GetPlayerStore` 来获得得分，而不是我们定义的本地函数（现在可以删除它了）。

下面是 server 的完整代码

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
}

type PlayerServer struct {
    store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]
    fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

### 解决问题

我们知道这些修改会让测试和程序不能编译，别急，先让编译器编译一下。

`./main.go:9:58: type PlayerServer is not an expression`

我们需要改一下测试，而不是创建一个新的 `PlayerServer` 实例，然后调用它的 `ServeHTTP` 方法。

```go
func TestGETPlayers(t *testing.T) {
    server := &PlayerServer{}

    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "10")
    })
}
```

请注意，我们目前仍不关心存储问题，我们只希望编译尽快通过。

你应该养成优先让编译通过再让测试通过的编程习惯。

通过添加更多功能（如存根存储）没使代码成功编译，我们正在面对更多潜在的编译问题。

现在 `main.go` 还是因同样的原因编译失败。

```go
func main() {
    server := &PlayerServer{}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

最后编译终于通过了，但测试还没通过

```
=== RUN   TestGETPlayers/returns_the_Pepper's_score
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
    panic: runtime error: invalid memory address or nil pointer dereference
```

这是因为在测试中我还没传入 `PlayerStore`，我们需要创建一个存根。

```go
type StubPlayerStore struct {
    scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}
```

使用 `map` 创建键/值存储是一种比较简便快捷的方式。现在让我们在测试中创建其中一个 `store` 并将其传给 `PlayerServer`。

```go
func TestGETPlayers(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{
            "Pepper": 20,
            "Floyd":  10,
        },
    }
    server := &PlayerServer{&store}

    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertResponseBody(t, response.Body.String(), "10")
    })
}
```

测试现在通过并且看起来更好一些了。由于 `store` 的引入，代码意图现在更清晰了。我们告诉 reader，在 `PlayerStore` 中有数据了，当你将它用在 `PlayerServer` 时，你应该得到正确的应答。

### 运行程序

现在测试通过了，要完成这个重构，我们要做的最后一件事是检查应用程序是否正常工作。该程序应该可以启动，但如果你尝试访问 `http://localhost:5000/players/Pepper`，你会得到一个异常的应答。

原因是我们没传入 `PlayerStore`。

我们要实现一个，但现在有点儿困难，因为我们还没存储有用的数据，所以我们先通过硬编码实现。

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return 123
}

func main() {
    server := &PlayerServer{&InMemoryPlayerStore{}}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

如果你再次运行 `go build` 并访问同一个 URL 你应该得到一个 `"123"` 的应答。尽管这样不太对，但在我们实现数据存储前已经是最好不过的了。

关于下一步该做什么，我们有几个选择

- 处理玩家不存在的场景
- 处理 `POST /players/{name}` 的场景
- 主应用程序能运行但实际上还不算真正能用。我们不得不手动测试才能看到问题。

虽然 `POST` 场景让我们更接近“测试通过”，但我觉得首先解决玩家不存在的情景会更容易，因为我们已经处于这种情况。我们稍后会讨论其余的事情。

## 先写测试

添加一个玩家不存在的测试用例

```go
t.Run("returns 404 on missing players", func(t *testing.T) {
    request := newGetScoreRequest("Apollo")
    response := httptest.NewRecorder()

    server.ServeHTTP(response, request)

    got := response.Code
    want := http.StatusNotFound

    if got != want {
        t.Errorf("got status %d want %d", got, want)
    }
})
```

## 尝试运行测试

```
=== RUN   TestGETPlayers/returns_404_on_missing_players
    --- FAIL: TestGETPlayers/returns_404_on_missing_players (0.00s)
        server_test.go:56: got status 200 want 404
```

## 编写最少量的代码让测试运行起来，然后检查错误输出

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    w.WriteHeader(http.StatusNotFound)

    fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

当 TDD 的拥护者说“确保你只需编写最少量的代码以使其通过”时，有时我也不敢苟同，因为这可能有些墨守陈规了。

但这个场景正好说明了这个例子。我已经做了最少量（明知不对）的修改，我让 **所有响应** 应答一个 `StatusNotFound` 但所有的测试却都通过了！

**小步修改让测试通过，可以突出测试中的问题**。在这个例子中，当玩家的确存在于 store 中时，我们并没有断言我们应该获得 `StatusOK`。

更新另外两个测试以断言返回状态并修复代码。

这是新的测试代码：

```go
func TestGETPlayers(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{
            "Pepper": 20,
            "Floyd":  10,
        },
    }
    server := &PlayerServer{&store}

    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusOK)
        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusOK)
        assertResponseBody(t, response.Body.String(), "10")
    })

    t.Run("returns 404 on missing players", func(t *testing.T) {
        request := newGetScoreRequest("Apollo")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusNotFound)
    })
}

func assertStatus(t *testing.T, got, want int) {
    t.Helper()
    if got != want {
        t.Errorf("did not get correct status, got %d, want %d", got, want)
    }
}

func newGetScoreRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
    return req
}

func assertResponseBody(t *testing.T, got, want string) {
    t.Helper()
    if got != want {
        t.Errorf("response body is wrong, got '%s' want '%s'", got, want)
    }
}
```

现在我们在测试中检查所有返回状态了，所以我创建了一个叫 `assertStatus` 的辅助函数来提高编码效率。

现在前两个测试失败了，因为返回状态是 404 而不是 200，所以如果得分为 0，我们可以修复 `PlayerServer` 只返回 “not found” 的问题。

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}
```

### 存储得分

现在我们可以从 store 中查询得分了，能够存储新的得分就很有意义了。

## 先写测试

```go
func TestStoreWins(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{},
    }
    server := &PlayerServer{&store}

    t.Run("it returns accepted on POST", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusAccepted)
    })
}
```

首先，我们检查使用 POST 访问指定路径时是否返回了正确的状态码。这迫使我们要实现可以接受不同类型请求的功能，以不同的方式处理 `GET /players/{name}`。一旦这个通过，我们就可以开始测试 handler 与 store 的交互。

## 尝试运行测试

```
=== RUN   TestStoreWins/it_returns_accepted_on_POST
    --- FAIL: TestStoreWins/it_returns_accepted_on_POST (0.00s)
        server_test.go:70: did not get correct status, got 404, want 202
```

## 编写足够的代码让它通过

注意我们故意写错，所以用一个 `if` 语句来测试请求方法就可以了。

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodPost {
        w.WriteHeader(http.StatusAccepted)
        return
    }

    player := r.URL.Path[len("/players/"):]

    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}
```

## 重构

现在 handler 看起来有点儿乱了。让我们用下面的代码重写，以便更容易理解并将不同的功能拆分到新函数中。

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    switch r.Method {
    case http.MethodPost:
        p.processWin(w)
    case http.MethodGet:
        p.showScore(w, r)
    }

}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter) {
    w.WriteHeader(http.StatusAccepted)
}
```

这令 `ServeHTTP` 的路由更加清晰，这意味着我们下一次存储迭代只能在 `processWin` 中。

接下来要检查当我们执行 `POST /players/{name}` 时 `PlayerStore` 被告知要做一次获胜记录。

## 先写测试

我们可以通过使用新的 `RecordWin` 方法扩展 `StubPlayerStore` 然后监视它的调用来实现这一点。

```go
type StubPlayerStore struct {
    scores   map[string]int
    winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}

func (s *StubPlayerStore) RecordWin(name string) {
    s.winCalls = append(s.winCalls, name)
}
```

现在在测试中扩展以检查启动的调用次数

```go
func TestStoreWins(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{},
    }
    server := &PlayerServer{&store}

    t.Run("it records wins when POST", func(t *testing.T) {
        request := newPostWinRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusAccepted)

        if len(store.winCalls) != 1 {
            t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
        }
    })
}

func newPostWinRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
    return req
}
```

## 尝试运行测试

```
./server_test.go:26:20: too few values in struct initializer
./server_test.go:65:20: too few values in struct initializer
```

## 编写最少量的代码让测试运行起来，然后检查错误输出

我们需要更新创建 `StubPlayerStore` 的代码，因为我们添加了一个新字段

```go
store := StubPlayerStore{
    map[string]int{},
    nil,
}
```

```
--- FAIL: TestStoreWins (0.00s)
    --- FAIL: TestStoreWins/it_records_wins_when_POST (0.00s)
        server_test.go:80: got 0 calls to RecordWin want 1
```

## 编写足够的代码让它通过

因为我们只是断言调用的次数而不是具体的值，所以它使我们的初始迭代值更小一些。

如果我们能够调用 `RecordWin`，我们需要修改 `PlayerStore` 接口来更新 `PlayerServer`。

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
}
```

修改后程序不能编译了

```
./main.go:17:46: cannot use InMemoryPlayerStore literal (type *InMemoryPlayerStore) as type PlayerStore in field value:
    *InMemoryPlayerStore does not implement PlayerStore (missing RecordWin method)
```

编译器告诉我们哪里出错了。我们给 `InMemoryPlayerStore` 加上那个方法。

```go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) RecordWin(name string) {}
```

尝试并运行测试，我们应该能编译代码了，但测试还是失败的。

既然 `PlayerStore` 有 `RecordWin` 方法，那我们可以在 `PlayerServer` 中调用它。

```go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
    p.store.RecordWin("Bob")
    w.WriteHeader(http.StatusAccepted)
}
```

运行测试，现应该通过了！显然 `"Bob"` 并不是我们想要发送给 `RecordWin` 的，所以让我们进一步完善测试。

## 先写测试

```go
t.Run("it records wins on POST", func(t *testing.T) {
    player := "Pepper"

    request := newPostWinRequest(player)
    response := httptest.NewRecorder()

    server.ServeHTTP(response, request)

    assertStatus(t, response.Code, http.StatusAccepted)

    if len(store.winCalls) != 1 {
        t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
    }

    if store.winCalls[0] != player {
        t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], player)
    }
})
```

现在我们知道 `winCalls` 切片中有一个元素，我们可以安全地引用第一个元素并检查它是否等于 `player`。

## 尝试运行测试

```
=== RUN   TestStoreWins/it_records_wins_on_POST
    --- FAIL: TestStoreWins/it_records_wins_on_POST (0.00s)
        server_test.go:86: did not store correct winner got 'Bob' want 'Pepper'
```

## 编写足够的代码让它通过

```go
func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]
    p.store.RecordWin(player)
    w.WriteHeader(http.StatusAccepted)
}
```

我们让 `processWin` 接收一个 `http.Request` 参数来从 URL 中获取玩家的名字。这样我们就可以用正确的值调用 `store` 来使测试通过。

## 重构

我们可以稍微精简一下这段代码，因为我们在两个地方以相同的方式获取玩家名称

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    switch r.Method {
    case http.MethodPost:
        p.processWin(w, player)
    case http.MethodGet:
        p.showScore(w, player)
    }
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
    score := p.store.GetPlayerScore(player)

    if score == 0 {
        w.WriteHeader(http.StatusNotFound)
    }

    fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
    p.store.RecordWin(player)
    w.WriteHeader(http.StatusAccepted)
}
```

即使测试显示通过了，但程序功能并没有真正完成。如果你尝试运行程序，你会发现它并不能按预期的那样工作，因为我们还没有正确实现 `PlayerStore`。没关系，通过专注 handler 我们已经确定了需要的接口，而不是妄图对它进行预先设计。

我们可以开始针对 `InMemoryPlayerStore` 编写一些测试，但在实现一种更强大的持久化存储玩家得分的方案（即数据库）之前，这只是暂时的。

我们现在要做的是针对 `PlayerServer` 和 `InMemoryPlayerStore` 编写一个集成测试来完成功能。这将让我们确保程序能正常工作，而无需直接测试 `InMemoryPlayerStore`。不仅如此，当我们开始使用数据库实现 `PlayerStore` 时，我们可以使用相同的集成测试来测试该实现。

### 集成测试

集成测试对较大型的测试很有用，但你必须牢记：

- 集成测试更难编写
- 测试失败时，可能很难知道原因（通常它是集成测试组件中的错误），因此可能更难修复
- 有时运行较慢（因为它们通常与“真实”组件一起使用，比如数据库）

因此，建议你研究一下 _金字塔测试_。

## 先写测试

长话短说，这是重构后的集成测试。

```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
    store := InMemoryPlayerStore{}
    server := PlayerServer{&store}
    player := "Pepper"

    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

    response := httptest.NewRecorder()
    server.ServeHTTP(response, newGetScoreRequest(player))
    assertStatus(t, response.Code, http.StatusOK)

    assertResponseBody(t, response.Body.String(), "3")
}
```

- 我们正在尝试集成两个组件：`InMemoryPlayerStore` 和 `PlayerServer`。
- 然后我们发起 3 个请求，为玩家记录 3 次获胜。我们并不太关心测试中的返回状态码，因为和集成得好不好无关。
- 我们真正关心的是下一个响应（所以我们用变量存储 `response`），因为我们要尝试并获得 `player` 的得分。

## 尝试运行测试

```
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:24: response body is wrong, got '123' want '3'
```

## 编写足够的代码让它通过

我会在这里采取一些自由，编写更多代码，而不是编写测试。

_这是允许的_！我们仍然有一个测试检查应该正常工作，但它不是我们正在使用的特定单元（`InMemoryPlayerStore`）。

如果我在这里犯难了，我会把我的更改恢复到测试失败时的状态，然后针对 `InMemoryPlayerStore` 编写更具体的单元测试来帮我找出解决方案。

```go
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
    return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct{
    store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
    i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return i.store[name]
}
```

- 我们需要存储数据，所以我在 `InMemoryPlayerStore` 结构中添加了 `map[string]int`
- 为方便起见，我已经让 `NewInMemoryPlayerStore` 初始化了 store，并更新了集成测试来使用它（`store := NewInMemoryPlayerStore()`）
- 代码的其余部分只是 `map` 相关的操作

集成测试通过，现在我们只需要将 `main` 改为使用 `NewInMemoryPlayerStore()`

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    server := &PlayerServer{NewInMemoryPlayerStore()}

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

构建并运行代码，然后使用 `curl` 来测试它。

- 运行几次这条命令 `curl -X POST http://localhost:5000/players/Pepper`，你换成别的玩家名称也可以
- 用 `curl http://localhost:5000/players/Pepper` 获取玩家得分

很好！你创建了一个 RESTful 风格的服务。为了实现这一目标，你需要选择一个数据存储来持久化得分数据。

- 选择一种存储机制（Bolt? Mongo? Postgres? File system?）
- 用一个 `PostgresPlayerStore` 函数来实现 `PlayerStore`
- 通过 TDD 来确保它能正常工作
- 接入集成测试中，检查它是否依然正常工作
- 最终接入到主程序中

## 总结

### `http.Handler`

- 通过实现这个接口来创建 web 服务器
- 用 `http.HandlerFunc` 把普通函数转化为 `http.Handler`
- 把 `httptest.NewRecorder` 作为一个 `ResponseWriter` 传进去，这样让你可以监视 handler 发送了什么响应
- 使用 `http.NewRequest` 构建对服务器的请求

### 接口，模拟和依赖注入

- 允许你以小步快速迭代的方式逐步构建系统
- 允许你开发需要存储 handler 而无需实际存储
- TDD 驱使你实现需要的接口

### 暂时提交有问题的代码，然后重构（然后提交到版本控制）

- 你需要将编译失败或测试失败视为一种红色警告状态，而且需要尽快摆脱它。
- 只编写必要的代码，然后重构优化代码。
- 在代码未编译或测试失败时尝试进行太多更改会让你面临陷入复杂问题的风险。
- 坚持这种小步快速迭代的方法编写测试，在处理复杂系统时，小的更改有助于提升系统可维护性。

---

作者：[Chris James](https://dev.to/quii)
译者：[pityonline](https://github.com/pityonline)
校对：[polaris1119](https://github.com/polaris1119)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
