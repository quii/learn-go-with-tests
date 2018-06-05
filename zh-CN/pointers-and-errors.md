# 指针和错误

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/pointers)**

我们在上一节中学习了结构体（structs），它可以组合与一个概念相关的一系列值。

你有时可能想用结构体来管理状态，通过将方法暴露给用户的方式，让他们在你可控的范围内修改状态。

**金融科技行业都喜欢 Go** 和比特币吧?那就来看看我们能创造出多么惊人的银行系统。

首先声明一个 `Wallet（钱包）` 结构体，我们可以利用它存放 `Bitcoin（比特币）`。

## 首先写测试

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()
	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

在前面的示例中，我们直接使用字段名称访问字段，但是在 *非常强调安全性的钱包* 中，我们不想暴露自己的内部状态，而是通过方法来控制访问的权限。

## 尝试运行测试

`./wallet_test.go:7:12: undefined: Wallet`

## 为测试的运行编写最少量的代码并检查失败测试的输出

编译器不知道 `Wallet` 是什么，所以让我们告诉它。

```go
type Wallet struct { }
```

现在我们已经生成了自己的钱包，尝试再次运行测试

```
./wallet_test.go:9:8: wallet.Deposit undefined (type Wallet has no field or method Deposit)
./wallet_test.go:11:15: wallet.Balance undefined (type Wallet has no field or method Balance)
```

正如所料，我们需要定义这些方法以使测试通过。

请记住，只做足够让测试运行的事情。 我们需要确保测试失败时，显示清晰的错误信息。

```go
func (w Wallet) Deposit(amount int) {

}

func (w Wallet) Balance() int {
	return 0
}
```

如果你对此语法不熟悉，请重新阅读结构体章节。

测试现在应该编译通过了，然后运行

`wallet_test.go:15: got 0 want 10`

## 编写足够的代码使其通过

结构中需要一些 `balance（余额）` 变量来存储状态

```go
type Wallet struct {
	balance int
}
```
在 Go 中，如果一个符号(例如变量、类型、函数等)是以小写符号开头，那么它在 *定义它的包之外* 就是私有的。

在我们的例子中，我们只想让自己的方法修改这个值，而其他的不可以。

记住，我们可以使用「receiver」变量访问结构体内部的 `balance` 字段。

```go
func (w Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w Wallet) Balance() int {
	return w.balance
}
```

现在我们的事业在金融科技的保护下，将会运行并轻松的通过测试

`wallet_test.go:15: got 0 want 10`

### 为什么报错了？

这让人很困惑，我们的代码看上去没问题，我们在余额中添加了新的金额，然后余额的方法应该返回它当前的状态值。

在 Go 中，**当调用一个函数或方法时，参数会被复制**。

当调用 `func (w Wallet) Deposit(amount int)` 时，`w` 是来自我们调用方法的副本。

不需要太过计算机化，当你创建一个值，例如一个 `wallet`，它就会被存储在内存的某处。你可以用 `&myval` 找到那块内存的*地址*。

通过在代码中添加一些 `prints` 来试验一下

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()

	fmt.Println("address of balance in test is", &wallet.balance)

	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

```go
func (w Wallet) Deposit(amount int) {
	fmt.Println("address of balance in Deposit is", &w.balance)
	w.balance += amount
}
```

现在重新运行测试

```text
address of balance in Deposit is 0xc420012268
address of balance in test is 0xc420012260
```

可以看出两个 `balance` 的地址是不同的。因此，当我们在代码中更改 `balance` 的值时，我们处理的是来自测试的副本。因此，`balance` 在测试中没有被改变。

我们可以用 *指针* 来解决这个问题。指针让我们 *指向* 某个值，然后修改它。所以，我们不是拿钱包的副本，而是拿一个指向钱包的指针，这样我们就可以改变它。

```go
func (w *Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w *Wallet) Balance() int {
	return w.balance
}
```

不同之处在于，接收者类型是 `*Wallet` 而不是 `Wallet`，你可以将其解读为「指向 `wallet` 的指针」。

尝试重新运行测试，它们应该可以通过了。

## 重构

我们曾说过我们正在制做一个比特币钱包，但到目前为止我们还没有提到它们。我们一直在使用 `int`，因为当用来计数时它是不错的类型!

为此创建一个结构体似乎有点过头了。就 `int` 的表现来说已经很好了，但问题是它不具有描述性。

Go 允许从现有的类型创建新的类型。

语法是 `type MyName OriginalType`

```go
type Bitcoin int

type Wallet struct {
	balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}
```

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(Bitcoin(10))

	got := wallet.Balance()

	want := Bitcoin(10)

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

