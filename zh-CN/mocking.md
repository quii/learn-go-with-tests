# Mocking

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/mocking)**

现在需要你写一个程序，从 3 开始依次向下，当到 0 时打印 「GO!」 并退出，要求每次打印从新的一行开始且打印间隔一秒的停顿。

```
3
2
1
Go!
```

我们将通过编写一个 `Countdown` 函数来处理这个问题，然后放入 `main` 程序，所以它看起来这样：

```go
package main

func main() {
    Countdown()
}
```

虽然这是一个非常简单的程序，但要完全测试它，我们需要像往常一样采用迭代的、测试驱动的方法。

所谓迭代是指：确保我们采取最小的步骤让软件可用。

我们不想花太多时间写那些在被攻击后理论上还能运行的代码，因为这经常导致开发人员陷入开发的无底深渊。**尽你所能拆分需求是一项很重要的技能，这样你就能拥有可以工作的软件**。

下面是我们如何划分工作和迭代的方法：
- 打印 3
- 打印 3 到 Go!
- 在每行中间等待一秒

## 先写测试

我们的软件需要将结果打印到标准输出界面。在 DI（依赖注入） 的部分，我们已经看到如何使用 DI 进行方便的测试。

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}

    Countdown(buffer)

    got := buffer.String()
    want := "3"

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

如果你对 `buffer` 不熟悉，请重新阅读前面的部分。

我们清楚，我们的目的是让 `Countdown` 函数将数据写到某处，`io.writer` 就是作为 Go 的一个接口来抓取数据的一种方式。
- 在 `main` 中，我们将信息发送到 `os.Stdout`，所以用户可以看到 `Countdown` 的结果打印到终端
- 在测试中，我们将发送到 `bytes.Buffer`，所以我们的测试能够抓取到正在生成的数据

## 尝试并运行测试

`./countdown_test.go:11:2: undefined: Countdown`

## 为测试的运行编写最少量的代码，并检查失败测试的输出

定义 `Countdown` 函数

```go
func Countdown() {}
```

再次尝试运行

```
./countdown_test.go:11:11: too many arguments in call to Countdown
    have (*bytes.Buffer)
    want ()
```

编译器正在告诉你函数的问题，所以更正它

```go
func Countdown(out *bytes.Buffer) {}
```

`countdown_test.go:17: got '' want '3'`

这样结果就完美了！

## 编写足够的代码使程序通过

```go
func Countdown(out *bytes.Buffer) {
    fmt.Fprint(out, "3")
}
```
我们正在使用 `fmt.Fprint` 传入一个 `io.Writer`（例如 `*bytes.Buffer`）并发送一个 `string`。这个测试应该可以通过。

## 重构代码

虽然我们都知道 `*bytes.Buffer` 可以运行，但最好使用通用接口代替。

```go
func Countdown(out io.Writer) {
    fmt.Fprint(out, "3")
}
```

重新运行测试他们应该就可以通过了。

为了完成任务，现在让我们将函数应用到 `main`中。这样的话，我们就有了一些可工作的软件来确保我们的工作正在取得进展。

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func Countdown(out io.Writer) {
    fmt.Fprint(out, "3")
}

func main() {
    Countdown(os.Stdout)
}
```

尝试运行程序，这些成果会让你感到神奇。

当然，这仍然看起来很简单，但是我建议任何项目都使用这种方法。**在测试的支持下，将功能切分成小的功能点，并使其首尾相连顺利的运行。**

接下来我们可以让它打印 2，1 然后输出「Go!」。

## 先写测试

通过花费一些时间让整个流程正确执行，我们就可以安全且轻松的迭代我们的解决方案。我们将不再需要停止并重新运行程序，要对它的工作充满信心因为所有的逻辑都被测试过了。

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}

    Countdown(buffer)

    got := buffer.String()
    want := `3
2
1
Go!`

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}
```

反引号语法是创建 `string` 的另一种方式，但是允许你放置东西例如放到新的一行，对我们的测试来说是完美的。

