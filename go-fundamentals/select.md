---
description: Select
---

# Select

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/select)**

あなたは`WebsiteRacer`という関数を作るように頼まれました。この関数は2つのURLを受け取り、それぞれにHTTP GETでアクセスして「競争」させ、先に応答が返ってきたURLを返します。どちらも10秒以内に応答しない場合は`error`を返します。

これには、次のものを使用します。

* HTTP呼び出しを行うための `net/http`。
* それらのテストを支援するための`net/http/httptest`。
* ゴルーチン。
* プロセスを同期するための`select`。

## 最初にテストを書く

まずは、単純なことから始めましょう。

```go
func TestRacer(t *testing.T) {
	slowURL := "http://www.facebook.com"
	fastURL := "http://www.quii.dev"

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

これは完璧ではなく、問題があることはわかっていますが、手始めとしては十分です。最初から完璧を目指すことにこだわりすぎないことが重要です。

## テストを実行してみます

`./racer_test.go:14:9: undefined: Racer`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func Racer(a, b string) (winner string) {
	return
}
```

`racer_test.go:25: got '', want 'http://www.quii.dev'`

## 成功させるのに十分なコードを書く

```go
func Racer(a, b string) (winner string) {
	startA := time.Now()
	resp, err := http.Get(a)
	if err == nil {
		resp.Body.Close()
	}
	aDuration := time.Since(startA)

	startB := time.Now()
	resp, err = http.Get(b)
	if err == nil {
		resp.Body.Close()
	}
	bDuration := time.Since(startB)

	if aDuration < bDuration {
		return a
	}

	return b
}
```

各URLについて

