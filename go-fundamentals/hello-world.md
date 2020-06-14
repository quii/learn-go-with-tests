---
description: 'Hello, World'
---

# Hello, World

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/hello-world)

新しい言語での最初のプログラムが [Hello, World](https://en.m.wikipedia.org/wiki/%22Hello,_World!%22_program)であるのは伝統的です。

[前の章](install-go.md#go-environment) では、ファイルの配置場所についてGoがどのように考えられているかを説明しました。

次のパス `$GOPATH/src/github.com/{your-user-id}/hello`にディレクトリを作成します。

したがって、UNIXベースのOSを使用していて、 `$GOPATH` を設定している場合は、次のコマンドでディレクトリを作成できます。 `mkdir -p $GOPATH/src/github.com/$USER/hello`.

以降の章では、コードを好きな名前で新しいフォルダーを作成できます。たとえば、次の章では `$GOPATH/src/github.com/{your-user-id}/integers` にコードを配置することをお勧めします。
このサイトの一部のユーザーは、`learn-go-with-tests/hello`. などのすべての作業用のフォルダーを作成することを好みます。
つまり、フォルダをどのように構成するかはあなた次第です。

このディレクトリに`hello.go`というファイルを作成し、このコードを記述します。実行するには`go run hello.go`と入力します。

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world")
}
```

## 使い方

Goでプログラムを作成すると、その中に `main`関数が定義された` main`パッケージが作成されます。
パッケージは、関連するGoコードをグループ化する方法です。

`func`キーワードは、名前と本体で関数を定義する方法です。

`import "fmt"`では、印刷に使用する `Println`関数を含むパッケージをインポートしています。

## テスト方法

どのようにテストすればよいと思いますか？
「ドメイン」コードを外界から分離することは良いことです。 `fmt.Println`は副作用であり、送信する文字列はドメインです。

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

`func`を使用して新しい関数を再度作成しましたが、今回は定義に別のキーワード`string`を追加しました。
つまり、この関数は `string`を返します。

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

説明する前に、コードを実行してみましょう。
端末で`go test`を実行します。合格したはずです！
確認するために、`want`文字列を変更して、意図的にテストを中断してみてください。

複数のテストフレームワークを選択する必要がなく、インストール方法を理解する必要がないことに注意してください。
必要なものはすべて言語に組み込まれており、構文は、これから記述する残りのコードと同じです。

### テストを書く

テストの作成は、関数の作成と同様であり、いくつかのルールがあります。

* `xxx_test.go`のような名前のファイルにある必要があります。
* テスト関数は`Test`という単語で始まる必要があります。
* テスト関数は1つの引数のみをとります。 `t *testing.T`

とりあえず、 `* testing.T`タイプの`t`がテストフレームワークへの`hook`(フック)であることを知っていれば十分なので、失敗したいときに `t.Fail()`のようなことを実行できます。

新しいトピックをいくつか取り上げました。

#### `if`

Goのステートメントが他のプログラミング言語とよく似ている場合。

#### Declaring variables

いくつかの変数を構文`varName := value`で宣言しています。これにより、読みやすくするためにテストでいくつかの値を再利用できます。

#### `t.Errorf`

メッセージを出力してテストに失敗する`t`で`Errorf` _method_ を呼び出しています。
`f`は、プレースホルダー値`％q`に値が挿入された文字列を作成できる形式を表します。
テストを失敗させたとき、それがどのように機能するかは明らかです。

プレースホルダー文字列の詳細については、 [fmt go doc](https://golang.org/pkg/fmt/#hdr-Printing)。
テストでは、 `%q` は値を二重引用符で囲むので非常に便利です。

メソッドと関数の違いについては後で説明します。

### Go ドキュメント

Goのもう1つの質機能はドキュメントです。
`godoc -http :8000`を実行すると、ローカルでドキュメントを起動できます。
[localhost:8000/pkg](http://localhost:8000/pkg) に移動すると、システムにインストールされているすべてのパッケージが表示されます。

標準ライブラリの大部分には、例を含む優れたドキュメントがあります。
[http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/) に移動すると、何が利用できるかを確認する価値があります。

`godoc` コマンドがない場合は、`godoc`を含まない新しいバージョンのGo（`1.14`以降）を使用している可能性があります [no longer including `godoc`](https://golang.org/doc/go1.14#godoc)。
`go get golang.org/x/tools/cmd/godoc`を使用して手動でインストールできます。

### Hello, YOU

これでテストが完了したので、ソフトウェアを安全に反復できます。

最後の例では、テストを記述した後、コードを記述したので、テストを記述して関数を宣言する方法の例を取得できます。この時点から、`テストを最初に作成`します。

次の要件は、挨拶の受信者を指定できるようにすることです。

これらの要件をテストに取り込むことから始めましょう。
これは基本的なテスト主導の開発であり、テストが希望どおりに実際にテストされていることを確認できます。テストをさかのぼって作成すると、コードが意図したとおりに機能しなくても、テストが引き続きパスする可能性があります。

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

Goのような静的に型付けされた言語を使用する場合、コンパイラーをリッスンすることが重要です。
コンパイラーは、コードがどのようにスナップして機能するかを理解しているので、そうする必要はありません。

この場合、コンパイラーは続行するために何をする必要があるかを指示しています。引数を受け付けるには、関数 `Hello`を変更する必要があります。

文字列型の引数を受け入れるように`Hello`関数を編集します

```go
func Hello(name string) string {
    return "Hello, world"
}
```

もう一度テストを実行すると、引数を渡していないため、 `main.go`はコンパイルに失敗します。それを通過させるために`"world"`を送ってください。

```go
func main() {
    fmt.Println(Hello("world"))
}
```

テストを実行すると、次のように表示されます。

```text
hello_test.go:10: got 'Hello, world' want 'Hello, Chris''
```

ようやくコンパイルプログラムができましたが、テストによると要件を満たしていません。

`name`引数を使用してテストに合格し、 `Hello,`と連結してみましょう

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

テストを実行すると、テストに合格するはずです。通常、TDDサイクルの一部として、 _refactor_ を実行する必要があります。

### ソース管理に関する注意

この時点で、ソース管理を使用している場合はそうする必要があります！。 コードをそのまま `commit` します。テストに裏打ちされた実用的なソフトウェアがあります。

ただし、次にリファクタリングする予定なので、マスターにプッシュしません。何らかの理由でリファクタリングで混乱に陥った場合に備えて、この時点でコミットすると便利です。いつでも作業バージョンに戻ることができます。

ここでリファクタリングすることは多くありませんが、 _constants_ という別の言語機能を導入できます。

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

定数は、 `Hello`が呼び出されるたびに`"Hello、"`文字列インスタンスを作成する手間を省くため、アプリケーションのパフォーマンスを向上させるはずです。

明確にするために、この例ではパフォーマンスの向上はごくわずかです。ただし、値の意味を把握するために、また場合によってはパフォーマンスを支援するために定数を作成することを検討する価値があります。

## Hello, world... もう一度

次の要件は、関数が空の文字列で呼び出されたときに、デフォルトで`"Hello、"`ではなく`"Hello、World"`を出力することです。

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

ここでは、テストの武器であるサブテストに別のツールを導入しています。 「もの」を中心にテストをグループ化し、さまざまなシナリオを説明するサブテストを作成すると便利な場合があります。

このアプローチの利点は、他のテストで使用できる共有コードを設定できることです。

メッセージが期待どおりかどうかを確認するときにコードが繰り返されます。

リファクタリングは、量産コードにとって `ちょうど`ではありません！

テストでは、コードが何をする必要があるのか​​を明確に指定することが重要です。

テストをリファクタリングすることができます。

```go
func TestHello(t *testing.T) {

    assertCorrectMessage := func(t *testing.T, got, want string) {
        t.Helper()
        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    }

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"
        assertCorrectMessage(t, got, want)
    })

    t.Run("empty string defaults to 'World'", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"
        assertCorrectMessage(t, got, want)
    })

}
```

さて、ここで何をしましたか？

アサーションを関数にリファクタリングしました。
これにより、重複が削減され、テストの可読性が向上します。
Goでは、他の関数内で関数を宣言して、変数に割り当てることができます。
その後、通常の関数と同じようにそれらを呼び出すことができます。 `t * testing.T`を渡す必要があるので、必要なときにテストコードを失敗させることができます。

このメソッドがヘルパーであることをテストスイートに伝えるには、 `t.Helper()`が必要です。失敗したときにこれを行うと、レポートされる行番号はテストヘルパー内ではなく、`関数呼び出し`内にあります。これにより、他の開発者が問題を簡単に追跡できるようになります。それでも理解できない場合は、コメントアウトし、テストを失敗させて、テスト出力を観察します。

うまく書かれた不合格のテストができたので、 `if`を使用してコードを修正しましょう。

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

### ソース管理に戻る

さて、前のコミットを修正するコードに満足しているので、コードの素敵なバージョンをテストでチェックインするだけです。

### 規律

サイクルをもう一度見てみましょう

* テストを書く
* コンパイラーをパスする
* テストを実行し、失敗することを確認し、エラーメッセージが意味があることを確認します
* テストに合格するのに十分なコードを記述します
* リファクタリング

一見面倒に見えるかもしれませんが、フィードバックループを守ることが重要です。

これは、`関連するテスト`があることを保証するだけでなく、テストの安全性を考慮してリファクタリングすることにより、`優れたソフトウェアを設計する`ことを保証するのに役立ちます。

テストが失敗したことを確認することは、エラーメッセージがどのように表示されるかを確認できるため、重要なチェックです。開発者としては、テストに失敗して問題が何であるかについて明確なアイデアが得られない場合、コードベースを操作するのは非常に困難です。

テストが`fast`であることを確認し、テストを簡単に実行できるようにツールを設定することで、コードを記述するときにフローの状態に入ることができます。

テストを記述しないことにより、フローの状態を壊すソフトウェアを実行して手動でコードをチェックすることを約束し、特に長期的には、時間を節約できなくなります。

## 立ち止まるな！ その他の要件

よかった、もっと要件があります。
次に、挨拶の言語を指定する2番目のパラメーターをサポートする必要があります。認識されない言語が渡された場合は、デフォルトで英語に設定されます。

TDDを使用してこの機能を簡単に具体化できると確信しているはずです。

スペイン語で合格するユーザーのテストを作成します。既存のスイートに追加します。

```go
    t.Run("in Spanish", func(t *testing.T) {
        got := Hello("Elodie", "Spanish")
        want := "Hola, Elodie"
        assertCorrectMessage(t, got, want)
    })
