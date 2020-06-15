---
description: 'Structs, methods & interfaces'
---

# 構造体、メソッド、インターフェース

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/structs)

高さと幅を指定して長方形の周囲を計算するために、いくつかのジオメトリコードが必要だとします。 `Perimeter(width float64, height float64)`関数を記述できます。
ここで、 `float64`は` 123.45`のような浮動小数点数用です。

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

新しいフォーマット文字列に注目してください。 `f`は` float64`用で、 `.2`は小数点以下2桁を出力することを意味します。

## テストを実行してみます

`./shapes_test.go:6:9: undefined: Perimeter`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func Perimeter(width float64, height float64) float64 {
    return 0
}
```

Results in `shapes_test.go:10: got 0.00 want 40.00`.

## 成功させるのに十分なコードを書く

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}
```

これまでのところ、とても簡単です。長方形の面積を返す `Area(width, height float64)`と呼ばれる関数を作成しましょう。

TDDサイクルに従って、自分で試してください。

おそらく、あなたはこのようなテストで終わるはずでしょう

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

## リファクタリング

私たちのコードはその役割を果たしますが、四角形について明示的なものは何も含まれていません。不注意な開発者は、三角形の幅と高さを間違った答えを返すことに気付かずにこれらの関数に提供しようとする場合があります。

`RectangleArea`のように、より具体的な名前を関数に付けることができます。より適切なソリューションは、この概念をカプセル化する`Rectangle`と呼ばれる独自の**型**を定義することです。


