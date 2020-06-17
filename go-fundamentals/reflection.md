---
description: Reflection
---

# リフレクション

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/reflection)

[Twitter](https://twitter.com/peterbourgon/status/1011403901419937792?s=09)

> golangチャレンジ：構造体`x`を受け取り、内部にあるすべての文字列フィールドに対して`fn`を呼び出す関数`walk(x interface{}, fn func(string))`を記述します。難易度：再帰的に。

これを行うには、リフレクション（_reflection_）を使用する必要があります。

> コンピューティングにおけるリフレクションは、プログラムが、特にタイプを通じて、独自の構造を調べる能力です。それは一種のメタプログラミングです。また、混乱の元にもなります。

[The Go Blog: Reflection](https://blog.golang.org/laws-of-reflection)抜粋

## 「インターフェース（`interface`）」とは何ですか？

Goでは、`string`、`int`などの既知の型や、`BankAccount`などの独自の型で機能する関数の点で、タイプの安全性を提供してきました。

つまり、自由な値（ドキュメント）を取得し、間違った型を関数に渡そうとするとコンパイラーが文句を言います。

コンパイル時に型がわからない関数を書きたいというシナリオに出くわすかもしれません。

Goでは、これを _any_ 型と考えることができる型`interface{}`で回避できます。

したがって、`walk(x interface{}, fn func(string))`は、`x`の任意の値を受け入れます。

### では、すべてに「インターフェース」を使用し、本当に柔軟な機能を持たないのはなぜでしょうか？
* 「インターフェース`interface`」をとる関数のユーザーとして、タイプの安全性を失います。タイプ`string`の`Foo.bar`を関数に渡すつもりでしたが、代わりに`int`である`Foo.baz`を渡した場合はどうなりますか？コンパイラーは間違いを通知できません。また、関数に渡すことが許可されている _what_ もわかりません。たとえば関数が `UserService`をとることを知ることは非常に便利です。
* そのような関数の書き方として、渡された _anything_ を検査して、型が何であり、それで何ができるのかを理解する必要があります。これは、リフレクション（_reflection_）を使用して行われます。これは非常に不格好で読みにくい場合があり、実行時にチェックを行う必要があるため一般的にパフォーマンスが低下します。

つまり、本当に必要な場合にのみ、リフレクションを使用してください。

ポリモーフィック関数が必要な場合は、ユーザーが関数を機能させるために必要なメソッドを実装している場合に、ユーザーが複数の型で関数を使用できるように、（`interface`ではなく混乱を防ぐために）の周囲で設計できるかどうか検討してください。

私たちの機能は、さまざまなことを処理できる必要があります。いつものように、サポートしたい新しいものごとにテストを作成し、完了するまでリファクタリングを繰り返すというアプローチをとります。

## 最初にテストを書く

文字列フィールドが（`x`）に含まれている構造体で関数を呼び出す必要があります。
次に、渡された関数（`fn`）をスパイして、呼び出されているかどうかを確認できます。

```go
func TestWalk(t *testing.T) {

    expected := "Chris"
    var got []string

    x := struct {
        Name string
    }{expected}

    walk(x, func(input string) {
        got = append(got, input)
    })

    if len(got) != 1 {
        t.Errorf("wrong number of function calls, got %d want %d", len(got), 1)
    }
}
```

* どの文字列が`walk`によって`fn`に渡されたかを格納する文字列のスライス（`got`）を格納したいと思います。多くの場合、前の章では、関数/メソッドの呼び出しをスパイするために専用の型を作成しましたが、この場合は、`got`を閉じる`fn`の匿名関数を渡すだけです。
* 最も単純な「ハッピー」パスを取得するために、文字列型の`Name`フィールドを持つ匿名の`struct`を使用します。
* 最後に、`walk`を`x`とスパイで呼び出し、今のところ`got`の長さを確認するだけです。非常に基本的な動作が得られたら、アサーションでより具体的になります。

## テストを実行してみます

```text
./reflection_test.go:21:2: undefined: walk
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`walk`を定義する必要があります

```go
func walk(x interface{}, fn func(input string)) {

}
```

テストを再試行してください

```text
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
    reflection_test.go:19: wrong number of function calls, got 0 want 1
FAIL
```

## 成功させるのに十分なコードを書く

このパスを作成するために、任意の文字列でスパイを呼び出すことができます。

```go
func walk(x interface{}, fn func(input string)) {
    fn("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}
```

これでテストに合格するはずです。
次に行う必要があるのは、`fn`の呼び出し対象をより具体的にアサートすることです。

## 最初にテストを書く

次のコードを既存のテストに追加して、`fn`に渡された文字列が正しいことを確認します

```go
if got[0] != expected {
    t.Errorf("got %q, want %q", got[0], expected)
}
```

## テストを実行してみます

```text
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
    reflection_test.go:23: got 'I still can't believe South Korea beat Germany 2-0 to put them last in their group', want 'Chris'
FAIL
```

## 成功させるのに十分なコードを書く

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)
    field := val.Field(0)
    fn(field.String())
}
```

このコードは _非常に安全でナイーブ_ ですが、「赤」（テスト失敗）にあるときの目標は、可能な限り最小限のコードを記述することです。次に、懸念に対処するためのテストをさらに記述します。

リフレクションを使用して`x`を確認し、そのプロパティを確認する必要があります。

[reflect package](https://godoc.org/reflect)には、指定された変数の`Value`を返す関数`ValueOf`があります。これには、次の行で使用するフィールドなど、値を検査する方法があります。

次に、渡された値について非常に楽観的な仮定を行います。

* 最初で唯一のフィールドを見て、パニックを引き起こすフィールドがまったくない場合があります。
* 次に、基になる値を文字列として返す`String()`を呼び出しますが、フィールドが文字列以外の場合は間違っていることがわかります。

## リファクタリング

私たちのコードは単純なケースに合格していますが、コードには多くの欠点があることを知っています。

さまざまな値を渡し、`fn`が呼び出された文字列の配列をチェックするいくつかのテストを作成します。

新しいシナリオのテストを続行しやすくするために、テストをテーブルベースのテストにリファクタリングする必要があります。

```go
func TestWalk(t *testing.T) {

    cases := []struct{
        Name string
        Input interface{}
        ExpectedCalls []string
    } {
        {
            "Struct with one string field",
            struct {
                Name string
            }{ "Chris"},
            []string{"Chris"},
        },
    }

    for _, test := range cases {
        t.Run(test.Name, func(t *testing.T) {
            var got []string
            walk(test.Input, func(input string) {
                got = append(got, input)
            })

            if !reflect.DeepEqual(got, test.ExpectedCalls) {
                t.Errorf("got %v, want %v", got, test.ExpectedCalls)
            }
        })
    }
}
```

これで、シナリオを簡単に追加して、複数の文字列フィールドがある場合にどうなるかを確認できます。

## 最初にテストを書く

次のシナリオを`cases`に追加します。

```go
{
    "Struct with two string fields",
    struct {
        Name string
        City string
    }{"Chris", "London"},
    []string{"Chris", "London"},
}
```

## テストを実行してみます

```text
=== RUN   TestWalk/Struct_with_two_string_fields
    --- FAIL: TestWalk/Struct_with_two_string_fields (0.00s)
        reflection_test.go:40: got [Chris], want [Chris London]
```

## 成功させるのに十分なコードを書く

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i:=0; i<val.NumField(); i++ {
        field := val.Field(i)
        fn(field.String())
    }
}
```

`val`には、値のフィールド数を返すメソッド`NumField`があります。
これにより、フィールドを反復処理し、テストに合格した`fn`を呼び出すことができます。

## リファクタリング

ここにコードを改善する明白なリファクターがあるようには見えないので、続けましょう。

`walk`の次の欠点は、すべてのフィールドが`string`であると想定していることです。
このシナリオのテストを書いてみましょう。

## 最初にテストを書く

次のケースを追加

```go
{
    "Struct with non string field",
    struct {
        Name string
        Age  int
    }{"Chris", 33},
    []string{"Chris"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/Struct_with_non_string_field
    --- FAIL: TestWalk/Struct_with_non_string_field (0.00s)
        reflection_test.go:46: got [Chris <int Value>], want [Chris]
```

## 成功させるのに十分なコードを書く

フィールドのタイプが`string`であることを確認する必要があります。

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }
    }
}
```

その[`Kind`](https://godoc.org/reflect#Kind)をチェックすることでそれを行うことができます。

## リファクタリング

繰り返しになりますが、コードは今のところ十分に妥当です。

次のシナリオは、「フラット」な「構造体」でない場合はどうなるのでしょうか。
言い換えると、いくつかのネストされたフィールドを持つ`struct`があるとどうなりますか？

## 最初にテストを書く

私たちは匿名構造体構文を使用して、テストのためにアドホックに型を宣言しているので、そのように続けることができます

```go
{
    "Nested fields",
    struct {
        Name string
        Profile struct {
            Age  int
            City string
        }
    }{"Chris", struct {
        Age  int
        City string
    }{33, "London"}},
    []string{"Chris", "London"},
},
```

しかし、内部の匿名構造体を取得すると、構文が少し乱雑になることがわかります。[構文を改善するために作成する提案があります](https://github.com/golang/go/issues/12854)。

このシナリオの既知のタイプを作成してこれをリファクタリングし、テストで参照してみましょう。
私たちのテストのコードの一部がテストの外にあるという点で少し間接的ですが、読者は初期化を見て、`struct`の構造を推測できるはずです。

次の型宣言をテストファイルのどこかに追加します。

```go
type Person struct {
    Name    string
    Profile Profile
}

type Profile struct {
    Age  int
    City string
}
```

これをケースに追加して、以前よりもはるかに明確に読み取ることができます。

```go
{
    "Nested fields",
    Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/Nested_fields
    --- FAIL: TestWalk/Nested_fields (0.00s)
        reflection_test.go:54: got [Chris], want [Chris London]
```

問題は、型の階層の最初のレベルのフィールドでのみ反復していることです。

## 成功させるのに十分なコードを書く

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }

        if field.Kind() == reflect.Struct {
            walk(field.Interface(), fn)
        }
    }
}
```

解決策は非常に簡単です。その`Kind`をもう一度調べ、それが「構造体`struct`」である場合は、その内部の「構造体`struct`」でもう一度`walk`を呼び出すだけです。

## リファクタリング

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

同じ値の比較を複数回行う場合、「一般的に」`switch`にリファクタリングすると、読みやすさが向上し、コードの拡張が容易になります。

渡された構造体の値がポインターの場合はどうなりますか？

## 最初にテストを書く

このケースを追加

```go
{
    "Pointers to things",
    &Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/Pointers_to_things
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
    panic: reflect: call of reflect.Value.NumField on ptr Value
```

## 成功させるのに十分なコードを書く

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

ポインター`Value`で`NumField`を使用することはできません。`Elem()`を使用する前に、基になる値を抽出する必要があります。

## リファクタリング

与えられた `interface{}`から `reflect.Value`を関数に抽出する責任をカプセル化しましょう。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}

func getValue(x interface{}) reflect.Value {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    return val
}
```

これは実際には _more_ コードを追加しますが、抽象化レベルは適切だと思います。

* 検査できるように、`x`の`reflect.Value`を取得します。方法は気にしません。
* フィールドを反復処理し、そのタイプに応じて必要なことをすべて実行します。

次に、スライスをカバーする必要があります。

## 最初にテストを書く

```go
{
    "Slices",
    []Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/Slices
panic: reflect: call of reflect.Value.NumField on slice Value [recovered]
    panic: reflect: call of reflect.Value.NumField on slice Value
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

これは以前のポインターシナリオに似ています。`reflect.Value`で `NumField`を呼び出そうとしていますが、構造体ではないため、これはありません。

## 成功させるのに十分なコードを書く

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    if val.Kind() == reflect.Slice {
        for i:=0; i< val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
        return
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

## リファクタリング

これは機能しますが、不幸です。心配はいりません。
テストに裏打ちされた実際のコードがあるので、好きなように自由に変更することができます。

少し抽象的に考えると、どちらかで`walk`と呼びたい

* 構造体の各フィールド
* スライス内の各 _thing_

現時点でのコードはこれを実行しますが、十分に反映していません。
最初に、それがスライス（残りのコードの実行を停止するための `return`付き）であるかどうかを確認し、そうでない場合は構造体であると想定します。

コードを書き直して、代わりにタイプ _first_ を確認してから作業を行います。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    switch val.Kind() {
    case reflect.Struct:
        for i:=0; i<val.NumField(); i++ {
            walk(val.Field(i).Interface(), fn)
        }
    case reflect.Slice:
        for i:=0; i<val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
    case reflect.String:
        fn(val.String())
    }
}
```

格好良く！

構造体またはスライスの場合は、それぞれの値に対して`walk`を呼び出してその値を繰り返し処理します。
それ以外の場合、`reflect.String`であれば、`fn`を呼び出すことができます。

それでも、私にはそれがより良いものになり得るような気がします。フィールド/値を反復して`walk`を呼び出すという操作の繰り返しがありますが、概念的には同じです。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

`value`が`reflect.String`の場合、通常のように`fn`を呼び出すだけです。

それ以外の場合、`switch`はタイプに応じて2つのものを抽出します

* いくつのフィールドがありますか
* `Value`を抽出する方法（`Field`または`Index`）

それらを決定したら、`getField`関数の結果を使用して、`numberOfValues`が `walk`を呼び出して反復できるようにします。

これで完了です。配列の処理は簡単です。

## 最初にテストを書く

ケースに追加

```go
{
    "Arrays",
    [2]Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/Arrays
    --- FAIL: TestWalk/Arrays (0.00s)
        reflection_test.go:78: got [], want [London Reykjavík]
```

## 成功させるのに十分なコードを書く

配列はスライスと同じように処理できるため、コンマを使用してケースに追加するだけです。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

次に処理するタイプは `map`です。

## 最初にテストを書く

```go
{
    "Maps",
    map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    },
    []string{"Bar", "Boz"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/Maps
    --- FAIL: TestWalk/Maps (0.00s)
        reflection_test.go:86: got [], want [Bar Boz]
```

## 成功させるのに十分なコードを書く

ここでも少し抽象的に考えると、`map`は`struct`に非常に似ていることがわかります。これは、コンパイル時にキーが不明であるということだけです。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walk(val.MapIndex(key).Interface(), fn)
        }
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

ただし、設計により、インデックスからマップから値を取得することはできません。これは _key_ によってのみ行われるため、抽象化を壊します。

## リファクタリング

今の気分はどうですか？

当初は素晴らしい抽象化のように思われたかもしれませんが、コードは少し不安定に感じられます。

これで問題ありません！
リファクタリングは道のりであり、間違いを犯すこともあります。
TDDの主要なポイントは、これらのことを試す自由を私たちに与えることです。

テストに裏打ちされた小さなステップを踏むことによって、これは決して不可逆的な状況ではありません。
リファクタリング前の状態に戻しましょう。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    walkValue := func(value reflect.Value) {
        walk(value.Interface(), fn)
    }

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        for i := 0; i< val.NumField(); i++ {
            walkValue(val.Field(i))
        }
    case reflect.Slice, reflect.Array:
        for i:= 0; i<val.Len(); i++ {
            walkValue(val.Index(i))
        }
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walkValue(val.MapIndex(key))
        }
    }
}
```

`walk`を導入しました。これは、`val`から`reflect.Value`を抽出するだけでよいように、`switch`内の`walk`への呼び出しを乾燥させます。

### 最後の問題

Goのマップは順序を保証するものではないことに注意してください。したがって、`fn`の呼び出しは特定の順序で行われると断言するため、テストが失敗することがあります。

これを修正するには、マップを含むアサーションを、順序を気にしない新しいテストに移動する必要があります。

```go
t.Run("with maps", func(t *testing.T) {
    aMap := map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    }

    var got []string
    walk(aMap, func(input string) {
        got = append(got, input)
    })

    assertContains(t, got, "Bar")
    assertContains(t, got, "Boz")
})
```

`assertContains`の定義方法は次のとおりです

```go
func assertContains(t *testing.T, haystack []string, needle string)  {
    t.Helper()
    contains := false
    for _, x := range haystack {
        if x == needle {
            contains = true
        }
    }
    if !contains {
        t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
    }
}
```

次に処理したい型は`chan`です。

## 最初にテストを書く

```go
t.Run("with channels", func(t *testing.T) {
        aChannel := make(chan Profile)

        go func() {
            aChannel <- Profile{33, "Berlin"}
            aChannel <- Profile{34, "Katowice"}
            close(aChannel)
        }()

        var got []string
        want := []string{"Berlin", "Katowice"}

        walk(aChannel, func(input string) {
            got = append(got, input)
        })

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v, want %v", got, want)
        }
    })
