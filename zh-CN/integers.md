# Integers

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/integers)**

整数（Integers）的工作方式会像你预期的一样，让我们编写一个 add 函数来尝试一下。创建一个名为 `adder_test.go` 的测试文件来写这段代码。

**注意：** 每个目录中的 Go 源文件只能属于一个包（package），请确保你的文件是分开组织的。[这里对此作了一个很好的解释](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)。

## 先写测试

```go
package integers

import "testing"

func TestAdder(t *testing.T) {
    sum := Add(2, 2)
    expected := 4

    if sum != expected {
        t.Errorf("expected '%d' but got '%d'", expected, sum)
    }
}
```

你会注意到，我们使用 `%d` 作为格式字符串，而不是 `%s`。这是因为我们希望它打印一个整数而不是字符串。

还需要注意的是，我们不再使用 main 包，而是定义了一个名为 integers 的包。顾名思义，它将聚集那些处理整数的函数，例如 Add。

## 尝试并运行测试

执行 `go test` 运行测试

检查编译错误

`./adder_test.go:6:9: undefined: Add`

## 为测试的运行编写最少量的代码并检查失败测试的输出

编写足够的代码来满足编译器，*仅此而已* -- 请记住，我们要检查的是我们的测试是否因为正确的原因而失败。

```go
package integers

func Add(x, y int) int {
    return 0
}
```

当你有多个相同类型的参数（在我们的例子中是两个整数），可以将它缩短为 `(x，y int)` 而不是 `(x int, y int)`。

现在运行测试，我们应该很高兴测试能够正确报告哪里错了。

`adder_test.go:10: expected '4' but got '0'`

如果你注意一下就会发现，我们在 [上一节](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/hello-world) 中学习了 *命名返回值*，但没有在这里使用。它通常应该在结果的含义在上下文中不明确时使用，在我们的例子中，`Add` 函数将参数相加是非常明显的。你可以参考 [这个](https://github.com/golang/go/wiki/CodeReviewComments#named-result-parameters) wiki 获取更多细节。

## 编写足够的代码使其通过

从测试驱动开发最严格的意义上来说，我们现在应该编写 *最少的代码以使测试通过*。学究式的程序员可能会这么做

```go
func Add(x, y int) int {
    return 4
}
```

啊哈！再次挫败，测试驱动开发就是一个骗局，对吗?

我们可以编写另一个测试，用一些不同的数字来迫使测试失败，但那感觉就像猫鼠游戏。

一旦我们更熟悉 Go 的语法后，我将介绍一种称为基于属性测试的技巧，它将不再打扰开发者，并帮助你查找程序漏洞。

现在我们正确的修改它

```go
func Add(x, y int) int {
    return x + y
}
```

如果你重新运行测试，它们应该可以通过了。

## 重构

在 *实际* 的代码中，我们并没有太多可以改进的地方。

我们在前面讨论了如何命名返回参数，它出现在文档中，也出现在大多数开发人员的文本编辑器中。

这很好，因为它有助于提高你编写的代码的可用性。最好是用户通过查看类型标记和文档就能了解你代码的用法。

你可以利用注释为函数添加文档，这些将出现在 Go Doc 中，就像你查看标准库的文档一样。

```go
// Add takes two integers and returns the sum of them
func Add(x, y int) int {
    return x + y
}
```

## 示例

如果你真想多测试的更加深入，可以写一些 [示例](https://blog.golang.org/examples)。你将在标准库的文档中找到许多示例。

通常代码示例与实际代码所做的工作相比是过时的，因为它们处于真实的代码之外并且不会被检查。

Go 示例执行起来就像测试一样，所以你可以对用示例来反映代码的实际功能有自信。

作为包的测试套件的一部分，示例会被编译（并可选择性地执行）。

与典型的测试一样，示例是存在于一个包的 \_test.go 文件中的函数。向 `adder_test.go` 文件添加以下 ExampleAdd 函数。

```go
func ExampleAdd() {
    sum := Add(1, 5)
    fmt.Println(sum)
    // Output: 6
}
```

如果你的代码更改了，导致示例不再有效，那么你的构建（build）也将失败。

运行这个包的测试套件，我们可以看到示例函数是在没有我们进一步安排的情况下执行的:

```bash
$ go test -v
=== RUN   TestAdder
--- PASS: TestAdder (0.00s)
=== RUN   ExampleAdd
--- PASS: ExampleAdd (0.00s)
```

请注意，如果删除注释 「//Output: 6」，示例函数将不会执行。虽然函数会被编译，但是它不会执行。

通过添加这段代码，示例将出现在 `godoc` 的文档中，这将使你的代码更容易理解。

## 总结

我们已经介绍了:

- 测试驱动开发工作流的更多实践
- 整数，加法
- 编写更好的文档，以便使用我们代码的人能够快速理解它的用法
- 如何使用我们代码的示例，这些代码作为测试的一部分被检查

---

作者：[Chris James](https://dev.to/quii)
译者：[Donng](https://github.com/Donng)
校对：[flw](https://github.com/flw-cn)、[polaris1119](https://github.com/polaris1119)、[pityonline](https://github.com/pityonline)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出