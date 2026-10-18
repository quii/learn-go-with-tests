---
description: コマンドライン、パッケージ構造
---

# コマンドライン、パッケージ構造

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/command-line)

プロダクトオーナーは、2つ目のアプリケーション（コマンドラインアプリケーション）を導入することで _ピボット_ したいと考えています。

今のところは、ユーザーが`Ruth wins`と入力したときに、プレーヤーの勝利を記録できればよいだけです。最終的には、ユーザーがポーカーをプレイするのを手助けするツールにすることを目指しています。

プロダクトオーナーは、2つのアプリケーション間でデータベースを共有して、新しいアプリケーションに記録された勝利に応じてリーグが更新されるようにしたいと考えています。

## コードのおさらい

HTTPサーバーを起動する`main.go`ファイルを含むアプリケーションがあります。今回の演習ではHTTPサーバーそのものには関心がありませんが、それが使っている抽象化には関心があります。それは`PlayerStore`に依存しています。

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}
```

前の章では、そのインターフェースを実装する`FileSystemPlayerStore`を作りました。これを新しいアプリケーションで再利用できるはずです。

## まずはプロジェクトのリファクタリングから

このプロジェクトでは、既存のWebサーバーとコマンドラインアプリの2つのバイナリを作成する必要があります。

新しい作業に取り掛かる前に、これに対応できるようにプロジェクトを構成する必要があります。

これまでのところ、すべてのコードは次のようなパスの1つのフォルダーにあります

`$GOPATH/src/github.com/your-name/my-app`

Goでアプリケーションを作成するには、`package main`内に`main`関数が必要です。これまでのところ、「ドメイン」コードはすべて`package main`内にあり、`func main`はすべてを参照できます。

これまではこれで問題ありませんでしたし、パッケージ構造をやり過ぎないのは良い習慣です。時間をかけて標準ライブラリを眺めてみると、たくさんのフォルダや構造といったものはほとんど見当たらないことがわかるでしょう。

ありがたいことに、 _必要になったときに_ 構造を追加するのは非常に簡単です。

既存のプロジェクト内に`cmd`ディレクトリを作成し、その中に`webserver`ディレクトリを作成します（例：`mkdir -p cmd/webserver`）。

`cmd`は、プロジェクトがビルドするアプリケーションの`main`パッケージを置くための、Goで広く使われている慣習です。これにより、それらをプロジェクトのルートにあるインポート可能なライブラリコードから分離しておけます。

そこに`main.go`を移動します。

`tree`がインストールされている場合は実行してみてください。構造は次のようになっているはずです。

```
.
|-- file_system_store.go
|-- file_system_store_test.go
|-- cmd
|   |-- webserver
|       |-- main.go
|-- league.go
|-- server.go
|-- server_integration_test.go
|-- server_test.go
|-- tape.go
|-- tape_test.go
```

これで、アプリケーションとライブラリコードが実質的に分離されましたが、一部のパッケージ名を変更する必要があります。 Goアプリケーションをビルドするときは、そのパッケージは`main`で _なければならない_ ことを思い出してください。

それ以外のすべてのコードのパッケージ名を`poker`に変更します。

最後に、このパッケージを`main.go`にインポートして、それを使用してWebサーバーを作成できるようにする必要があります。そうすれば、`poker.FunctionName`のようにしてライブラリコードを使えます。

パスはコンピューターによって異なりますが、次のようになります。

```go
// cmd/webserver/main.go
package main

import (
	"github.com/quii/learn-go-with-tests/command-line/v1"
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

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := poker.NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", server))
}
```

`dbFileName`は相対パスなので、`game.db.json`はバイナリが置かれているディレクトリではなく、ビルドしたバイナリを _実行した_ ディレクトリを基準に作成（または読み込み）されます。プロダクトオーナーはCLIとWebサーバーで同じデータベースを共有したいと考えているので、これは重要です。両者を異なるディレクトリから実行すると、それぞれが別々の`game.db.json`を持つことになり、しかもそのことに気付けません。これについては後述の「最終チェック」で改めて触れます。

フルパスは少し違和感があるかもしれませんが、これが、公開されている _あらゆる_ ライブラリをコードにインポートする方法なのです。

ドメインコードを個別のパッケージに分離し、それをGitHubのようなパブリックリポジトリにコミットすることで、どのGo開発者でも、そのパッケージをインポートして、私たちが書いた機能を利用する独自のコードを書けるようになります。初めて実行しようとすると、パッケージが存在しないというエラーになりますが、`go get`を実行するだけで済みます。

さらに、ユーザーは[`pkg.go.dev`でドキュメント](https://pkg.go.dev/github.com/quii/learn-go-with-tests/command-line/v1)を表示できます。

### 最終チェック

* ルートで`go test`を実行し、テストがまだ成功することを確認します
* `cmd/webserver`内に移動し、`go run main.go`を実行します
  * `http://localhost:5000/league` にアクセスすると、まだ機能していることがわかります

