# Hello World

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/hello-world)**

当使用新的语言写你的第一个程序时，有一个传统就是写 Hello,world。创建一个 `hello.go` 的文件并写入这些代码。输入 `go run hello.go` 去运行它。

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world")
}
```

## 它是如何工作的

当你使用 Go 编写程序时，你将会有一个 `main` 的包，其中定义了 `main` 的函数。 `func` 关键字用来定义带有名称和主体的函数。

通过 `import "fmt"` 我们导入一个包含 `Println` 函数的包，我们用它来打印输出。

## 如何测试

怎么测试这个？
把你「领域」（domain）的代码与外界分开可以避免副作用（side-effects）。 `fmt.Println` 是一种副作用（打印到 stdout），我们发送的字符串是我们的领域。

所以我们把这些问题分开，就更容易测试了

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

我们再次使用 `func` 创建了一个新函数，但是这次我们在定义中添加了另一个关键字 `string`。这意味着这个函数返回一个 `string`。

现在创建一个名为 `hello_test.go` 的新文件，我们将在这里为 `Hello` 函数编写一个测试

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello()
    want := "Hello, world"

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

在解释之前，让我们先运行一下代码。在终端运行 `go test`，它应该已经通过了！为了检验测试，可以通过尝试改变 `want` 字符串来破坏测试。

请注意，你不必在多个测试框架之间进行选择，也不必破译测试 DSL 来编写测试。你需要的一切都内建在语言中，语法与你将要编写的其余代码相同。

### 编写测试

编写测试和写函数很类似，其中有一些规则

- 它需要在一个名为 `xxx_test.go` 的文件中编写
- 测试函数的命名必须从单词 `Test` 开始
- 测试函数只接受一个参数 `t *testing.T`

现在这些信息足以让我们明白，类型为 `*testing.T` 的变量 `t` 是你在测试框架中的 「钩子」，所以你可以在想要失败时执行 `t.Fail()` 之类的操作。

#### 新的知识

`if`

Go 的 `if` 语句非常类似于其他编程语言。

**声明变量**

我们使用语法 `varName := value` 声明了一些变量，它允许我们在测试中重用一些值以获得可读性。

`t.Errorf`

我们正在调用 `t` 的 `Errorf` 方法，该方法将打印一条消息并使测试失败。`f` 表示格式化，它允许我们构建一个字符串，并将值插入占位符值 `%s` 中。当你测试失败时，它能够让你清楚测试是如何工作的。

稍后我们将探讨方法和函数之间的区别。

### Go 文档

Go 的另一个高质量特征是文档化。通过运行 `godoc -http:8000`，可以在本地启动文档。如果你访问 [localhost:8000/pkg](localhost:8000/pkg)，您将看到系统上安装的所有包。

大多数标准库都有优秀的文档和示例。浏览 [http://localhost:8000/pkg/testing/](http://localhost:8000/pkg/testing/) 是非常值得的，去看一下你有什么可以用的。

### Hello, YOU

现在有了测试，就可以安全地迭代我们的软件了。

在上一个示例中，我们在编写代码 *之后* 编写了测试，以便你可以获得如何编写测试和声明函数的示例。从此刻起，我们将 *首先编写测试*。

我们的下一个要求是让我们指定 `greeting` 的接受者。

让我们从在测试中捕获这些需求开始。这是基本的测试驱动开发，允许我们确保我们的测试 *确实* 是测试我们想要的。当你回顾编写测试时，存在一个风险：即使代码没有按照预期工作，测试也可能继续通过。

```go
package main

import "testing"

