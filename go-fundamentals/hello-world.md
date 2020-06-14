---
description: 'Hello, World'
---

# Hello, World

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/hello-world)

新しい言語での最初のプログラムが [Hello, World](https://en.m.wikipedia.org/wiki/%22Hello,_World!%22_program)であるのは伝統的です。

[前の章](install-go.md#go-environment) では、ファイルの配置場所についてGoがどのように考えられているかを説明しました。

次のパス `$GOPATH/src/github.com/{your-user-id}/hello`にディレクトリを作成します。

したがって、UNIXベースのOSを使用していて、 `$GOPATH` を設定している場合は、次のコマンドでディレクトリを作成できます。 `mkdir -p $GOPATH/src/github.com/$USER/hello`.

以降の章では、コードを好きな名前で新しいフォルダーを作成できます。たとえば、次の章では `$GOPATH/src/github.com/{your-user-id}/integers` にコードを配置することをお勧めします。
このサイトの一部のユーザーは、`learn-go-with-tests/hello`. などのすべての作業用のフォルダーを作成することを好みます。
つまり、フォルダをどのように構成するかはあなた次第です。

このディレクトリに`hello.go`というファイルを作成し、このコードを記述します。実行するには`go run hello.go`と入力します。

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world")
}
```

## 使い方

Goでプログラムを作成すると、その中に `main`関数が定義された` main`パッケージが作成されます。
パッケージは、関連するGoコードをグループ化する方法です。

`func`キーワードは、名前と本体で関数を定義する方法です。

`import "fmt"`では、印刷に使用する `Println`関数を含むパッケージをインポートしています。

## テスト方法

どのようにテストすればよいと思いますか？
「ドメイン」コードを外界から分離することは良いことです。 `fmt.Println`は副作用であり、送信する文字列はドメインです。

テストを簡単にするために、これらの懸念事項を分離しましょう

```go
package main

import "fmt"

func Hello() string {
    return "Hello, world"
}

func main() {
    fmt.Println(Hello())
}
```

`func`を使用して新しい関数を再度作成しましたが、今回は定義に別のキーワード`string`を追加しました。
つまり、この関数は `string`を返します。

ここで、 `hello_test.go`という新しいファイルを作成します。ここで、`Hello`関数のテストを記述します

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    want := "Hello, world"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

説明する前に、コードを実行してみましょう。
端末で`go test`を実行します。合格したはずです！
確認するために、`want`文字列を変更して、意図的にテストを中断してみてください。

複数のテストフレームワークを選択する必要がなく、インストール方法を理解する必要がないことに注意してください。
必要なものはすべて言語に組み込まれており、構文は、これから記述する残りのコードと同じです。

### テストを書く

テストの作成は、関数の作成と同様であり、いくつかのルールがあります。

* `xxx_test.go`のような名前のファイルにある必要があります。
* テスト関数は`Test`という単語で始まる必要があります。
* テスト関数は1つの引数のみをとります。 `t *testing.T`

とりあえず、 `* testing.T`タイプの`t`がテストフレームワークへの`hook`(フック)であることを知っていれば十分なので、失敗したいときに `t.Fail()`のようなことを実行できます。

新しいトピックをいくつか取り上げました。

#### `if`

Goのステートメントが他のプログラミング言語とよく似ている場合。

#### Declaring variables

いくつかの変数を構文`varName := value`で宣言しています。これにより、読みやすくするためにテストでいくつかの値を再利用できます。

#### `t.Errorf`

メッセージを出力してテストに失敗する`t`で`Errorf` _method_ を呼び出しています。
`f`は、プレースホルダー値`％q`に値が挿入された文字列を作成できる形式を表します。
テストを失敗させたとき、それがどのように機能するかは明らかです。

プレースホルダー文字列の詳細については、 [fmt go doc](https://golang.org/pkg/fmt/#hdr-Printing)。
テストでは、 `%q` は値を二重引用符で囲むので非常に便利です。

メソッドと関数の違いについては後で説明します。

### Go ドキュメント

Goのもう1つの質機能はドキュメントです。
`godoc -http :8000`を実行すると、ローカルでドキュメントを起動できます。
[localhost:8000/pkg](http://localhost:8000/pkg) に移動すると、システムにインストールされているすべてのパッケージが表示されます。

標準ライブラリの大部分には、例を含む優れたドキュメントがあります。
[http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/) に移動すると、何が利用できるかを確認する価値があります。

`godoc` コマンドがない場合は、`godoc`を含まない新しいバージョンのGo（`1.14`以降）を使用している可能性があります [no longer including `godoc`](https://golang.org/doc/go1.14#godoc)。
`go get golang.org/x/tools/cmd/godoc`を使用して手動でインストールできます。

### Hello, YOU

これでテストが完了したので、ソフトウェアを安全に反復できます。

In the last example we wrote the test _after_ the code had been written just so you could get an example of how to write a test and declare a function. From this point on we will be _writing tests first_.

Our next requirement is to let us specify the recipient of the greeting.

Let's start by capturing these requirements in a test. This is basic test driven development and allows us to make sure our test is _actually_ testing what we want. When you retrospectively write tests there is the risk that your test may continue to pass even if the code doesn't work as intended.

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Chris")
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

Now run `go test`, you should have a compilation error

```text
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

When using a statically typed language like Go it is important to _listen to the compiler_. The compiler understands how your code should snap together and work so you don't have to.

In this case the compiler is telling you what you need to do to continue. We have to change our function `Hello` to accept an argument.

Edit the `Hello` function to accept an argument of type string

```go
func Hello(name string) string {
    return "Hello, world"
}
```

If you try and run your tests again your `main.go` will fail to compile because you're not passing an argument. Send in "world" to make it pass.

```go
func main() {
    fmt.Println(Hello("world"))
}
```

Now when you run your tests you should see something like

```text
hello_test.go:10: got 'Hello, world' want 'Hello, Chris''
```

We finally have a compiling program but it is not meeting our requirements according to the test.

Let's make the test pass by using the name argument and concatenate it with `Hello,`

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

When you run the tests they should now pass. Normally as part of the TDD cycle we should now _refactor_.

### A note on source control

At this point, if you are using source control \(which you should!\) I would `commit` the code as it is. We have working software backed by a test.

I _wouldn't_ push to master though, because I plan to refactor next. It is nice to commit at this point in case you somehow get into a mess with refactoring - you can always go back to the working version.

There's not a lot to refactor here, but we can introduce another language feature, _constants_.

### Constants

Constants are defined like so

```go
const englishHelloPrefix = "Hello, "
```

We can now refactor our code

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
    return englishHelloPrefix + name
}
```

After refactoring, re-run your tests to make sure you haven't broken anything.

Constants should improve performance of your application as it saves you creating the `"Hello, "` string instance every time `Hello` is called.

To be clear, the performance boost is incredibly negligible for this example! But it's worth thinking about creating constants to capture the meaning of values and sometimes to aid performance.

## Hello, world... again

The next requirement is when our function is called with an empty string it defaults to printing "Hello, World", rather than "Hello, ".

Start by writing a new failing test

```go
func TestHello(t *testing.T) {

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })

    t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })

}
```

Here we are introducing another tool in our testing arsenal, subtests. Sometimes it is useful to group tests around a "thing" and then have subtests describing different scenarios.

A benefit of this approach is you can set up shared code that can be used in the other tests.

There is repeated code when we check if the message is what we expect.

Refactoring is not _just_ for the production code!

It is important that your tests _are clear specifications_ of what the code needs to do.

We can and should refactor our tests.

```go
func TestHello(t *testing.T) {

    assertCorrectMessage := func(t *testing.T, got, want string) {
        t.Helper()
        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    }

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"
        assertCorrectMessage(t, got, want)
    })

    t.Run("empty string defaults to 'World'", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"
        assertCorrectMessage(t, got, want)
    })

}
```

What have we done here?

We've refactored our assertion into a function. This reduces duplication and improves readability of our tests. In Go you can declare functions inside other functions and assign them to variables. You can then call them, just like normal functions. We need to pass in `t *testing.T` so that we can tell the test code to fail when we need to.

`t.Helper()` is needed to tell the test suite that this method is a helper. By doing this when it fails the line number reported will be in our _function call_ rather than inside our test helper. This will help other developers track down problems easier. If you still don't understand, comment it out, make a test fail and observe the test output.

Now that we have a well-written failing test, let's fix the code, using an `if`.

```go
const englishHelloPrefix = "Hello, "

