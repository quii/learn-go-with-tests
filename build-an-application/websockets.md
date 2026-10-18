---
description: ウェブソケット
---

# ウェブソケット

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/websockets)

この章では、アプリケーションを改善するためにWebSocketを使用する方法を学びます。

## プロジェクトの要約

ポーカーコードベースには2つのアプリケーションがあります

* _コマンドラインアプリ_。ゲームのプレーヤー数を入力するようユーザーに求めます。それ以降は、時間とともに増加する「ブラインドベット」の値をプレーヤーに知らせます。ユーザーはいつでも`"{Playername} wins"`と入力してゲームを終了し、ストアに勝者を記録できます。
* _Webアプリ_。ユーザーがゲームの勝者を記録し、リーグテーブルを表示できるようにします。コマンドラインアプリと同じストアを共有します。

## 次のステップ

プロダクトオーナーはコマンドラインアプリケーションに大喜びしていますが、その機能をブラウザでも使えるようにしてほしいと考えています。彼女が思い描いているのは、ユーザーがプレーヤーの数を入力できるテキストボックスを備えたWebページで、フォームを送信するとページにブラインドの値が表示され、適切なタイミングで自動的に更新されるというものです。コマンドラインアプリケーションと同様に、ユーザーは勝者を宣言でき、それはデータベースに保存されます。

一見、非常に単純に聞こえますが、いつものように、ソフトウェアを書く際には _反復的_ なアプローチを取ることを強調しなければなりません。

