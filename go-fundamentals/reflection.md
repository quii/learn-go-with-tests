---
description: リフレクション
---

# リフレクション

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/reflection)

[X より](https://x.com/peterbourgon/status/1011403901419937792?s=09)

> golangチャレンジ：構造体`x`を受け取り、内部にあるすべての文字列フィールドに対して`fn`を呼び出す関数`walk(x interface{}, fn func(string))`を記述します。難易度：再帰的に。

これを行うには、リフレクション（_reflection_）を使用する必要があります。

> コンピューティングにおけるリフレクションは、プログラムが、特に型を通じて、自身の構造を調べる能力のことです。それは一種のメタプログラミングです。また、混乱の元にもなります。

[The Go Blog: Reflection](https://blog.golang.org/laws-of-reflection)抜粋

## 「インターフェース（`interface{}`）」とは何ですか？

これまで私たちは、`string`や`int`のような既知の型や、`BankAccount`のような独自の型を扱う関数において、Goが提供する型安全性の恩恵を受けてきました。

これにより、ドキュメントの役割を無料で得られ、間違った型を関数に渡そうとするとコンパイラが文句を言ってくれます。

コンパイル時に型がわからない関数を書きたいというシナリオに出くわすかもしれません。

Goでは、これを _any_ 型と考えることができる型`interface{}`で回避できます（実際、Goでは`any`は`interface{}`の[エイリアス](https://cs.opensource.google/go/go/+/master:src/builtin/builtin.go;drc=master;l=95)です）。

したがって、`walk(x interface{}, fn func(string))`は、`x`の任意の値を受け入れます。

### では、すべてに`interface{}`を使って、本当に柔軟な関数にしないのはなぜでしょうか？
* `interface{}`をとる関数の利用者としては、型安全性を失います。`string`型の`Herd.species`を関数に渡すつもりでしたが、代わりに`int`である`Herd.count`を渡した場合はどうなりますか？コンパイラは間違いを教えてくれません。また、関数に _何を_ 渡してよいのかもわかりません。たとえば関数が `UserService`をとることを知ることは非常に便利です。
* そのような関数を書く側としては、渡された _あらゆるもの_ を検査して、型が何であり、それで何ができるのかを理解する必要があります。これは、リフレクション（_reflection_）を使用して行われます。これは非常に不格好で読みにくい場合があり、実行時にチェックを行う必要があるため一般的にパフォーマンスが低下します。

つまり、本当に必要な場合にのみ、リフレクションを使用してください。

ポリモーフィック関数が必要な場合は、ユーザーが関数を機能させるために必要なメソッドを実装している場合に、ユーザーが複数の型で関数を使用できるように、（紛らわしいですが、`interface{}`ではなく）インターフェースを中心に設計できるかどうか検討してください。

私たちの関数は、さまざまなものを扱える必要があります。いつものように反復的なアプローチをとり、サポートしたい新しいものごとにテストを書き、完成するまで途中でリファクタリングを重ねていきます。

## 最初にテストを書く

文字列フィールドを含む構造体（`x`）を渡して関数を呼び出したいと思います。
そうすれば、渡された関数（`fn`）をスパイして、呼び出されたかどうかを確認できます。

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

* どの文字列が`walk`によって`fn`に渡されたかを格納する文字列のスライス（`got`）を格納したいと思います。これまでの章では、関数/メソッドの呼び出しをスパイするために専用の型を作ることが多かったですが、今回は`got`をクロージャとして捕捉する匿名関数を`fn`として渡すだけで済みます。
* 最も単純な「ハッピー」パスとして、文字列型の`Name`フィールドを持つ匿名の`struct`を使用します。
* 最後に、`walk`を`x`とスパイで呼び出し、今のところ`got`の長さを確認するだけです。ごく基本的な動作ができたら、アサーションをより具体的にしていきます。

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

これを成功させるには、任意の文字列でスパイを呼び出すだけで十分です。

```go
func walk(x interface{}, fn func(input string)) {
    fn("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}
```

これでテストに合格するはずです。
次に行う必要があるのは、`fn`が何を引数に呼び出されているかを、より具体的にアサートすることです。

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

このコードは _非常に安全でなくナイーブ_ ですが、覚えておいてください。「赤」（テスト失敗）にあるときの目標は、可能な限り最小限のコードを記述することです。次に、懸念に対処するためのテストをさらに記述します。

リフレクションを使用して`x`を確認し、そのプロパティを確認する必要があります。

[reflectパッケージ](https://pkg.go.dev/reflect)には、指定された変数の`Value`を返す関数`ValueOf`があります。これには、次の行で使用するフィールドなど、値を検査する方法があります。

次に、渡された値について非常に楽観的な仮定を行います。

* 最初で唯一のフィールドを見ます。しかし、フィールドがまったくない場合もあり、その場合はパニックを引き起こします。
* 次に、基になる値を文字列として返す`String()`を呼び出します。しかし、フィールドが文字列以外の場合、これは間違ったものになります。

## リファクタリング♪

私たちのコードは単純なケースに合格していますが、コードには多くの欠点があることを知っています。

さまざまな値を渡して、`fn`の呼び出しに使われた文字列の配列をチェックするテストをいくつも書いていきます。

新しいシナリオのテストを続行しやすくするために、テストをテーブルベースのテストにリファクタリングする必要があります。

```go
func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"struct with one string field",
			struct {
				Name string
			}{"Chris"},
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

```
{
    "struct with two string fields",
    struct {
        Name string
        City string
    }{"Chris", "London"},
    []string{"Chris", "London"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/struct_with_two_string_fields
    --- FAIL: TestWalk/struct_with_two_string_fields (0.00s)
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
これにより、フィールドを反復処理して`fn`を呼び出すことができ、テストに合格します。

## リファクタリング♪

コードを改善できる明らかなリファクタリングはなさそうなので、先に進みましょう。

`walk`の次の欠点は、すべてのフィールドが`string`であると想定していることです。
このシナリオのテストを書いてみましょう。

## 最初にテストを書く

次のケースを追加

```
{
    "struct with non string field",
    struct {
        Name string
        Age  int
    }{"Chris", 33},
    []string{"Chris"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/struct_with_non_string_field
    --- FAIL: TestWalk/struct_with_non_string_field (0.00s)
        reflection_test.go:46: got [Chris <int Value>], want [Chris]
```

## 成功させるのに十分なコードを書く

フィールドの型が`string`であることを確認する必要があります。

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

その[`Kind`](https://pkg.go.dev/reflect#Kind)をチェックすることでそれを行うことができます。

## リファクタリング♪

繰り返しになりますが、コードは今のところ十分に妥当です。

次のシナリオは、「フラット」な「構造体」でない場合はどうなるのでしょうか。
言い換えると、いくつかのネストされたフィールドを持つ`struct`があるとどうなりますか？

## 最初にテストを書く

私たちは匿名構造体構文を使用して、テストのためにアドホックに型を宣言しているので、そのように続けることができます

```
{
    "nested fields",
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

しかし、匿名構造体を入れ子にすると、構文が少し乱雑になることがわかります。[構文をより良くしようという提案があります](https://github.com/golang/go/issues/12854)。

このシナリオ用に名前付きの型を作ってリファクタリングし、テストからそれを参照するようにしましょう。
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

これでケースを追加でき、以前よりもはるかに読みやすくなります。

```
{
    "nested fields",
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
    --- FAIL: TestWalk/nested_fields (0.00s)
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

## リファクタリング♪

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

渡された構造体の値がポインタの場合はどうなりますか？

## 最初にテストを書く

このケースを追加

```
{
    "pointers to things",
    &Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/pointers_to_things
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
    panic: reflect: call of reflect.Value.NumField on ptr Value
```

## 成功させるのに十分なコードを書く

```go
func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Pointer {
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

ポインタの`Value`に対して`NumField`を使用することはできません。その前に、`Elem()`を使って基になる値を取り出す必要があります。

## リファクタリング♪

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

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	return val
}
```

これは実際にはコードを _増やす_ ことになりますが、抽象化レベルは適切だと思います。

* 検査できるように、`x`の`reflect.Value`を取得します。方法は気にしません。
* フィールドを反復処理し、その型に応じて必要なことを何でも実行します。

次に、スライスをカバーする必要があります。

## 最初にテストを書く

```
{
    "slices",
    []Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/slices
panic: reflect: call of reflect.Value.NumField on slice Value [recovered]
    panic: reflect: call of reflect.Value.NumField on slice Value
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

これは以前のポインタシナリオに似ています。`reflect.Value`に対して`NumField`を呼び出そうとしていますが、構造体ではないため、そのメソッドがありません。

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

## リファクタリング♪

これは動きますが、イマイチです。でも心配はいりません。
テストに裏打ちされた動くコードがあるので、好きなだけ自由にいじることができます。

少し抽象的に考えると、私たちは次のどちらかに対して`walk`を呼び出したいのです。

* 構造体の各フィールド
* スライス内の各 _要素_

現時点でのコードはこれを実行していますが、それをうまく表現できていません。
最初にスライス（残りのコードの実行を停止するための `return`付き）であるかどうかを確認し、そうでなければ構造体であると決めつけているだけです。

コードを書き直して、代わりに _最初に_ 型を確認してから処理を行うようにしましょう。

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

ずっと良くなりました！

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

それ以外の場合、`switch`は型に応じて2つのものを取り出します

* フィールドがいくつあるか
* `Value`を抽出する方法（`Field`または`Index`）

それらが決まれば、`numberOfValues`回反復しながら、`getField`関数の結果を使って`walk`を呼び出せます。

ここまでできれば、配列の処理は簡単なはずです。

## 最初にテストを書く

ケースに追加

```
{
    "arrays",
    [2]Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/arrays
    --- FAIL: TestWalk/arrays (0.00s)
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

次に処理したい型は`map`です。

## 最初にテストを書く

```
{
    "maps",
    map[string]string{
        "Cow": "Moo",
        "Sheep": "Baa",
    },
    []string{"Moo", "Baa"},
},
```

## テストを実行してみます

```text
=== RUN   TestWalk/maps
    --- FAIL: TestWalk/maps (0.00s)
        reflection_test.go:86: got [], want [Moo Baa]
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

ただし、マップは設計上、インデックスで値を取り出すことはできません。 _キー_ でしか取り出せないため、私たちの抽象化が崩れてしまいます。残念。

## リファクタリング♪

今の気分はどうですか？

そのときは良い抽象化に思えたかもしれませんが、今ではコードが少しいびつに感じられます。

_これで良いのです！_
リファクタリングは道のりであり、間違いを犯すこともあります。
TDDの主要なポイントは、これらのことを試す自由を私たちに与えることです。

テストに裏打ちされた小さなステップを踏んでいるので、これは決して取り返しのつかない状況ではありません。
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

`walkValue`を導入しました。これにより`switch`内の`walk`呼び出しの重複がなくなり（DRY）、各ケースでは`val`から`reflect.Value`を取り出すだけで済むようになります。

### 最後の問題

Goのマップは順序を保証するものではないことに注意してください。`fn`の呼び出しが特定の順序で行われることをアサートしているため、テストがときどき失敗します。

これを修正するには、マップを含むアサーションを、順序を気にしない新しいテストに移動する必要があります。

```go
t.Run("with maps", func(t *testing.T) {
	aMap := map[string]string{
		"Cow":   "Moo",
		"Sheep": "Baa",
	}

	var got []string
	walk(aMap, func(input string) {
		got = append(got, input)
	})

	assertContains(t, got, "Moo")
	assertContains(t, got, "Baa")
})
```

`assertContains`の定義方法は次のとおりです

```go
func assertContains(t testing.TB, haystack []string, needle string) {
	t.Helper()
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}
	if !contains {
		t.Errorf("expected %v to contain %q but it didn't", haystack, needle)
	}
}
```

マップを新しいテストに切り出したので、まだ失敗メッセージを見ていません。ここで意図的に`with maps`テストを壊してエラーメッセージを確認し、その後すべてのテストが成功するように元に戻してください。

マップは順序を保証しないため`got`の _順序_ のチェックはあきらめましたが、だからといってその形についてのすべてのチェックをあきらめる必要はありません。`assertContains`だけでは、`walk`がマップのエントリを2回訪れたり、1つを見落として残りのエントリだけをチェックしたりしても、検索する特定の値がどこかに（何回であれ）存在する限り、テストは成功してしまいます。順序とは異なり、`got`の _長さ_ はマップがどの順序で反復されても完全に予測可能なので、それもアサートしましょう。

```go
func assertLength(t testing.TB, got []string, want int) {
	t.Helper()
	if len(got) != want {
		t.Errorf("got %d values but expected %d", len(got), want)
	}
}
```

`with maps`テストの先頭でこれを呼び出すように追加します。

```go
t.Run("with maps", func(t *testing.T) {
	aMap := map[string]string{
		"Cow":   "Moo",
		"Sheep": "Baa",
	}

	var got []string
	walk(aMap, func(input string) {
		got = append(got, input)
	})

	assertLength(t, got, len(aMap))
	assertContains(t, got, "Moo")
	assertContains(t, got, "Baa")
})
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

`Recv()`を使うと、チャネルが閉じられるまで、チャネルを通じて送信されたすべての値を反復処理できます

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
		for i := 0; i < val.NumField(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}
	case reflect.Chan:
		for {
			if v, ok := val.Recv(); ok {
				walkValue(v)
			} else {
				break
			}
		}
	}
}
```

次に処理したい型は`func`です。

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

このシナリオでは、引数を取る関数はあまり意味がなさそうです。ただし、任意の戻り値は許可する必要があります。

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
		for i := 0; i < val.NumField(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}
	case reflect.Chan:
		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
			walkValue(v)
		}
	case reflect.Func:
		valFnResult := val.Call(nil)
		for _, res := range valFnResult {
			walkValue(res)
		}
	}
}
```

## まとめ

* `reflect`パッケージのいくつかの概念を紹介しました。
* 任意のデータ構造をたどるために再帰を使用しました。
* 振り返ってみると悪いリファクタリングをしてしまいましたが、それほど動揺はしませんでした。テストとともに反復的に作業していれば、大したことではありません。
* これは、リフレクションの小さな側面だけをカバーしています。[Goブログでは、詳細を網羅した優れた記事を掲載しています](https://blog.golang.org/laws-of-reflection)。
* リフレクションについて理解したので、使用しないように最善を尽くしてください。

### 既知の制限: 循環参照

自分自身へのポインタを含む構造体（あるいは任意のポインタの循環）を渡すと、`walk`はスタックオーバーフローを起こします。例えば次のとおりです。

```go
type Person struct {
    Name   string
    Friend *Person
}
p := Person{Name: "Alice"}
p.Friend = &p  // cycle!
walk(p, fn)  // fatal error: stack overflow
```

これを修正することは演習として残しておきます。標準的なアプローチは、すでに訪れたポインタのアドレスを記録しておき、再び入ったときにはスキップすることです。`map[uintptr]bool`が必要になるでしょう。`uintptr`はポインタの生の数値アドレスで、ポインタ型の場合は`reflect.Value.Pointer()`で取得できます。このマップは走査全体を通して存在する必要があるため、それをパラメータとして持ち回る内部ヘルパー関数を用意し、公開の`walk`がマップを作成してそれを呼び出すようにするとよいでしょう。
