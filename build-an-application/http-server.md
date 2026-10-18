---
description: HTTP server
---

# HTTPサーバー

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/http-server)

あなたは、プレーヤーが勝ったゲームの数をユーザーが追跡できるWebサーバーを作成するように求められました。

* `GET /players/{name}`は、勝利の合計数を示す数値を返す必要があります
* `POST /players/{name}`は、その名前の勝利を記録し、後続の`POST`ごとに増分する必要があります

TDDアプローチに従い、できる限り迅速にソフトウェアを動作させ、解決策が見つかるまで小さな反復的な改善を行います。このアプローチを取ることによって

* 問題のあるスペースを常に小さく保つ
* なかなか抜け出すことができない状況に陥ってはいけません
* 行き詰まったり失われたりした場合でも、元に戻しても負荷は減りません。

## レッド、グリーン、リファクタリング（Red, green, refactor）

この本全体を通して、テストを作成して失敗するのを監視するTDDプロセスを強調し、（red）、それを機能させるための _minimal_ 量のコードを記述し、（green）してリファクタリングします。

最小限のコードを書くというこの規律は、TDDが与える安全性の観点から重要です。できるだけ早く「赤（red）」から抜け出すように努力する必要があります。

ケント・ベックは次のように説明しています。

> テストを迅速に実行し、プロセスで必要なあらゆる罪を犯します。

テストの安全性に後押しされてリファクタリングされるため、これらの罪を犯すことができます。

### これを行わないとどうなりますか？

赤で表示されている変更が多いほど、テストではカバーされない問題を追加する可能性が高くなります。

アイデアは、ウサギの穴に何時間も陥らないように、テストによって駆動される小さなステップで有用なコードを繰り返し書くことです。

### 鶏肉と卵

これを段階的に構築するにはどうすればよいですか？何かを保存せずにプレーヤーを`GET`することはできず、`GET`エンドポイントがすでに存在しない状態で`POST`が機能したかどうかを知るのは難しいようです。

これが _mocking_ の輝きです。

* プレーヤーのスコアを取得するには、`GET`に`PlayerStore` _thing_ が必要です。これはインターフェースである必要があるので、テストするときに、実際のストレージコードを実装する必要なく、コードをテストするための簡単なスタブを作成できます。
* `POST`の場合、`PlayerStore`への呼び出しを _spy_ して、プレーヤーが正しく保存されていることを確認できます。保存の実装は検索と連動しません。
* 機能するソフトウェアをすばやく用意するために、非常にシンプルなインメモリ実装を作成し、その後、任意のストレージメカニズムに基づく実装を作成できます。

## 最初にテストを書く

テストを作成し、ハードコードされた値を返すことでテストを成功させることができます。ケントベックはこれを「偽造（`Faking it`）」と呼んでいます。動作するテストができたら、その定数を削除するのに役立つテストをさらに記述できます。

この非常に小さなステップを実行することで、アプリケーションロジックをあまり気にすることなく、プロジェクト全体の構造を正しく機能させる重要な出発点を作ることができます。