まず、HTMLを提供する必要があります。
これまでのところ、すべてのHTTPエンドポイントはプレーンテキストかJSONを返してきました。（最終的にはすべて文字列なので）すでに知っている手法を使うことも _できます_ が、[html/template](https://golang.org/pkg/html/template/)パッケージを使えば、よりクリーンに解決できます。

また、ブラウザを更新しなくても、`The blind is now *y*`（ブラインドが *y* になりました）というメッセージを非同期でユーザーに送信できる必要があります。これを実現するために、[WebSockets](https://en.wikipedia.org/wiki/WebSocket)を使用できます。

> WebSocketはコンピューター通信プロトコルであり、単一のTCP接続で全二重通信チャネルを提供します。

いくつもの手法に取り組むことになるので、まずは役に立つ最小限の作業を行い、それから反復していくことがいっそう重要です。

そのため、最初に行うことは、ユーザーが勝者を記録するためのフォームを含むWebページを作成することです。
プレーンフォームを使用するのではなく、WebSocketを使用して、そのデータをサーバーに送信して記録します。

その後でブラインドアラートに取り組みます。
その頃には、インフラストラクチャのコードがある程度整っているでしょう。

### JavaScriptのテストについてはどうですか？

これを実現するためにJavaScriptをいくらか書きますが、そのテストを書くことまではしません。

もちろんテストを書くことは可能ですが、簡潔さを優先して、その説明は省きます。

皆さん、ごめんなさい。
私にお金を払って「Learn JavaScript with tests」を書かせるよう、O'Reillyに働きかけてください。

## 最初にテストを書く

まず、ユーザーが`/game`にアクセスしたときに、HTMLを返す必要があります。

おさらいとして、Webサーバーの関連するコードを示します。

```go
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}
```

今のところできる _最も簡単な_ ことは、`GET /game`したときに`200`が返ることを確認することです。

```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}
```

## テストを実行してみます

```
--- FAIL: TestGame (0.00s)
=== RUN   TestGame/GET_/game_returns_200
    --- FAIL: TestGame/GET_/game_returns_200 (0.00s)
    	server_test.go:109: did not get correct status, got 404, want 200
```

## 成功させるのに十分なコードを書く

私たちのサーバーにはルーターの設定があるので、修正は比較的簡単です。

私たちのルーターに追加します。

```go
router.Handle("/game", http.HandlerFunc(p.game))
```

そして`game`メソッドを書きます

```go
func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
```

## リファクタリング♪

サーバーのコードは、うまく構成された既存のコードにすんなりとコードを追加できたので、すでに問題ありません。

`/game`にリクエストを送信するためのテストヘルパー関数`newGameRequest`を追加することで、テストを少し片付けることができます。これを自分で書いてみてください。

```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})
}
```

また、そのほうが読みやすいと感じたので、`assertStatus`が`response.Code`ではなく`response`を受け取るように変更したことにも気づくでしょう。

次に、エンドポイントがHTMLを返すようにする必要があります。HTMLは次のとおりです。

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Let's play poker</title>
</head>
<body>
<section id="game">
    <div id="declare-winner">
        <label for="winner">Winner</label>
        <input type="text" id="winner"/>
        <button id="winner-button">Declare winner</button>
    </div>
</section>
</body>
<script type="application/javascript">

    const submitWinnerButton = document.getElementById('winner-button')
    const winnerInput = document.getElementById('winner')

    if (window['WebSocket']) {
        const conn = new WebSocket('ws://' + document.location.host + '/ws')

        submitWinnerButton.onclick = event => {
            conn.send(winnerInput.value)
        }
    }
</script>
</html>
```

とてもシンプルなウェブページがあります。

* ユーザーが勝者を入力するためのテキスト入力
* 勝者を宣言するためにクリックするボタン
* サーバーへのWebSocket接続を開き、送信ボタンが押されたときの処理を行うJavaScript

`WebSocket`はほとんどの最新のブラウザに組み込まれているので、ライブラリを導入する心配はありません。古いブラウザではこのWebページは動作しませんが、今回のシナリオではそれで問題ありません。

### 正しいマークアップを返すことをテストするにはどうすればよいですか？

いくつかの方法があります。本全体を通して強調されているように、あなたが書くテストはコストを正当化するのに十分な価値を持つことが重要です。

1. Seleniumなどを使用して、ブラウザベースのテストを記述します。これらのテストは、ある種の実際のWebブラウザーを起動して、ユーザーとの対話をシミュレートするため、すべてのアプローチの中で最も「現実的な」ものです。これらのテストは、システムの動作に大きな自信を与えることができますが、単体テストよりも作成が難しく、実行速度がはるかに遅くなります。私たちの製品の目的上、これはやりすぎです。
2. 文字列が完全に一致するかを確認します。これでも _問題ない_ 場合はありますが、この種のテストは非常に壊れやすくなりがちです。誰かがマークアップを変更した途端、実際には何も _壊れていない_ のにテストが失敗することになります。
3. 正しいテンプレートを呼び出していることを確認します。HTMLを提供するために標準ライブラリのテンプレートライブラリを使用する予定です（すぐに説明します）。HTMLを生成する _もの_ を注入し、その呼び出しをスパイして、正しく処理していることを確認できます。これはコードの設計に影響を与えますが、実際にはそれほど多くのことをテストしていません。確認できるのは、正しいテンプレートファイルで呼び出しているかどうかくらいです。プロジェクトにはテンプレートが1つしかないので、ここで失敗する可能性は低そうです。

そこで、「Learn Go with Tests」という本の中で初めて、テストを書かないことにします。

マークアップを`game.html`というファイルに入れます

次に、先ほど書いたエンドポイントを次のように変更します

```go
func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("game.html")

	if err != nil {
		http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
```

[`html/template`](https://golang.org/pkg/html/template/)は、HTMLを作成するためのGoパッケージです。
私たちの場合、`template.ParseFiles`を呼び出して、htmlファイルのパスを指定します。
エラーがなければ、テンプレートを`Execute`（実行）できます。これにより、テンプレートが`io.Writer`に書き込まれます。今回はインターネットに`Write`（書き込み）したいので、`http.ResponseWriter`を渡します。

テストを作成していないので、希望どおりに機能していることを確認するためだけにWebサーバーを手動でテストするのが賢明です。 `cmd/webserver`に移動し、`main.go`ファイルを実行します。`http://localhost:5000/game`にアクセスしてください。

テンプレートが見つからないというエラーが発生した _はず_ です。パスを自分のフォルダーからの相対パスに変更するか、`cmd/webserver`ディレクトリに`game.html`のコピーを置くことができます。私は、プロジェクトのルートにあるファイルへのシンボリックリンク（`ln -s ../../game.html game.html`）を作成することにしました。こうすれば、ファイルに変更を加えたときに、サーバーの実行時にも反映されます。

この変更を加えて再度実行すると、UIが表示されます。

次に、サーバーへのWebSocket接続を介して文字列を受け取ったときに、それをゲームの勝者として宣言することをテストする必要があります。

## 最初にテストを書く

WebSocketを扱えるようにするために、初めて外部ライブラリを使用します。

`go get github.com/gorilla/websocket`を実行します

これにより、優れた[Gorilla WebSocket](https://github.com/gorilla/websocket) ライブラリのコードがフェッチされます。これで、新しい要件に合わせてテストを更新できます。

```go
t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
	store := &StubPlayerStore{}
	winner := "Ruth"
	server := httptest.NewServer(NewPlayerServer(store))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}

	AssertPlayerWin(t, store, winner)
})
```

`websocket`ライブラリのインポートがあることを確認してください。
私のIDEは自動的にやってくれました。皆さんのIDEもそうなるはずです。

ブラウザから起こることをテストするには、自分でWebSocket接続を開いて、そこに書き込む必要があります。

サーバーに関するこれまでのテストではサーバーのメソッドを呼び出すだけでしたが、今回はサーバーへの永続的な接続が必要です。そのためには、`http.Handler`を受け取り、それを起動して接続を待ち受ける`httptest.NewServer`を使用します。

`websocket.DefaultDialer.Dial`を使用してサーバーに接続し、`winner`を含むメッセージを送信します。

最後に、勝者が記録されたことを確認するために、プレーヤーストアに対してアサーションを行います。

## テストを実行してみます

```
=== RUN   TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game
    --- FAIL: TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game (0.00s)
        server_test.go:124: could not open a ws connection on ws://127.0.0.1:55838/ws websocket: bad handshake
```

`/ws`でWebSocket接続を受け入れるようにサーバーを変更していないため、まだハンドシェイクが行われていません。

## 成功させるのに十分なコードを書く

ルーターにもう1つルートを追加します

```go
router.Handle("/ws", http.HandlerFunc(p.webSocket))
```

次に、新しい`webSocket`ハンドラーを追加します

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.Upgrade(w, r, nil)
}
```

WebSocket接続を受け入れるには、リクエストを`Upgrade`します。ここでテストを再実行すると、次のエラーに進むはずです。

```
=== RUN   TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game
    --- FAIL: TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game (0.00s)
        server_test.go:132: got 0 calls to RecordWin want 1
```

接続を開いたので、メッセージをリッスンして、それを勝者として記録します。

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _ := upgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.store.RecordWin(string(winnerMsg))
}
```

（はい、現在、多くのエラーを無視しています！）

`conn.ReadMessage()`は、接続にメッセージが届くまでブロックして待機します。メッセージを受け取ったら、それを使って`RecordWin`します。これで最終的にWebSocket接続が閉じられます。

テストを実行すると、まだ失敗します。

問題はタイミングです。WebSocket接続がメッセージを読み取って勝利を記録するまでには遅延があり、それが起こる前にテストが終了してしまいます。これは、最後のアサーションの前に短い`time.Sleep`を入れることで確かめられます。

とりあえずはそれで進めましょう。ただし、テストに任意のスリープを入れるのは**非常に悪い習慣**であることは認識しておいてください。

```go
time.Sleep(10 * time.Millisecond)
AssertPlayerWin(t, store, winner)
```

## リファクタリング♪

このテストをサーバーコードとテストコードの両方で機能させるために多くの罪を犯しましたが、これが私たちにとって最も簡単な方法であることを覚えておいてください。

テストに裏打ちされた、ひどくて汚いけれど _動く_ ソフトウェアができました。これで、うっかり何かを壊してしまうことはないと分かったうえで、自由にきれいにしていけます。

サーバーコードから始めましょう。

すべてのWebSocket接続リクエストで再宣言する必要がないため、`upgrader`をパッケージ内のプライベート値に移動できます。

```go
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.store.RecordWin(string(winnerMsg))
}
```

`template.ParseFiles("game.html")`への呼び出しは、すべての `GET /game`で実行されます。
つまり、テンプレートを再解析する必要がないにもかかわらず、リクエストのたびにファイルシステムにアクセスすることになります。代わりに`NewPlayerServer`でテンプレートを一度だけ解析するようにリファクタリングしましょう。
ディスクからテンプレートを取得したり解析したりする際に問題が発生した場合に備えて、この関数がエラーを返せるようにする必要があります。

`PlayerServer`に関連する変更は次のとおりです

```go
type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store PlayerStore) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.game))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}
```

`NewPlayerServer`のシグネチャを変更したので、コンパイルエラーが発生します。自分で修正してみるか、うまくいかない場合はソースコードを参照してください。

テストコードでは、 `mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer`というヘルパーを作成して、エラーノイズをテストから隠せるようにしました。

```go
func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}
```

同様に、別のヘルパー`mustDialWS`を作成して、WebSocket接続の作成時に厄介なエラーノイズを隠すことができるようにしました。

```go
func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}
```

最後に、テストコードで、メッセージの送信を整理するヘルパーを作成できます。

```go
func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
```

テストが成功するようになったので、サーバーを起動して`/game`で何人か勝者を宣言してみてください。`/league`に記録されているはずです。勝者を受け取るたびに _接続を閉じている_ ので、再び接続を開くにはページを更新する必要があることを忘れないでください。

ユーザーがゲームの勝者を記録できる簡単なWebフォームを作成しました。これを反復して、ユーザーがプレーヤーの数を入力してゲームを開始でき、時間の経過とともにサーバーがブラインドの値をクライアントにプッシュして知らせるようにしましょう。

まず、`game.html`を更新して、新しい要件に合わせてクライアント側のコードを更新します

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Lets play poker</title>
</head>
<body>
<section id="game">
    <div id="game-start">
        <label for="player-count">Number of players</label>
        <input type="number" id="player-count"/>
        <button id="start-game">Start</button>
    </div>

    <div id="declare-winner">
        <label for="winner">Winner</label>
        <input type="text" id="winner"/>
        <button id="winner-button">Declare winner</button>
    </div>

    <div id="blind-value"/>
</section>

<section id="game-end">
    <h1>Another great game of poker everyone!</h1>
    <p><a href="/league">Go check the league table</a></p>
</section>

</body>
<script type="application/javascript">
    const startGame = document.getElementById('game-start')

    const declareWinner = document.getElementById('declare-winner')
    const submitWinnerButton = document.getElementById('winner-button')
    const winnerInput = document.getElementById('winner')

    const blindContainer = document.getElementById('blind-value')

    const gameContainer = document.getElementById('game')
    const gameEndContainer = document.getElementById('game-end')

    declareWinner.hidden = true
    gameEndContainer.hidden = true

    document.getElementById('start-game').addEventListener('click', event => {
        startGame.hidden = true
        declareWinner.hidden = false

        const numberOfPlayers = document.getElementById('player-count').value

        if (window['WebSocket']) {
            const conn = new WebSocket('ws://' + document.location.host + '/ws')

            submitWinnerButton.onclick = event => {
                conn.send(winnerInput.value)
                gameEndContainer.hidden = false
                gameContainer.hidden = true
            }

            conn.onclose = evt => {
                blindContainer.innerText = 'Connection closed'
            }

            conn.onmessage = evt => {
                blindContainer.innerText = evt.data
            }

            conn.onopen = function () {
                conn.send(numberOfPlayers)
            }
        }
    })
</script>
</html>
```

主な変更点は、プレーヤー数を入力するセクションとブラインドの値を表示するセクションを追加したことです。ゲームの段階に応じて、ユーザーインターフェースを表示/非表示にするちょっとしたロジックもあります。

`conn.onmessage`を介して受信するメッセージはすべてブラインドアラートであると想定するため、それに応じて`blindContainer.innerText`を設定します。

ブラインドアラートはどのように送信すればよいでしょうか？前の章では`Game`という考え方を導入し、CLIのコードが`Game`を呼び出すだけで、ブラインドアラートのスケジューリングを含むそれ以外のすべてが処理されるようにしました。これは良い関心の分離であることがわかりました。

```go
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}
```

CLIでは、ユーザーがプレーヤーの数を入力するとゲームを`Start`してブラインドアラートを開始し、ユーザーが勝者を宣言すると`Finish`していました。今回の要件もこれと同じで、入力の受け取り方が違うだけです。ですから、可能であればこの概念を再利用すべきです。

`Game`の「本物の」実装は`TexasHoldem`です。

```go
type TexasHoldem struct {
	alerter BlindAlerter
	store   PlayerStore
}
```

`BlindAlerter`を渡すことで、`TexasHoldem`はブラインドアラートを _どこにでも_ 送信するようにスケジュールできます

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}
```

また、注意点として、CLIで使用する `BlindAlerter` の実装を以下に示します。

```go
func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
```

これがCLIで機能するのは、 _常にアラートを`os.Stdout`に送信したい_ からですが、Webサーバーではうまくいきません。リクエストごとに新しい`http.ResponseWriter`を受け取り、それを`*websocket.Conn`にアップグレードします。
そのため、依存関係を構築する時点では、アラートをどこに送る必要があるのかを知ることができません。

そのため、Webサーバーで再利用できるように、`BlindAlerter.ScheduleAlertAt`がアラートの宛先を受け取るように変更する必要があります。

`blind_alerter.go`を開き、パラメータ`to io.Writer`を追加します

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}
```

`StdoutAlerter`のアイデアは新しいモデルに適合しないため、名前を`Alerter`に変更します

```go
func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}
```

コンパイルしようとすると、`TexasHoldem`が宛先なしで`ScheduleAlertAt`を呼び出しているために失敗します。 _ひとまず_ コンパイルを通すために、`os.Stdout`をハードコードしておきましょう。

テストを実行してみると、`SpyBlindAlerter`がもはや`BlindAlerter`を実装していないために失敗します。`ScheduleAlertAt`のシグネチャを更新してこれを修正し、テストを実行すれば、グリーンのままのはずです。

`TexasHoldem`がブラインドアラートの送信先を知っているのは筋が通りません。
そこで`Game`を更新して、ゲームを開始するときにアラートを _どこに_ 送るかを指定するようにしましょう。

```go
type Game interface {
	Start(numberOfPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}
```

何を修正する必要があるかはコンパイラに教えてもらいましょう。変更はそれほど大変ではありません。

* `TexasHoldem`を更新して、`Game`を適切に実装します。
* `CLI`でゲームを開始するときに、`out`プロパティを渡します（`cli.game.Start(numberOfPlayers, cli.out)`）
* `TexasHoldem`のテストでは、`game.Start(5, io.Discard)`を使用してコンパイルの問題を修正し、アラート出力が破棄されるように設定します

すべて正しくできていれば、すべてグリーンになるはずです！
これで、`Server`の中で`Game`を使ってみることができます。

## 最初にテストを書く

`CLI`と`Server`の要件は同じです！届け方が違うだけです。

インスピレーションを得るために、 `CLI`テストを見てみましょう。

```go
t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
	game := &GameSpy{}

	out := &bytes.Buffer{}
	in := userSends("3", "Chris wins")

	poker.NewCLI(in, out, game).PlayPoker()

	assertMessagesSentToUser(t, out, poker.PlayerPrompt)
	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, "Chris")
})
```

`GameSpy`を使えば、同じような結果をテスト駆動で導き出せそうです

古いwebsocketテストを次のものに置き換えます

```go
t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
	game := &poker.GameSpy{}
	winner := "Ruth"
	server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

	defer server.Close()
	defer ws.Close()

	writeWSMessage(t, ws, "3")
	writeWSMessage(t, ws, winner)

	time.Sleep(10 * time.Millisecond)
	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, winner)
})
```

* 議論したように、スパイの`Game`を作成して`mustMakePlayerServer`に渡します（これをサポートするためにヘルパーを更新してください）
* 次に、ゲームのためのWebSocketメッセージを送信します。
* 最後に、ゲームが期待どおりに開始・終了したことをアサートします。

## テストを実行してみます

他のテストでは、`mustMakePlayerServer`の周りにいくつかのコンパイルエラーが発生します。エクスポートされていない変数`dummyGame`を導入し、コンパイルされていないすべてのテストで使用します

```go
var (
	dummyGame = &GameSpy{}
)
```

最後のエラーは、`Game`を`NewPlayerServer`に渡そうとしているところですが、まだサポートしていません。

```
./server_test.go:21:38: too many arguments in call to "github.com/quii/learn-go-with-tests/WebSockets/v2".NewPlayerServer
	have ("github.com/quii/learn-go-with-tests/WebSockets/v2".PlayerStore, "github.com/quii/learn-go-with-tests/WebSockets/v2".Game)
	want ("github.com/quii/learn-go-with-tests/WebSockets/v2".PlayerStore)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

