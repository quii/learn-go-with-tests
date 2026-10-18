---
description: 構造体、メソッド、インターフェース
---

# 構造体、メソッド、インターフェース

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/structs)**

高さと幅を指定して長方形の周囲を計算するために、いくつかのジオメトリコードが必要だとします。 `Perimeter(width float64, height float64)`関数を記述できます。
ここで、 `float64`は`123.45`のような浮動小数点数用です。

テスト駆動開発（TDD）のサイクルはもうお馴染みのものになっているはずです。

## 最初にテストを書く

```go
func TestPerimeter(t *testing.T) {
	got := Perimeter(10.0, 10.0)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

新しいフォーマット文字列に注目してください。 `f`は`float64`用で、 `.2`は小数点以下2桁を出力することを意味します。

## テストを実行してみます

`./shapes_test.go:6:9: undefined: Perimeter`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func Perimeter(width float64, height float64) float64 {
	return 0
}
```

結果は `shapes_test.go:10: got 0.00 want 40.00` となります。

## 成功させるのに十分なコードを書く

```go
func Perimeter(width float64, height float64) float64 {
	return 2 * (width + height)
}
```

これまでのところ、とても簡単です。長方形の面積を返す `Area(width, height float64)`と呼ばれる関数を作成しましょう。

TDDサイクルに従って、自分で試してください。

最終的に、次のようなテストができあがるはずです

```go
func TestPerimeter(t *testing.T) {
	got := Perimeter(10.0, 10.0)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	got := Area(12.0, 6.0)
	want := 72.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

そして、このようなコード

```go
func Perimeter(width float64, height float64) float64 {
	return 2 * (width + height)
}

func Area(width float64, height float64) float64 {
	return width * height
}
```

浮動小数点数を `!=`/`==` で比較するのは、メモリ上での表現方法のせいで良くない考えだと聞いたことがあるかもしれません。

```go
func Example_floatComparison() {
	fmt.Println(0.1+0.2 == 0.3)

	var a, b, c float64 = 0.1, 0.2, 0.3
	fmt.Println(a+b == c)

	// Output:
	// true
	// false
}
```

2つ目の比較が `false` になるのは、`0.1` と `0.2` が `float64` では正確に表現できないため、それらを足し合わせても正確に `0.3` にはならないからです。1つ目が `true` になるのは、リテラル式として書かれているからにすぎません。Goはリテラル式を、何かに代入されるまでは `float64` ではなく任意精度で評価します。

ただし、この不正確さは上の `Area` のテストには当てはまりません。使っている値（`12.0`、`6.0`、`72.0` など）はすべて整数であり、`float64` で正確に表現 _できます_ 。そして、正確に表現された2つの数を掛け合わせた結果は、範囲内に収まっている限り、やはり正確に表現できる値になります。厳密な比較が危険になるのは、`0.1` のような正確に表現できない値や、除算の結果のような計算を持ち込んだときだけです。そのようなテストを書くことになったら、代わりに許容誤差を設けた比較を使いましょう。例えば `go-cmp` の [`cmp.Diff` と `cmpopts.EquateApprox`](https://pkg.go.dev/github.com/google/go-cmp/cmp/cmpopts#EquateApprox) などです。


## リファクタリング♪

私たちのコードは役割を果たしてはいますが、長方形であることを明示するものは何も含まれていません。不注意な開発者が、間違った答えが返ることに気付かずに、三角形の幅と高さをこれらの関数に渡してしまうかもしれません。

`RectangleArea`のように、より具体的な名前を関数に付けることができます。より適切なソリューションは、この概念をカプセル化する`Rectangle`と呼ばれる独自の**型**を定義することです。


**struct**を使用して単純な型を作成できます。[構造体](https://golang.org/ref/spec#Struct_types)は、データを保存できるフィールドの名前付きコレクションです。

`shapes.go`ファイルで次のように構造体を宣言します

```go
type Rectangle struct {
	Width  float64
	Height float64
}
```

では、プレーンな`float64`ではなく、`Rectangle`を使用するようにテストをリファクタリングしましょう。

```go
func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	rectangle := Rectangle{12.0, 6.0}
	got := Area(rectangle)
	want := 72.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}
```

修正を試みる前に必ずテストを実行してください。
テストは次のような有用なエラーを表示するはずです。

```text
./shapes_test.go:7:18: not enough arguments in call to Perimeter
    have (Rectangle)
    want (float64, float64)
```

`myStruct.field`の構文で構造体のフィールドにアクセスできます。

2つの関数を変更してテストを修正します。

```go
func Perimeter(rectangle Rectangle) float64 {
	return 2 * (rectangle.Width + rectangle.Height)
}