GoでWebサーバーを作成するには、通常[ListenAndServe](https://golang.org/pkg/net/http/#ListenAndServe)を呼び出します。

```go
func ListenAndServe(addr string, handler Handler) error
```

これにより、ポートでリッスンするWebサーバーが起動し、すべてのリクエストに対してゴルーチンが作成され、[`Handler`](https://golang.org/pkg/net/http/#Handler)に対して実行されます。

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

タイプは、2つの引数を期待する`ServeHTTP`メソッドを実装することにより、ハンドラーインターフェースを実装します。1つ目は、レスポンスを書き込む場所で、2つ目はサーバーに送信されたHTTPリクエストです。

`server_test.go`というファイルを作成し、これらの2つの引数を受け取る関数`PlayerServer`のテストを書いてみましょう。送信されるリクエストは、プレーヤーのスコアを取得することです。これは`"20"`であると予想されます。

```go
func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
```

サーバーをテストするには、送信する`Request`が必要であり、ハンドラーが`ResponseWriter`に書き込む内容を _spy_ する必要があります。

* `http.NewRequest`を使用してリクエストを作成します。最初の引数はリクエストのメソッドで、2番目はリクエストのパスです。`nil`引数はリクエストの本文を参照します。この場合、設定する必要はありません。
* `net/http/httptest`には、`ResponseRecorder`というスパイが既に作成されているので、それを使用できます。応答として書き込まれた内容を検査するための多くの便利な方法があります。

## テストを実行してみます

`./server_test.go:13:2: undefined: PlayerServer`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

コンパイラーは正常に動いています。耳を傾けてください。

`server.go`というファイルを作成し、`PlayerServer`を定義します

```go
func PlayerServer() {}
```

再試行

```
./server_test.go:13:14: too many arguments in call to PlayerServer
    have (*httptest.ResponseRecorder, *http.Request)
    want ()
```

関数に引数を追加します

```go
import "net/http"

func PlayerServer(w http.ResponseWriter, r *http.Request) {

}
```

コードがコンパイルされ、テストが失敗します

```
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- FAIL: TestGETPlayers/returns_Pepper's_score (0.00s)
        server_test.go:20: got '', want '20'
```

## 成功させるのに十分なコードを書く

DIの章から、`Greet`関数を使用してHTTPサーバーに触れました。 net/httpの`ResponseWriter`もio`Writer`を実装しているため、`fmt.Fprint`を使用して文字列をHTTP応答として送信できることがわかりました。

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
```

これでテストに成功するはずです。

## 足場を完成させましょう

これをアプリケーションに結び付けたいと思います。これは重要です

* 実際に動作するソフトウェアを用意します。そのためのテストを記述したくありません。コードの動作を確認することをお勧めします。
* コードをリファクタリングすると、プログラムの構造が変更される可能性があります。これは、インクリメンタルアプローチの一部として、アプリケーションにも反映されるようにしたいと考えています。

アプリケーション用の新しいファイル`main.go`を作成し、このコードを配置します。

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(PlayerServer)
	log.Fatal(http.ListenAndServe(":5000", handler))
}
```

これまでのところ、すべてのアプリケーションコードが1つのファイルに含まれていますが、これは、物事を異なるファイルに分離する必要がある大規模なプロジェクトでは、ベストプラクティスではありません。

これを実行するには、ディレクトリ内のすべての`.go`ファイルを取得してプログラムをビルドする`go build`を実行します。その後、`./myprogram`で実行できます。

### `http.HandlerFunc`

以前に`Handler`インターフェースがサーバーを作るために実装する必要があるものであることを探りました。 通常は、`struct`を作成してそれを行い、独自のServeHTTPメソッドを実装してインターフェースを実装します。ただし、構造体のユースケースはデータを保持するためのものですが、_currently_ には状態がないため、データを作成するのは適切ではありません。

[HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)を使用すると、これを回避できます。

> HandlerFuncタイプは、通常の関数をHTTPハンドラーとして使用できるようにするアダプターです。fが適切なシグネチャを持つ関数である場合、HandlerFunc\(f\) はfを呼び出すハンドラーです。

```go
type HandlerFunc func(ResponseWriter, *Request)
```

ドキュメントから、タイプ`HandlerFunc`がすでに`ServeHTTP`メソッドを実装していることがわかります。`PlayerServer`関数をタイプキャストすることで、必要な`Handler`を実装しました。

### `http.ListenAndServe(":5000"...)`

`ListenAndServe` はリッスンするポートを `Handler` に受け取ります。問題がある場合、ウェブサーバーはエラーを返します。エラーの一例として、ポートがすでにリッスンされていることが考えられます。そのため、この呼び出しを `log.Fatal` でラップして、ユーザーのために、エラーをログに出力します。

これから行うのは、ハードコーディングされた値から離れるようにポジティブな変更を強制する _another_ テストを作成することです。

## 最初にテストを書く

別のサブテストをスイートに追加して、別のプレーヤーのスコアを取得しようとします。これにより、ハードコーディングされたアプローチが壊れます。

```go
t.Run("returns Floyd's score", func(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
	response := httptest.NewRecorder()

	PlayerServer(response, request)

	got := response.Body.String()
	want := "10"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
})
```

あなたは考えていたかもしれません。

> 確かに、どのプレイヤーがどのスコアを獲得するかを制御するために、何らかのストレージの概念が必要です。テストで値が非常に恣意的に見えるのは奇妙です。

できる限り小さなステップを実行するように心がけていることを忘れないでください。したがって、今は定数を壊そうとしているだけです。

## テストを実行してみます

```
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- PASS: TestGETPlayers/returns_Pepper's_score (0.00s)
=== RUN   TestGETPlayers/returns_Floyd's_score
    --- FAIL: TestGETPlayers/returns_Floyd's_score (0.00s)
        server_test.go:34: got '20', want '10'
```

## 成功させるのに十分なコードを書く

```go
//server.go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	if player == "Pepper" {
		fmt.Fprint(w, "20")
		return
	}

	if player == "Floyd" {
		fmt.Fprint(w, "10")
		return
	}
}
```

このテストにより、リクエストのURLを実際に確認して決定を迫られました。したがって、頭の中では、プレイヤーのストアとインターフェースについて心配している可能性があります。次の論理的なステップは、実際には _routing_ のようです。

店舗コードから始めた場合、必要な変更の量はこれに比べて非常に大きくなります。 **これは最終目標に向けたより小さなステップであり、テストによって推進されました**。

現在、ルーティングライブラリを使用するという誘惑に抵抗しています。テストに合格するための最小のステップにすぎません。

`r.URL.Path`はリクエストのパスを返すので、[`strings.TrimPrefix`](https://golang.org/pkg/strings/#TrimPrefix)を使用して、 `/players/`を削除します。要求されたプレーヤーを取得します。それほど堅牢ではありませんが、とりあえずはうまくいくでしょう。

## リファクタリング♪

スコアの取得を関数に分離することで、`PlayerServer`を簡略化できます

```go
//server.go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(name string) string {
	if name == "Pepper" {
		return "20"
	}

	if name == "Floyd" {
		return "10"
	}

	return ""
}
```

そして、いくつかのヘルパーを作成することで、テストのコードの一部をDRYにすることができます。

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
```

しかし、私たちはまだ幸せであってはなりません。私たちのサーバーがスコアを知っていることは正しくありません。

私たちのリファクタリングは何をすべきかをかなり明確にしました。

スコア計算をハンドラーの本体から関数`GetPlayerScore`に移動しました。これは、インターフェースを使用して懸念事項を分離するのに適切な場所のように感じます。

代わりにリファクタリングした関数をインターフェイスに移動してみましょう。

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
}
```

`PlayerServer`が`PlayerStore`を使用できるようにするには、それを参照する必要があります。これで、アーキテクチャを変更して、`PlayerServer`が`struct`になるようにする適切なタイミングのように感じられます。

```go
type PlayerServer struct {
	store PlayerStore
}
```

最後に、新しい構造体にメソッドを追加して既存のハンドラーコードを挿入することにより、`Handler`インターフェースを実装します。

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

他の唯一の変更は、定義したローカル関数（これで削除できます）ではなく、`store.GetPlayerScore`を呼び出してスコアを取得することです。

サーバーの完全なコードリストは次のとおりです。

```go
//server.go
type PlayerStore interface {
	GetPlayerScore(name string) int
}

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

### 問題を修正する

これはかなりの数の変更であり、テストとアプリケーションがコンパイルされなくなることがわかっています。リラックスして、コンパイラーにそれを実行させてください。

`./main.go:9:58: type PlayerServer is not an expression`

テストを変更して、代わりに`PlayerServer`の新しいインスタンスを作成し、そのメソッド`ServeHTTP`を呼び出す必要があります。

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	server := &PlayerServer{}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}
```

まだストアを作成することについてはまだ心配していないことに注意してください。できるだけ早くコンパイラーを渡したいだけです。

コンパイルするコードを優先し、次にテストに合格するコードを優先する習慣を身に付ける必要があります。

コードがコンパイルされていないときに（スタブストアのような）機能を追加することにより、潜在的に _more_ コンパイルの問題に直面することになります。

同じ理由で`main.go`はコンパイルされません。

```go
func main() {
	server := &PlayerServer{}
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

最後に、すべてがコンパイルされていますが、テストは失敗しています

```
=== RUN   TestGETPlayers/returns_the_Pepper's_score
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
    panic: runtime error: invalid memory address or nil pointer dereference
```

これは、テストで`PlayerStore`を渡していないためです。スタブを1つ作成する必要があります。

```go
//server_test.go
type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}
```

`map`は、テスト用のスタブ キー/値（key/value）ストアを作成する迅速で簡単な方法です。次に、テスト用にこれらのストアの1つを作成して、`PlayerServer`に送信します。

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := &PlayerServer{&store}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}
```

テストは成功し、見た目も良くなっています。ストアの導入により、コードの背後にある _intent_ がより明確になりました。`PlayerStore`にこのデータがあるので、それを`PlayerServer`で使用すると、次の応答が得られるはずであることを読者に伝えています。

### アプリケーションを実行します

これで、このリファクタリングを完了するために必要な最後のことは、アプリケーションの動作を確認することです。プログラムは起動するはずですが、 `http://localhost:5000/players/Pepper`でサーバーにアクセスしようとすると、恐ろしい応答が返されます。

これは、`PlayerStore`を渡していないためです。

1つの実装を作成する必要がありますが、意味のあるデータを格納していないため、当面はハードコーディングする必要があるため、現時点ではそれは困難です。

```go
//main.go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

`go build`を再度実行して同じURLにアクセスすると、`"123"`が表示されます。すばらしいとは言えませんが、データを保存するまでは、私たちができる最高のことです。
また、メインのアプリケーションが起動しても実際には動かないというのも、あまり気分のいいものではありませんでした。問題を確認するために、手動でテストする必要がありました。


私たちは次に何をすべきかについていくつかのオプションがあります

* プレーヤーが存在しないシナリオを処理します
* `POST /players/{name}`シナリオを処理します

`POST`シナリオは「ハッピーパス」に近づきますが、すでにそのコンテキストにいるため、最初に不足しているプレーヤーシナリオに取り組む方が簡単だと思います。残りは後で行います。

## 最初にテストを書く

不足しているプレーヤーのシナリオを既存のスイートに追加する

```go
//server_test.go
t.Run("returns 404 on missing players", func(t *testing.T) {
	request := newGetScoreRequest("Apollo")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusNotFound

	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}
})
```

## テストを実行してみます

```
=== RUN   TestGETPlayers/returns_404_on_missing_players
    --- FAIL: TestGETPlayers/returns_404_on_missing_players (0.00s)
        server_test.go:56: got status 200 want 404
```

## 成功させるのに十分なコードを書く

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

TDDの提唱者が「コードを最小限にするだけでコードをパスできるようにする」と言ったとき、私はときどき目を凝らします。

しかし、このシナリオは例をよく示しています。私は最低限の（正しくないことを知っている）を実行しました。これは**すべての応答**に`StatusNotFound`を書き込むことですが、すべてのテストに成功しています！

**テストに合格するために最低限必要なことを行うことで、テストのギャップを強調できます**。今回のケースでは、プレイヤーがストアに存在するときに`StatusOK`を取得する必要があることを表明していません。

他の2つのテストを更新してステータスを評価し、コードを修正します。

これが新しいテストです

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := &PlayerServer{&store}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
```

現在、すべてのテストでステータスをチェックしているので、これを容易にするヘルパー`assertStatus`を作成しました。

これで、最初の2つのテストは200ではなく404が原因で失敗します。そのため、スコアが0の場合にのみ見つからないことを返すように`PlayerServer`を修正できます。

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}
```

### スコアを保存する

ストアからスコアを取得できるようになったので、新しいスコアを格納できるようになりました。

## 最初にテストを書く

```go
//server_test.go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
	}
	server := &PlayerServer{&store}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}
