# JSON，路由和嵌入

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/json)**

在[上一章](http-server.md)中，我们创建了一个 web 服务器用来储存多少玩家赢了（游戏）。

我们的项目负责人有个新的需求，要求有一个叫 `/league`（联盟）的新的端点（endpoint），它可以返回一个玩家清单。她想让它返回一个 JSON 格式的数据。

## 这是我们目前已有的代码

```go
// server.go
package main

import (
    "fmt"
    "net/http"
)

type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
}

type PlayerServer struct {
    store PlayerStore
}

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

```go
// InMemoryPlayerStore.go
package main

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
    return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
    store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
    i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return i.store[name]
}

```

```go
// main.go
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

你可以在本章开始位置的链接中找到对应的测试用例。

我们将从创建一个联盟成员表端点开始。

## 先写测试

我们会基于已有的测试用例和模拟的 `PlayerStore` 进行扩充。

```go
func TestLeague(t *testing.T) {
    store := StubPlayerStore{}
    server := &PlayerServer{&store}

    t.Run("it returns 200 on /league", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/league", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response.Code, http.StatusOK)
    })
}
```

在开始为真实得分和 JSON 操心之前，我们先尝试对实现目标做最小的修改。从最简单的开始，检查我们访问 `/league` 能否返回一个 `OK` 的响应。

## 尝试运行测试

```
=== RUN   TestLeague/it_returns_200_on_/league
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range

goroutine 6 [running]:
testing.tRunner.func1(0xc42010c3c0)
    /usr/local/Cellar/go/1.10/libexec/src/testing/testing.go:742 +0x29d
panic(0x1274d60, 0x1438240)
    /usr/local/Cellar/go/1.10/libexec/src/runtime/panic.go:505 +0x229
github.com/quii/learn-go-with-tests/json-and-io/v2.(*PlayerServer).ServeHTTP(0xc420048d30, 0x12fc1c0, 0xc420010940, 0xc420116000)
    /Users/quii/go/src/github.com/quii/learn-go-with-tests/json-and-io/v2/server.go:20 +0xec
```

你的 `PlayerServer` 应该产生一个这样的 panic。通过栈跟踪找到指向 `server.go` 的那行代码。

```go
player := r.URL.Path[len("/players/"):]
```

在上一章中我们提到这是一种天真幼稚的路由处理方式。事实上，它试图把 `/league` 字符串进行路径分割，因此会报错 `slice bounds out of range`（切片超出范围）。

## 编写足够的代码让它通过

Go 有一个内置的路由机制叫做 [`ServeMux`](https://golang.org/pkg/net/http/#ServeMux)（request multiplexer，多路请求复用器），它允许你将 `http.Handler` 附加到特定的请求路径。

让我们暂时忽略一些已知问题，先以尽可能最快的方式让测试通过，因为一旦我们知道测试通过就可以安全地重构它。

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    router := http.NewServeMux()

    router.Handle("/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    router.Handle("/players/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        player := r.URL.Path[len("/players/"):]

        switch r.Method {
        case http.MethodPost:
            p.processWin(w, player)
        case http.MethodGet:
            p.showScore(w, player)
        }
    }))

    router.ServeHTTP(w, r)
}
```

- 当请求开始时，我们创建了一个路由，然后我们告诉它 `x` 路径使用 `y` handler。
- 那么对于我们的新端点 `/league` 被请求时，我们用 `http.HandlerFunc` 和一个匿名函数来响应 `w.WriteHeader(http.StatusOK)` 使测试通过。
- 对于 `/players/` 路由我们只需剪贴代码并粘贴到另一个 `http.HandlerFunc`。
- 最终，我们通过调用新路由的 `ServeHTTP` 方法处理到来的请求（注意到 `ServeMux` 也是一个 `http.Handler` 了吗？）

现在测试应该可以通过了。

## 重构

`ServeHTTP` 看起来很大，我们可以通过重构 handlers 分离成独立的方法。

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    router := http.NewServeMux()
    router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    router.ServeHTTP(w, r)
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

    switch r.Method {
    case http.MethodPost:
        p.processWin(w, player)
    case http.MethodGet:
        p.showScore(w, player)
    }
}
```

