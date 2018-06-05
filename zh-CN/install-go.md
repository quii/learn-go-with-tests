# 安装 Go，搭建开发环境

Go 语言官方安装指引可参见[这里](https://golang.org/doc/install)。

这个指引假设你正在使用其中一种包管理工具，如 [Homebrew](https://brew.sh)，[Chocolatey](https://chocolatey.org)，[Apt](https://help.ubuntu.com/community/AptGet/Howto) 或 [yum](https://access.redhat.com/solutions/9934)。

出于示范的目的，我们将展示在 OSX 系统中使用 Homebrew 安装的步骤。

## 安装

安装过程非常容易。首先，你需要先运行下面这个命令来安装 homebrew(brew)。Brew 依赖 Xcode，所以先确保 Xcode 已安装。

```sh
xcode-select --install
```

然后运行下面的命令来安装 homebrew：

```sh
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

现在可以安装 Go 了：

```sh
brew install go
```

如果你想把你的程序部署到 Linux 服务器，你应该启用交叉编译的特性。需要的话可用下面的命令安装：

```sh
brew install go --cross-compile-common
```

*你应该参照你所用的包管理工具给出的建议，不要依赖这种特定操作系统的安装方式。*

你可以用下面的命令验证安装是否成功：

```sh
$ go version
go version go1.10 darwin/amd64
```

## Go 环境

Go 有一套特别的惯例。

按照惯例，所有 Go 代码都存在于一个工作区（文件夹）中。这个工作区可以在你的机器的任何地方。如果你不指定，Go 将假定 $HOME/go 为默认工作区。工作区由环境变量 [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable) 标识并修改。

你应该设置环境变量，便于以后可在脚本或 shell 中使用它。

在你的 .bash_profile 中添加以下导出语句：

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

提示：你应该打开一个新的 shell 来使这些变量生效。

Go 假设你的工作区包含一个特定的目录结构。

Go 把文件放到三个目录中：所有源代码位于 src，包对象位于 pkg，编译好的程序位于 bin。你可以参照以下方式创建目录。

```sh
mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
```

此时，你可以用 _go get_，它会把下载到的资源按 src/pkg/bin 结构正确安装在相应的 $GOPATH/xxx 目录中。

## Go 编辑器

编辑器是可根据个人口味定制化的，可能你已经对你的编辑器配置了 Go 的支持。如果还没有，你可以考虑一下像 [Visual Studio Code](https://code.visualstudio.com) 这类对 Go 有良好支持的编辑器。

因为 VS Code 是一个图形界面程序，用 brew 安装 VS Code 需要安装一个叫 cask 的扩展。

```sh
brew tap caskroom/cask
```

此时你可以使用 brew 安装 VS Code 了。

```sh
brew cask install visual-studio-code
```

然后可以运行以下 shell 命令来验证 VS Code 是否安装正确。

```sh
code .
```

默认 VS Code 只自带了少量功能集成，你可以通过安装扩展以集成更多功能。如果想添加对 Go 的支持，你需要安装一些扩展。VS Code 有大量这类扩展，其中一个很棒的扩展是 [Luke Hoban 的 vscode-go](https://github.com/Microsoft/vscode-go)。可通过以下方法安装：

```sh
code --install-extension lukehoban.Go
```

当你第一次打开一个 Go 文件时，它会提示缺少分析工具，你应该点击按钮来安装它们。由 VS Code 安装（和使用）的工具列表可在[这里](https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on)找到。

## Go 调试

在 VS Code 中调试 Go 代码的一个好的选择是 Delve。可通过以下 go get 命令安装：

```sh
go get -u github.com/derekparker/delve/cmd/dlv
```

## Go 语法检查

对默认的语法检查进行增强可以用 [Gometalinter](https://github.com/alecthomas/gometalinter)。可通过以下方式安装：

```sh
go get -u github.com/alecthomas/gometalinter
gometalinter --install
```

## 重构和工具

本书的重心在于重构的重要性。

好的工具可以帮你放心地进行大型的重构。

你应该足够熟悉你的编辑器，以便使用简单的组合键执行以下操作：

- **提取/内联变量**。能够找出魔法值（magic value）并给他们一个名字可以让你快速简化你的代码。
- **提取方法/功能**。能够获取一段代码并提取函数/方法至关重要。
- **改名**。你应该能够自信地对多个文件内的符号批量重命名。
- **格式化**。Go 有一个名为 `go fmt` 的专有格式化程序。你的编辑器应该在每次保存文件时都运行它。
- **运行测试**。毫无疑问，你应该能够做到以上任何一点，然后快速重新运行你的测试，以确保你的重构没有破坏任何东西。

另外，为了对你处理代码更有帮助，你应该能够：

- **查看函数签名** - 在 Go 中调用某个函数时，你应该了解并熟悉它。你的 IDE 应根据其文档，参数以及返回的内容描述一个函数。
- **查看函数定义** - 如果仍不清楚函数的功能，你应该能够跳转到源代码并尝试自己弄清楚。
- **查找符号的用法** - 能够查看被调用函数的上下文可以在重构时帮你做出决定。

运用好你的工具将帮助你专注于代码并减少上下文切换。

## 总结

此时你应该已经安装了 Go，一个可用的编辑器和一些基本的工具。Go 拥有非常庞大的第三方生态系统。我们在本章提到了一些有用的组件，有关更完整的列表，请参阅 https://awesome-go.com 。

---

作者：[Chris James](https://dev.to/quii)
译者：[pityonline](https://github.com/pityonline)
校对：[校对 ID](https://github.com/校对ID)

本文由 [GCTT](https://github.com/studygolang/GCTT) 原创编译，[Go 中文网](https://studygolang.com/) 荣誉推出