この章の後半では、Webサーバーと同じ`game.db.json`を共有することを想定した2つ目のアプリケーション`cmd/cli`を作成します。`dbFileName`は現在の作業ディレクトリを基準に解決されるため、両方のバイナリがお互いの更新を参照できるようにするには、 _同じディレクトリから_ 実行する必要があります。たとえば、各`cmd`サブフォルダ内で`go run main.go`を使うのではなく、先にビルドしてからプロジェクトのルートでバイナリを実行します（`go build -o webserver ./cmd/webserver && go build -o cli ./cmd/cli`の後に、`./webserver`と`./cli`を実行）。

### ウォーキングスケルトン

テストの作成に取り掛かる前に、プロジェクトで構築する新しいアプリケーションを追加しましょう。`cmd`内に`cli`（コマンドラインインターフェース）という別のディレクトリを作成し、次のように`main.go`を追加します

```go
// cmd/cli/main.go
package main

import "fmt"

func main() {
	fmt.Println("Let's play poker")
}
```

私たちが取り組む最初の要件は、ユーザーが`{PlayerName} wins`と入力したときに勝利を記録することです。

## 最初にテストを書く

`CLI`と呼ばれるものを作ってポーカーを`Play`できるようにする必要があることはわかっています。ユーザー入力を読み取り、勝利を`PlayerStore`に記録する必要があります。

ただし、先走りすぎる前に、まずは`PlayerStore`と期待どおりに連携できるかを確認するテストを書いてみましょう。

`CLI_test.go`内（`cmd`内ではなくプロジェクトのルート内）

```go
// CLI_test.go
package poker

import "testing"

func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	cli := &CLI{playerStore}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}
}
```

* 他のテストの`StubPlayerStore`を使用できます
* 依存関係をまだ存在していない`CLI`型に渡します
* 未作成の`PlayPoker`メソッドでゲームをトリガーします
* 勝利が記録されていることを確認する

## テストを実行してみます

```
# github.com/quii/learn-go-with-tests/command-line/v2
./cli_test.go:25:10: undefined: CLI
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

この時点で、依存関係用のフィールドを持つ新しい`CLI`構造体を作成し、メソッドを追加するくらいは、もう慣れたものでしょう。

最終的に次のようなコードになるはずです。

```go
// CLI.go
package poker

type CLI struct {
	playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {}
```

今はテストを実行できるようにして、期待どおりにテストが失敗することを確認しようとしているだけだということを忘れないでください。

```
--- FAIL: TestCLI (0.00s)
    cli_test.go:30: expected a win call but didn't get any
FAIL
```

## 成功させるのに十分なコードを書く

```go
//CLI.go
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
```

これでテストが成功するはずです。

次に、特定のプレーヤーの勝利を記録できるように、`Stdin`（ユーザーからの入力）からの読み取りをシミュレートする必要があります。

これを確かめられるようにテストを拡張しましょう。

## 最初にテストを書く

```go
//CLI_test.go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}

	got := playerStore.winCalls[0]
	want := "Chris"

	if got != want {
		t.Errorf("didn't record correct winner, got %q, want %q", got, want)
	}
}
```

`os.Stdin`は、ユーザーの入力をキャプチャするために`main`で使用するものです。これは内部的には`*File`であり、つまり`io.Reader`を実装しています。ご存じのとおり、`io.Reader`はテキストを取り込むのに便利な手段です。

テストでは、便利な`strings.NewReader`を使用して`io.Reader`を作成し、ユーザーが入力すると想定される内容を詰めておきます。

## テストを実行してみます

`./CLI_test.go:12:32: too many values in struct initializer`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

新しい依存関係を`CLI`に追加する必要があります。

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}
```

```
--- FAIL: TestCLI (0.00s)
    CLI_test.go:23: didn't record the correct winner, got 'Cleo', want 'Chris'
FAIL
```

## 成功させるのに十分なコードを書く

まずは厳密に最も簡単なことをするのを忘れないでください。

```go
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Chris")
}
```

テストに成功しました。次に、実際のコードを書くように強制する別のテストを追加しますが、最初にリファクタリングしましょう。

## リファクタリング♪

`server_test`では、ここでのように勝利が記録されているかどうかを以前に確認しました。そのアサーションをヘルパーにまとめてDRYにしましょう

