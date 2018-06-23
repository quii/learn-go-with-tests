# 结构体，方法和接口

**[你可以在这里找到本章的所有代码](https://github.com/quii/learn-go-with-tests/tree/master/structs)**

假设我们需要编程计算一个给定高和宽的长方形的周长。我们可以写一个函数如下：

`Premeter(width float64, height float64)`

其中 `float64` 是形如 `123.45` 的浮点数。

现在我们应该很熟悉 TDD 的方式了。

## 先写测试函数

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

注意到新的格式化字符串了吗？这里的 `f` 对应 `float64`，`.2` 表示输出 2 位小数。

## 运行测试

`./shapes_test.go:6:9: undefined: Perimeter`

## 为运行测试函数编写最少的代码并检查失败时的输出

```go
func Perimeter(width float64, height float64) float64 {
    return 0
}
```

运行结果是：`shapes_test.go:10: got 0 want 40`

## 编写正确的代码让测试函数通过

```go
func Perimeter(width float64, height float64) float64 {
    return 2*(width + height)
}
```

到目前为止还很简单。现在让我们来编写一个函数 Area(width, height float64) 来返回长方形的面积。

你可以先自己按照 TDD 的方式尝试一下。

相应的测试函数如下所示：

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    got := Area(12.0, 6.0)
    want := 72.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

相应的代码如下：

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}

func Area(width float64, height float64) float64 {
    return width * height
}
```

## 重构

我们的代码能正常工作，但是其中不包含任何显式的信息表示计算的是长方形。粗心的开发者可能会错误的调用这些函数来计算三角形的周长和面积而没有意识到错误的结果。

我们可以仅仅给这些函数命名成像 RectangleArea 一样更具体的名字。但是更简洁的方案是定义我们自己的类型 Rectangle，它可以封装长方形的信息。

我们可以使用保留字 struct 来定义自己的类型。一个通过 struct 定义出来的类型是一些已命名的域的集合，这些域用来保存数据。

一个 struct 的声明如下：

```go
type Rectangle struct {
    Width float64
    Height float64
}
```

现在让我们用类型 Rectangle 代替简单的 float64 来重构这些测试函数。

```go
func TestPerimeter(t *testing.T) {
    rectangle := Rectangle{10.0, 10.0}
    got := Perimeter(rectangle)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    rectangle := Rectangle{12.0, 6.0}
    got := Area(rectangle)
    want := 72.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

记住先运行这些测试函数再尝试修复问题，因为运行后我们能获得有用的错误信息：

```
./shapes_test.go:7:18: not enough arguments in call to Perimeter
    have (Rectangle)
    want (float64, float64)
```

我们可以通过下面的语法来访问一个 struct 中的域： `myStruct.field`

代码需要调整如下：

```go
func Perimeter(rectangle Rectangle) float64 {
    return 2 * (rectangle.Width + rectangle.Height)
}

func Area(rectangle Rectangle) float64 {
    return rectangle.Width * rectangle.Height
}
```

我希望你同意通过传递一个类型为 Rectangle 的参数给这些函数更能表达我们的用意。并且这样做我们将会看到有更多的好处。

我们的下一个需求是为圆形写一个类似的函数。

## 先写测试函数

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := Area(rectangle)
        want := 72.0

        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := Area(circle)
        want := 314.16

        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    })

}
```

## 运行测试

`./shapes_test.go:28:13: undefined: Circle`

## 为运行测试函数编写最少的代码并检查失败时的输出

我们需要定义一个 Circle 类型：

```go
type Circle struct {
    Radius float64
}
```

现在我们重新运行测试：

`./shapes_test.go:29:14: cannot use circle (type Circle) as type Rectangle in argument to Area`

有些编程语言中我们可以这样做：

```go
func Area(circle Circle) float64 { ... }
func Area(rectangle Rectangle) float64 { ... }
```

但是在 Go 语言中你不能这么做：

`./shapes.go:20:32: Area redeclared in this block`

我们有以下两个选择：

* 不同的包可以有函数名相同的函数。所以我们可以在一个新的包里创建函数 Area(Circle)。但是感觉有点大才小用了
* 我们可以为新类型定义方法。

## 什么是方法？

到目前为止我们只编写过函数但是我们已经使用过方法。当我们调用 t.Errorf 时我们调用了 t(testing.T) 这个实例的方法 ErrorF。

方法和函数很相似但是方法是通过一个特定类型的实例调用的。函数可以随时被调用，比如 Area(rectangle)。不像方法需要在某个事物上调用。

示例会帮助我们理解。让我们通过方法调用的方式来改写测试函数并尝试修复代码。

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := rectangle.Area()
        want := 72.0

        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := circle.Area()
        want := 314.1592653589793

        if got != want {
            t.Errorf("got %f want %f", got, want)
        }
    })

}
```

尝试运行测试函数，我们会得到如下结果：

