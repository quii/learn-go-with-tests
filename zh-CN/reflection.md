# 反射

[来自推特](https://twitter.com/peterbourgon/status/1011403901419937792?s=09)

> golang 挑战：编写函数 `walk(x interface{}, fn func(string))`，参数为结构体 `x`，并对 `x` 中的所有字符串字段调用 `fn` 函数。难度级别：递归。

为此，我们需要使用 _反射_。

> 计算中的反射提供了程序检查自身结构体的能力，特别是通过类型，这是元编程的一种形式。这也是造成困惑的一个重要原因。

以上引用自 [Go 博客：反射](https://blog.golang.org/laws-of-reflection)

## 什么是 `interface`？

由于函数使用已知的类型，例如 `string`，`int` 以及我们自己定义的类型，如 `BankAccount`，我们享受到了 Go 为我们提供的类型安全。

这意味着我们可以免费获得一些文档，如果你试图向函数传递错误的类型，编译器就会报错。

但是，你可能会遇到这样的情况，即你不知道要编写的函数参数在编译时是什么类型的。

Go 允许我们使用类型 `interface{}` 来解决这个问题，你可以将其视为 _任意_ 类型。

所以 `walk(x interface{}, fn func(string))` 的 `x` 参数可以接收任何的值。

### 那么为什么不通过将所有参数都定义为 `interface` 类型来得到真正灵活的函数呢？

- 作为函数的使用者，使用 `interface` 将失去对类型安全的检查。如果你想传入 `string` 类型的 `Foo.bar` 但是传入的是 `int` 类型的 `Foo.baz`，编译器将无法通知你这个错误。你也搞不清楚函数允许传递什么类型的参数。知道一个函数接收什么类型，例如 `UserService`，是非常有用的。
- 作为这样一个函数的作者，你必须检查传入的 _所有_ 参数，并尝试断定参数类型以及如何处理它们。这是通过 _反射_ 实现的。这种方式可能相当笨拙且难以阅读，而且一般性能比较差（因为程序必须在运行时执行检查）。

简而言之，除非真的需要否则不要使用反射。

如果你想实现函数的多态性，请考虑是否可以围绕接口（不是 `interface` 类型，这里容易让人困惑）设计它，以便用户可以用多种类型来调用你的函数，这些类型实现了函数工作所需要的任何方法。

我们的函数需要能够处理很多不同的东西。和往常一样，我们将采用迭代的方法，为我们想要支持的每一件新事物编写测试，并一路进行重构，直到完成。

## 首先编写测试

我们想用一个 `struct` 来调用我们的函数，这个 `struct` 中有一个字符串字段（`x`），然后我们可以监视传入的函数（`fn`），看看它是否被调用。

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

- 我们想存储一个字符串切片（`got`），字符串通过 `walk` 传递到 `fn`。在前面的章节中，通常我们会专门为函数或方法调用指定类型，但在这种情况下，我们可以传递一个匿名函数给 `fn`，它会隐藏 `got`。
- 我们使用带有 `string` 类型的 `Name` 字段的匿名 `struct`，以此得到最简单的实现路径。
- 最后调用 `walk` 并传入 `x` 参数，现在只检查 `got` 的长度，一旦有了基本的可以运行的程序，我们的断言就会更加具体。

## 尝试运行测试

```
./reflection_test.go:21:2: undefined: walk
```

## 为测试的运行编写最小量的代码，并检查测试的失败输出

我们需要定义 `walk` 函数

```go
func walk(x interface{}, fn func(input string)) {

}
```

再次尝试运行测试

```
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
    reflection_test.go:19: wrong number of function calls, got 0 want 1
FAIL
```

## 编写足够的代码使测试通过

我们可以使用任意的字符串调用 `fn` 函数来使测试通过。

```go
func walk(x interface{}, fn func(input string)) {
    fn("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}
```

现在测试应该通过了。接下来我们需要做的是对我们的 `fn` 是如何被调用的做一个更具体的断言。

## 首先编写测试

在之前的测试中添加以下代码，检查传入 `fn` 函数的字符串是否正确。

```go
if got[0] != expected {
    t.Errorf("got '%s', want '%s'", got[0], expected)
}
```

## 尝试运行测试

```
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
    reflection_test.go:23: got 'I still can't believe South Korea beat Germany 2-0 to put them last in their group', want 'Chris'
FAIL
```

## 编写足够的代码使测试通过

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)
    field := val.Field(0)
    fn(field.String())
}
```

这段代码 _非常不安全，也非常幼稚_，但请记住，当我们处于「红色」状态（测试失败）时，我们的目标是编写尽可能少的代码。然后我们编写更多的测试来解决我们的问题。

我们需要使用反射来查看 `x` 并尝试查看它的属性。

[反射包](https://godoc.org/reflect)有一个函数 `ValueOf`，该函数值返回一个给定变量的 `Value`。这为我们提供了检查值的方法，包括我们在下一行中使用的字段。

然后我们对传入的值做了一些非常乐观的假设：

- 我们只看第一个也是唯一的字段，可能根本就没有字段会引起 `panic`
- 然后我们调用 `String()`，它以字符串的形式返回底层值，但是我们知道，如果这个字段不是字符串，程序就会出错。

## 重构

我们的代码在简单的测试中可以通过，但是我们也知道代码有很多缺点。

我们将编写一些测试，在这些测试中我们传入不同的值并检查 `fn` 调用的字符串数组。

我们应该将我们的测试重构到一个基于表的测试中，以便更容易地继续测试新的场景。

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

现在，我们可以很容易地添加一个场景，看看如果有多个字符串字段会发生什么。

## 首先编写测试

为测试用例添加以下场景。

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

## 尝试运行测试

```
=== RUN   TestWalk/Struct_with_two_string_fields
    --- FAIL: TestWalk/Struct_with_two_string_fields (0.00s)
        reflection_test.go:40: got [Chris], want [Chris London]