func Area(rectangle Rectangle) float64 {
	return rectangle.Width * rectangle.Height
}
```

`Rectangle`を関数に渡すと意図がより明確に伝わることに同意していただければ幸いです。構造体を使用することには他にもメリットがあり、それについては後で説明します。

次の要件は、円の`Area`関数を記述することです。

## 最初にテストを書く

```go
func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		got := Area(rectangle)
		want := 72.0

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := Area(circle)
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

}
```

ご覧のとおり、 `f` は `g`に置き換えられています。これにはちゃんとした理由があります。
`g`を使用すると、エラーメッセージでより精度の高い10進数が表示されます\([fmt options](https://golang.org/pkg/fmt/)\)。
例えば、円の面積の計算で半径1.5を使った場合、`f`では`7.068583`と表示されますが、`g`では`7.0685834705770345`と表示されます。

## テストを実行してみます

`./shapes_test.go:28:13: undefined: Circle`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`Circle`タイプを定義する必要があります。

```go
type Circle struct {
	Radius float64
}
```

もう一度テストを実行してみてください

`./shapes_test.go:29:14: cannot use circle (type Circle) as type Rectangle in argument to Area`

一部のプログラミング言語では、次のようなことができます。

```go
func Area(circle Circle) float64       {}
func Area(rectangle Rectangle) float64 {}
```

しかし、**Go**ではできません

`./shapes.go:20:32: Area redeclared in this block`

2つの選択肢があります。

* 同じ名前の関数を異なる _パッケージ_ で宣言することができます。新しいパッケージで `Area(Circle)`を作成することはできますが、ここではやりすぎだと感じます。
* 代わりに、新しく定義した型に[メソッド](https://golang.org/ref/spec#Method_declarations)を定義できます。

### メソッドとは?

これまでは _関数_ だけを書いてきましたが、メソッドもいくつか使ってきました。
`t.Errorf`を呼び出すときは、`t` \(`testing.T`\)のインスタンスでメソッド`Errorf`を呼び出しています。

メソッドは、レシーバーを持つ関数です。
メソッド宣言は、識別子（メソッド名）をメソッドにバインドし、メソッドをレシーバーの基底型に関連付けます。

メソッドは関数と非常によく似ていますが、特定の型のインスタンスに対して呼び出すという点が異なります。関数は`Area(rectangle)`のように好きな場所で呼び出せますが、メソッドは「もの」に対してしか呼び出せません。

例が役立つので、まずテストを変更して、代わりにメソッドを呼び出し、次にコードを修正しましょう。

```go
func TestArea(t *testing.T) {

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		got := rectangle.Area()
		want := 72.0

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		got := circle.Area()
		want := 314.1592653589793

		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	})

}
```

テストを実行しようとすると、次のようになります。

```text
./shapes_test.go:19:19: rectangle.Area undefined (type Rectangle has no field or method Area)
./shapes_test.go:29:16: circle.Area undefined (type Circle has no field or method Area)
```

> Circle型にはAreaというフィールドもメソッドもありません（type Circle has no field or method Area）

ここで改めて、コンパイラがいかに優れているかを強調しておきたいと思います。表示されたエラーメッセージを時間をかけてじっくり読むことは非常に重要です。長い目で見て必ず役に立ちます。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

タイプにいくつかのメソッドを追加しましょう

```go
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return 0
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 0
}
```

メソッドを宣言するための構文は、関数とほとんど同じです。
これは、メソッドが非常に似ているためです。
唯一の違いは、メソッドレシーバー `func (receiverName ReceiverType) MethodName(args)`の構文です。

そのタイプの変数でメソッドが呼び出されると、 `receiverName`変数を介してそのデータへの参照が取得されます。他の多くのプログラミング言語では、これは暗黙的に行われ、 `this`を介してレシーバーにアクセスします。

Goの慣例では、レシーバー変数をタイプの最初の文字にします。

```
r Rectangle
```

テストを再実行しようとすると、テストがコンパイルされ、失敗した出力がいくつか表示されます。

## 成功させるのに十分なコードを書く

新しいメソッドを修正して、長方形のテストに成功させましょう

```go
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
```

テストを再実行すると、長方形のテストは成功するはずですが、円のテストはまだ失敗しているはずです。

円の `Area`関数のテストを通すために、`math`パッケージから `Pi`定数を借ります（インポートすることを忘れないでください）。

```go
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
```

## リファクタリング♪

テストに重複があります。

やりたいことは、`shapes`のコレクションを取得し、それらの `Area()`メソッドを呼び出して、結果を確認することだけです。

`Rectangle`と`Circle`の両方を渡せる `checkArea`のような関数を書けるようにしつつ、形状ではないものを渡そうとしたらコンパイルに失敗するようにしたいのです。

Goでは、この意図を**インターフェース**で体系化できます。

[インターフェース](https://golang.org/ref/spec#Interface_types)は、Goなどの静的型付き言語で非常に強力な概念です。
さまざまな型で使用できる関数を作成でき、型安全性を保ったまま高度に疎結合なコードを書けるからです。

テストをリファクタリングしてこれを紹介しましょう。

```go
func TestArea(t *testing.T) {

	checkArea := func(t testing.TB, shape Shape, want float64) {
		t.Helper()
		got := shape.Area()
		if got != want {
			t.Errorf("got %g want %g", got, want)
		}
	}

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{12, 6}
		checkArea(t, rectangle, 72.0)
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10}
		checkArea(t, circle, 314.1592653589793)
	})

}
```

他の演習と同様にヘルパー関数を作成していますが、今回は`Shape`が渡されるように要求しています。これを形状ではないもので呼び出そうとすると、コンパイルされません。

では、どうすれば何かが形状（shape）になるのでしょうか？
インターフェース宣言を使って、`Shape`が何であるかをGoに伝えるだけです

```go
type Shape interface {
	Area() float64
}
```

`Rectangle`と`Circle`で行ったように新しい `type`を作成していますが、今回は`struct`ではなく `interface`です。

これをコードに追加すると、テストに合格します。

### ちょっと待って、どういうこと？

これは、他のほとんどのプログラミング言語のインターフェースとはかなり異なります。通常、`My type Foo implements interface Bar`と言うコードを書く必要があります。

しかし、私たちの場合

* `Rectangle`には`Area`というメソッドがあり、 `float64`を返すため、`Shape`インターフェースを満たします
* `Circle`には`Area`というメソッドがあり、 `float64`を返すため、`Shape`インターフェースを満たします
* `string`にはそのようなメソッドがないため、インターフェースを満たしていません
* など

Goでは、**インターフェースの解決は暗黙的です**。
渡したタイプがインターフェースが要求するものと一致する場合、それはコンパイルされます。

### 切り離し（Decoupling）

ヘルパーが形状が `Rectangle`、`Circle`、または `Triangle`のどちらであるかを気にする必要がないことに注意してください。
インターフェースを宣言することにより、ヘルパーは具象型から _切り離され_ （疎結合になり）、その機能を実行するために必要なメソッドのみを持ちます。

インターフェースを使用して**必要なもののみ**を宣言するこの種のアプローチは、ソフトウェア設計において非常に重要であり、後のセクションでより詳細に説明します。

## さらにリファクタリング

構造体についてある程度理解できたので、「**テーブル駆動テスト**」を紹介します。

[テーブル駆動テスト](https://go.dev/wiki/TableDrivenTests)は、同じ方法でテストできるテストケースのリストを作成する場合に役立ちます。

```go
func TestArea(t *testing.T) {

	areaTests := []struct {
		shape Shape
		want  float64
	}{
		{Rectangle{12, 6}, 72.0},
		{Circle{10}, 314.1592653589793},
	}

	for _, tt := range areaTests {
		got := tt.shape.Area()
		if got != tt.want {
			t.Errorf("got %g want %g", got, tt.want)
		}
	}

}
```

ここでの唯一の新しい構文は、「匿名の構造体」`areaTests`を作成することです。 2つのフィールド、 `shape`と`want`で `[]struct`を使用して、構造体のスライスを宣言しています。次に、スライスをケースで埋めます。

次に、構造体フィールドを使用してテストを実行し、他のスライスと同じようにそれらを繰り返します。

開発者が新しい形状を導入し、 `Area`を実装してテストケースに追加するのが、いかに簡単かがわかるでしょう。
さらに、`Area`でバグが見つかった場合も、修正する前にそのバグを確かめる新しいテストケースを追加するのは非常に簡単です。

テーブル駆動テストは、ツールボックスの優れた道具になり得ますが、テストに余分なノイズを加える必要が本当にあるのかを確認してください。
インターフェースのさまざまな実装をテストしたい場合、または関数に渡されるデータに、テストを必要とするさまざまな要件がたくさんある場合、それらは非常に適しています。

別の形状を追加してテストすることで、これらすべてを実際に示してみましょう。三角形です。

## 最初にテストを書く

新しい形状の新しいテストを追加するのはとても簡単です。リストに`{Triangle{12, 6}, 36.0},`を追加するだけです。

```go
func TestArea(t *testing.T) {

	areaTests := []struct {
		shape Shape
		want  float64
	}{
		{Rectangle{12, 6}, 72.0},
		{Circle{10}, 314.1592653589793},
		{Triangle{12, 6}, 36.0},
	}

	for _, tt := range areaTests {
		got := tt.shape.Area()
		if got != tt.want {
			t.Errorf("got %g want %g", got, tt.want)
		}
	}

}
```

## テストを実行してみます

テストの実行を繰り返し、コンパイラに解決策へと導いてもらうことを忘れないでください。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`./shapes_test.go:25:4: undefined: Triangle`

`Triangle`はまだ定義していません

```go
type Triangle struct {
	Base   float64
	Height float64
}
```

再試行

```text
./shapes_test.go:25:8: cannot use Triangle literal (type Triangle) as type Shape in field value:
    Triangle does not implement Shape (missing Area method)
