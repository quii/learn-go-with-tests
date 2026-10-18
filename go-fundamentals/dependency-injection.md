---
description: Dependency Injection
---

# 依存性注入

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/di)**

これにはインターフェースの理解が必要になるため、[構造体のセクション](./structs-methods-and-interfaces.md)をすでに読んでいることが前提です。

プログラミングコミュニティには、依存性注入に関する誤解が _たくさん_ あります。このガイドを通じて、次のことを示せればと思います。

* フレームワークは必要ありません
* デザインが複雑になりすぎない
* テストを容易にします
* 優れた汎用関数を作成できます

`hello-world`の章で行ったように、誰かに挨拶する関数を書きたいのですが、今回は _実際の表示_ をテストします。

要約すると、その関数は次のようになります。

```go
func Greet(name string) {
	fmt.Printf("Hello, %s", name)
}
```

しかし、これをどのようにテストできますか？
`fmt.Printf`を呼び出すと stdout に出力されますが、テストフレームワークを使用してキャプチャするのはかなり困難です。

私たちがする必要があるのは、表示の依存関係を**注入** \(これは「渡す」を格好良く言っただけの言葉です\) できるようにすることです。

**関数は表示が _どこで_ 、 _どのように_ 行われるかを気にする必要がないので、具体的な型ではなく _インターフェース_ を受け入れるべきです。**

そうすれば、実装を変更して自分たちが制御できるものに表示させ、テストできるようになります。「実際の運用」では、stdoutに書き込むものを注入します。

[`fmt.Printf`](https://pkg.go.dev/fmt#Printf)のソースコードを見ると、フックする方法がわかります。

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

ここから、`os.Stdout`は`io.Writer`を実装していると推測できます。`Printf`は、`io.Writer`を期待する`Fprintf`に`os.Stdout`を渡しているからです。

さらに多くのGoコードを書くと、このインターフェイスが「このデータをどこかに置く」ための優れた汎用インターフェイスであるため、頻繁に登場することに気づくでしょう。

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

`bytes`パッケージの`Buffer`型は、`Write(p []byte) (n int, err error)`メソッドを持っているので、`Writer`インターフェースを実装しています。

テストでこれを使用して`Writer`として送信し、`Greet`を呼び出した後に何が書き込まれたかを確認できます。

## テストを実行してみます

テストはコンパイルされません

```text
./di_test.go:10:2: undefined: Greet
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

_コンパイラの言うことを聞いて_ 、問題を修正してください。

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Printf("Hello, %s", name)
}
```

`Hello, Chris di_test.go:16: got '' want 'Hello, Chris'`

テストは失敗します。名前は出力されますが、標準出力になることに注意してください。

## 成功させるのに十分なコードを書く

テストでは、ライターを使用して挨拶をバッファに送信します。`fmt.Fprintf`は`fmt.Printf`に似ていますが、代わりに `Writer`を受け取って文字列を送信します。一方、`fmt.Printf`のデフォルトはstdoutです。

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}
```

テストに合格しました。

## リファクタリング♪

以前のコンパイラーは、`bytes.Buffer`へのポインターを渡すように指示しました。これは技術的には正しいですが、あまり役に立ちません。

これを実証するために、`Greet`関数を標準出力に出力するGoアプリケーションに接続してみてください。

```go
func main() {
	Greet(os.Stdout, "Elodie")
}
```

`./di.go:14:7: cannot use os.Stdout (type *os.File) as type *bytes.Buffer in argument to Greet`

前に説明したように、`fmt.Fprintf`を使用すると、`os.Stdout`と `bytes.Buffer`の両方が実装していることがわかっている`io.Writer`を渡すことができます。

より汎用的なインターフェースを使用するようにコードを変更すると、テストとアプリケーションの両方で使用できるようになります。

```go
package main

import (
	"fmt"
	"io"
	"os"
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
`Greet`関数はどれほど汎用的なのでしょうか？

### インターネット

以下を実行します

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "world")
}

func main() {
	log.Fatal(http.ListenAndServe(":5001", http.HandlerFunc(MyGreeterHandler)))
}
```

プログラムを実行し、[http://localhost:5001](http://localhost:5001)に移動します。グリーティング機能が使用されているのがわかります。

HTTPサーバーについては後の章で説明しますので、詳細についてはあまり気にしないでください。

HTTPハンドラーを作成すると、 `http.ResponseWriter`と、リクエストの作成に使用された`http.Request`が与えられます。サーバーを実装するときは、ライターを使用して応答を _書き込み_ ます。

`http.ResponseWriter`も`io.Writer`を実装していると推測できるでしょう。だからこそ、ハンドラー内で `Greet`関数を再利用できたのです。

## まとめ

最初のコードは、制御できない場所にデータを書き込んだため、簡単にテストできませんでした。

_テストに動機付けられて_ 、**依存関係を注入する**ことでデータが _どこに_ 書き込まれるかを制御できるようにコードをリファクタリングしました。これにより次のことが可能になりました。

* **コードをテストする** 関数を _簡単に_ テストできない場合は、通常、依存関係が関数に組み込まれているか、 _または_ グローバルな状態のためです。たとえば、ある種のサービス層で使用されているグローバルデータベース接続プールがある場合、テストが困難になる可能性が高く、実行が遅くなります。DIは、\(インターフェイスを介して\)データベースの依存関係を注入するように動機付けし、テストで制御できるものでモックアウトできます。
* **関心事を分離する** 「 _データの行き先_ 」と「 _データの生成方法_ 」を分離します。メソッド/関数の責任が多すぎると感じた場合は\(データを生成 _して_ DBに書き込む？HTTPリクエストを処理 _して_ ドメインレベルのロジックも実行する？\)、おそらくDIが必要なツールになるでしょう。
* **コードをさまざまなコンテキストで再利用できるようにする** コードを使用できる最初の「新しい」コンテキストは、テスト内です。しかし、さらに誰かがあなたの関数で何か新しいことを試したい場合、彼らは彼ら自身の依存関係を注入することができます。

### モックにするのはどうなの？ DIにも必要だそうですが、それも悪だそうです。

モックについては後で詳しく説明します\(そしてそれは悪ではありません\)。
モックを使用して、実際に注入するものを、テストで制御および検査できる偽バージョンに置き換えます。
ただ私たちの場合は、標準ライブラリに使用する準備ができているものがありました。

### Go標準ライブラリは本当に良いです。時間をかけて勉強してください。

このように`io.Writer`インターフェースにある程度慣れていることで、テストで`bytes.Buffer`を `Writer`として使うことができ、標準ライブラリの他の`Writer`を使ってコマンドラインアプリやウェブサーバで関数を使うことができます。

標準ライブラリに慣れるほど、これらの汎用インターフェイスが表示され、独自のコードで再利用して、ソフトウェアをさまざまなコンテキストで再利用可能にすることができます。

この例は、[The Go Programming language](https://www.amazon.co.uk/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440)\(邦訳『[プログラミング言語Go](https://www.amazon.co.jp/dp/4621300253)』\)の章に大きく影響されているため、これを楽しんだ場合、是非買ってみてください！