```

まず、POSTで特定のルートに到達した場合に正しいステータスコードを取得することを確認します。これにより、異なる種類のリクエストを受け入れ、それを `GET /players/{name}`とは異なる方法で処理する機能を実行できます。これがうまくいったら、ハンドラーとストアの相互作用を評価し始めることができます。

## テストを実行してみます

```
=== RUN   TestStoreWins/it_returns_accepted_on_POST
    --- FAIL: TestStoreWins/it_returns_accepted_on_POST (0.00s)
        server_test.go:70: did not get correct status, got 404, want 202
```

## 成功させるのに十分なコードを書く

意図的に罪を犯しているので、リクエストのメソッドに基づく`if`ステートメントでうまくいくことを覚えておいてください。

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}
```

## リファクタリング♪

ハンドラーは少し混乱しています。コードを分割して、さまざまな機能を簡単に追跡して分離し、新しい機能に分離しましょう。

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		p.processWin(w)
	case http.MethodGet:
		p.showScore(w, r)
	}

}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}
```

これにより、`ServeHTTP`のルーティングの側面が少し明確になり、格納に関する次の反復が`processWin`の内部に収まるようになります

次に、`POST /players/{name}`を実行するときに、`PlayerStore`が勝利を記録するように指示されていることを確認します。

## 最初にテストを書く

これは、`StubPlayerStore`を新しい`RecordWin`メソッドで拡張し、その呼び出しをスパイすることで実現できます。

```go
//server_test.go
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
```

テストを拡張して、開始の呼び出しの数を確認します。

```go
//server_test.go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
	}
	server := &PlayerServer{&store}

	t.Run("it records wins when POST", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}