**struct**を使用して単純なタイプを作成できます。[構造体](https://golang.org/ref/spec#Struct_types)は、データを保存できるフィールドの名前付きコレクションです。

このような構造体を宣言します

```go
type Rectangle struct {
    Width float64
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
次のような有用なエラーが表示されるはずです。

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

`Rectangle`を関数に渡すと、意図がより明確に伝わるが、構造体を使用することで得られるメリットが増えることに同意していただければ幸いです。

次の要件は、サークルの`Area`関数を記述することです。

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

ご覧のとおり、 `f` は `g`に置き換えられています。`f`を使用すると、正確な10進数を知るのが難しい場合があります。`g`を使用すると、エラーメッセージで完全な10進数が表示されます。\([fmt options](https://golang.org/pkg/fmt/)\).

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
func Area(circle Circle) float64 { ... }
func Area(rectangle Rectangle) float64 { ... }
```

しかし、**Go**ではできません

`./shapes.go:20:32: Area redeclared in this block`

2つの選択肢があります。

* 同じ名前の関数を異なる`packages`で宣言することができます。新しいパッケージで `Area(Circle)`を作成することはできますが、ここではやりすぎだと感じます。
* 代わりに、新しく定義した型に[メソッド](https://golang.org/ref/spec#Method_declarations)を定義できます。

### メソッドとは?

これまでは _functions_ のみを記述してきましたが、いくつかのメソッドを使用しています。
`t.Errorf`を呼び出すときは、`t` \(`testing.T`\)のインスタンスでメソッド`Errorf`を呼び出しています。

メソッドは、レシーバーを持つ関数です。
メソッド宣言は、識別子（メソッド名）をメソッドにバインドし、メソッドをレシーバーの基本タイプに関連付けます。

メソッドは関数と非常に似ていますが、特定のタイプのインスタンスで呼び出すことによって呼び出されます。 `Area(rectangle)`など、好きな場所で関数を呼び出すことができる場所では、「もの」のメソッドのみを呼び出すことができます。

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

テストを実行しようとすると、

```text
./shapes_test.go:19:19: rectangle.Area undefined (type Rectangle has no field or method Area)
./shapes_test.go:29:16: circle.Area undefined (type Circle has no field or method Area)
```

> タイプCircleにはフィールドまたはメソッドエリアがありません（type Circle has no field or method Area）

ここでコンパイラがどれほど優れているかを繰り返し説明します。時間をかけてゆっくりと表示されるエラーメッセージを読むことは非常に重要です。それは長期的には役立ちます。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

タイプにいくつかのメソッドを追加しましょう

```go
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64  {
    return 0
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64  {
    return 0
}
```

メソッドを宣言するための構文は、関数とほとんど同じです。
これは、メソッドが非常に似ているためです。
唯一の違いは、メソッドレシーバー `func (receiverName ReceiverType) MethodName(args)`の構文です。

そのタイプの変数でメソッドが呼び出されると、 `receiverName`変数を介してそのデータへの参照が取得されます。他の多くのプログラミング言語では、これは暗黙的に行われ、 `this`を介してレシーバーにアクセスします。

Goの慣例では、レシーバー変数をタイプの最初の文字にします。

```go
r Rectangle
```

テストを再実行しようとすると、テストがコンパイルされ、失敗した出力がいくつか表示されます。

## 成功させるのに十分なコードを書く

新しいメソッドを修正して、長方形のテストに成功させましょう

```go
func (r Rectangle) Area() float64  {
    return r.Width * r.Height
}
```

テストを再実行すると、四角形テストはパスするはずですが、円はまだ失敗しているはずです。

サークルの `Area`関数を渡すために、`math`パッケージから `Pi`定数を借ります（インポートすることを忘れないでください\）。

```go
func (c Circle) Area() float64  {
    return math.Pi * c.Radius * c.Radius
}
```

## リファクタリング

テストに重複があります。

やりたいことは、`shapes`のコレクションを取得し、それらの `Area()`メソッドを呼び出して、結果を確認することだけです。

`Rectangle`と` Circle`の両方を渡すことができるある種の `checkArea`関数を記述できるようにしたいが、形状ではないものを渡そうとするとコンパイルに失敗します。

Goでは、この意図を**インターフェース**で体系化できます。

[インターフェイス](https://golang.org/ref/spec#Interface_types)は、Goなどの静的型付き言語で非常に強力な概念です。
これにより、さまざまな型で使用できる関数を作成し、高度に分離されたコードを作成できます。まだタイプセーフを維持しています。

テストをリファクタリングしてこれを紹介しましょう。

```go
func TestArea(t *testing.T) {

    checkArea := func(t *testing.T, shape Shape, want float64) {
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

どのようにして何かが形になりますか？
`Shape`がインターフェース宣言を使用しているものをGoに伝えるだけです

```go
type Shape interface {
    Area() float64
}
```

`Rectangle`と` Circle`で行ったように新しい `type`を作成していますが、今回は`struct`ではなく `interface`です。

これをコードに追加すると、テストに合格します。

### ちょ待って、なぜ？

これは、他のほとんどのプログラミング言語のインターフェースとはかなり異なります。通常、`My type Foo implements interface Bar`と言うコードを書く必要があります。

しかし、私たちの場合

* `Rectangle`には`Area`というメソッドがあり、 `float64`を返すため、`Shape`インターフェースを満たします
* `Circle`には` Area`というメソッドがあり、 `float64`を返すため、`Shape`インターフェースを満たします
* `string`にはそのようなメソッドがないため、インターフェースを満たしていません
* など

Goでは、**インターフェースの解決は暗黙的です**。
渡したタイプがインターフェースが要求するものと一致する場合、それはコンパイルされます。

### 切り離し（Decoupling）

ヘルパーが形状が `Rectangle`、`Circle`、または `Triangle`のどちらであるかを気にする必要がないことに注意してください。
インターフェースを宣言することにより、ヘルパーは具象型から切り離（Decoupling）され、その機能を実行するために必要なメソッドのみを持ちます。

インターフェイスを使用して**必要なもののみ**を宣言するこの種のアプローチは、ソフトウェア設計において非常に重要であり、後のセクションでより詳細に説明します。

## さらにリファクタリング

構造体についてある程度理解できたので、「**テーブル駆動テスト**」を紹介します。

[テーブル駆動テスト](https://github.com/golang/go/wiki/TableDrivenTests)は、同じ方法でテストできるテストケースのリストを作成する場合に役立ちます。

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

ここでの唯一の新しい構文は、「匿名の構造体」`areaTests`を作成することです。 2つのフィールド、 `shape`と` want`で `[]struct`を使用して、構造体のスライスを宣言しています。次に、スライスをケースで埋めます。

次に、構造体フィールドを使用してテストを実行し、他のスライスと同じようにそれらを繰り返します。

開発者が新しい形状を導入し、 `Area`を実装してテストケースに追加するのが非常に簡単であることを確認できます。
さらに、`Area`でバグが見つかった場合、修正する前に新しいテストケースを追加して実行するのは非常に簡単です。

テーブルベースのテストは、ツールボックスの優れた項目になる可能性がありますが、テストで余分なノイズが必要であることを確認してください。
インターフェースのさまざまな実装をテストしたい場合、または関数に渡されるデータに、テストを必要とするさまざまな要件がたくさんある場合、それらは非常に適しています。

別の形状を追加してテストすることで、これらすべてを実証してみましょう。三角形含めて。

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

忘れずに、テストを実行し続けて、コンパイラーに解決策を導きましょう。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`./shapes_test.go:25:4: undefined: Triangle`

三角形はまだ定義していません

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

Triangleは `Area()`メソッドがないため、形状として使用できないので、テストを機能させるために空の実装を追加します

```go
func (t Triangle) Area() float64 {
    return 0
}
```

最後にコードがコンパイルされ、エラーが発生します

`shapes_test.go:31: got 0.00 want 36.00`

## 成功させるのに十分なコードを書く

```go
func (t Triangle) Area() float64 {
    return (t.Base * t.Height) * 0.5
}
```

そして、テストは成功です！

## リファクタリング

繰り返しになりますが、実装は問題ありませんが、テストでは多少の改善が見込めます。

これを見直すと

```go
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

すべての数値が何を表しているのかすぐには明確ではなく、テストを簡単に理解できるようにする必要があります。

ここまでは、 `MyStruct{val1、val2}`構造体のインスタンスを作成するための構文だけを示してきましたが、オプションでフィールドに名前を付けることができます。

それがどのように見えるか見てみましょう

```go
        {shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
        {shape: Circle{Radius: 10}, want: 314.1592653589793},
        {shape: Triangle{Base: 12, Height: 6}, want: 36.0},
```

[例によるテスト駆動開発](https://g.co/kgs/yCzDLF) で、Mr. Kent Beckは、いくつかのテストをある程度までリファクタリングして評価します。

> テストは、それが真実の主張であるかのように、より明確に私たちに話しかけます。**一連の操作ではありません**

今度は、少なくともケースのリストのテストで、形状とその領域について真実を主張します。

## テスト出力が役立つことを確認する

以前に`Triangle`を実装していて、失敗したテストがあったことを覚えていますか？ 

`shapes_test.go:31: got 0.00 want 36.00`と表示されました。

これが`Triangle`に関連していることはわかっていましたが、それを扱っているだけでしたが、表の20のケースのいずれかでバグがシステムに侵入した場合はどうなりますか？
開発者はどのケースが失敗したかをどのようにして知るのでしょうか？
これは開発者にとって素晴らしい経験ではありません。

実際に失敗したケースを見つけるために、ケースを手動で調べる必要があります。

エラーメッセージを `%#v got %.2f want %.2f`に変更できます。 `%#v`形式の文字列は、フィールドの値を含む構造体を出力するため、開発者はテストされているプロパティを一目で確認できます。

テストケースを読みやすくするために、 `want`フィールドの名前を`hasArea`のようなわかりやすい名前に変更できます。

テーブル駆動テストの最後のヒントは、 `t.Run`を使用してテストケースに名前を付けることです。

各ケースを `t.Run`でラップすることで、ケースの名前が出力されるため、失敗時のテスト出力がより明確になります。

```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
```

また、 `go test -run TestArea/Rectangle`を使用して、テーブル内で特定のテストを実行できます。

これを捉えた最終テストコードは次のとおりです

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

これはより基本的な数学の問題の解決策を繰り返し、テストによって動機付けされた新しい言語機能を学習する、よりTDDの実践でした。

* 構造体を宣言して独自のデータ型を作成し、関連するデータをまとめてコードの意図を明確にする
* さまざまなタイプで使用できる関数を定義できるようにインターフェイスを宣言する \([parametric polymorphism](https://en.wikipedia.org/wiki/Parametric_polymorphism)\)
* データ型に機能を追加したり、インターフェースを実装したりできるようにメソッドを追加する
* アサーションをより明確にし、スイートを拡張および保守しやすくするためのテーブルベースのテスト

これは重要な章でした。私たちは今、独自の型を定義し始めているからです。 Goのような静的に型付けされた言語では、理解しやすく、つなぎ合わせてテストできるソフトウェアを構築するために、独自の型を設計できることが不可欠です。

インターフェイスは、システムの他の部分から複雑さを隠すための優れたツールです。私たちの場合、テストヘルパーは、それがアサートしている正確な形状を知る必要はなく、その領域を`尋ねる`方法を知るだけでした。

Goに慣れるにつれて、インターフェースと標準ライブラリの本当の強みを理解し始めることができます。
`everywhere`で使用される標準ライブラリで定義されたインターフェイスについて学び、独自のタイプに対してそれらを実装することにより、多くの優れた機能を非常に迅速に再利用できます。