把一个路由作为一个请求来处理并调用它挺奇怪的（并且效率低下）。我们想要的理想情况是有一个 `NewPlayerServer` 这样的函数，它可以取得依赖并进行一次创建路由的设置。每个请求都可以使用该路由的一个实例。

```go
type PlayerServer struct {
    store  PlayerStore
    router *http.ServeMux
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
    p := &PlayerServer{
        store,
        http.NewServeMux(),
    }

    p.router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    p.router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    return p
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    p.router.ServeHTTP(w, r)
}
```

- `PlayerServer` 现在需要储存一个路由。
- 我们已经把创建 `ServeHTTP` 路由的动作移到 `NewPlayerServer`，这样只需要完成一次，而不是每次请求都要做。
- 在所有测试和程序代码中，用到 `PlayerServer{&store}` 的地方你需要更新为 `NewPlayerServer{&store}`。

### 最后一次重构

试着把代码改成以下这样：

```go
type PlayerServer struct {
    store  PlayerStore
    http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
    p := new(PlayerServer)

    p.store = store

    router := http.NewServeMux()
    router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    p.Handler = router

    return p
}
```

最后确定你 **删除** 了 `func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request)`，因为不再需要它了。

## 嵌入

我们更改了 `PlayerServer` 的第二个属性，删除了命名属性 `router http.ServeMux`，并用 `http.Handler` 替换了它；这被称为 _嵌入_。

> Go 没有提供典型的，类型驱动的子类化概念，但它具有通过在结构或接口中嵌入类型来“借用”一部分实现的能力。

[高效 Go - 嵌入](https://golang.org/doc/effective_go.html#embedding)

这意味着我们的 `PlayerServer` 现在已经有了 `http.Handler` 所有的方法，也就是 `ServeHTTP`。

为了“填充” `http.Handler`，我们将它分配给我们在 `NewPlayerServer` 中创建的 `router`。我们可以这样做是因为 `http.ServeMux` 具有 `ServeHTTP` 方法。

这允许我们删除我们的 `ServeHTTP` 方法，因为我们已经通过嵌入类型公开了它。

嵌入是一个非常有意思的语法特性。你可以用它将接口组成新的接口。

```go
type Animal interface {
    Eater
    Sleeper
}
```

你也可以用它混合类型，而不仅是接口，正如你所期望的，如果你嵌入一个混合类型，你将可以访问它所有的公共接口和字段。

### 有任何缺点吗？

你必须小心使用嵌入类型，因为你将公开所有嵌入类型的公共方法和字段。在我们的例子中它是可以的，因为我们只是嵌入了 `http.Handler` 这个 _接口_。

如果我们懒一点，嵌入了 `http.ServeMux`（混合类型），它仍然可以工作 _但_ `PlayerServer` 的用户就可以给我们的服务器添加新路由了，因为 `Handle(path, handler)` 会公开。

**嵌入类型时，真正要考虑的是对你公开的 API 有什么影响。**

滥用嵌入最终会污染你的 API，并暴露你的类型的内部信息，这是个常见的错误。

现在我们重新构建了我们的应用，我们可以轻易地添加新的路由，并让 `/league` 端点有了一个新的开始。现我我们需要让它返回一些有用的信息。

我们应该返回一些类似这样的 JSON 数据：

```json
[
   {
      "Name":"Bill",
      "Wins":10
   },
   {
      "Name":"Alice",
      "Wins":15
   }
]
```

## 先写测试

我们先尝试把响应解析为有意义的信息。

```go
func TestLeague(t *testing.T) {
    store := StubPlayerStore{}
    server := NewPlayerServer(&store)

    t.Run("it returns 200 on /league", func(t *testing.T) {
        request, _ := http.NewRequest(http.MethodGet, "/league", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        var got []Player

        err := json.NewDecoder(response.Body).Decode(&got)

        if err != nil {
            t.Fatalf ("Unable to parse response from server '%s' into slice of Player, '%v'", response.Body, err)
        }

        assertStatus(t, response.Code, http.StatusOK)
    })
}
```

### 为什么不测试 JSON 字符串？

你可以争论一个更简单的初始步骤就是断言响应体有一个特定的 JSON 字符串。

以我的经验，断言 JSON 字符串的测试有以下问题。

- *脆弱*。如果你改了数据模型，测试将会失败。
- *难以调试*。在比较两个 JSON 字符串时，很难理解真正的问题是什么。
- *意图不佳*。当输出应该是 JSON 时，真正重要的是数据究竟是什么，而不是它的编码方式。
- *重复测试标准库*。没有必要测试标准库如何输出 JSON，它已经过测试。不要测试别人的代码。

相反，我们应该将 JSON 解析为与我们测试相关的数据结构。

### 数据建模

鉴于 JSON 数据模型，它看起来像我们需要一个带有一些字段的 `Player` 数组，因此我们创建了一个新类型来捕获这个数据模型。

```go
type Player struct {
    Name string
    Wins int
}
```

### JSON 解码

```go
var got []Player
err := json.NewDecoder(response.Body).Decode(&got)
```

我们创建一个 `Decoder`（解码器），然后调用 `encoding/json` 包里的 `Decode` 方法来把 JSON 解析为我们的数据模型。在我们的例子中，`Decoder` 需要一个 `io.Reader` 来读取响应体。

`Decode` 取到我们正在尝试解析的东西的地址，这就是为什么我们要在之前的行中声明 `Player` 的空切片。

解析 JSON 可能会失败，所以 `Decode` 可以返回一个 `error`。如果失败了，继续测试没有意义，如果发生错误，用 `t.Fatalf` 停止测试并检查错误。请注意，我们打印了响应正文以及错误，因为对于运行测试的人来说，看看哪些字符串不能被解析很重要。

## 尝试运行测试

```
=== RUN   TestLeague/it_returns_200_on_/league
    --- FAIL: TestLeague/it_returns_200_on_/league (0.00s)
        server_test.go:107: Unable to parse response from server '' into slice of Player, 'unexpected end of JSON input'
```

我们的端点现在不能返回一个响应实体，所以不能被解析为 JSON。

## 编写足够的代码让它通过

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    leagueTable := []Player{
        {"Chris", 20},
    }

    json.NewEncoder(w).Encode(leagueTable)

    w.WriteHeader(http.StatusOK)
}
```

现在测试通过。

### 编码和解码

请注意标准库中的的有意思的对称性。

- 要创建一个 `Encoder`，你需要一个 `http.ResponseWriter` 实现的 `io.Writer`。
- 要创建一个 `Decoder`，你需要一个 `io.Reader`，由我们的响应的 `Body` 字段实现。

在本书中，我们使用了 `io.Writer`，这是它在标准库中流行及许多库如何轻松使用它的的另一个示范。

## 重构

在 handler 和获得 `leagueTable` 之间考虑引入一个拆分是很好的，因为我们知道很快就不会硬编码了。

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(p.getLeagueTable())
    w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) getLeagueTable() []Player{
    return []Player{
        {"Chris", 20},
    }
}
```