今のところ、テストを実行するために引数として追加するだけです

```go
func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error)
```

ようやくです！

```
=== RUN   TestGame/start_a_game_with_3_players_and_declare_Ruth_the_winner
--- FAIL: TestGame (0.01s)
    --- FAIL: TestGame/start_a_game_with_3_players_and_declare_Ruth_the_winner (0.01s)
    	server_test.go:146: wanted Start called with 3 but got 0
    	server_test.go:147: expected finish called with 'Ruth' but got ''
FAIL
```

## 成功させるのに十分なコードを書く

リクエストを受け取ったときに使用できるように、フィールドとして`Game`を`PlayerServer`に追加する必要があります。

```go
type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
	game     Game
}
```

（すでに`game`というメソッドがあるので、名前を`playGame`に変更します）

次に、コンストラクタで代入しましょう

```go
func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.game = game

	// etc
}
```

これで、`webSocket`内で`Game`を使用できます。

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.game.Start(numberOfPlayers, io.Discard) //todo: Don't discard the blinds messages!

	_, winner, _ := conn.ReadMessage()
	p.game.Finish(string(winner))
}
```

やったー！テストに成功しました。

ブラインドメッセージの送信先については少し考える必要があるので、 _まだ_ どこにも送信しません。`game.Start`を呼び出すときには、書き込まれたメッセージを単に破棄する`io.Discard`を渡します。

とりあえずWebサーバーを起動してみましょう。`Game`を`PlayerServer`に渡すには、`main.go`を更新する必要があります

```go
func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)

	server, err := poker.NewPlayerServer(store, game)

	if err != nil {
		log.Fatalf("problem creating player server %v", err)
	}

	log.Fatal(http.ListenAndServe(":5000", server))
}
```

ブラインドアラートをまだ受け取れていないことを除けば、アプリは動作します！`PlayerServer`で`Game`を再利用することに成功し、`Game`がすべての詳細を処理してくれました。ブラインドアラートを破棄するのではなくWebSocketに送信する方法がわかれば、すべて動作する _はず_ です。

その前に、いくつかのコードを整理しましょう。

## リファクタリング♪

WebSocketの使い方はかなり基本的で、エラー処理もかなり素朴なので、サーバーのコードからその煩雑さを取り除くためだけに、型にカプセル化したいと思いました。後で見直すことになるかもしれませんが、今のところはこれで少し整理されます

```go
type playerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}
```

これでサーバーのコードが少しシンプルになりました

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, io.Discard) //todo: Don't discard the blinds messages!

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
```

