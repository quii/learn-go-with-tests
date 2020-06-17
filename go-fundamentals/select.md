---
description: Select
---

# 選択

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/select)

2つのURLを取得し、それらをHTTP GETでヒットして最初に返されたURLを返すことで「競合」する`WebsiteRacer`と呼ばれる関数を作成するように求められました。 10秒以内に戻らない場合は、「エラー（`error`）」を返します。

これには

* HTTP呼び出しを行うための `net/http`。
* `net/http/httptest`は、それらをテストするのに役立ちます。
* ゴルーチン。
* プロセスを同期するための`select`。

## 最初にテストを書く

まずは、単純なことから始めましょう。

```go
func TestRacer(t *testing.T) {
    slowURL := "http://www.facebook.com"
    fastURL := "http://www.quii.co.uk"

    want := fastURL
    got := Racer(slowURL, fastURL)

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

これは完璧ではなく、問題があることはわかっていますが、うまくいくでしょう。
物事を最初から完璧にするのにあまり夢中にならないようにすることが重要です。

## テストを実行してみます

`./racer_test.go:14:9: undefined: Racer`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func Racer(a, b string) (winner string) {
    return
}
```

`racer_test.go:25: got '', want 'http://www.quii.co.uk'`

## 成功させるのに十分なコードを書く

```go
func Racer(a, b string) (winner string) {
    startA := time.Now()
    http.Get(a)
    aDuration := time.Since(startA)

    startB := time.Now()
    http.Get(b)
    bDuration := time.Since(startB)

    if aDuration < bDuration {
        return a
    }

    return b
}
```

各URLについて

