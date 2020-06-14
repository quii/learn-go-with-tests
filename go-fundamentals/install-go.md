---
description: Install Go
---

# Goをインストールする

Goの公式インストール手順は、[こちら](https://golang.org/doc/install)から入手できます。

このガイドでは、パッケージマネージャーを使用していることを前提としています。 [Homebrew](https://brew.sh), [Chocolatey](https://chocolatey.org), [Apt](https://help.ubuntu.com/community/AptGet/Howto) or [yum](https://access.redhat.com/solutions/9934).

デモのために、Homebrewを使用したOSXのインストール手順を紹介します。

## インストール

インストールのプロセスは非常に簡単です。
まず、このコマンドを実行して自作をインストールする必要があります。これはXcodeに依存しているため、これが最初にインストールされていることを確認する必要があります。

```bash
xcode-select --install
```

次に、以下を実行して自作をインストールします。

```bash
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

この時点で、次のようにGoをインストールできます。

```bash
brew install go
```

※ パッケージマネージャーが推奨する指示に従ってください。 **注**これらはホストos固有の場合があります。

インストールは次の方法で確認できます。

```bash
$ go version
go version go1.14 darwin/amd64
```

## Go環境

### $GOPATH

Go は意欲的です。

慣例により、すべてのGoコードは単一のワークスペース \(folder\)内にあります。このワークスペースは、マシンのどこにあってもかまいません。指定しない場合、Goはデフォルトのワークスペースとして  `$HOME/go` を想定します。ワークスペースは、 特定された \(and modified\) 環境変数 [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable)によって、変更されています。

後でスクリプトやシェルなどで使用できるように、環境変数を設定する必要があります。

次のエクスポートが含まれるように `.bash_profile`を更新します。

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

_注_これらの環境変数を取得するには、新しいシェルを開く必要があります。

Goは、ワークスペースに特定のディレクトリ構造が含まれていることを前提としています。

Goはそのファイルを3つのディレクトリに配置します。すべてのソースコードは`src`にあり、パッケージオブジェクトは`pkg`にあり、コンパイルされたプログラムは`bin`にあります。これらのディレクトリは次のように作成できます。

```bash
mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
```

この時点で、 `go get` を実行すると、 `src/package/bin` が適切な `$GOPATH/xxx` ディレクトリに正しくインストールされます。

### Go モジュール

Go 1.11 は [モジュール](https://github.com/golang/go/wiki/Modules), 導入し、代替ワークフローを可能にしました。この新しいアプローチは徐々にデフォルトになります。 [become the default](https://blog.golang.org/modules2019) モードになり、 `GOPATH`の使用は廃止されます。

モジュールは、依存関係の管理、バージョンの選択、および再現可能なビルドに関連する問題の解決を目的としています。
また、ユーザーが`GOPATH`の外でGoコードを実行できるようにします。

モジュールの使用は非常に簡単です。プロジェクトのルートとして`GOPATH`以外のディレクトリを選択し、`go mod init`コマンドで新しいモジュールを作成します。

モジュールパス、Goバージョン、およびその依存関係要件を含む `go.mod`ファイルが生成されます。これらは、正常なビルドに必要な他のモジュールです。

`<modulepath>`が指定されていない場合、 `go mod init`はディレクトリ構造からモジュールパスを推測しようとしますが、引数を指定してオーバーライドすることもできます。

```bash
mkdir my-project
cd my-project
go mod init <modulepath>
```

`go.mod`ファイルは次のようになります。

```text
module cmd

go 1.12

require (
        github.com/google/pprof v0.0.0-20190515194954-54271f7e092f
        golang.org/x/arch v0.0.0-20190815191158-8a70ba74b3a1
        golang.org/x/tools v0.0.0-20190611154301-25a4f137592f
)
```

組み込みのドキュメントは、利用可能なすべての `go mod`コマンドの概要を提供します。

```bash
go help mod
go help mod init
```

## Go エディター

エディターの設定は非常に個性的です。Goをサポートする設定がすでにある可能性があります。そうでない場合は、[Visual Studio Code]（https://code.visualstudio.com）などの優れたGoサポートを備えたエディターを検討する必要があります。

次のコマンドを使用してインストールできます。

```bash
brew cask install visual-studio-code
```

VS Codeが正しくインストールされていることを確認できます。シェルで以下を実行できます。

```bash
code .
```

VS Codeはほとんどソフトウェアが有効になっていない状態で出荷されます。拡張機能をインストールすることで新しいソフトウェアを有効にできます。 Goサポートを追加するには、拡張機能をインストールする必要があります。VSCodeにはさまざまなものがありますが、例外は[Luke Hobanのパッケージ]（https://github.com/Microsoft/vscode-go）です。これは次のようにインストールできます。

```bash
code --install-extension ms-vscode.go
```

VS Codeで初めてGoファイルを開くと、分析ツールが見つからないことが示されます。これらをインストールするには、ボタンをクリックする必要があります。 VS Codeによってインストールされるツールのリストは、[こちら]（https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-依存します）。

## Go デバッガー

Go（VS Codeに統合されている）のデバッグに適したオプションは`Delve`です。これは次のようにインストールできます。

```bash
go get -u github.com/go-delve/delve/cmd/dlv
```

VS CodeでGoデバッガーを構成および実行するための追加のヘルプについては、[VS Codeデバッグドキュメント]（https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code）を参照してください）。

## Go リンター

デフォルトのリンターに対する改良は、[GolangCI-Lint]（https://github.com/golangci/golangci-lint）を使用して構成できます。

これは次のようにインストールできます。

```bash
brew install golangci/tap/golangci-lint
```

## リファクタリングとツール

この本では、リファクタリングの重要性を強調しています。

ツールは、自信を持って大きなリファクタリングを行うのに役立ちます。

簡単なキーの組み合わせで次のことを実行できるように、エディタを十分に理解している必要があります。

* **Extract/Inline variable**. 変数値を取り、それらに名前を付けることができることで、コードをすばやく単純化できます。
* **Extract method/function**. コードのセクションを取得し、関数/メソッドを抽出できることが重要です。
* **Rename**. ファイル間でシンボルの名前を自信を持って変更できるはずです。
* **go fmt**. Goには、 `go fmt`と呼ばれる定型化されたフォーマッターがあります。エディターは、ファイルを保存するたびにこれを実行する必要があります。
* **Run tests**. 上記のいずれかを実行してから、テストをすばやく再実行して、リファクタリングによって何も壊れていないことを確認する必要があります。

さらに、コードの操作を支援するために、次のことができる必要があります。

* **View function signature** - Goで関数を呼び出す方法がわからないはずです。 IDEは、そのドキュメント、パラメータ、および返されるものの観点から関数を記述する必要があります。
* **View function definition** - 関数の機能がまだ明確でない場合は、ソースコードにジャンプして、自分でそれを理解できるようにする必要があります。
* **Find usages of a symbol** - 呼び出されている関数のコンテキストを確認できると、リファクタリング時の決定プロセスに役立ちます。

ツールを習得することで、コードに集中し、コンテキストの切り替えを減らすことができます。

## まとめ

この時点で、Goがインストールされ、エディターが利用可能で、いくつかの基本的なツールが整っているはずです。 Goには、サードパーティ製品の非常に大きなエコシステムがあります。ここでいくつかの有用なコンポーネントを特定しました。より完全なリストについては、[https://awesome-go.com]（https://awesome-go.com）を参照してください。
