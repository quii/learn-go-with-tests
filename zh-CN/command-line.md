# 命令行和项目结构

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/command-line)**

我们的项目负责人现在想再创建一个命令行应用。

现在，当用户输入 `Ruth wins` 时，它只需要能够记录玩家的胜出情况，最终目的是作为一个帮助用户玩扑克的工具。

产品负责人希望在两个应用程序之间共享数据库，以便玩家联盟根据新程序中记录的胜负情况进行更新。

## 这是我们目前已有的代码

我们已经有了一个用于启动 HTTP 服务器的 `main.go` 文件。在这个练习中，我们对 HTTP 服务器不感兴趣，但对它使用的抽象方法感兴趣。这取决于 `PlayerStore`。

```go
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
    GetLeague() League
}
```

上一章中我们创建了一个 `FileSystemPlayerStore` 的接口实现。我们应该可以在新的程序中重用它。

## 先做一些重构

我们的项目现在需要创建两个二进制文件：现有的 Web 服务器和命令行应用程序。

在我们投入新工作之前，我们应该构建一个项目结构来适应这一点。

目前所有代码都在同一个目录里，类似这样的：

`$GOPATH/src/github.com/your-name/my-app`

为了在 Go 中创建一个应用程序，你需要在 `package main` 中有一个 `main` 函数。到目前为止，我们所有的「域」代码都在 `package main` 中，而 `func main` 可以引用所有内容。

目前这样还好，最好不要过度使用包结构。如果你花些时间浏览标准库，你很少会看到很多文件夹和结构的形式。

庆幸的是，当你需要时，添加项目结构非常简单。

在现有项目内部创建一个 `cmd` 目录，其中包含一个 `webserver` 目录（例如 `mkdir -p cmd/webserver`）。

把 `main.go` 移到上面的目录中。

如果你安装了 `tree` 这个工具，可以运行一下看看，你的目录结构看起来应该像下面这样：

```
.
├── FileSystemStore.go
├── FileSystemStore_test.go
├── cmd
│   └── webserver
│       └── main.go
├── league.go
├── server.go
├── server_integration_test.go
├── server_test.go
├── tape.go
└── tape_test.go
```

我们现在有效地将应用程序和库代码分开了，但还需要更改一些包名。记住当你构建一个 Go 应用程序时，它的包 _必须_ 是 `main`。

将所有其它代码更改为包含名为 `poker` 的包。

最后，我们需要将此包导入 `main.go`，以便我们可以使用它来创建 web 服务器。然后可以通过 `poker.FunctionName` 来调用库代码。

你电脑上的路径可能会有所不同，但它应该类似这样：

```go
package main

import (
    "log"
    "net/http"
    "os"
    "github.com/quii/learn-go-with-tests/command-line/v1"
)

const dbFileName = "game.db.json"

func main() {
    db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("problem opening %s %v", dbFileName, err)
    }

    store, err := poker.NewFileSystemPlayerStore(db)

    if err != nil {
        log.Fatalf("problem creating file system player store, %v ", err)
    }

    server := poker.NewPlayerServer(store)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

这里的完整路径可能看起来有点别扭，但这就是将 _任何_ 可用的公共库导入到代码中的方法。

通过将域代码分离到一个单独的包中并将其提交到 GitHub 这样的公共仓库，任何 Go 开发者都可以将我们编写的功能导入并编写自己的代码。第一次尝试运行它会抱怨包不存在，但你只要运行 `go get` 就行。

[此外，用户可以在 godoc.org 上查看文档](https://godoc.org/github.com/quii/learn-go-with-tests/command-line/v1)。

### 最终检查

- 在项目根目录里面运行 `go test` 并检查它们是否仍能通过
- 进入 `cmd/webserver` 并执行 `go run main.go`
  - 访问 http://localhost:5000/league 你应该可以看到它仍然有效

### 项目框架

在开始编写测试之前，我们先添加一个项目将要构建的新应用程序。在 `cmd` 中创建另一个名为 `cli`（<ruby>命令行界面<rt>command line interface</rt><ruby>）的目录，并添加一个带有以下内容的 `main.go`：

```go
package main

import "fmt"