1. `time.Now()`を使用して、 `URL`を取得しようとする直前に記録します。
2. 次に、[`http.Get`](https://golang.org/pkg/net/http/#Client.Get)を使用して、`URL`のコンテンツを取得します。この関数は[`http.Response`](https://golang.org/pkg/net/http/#Response)と`error`を返しますが、今のところこれらの値には興味がありません。
3. `time.Since`は開始時間を取り、差の`time.Duration`を返します。

これを実行したら、期間を比較してどちらが最も速いかを確認します。

### 問題

これにより、テストに合格する場合と合格しない場合があります。
問題は、実際のWebサイトに連絡して、独自のロジックをテストしていることです。

HTTPを使用するコードのテストは非常に一般的であるため、Goの標準ライブラリには、テストに役立つツールがあります。

モックと依存性注入の章では、コードをテストするために外部サービスに依存したくないという理想的な方法について説明しました。

* スロー（Slow）
* フレーク状（Flaky）
* エッジケースをテストできません（Can't test edge cases）

標準ライブラリには、[`net/http/httptest`](https://golang.org/pkg/net/http/httptest/)というパッケージがあり、模擬HTTPサーバーを簡単に作成できます。

テストをモックを使用するように変更して、制御できる信頼性の高いサーバーをテストできるようにします。

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

構文は少しせわしなく見えるかもしれませんが、時間をかけてください。

`httptest.NewServer`は、_anonymous function_ を介して送信する`http.HandlerFunc`を受け取ります。

`http.HandlerFunc`は、`type HandlerFunc func(ResponseWriter, *Request)`のようなタイプです。

実際に言っているのは、`ResponseWriter`と`Request`を受け取る関数が必要なことだけです。
これは、HTTPサーバーにとってそれほど驚くべきことではありません。

ここには特別な魔法はありません。**これは、Goで** _ **実際に** _ **HTTPサーバーを作成する方法でもあります**。唯一の違いは、それを`httptest.NewServer`でラップすることです。これにより、リッスンする開いているポートが見つかり、テストが完了したら閉じることができるため、テストでの使用が簡単になります。

2つのサーバー内では、遅いサーバーに他のサーバーよりも遅いリクエストを受け取ったときに、遅い方の`time.Sleep`を作成します。

次に、両方のサーバーが、`w.WriteHeader(http.StatusOK)`を使用して`OK`応答を呼び出し元に返します。

テストを再実行すると、テストは確実に成功し、より高速になるはずです。
これらの睡眠を試して、意図的にテストを中断します。


## リファクタリング

製品コードとテストコードの両方に重複があります。

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
    http.Get(url)
    return time.Since(start)
}
```

このドライアップ（DRY-ing up）により、`Racer`コードが非常に読みやすくなります。

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

偽のサーバーの作成を`makeDelayedServer`という関数にリファクタリングし、興味のないコードをテストから除外して繰り返しを減らしました。

### `defer`

関数呼び出しの前に`defer`を付けることで、その関数を含まれている関数の最後に呼び出します。

場合によっては、ファイルを閉じるなどのリソースをクリーンアップする必要があります。この場合、サーバーがポートをリッスンし続けないようにサーバーを閉じる必要があります。

これを関数の最後に実行したいが、将来のコードの読む人のために、サーバーを作成した場所の近くに命令を置いておきます。

私たちのリファクタリングは改善であり、これまでに取り上げたGo機能を考えると合理的なソリューションですが、ソリューションをよりシンプルにすることができます。

### プロセスの同期

* Goが同時実行性に優れているのに、なぜWebサイトの速度を次々にテストするのですか？両方を同時にチェックできるはずです。
* リクエストの「正確な応答時間」については特に気にしません。どちらが最初に返されるかを知りたいだけです。

これを行うために、プロセスを非常に簡単かつ明確に同期するのに役立つ`select`と呼ばれる新しい構成を導入します。

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
        http.Get(url)
        close(ch)
    }()
    return ch
}
```

#### `ping`

`chan struct{}`を作成して返す関数`ping`を定義しました。

私たちのケースでは、チャネルに送信されるタイプを _ケア_ するのではなく、**完了したことを通知したいだけです**。
チャネルを閉じることは完全に機能します！

なぜ`struct{}`で、`bool`のような別の型ではないのですか？まあ、`chan struct{}`はメモリの観点から利用できる最小のデータ型なので、`bool`に対して割り当てはありません。
ちゃんと閉じて何も送信しないので、なぜ何かを割り当てるのですか？

同じ関数内で、`http.Get(url)`を完了すると、そのチャネルに信号を送信するゴルーチン（`goroutine`）を開始します。

**常にチャネルを作成する**

チャネルを作成するときに、`make`を使用する方法に注意してください。
「`var ch chan struct{}`」と言うのではなく。
`var`を使用すると、変数は型の「ゼロ」値で初期化されます。したがって、`string`の場合は`""`、`int`の場合は`0`になります。

チャネルの場合、ゼロ値は`nil`であり、`<-`で送信しようとすると、`nil`チャネルに送信できないため、永久にブロックされます。

[これは Go Playground で実際に見ることができます](https://play.golang.org/p/IIbeAox5jKA)

#### `select`

同時実行の章を思い出すと、`myVar := <-ch`を使用して値がチャネルに送信されるのを待つことができます。値を待っているので、これは _blocking_ 呼び出しです。

`select`でできることは、_multiple_ チャネルで待機することです。
値を送信する最初のものは「勝ち」、`case`の下のコードが実行されます。

`select`で`ping`を使用して、`URL`ごとに2つのチャネルを設定します。
最初にチャネルに書き込む方は、コードが`select`で実行され、その結果、`URL`が返されます（勝者となります）。

これらの変更後、コードの背後にある意図は非常に明確になり、実装は実際にはより単純になります。

### タイムアウト

最後の要件は、`Racer`に10秒以上かかる場合にエラーを返すことでした。

## 最初にテストを書く

```go
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
```

テストサーバーがこのシナリオを実行するために戻るまでに10秒以上かかるようにしました。
ここでは、`Racer`が2つの値を返すことを期待しています。勝つURL（このテストでは`_`で無視）と `error`。

## テストを実行してみます

`./racer_test.go:37:10: assignment mismatch: 2 variables but 1 values`

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

勝者と`error`を返すように`Racer`の署名を変更します。ハッピーケースの場合は`nil`を返します。

コンパイラーは _first test_ が1つの値しか検索しないと文句を言うので、この行を`got, _ := Racer(slowURL, fastURL)`に変更します。
確認すると、私たちの幸せなシナリオでエラーが発生しないことを確認する必要があります。

11秒後に実行すると、失敗します。

```text
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

`time.After`は、`select`を使用する場合に非常に便利な関数です。

今回のケースでは発生しませんでしたが、リッスンしているチャネルが値を返さない場合、永久にブロックするコードを書く可能性があります。`time.After`は、`chan`（ `ping`のように）を返し、指定した時間が経過すると信号を送ります。

私たちにとってこれは完璧です。`a`または`b`が戻って成功した場合、10秒に到達すると、`time.After`がシグナルを送信し、`error`を返します。

### 遅いテスト

問題は、このテストの実行に10秒かかることです。そのような単純なロジックの場合、これは気分が良くありません。

私たちができることは、タイムアウトを構成可能にすることです。したがって、テストでは非常に短いタイムアウトを設定できます。コードを実際に使用する場合は、10秒に設定できます。

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

急いでこのデフォルト値を両方のテストに追加する前に、リッスンしてみましょう_。

* 「ハッピー」テストのタイムアウトを気にしますか？
* タイムアウトに関する要件は明示的でした

この知識を踏まえて、テストとコードのユーザーの両方に同情するように少しリファクタリングしてみましょう。

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

ユーザーと最初のテストでは、`Racer`（これは内部で`ConfigurableRacer`を使用します）を使用でき、悲しいパステストでは`ConfigurableRacer`を使用できます。

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

    t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
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
* 場合によっては、`case.`の1つに`time.After`を含めて、システムが永久にブロックされるのを防ぐ必要があります。

### `httptest`

* テストサーバーを作成して、信頼性の高い制御可能なテストを作成できる便利な方法。
* 「実際の」`net/http`サーバーと同じインターフェースを使用します。これは一貫性があり、習得するのに時間がかかります。
