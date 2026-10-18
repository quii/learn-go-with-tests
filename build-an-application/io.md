---
description: IO and sorting
---

# IO、並び替え

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/io)

[前の章](json.md) 新しいエンドポイント `/league`を追加することで、アプリケーションを繰り返し続けました。途中で、JSONの処理方法、埋め込みタイプ、ルーティングについて学習しました。

私たちの製品の所有者は、サーバーが再起動されたときにソフトウェアがスコアを失うことにいくらか混乱しています。これは、ストアの実装がインメモリであるためです。彼女はまた、「`/league`」エンドポイントが勝ちの数で順序付けられたプレーヤーを返す必要があると解釈しなかったことに不満を持っています！

## これまでのコード

```go
// server.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

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
// in_memory_player_store.go
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

```go
// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewPlayerServer(NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

章の上部にあるリンクで対応するテストを見つけることができます。

## データを保存する

これを使用できるデータベースは数十ありますが、ここでは非常にシンプルなアプローチを採用します。このアプリケーションのデータをJSONとしてファイルに保存します。

これにより、データの移植性が非常に高くなり、実装が比較的簡単になります。

特にうまくスケーリングしませんが、これはプロトタイプなので、今のところは問題ありません。私たちの状況が変化し、それが適切でなくなった場合、私たちが使用した`PlayerStore`抽象化のため、それを別のものと交換するのは簡単です。

新しいストアを開発するときに統合テストに合格し続けるために、今のところ`InMemoryPlayerStore`を保持します。新しい実装が統合テストに合格するのに十分であると確信したら、それを入れ替えてから、`InMemoryPlayerStore`を削除します。

## 最初にテストを書く

これで、データの読み取り（`io.Reader`）、データの書き込み（`io.Writer`）、および標準ライブラリを使用してこれらの関数をテストすることなく標準ライブラリをテストする方法について、標準ライブラリに関連するインターフェースに精通しているはずです。実際のファイルを使用する必要があります。

この作業を完了するには、 `PlayerStore`を実装する必要があるため、実装する必要のあるメソッドを呼び出すストアのテストを記述します。`GetLeague`から始めます。

```go
//file_system_store_test.go
func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)
	})
}
```

私たちは`Reader`を返す`strings.NewReader`を使用しています。
これは、`FileSystemPlayerStore`がデータを読み取るために使用するものです。`main`では、`Reader`でもあるファイルを開きます。

## テストを実行してみます

```
# github.com/quii/learn-go-with-tests/io/v1
./file_system_store_test.go:15:12: undefined: FileSystemPlayerStore
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

新しいファイルで`FileSystemPlayerStore`を定義しましょう

```go
//file_system_store.go
type FileSystemPlayerStore struct{}
```

再試行

```
# github.com/quii/learn-go-with-tests/io/v1
./file_system_store_test.go:15:28: too many values in struct initializer
./file_system_store_test.go:17:15: store.GetLeague undefined (type FileSystemPlayerStore has no field or method GetLeague)
```

`Reader`を渡していますが、予期しておらず、まだ`GetLeague`が定義されていないため、問題があります。

```go
//file_system_store.go
type FileSystemPlayerStore struct {
	database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	return nil
}
```

もう一回試してみる...

```
=== RUN   TestFileSystemStore//league_from_a_reader
    --- FAIL: TestFileSystemStore//league_from_a_reader (0.00s)
        file_system_store_test.go:24: got [] want [{Cleo 10} {Chris 33}]
```

## 成功させるのに十分なコードを書く

以前にリーダーからJSONを読み取りました

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetLeague() []Player {
	var league []Player
	json.NewDecoder(f.database).Decode(&league)
	return league
}
```

テストは成功するはずです。

## リファクタリング♪

これは以前に行ったことです！
サーバーのテストコードは、応答からJSONをデコードする必要がありました。

この関数をDRYにしてみましょう。

`league.go`と呼ばれる新しいファイルを作成し、これを中に入れます。

```go
//league.go
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}
```

