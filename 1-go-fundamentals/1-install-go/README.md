# Install Go, set up environment for productivity

Go 공식 설치 방법은 [여기](https://golang.org/doc/install)에서 찾을 수 있습니다.

## Go Environment

### Go Modules
Go 1.11 버전에서 [모듈](https://github.com/golang/go/wiki/Modules)이 추가되었습니다. 모듈은 Go 1.16 버전부터 기본 방식으로 사용됨으로 `GOPATH`의 사용은 지양됩니다.

모듈은 의존성 관리, 버전 선택, 재현 가능한 빌드와 같은 문제들을 해결하기 위해 사용되며, `GOPATH` 밖에서도 Go 코드를 실행할 수 있도록 합니다.

모듈을 사용하기 매우 간단합니다. 먼저 `GOPATH` 밖의 아무 디렉토리를 프로젝트의 루트 디렉토리로 선택하고, `go mod init` 명령어로 새로운 모듈을 생성하면 됩니다.

명령어를 실행하면 모듈 경로와, Go 버전, 그리고 의존성 필요조건이 담긴 `go.mod` 파일이 생성되며 이 파일은 다른 모듈이 성공적으로 빌드되기 위해 필요합니다.

만약 아무 `<modulepath>`가 명시되지 않으면, `go mod init` 명령어는 디렉토리 구조로부터 모듈 경로를 추측하여 생성합니다. 인자를 제공함으로 해당 방식을 오버라이딩 하는 것도 가능합니다.


```sh
mkdir my-project
cd my-project
go mod init <modulepath>
```

`go.mod` 파일은 다음과 같은 형태로 생성됩니다.

```
module learn-go-with-tests

go 1.19
```
내장된 문서는 사용 가능한 모든 `go mod` 명령어 설명을 제공합니다.


```sh
go help mod
go help mod init
```

## Go Linting

An improvement over the default linter can be configured using [GolangCI-Lint](https://golangci-lint.run).

This can be installed as follows:

```sh
brew install golangci-lint
```

## Refactoring and your tooling

A big emphasis of this book is the importance of refactoring.

Your tools can help you do bigger refactoring with confidence.

You should be familiar enough with your editor to perform the following with a simple key combination:

- **Extract/Inline variable**. Being able to take magic values and give them a name lets you simplify your code quickly.
- **Extract method/function**. It is vital to be able to take a section of code and extract functions/methods
- **Rename**. You should be able to confidently rename symbols across files.
- **go fmt**. Go has an opinioned formatter called `go fmt`. Your editor should be running this on every file save.
- **Run tests**. You should be able to do any of the above and then quickly re-run your tests to ensure your refactoring hasn't broken anything.

In addition, to help you work with your code you should be able to:

- **View function signature**. You should never be unsure how to call a function in Go. Your IDE should describe a function in terms of its documentation, its parameters and what it returns.
- **View function definition**. If it's still not clear what a function does, you should be able to jump to the source code and try and figure it out yourself.
- **Find usages of a symbol**. Being able to see the context of a function being called can help your decision process when refactoring.

Mastering your tools will help you concentrate on the code and reduce context switching.

## Wrapping up

At this point you should have Go installed, an editor available and some basic tooling in place. Go has a very large ecosystem of third party products. We have identified a few useful components here. For a more complete list, see [https://awesome-go.com](https://awesome-go.com).
