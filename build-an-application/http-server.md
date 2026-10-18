---
description: HTTPサーバー
---

# HTTPサーバー

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/http-server)

あなたは、プレーヤーが勝ったゲームの数をユーザーが追跡できるWebサーバーを作成するように求められました。

* `GET /players/{name}`は、勝利の合計数を示す数値を返す必要があります
* `POST /players/{name}`は、その名前の勝利を記録し、後続の`POST`ごとに増分する必要があります

TDDアプローチに従い、できる限り迅速にソフトウェアを動作させ、解決策が見つかるまで小さな反復的な改善を行います。このアプローチを取ることによって

* どの時点でも問題の範囲を小さく保つ
* 深みにはまり込まない
* 行き詰まったり迷ったりしても、変更を元に戻して失う作業はわずかで済む

## レッド、グリーン、リファクタリング（Red, green, refactor）

この本全体を通して、テストを書いて失敗するのを確認し（レッド）、それを動かすための _最小限_ のコードを書き（グリーン）、そしてリファクタリングするというTDDのプロセスを強調してきました。

最小限のコードを書くというこの規律は、TDDが与える安全性の観点から重要です。できるだけ早く「赤（red）」から抜け出すように努力する必要があります。

ケント・ベックは次のように説明しています。

> テストを素早く通せ。そのためならどんな罪を犯しても構わない。

後でテストの安全性に支えられながらリファクタリングするので、これらの罪を犯しても構わないのです。

### これを行わないとどうなりますか？

レッドの状態のまま行う変更が多いほど、テストでカバーされていない問題を持ち込む可能性が高くなります。

大事なのは、何時間も深みにはまり込まないように、テストに導かれた小さなステップで有用なコードを繰り返し書いていくことです。

### 鶏が先か卵が先か

これを段階的に構築するにはどうすればよいですか？何かを保存せずにプレーヤーを`GET`することはできず、`GET`エンドポイントがすでに存在しない状態で`POST`が機能したかどうかを知るのは難しいようです。

ここで _モック_ が真価を発揮します。

* `GET`では、プレーヤーのスコアを取得するために`PlayerStore`という _もの_ が必要です。これはインターフェースである必要があります。そうすれば、テストするときに、実際のストレージのコードを実装しなくても、コードをテストするための簡単なスタブを作成できます。
* `POST`では、`PlayerStore`への呼び出しを _スパイ_ して、プレーヤーが正しく保存されていることを確認できます。保存の実装が取得の実装と結合することはありません。
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

型は、2つの引数を受け取る`ServeHTTP`メソッドを実装することで`Handler`インターフェースを実装します。1つ目は _レスポンスを書き込む_ 場所で、2つ目はサーバーに送信されたHTTPリクエストです。

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

サーバーをテストするには、送信する`Request`が必要であり、ハンドラーが`ResponseWriter`に書き込む内容を _スパイ_ する必要があります。

* `http.NewRequest`を使用してリクエストを作成します。最初の引数はリクエストのメソッドで、2番目はリクエストのパスです。`nil`引数はリクエストのボディで、今回は設定する必要はありません。
* `net/http/httptest`には、`ResponseRecorder`というスパイが既に用意されているので、それを使用できます。レスポンスとして書き込まれた内容を検査するための便利なメソッドがたくさんあります。

## テストを実行してみます

`./server_test.go:13:2: undefined: PlayerServer`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

コンパイラは私たちを助けてくれる存在です。その声に耳を傾けましょう。

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