func main() {
    fmt.Println("Let's play poker")
}
```

我们要解决的第一个需求就是当用户输入 `{PlayerName} wins` 时记录一次胜利。

## 先写测试

我们需要创建一个名为 `CLI` 的东西，它允许我们 `Play` 扑克。它需要读取用户输入，然后将胜利记录到 `PlayerStore`。

在跑得太远之前，我们先写一个测试来检查它是否能与我们想要的 `PlayerStore` 集成。

在 `CLI_test.go` 中（在项目根目录里，不是在 `cmd` 目录中）添加以下代码：

```go
func TestCLI(t *testing.T) {
    playerStore := &StubPlayerStore{}
    cli := &CLI{playerStore}
    cli.PlayPoker()

    if len(playerStore.winCalls) !=1 {
        t.Fatal("expected a win call but didn't get any")
    }
}
```

- 我们可以使用其它测试中的 `StubPlayerStore`
- 我们将依赖关系传递给尚未存在的 `CLI` 类型
- 通过还未编写的 `PlayPoker` 方法触发游戏
- 检查是否记录了胜利

## 尝试运行测试

```
# github.com/quii/learn-go-with-tests/command-line/v2
./cli_test.go:25:10: undefined: CLI
```

## 编写最少量的代码让测试运行起来，然后检查错误输出

此时，你应该能相当自如地创建新的 `CLI` 结构，其中包含依赖项的相应字段并添加方法。

你最终应该得到这样的代码：

```go
type CLI struct {
    playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {}
```

记住我们只是试图让测试运行，所以可以按期望的方式检查测试失败：

```
--- FAIL: TestCLI (0.00s)
    cli_test.go:30: expected a win call but didn't get any
FAIL
```

## 编写足够的代码让它通过

```go
func (cli *CLI) PlayPoker() {
    cli.playerStore.RecordWin("Cleo")
}
```

这应该可以让测试通过。

接下来我们需要模拟从 `Stdin`（来自用户的输入）读取，以便记录特定玩家的胜利。

让我们扩展测试来练习一下。

## 先写测试

```go
func TestCLI(t *testing.T) {
    in := strings.NewReader("Chris wins\n")
    playerStore := &StubPlayerStore{}

    cli := &CLI{playerStore, in}
    cli.PlayPoker()

    if len(playerStore.winCalls) < 1 {
        t.Fatal("expected a win call but didn't get any")
    }

    got := playerStore.winCalls[0]
    want := "Chris"

    if got != want {
        t.Errorf("didn't record correct winner, got '%s', want '%s'", got, want)
    }
}
```

`os.Stdin` 是我们在 `main` 中用来捕获用户输入的。它实际上是一个 `*File` 类型，这意味着它实现了 `io.Reader`，现在我们知道它是一种获得文本的便捷方式。

我们在测试中使用 `strings.NewReader` 方法创建一个 `io.Reader`，用期望用户输入的内容填充它。

## 尝试运行测试

`./CLI_test.go:12:32: too many values in struct initializer`

## 编写最少量的代码让测试运行起来，然后检查错误输出

我们需要将新的依赖添加到 `CLI` 中。

```go
type CLI struct {
    playerStore PlayerStore
    in io.Reader
}
```

## 编写足够的代码让它通过

```
--- FAIL: TestCLI (0.00s)
    CLI_test.go:23: didn't record correct winner, got 'Cleo', want 'Chris'
FAIL
```

记得先做最简单的测试：

```go
func (cli *CLI) PlayPoker() {
    cli.playerStore.RecordWin("Chris")
}
```

测试通过。我们将添加另一个测试，迫使我们接下来写一些真正的代码，但首先让我们重构一下。

## 重构

我们之前在 `server_test` 中检查过是否记录了胜利，和这里一样。我们把这个断言改成一个辅助函数：

```go
func assertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
    t.Helper()

    if len(store.winCalls) != 1 {
        t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
    }

    if store.winCalls[0] != winner {
        t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], winner)
    }
}
```

现在在 `server_test.go` 和 `CLI_test.go` 中把断言都替换掉。

现在测试看起来应该类似这样：

```go
func TestCLI(t *testing.T) {
    in := strings.NewReader("Chris wins\n")
    playerStore := &StubPlayerStore{}

    cli := &CLI{playerStore, in}
    cli.PlayPoker()

    assertPlayerWin(t, playerStore, "Chris")
}
```

现在写另一个不同用户输入的测试来确保我们真正能读到它。

## 先写测试

```go
func TestCLI(t *testing.T) {

    t.Run("record chris win from user input", func(t *testing.T) {
        in := strings.NewReader("Chris wins\n")
        playerStore := &StubPlayerStore{}

        cli := &CLI{playerStore, in}
        cli.PlayPoker()

        assertPlayerWin(t, playerStore, "Chris")
    })

    t.Run("record cleo win from user input", func(t *testing.T) {
        in := strings.NewReader("Cleo wins\n")
        playerStore := &StubPlayerStore{}

        cli := &CLI{playerStore, in}
        cli.PlayPoker()

        assertPlayerWin(t, playerStore, "Cleo")
    })

}
```

## 尝试运行测试

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/record_chris_win_from_user_input
    --- PASS: TestCLI/record_chris_win_from_user_input (0.00s)
=== RUN   TestCLI/record_cleo_win_from_user_input
    --- FAIL: TestCLI/record_cleo_win_from_user_input (0.00s)
        CLI_test.go:27: did not store correct winner got 'Chris' want 'Cleo'
FAIL
```