```go
//server_test.go
func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}
```

ここで、`server_test.go`と`CLI_test.go`の両方のアサーションを置き換えます。

テストは次のようになります。

```go
//CLI_test.go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	assertPlayerWin(t, playerStore, "Chris")
}
```

次に、異なるユーザー入力を使って _別の_ テストを書き、実際に入力を読み取らざるを得ないようにします。

## 最初にテストを書く

```go
//CLI_test.go
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

## テストを実行してみます

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

## 成功させるのに十分なコードを書く

[`bufio.Scanner`](https://golang.org/pkg/bufio/)を使用して、`io.Reader`から入力を読み取ります。

> パッケージbufioはバッファI/Oを実装しています。io.Readerまたはio.Writerオブジェクトをラップし、別のオブジェクト（ReaderまたはWriter）を作成します。このオブジェクトも同じインターフェースを実装しますが、バッファリングとテキストI/Oのための補助機能を提供します。

コードを次のように更新します。

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
```

テストに成功します。

* `Scanner.Scan()`は改行まで読み込みます。
* 次に、`Scanner.Text()`を使用して、スキャナーが読み取った`string`を返します。

成功するテストがいくつかできたので、これを`main`に組み込みましょう。できるだけ早く、完全に統合された動作するソフトウェアを手に入れるよう常に努めるべきだということを忘れないでください。

`main.go`に以下を追加して実行します。（2番目の依存関係のパスをコンピューター上のものと一致するように調整する必要がある場合があります）

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

エラーが発生するはずです。

```
command-line/v3/cmd/cli/main.go:32:25: implicit assignment of unexported field 'playerStore' in poker.CLI literal
command-line/v3/cmd/cli/main.go:32:34: implicit assignment of unexported field 'in' in poker.CLI literal
```

ここで起こっているのは、`CLI`のフィールド`playerStore`および`in`に割り当てようとしているためです。これらは、エクスポートされていない（private）フィールドです。テストは`CLI`（`poker`）と同じパッケージにあるため、テストコードでこれを行うことができました。しかし、`main`はパッケージ`main`にあるため、アクセスできません。

これは _作業を統合すること_ の重要性を浮き彫りにしています。私たちは`CLI`の依存関係を（`CLI`のユーザーに公開したくないので）正当な理由でプライベートにしましたが、ユーザーが`CLI`を構築する方法を用意していませんでした。

この問題を早期に発見する方法はありますか？

### `package mypackage_test`

これまでの他のすべての例では、テストファイルを作成するときに、テストしているのと同じパッケージ内にあることを宣言しています。

これはこれで問題ありませんし、まれにパッケージの内部にあるものをテストしたくなったときにも、エクスポートされていない型にアクセスできます。

しかし、 _一般的には_ 内部的なものをテスト _しない_ ことを推奨してきたのですから、Goでそれを強制できないでしょうか？（`main`のように）エクスポートされた型にしかアクセスできない状態でコードをテストできたらどうでしょう？

複数のパッケージを含むプロジェクトを作成している場合、テストパッケージ名の最後に`_test`を付けることを強くお勧めします。こうすると、パッケージ内のパブリックな型にしかアクセスできなくなります。今回のケースに役立つだけでなく、パブリックAPIのみをテストするという規律を守らせるのにも役立ちます。それでも内部をテストしたい場合は、テストしたいパッケージと同じパッケージ名で別のテストを作成できます。

TDDの格言に、「コードをテストできないなら、おそらくコードのユーザーにとってもそれを組み込むのは難しい」というものがあります。`package foo_test`を使うと、パッケージのユーザーと同じようにコードをインポートしてテストすることを強いられるので、この点で役立ちます。

`main`を修正する前に、`CLI_test.go`内のテストのパッケージを`poker_test`に変更しましょう。

適切に設定されたIDEを使っていれば、突然真っ赤になるはずです！
コンパイラを実行すると、次のエラーが発生します

```
./CLI_test.go:12:19: undefined: StubPlayerStore
./CLI_test.go:17:3: undefined: assertPlayerWin
./CLI_test.go:22:19: undefined: StubPlayerStore
./CLI_test.go:27:3: undefined: assertPlayerWin
```

ここで、パッケージ設計に関するさらなる疑問に突き当たりました。ソフトウェアをテストするために、エクスポートされていないスタブとヘルパー関数を作成しました。ヘルパーは、`poker`パッケージの`_test.go`ファイルで定義されているため、`CLI_test`では使用できなくなりました。

#### スタブとヘルパーを「公開」したいですか？

