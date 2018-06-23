# 数组与切片

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/arrays)**

数组允许你以特定的顺序在变量中存储相同类型的多个元素。

对于数组来说，最常见的就是迭代数组中的元素。我们创建一个 `Sum` 函数，它使用 [`for`](for) 来循环获取数组中的元素并返回所有元素的总和。

让我们使用 TDD 思想。

## 先写测试函数

在 `sum_test.go` 中：

```go
package main

import "testing"

func TestSum(t *testing.T) {

    numbers := [5]int{1, 2, 3, 4, 5}

    got := Sum(numbers)
    want := 15

    if want != got {
        t.Errorf("got %d want %d given, %v", got, want, numbers)
    }
}
```

数组的容量是我们在声明它时指定的固定值。我们可以通过两种方式初始化数组：

* \[N\]type{value1, value2, ..., valueN} e.g. `numbers := [5]int{1, 2, 3, 4, 5}`
* \[...\]type{value1, value2, ..., valueN} e.g. `numbers := [...]int{1, 2, 3, 4, 5}`

在错误信息中打印函数的输入有时很有用。我们使用 `%v`（默认输出格式）占位符来打印输入，它非常适用于展示数组。

[了解更多关于格式化字符串的信息](https://golang.org/pkg/fmt/)

## 运行测试

使用 `go test` 运行测试将会报编译时错误：`./sum_test.go:10:15: undefined: Sum`。

## 先使用最少的代码来让失败的测试先跑起来：

`Sum.go`

```go
package main

func Sum(numbers [5]int) (sum int) {
    return 0
}
```

这时测试还会失败，不过会返回明确的错误信息：

`sum_test.go:13: got 0 want 15 given, [1 2 3 4 5]`

## 把代码补充完整，使得它能够通过测试：

```go
func Sum(numbers [5]int) int {
    sum := 0
    for i := 0; i < 5; i++ {
        sum += numbers[i]
    }
    return sum
}
```

可以使用 `array[index]` 语法来获取数组中指定索引对应的值。在本例中我们使用 `for` 循环分 5 次取出数组中的元素并与 `sum` 变量求和。

### 一个源码版本控制的小贴士

此时如果你正在使用源码的版本控制工具（你应该使用它！），我会在此刻先提交一次代码。因为我们已经拥有了一个有测试支持的程序。

但我 *不会* 将它推送到远程的 master 分支，因为我马上就会重构它。在此时提交一次代码是一种很好的习惯。因为你可以在之后重构导致的代码乱掉时回退到当前版本。

* 你总是能够回到这个可用的版本。

## 重构

我们可以使用 [`range`](https://gobyexample.com/range) 语法来让函数变得更加整洁。

```go
func Sum(numbers [5]int) int {
    sum := 0
    for _, number := range numbers {
        sum += number
    }
    return sum
}
```

`range` 会迭代数组，每次迭代都会返回数组元素的索引和值。我们选择使用 `_` [空白标志符](https://golang.org/doc/effective_go.html#blank) 来忽略索引。

### 回到版本控制

现在我们已经重构了之前版本的代码，我们只需要使用之前版本的测试来检查它是否能够通过测试。

### 数组和它的类型

数组有一个有趣的属性，它的大小也属于类型的一部分，如果你尝试将 `[4]int` 作为 `[5]int` 类型的参数传入函数，是不能通过编译的。它们是不同的类型，就像尝试将 `string` 当做 `int` 类型的参数传入函数一样。

因为这个原因，所以数组比较笨重，大多数情况下我们都不会使用它。

Go 的切片（slice）类型不会将集合的长度保存在类型中，因此它的尺寸可以是不固定的。

下面我们会完成一个动态长度的 `Sum` 函数。

## 先写测试

我们会使用 [切片类型](slice)，它可以接收不同大小的切片集合。语法上和数组非常相似，只是在声明的时候不指定长度：

`mySlice := []int{1,2,3}` 而不是 `mySlice := [3]int{1,2,3}`

```go
func TestSum(t *testing.T) {

    t.Run("collection of 5 numbers", func(t *testing.T) {
        numbers := [5]int{1, 2, 3, 4, 5}

        got := Sum(numbers)
        want := 15

        if got != want {
            t.Errorf("got %d want %d given, %v", got, want, numbers)
        }
    })

    t.Run("collection of any size", func(t *testing.T) {
        numbers := []int{1, 2, 3}

        got := Sum(numbers)
        want := 6

        if got != want {
            t.Errorf("got %d want %d given, %v", got, want, numbers)
        }
    })

}
```

## 运行测试

编译出错：

`./sum_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Sum`。

## 先使用最少的代码来让失败的测试先跑起来

这时我们可以选择一种解决方案：

* 修改现有的 API，将 `Sum` 函数的参数从数组改为切片。如果这么做我们就有可能会影响使用这个 API 的人，因为我们的 *其他* 测试不能编译通过。
* 创建一个新函数。

根据目前的情况，并没有人使用我们的函数，所以选择修改原来的函数。

```go
func Sum(numbers []int) int {
    sum := 0
    for _, number := range numbers {
        sum += number
    }
    return sum
}
```

如果你运行测试，它们还是不能编译通过，你必须把之前测试代码中的数组换成切片。

## 把 `Sum` 补充完整，使得它能够通过测试：

事实证明，这里需要我们做的只是修复编译器错误，然后测试就通过了。

## 重构

我们已经重构了 `Sum` 函数把参数从数组改为切片。注意不要在重构以后忘记维护你的测试代码。

```go
func TestSum(t *testing.T) {

    t.Run("collection of 5 numbers", func(t *testing.T) {
        numbers := []int{1, 2, 3, 4, 5}

        got := Sum(numbers)
        want := 15

        if got != want {
            t.Errorf("got %d want %d given, %v", got, want, numbers)
        }
    })

    t.Run("collection of any size", func(t *testing.T) {
        numbers := []int{1, 2, 3}

        got := Sum(numbers)
        want := 6

        if got != want {
            t.Errorf("got %d want %d given, %v", got, want, numbers)
        }
    })

}
```

质疑测试的价值是非常重要的。测试并不是越多越好，而是尽可能的使你的代码更加健壮。太多的测试会增加维护成本，因为 **维护每个测试都是需要成本的**。

在本例中，针对该函数写两个测试其实是多余的，因为切片尺寸并不影响函数的运行。

Go 有内置的计算测试 [覆盖率的工具](https://blog.golang.org/cover)，它能帮助你发现没有被测试过的区域。我们不需要追求 100% 的测试覆盖率，它只是一个供你获取测试覆盖率的方式。只要你严格遵循 TDD 规范，那你的测试覆盖率就会很接近 100%。

运行：

`go test -cover`

你会看到：

```
PASS
coverage: 100.0% of statements
```

现在删除一个测试，然后再次运行。

你现在可以提交一次代码，然后再进行接下来的事情。

这回我们需要一个 `SumAll` 函数，它接受多个切片，并返回由每个切片元素的总和组成的新切片。

例如：

`SumAll([]int{1,2}, []int{0,9})` would return `[]int{3, 9}`

或者

`SumAll([]int{1,1,1})` would return `[]int{3}`

## 先写测试

```go
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1,2}, []int{0,9})
    want := []int{3, 9}

    if got != want {
        t.Errorf("got %v want %v", got, want)
    }
}
```

## 运行测试

`./sum_test.go:23:9: undefined: SumAll`

## 先使用最少的代码来让失败的测试先跑起来

我们需要定义满足测试要求的 `SumAll`。

我们可以写一个 [可变参数](https://gobyexample.com/variadic-functions) 的函数：

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
    return
}
```

这时运行测试会报编译时错误：

`./sum_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

在 Go 中不能对切片使用等号运算符。你可以写一个函数迭代每个元素来检查它们的值。但是一种比较简单的办法是使用 [`reflect.DeepEqual`](deepEqual)，它在判断两个变量是否相等时十分有用。

```go
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1,2}, []int{0,9})
    want := []int{3, 9}

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

确保你已经在文件头部 `import reflect`，这样你才能使用 `DeepEqual` 方法。

需要注意的是 `reflect.DeepEqual` 不是「类型安全」的，所以有时候会发生比较怪异的行为。为了看到这种行为，暂时将测试修改为：

```go
func TestSumAll(t *testing.T) {

    got := SumAll([]int{1,2}, []int{0,9})
    want := "bob"

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

这里我们尝试比较 `slice` 和 `string`。这显然是不合理的，但是却通过了测试！所以使用 `reflect.DeepEqual` 比较简洁但是在使用时需多加小心。

回到我们的测试中。运行测试会得到以下信息：

`sum_test.go:30: got [] want [3 9]`

## 将代码补充完整使函数能够测试通过

我们需要做的就是迭代可变参数，使用 `Sum` 计算每个参数的总和并把结果放入函数返回的切片中。

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
    lengthOfNumbers := len(numbersToSum)
    sums = make([]int, lengthOfNumbers)

    for i, numbers := range numbersToSum {
        sums[i] = Sum(numbers)
    }

    return
}
```

我们学到了很多新东西。

这里有一种创建切片的新方式。`make` 可以在创建切片的时候指定我们需要的长度和容量。

我们可以使用切片的索引访问切片内的元素，使用 `=` 对切片元素进行赋值。

现在应该可以测试通过。

## 重构

顺便说一下，切片有容量的概念。如果你有一个容量为 2 的切片，但使用 `mySlice[10]=1` 进行赋值，会报运行时错误。

不过你可以使用 `append` 函数，它能为切片追加一个新值。

```go
func SumAll(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
        sums = append(sums, Sum(numbers))
    }

    return sums
}
```

在这个实现中，我们不用担心切片元素会超过容量。我们开始使用空切片（在函数签名中定义），在每次计算完切片的总和后将结果添加到切片中。

接下来的工作是把 `SumAll` 变成 `SumAllTails`。它会把每个切片的尾部元素想加（尾部的意思就是出去第一个元素以外的其他元素）。

## 先写测试

```go
func TestSumAllTails(t *testing.T) {
    got := SumAllTails([]int{1,2}, []int{0,9})
    want := []int{2, 9}

    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}
```

## 运行测试

`./sum_test.go:26:9: undefined: SumAllTails`

## 先使用最少的代码来让失败的测试先跑起来

把函数名称改为 `SumAllTails` 并重新运行测试

`sum_test.go:30: got [3 9] want [2 9]`

## 将代码补充完整使函数能够测试通过

```go
func SumAllTails(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
        tail := numbers[1:]
        sums = append(sums, Sum(tail))
    }

    return sums
}
```

我们可以使用语法 `slice[low:high]` 获取部分切片。如果在冒号的一侧没有数字就会一直取到最边缘的元素。在我们的函数中，我们使用 `numbers[1:]` 取到从索引 1 到最后一个元素。你可能需要花费一些时间才能熟悉切片的操作。

## 重构

这次没有太多需要重构的地方。

如果传入一个空切片会怎样？空切片的「尾部」是什么呢，如果我们在空数组上使用 `myEmptySlice[1:]` 会发生什么？

## 先写测试

```go
func TestSumAllTails(t *testing.T) {

    t.Run("make the sums of some slices", func(t *testing.T) {
        got := SumAllTails([]int{1,2}, []int{0,9})
        want := []int{2, 9}

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    })

    t.Run("safely sum empty slices", func(t *testing.T) {
        got := SumAllTails([]int{}, []int{3, 4, 5})
        want :=[]int{0, 9}

        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    })

}
```

## 运行测试

```text
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range
```

值得注意的是，该函数 *编译通过* 了，但是在运行时出现错误。

编译时错误是我们的朋友，因为它帮助我们让程序可以工作。运行时错误是我们的敌人，因为它影响我们的用户。

## 将代码补充完整使函数能够测试通过

```go
func SumAllTails(numbersToSum ...[]int) []int {
    var sums []int
    for _, numbers := range numbersToSum {
        if len(numbers) == 0 {
            sums = append(sums, 0)
        } else {
            tail := numbers[1:]
            sums = append(sums, Sum(tail))
        }
    }

    return sums
}
```

## 重构

我们的测试代码有一部分是重复的，我们可以把它放到另一个函数中复用。

```go
func TestSumAllTails(t *testing.T) {

    checkSums := func(t *testing.T, got, want []int) {
        if !reflect.DeepEqual(got, want) {
            t.Errorf("got %v want %v", got, want)
        }
    }

    t.Run("make the sums of tails of", func(t *testing.T) {
        got := SumAllTails([]int{1, 2}, []int{0, 9})
        want := []int{2, 9}
        checkSums(t, got, want)
    })

    t.Run("safely sum empty slices", func(t *testing.T) {
        got := SumAllTails([]int{}, []int{3, 4, 5})
        want := []int{0, 9}
        checkSums(t, got, want)
    })

}
```

这样使用起来更加方便，而且还能增加代码的类型安全性。如果一个粗心的开发者使用 `checkSums(t, got, "dave")` 是不能通过编译的。

```
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## 总结

我们学习了：

* 数组
* 切片
  * 多种方式的切片初始化
  * 切片的容量是 *固定* 的，但是你可以使用 `append` 从原来的切片中创建一个新切片
  * 如何获取部分切片
* 使用 `len` 获取数组和切片的长度
* 使用测试代码覆盖率的工具
* `reflect.DeepEqual` 的妙用和对代码类型安全性的影响

数组和切片的元素可以是任何类型，包括数组和切片自己。如果需要你可以定义 `[][]string` 。

[Go 官网博客中关于切片的文章](blog-slice) 可以让你更加深入的了解切片。尝试写更多的测试来从中学到东西。

另一种练习 Go 的方式是在 Go 的在线编译器中写代码。几乎所有东西都可以写在上面，而且如果你想问问题，它可以让你的代码很容易分享给其他人。[为了你方便试验，我已经在 go playground 中写好了一个 slice 的示例](https://play.golang.org/p/ICCWcRGIO68)

[for]: https://github.com/quii/learn-go-with-tests/blob/master/iteration.md
[blog-slice]: https://blog.golang.org/go-slices-usage-and-internals
[deepEqual]: https://golang.org/pkg/reflect/#DeepEqual
[slice]: https://golang.org/doc/effective_go.html#slices

---

作者：[Chris James](https://dev.to/quii)
译者：[saberuster](https://github.com/saberuster)
校对：[polaris1119](https://github.com/polaris1119)、[Donng](https://github.com/Donng)、[pityonline](https://github.com/pityonline)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