DIの章では、`Greet`関数でHTTPサーバーに少し触れました。net/httpの`ResponseWriter`はioの`Writer`も実装しているので、`fmt.Fprint`を使って文字列をHTTPレスポンスとして送信できることを学びました。

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
```

これでテストに成功するはずです。

## 足場を完成させましょう

これをアプリケーションに組み込みたいと思います。これが重要なのは、次の理由からです。

* _実際に動作するソフトウェア_ が手に入ります。テストのためだけにテストを書きたいわけではありません。コードが実際に動くのを見るのは良いことです。
* コードをリファクタリングすると、プログラムの構造を変更することになるでしょう。インクリメンタルなアプローチの一環として、それがアプリケーションにも反映されるようにしたいのです。

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

以前、サーバーを作るには`Handler`インターフェースを実装する必要があることを見てきました。_通常_ は、`struct`を作成し、独自のServeHTTPメソッドを実装することでインターフェースを実装します。しかし、構造体の用途はデータを保持することであり、_現時点では_ 状態を持っていないので、構造体を作るのはしっくりきません。

[HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)を使用すると、これを回避できます。

> HandlerFunc型は、通常の関数をHTTPハンドラーとして使用できるようにするアダプターです。fが適切なシグネチャを持つ関数である場合、HandlerFunc(f) はfを呼び出すハンドラーです。

```go
type HandlerFunc func(ResponseWriter, *Request)
```

ドキュメントから、`HandlerFunc`型がすでに`ServeHTTP`メソッドを実装していることがわかります。`PlayerServer`関数をこの型に型変換することで、必要な`Handler`を実装したことになります。

### `http.ListenAndServe(":5000"...)`

`ListenAndServe` は、リッスンするポートと `Handler` を受け取ります。問題がある場合、Webサーバーはエラーを返します。エラーの一例として、ポートがすでにリッスンされていることが考えられます。そのため、この呼び出しを `log.Fatal` でラップして、ユーザーのために、エラーをログに出力します。

これから行うのは、ハードコーディングされた値から離れるための前向きな変更を強制する _別の_ テストを書くことです。

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

> 確かに、どのプレーヤーがどのスコアを獲得するかを制御するために、何らかのストレージの概念が必要です。テストで値が非常に恣意的に見えるのは奇妙です。

できる限り小さなステップで進もうとしていることを忘れないでください。今は定数を壊そうとしているだけです。

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

このテストによって、実際にリクエストのURLを見て判断せざるを得なくなりました。頭の中ではプレーヤーのストアやインターフェースのことを心配していたかもしれませんが、次の論理的なステップは実は _ルーティング_ に関するもののようです。

ストアのコードから始めていたら、必要な変更の量はこれに比べて非常に大きくなっていたでしょう。 **これは最終目標に向けたより小さなステップであり、テストによって導かれたものです**。

今のところ、ルーティングライブラリを使いたいという誘惑には抵抗し、テストに合格するための最小のステップにとどめています。

`r.URL.Path`はリクエストのパスを返すので、[`strings.TrimPrefix`](https://golang.org/pkg/strings/#TrimPrefix)を使用して`/players/`を取り除き、要求されたプレーヤーを取得できます。それほど堅牢ではありませんが、とりあえずはこれで十分です。

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

しかし、まだ満足すべきではありません。サーバーがスコアを知っているのはしっくりきません。

リファクタリングによって、何をすべきかがかなり明確になりました。

スコアの計算をハンドラーの本体から関数`GetPlayerScore`に移動しました。ここが、インターフェースを使って関心事を分離するのに適切な場所のように感じます。

リファクタリングで切り出した関数を、代わりにインターフェースにしてみましょう。

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

かなりの量の変更をしたので、テストもアプリケーションもコンパイルできなくなっていることはわかっています。でも、リラックスして、コンパイラに導いてもらいながら一つずつ直していきましょう。

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

_まだ_ ストアを作ることは気にしていない点に注意してください。できるだけ早くコンパイルを通したいだけです。

まずはコンパイルが通るコードを、次にテストが通るコードを優先する習慣を身につけましょう。

コードがコンパイルできない状態で（スタブのストアのような）機能を追加すると、 _さらに多くの_ コンパイルの問題を抱え込むことになりかねません。

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

テストは成功し、見た目も良くなりました。ストアを導入したことで、コードの背後にある _意図_ がより明確になりました。_`PlayerStore`にこのデータがある_ ので、それを`PlayerServer`で使用すると次のレスポンスが得られるはずだ、ということを読者に伝えています。

### アプリケーションを実行します

テストが成功するようになったので、このリファクタリングを完了するために最後に必要なのは、アプリケーションが動作するかを確認することです。プログラムは起動するはずですが、`http://localhost:5000/players/Pepper`でサーバーにアクセスしようとすると、ひどいレスポンスが返ってきます。

これは、`PlayerStore`を渡していないためです。

実装を1つ作る必要がありますが、まだ意味のあるデータを保存していないので現時点では難しく、当面はハードコーディングするしかありません。

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