```

`Triangle`には `Area()`メソッドがないため形状として使用できない、と言われています。テストが動くように空の実装を追加しましょう

```go
func (t Triangle) Area() float64 {
	return 0
}
```

ようやくコードがコンパイルされ、エラーが表示されます

`shapes_test.go:31: got 0.00 want 36.00`

## 成功させるのに十分なコードを書く

```go
func (t Triangle) Area() float64 {
	return (t.Base * t.Height) * 0.5
}
```

そして、テストは成功です！

## リファクタリング♪

繰り返しになりますが、実装は問題ありませんが、テストでは多少の改善が見込めます。

これを見直すと

```
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

すべての数値が何を表しているのかすぐには明確ではなく、テストを簡単に理解できるようにする必要があります。

ここまでは、構造体のインスタンスを作成するための `MyStruct{val1, val2}` という構文だけを示してきましたが、オプションでフィールドに名前を付けることができます。

それがどのように見えるか見てみましょう

```
        {shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
        {shape: Circle{Radius: 10}, want: 314.1592653589793},
        {shape: Triangle{Base: 12, Height: 6}, want: 36.0},
```

[例によるテスト駆動開発](https://g.co/kgs/yCzDLF) で、Kent Beckはいくつかのテストをある程度までリファクタリングしたうえで、次のように述べています。

> テストは、**一連の操作ではなく**、真実の表明であるかのように、より明確に私たちに語りかけてくる

\(引用中の強調は私によるものです\)

これで、テスト、というよりテストケースのリストが、形状とその面積についての真実を主張するようになりました。

## テスト出力が役立つことを確認する

以前に`Triangle`を実装していて、失敗したテストがあったことを覚えていますか？ 

`shapes_test.go:31: got 0.00 want 36.00`と表示されました。

ちょうど`Triangle`を扱っていたので、これが`Triangle`に関連していることはわかっていました。
しかし、表の20のケースのいずれかでバグがシステムに紛れ込んだ場合はどうなるでしょうか？
開発者はどのケースが失敗したかをどのようにして知るのでしょうか？
これは開発者にとって良い体験ではありません。実際に失敗したケースを見つけるために、ケースを手動で調べる必要があります。

エラーメッセージを `%#v got %g want %g`に変更できます。 `%#v`形式の文字列は、フィールドの値を含む構造体を出力するため、開発者はテストされているプロパティを一目で確認できます。

テストケースを読みやすくするために、 `want`フィールドの名前を`hasArea`のようなわかりやすい名前に変更できます。

テーブル駆動テストの最後のヒントは、 `t.Run`を使用してテストケースに名前を付けることです。

各ケースを `t.Run`でラップすることで、ケースの名前が出力されるため、失敗時のテスト出力がより明確になります。

```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
```

また、 `go test -run TestArea/Rectangle`を使用して、テーブル内で特定のテストを実行できます。

これらを反映した最終的なテストコードは次のとおりです

```go
func TestArea(t *testing.T) {

	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
		{name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
	}

	for _, tt := range areaTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.hasArea {
				t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
			}
		})

	}

}
```

## まとめ

今回もTDDの実践でした。基本的な数学の問題に対する解決策を反復的に改善しながら、テストをきっかけに新しい言語機能を学びました。

* 構造体を宣言して独自のデータ型を作成し、関連するデータをまとめてコードの意図を明確にする
* さまざまな型で使用できる関数を定義できるようにインターフェースを宣言する \([アドホック多相](https://en.wikipedia.org/wiki/Ad_hoc_polymorphism)\)
* データ型に機能を追加したり、インターフェースを実装したりできるようにメソッドを追加する
* アサーションをより明確にし、テストスイートを拡張および保守しやすくするためのテーブル駆動テスト

これは重要な章でした。私たちは今、独自の型を定義し始めているからです。 Goのような静的に型付けされた言語では、理解しやすく、つなぎ合わせてテストできるソフトウェアを構築するために、独自の型を設計できることが不可欠です。

インターフェースは、システムの他の部分から複雑さを隠すための優れたツールです。私たちの場合、テストヘルパーの _コード_ は、それがアサートしている正確な形状を知る必要はなく、その面積を「尋ねる」方法を知るだけでした。

Goに慣れるにつれて、インターフェースと標準ライブラリの本当の強みがわかり始めるでしょう。
 _あらゆる場所_ で使用される、標準ライブラリで定義されたインターフェースについて学び、独自の型に対してそれらを実装することで、多くの優れた機能を非常に迅速に再利用できます。
