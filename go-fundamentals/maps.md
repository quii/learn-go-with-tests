---
description: Maps
---

# マップ

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/maps)

[配列とスライス](arrays-and-slices.md)では、値を順番に格納する方法を見ました。
では、`key`でアイテムを保存し、すばやく検索する方法を見てみましょう。

マップを使用すると、辞書と同じようにアイテムを保存できます。
`key`は単語、`value`は定義と考えることができます。
そして、独自の辞書を構築するよりも、マップについて学ぶより良い方法は何でしょうか？

まず、辞書に定義された単語がすでにあると仮定すると、単語を検索すると、その単語の定義が返されます。

## 最初にテストを書く

`dictionary_test.go`

```go
package main

import "testing"

func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    if got != want {
        t.Errorf("got %q want %q given, %q", got, want, "test")
    }
}
```

マップの宣言は、配列と多少似ています。
例外として、`map`キーワードで始まり、2つのタイプが必要です。
1つはキーのタイプで、`[]`内に記述されます。
2番目は値のタイプで、 `[]`の直後に続きます。

キーのタイプは特別です。 2つのキーが等しいかどうかを判別できないと、正しい値が取得されていることを確認する方法がないため、比較可能な型にしかできません。比較可能な型については、[言語仕様](https://golang.org/ref/spec#Comparison_operators)で詳しく説明しています。

一方、値タイプは任意のタイプにすることができます。別のマップにすることもできます。

このテストの他のすべてはよく知っている必要があります。

## テストを実行してみます

`go test`を実行すると、コンパイラーは「`./dictionary_test.go:8:9: undefined: Search`」で失敗します。

## テストを実行して出力を確認するための最小限のコードを記述します

`dictionary.go`

```go
package main

func Search(dictionary map[string]string, word string) string {
    return ""
}
```

テストは _clearエラーメッセージ_ で失敗するはずです。

`dictionary_test.go:12: got '' want 'this is just a test' given, 'test'`.

## 成功させるのに十分なコードを書く

```go
func Search(dictionary map[string]string, word string) string {
    return dictionary[word]
}
```

マップから値を取得することは、配列 `map[key]`から値を取得することと同じです。

## リファクタリング

```go
func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    assertStrings(t, got, want)
}

func assertStrings(t *testing.T, got, want string) {
    t.Helper()

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

実装をより一般的なものにするために、 `assertStrings`ヘルパーを作成することにしました。

### カスタムタイプを使用する

マップの周りに新しいタイプを作成し、 `Search`をメソッドにすることで、辞書の使用法を改善できます。

`dictionary_test.go`:

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    got := dictionary.Search("test")
    want := "this is just a test"

    assertStrings(t, got, want)
}
```

まだ定義していない `Dictionary`タイプを使い始めました。次に、 `Dictionary`インスタンスで`Search`を呼び出します。

`assertStrings`を変更する必要はありませんでした。

`dictionary.go`:

```go
type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
    return d[word]
}
```

ここでは、 `map`の薄いラッパーとして機能する`Dictionary`タイプを作成しました。
カスタムタイプが定義されたら、`Search`メソッドを作成できます。

## 最初にテストを書く

基本的な検索は非常に簡単に実装できましたが、辞書にない単語を指定するとどうなりますか？

実際には何も返されません。
プログラムは実行し続けることができるのでこれは良いですが、より良いアプローチがあります。
関数は、単語が辞書にないことを報告できます。
このように、ユーザーは単語が存在しないのか、それとも定義がないのか疑問に思うことはありません（これは、辞書にとってはあまり役に立たないように思われるかもしれません。ただし、他のユースケースで重要になる可能性があるシナリオです）。

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    t.Run("known word", func(t *testing.T) {
        got, _ := dictionary.Search("test")
        want := "this is just a test"

        assertStrings(t, got, want)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dictionary.Search("unknown")
        want := "could not find the word you were looking for"

        if err == nil {
            t.Fatal("expected to get an error.")
        }

        assertStrings(t, err.Error(), want)
    })
}
```

Goでこのシナリオを処理する方法は、 `Error`タイプである2番目の引数を返すことです。

`Error`は、`.Error()`メソッドで文字列に変換できます。
これは、アサーションに渡すときに行います。
また、`nil`で`.Error()`を呼び出さないように、`assertStrings`を`if`で保護しています。

## テストを試して実行する

これはコンパイルされません

```text
./dictionary_test.go:18:10: assignment mismatch: 2 variables but 1 values
```

## テストを実行して出力を確認するための最小限のコードを記述します

```go
func (d Dictionary) Search(word string) (string, error) {
    return d[word], nil
}
```

テストは失敗し、より明確なエラーメッセージが表示されます。

`dictionary_test.go:22: expected to get an error.`

## 成功させるのに十分なコードを書く

```go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", errors.New("could not find the word you were looking for")
    }

    return definition, nil
}
```

このパスを作成するために、マップルックアップの興味深いプロパティを使用しています。 2つの値を返すことができます。 2番目の値は、キーが正常に検出されたかどうかを示すブール値です。

このプロパティにより、存在しない単語と定義がない単語を区別できます。

## リファクタリング

```go
var ErrNotFound = errors.New("could not find the word you were looking for")