接下来我们将扩展我们的测试，以便我们可以准确地控制想要的数据。

## 先写测试

我们可以更新测试，以断言联盟表中包含一些我们将存储在商店中的玩家。

更新 `StubPlayerStore` 让它存储一个联盟，这只是一个“玩家”的切片类型。我们将存储我们的预期数据。

```go
type StubPlayerStore struct {
    scores   map[string]int
    winCalls []string
    league []Player
}
```

接下来更新我们目前的测试，将一些玩家放入我们存放的联盟属性中，并声明他们从我们的服务器返回。

```go
func TestLeague(t *testing.T) {

    t.Run("it returns the league table as JSON", func(t *testing.T) {
        wantedLeague := []Player{
            {"Cleo", 32},
            {"Chris", 20},
            {"Tiest", 14},
        }

        store := StubPlayerStore{nil, nil, wantedLeague,}
        server := NewPlayerServer(&store)

        request, _ := http.NewRequest(http.MethodGet, "/league", nil)
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        var got []Player

        err := json.NewDecoder(response.Body).Decode(&got)

        if err != nil {
            t.Fatalf("Unable to parse response from server '%s' into slice of Player, '%v'", response.Body, err)
        }

        assertStatus(t, response.Code, http.StatusOK)

        if !reflect.DeepEqual(got, wantedLeague) {
            t.Errorf("got %v want %v", got, wantedLeague)
        }
    })
}
```

## 尝试运行测试

```
./server_test.go:33:3: too few values in struct initializer
./server_test.go:70:3: too few values in struct initializer
```