`POST`のシナリオのほうが「ハッピーパス」に近づけますが、すでにその文脈にいるので、まずは存在しないプレーヤーのシナリオに取り組むほうが簡単だと思います。残りは後で行います。

## 最初にテストを書く

存在しないプレーヤーのシナリオを既存のスイートに追加します

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

TDDの提唱者が「テストを通すための最小限のコードだけを書くようにしてください」と言うのを聞くと、あまりに杓子定規に感じて、思わずうんざりしてしまうことがあります。

しかし、このシナリオはその良い例です。私は（正しくないと知りながら）最低限のことだけを行いました。つまり、**すべてのレスポンス**に`StatusNotFound`を書き込んだのですが、それでもすべてのテストが成功しています！

**テストを通すために最低限のことだけを行うと、テストの抜けが浮き彫りになることがあります**。今回のケースでは、プレーヤーがストアに _存在する_ ときに`StatusOK`が返されるべきことをアサートしていません。

他の2つのテストを更新してステータスをアサートし、コードを修正します。

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

これで、最初の2つのテストは200ではなく404が返るために失敗します。そこで、スコアが0の場合にのみ not found を返すように`PlayerServer`を修正します。

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

ストアからスコアを取得できるようになったので、次は新しいスコアを保存できるようにするのが理にかなっています。

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

まずは、POSTで特定のルートにアクセスしたときに正しいステータスコードが返ることだけを確認しましょう。これにより、異なる種類のリクエストを受け付け、`GET /players/{name}`とは異なる方法で処理する機能を導き出せます。これがうまくいったら、ハンドラーとストアのやり取りについてアサーションを書き始めることができます。

## テストを実行してみます

```
=== RUN   TestStoreWins/it_returns_accepted_on_POST
    --- FAIL: TestStoreWins/it_returns_accepted_on_POST (0.00s)
        server_test.go:70: did not get correct status, got 404, want 202
```

## 成功させるのに十分なコードを書く

意図的に罪を犯していることを思い出してください。リクエストのメソッドに基づく`if`文で十分です。

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

ハンドラーが少しごちゃごちゃしてきました。コードを分割して読みやすくし、異なる機能を新しい関数に分離しましょう。

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

これにより、`ServeHTTP`のルーティングの側面が少し明確になり、保存に関する次の反復は`processWin`の内部だけで済むようになります。

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

まずは呼び出しの回数を確認するようにテストを拡張します。

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

特定の値ではなく呼び出しの回数だけをアサートしているので、最初の反復が少し小さくなります。

`RecordWin`を呼び出せるようにするには、インターフェースを変更して、`PlayerServer`が考える`PlayerStore`の定義を更新する必要があります。

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

テストを実行してみると、コンパイルは通るようになりますが、テストはまだ失敗しています。

`PlayerStore`に`RecordWin`があるので、`PlayerServer`内で呼び出すことができます。

```go
//server.go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
	p.store.RecordWin("Bob")
	w.WriteHeader(http.StatusAccepted)
}
```

テストを実行すると成功するはずです！もちろん、`"Bob"`は`RecordWin`に渡したい値ではないので、テストをさらに改良してみましょう。

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

`processWin`が`http.Request`を受け取るように変更して、URLを見てプレーヤーの名前を抽出できるようにしました。それができたら、正しい値で`store`を呼び出してテストに合格させることができます。

## リファクタリング♪

2か所で同じ方法でプレーヤー名を抽出しているので、このコードを少しDRYにすることができます。

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

テストは成功していますが、実際に動くソフトウェアにはなっていません。`main`を実行して意図したとおりにソフトウェアを使おうとしても、まだ`PlayerStore`をきちんと実装していないので動作しません。ですが、これで問題ありません。ハンドラーに焦点を当てることで、事前に設計しようとするのではなく、必要なインターフェースを特定できたのです。

`InMemoryPlayerStore`に対するテストを書き始めることも _できます_ が、これはプレーヤーのスコアを永続化するより堅牢な方法（つまりデータベース）を実装するまでの一時的なものにすぎません。