これを実装と、`server_test.go`のテストヘルパー`getLeagueFromResponse`で呼び出します

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetLeague() []Player {
	league, _ := NewLeague(f.database)
	return league
}
```

構文解析エラーに対処するための戦略はまだありませんが、続けましょう。

### 問題を探す

私たちの実装に欠陥があります。まず、`io.Reader`の定義を思い出してみましょう。

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

私たちのファイルを使用すると、最後までバイト単位で読み取ることを想像できます。
もう一度「読み取り`Read`」しようとするとどうなりますか？

現在のテストの最後に以下を追加します。

```go
//file_system_store_test.go

// read again
got = store.GetLeague()
assertLeague(t, got, want)
```

これは成功させたいのですが、テストを実行すると成功しません。

問題は、`Reader`が最後に到達したため、これ以上読むものがないことです。最初に戻るように指示する方法が必要です。

[ReadSeeker](https://golang.org/pkg/io/#ReadSeeker)は、標準ライブラリにあるもう1つのインターフェースで、役立ちます。

```go
type ReadSeeker interface {
	Reader
	Seeker
}
```

埋め込みを覚えていますか？これは、`Reader`と[`Seeker`](https://golang.org/pkg/io/#Seeker)で構成されるインターフェースです。

```go
type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}
```

これはいいですね、代わりにこのインターフェイスを取るように`FileSystemPlayerStore`を変更できますか？

```go
//file_system_store.go
type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	f.database.Seek(0, io.SeekStart)
	league, _ := NewLeague(f.database)
	return league
}
```

テストを実行してみてください。テストは成功しました！
テストで使用した`strings.NewReader`は、`ReadSeeker`も実装しているため、他の変更を行う必要はありませんでした。

次に、`GetPlayerScore`を実装します。

## 最初にテストを書く

```go
//file_system_store_test.go
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

## テストを実行してみます

```
./file_system_store_test.go:38:15: store.GetPlayerScore undefined (type FileSystemPlayerStore has no field or method GetPlayerScore)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

テストをコンパイルするには、メソッドを新しい型に追加する必要があります。

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	return 0
}
```

これでコンパイルされ、テストは失敗します

```
=== RUN   TestFileSystemStore/get_player_score
    --- FAIL: TestFileSystemStore//get_player_score (0.00s)
        file_system_store_test.go:43: got 0 want 33
```

## 成功させるのに十分なコードを書く

リーグを繰り返してプレーヤーを見つけ、スコアを返すことができます。

