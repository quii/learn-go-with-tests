# Maps

[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/maps)

在[数组和切片](arrays-and-slices.md)的章节中，你学会了如何按顺序存储值。现在，我们再来看看如何通过`键`存储值，并快速查找它们。

Maps 允许你以类似于字典的方式存储值。你可以将`键`视为单词，将`值`视为定义。 难道还有比构建我们自己的字典更好的学习 Maps 的方式吗？

## 首先编写测试

在 `dictionary_test.go` 中编写代码：

```go
package main

import "testing"

func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    if got != want {
        t.Errorf("got '%s' want '%s' given, '%s'", got, want, "test")
    }
}
```

声明 Map 的方式有点儿类似于数组。不同之处是，它以 `map` 关键字开头，需要两种类型。 第一个是键，它写在 `[]` 中。 第二个是值，它在 `[]` 之后。

键很特别， 它只能是一个可比较的类型，因为如果不能判断两个键是否相等，我们就无法确保我们得到的是正确的值。可比类型在[语言规范](https://golang.org/ref/spec#Comparison_operators)中有详细解释。

另一方面，值可以是任意类型，它甚至可以是另一个 Map。

测试中的其他内容应该都很熟悉了。

## 尝试运行测试

运行 `go test` 后编译器会提示失败信息 `./dictionary_test.go:8:9: undefined: Search`。

## 编写最少量的代码让测试运行并检查输出

在 `dictionary.go` 中：

```go
package main

func Search(dictionary map[string]string, word string) string {
    return ""
}
```

测试应该失败并显示*明确的错误信息*：

`dictionary_test.go:12: got '' want 'this is just a test' given, 'test'`。

## 编写足够的代码使测试通过

```go
func Search(dictionary map[string]string, word string) string {
    return dictionary[word]
}
```

从 Map 中获取值和数组相同，都是通过 `map[key]` 的方式。

## 重构

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
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

我决定创建一个 `assertStrings` 辅助函数并删除 `given` 的部分以使实现更通用。

## 使用自定义的类型

我们可以通过为 Map 创建新的类型并使用 `Search` 方法改进字典的使用。

在 `dictionary_test.go` 中：

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    got := dictionary.Search("test")
    want := "this is just a test"

    assertStrings(t, got, want)
}
```

我们已经开始使用 `Dictionary` 类型了，但是我们还没有定义它。然后要在 `Dictionary` 实例上调用 `Search` 方法。

我们不需要更改 `assertStrings`。

在 `dictionary.go` 中：

```go
type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
    return d[word]
}
```

在这里，我们创建了一个 `Dictionary` 类型，它是对 `map` 的简单封装。定义了自定义类型后，我们可以创建`Search` 方法。

## 首先编写测试

基本的搜索很容易实现，但是如果我们提供一个不在我们字典中的单词，会发生什么呢？

我们实际上得不到任何返回。这很好，因为程序可以继续运行，但还有更好的方法。这个函数可以证明该单词不在字典中。这样，用户不会好奇这个单词不存在还是未定义（这看起来可能对于字典没有用。但是，这可能是其他用例的关键场景）。

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    t.Run("known word", func(t *testing.T) {
        got, _ := dictionary.Search("test")
        want := "this is just a test"

        assertStrings(t, got, want)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dictionary.Search("test")
        want := "could not find the word you were looking for"

        if err == nil {
            t.Error("expected to get an error.")
        }

        assertStrings(t, got.Error(), want)
    })
}
```

在 Go 中处理这种情况的方法是返回第二个参数，它是一个 `Error` 类型。

`Error`s 可以使用 `.Error()` 方法转换为字符串，我们将其传递给断言时会执行此操作。我们也用 `if` 来保护 `assertStrings`，以确保我们不在 `nil` 上调用 `.Error()`。

## 尝试运行测试

这不会通过编译

```
./dictionary_test.go:18:10: assignment mismatch: 2 variables but 1 values
```

## 编写最少量的代码让测试运行并检查输出

```go
func (d Dictionary) Search(word string) (string, error) {
    return d[word], nil
}
```

现在你的测试将会失败，并显示更加清晰的错误信息

`dictionary_test.go:22: expected to get an error.`

## 编写足够的代码使测试通过

```go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", errors.New("could not find the word you were looking for")
    }

    return definition, nil
}
```

为了使测试通过，我们使用了一个 Map 查找的有趣属性。它可以返回两个值。第二个值是一个布尔值，表示是否成功找到 `key`。

此属性允许我们区分单词不存在还是未定义。

## 重构

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

我们通过将错误提取为变量的方式，摆脱 `Search` 中魔术错误（magic error）。这也会使我们获得更好的测试。

