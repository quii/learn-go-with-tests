---
description: 反復、繰り返し
---

# 反復、繰り返し

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/for)**

Goで繰り返し作業を行うには、 `for`が必要です。 Goには `while`、`do`、 `until`キーワードはなく、`for`のみ使用できます。これは良いことです！

文字を5回繰り返す関数のテストを書いてみましょう。

ここまでは新しいことは何もないので、練習のために自分で書いてみてください。

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

## テストを実行してみます

`./repeat_test.go:6:14: undefined: Repeat`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

_規律を守りましょう！_ テストを適切に失敗させるために、今すぐ新しいことを知る必要はありません。

今必要なのは、コンパイルできるようにすることだけです。そうすれば、テストがうまく書かれていることを確認できます。

```go
package iteration

func Repeat(character string) string {
	return ""
}
```

いくつかの基本的な問題のテストを書くのに、十分なGoをすでに知っているのはいいことではありませんか？ つまり、これからは本番コードを好きなだけいじっても、期待どおりに動作しているかどうかを確かめられるということです。

`repeat_test.go:10: expected 'aaaaa' but got ''`

## 成功させるのに十分なコードを書く

`for`の構文はごく平凡で、ほとんどのC系の言語と同じです。

```go
func Repeat(character string) string {
	var repeated string
	for i := 0; i < 5; i++ {
		repeated = repeated + character
	}
	return repeated
}
```

C、Java、JavaScriptなどの他の言語とは異なり、forステートメントの3つのコンポーネントを囲む括弧はなく、中括弧 `{ }` は常に必要です。 次の行で何が起こっているのか不思議に思うかもしれません。

```go
	var repeated string
```

これまで変数の宣言と初期化に `:=`を使用してきました。ただし、 `:=`は単に [両方のステップの省略形](https://gobyexample.com/variables)です。 ここでは `string`変数の宣言のみを行っています。したがって、明示的なバージョンを使っています。後で説明するように、`var`を使用して関数を宣言することもできます。

テストを実行すれば合格です。

forループのその他のバリアントについては、 [こちら](https://gobyexample.com/for)をご覧ください。

## リファクタリング♪

次はリファクタリングをして、新しい構文である `+=` 代入演算子を導入しましょう。

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

`+=` は _「加算代入演算子」_ と呼ばれ、右のオペランドを左のオペランドに加算し、結果を左のオペランドに代入します。 整数のような他の型でも動作します。

### ベンチマーク

Goでの [ベンチマーク](https://golang.org/pkg/testing/#hdr-Benchmarks) の記述は、テストと同様に言語に最初から組み込まれた（ファーストクラスの）機能であり、テストの記述とよく似ています。

```go
func BenchmarkRepeat(b *testing.B) {
	for b.Loop() {
		Repeat("a")
	}
}
```

コードがテストに非常に似ていることがわかります。

`testing.B`を使うと、ループ関数にアクセスできます。`Loop()` は、ベンチマークを続けて実行すべき間は true を返します。

ベンチマークコードが実行されると、かかる時間を測定します。`Loop()` が false を返した後、`b.N` には実行された反復の総数が入っています。

コードが実行される回数を気にする必要はありません。フレームワークは、まともな結果を得るために「適切な」値を決定します。

ベンチマークを実行するには、`go test -bench=.` を実行します（Windows PowerShell を使用している場合は `go test -bench="."`）。

```text
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000           136 ns/op
PASS
```

`136 ns/op`が意味することは、関数が\(私のコンピュータ上で\)実行するのに平均で136ナノ秒かかることです。なかなかの結果です！ この計測のために、関数は10000000回実行されました。

**注:** デフォルトでは、ベンチマークは順次実行されます。

時間が計測されるのはループの本体だけです。セットアップやクリーンアップのコードは、自動的にベンチマークの計測時間から除外されます。典型的なベンチマークは次のような構造になります。

```go
func Benchmark(b *testing.B) {
	//... setup ...
	for b.Loop() {
		//... code to measure ...
	}
	//... cleanup ...
}
```

Goの文字列はイミュータブル（不変）です。つまり、`Repeat` 関数で行っているような連結のたびに、新しい文字列を収めるためのメモリのコピーが発生します。これは、特に大量の文字列連結を行う場合に、パフォーマンスに影響します。

標準ライブラリには、メモリのコピーを最小限に抑える `strings.Builder`[stringsBuilder] 型が用意されています。
これは文字列の連結に使える `WriteString` メソッドを実装しています。

```go
const repeatCount = 5

func Repeat(character string) string {
	var repeated strings.Builder
	for i := 0; i < repeatCount; i++ {
		repeated.WriteString(character)
	}
	return repeated.String()
}
```

**注**: 最終的な結果を取得するには `String` メソッドを呼び出す必要があります。

`BenchmarkRepeat` を使って、`strings.Builder` によってパフォーマンスが大幅に向上することを確認できます。
`go test -bench=. -benchmem` を実行します。

```text
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000           25.70 ns/op           8 B/op           1 allocs/op
PASS
```

`-benchmem` フラグは、メモリ割り当てに関する情報を報告します。

* `B/op`: 1回の反復あたりに割り当てられたバイト数
* `allocs/op`: 1回の反復あたりのメモリ割り当て回数

## 練習問題

* 呼び出し側が文字を繰り返す回数を指定できるようにテストを変更し、コードを修正する
* 関数をドキュメント化するために `ExampleRepeat` を記述します
* [strings](https://golang.org/pkg/strings) パッケージをご覧ください。便利だと思われる関数を見つけて、ここにあるようなテストを作成して実験してください。標準ライブラリの学習に時間を投資すれば、長い目で見て必ず報われます。

## まとめ

* TDDのさらなる練習
* `for` を学びました
* ベンチマークの書き方を学びました

[stringsBuilder]: https://pkg.go.dev/strings#Builder
