# Install Go, set up environment for productivity

The official installation instructions for Go are available [here](https://golang.org/doc/install).

This guide will assume that you are using a package manager for e.g. [Homebrew](https://brew.sh), [chocolatey](https://chocolatey.org), [Apt](https://help.ubuntu.com/community/AptGet/Howto) or [Yum](https://access.redhat.com/solutions/9934).

For demonstration purposes we will show the installation proceddure for OSX using Homebrew.

## Installation

The process of installation is very easy. First, what you have to do is to run this command to install homebrew (brew).  Brew has a dependency on xcode so you should ensure this is installed first.

```sh
xcode-select --install
```

Then you run the following to install homebrew:

```sh
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

At this point you can now install Go:

```sh
brew install go
```

If you are going to deploy your programs to linux based servers, you should enable cross compilation feature. If so, install using the following command:

```sh
brew install go --cross-compile-common
```

*You should follow any instructions reccomended by your package manager, not these may be host os specific*.

You can verify the installation with:

```sh
$ go version
go version go1.10 darwin/amd64
```

## Go Environment

Go is opinionated.

By convention, all Go code lives within a single workspace (folder). This workspace could be anywhere in your machine. If you don't specify, Go will assume $HOME/go as the default workspace.  The workspace is identified (and modified) by the environment variable [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable).

You should set the evnironment variable so that you can use it later in scripts, shells, etc.

Update your bash_profile to contain the following exports:

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

*Note* you should open a new shell to pickup these environment variables.

Go assumes that your workspace contains a specific directory structure.

Go places its files in three directories: All source code lives in src, package objects lives in pkg, and the compiled programs live in bin. You can create these directories as follows.

```sh
mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
```

At this point you can _go get_ and the src/package/bin will be installed correctly in the appropriate $GOPATH/xxx directory.

## Go Editor

Editor preference is very individualistic, you may already have a preference that supports Go.  If you don't you should consider an Editor such as [Visual Studio Code](https://code.visualstudio.com), which has execptional Go support.

To install vs code using brew, because this is a GUI application you need an extension to homebrew called cask to support install vs code.

```sh
brew tap caskroom/cask
```

At this point you can now use brew to install vs code.

```sh
brew cask install visual-studio-code
```

You can confirm vs code installed correctly you can run the following in your shell.

```sh
code . 
```

vs code is shipped with very little software enabled, you enable new software by installing extensions.  To add Go support you must install a extension, there are a variety available for vs code, an exceptional one is [Luke Hobans package](https://github.com/Microsoft/vscode-go).  This can be installed as follows:

```sh
code --install-extension lukehoban.Go
``` 

When you open a Go file for the first time in vs code, it will indicate that the Analysis tools are missing, you should click the button to install these. The list of tools that gets installed (and used) by vs code are available [here](https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on).

## Go Debugger

A good option for debugging Go (that's integrated with vs code) is Delve. This can be installed as follows using go get:

```sh
go get -u github.com/derekparker/delve/cmd/dlv
```

## Go Linting

An improvement over the default linter can be configured using [Gometalinter](https://github.com/alecthomas/gometalinter).

This can be installed as follows:

```sh
go get -u github.com/alecthomas/gometalinter
gometalinter --install
```

## Wrapping up

At this point you should have Go installed, an editor available and some basic tooling in place.  Go has a very large ecosystem of third party products.  We have identified a few useful components here, for a more complete list see https://awesome-go.com.