## 编写最少量的代码让测试运行起来，然后检查错误输出

你需要更新其它测试，因为我们在 `StubPlayerStore` 中有了一个新字段；将其设置为 nil 以进行其它测试。

尝试再次运行测试，你应该得到

```
=== RUN   TestLeague/it_returns_the_league_table_as_JSON
    --- FAIL: TestLeague/it_returns_the_league_table_as_JSON (0.00s)
        server_test.go:124: got [{Chris 20}] want [{Cleo 32} {Chris 20} {Tiest 14}]
```

## 编写足够的代码让它通过

我们知道数据是存储在 `StubPlayerStore` 里，我们已经把它抽象为 `PlayerStore` 接口。我们需要更新它以便任何人传入一个 `PlayerStore` 可以提供联盟的数据。

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
    GetLeague() []Player
}
```

现在我们可以更新 handler 代码来调用它，而不是返回一个硬编码列表。删除 `getLeagueTable()` 方法，然后更新 `leagueHandler` 来调用 `GetLeague()`。

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(p.store.GetLeague())
    w.WriteHeader(http.StatusOK)
}
```

尝试重新运行测试。

```
# github.com/quii/learn-go-with-tests/json-and-io/v4
./main.go:9:50: cannot use NewInMemoryPlayerStore() (type *InMemoryPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *InMemoryPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_integration_test.go:11:27: cannot use store (type *InMemoryPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *InMemoryPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_test.go:36:28: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *StubPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_test.go:74:28: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *StubPlayerStore does not implement PlayerStore (missing GetLeague method)
./server_test.go:106:29: cannot use &store (type *StubPlayerStore) as type PlayerStore in argument to NewPlayerServer:
    *StubPlayerStore does not implement PlayerStore (missing GetLeague method)
```

编译器报怨因为 `InMemoryPlayerStore` 和 `StubPlayerStore` 没有我们刚添加到接口的新方法。

对于 `StubPlayerStore` 很简单，只要返回我们之前添加的 `league` 字段即可。

```go
func (s *StubPlayerStore) GetLeague() []Player {
    return s.league
}
```

这里提示一下 `InMemoryStore` 是如何实现的。

```go
type InMemoryPlayerStore struct {
    store map[string]int
}
```

此时通过遍历映射来正确实现 `GetLeague` 非常简单，但记住，我们只是试图 _编写最少量的代码来使测试通过_。

所以让我们现在就先让编译通过，暂放一下未完整实现 `InMemoryStore` 的问题。

```go
func (i *InMemoryPlayerStore) GetLeague() []Player {
    return nil
}
```

这实际上告诉我们的是 _稍后_ 我们要测试这个，但现在先不管它。

尝试重新运行测试，编译器和测试应该都通过了！

## 重构

测试代码不能很好地表达我们的意图，并且有很多我们可以重构的样板文件。

```go
t.Run("it returns the league table as JSON", func(t *testing.T) {
    wantedLeague := []Player{
        {"Cleo", 32},
        {"Chris", 20},
        {"Tiest", 14},
    }

    store := StubPlayerStore{nil, nil, wantedLeague,}
    server := NewPlayerServer(&store)

    request := newLeagueRequest()
    response := httptest.NewRecorder()

    server.ServeHTTP(response, request)

    got := getLeagueFromResponse(t, response.Body)
    assertStatus(t, response.Code, http.StatusOK)
    assertLeague(t, got, wantedLeague)
})
```

这些是新的辅助函数

```go
func getLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
    t.Helper()
    err := json.NewDecoder(body).Decode(&league)

    if err != nil {
        t.Fatalf("Unable to parse response from server '%s' into slice of Player, '%v'", body, err)
    }

    return
}

func assertLeague(t *testing.T, got, want []Player) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}

func newLeagueRequest() *http.Request {
    req, _ := http.NewRequest(http.MethodGet, "/league", nil)
    return req
}
```

我们需要为服务器工作做的最后一件事是确保我们在响应中返回一个 `content-type` 头（HTTP header），这样机器就能识别出我们正在返回 `JSON`。

## 先写测试

向已有测试中添加这个断言

```go
if response.Header().Get("content-type") != "application/json" {
    t.Errorf("response did not have content-type of application/json, got %v", response.HeaderMap)
}
```

