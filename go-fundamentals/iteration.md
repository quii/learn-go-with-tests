---
description: Iteration
---

# 反復、繰り返し

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/for)

Goで繰り返し作業を行うには、 `for`が必要です。 Goには `while`、`do`、 `until`キーワードはなく、`for`のみ使用できます。これは良いことです！

文字を5回繰り返す関数のテストを書いてみましょう。

これまでに新しいことは何もないので、練習のために自分で書いてみてください。

## 最初にテストを書く

```go
package iteration

import "testing"

func TestRepeat(t *testing.T) {
    repeated := Repeat("a")
    expected := "aaaaa"

    if repeated != expected {
        t.Errorf("expected %q but got %q", expected, repeated)
    }
}
```

## テストを試して実行する

`./repeat_test.go:6:14: undefined: Repeat`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

規律を守りましょう！ テストを適切に失敗させるために、今すぐ新しいことを知る必要はありません。

あなたが今する必要があるのはそれをコンパイルするのに十分なので、あなたのテストがうまく書かれていることを確認することができます。

```go
package iteration

func Repeat(character string) string {
    return ""
}
```

いくつかの基本的な問題のテストを書くのに、十分なGoをすでに知っているのはいいことではありませんか？ つまり、これで本番コードを好きなだけ動かしても、希望どおりに動作していることがわかるようになります。

`repeat_test.go:10: expected 'aaaaa' but got ''`

## 成功させるのに十分なコードを書く

`for`構文は非常に目立たず、ほとんどのC言語のような言語に従います。

```go
func Repeat(character string) string {
    var repeated string
    for i := 0; i < 5; i++ {
        repeated = repeated + character
    }
    return repeated
}
```

C、Java、JavaScriptなどの他の言語とは異なり、forステートメントの3つのコンポーネントを囲む括弧はなく、中括弧 `{ }` は常に必要です。 行で何が起こっているのか不思議に思うかもしれません。

```go
    var repeated string
```

これまで変数の宣言と初期化に `:=`を使用してきました。ただし、 `:=`は単に [変数代入の省略形](https://gobyexample.com/variables)です。 ここでは `string`変数のみを宣言しています。したがって、明示的なバージョンです。後で説明するように、`var`を使用して関数を宣言することもできます。

テストを実行すれば合格です。

forループのその他のバリアントについては、 [こちら](https://gobyexample.com/for)をご覧ください。

## リファクタリング

次に、リファクタリングして、別の構成体 `+=`代入演算子を導入します。

```go
const repeatCount = 5

func Repeat(character string) string {
    var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

`+=` は `加算代入演算子`と呼ばれ、右のオペランドを左のオペランドに追加し、結果を左のオペランドに割り当てます。 整数のような他の型で動作します。

### ベンチマーク

Goでの [ベンチマーク](https://golang.org/pkg/testing/#hdr-Benchmarks) の記述は、言語のもう1つの優れた機能であり、テストの記述とよく似ています。

```go
func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repeat("a")
    }
}
```

コードがテストに非常に似ていることがわかります。

`testing.B`は、暗号的に命名された`b.N`にアクセスできるようになります。

ベンチマークコードが実行されると、`b.N`回実行され、かかる時間を測定します。

コードが実行される回数は重要ではありません。フレームワークは、適切な結果を得るために、その`適切な`値を決定します。

ベンチマークを実行するには、`go test -bench=.` を実行します （ Windows Powershellを使用している場合は、 `go test -bench="."`）

```text
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000           136 ns/op
PASS
```

`136 ns/op`が意味することは、関数がコンピュータ上を実行するのに平均で`136ナノ秒`かかることです。かなり大丈夫です！ これをテストするために、10000000回実行しました。

デフォルトでは、ベンチマークは順次実行されます。

## 練習問題

* テストを変更して、発信者が文字が繰り返される回数を指定し、コードを修正できるようにする
* 関数をドキュメント化するために `ExampleRepeat` を記述します
* [strings](https://golang.org/pkg/strings) パッケージをご覧ください。便利だと思われる関数を見つけて、ここにあるようなテストを作成して実験してください。標準ライブラリの学習に時間を費やすことは、時間の経過とともに本当に実を結びます。

## まとめ

* より多くのテスト駆動開発（TDD）プラクティス
* `for`文の学び
* ベンチマークの書き方を学びました

