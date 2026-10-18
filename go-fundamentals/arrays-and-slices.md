---
description: 配列とスライス
---

# 配列とスライス

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/arrays)**

配列を使うと、同じ型の複数の要素を特定の順番で変数に格納することができます。

配列を持っていると、それらの要素を繰り返し処理しなければならないことがよくあります。そこで、[新しく知った`for`の知識](iteration.md) を使って`Sum`関数を作ってみましょう。`Sum`は数値の配列を受け取り、その合計を返します。

テスト駆動開発（TDD）の技術を使ってみましょう

## 最初にテストを書く

作業用の新しいフォルダを作成します。`sum_test.go` という名前の新しいファイルを作成し、次の内容を記述します。

```go
package main

import "testing"

func TestSum(t *testing.T) {

	numbers := [5]int{1, 2, 3, 4, 5}

	got := Sum(numbers)
	want := 15

	if got != want {
		t.Errorf("got %d want %d given, %v", got, want, numbers)
	}
}
```

配列には _固定容量_ があり、これは変数を宣言する際に定義します。配列を初期化するには、以下の2つの方法があります。

* \[N\]type{value1, value2, ..., valueN} 例: `numbers := [5]int{1, 2, 3, 4, 5}`
* \[...\]type{value1, value2, ..., valueN} 例: `numbers := [...]int{1, 2, 3, 4, 5}`

エラーメッセージに関数への入力も表示すると便利な場合があります。
ここでは、`%v` プレースホルダを使用して「デフォルト」の書式で出力しています。これは配列に対してうまく機能します。