```go
//file_system_store.go
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

## リファクタリング♪

何十ものテストヘルパーリファクタリングを確認したので、これを機能させるためにこれをあなたにお任せします。

```go
//file_system_store_test.go
t.Run("get player score", func(t *testing.T) {
	database := strings.NewReader(`[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)

	store := FileSystemPlayerStore{database}

	got := store.GetPlayerScore("Chris")
	want := 33
	assertScoreEquals(t, got, want)
})
```

最後に、`RecordWin`でスコアの記録を開始する必要があります。

## 最初にテストを書く

私たちのアプローチは、書き込みに対してかなり近視眼的です。ファイル内のJSONの「行」を1つだけ更新することは簡単にはできません。すべての書き込みで、データベースの _whole_ 新しい表現を保存する必要があります。

どうやって書くの？通常は`Writer`を使用しますが、すでに`ReadSeeker`があります。潜在的に2つの依存関係が存在する可能性がありますが、標準ライブラリにはすでに`ReadWriteSeeker`用のインターフェースがあり、ファイルで必要なすべてのことを実行できます。

タイプを更新しましょう

```go
//file_system_store.go
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}
```

コンパイルされるかどうかを確認します

```
./file_system_store_test.go:15:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
./file_system_store_test.go:36:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
```

`strings.Reader`が` ReadWriteSeeker`を実装しないことはそれほど驚くべきことではないので、何をしますか？

2つの選択肢があります

* テストごとに一時ファイルを作成します。`*os.File`は`ReadWriteSeeker`を実装します。これの長所は、それが統合テストになり、実際にファイルシステムからの読み取りと書き込みを行っているため、非常に高い信頼性が得られることです。短所は、ユニットテストの方が速く、一般にシンプルであるためです。また、一時ファイルを作成し、テスト後にそれらが確実に削除されるように、さらに作業を行う必要があります。
* サードパーティのライブラリを使用できます。[Mattetti](https://github.com/mattetti)は、必要なインターフェイスを実装し、ファイルシステムに触れないライブラリ[filebuffer](https://github.com/mattetti/filebuffer)を作成しました。

ここでは特に間違った答えはないと思いますが、サードパーティのライブラリを使用することを選択することで、依存関係の管理について説明する必要があります。そのため、代わりにファイルを使用します。

テストを追加する前に、`strings.Reader`を`os.File`で置き換えることにより、他のテストをコンパイルする必要があります。

内部にいくつかのデータを含む一時ファイルを作成し、スコアのテストを抽象化するヘルパー関数をいくつか作成しましょう

```go
//file_system_store_test.go
func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

[CreateTemp](https://pkg.go.dev/os#CreateTemp)は、使用する一時ファイルを作成します。渡した `"db"`値は、作成するランダムなファイル名に付けられるプレフィックスです。これは、誤って他のファイルと衝突しないようにするためです。

`ReadWriteSeeker`（ファイル）だけでなく、関数も返すことに気づくでしょう。
テストが終了したら、ファイルを確実に削除する必要があります。エラーが発生しやすく、読者にとっても面白いものではないため、ファイルの詳細をテストに漏らしたくありません。`removeFile`関数を返すことで、ヘルパーの詳細を処理でき、呼び出し側が実行する必要があるのは`defer removeFile()`を実行することだけです。

```go
//file_system_store_test.go
func TestFileSystemStore(t *testing.T) {

	t.Run("league from a reader", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer removeFile()

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
		database, removeFile := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer removeFile()

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Chris")
		want := 33
		assertScoreEquals(t, got, want)
	})
}
```

テストを実行すると、テストに成功するはずです。かなりの量の変更がありましたが、インターフェースの定義が完了したように感じられ、これから新しいテストを簡単に追加できるはずです。

既存のプレイヤーの勝利を記録する最初の反復を取得しましょう

```go
//file_system_store_test.go
t.Run("store wins for existing players", func(t *testing.T) {
	database, removeFile := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer removeFile()

	store := FileSystemPlayerStore{database}

	store.RecordWin("Chris")

	got := store.GetPlayerScore("Chris")
	want := 34
	assertScoreEquals(t, got, want)
})
```

## テストを実行してみます

`./file_system_store_test.go:67:8: store.RecordWin undefined (type FileSystemPlayerStore has no field or method RecordWin)`


## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

新しいメソッドを追加する

```go
//file_system_store.go
func (f *FileSystemPlayerStore) RecordWin(name string) {

}
```

```
=== RUN   TestFileSystemStore/store_wins_for_existing_players
    --- FAIL: TestFileSystemStore/store_wins_for_existing_players (0.00s)
        file_system_store_test.go:71: got 33 want 34
```

私たちの実装は空なので、古いスコアが返されます。

## 成功させるのに十分なコードを書く

```go
//file_system_store.go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()

	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
		}
	}

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(league)
}
```

なぜ私が「`player.Wins++`」ではなく「`league[i].Wins++`」をしているのか、疑問に思われるかもしれません。

スライス上で「範囲`range`」を指定すると、ループの現在のインデックス（この場合は`i`）とそのインデックスにある要素の _copy_ が返されます。コピーの`Wins`値を変更しても、繰り返し処理する`league`スライスには影響しません。そのため、`league[i]`を実行して実際の値への参照を取得し、代わりにその値を変更する必要があります。

テストを実行すると、テストは成功するはずです。

## リファクタリング♪

`GetPlayerScore`と`RecordWin`では、名前でプレーヤーを見つけるために`[]Player`を繰り返し処理しています。

`FileSystemStore`の内部でこの共通コードをリファクタリングすることもできますが、私には、これが新しいタイプに引き上げることができるおそらく有用なコードであると感じています。これまで「リーグ`"League"`」での作業は常に`[]Player`で行っていましたが、`League`という新しいタイプを作成できます。これは、他の開発者が理解しやすくなり、そのタイプに便利なメソッドをアタッチして使用できるようになります。

`league.go`内に以下を追加します

```go
//league.go
type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
```

誰かが`League`を持っている場合、彼らは与えられたプレイヤーを簡単に見つけることができます。

`PlayerStore`インターフェイスを変更して、`[]Player`ではなく`League`を返すようにします。テストを再実行してみてください。インターフェイスを変更したためコンパイルの問題が発生しますが、修正は非常に簡単です。戻り値の型を `[]Player`から`League`に変更するだけです。

これにより、`file_system_store`のメソッドを簡略化できます。

```go
//file_system_store.go
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

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(league)
}
```

これはかなり見栄えがよくなりました。`League`の周りでリファクタリングできる他の便利な機能も見つけられそうです。

新しいプレーヤーの勝利を記録するシナリオを処理する必要があります。

## 最初にテストを書く

```go
//file_system_store_test.go
t.Run("store wins for new players", func(t *testing.T) {
	database, removeFile := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer removeFile()

	store := FileSystemPlayerStore{database}

	store.RecordWin("Pepper")

	got := store.GetPlayerScore("Pepper")
	want := 1
	assertScoreEquals(t, got, want)
})
```

## テストを実行してみます

```
=== RUN   TestFileSystemStore/store_wins_for_new_players#01
    --- FAIL: TestFileSystemStore/store_wins_for_new_players#01 (0.00s)
        file_system_store_test.go:86: got 0 want 1