func TestHello(t *testing.T) {
    got := Hello("Chris")
    want := "Hello, Chris"

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

这时运行 `go test`，你应该会获得一个编译错误

```
./hello_test.go:6:18: too many arguments in call to Hello
    have (string)
    want ()
```

当使用像 Go 这样的静态类型语言时，*聆听编译器* 是很重要的。编译器理解你的代码应该如何合并并工作，这样你就不必再做这些了。

在这种情况下，编译器告诉你需要做什么才能继续。我们必须修改我们的函数 `Hello` 来接受一个参数。

编辑 `Hello` 函数以接受字符串类型的参数

```go
func Hello(name string) string {
    return "Hello, world"
}
```

如果你尝试再次运行测试，`main.go` 将无法编译，因为你没有传递参数。发送「world」让它通过。

```go
func main() {
    fmt.Println(Hello("world"))
}
```

现在，当你运行测试时，你应该看到类似的内容

```
hello_test.go:10: got 'Hello, world' want 'Hello, Chris''
```

我们终于有了一个编译通过的程序，但是根据测试它并没有达到我们的要求。

为了使测试通过，我们使用 name 参数并用 `Hello` 字符串连接它，

```go
func Hello(name string) string {
    return "Hello, " + name
}
```

当你运行测试时，现在应该通过了。通常作为 TDD 周期的一部分，我们现在应该 *重构* 测试。

这里没有太多可重构的，但我们可以介绍另一种语言特性 *常量*。

### 常量

常量的定义如下

```go
const helloPrefix = "Hello, "
```

现在我们可以重构代码

```go
const helloPrefix = "Hello, "

func Hello(name string) string {
    return helloPrefix + name
}
```

重构之后，重新测试，以确保没有破坏任何东西。

常量应该可以提高应用程序的性能，它避免了每次调用 `Hello` 时创建 `"Hello, "` 字符串实例。

显然，对于这个例子来说，性能提升是微不足道的！但是值得考虑的是创建常量来捕获值的含义，有时还可以帮助提高性能。

## 再次回到 Hello, world

下一个需求是当我们的函数用空字符串调用时，它默认为打印 「Hello, World」 而不是 「Hello，」

首先编写一个新的失败测试

```go
func TestHello(t *testing.T) {

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"

        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    })

    t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"

        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    })

}
```

在这里，我们将在我们的测试库中引入另一个工具 -- 子测试。有时，对一个 「事情」进行分组测试是很有用的，然后进行描述不同场景的子测试。

这种方法的好处是，你可以设置在其他测试中也能够使用的共享代码。

当我们检查信息是否符合预期时，会有重复的代码。

重构不 *仅仅* 是为了生产代码!

重要的是，你的测试 *清楚地说明* 了代码需要做什么。

我们可以并且应该重构我们的测试。

```go
func TestHello(t *testing.T) {

    assertCorrectMessage := func(t *testing.T, got, want string) {
        t.Helper()
        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    }

    t.Run("saying hello to people", func(t *testing.T) {
        got := Hello("Chris")
        want := "Hello, Chris"
        assertCorrectMessage(t, got, want)
    })

    t.Run("empty string defaults to 'world'", func(t *testing.T) {
        got := Hello("")
        want := "Hello, World"
        assertCorrectMessage(t, got, want)
    })

}
```

我们在这里做了什么?

我们将断言重构为函数。这减少了重复，提高了测试的可读性。在 go 中，你可以在其他函数中声明函数并将它们分配给变量。你可以像普通函数一样调用它们。我们需要传入 `t *testing.T`，这样我们就可以在需要的时候令测试代码失败。

`t.Helper()` 需要告诉测试套件这个方法是助手。通过这样做，当助手失败时所报告的行号将在函数调用中而不是在测试助手内部。这将帮助其他开发人员更容易地跟踪问题。如果你仍然不理解，请注释掉它，使测试失败并观察测试输出。

现在我们有了一个写得很好的失败测试，让我们修复代码。

```go
const helloPrefix = "Hello, "

func Hello(name string) string {
    if name == "" {
        name = "World"
    }
    return helloPrefix + name
}
```

如果我们运行测试，应该看到它满足了新的要求，并且我们没有意外地破坏其他功能。

### 规律

让我们再次回顾一下这个周期

- 写一个测试
- 让编译器通过
- 运行测试，查看失败原因并检查错误消息是很有意义的
- 编写足够的代码以使测试通过
- 重构

从表面上看，这可能看起来很乏味，但坚持反馈循环非常重要。

它不仅确保你有 *相关的测试*，还可以确保你通过重构测试的安全性来 *设计优秀的软件*。

看到测试失败是一个重要的检查手段，因为它还可以让你看到错误信息。作为一名开发人员，如果测试失败时不能清楚地说明问题所在，那么使用这个代码库可能会非常困难。

通过确保你的测试 *快速* 并设置你的工具，以便运行测试足够简单，你在编写代码时就可以进入流畅的状态。

如果不写测试，你承诺通过运行你的软件来手动检查你的代码，这会打破你的流畅状态，而且你任何时候都无法将自己从这种状态中拯救出来，尤其是从长远来看。

## 继续前进！更多需求

天呐，我们有更多的需求了。 我们现在需要支持第二个参数，指定问候的语言。 如果一种我们不能识别的语言被传进来，就默认为英语。

我们应该确信，我们可以使用 TDD 轻松实现这一功能！

为通过西班牙语的用户编写测试，将其添加到现有套件。

```go
 t.Run("in Spanish", func(t *testing.T) {
        got := Hello("Elodie", "Spanish")
        want := "Hola, Elodie"
        assertCorrectMessage(t, got, want)
    })
```

记住不要作弊！*先写测试*。 当你尝试运行测试时，编译器 *应该* 会出错，因为你用两个参数而不是一个来调用 `Hello`。

```
./hello_test.go:27:19: too many arguments in call to Hello
    have (string, string)
    want (string)
```

通过向 `Hello` 添加另一个字符串参数来修复编译问题

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }
    return helloPrefix + name
}
```

当你尝试再次运行测试时，它会抱怨在其他测试和 `main.go` 中没有传递足够的参数给 `Hello`

```
./hello.go:15:19: not enough arguments in call to Hello
    have (string)
    want (string, string)
