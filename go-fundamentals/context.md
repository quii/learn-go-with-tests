---
description: Context (長期実行プロセスの管理に役立つパッケージ)
---

# コンテキスト

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/context)

ソフトウェアは、多くの場合、長時間実行され、リソースを大量に消費するプロセスを開始します（多くの場合、ゴルーチンで）。これを引き起こしたアクションがキャンセルされるか、何らかの理由で失敗した場合は、アプリケーションを通じてこれらのプロセスを一貫した方法で停止する必要があります。

これを管理しないと、非常に誇りに思っているキレの良いGoアプリケーションは、パフォーマンスの問題のデバッグが困難になる可能性があります。

この章では、`context`パッケージを使用して、実行時間の長いプロセスを管理します。

まず、ヒットしたときに長時間実行される可能性のあるプロセスを開始して、データをフェッチして応答で返すWebサーバーの古典的な例から始めます。

データを取得する前にユーザーがリクエストをキャンセルするシナリオを実行し、プロセスが中止されるように指示します。

私たちは幸せなパスにいくつかのコードを設定して始めました。
これがサーバーコードです。

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, store.Fetch())
    }
}
```

関数`Server`は`Store`を受け取り、`http.HandlerFunc`を返します。

`Store`は次のように定義されます

```go
type Store interface {
    Fetch() string
}
```

返された関数は、`store`の` Fetch`メソッドを呼び出してデータを取得し、それを応答に書き込みます。

テストで使用する`Store`に対応するスタブがあります。

```go
type StubStore struct {
    response string
}

func (s *StubStore) Fetch() string {
    return s.response
}

func TestHandler(t *testing.T) {
    data := "hello, world"
    svr := Server(&StubStore{data})

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }
}
```

幸せなパスができたので、ユーザーがリクエストをキャンセルする前に`Store`が`Fetch`を完了できない、より現実的なシナリオを作成したいと思います。

## 最初にテストを書く

私たちのハンドラーは、`Store`に作業をキャンセルしてインターフェースを更新するように指示する方法が必要になります。

```go
type Store interface {
    Fetch() string
    Cancel()
}
```

スパイを調整する必要があるので、`data`を返すには時間がかかり、キャンセルするように指示されたことを知る方法があります。呼び出し方法を確認しているので、名前を`SpyStore`に変更します。

`Store`インターフェースを実装するメソッドとして`Cancel`を追加する必要があります。

```go
type SpyStore struct {
    response string
    cancelled bool
}

func (s *SpyStore) Fetch() string {
    time.Sleep(100 * time.Millisecond)
    return s.response
}

func (s *SpyStore) Cancel() {
    s.cancelled = true
}
```

100ミリ秒前にリクエストをキャンセルする新しいテストを追加して、ストアがキャンセルされるかどうかを確認します。

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
      store := &SpyStore{response: data}
      svr := Server(store)

      request := httptest.NewRequest(http.MethodGet, "/", nil)

      cancellingCtx, cancel := context.WithCancel(request.Context())
      time.AfterFunc(5 * time.Millisecond, cancel)
      request = request.WithContext(cancellingCtx)

      response := httptest.NewRecorder()

      svr.ServeHTTP(response, request)

      if !store.cancelled {
          t.Errorf("store was not told to cancel")
      }
  })
```