```go
t.Run("unknown word", func(t *testing.T) {
    _, got := dictionary.Search("unknown")

    assertError(t, got, ErrNotFound)
})

func assertError(t *testing.T, got, want error) {
    t.Helper()

    if got != want {
        t.Errorf("got error '%s' want '%s'", got, want)
    }
}
```

通过创建一个新的辅助函数，我们能够简化测试，并使用 `ErrNotFound` 变量，如果我们将来更改显示错误的文本，测试也不会失败。

## 首先编写测试

我们现在有很好的方法来搜索字典。但是，我们无法在字典中添加新单词。

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    dictionary.Add("test", "this is just a test")

    want := "this is just a test"
    got, err := dictionary.Search("test")
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if want != got {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

在这个测试中，我们利用 `Search` 方法使字典的验证更加容易。

## 编写最少量的代码让测试运行并检查输出

在 `dictionary.go` 中：

```go
func (d Dictionary) Add(word, definition string) {
}
```

测试现在应该会失败：

```
dictionary_test.go:31: should find added word: could not find the word you were
looking for
```

## 编写足够的代码使测试通过

```go
func (d Dictionary) Add(word, definition string) {
    d[word] = definition
}
```

Map 添加元素也类似于数组。你只需指定键并将其等于值。

### 引用类型

Maps 有一个有趣的属性，不使用指针传递你就可以修改它们。这是因为 `map` 是引用类型。这意味着它拥有对底层数据结构的引用，就像指针一样。它底层的数据结构是 `hash table` 或 `hash map`，你可以在[这里](https://en.wikipedia.org/wiki/Hash_table)阅读有关 `hash tables` 的更多信息。

Maps 作为引用类型是非常好的，因为无论 Map 有多大，都只会有一个副本。

引用类型引入的问题是 `maps` 可以是 `nil` 值。 如果你尝试使用 `nil` Map，你会得到一个 `nil 指针异常`，这将导致程序出错。

由于 `nil 指针异常`，你永远不应该初始化一个空的 Map 变量：

```go
var m map[string]string
```

相反，你可以像我们上面那样初始化空 `map`，或使用 `make` 关键字创建 `map`：

```go
dictionary = map[string]string{}

// OR

dictionary = make(map[string]string)
```

两种方法都在其上创建空的 `hash map` 和指针 `dictionary`。 以确保永远不会获得 `nil 指针异常`。

## 重构

在我们的实现中没有太多可以重构的地方，但测试可以简化一点。

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
        t.Errorf("got '%s' want '%s'", got, definition)
    }
}
```

我们为单词和定义创建了变量，并将定义断言移到了自己的辅助函数中。

我们的 `Add` 看起来不错。除此之外，我们没有考虑当我们尝试添加的值已经存在时会发生什么！

如果值已存在，Map 不会抛出错误。相反，它们将继续并使用新提供的值覆盖该值。这在实践中很方便，但会导致我们的函数名称不准确。`Add` 不应修改现有值。它应该只在我们的字典中添加新单词。

## 首先编写测试

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

        assertError(t, err, WordExistsError)
        assertDefinition(t, dictionary, word, definition)
    })
}
```

对于此测试，我们修改了 `Add` 以返回错误，我们将针对新的错误变量 `WordExistsError` 进行验证。我们还修改了之前的测试以检查是否为 `nil` 错误。

## 尝试运行测试

编译将失败，因为我们没有为 `Add` 返回值。

```
./dictionary_test.go:30:13: dictionary.Add(word, definition) used as value
./dictionary_test.go:41:13: dictionary.Add(word, "new test") used as value
```

## 编写最少量的代码让测试运行并检查输出

在 `dictionary.go` 中：

```go
var (
    ErrNotFound   = errors.New("could not find the word you were looking for")
    WordExistsError = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
    d[word] = definition
    return nil
}
```

现在我们又得到两个错误。我们仍在修改值，并返回 `nil` 错误。

```
dictionary_test.go:43: got error '%!s(<nil>)' want 'cannot add word because
it already exists'
dictionary_test.go:44: got 'new test' want 'this is just a test'
```

## 编写足够的代码使测试通过

```go
func (d Dictionary) Add(word, definition string) error {
    _, err := d.Search(word)
    switch err {
    case ErrNotFound:
        d[word] = definition
    case nil:
        return WordExistsError
    default:
        return err

    }

    return nil
}
```

这里我们使用 `switch` 语句来匹配错误。如上使用 `switch` 提供了额外的安全，以防 `Search` 返回错误而不是 `ErrNotFound`。

## 重构

我们没有太多需要重构的地方了，但随着对错误使用的增多，我们还可以做一些修改。

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

