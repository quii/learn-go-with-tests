---
description: Arrays and slices
---

# 配列とスライス

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/arrays)

配列を使うと、同じ型の複数の要素を特定の順番で変数に格納することができます。

配列を持っていると、それらの要素を繰り返し処理しなければならないことがよくあります。そこで、[新しく知った`for`の知識](iteration.md) を使って`Sum`関数を作ってみましょう。`Sum`は数値の配列を受け取り、その合計を返します。

テスト駆動開発（TDD）の技術を使ってみましょう

## 最初にテストを書く

In `sum_test.go`

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

配列には`固定容量`があり、これは変数を宣言する際に定義します。配列を初期化するには、以下の2つの方法があります。

* `numbers := [5]int{1, 2, 3, 4, 5}`
* `numbers := [...]int{1, 2, 3, 4, 5}`

エラーメッセージに関数への入力を表示するのも便利な場合がありますが、ここでは `%v` プレースホルダを使用しています。

[文字列の書式についての詳細はこちら](https://golang.org/pkg/fmt/)

## テストを実行してみてください

`go test`を実行すると、コンパイラは失敗します。 `./sum_test.go:10:15: undefined: Sum`

## テストを実行するための最小限のコードを書き、失敗したテストの出力をチェックする

In `sum.go`

```go
package main

func Sum(numbers [5]int) int {
    return 0
}
```

これで、テストは明確なエラーメッセージが表示されて失敗するはずです。

`sum_test.go:13: got 0 want 15 given, [1 2 3 4 5]`

## テストをパスするのに十分なコードを書く

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

## リファクタリング

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

`range`は配列の反復処理を行うことができる。呼び出されるたびにインデックスと値の2つの値を返します。ここではインデックスの値を無視して `_` [空（スペース）の識別子](https://golang.org/doc/effective_go.html#blank)を使用しています。

### 配列とその型

配列の興味深い特性として、サイズが型でエンコードされていることが挙げられます。
もし `[5]int` を期待する関数に `[4]int` を渡そうとしてもコンパイルできません。
これらは異なる型なので、`int` を求める関数に `string` を渡そうとするのと同じです。

配列の長さが固定されているのは非常に面倒だと思うかもしれませんし、ほとんどの場合、配列を使うことはないでしょう。

Goには`slices`がありますが、これはコレクションのサイズをエンコードするものではなく、任意のサイズを持つことができます。

次の要件は、さまざまなサイズのコレクションを合計することです。

## 最初にテストを書く

ここでは、任意のサイズのコレクションを持つことができる [slice type](https://golang.org/doc/effective_go.html#slices) を使用します。構文は配列と非常に似ていますが、宣言時にサイズを省略するだけです。

`myArray := [3]int{1,2,3}`ではなく、`mySlice := []int{1,2,3}`こちらです。

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

## テストを実行してみてください

これはコンパイルされません

`./sum_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Sum`

## テストを実行するための最小限のコードを書き、失敗したテストの出力をチェックする

ここで問題なのは

* 引数`Sum`をスライスに変更することで、既存のAPIを壊します。

  配列に対して、このようなことをすると、破壊された可能性を潜在的発見するので他のテストがコンパイルされません!

* 新しい関数を作成しましょう

私たちの場合は誰もこの関数を使っていないので、二つの関数を持つよりも一つの関数を持つことにしましょう。

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

## テストをパスさせるのに十分なコードを書いてください

コンパイラの問題を修正するだけで、次のことができることがわかりました。

## リファクタリング

すでに`Sum`をリファクタリングしており、配列からスライスに変更しただけなので、ここでやることはそれほど多くありません。リファクタリングの段階でテストコードをおろそかにしてはいけないことを覚えておいてください。

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

テストの価値を問うことは重要です。可能な限り多くのテストを行うことが目標ではなく、むしろコードベースに対して可能な限り多くの`信頼`を持つことが目標です。あまりにも多くのテストを持っていると、実際に問題になることがありますし、メンテナンスのオーバーヘッドを増やすだけです。 **すべてのテストにはコストがある**.

私たちのケースでは、この関数のために2つのテストを持つことが冗長であることがわかります。あるサイズのスライスで動作するなら、どのサイズのスライスでも動作する可能性が高いです。

Goの組み込みテストツールキットには、 [カバレッジツール](https://blog.golang.org/cover)があり、あなたがカバーしていないコードの領域を特定するのに役立ちます。私が強調したいのは、100%のカバレッジを持つことがゴールではないということです。TDD を厳格に行っていれば、100% に近いカバレッジが得られる可能性が高いでしょう。

実行してみてください。

`go test -cover`

見てください。

```bash
PASS
coverage: 100.0% of statements
```

テストの一つを削除して、カバレッジをもう一度確認してください。

これで、十分にテストされた関数を手に入れたことに満足しているので、次の課題に挑戦する前に、あなたの素晴らしい作品をコミットしてください。

次に様々なスライスの数を受け取り、渡された各スライスの `SumAll` を含む新しいスライスを返します。

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

## テストを実行してみてください

`./sum_test.go:23:9: undefined: SumAll`

## テストを実行するための最小限のコードを書き、失敗したテストの出力をチェックする

テストの目的に応じて`SumAll`を定義する必要があります。

Goを使えば、可変数の引数を取ることができる [可変関数](https://gobyexample.com/variadic-functions) を書くことができます。

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
    return
}
```

コンパイルしようとしても、テストはまだコンパイルされません。

`./sum_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

Goでは、スライスで等号演算子を使うことはできません。`got` と `want` の各スライスを繰り返し処理して値を確認する関数を書くこともできますが、便利のために [`reflect.DeepEqual`](https://golang.org/pkg/reflect/#DeepEqual) を使うと、2つの変数が同じであるかどうかを確認するのに便利です。

```go
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1, 2}, []int{0, 9})
    want := []int{3, 9}

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

（`DeepEqual`にアクセスできるようにするには、ファイルの先頭で`import reflect`を作成します）

重要なのは、`reflect.DeepEqual`は「型安全」ではないことに注意することです。これを確認するには、テストを一時的に変更してください。

```go
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1, 2}, []int{0, 9})
    want := "bob"

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

ここで行ったことは、`slice` と `string` を比較しようとしていることです。これでは意味がありませんが、テストはコンパイルされます。だから、`reflect.DeepEqual`を使うのはスライスを比較するのに便利な方法ですが、使うときは注意が必要です。

テストを再度変更して実行すると、以下のようなテスト出力が得られるはずです。

`sum_test.go:30: got [] want [3 9]`

## テストをパスするのに十分なコードを書く

必要なのは可変長引数を繰り返し処理して、前に使った`Sum`関数を使って合計を計算し、それを返すスライスに追加することです。

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

スライスを作成する新しい方法があります。`make` を使うと、`numbersToSum`の開始容量が `len`であるスライスを作成できるようになります。

配列のように `mySlice[N]` でスライスをインデックス化して値を取得したり、`=` で新しい値を代入したりすることができます。

これでテストは合格するはずです。

## リファクタリング

前述したように、スライスには容量があります。容量2のスライスを持っていて `mySlice[10] = 1` を実行しようとすると、ランタイムエラーが発生します。

しかし、`append`関数を使えば、スライスと新しい値を受け取り、その中にあるすべての項目を含む新しいスライスを返すことができます。

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

次の要件は `SumAll` を `SumAllTails` に変更することです。コレクションの末尾とは、最初のものを除いたすべてのアイテムのことです。

## 最初にテストを書く

```go
func TestSumAllTails(t *testing.T) {
    got := SumAllTails([]int{1, 2}, []int{0, 9})
    want := []int{2, 9}

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

## テストを実行してみてください

`./sum_test.go:26:9: undefined: SumAllTails`

## テストを実行するための最小限のコードを書き、失敗したテストの出力をチェックする

関数名を `SumAllTails` に変更し、テストを再実行します。

`sum_test.go:30: got [3 9] want [2 9]`

## テストがパスするのに十分なコードを書く

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

スライスはスライスすることができます。構文は `slice[low:high]` で、`:` の片方の辺の値を省略すると、その辺までのすべての値をキャプチャします。この例では、`numbers[1:]` を使って「1から最後まで取る」と言っています。スライスを使った他のテストを書いたり、スライス演算子に慣れるために実験をしたりすることに時間を投資したほうがいいかもしれません。

## リファクタリング

今回はリファクタリングすることはあまりありません。

空のスライスを関数に渡すとどうなると思いますか？

空のスライスの「末尾」とは何ですか？

Goに `myEmptySlice[1:]` からすべての要素をキャプチャするように指示するとどうなるか?

## 最初にテストを書く

```go
func TestSumAllTails(t *testing.T) {

    t.Run("make the sums of some slices", func(t *testing.T) {
        got := SumAllTails([]int{1, 2}, []int{0, 9})
        want := []int{2, 9}

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    })

    t.Run("safely sum empty slices", func(t *testing.T) {
        got := SumAllTails([]int{}, []int{3, 4, 5})
        want := []int{0, 9}

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    })

}
```

## テストを実行してみてください

```text
panic: runtime error: slice bounds out of range [recovered]
panic: runtime error: slice bounds out of range
```

これはランタイムエラーです。
コンパイル時のエラーは、動作するソフトウェアを書くのに役立ちますが、ランタイムエラーはユーザーに影響を与えます。

## テストをパスするのに十分なコードを書く

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

## リファクタリング

このテストでは、アサーションの周りに繰り返しコードがあるので、それを関数に抽出してみましょう。

```go
func TestSumAllTails(t *testing.T) {

    checkSums := func(t *testing.T, got, want []int) {
        t.Helper()
        if !reflect.DeepEqual(got, want) {
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

これは便利な副作用として、コードに少しだけ型の安全性が追加されます。
アホな開発者が `checkSums(t, got, "dave")` で新しいテストを追加しても、コンパイラはその場で止めてくれます。

```bash
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## まとめ

ここで学んだこと

* 配列
* スライス
  * 作り方いろいろ
  * 固定容量を持っていますが、古いものから新しいスライスを作成することができます。

    `append`を使えば追加できます

  * スライスからスライス
* 配列やスライスの長さを取得するために `len` を使用します。
* テストカバレッジツール
* `reflect.DeepEqual` と、なぜそれが便利なのか、しかしコードの型安全性を低下させる可能性があるのか。

ここまでは整数のスライスや配列を使ってきましたが、配列やスライス自体を含め、他の型でも動作します。ですから、必要に応じて `[][]string` の変数を宣言することができます。

スライスの詳細については、[スライスに関するGoブログ記事](https://blog.golang.org/go-slices-usage-and-internals)を参照してください。これを読んで学んだことを実証するために、より多くのテストを書いてみてください。

テストを書く以外にGoを使って実験するもう一つの便利な方法は、Goの遊び場です。ほとんどのことを試すことができますし、質問が必要な場合は簡単にコードを共有することができます。 [Go playgroundにスライスを入れて実験できるようにしてみました。](https://play.golang.org/p/ICCWcRGIO68)

[配列をスライスした例](https://play.golang.org/p/bTrRmYfNYCp) では、配列をスライスして、スライスを変更すると元の配列にどのように影響するかを説明していますが、スライスの「コピー」は元の配列には影響しません。

[別の例](https://play.golang.org/p/Poth8JS28sc) 非常に大きなスライスをスライスした後にコピーを作るのが良い理由があります。