```

## 成功させるのに十分なコードを書く

プレーヤーが見つからなかったために`Find`が`nil`を返すシナリオを処理する必要があるだけです。

```go
//file_system_store.go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		league = append(league, Player{name, 1})
	}

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(league)
}
```

ハッピーパスは問題なく見えたので、統合テストで新しい`Store`を使用してみることができます。これにより、ソフトウェアが動作するという確信が高まり、冗長な`InMemoryPlayerStore`を削除できます。

`TestRecordingWinsAndRetrievingThem`で古いストアを置き換えます。

```go
//server_integration_test.go
database, removeFile := createTempFile(t, "")
defer removeFile()
store := &FileSystemPlayerStore{database}
```

テストを実行すると、テストに成功し、`InMemoryPlayerStore`を削除できます。
`main.go`にコンパイルの問題が発生し、「実際の」コードで新しいストアを使用するようになります。

```go
// main.go
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

* データベース用のファイルを作成します。
* `os.OpenFile`の2番目の引数では、ファイルを開くための権限を定義できます。この場合、`O_RDWR`は読み取りと書き込みを行うことを意味し、`os.O_CREATE`はファイルが存在しない場合にファイルを作成することを意味します。
* 3番目の引数は、ファイルのアクセス権を設定することを意味します。この場合、すべてのユーザーがファイルの読み取りと書き込みを行うことができます。[（詳細な説明は`superuser.com`を参照）](https://superuser.com/questions/295591/what-is-the-meaning-of-chmod-666).

プログラムを実行すると、再起動の間にデータがファイルに永続化されるようになりました。

## より多くのリファクタリングとパフォーマンスの懸念

誰かが`GetLeague()`または`GetPlayerScore()`を呼び出すたびに、ファイル全体を読み取ってJSONに解析します。`FileSystemStore`はリーグの状態に完全に責任があるため、これを行う必要はありません。プログラムの起動時にファイルを読み取るだけで、データが変更されたときにファイルを更新するだけです。

この初期化の一部を実行できるコンストラクターを作成し、代わりに読み取りで使用するためにリーグを`FileSystemStore`の値として保存できます。

```go
//file_system_store.go
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league   League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database: database,
		league:   league,
	}
}
```

この方法では、ディスクから一度だけ読み取る必要があります。ディスクからリーグを取得するための以前の呼び出しをすべて置き換えて、代わりに`f.league`を使用できます。

```go
//file_system_store.go
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

	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(f.league)
}
```

テストを実行しようとすると、`FileSystemPlayerStore`の初期化について不平を言うので、新しいコンストラクターを呼び出して修正するだけです。

### 別の問題

非常に厄介なバグを発生させる可能性があるファイルを処理する方法には、いくつかの単純な点があります。

`RecordWin`を実行すると、ファイルの先頭に`Seek`して新しいデータを書き込みますが、新しいデータが以前のデータよりも小さい場合はどうなるでしょうか。

私たちの現在のケースでは、これは不可能です。スコアを編集または削除することはないため、データが大きくなるだけです。ただし、コードをこのようにしておくのは無責任です。削除シナリオが発生することは考えられないことではありません。

しかし、これをどのようにテストしますか？最初にコードをリファクタリングする必要があるので、書き込むデータの種類の懸念を書き込みから分離します。次に、それを個別にテストして、期待どおりに動作することを確認できます。

新しいタイプを作成して、「最初から書き始める」機能をカプセル化します。これを「テープ`Tape`」と呼びます。以下を使用して新しいファイルを作成します。

```go
// tape.go
package main

import "io"

type tape struct {
	file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}
```

ここでは、`Seek`部分をカプセル化しているため、ここでは`Write`のみを実装していることに注意してください。つまり、`FileSystemStore`は、代わりに`Writer`への参照のみを持つことができます。

```go
//file_system_store.go
type FileSystemPlayerStore struct {
	database io.Writer
	league   League
}
```

`Tape`を使用するようにコンストラクタを更新します

```go
//file_system_store.go
func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, _ := NewLeague(database)

	return &FileSystemPlayerStore{
		database: &tape{database},
		league:   league,
	}
}
```

最後に、`RecordWin`から`Seek`呼び出しを削除することで、私たちが望んでいた驚くべき見返りを得ることができます。はい、あまり感じませんが、少なくとも他の種類の書き込みを行う場合、`Write`を使用して必要な動作を実行できることを意味します。さらに、潜在的に問題のあるコードを個別にテストして修正できるようになります。

ファイルの内容全体を元の内容よりも小さいもので更新するテストを書いてみましょう。

## 最初にテストを書く

テストでは、コンテンツを含むファイルを作成し、`tape`を使用してファイルに書き込み、もう一度読み取って、ファイルの内容を確認します。`tape_test.go`では次のようにします。

```go
//tape_test.go
func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, io.SeekStart)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

## テストを実行してみます

```
=== RUN   TestTape_Write
--- FAIL: TestTape_Write (0.00s)
    tape_test.go:23: got 'abc45' want 'abc'
```

思った通り！必要なデータを書き込みますが、元のデータの残りが残ったままです。

## 成功させるのに十分なコードを書く

`os.File`には、ファイルを効率的に空にできる切り捨て関数があります。これを呼び出して、必要なものを取得できます。

`tape`を次のように変更します。

```go
//tape.go
type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}
```

コンパイラーは、`io.ReadWriteSeeker`を想定している多くの場所で失敗しますが、`*os.File`で送信しています。これらの問題は自分で修正できるはずですが、行き詰まった場合はソースコードを確認してください。

それを取得したら、`TestTape_Write`テストはパスするはずです！

### もう1つの小さなリファクタリング

`RecordWin`には、`json.NewEncoder(f.database).Encode(f.league)`という行があります。

書き込むたびに新しいエンコーダを作成する必要はありません。コンストラクタでエンコーダを初期化して、代わりに使用できます。

タイプに「エンコーダー`Encoder`」への参照を保存し、コンストラクターで初期化します。

```go
//file_system_store.go
type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
	file.Seek(0, io.SeekStart)
	league, _ := NewLeague(file)

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}
}
```

`tape.Write` は _毎回_ の呼び出しでファイルの先頭にシークするので、この組み合わせが本当に安全かどうか立ち止まって確認しておく価値があります。もし `Encode` が1回の呼び出しで `Write` を複数回呼び出すことがあれば、後続の書き込みのたびに先頭へシークし直して前の書き込みを上書きしてしまい、出力が壊れてしまいます。しかし実際にはそうはなりません。`Encoder.Encode` は常に値全体をまずメモリ上でマーシャルし、それを1回の `Write` 呼び出しで下層の `Writer` に渡します。そのため `tape` は `Encode` ごとに1回しかシークせず、これはまさに私たちが望む「常に先頭から上書きする」動作です。

`RecordWin`で使用します。

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)
}
```