我们将错误声明为常量，这需要我们创建自己的 `DictionaryErr` 类型来实现 `error` 接口。你可以在 [Dave Cheney 的这篇优秀文章](https://dave.cheney.net/2016/04/07/constant-errors)中了解更多相关的细节。 简而言之，它使错误更具可重用性和不可变性。

## 首先编写测试

```go
func TestUpdate(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    dictionary.Update(dictionaryword, newDefinition)

    assertDefinition(t, dictionary, word, newDefinition)
}
```

`Update` 与 `Create` 密切相关，这是下一个需要我们实现的方法。

## 尝试运行测试

```
./dictionary_test.go:53:2: dictionary.Update undefined (type Dictionary has no field or method Update)
```

## 编写最少量的代码让测试运行并检查输出

我们已经知道如何处理这样的错误。我们需要定义我们的函数。

```go
func (d Dictionary) Update(word, definition string) {}
```

有了这个，就会看到我们需要改变这个词的定义。

```
 dictionary_test.go:55: got 'this is just a test' want 'new definition'
```

## 编写足够的代码使测试通过

当我们用 `create` 解决问题时就明白了如何处理这个问题。所以让我们实现一个与 `create` 非常相似的方法。

```go
func (d Dictionary) Update(word, definition string) {
    d[word] = definition
}
```

我们不需要对此进行重构，因为更改很简单。但是，我们现在遇到与 `create` 相同的问题。如果我们传入一个新单词，`Update` 会将它添加到字典中。

## 首先编写测试

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

    err := dictionary.Update(dictionaryword, definition)

    assertError(t, err, ErrWordDoesNotExist)
})
```

我们在单词不存在时添加了另一种错误类型。我们还修改了 `Update` 以返回 `error` 值。

## 尝试运行测试

```
./dictionary_test.go:53:16: dictionary.Update(word, "new test") used as value
./dictionary_test.go:64:16: dictionary.Update(word, definition) used as value
./dictionary_test.go:66:23: undefined: ErrWordDoesNotExists
```

这次我们得到 3 个错误，但我们知道如何处理这些错误。

## 编写最少量的代码让测试运行并检查输出

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

我们添加了自己的错误类型并返回 `nil` 错误。

通过这些更改，我们现在得到一个非常明确的错误：

```
dictionary_test.go:66: got error '%!s(<nil>)' want 'cannot update word because it does not exist'
```

## 编写足够的代码使测试通过

```go
func (d Dictionary) Update(word, definition string) error {
    _, err := dictionary.Search(word)
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

除了在更新 `dictionary` 和返回错误时切换之外，这个函数看起来几乎与Add完全相同。

### 关于声明 `Update` 的新错误的注意事项

我们可以重用 `ErrNotFound` 而不添加新错误。但是，更新失败时有更精确的错误通常更好。

具有特定错误可以为你提供有关错误的更多信息。以下是 Web 应用程序中的示例

> 遇到 `ErrNotFound` 时可以重定向用户，但遇到 `ErrWordDoesNotExist` 时会显示错误消息。

## 首先编写测试

```go
func TestDelete(t *testing.T) {
    word := "test"
    dictionary := Dictionary{word: "test definition"}

    dictionary.Delete(word)

    _, err := dictionary.Search(word)
    if err != ErrNotFound {
        t.Errorf("Expected '%s' to be deleted", word)
    }
}
```

我们的测试创建一个带有单词的 `Dictionary`，然后检查该单词是否已被删除。

## 尝试运行测试

通过运行 `go test` 我们得到：

```
./dictionary_test.go:74:6: dictionary.Delete undefined (type Dictionary has no field or method Delete)
```

## 编写最少量的代码让测试运行并检查输出

```go
func (d Dictionary) Delete(word string) {

}
```

添加这个之后，测试告诉我们没有删除这个单词。

```
dictionary_test.go:78: Expected 'test' to be deleted
```

## 编写足够的代码使测试通过

```go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

Go 的 maps 有一个内置函数 `delete`。它需要两个参数。第一个是 map，第二个是要删除的键。

`delete` 函数不返回任何内容，我们基于相同的概念构建 `Delete` 方法。由于删除一个不存在的值是没有影响的，与我们的 `Update` 和 `Create` 方法不同，我们不需要用错误复杂化 API。

## 总结

在本节中，我们介绍了很多内容。我们为 `dictionary` 创建了完整的 CRUD API。在整个过程中，我们学会了如何：

- 创建 maps
- 在 maps 中搜索值
- 向 maps 添加新值
- 更新 maps 中的值
- 从 maps 中删除值
- 了解更多有关**如何创建常量错误**的内容，编写对错误的封装