[Goブログ: コンテキスト`Context`](https://blog.golang.org/context)

> コンテキストパッケージは、既存のコンテキスト値から新しいコンテキスト値を導出する関数を提供します。これらの値はツリーを形成します。コンテキストが取り消されると、それから派生したすべてのコンテキストも取り消されます。

キャンセルが特定のリクエストのコールスタック全体に伝播されるように、コンテキストを派生させることが重要です。

私たちがすることは、`cancel`関数を返す`request`から新しい`cancellingCtx`を派生させることです。次に、`time.AfterFunc`を使用して、その関数が5ミリ秒で呼び出されるようにスケジュールします。最後に、`request.WithContext`を呼び出して、この新しいコンテキストをリクエストで使用します。

## テストを実行してみます

テストは予想通り失敗します。

```go
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.00s)
        context_test.go:62: store was not told to cancel
```

## 成功させるのに十分なコードを書く

TDDの訓練を受けることを忘れないでください。テストに合格するための _minimal_ 量のコードを記述します。

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        store.Cancel()
        fmt.Fprint(w, store.Fetch())
    }
}
```

これはこのテストに合格しますが、気分が良くありません！
あらゆるリクエストをフェッチする前に、`Store`をキャンセルしてはいけません。

懲戒処分を受けることで、テストの欠陥が明らかになり、これは良いことです！

キャンセルされないことを確認するために、幸せなパステストを更新する必要があります。

```go
t.Run("returns data from store", func(t *testing.T) {
    store := &SpyStore{response: data}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }

    if store.cancelled {
        t.Error("it should not have cancelled the store")
    }
})
```

両方のテストを実行すると、ハッピーパステストが失敗し、より賢明な実装を実行する必要があります。

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        data := make(chan string, 1)

        go func() {
            data <- store.Fetch()
        }()

        select {
        case d := <-data:
            fmt.Fprint(w, d)
        case <-ctx.Done():
            store.Cancel()
        }
    }
}
```

ここで何をしましたか？

`context`にはメソッド`Done()`があり、コンテキストが「完了」または「キャンセル」されたときに信号を送信するチャネルを返します。そのシグナルをリッスンし、それを取得した場合は`store.Cancel`を呼び出しますが、`Store`がその前に`Fetch`を実行した場合は無視します。

これを管理するには、ゴルーチンで`Fetch`を実行し、結果を新しいチャネル`data`に書き込みます。次に、`select`を使用して2つの非同期プロセスに効率的に競合し、応答または`Cancel`を書き込みます。

## リファクタリング

スパイでアサーションメソッドを作成することで、テストコードを少しリファクタリングできます。

```go
func (s *SpyStore) assertWasCancelled() {
    s.t.Helper()
    if !s.cancelled {
        s.t.Errorf("store was not told to cancel")
    }
}

func (s *SpyStore) assertWasNotCancelled() {
    s.t.Helper()
    if s.cancelled {
        s.t.Errorf("store was told to cancel")
    }
}
```

スパイを作成するときは、`*testing.T`を渡すことを忘れないでください。

```go
func TestServer(t *testing.T) {
    data := "hello, world"

    t.Run("returns data from store", func(t *testing.T) {
        store := &SpyStore{response: data, t: t}
        svr := Server(store)

        request := httptest.NewRequest(http.MethodGet, "/", nil)
        response := httptest.NewRecorder()

        svr.ServeHTTP(response, request)

        if response.Body.String() != data {
            t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
        }

        store.assertWasNotCancelled()
    })

    t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
        store := &SpyStore{response: data, t: t}
        svr := Server(store)

        request := httptest.NewRequest(http.MethodGet, "/", nil)

        cancellingCtx, cancel := context.WithCancel(request.Context())
        time.AfterFunc(5*time.Millisecond, cancel)
        request = request.WithContext(cancellingCtx)

        response := httptest.NewRecorder()

        svr.ServeHTTP(response, request)

        store.assertWasCancelled()
    })
}
```

このアプローチは大丈夫ですが、慣用的ですか？

私たちのウェブサーバーが手動で`Store`をキャンセルすることに関心を持つことは理にかなっていますか？`Store`が他の実行速度の遅いプロセスに依存している場合はどうなりますか？
`Store.Cancel`がキャンセルを依存するすべてに正しくキャンセルすることを確認する必要があります。

`context`の主なポイントの1つは、キャンセルを提供する一貫した方法であることです。

