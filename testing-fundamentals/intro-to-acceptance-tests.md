---
description: 受け入れテスト入門
---

# 受け入れテスト入門

`$WORK`（私の職場）では、サービスに「グレースフルシャットダウン」が必要になる場面に出くわしてきました。グレースフルシャットダウンとは、システムが終了させられる前に、処理中の作業をきちんと終わらせることです。現実世界でたとえるなら、次の会議に移る前に、話の途中でいきなり電話を切るのではなく、きちんと通話を締めくくろうとするようなものです。

この章では、HTTPサーバーにおけるグレースフルシャットダウンの入門と、コードの振る舞いに自信を持つための「受け入れテスト」の書き方を紹介します。

この章を読み終えるころには、優れたテストを備えたパッケージを共有する方法、保守の手間を減らす方法、そして自分の仕事の品質への自信を高める方法がわかるようになります。

## Kubernetesについて最低限知っておくこと

私たちはソフトウェアを[Kubernetes](https://kubernetes.io/)（K8s）上で動かしています。K8sはさまざまな理由で「ポッド」（実際には私たちのソフトウェア）を終了させます。よくある理由の1つは、デプロイしたい新しいコードをプッシュしたときです。

私たちは[DORAメトリクス](https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance)について高い基準を設けているので、小さく段階的な改善や機能を1日に何度も本番環境にデプロイするという働き方をしています。

k8sがポッドを終了させたいとき、k8sは[「終了ライフサイクル」](https://cloud.google.com/blog/products/containers-kubernetes/kubernetes-best-practices-terminating-with-grace)を開始します。その一環として、私たちのソフトウェアにSIGTERMシグナルを送ります。これはk8sが私たちのコードにこう伝えているのです。

> 自分自身をシャットダウンしなさい。やっている作業は何でも終わらせなさい。一定の「猶予期間」が過ぎたら`SIGKILL`を送るから、そうしたら一巻の終わりだよ。

`SIGKILL`を受け取ると、プログラムが行っていた作業はすべて即座に停止されます。

## 猶予がない場合

ソフトウェアの性質によっては、`SIGTERM`を無視すると問題が起きることがあります。

私たちが抱えていた具体的な問題は、処理中（in-flight）のHTTPリクエストでした。自動テストがAPIを実行しているときに、k8sがポッドを停止することにすると、サーバーが落ち、テストはサーバーからレスポンスを受け取れず、テストが失敗してしまうのです。

これによってインシデントチャンネルでアラートが発火し、開発者は今やっていることを中断して問題に対処しなければなりません。こうした断続的な失敗は、チームにとって煩わしい邪魔になります。

こうした問題は私たちのテストに限ったことではありません。ユーザーがシステムにリクエストを送り、その処理中にプロセスが終了させられると、ユーザーはおそらく5xx系のHTTPエラーに迎えられることになります。これは提供したいユーザー体験ではありません。

## 猶予がある場合

私たちがやりたいのは、`SIGTERM`を待ち受け、サーバーを即座に終了させるのではなく、次のようにすることです。

- それ以上のリクエストの受け付けを停止する
- 処理中のリクエストを完了させる
- *その後で*プロセスを終了する

## 猶予を持つ方法

ありがたいことに、Goには[net/http/Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)というサーバーをグレースフルにシャットダウンするための仕組みがすでにあります。

> Shutdownは、アクティブな接続を中断することなく、サーバーをグレースフルにシャットダウンします。Shutdownは、まず開いているすべてのリスナーを閉じ、次にすべてのアイドル状態の接続を閉じ、その後、接続がアイドル状態に戻るのを無期限に待ってからシャットダウンします。シャットダウンが完了する前に提供されたコンテキストの期限が切れた場合、Shutdownはコンテキストのエラーを返します。そうでない場合は、サーバーの基になるリスナーを閉じる際に返されたエラーを返します。

`SIGTERM`を処理するには、[os/signal.Notify](https://pkg.go.dev/os/signal#Notify)を使うことができます。これは受信したシグナルを、私たちが渡したチャネルに送ってくれます。

標準ライブラリのこの2つの機能を使えば、`SIGTERM`を待ち受けてグレースフルにシャットダウンできます。

## グレースフルシャットダウンのパッケージ

そのために、私は[https://pkg.go.dev/github.com/quii/go-graceful-shutdown](https://pkg.go.dev/github.com/quii/go-graceful-shutdown)を書きました。これは`*http.Server`のためのデコレーター関数を提供し、`SIGTERM`シグナルが検出されたときに`Shutdown`メソッドを呼び出します。

```go
func main() {
	var (
		ctx        = context.Background()
		httpServer = &http.Server{Addr: ":8080", Handler: http.HandlerFunc(acceptancetests.SlowHandler)}
		server     = gracefulshutdown.NewServer(httpServer)
	)

	if err := server.ListenAndServe(ctx); err != nil {
		// this will typically happen if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didn't shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
```

コードの詳細はこの章を読むうえではそれほど重要ではありませんが、先に進む前にざっと目を通しておく価値はあります。

## テストとフィードバックループ

`gracefulshutdown`パッケージを書いたとき、私たちはそれが正しく振る舞うことを証明するユニットテストを持っていたので、積極的にリファクタリングする自信がありました。しかし、それが**本当に**動くという「自信」はまだ感じられませんでした。

そこで`cmd`パッケージを追加し、書いているパッケージを使う実際のプログラムを作りました。手動でそれを起動し、HTTPリクエストを送り、そして`SIGTERM`を送って何が起きるかを確認していたのです。

**エンジニアであるあなたは、手動テストに居心地の悪さを感じるべきです**。
手動テストは退屈で、スケールせず、不正確で、無駄が多いのです。共有するつもりのパッケージを書いていて、しかもそれをシンプルで安価に変更できる状態に保ちたいなら、手動テストでは不十分です。

## 受け入れテスト

この本の他の章を読んできたなら、これまでに書いてきたのはほとんどが「ユニットテスト」でしょう。ユニットテストは、恐れることなくリファクタリングできるようにし、優れたモジュール設計を促し、リグレッションを防ぎ、素早いフィードバックを可能にする素晴らしいツールです。

その性質上、ユニットテストはシステムの小さな部分しかテストしません。通常、ユニットテストだけでは効果的なテスト戦略として*十分ではありません*。私たちはシステムを**常に出荷可能な状態**にしておきたいことを思い出してください。手動テストに頼ることはできないので、別の種類のテストが必要です。それが**受け入れテスト**です。

### 受け入れテストとは何か

受け入れテストは「ブラックボックステスト」の一種です。「機能テスト」と呼ばれることもあります。受け入れテストは、システムのユーザーと同じようにシステムを実行するべきです。

「ブラックボックス」という用語は、テストコードがシステムの内部にアクセスできず、公開されたインターフェースだけを使い、観察した振る舞いに対してアサーションを行うことしかできない、という考え方を指しています。つまり、システム全体としてしかテストできないということです。

これは有利な特性です。テストがユーザーと同じようにシステムを実行するので、テストを通すことはできても証明すべきことを実際には証明していない、というような特別な回避策を使えないからです。これは、ユニットテストのファイルを別のテストパッケージ、たとえば`package mypkg`ではなく`package mypkg_test`に置くことを好むという原則と似ています。

### 受け入れテストの利点

- テストが成功すれば、システム全体が望みどおりに振る舞うことがわかります。
- 手動テストよりも正確で、速く、労力も少なくて済みます。
- うまく書かれていれば、システムの正確で検証済みのドキュメントとして機能します。システムの実際の振る舞いと乖離していくドキュメントという罠に陥りません。
- モックは不要です！すべて本物です。

### ユニットテストと比べたときの潜在的な欠点

- 書くのにコストがかかります。
- 実行に時間がかかります。
- システムの設計に依存します。
- 失敗したとき、通常は根本原因を教えてくれず、デバッグが難しいことがあります。
- システムの内部品質についてのフィードバックを与えてくれません。まったくひどいコードを書いても、受け入れテストを通すことはできてしまいます。
- ブラックボックスという性質上、すべてのシナリオを現実的に実行できるわけではありません。

このため、受け入れテストだけに頼るのは愚かなことです。受け入れテストはユニットテストが持つ多くの性質を持っておらず、大量の受け入れテストを抱えたシステムは、保守コストやリードタイムの悪化に苦しむ傾向があります。

#### リードタイムとは？

リードタイムとは、コミットがメインブランチにマージされてから本番環境にデプロイされるまでにかかる時間のことです。この数値は、チームによって数週間、さらには数か月から、数分までさまざまです。繰り返しになりますが、`$WORK`ではDORAの知見を重視しており、リードタイムを10分未満に保ちたいと考えています。

信頼性が高く、優れたリードタイムを持つシステムには、バランスの取れたテストアプローチが必要です。これは通常、[テストピラミッド](https://martinfowler.com/articles/practical-test-pyramid.html)という観点で説明されます。

## 基本的な受け入れテストの書き方

これは元の問題とどう関係するのでしょうか？ここで書いたのはパッケージであり、それは完全にユニットテスト可能です。

先ほど述べたように、ユニットテストだけでは必要な自信が十分に得られませんでした。実際に動いているプログラムと統合したときにパッケージが動作することを、*本当に*確かめたいのです。手動で行っていた確認を自動化できるはずです。

テスト用のプログラムを見てみましょう。

```go
func main() {
	var (
		ctx        = context.Background()
		httpServer = &http.Server{Addr: ":8080", Handler: http.HandlerFunc(acceptancetests.SlowHandler)}
		server     = gracefulshutdown.NewServer(httpServer)
	)

	if err := server.ListenAndServe(ctx); err != nil {
		// this will typically happen if our responses aren't written before the ctx deadline, not much can be done
		log.Fatalf("uh oh, didn't shutdown gracefully, some responses may have been lost %v", err)
	}

	// hopefully, you'll always see this instead
	log.Println("shutdown gracefully! all responses were sent")
}
```

お察しのとおり、`SlowHandler`はレスポンスを遅らせるために`time.Sleep`を持っています。そのおかげで、`SIGTERM`を送って何が起きるかを確認する時間がありました。残りはほぼ定型的なコードです。

- `net/http/Server`を作成する
- それをライブラリでラップする（参照: [Decoratorパターン](https://en.wikipedia.org/wiki/Decorator_pattern)）
- ラップしたものを使って`ListenAndServe`する

### 受け入れテストの大まかな手順

- プログラムをビルドする
- それを実行する（そして`8080`でリッスンするのを待つ）
- サーバーにHTTPリクエストを送る
- サーバーがHTTPレスポンスを送る前に、`SIGTERM`を送る
- それでもレスポンスが得られるかを確認する

### プログラムのビルドと実行

```go
package acceptancetests

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	baseBinName = "temp-testbinary"
)

func LaunchTestProgram(port string) (cleanup func(), sendInterrupt func() error, err error) {
	binName, err := buildBinary()
	if err != nil {
		return nil, nil, err
	}

	sendInterrupt, kill, err := runServer(binName, port)

	cleanup = func() {
		if kill != nil {
			kill()
		}
		os.Remove(binName)
	}

	if err != nil {
		cleanup() // even though it's not listening correctly, the program could still be running
		return nil, nil, err
	}

	return cleanup, sendInterrupt, nil
}

func buildBinary() (string, error) {
	binName := randomString(10) + "-" + baseBinName

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		return "", fmt.Errorf("cannot build tool %s: %s", binName, err)
	}
	return binName, nil
}

func runServer(binName string, port string) (sendInterrupt func() error, kill func(), err error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	cmdPath := filepath.Join(dir, binName)

	cmd := exec.Command(cmdPath)

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("cannot run temp converter: %s", err)
	}

	kill = func() {
		_ = cmd.Process.Kill()
	}

	sendInterrupt = func() error {
		return cmd.Process.Signal(syscall.SIGTERM)
	}

	err = waitForServerListening(port)

	return
}

func waitForServerListening(port string) error {
	for i := 0; i < 30; i++ {
		conn, _ := net.Dial("tcp", net.JoinHostPort("localhost", port))
		if conn != nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("nothing seems to be listening on localhost:%s", port)
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
```

`LaunchTestProgram`は次のことを担当します。
- プログラムをビルドする
- プログラムを起動する
- ポート`8080`でリッスンするのを待つ
- テスト終了時にクリーンな状態が残るよう、プログラムを終了させて削除する`cleanup`関数を提供する
- 振る舞いをテストできるよう、プログラムに`SIGTERM`を送る`interrupt`関数を提供する

正直なところ、これは世界一きれいなコードではありません。しかし、エクスポートされた関数`LaunchTestProgram`にだけ注目してください。そこから呼び出されているエクスポートされていない関数は、面白みのない定型コードです。

先に述べたように、受け入れテストはセットアップが厄介になりがちです。このコードのおかげで*テスト*コードはかなり読みやすくなります。そして受け入れテストでは、こうした儀式的なコードを一度書いてしまえばそれで終わりで、あとは忘れてしまえることがよくあります。

### 受け入れテスト

私たちは、グレースフルシャットダウンありのプログラムとなしのプログラムの2つに対して、2つの受け入れテストを用意したいと考えました。そうすれば、私たちも読者も振る舞いの違いを確認できます。プログラムをビルドして実行する`LaunchTestProgram`があるので、両方の受け入れテストを書くのはとても簡単で、いくつかのヘルパー関数を再利用する恩恵も受けられます。

これはグレースフルシャットダウン*あり*のサーバーのテストです。[なしのテストはGitHubで見られます](https://github.com/quii/go-graceful-shutdown/blob/main/acceptancetests/withoutgracefulshutdown/main_test.go)。

```go
package main

import (
	"testing"
	"time"

	"github.com/quii/go-graceful-shutdown/acceptancetests"
	"github.com/quii/go-graceful-shutdown/assert"
)

const (
	port = "8080"
	url  = "<http://localhost:" + port
)

func TestGracefulShutdown(t *testing.T) {
	cleanup, sendInterrupt, err := acceptancetests.LaunchTestProgram(port)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)

	// just check the server works before we shut things down
	assert.CanGet(t, url)

	// fire off a request, and before it has a chance to respond send SIGTERM.
	time.AfterFunc(50*time.Millisecond, func() {
		assert.NoError(t, sendInterrupt())
	})
	// Without graceful shutdown, this would fail
	assert.CanGet(t, url)

	// after interrupt, the server should be shutdown, and no more requests will work
	assert.CantGet(t, url)
}
```

セットアップがカプセル化されているので、テストは網羅的で、振る舞いを説明しており、比較的追いやすくなっています。

`assert.CanGet/CantGet`は、このテストスイートでよく使うアサーションをDRYにするために私が作ったヘルパー関数です。

```go
func CanGet(t testing.TB, url string) {
	errChan := make(chan error)

	go func() {
		res, err := http.Get(url)
		if err != nil {
			errChan <- err
			return
		}
		res.Body.Close()
		errChan <- nil
	}()

	select {
	case err := <-errChan:
		NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Errorf("timed out waiting for request to %q", url)
	}
}
```

これはゴルーチン上で`URL`に`GET`を送り、3秒以内にエラーなしで応答すれば失敗しません。`CantGet`は簡潔さのために省略しますが、[GitHubで見ることができます](https://github.com/quii/go-graceful-shutdown/blob/main/assert/assert.go#L61)。

改めて強調しておきますが、Goには受け入れテストを書くために必要なツールが最初からすべて揃っています。受け入れテストを作るのに特別なフレームワークが*必要*なわけではありません。

### 小さな投資で大きな見返り

これらのテストがあれば、読者はサンプルプログラムを見て、そのサンプルが*実際に*動くことを確信できるので、パッケージの主張を信頼できます。

重要なのは、作者である私たちが、パッケージが現実の環境で動作するという**素早いフィードバック**と**大きな自信**を得られることです。

```shell
go test -count=1 ./...
ok  	github.com/quii/go-graceful-shutdown	0.196s
?   	github.com/quii/go-graceful-shutdown/acceptancetests	[no test files]
ok  	github.com/quii/go-graceful-shutdown/acceptancetests/withgracefulshutdown	4.785s
ok  	github.com/quii/go-graceful-shutdown/acceptancetests/withoutgracefulshutdown	2.914s
?   	github.com/quii/go-graceful-shutdown/assert	[no test files]
```

## まとめ

この章では、あなたのテストの道具箱に受け入れテストを加えました。受け入れテストは実際のシステムを構築し始めたときに非常に価値があり、ユニットテストを補完する重要なものです。

受け入れテストを*どのように*書くかは構築しているシステムによって異なりますが、原則は同じです。システムを「ブラックボックス」として扱うことです。Webサイトを作っているなら、テストはユーザーのように振る舞うべきなので、[Selenium](https://www.selenium.dev/)のようなヘッドレスWebブラウザを使って、リンクをクリックしたり、フォームに入力したりしたくなるでしょう。RESTful APIであれば、クライアントを使ってHTTPリクエストを送ることになります。

### より複雑なシステムに向けてさらに進める

複雑なシステムは、ここで扱ったような単一プロセスのアプリケーションであることはあまりありません。通常は、データベースなどの他のシステムに依存しています。こうしたシナリオでは、テストに使うローカル環境を自動化する必要があります。[docker-compose](https://docs.docker.com/compose/)のようなツールは、システムをローカルで動かすのに必要な環境のコンテナを立ち上げるのに役立ちます。

### 次の章

この章では、受け入れテストを後から書きました。しかし、[Growing Object-Oriented Software](http://www.growing-object-oriented-software.com)（邦題『実践テスト駆動開発』）の著者たちは、受け入れテストをテスト駆動のアプローチで使い、私たちの取り組みを導く「北極星」として機能させられることを示しています。

システムが複雑になるにつれて、受け入れテストを書いて保守するコストはあっという間に手に負えなくなることがあります。高コストな受け入れテストスイートのせいで身動きが取れなくなった開発チームの話は数え切れないほどあります。

次の章では、受け入れテストを使って設計を導く方法と、受け入れテストのコストを管理するための原則とテクニックを紹介します。

### オープンソースの品質を高める

共有するつもりのパッケージを書いているなら、パッケージが何をするのかを示すシンプルなサンプルプログラムを作り、わかりやすい受け入れテストに時間を投資することをお勧めします。そうすることで、あなた自身にも、あなたの成果物を使う可能性のあるユーザーにも自信を与えられます。

[Testable Examples](https://go.dev/blog/examples)と同じように、開発者体験に対するこうしたちょっとした追加の努力が見えることは、あなたの仕事への信頼を築くうえで大いに役立ち、あなた自身の保守コストも減らしてくれます。

## `$WORK`の採用宣伝

他のエンジニアと一緒に興味深い問題を解決する環境で働きたい方、ロンドンやポルトの近郊に住んでいる方、そしてこの章やこの本の内容を楽しんでいただけた方は、ぜひ[Twitterで私に連絡してください](https://twitter.com/quii)。近いうちに一緒に働けるかもしれません！