func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", ErrNotFound
    }

    return definition, nil
}
```

変数に抽出することで、`Search`関数の魔法のエラーを取り除くことができます。
これにより、より良いテストを行うことができます。

```go
t.Run("unknown word", func(t *testing.T) {
    _, got := dictionary.Search("unknown")

    assertError(t, got, ErrNotFound)
})
}

func assertError(t *testing.T, got, want error) {
    t.Helper()

    if got != want {
        t.Errorf("got error %q want %q", got, want)
    }
}
```

新しいヘルパーを作成することで、テストを簡素化し、 `ErrNotFound`変数の使用を開始できるため、将来エラーテキストを変更してもテストが失敗しません。

## Write the test first

辞書を検索するには素晴らしい方法があります。ただし、新しい単語を辞書に追加する方法はありません。

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    dictionary.Add("test", "this is just a test")

    want := "this is just a test"
    got, err := dictionary.Search("test")
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

このテストでは、 `Search`関数を使用して、辞書の検証を少し簡単にします。

## テストを実行して出力を確認するための最小限のコードを記述します

`dictionary.go`

```go
func (d Dictionary) Add(word, definition string) {
}
```

これでテストは失敗するはずです

```text
dictionary_test.go:31: should find added word: could not find the word you were looking for
```

## 成功させるのに十分なコードを書く

```go
func (d Dictionary) Add(word, definition string) {
    d[word] = definition
}
```

マップへの追加も配列に似ています。キーを指定して、値に設定するだけです。

### 参照型

マップの興味深い特性は、マップをポインタとして渡さなくても変更できることです。これは、 `map`が参照型であるためです。つまり、ポインタのように、基礎となるデータ構造への参照を保持します。
基本的なデータ構造は`hash tables`または`hash map`であり、`hash tables`の詳細については[こちら](https://en.wikipedia.org/wiki/Hash_table)を参照してください。

マップがどれほど大きくても、コピーは1つしかないので、マップは参照として非常に適しています。

参照型がもたらす落とし穴は、マップが`nil`値になる可能性があることです。 `nil`マップは読み取り時に空のマップのように動作しますが、`nil`マップに書き込もうとすると、ランタイムパニックが発生します。
マップの詳細については、[こちら](https://blog.golang.org/go-maps-in-action)をご覧ください。

したがって、空のマップ変数を初期化しないでください。

```go
var m map[string]string
```

代わりに、上記のように空のマップを初期化するか、`make`キーワードを使用してマップを作成できます。

```go
var dictionary = map[string]string{}

// OR

var dictionary = make(map[string]string)
```

どちらのアプローチでも空の`hash map`を作成し、`dictionary`を指し示します。これにより、ランタイムパニックが発生することはありません。

## リファクタリング

私たちの実装ではリファクタリングするものは多くありませんが、テストでは少し単純化を使用できます。

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    word := "test"
    definition := "this is just a test"

    dictionary.Add(word, definition)

    assertDefinition(t, dictionary, word, definition)
}

func assertDefinition(t *testing.T, dictionary Dictionary, word, definition string) {
    t.Helper()

    got, err := dictionary.Search(word)
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if definition != got {
        t.Errorf("got %q want %q", got, definition)
    }
}
```

