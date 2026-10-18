---
description: ジェネリクス
---

# ジェネリクス

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/generics)**

この章では、ジェネリクスを紹介し、ジェネリクスに対して抱いているかもしれない懸念を払拭し、今後コードをどのようにシンプルにできるかのアイデアを提供します。この章を読めば、次のものを書けるようになります。

- ジェネリックな引数を取る関数
- ジェネリックなデータ構造


## 独自のテストヘルパー（`AssertEqual`、`AssertNotEqual`）

ジェネリクスを探求するために、いくつかのテストヘルパーを書いてみましょう。

### 整数に対するアサーション

基本的なところから始めて、ゴールに向かって反復していきましょう。

```go
import "testing"

func TestAssertFunctions(t *testing.T) {
	t.Run("asserting on integers", func(t *testing.T) {
		AssertEqual(t, 1, 1)
		AssertNotEqual(t, 1, 2)
	})
}

func AssertEqual(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func AssertNotEqual(t *testing.T, got, want int) {
	t.Helper()
	if got == want {
		t.Errorf("didn't want %d", got)
	}
}
```


### 文字列に対するアサーション

整数が等しいことをアサーションできるのは素晴らしいことですが、`string` に対してアサーションしたい場合はどうでしょうか？

```go
t.Run("asserting on strings", func(t *testing.T) {
	AssertEqual(t, "hello", "hello")
	AssertNotEqual(t, "hello", "Grace")
})
```

エラーが発生します。

```
# github.com/quii/learn-go-with-tests/generics [github.com/quii/learn-go-with-tests/generics.test]
./generics_test.go:12:18: cannot use "hello" (untyped string constant) as int value in argument to AssertEqual
./generics_test.go:13:21: cannot use "hello" (untyped string constant) as int value in argument to AssertNotEqual
./generics_test.go:13:30: cannot use "Grace" (untyped string constant) as int value in argument to AssertNotEqual
```

時間をかけてエラーを読めば、`integer` を期待する関数に `string` を渡そうとしていることについて、コンパイラが文句を言っていることがわかるでしょう。

#### 型安全性のおさらい

この本のこれまでの章を読んできた方や、静的型付け言語の経験がある方であれば、これは驚くことではないでしょう。Go のコンパイラは、関数や構造体などを書く際に、どの型を扱いたいのかを記述することを求めます。

`integer` を期待する関数に `string` を渡すことはできません。

これは形式的な手続きのように感じられるかもしれませんが、非常に役に立つことがあります。これらの制約を記述することで、

- 関数の実装がシンプルになります。どの型を扱うのかをコンパイラに記述することで、**有効な実装の可能性の数を制約** します。`Person` と `BankAccount` を「足す」ことはできません。`integer` を大文字にすることもできません。ソフトウェアにおいて、制約は多くの場合とても役に立ちます。
- 意図しないデータを誤って関数に渡すことを防げます。

Go では、[インターフェース](./structs-methods-and-interfaces.md)によって型をより抽象的に扱う方法が提供されています。これにより、具体的な型ではなく、必要な振る舞いを提供する型を受け取る関数を設計できます。これは、型安全性を維持しながらある程度の柔軟性を与えてくれます。

### 文字列または整数を受け取る関数？（あるいは、実際には他のものも）

関数をより柔軟にするために Go が持っているもう 1 つの選択肢は、引数の型を `interface{}` として宣言することです。これは「何でも」を意味します。

シグネチャをこの型を使うように変更してみてください。

```go
func AssertEqual(got, want interface{})

func AssertNotEqual(got, want interface{})

```

これでテストはコンパイルされ、成功するはずです。テストを失敗させてみると、メッセージの出力に整数用の `%d` フォーマット文字列を使っているため、出力が少しお粗末なことがわかるでしょう。どんな種類の値でもより良く出力できるように、汎用の `%+v` フォーマットに変更しましょう。

### `interface{}` の問題点