```

## テストを実行してみます

```
./server_test.go:26:20: too few values in struct initializer
./server_test.go:65:20: too few values in struct initializer
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

新しいフィールドを追加したので、`StubPlayerStore`を作成するコードを更新する必要があります

```go
//server_test.go
store := StubPlayerStore{
	map[string]int{},
	nil,
}
```

```
--- FAIL: TestStoreWins (0.00s)
    --- FAIL: TestStoreWins/it_records_wins_when_POST (0.00s)
        server_test.go:80: got 0 calls to RecordWin want 1
```

## 成功させるのに十分なコードを書く

特定の値ではなく呼び出しの数のみを評価しているため、最初の反復が少し小さくなります。

`RecordWin`を呼び出せるようにするには、インターフェイスを変更して、`PlayerStore`が何であるかについての`PlayerServer`の考えを更新する必要があります。

```go
//server.go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}
```

これを行うことにより、`main`はコンパイルされなくなります

```
./main.go:17:46: cannot use InMemoryPlayerStore literal (type *InMemoryPlayerStore) as type PlayerStore in field value:
    *InMemoryPlayerStore does not implement PlayerStore (missing RecordWin method)
```

コンパイラは何が悪いのかを教えてくれます。そのメソッドを持つように`InMemoryPlayerStore`を更新しましょう。