単語と定義の変数を作成し、定義アサーションを独自のヘルパー関数に移動しました。

`Add`は見栄えがしました。ただし、追加しようとしている値が既に存在する場合に何が起こるかは考慮しませんでした。

値がすでに存在する場合、マップはエラーをスローしません。
代わりに、先に進み、新しく提供された値で値を上書きします。これは実際には便利ですが、関数名が正確ではありません。
`Add`は既存の値を変更しません。辞書に新しい単語を追加するだけです。

## 最初にテストを書く

```go
func TestAdd(t *testing.T) {
    t.Run("new word", func(t *testing.T) {
        dictionary := Dictionary{}
        word := "test"
        definition := "this is just a test"

        err := dictionary.Add(word, definition)

        assertError(t, err, nil)
        assertDefinition(t, dictionary, word, definition)
    })

    t.Run("existing word", func(t *testing.T) {
        word := "test"
        definition := "this is just a test"
        dictionary := Dictionary{word: definition}
        err := dictionary.Add(word, "new test")

        assertError(t, err, ErrWordExists)
        assertDefinition(t, dictionary, word, definition)
    })
}
...
func assertError(t *testing.T, got, want error) {
    t.Helper()
    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
    if got == nil {
        if want == nil {
            return
        }
        t.Fatal("expected to get an error.")
    }
}
```

このテストでは、エラーを返すように`Add`を変更しました。
これは、新しいエラー変数`ErrWordExists`に対して検証しています。また、前のテストを変更して、`nil`エラーと`assertError`関数をチェックしました。

## テストを実行してみます

`Add`の値を返さないため、コンパイラは失敗します。

```text
./dictionary_test.go:30:13: dictionary.Add(word, definition) used as value
./dictionary_test.go:41:13: dictionary.Add(word, "new test") used as value
```

## テストを実行して出力を確認するための最小限のコードを記述します

`dictionary.go`

```go
var (
    ErrNotFound   = errors.New("could not find the word you were looking for")
    ErrWordExists = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
    d[word] = definition
    return nil
}
```

これで、さらに2つのエラーが発生します。まだ値を変更しており、 `nil`エラーを返しています。

```text
dictionary_test.go:43: got error '%!q(<nil>)' want 'cannot add word because it already exists'
dictionary_test.go:44: got 'new test' want 'this is just a test'
```

## 成功させるのに十分なコードを書く

```go
func (d Dictionary) Add(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        d[word] = definition
    case nil:
        return ErrWordExists
    default:
        return err
    }

    return nil
}
```

ここでは、エラーを照合するために `switch`ステートメントを使用しています。このような`switch`があると、`Search`が`ErrNotFound`以外のエラーを返す場合に備えて、追加の安全策が提供されます。

## リファクタリング

リファクタリングするものはあまりありませんが、エラーの使用が増えるにつれて、いくつかの変更を加えることができます。

```go
const (
    ErrNotFound   = DictionaryErr("could not find the word you were looking for")
    ErrWordExists = DictionaryErr("cannot add word because it already exists")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}
```

エラーを一定にしました。これには、 `error`インターフェースを実装する独自の` DictionaryErr`タイプを作成する必要がありました。詳細については、[Dave Cheneyによるこの優れた記事](https://dave.cheney.net/2016/04/07/constant-errors)を参照してください。
簡単に言うと、エラーが再利用可能で不変になります。

次に、単語の定義を`Update`する関数を作成しましょう。

## 最初にテストを書く

```go
func TestUpdate(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    dictionary.Update(word, newDefinition)

    assertDefinition(t, dictionary, word, newDefinition)
}
```

`Update`は`Add`と非常に密接に関連しており、次の実装になります。

## テストを試して実行する

```text
./dictionary_test.go:53:2: dictionary.Update undefined (type Dictionary has no field or method Update)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認する

このようなエラーに対処する方法はすでに知っています。関数を定義する必要があります。

```go
func (d Dictionary) Update(word, definition string) {}
```

これを実行すると、単語の定義を変更する必要があることがわかります。

```text
dictionary_test.go:55: got 'this is just a test' want 'new definition'
```

## 成功させるのに十分なコードを書く

`Add`で問題を修正したときに、これを行う方法はすでに見ました。
それでは、`Add`に本当に似たものを実装しましょう。

```go
func (d Dictionary) Update(word, definition string) {
    d[word] = definition
}
```

これは単純な変更だったので、これに必要なリファクタリングはありません。ただし、`Add`と同じ問題が発生しました。
新しい単語を渡すと、 `Update`はそれを辞書に追加します。

## 最初にテストを書く

```go
t.Run("existing word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    newDefinition := "new definition"
    dictionary := Dictionary{word: definition}

    err := dictionary.Update(word, newDefinition)

    assertError(t, err, nil)
    assertDefinition(t, dictionary, word, newDefinition)
})