## ルールを破っただけではありませんか？プライベートなものをテストしていますか？インターフェイスはありませんか？

### プライベート型のテストについて

一般的にプライベートなものをテストしないほうがよい場合があります。テストが実装に密に結合しすぎて、将来のリファクタリングを妨げる可能性があるためです。

ただし、テストによって _確信_ が得られることを忘れてはなりません。

なんらかの編集または削除機能を追加した場合、実装が機能するかどうか確信が持てませんでした。特に、最初のアプローチの欠点を認識していない複数の人が作業している場合は、コードをそのままにしたくありませんでした。

最後に、これは1つのテストにすぎません。動作方法を変更することを決定した場合、テストを削除するだけで障害になることはありませんが、少なくとも将来のメンテナの要件を把握しています。

### インターフェース

新しい`PlayerStore`を単体テストするための最も簡単な方法であったため、`io.Reader`を使用してコードを開始しました。コードを開発したとき、`io.ReadWriter`に移動し、次に`io.ReadWriteSeeker`に移動しました。次に、標準ライブラリには、`*os.File`以外に実際にそれを実装したものは何もないことがわかりました。独自に作成するか、オープンソースを使用するかを決定することもできましたが、テスト用の一時ファイルを作成するだけで実用的でした。

最後に、`*os.File`にもある`Truncate`が必要でした。これらの要件を取り込んだ独自のインターフェースを作成することはオプションでした。