## 尝试并运行测试

```
countdown_test.go:21: got '3' want '3
        2
        1
        Go!'
```

## 写足够的代码令测试通过

```go
func Countdown(out io.Writer) {
    for i := 3; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, "Go!")
}
```

使用 `for` 循环与 `i--` 反向计数，并且用 `fmt.println` 打印我们的数字到 `out`，后面跟着一个换行符。最后用 `fmt.Fprint` 发送 「Go!」。

## 重构代码

这里已经没有什么可以重构的了，只需要将变量重构为命名常量。

```go
const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, finalWord)
}
```

如果你现在运行程序，你应该可以获得想要的输出，但是向下计数的输出没有 1 秒的暂停。

Go 可以通过 `time.Sleep` 实现这个功能。尝试将其添加到我们的代码中。

```go
func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(out, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(out, finalWord)
}
```

如果你运行程序，它会以我们期望的方式工作。

## Mocking

测试可以通过，软件按预期的工作。但是我们有一些问题：

- 我们的测试花费了 4 秒的时间运行
  - 每一个关于软件开发的前沿思考性文章，都强调快速反馈循环的重要性。
  - **缓慢的测试会破坏开发人员的生产力。**
  - 想象一下，如果需求变得更复杂，将会有更多的测试。对于每一次新的 `Countdown` 测试，我们是否会对被添加到测试运行中 4 秒钟感到满意呢？
- 我们还没有测试这个函数的一个重要属性。

我们有个 `Sleep`ing 的注入，需要抽离出来然后我们才可以在测试中控制它。

如果我们能够 *mock* `time.Sleep`，我们可以用 *依赖注入* 的方式去来代替「真正的」`time.Sleep`，然后我们可以使用断言 **监视调用**。

## 先写测试

让我们将依赖关系定义为一个接口。这样我们就可以在 `main` 使用 *真实的* `Sleeper`，并且在我们的测试中使用 *spy sleeper*。通过使用接口，我们的 `Countdown` 函数忽略了这一点，并为调用者增加了一些灵活性。

```go
type Sleeper interface {
    Sleep()
}
```

我做了一个设计的决定，我们的 `Countdown` 函数将不会负责 `sleep` 的时间长度。
这至少简化了我们的代码，也就是说，我们函数的使用者可以根据喜好配置休眠的时长。

现在我们需要为我们使用的测试生成它的 *mock*。

```go
type SpySleeper struct {
    Calls int
}

func (s *SpySleeper) Sleep() {
    s.Calls++
}
```

*监视器（spies）*是一种 *mock*，它可以记录依赖关系是怎样被使用的。它们可以记录被传入来的参数，多少次等等。在我们的例子中，我们跟踪记录了 `Sleep()` 被调用了多少次，这样我们就可以在测试中检查它。

更新测试以注入对我们监视器的依赖，并断言 `sleep` 被调用了 4 次。

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}
    spySleeper := &SpySleeper{}

    Countdown(buffer, spySleeper)

    got := buffer.String()
    want := `3
2
1
Go!`

    if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }

    if spySleeper.Calls != 4 {
        t.Errorf("not enough calls to sleeper, want 4 got %d", spySleeper.Calls)
    }
}
```
## 尝试并运行测试

```
too many arguments in call to Countdown
    have (*bytes.Buffer, Sleeper)
    want (io.Writer)
```

## 为测试的运行编写最少量的代码，并检查失败测试的输出

我们需要更新 `Countdow` 来接受我们的 `Sleeper`。

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(out, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(out, finalWord)
}
```

如果您再次尝试，你的 `main` 将不会出现相同编译错误的原因

```
./main.go:26:11: not enough arguments in call to Countdown
    have (*os.File)
    want (io.Writer, Sleeper)
```
让我们创建一个 *真正的* sleeper 来实现我们需要的接口

```go
type ConfigurableSleeper struct {
    duration time.Duration
}

func (o *ConfigurableSleeper) Sleep() {
    time.Sleep(o.duration)
}
```