```text
./shapes_test.go:19:19: rectangle.Area undefined (type Rectangle has no field or method Area)
./shapes_test.go:29:16: circle.Area undefined (type Circle has no field or method Area)
```

> type Circle has no field or method Area

大家可以看到编译器的伟大之处。花些时间慢慢阅读这个错误信息是很重要的，这种习惯将对你长期有用。

## 为运行测试函数编写最少的代码并检查失败时的输出

我们给这些类型加一些方法：

```go
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return 0
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 0
}
```

声明方法的语法跟函数差不多，因为他们本身就很相似。唯一的不同是方法接收者的语法 `func(receiverName ReceiverType) MethodName(args)`。

当方法被这种类型的变量调用时，数据的引用通过变量 `receiverName` 获得。在其他许多编程语言中这些被隐藏起来并且通过 `this` 来获得接收者。

把类型的第一个字母作为接收者变量是 Go 语言的一个惯例。

```go
r Rectangle
```

现在尝试重新运行测试，编译通过了但是会有一些错误输出。

## 编写足够的代码让测试函数通过

现在让我们修改我们的新方法以让矩形测试通过：

```go
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}
```

现在重跑测试，矩形测试应该通过了但是圆的测试还是失败的。

为使 Circle 测试通过我们需要从 math 包中借用常数 Pi（记得引入 math 包）。

```go
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}
```

## 重构

我们的测试有些重复。

我们想做的是给定一些几何形状，调用 `Area()` 方法并检查结果。

我们想写一个这样的函数 `CheckArea`，其参数是任何类型的几何形状。如果参数不是几何形状的类型，那么编译应该报错。
Go 语言中我们可以通过接口实现这一目的。

[接口](https://golang.org/ref/spec#Interface_types)在 Go 这种静态类型语言中是一种非常强有力的概念。因为接口可以让函数接受不同类型的参数并能创造类型安全且高解耦的代码。

让我们引入接口来重构我们的测试代码：

```go
func TestArea(t *testing.T) {

    checkArea := func(t *testing.T, shape Shape, want float64) {
        t.Helper()
        got := shape.Area()
        if got != want {
            t.Errorf("got %.2f want %.2f", got, want)
        }
    }

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        checkArea(t, rectangle, 72.0)
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        checkArea(t, circle, 314.1592653589793)
    })

}
```

像其他练习一样我们创建了一个辅助函数，但不同的是我们传入了一个 Shape 类型。如何没有定义 Shape 类型编译会报错。

怎样定义 Shape 类型呢？我们用一个 Go 语言的接口定义来声明 Shape 类型：

```go
type Shape interface {
    Area() float64
}
```

这样我们就像创建 `Rectangle` 和 `Circle` 一样创建了一个新类型，不过这次是 interface 而不是 struct。

加了这个代码后测试运行通过了。

### 稍等，什么？

这种定义 `interface` 的方式与大部分其他编程语言不同。通常接口定义需要这样的代码 `My type Foo implements interface Bar`。

但是在我们的例子里，

* `Rectangle` 有一个返回值类型为 `float64` 的方法 `Area`，所以它满足接口 `Shape`
* `Circle` 有一个返回值类型为 `float64` 的方法 `Area`，所以它满足接口 `Shape`
* `string` 没有这种方法，所以它不满足这个接口
* 等等

在 Go 语言中 **interface resolution 是隐式的**。如果传入的类型匹配接口需要的，则编译正确。

### 解耦

请注意我们的辅助函数是怎样实现不需要关心参数是矩形，圆形还是三角形的。通过声明一个接口，辅助函数能从具体类型解耦而只关心方法本身需要做的工作。

这种方法使用接口来声明我们仅仅需要的。这种方法在软件设计中非常重要，我们以后在后续部分中还是涉及到更多细节。

## 进一步重构

现在我们对结构体有一定的理解了，我们可以引入「表格驱动测试」。

[表格驱动测试](https://github.com/golang/go/wiki/TableDrivenTests)在我们要创建一系列相同测试方式的测试用例时很有用。

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.want {
            t.Errorf("got %.2f want %.2f", got, tt.want)
        }
    }

}
```

这里唯一的新语法是创建了一个匿名的结构体。我们用含有两个域 shape 和 want 的 []struct 声明了一个结构体切片。然后我们用测试用例去填充这个数组了。

我们可以像遍历任何其他切片一样来遍历这个数组，进而用这个结构体的域来做我们的测试。

你会看到开发人员能方便的引入一个新的几何形状，只需实现 Area 方法并把新的类型加到测试用例中。另外发现 Area 方法有错误，我们可以在修复这个错误之前非常容易的添加新的测试用例。

列表驱动测试可以成为你工具箱中的得力武器。但是确保你在测试中真的需要使用它。如果你要测试一个接口的不同实现，或者传入函数的数据有很多不同的测试需求，这个武器将非常给力。

让我们通过再添加一个三角形并测试它来演示所有这些技术。

## 先写测试函数

为我们的新类型添加测试用例非常容易，只需添加 "{Triangle{12,6},36.0}," 到我们的列表中去就行了。

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
        {Triangle{12, 6}, 36.0},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.want {
            t.Errorf("got %.2f want %.2f", got, tt.want)
        }
    }

}
```