```go
type ReadWriteSeekTruncate interface {
	io.ReadWriteSeeker
	Truncate(size int64) error
}
```

しかし、これは本当に私たちに何を与えているのでしょうか？私たちは _モックではない_ ことを覚えておいてください。**ファイルシステム**ストアが `*os.File`以外のタイプを取ることは非現実的であるため、インターフェイスが提供するポリモーフィズムは必要ありません。

ここにあるように、タイプを切り刻んで変更し、実験することを恐れないでください。静的に型付けされた言語を使用することの素晴らしい点は、コンパイラーがすべての変更を支援することです。

## エラー処理

並べ替えに取り掛かる前に、現在のコードに満足していることを確認し、技術的な負債をすべて取り除く必要があります。ソフトウェアをできるだけ早く（赤の状態から抜け出す）ことは重要な原則ですが、それはエラーのケースを無視する必要があるという意味ではありません。

`file_system_store.go`に戻ると、コンストラクターに`league, _ := NewLeague(file)`があります。

私たちが提供する`*os.File`からリーグを解析できない場合、`NewLeague`はエラーを返す可能性があります。

すでに失敗したテストがあったため、その時点でそれを無視するのは実用的でした。同時にそれに取り組んだとしたら、2つのことを同時に処理していたことになります。

コンストラクタがエラーを返すことができるようにしましょう。

```go
//file_system_store.go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	file.Seek(0, io.SeekStart)
	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}
```

（テストと同じように）役立つエラーメッセージを表示することが非常に重要であることを忘れないでください。インターネットの人々は冗談めかして、ほとんどのGoコードは次のとおりだと言っています

```go
if err != nil {
	return err
}
```

**それは100％慣用的ではありません** エラーメッセージにコンテキスト情報（つまり、エラーを発生させるために実行していたこと）を追加すると、ソフトウェアの操作がはるかに簡単になります。

コンパイルしようとすると、いくつかのエラーが発生します。

```
./main.go:18:35: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:35:36: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:57:36: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:70:36: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:85:36: multiple-value NewFileSystemPlayerStore() in single-value context
./server_integration_test.go:12:35: multiple-value NewFileSystemPlayerStore() in single-value context
```

メインでは、プログラムを終了してエラーを出力します。

```go
//main.go
store, err := NewFileSystemPlayerStore(db)

if err != nil {
	log.Fatalf("problem creating file system player store, %v ", err)
}
```

テストでは、エラーがないことを確認する必要があります。これを手伝ってくれるヘルパーを作ることができます。

```go
//file_system_store_test.go
func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
```

このヘルパーを使用して、他のコンパイル問題を処理します。
最後に、失敗するテストがあるはずです。

```
=== RUN   TestRecordingWinsAndRetrievingThem
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:14: didn't expect an error but got one, problem loading player store from file /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db841037437, problem parsing league, EOF
```

ファイルが空であるため、リーグを解析できません。以前は常にエラーを無視していたため、以前はエラーが発生していませんでした。

有効なJSONをいくつか入れて、大きな統合テストを修正しましょう。