func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return englishHelloPrefix + name
}
```

If we run our tests we should see it satisfies the new requirement and we haven't accidentally broken the other functionality.

### Back to source control

Now we are happy with the code I would amend the previous commit so we only check in the lovely version of our code with its test.

### Discipline

Let's go over the cycle again

* Write a test
* Make the compiler pass
* Run the test, see that it fails and check the error message is meaningful
* Write enough code to make the test pass
* Refactor

On the face of it this may seem tedious but sticking to the feedback loop is important.

Not only does it ensure that you have _relevant tests_, it helps ensure _you design good software_ by refactoring with the safety of tests.

Seeing the test fail is an important check because it also lets you see what the error message looks like. As a developer it can be very hard to work with a codebase when failing tests do not give a clear idea as to what the problem is.

By ensuring your tests are _fast_ and setting up your tools so that running tests is simple you can get in to a state of flow when writing your code.

By not writing tests you are committing to manually checking your code by running your software which breaks your state of flow and you won't be saving yourself any time, especially in the long run.

## Keep going! More requirements

Goodness me, we have more requirements. We now need to support a second parameter, specifying the language of the greeting. If a language is passed in that we do not recognise, just default to English.

We should be confident that we can use TDD to flesh out this functionality easily!

Write a test for a user passing in Spanish. Add it to the existing suite.

```go
    t.Run("in Spanish", func(t *testing.T) {
        got := Hello("Elodie", "Spanish")
        want := "Hola, Elodie"
        assertCorrectMessage(t, got, want)
    })