これは主観的な議論です。テストを容易にするためのコードでパッケージのAPIを汚染したくない、と主張する人もいるでしょう。

Mitchell Hashimotoによるプレゼンテーション[Goを使った高度なテスト](https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=53)では、HashiCorpではこれを行うことを推奨していると説明されています。
そうすれば、パッケージのユーザーはスタブを書くという車輪の再発明をせずにテストを書けるからです。私たちの場合、これは、`poker`パッケージを使用する誰もが、コードを操作したい場合、独自のスタブ`PlayerStore`を作成する必要がないことを意味します。

個人的な経験ですが、私は他の共有パッケージでもこの手法を使ったことがあり、ユーザーが私たちのパッケージと統合するときに時間を節約するという点で非常に有用であることがわかりました。

それでは、`testing.go`というファイルを作成して、スタブとヘルパーを追加しましょう。

```go
// testing.go
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

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}

// todo for you - the rest of the helpers
```

パッケージをインポートする側にヘルパーを公開したいなら、ヘルパーをパブリックにする必要があります（エクスポートするには先頭を大文字にすることを思い出してください）。

`CLI`テストでは、別のパッケージ内でコードを使用しているかのようにコードを呼び出す必要があります。

```go
//CLI_test.go
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

これで、`main`と同じ問題が発生することがわかります。

```
./CLI_test.go:15:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:15:39: implicit assignment of unexported field 'in' in poker.CLI literal
./CLI_test.go:25:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:25:39: implicit assignment of unexported field 'in' in poker.CLI literal
```

これを回避する最も簡単な方法は、他の型と同じようにコンストラクターを作成することです。また、リーダーは構築時に自動的にラップされるようになるので、`CLI`がリーダーの代わりに`bufio.Scanner`を保持するように変更します。

```go
//CLI.go
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
```

これにより、読み取りコードを簡素化してリファクタリングできます。

```go
//CLI.go
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

代わりにコンストラクタを使用するようにテストを変更すれば、テストは再び成功するはずです。

最後に、新しい`main.go`に戻り、先ほど作成したコンストラクタを使用できます。

```go
//cmd/cli/main.go
game := poker.NewCLI(store, os.Stdin)
```

実行して、"Bob wins"と入力してみてください。

### リファクタリング♪

それぞれのアプリケーションで、ファイルを開いてその内容から`file_system_store`を作成するという繰り返しがあります。
これはパッケージの設計の小さな弱点のように感じられるので、パスからファイルを開いて`PlayerStore`を返す処理をカプセル化する関数をパッケージに作成すべきでしょう。

```go
//file_system_store.go
func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, closeFunc, nil
}
```

次に、両方のアプリケーションをリファクタリングして、この関数を使用してストアを作成します。

#### CLIアプリケーションコード

```go
// cmd/cli/main.go
package main

import (
	"fmt"
	"github.com/quii/learn-go-with-tests/command-line/v3"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	poker.NewCLI(store, os.Stdin).PlayPoker()
}
```

#### Webサーバーアプリケーションコード

```go
// cmd/webserver/main.go
package main

import (
	"github.com/quii/learn-go-with-tests/command-line/v3"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

対称性に注目してください。

ユーザーインターフェースは異なるのに、セットアップはほぼ同じです。
これは、ここまでの設計が妥当であることの良い裏付けだと言えるでしょう。
また、`FileSystemPlayerStoreFromFile`はファイルを閉じるための関数を返すため、ストアの使用が終了したら、基になるファイルを閉じることができます。

## まとめ

### パッケージ構造

この章では、これまでに作成したドメインコードを再利用して、2つのアプリケーションを作成したいと考えました。
これを行うには、パッケージ構造を更新して、それぞれの`main`に個別のフォルダーを用意する必要がありました。

これを行うと、エクスポートされない値が原因で統合の問題が発生したため、小さな「スライス」で作業し、頻繁に統合することの価値をさらに示しています。

`mypackage_test`を使うと、コードを組み込む他のパッケージと同じ体験ができるテスト環境を作れることを学びました。これにより、統合の問題を発見したり、コードがどれほど扱いやすいか（あるいは扱いにくいか！）を確かめたりできます。

### ユーザー入力を読み取る

`os.Stdin`は、`io.Reader`を実装しているため、読み取りが非常に簡単です。
`bufio.Scanner`を使用して、ユーザー入力を1行ずつ簡単に読み取ることができました。

### 単純な抽象化により、コードの再利用が簡単になります

`PlayerStore`を新しいアプリケーションに統合するのはほとんど手間がかかりませんでした（パッケージの調整を行った後）。
スタブバージョンも公開することにしたため、テストも非常に簡単でした。