我决定做点额外的努力，让它成为我们真正的可配置的 sleeper。但你也可以在 1 秒内毫不费力地编写它。

我们可以在实际应用中使用它，就像这样：

```go
func main() {
    sleeper := &ConfigurableSleeper{1 * time.Second}
    Countdown(os.Stdout, sleeper)
}
```

## 足够的代码令测试通过

现在测试正在编译但是没有通过，因为我们仍然在调用 `time.Sleep` 而不是依赖注入。让我们解决这个问题。

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.sleep()
        fmt.Fprintln(out, i)
    }

    sleeper.sleep()
    fmt.Fprint(out, finalWord)
}
```

测试应该可以该通过，并且不再需要 4 秒。

### 仍然还有一些问题

还有一个重要的特性，我们还没有测试过。

`Countdown` 应该在第一个打印之前 sleep，然后是直到最后一个前的每一个，例如：
- `Sleep`
- `Print N`
- `Sleep`
- `Print N-1`
- `Sleep`
- `etc`

我们最新的修改只断言它已经 `sleep` 了 4 次，但是那些 `sleeps` 可能没按顺序发生。

当你在写测试的时候，如果你没有信心，你的测试将给你足够的信心，尽管推翻它！（不过首先要确定你已经将你的更改提交给了源代码控制）。将代码更改为以下内容。

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
    }

    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
```

如果你运行测试，它们仍然应该通过，即使实现是错误的。

让我们再用一种新的测试来检查操作的顺序是否正确。

我们有两个不同的依赖项，我们希望将它们的所有操作记录到一个列表中。所以我们会为它们俩创建 *同一个监视器*。

```go
type CountdownOperationsSpy struct {
    Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
    s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
    s.Calls = append(s.Calls, write)
    return
}

const write = "write"
const sleep = "sleep"
```

我们的 `CountdownOperationsSpy` 同时实现了 `io.writer` 和 `Sleeper`，把每一次调用记录到 `slice`。在这个测试中，我们只关心操作的顺序，所以只需要记录操作的代名词组成的列表就足够了。

现在我们可以在测试套件中添加一个子测试。

```go
t.Run("sleep after every print", func(t *testing.T) {
    spySleepPrinter := &CountdownOperationsSpy{}
    Countdown(spySleepPrinter, spySleepPrinter)

    want := []string{
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
    }

    if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
        t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
    }
})
```

现在这个测试应该会失败。恢复原状新测试应该又可以通过。

我们现在在 `Sleeper` 上有两个测试监视器，所以我们现在可以重构我们的测试，一个测试被打印的内容，另一个是确保我们在打印时间 `sleep`。最后我们可以删除第一个监视器，因为它已经不需要了。

```go
func TestCountdown(t *testing.T) {

    t.Run("prints 3 to Go!", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        Countdown(buffer, &CountdownOperationsSpy{})

        got := buffer.String()
        want := `3
2
1
Go!`

        if got != want {
            t.Errorf("got '%s' want '%s'", got, want)
        }
    })

    t.Run("sleep after every print", func(t *testing.T) {
        spySleepPrinter := &CountdownOperationsSpy{}
        Countdown(spySleepPrinter, spySleepPrinter)

        want := []string{
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
        }

        if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
            t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
        }
    })
}
```

我们现在有了自己的函数，并且它的两个重要的属性已经通过合理的测试。

## 难道 mocking 不是在作恶（evil）吗？

