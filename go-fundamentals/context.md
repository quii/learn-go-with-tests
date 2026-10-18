---
description: コンテキスト
---

# コンテキスト

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/context)

ソフトウェアは、長時間実行され、リソースを大量に消費するプロセスを（多くの場合ゴルーチンで）開始することがよくあります。そのきっかけとなったアクションがキャンセルされたり、何らかの理由で失敗したりした場合は、アプリケーション全体で一貫した方法でこれらのプロセスを停止する必要があります。

これをきちんと管理しないと、あなたが誇りに思っているキビキビとしたGoアプリケーションで、デバッグが難しいパフォーマンスの問題が起こり始めるかもしれません。

この章では、`context`パッケージを使用して、実行時間の長いプロセスを管理できるようにします。

まずは古典的な例として、アクセスされると、レスポンスで返すデータを取得するために長時間かかる可能性のあるプロセスを開始するWebサーバーから始めます。

データを取得できる前にユーザーがリクエストをキャンセルするシナリオを試し、プロセスに中止するよう確実に伝わるようにします。

手始めに、正常系（ハッピーパス）のコードを用意しておきました。これがサーバーコードです。

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, store.Fetch())
    }
}
```

関数`Server`は`Store`を受け取り、`http.HandlerFunc`を返します。`Store`は次のように定義されています。

```go
type Store interface {
    Fetch() string
}
```

返された関数は、`store`の`Fetch`メソッドを呼び出してデータを取得し、それをレスポンスに書き込みます。

テストで使用する、`Store`に対応するスパイもあります。

```go
type SpyStore struct {
	response string
}

func (s *SpyStore) Fetch() string {
	return s.response
}

func TestServer(t *testing.T) {
	data := "hello, world"
	svr := Server(&SpyStore{data})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	svr.ServeHTTP(response, request)

	if response.Body.String() != data {
		t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
	}
}
```

ハッピーパスができたので、ユーザーがリクエストをキャンセルするまでに`Store`が`Fetch`を完了できない、より現実的なシナリオを作りたいと思います。

## 最初にテストを書く

ハンドラーには`Store`に作業のキャンセルを伝える方法が必要になるので、インターフェースを更新します。

```go
type Store interface {
    Fetch() string
    Cancel()
}
```

スパイを調整して、`data`を返すまでに時間がかかるようにし、キャンセルを指示されたことを知る方法も持たせる必要があります。`Store`インターフェースを実装するために、`Cancel`メソッドも追加しなければなりません。

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

100ミリ秒が経過する前にリクエストをキャンセルする新しいテストを追加して、ストアがキャンセルされるかどうかを確認しましょう。

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
	data := "hello, world"
	store := &SpyStore{response: data}
	svr := Server(store)

	request := httptest.NewRequest(http.MethodGet, "/", nil)

	cancellingCtx, cancel := context.WithCancel(request.Context())
	time.AfterFunc(5*time.Millisecond, cancel)
	request = request.WithContext(cancellingCtx)

	response := httptest.NewRecorder()

	svr.ServeHTTP(response, request)

	if !store.cancelled {
		t.Error("store was not told to cancel")
	}
})
```