1. `time.Now()`を使用して、 `URL`を取得しようとする直前に記録します。
2. 次に、[`http.Get`](https://golang.org/pkg/net/http/#Client.Get)を使用して、`URL`に対してHTTP `GET`リクエストを実行します。この関数は[`http.Response`](https://golang.org/pkg/net/http/#Response)と`error`を返します。ファイルディスクリプタのリークを避けるため、レスポンスボディを閉じます。閉じないと、リクエストごとに接続が開いたまま残ってしまいます。
3. `time.Since`は開始時間を取り、差の`time.Duration`を返します。

これができたら、あとは所要時間を比較してどちらが速いかを確認するだけです。

### 問題

これでテストに合格するかどうかは環境次第です。問題は、自分たちのロジックをテストするために実際のWebサイトにアクセスしていることです。

HTTPを使用するコードのテストは非常に一般的であるため、Goの標準ライブラリには、テストに役立つツールがあります。

モックと依存性注入の章では、コードをテストするために外部サービスに依存するのは理想的ではないことを説明しました。外部サービスには次のような問題があるからです。

* 遅い
* 不安定（Flaky）
* エッジケースをテストできない

標準ライブラリには、[`net/http/httptest`](https://golang.org/pkg/net/http/httptest/)というパッケージがあり、これを使うとモックのHTTPサーバーを簡単に作成できます。

テストを変更してモックを使うようにし、自分たちで制御できる信頼性の高いサーバーに対してテストできるようにしましょう。

```go
func TestRacer(t *testing.T) {

	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	slowServer.Close()
	fastServer.Close()
}
```

構文が少しごちゃごちゃして見えるかもしれませんが、焦らずに読んでください。

`httptest.NewServer`は`http.HandlerFunc`を受け取ります。ここではそれを _無名関数_ として渡しています。

`http.HandlerFunc`は、`type HandlerFunc func(ResponseWriter, *Request)`のような型です。

これが言っているのは、`ResponseWriter`と`Request`を受け取る関数が必要だということだけです。HTTPサーバーとしては、それほど驚くことではないでしょう。

実のところ、ここには特別な魔法はありません。**これは、Goで _本物の_ HTTPサーバーを書くときの方法でもあります**。唯一の違いは、それを`httptest.NewServer`でラップしていることです。`httptest.NewServer`は待ち受け用の空いているポートを見つけてくれて、テストが終わったら閉じることができるので、テストで使いやすくなります。

2つのサーバーのうち遅い方には、リクエストを受け取ったときに短い`time.Sleep`を入れて、もう一方より遅くなるようにしています。そして両方のサーバーとも、`w.WriteHeader(http.StatusOK)`で`OK`レスポンスを呼び出し元に返します。

テストを再実行すると、今度は確実に成功し、実行も速くなるはずです。これらのスリープ時間をいじって、意図的にテストを失敗させてみてください。


## リファクタリング♪

プロダクションコードとテストコードの両方に重複があります。

```go
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	resp, err := http.Get(url)
	if err == nil {
		resp.Body.Close()
	}
	return time.Since(start)
}
```

このように重複を取り除く（DRYにする）ことで、`Racer`のコードがずっと読みやすくなります。

```go
func TestRacer(t *testing.T) {

	slowServer := makeDelayedServer(20 * time.Millisecond)
	fastServer := makeDelayedServer(0 * time.Millisecond)

	defer slowServer.Close()
	defer fastServer.Close()

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
```

偽のサーバーの作成を`makeDelayedServer`という関数にリファクタリングし、本質的でないコードをテストの外に移して繰り返しを減らしました。

### `defer`

関数呼び出しの前に`defer`を付けることで、その関数を _それを含む関数の最後に_ 呼び出します。

ファイルを閉じたり、今回のようにサーバーがポートで待ち受け続けないようにサーバーを閉じたりと、リソースをクリーンアップする必要がある場合があります。

これは関数の最後に実行したいのですが、将来コードを読む人のために、その命令はサーバーを作成した場所の近くに置いておきたいわけです。

このリファクタリングは改善であり、これまでに取り上げたGoの機能からすれば妥当な解決策ですが、もっとシンプルにすることができます。

### プロセスの同期

* Goは並行処理が得意なのに、なぜWebサイトの速度を1つずつ順番にテストしているのでしょうか？両方を同時にチェックできるはずです。
* リクエストの _正確な応答時間_ はあまり気にしていません。どちらが先に返ってくるかを知りたいだけです。

そのために、プロセスをとても簡単かつ明確に同期できる`select`という新しい構文を導入します。

```go
func Racer(a, b string) (winner string) {
	select {
	case <-ping(a):
		return a
	case <-ping(b):
		return b
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
		}
		close(ch)
	}()
	return ch
}
```

#### `ping`

`chan struct{}`を作成して返す関数`ping`を定義しました。

今回のケースでは、チャネルにどんな型が送られるかは _気にしません_ 。 _完了したことを通知したいだけ_ なので、チャネルを閉じるだけで完璧に用が足ります！

なぜ`bool`のような別の型ではなく`struct{}`なのでしょうか？`chan struct{}`はメモリの観点から利用できる最小のデータ型なので、`bool`と違ってメモリ割り当てが発生しません。チャネルを閉じるだけで何も送信しないのに、わざわざ何かを割り当てる必要はないでしょう？

同じ関数内で、`http.Get(url)`を完了すると、そのチャネルに信号を送信するゴルーチン（`goroutine`）を開始します。レスポンスボディはすぐに閉じます。私たちが気にしているのはリクエストが完了したことだけで、レスポンスの内容ではありません。また、開いたままにするとファイルディスクリプタがリークしてしまいます。

##### 常にチャネルを`make`する

チャネルを作成するときに、`var ch chan struct{}`のように書くのではなく、`make`を使わなければならないことに注意してください。`var`を使用すると、変数は型の「ゼロ」値で初期化されます。したがって、`string`の場合は`""`、`int`の場合は`0`になります。

チャネルの場合、ゼロ値は`nil`であり、`<-`で送信しようとすると、`nil`チャネルに送信できないため、永久にブロックされます。

[これは Go Playground で実際に見ることができます](https://play.golang.org/p/IIbeAox5jKA)

#### `select`

並行性の章で学んだように、`myVar := <-ch`を使用して値がチャネルに送信されるのを待つことができます。値を待っているので、これは _ブロッキング_ 呼び出しです。

`select`を使うと、 _複数の_ チャネルで待機できます。最初に値を送ってきたものが「勝ち」となり、その`case`の下のコードが実行されます。

`select`の中で`ping`を使い、`URL`ごとに1つずつ、2つのチャネルを用意します。先に自分のチャネルに書き込んだ方の`case`のコードが`select`で実行され、その結果、その`URL`が返されます（つまり勝者となります）。

これらの変更によって、コードの意図が非常に明確になり、実装も実際にはより単純になりました。

### タイムアウト

最後の要件は、`Racer`に10秒以上かかる場合にエラーを返すことでした。

## 最初にテストを書く

```go
func TestRacer(t *testing.T) {
	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, _ := Racer(slowURL, fastURL)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		serverA := makeDelayedServer(11 * time.Second)
		serverB := makeDelayedServer(12 * time.Second)

		defer serverA.Close()
		defer serverB.Close()

		_, err := Racer(serverA.URL, serverB.URL)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}
```

このシナリオを試すために、テストサーバーが応答を返すまでに10秒以上かかるようにしました。また、`Racer`が2つの値、つまり勝者のURL（このテストでは`_`で無視しています）と`error`を返すことを期待するようにしています。

元のテストでもエラーの戻り値を扱うようにしていることに注意してください。今のところは、テストが実行できるように`_`を使っています。

## テストを実行してみます

`./racer_test.go:37:10: assignment mismatch: 2 variables but Racer returns 1 value`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	}
}
```

勝者と`error`を返すように`Racer`のシグネチャを変更します。正常系のケースでは`nil`を返します。

今実行すると、11秒後に失敗します。

```
--- FAIL: TestRacer (12.00s)
    --- FAIL: TestRacer/returns_an_error_if_a_server_doesn't_respond_within_10s (12.00s)
        racer_test.go:40: expected an error but didn't get one
```

## 成功させるのに十分なコードを書く

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

`time.After`は、`select`を使用する場合に非常に便利な関数です。今回のケースでは起きませんでしたが、待ち受けているチャネルが値を返さないと、永久にブロックするコードを書いてしまう可能性があります。`time.After`は（`ping`のように）`chan`を返し、指定した時間が経過するとそこにシグナルを送ります。

これは今回の用途にぴったりです。`a`か`b`が先に応答すればそれが勝ちますが、10秒が経過すると`time.After`がシグナルを送り、`error`を返します。

### 遅いテスト

問題は、このテストの実行に10秒かかることです。これほど単純なロジックにしては、あまり良いとは言えません。

そこで、タイムアウトを設定可能にします。そうすれば、テストではとても短いタイムアウトにして、実際にコードを使うときには10秒に設定できます。

```go
func Racer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

タイムアウトを指定していないため、テストはコンパイルされません。

急いでこのデフォルト値を両方のテストに追加する前に、 _テストの声に耳を傾けて_ みましょう。

* 「ハッピー」なテストでタイムアウトを気にする必要があるでしょうか？
* タイムアウトについては要件で明示されていました

これを踏まえて、テストとコードの利用者の両方に配慮した形に少しリファクタリングしてみましょう。

```go
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

利用者と最初のテストは`Racer`（内部では`ConfigurableRacer`を使っています）を使い、異常系（サッドパス）のテストでは`ConfigurableRacer`を使えます。

```go
func TestRacer(t *testing.T) {

	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within the specified time", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}
```

最初のテストに最後のチェックを1つ追加して、`error`が発生しないことを確認しました。

## まとめ

### `select`

* 複数のチャネルで待機するのに役立ちます。
* システムが永久にブロックされるのを防ぐために、`case`の1つに`time.After`を含めたい場合もあるでしょう。

### `httptest`

* テストサーバーを作成する便利な方法で、信頼性が高く制御可能なテストを書けるようになります。
* 「実際の」`net/http`サーバーと同じインターフェースを使用します。これは一貫性があり、学ぶことが少なくて済みます。