```

## 编写足够的代码使测试通过

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i:=0; i<val.NumField(); i++ {
        field := val.Field(i)
        fn(field.String())
    }
}
```

`value` 有一个方法 `NumField`，它返回值中的字段数。这让我们遍历字段并调用 `fn` 通过我们的测试。

## 重构

这里似乎没有任何明显的重构可以改进代码，让我们继续。

`walk` 的另一个缺点是它假设每个字段都是 `string`。让我们为这个场景编写一个测试。

## 首先编写测试

添加一下测试用例

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

## 尝试运行测试

```
=== RUN   TestWalk/Struct_with_non_string_field
    --- FAIL: TestWalk/Struct_with_non_string_field (0.00s)
        reflection_test.go:46: got [Chris <int Value>], want [Chris]
```

## 编写足够的代码使测试通过

我们需要检查字段的类型是 `string`。

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

我们可以通过检查它的 [Kind](https://godoc.org/reflect#Kind) 来实现这个功能。

## 重构

现在看起来代码已经足够合理了。

下一个场景是，如果它不是一个「平的」`struct` 怎么办？换句话说，如果我们有一个包含嵌套字段的 `struct` 会发生什么？

## 首先编写测试

我们临时使用匿名结构体语法为我们的测试声明类型，所以我们可以继续这样做

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

但我们可以看到，当你得到内部匿名结构时语法会有点混乱。[这里有一个建议可以使它的语法更好](https://github.com/golang/go/issues/12854)。

让我们通过为这个场景创建一个已知类型并在测试中引用它来重构它。有一点间接的地方是，我们测试的一些代码在测试之外，但是读者应该能够通过观察初始化来推断 `struct` 的结构。

在测试文件中添加以下类型声明

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

现在我们将这些添加到测试用例中，它提高了代码的可读性

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

## 尝试运行测试

```
=== RUN   TestWalk/Nested_fields
    --- FAIL: TestWalk/Nested_fields (0.00s)
        reflection_test.go:54: got [Chris], want [Chris London]
```

这个问题是我们只在类型层次结构的第一级上迭代字段导致的。

## 编写足够的代码使测试通过

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

解决方法很简单，我们再次检查它的 `Kind` 如果它碰巧是一个 `struct` 我们就在内部 `struct` 上再次调用 `walk`。

## 重构

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

当你多次对相同的值进行比较时，通常情况下，将代码重构为 `switch` 会提高可读性，使代码更易于扩展。

如果传递进来的结构的值是一个指针呢？

## 首先编写测试

添加这个测试用例

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

## 尝试运行测试

```
=== RUN   TestWalk/Pointers_to_things
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
    panic: reflect: call of reflect.Value.NumField on ptr Value
```

## 编写足够的代码使测试通过

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

指针类型的 `Value` 不能使用 `NumField` 方法，在执行此方法前需要调用 `Elem()` 提取底层值。

## 重构

让我们封装一个获得 `reflect.Value` 的功能，将 `interface{}` 传入函数并返回这个值

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

这实际上增加了更多的代码，但我觉得抽象层是没有问题的

- 得到 `x` 的 `reflect.Value`，这样我就可以检查它，我不在乎怎么做。
- 遍历字段，根据其类型执行任何需要执行的操作。

接下来我们需要覆盖切片。

## 首先编写测试

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

## 尝试运行测试

```
=== RUN   TestWalk/Slices
panic: reflect: call of reflect.Value.NumField on slice Value [recovered]
    panic: reflect: call of reflect.Value.NumField on slice Value
```

## 为测试的运行编写最小量的代码，并检查测试的失败输出

这与前面的指针场景类似，我们试图在 `reflect.Value` 中调用 `NumField`。但它没有，因为它不是结构体。

## 编写足够的代码使测试通过

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

## 重构

这招很管用，但很恶心。不过不用担心，我们有测试支持的工作代码，所以我们可以随意修改我们喜欢的代码。

如果你抽象地想一下，我们想要针对下面的对象调用 `walk`

- 结构体中的每个字段
- 切片中的每一项

我们目前的代码可以做到这一点，但反射用得不太好。我们只是在一开始检查它是否是切片（通过 `return` 来停止执行剩余的代码），如果不是，我们就假设它是 `struct`。

让我们重新编写代码，先检查类型，再执行我们的逻辑代码。

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

看起来好多了！如果是 `struct` 或切片，我们会遍历它的值，并对每个值调用 `walk` 函数。如果是 `reflect.String`，我们就调用 `fn`。

不过，对我来说，感觉还可以更好。这里有遍历字段、值，然后调用 `walk` 的重复操作，但概念上它们是相同的。

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

如果 `value` 是一个 `reflect.String`，我们就像平常一样调用 `fn`。

否则，我们的 `switch` 将根据类型提取两个内容

- 有多少字段
- 如何提取 `Value`（`Field` 或 `Index`）

一旦确定了这些东西，我们就可以遍历 `numberOfValues`，使用 `getField` 函数的结果调用 `walk` 函数。

现在我们已经完成了，处理数组应该很简单了。

## 首先编写测试

添加以下代码到测试用例中：

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

## 尝试运行测试

```
=== RUN   TestWalk/Arrays
    --- FAIL: TestWalk/Arrays (0.00s)
        reflection_test.go:78: got [], want [London Reykjavík]
```

## 编写足够的代码使测试通过

数组的处理方式与切片处理方式相同，因此只需用逗号将其添加到测试用例中

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

我们想处理的最后一个类型是 `map`。

## 首先编写测试

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

## 尝试运行测试

```
=== RUN   TestWalk/Maps
    --- FAIL: TestWalk/Maps (0.00s)
        reflection_test.go:86: got [], want [Bar Boz]
```

## 编写足够的代码使测试通过

如果你抽象地想一下你会发现 `map` 和 `struct` 很相似，只是编译时的键是未知的。

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

然而通过设计，你无法通过索引从 `map` 中获取值。它只能通过 _键_ 来完成，这样就打破了我们的抽象，该死。

## 重构

你现在感觉怎么样？这在当时可能是一个很好的抽象，但现在代码感觉有点不稳定。

_这没问题！_ 重构是一段旅程，有时我们会犯错误。TDD 的一个主要观点是它给了我们尝试这些东西的自由。

通过步步为营的原则就不会发生不可逆转的局面。让我们把它恢复到重构之前的状态。

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

我们已经介绍了 `walkValue`，它依照「Don't repeat yourself」的原则在 `switch` 中调用 `walk` 函数，这样它们就只需要从 `val` 中提取 `reflect.Value` 即可。

### 最后一个问题

记住，Go 中的 `map` 不能保证顺序一致。因此，你的测试有时会失败，因为我们断言对 `fn` 的调用是以特定的顺序完成的。

为了解决这个问题，我们需要将带有 `map` 的断言移动到一个新的测试中，在这个测试中我们不关心顺序。

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

下面是 `assertContains` 是如何定义的

```go
func assertContains(t *testing.T, haystack []string, needle string)  {
    contains := false
    for _, x := range haystack {
        if x == needle {
            contains = true
        }
    }
    if !contains {
        t.Errorf("expected %+v to contain '%s' but it didnt", haystack, needle)
    }
}
```

## 总结

- 介绍了 `reflect` 包中的一些概念。
- 使用递归遍历任意数据结构。
- 在回顾中做了一个糟糕的重构，但不用对此感到太沮丧。通过迭代地进行测试，这并不是什么大问题。
- 这只是 `reflection` 的一个小方面。[Go 博客上有一篇精彩的文章介绍了更多细节](https://blog.golang.org/laws-of-reflection)。
- 现在你已经了解了反射，请尽量避免使用它。

---

作者：[Chris James](https://dev.to/quii)
译者：[Donng](https://github.com/Donng)
校对：[pityonline](https://github.com/pityonline)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