## 尝试运行测试函数

记住，不断尝试运行这些测试函数并让编译器引导你找到正确的方案。

## 为运行测试函数编写最少的代码并检查失败时的输出

`./shapes_test.go:25:4: undefined: Triangle`

我们还没有定义 Triangle 类型：

```go
type Triangle struct {
    Base   float64
    Height float64
}
```

在运行一次测试函数：

```text
./shapes_test.go:25:8: cannot use Triangle literal (type Triangle) as type Shape in field value:
    Triangle does not implement Shape (missing Area method)
```

编译器告诉我们不能把 Triangle 当作一个类型因为它没有方法 Area()。所以我们添加一个空的实现让测试函数能工作：

```go
func (c Triangle) Area() float64 {
    return 0
}
```

最后代码编译通过。运行后得到如下错误：

`shapes_test.go:31: got 0.00 want 36.00`

## 编写正确的代码让测试函数通过

```go
func (c Triangle) Area() float64 {
    return (c.Base * c.Height) * 0.5
}
```

最后测试通过了！

## 重构

虽然实现很好但我们的测试函数还能够改进。

注意如下代码：

```go
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

这些数字代表什么并不一目了然，我们应该让我们的测试函数更容易理解。

到目前为止我们仅仅学到一种创建结构体 MyStruct{val1, val2} 的方法，但是我们可以选择命名这些域。就像如下代码所示：

```go
    {shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
    {shape: Circle{Radius: 10}, want: 314.1592653589793},
    {shape: Triangle{Base: 12, Height: 6}, want: 36.0},
```

在 Kent Beck 的这边题为 [测试驱动开发实例](https://g.co/kgs/yCzDLF) 的帖子中把测试用例重构成要点和断言：

>当测试用例不是一系列操作，而是事实的断言时，测试才清晰明了。

（认证了我的观点）

现在我们的测试用例是关于几何图形的面积这些事实的断言了。

## 确保测试输出有效

记得当时我们实现三角形时的错误输出吗？它输出 `shapes_test.go:31: got 0.00 want 36.00`

我们知道它仅仅和三角形有关，但是如果在一个包含二十个测试用例的系统里出现类似的错误呢？开发人员如何指导是哪个测试用例失败了呢？这对于开发人员来说不是一个好的体验，他们需要手工检查所有的测试用例以定位到哪个用例失败了。

我们可以改进我们的错误输出为 "%#v got %.2f want %.2f. %#v"，这样会打印结构体中域的值。这样开发人员能一眼看出被测试的属性。

关于列表驱动测试的最后一点提示是使用 t.Run。

在每个用例中使用 t.Run，测试用例的错误输出中会包含用例的名字：

```text
-------- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
```

我们可以通过如下命令来运行列表中指定的测试用例： `go test -run TestArea/Rectangle`

下面是满足要求的最终测试代码：

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        name    string
        shape   Shape
        hasArea float64
    }{
        {name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
        {name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
        {name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
    }

    for _, tt := range areaTests {
        // using tt.name from the case to use it as the `t.Run` test name
        t.Run(tt.name, func(t *testing.T) {
            got := tt.shape.Area()
            if got != tt.hasArea {
                t.Errorf("%#v got %.2f want %.2f", tt.shape, got, tt.hasArea)
            }
        })

    }

}
```

## 总结

这是进一步的 TDD 实践。我们在对一个基本数学问题的解决方案的迭代中，通过测试学习了语言的新特性。

* 声明结构体以创建我们自己的类型，让我们把数据集合在一起并达到简化代码的目地
* 声明接口，这样我们可以定义适合不同参数类型的函数（参数多态）
* 在自己的数据类型中添加方法以实现接口
* 列表驱动测试让断言更清晰，这样可以使测试文件更易于扩展和维护

这是重要的一课。因为我们开始定义自己的类型。在像 Go 这样的静态语言中，能定义自己的类型是开发易维护，低耦合，好测试的软件的基础。

接口是把负责从系统的其他部分隐藏起来的伟大工具。在我们的测试中，辅助函数的代码不需要知道具体的几何形状，只需要知道获取它的面积即可。

当你以后更熟悉 Go 后你会发现接口和标准库的真正威力。你会看到标准库中的随处可见的接口。通过在你自己的类型中实现这些接口你能很快的重用大量的伟大功能。

---

作者：[Chris James](https://dev.to/quii)
译者：[hzpfly](https://github.com/hzpfly)
校对：[polaris1119](https://github.com/polaris1119)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