t.Run("new word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{}

    err := dictionary.Update(word, definition)

    assertError(t, err, ErrWordDoesNotExist)
})
```

単語が存在しない場合のエラータイプをさらに追加しました。また、`Update`を変更して`error`値を返すようにしました。

## テストを試して実行する

```text
./dictionary_test.go:53:16: dictionary.Update(word, "new test") used as value
./dictionary_test.go:64:16: dictionary.Update(word, definition) used as value
./dictionary_test.go:66:23: undefined: ErrWordDoesNotExist
```

今回は3つのエラーが発生しますが、対処方法はわかっています。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
const (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

func (d Dictionary) Update(word, definition string) error {
    d[word] = definition
    return nil
}
```

独自のエラータイプを追加し、`nil`エラーを返しています。

これらの変更により、非常に明確なエラーが発生します。

```text
dictionary_test.go:66: got error '%!q(<nil>)' want 'cannot update word because it does not exist'
```

## 成功させるのに十分なコードを書く

```go
func (d Dictionary) Update(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        d[word] = definition
    default:
        return err
    }

    return nil
}
```

この関数は、`dictionary`を更新したときとエラーを返したときを除いて、`Add`とほとんど同じように見えます。

### 更新（Update）の新しいエラーの宣言に関する注意

`ErrNotFound`を再利用して、新しいエラーを追加することはできません。ただし、更新が失敗したときに正確なエラーを表示する方がよい場合がよくあります。

特定のエラーがあると、何が問題だったかに関する詳細情報が得られます。以下はWebアプリの例です。

> `ErrNotFound`が発生したときにユーザーをリダイレクトできますが、`ErrWordDoesNotExist`が発生したときにエラーメッセージを表示できます。

次に、辞書の単語を削除（`Delete`）する関数を作成しましょう。

## 最初にテストを書く

```go
func TestDelete(t *testing.T) {
    word := "test"
    dictionary := Dictionary{word: "test definition"}

    dictionary.Delete(word)

    _, err := dictionary.Search(word)
    if err != ErrNotFound {
        t.Errorf("Expected %q to be deleted", word)
    }
}
```

このテストでは、単語を含む`Dictionary`を作成し、単語が削除されているかどうかを確認します。

## Try to run the test

`go test`を実行すると、次のようになります。

```text
./dictionary_test.go:74:6: dictionary.Delete undefined (type Dictionary has no field or method Delete)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func (d Dictionary) Delete(word string) {

}
```

これを追加した後、テストは単語を削除しないことを通知します。

```text
dictionary_test.go:78: Expected 'test' to be deleted
```

## 成功させるのに十分なコードを書く

```go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

Goには、マップで機能する組み込み関数`delete`があります。
2つの引数を取ります。1つ目はマップで、2つ目は削除するキーです。

`delete`関数は何も返さず、同じ概念に基づいて`Delete`メソッドを作成しました。存在しない値を削除しても効果がないため、`Update`や`Add`メソッドとは異なり、APIを複雑にしてエラーを発生させる必要はありません。

## まとめ

このセクションでは、多くのことを取り上げました。辞書用に完全なCRUD（作成、読み取り、更新、削除）APIを作成しました。プロセス全体を通じて、次の方法を学びました。

* マップを作成する
* マップ内のアイテムを検索
* マップに新しいアイテムを追加する
* マップのアイテムを更新する
* マップからアイテムを削除する
* エラーの詳細
  * 定数であるエラーを作成する方法
  * エラーラッパーを書く