ブラインドメッセージを破棄しない方法を見つけたら、完了です。

### テストを _書かない_ ことにしよう！

どうすればよいかわからないときは、とにかくいろいろ試してみるのが一番です！ただし、まずは作業内容をコミットしておいてください。方法がわかったら、それをテストで駆動するべきだからです。

問題となっているコードの行は次のとおりです

```go
p.game.Start(numberOfPlayers, io.Discard) //todo: Don't discard the blinds messages!
```

ゲームがブラインドアラートを書き込むには、`io.Writer`を渡す必要があります。

先ほどの`playerServerWS`を渡せたら素敵だと思いませんか？これはWebSocketのラッパーなので、`Game`に渡してメッセージを送ってもらえそうな _気がします_ 。

試してみましょう。

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ws)
	//etc...
}
```

コンパイラは以下のような問題を指摘しています。

```
./server.go:71:14: cannot use ws (type *playerServerWS) as type io.Writer in argument to p.game.Start:
	*playerServerWS does not implement io.Writer (missing Write method)
```

当然の対応は、`playerServerWS`が`io.Writer`を _実装する_ ようにすることでしょう。そのためには、内部の`*websocket.Conn`の`WriteMessage`を使って、メッセージをWebSocketに送信します。

```go
func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
```

簡単すぎるように思えます！アプリケーションを実行して、動作するか確認してみましょう。

その前に`TexasHoldem`を編集してブラインドの増加間隔を短くし、動作を確認できるようにしてください

```go
blindIncrement := time.Duration(5+numberOfPlayers) * time.Second // (rather than a minute)
```

動作しているのが確認できるはずです！ブラインドの額がまるで魔法のようにブラウザ上で増えていきます。

では、コードを元に戻して、テスト方法を考えましょう。これを _実装する_ ためにやったことは、`StartGame`に`io.Discard`ではなく`playerServerWS`を渡しただけなので、その呼び出しをスパイして動作を検証すべきだと思うかもしれません。

スパイは素晴らしく、実装の詳細を確認するのに役立ちますが、可能であれば常に _本物の_ 振る舞いをテストすることを優先すべきです。リファクタリングしようとしたときに失敗し始めるのは、たいていスパイを使ったテストです。スパイのテストは通常、変更しようとしている実装の詳細を確認しているからです。

現在のテストは、実行中のサーバーへのWebSocket接続を開き、メッセージを送信してサーバーに処理をさせています。同様に、サーバーがWebSocket接続を介して送り返すメッセージもテストできるはずです。

## 最初にテストを書く

既存のテストを編集します。

現在、`GameSpy`は`Start`が呼ばれても`out`に何もデータを送信しません。あらかじめ用意したメッセージを送信するよう設定できるように変更すれば、そのメッセージがWebSocketに送信されることを確認できます。これにより、望んでいる本物の振る舞いを確かめつつ、正しく設定できていることに自信が持てるはずです。

```go
type GameSpy struct {
	StartCalled     bool
	StartCalledWith int
	BlindAlert      []byte

	FinishedCalled   bool
	FinishCalledWith string
}
```

`BlindAlert`フィールドを追加します。

`GameSpy`の`Start`を更新して、あらかじめ用意したメッセージを`out`に送信するようにします。

```go
func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartCalledWith = numberOfPlayers
	out.Write(g.BlindAlert)
}
```

これで、`PlayerServer`を動かしてゲームを`Start`しようとすると、正しく動作していればWebSocketを介してメッセージが送信されることになります。

最後にテストを更新できます

```go
t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
	wantedBlindAlert := "Blind is 100"
	winner := "Ruth"

	game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
	server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

	defer server.Close()
	defer ws.Close()

	writeWSMessage(t, ws, "3")
	writeWSMessage(t, ws, winner)

	time.Sleep(10 * time.Millisecond)
	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, winner)

	_, gotBlindAlert, _ := ws.ReadMessage()

	if string(gotBlindAlert) != wantedBlindAlert {
		t.Errorf("got blind alert %q, want %q", string(gotBlindAlert), wantedBlindAlert)
	}
})
```

* `wantedBlindAlert`を追加し、`Start`が呼び出された場合に`out`に送信するように`GameSpy`を設定しました。
* メッセージがWebSocket接続で送信されることを期待して、`ws.ReadMessage()`の呼び出しを追加し、メッセージが送信されるのを待ってから、それが期待どおりのものかを確認します。

## テストを実行してみます

テストが永遠にハングしてしまうはずです。これは、`ws.ReadMessage()`がメッセージを受け取るまでブロックするためで、メッセージは決して届きません。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

テストがハングするようなことは決してあってはならないので、タイムアウトさせたいコードを扱う方法を導入しましょう。

```go
func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}
```

`within`は、関数`assert`を引数として受け取り、それをゴルーチンで実行します。関数が終了した場合（終了したとき）には、`done`チャネルを通じて完了したことを通知します。

その間、チャネルがメッセージを送信するのを待てる`select`文を使います。ここからは、`assert`関数と、指定した時間が経過したときにシグナルを送る`time.After`との競争になります。

最後に、少しすっきりさせるために、アサーション用のヘルパー関数を作りました

```go
func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}
```

テストは次のようになりました

```go
t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
	wantedBlindAlert := "Blind is 100"
	winner := "Ruth"

	game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
	server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

	defer server.Close()
	defer ws.Close()

	writeWSMessage(t, ws, "3")
	writeWSMessage(t, ws, winner)

	time.Sleep(tenMS)

	assertGameStartedWith(t, game, 3)
	assertFinishCalledWith(t, game, winner)
	within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
})
```

テストを実行すると...

```
=== RUN   TestGame
=== RUN   TestGame/start_a_game_with_3_players,_send_some_blind_alerts_down_WS_and_declare_Ruth_the_winner
--- FAIL: TestGame (0.02s)
    --- FAIL: TestGame/start_a_game_with_3_players,_send_some_blind_alerts_down_WS_and_declare_Ruth_the_winner (0.02s)
    	server_test.go:143: timed out
    	server_test.go:150: got "", want "Blind is 100"
