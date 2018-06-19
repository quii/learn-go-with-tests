# 迭代

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/for)**

在 Go 中 `for` 用来循环和迭代，Go 语言没有 `while`，`do`，`until` 这几个关键字，你只能使用 `for`。这也算是件好事！

让我们来为一个重复字符 5 次的函数写一个测试。

目前这里没什么新知识，所以你可以自己尝试去写。

## 先写测试

```go
package iteration

import "testing"

func TestRepeat(t *testing.T) {
    repeated := Repeat("a")
    expected := "aaaaa"

    if repeated != expected {
        t.Errorf("expected '%s' but got '%s'", expected, repeated)
    }
}
```

## 尝试运行测试

`./repeat_test.go:6:14: undefined: Repeat`

## 先使用最少的代码来让失败的测试先跑起来

_遵守原则！_你现在不需要学习任何新知识就可以让测试恰当地失败。

现在你只需让代码可编译，这样你就可以检查测试用例是否写成。

```go
package iteration

func Repeat(character string) string {
    return ""
}
```

很高兴现在你已经了解足够的 Go 知识来给一些基本的问题写测试，是吧？这意味着你可以随心所欲地使用生产代码，并知道它的行为如你所愿。

`repeat_test.go:10: expected 'aaaaa' but got ''`

## 把代码补充完整，使得它能够通过测试

就像大多数类 C 的语言一样，`for` 语法很不起眼。

```go
func Repeat(character string) string {
    var repeated string
    for i := 0; i < 5; i++ {
        repeated = repeated + character
    }
    return repeated
}
```

与其它语言如 C，Java 或 JavaScript 不同，在 Go 中 `for` 语句前导条件部分并没有圆括号，而且大括号 { } 是必须的。

运行测试应该是通过的。

关于 `for` 循环其它变体请参考[这里](https://gobyexample.com/for)。

## 重构

现在是时候重构并引入另一个构造（construct）`+=` 赋值运算符。

```go
const repeatCount = 5

func Repeat(character string) string {
    var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

`+=` 是自增赋值运算符（Add AND assignment operator），它把运算符右边的值加到左边并重新赋值给左边。它在其它类型也可以使用，比如整数类型。

### 基准测试

在 Go 中编写[基准测试](https://golang.org/pkg/testing/#hdr-Benchmarks)（benchmarks）是该语言的另一个一级特性，它与编写测试非常相似。

```go
func BenchmarkRepeat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Repeat("a")
    }
}
```

你会看到上面的代码和写测试差不多。

`testing.B` 可使你访问隐性命名（cryptically named）`b.N`。

基准测试运行时，代码会运行 `b.N` 次，并测量需要多长时间。

代码运行的次数应该不影响你，框架将决定什么是「好」的值，以便让你获得一些得体的结果。

用 `go test -bench=.` 来运行基准测试。

```text
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000           136 ns/op
PASS
```

以上结果说明运行一次这个函数需要 136 纳秒（在我的电脑上）。这挺不错的。

注意：基准测试默认是顺序运行的。

## 练习

* 修改测试代码，以便调用者可以指定字符重复的次数，然后修复代码
* 写一个 `ExampleRepeat` 来完善你的函数文档
* 看一下 [strings 包](https://golang.org/pkg/strings)。找到你认为可能有用的函数，并对它们编写一些测试。投入时间学习标准库会慢慢得到回报。

## 总结

* 更多的 TDD 练习
* 学习了 `for` 循环
* 学习了如何编写基准测试

---

作者：[Chris James](https://dev.to/quii)
译者：[pityonline](https://github.com/pityonline)
校对：[Donng](https://github.com/Donng)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