```

## テストを実行してみます

```text
--- FAIL: TestWalk (0.00s)
    --- FAIL: TestWalk/with_channels (0.00s)
        reflection_test.go:115: got [], want [Berlin Katowice]
```

## 成功させるのに十分なコードを書く

`Recv()`で閉じられるまで、チャネルを通じて送信されたすべての値を反復処理できます

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walk(val.MapIndex(key).Interface(), fn)
        }
    case reflect.Chan:
        for v, ok := val.Recv(); ok; v, ok = val.Recv() {
            walk(v.Interface(), fn)
        }
    }
}
```

次に処理するタイプは`func`です。

## 最初にテストを書く

```go
t.Run("with function", func(t *testing.T) {
        aFunction := func() (Profile, Profile) {
            return Profile{33, "Berlin"}, Profile{34, "Katowice"}
        }

        var got []string
        want := []string{"Berlin", "Katowice"}

        walk(aFunction, func(input string) {
            got = append(got, input)
        })

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v, want %v", got, want)
        }
    })
```

## テストを実行してみます

```text
--- FAIL: TestWalk (0.00s)
    --- FAIL: TestWalk/with_function (0.00s)
        reflection_test.go:132: got [], want [Berlin Katowice]
```

## 成功させるのに十分なコードを書く

このシナリオでは、引数のない関数はあまり意味がありません。ただし、任意の戻り値を許可する必要があります。

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walk(val.MapIndex(key).Interface(), fn)
        }
    case reflect.Chan:
        for v, ok := val.Recv(); ok; v, ok = val.Recv() {
            walk(v.Interface(), fn)
        }
    case reflect.Func:
        valFnResult := val.Call(nil)
        for _, res := range valFnResult {
            walk(res.Interface(), fn)
        }
    }
}
```

## まとめ

* `reflect`パッケージのいくつかの概念を導入しました。
* 任意のデータ構造をたどるために再帰を使用しました。
* 振り返ってみると、悪いリファクタリングをしましたが、それについてはあまり動揺はありません。テストを反復的に行うことで、それほど大したことではありません。
* これは、リフレクションの小さな側面だけをカバーしています。[GOブログでは、詳細を網羅した優れた記事を掲載しています](https://blog.golang.org/laws-of-reflection)。
* リフレクションについて理解したので、使用しないように最善を尽くしてください。