[文字列の書式についての詳細はこちら](https://golang.org/pkg/fmt/)

## テストを実行してみます

`go mod init main` で go mod を初期化していた場合、`_testmain.go:13:2: cannot import "main"` というエラーが表示されます。これは、一般的な慣習として、mainパッケージには他のパッケージを統合するコードだけが含まれ、ユニットテスト可能なコードは含まれないためです。そのため、Goでは `main` という名前のパッケージをインポートすることができません。

これを修正するには、`go.mod` のメインモジュールの名前を他の任意の名前に変更します。

上記のエラーを修正したら、`go test`を実行すると、おなじみの `./sum_test.go:10:15: undefined: Sum` というエラーでコンパイラが失敗します。これで、テスト対象の実際のメソッドを書き始めることができます。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`sum.go` に次を記述します。

```go
package main

func Sum(numbers [5]int) int {
	return 0
}
```

これで、テストは _明確なエラーメッセージ_ を表示して失敗するはずです。

`sum_test.go:13: got 0 want 15 given, [1 2 3 4 5]`

## 成功させるのに十分なコードを書く

```go
func Sum(numbers [5]int) int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += numbers[i]
	}
	return sum
}
```

特定のインデックスの配列から値を取り出すには、`array[index]`構文を使えば大丈夫です。この例では`for`を使って配列を5回繰り返し、各項目を`sum`に加算しています。

## リファクタリング♪

コードをきれいにするために[`range`](https://gobyexample.com/range)を導入してみましょう。

```go
func Sum(numbers [5]int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```

`range`を使うと配列の反復処理を行うことができます。各反復で、`range`はインデックスと値の2つの値を返します。
ここでは `_` [ブランク識別子](https://golang.org/doc/effective_go.html#blank)を使用して、インデックスの値を無視しています。

### 配列とその型

配列の興味深い特性として、サイズが型でエンコードされていることが挙げられます。 もし `[5]int` を期待する関数に `[4]int` を渡そうとしてもコンパイルできません。 これらは異なる型なので、`int` を求める関数に `string` を渡そうとするのと同じです。

配列の長さが固定されているのは非常に面倒だと思うかもしれませんし、ほとんどの場合、配列を使うことはないでしょう。

Goには _スライス_ がありますが、これはコレクションのサイズをエンコードするものではなく、任意のサイズを持つことができます。

次の要件は、さまざまなサイズのコレクションを合計することです。

## 最初にテストを書く

ここでは、任意のサイズのコレクションを持つことができる [スライス型][slice] を使用します。構文は配列と非常に似ていますが、宣言時にサイズを省略するだけです。

`myArray := [3]int{1,2,3}`ではなく、`mySlice := []int{1,2,3}`のように書きます。

```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := [5]int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```

## テストを実行してみます

これはコンパイルされません

`./sum_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Sum`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

ここでの問題は、次のどちらかを選ぶ必要があることです。

* `Sum`の引数を配列ではなくスライスに変更して、既存のAPIを壊す。これを行うと、 _もう一方の_ テストがコンパイルできなくなるため、誰かの一日を台無しにしてしまう可能性があります！
* 新しい関数を作成する

私たちの場合は誰もこの関数を使っていないので、二つの関数をメンテナンスするよりも一つの関数を持つことにしましょう。

```go
func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```

テストを実行しようとしてもコンパイルされないので、最初のテストを配列ではなくスライスで渡すように変更しなければなりません。

## 成功させるのに十分なコードを書く

ここでは、コンパイラの問題を修正するだけで十分で、テストはパスすることがわかりました！

## リファクタリング♪

すでに`Sum`をリファクタリングしています。配列をスライスに置き換えただけなので、追加の変更は必要ありません。
リファクタリングの段階でテストコードをおろそかにしてはいけないことを覚えておいてください。`Sum`のテストをさらに改善できます。

```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```

テストの価値を問うことは重要です。可能な限り多くのテストを行うことが目標ではなく、むしろコードベースに対して可能な限り多くの _信頼_ を持つことが目標です。あまりにも多くのテストを持っていると、実際に問題になることがありますし、メンテナンスのオーバーヘッドを増やすだけです。 **すべてのテストにはコストがあります**。

私たちのケースでは、この関数のために2つのテストを持つことが冗長であることがわかります。あるサイズのスライスで動作するなら、\(常識の範囲内で\)どのサイズのスライスでも動作する可能性が高いです。

Goの組み込みテストツールキットには、 [カバレッジツール](https://blog.golang.org/cover)があります。
100%のカバレッジを目指すことを最終目標にすべきではありませんが、カバレッジツールはテストでカバーされていないコードの領域を特定するのに役立ちます。TDD を厳格に行っていれば、いずれにせよ 100% に近いカバレッジが得られる可能性が高いでしょう。

実行してみてください。

`go test -cover`

次のように表示されるはずです。

```bash
PASS
coverage: 100.0% of statements
```

テストの一つを削除して、カバレッジをもう一度確認してください。

十分にテストされた関数ができて満足したところで、次の課題に取り組む前に、素晴らしい成果をコミットしておきましょう。

`SumAll`という新しい関数が必要です。これは可変個のスライスを受け取り、渡された各スライスの合計を含む新しいスライスを返します。

例えば、以下のようになります。

`SumAll([]int{1,2}, []int{0,9})` は `[]int{3, 9}`を返します。

もしくは

`SumAll([]int{1,1,1})` は `[]int{3}`を返します。

## 最初にテストを書く

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## テストを実行してみます

`./sum_test.go:23:9: undefined: SumAll`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

テストが求めるとおりに`SumAll`を定義する必要があります。

Goを使えば、可変数の引数を取ることができる [_可変長引数関数_](https://gobyexample.com/variadic-functions) を書くことができます。

```go
func SumAll(numbersToSum ...[]int) []int {
	return nil
}
```

これは正しいコードですが、テストはまだコンパイルされません！

`./sum_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

Goでは、スライスで等号演算子を使うことはできません。`got` と `want` の各スライスを繰り返し処理して値を確認する関数を書くことも _できます_ が、もっと便利な方法があったらどうでしょうか？

Go 1.21 からは標準パッケージ [slices](https://pkg.go.dev/slices#pkg-overview) が利用できます。このパッケージには、スライスを単純に浅く比較する [slices.Equal](https://pkg.go.dev/slices#Equal) 関数があり、上記のように型について気にする必要はありません。
この関数は、要素が [comparable](https://pkg.go.dev/builtin#comparable)（比較可能）であることを前提としていることに注意してください。
そのため、2次元スライスのような比較できない要素を持つスライスには適用できません。

早速実践してみましょう！

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

以下のようなテスト出力が得られるはずです。
`sum_test.go:30: got [] want [3 9]`

## 成功させるのに十分なコードを書く

必要なのは可変長引数を繰り返し処理して、既存の`Sum`関数を使って合計を計算し、それを返すスライスに追加することです。

```go
func SumAll(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)
	sums := make([]int, lengthOfNumbers)

	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}

	return sums
}
```

新しい学びがたくさんあります!

スライスを作成する新しい方法があります。`make` を使うと、処理する必要がある`numbersToSum`の `len` を開始容量とするスライスを作成できます。スライスの長さ\(length\)はスライスが保持している要素の数 `len(mySlice)` であり、容量\(capacity\)は基になる配列に保持できる要素の数 `cap(mySlice)` です。例えば、`make([]int, 0, 5)` は長さ0、容量5のスライスを作成します。

配列のように `mySlice[N]` でスライスをインデックス化して値を取得したり、`=` で新しい値を代入したりすることができます。

これでテストは合格するはずです。

## リファクタリング♪

前述したように、スライスには容量があります。容量2のスライスを持っていて `mySlice[10] = 1` を実行しようとすると、 _ランタイム_ エラーが発生します。

しかし、`append`関数を使うことができます。この関数はスライスと新しい値を受け取り、すべての要素を含む新しいスライスを返します。

```go
func SumAll(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}
```

この実装では、容量についてはあまり気にしていません。空のスライス `sums` から始め、 可変長引数を処理しながら `Sum` の結果をそれに追加します。

次の要件は `SumAll` を `SumAllTails` に変更することです。これは各スライスの「末尾\(tail\)」の合計を計算します。コレクションの末尾とは、最初の要素\(「先頭\(head\)」\)を除いたすべての要素のことです。

## 最初にテストを書く

```go
func TestSumAllTails(t *testing.T) {
	got := SumAllTails([]int{1, 2}, []int{0, 9})
	want := []int{2, 9}

	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## テストを実行してみます

`./sum_test.go:26:9: undefined: SumAllTails`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

関数名を `SumAllTails` に変更し、テストを再実行します。

`sum_test.go:30: got [3 9] want [2 9]`

## 成功させるのに十分なコードを書く

```go
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		tail := numbers[1:]
		sums = append(sums, Sum(tail))
	}

	return sums
}
```

スライスはスライスすることができます！構文は `slice[low:high]` で、`:` の片側の値を省略すると、その側の端までのすべてを取り出します。この例では、`numbers[1:]` を使って「1から最後まで取る」と言っています。スライスを使った他のテストを書いたり、スライス演算子に慣れるために実験をしたりすることに時間を使うとよいかもしれません。

## リファクタリング♪

今回はリファクタリングすることはあまりありません。

空のスライスを関数に渡すとどうなると思いますか？空のスライスの「末尾」とは何でしょうか？Goに `myEmptySlice[1:]` からすべての要素をキャプチャするように指示するとどうなるでしょうか？

## 最初にテストを書く

```go
func TestSumAllTails(t *testing.T) {

	t.Run("make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}

		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}

		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

}
```

## テストを実行してみます

```text
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range
```

なんてことだ！テストは _コンパイルできた_ ものの、 _ランタイムエラーが発生している_ ことに注意してください。

コンパイル時のエラーは、動作するソフトウェアを書くのに役立つので私たちの味方ですが、ランタイムエラーはユーザーに影響を与えるので私たちの敵です。

## 成功させるのに十分なコードを書く

```go
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
}
```

## リファクタリング♪

このテストでは、再びアサーションの周りに繰り返しコードがあるので、それを関数に抽出してみましょう。

```go
func TestSumAllTails(t *testing.T) {

	checkSums := func(t *testing.T, got, want []int) {
		t.Helper()
		if !slices.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("make the sums of tails of", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})

}
```

いつものように新しい関数 `checkSums` を作成することもできましたが、ここでは新しいテクニック、つまり関数を変数に代入する方法を紹介しています。奇妙に見えるかもしれませんが、`string` や `int` を変数に代入するのと何も変わりません。関数も実質的には値なのです。

ここでは示していませんが、このテクニックは、関数を「スコープ」\(例えば `{}` の間\)内の他のローカル変数に結びつけたい場合に便利です。また、APIの表面積を減らすこともできます。

この関数をテストの中で定義することで、このパッケージ内の他の関数からは使用できなくなります。エクスポートする必要のない変数や関数を隠すことは、設計上の重要な考慮事項です。

これによる便利な副作用として、コードに少しだけ型安全性が加わります。開発者が誤って `checkSums(t, got, "dave")` で新しいテストを追加しても、コンパイラはその場で止めてくれます。

```bash
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## まとめ

ここで学んだこと

* 配列
* スライス
  * 作り方いろいろ
  * _固定_ 容量を持っていますが、`append` を使って古いスライスから新しいスライスを作成することができます
  * スライスをスライスする方法！
* 配列やスライスの長さを取得するために `len` を使用します。
* テストカバレッジツール
* `slices.Equal` と、なぜ通常の等号演算子の代わりにそれが必要なのか

ここまでは整数のスライスや配列を使ってきましたが、配列やスライス自体を含め、他の型でも動作します。ですから、必要に応じて `[][]string` の変数を宣言することができます。

スライスの詳細については、[スライスに関するGoブログ記事][blog-slice]を参照してください。これを読んで学んだことを定着させるために、より多くのテストを書いてみてください。

テストを書く以外にGoを使って実験するもう一つの便利な方法は、Go playgroundです。ほとんどのことを試すことができますし、質問したいときには簡単にコードを共有することができます。 [Go playgroundにスライスを入れて実験できるようにしてみました。](https://play.golang.org/p/ICCWcRGIO68)

[配列をスライスした例](https://play.golang.org/p/bTrRmYfNYCp) では、配列をスライスして、スライスを変更すると元の配列にどのように影響するかを説明していますが、スライスの「コピー」は元の配列には影響しません。
[別の例](https://play.golang.org/p/Poth8JS28sc) では、非常に大きなスライスをスライスした後にそのコピーを作るのが良い理由を示しています。

[for]: iteration.md
[blog-slice]: https://blog.golang.org/go-slices-usage-and-internals
[deepEqual]: https://golang.org/pkg/reflect/#DeepEqual
[slice]: https://golang.org/doc/effective_go.html#slices