[go doc から](https://golang.org/pkg/context/)

> サーバーへの着信要求はコンテキストを作成し、サーバーへの発信呼び出しはコンテキストを受け入れる必要があります。それらの間の関数呼び出しのチェーンは、コンテキストを伝播する必要があり、オプションで、`WithCancel`、`WithDeadline`、`WithTimeout`、または`WithValue`を使用して作成された派生コンテキストに置き換えます。コンテキストがキャンセルされると、そのコンテキストから派生したすべてのコンテキストもキャンセルされます。

再び[Goブログ: コンテキスト`Context`](https://blog.golang.org/context)から

> Googleでは、Goプログラマーが、最初の引数として、着信要求と発信要求の間の呼び出しパス上のすべての関数にContextパラメーターを渡す必要があります。これにより、多くの異なるチームが開発したGoコードを適切に相互運用できます。タイムアウトとキャンセルを簡単に制御し、セキュリティ認証情報などの重要な値がGoプログラムを適切に通過するようにします。

（少し間を置いて、コンテキストで送信する必要があるすべての機能の影響と、その人間工学について考えてください。）

少し不安ですか？大丈夫です。
そのアプローチを試してみましょう。代わりに、`context`を介して私たちの`Store`に渡し、責任を持たせましょう。そうすることで、`context`をその依存関係に渡すこともでき、それらも依存を停止する責任があります。

## 最初にテストを書く

責任が変化するにつれて、既存のテストを変更する必要があります。ハンドラーが今担当する唯一のことは、コンテキストがダウンストリームの`Store`に送信されることと、キャンセルされたときに`Store`から発生するエラーを処理することです。

`Store`インターフェースを更新して、新しい責任を示しましょう。

```go
type Store interface {
    Fetch(ctx context.Context) (string, error)
}
```

とりあえずハンドラー内のコードを削除してください

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    }
}
```

`SpyStore`を更新します

```go
type SpyStore struct {
    response string
    t        *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
    data := make(chan string, 1)

    go func() {
        var result string
        for _, c := range s.response {
            select {
            case <-ctx.Done():
                s.t.Log("spy store got cancelled")
                return
            default:
                time.Sleep(10 * time.Millisecond)
                result += string(c)
            }
        }
        data <- result
    }()

    select {
    case <-ctx.Done():
        return "", ctx.Err()
    case res := <-data:
        return res, nil
    }
}
```

スパイを`context`で機能する実際の方法のように動作させる必要があります。

遅いプロセスをシミュレートしていて、ゴルーチンで文字ごとに文字列を追加して、結果をゆっくり構築しています。ゴルーチンが作業を終了すると、文字列を`data`チャネルに書き込みます。ゴルーチンは`ctx.Done`をリッスンし、そのチャネルでシグナルが送信されると作業を停止します。

最後に、コードは別の`select`を使用して、そのゴルーチンが作業を完了するか、キャンセルが発生するのを待ちます。

これは以前のアプローチに似ています。Goの同時実行プリミティブを使用して、2つの非同期プロセスが互いに競合して、何を返すかを決定します。

`context`を受け入れる独自の関数とメソッドを記述する場合も同様のアプローチをとるので、何が起こっているのかを確実に理解してください。

最後に、テストを更新できます。幸せなパステストを最初に修正できるように、キャンセルのテストをコメント化します。

```go
t.Run("returns data from store", func(t *testing.T) {
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }
})
```

## テストを実行してみます

```text
=== RUN   TestServer/returns_data_from_store
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/returns_data_from_store (0.00s)
        context_test.go:22: got "", want "hello, world"
```

## 成功させるのに十分なコードを書く

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, _ := store.Fetch(r.Context())
        fmt.Fprint(w, data)
    }
}
```

私たちの幸せな道は...幸せでなければなりません。これで、他のテストを修正できます。

## 最初にテストを書く

エラーの場合には、いかなる種類の応答も書かないことをテストする必要があります。悲しいことに、`httptest.ResponseRecorder`にはこれを理解する方法がないため、これをテストするために私たち自身のスパイをロールする必要があります。

