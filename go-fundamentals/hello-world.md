---
description: 'Hello, World'
---

# Hello, World

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/hello-world)**

新しい言語での最初のプログラムが [Hello, World](https://en.m.wikipedia.org/wiki/%22Hello,_World!%22_program)であるのは伝統的です。

- 好きな場所にフォルダを作成します。
- その中に `hello.go` という名前の新しいファイルを作り、その中に以下のコードを記述します。

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
```

実行するには、`go run hello.go` と入力します。

## 使い方

Goでプログラムを作成すると、その中に `main`関数が定義された`main`パッケージが作成されます。 パッケージは、関連するGoコードをグループ化する方法です。

`func`キーワードは、名前と本体を持つ関数を定義します。

`import "fmt"`では、印刷に使用する `Println`関数を含むパッケージをインポートしています。

## テスト方法

どのようにテストすればよいと思いますか？ 「ドメイン」コードを外界\(副作用\)から分離することは良いことです。 `fmt.Println`は副作用\(標準出力への出力\)であり、送信する文字列はドメインです。

テストを簡単にするために、これらの懸念事項を分離しましょう

```go
package main

import "fmt"

func Hello() string {
	return "Hello, world"
}

func main() {
	fmt.Println(Hello())
}
```

`func`を使用して新しい関数を作成しましたが、今回は定義に別のキーワード`string`を追加しました。 つまり、この関数は `string`を返します。

ここで、 `hello_test.go`という新しいファイルを作成します。ここで、`Hello`関数のテストを記述します

```go
package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	want := "Hello, world"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

## Go Modulesってなに？

次は、テストを実行します。ターミナルで `go test` と入力してください。テストが成功したら、おそらく以前のバージョンの Go を使っているのでしょう。しかし、Go 1.16 以降を使っている場合は、テストが実行されない可能性が高いです。その代わり、ターミナルに次のようなエラーメッセージが表示されます。

```shell
$ go test
go: cannot find main module; see 'go help modules'
```

何が問題なのか？一言で言えば、[モジュール](https://blog.golang.org/go116-module-changes)です。幸いなことに、この問題は簡単に解決することができます。ターミナルで `go mod init example.com/hello` と入力してください。そうすると、次のような内容の新しいファイルが作成されます。

```
module example.com/hello

go 1.16
```

このファイルは `go` ツールに、あなたのコードに関する重要な情報を伝えます。もし、アプリケーションを配布する予定があるなら、コードがダウンロードできる場所や、依存関係の情報も含めるでしょう。モジュールの名前である example\.com\/hello は、通常、そのモジュールを見つけてダウンロードできる URL を指します。このあとすぐに使い始めるツールとの互換性のために、example\.com/hello の .com のドットのように、モジュール名のどこかにドットを含めるようにしてください。今のところ、モジュールファイルは最小限のものであり、そのままにしておいても構いません。モジュールについて詳しく知りたいなら、[Golangドキュメントのリファレンスを参照してください](https://golang.org/doc/modules/gomod-ref)。これでGo 1.16 でもテストは実行されるはずなので、テストと Go の学習に戻ることができます。

今後の章では、`go test` や `go build` といったコマンドを実行する前に、それぞれの新しいフォルダで `go mod init SOMENAME` を実行する必要があります。

はっきりさせておくと、モジュールの名前である `SOMENAME` は、`package main`（この本でこれまでに出てきたすべての `.go` ファイルが先頭で宣言しているもの）とは無関係です。モジュール名はプロジェクト全体を表す単なる識別子であり、`main` という名前にする必要はありません。どのような名前を選んでも、同じフォルダの中から実行すれば `go run`、`go test`、`go build` はすべて正しく動作します。

## テストに戻ります

ターミナルで`go test`を実行します。合格したはずです！ 確認するために、`want`文字列を変更して、意図的にテストを壊してみてください。

複数のテストフレームワークから選択したり、そのインストール方法を理解したりする必要がないことに注目してください。 必要なものはすべて言語に組み込まれており、構文は、これから記述する残りのコードと同じです。

### テストを書く

テストの作成は、関数の作成と同様であり、いくつかのルールがあります。

* `xxx_test.go`のような名前のファイルにある必要があります。
* テスト関数は`Test`という単語で始まる必要があります。
* テスト関数は1つの引数のみをとります。 `t *testing.T`
* `*testing.T` 型を使うには、他のファイルの `fmt` と同じように `import "testing"` が必要です。

とりあえず、 `*testing.T`型の`t`がテストフレームワークへの「フック」であることを知っていれば十分です。これにより、失敗させたいときに `t.Fail()`のようなことを実行できます。

新しいトピックをいくつか取り上げました。

#### `if`
Goのif文は、他のプログラミング言語のものとよく似ています。

#### 変数の宣言

いくつかの変数を構文`varName := value`で宣言しています。これにより、読みやすくするためにテストでいくつかの値を再利用できます。

#### `t.Errorf`

`t`の`Errorf` _メソッド_ を呼び出しています。これはメッセージを出力してテストを失敗させます。 `f`はフォーマット\(format\)を表し、プレースホルダー値`%q`に値が挿入された文字列を作成できます。 テストを失敗させてみれば、それがどのように機能するかは明らかでしょう。

プレースホルダー文字列の詳細については、[fmtのドキュメント](https://pkg.go.dev/fmt#hdr-Printing)を参照してください。 テストでは、 `%q` は値を二重引用符で囲むので非常に便利です。

メソッドと関数の違いについては後で説明します。

### Goのドキュメント

Goのもう1つの便利な機能はドキュメントです。先ほど、公式のパッケージ閲覧サイトでfmtパッケージのドキュメントを見ましたが、Goにはオフラインでドキュメントにすばやくアクセスする方法も用意されています。

Goには doc という組み込みツールがあり、システムにインストールされている任意のパッケージや、現在作業中のモジュールを調べることができます。先ほどと同じ、Printing の書式指定子\(verbs\)のドキュメントを表示するには次のようにします。

```
$ go doc fmt
package fmt // import "fmt"

Package fmt implements formatted I/O with functions analogous to C's printf and
scanf. The format 'verbs' are derived from C's but are simpler.

# Printing

The verbs:

General:

    %v	the value in a default format
    	when printing structs, the plus flag (%+v) adds field names
    %#v	a Go-syntax representation of the value
    %T	a Go-syntax representation of the type of the value
    %%	a literal percent sign; consumes no value
...
```

ドキュメントを閲覧するためのGoの2つ目のツールは pkgsite コマンドで、Goの公式パッケージ閲覧サイトを動かしているものです。`go install golang.org/x/pkgsite/cmd/pkgsite@latest` で pkgsite をインストールし、`pkgsite -open .` で実行できます。Goの install コマンドは、そのリポジトリからソースファイルをダウンロードし、実行可能なバイナリにビルドします。Goをデフォルトでインストールした場合、その実行ファイルは Linux と macOS では `$HOME/go/bin`、Windows では `%USERPROFILE%\go\bin` に置かれます。まだこれらのパスを $PATH 変数に追加していない場合は、go install したコマンドを実行しやすくするために追加しておくとよいでしょう。

標準ライブラリの大部分には、例を含む優れたドキュメントがあります。 [http://localhost:8080/testing](http://localhost:8080/testing) に移動して、何が利用できるかを確認してみる価値があります。


### Hello, YOU

これでテストができたので、ソフトウェアを安全に反復開発できます。

前の例では、テストの書き方と関数の宣言方法の例を示すために、コードを書いた _後で_ テストを書きました。ここからは、 _テストを最初に書いて_ いきます。

次の要件は、挨拶の受信者を指定できるようにすることです。

これらの要件をテストに取り込むことから始めましょう。 これは基本的なテスト駆動開発であり、テストが _実際に_ 望むものをテストしていることを確認できます。テストを後から作成すると、コードが意図したとおりに機能しなくても、テストが引き続きパスしてしまうリスクがあります。

```go
package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Chris")
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

`go test`を実行すると、コンパイルエラーが発生するはずです

```text
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

Goのような静的型付け言語を使用する場合、 _コンパイラの声に耳を傾ける_ ことが重要です。 コンパイラは、コードがどのように組み合わさって動作すべきかを理解しているので、あなたがそれを考える必要はありません。

この場合、コンパイラは続行するために何をする必要があるかを指示しています。引数を受け付けるには、関数 `Hello`を変更する必要があります。

文字列型の引数を受け入れるように`Hello`関数を編集します

```go
func Hello(name string) string {
	return "Hello, world"
}
```

もう一度テストを実行すると、引数を渡していないため、 `hello.go`はコンパイルに失敗します。コンパイルできるように`"world"`を渡してください。

```go
func main() {
	fmt.Println(Hello("world"))
}
```

テストを実行すると、次のように表示されます。

```text
hello_test.go:10: got 'Hello, world' want 'Hello, Chris''
```

ようやくコンパイルできるプログラムができましたが、テストによると要件を満たしていません。

`name`引数を使用してテストに合格し、 `Hello,`と連結してみましょう

```go
func Hello(name string) string {
	return "Hello, " + name
}
```

テストを実行すると、テストに合格するはずです。通常、TDDサイクルの一部として、ここで _リファクタリング_ を行うべきです。

### ソース管理に関する注意

この時点で、ソース管理を使用している場合は\(使用すべきです！\)、コードをそのまま `commit` します。テストに裏打ちされた動作するソフトウェアがあります。

ただし、次にリファクタリングする予定なので、mainにプッシュすることは _しません_ 。何らかの理由でリファクタリングで混乱に陥った場合に備えて、この時点でコミットしておくと便利です。いつでも動作するバージョンに戻ることができます。

ここでリファクタリングすることは多くありませんが、 _定数_ という別の言語機能を導入できます。

### 定数

定数は次のように定義できます。

```go
const englishHelloPrefix = "Hello, "
```

コードをリファクタリングできるようになりました。

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
	return englishHelloPrefix + name
}
```

リファクタリング後、テストを再実行して、何も壊れていないことを確認します。

値の意味を表現するために、また場合によってはパフォーマンスを向上させるために、定数を作成することを検討する価値があります。

## Hello, world... もう一度

次の要件は、関数が空の文字列で呼び出されたときに、`"Hello, "`ではなくデフォルトで`"Hello, World"`を出力することです。

新しい失敗するテストを書くことから始めます。

```go
func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris")
		want := "Hello, Chris"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
```

ここでは、テストの武器に別のツール、サブテストを導入しています。 「もの」を中心にテストをグループ化し、さまざまなシナリオを説明するサブテストを作成すると便利な場合があります。

このアプローチの利点は、他のテストで使用できる共有コードを設定できることです。

テストが失敗している間に、 `if`を使用してコードを修正しましょう。

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return englishHelloPrefix + name
}
```

テストを実行すると、新しい要件を満たし、他の機能を誤って壊していないことがわかります。

テストが、コードが何をする必要があるのかの _明確な仕様_ であることは重要です。しかし、メッセージが期待どおりかどうかを確認するところでコードが繰り返されています。

リファクタリングは、プロダクションコード _だけ_ のものではありません！

テストが通るようになったので、テストもリファクタリングすることができますし、そうすべきです。

```go
func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("empty string defaults to 'world'", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})

}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

さて、ここで何をしましたか？

アサーションを新しい関数にリファクタリングしました。 これにより、重複が削減され、テストの可読性が向上します。 `t *testing.T`を渡す必要があるので、必要なときにテストコードを失敗させることができます。

ヘルパー関数については、 `*testing.T` と `*testing.B` の両方が満たすインターフェースである `testing.TB` を受け入れるのがよいでしょう。そうすれば、テストからでもベンチマークからでもヘルパー関数を呼び出すことができます\(「インターフェース」のような言葉が今はまったくわからなくても心配はいりません。後で説明します\)。

`t.Helper()`は、このメソッドがヘルパーであることをテストスイートに伝えるために必要です。こうすることで、テストが失敗したときに報告される行番号は、テストヘルパーの中ではなく _関数呼び出し_ の箇所を示します。これにより、他の開発者が問題を追跡しやすくなります。それでも理解できない場合は、コメントアウトし、テストを失敗させて、テスト出力を観察してください。コメントは、コードに追加情報を加えるのに最適な方法です。この場合は、コンパイラに特定の行を無視するように指示する手っ取り早い方法です。`t.Helper()` のコードをコメントアウトするには、行頭に2つのフォワードスラッシュ `//` を追加してください。その行がグレーになるか、他のコードとは別の色に変わることで、コメントアウトされたことがわかるはずです。

同じ型の引数が複数ある場合\(この例では2つの文字列\)、`(got string, want string)` と書く代わりに `(got, want string)` と短縮できます。

### ソース管理に戻る

コードに満足できたので、前のコミットを修正（amend）して、テスト付きの素敵なバージョンのコードだけをチェックインします。

### 規律

サイクルをもう一度見てみましょう

* テストを書く
* コンパイラーをパスする
* テストを実行し、失敗することを確認し、エラーメッセージが意味があることを確認します
* テストに合格するのに十分なコードを記述します
* リファクタリング

一見面倒に見えるかもしれませんが、フィードバックループを守ることが重要です。

これは、 _関連するテスト_ があることを保証するだけでなく、テストの安全性を考慮してリファクタリングすることにより、 _優れたソフトウェアを設計する_ ことを保証するのに役立ちます。

テストが失敗したことを確認することは、エラーメッセージがどのように表示されるかを確認できるため、重要なチェックです。開発者としては、テストに失敗して問題が何であるかについて明確なアイデアが得られない場合、コードベースを操作するのは非常に困難です。

テストが _高速_ であることを確認し、テストを簡単に実行できるようにツールを設定することで、コードを記述するときにフローの状態に入ることができます。

テストを書かないということは、ソフトウェアを実行して手動でコードをチェックし続けることを意味し、それはフローの状態を壊します。特に長期的には、時間を節約することにはなりません。

## 立ち止まるな！ その他の要件

なんと、さらに要件があります。 次に、挨拶の言語を指定する2番目のパラメーターをサポートする必要があります。認識されない言語が渡された場合は、デフォルトで英語にします。

TDDを使用してこの機能を簡単に具体化できると確信しているはずです。

スペイン語を渡すユーザーのテストを作成します。既存のスイートに追加します。

```go
	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})
```

ズルをしないことを忘れないでください！ _最初にテスト_ 。テストを実行しようとすると、1つではなく2つの引数を指定して `Hello`を呼び出しているため、コンパイラは文句を言う _はず_ です。

```text
./hello_test.go:27:19: too many arguments in call to Hello
    have (string, string)
    want (string)
```

`Hello`に別の文字列引数を追加して、コンパイルの問題を修正します

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}
	return englishHelloPrefix + name
}
```

テストをもう一度実行すると、他のテストと `hello.go`で`Hello`に十分な引数を渡していないというメッセージが表示されます。

```text
./hello.go:15:19: not enough arguments in call to Hello
    have (string)
    want (string, string)
```

空の文字列を渡すことで修正します。これで、新しいシナリオを除いて、すべてのテストがコンパイルされ、 _かつ_ 合格するはずです。

```text
hello_test.go:29: got 'Hello, Elodie' want 'Hola, Elodie'
```

ここで`if`を使用して、言語が`"Spanish"`に等しいことを確認し、そうであればメッセージを変更します

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language == "Spanish" {
		return "Hola, " + name
	}
	return englishHelloPrefix + name
}
```

これでテストに合格するはずです。

さて、 _リファクタリング_ の時間です。コードにいくつかの問題が見られるはずです。「マジック」文字列があり、その一部は繰り返されています。自分でリファクタリングを試してください。変更を加えるたびに、テストを再実行して、リファクタリングが何も壊していないことを確認してください。

```go
	const spanish = "Spanish"
	const englishHelloPrefix = "Hello, "
	const spanishHelloPrefix = "Hola, "

	func Hello(name string, language string) string {
		if name == "" {
			name = "World"
		}

		if language == spanish {
			return spanishHelloPrefix + name
		}
		return englishHelloPrefix + name
	}
```

### フランス語

* `"French"`を渡すと、 `"Bonjour, "`が得られることを表明するテストを作成します
* それが失敗するのを見て、エラーメッセージが読みやすいことを確認してください
* コードに最小限の合理的な変更を加える

大体このようなものを書いたかもしれません。

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	if language == spanish {
		return spanishHelloPrefix + name
	}
	if language == french {
		return frenchHelloPrefix + name
	}
	return englishHelloPrefix + name
}
```

## `switch`

特定の値をチェックする多くの `if`ステートメントがある場合、代わりに`switch`ステートメントを使用するのが一般的です。後で言語サポートを追加したい場合は、 `switch`を使用してコードをリファクタリングして読みやすくし、拡張性を高めることができます。

```go
func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	prefix := englishHelloPrefix

	switch language {
	case spanish:
		prefix = spanishHelloPrefix
	case french:
		prefix = frenchHelloPrefix
	}

	return prefix + name
}
```

`prefix` に対しては `:=` を一度だけ使い（新しい変数を宣言して初期値を与える）、それ以外で値を変更している箇所では、上の `name` も `switch` 内の `prefix` も、すべて単なる `=` を使っていることに注目してください。`:=` はGoの[短い変数宣言](https://go.dev/ref/spec#Short_variable_declarations)で、新しい変数を作成します。`=` は単なる[代入](https://go.dev/ref/spec#Assignment_statements)で、すでに存在する変数の値を変更します（`name` は引数としてすでに存在しており、`prefix` は数行上で `:=` によってすでに宣言されています）。同じスコープにすでに存在する変数に `:=` を使ったり、まだ存在しない変数に `=` を使ったりすると、コンパイルエラーになります。

選択した言語での挨拶を含めるテストを作成してみてください。私たちの _素晴らしい_ 関数を拡張するのがいかに簡単かがわかるはずです。

### あとひとつ...最後に...リファクター？

関数が少し大きくなってきていると主張することもできるでしょう。このための最も簡単なリファクタリングは、一部の機能を別の関数に抽出することです。

```go