```go
//server_integration_test.go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, removeFile := createTempFile(t, `[]`)
	//etc...
}
```

すべてのテストに合格したので、ファイルが空であるシナリオを処理する必要があります。

## 最初にテストを書く

```go
//file_system_store_test.go
t.Run("works with an empty file", func(t *testing.T) {
	database, removeFile := createTempFile(t, "")
	defer removeFile()

	_, err := NewFileSystemPlayerStore(database)

	assertNoError(t, err)
})
```

## テストを実行してみます

```
=== RUN   TestFileSystemStore/works_with_an_empty_file
    --- FAIL: TestFileSystemStore/works_with_an_empty_file (0.00s)
        file_system_store_test.go:108: didn't expect an error but got one, problem loading player store from file /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db019548018, problem parsing league, EOF
```

## 成功させるのに十分なコードを書く

コンストラクタを次のように変更します。

```go
//file_system_store.go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}
```

`file.Stat`はファイルの統計を返し、ファイルのサイズを確認できます。空の場合は、空のJSON配列を`Write`、`Seek`を最初に戻して、残りのコードの準備をします。

## リファクタリング♪

コンストラクターは少し面倒なので、初期化コードを関数に抽出してみましょう。

```go
//file_system_store.go
func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}
```

```go
//file_system_store.go
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
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}
```

## 並べ替え

私たちの製品所有者は、`/league`がプレーヤーをスコアの高い順に並べ替えることを望んでいます。

ここでの主な決定は、ソフトウェアでこれが発生する場所です。「実際の」データベースを使用している場合は、「`ORDER BY`」のようなものを使用するので、ソートは超高速です。そのため、`PlayerStore`の実装に責任があるように思えます。

## 最初にテストを書く

`TestFileSystemStore`の最初のテストでアサーションを更新できます。

```go
//file_system_store_test.go
t.Run("league sorted", func(t *testing.T) {
	database, removeFile := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer removeFile()

	store, err := NewFileSystemPlayerStore(database)

	assertNoError(t, err)

	got := store.GetLeague()

	want := League{
		{"Chris", 33},
		{"Cleo", 10},
	}

	assertLeague(t, got, want)

	// read again
	got = store.GetLeague()
	assertLeague(t, got, want)
})
```

入ってくるJSONの順序が間違っており、`want`が正しい順序で呼び出し元に返されることを確認します。

## テストを実行してみます

```
=== RUN   TestFileSystemStore/league_from_a_reader,_sorted
    --- FAIL: TestFileSystemStore/league_from_a_reader,_sorted (0.00s)
        file_system_store_test.go:46: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
        file_system_store_test.go:51: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
```

## 成功させるのに十分なコードを書く

```go
func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}
```

[`sort.Slice`](https://golang.org/pkg/sort/#Slice)

> Sliceは、提供されたless関数を指定して、提供されたスライスをソートします。

かんたん！

## まとめ

### 学んだこと

* `Seeker`インターフェースと、`Reader`および`Writer`との関係。
* ファイルの操作。
* すべての面倒なものを隠すファイルでテストするための使いやすいヘルパーを作成します。
* スライスをソートするための`sort.Slice`。
* コンパイラを使用して、アプリケーションの構造を安全に変更できるようにします。

### ルール違反

* ソフトウェアエンジニアリングのほとんどのルールは実際にはルールではなく、80％の時間で機能するベストプラクティスにすぎません。
* 内部関数をテストしないという以前の「ルール」の1つが役に立たないというシナリオを発見したため、ルールを破りました。
* ルールを破るときは、自分のトレードオフを理解することが重要です。私たちの場合、それは1つのテストにすぎず、それ以外の場合はシナリオを実行することが非常に困難であったため、問題はありませんでした。
* ルールを破ることができるようにするには、**最初にそれらを理解する必要があります**。例えはギターを学ぶことです。自分がどれほどクリエイティブであるかは関係ありません。基本を理解し、実践する必要があります。

### ソフトウェアのある場所

* プレーヤーを作成してスコアを増やすことができるHTTP APIがあります。
* 全員のスコアのリーグをJSONとして返すことができます。
* データはJSONファイルとして永続化されます。