## 编写足够的代码让它通过

我们将使用 [`bufio.Scanner`](https://golang.org/pkg/bufio/) 从 `io.Reader` 读取输入。

> bufio 包实现了 I/O 缓冲。它封装了一个 io.Reader 或 io.Writer 对象，创建了另一个对象（Reader 或 Writer），也实现了接口，并为文本 I/O 提供了缓冲和一些帮助。

把代码改成以下这样：

```go
type CLI struct {
    playerStore PlayerStore
    in          *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
    return &CLI{
        playerStore: store,
        in:          bufio.NewScanner(in),
    }
}

func (cli *CLI) PlayPoker() {
    userInput := cli.readLine()
    cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
    return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
    cli.in.Scan()
    return cli.in.Text()
}
```

现在测试应该通过了。

- `Scanner.Scan()` 会逐行读取内容。
- 然后使用 `Scanner.Text()` 来返回 scanner 读取的 `string`。
- 我们将它封装到一个名为 `readLine()` 的函数中。

现在一些测试通过了，我们应该把它嵌入到 `main` 中。记住，我们应该尽可能快地使用完全集成的工作软件。

在 `main.go` 中添加以下内容并运行它。（你可能必须调整第二个依赖项的路径来适配你的环境）

```go
package main

import (
    "fmt"
    "github.com/quii/learn-go-with-tests/command-line/v3"
    "log"
    "os"
)

const dbFileName = "game.db.json"

func main() {
    fmt.Println("Let's play poker")
    fmt.Println("Type {Name} wins to record a win")

    db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("problem opening %s %v", dbFileName, err)
    }

    store, err := poker.NewFileSystemPlayerStore(db)

    if err != nil {
        log.Fatalf("problem creating file system player store, %v ", err)
    }

    game := poker.CLI{store, os.Stdin}
    game.PlayPoker()
}
```

你应该会得到一个报错：

```
command-line/v3/cmd/cli/main.go:32:25: implicit assignment of unexported field 'playerStore' in poker.CLI literal
command-line/v3/cmd/cli/main.go:32:34: implicit assignment of unexported field 'in' in poker.CLI literal
```

这是因为我们试图在 `CLI` 中分配 `playerStore` 和 `in` 字段。这些是未导出的（私有）字段。我们可以在测试代码中执行此操作，因为测试与 `CLI`（`poker`）在同一个包中。但是 `main` 是在 `main` 包中，所以它没有访问权限。

这突出了 _整合_ 的重要性。我们理所当然地将 `CLI` 的依赖关系变为私有（因为我们不希望它们暴露给 `CLI` 的用户）但是没有为用户构建它的方法。

有没有办法早点发现这个问题？

### `package mypackage_test`

在目前为止的所有其它示例中，当我们创建一个测试文件时，我们将其声明为与我们正在测试的同一个包中。

这是可以的，这意味着在某些测试包内部功能的场合，可以访问未导出的类型。

但鉴于我们通常主张不测试包内部功能，Go 可以帮助强制执行吗？如果可以测试只能访问导出类型的代码（比如 `main`）怎么办？

当你编写包含多个包的项目时，我强烈建议测试包名称最后包含 `_test`。这样你将只能访问包中的公共类型。这有助于解决这一特定情况，也有助于强制执行仅测试公共 API 的规则。如果你仍希望测试包内部，则可以使用要测试的包进行单独测试。

TDD 的一句格言是，如果你无法测试代码，那么你的代码用户可能很难与其集成。使用 `package foo_test` 可以帮助你测试你的代码，就好像包的使用者一样导入它。

在修复 `main` 之前，让我们将 `CLI_test.go` 中的测试包更改为 `poker_test`。

如果你的 IDE 配置得好，你会突然看到很多红色提示！如果你编译它，你将得到以下错误：

```
./CLI_test.go:12:19: undefined: StubPlayerStore
./CLI_test.go:17:3: undefined: assertPlayerWin
./CLI_test.go:22:19: undefined: StubPlayerStore
./CLI_test.go:27:3: undefined: assertPlayerWin
```

我们现在遇到了关于包装设计的更多问题。为了测试，我们创建了未导出的存根和辅助函数，这些函数在 `CLI_test` 中不再可用，因为辅助函数是在 `poker` 包中的 `_test.go` 文件中定义的。

#### 我们想让存根和辅助函数都公开吗？

这是一个主观的讨论。有人可能会争辩说，你不能为方便测试而污染 API。

在 Mitchell Hashimoto 的演示文稿 [“Advanced Testing with Go”](https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=53) 中描述了 HashiCorp 如何提倡这样做以便用户可以在此基础上编写测试而无需重新发明轮子。在我们的例子中，这意味着任何使用 `poker` 包的人如果希望使用我们的代码，就不必创建自己的存根 `PlayerStore`。

我在其它项目中使用过这种技术，事实证明，在用户与软件包集成时可以非常有效地节省时间。

让我们创建一个叫 `testing.go` 的文件，并把存根和辅助函数放进去。

```go
package poker

import "testing"

type StubPlayerStore struct {
    scores   map[string]int
    winCalls []string
    league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}

func (s *StubPlayerStore) RecordWin(name string) {
    s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
    return s.league
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
    t.Helper()

    if len(store.winCalls) != 1 {
        t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
    }

    if store.winCalls[0] != winner {
        t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], winner)
    }
}

// todo for you - the rest of the helpers
```

如果你希望将辅助助程序公开给包的导入程序，你需要导出辅助程序（记住在开始时使用首字母大写的方式完成导出）。

在 `CLI` 测试中，你需要像在不同的包中使用它一样调用代码。

```go
func TestCLI(t *testing.T) {

    t.Run("record chris win from user input", func(t *testing.T) {
        in := strings.NewReader("Chris wins\n")
        playerStore := &poker.StubPlayerStore{}

        cli := &poker.CLI{playerStore, in}
        cli.PlayPoker()

        poker.AssertPlayerWin(t, playerStore, "Chris")
    })

    t.Run("record cleo win from user input", func(t *testing.T) {
        in := strings.NewReader("Cleo wins\n")
        playerStore := &poker.StubPlayerStore{}

        cli := &poker.CLI{playerStore, in}
        cli.PlayPoker()

        poker.AssertPlayerWin(t, playerStore, "Cleo")
    })

}
```

你会遇到和 `main` 一样的问题：

```
./CLI_test.go:15:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:15:39: implicit assignment of unexported field 'in' in poker.CLI literal
./CLI_test.go:25:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:25:39: implicit assignment of unexported field 'in' in poker.CLI literal
```

解决这个问题最简单的办法是创建一个构造函数，就像我们对其它类型一样：

```go
func NewCLI(store PlayerStore, in io.Reader) *CLI {
    return &CLI{
        playerStore: store,
        in:          in,
    }
}
```

使用构造函数后，测试应该可以通过了。

最后，我们可以回到新的 `main.go` 并使用刚刚创建的构造函数：

```go
game := poker.NewCLI(store, os.Stdin)
```

尝试运行一下，输入“Bob wins”。

### 重构

现在在打开文件并从其内容创建 `FileSystemStore` 的代码中有些重复的地方。这是当前设计中的一个小瑕疵，所以我们应该创建一个函数来封装从路径打开文件并返回 `PlayerStore`。

```go
func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, error) {
    db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        return nil, fmt.Errorf("problem opening %s %v", path, err)
    }

    store, err := NewFileSystemPlayerStore(db)

    if err != nil {
        return nil, fmt.Errorf("problem creating file system player store, %v ", err)
    }

    return store, nil
}
```

现在重构我们的两个应用程序以使用此函数来创建 store。

#### CLI 程序代码

```go
package main

import (
    "log"
    "os"
    "fmt"

    "github.com/quii/learn-go-with-tests/command-line/v3"
)

const dbFileName = "game.db.json"

func main() {
    store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Let's play poker")
    fmt.Println("Type {Name} wins to record a win")
    poker.NewCLI(store, os.Stdin).PlayPoker()
}
```

#### Web 服务器代码

```go
package main

import (
    "github.com/quii/learn-go-with-tests/command-line/v3"
    "log"
    "net/http"
)

const dbFileName = "game.db.json"

func main() {
    store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

    if err != nil {
        log.Fatal(err)
    }

    server := poker.NewPlayerServer(store)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

注意对称性：尽管用户界面不同，但设置几乎相同。

## 总结

### 包结构

本章讲述了我们如何重用已有的代码来创建两个应用程序。为了做到这一点，我们需要更新包结构，以便为各自的 `main` 包提供单独的目录。

在此过程中，我们遇到了由于未导出的值导致的集成问题，因此这进一步证明了在小步重构中并经常进行测试的价值。

我们学习了如何借助 `mypackage_test` 这种形式创建一个测试环境，这与别人集成你的代码体验是一样的，可以帮助你捕获集成问题并查看代码是否易用。

### 读取用户输入

我们看到了用 `os.Stdin` 读取输入多么简单，因为它实现了 `io.Reader`。我们使用 `bufio.Scanner` 轻松地逐行读取用户输入。

### 简单抽象让代码复用更容易

几乎不费任何力气就将 `PlayerStore` 集成到我们的新应用程序中（同时我们对包进行了调整），因为我们决定公开我们的存根版本，随后测试也非常简单。

---

作者：[Chris James](https://dev.to/quii)
译者：[pityonline](https://github.com/pityonline)
校对：[译者ID](https://github.com/译者ID)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