你可能听过 mocking 是在作恶。就像软件开发中的任何东西一样，它可以被用来作恶，就像 DRY(Don't repeat yourself) 一样。

当人们 *不听从他们的测试* 并且 *不尊重重构阶段时*，他们通常会陷入糟糕的境地。

如果你的模拟代码变得很复杂，或者你需要模拟很多东西来测试一些东西，那么你应该 *倾听* 那种糟糕的感觉，并考虑你的代码。通常这是一个征兆：

- 你正在进行的测试需要做太多的事情
  - 把模块分开就会减少测试内容
- 它的依赖关系太细致
  - 考虑如何将这些依赖项合并到一个有意义的模块中
- 你的测试过于关注实现细节
  - 最好测试预期的行为，而不是功能的实现

通常，在你的代码中有大量的 mocking 指向 *错误的抽象*。

**人们在这里看到的是测试驱动开发的弱点，但它实际上是一种力量**，通常情况下，糟糕的测试代码是糟糕设计的结果，而设计良好的代码很容易测试。

### 但是模拟和测试仍然让我举步维艰！

曾经遇到过这种情况吗？

- 你想做一些重构
- 为了做到这一点，你最终会改变很多测试
- 你对测试驱动开发提出质疑，并在媒体上发表一篇文章，标题为「Mocking 是有害的」

这通常是您测试太多 *实现细节* 的标志。尽力克服这个问题，所以你的测试将测试 *有用的行为*，除非这个实现对于系统运行非常重要。

有时候很难知道到底要测试到 *什么级别*，但是这里有一些我试图遵循的思维过程和规则。

- *重构的定义是代码更改，但行为保持不变。* 如果您已经决定在理论上进行一些重构，那么你应该能够在没有任何测试更改的情况下进行提交。所以，在写测试的时候问问自己。
  - 我是在测试我想要的行为还是实现细节？
  - 如果我要重构这段代码，我需要对测试做很多修改吗？
- 虽然 Go 允许你测试私有函数，但我将避免它作为私有函数与实现有关。
- 我觉得如果一个测试 **超过 3 个模拟，那么它就是警告** —— 是时候重新考虑设计。
- 小心使用监视器。监视器让你看到你正在编写的算法的内部细节，这是非常有用的，但是这意味着你的测试代码和实现之间的耦合更紧密。**如果你要监视这些细节，请确保你真的在乎这些细节。**

和往常一样，软件开发中的规则并不是真正的规则，也有例外。Uncle Bob 的文章 [「When to mock」](https://8thlight.com/blog/uncle-bob/2014/05/10/WhenToMock.html) 有一些很好的指南。

## 总结

### 更多关于测试驱动开发的方法

- 当面对不太简单的例子，把问题分解成「简单的模块」。试着让你的工作软件尽快得到测试的支持，以避免掉进兔子洞（rabbit holes，意指未知的领域）和采取「最终测试（Big bang）」的方法。
- 一旦你有一些正在工作的软件，*小步迭代* 应该是很容易的，直到你实现你所需要的软件。

### Mocking

- 没有对代码中重要的区域进行 mock 将会导致难以测试。在我们的例子中，我们不能测试我们的代码在每个打印之间暂停，但是还有无数其他的例子。调用一个 *可能* 失败的服务？想要在一个特定的状态测试您的系统？在不使用 mocking 的情况下测试这些场景是非常困难的。
- 如果没有 mock，你可能需要设置数据库和其他第三方的东西来测试简单的业务规则。你可能会进行缓慢的测试，从而导致 **缓慢的反馈循环**。
- 当不得不启用一个数据库或者 webservice 去测试某个功能时，由于这种服务的不可靠性，你将会得到的是一个 **脆弱的测试**。

一旦开发人员学会了 mocking，就很容易对系统的每一个方面进行过度测试，按照 *它工作的方式* 而不是 *它做了什么*。始终要注意 **测试的价值**，以及它们在将来的重构中会产生什么样的影响。

在这篇关于 mocking 的文章中，我们只提到了 **监视器（Spies）**，他们是一种 mock。也有不同类型的 mocks[。Uncle Bob 的一篇极易阅读的文章中解释了这些类型](https://8thlight.com/blog/uncle-bob/2014/05/14/TheLittleMocker.html)。在后面的章节中，我们将需要编写依赖于其他数据的代码，届时我们将展示 **Stubs** 行为。

---

作者：[Chris James](https://dev.to/quii)
译者：[Donng](https://github.com/Donng)
校对：[rxcai](https://github.com/rxcai)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