## 尝试运行测试

```
=== RUN   TestLeague/it_returns_the_league_table_as_JSON
    --- FAIL: TestLeague/it_returns_the_league_table_as_JSON (0.00s)
        server_test.go:124: response did not have content-type of application/json, got map[Content-Type:[text/plain; charset=utf-8]]
```

## 编写足够的代码让它通过

更新 `leagueHandler`

```go
func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(p.store.GetLeague())
}
```

测试应该可以通过。

## 重构

给 `assertContentType` 添加一个辅助函数。

```go
const jsonContentType = "application/json"

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
    t.Helper()
    if response.Header().Get("content-type") != want {
        t.Errorf("response did not have content-type of %s, got %v", want, response.HeaderMap)
    }
}
```

在测试中使用它。

```go
assertContentType(t, response, jsonContentType)
```

现在我们已经对 `PlayerServer` 进行了整理，现在我们可以把注意力转向 `InMemoryPlayerStore` 因为现在如果我们试图给产品负责人演示 `/league` 将不起作用。

获得信心的最快方法是增加集成测试，我们可以访问新的端点并检查我们是否从 `/league` 获得了正确的响应。

## 先写测试

我们可以使用 `t.Run` 来分解这个测试，我们可以重用服务器测试中的助手 —— 再次显示重构测试的重要性。

```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
    store := NewInMemoryPlayerStore()
    server := NewPlayerServer(store)
    player := "Pepper"

    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
    server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

    t.Run("get score", func(t *testing.T) {
        response := httptest.NewRecorder()
        server.ServeHTTP(response, newGetScoreRequest(player))
        assertStatus(t, response.Code, http.StatusOK)

        assertResponseBody(t, response.Body.String(), "3")
    })

    t.Run("get league", func(t *testing.T) {
        response := httptest.NewRecorder()
        server.ServeHTTP(response, newLeagueRequest())
        assertStatus(t, response.Code, http.StatusOK)

        got := getLeagueFromResponse(t, response.Body)
        want := []Player{
            {"Pepper", 3},
        }
        assertLeague(t, got, want)
    })
}
```

## 尝试运行测试

```
=== RUN   TestRecordingWinsAndRetrievingThem/get_league
    --- FAIL: TestRecordingWinsAndRetrievingThem/get_league (0.00s)
        server_integration_test.go:35: got [] want [{Pepper 3}]
```

## 编写足够的代码让它通过

当你调用 `GetLeague()` 时，`InMemoryPlayerStore` 返回 `nil`，所以我们需要修复它。

```go
func (i *InMemoryPlayerStore) GetLeague() []Player {
    var league []Player
    for name, wins := range i.store {
        league = append(league, Player{name, wins})
    }
    return league
}
```

我们所需要做的就是遍历映射并将每个键 / 值对转换为一个 `Player`。

现在测试应该可以通过。

## 总结

我们继续使用 TDD 安全地迭代了我们的程序，使其能够通过路由以可维护的方式支持新端点，现在它可以为我们的客户返回 JSON。在下一章中，我们将介绍持久化数据和对联盟排序。

我们所涵盖的内容：

- **路由**。标准库为你提供了易于使用的类型来处理路由。它完全支持 `http.Handler` 接口，因为你可以将路由分配给 `Handler`，而路由本身也是 `Handler`。它没有你可能期望的某些特性，例如路径变量（例如 `/users/{id}`）。你可以自己轻易地解析这些信息，但如果它成了负担，你可能会考虑查看其它路由库。大多数流行的库都坚持标准库的实现 `http.Handler` 的理念。
- **类型嵌入**。我们对这项技术略有提及，但你可以[从 Effective Go 了解更多信息](https://golang.org/doc/effective_go.html#embedding)。如果你应该从中得到一个收获，那就是它极其有用，但是 _总是考虑你的公开 API，只有适合被公开的才公开_。
- ** JSON 反序列化和序列化**。标准库使得序列化和反序列化数据变得非常简单。它也是开放的配置，你可以根据需要自定义这些数据转换的工作方式。

---

作者：[Chris James](https://dev.to/quii)
译者：[pityonline](https://github.com/pityonline)
校对：[polaris1119](https://github.com/polaris1119)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