要生成 `Bitcoin（比特币）`，你只需要用 `Bitcoin(999)` 的语法就可以了。

类型别名有一个有趣的特性，你还可以对它们声明 *方法*。当你希望在现有类型之上添加一些领域内特定的功能时，这将非常有用。

[让我们实现 Bitcoin 的 Stringer 方法](https://golang.org/pkg/fmt/#Stringer)

```go
type Stringer interface {
	String() string
}
```

这个接口是在 `fmt` 包中定义的。当使用 `%s` 打印格式化的字符串时，你可以定义此类型的打印方式。

```go
func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

如你所见，在类型别名上创建方法的语法与结构上的语法相同。

接下来，我们需要更新测试中的格式化字符串，以便它们将使用 `String()` 方法。

```go
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
```

为了看到这一点，故意令测试失败我们就能看到

`wallet_test.go:18: got 10 BTC want 20 BTC`

这使得我们的测试更加清晰。

下一个需求是 `Withdraw（提取）` 函数。

## 先写测试

几乎跟 `Deposit()` 相反

```go
func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}

		wallet.Deposit(Bitcoin(10))

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}

		wallet.Withdraw(10)

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

}
```

## 尝试运行测试

`./wallet_test.go:26:9: wallet.Withdraw undefined (type Wallet has no field or method Withdraw)`

## 为测试的运行编写最少量的代码并检查失败测试的输出

```go
func (w *Wallet) Withdraw(amount Bitcoin) {

}
```

`wallet_test.go:33: got 20 BTC want 10 BTC`

## 编写足够的代码使其通过

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
	w.balance -= amount
}
```

## 重构

在我们的测试中有一些重复部分，我们来重构一下。

```go
func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

}
```

如果你试图从账户中取出更多的钱，会发生什么？目前，我们的要求是假定没有透支设备。

我们如何在使用 `Withdraw` 时标记出现的问题呢？

在 Go 中，如果你想指出一个错误，通常你的函数要返回一个 `err`，以便调用者检查并执行相应操作。

让我们在测试中试试。

## 先写测试

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
	startingBalance := Bitcoin(20)
	wallet := Wallet{startingBalance}
	err := wallet.Withdraw(Bitcoin(100))

	assertBalance(t, wallet, startingBalance)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
})
```

*如果* 你尝试取出超过你余额的比特币，我们想让 `Withdraw` 返回一个错误，而余额应该保持不变。

然后，如果测试失败，我们检查错误是否为 `nil`。

`nil` 是其他编程语言的 `null`。错误可以是 `nil`，因为返回类型是 `error`，这是一个接口。如果你看到一个函数，它接受参数或返回值的类型是接口，它们就可以是 `nil`。

如果你尝试访问一个值为 `nil` 的值，它将会引发 **运行时的 panic**。这很糟糕!你应该确保你检查了 `nil` 的值。

## 尝试运行测试

`./wallet_test.go:31:25: wallet.Withdraw(Bitcoin(100)) used as value`

错误信息可能不太清楚，但我们之前对于 `Withdraw` 的意图只是调用它，它永远不会返回一个值。为了使它编译通过，我们需要更改它，以便它有一个返回类型。

## 为测试的运行编写最少量的代码并检查失败测试的输出

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
	w.balance -= amount
	return nil
}
```

再次强调，编写足够的代码来满足编译器的要求是非常重要的。我们纠正了自己的 `Withdraw` 方法返回 `error`，现在我们必须返回 *一些东西*，所以我们就返回 `nil` 好了。

## 编写足够的代码使其通过

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("oh no")
	}

	w.balance -= amount
	return nil
}
```

记住要将 `errors` 导入到代码中。

`errors.New` 创建了一个新的 `error`，并带有你选择的消息。

## 重构

让我们为错误检查做一个快速测试的助手方法，以帮助我们的测试读起来更清晰。

```go
assertError := func(t *testing.T, err error) {
	if err == nil {
		t.Error("wanted an error but didnt get one")
	}
}
```

并且在我们的测试中

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
	wallet := Wallet{Bitcoin(20)}
	err := wallet.Withdraw(Bitcoin(100))

	assertBalance(t, wallet, Bitcoin(20))
	assertError(t, err)
})
```

希望在返回「oh no」的错误时，你会认为我们可能会迭代这个问题，因为它作为返回值看起来没什么用。

假设错误最终会返回给用户，让我们更新测试以断言某种错误消息，而不只是让错误存在。

## 先写测试

更新一个 `string` 的助手方法来比较。

```go
assertError := func(t *testing.T, got error, want string) {
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got.Error() != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
```