```go
//main.go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) RecordWin(name string) {}
```

テストを試して実行すると、コードのコンパイルに戻るはずですが、テストはまだ失敗しています。

`PlayerStore`に`RecordWin`があるので、`PlayerServer`内で呼び出すことができます。

```go
//server.go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
	p.store.RecordWin("Bob")
	w.WriteHeader(http.StatusAccepted)
}
```

テストを実行すれば合格です。明らかに、`"Bob"`は、`RecordWin`に送信したいものではないので、テストをさらに改良してみましょう。

## 最初にテストを書く

```go
//server_test.go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
	}
	server := &PlayerServer{&store}

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}
```

`winCalls`スライスに1つの要素があることがわかったので、最初の要素を安全に参照して、それが`player`と等しいことを確認できます。

## テストを実行してみます

```
=== RUN   TestStoreWins/it_records_wins_on_POST
    --- FAIL: TestStoreWins/it_records_wins_on_POST (0.00s)
        server_test.go:86: did not store correct winner got 'Bob' want 'Pepper'
```

## 成功させるのに十分なコードを書く

```go
//server.go
func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
```

`processWin`を`http.Request`に変更して、URLを見てプレーヤーの名前を抽出できるようにしました。それができたら、正しい値で`store`を呼び出してテストに合格することができます。

## リファクタリング♪