[Goブログ: Context](https://blog.golang.org/context)より

> context パッケージは、既存の Context 値から新しい Context 値を派生させる関数を提供します。これらの値はツリーを形成します。Context がキャンセルされると、そこから派生したすべての Context もキャンセルされます。

キャンセルが特定のリクエストのコールスタック全体に伝播されるように、コンテキストを派生させることが重要です。

ここでは、`request`から新しい`cancellingCtx`を派生させています。このとき`cancel`関数が返されます。次に、`time.AfterFunc`を使用して、その関数が5ミリ秒後に呼び出されるようにスケジュールします。最後に、`request.WithContext`を呼び出して、この新しいコンテキストをリクエストで使用します。

## テストを実行してみます

テストは予想通り失敗します。

```text
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.00s)
        context_test.go:62: store was not told to cancel
```

## 成功させるのに十分なコードを書く

TDDの規律を守ることを忘れないでください。テストに合格するための _最小限_ のコードを書きます。

```go
func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        store.Cancel()
        fmt.Fprint(w, store.Fetch())
    }
}
```

これでテストは通りますが、気持ちの良いものではありませんよね？ _すべてのリクエスト_ でフェッチの前に`Cancel()`を呼び出すべきでないのは明らかです。

規律を守ったことで、テストの欠陥が明らかになりました。これは良いことです！

キャンセルされないことを確認するように、ハッピーパスのテストを更新する必要があります。

```go
t.Run("returns data from store", func(t *testing.T) {
	data := "hello, world"
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

両方のテストを実行すると、今度はハッピーパスのテストが失敗するはずです。これで、より理にかなった実装をせざるを得なくなります。

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

`context`にはメソッド`Done()`があり、コンテキストが「完了」または「キャンセル」されたときにシグナルが送られるチャネルを返します。そのシグナルを待ち受け、受け取ったら`store.Cancel`を呼び出したいのですが、`Store`がそれより先に`Fetch`を完了できた場合は無視したいわけです。

これを実現するために、ゴルーチンで`Fetch`を実行し、その結果を新しいチャネル`data`に書き込みます。そして`select`を使って2つの非同期プロセスを実質的に競争させ、その結果に応じてレスポンスを書き込むか、`Cancel`を呼び出します。

## リファクタリング♪

スパイにアサーションメソッドを作ることで、テストコードを少しリファクタリングできます。

```go
type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Error("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Error("store was told to cancel")
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

このアプローチでも問題はありませんが、Goらしい書き方でしょうか？

Webサーバーが`Store`を手動でキャンセルすることまで気にかけるのは理にかなっているでしょうか？`Store`がさらに他の遅いプロセスに依存していたらどうでしょう？`Store.Cancel`が、自身が依存するすべてのものにキャンセルを正しく伝播するようにしなければなりません。

`context`の主な利点の1つは、キャンセルを提供するための一貫した方法であることです。

[go doc から](https://golang.org/pkg/context/)

> サーバーへの受信リクエストは Context を作成し、サーバーへの送信呼び出しは Context を受け取るべきです。その間にある関数呼び出しの連鎖は Context を伝播しなければならず、必要に応じて`WithCancel`、`WithDeadline`、`WithTimeout`、`WithValue`で作成した派生 Context に置き換えることができます。コンテキストがキャンセルされると、そのコンテキストから派生したすべてのコンテキストもキャンセルされます。

再び[Goブログ: Context](https://blog.golang.org/context)から

> Googleでは、受信リクエストと送信リクエストの間の呼び出しパス上にあるすべての関数に、最初の引数としてContextパラメータを渡すことをGoプログラマーに義務付けています。これにより、多くの異なるチームが開発したGoコードがうまく連携できます。タイムアウトとキャンセルを簡単に制御でき、セキュリティ認証情報のような重要な値がGoプログラム内を適切に受け渡されることを保証します。

（少し手を止めて、すべての関数がコンテキストを渡さなければならないことの影響と、その使い勝手について考えてみてください。）

少し不安になりましたか？それで良いのです。それでもこのアプローチに従ってみましょう。代わりに`context`を`Store`に渡し、`Store`に責任を持たせます。そうすれば`Store`も`context`を自身の依存先に渡すことができ、それらもまた自分自身を停止する責任を持てるようになります。

## 最初にテストを書く

責任が変わるので、既存のテストを変更する必要があります。ハンドラーが今担当するのは、下流の`Store`にコンテキストを確実に渡すことと、キャンセルされたときに`Store`から返されるエラーを処理することだけです。

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
				log.Println("spy store got cancelled")
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

スパイを、`context`と連携する本物のメソッドのように動作させる必要があります。

ここでは遅いプロセスをシミュレートしており、ゴルーチンの中で1文字ずつ文字列を追加して、結果をゆっくり組み立てています。ゴルーチンが作業を終了すると、文字列を`data`チャネルに書き込みます。ゴルーチンは`ctx.Done`を待ち受けており、そのチャネルにシグナルが送られると作業を停止します。

最後に、コードは別の`select`を使用して、そのゴルーチンが作業を完了するか、キャンセルが発生するのを待ちます。

これは以前のアプローチに似ています。Goの並行処理プリミティブを使用して2つの非同期プロセスを競争させ、何を返すかを決めています。

`context`を受け入れる独自の関数とメソッドを記述する場合も同様のアプローチをとるので、何が起こっているのかを確実に理解してください。

ようやくテストを更新できます。まずハッピーパスのテストを修正できるように、キャンセルのテストをコメントアウトしておきます。

```go
t.Run("returns data from store", func(t *testing.T) {
	data := "hello, world"
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

ハッピーパスは…ハッピーになったはずです。これで、もう一方のテストを修正できます。

## 最初にテストを書く

エラーの場合には、いかなるレスポンスも書き込まないことをテストする必要があります。残念ながら、`httptest.ResponseRecorder`にはこれを確認する方法がないため、テスト用に独自のスパイを自作する必要があります。

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

`SpyResponseWriter`は`http.ResponseWriter`を実装しているので、テストで使用できます。

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
	data := "hello, world"
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

サーバーコードはもはやキャンセルに対して明示的な責任を負わなくなったので、単純になったことがわかります。単に`context`を渡し、起こりうるキャンセルを下流の関数が尊重してくれることに頼っています。

## まとめ

### カバーしたこと

* クライアントによってリクエストがキャンセルされたHTTPハンドラーをテストする方法。
* キャンセルを管理するためのコンテキストの使用方法。
* `context`を受け取り、ゴルーチン、`select`、チャネルを使ってそれを利用し、自分自身をキャンセルする関数の書き方。
* リクエストスコープのコンテキストをコールスタック全体に伝播させてキャンセルを管理するという、Googleのガイドラインに従うこと。
* 必要に応じて`http.ResponseWriter`の独自のスパイを自作する方法。

### `context.Value`はどうですか？

[Michal Štrba](https://faiface.github.io/post/context-should-go-away-go2/)と私は同様の意見を持っています。

> 私の（存在しない）会社で`ctx.Value`を使ったら、クビです

一部のエンジニアは、 _便利に感じる_ という理由で`context`を介して値を渡すことを提唱しています。

利便性は、しばしば悪いコードの原因になります。

`context.Values`の問題は、それが単なる型付けされていないマップであるため、型安全性がなく、値が実際には含まれていない場合にも対処しなければならないことです。あるモジュールから別のモジュールへとマップのキーによる結合を作る必要があり、誰かが何かを変更すると、いろいろなものが壊れ始めます。

要するに、**関数が何らかの値を必要とするなら、`context.Value`から取り出そうとするのではなく、型付きパラメータとして渡してください**。そうすれば静的にチェックされ、誰もが見られるようにドキュメント化されます。

#### しかし...

一方で、トレースIDのようなリクエストとは直交する情報をコンテキストに含めるのは役立つ場合があります。この情報はコールスタック上のすべての関数で必要になるわけではなく、パラメータにすると関数シグネチャが非常に乱雑になってしまうでしょう。

[Jack Lindamoodは、**Context.Valueは制御するためではなく、情報を伝えるためのもの**だと述べています](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)

> `context.Value`の内容は、ユーザーではなくメンテナーのためのものです。ドキュメントに記載された結果や期待される結果を得るために必須の入力であってはなりません。

### 追加資料

* [Michal Štrbaによる「Context should go away for Go 2」](https://faiface.github.io/post/context-should-go-away-go2/)はとても楽しく読めました。彼の主張は、どこにでも`context`を渡さなければならないのは悪い兆候（スメル）であり、キャンセルに関する言語の欠陥を示しているというものです。ライブラリのレベルではなく、言語のレベルで何らかの形で解決されるほうが良いと彼は述べています。それが実現するまでは、実行時間の長いプロセスを管理したいなら`context`が必要です。
* [Goブログでは、`context`を使用する動機についてさらに説明し、いくつかの例を挙げています](https://blog.golang.org/context)