```

不正行為をしないことを忘れないでください！ `最初にテスト`。テストを実行しようとすると、1つではなく2つの引数を指定して `Hello`を呼び出すため、コンパイラは文句を言うべきです_。

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

テストをもう一度実行すると、他のテストと `hello.go`で` Hello`に十分な引数を渡さないというメッセージが表示されます。

```text
./hello.go:15:19: not enough arguments in call to Hello
    have (string)
    want (string, string)
```

空の文字列を渡すことで修正します。これで、新しいシナリオを除いて、すべてのテストで`and`がコンパイルされます。

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

 さて、リファクタリングのお時間です。コードにいくつかの問題が見られるはずです。`"magic"`文字列は、その一部が繰り返されます。自分で試してリファクタリングしてください。変更を加えるたびに、テストを再実行して、リファクタリングが何も壊していないことを確認してください。

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

* `"French"`を渡すと、 `"Bonjour、"`が得られることを表明するテストを作成します
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
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    }

    return prefix + name
}
```

選択した言語で挨拶を含めるためのテストを作成すると、`amazing`関数を拡張するのがいかに簡単かがわかります。

### あとひとつ...最後に...リファクター？

多分私たちの機能が少し大きくなっていると主張することができます。このための最も簡単なリファクタリングは、一部の機能を別の関数に抽出することです。

