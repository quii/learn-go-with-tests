---
description: Dependency Injection
---

# 依存性注入

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/di)

これにはインターフェースの理解が必要になるため、構造体のセクションをすでに読んでいることが前提です。

プログラミングコミュニティには、依存性注入に関する誤解がたくさんあります。

このガイドでは、

* フレームワークは必要ありません
* デザインが複雑になりすぎない
* テストを容易にします
* 優れた汎用関数を作成できます

`hello-world`の章で行ったように、誰かに挨拶する関数を書きたいのですが、今回は _actual Printing_ をテストします。

要約すると、その関数は次のようになります。

```go
func Greet(name string) {
    fmt.Printf("Hello, %s", name)
}
```

しかし、これをどのようにテストできますか？
`fmt.Printf`を呼び出すと _stdout_ に出力されますが、テストフレームワークを使用してキャプチャするのはかなり困難です。

私たちがする必要があるのは、印刷の依存関係を**注入** （過ぎたるは及ばざるが如し）をできるようにすることです。
**関数は気にする必要はありません** _ **場所** _ **または** _ **方法** _ **印刷が行われるため、** _ **インターフェース**を受け入れる必要があります_ **具体的なタイプではありません。**

その場合は、実装を変更して印刷するように制御し、テストできるようにします。 
実際では、stdoutに書き込むものを注入します。

`fmt.Printf`のソースコードを見ると、フックする方法がわかります。

```go
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
    return Fprintf(os.Stdout, format, a...)
}
```

面白いですね！
内部では、 `Printf`は`os.Stdout`を渡して `Fprintf`を呼び出しているだけです。