ここでは、`PlayerServer`と`InMemoryPlayerStore`の間の _統合テスト_ を書いて、機能を完成させます。これにより、`InMemoryPlayerStore`を直接テストしなくても、アプリケーションが動作していると確信できるという目標を達成できます。それだけでなく、データベースを使って`PlayerStore`を実装する段になったら、同じ統合テストでその実装をテストできます。

### 統合テスト

統合テストは、システムのより広い領域が機能することをテストするのに役立ちますが、次の点に注意する必要があります。

* 書くのが難しい
* 失敗したときに原因がわかりにくく（通常は統合テストのコンポーネント内のバグです）、そのため修正が難しくなることがある
* 実行に時間がかかる場合がある（データベースなどの「本物の」コンポーネントと一緒に使われることが多いため）

そのためにも、_テストピラミッド（The Test Pyramid）_ をリサーチしておくことをお勧めします。

## 最初にテストを書く

簡潔にするために、最終的にリファクタリングした統合テストをお見せします。

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
* 次に、`player`の勝利を3回記録するために3つのリクエストを送ります。このテストでは、ステータスコードはコンポーネントがうまく統合されているかどうかとは関係がないので、あまり気にしていません。
* その次のレスポンスは重要なので（変数`response`に格納しています）、これで`player`のスコアを取得しようとしています。

## テストを実行してみます

```
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:24: response body is wrong, got '123' want '3'
```

## 成功させるのに十分なコードを書く

ここでは少し自由にやらせてもらい、テストを書かずに、皆さんが落ち着かないと感じるかもしれないほど多くのコードを書きます。

_これは許されることです！_ 正しく動作していることを確認するテストはまだありますが、それは今作業している特定のユニット（`InMemoryPlayerStore`）を対象にしたものではありません。

このシナリオで行き詰まった場合は、失敗するテストの時点まで変更を元に戻し、`InMemoryPlayerStore`に対するより具体的な単体テストを書いて、解決策を導き出す助けにします。

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

統合テストに合格したので、あとは`main`を変更して`NewInMemoryPlayerStore()`を使用するだけです。

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

* `curl -X POST http://localhost:5000/players/Pepper` を何回か実行します。お好みでプレーヤー名を変えても構いません
* `curl http://localhost:5000/players/Pepper`でスコアを確認してください

すごい！ REST風のサービスを作成しました。さらに先へ進めるには、プログラムの実行が終わった後もスコアを保持できるように、データストアを選ぶ必要があるでしょう。

* ストアを選択してください （Bolt? Mongo? Postgres? File system?）
* `PostgresPlayerStore`に`PlayerStore`を実装させる
* その機能をTDDで実装し、確実に動作するようにする
* 統合テストに接続し、問題がないことを確認します
* 最後に`main`に接続します

## リファクタリング♪

あと少しです！このような並行処理のエラーを防ぐために、もうひと手間かけましょう

```
fatal error: concurrent map read and map write
```

ミューテックスを追加することで、特に`RecordWin`関数のカウンターに対して並行処理の安全性を確保します。ミューテックスの詳細については、同期（sync）の章をご覧ください。

## まとめ

### `http.Handler`

* このインターフェースを実装してWebサーバーを作成する
* 通常の関数を`http.Handler`に変換するには、`http.HandlerFunc`を使用します
* `httptest.NewRecorder`を使用して`ResponseWriter`として渡し、ハンドラーが送信する応答をスパイできるようにします
* `http.NewRequest`を使用して、システムに入ると予想されるリクエストを作成します

### インターフェース、モッキング、DI

* 小さなチャンクでシステムを繰り返し構築できます
* 実際のストレージを必要とせずにストレージを必要とするハンドラーを開発できます
* TDDで必要なインターフェースを導き出す

### 罪を犯して、リファクタリング（そして、ソースコントロールにコミットして）

* コンパイルに失敗したり、テストに失敗したりすることは、できるだけ早く脱出しなければならない赤信号として扱う必要があります。
* そのために必要なコードだけを書いてください。その後、リファクタリングをしてコードを良くしてください。
* コードがコンパイルされていなかったり、テストが失敗している間に多くの変更をしようとすると、問題を悪化させる危険性があります。
* このアプローチを守ると、小さなテストを書かざるを得なくなります。それは小さな変更を意味し、複雑なシステムでの作業を管理しやすく保つのに役立ちます。