const (
	spanish = "Spanish"
	french  = "French"

	englishHelloPrefix = "Hello, "
	spanishHelloPrefix = "Hola, "
	frenchHelloPrefix  = "Bonjour, "
)

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}
```

いくつかの新しい概念

* 関数のシグネチャでは、 _名前付き戻り値_ `(prefix string)`を作成しました。
* これにより、関数に `prefix` という変数が作成されます。
  * 「ゼロ」値が割り当てられます。これは型によって異なります。たとえば、`int`は`0`で、`string`の場合は`""`です。
    * `return prefix`ではなく`return`を呼び出すだけで、設定されているものを返すことができます。
  * これは関数のGo Docに表示されるので、コードの意図をより明確にすることができます。
* switchケースの`default` は、他の`case`ステートメントのいずれも一致しない場合に分岐します。
* 関数名は小文字で始まります。 Goでは、パブリック関数は大文字で始まり、プライベート関数は小文字で始まります。アルゴリズムの内部を世界に公開したくないので、この関数をプライベートにしました。
* また、定数はそれぞれ別の行で宣言する代わりに、ブロックにまとめることができます。読みやすさのために、関連する定数のまとまりの間に空行を入れるとよいでしょう。

## まとめ

`Hello, world`からこんなに多くを得ることができると誰が知ってたのでしょう。

これで、次のことをある程度理解できたはずです。

### Goの構文のいくつか

* テストを書く
* 引数と戻り値の型を使用した関数の宣言
* `if`、`const` および `switch`
* 変数と定数の宣言

### TDDプロセスとそのステップが _なぜ_ 重要なのか

* _失敗するテストを書き、失敗するのを確認する_ ことで、要件に対して _関連する_ テストを書いたこと、そしてそれが _理解しやすい失敗の説明_ を出力することを確認しました。
* 機能するソフトウェアがあることを確認するために、最小限のコードを記述して合格させる。
* _その後_ リファクタリングする。テストの安全性に裏打ちされているので、扱いやすく巧みに作られたコードであることを確認できます。

今回のケースでは、 `Hello()` から `Hello("name")` へ、そして `Hello("name", "French")` へと、小さくて理解しやすいステップで進みました。

もちろん、これは「現実の世界」のソフトウェアに比べれば取るに足らないことですが、原則は変わりません。 TDDは、身につけるのに練習が必要なスキルですが、問題をテスト可能な小さなコンポーネントに分解することで、ソフトウェアの作成がはるかに簡単になります。