```go
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

* 関数のシグネチャでは、 _named return value_ `(prefix string)`を作成しました。
* これにより、関数に `prefix` という変数が作成されます。
  * "zero" 値が割り当てられます。これはタイプによって異なります。たとえば、`int`は`0`で、文字列の場合は`""`です。
    * `return prefix`ではなく`return`を呼び出すだけで、設定されているものを返すことができます。
  * これは関数のGo Docに表示されるので、コードの意図をより明確にすることができます。
* switchケースの`default` は、他の`case`ステートメントのいずれも一致しない場合に分岐します。
* 関数名は小文字で始まります。 Goでは、パブリック関数は大文字で始まり、プライベート関数は小文字で始まります。アルゴリズムの内部を世界に公開したくないので、この関数をプライベートにしました。

## まとめ

`Hello, world`からこんなに多くを得ることができると誰が知ってたのでしょう。

これで、次のことをある程度理解できたはずです。

### Goの構文のいくつか

* テストを書く
* 引数と戻り値の型を使用した関数の宣言
* `if`、`const` および `switch`
* 変数と定数の宣言

### TDDプロセスとそのステップが重要である理由

* `失敗するテストを作成してそれを確認する`要件に対応する`relevant`テストを作成し、失敗の説明を簡単に理解できることを確認しました。
* 機能するソフトウェアがあることを確認するために、最小限のコードを記述して合格させる。
* リファクタリング、テストの安全性に裏打ちされており、操作が簡単な巧妙に作成されたコードがあることを確認します

今回のケースでは、 `Hello()` から `Hello("name")` から `Hello("name", "French")` に、小さくて簡単な手順で進みました。

もちろん、これは「現実の世界」のソフトウェアに比べれば取るに足らないことですが、原則は変わりません。
TDDは、開発するための練習が必要なスキルですが、問題をより小さなコンポーネントに分解してテストできるため、ソフトウェアの作成がはるかに簡単になります。