```go
type SpyResponseWriter struct {
    written bool
}

func (s *SpyResponseWriter) Header() http.Header {
    s.written = true
    return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
    s.written = true
    return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
    s.written = true
}
```

テストで使用できるように、`SpyResponseWriter`は`http.ResponseWriter`を実装します。

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)

    cancellingCtx, cancel := context.WithCancel(request.Context())
    time.AfterFunc(5*time.Millisecond, cancel)
    request = request.WithContext(cancellingCtx)

    response := &SpyResponseWriter{}

    svr.ServeHTTP(response, request)

    if response.written {
        t.Error("a response should not have been written")
    }
})
```

## テストを実行してみます

```text
=== RUN   TestServer
=== RUN   TestServer/tells_store_to_cancel_work_if_request_is_cancelled
--- FAIL: TestServer (0.01s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.01s)
        context_test.go:47: a response should not have been written
```

## 成功させるのに十分なコードを書く

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        data, err := store.Fetch(r.Context())

        if err != nil {
            return // todo: log error however you like
        }

        fmt.Fprint(w, data)
    }
}
```

この後、サーバーコードはキャンセルの明示的な責任がなくなり、単純化されていることがわかります。サーバーコードは単に`context`を通過し、発生する可能性のあるすべてのキャンセルを下流の関数に依存しています。

## まとめ

### カバーしたこと

* リクエストがクライアントによってキャンセルされたHTTPハンドラをテストする方法。
* キャンセルを管理するためのコンテキストの使用方法。
* `context`を受け入れ、それを使ってゴルーチン、`select`、およびチャネルを使用してそれ自体をキャンセルする関数の作成方法。
* コールスタックを通じてリクエストスコープのコンテキストを伝播してキャンセルを管理する方法については、Googleのガイドラインに従ってください。
* 必要に応じて、`http.ResponseWriter`の独自のスパイをロールする方法。

### `context.Value`はどうですか？

[Michal Štrba](https://faiface.github.io/post/context-should-go-away-go2/)と私は同様の意見を持っています。

> 私の（non-existent）会社で`ctx.Value`を使用すると、解雇されます

一部のエンジニアは、「便利」と感じて`context`を介して値を渡すことを提唱しています。

多くの場合、利便性が悪いコードの原因です。

`context.Values`の問題は、型付けされていないマップであるため、タイプの安全性がなく、値を実際に含まないように処理する必要があることです。あるモジュールから別のモジュールへのマップキーの結合を作成する必要があり、誰かが何かを変更すると、何かが壊れ始めます。

要するに、**関数がいくつかの値を必要とする場合、 `context.Value`からそれらをフェッチしようとするのではなく、型付きパラメーターとしてそれらを置きます**。これは静的にチェックされ、誰もが見ることができるように文書化されます。

#### しかし...

一方、リクエストに直交する情報（トレースIDなど）をコンテキストに含めると役立つ場合があります。潜在的に、この情報はコールスタックのすべての関数で必要とされるわけではなく、関数シグネチャが非常に乱雑になります。

[Jack Lindamoodによると、**Context.Valueは制御ではなく通知**](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)

> `context.Value`のコンテンツは、ユーザーではなくメンテナー向けです。文書化された結果または期待される結果の入力が必要になることはありません。

### 追加資料

* [MichalŠtrbaによるGo2のコンテキストはなくなるはずです](https://faiface.github.io/post/context-should-go-away-go2/)を読んで本当に楽しんでいました。彼の主張は、どこでも`context`を渡す必要があることは匂いであり、キャンセルに関する言語の欠陥を指摘しているということです。ライブラリレベルではなく、言語レベルで何らかの方法でこれを解決した方が良いと彼は言います。それが発生するまで、実行時間の長いプロセスを管理する場合は、`context`が必要になります。
* [Goブログでは、`context`を使用する動機についてさらに説明し、いくつかの例を挙げています](https://blog.golang.org/context)