`os.Stdout`とは正確に何ですか？
`Fprintf`は第1引数として何が渡されることを期待していますか？

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
    p := newPrinter()
    p.doPrintf(format, a)
    n, err = w.Write(p.buf)
    p.free()
    return
}
```

`io.Writer`

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

さらに多くのGoコードを書くと、このインターフェイスが「このデータをどこかに置く」ための優れた汎用インターフェイスであるため、多くのポップアップが表示されます。

つまり、私たちは最終的に `Writer`を使用して挨拶をどこかに送信していることを知っています。この既存の抽象化を使用して、コードをテスト可能にし、再利用可能にします。

## 最初にテストを書く

```go
func TestGreet(t *testing.T) {
    buffer := bytes.Buffer{}
    Greet(&buffer, "Chris")

    got := buffer.String()
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

`bytes`パッケージの`buffer`タイプは `Writer`インターフェースを実装しています。

テストでこれを使用して`Writer`として送信し、`Greet`を呼び出した後に何が書き込まれたかを確認できます。

## テストを試して実行する

テストはコンパイルされません

```text
./di_test.go:10:7: too many arguments in call to Greet
    have (*bytes.Buffer, string)
    want (string)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

コンパイラを読んで、問題を修正してください。

```go
func Greet(writer *bytes.Buffer, name string) {
    fmt.Printf("Hello, %s", name)
}
```

`Hello, Chris di_test.go:16: got '' want 'Hello, Chris'`

テストは失敗します。名前は出力されますが、標準出力になることに注意してください。

## 成功させるのに十分なコードを書く

テストでは、ライターを使用して挨拶をバッファに送信します。`fmt.Fprintf`は`fmt.Printf`に似ていますが、代わりに `Writer`を使用して文字列を送信しますが、`fmt.Printf`のデフォルトはstdoutです。

```go
func Greet(writer *bytes.Buffer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}
```

テストに合格しました。

## リファクタリング

以前のコンパイラーは、`bytes.Buffer`へのポインターを渡すように指示しました。これは技術的には正しいですが、あまり役に立ちません。

これを実証するために、`Greet`関数を標準出力に出力するGoアプリケーションに接続してみてください。

```go
func main() {
    Greet(os.Stdout, "Elodie")
}
```

`./di.go:14:7: cannot use os.Stdout (type *os.File) as type *bytes.Buffer in argument to Greet`

前に説明したように、`fmt.Fprintf`を使用すると、`os.Stdout`と `bytes.Buffer`の両方の実装がわかっている`io.Writer`を渡すことができます。

より汎用的なインターフェースを使用するようにコードを変更すると、テストとアプリケーションの両方で使用できるようになります。

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

## io.Writerの詳細

`io.Writer`を使用してデータを書き込むことができる他の場所は何ですか？
`Greet`関数はどれほど一般的な目的ですか？

### インターネット

以下を実行します

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

プログラムを実行し、[http://localhost:5000](http://localhost:5000)に移動します。グリーティング機能が使用されているのがわかります。

HTTPサーバーについては後の章で説明しますので、詳細についてはあまり気にしないでください。

HTTPハンドラーを作成すると、 `http.ResponseWriter`と、リクエストの作成に使用された` http.Request`が与えられます。サーバーを実装するときは、ライターを使用して応答を _write_ します。

`http.ResponseWriter`も`io.Writer`を実装していると思われるので、ハンドラー内で `Greet`関数を再利用できます。

## まとめ

最初のコードは、制御できない場所にデータを書き込んだため、簡単にテストできませんでした。

テストによって動機付けされたコードをリファクタリングして、制御できるようにしました。
データを、**依存関係を注入する**ことによって書き込まれ、次のことが可能になりました：

_Motivated by our tests_ we refactored the code so we could control _where_ the data was written by **injecting a dependency** which allowed us to:

* **コードをテストする**関数を簡単にテストできない場合は、通常、依存関係が関数またはグローバルな状態に組み込まれているためです。たとえば、ある種のサービス層で使用されているグローバルデータベース接続プールがある場合、テストが困難になる可能性が高く、実行が遅くなります。DIは、（インターフェイスを介して）データベースの依存関係を挿入するように動機付けし、テストで制御できるものでモックアウトできます。
* **懸念事項を分離**して、「データの移動先」と「生成方法」を分離します。メソッド/関数の責任が多すぎると感じた場合は、（データの生成、およびデータベースへの書き込み、HTTPリクエストの処理、およびドメインレベルのロジックの実行）おそらくDIが必要なツールになるでしょう。
* **コードをさまざまなコンテキストで再利用できるようにする**コードを使用できる最初の「新しい」コンテキストは、テスト内です。しかし、さらに誰かがあなたの関数で何か新しいことを試したい場合、彼らは彼ら自身の依存関係を注入することができます。

### モックにするのはどうなの？ DIにも必要だそうですが、それも悪だそうです。

モックについては後で詳しく説明します（そしてそれは悪ではありません）。
モックを使用して、実際に注入するものを、テストで制御および検査できる偽バージョンに置き換えます。
私たちの場合でも、標準ライブラリには、使用する準備ができています。

### Go標準ライブラリは本当に良いです。時間をかけて勉強してください。

このように`io.Writer`インターフェースにある程度慣れていることで、テストで`bytes.Buffer`を `Writer`として使うことができ、標準ライブラリの他の`Writer`を使ってコマンドラインアプリやウェブサーバで関数を使うことができます。

標準ライブラリに慣れるほど、これらの汎用インターフェイスが表示され、独自のコードで再利用して、ソフトウェアをさまざまなコンテキストで再利用可能にすることができます。

この例は、[プログラミング言語Go](https://www.amazon.co.jp/%E3%83%97%E3%83%AD%E3%82%B0%E3%83%A9%E3%83%9F%E3%83%B3%E3%82%B0%E8%A8%80%E8%AA%9EGo-ADDISON-WESLEY-PROFESSIONAL-COMPUTING-Donovan/dp/4621300253/ref=sr_1_6?__mk_ja_JP=%E3%82%AB%E3%82%BF%E3%82%AB%E3%83%8A&dchild=1&keywords=Go+Programming+Language&qid=1592323254&sr=8-6), の章に大きく影響されているため、これを楽しんだ場合、是非買ってみてください！
