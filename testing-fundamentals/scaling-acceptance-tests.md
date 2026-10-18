---
description: Learn Go with Tests - 受け入れテストのスケーリング（とgRPC入門）
---

# Learn Go with Tests - 受け入れテストのスケーリング（とgRPC入門）

この章は[受け入れテスト入門](intro-to-acceptance-tests.md)の続編です。[この章の完成したコードはGitHubにあります](https://github.com/quii/go-specs-greet)。

受け入れテストは不可欠なものであり、妥当な変更コストでシステムを自信を持って進化させ続けられるかどうかに直接影響します。

受け入れテストは、レガシーコードを扱う際にも素晴らしいツールです。テストがまったくない質の悪いコードベースに直面したときは、リファクタリングを始めたい誘惑に抗ってください。代わりに、いくつかの受け入れテストを書いて、外部から見た機能的な振る舞いに影響を与えることなく、システムの内部を自由に変更できるセーフティネットを用意しましょう。受け入れテスト（AT）は内部品質を気にする必要がないので、このような状況にぴったりです。

この章を読み終えるころには、受け入れテストが検証に役立つだけでなく、システムをより慎重かつ系統立てて変更する手助けをすることで、開発プロセスにおいても活用でき、無駄な労力を減らせることがわかるでしょう。

## 前提となる資料

この章は、長年にわたる受け入れテストへのフラストレーションから生まれました。ぜひ見ていただきたい動画が2つあります。

- Dave Farley - [How to write acceptance tests](https://www.youtube.com/watch?v=JDD5EEJgpHU)
- Nat Pryce - [E2E functional tests that can run in milliseconds](https://www.youtube.com/watch?v=Fk4rCn4YLLU)

『Growing Object Oriented Software』（GOOS）は、私を含め多くのソフトウェアエンジニアにとって非常に重要な本です。この本が説くアプローチは、私が一緒に働くエンジニアに従うよう指導しているものです。

- [GOOS](http://www.growing-object-oriented-software.com) - Nat Pryce & Steve Freeman

最後に、[Riya Dattani](https://twitter.com/dattaniriya)と私は、BDDの文脈でこのトピックについて[Acceptance tests, BDD and Go](https://www.youtube.com/watch?v=ZMWJCk_0WrY)という講演で話しました。

## おさらい

ここで話しているのは、システムが外部から、つまり「**ビジネスの観点**」から期待どおりに振る舞うことを検証する「ブラックボックス」テストのことです。テストはテスト対象のシステムの内部にアクセスできません。システムが**どのように**動くかではなく、**何を**するかだけに関心を持ちます。

## 悪い受け入れテストの解剖

長年にわたり、私はいくつもの会社やチームで働いてきました。どのチームも受け入れテスト、つまりユーザーの視点からシステムをテストし、意図したとおりに動くことを検証する何らかの方法の必要性を認識していました。しかし、ほぼ例外なく、これらのテストのコストがチームにとって本当の問題になっていました。

- 実行が遅い
- 壊れやすい
- 不安定（フレーキー）
- メンテナンスにコストがかかり、ソフトウェアの変更を本来よりも難しくしているように思える
- 特定の環境でしか実行できず、フィードバックループが遅く質も悪くなる

構築中のWebサイトに対して受け入れテストを書くとしましょう。ヘッドレスWebブラウザ（[Selenium](https://www.selenium.dev)など）を使って、ユーザーがWebサイトのボタンをクリックする動作をシミュレートし、サイトが必要なことをしているかを検証することにしました。

時間が経つにつれ、新しい機能が見つかるとWebサイトのマークアップを変更しなければならなくなり、エンジニアたちは何かを`<article>`にすべきか`<section>`にすべきかについて、もう10億回目になる自転車置き場の議論（bike-shed）を繰り広げます。

チームが実際のユーザーにはほとんど気づかれないような小さな変更しかシステムに加えていないにもかかわらず、受け入れテストの更新に多くの時間を浪費していることに気づくでしょう。

### 密結合

受け入れテストを変更するきっかけとなるものについて考えてみましょう。

- 外部の振る舞いの変更。システムが何をするかを変えたいのであれば、受け入れテストスイートを変更するのは妥当であり、むしろ望ましいことでさえあります。
- 実装の詳細の変更 / リファクタリング。理想的には、これは変更のきっかけになるべきではありません。なるとしても、ごく小さな変更であるべきです。

しかし、受け入れテストを変更しなければならない理由が後者であることがあまりにも多いのです。テストを更新する手間を想像して、エンジニアがシステムを変更することをためらうようになるほどです！

![Riya and myself talking about separating concerns in our tests](https://i.imgur.com/bbG6z57.png)

これらの問題は、先に挙げた著者たちが書いた、確立され実践されてきたエンジニアリングの習慣を適用していないことに起因しています。**受け入れテストを単体テストのように書くことはできません**。より多くの思考と、異なるプラクティスが必要なのです。

## 良い受け入れテストの解剖

振る舞いを変更したときだけ変更が必要になり、実装の詳細を変更したときには変更が不要な受け入れテストが欲しいのであれば、それらの関心事を分離する必要があるのは当然のことです。

### 複雑さの種類について

ソフトウェアエンジニアとして、私たちは2種類の複雑さに対処しなければなりません。

- **偶有的な複雑さ（Accidental complexity）** は、コンピュータを扱っているがゆえに対処しなければならない複雑さです。ネットワーク、ディスク、APIなどがこれにあたります。

- **本質的な複雑さ（Essential complexity）** は、「ドメインロジック」と呼ばれることもあります。ドメイン内の特定のルールや真理のことです。
  - たとえば、「口座の所有者が利用可能な額よりも多くのお金を引き出した場合、その口座は当座貸越になる」というものです。この文はコンピュータについて何も言っていません。この文は、銀行でコンピュータが使われるようになる前から真実でした！

本質的な複雑さは技術者ではない人にも説明できるものであるべきで、それを「ドメイン」コードと受け入れテストの中でモデル化しておくことには価値があります。

### 関心の分離

先ほどの動画でDave Farleyが提案したこと、そしてRiyaと私も議論したことは、**仕様（specification）** という考え方を持つべきだということです。仕様は、偶有的な複雑さや実装の詳細と結合することなく、私たちが望むシステムの振る舞いを記述します。

この考え方は妥当に感じられるはずです。本番コードでは、私たちは頻繁に関心事を分離し、作業の単位を疎結合にしようと努めています。`HTTP`ハンドラーをHTTP以外の関心事から切り離すために、ためらうことなく`interface`を導入するのではないでしょうか？同じ考え方を受け入れテストにも適用してみましょう。

Dave Farleyは特定の構造について説明しています。

![Dave Farley on Acceptance Tests](https://i.imgur.com/nPwpihG.png)

GopherconUKで、Riyaと私はこれをGoの用語に置き換えました。

![Separation of concerns](https://i.imgur.com/qdY4RJe.png)

### テストを強化する

仕様をどのように実行するかを切り離すことで、さまざまなシナリオで仕様を再利用できるようになります。次のことが可能です。

#### ドライバーを設定可能にする

これにより、受け入れテストをローカル、ステージング環境、そして（理想的には）本番環境で実行できます。
- あまりにも多くのチームが、受け入れテストをローカルで実行できないようにシステムを設計しています。これは耐えがたいほど遅いフィードバックループをもたらします。コードを統合する _前に_ 受け入れテストが通ると確信できたほうがよいと思いませんか？テストが壊れ始めたとき、その失敗をローカルで再現できず、代わりに変更をコミットして、20分後に別の環境でテストが通ることを祈るしかないというのは許容できることでしょうか？
- 覚えておいてください。ステージングでテストが通ったからといって、システムが動くとは限りません。開発環境と本番環境の一致（Dev/Prod parity）は、よく言っても罪のない嘘です。[私は本番でテストします](https://increment.com/testing/i-test-in-production/)。
- 環境間には常に、システムの*振る舞い*に影響を与えうる違いがあります。CDNのキャッシュヘッダーが誤って設定されているかもしれませんし、依存している下流のサービスの振る舞いが異なるかもしれませんし、設定値が間違っているかもしれません。しかし、仕様を本番環境で実行してこれらの問題をすばやく捉えられたら素晴らしいと思いませんか？

#### _異なる_ ドライバーを差し込んで、システムの他の部分をテストする

この柔軟性により、異なる抽象化レイヤーやアーキテクチャレイヤーで振る舞いをテストでき、ブラックボックステストを超えた、より焦点を絞ったテストが可能になります。
- たとえば、背後にAPIがあるWebページがあるとします。同じ仕様を使って両方をテストしてみてはどうでしょうか？Webページにはヘッドレスブラウザを、APIにはHTTP呼び出しを使えます。
- この考え方をさらに進めると、理想的には**コードで本質的な複雑さをモデル化**したい（「ドメイン」コードとして）ので、仕様を単体テストにも使えるはずです。これにより、システムの本質的な複雑さがモデル化され、正しく振る舞っていることについて素早いフィードバックが得られます。


### 正しい理由で変更される受け入れテスト

このアプローチでは、仕様が変更される唯一の理由はシステムの振る舞いが変わった場合であり、それは妥当なことです。

- HTTP APIを変更しなければならない場合、更新すべき明白な場所が1つあります。それはドライバーです。
- マークアップが変更された場合も、やはりその特定のドライバーを更新します。

システムが成長するにつれ、複数のテストでドライバーを再利用するようになります。つまり、実装の詳細が変わっても、通常は明白な1か所を更新するだけで済むということです。

正しく行えば、このアプローチは実装の詳細における柔軟性と、仕様における安定性をもたらします。重要なのは、変更を管理するためのシンプルで明白な構造を提供してくれることです。これはシステムとそのチームが成長するにつれて不可欠になります。

### ソフトウェア開発の手法としての受け入れテスト

講演の中で、Riyaと私は受け入れテストとBDDとの関係について議論しました。 _解決しようとしている問題を理解する_ ことから作業を始め、それを仕様として表現することが、意図を明確にするのに役立ち、作業を始めるための素晴らしい方法であることについて話しました。

私が初めてこの働き方に出会ったのはGOOSでした。少し前に、そのアイデアを自分のブログにまとめました。以下は私の投稿[Why TDD](https://quii.dev/The_Why_of_TDD)からの抜粋です。

---

TDDは、必要な振る舞いを正確に、反復的に設計できるようにすることに重点を置いています。新しい領域に取り組み始めるときは、鍵となる必要な振る舞いを特定し、スコープを積極的に削らなければなりません。

「トップダウン」アプローチに従い、外部から振る舞いを実行する受け入れテスト（AT）から始めましょう。これがあなたの取り組みの北極星として機能します。集中すべきなのは、そのテストを通すことだけです。テストを通すのに十分なコードを開発する間、このテストはしばらく失敗し続けるでしょう。

![](https://i.imgur.com/pxTaYu4.png)

受け入れテストを準備したら、TDDのプロセスに入り、受け入れテストを通すのに十分なユニットを導き出していきます。コツは、この時点では設計についてあまり心配しすぎないことです。まだ問題について学び、探索している段階なので、受け入れテストを通すのに十分なコードを書きましょう。

この最初の一歩は、思っている以上に大がかりになることがよくあります。Webサーバーのセットアップ、ルーティング、設定などがあるからです。だからこそ、作業のスコープを小さく保つことが不可欠なのです。真っ白なキャンバスに最初の前向きな一歩を踏み出し、それを通過する受け入れテストで裏付けることで、迅速かつ安全に反復を続けられるようにしたいのです。

![](https://i.imgur.com/t5y5opw.png)

開発を進めながら、テストに耳を傾けてください。テストは、設計をより良い方向へ押し進めるのに役立つシグナルを与えてくれるはずです。ただし、ここでも想像ではなく振る舞いに根ざして進めます。

通常、受け入れテストを通すために大変な仕事をする最初の「ユニット」は、このわずかな振る舞いに対してさえ、快適とは言えないほど大きくなるでしょう。このときこそ、問題をどのように分解し、新しい協力者（コラボレーター）を導入するかを考え始めるタイミングです。

![](https://i.imgur.com/UYqd7Cq.png)

ここでテストダブル（フェイクやモックなど）が役に立ちます。ソフトウェア内部に存在する複雑さの大部分は、通常、実装の詳細ではなく、ユニットの「間」、つまりユニット同士がどのように相互作用するかにあるからです。

#### ボトムアップの危険性

これは「ボトムアップ」ではなく「トップダウン」のアプローチです。ボトムアップにも用途はありますが、リスクの要素を伴います。アプリケーションにすばやく統合することなく、また高レベルのテストで検証することなく「サービス」やコードを構築すると、**検証されていないアイデアに多くの労力を無駄に費やすリスクがあります**。

これは受け入れテスト駆動のアプローチの重要な特性であり、テストを使ってコードの本当の検証を得るということです。

ボトムアップで、単独で、何かの仕事を解決してくれると思ってコードの塊を作ったエンジニアに、私はあまりにも何度も出会ってきました。しかし、そのコードは次のようなものでした。

- 私たちが望むように動かない
- 必要のないことをする
- 簡単に統合できない
- 結局、大量の書き直しが必要になる

これは無駄です。

## 話はここまでにして、コードを書きましょう

他の章とは異なり、アプリケーションをコンテナで実行するため、[Docker](https://www.docker.com)をインストールしておく必要があります。本書のこの時点では、Goのコードを書いたり、異なるパッケージからインポートしたりすることには慣れているものと想定しています。

`go mod init github.com/quii/go-specs-greet`で新しいプロジェクトを作成します（ここには好きなものを指定できますが、パスを変更した場合は、内部のすべてのインポートをそれに合わせて変更する必要があります）。

仕様を格納するための`specifications`フォルダを作成し、`greet.go`ファイルを追加します。

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet() (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet()
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, world")
}
```

私のIDE（Goland）は依存関係の追加という面倒な作業を代わりにやってくれますが、手動で行う必要がある場合は次のようにします。

`go get github.com/alecthomas/assert/v2`

Farleyの受け入れテストの設計（仕様 -> DSL -> ドライバー -> システム）に従い、これで実装から切り離された仕様ができました。仕様は、私たちが _どのように_ `Greet`するかを知りませんし、気にもしません。ドメインの本質的な複雑さだけに関心を持っています。確かに今のところこの複雑さは大したものではありませんが、反復を重ねるにつれて仕様を拡張し、機能を追加していきます。小さく始めることは常に重要です！

このインターフェースは、DSLの最初の一歩と見ることもできます。プロジェクトが成長するにつれて、異なる形で抽象化する必要が出てくるかもしれませんが、今のところはこれで問題ありません。

この時点で、仕様を実装から切り離すためのこの程度の儀式的な手続きを見て、「過度な抽象化だ」と非難する人もいるかもしれません。**実装と結合しすぎた受け入れテストは、エンジニアリングチームにとって本当の重荷になることを約束します**。世の中にあるほとんどの受け入れテストは、過度に抽象化されているからではなく、この不適切な結合のせいでメンテナンスにコストがかかっているのだと私は確信しています。

この仕様を使って、`Greet`できるあらゆる「システム」を検証できます。

### 最初のシステム: HTTP API

私たちは「挨拶サービス」をHTTPで提供する必要があります。そのため、次のものを作成する必要があります。

1. **ドライバー**。この場合、**HTTPクライアント**を使ってHTTPシステムを操作します。このコードは、私たちのAPIをどのように扱えばよいかを知っています。ドライバーはDSLをシステム固有の呼び出しに変換します。私たちの場合、ドライバーは仕様が定義するインターフェースを実装します。
2. 挨拶APIを持つ**HTTPサーバー**
3. **テスト**。サーバーを起動するライフサイクルを管理し、ドライバーを仕様に差し込んでテストとして実行する責任を持ちます。

## 最初にテストを書く

プログラムをコンパイルして実行し、テストを実行し、その後すべてをクリーンアップするブラックボックステストを作成する最初のプロセスは、かなり手間がかかることがあります。だからこそ、プロジェクトの開始時に最小限の機能で行うのが望ましいのです。私は通常、すべてのプロジェクトを「hello world」サーバーの実装から始め、すべてのテストを準備しておくことで、実際の機能をすばやく構築できるようにしています。

「仕様」「ドライバー」「受け入れテスト」というメンタルモデルに慣れるには少し時間がかかるかもしれないので、注意深くついてきてください。最初に仕様を呼び出そうとすることで「逆算して進める」と役に立つことがあります。

出荷する予定のプログラムを格納するための構造を作成します。

`mkdir -p cmd/httpserver`

新しいフォルダの中に新しいファイル`greeter_server_test.go`を作成し、次の内容を追加します。

```go
package main_test

import (
	"testing"

	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	specifications.GreetSpecification(t, nil)
}
```

私たちは仕様をGoのテストで実行したいと考えています。すでに`*testing.T`にはアクセスできるので、それが最初の引数になりますが、2番目の引数はどうでしょうか？

`specifications.Greeter`はインターフェースであり、新しいTestGreeterServerのコードを次のように変更して、`Driver`で実装します。

```go
import (
	go_specs_greet "github.com/quii/go-specs-greet"
)

func TestGreeterServer(t *testing.T) {
	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080"}
	specifications.GreetSpecification(t, driver)
}
```

ローカルを含むさまざまな環境に対して実行できるように、`Driver`を設定可能にしておくのが望ましいので、`BaseURL`フィールドを追加しました。

## テストを実行してみます

```
./greeter_server_test.go:46:12: undefined: go_specs_greet.Driver
```

ここでもまだTDDを実践しています！これは私たちが踏み出さなければならない大きな最初の一歩です。いくつかのファイルを作成し、普段よりも多くのコードを書く必要があるかもしれませんが、始めたばかりのときにはよくあることです。レッドのステップのルールを忘れないようにすることがとても重要です。

> テストを通すために必要なだけ罪を犯しましょう

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

鼻をつまんでください。テストが通ったらリファクタリングできることを忘れずに。以下は`driver.go`に書くドライバーのコードで、プロジェクトのルートに配置します。

```go
package go_specs_greet

import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
}

func (d Driver) Greet() (string, error) {
	res, err := http.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```


注意点:

- さまざまな`if err != nil`を導き出すためにテストを書くべきだと主張することもできますが、私の経験では、`err`に対して何もしていない限り、「受け取ったエラーを返す」ことを確認するテストの価値は比較的低いです。
- **デフォルトのHTTPクライアントを使うべきではありません**。後でHTTPクライアントを渡してタイムアウトなどを設定できるようにしますが、今はとにかくテストを通すことを目指しています。
-  `greeter_server_test.go`では、作成したばかりの`go_specs_greet`パッケージのDriver関数を呼び出しているので、インポートに`github.com/quii/go-specs-greet`を追加するのを忘れないでください。
テストを再実行してみてください。コンパイルは通るようになりましたが、テストは通らないはずです。

```
Get "http://localhost:8080/greet": dial tcp [::1]:8080: connect: connection refused
```

`Driver`はありますが、まだアプリケーションを起動していないので、HTTPリクエストを送ることができません。テストを実行するために、受け入れテストでシステムのビルド、実行、そして最後に終了までを調整する必要があります。

### アプリケーションの実行

チームがデプロイのためにシステムのDockerイメージをビルドするのは一般的なことなので、テストでも同じことをします。

テストでDockerを使いやすくするために、[Testcontainers](https://golang.testcontainers.org)を使います。Testcontainersを使うと、プログラムからDockerイメージをビルドし、コンテナのライフサイクルを管理できます。

`go get github.com/testcontainers/testcontainers-go`

これで、`cmd/httpserver/greeter_server_test.go`を次のように編集できます。

```go
package main_test

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGreeterServer(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "./cmd/httpserver/Dockerfile",
			// set to false if you want less spam, but this is helpful if you're having troubles
			PrintBuildLog: true,
		},
		ExposedPorts: []string{"8080:8080"},
		WaitingFor:   wait.ForHTTP("/").WithPort("8080"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})

	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080"}
	specifications.GreetSpecification(t, driver)
}
```

テストを実行してみてください。

```
=== RUN   TestGreeterHandler
2022/09/10 18:49:44 Starting container id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Waiting for container id 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
2022/09/10 18:49:45 Container is ready id: 03e8588a1be4 image: docker.io/testcontainers/ryuk:0.3.3
    greeter_server_test.go:32: Did not expect an error but got:
        Error response from daemon: Cannot locate specified Dockerfile: ./cmd/httpserver/Dockerfile: failed to create container
--- FAIL: TestGreeterHandler (0.59s)
```

プログラム用のDockerfileを作成する必要があります。`httpserver`フォルダの中に`Dockerfile`を作成し、次の内容を追加します。

```dockerfile
# Make sure to specify the same Go version as the one in the go.mod file.
# For example, golang:1.22.1-alpine.
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/httpserver/*.go

EXPOSE 8080
CMD [ "./svr" ]
```

ここでの詳細はあまり気にしないでください。洗練させたり最適化したりできますが、この例ではこれで十分です。このアプローチの利点は、後でDockerfileを改善しても、意図どおりに動くことを証明するテストがあることです。これこそがブラックボックステストを持つことの真の強みです！

テストを再実行してみてください。イメージをビルドできないというエラーが出るはずです。もちろん、それはまだビルドするプログラムを書いていないからです！

テストを完全に実行するには、`8080`でリッスンするプログラムを作成する必要がありますが、**それだけです**。TDDの規律を守り、テストが期待どおりに失敗することを確認するまでは、テストを通すような本番コードを書かないでください。

`httpserver`フォルダの中に、次の内容で`main.go`を作成します。

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

もう一度テストを実行してみてください。次のように失敗するはずです。

```
    greet.go:16: Expected values to be equal:
        +Hello, World
        \ No newline at end of file
--- FAIL: TestGreeterHandler (2.09s)
```

## 成功させるのに十分なコードを書く

仕様が望むように振る舞うよう、ハンドラーを更新します。

```go
import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "Hello, world")
	})
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
```

## リファクタリング♪

厳密にはリファクタリングではありませんが、デフォルトのHTTPクライアントに頼るべきではないので、Driverを変更してクライアントを渡せるようにし、テストからそれを与えるようにしましょう。

```go
import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
	Client  *http.Client
}

func (d Driver) Greet() (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

`cmd/httpserver/greeter_server_test.go`のテストで、ドライバーの生成部分を更新してクライアントを渡すようにします。

```go
client := http.Client{
	Timeout: 1 * time.Second,
}

driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080", Client: &client}
specifications.GreetSpecification(t, driver)
```

`main.go`はできるだけシンプルに保つのが良いプラクティスです。`main.go`は、作成した構成要素を組み合わせてアプリケーションにすることだけに関心を持つべきです。

プロジェクトのルートに`handler.go`というファイルを作成し、コードをそこに移動します。

```go
package go_specs_greet

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world")
}
```

代わりにそのハンドラーをインポートして使うように`main.go`を更新します。

```go
package main

import (
	"net/http"

	go_specs_greet "github.com/quii/go-specs-greet"
)

func main() {
	handler := http.HandlerFunc(go_specs_greet.Handler)
	http.ListenAndServe(":8080", handler)
}
```

## 振り返り

最初のステップは大変に感じられました。ハードコードされた文字列を返すHTTPハンドラーを作成してテストするために、いくつもの`go`ファイルを作りました。この「イテレーション0」の儀式とセットアップは、今後のイテレーションで大いに役立つでしょう。

機能の変更は、仕様を通じて駆動し、それによって強いられる変更に対処することで、シンプルかつ制御されたものになるはずです。これで受け入れテストのための`DockerFile`と`testcontainers`の準備ができたので、アプリケーションの構築方法が変わらない限り、これらのファイルを変更する必要はないはずです。

次の要件である、特定の人に挨拶することで、これを確認してみましょう。

## 最初にテストを書く

仕様を編集します。

```go
package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet(name string) (string, error)
}

func GreetSpecification(t testing.TB, greeter Greeter) {
	got, err := greeter.Greet("Mike")
	assert.NoError(t, err)
	assert.Equal(t, got, "Hello, Mike")
}
```

特定の人に挨拶できるようにするには、システムへのインターフェースを変更して`name`パラメータを受け取れるようにする必要があります。

## テストを実行してみます

```
./greeter_server_test.go:48:39: cannot use driver (variable of type go_specs_greet.Driver) as type specifications.Greeter in argument to specifications.GreetSpecification:
	go_specs_greet.Driver does not implement specifications.Greeter (wrong type for Greet method)
		have Greet() (string, error)
		want Greet(name string) (string, error)
```

仕様の変更により、ドライバーを更新する必要が生じました。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

特定の`name`に挨拶するよう求めるために、リクエストで`name`クエリの値を指定するようドライバーを更新します。

```go
import "io"

func (d Driver) Greet(name string) (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/greet?name=" + name)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	greeting, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(greeting), nil
}
```

これでテストは実行され、失敗するはずです。

```
    greet.go:16: Expected values to be equal:
        -Hello, world
        \ No newline at end of file
        +Hello, Mike
        \ No newline at end of file
--- FAIL: TestGreeterHandler (1.92s)
```

## 成功させるのに十分なコードを書く

リクエストから`name`を取り出して挨拶します。

```go
import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", r.URL.Query().Get("name"))
}
```

これでテストは通るはずです。

## リファクタリング♪

[HTTPハンドラーの再検討](../questions-and-answers/http-handlers-revisited.md)では、HTTPハンドラーはHTTPの関心事の処理だけに責任を持つべきであり、「ドメインロジック」はハンドラーの外に置くべきだということがいかに重要かについて議論しました。こうすることで、ドメインロジックをHTTPから切り離して開発でき、テストや理解がより簡単になります。

これらの関心事を分離しましょう。

`./handler.go`のハンドラーを次のように更新します。

```go
func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, Greet(name))
}
```

新しいファイル`./greet.go`を作成します。
```go
package go_specs_greet

import "fmt"

func Greet(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}
```

## 「アダプター」デザインパターンへのちょっとした寄り道

人に挨拶するというドメインロジックを別の関数に分離したので、greet関数の単体テストを自由に書けるようになりました。これは、文字列を得るために、Webサーバーにアクセスするドライバーを経由する仕様を通してテストするよりも、間違いなくずっとシンプルです！

ここでも仕様を再利用できたら素晴らしいと思いませんか？結局のところ、仕様のポイントは実装の詳細から切り離されていることです。仕様が**本質的な複雑さ**を捉えていて、「ドメイン」コードがそれをモデル化するものであるなら、それらを一緒に使えるはずです。

次のように`./greet_test.go`を作成して試してみましょう。

```go
package go_specs_greet_test

import (
	"testing"

	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecification(t, go_specs_greet.Greet)
}

```

これができれば素晴らしいのですが、うまくいきません。

```
./greet_test.go:11:39: cannot use go_specs_greet.Greet (value of type func(name string) string) as type specifications.Greeter in argument to specifications.GreetSpecification:
	func(name string) string does not implement specifications.Greeter (missing Greet method)
```

仕様は、関数ではなく`Greet()`メソッドを持つものを求めています。

このコンパイルエラーはもどかしいものです。`Greeter`であると「わかっている」ものがあるのに、コンパイラが使用を許可するための正しい**形**になっていないのです。これこそが**アダプター**パターンが対応するものです。

> [ソフトウェア工学](https://en.wikipedia.org/wiki/Software_engineering)において、**アダプターパターン**は[ソフトウェアデザインパターン](https://en.wikipedia.org/wiki/Software_design_pattern)の一つであり（[ラッパー](https://en.wikipedia.org/wiki/Wrapper_function)とも呼ばれ、これは[デコレーターパターン](https://en.wikipedia.org/wiki/Decorator_pattern)と共有される別名です）、既存の[クラス](https://en.wikipedia.org/wiki/Class_(computer_science))の[インターフェース](https://en.wikipedia.org/wiki/Interface_(computer_science))を別のインターフェースとして使えるようにするものです。[[1\]](https://en.wikipedia.org/wiki/Adapter_pattern#cite_note-HeadFirst-1) 既存のクラスの[ソースコード](https://en.wikipedia.org/wiki/Source_code)を変更することなく、他のクラスと連携させるためによく使われます。

比較的シンプルなものに対して、ずいぶんと仰々しい言葉が並んでいます。デザインパターンにはよくあることで、だからこそ人々はデザインパターンに呆れた顔をしがちなのです。デザインパターンの価値は特定の実装にあるのではなく、エンジニアが直面する一般的な問題に対する特定の解決策を説明するための言語にあります。共通の語彙を持つチームであれば、コミュニケーションの摩擦が減ります。

このコードを`./specifications/adapters.go`に追加します。

```go
type GreetAdapter func(name string) string

func (g GreetAdapter) Greet(name string) (string, error) {
	return g(name), nil
}
```

これで、テストでアダプターを使って`Greet`関数を仕様に差し込めるようになりました。

```go
package go_specs_greet_test

import (
	"testing"

	gospecsgreet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		specifications.GreetAdapter(gospecsgreet.Greet),
	)
}
```

アダプターパターンは、インターフェースが求める振る舞いを示す型があるものの、それが正しい形になっていない場合に便利です。

## 振り返り

振る舞いの変更はシンプルに感じられたのではないでしょうか？まあ、単に問題の性質によるものだったのかもしれませんが、この作業方法は規律と、システムを上から下まで変更するためのシンプルで繰り返し可能な方法を与えてくれます。

- 問題を分析し、正しい方向へ進むためのシステムへのわずかな改善を特定する
- 新しい本質的な複雑さを仕様に捉える
- 受け入れテストが実行されるまでコンパイルエラーに従う
- 仕様に従ってシステムが振る舞うように実装を更新する
- リファクタリングする

最初のイテレーションの苦労の後は、仕様、ドライバー、実装が分離されているおかげで、受け入れテストのコードを編集する必要はありませんでした。仕様を変更したことでドライバーと最終的に実装を更新する必要がありましたが、システムをコンテナとして _どのように_ 起動するかに関するボイラープレートコードは影響を受けませんでした。

アプリケーションのDockerイメージをビルドしてコンテナを起動するオーバーヘッドがあっても、アプリケーション**全体**をテストするためのフィードバックループは非常に短いです。

```
quii@Chriss-MacBook-Pro go-specs-greet % go test ./...
ok  	github.com/quii/go-specs-greet	0.181s
ok  	github.com/quii/go-specs-greet/cmd/httpserver	2.221s
?   	github.com/quii/go-specs-greet/specifications	[no test files]
```

さて、CTOがgRPCこそが _未来_ だと決めたと想像してください。彼女は、既存のHTTPサーバーを維持しつつ、同じ機能をgRPCサーバーでも公開してほしいと考えています。

これは**偶有的な複雑さ**の一例です。偶有的な複雑さとは、コンピュータを扱っているがゆえに対処しなければならない複雑さ、つまりネットワーク、ディスク、APIなどのことだったのを思い出してください。**本質的な複雑さは変わっていない**ので、仕様を変更する必要はないはずです。

多くのリポジトリ構造やデザインパターンは、主に複雑さの種類を分離することを扱っています。たとえば、「ポートとアダプター」は、偶有的な複雑さに関わるあらゆるものからドメインコードを分離することを求めます。そのようなコードは「adapters」フォルダに置かれます。

### 変更を簡単にする

変更を加える _前に_ 、いくらかリファクタリングしておくのが理にかなっていることもあります。

> まず変更を簡単にし、それから簡単な変更をせよ

~Kent Beck

そのため、`http`のコード（`driver.go`と`handler.go`）を`adapters`フォルダ内の`httpserver`というパッケージに移動し、パッケージ名を`httpserver`に変更しましょう。

Greetメソッドを参照するために、`handler.go`にルートパッケージをインポートする必要があります。

```go
package httpserver

import (
	"fmt"
	"net/http"

	go_specs_greet "github.com/quii/go-specs-greet/domain/interactions"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, go_specs_greet.Greet(name))
}

```

httpserverアダプターをmain.goにインポートします。

```go
package main

import (
	"net/http"

	"github.com/quii/go-specs-greet/adapters/httpserver"
)

func main() {
	handler := http.HandlerFunc(httpserver.Handler)
	http.ListenAndServe(":8080", handler)
}
```

そして、greeter_server_test.goのインポートと`Driver`への参照を更新します。

```go
driver := httpserver.Driver{BaseURL: "http://localhost:8080", Client: &client}
```

最後に、ドメインレベルのコードも独自のフォルダにまとめておくと便利です。怠けて、プロジェクトに何百もの無関係な型や関数が入った`domain`フォルダを作るようなことはしないでください。ドメインについて考える努力をして、一緒にあるべきアイデアは一緒にまとめましょう。そうすることでプロジェクトが理解しやすくなり、インポートの質も向上します。

次のようなものを見るよりも

```go
domain.Greet
```

これは少し奇妙なので、代わりに次のようにするのが良いでしょう。

```go
interactions.Greet
```

すべてのドメインコードを格納する`domain`フォルダを作成し、その中に`interactions`フォルダを作成します。使っているツールによっては、いくつかのインポートやコードを更新する必要があるかもしれません。

プロジェクトツリーは次のようになっているはずです。

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── Makefile
├── README.md
├── adapters
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   └── httpserver
|       ├── Dockerfile
│       ├── greeter_server_test.go
│       └── main.go
├── domain
│   └── interactions
│       ├── greet.go
│       └── greet_test.go
├── go.mod
├── go.sum
└── specifications
    └── adapters.go
    └── greet.go

```

ドメインコード、つまり**本質的な複雑さ**はGoモジュールのルートに置かれ、それらを「現実世界」で使えるようにするコードは**アダプター**として整理されています。`cmd`フォルダは、これらの論理的なまとまりを組み合わせて実用的なアプリケーションにする場所であり、すべてが動くことを検証するブラックボックステストを備えています。いいですね！

最後に、受け入れテストを _ほんの少し_ 整理しましょう。受け入れテストの高レベルなステップを考えてみると、次のようになります。

- Dockerイメージをビルドする
- _何らかの_ ポートでリッスンするまで待つ
- DSLをシステム固有の呼び出しに変換する方法を理解しているドライバーを作成する
- ドライバーを仕様に差し込む

…gRPCサーバーの受け入れテストにも同じ要件があることに気づくでしょう！

`adapters`フォルダが置き場所として適切そうなので、`docker.go`というファイルの中で、最初の2つのステップを、次に再利用する関数にカプセル化します。

```go
package adapters

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartDockerServer(
	t testing.TB,
	port string,
	dockerFilePath string,
) {
	ctx := context.Background()
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../../.",
			Dockerfile:    dockerFilePath,
			PrintBuildLog: true,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

これで受け入れテストを少し整理する機会が得られます。

```go
func TestGreeterServer(t *testing.T) {
	var (
		port           = "8080"
		dockerFilePath = "./cmd/httpserver/Dockerfile"
		baseURL        = fmt.Sprintf("http://localhost:%s", port)
		driver         = httpserver.Driver{BaseURL: baseURL, Client: &http.Client{
			Timeout: 1 * time.Second,
		}}
	)

	adapters.StartDockerServer(t, port, dockerFilePath)
	specifications.GreetSpecification(t, driver)
}
```

これで _次の_ テストを書くのが簡単になるはずです。

## 最初にテストを書く

この新しい機能は、ドメインコードとやり取りする新しい`adapter`を作成することで実現できます。そのため、私たちは次のようにできるはずです。

- 仕様を変更する必要はない
- 仕様を再利用できる
- ドメインコードを再利用できる

新しいプログラムと対応する受け入れテストを格納するために、`cmd`の中に新しいフォルダ`grpcserver`を作成します。`cmd/grpc_server/greeter_server_test.go`の中に受け入れテストを追加します。これはHTTPサーバーのテストとよく似ていますが、それは偶然ではなく設計によるものです。

```go
package main_test

import (
	"fmt"
	"testing"

	"github.com/quii/go-specs-greet/adapters"
	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	var (
		port           = "50051"
		dockerFilePath = "./cmd/grpcserver/Dockerfile"
		driver         = grpcserver.Driver{Addr: fmt.Sprintf("localhost:%s", port)}
	)

	adapters.StartDockerServer(t, port, dockerFilePath)
	specifications.GreetSpecification(t, &driver)
}
```

違いは次の点だけです。

- 別のプログラムをビルドするので、別のDockerファイルを使う
- つまり、新しいプログラムとやり取りするために`gRPC`を使う新しい`Driver`が必要になる

## テストを実行してみます

```
./greeter_server_test.go:26:12: undefined: grpcserver
```

まだ`Driver`を作成していないので、コンパイルできません。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`adapters`の中に`grpcserver`フォルダを作成し、その中に`driver.go`を作成します。

```go
package grpcserver

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	return "", nil
}
```

もう一度実行すると、 _コンパイル_ は通るようになりますが、Dockerfileと実行する対応プログラムをまだ作成していないので、テストは通りません。

`cmd/grpcserver`の中に新しい`Dockerfile`を作成します。

```dockerfile
# Make sure to specify the same Go version as the one in the go.mod file.
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/grpcserver/*.go

EXPOSE 50051
CMD [ "./svr" ]
```

そして`main.go`も作成します。

```go
package main

import "fmt"

func main() {
	fmt.Println("implement me")
}
```

サーバーがポートでリッスンしていないため、テストが失敗することがわかるはずです。いよいよgRPCでクライアントとサーバーの構築を始めるときです。

## 成功させるのに十分なコードを書く

### gRPC

gRPCに馴染みがない場合は、まず[gRPCのWebサイト](https://grpc.io)を見てみることをお勧めします。とはいえ、この章においては、gRPCはシステムへのアダプターの一種にすぎず、他のシステムが私たちの優れたドメインコードを呼び出す（**r**emote **p**rocedure **c**all、リモートプロシージャコール）ための手段です。

特徴的なのは、Protocol Buffersを使って「サービス定義」を定義することです。そして、その定義からサーバーとクライアントのコードを生成します。これはGoだけでなく、ほとんどの主要な言語でも機能します。つまり、Goを書いていないかもしれない社内の他のチームとも定義を共有でき、サービス間の通信をスムーズに行えるということです。

gRPCを使ったことがない場合は、**Protocol Bufferコンパイラ**といくつかの**Goプラグイン**をインストールする必要があります。[gRPCのWebサイトに、その方法についてのわかりやすい説明があります](https://grpc.io/docs/languages/go/quickstart/)。

新しいドライバーと同じフォルダの中に、次の内容で`greet.proto`ファイルを追加します。

```protobuf
syntax = "proto3";

option go_package = "github.com/quii/adapters/grpcserver";

package grpcserver;

service Greeter {
  rpc Greet (GreetRequest) returns (GreetReply) {}
}

message GreetRequest {
  string name = 1;
}

message GreetReply {
  string message = 1;
}
```

この定義を理解するのに、Protocol Buffersの専門家である必要はありません。Greetメソッドを持つサービスを定義し、入力と出力のメッセージ型を記述しています。

`adapters/grpcserver`の中で次のコマンドを実行して、クライアントとサーバーのコードを生成します。

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet.proto
```

うまくいけば、使えるコードが生成されているはずです。まずは、生成されたクライアントコードを`Driver`の中で使うことから始めましょう。

```go
package grpcserver

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string
}

func (d Driver) Greet(name string) (string, error) {
	//todo: we shouldn't redial every time we call greet, refactor out when we're green
	conn, err := grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := NewGreeterClient(conn)
	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}
```

クライアントができたので、サーバーを作成するために`main.go`を更新する必要があります。この時点では、テストを通すことだけを目指しており、コードの品質については気にしていないことを思い出してください。

```go
package main

import (
	"context"
	"log"
	"net"

	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	grpcserver.RegisterGreeterServer(s, &GreetServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: "fixme"}, nil
}
```

gRPCサーバーを作成するには、生成されたインターフェースを実装する必要があります。

```go
// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	Greet(context.Context, *GreetRequest) (*GreetReply, error)
	mustEmbedUnimplementedGreeterServer()
}
```

`main`関数は次のことを行います。

- ポートでリッスンする
- インターフェースを実装する`GreetServer`を作成し、`grpc.Server`とともに`grpcServer.RegisterGreeterServer`で登録する
- リスナーを使ってサーバーを動かす

メッセージに`fix-me`をハードコードする代わりに、`greetServer.Greet`の中でドメインコードを呼び出すのは大した手間ではありませんが、まずは受け入れテストを実行して、トランスポートレベルですべてが機能しているかを確認し、失敗するテストの出力を検証したいと思います。

```
greet.go:16: Expected values to be equal:
-fixme
\ No newline at end of file
+Hello, Mike
\ No newline at end of file
```

いいですね！テストの中で、ドライバーがgRPCサーバーに接続できていることがわかります。

では、`GreetServer`の中でドメインコードを呼び出しましょう。

```go
type GreetServer struct {
	grpcserver.UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *grpcserver.GreetRequest) (*grpcserver.GreetReply, error) {
	return &grpcserver.GreetReply{Message: interactions.Greet(request.Name)}, nil
}
```

ついにテストが通りました！gRPCの挨拶サーバーが私たちの望むように振る舞うことを証明する受け入れテストができました。

## リファクタリング♪

テストを通すためにいくつかの罪を犯しましたが、テストが通った今、リファクタリングのためのセーフティネットがあります。

### mainをシンプルにする

以前と同様に、`main`にあまり多くのコードを持たせたくありません。新しい`GreetServer`は`adapters/grpcserver`に置くべきなので、そこに移動できます。凝集度の観点から言えば、サービス定義を変更した場合、変更の「影響範囲（blast-radius）」をコードのその領域に限定したいのです。

### ドライバーで毎回再接続しない

テストは1つしかありませんが、仕様を拡張した場合（そうするつもりです）、RPC呼び出しのたびにDriverが再接続するのは理にかなっていません。

```go
package grpcserver

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Driver struct {
	Addr string

	connectionOnce sync.Once
	conn           *grpc.ClientConn
	client         GreeterClient
}

func (d *Driver) Greet(name string) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}

	greeting, err := client.Greet(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}

func (d *Driver) getClient() (GreeterClient, error) {
	var err error
	d.connectionOnce.Do(func() {
		d.conn, err = grpc.Dial(d.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		d.client = NewGreeterClient(d.conn)
	})
	return d.client, err
}
```

ここでは、[`sync.Once`](https://pkg.go.dev/sync#Once)を使って、`Driver`がサーバーへの接続を一度だけ作成しようとするようにする方法を示しています。

先に進む前に、プロジェクト構造の現在の状態を見てみましょう。

```
quii@Chriss-MacBook-Pro go-specs-greet % tree
.
├── Makefile
├── README.md
├── adapters
│   ├── docker.go
│   ├── grpcserver
│   │   ├── driver.go
│   │   ├── greet.pb.go
│   │   ├── greet.proto
│   │   ├── greet_grpc.pb.go
│   │   └── server.go
│   └── httpserver
│       ├── driver.go
│       └── handler.go
├── cmd
│   ├── grpcserver
│   │   ├── Dockerfile
│   │   ├── greeter_server_test.go
│   │   └── main.go
│   └── httpserver
│       ├── Dockerfile
│       ├── greeter_server_test.go
│       └── main.go
├── domain
│   └── interactions
│       ├── greet.go
│       └── greet_test.go
├── go.mod
├── go.sum
└── specifications
    └── greet.go
```

- `adapters`には、凝集した機能の単位がまとめられています
- `cmd`には、アプリケーションと対応する受け入れテストが置かれています
- コードは、あらゆる偶有的な複雑さから完全に切り離されています

### `Dockerfile`の統合

2つの`Dockerfile`が、ビルドしたいバイナリへのパス以外はほぼ同じであることに気づいたかもしれません。

`Dockerfile`は引数を受け取ることができ、異なるコンテキストで再利用できるようになります。これはまさにうってつけです。2つのDockerfileを削除し、代わりにプロジェクトのルートに次の内容のDockerfileを1つ置くことができます。

```dockerfile
# Make sure to specify the same Go version as the one in the go.mod file.
FROM golang:1.18-alpine

WORKDIR /app

ARG bin_to_build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/${bin_to_build}/main.go

CMD [ "./svr" ]
```

イメージをビルドするときに引数を渡すように、`StartDockerServer`関数を更新する必要があります。

```go
func StartDockerServer(
	t testing.TB,
	port string,
	binToBuild string,
) {
	ctx := context.Background()
	t.Helper()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../.",
			Dockerfile: "Dockerfile",
			BuildArgs: map[string]*string{
				"bin_to_build": &binToBuild,
			},
			PrintBuildLog: true,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})
}
```

最後に、ビルドするイメージを渡すようにテストを更新します（もう一方のテストでも同様にし、`grpcserver`を`httpserver`に変更してください）。

```go
func TestGreeterServer(t *testing.T) {
	var (
		port   = "50051"
		driver = grpcserver.Driver{Addr: fmt.Sprintf("localhost:%s", port)}
	)

	adapters.StartDockerServer(t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
}
```

### 異なる種類のテストを分離する

受け入れテストは、純粋にユーザーに面した振る舞いの観点からシステム全体が機能することをテストするという点で素晴らしいものですが、単体テストと比べると欠点もあります。

- 遅い
- フィードバックの質が、単体テストほど焦点の絞られたものにならないことが多い
- 内部品質や設計の助けにはならない

[テストピラミッド](https://martinfowler.com/articles/practical-test-pyramid.html)は、テストスイートに求める種類の組み合わせについての指針を与えてくれます。詳しくはFowlerの投稿を読むべきですが、この記事のためのごく単純な要約は「たくさんの単体テストと、少しの受け入れテスト」です。

そのため、プロジェクトが成長するにつれて、受け入れテストの実行に数分かかるような状況になることがよくあります。プロジェクトをチェックアウトする人々に親切な開発者体験を提供するために、開発者が異なる種類のテストを別々に実行できるようにしておくとよいでしょう。

`go test ./...`は、Goコンパイラ（当然）や、場合によってはDockerといった、いくつかの主要な依存関係以外には、エンジニアがそれ以上のセットアップをしなくても実行できることが望ましいです。

Goは、[shortフラグ](https://pkg.go.dev/testing#Short)を使って、エンジニアが「短い」テストだけを実行できる仕組みを提供しています。

`go test -short ./...`

受け入れテストに、このフラグの値を調べてユーザーが受け入れテストを実行したいかどうかを確認する処理を追加できます。

```go
if testing.Short() {
	t.Skip()
}
```

この使い方を示すために`Makefile`を作成しました。

```makefile
build:
	golangci-lint run
	go test ./...

unit-tests:
	go test -short ./...
```

### 受け入れテストはいつ書くべきか？

ベストプラクティスは、高速に実行される単体テストをたくさん持ち、受け入れテストは少しにすることですが、受け入れテストを書くべきか、単体テストを書くべきかをどのように判断すればよいのでしょうか？

具体的なルールを示すのは難しいですが、私が普段自問している質問は次のとおりです。

- これはエッジケースか？そうであれば、単体テストでテストしたい
- これはコンピュータに詳しくない人たちがよく話題にすることか？その重要なことが「本当に」機能しているという確信を強く持ちたいので、受け入れテストを追加する
- 特定の関数ではなく、ユーザージャーニーを記述しているか？そうであれば受け入れテスト
- 単体テストで十分な確信が得られるか？すでに受け入れテストがある既存のジャーニーに対して、異なる入力による異なるシナリオに対処するための機能を追加している場合があります。この場合、別の受け入れテストを追加するとコストがかかる割に得られる価値は小さいので、いくつかの単体テストを書くほうが好ましいです。

## 作業の反復

これだけの労力をかけたのですから、システムの拡張は簡単になっていてほしいものです。作業しやすいシステムを作ることは必ずしも簡単ではありませんが、時間をかける価値があり、プロジェクトの開始時に行えばずっと簡単です。

APIを拡張して「悪態（curse）」機能を追加しましょう。

## 最初にテストを書く

これはまったく新しい振る舞いなので、受け入れテストから始めるべきです。仕様ファイルに次の内容を追加します。

```go
type MeanGreeter interface {
	Curse(name string) (string, error)
}

func CurseSpecification(t *testing.T, meany MeanGreeter) {
	got, err := meany.Curse("Chris")
	assert.NoError(t, err)
	assert.Equal(t, got, "Go to hell, Chris!")
}
```

受け入れテストのどれか1つを選んで、この仕様を使ってみましょう。

```go
func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port   = "50051"
		driver = grpcserver.Driver{Addr: fmt.Sprintf("localhost:%s", port)}
	)

	t.Cleanup(driver.Close)
	adapters.StartDockerServer(t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
	specifications.CurseSpecification(t, &driver)
}
```

## テストを実行してみます

```
# github.com/quii/go-specs-greet/cmd/grpcserver_test [github.com/quii/go-specs-greet/cmd/grpcserver.test]
./greeter_server_test.go:27:39: cannot use &driver (value of type *grpcserver.Driver) as type specifications.MeanGreeter in argument to specifications.CurseSpecification:
	*grpcserver.Driver does not implement specifications.MeanGreeter (missing Curse method)
```

`Driver`はまだ`Curse`をサポートしていません。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

テストを実行できるようにすることだけを目指していることを思い出して、`Driver`にメソッドを追加します。

```go
func (d *Driver) Curse(name string) (string, error) {
	return "", nil
}
```

もう一度試すと、テストはコンパイルされ、実行され、そして失敗するはずです。

```
greet.go:26: Expected values to be equal:
+Go to hell, Chris!
\ No newline at end of file
```

## 成功させるのに十分なコードを書く

Protocol Bufferの仕様を更新して`Curse`メソッドを持たせ、コードを再生成する必要があります。

```protobuf
service Greeter {
  rpc Greet (GreetRequest) returns (GreetReply) {}
  rpc Curse (GreetRequest) returns (GreetReply) {}
}
```

`GreetRequest`と`GreetReply`の型を再利用するのは不適切な結合だと主張することもできますが、それはリファクタリングの段階で対処できます。何度も強調しているように、今はテストを通すことだけを目指しています。そうしてソフトウェアが機能することを検証し、 _それから_ きれいにすればよいのです。

（`adapters/grpcserver`の中で）次のコマンドでコードを再生成します。

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greet.proto
```

### ドライバーの更新

クライアントコードが更新されたので、`Driver`で`Curse`を呼び出せるようになりました。

```go
func (d *Driver) Curse(name string) (string, error) {
	client, err := d.getClient()
	if err != nil {
		return "", err
	}

	greeting, err := client.Curse(context.Background(), &GreetRequest{
		Name: name,
	})
	if err != nil {
		return "", err
	}

	return greeting.Message, nil
}
```

### サーバーの更新

最後に、`Server`に`Curse`メソッドを追加する必要があります。

```go
package grpcserver

import (
	"context"
	"fmt"

	"github.com/quii/go-specs-greet/domain/interactions"
)

type GreetServer struct {
	UnimplementedGreeterServer
}

func (g GreetServer) Curse(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: fmt.Sprintf("Go to hell, %s!", request.Name)}, nil
}

func (g GreetServer) Greet(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: interactions.Greet(request.Name)}, nil
}
```

これでテストが通るはずです。

## リファクタリング♪

これは自分でやってみてください。

- `Greet`で行ったように、`Curse`の「ドメインロジック」をgRPCサーバーから抽出しましょう。仕様を、ドメインロジックに対する単体テストとして使いましょう
- protobufで異なる型を持たせ、`Greet`と`Curse`のメッセージ型が切り離されるようにしましょう

## HTTPサーバーに`Curse`を実装する

これも読者であるあなたへの演習です。ドメインレベルの仕様とドメインレベルのロジックがきれいに分離されています。この章をここまで進めてきたのであれば、とても簡単なはずです。

- HTTPサーバーの既存の受け入れテストに仕様を追加する
- `Driver`を更新する
- サーバーに新しいエンドポイントを追加し、ドメインコードを再利用して機能を実装する。別々のエンドポイントへのルーティングを処理するために`http.NewServeMux`を使うとよいかもしれません。

小さなステップで作業し、頻繁にコミットしてテストを実行することを忘れないでください。本当に行き詰まったら、[GitHubで私の実装を見ることができます](https://github.com/quii/go-specs-greet)。

## 単体テストでドメインロジックを更新して、両方のシステムを強化する

前述のとおり、システムへのすべての変更を受け入れテストで駆動すべきではありません。関心事をうまく分離していれば、ビジネスルールの組み合わせやエッジケースは単体テストで簡単に駆動できるはずです。

`Greet`関数に、`name`が空の場合はデフォルトで`World`にする単体テストを追加しましょう。これがいかにシンプルかがわかるはずです。そして、そのビジネスルールは「タダで」両方のアプリケーションに反映されます。

## まとめ

妥当な変更コストでシステムを構築するには、メンテナンスの重荷になるのではなく、助けとなるように設計された受け入れテストが必要です。受け入れテストは、ソフトウェアを導く手段として、あるいはGOOSが言うように、ソフトウェアを系統立てて「育てる」手段として使うことができます。

この例を通して、変更を駆動するためのアプリケーションの予測可能で構造化されたワークフローと、それを自分の仕事にどのように活用できるかがわかっていただけたなら幸いです。

あなたが取り組んでいるシステムを何らかの形で拡張したいと考えているステークホルダーと話している場面を想像してみてください。それをドメイン中心で、実装に依存しない形で仕様に捉え、取り組みの北極星として使いましょう。Riyaと私は、[GopherconUKでの講演](https://www.youtube.com/watch?v=ZMWJCk_0WrY)で、「Example Mapping」のようなBDDのテクニックを活用して本質的な複雑さをより深く理解し、より詳細で意味のある仕様を書けるようにする方法について説明しています。

本質的な複雑さと偶有的な複雑さという関心事を分離することで、作業は場当たり的なものではなくなり、より構造化された慎重なものになります。これにより受け入れテストの回復力が確保され、メンテナンスの重荷になりにくくなります。

Dave Farleyは素晴らしいヒントを与えてくれています。

> 問題領域を理解している人の中で、思いつく限り最も技術に詳しくない人が、あなたの受け入れテストを読んでいるところを想像してください。テストはその人にとって意味の通じるものであるべきです。

そうすれば、仕様はドキュメントとしても機能するはずです。仕様は、システムがどのように振る舞うべきかを明確に規定すべきです。この考え方は[Cucumber](https://cucumber.io)のようなツールの根底にある原則です。Cucumberは振る舞いをコードとして捉えるためのDSLを提供し、私たちがここで行ったのと同じように、そのDSLをシステム呼び出しに変換します。

### 取り上げた内容

- 抽象的な仕様を書くことで、解決しようとしている問題の本質的な複雑さを表現し、偶有的な複雑さを取り除くことができます。これにより、仕様を異なるコンテキストで再利用できるようになります。
- [Testcontainers](https://golang.testcontainers.org)を使って、受け入れテストのためにシステムのライフサイクルを管理する方法。これにより、出荷する予定のイメージを自分のコンピュータで徹底的にテストでき、素早いフィードバックと確信が得られます。
- Dockerを使ったアプリケーションのコンテナ化の簡単な紹介
- gRPC
- 出来合いのフォルダ構造を追いかけるのではなく、自分自身のニーズに基づいて、開発のアプローチを通じてアプリケーションの構造を自然に導き出すことができます

### さらなる資料

- この例では、「DSL」はDSLと呼べるほどのものではありません。インターフェースを使って仕様を現実世界から切り離し、ドメインロジックをきれいに表現できるようにしただけです。システムが成長するにつれて、このレベルの抽象化では扱いにくく、わかりにくくなるかもしれません。仕様をどのように構造化するかについてもっとアイデアを得たい場合は、[「Screenplayパターン」について読んでみてください](https://cucumber.io/blog/bdd/understanding-screenplay-part-1/)。
- 強調しておくと、[Growing Object-Oriented Software, Guided by Tests](http://www.growing-object-oriented-software.com)は古典です。この「ロンドンスタイル」、「トップダウン」のアプローチをソフトウェアの作成に適用する方法を示しています。Learn Go with Testsを楽しんだ人なら誰でも、GOOSを読むことで多くの価値を得られるはずです。
- [サンプルコードのリポジトリ](https://github.com/quii/go-specs-greet)には、マルチステージのDockerビルドなど、ここでは書いていないコードやアイデアがさらにあるので、確認してみるとよいでしょう。
  - 特に、*お遊びで*、`Greet`と`Curse`のためのHTMLフォームを備えたWebサイトという**3つ目のプログラム**を作りました。その`Driver`は、見るからに優れた[https://github.com/go-rod/rod](https://github.com/go-rod/rod)モジュールを活用しており、ユーザーと同じようにブラウザでWebサイトを操作できます。Gitの履歴を見ると、「とにかく動かすために」テンプレートツールを使わずに始めた様子がわかります。そして、受け入れテストが通ってからは、何かを壊す心配をせずに自由にテンプレートツールを導入できました。