```

## 成功させるのに十分なコードを書く

最後に、ゲームの開始時にWebSocket接続をゲームに渡すように、サーバーのコードを変更します

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ws)

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
```

## リファクタリング♪

サーバーコードは非常に小さな変更でしたので、ここで変更することは多くありませんが、サーバーが非同期で動作するのを待つ必要があるため、テストコードにはまだ`time.Sleep`呼び出しがあります。

ヘルパーの`assertGameStartedWith`と`assertFinishCalledWith`をリファクタリングして、失敗とみなす前に短い間アサーションを再試行できるようにします。

`assertFinishCalledWith`でこれを行う方法を次に示します。
他のヘルパーにも同じアプローチを使用できます。

```go
func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishCalledWith == winner
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}
```

`retryUntil`の定義方法は次のとおりです

```go
func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
```

## まとめ

これでアプリケーションが完成しました。ポーカーのゲームはWebブラウザから開始でき、ユーザーは時間の経過とともにWebSocketを介してブラインドベットの値を知らされます。ゲームが終了すると勝者を記録でき、それは数章前に書いたコードを使って永続化されます。
プレーヤーは、Webサイトの`/league`エンドポイントを使って、誰が最高の（あるいは最も幸運な）ポーカープレーヤーかを知ることができます。

道中では間違いも犯しましたが、TDDの流れのおかげで、動くソフトウェアから大きく離れることはありませんでした。自由に反復や実験を続けることができました。

最後の章では、このアプローチと、たどり着いた設計を振り返り、やり残したことを片付けます。

この章ではいくつかのことを説明しました

### WebSockets

* クライアントがサーバーをポーリングし続ける必要がないクライアントとサーバー間でメッセージを送信する便利な方法。私たちが持っているクライアントとサーバーのコードはどちらも非常に単純です。
* テストは簡単ですが、テストの非同期の性質に注意する必要があります

### 遅延したり終了しなかったりする可能性のあるコードをテストで扱う

* アサーションを再試行し、タイムアウトを追加するヘルパー関数を作成します。
* ゴルーチンを使用して、アサーションが何もブロックしないようにし、チャネルを使用して、終了したかどうかを通知できます。
* `time`パッケージには、時間に関するイベントをチャネル経由で通知する便利な関数がいくつかあり、それを使ってタイムアウトを設定できます