```

通过传递空字符串来修复它们。现在，除了我们的新场景外，你的所有测试都应该编译并通过

```
hello_test.go:29: got 'Hola, Elodie' want 'Hello, Elodie'
```

这里我们可以使用 `if` 检查语言是否是「西班牙语」，如果是就修改信息

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == "Spanish" {
        return "Hola, " + name
    }

    return helloPrefix + name
}
```

测试现在应该通过了。

现在是 *重构* 的时候了。你应该在代码中看出了一些问题，其中有一些重复的「魔术」字符串。 自己尝试重构它，每次更改都要重新运行测试，以确保重构不会破坏任何内容。

```go
const spanish = "Spanish"
const helloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "

func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    if language == spanish {
        return spanishHelloPrefix + name
    }

    return helloPrefix + name
}
```

### 法语

- 写一个测试，断言如果你传递 `"French"` 你会得到 `"Bonjour, "`
- 看到它失败，检查易读的错误消息
- 在代码中进行最小的合理更改

你可能写了一些看起来大致如此的东西

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

    return helloPrefix + name
}
```

### `switch`

当你有很多 `if` 语句检查一个特定的值时，通常使用 `switch` 语句来代替。 如果我们希望稍后添加更多的语言支持，我们可以使用 `switch` 来重构代码，以便更易于阅读和扩展

```go
func Hello(name string, language string) string {
    if name == "" {
        name = "World"
    }

    prefix := helloPrefix

    switch language {
    case french:
        prefix = frenchHelloPrefix
    case spanish:
        prefix = spanishHelloPrefix
    }

    return prefix + name
}
```

写一个测试，添加一个用你选择的语言写的问候，你应该看到它是多么简单，去扩展我们 *惊人* 的功能。

### 最后一次重构？

你可能会争辩说，也许我们的功能有点大。对此最简单的重构是将一些功能提取到另一个函数中，并且你已经知道如何声明函数。

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
        prefix = englishPrefix
    }
    return
}
```

一些新的概念：

- 在我们的函数签名中，我们做了一个 *命名返回值*（`prefix string`）。
- 这将在你的函数中创建一个名为 `prefix` 的变量
  - 它将被分配「零」值。 这取决于类型，例如 int 是 0，对于字符串它是`""`
    - 你只需调用 `return` 而不是 `return prefix` 即可返回所设置的值。

  - 这将显示在 Go Doc 中，以便你的代码更清晰。

- 如果没有其他 `case` 语句匹配，将会执行 `default` 分支

- 函数名称以小写字母开头。在 Go 中，公共函数以大写字母开始，私有函数以小写字母开头。我们不希望我们算法的内部结构暴露给外部，所以我们将这个功能私有化。

## 总结

谁会知道你可以从 `Hello, world` 中学到这么多东西呢？

现在你应该有一些了解

### Go 的一些语法

- 编写测试
- 用参数和返回类型声明函数
- `if`，`else`，`switch`
- 声明变量和常量

### 了解 TDD 过程以及步骤的重要性

- *编写一个失败的测试，并查看失败信息*，可以看到我们已经为需求写了一个 *相关* 的测试，并且看到它产生了一个 *易于理解的失败描述*
- 编写最少量的代码以使其通过，因此我们知道我们有可工作软件
- *然后* 重构，支持我们测试的安全性，以确保我们拥有易于使用的精心制作的代码

在我们的例子中，我们通过小巧易懂的步骤从 `Hello()` 到 `Hello("name")`，到 `Hello("name", "french")`。

与「现实世界」软件相比，这当然是微不足道的，但原则依然成立。TDD 是一门需要通过开发去实践的技能，但通过将问题分解成更小的可测试的组件，编写软件的时间将更加轻松。

----------------

作者：[Chris James](https://dev.to/quii)
译者：[Donng](https://github.com/Donng)
校对：[polaris1119](https://github.com/polaris1119)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