```

Remember not to cheat! _Test first_. When you try and run the test, the compiler _should_ complain because you are calling `Hello` with two arguments rather than one.

```text
./hello_test.go:27:19: too many arguments in call to Hello
    have (string, string)
    want (string)
```

Fix the compilation problems by adding another string argument to `Hello`

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }
    return englishHelloPrefix + name
}
```

When you try and run the test again it will complain about not passing through enough arguments to `Hello` in your other tests and in `hello.go`

```text
./hello.go:15:19: not enough arguments in call to Hello
    have (string)
    want (string, string)
```

Fix them by passing through empty strings. Now all your tests should compile _and_ pass, apart from our new scenario

```text
hello_test.go:29: got 'Hello, Elodie' want 'Hola, Elodie'
```

We can use `if` here to check the language is equal to "Spanish" and if so change the message

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == "Spanish" {
        return "Hola, " + name
    }

    return englishHelloPrefix + name
}
```

The tests should now pass.

Now it is time to _refactor_. You should see some problems in the code, "magic" strings, some of which are repeated. Try and refactor it yourself, with every change make sure you re-run the tests to make sure your refactoring isn't breaking anything.

```go
const spanish = "Spanish"
const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "

func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == spanish {
        return spanishHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

### French

* Write a test asserting that if you pass in `"French"` you get `"Bonjour, "`
* See it fail, check the error message is easy to read
* Do the smallest reasonable change in the code

You may have written something that looks roughly like this

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == spanish {
        return spanishHelloPrefix + name
    }

    if language == french {
        return frenchHelloPrefix + name
    }

    return englishHelloPrefix + name
}
```

## `switch`

When you have lots of `if` statements checking a particular value it is common to use a `switch` statement instead. We can use `switch` to refactor the code to make it easier to read and more extensible if we wish to add more language support later

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    prefix := englishHelloPrefix

    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    }

    return prefix + name
}
```

Write a test to now include a greeting in the language of your choice and you should see how simple it is to extend our _amazing_ function.

### one...last...refactor?

You could argue that maybe our function is getting a little big. The simplest refactor for this would be to extract out some functionality into another function.

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    default:
        prefix = englishHelloPrefix
    }
    return
}
```

A few new concepts:

* In our function signature we have made a _named return value_ `(prefix string)`.
* This will create a variable called `prefix` in your function.
  * It will be assigned the "zero" value. This depends on the type, for example `int`s are 0 and for strings it is `""`.
    * You can return whatever it's set to by just calling `return` rather than `return prefix`.
  * This will display in the Go Doc for your function so it can make the intent of your code clearer.
* `default` in the switch case will be branched to if none of the other `case` statements match.
* The function name starts with a lowercase letter. In Go public functions start with a capital letter and private ones start with a lowercase. We don't want the internals of our algorithm to be exposed to the world, so we made this function private.

## Wrapping up

Who knew you could get so much out of `Hello, world`?

By now you should have some understanding of:

### Some of Go's syntax around

* Writing tests
* Declaring functions, with arguments and return types
* `if`, `const` and `switch`
* Declaring variables and constants

### The TDD process and _why_ the steps are important

* _Write a failing test and see it fail_ so we know we have written a _relevant_ test for our requirements and seen that it produces an _easy to understand description of the failure_
* Writing the smallest amount of code to make it pass so we know we have working software
* _Then_ refactor, backed with the safety of our tests to ensure we have well-crafted code that is easy to work with

In our case we've gone from `Hello()` to `Hello("name")`, to `Hello("name", "French")` in small, easy to understand steps.

This is of course trivial compared to "real world" software but the principles still stand. TDD is a skill that needs practice to develop but by being able to break problems down into smaller components that you can test you will have a much easier time writing software.

