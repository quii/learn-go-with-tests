# IO 和排序

[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/io)

在[上一章](json.md)中，我们通过添加新的服务器访问地址 `/league` 来迭代我们的应用程序。在此过程中，我们学习了如何处理 JSON、嵌入类型和路由。

服务器重启后软件会丢失所有得分，产品负责人对此感到不安。这是因为我们存储的实现是在内存里。对于我们没有解释 `/league` 的访问地址应该按赢的次数排序返回玩家列表，她也很不满意。

## 目前为止的代码

```go
// server.go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// PlayerStore stores score information about players
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
    GetLeague() []Player
}

// Player stores a name with a number of wins
type Player struct {
    Name string
    Wins int
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
    store PlayerStore
    http.Handler
}

const jsonContentType = "application/json"

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(store PlayerStore) *PlayerServer {
    p := new(PlayerServer)

    p.store = store

    router := http.NewServeMux()
    router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    p.Handler = router

    return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(p.store.GetLeague())
    w.Header().Set("content-type", jsonContentType)
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

func (i *InMemoryPlayerStore) GetLeague() []Player {
    var league []Player
    for name, wins := range i.store {
        league = append(league, Player{name, wins})
    }
    return league
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
    i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return i.store[name]
}
```

```
// main.go
package main

import (
    "log"
    "net/http"
)

func main() {
    server := NewPlayerServer(NewInMemoryPlayerStore())

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

你可以在本章顶部的链接中找到相应的测试。

## 存储数据

满足这个需求的数据库有很多，但我们会使用一种非常简单的方法。我们将把这个应用程序的数据以 JSON 的格式存储到文件中。

这使得数据具有很强的可移植性，并且实现起来相对简单。

它的伸缩性不高，但考虑到这是一个原型，至少现在是没问题的。如果我们的环境变得不再合适，换成其它的存储方式也会非常简单，因为我们使用的是 `PlayerStore` 的抽象。

我们将暂时保留 `InMemoryPlayerStore`，以便在开发新的存储实现时还能通过集成测试。一旦我们确信新实现足以通过集成测试，我们会替换然后删除 `InMemoryPlayerStore`。

## 首先编写测试

现在，你应该已经熟悉以下标准库相关的接口，用于读取数据（`io.Reader`）、写入数据（`io.Writer`）的接口，以及如何使用标准库来测试这些函数，而不必使用真正的文件。

为了完成这项工作，我们需要实现 `PlayerStore`，因此我们调用需要实现的方法来编写测试。我们将从 `GetLeague` 开始。

```go
func TestFileSystemStore(t *testing.T) {

    t.Run("/league from a reader", func(t *testing.T) {
        database := strings.NewReader(`[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)

        store := FileSystemStore{database}

        got := store.GetLeague()

        want := []Player{
            {"Cleo", 10},
            {"Chris", 33},
        }

        assertLeague(t, got, want)
    })
}
```

我们使用 `strings.NewReader` 会返回一个 `Reader`，这是我们的 `FileSystemStore` 函数中用来读取数据的。在 `main` 中我们将打开一个文件，它也是一个 `Reader`。

## 尝试运行测试

```
# github.com/quii/learn-go-with-tests/json-and-io/v7
./FileSystemStore_test.go:15:12: undefined: FileSystemStore
```

## 编写最少量的代码让测试运行起来，然后检查错误输出

让我们在新的文件中定义 `FileSystemStore`

```go
type FileSystemStore struct {}
```

再次尝试运行测试

```
# github.com/quii/learn-go-with-tests/json-and-io/v7
./FileSystemStore_test.go:15:28: too many values in struct initializer
./FileSystemStore_test.go:17:15: store.GetLeague undefined (type FileSystemStore has no field or method GetLeague)
```

报错是因为我们传入了不需要的 `Reader` 参数，并且 `GetLeague` 函数还没有定义。

```go
type FileSystemStore struct {
    database io.Reader
}

func (f *FileSystemStore) GetLeague() []Player {
    return nil
}
```

再试一次...

```
=== RUN   TestFileSystemStore//league_from_a_reader
    --- FAIL: TestFileSystemStore//league_from_a_reader (0.00s)
        FileSystemStore_test.go:24: got [] want [{Cleo 10} {Chris 33}]
```

## 编写足够的代码使测试通过

我们之前已经从 `reader` 中读取了 JSON 数据

```go
func (f *FileSystemStore) GetLeague() []Player {
    var league []Player
    json.NewDecoder(f.database).Decode(&league)
    return league
}
```

现在测试应该通过了。

## 重构

我们以前就这样做过！服务器的测试代码必须从响应中解码 JSON 数据。

我们试着把它提炼为一个函数。

创建一个名为 `league.go` 的新文件，输入以下代码。

```go
func NewLeague(rdr io.Reader) ([]Player, error) {
    var league []Player
    err := json.NewDecoder(rdr).Decode(&league)
    if err != nil {
        err = fmt.Errorf("problem parsing league, %v", err)
    }

    return league, err
}
```

在我们的实现和 `server_test.go` 的辅助函数 `getLeagueFromResponse` 中调用这个函数

```go
func (f *FileSystemStore) GetLeague() []Player {
    league, _ := NewLeague(f.database)
    return league
}
```

我们还没有处理解析错误的方法，但是我们还是继续向前推进吧。

## 寻找问题

我们的实现中有一个缺陷。首先注意 `io.Reader` 是如何定义的。

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

你可以想象它一个一个字节读取文件直到结束。如果你再读一遍会发生什么？

在当前测试的末尾添加以下内容。

```go
// read again
got = store.GetLeague()
assertLeague(t, got, want)
```

我们希望它通过测试，但是如果你运行会发现它并没有通过。

这里的问题是我们的 `Reader` 已经到了结尾，没什么可读的了。我们需要一种方法让它回到开始位置。

[ReadSeeker](https://golang.org/pkg/io/#ReadSeeker) 是标准库中的另一个可以提供帮助的接口。

```go
type ReadSeeker interface {
    Reader
    Seeker
}
```

还记得嵌入吗？这是由 `Reader` 和 [`Seeker`](https://golang.org/pkg/io/#Seeker) 组成的接口

```go
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
```

这感觉不错，我们可以更改 `FileSystemStore` 来替代这个接口吗？

```go
type FileSystemStore struct {
    database io.ReadSeeker
}

func (f *FileSystemStore) GetLeague() []Player {
    f.database.Seek(0, 0)
    league, _ := NewLeague(f.database)
    return league
}
```

尝试运行测试，它现在通过了！很高兴我们在测试中使用的 `string.NewReader` 也实现了 `ReadSeeker`，所以我们不需要做任何其他的改变。

接下来我们将实现 `GetPlayerScore`。

## 首先编写测试

```go
t.Run("get player score", func(t *testing.T) {
    database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)

    store := FileSystemPlayerStore{database}

    got := store.GetPlayerScore("Chris")

    want := 33

    if got != want {
        t.Errorf("got %d want %d", got, want)
    }
})
```

## 尝试运行测试

`./FileSystemStore_test.go:38:15: store.GetPlayerScore undefined (type FileSystemPlayerStore has no field or method GetPlayerScore)`

## 编写最少量的代码让测试运行起来，然后检查错误输出

我们需要将方法添加到新类型中，以便编译测试。

```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
    return 0
}
```

现在它可以编译并且测试失败

```
=== RUN   TestFileSystemStore/get_player_score
    --- FAIL: TestFileSystemStore//get_player_score (0.00s)
        FileSystemStore_test.go:43: got 0 want 33
```

## 编写足够的代码使测试通过

我们可以遍历 `league` 寻找玩家并返回他们的得分

```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

    var wins int

    for _, player := range f.GetLeague() {
        if player.Name == name {
            wins = player.Wins
            break
        }
    }

    return wins
}
```

## 重构

你会看到许多辅助函数需要重构，这些将留给你来实现

 ```go
 t.Run("/get player score", func(t *testing.T) {
    database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)

    store := FileSystemPlayerStore{database}

    got := store.GetPlayerScore("Chris")
    want := 33
    assertScoreEquals(t, got, want)
})
 ```

最后，我们需要用 `RecordWin` 来记录得分。

## 首先编写测试

我们的方法写的相当短视。我们不能（很容易地）只更新文件中 JSON 的一「行」。我们需要在每次写入时存储整个数据新的表现形式。

我们应该怎么编写？我们通常会使用一个 `Writer`，但我们已经有了 `ReadSeeker`。我们可能有两个依赖项，但是标准库已经为我们提供了一个接口 `ReadWriteSeeker`，我们需要对文件做的处理它都可以满足。

我们来更新一下类型

```go
type FileSystemPlayerStore struct {
    database io.ReadWriteSeeker
}
```

查看是否通过编译

```
./FileSystemStore_test.go:15:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
./FileSystemStore_test.go:36:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
```

`strings.Reader` 没有实现 `ReadWriteSeeker` 并不奇怪，这时我们该怎么办呢？

我们有两个选择

- 为每个测试创建一个临时文件。`*os.File` 实现 `ReadWriteSeeker`。好处是它变得更像集成测试，我们真的是从文件系统中读取和写入，所以我们对此更有信心。缺点是我们更喜欢单元测试，因为它们更快而且通常更简单。我们还需要做更多关于创建临时文件的工作，然后确保在测试之后删除它们。
- 使用第三方库。[github.com/mattetti](https://github.com/quii/learn-go-with-tests/tree/cc354086d4d8a9304914be9ac6c5244c6d32d510/Mattetti/README.md) 已经编写了一个 [filebuffer](https://github.com/mattetti/filebuffer) 库，它实现了我们需要的接口，并且不触及文件系统。

这两种选择都没有问题，但是如果选择使用第三方库，我将不得不解释依赖管理！所以还是用文件代替吧。

在添加测试之前，我们需要通过用 `os.File` 替换 `strings.Reader` 来使其他测试编译通过。

让我们创建一个辅助函数，它将创建包含一些数据的临时文件

```go
func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
    t.Helper()

    tmpfile, err := ioutil.TempFile("", "db")

    if err != nil {
        t.Fatalf("could not create temp file %v", err)
    }

    tmpfile.Write([]byte(initialData))

    removeFile := func() {
        os.Remove(tmpfile.Name())
    }

    return tmpfile, removeFile
}
```

[TempFile](https://golang.org/pkg/io/ioutil/#TempDir) 创建一个临时文件供我们使用。我们传入的 `"db"` 值是在它将创建的随机文件名上加上的前缀。这是为了确保它不会与其他文件发生意外冲突。

你会注意到，我们不仅返回 `ReadWriteSeeker`（文件），而且还返回一个函数。我们需要确保在测试完成后删除该文件。我们不希望将文件的细节泄露到测试中，因为它很容易出错，对读者来说也没什么意思。通过返回 `removeFile` 函数，我们可以处理辅助函数中的细节，调用者只需运行 `deferred cleanDatabase()`。

```go
func TestFileSystemStore(t *testing.T) {

    t.Run("league from a reader", func(t *testing.T) {
        database, cleanDatabase := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
        defer cleanDatabase()

        store := FileSystemPlayerStore{database}

        got := store.GetLeague()

        want := []Player{
            {"Cleo", 10},
            {"Chris", 33},
        }

        assertLeague(t, got, want)

        // read again
        got = store.GetLeague()
        assertLeague(t, got, want)
    })

    t.Run("get player score", func(t *testing.T) {
        database, cleanDatabase := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
        defer cleanDatabase()

        store := FileSystemPlayerStore{database}

        got := store.GetPlayerScore("Chris")
        want := 33
        assertScoreEquals(t, got, want)
    })
}
```

运行测试，他们应该可以通过了！这里有大量的更改，但是现在感觉我们已经完成了接口定义，从现在开始添加新的测试应该非常容易了。

让我们执行第一次迭代，为现有的玩家记录一次胜利

```go
t.Run("store wins for existing players", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    store.RecordWin("Chris")

    got := store.GetPlayerScore("Chris")
    want := 34
    assertScoreEquals(t, got, want)
})
```

## 尝试运行测试

`./FileSystemStore_test.go:67:8: store.RecordWin undefined (type FileSystemPlayerStore has no field or method RecordWin)`

## 编写最少量的代码让测试运行起来，然后检查错误输出

添加新的方法

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {

}
```

```
=== RUN   TestFileSystemStore/store_wins_for_existing_players
    --- FAIL: TestFileSystemStore/store_wins_for_existing_players (0.00s)
        FileSystemStore_test.go:71: got 33 want 34
```

我们的实现是空的，因此旧的得分将会返回。

## 编写足够的代码使测试通过

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
    league := f.GetLeague()

    for i, player := range league {
        if player.Name == name {
            league[i].Wins++
        }
    }

    f.database.Seek(0,0)
    json.NewEncoder(f.database).Encode(league)
}
```

你可能会问，为什么我要用 `league[i].Wins++` 而不是 `player.Wins++`。

当你在一个切片上取值时，将返回当前循环的索引（我们示例中的 `i`）和该索引中的元素的副本。更改副本 `Wins` 的值不会对我们迭代的 `league` 产生任何影响。因此，我们需要通过使用 `league[i]` 来获取对实际值的引用，然后更改该值。

如果你运行这些测试，它们应该可以通过了。

## 重构

在 `GetPlayerScore` 和 `RecordWin` 中，我们遍历 `[]Player`，按名称查找 `player`。

我们可以在 `FileSystemStore` 的内部重构这个公共代码，但对我来说，它可能还有用，我们可以将其提升为新的类型。到目前为止，操作「League」都是用 `[]Player`，但我们可以创造一种新的类型 `League`。这使其他开发人员更容易理解，然后我们可以将有用的方法附加到该类型上供我们使用。

在 `league.go` 添加一下代码

```go
type League []Player

func (l League) Find(name string) *Player {
    for i, p := range l {
        if p.Name==name {
            return &l[i]
        }
    }
    return nil
}
```

现在如果任何有 `League` 的人都可以很容易找到给定的玩家。

更改我们的 `PlayerStore` 接口以返回 `League` 而不是 `[]Player`。试着重新运行测试，你会遇到编译问题，因为我们修改了接口。但是这很容易修复，只要将返回类型从 `[]Player` 改为 `League` 就行了。

这使我们可以简化 `FileSystemStore` 的方法。

```go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

    player := f.GetLeague().Find(name)

    if player != nil {
        return player.Wins
    }

    return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
    league := f.GetLeague()
    player := league.Find(name)

    if player != nil {
        player.Wins++
    }

    f.database.Seek(0, 0)
    json.NewEncoder(f.database).Encode(league)
}
```

这看起来好多了，我们可以在 `League` 中找到其他可以被重构的功能。

我们现在需要处理记录新玩家获胜的场景。

## 首先编写测试

```go
t.Run("store wins for existing players", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    store.RecordWin("Pepper")

    got := store.GetPlayerScore("Pepper")
    want := 1
    assertScoreEquals(t, got, want)
})
```

## 尝试并运行测试

```
=== RUN   TestFileSystemStore/store_wins_for_existing_players#01
    --- FAIL: TestFileSystemStore/store_wins_for_existing_players#01 (0.00s)
        FileSystemStore_test.go:86: got 0 want 1
```

## 编写足够的代码使测试通过

我们只需要处理查找返回 `nil` 的情况因为它找不到 `player`。

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
    league := f.GetLeague()
    player := league.Find(name)

    if player != nil {
        player.Wins++
    } else {
        league = append(league, Player{name, 1})
    }

    f.database.Seek(0, 0)
    json.NewEncoder(f.database).Encode(league)
}
```

效果看起来不错，因此我们现在可以在集成测试中使用我们的新的 `Store`。这将使我们对软件的工作更有信心，然后我们可以删除冗余的 `InMemoryPlayerStore`。

在 `TestRecordingWinsAndRetrievingThem` 中，替换之前的记录。

```go
database, cleanDatabase := createTempFile(t, "")
defer cleanDatabase()
store := &FileSystemPlayerStore{database}
```

测试通过后就可以删除 `InMemoryPlayerStore` 了。`main.go` 现在会出现编译问题，这将促使我们现在在「真实」代码中使用我们的新存储。

```go
package main

import (
    "log"
    "net/http"
    "os"
)

const dbFileName = "game.db.json"

func main() {
    db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

    if err != nil {
        log.Fatalf("problem opening %s %v", dbFileName, err)
    }

    store := &FileSystemPlayerStore{db}
    server := NewPlayerServer(store)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

- 我们创建了一个文件作为数据库。
- 第 2 个参数 `os.OpenFile` 允许你定义打开文件的权限，在我们的例子中，`O_RDWR` 意味着我们想要读写权限，`os.O_CREATE` 是指如果文件不存在，则创建该文件。
- 第 3 个参数表示设置文件的权限，在我们的示例中，所有用户都可以读写文件。[（详情请参阅 superuser.com）](https://superuser.com/questions/295591/what-is-the-meaning-of-chmod-666)。

重启运行程序，现在将持久化数据到文件中。

## 更多的重构和性能问题

每当有人调用 `GetLeague()` 或 `GetPlayerScore()` 时，我们就从头读取该文件，并将其解析为 JSON。我们不应该这样做，因为 `FileSystemStore` 完全负责 league 的状态。我们只是希望在开始时使用该文件来获取当前状态，并在数据更改时更新它。

我们可以创建一个构造函数，该构造函数可以为我们执行一些初始化操作，并将 league 作为值存储在我们的 `FileSystemStore` 中，以便在读取中使用。

```go
type FileSystemPlayerStore struct {
    database io.ReadWriteSeeker
    league League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
    database.Seek(0, 0)
    league, _ := NewLeague(database)
    return &FileSystemPlayerStore{
        database:database,
        league:league,
    }
}
```

这样，我们只需从磁盘读取一次。我们现在可以替换以前的所有从磁盘上获得 league 的调用，并且只使用 `f.league`。

```go
func (f *FileSystemPlayerStore) GetLeague() League {
    return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

    player := f.league.Find(name)

    if player != nil {
        return player.Wins
    }

    return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
    player := f.league.Find(name)

    if player != nil {
        player.Wins++
    } else {
        f.league = append(f.league, Player{name, 1})
    }

    f.database.Seek(0, 0)
    json.NewEncoder(f.database).Encode(f.league)
}
```

运行测试将会提示初始化 `FileSystemPlayerStore`，因此只需通过调用我们新的构造函数来修复它们。

### 另一个问题

在我们处理文件的过程中有一些非常天真的行为，这可能会在以后产生非常严重的错误。

当我们 `Recordwin` 时，我们返回到文件的开头，然后写入新的数据，但是如果新的数据比之前的数据要小怎么办?

在我们目前的情况下，这是不可能的。我们从不编辑或删除得分，因此数据只会变得更大，但是这样的代码是不负责任的，出现删除场景的结果是不可想象的。

但是我们要怎么测试这种问题呢？我们需要做的是首先重构我们的代码，这样就可以将我们所编写的数据和正在写入的分开。然后我们可以分别测试它是否以我们期望的方式运行。

我们将创建一个新类型来封装我们的「当写入时，从头部开始」功能。我把它叫做 `Tape`。创建一个包含以下内容的新文件

```go
package main

import "io"

type tape struct {
    file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
    t.file.Seek(0, 0)
    return t.file.Write(p)
}
```

注意，我们现在只实现了 `Write`，因为它封装了 `Seek` 部分。这意味着我们的 `FileSystemStore` 可以只具有对 `Writer` 的引用。

```go
type FileSystemPlayerStore struct {
    database io.Writer
    league   League
}
```

更新构造函数以使用 `Tape`

```go
func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
    database.Seek(0, 0)
    league, _ := NewLeague(database)

    return &FileSystemPlayerStore{
        database: &tape{database},
        league:   league,
    }
}
```

最后，我们可以通过从 `RecordWin` 中删除 `Seek` 调用来获得我们想要的惊人回报。是的，这感觉并不多，但至少这意味着如果我们做任何其它类型的写入操作，我们可以依赖 `write` 来表达我们对它的需求。此外，它现在将允许我们分别测试可能存在问题的代码并修复它。

让我们编写一个测试，我们想用比原始内容更小的东西来更新文件的整个内容。在 `tape_test.go` 中：

## 首先编写测试

我们只需要创建一个文件，尝试用我们的 `tape` 来写，再读一遍，看看文件里有什么。

```go
func TestTape_Write(t *testing.T) {
    file, clean := createTempFile(t, "12345")
    defer clean()

    tape := &tape{file}

    tape.Write([]byte("abc"))

    file.Seek(0, 0)
    newFileContents, _ := ioutil.ReadAll(file)

    got := string(newFileContents)
    want := "abc"

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

## 尝试运行测试

```
=== RUN   TestTape_Write
--- FAIL: TestTape_Write (0.00s)
    tape_test.go:23: got 'abc45' want 'abc'
```

就像我们想的一样！它只写我们想要的数据，而不写其他数据。

## 编写足够的代码使测试通过

`os.File` 文件有一个 truncate 函数，可以让我们有效地清空文件。我们应该能够调用它来得到我们想要的功能。

修改 `tape` 为以下内容

```
type tape struct {
    file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
    t.file.Truncate(0)
    t.file.Seek(0, 0)
    return t.file.Write(p)
}
```

编译器会在许多我们期望一个 `io.ReadWriteSeeker` 类型但是我们传入 `*os.File` 的地方失败。你现在应该可以自己修复这些问题了，但是如果你遇到困难，请检查源代码。

一旦重构完成，我们 `TestTape_Write` 的测试就应该通过了！

### 一个另外的小重构

在 `RecordWin` 中，我们有行 `json.NewEncoder(f.database).Encode(f.league)`。

我们不需要在每次编写代码时创建一个新的编码器，我们可以在构造函数中初始化一个编码器并使用它。

在我们的类型中存储对编码器的引用。

```go
type FileSystemPlayerStore struct {
    database *json.Encoder
    league   League
}
```

在构造器中初始化它

```go
func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
    file.Seek(0, 0)
    league, _ := NewLeague(file)

    return &FileSystemPlayerStore{
        database: json.NewEncoder(&tape{file}),
        league:   league,
    }
}
```

在 `RecordWin` 中使用它。

## 刚刚我们不是打破了一些规则？测试私有的东西？没有接口？

### 测试私有的类型

的确，一般来说，你不应该测试私有的东西，因为这有时会导致你的测试与实现的耦合过于紧密，这可能会阻碍将来的重构。

然而，我们不能忘记测试应该给我们信心。

如果添加任何类型的编辑或删除功能，我们对这些实现是否能运行就没有信心了。我们不想留下这样的代码，特别是如果有不止一个人在处理这些代码，他们可能不知道我们最初的方法有什么缺点。

最后，这只是一个测试！如果我们决定改变它的工作方式，仅仅删除测试并不是什么灾难，但是我们至少实现了对未来维护者的要求。

### 接口

我们从使用 `io.Reader` 开始编写代码。因为那是对我们新的 `PlayerStore` 进行单元测试最简单的方法。当我们开发代码时，我们转而使用 `io.ReadWriter` 然后是 `io.ReadWriteSeeker`。然后我们发现，除了 `*os.File` 之外，标准库中没有任何实际实现的东西。我们本来决定编写自己的或者使用开源的库，但是仅仅为测试使用临时文件就显得很实用了。

最后我们需要 `Truncate`，它也在 `*os.File` 中。我们可以选择创建自己的接口实现这些需求。

```go
type ReadWriteSeekTruncate interface {
    io.ReadWriteSeeker
    Truncate(size int64) error
}
```

但这有什么好处呢？请记住，我们并不是在模拟，**文件系统**存储采取除 `*os.File` 之外的任何类型都是不现实的。所以我们不需要接口给我们的多态性。

不要害怕像我们这里所做的那样去改变类型和做新的实验。使用静态类型语言的好处是编译器可以帮助你完成每一个更改。

## 错误处理

在开始排序之前，我们应该确保对当前代码感到满意，并删除可能存在的任何技术债务。尽可能快地使用软件（脱离红色状态）是一个重要的原则，但这并不意味着我们应该忽略出错的场景！

如果我们回到 `FileSystemStore.go`。我们在构造函数中有 `league, _:= NewLeague(f.database)`。

如果 `NewLeague` 无法从我们提供的`io.Reader` 中解析 league，它会返回一个错误。

在我们测试失败的时候，忽略这一点是很实际的。如果我们同时处理它，我们将同时处理两件事。

如果我们的构造函数能够返回一个错误，我们就这样做。

```go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
    file.Seek(0, 0)
    league, err := NewLeague(file)

    if err != nil {
        return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
    }

    return &FileSystemPlayerStore{
        database:&tape{file},
        league:league,
    }, nil
}
```

请记住，提供有用的错误信息非常重要（就像你写的测试一样）。人们在网上开玩笑说大多数 Go 代码都是

```go
if err != nil {
    return err
}
```

**这绝对不是习惯用语。** 为你的错误添加上下文信息（例如你正在做什么导致的错误）使操作你的软件更加容易。

如果你尝试编译，将会得到一些错误。

```
./main.go:18:35: multiple-value NewFileSystemPlayerStore() in single-value context
./FileSystemStore_test.go:35:36: multiple-value NewFileSystemPlayerStore() in single-value context
./FileSystemStore_test.go:57:36: multiple-value NewFileSystemPlayerStore() in single-value context
./FileSystemStore_test.go:70:36: multiple-value NewFileSystemPlayerStore() in single-value context
./FileSystemStore_test.go:85:36: multiple-value NewFileSystemPlayerStore() in single-value context
./server_integration_test.go:12:35: multiple-value NewFileSystemPlayerStore() in single-value context
```

在 main 中，我们要退出程序，打印错误。

```go
store, err := NewFileSystemPlayerStore(db)

if err != nil {
    log.Fatalf("problem creating file system player store, %v ", err)
}
```

在测试中，我们应该断言没有错误。我们可以编写辅助函数来协助处理。

```go
func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("didnt expect an error but got one, %v", err)
    }
}
```

使用这个辅助函数处理其他编译问题。最后，你应该得到一个失败的测试。

```go
=== RUN   TestRecordingWinsAndRetrievingThem
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:14: didnt expect an error but got one, problem loading player store from file /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db841037437, problem parsing league, EOF
```

我们不能解析 league，因为文件是空的。我们以前没有出错，因为我们一直都忽略了它们。

通过将一些有效的 JSON 数据放入其中来修复我们的大型集成测试，然后我们可以为这个场景编写一个特定的测试。

```go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[]`)
    //etc...
```

现在所有的测试都通过了，我们需要处理文件为空的场景。

## 首先编写测试

```go
t.Run("works with an empty file", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, "")
    defer cleanDatabase()

    _, err := NewFileSystemPlayerStore(database)

    assertNoError(t, err)
})
```

## 尝试运行测试

```go
=== RUN   TestFileSystemStore/works_with_an_empty_file
    --- FAIL: TestFileSystemStore/works_with_an_empty_file (0.00s)
        FileSystemStore_test.go:108: didnt expect an error but got one, problem loading player store from file /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db019548018, problem parsing league, EOF
```

## 编写足够的代码使测试通过

将构造函数更改为以下内容

```go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

    file.Seek(0, 0)

    info, err := file.Stat()

    if err != nil {
        return nil, fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
    }

    if info.Size()==0 {
        file.Write([]byte("[]"))
        file.Seek(0, 0)
    }

    league, err := NewLeague(file)

    if err != nil {
        return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
    }

    return &FileSystemPlayerStore{
        database:&tape{file},
        league:league,
    }, nil
}
```

`file.Stat` 返回我们的文件的统计数据。我们可以检查文件的大小，如果它是空的，我们就会编写一个空的 JSON 数组，然后 `Seek` 到开始位置，为剩下的代码做准备。

## 重构

我们的构造函数现在有点混乱，我们可以将初始化代码提取到函数中

```go
func initialisePlayerDBFile(file *os.File) error {
    file.Seek(0, 0)

    info, err := file.Stat()

    if err != nil {
        return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
    }

    if info.Size()==0 {
        file.Write([]byte("[]"))
        file.Seek(0, 0)
    }

    return nil
}
```

```
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

    err := initialisePlayerDBFile(file)

    if err != nil {
        return nil, fmt.Errorf("problem initialising player db file, %v", err)
    }

    league, err := NewLeague(file)

    if err != nil {
        return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
    }

    return &FileSystemPlayerStore{
        database:&tape{file},
        league:league,
    }, nil
}
```

## 排序

我们的产品负责人想让 `/league` 返回按得分排序的玩家。

这里主要要做的决定是，在软件的什么位置处理这个问题。如果我们使用的是「真实的」数据库，我们会使用像 `ORDER BY` 这样的东西，所以排序非常快，所以出于这个原因，应该由 `PlayerStore` 的实现负责。

## 首先编写测试

我们可以在 `TestFileSystemStore` 中的第一个测试上更新断言

```go
t.Run("league sorted", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    got, err := store.GetLeague()
    assertNoError(t, err)

    want := []Player{
        {"Chris", 33},
        {"Cleo", 10},
    }

    assertLeague(t, got, want)

    // read again
    got, err = store.GetLeague()
    assertNoError(t, err)
    assertLeague(t, got, want)
})
```

JSON 输入的顺序是错误的，我们的 `want` 将检查它是否以正确的顺序返回给调用者。

## 尝试运行测试

```
=== RUN   TestFileSystemStore/league_from_a_reader,_sorted
    --- FAIL: TestFileSystemStore/league_from_a_reader,_sorted (0.00s)
        FileSystemStore_test.go:46: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
        FileSystemStore_test.go:51: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
```

## 编写足够的代码使测试通过

```go
func (f *FileSystemPlayerStore) GetLeague() League {
    sort.Slice(f.league, func(i, j int) bool {
        return f.league[i].Wins > f.league[j].Wins
    })
    return f.league
}
```

[`sort.Slice`](https://golang.org/pkg/sort/#Slice)

> 根据给定的比较函数，Slice 对提供的切片进行排序。

真的很简单！

## 总结

### 讨论的内容

- `Seeker` 接口以及它与 `Reader` 和 `Writer` 的关系。
- 处理文件读写。
- 为测试创建辅助函数，隐藏文件中所有杂乱的内容。
- 使用 `sort.Slice` 对切片排序。
- 利用编译器帮助我们安全地对应用程序进行结构更改。

### 打破规则

- 软件工程中的大多数规则并非铁律，只是 80% 的时间在工作中都是最佳实践。
- 当我们发现以前不测试内部函数的「规则」对我们没有帮助时，我们就打破了这个规则。
- 当打破规则时，了解你所做的权衡是很重要的。在我们的例子中这样做没有问题，因为它只是一个测试，如果不这样做的话，将很难执行这个场景。
- 为了能够打破规则，**你首先必须理解它们**。可以跟学习吉他做类比，不管你认为自己多有创意，你都必须理解和练习基础。

### 软件的功能

- 我们创建了一个 HTTP API，你可以利用它创建玩家并增加他们的得分。
- 我们可以将包含每个人得分的联盟数据作为 JSON 返回。
- 数据以 JSON 文件的形式存储。

---

作者：[Chris James](https://dev.to/quii)
译者：[Donng](https://github.com/Donng)
校对：[pityonline](https://github.com/pityonline)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