私たちの `AssertX` 関数はかなり素朴なものですが、概念的には[他の人気ライブラリがこの機能を提供している方法](https://github.com/matryer/is/blob/master/is.go#L160)とそれほど変わりません。

```go
func (is *I) Equal(a, b interface{})
```

では、何が問題なのでしょうか？

`interface{}` を使うと、関数に渡されるものの型について何も有用な情報をコンパイラに伝えていないため、コードを書く際にコンパイラが助けてくれません。2 つの異なる型を比較してみてください。

```go
AssertEqual(1, "1")
```

この場合はなんとかなります。テストはコンパイルされ、期待どおりに失敗します。ただし、エラーメッセージ `got 1, want 1` はわかりにくいです。しかし、そもそも文字列と整数を比較できるようにしたいのでしょうか？ `Person` と `Airport` を比較するのはどうでしょう？

`interface{}` を受け取る関数を書くことは、非常に難しくバグを生みやすいものになりえます。なぜなら、制約を _失って_ しまい、どのような種類のデータを扱っているのかについての情報がコンパイル時に何もないからです。

これは、**コンパイラが助けてくれない** ということを意味し、代わりに **実行時エラー** が起こりやすくなります。実行時エラーはユーザーに影響を与えたり、障害を引き起こしたり、さらに悪い事態を招いたりする可能性があります。

開発者は、こうした *ゴホン* ジェネリックな関数を実装するために、リフレクションを使わなければならないことがよくあります。リフレクションは読み書きが複雑になりがちで、プログラムのパフォーマンスを損なう可能性もあります。

## ジェネリクスを使った独自のテストヘルパー

理想的には、扱うすべての型ごとに個別の `AssertX` 関数を作る必要はないようにしたいものです。_あらゆる_ 型で動作するが、[リンゴとオレンジ](https://en.wikipedia.org/wiki/Apples_and_oranges)を比較することはできない、_1 つの_ `AssertEqual` 関数を持てるようにしたいのです。

ジェネリクスは、**制約を記述する** ことによって、（インターフェースのような）抽象化を作る方法を提供してくれます。ジェネリクスを使うと、`interface{}` が提供するのと同程度の柔軟性を持ちつつ、型安全性を保ち、呼び出し側にとってより良い開発者体験を提供する関数を書くことができます。

```go
func TestAssertFunctions(t *testing.T) {
	t.Run("asserting on integers", func(t *testing.T) {
		AssertEqual(t, 1, 1)
		AssertNotEqual(t, 1, 2)
	})

	t.Run("asserting on strings", func(t *testing.T) {
		AssertEqual(t, "hello", "hello")
		AssertNotEqual(t, "hello", "Grace")
	})

	// AssertEqual(t, 1, "1") // uncomment to see the error
}

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("didn't want %v", got)
	}
}
```

Go でジェネリックな関数を書くには、「型パラメータ」を指定する必要があります。これは「ジェネリックな型を記述して、それにラベルを付ける」ということを小難しく言っているだけです。

私たちの場合、型パラメータの型は `comparable` で、それに `T` というラベルを付けました。このラベルを使って、関数の引数の型を記述できます（`got, want T`）。

`comparable` を使っているのは、関数内で `T` 型のものに対して `==` と `!=` 演算子を使いたいということをコンパイラに伝えたいからです。比較したいのです！ 型を `any` に変更してみると、

```go
func AssertNotEqual[T any](got, want T)
```

次のエラーが発生します。

```
prog.go2:15:5: cannot compare got != want (operator != not defined for T)
```

これはとても理にかなっています。なぜなら、すべての（`any` の）型に対してそれらの演算子を使えるわけではないからです。

### [`T any`](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md#the-constraint) を使ったジェネリックな関数は `interface{}` と同じなのか？

2 つの関数を考えてみましょう。

```go
func GenericFoo[T any](x, y T)
```

```go
func InterfaceyFoo(x, y interface{})
```

ここでジェネリクスを使う意味は何でしょうか？ `any` は... 何でも、を表すのではないのでしょうか？

制約という観点では、`any` は確かに「何でも」を意味し、`interface{}` も同様です。実際、`any` は 1.18 で追加されたもので、_`interface{}` の単なるエイリアス_ です。

ジェネリック版との違いは、_依然として特定の型を記述している_ ということであり、それはつまり、この関数が _1 つの_ 型でのみ動作するように制約しているということです。

これが意味するのは、`InterfaceyFoo` はどんな型の組み合わせでも呼び出せる（例: `InterfaceyFoo(apple, orange)`）ということです。しかし、`GenericFoo` は _1 つの_ 型 `T` でのみ動作すると宣言しているため、依然としていくつかの制約を提供します。

有効なもの:

- `GenericFoo(apple1, apple2)`
- `GenericFoo(orange1, orange2)`
- `GenericFoo(1, 2)`
- `GenericFoo("one", "two")`

有効でないもの（コンパイルに失敗する）:

- `GenericFoo(apple1, orange1)`
- `GenericFoo("1", 1)`

関数がジェネリックな型を返す場合、呼び出し側はその型をそのまま使うこともできます。関数が `interface{}` を返す場合はコンパイラが型について何も保証できないため型アサーションを行う必要がありますが、その必要はありません。

## 次へ: ジェネリックなデータ型

[スタック](https://en.wikipedia.org/wiki/Stack_(abstract_data_type))というデータ型を作成します。スタックは、要件の観点からはかなり理解しやすいはずです。スタックはアイテムのコレクションで、アイテムを「一番上」に `Push` し、アイテムを取り戻すには一番上からアイテムを `Pop` します（LIFO - last in, first out、後入れ先出し）。

簡潔にするため、`int` のスタックと `string` のスタックについて、次のコードにたどり着くまでの TDD プロセスは省略しました。

```go
type StackOfInts struct {
	values []int
}

func (s *StackOfInts) Push(value int) {
	s.values = append(s.values, value)
}

func (s *StackOfInts) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *StackOfInts) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}

type StackOfStrings struct {
	values []string
}

func (s *StackOfStrings) Push(value string) {
	s.values = append(s.values, value)
}

func (s *StackOfStrings) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *StackOfStrings) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}
```

手助けとなるアサーション関数をいくつか追加で作成しました。

```go
func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got %v, want true", got)
	}
}

func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("got %v, want false", got)
	}
}
```

そして、テストは次のとおりです。

```go
func TestStack(t *testing.T) {
	t.Run("integer stack", func(t *testing.T) {
		myStackOfInts := new(StackOfInts)

		// check stack is empty
		AssertTrue(t, myStackOfInts.IsEmpty())

		// add a thing, then check it's not empty
		myStackOfInts.Push(123)
		AssertFalse(t, myStackOfInts.IsEmpty())

		// add another thing, pop it back again
		myStackOfInts.Push(456)
		value, _ := myStackOfInts.Pop()
		AssertEqual(t, value, 456)
		value, _ = myStackOfInts.Pop()
		AssertEqual(t, value, 123)
		AssertTrue(t, myStackOfInts.IsEmpty())
	})

	t.Run("string stack", func(t *testing.T) {
		myStackOfStrings := new(StackOfStrings)

		// check stack is empty
		AssertTrue(t, myStackOfStrings.IsEmpty())

		// add a thing, then check it's not empty
		myStackOfStrings.Push("123")
		AssertFalse(t, myStackOfStrings.IsEmpty())

		// add another thing, pop it back again
		myStackOfStrings.Push("456")
		value, _ := myStackOfStrings.Pop()
		AssertEqual(t, value, "456")
		value, _ = myStackOfStrings.Pop()
		AssertEqual(t, value, "123")
		AssertTrue(t, myStackOfStrings.IsEmpty())
	})
}
```

### 問題点

- `StackOfStrings` と `StackOfInts` のコードはほぼ同じです。重複は必ずしも世界の終わりではありませんが、読み、書き、保守すべきコードが増えます。
- 2 つの型にまたがってロジックを重複させているため、テストも重複させなければなりませんでした。

私たちが本当にしたいのは、スタックという _考え_ を 1 つの型で表現し、それに対して 1 組のテストを持つことです。今はリファクタリングの帽子をかぶっているはずです。つまり、同じ振る舞いを維持したいので、テストは変更すべきではありません。

ジェネリクスがなければ、次のようにすることが _できます_。

```go
type StackOfInts = Stack
type StackOfStrings = Stack

type Stack struct {
	values []interface{}
}

func (s *Stack) Push(value interface{}) {
	s.values = append(s.values, value)
}

func (s *Stack) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		var zero interface{}
		return zero, false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}
```

- 以前の `StackOfInts` と `StackOfStrings` の実装を、新しく統一された型 `Stack` のエイリアスにしています
- `values` を `interface{}` の[スライス](arrays-and-slices.md)にすることで、`Stack` から型安全性を取り除いています

このコードを試すには、アサーション関数から型の制約を取り除く必要があります。

```go
func AssertEqual(t *testing.T, got, want interface{})
```

こうすると、テストは依然として成功します。ジェネリクスなんて誰が必要とするのでしょう？

### 型安全性を捨てることの問題点

最初の問題は、`AssertEquals` で見たのと同じです。型安全性を失ってしまいました。今ではオレンジのスタックにリンゴを `Push` できてしまいます。

たとえそうしないだけの規律を持っていたとしても、このコードは依然として扱いにくいものです。なぜなら、**`interface{}` を返すメソッドは扱うのがひどく面倒** だからです。

次のテストを追加してみてください。

```go
t.Run("interface stack DX is horrid", func(t *testing.T) {
	myStackOfInts := new(StackOfInts)

	myStackOfInts.Push(1)
	myStackOfInts.Push(2)
	firstNum, _ := myStackOfInts.Pop()
	secondNum, _ := myStackOfInts.Pop()
	AssertEqual(t, firstNum+secondNum, 3)
})
```

コンパイラエラーが発生し、型安全性を失うことの弱点が示されます。

```
invalid operation: operator + not defined on firstNum (variable of type interface{})
```

`Pop` が `interface{}` を返すということは、コンパイラがそのデータが何であるかについての情報を持たないことを意味し、そのため私たちにできることが大きく制限されます。コンパイラはそれが整数であるべきだと知ることができないので、`+` 演算子を使わせてくれません。

これを回避するには、呼び出し側がそれぞれの値に対して[型アサーション](https://golang.org/ref/spec#Type_assertions)を行う必要があります。

```go
t.Run("interface stack dx is horrid", func(t *testing.T) {
	myStackOfInts := new(StackOfInts)

	myStackOfInts.Push(1)
	myStackOfInts.Push(2)
	firstNum, _ := myStackOfInts.Pop()
	secondNum, _ := myStackOfInts.Pop()

	// get our ints from out interface{}
	reallyFirstNum, ok := firstNum.(int)
	AssertTrue(t, ok) // need to check we definitely got an int out of the interface{}

	reallySecondNum, ok := secondNum.(int)
	AssertTrue(t, ok) // and again!

	AssertEqual(t, reallyFirstNum+reallySecondNum, 3)
})
```

このテストからにじみ出る不快さは、`Stack` 実装を使う可能性のあるすべてのユーザーにとって繰り返されることになります。うんざりですね。

### ジェネリックなデータ構造が助けに来る

関数に対してジェネリックな引数を定義できるのと同じように、ジェネリックなデータ構造を定義することもできます。

これが、ジェネリックなデータ型を使った新しい `Stack` の実装です。

```go
type Stack[T any] struct {
	values []T
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}
```

テストは次のとおりで、完全な型安全性を保ちながら、私たちが望むとおりに動作していることがわかります。

```go
func TestStack(t *testing.T) {
	t.Run("integer stack", func(t *testing.T) {
		myStackOfInts := new(Stack[int])

		// check stack is empty
		AssertTrue(t, myStackOfInts.IsEmpty())

		// add a thing, then check it's not empty
		myStackOfInts.Push(123)
		AssertFalse(t, myStackOfInts.IsEmpty())

		// add another thing, pop it back again
		myStackOfInts.Push(456)
		value, _ := myStackOfInts.Pop()
		AssertEqual(t, value, 456)
		value, _ = myStackOfInts.Pop()
		AssertEqual(t, value, 123)
		AssertTrue(t, myStackOfInts.IsEmpty())

		// can get the numbers we put in as numbers, not untyped interface{}
		myStackOfInts.Push(1)
		myStackOfInts.Push(2)
		firstNum, _ := myStackOfInts.Pop()
		secondNum, _ := myStackOfInts.Pop()
		AssertEqual(t, firstNum+secondNum, 3)
	})
}
```

ジェネリックなデータ構造を定義する構文が、関数のジェネリックな引数を定義する構文と一貫していることに気づくでしょう。

```go
type Stack[T any] struct {
	values []T
}
```

以前と _ほぼ_ 同じですが、ここで言っているのは、**スタックの型が、扱える値の型を制約する** ということです。

`Stack[Orange]` や `Stack[Apple]` を作成すると、スタックに定義されたメソッドは、扱っているスタックの特定の型しか渡すことができず、またその型しか返しません。

```go
func (s *Stack[T]) Pop() (T, bool)
```

作成するスタックの型に応じて、実装の型が何らかの形で自動生成されると想像するとよいでしょう。

```go
func (s *Stack[Orange]) Pop() (Orange, bool)
```

```go
func (s *Stack[Apple]) Pop() (Apple, bool)
```

このリファクタリングを行ったので、同じロジックを何度も証明する必要はないため、文字列スタックのテストは安全に削除できます。

ここまでのジェネリックな関数を呼び出す例では、ジェネリックな型を指定する必要がなかったことに注意してください。たとえば、`AssertEqual[T]` を呼び出すとき、`T` は引数から推論できるので、`T` が何の型であるかを指定する必要はありません。ジェネリックな型を推論できない場合は、関数を呼び出すときに型を指定する必要があります。構文は関数を定義するときと同じで、引数の前の角括弧の中に型を指定します。

具体的な例として、`Stack[T]` のコンストラクタを作ることを考えてみましょう。
```go
func NewStack[T any]() *Stack[T] {
	return new(Stack[T])
}
```
このコンストラクタを使って、たとえば int のスタックと string のスタックを作成するには、次のように呼び出します。
```go
myStackOfInts := NewStack[int]()
myStackOfStrings := NewStack[string]()
```

コンストラクタを追加した後の `Stack` の実装とテストは次のとおりです。

```go
type Stack[T any] struct {
	values []T
}

func NewStack[T any]() *Stack[T] {
	return new(Stack[T])
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}

	index := len(s.values) - 1
	el := s.values[index]
	s.values = s.values[:index]
	return el, true
}
```

```go
func TestStack(t *testing.T) {
	t.Run("integer stack", func(t *testing.T) {
		myStackOfInts := NewStack[int]()

		// check stack is empty
		AssertTrue(t, myStackOfInts.IsEmpty())

		// add a thing, then check it's not empty
		myStackOfInts.Push(123)
		AssertFalse(t, myStackOfInts.IsEmpty())

		// add another thing, pop it back again
		myStackOfInts.Push(456)
		value, _ := myStackOfInts.Pop()
		AssertEqual(t, value, 456)
		value, _ = myStackOfInts.Pop()
		AssertEqual(t, value, 123)
		AssertTrue(t, myStackOfInts.IsEmpty())

		// can get the numbers we put in as numbers, not untyped interface{}
		myStackOfInts.Push(1)
		myStackOfInts.Push(2)
		firstNum, _ := myStackOfInts.Pop()
		secondNum, _ := myStackOfInts.Pop()
		AssertEqual(t, firstNum+secondNum, 3)
	})
}
```


ジェネリックなデータ型を使うことで、次のことを実現しました。

- 重要なロジックの重複を減らしました。
- `Pop` が `T` を返すようにしたので、`Stack[int]` を作成すれば、実際に `Pop` から `int` が返ってきます。型アサーションの曲芸をする必要なく `+` を使えるようになりました。
- コンパイル時に誤用を防げるようになりました。リンゴのスタックにオレンジを `Push` することはできません。

## まとめ

この章で、ジェネリクスの構文を味わい、ジェネリクスがなぜ役に立つのかについてのアイデアを得られたはずです。独自の `Assert` 関数を書きました。これは、ジェネリクスに関する他のアイデアを試す際にも安全に再利用できます。また、好きな型のデータを型安全な方法で格納できる、シンプルなデータ構造を実装しました。

### ほとんどの場合、ジェネリクスは `interface{}` を使うよりもシンプル

静的型付け言語の経験が浅い場合、ジェネリクスの意義はすぐには明らかでないかもしれません。しかし、この章の例によって、Go 言語が私たちの望むほど表現力豊かではない部分が示せたのではないかと思います。特に `interface{}` を使うと、コードは次のようになります。

- 安全性が低くなり（リンゴとオレンジが混ざる）、より多くのエラー処理が必要になる
- 表現力が低くなる。`interface{}` はデータについて何も教えてくれない
- [リフレクション](reflection.md)や型アサーションなどに頼る可能性が高くなる。これにより、チェックがコンパイル時から実行時に押しやられるため、コードが扱いにくくなり、エラーも起こりやすくなる

静的型付け言語を使うということは、制約を記述するという行為です。それをうまく行えば、安全で使いやすいだけでなく、解の空間が狭くなるため書くのもシンプルなコードを作れます。

ジェネリクスは、コードの中で制約を表現する新しい方法を与えてくれます。これまで見てきたように、Go 1.18 まではできなかったコードの統合と簡素化が可能になります。

### ジェネリクスは Go を Java に変えてしまうのか？

- いいえ。

Go コミュニティには、ジェネリクスが悪夢のような抽象化や、不可解なコードベースにつながるという多くの [FUD（恐怖、不確実性、疑念）](https://en.wikipedia.org/wiki/Fear,_uncertainty,_and_doubt)があります。これには通常「注意して使わなければならない」という但し書きが付きます。

これは真実ではありますが、どの言語機能についても言えることなので、特に役に立つアドバイスではありません。

インターフェースを定義できることについて文句を言う人はあまりいません。インターフェースもジェネリクスと同様に、コード内の制約を記述する方法です。インターフェースを記述するとき、あなたは _まずいものになりうる_ 設計上の選択をしています。わかりにくく使いにくいコードを作りうるのは、ジェネリクスに限った話ではありません。

### あなたはすでにジェネリクスを使っている

配列、スライス、マップを使ったことがあるなら、_すでにジェネリックなコードの利用者だった_ ということを考えてみてください。

```
var myApples []Apple
// You can't do this!
append(myApples, Orange{})
```

### 抽象化は汚い言葉ではない

[AbstractSingletonProxyFactoryBean](https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/aop/framework/AbstractSingletonProxyFactoryBean.html) をこき下ろすのは簡単ですが、抽象化がまったくないコードベースもまた悪いものだということから目を背けないようにしましょう。適切なときに関連する概念を _まとめ_、システムを理解しやすく変更しやすくするのはあなたの仕事です。そうしなければ、システムは明確さに欠けるバラバラの関数と型の寄せ集めになってしまいます。

### [動くようにし、正しくし、速くする（Make it work, make it right, make it fast）](https://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast#:~:text=%22Make%20it%20work%2C%20make%20it,to%20DesignForPerformance%20ahead%20of%20time.)

ジェネリクスで問題に陥るのは、良い設計判断を下すための十分な情報がないまま、早すぎる段階で抽象化してしまうときです。

レッド、グリーン、リファクタリングという TDD のサイクルは、**事前に抽象化を想像するのではなく**、振る舞いを実現するために _実際に必要な_ コードについて、より多くの指針を与えてくれます。それでも注意は必要です。

ここに厳格なルールはありませんが、有用な一般化が見えてくるまでは、ものごとをジェネリックにするのは控えましょう。さまざまな `Stack` の実装を作ったとき、重要なのは、テストに裏付けられた `StackOfStrings` や `StackOfInts` のような _具体的な_ 振る舞いから始めたことです。_実際の_ コードから本当のパターンが見え始め、テストに裏付けられた状態で、より汎用的な解決策に向けたリファクタリングを探求することができました。

同じコードを 3 回見たときにだけ一般化せよ、とよくアドバイスされますが、これは手始めの経験則として良さそうです。

私が他のプログラミング言語でよくたどってきた道筋は次のとおりです。

- ある振る舞いを駆動するための 1 回目の TDD サイクル
- 関連する他のシナリオを試すための、もう 1 回の TDD サイクル

> うーん、これらは似ているな。でも、少しの重複は、まずい抽象化に結合するよりはましだ

- 一晩寝かせる
- さらにもう 1 回の TDD サイクル

> よし、これを一般化できるか試してみたい。ありがたいことに、私は TDD を使っているのでとても賢くてかっこいい。だからいつでも好きなときにリファクタリングできるし、このプロセスのおかげで、設計しすぎる前に実際に必要な振る舞いを理解できた。

- この抽象化はいい感じだ！ テストは依然として成功しているし、コードはよりシンプルになった
- これで多くのテストを削除できる。振る舞いの _本質_ を捉え、不要な詳細を取り除いた。