2つの場所で同じ方法でプレイヤー名を抽出しているので、このコードを少しDRYにすることができます。

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
```

テストは成功していますが、実際に機能するソフトウェアはありません。`main`を実行して、意図したとおりにソフトウェアを使用すると、`PlayerStore`を正しく実装するためのラウンドがないため、機能しません。これは問題ありません。ハンドラーに焦点を当てることで、事前に設計するのではなく、必要なインターフェースを特定しました。

`InMemoryPlayerStore`の周りにいくつかのテストを書き始めることができましたが、これは、プレーヤーのスコアを永続化するためのより堅牢な方法を実装するまで一時的にのみです（つまり、データベース）。

ここでは、`PlayerServer`と`InMemoryPlayerStore`の間に _統合テスト_ を記述して、機能を完成させます。これにより、`InMemoryPlayerStore`を直接テストする必要なく、アプリケーションが機能していると確信できるという目標を達成できます。それだけでなく、データベースでの`PlayerStore`の実装に取り​​掛かると、同じ統合テストでその実装をテストできます。

### 統合テスト

統合テストは、システムのより広い領域が機能することをテストするのに役立ちますが、次の点に注意する必要があります。

* 書くのが難しい
* 失敗すると、なぜ（通常、統合テストのコンポーネント内のバグであるか）を理解するのが難しくなるため、修正が困難になる可能性があります。
* 実行に時間がかかる場合があります（データベースなどの「実際の」コンポーネントで使用されることが多いため）。

そのためにも、_テストピラミッド（The Test Pyramid）_ をリサーチしておくことをお勧めします。

## 最初にテストを書く

簡潔にするために、最後のリファクタリングされた統合テストを紹介します。

```go
// server_integration_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := InMemoryPlayerStore{}
	server := PlayerServer{&store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
```

* 統合しようとしている2つのコンポーネント、`InMemoryPlayerStore`と`PlayerServer`を作成しています。
* 次に、`player`の3つの勝利を記録するために3つのリクエストを発行します。このテストのステータスコードは、それらがうまく統合されているかどうかには関係がないので、あまり心配していません。
* 次に注意するのは変数`response`を格納することなので、`player`のスコアを取得しようとするためです。

## テストを実行してみます

```
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:24: response body is wrong, got '123' want '3'
```

## 成功させるのに十分なコードを書く

私はここでいくつかの自由を取り、テストを書かずに慣れるよりも多くのコードを書きます。

_これは許可されています!_ 正常に機能していることを確認するテストはまだありますが、`InMemoryPlayerStore`で使用している特定のユニットの周りではありません。

このシナリオで行き詰まった場合は、変更を失敗したテストに戻し、`InMemoryPlayerStore`に関連するより具体的な単体テストを記述して、ソリューションを実行できるようにします。

```go
//in_memory_player_store.go
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}
```

* データを保存する必要があるので、`map[string]int`を`InMemoryPlayerStore`構造体に追加しました
* 便宜上、ストアを初期化するために`NewInMemoryPlayerStore`を作成し、それを使用するように統合テストを更新しました。
    ```go
    //server_integration_test.go
    store := NewInMemoryPlayerStore()
    server := PlayerServer{store}
    ```
* 残りのコードは、`map`をラップするだけです

統合テストに合格したので、`main`を変更して` NewInMemoryPlayerStore()`を使用するだけです。

```go
// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}
	log.Fatal(http.ListenAndServe(":5000", server))
}
```

ビルドして実行し、`curl`を使用してテストします。

* これを数回実行し、 `curl -X POST http://localhost:5000/players/Pepper`のようにプレーヤー名を変更します
* `curl http://localhost:5000/players/Pepper`でスコアを確認してください

すごい！ REST風のサービスを作成しました。これを進めるには、データストアを選択して、プログラムの実行時間よりも長いスコアを保持する必要があります。

* ストアを選択してください （Bolt? Mongo? Postgres? File system?）
* `PostgresPlayerStore`に` PlayerStore`を実装させる
* TDDの機能により、確実に機能します
* 統合テストに接続し、問題がないことを確認します
* 最後に`main`に接続します

## リファクタリング♪

私たちは、ほぼ、そこにいる！このような同時実行エラーを防ぐために少し努力しましょう

```
fatal error: concurrent map read and map write
```

ミューテックスを追加することで、特に`RecordWin`関数のカウンターに同時実行の安全性を適用します。 mutexの詳細については、同期の章をご覧ください。

## まとめ

### `http.Handler`

* このインターフェースを実装してWebサーバーを作成する
* 通常の関数を`http.Handler`に変換するには、`http.HandlerFunc`を使用します
* `httptest.NewRecorder`を使用して`ResponseWriter`として渡し、ハンドラーが送信する応答をスパイできるようにします
* `http.NewRequest`を使用して、システムに入ると予想されるリクエストを作成します

### インターフェース、モッキング、DI

* 小さなチャンクでシステムを繰り返し構築できます
* 実際のストレージを必要とせずにストレージを必要とするハンドラーを開発できます
* TDDは必要なインターフェースを駆動します

### 罪を犯して、リファクタリング（そして、ソースコントロールにコミットして）

* コンパイルに失敗したり、テストに失敗したりすることは、できるだけ早く脱出しなければならない赤信号として扱う必要があります。
* そのために必要なコードだけを書いてください。その後、リファクタリングをしてコードを良くしてください。
* コードがコンパイルされていなかったり、テストが失敗している間に多くの変更をしようとすると、問題を悪化させる危険性があります。
* このアプローチに固執すると、小さなテストを書くことを強制され、小さな変更を意味し、複雑なシステムでの作業を管理しやすくします。