同时再更新调用者

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
	startingBalance := Bitcoin(20)
	wallet := Wallet{startingBalance}
	err := wallet.Withdraw(Bitcoin(100))

	assertBalance(t, wallet, startingBalance)
	assertError(t, err, "cannot withdraw, insufficient funds")
})
```

我们已经介绍了 `t.Fatal`。如果它被调用，它将停止测试。这是因为我们不希望对返回的错误进行更多断言。如果没有这个，测试将继续进行下一步，并且因为一个空指针而引起 panic。

## 尝试运行测试

`wallet_test.go:61: got err 'oh no' want 'cannot withdraw, insufficient funds'`

## 编写足够的代码使其通过

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("cannot withdraw, insufficient funds")
	}

	w.balance -= amount
	return nil
}
```

## 重构

我们在测试代码和 `Withdraw` 代码中都有重复的错误消息。

如果有人想要重新定义这个错误，那么测试就会失败，这将是非常恼人的，而对于我们的测试来说，这里有太多的细节了。我们并不关心具体的措辞是什么，只是在给定条件的情况下返回一些有意义的错误。

在Go中，错误是值，因此我们可以将其重构为一个变量，并为其提供一个单一的事实来源。

```go
var InsufficientFundsError = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return InsufficientFundsError
	}

	w.balance -= amount
	return nil
}
```

`var` 关键字允许我们定义包的全局值。

这是一个积极的变化，因为现在我们的 `Withdraw` 函数看起来很清晰。

接下来，我们可以重构我们的测试代码来使用这个值而不是特定的字符串。

```go
func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertBalance(t, wallet, Bitcoin(20))
		assertError(t, err, InsufficientFundsError)
	})
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	got := wallet.Balance()

	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertError(t *testing.T, got error, want error) {
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
```

现在这个测试也更容易理解了。

我已经将助手函数从主测试函数中移出，这样当某人打开一个文件时，他们就可以开始读取我们的断言，而不是一些助手函数。

测试的另一个有用的特性是，它帮助我们理解代码的真实用途，从而使我们的代码更具交互性。我们可以看到，开发人员可以简单地调用我们的代码，并对 `InsufficientFundsError` 进行相等的检查，并采取相应的操作。

### 未经检查的错误

虽然 Go 编译器对你有很大帮助，但有时你仍然会忽略一些事情，错误处理有时会很棘手。

有一种情况我们还没有测试过。要找到它，在一个终端中运行以下命令来安装 `errcheck`，这是许多可用的 linters（代码检测工具） 之一。

`go get -u github.com/kisielk/errcheck`

然后，在你的代码目录中运行 `errcheck .`。

你应该会得到如下类似的内容：

`wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))`

这告诉我们的是，我们没有检查在代码行中返回的错误。我的计算机上的这行代码与我们的正常 `withdraw` 的场景相对应，因为我们没有检查 `Withdraw` 是否成功，因此没有返回错误。

这是最终的测试代码。

```go
func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
		assertNoError(t, err)
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertBalance(t, wallet, Bitcoin(20))
		assertError(t, err, InsufficientFundsError)
	})
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	if got != nil {
		t.Fatal("got an error but didnt want one")
	}
}

func assertError(t *testing.T, got error, want error) {
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
```

## 总结

### 指针

* 当你传值给函数或方法时，Go 会复制这些值。因此，如果你写的函数需要更改状态，你就需要用指针指向你想要更改的值
* Go 取值的副本在大多数时候是有效的，但是有时候你不希望你的系统只使用副本，在这种情况下你需要传递一个引用。例如，非常庞大的数据或者你只想有一个实例（比如数据库连接池）

### nil

* 指针可以是 nil
* 当函数返回一个的指针，你需要确保检查过它是否为 nil，否则你可能会抛出一个执行异常，编译器在这里不能帮到你
* nil 非常适合描述一个可能丢失的值

### 错误

* 错误是在调用函数或方法时表示失败的
* 通过测试我们得出结论，在错误中检查字符串会导致测试不稳定。因此，我们用一个有意义的值重构了，这样就更容易测试代码，同时对于我们 API 的用户来说也更简单。
* 错误处理的故事远远还没有结束，你可以做更复杂的事情，这里只是抛砖引玉。后面的部分将介绍更多的策略。
* [不要只是检查错误，要优雅地处理它们](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### 从现有的类型中创建新的类型。

* 用于为值添加更多的领域内特定的含义
* 可以让你实现接口

指针和错误是 Go 开发中重要的组成部分，你需要适应这些。幸运的是，如果你做错了，编译器通常会帮你解决问题，你只需要花点时间读一下错误信息。

----------------

作者：[Chris James](https://dev.to/quii)
译者：[Donng](https://github.com/Donng)
校对：[polaris1119](https://github.com/polaris1119)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
