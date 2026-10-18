---
description: Revisiting arrays and slices with generics
---

# ジェネリクスで配列とスライスを再訪する

**[この章のコードは「配列とスライス」の続きで、ここにあります](https://github.com/quii/learn-go-with-tests/tree/main/arrays)**

[配列とスライス](arrays-and-slices.md)で書いた `SumAll` と `SumAllTails` の両方を見てみましょう。手元に自分のバージョンがない場合は、[配列とスライス](arrays-and-slices.md)の章からテストと一緒にコードをコピーしてください。

```go
// Sum calculates the total from a slice of numbers.
func Sum(numbers []int) int {
	var sum int
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// SumAllTails calculates the sums of all but the first number given a collection of slices.
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

繰り返し現れるパターンに気づきましたか？

- 何らかの「初期」結果値を作成する。
- コレクションを反復処理し、結果とスライスの次の要素に何らかの演算（または関数）を適用して、結果に新しい値を設定する。
- 結果を返す。

この考え方は関数型プログラミングの世界でよく語られるもので、多くの場合「reduce」や [fold](https://en.wikipedia.org/wiki/Fold_(higher-order_function)) と呼ばれます。

> In functional programming, fold (also termed reduce, accumulate, aggregate, compress, or inject) refers to a family of higher-order functions that analyze a recursive data structure and through use of a given combining operation, recombine the results of recursively processing its constituent parts, building up a return value. Typically, a fold is presented with a combining function, a top node of a data structure, and possibly some default values to be used under certain conditions. The fold then proceeds to combine elements of the data structure's hierarchy, using the function in a systematic way.
>
> （関数型プログラミングにおいて、fold（reduce、accumulate、aggregate、compress、inject とも呼ばれる）は、再帰的なデータ構造を解析し、与えられた結合演算を用いて、その構成要素を再帰的に処理した結果を再結合し、戻り値を構築する高階関数の一群を指します。通常、fold には結合関数、データ構造の最上位ノード、そして場合によっては特定の条件下で使われるデフォルト値が与えられます。そして fold は、その関数を体系的に用いて、データ構造の階層の要素を結合していきます。）

Go には以前から高階関数があり、バージョン 1.18 からは[ジェネリクス](./generics.md)もあるので、より広い分野で議論されているこうした関数のいくつかを定義できるようになりました。現実から目を背けても仕方ありません。これは Go のエコシステムの外では非常によく使われる抽象化であり、理解しておくと役に立ちます。

さて、これを見て顔をしかめている人もいるでしょう。

> Go はシンプルであるべきだ

**簡単さ（easiness）とシンプルさ（simplicity）を混同しないでください**。ループを書いたりコードをコピー＆ペーストしたりするのは簡単ですが、必ずしもシンプルではありません。シンプルと簡単の違いについて詳しくは、[Rich Hickey の傑作トーク「Simple Made Easy」](https://www.youtube.com/watch?v=SxdOUGdseq4)を観てください。

**不慣れさと複雑さを混同しないでください**。fold/reduce は最初は恐ろしく、いかにもコンピューターサイエンス的に聞こえるかもしれませんが、実際には非常によくある操作を抽象化したものにすぎません。コレクションを受け取り、それを 1 つの要素にまとめる、というものです。一歩引いて考えてみれば、おそらくあなたはこれを _頻繁に_ 行っていることに気づくでしょう。

## ジェネリックなリファクタリング

ピカピカの新しい言語機能について人々がよく犯す間違いは、具体的なユースケースがないまま使い始めてしまうことです。推測に頼って取り組みを進めてしまうのです。

ありがたいことに、私たちは「役に立つ」関数を書き、その周りにテストを用意しているので、TDD のリファクタリング段階で自由にアイデアを試すことができます。そして何を試しても、その価値はユニットテストによって検証されます。

リファクタリングのステップでコードをシンプルにするツールとしてジェネリクスを使うほうが、時期尚早な抽象化に陥るよりも、有用な改善へと導かれる可能性がはるかに高くなります。

私たちは安心していろいろ試し、テストを再実行し、変更が気に入ればコミットできます。気に入らなければ、変更を元に戻せばいいだけです。この実験の自由こそが、TDD の本当に大きな価値の 1 つです。

[前の章](generics.md)でジェネリクスの構文には慣れているはずなので、自分で `Reduce` 関数を書いて、`Sum` と `SumAllTails` の中で使ってみてください。

### ヒント

まず関数の引数について考えれば、有効な解決策はごく少数に絞られるはずです。
  - reduce したい配列
  - 何らかの結合関数

「Reduce」は非常によく文書化されたパターンなので、車輪の再発明をする必要はありません。[wiki、特にリストのセクションを読んでください](https://en.wikipedia.org/wiki/Fold_(higher-order_function))。必要になるもう 1 つの引数に気づくはずです。

> In practice, it is convenient and natural to have an initial value
>
> （実際には、初期値があると便利で自然です）

### `Reduce` の最初の実装

```go
func Reduce[A any](collection []A, f func(A, A) A, initialValue A) A {
	var result = initialValue
	for _, x := range collection {
		result = f(result, x)
	}
	return result
}
```

Reduce はこのパターンの _本質_ を捉えています。コレクション、累積関数、初期値を受け取り、単一の値を返す関数です。具体的な型にまつわるごちゃごちゃした雑音はありません。

ジェネリクスの構文を理解していれば、この関数が何をしているかを理解するのに問題はないはずです。`Reduce` という広く認知された用語を使うことで、他の言語のプログラマーにも意図が伝わります。

### 使い方

```go
// Sum calculates the total from a slice of numbers.
func Sum(numbers []int) int {
	add := func(acc, x int) int { return acc + x }
	return Reduce(numbers, add, 0)
}

// SumAllTails calculates the sums of all but the first number given a collection of slices.
func SumAllTails(numbers ...[]int) []int {
	sumTail := func(acc, x []int) []int {
		if len(x) == 0 {
			return append(acc, 0)
		} else {
			tail := x[1:]
			return append(acc, Sum(tail))
		}
	}

	return Reduce(numbers, sumTail, []int{})
}
```

`Sum` と `SumAllTails` は、それぞれの最初の行で宣言された関数によって、計算の振る舞いを記述するようになりました。コレクションに対して計算を実行する処理は `Reduce` の中に抽象化されています。

## reduce のさらなる応用

テストを使って reduce 関数をいろいろ試し、どれだけ再利用できるかを確かめられます。前の章のジェネリックなアサーション関数をコピーしてきました。

```go
func TestReduce(t *testing.T) {
	t.Run("multiplication of all elements", func(t *testing.T) {
		multiply := func(x, y int) int {
			return x * y
		}

		AssertEqual(t, Reduce([]int{1, 2, 3}, multiply, 1), 6)
	})

	t.Run("concatenate strings", func(t *testing.T) {
		concatenate := func(x, y string) string {
			return x + y
		}

		AssertEqual(t, Reduce([]string{"a", "b", "c"}, concatenate, ""), "abc")
	})
}
```

### ゼロ値

掛け算の例では、`Reduce` の引数としてデフォルト値を持つ理由を示しています。もし Go の `int` のデフォルト値である 0 に頼ると、初期値に 0 を掛け、その後の値にも掛けていくことになるので、結果は常に 0 になってしまいます。1 に設定することで、スライスの最初の要素はそのまま残り、残りは次の要素と掛け合わされていきます。

オタク仲間に賢そうに聞こえたいなら、これを[単位元（The Identity Element）](https://en.wikipedia.org/wiki/Identity_element)と呼ぶとよいでしょう。

> In mathematics, an identity element, or neutral element, of a binary operation operating on a set is an element of the set which leaves unchanged every element of the set when the operation is applied.
>
> （数学において、集合上の二項演算の単位元（または中立元）とは、その演算を適用しても集合のすべての要素を変化させない、その集合の要素のことです。）

足し算では、単位元は 0 です。

`1 + 0 = 1`

掛け算では、1 です。

`1 * 1 = 1`

## `A` とは異なる型に reduce したい場合は？

取引 `Transaction` のリストがあり、それと名前を受け取って銀行残高を計算する関数が欲しいとしましょう。

TDD のプロセスに従いましょう。

## 最初にテストを書く

```go
func TestBadBank(t *testing.T) {
	transactions := []Transaction{
		{
			From: "Chris",
			To:   "Riya",
			Sum:  100,
		},
		{
			From: "Adil",
			To:   "Chris",
			Sum:  25,
		},
	}

	AssertEqual(t, BalanceFor(transactions, "Riya"), 100)
	AssertEqual(t, BalanceFor(transactions, "Chris"), -75)
	AssertEqual(t, BalanceFor(transactions, "Adil"), -25)
}
```

## テストを実行してみます
```
# github.com/quii/learn-go-with-tests/arrays/v8 [github.com/quii/learn-go-with-tests/arrays/v8.test]
./bad_bank_test.go:6:20: undefined: Transaction
./bad_bank_test.go:18:14: undefined: BalanceFor
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

まだ型も関数もないので、テストが実行できるように追加します。

```go
type Transaction struct {
	From string
	To   string
	Sum  float64
}

func BalanceFor(transactions []Transaction, name string) float64 {
	return 0.0
}
```

テストを実行すると、次のように表示されるはずです。

```
=== RUN   TestBadBank
    bad_bank_test.go:19: got 0, want 100
    bad_bank_test.go:20: got 0, want -75
    bad_bank_test.go:21: got 0, want -25
--- FAIL: TestBadBank (0.00s)
```

## 成功させるのに十分なコードを書く

まずは `Reduce` 関数がないものとしてコードを書いてみましょう。

```go
func BalanceFor(transactions []Transaction, name string) float64 {
	var balance float64
	for _, t := range transactions {
		if t.From == name {
			balance -= t.Sum
		}
		if t.To == name {
			balance += t.Sum
		}
	}
	return balance
}
```

## リファクタリング♪

この時点で、ソース管理の規律を守って作業をコミットしましょう。動くソフトウェアができました。Monzo や Barclays などに挑む準備は万端です。

作業をコミットしたので、自由にいじったり、リファクタリングの段階でいくつかの異なるアイデアを試したりできます。正直なところ、今のコードは決して悪くはありませんが、この演習のために、同じコードを `Reduce` を使って書いてみせたいと思います。

```go
func BalanceFor(transactions []Transaction, name string) float64 {
	adjustBalance := func(currentBalance float64, t Transaction) float64 {
		if t.From == name {
			return currentBalance - t.Sum
		}
		if t.To == name {
			return currentBalance + t.Sum
		}
		return currentBalance
	}
	return Reduce(transactions, adjustBalance, 0.0)
}
```

しかし、これはコンパイルできません。

```
./bad_bank.go:19:35: type func(acc float64, t Transaction) float64 of adjustBalance does not match inferred type func(Transaction, Transaction) Transaction for func(A, A) A
```

その理由は、コレクションの型とは _異なる_ 型に reduce しようとしているからです。恐ろしく聞こえますが、実際には `Reduce` の型シグネチャを調整するだけで動くようになります。関数本体を変更する必要はなく、既存の呼び出し側を変更する必要もありません。

```go
func Reduce[A, B any](collection []A, f func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, x := range collection {
		result = f(result, x)
	}
	return result
}
```

2 つ目の型制約を追加したことで、`Reduce` の制約を緩めることができました。これにより、`A` のコレクションを `B` に `Reduce` できるようになります。今回の場合は `Transaction` から `float64` へです。

これで `Reduce` はより汎用的で再利用しやすくなり、それでいて型安全性も保たれています。もう一度テストを実行してみると、コンパイルが通り、テストも成功するはずです。

## 銀行を拡張する

遊びとして、銀行のコードの使い勝手を改善してみたくなりました。簡潔にするため、TDD のプロセスは省略しています。

```go
func TestBadBank(t *testing.T) {
	var (
		riya  = Account{Name: "Riya", Balance: 100}
		chris = Account{Name: "Chris", Balance: 75}
		adil  = Account{Name: "Adil", Balance: 200}

		transactions = []Transaction{
			NewTransaction(chris, riya, 100),
			NewTransaction(adil, chris, 25),
		}
	)

	newBalanceFor := func(account Account) float64 {
		return NewBalanceFor(account, transactions).Balance
	}

	AssertEqual(t, newBalanceFor(riya), 200)
	AssertEqual(t, newBalanceFor(chris), 0)
	AssertEqual(t, newBalanceFor(adil), 175)
}
```

そして、更新したコードがこちらです。

```go
package main

type Transaction struct {
	From string
	To   string
	Sum  float64
}

func NewTransaction(from, to Account, sum float64) Transaction {
	return Transaction{From: from.Name, To: to.Name, Sum: sum}
}

type Account struct {
	Name    string
	Balance float64
}

func NewBalanceFor(account Account, transactions []Transaction) Account {
	return Reduce(
		transactions,
		applyTransaction,
		account,
	)
}

func applyTransaction(a Account, transaction Transaction) Account {
	if transaction.From == a.Name {
		a.Balance -= transaction.Sum
	}
	if transaction.To == a.Name {
		a.Balance += transaction.Sum
	}
	return a
}
```

これは `Reduce` のような概念を使うことの威力をよく示していると思います。`NewBalanceFor` は、 _どのように_ ではなく _何が_ 起こるかを記述しており、より _宣言的_ に感じられます。コードを読むとき、私たちはたくさんのファイルを飛び回りながら、 _どのように_ ではなく _何が_ 起こっているのかを理解しようとしていることが多く、このスタイルのコードはそれをうまく助けてくれます。

詳細を掘り下げたければそうすることもでき、ループや状態の変更を気にすることなく `applyTransaction` の _ビジネスロジック_ を確認できます。それらは `Reduce` が別途面倒を見てくれるからです。


### fold/reduce はかなり普遍的

`Reduce`（または `Fold`）の可能性は無限大™️です。よく使われるパターンなのには理由があり、算術や文字列連結だけのためのものではありません。他の応用もいくつか試してみてください。

- いくつかの `color.RGBA` を混ぜて 1 つの色にしてみてはどうでしょう？
- 投票数や、買い物かごの商品数を合計する。
- リストの処理を伴うことなら、だいたい何でも。

## Find

Go にジェネリクスが加わったことで、それを高階関数と組み合わせれば、プロジェクト内の多くのボイラープレートコードを減らし、システムを理解・管理しやすくすることができます。

検索したいコレクションの型ごとに専用の `Find` 関数を書く必要はもうありません。代わりに `Find` 関数を再利用するか、1 つ書けばよいのです。上の `Reduce` 関数を理解できたなら、`Find` 関数を書くのは造作もないことでしょう。

テストはこちらです。

```go
func TestFind(t *testing.T) {
	t.Run("find first even number", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		firstEvenNumber, found := Find(numbers, func(x int) bool {
			return x%2 == 0
		})
		AssertTrue(t, found)
		AssertEqual(t, firstEvenNumber, 2)
	})
}
```

そして実装はこちらです。

```go
func Find[A any](items []A, predicate func(A) bool) (value A, found bool) {
	for _, v := range items {
		if predicate(v) {
			return v, true
		}
	}
	return
}
```

繰り返しになりますが、ジェネリックな型を受け取るので、さまざまな方法で再利用できます。

```go
type Person struct {
	Name string
}

t.Run("Find the best programmer", func(t *testing.T) {
	people := []Person{
		Person{Name: "Kent Beck"},
		Person{Name: "Martin Fowler"},
		Person{Name: "Chris James"},
	}

	king, found := Find(people, func(p Person) bool {
		return strings.Contains(p.Name, "Chris")
	})

	AssertTrue(t, found)
	AssertEqual(t, king, Person{Name: "Chris James"})
})
```

ご覧の通り、このコードには非の打ち所がありません。

## まとめ

センスよく使えば、このような高階関数はコードを読みやすく保守しやすいものにしてくれますが、次の経験則を忘れないでください。

TDD のプロセスを使って、実際に必要な、現実の具体的な振る舞いを導き出しましょう。そうすればリファクタリングの段階で、コードを整理するのに役立つ有用な抽象化を発見できる _かもしれません_ 。

TDD と良いソース管理の習慣を組み合わせて実践しましょう。テストが通ったら、リファクタリングを試みる _前に_ 作業をコミットしてください。そうすれば、めちゃくちゃにしてしまっても、簡単に動く状態に戻ることができます。

### 名前は重要

Go の外の世界についても調べる努力をしましょう。そうすれば、すでに確立した名前を持つ既存のパターンを再発明せずに済みます。

`A` のコレクションを受け取って `B` に変換する関数を書いていますか？それを `Convert` と呼ばないでください。それは [`Map`](https://en.wikipedia.org/wiki/Map_(higher-order_function)) です。こうしたものに「正しい」名前を使うことで、他の人の認知的負担が減り、さらに学ぶために検索エンジンで調べやすくもなります。

### これはイディオムに沿っていないように感じる？

心を開いてみてください。

ジェネリクスがリリースされたからといって、Go のイディオムが _根本的に_ 変わることはありませんし、変わるべきでもありませんが、言語が変わるのですから、イディオムは _変わる_ でしょう！これは議論の余地のある点ではないはずです。

それ以上の詳細なしに

> これはイディオムに沿っていない

と言うことは、特に新しい言語機能について議論しているときには、行動に移せることでも有用なことでもありません。

同僚とは、ドグマではなく、それぞれの長所に基づいてコードのパターンやスタイルについて議論しましょう。よく設計されたテストがある限り、自分やチームにとって何がうまくいくかを理解するにつれて、いつでもリファクタリングして物事を変えていくことができます。

### 参考資料

fold はコンピューターサイエンスにおける真の基礎です。さらに深く掘り下げたい場合は、次の興味深い資料を参照してください。
- [Wikipedia: Fold](https://en.wikipedia.org/wiki/Fold)
- [A tutorial on the universality and expressiveness of fold](http://www.cs.nott.ac.uk/~pszgmh/fold.pdf)
