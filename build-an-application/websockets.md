---
description: WebSockets
---

# ウェブソケット

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/websockets)

この章では、アプリケーションを改善するためにWebSocketを使用する方法を学びます。

## プロジェクトの要約

ポーカーコードベースには2つのアプリケーションがあります

* _コマンドラインアプリ_。ゲームのプレーヤー数を入力するようにユーザーに求めます。それ以降は、「ブラインドベット」の値が何であるかをプレイヤーに知らせます。ユーザーはいつでも`" {Playername} wins"`を入力してゲームを終了し、ストアに勝利者を記録できます。
* _Webアプリ_。ユーザーがゲームの勝者を記録し、リーグテーブルを表示できるようにします。コマンドラインアプリと同じストアを共有します。

## 次のステップ

製品の所有者はコマンドラインアプリケーションに興奮していますが、ブラウザーにその機能を提供できればそれを好むでしょう。彼女は、ユーザーがプレーヤーの数を入力できるテキストボックスを備えたWebページを想像し、ユーザーがフォームを送信すると、ページにブラインド値が表示され、必要に応じて自動的に更新されます。コマンドラインアプリケーションと同様に、ユーザーは勝者を宣言でき、データベースに保存されます。

一見、それは非常に単純に聞こえますが、いつものように、ソフトウェアを書くために反復的アプローチを取ることを強調しなければなりません。

まず、HTMLを提供する必要があります。
これまでのところ、すべてのHTTPエンドポイントはプレーンテキストまたはJSONを返しています。（最終的にはすべて文字列なので）知っているのと同じ手法を使用できますが、よりクリーンな[html/template](https://golang.org/pkg/html/template/)パッケージを使用することもできます解決。

また、ブラウザーを更新せずに、「`ブラインドが *y* になりました`」というメッセージを非同期でユーザーに送信できる必要もあります。これを容易にするために、[WebSockets](https://en.wikipedia.org/wiki/WebSocket)を使用できます。

> WebSocketはコンピューター通信プロトコルであり、単一のTCP接続で全二重通信チャネルを提供します。

多くの手法を採用しているため、最初に可能な限り最小限の有用な作業を行ってから繰り返すことがさらに重要です。

そのため、最初に行うことは、ユーザーが勝者を記録するためのフォームを含むWebページを作成することです。
プレーンフォームを使用するのではなく、WebSocketを使用して、そのデータをサーバーに送信して記録します。

その後、ブラインドアラートに取り組みます。
その時点で、インフラストラクチャコードが少し設定されます。

### JavaScriptのテストについてはどうですか？

これを行うために記述されたJavaScriptがいくつかありますが、テストの記述まで行いません。

ごめんなさい。
O'Reillyに「テストでJavaScriptを学ぶ」を作ってくれるようにロビーをお願いします。

## 最初にテストを書く

まず、ユーザーが「`/game`」を押したときにHTMLを提供する必要があります。

ここに私たちのウェブサーバーの関連コードのリマインダーがあります。

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

ここでできる最も簡単なことは、「`GET /game`」を実行したときに「`200`」が得られることを確認することです。

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

サーバーコードは、既存の適切に因数分解されたコードに非常に簡単にコードを追加できるため、既に問題ありません。

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

また、「`assertStatus`」を変更して、「`response.Code`」ではなく「`response`」を受け入れるように変更しました。

次に、エンドポイントにHTMLを返すようにする必要があります。

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
* 彼らが勝者を宣言するためにクリックできるボタン。
* サーバーへのWebSocket接続を開き、送信ボタンが押されたときに処理するためのJavaScriptをいくつか用意しました。

`WebSocket`はほとんどの最新のブラウザーに組み込まれているので、ライブラリーの持ち込みについて心配する必要はありません。 Webページは古いブラウザーでは機能しませんが、このシナリオではそれで問題ありません。

### 正しいマークアップを返すことをテストするにはどうすればよいですか？

いくつかの方法があります。本全体を通して強調されているように、あなたが書くテストはコストを正当化するのに十分な価値を持つことが重要です。

1. Seleniumなどを使用して、ブラウザベースのテストを記述します。これらのテストは、ある種の実際のWebブラウザーを起動して、ユーザーとの対話をシミュレートするため、すべてのアプローチの中で最も「現実的な」ものです。これらのテストは、システムの動作に大きな自信を与えることができますが、単体テストよりも作成が難しく、実行速度がはるかに遅くなります。私たちの製品の目的上、これはやりすぎです。
2. 文字列を完全に一致させます。これは大丈夫ですが、この種のテストは非常に壊れやすくなります。誰かがマークアップを変更した瞬間に、実際には何も実際に壊れていない場合、テストは失敗します。
3. 正しいテンプレートを呼び出すことを確認します。標準のlibのテンプレートライブラリを使用してHTMLを提供します（すぐに説明します）。_thing_ を挿入してHTMLを生成し、その呼び出しをスパイして、正しく実行されていることを確認できます。これはコードの設計に影響を与えますが、実際にはあまりテストしていません。正しいテンプレートファイルを使用して呼び出す場合を除きます。プロジェクトにはテンプレートが1つしかないため、ここで失敗する可能性は低いようです。

そこで、「テスト駆動開発でGo言語を学びましょう」というサイトの中で、初めてテストを書くことになります。

マークアップを「`game.html`」というファイルに入れます

次に、先ほど書き込んだエンドポイントを次のように変更します

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
エラーがなければ、テンプレートを「実行`Execute`」して、それを「`io.Writer`」に書き込みます。この例では、インターネットに「書き込み`Write`」したいので、「`http.ResponseWriter`」を指定します。

テストを作成していないので、希望どおりに機能していることを確認するためだけにWebサーバーを手動でテストするのが賢明です。 `cmd/webserver`に移動し、`main.go`ファイルを実行します。`http://localhost:5000/game`にアクセスしてください。

テンプレートが見つからないというエラーが発生するはずです。フォルダーに相対的なパスに変更するか、`cmd/webserver`ディレクトリに`game.html`のコピーを置くことができます。プロジェクトのルート内のファイルへのシンボリックリンク（`ln -s ../../game.html game.html`）を作成することを選択したので、変更を加えると、サーバーの実行時に変更が反映されます。

この変更を加えて再度実行すると、UIが表示されます。

次に、サーバーへのWebSocket接続を介して文字列を取得したときに、ゲームの勝者として宣言することをテストする必要があります。

## 最初にテストを書く

初めて、WebSocketで作業できるように外部ライブラリを使用します。

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
私のIDEはそれを自動的に行いました。

ブラウザーから何が起こるかをテストするには、独自のWebSocket接続を開いて、それに書き込む必要があります。

サーバーに関する以前のテストではサーバーのメソッドを呼び出しましたが、サーバーへの永続的な接続が必要です。そのためには、`http.Handler`を受け取り、それを起動して接続を待機する`httptest.NewServer`を使用します。

`websocket.DefaultDialer.Dial`を使用してサーバーにダイヤルインしてから、`winner`でメッセージを送信しようとします。

最後に、勝者が記録されたことを確認するためにプレイヤーストアを評価します。

## テストを実行してみます

```
=== RUN   TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game
    --- FAIL: TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game (0.00s)
        server_test.go:124: could not open a ws connection on ws://127.0.0.1:55838/ws websocket: bad handshake
```

`/ws`でWebSocket接続を受け入れるようにサーバーを変更していないため、まだ握手はしていません。

## 成功させるのに十分なコードを書く

ルーターに別のリストを追加する

```go
router.Handle("/ws", http.HandlerFunc(p.webSocket))
```

次に、新しい「`webSocket`」ハンドラを追加します

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.Upgrade(w, r, nil)
}
```

WebSocket接続を受け入れるには、リクエストを「`Upgrade`」します。ここでテストを再実行する場合は、次のエラーに進む必要があります。

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

`conn.ReadMessage()`は、接続でのメッセージの待機をブロックします。取得したら、それを `RecordWin`に使用します。これは最終的にWebSocket接続を閉じます。

テストを実行すると、まだ失敗します。

問題はタイミングです。メッセージを読み取るWebSocket接続と勝利の記録の間には遅延があり、テストが完了する前にテストが終了します。これをテストするには、最後のアサーションの前に短い「`time.Sleep`」を配置します。

とりあえずそれを続けましょうが、任意のスリープをテストに入れることは**非常に悪い習慣**であることを認めます。

```go
time.Sleep(10 * time.Millisecond)
AssertPlayerWin(t, store, winner)
```

## リファクタリング♪

このテストをサーバーコードとテストコードの両方で機能させるために多くの罪を犯しましたが、これが私たちにとって最も簡単な方法であることを覚えておいてください。

テストに裏打ちされた、厄介で恐ろしい _working_ ソフトウェアがあるので、これを自由に変更して、誤って何も壊さないようにすることができます。

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
つまり、テンプレートを再解析する必要がない場合でも、すべてのリクエストでファイルシステムに移動します。代わりにテンプレートを「`NewPlayerServer`」で一度解析するようにコードをリファクタリングしましょう。
ディスクからテンプレートをフェッチしたり解析したりするときに問題が発生した場合に、この関数がエラーを返すことができるようにする必要があります。

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

`NewPlayerServer`のシグネチャを変更すると、コンパイルの問題が発生します。自分で試して修正するか、苦労している場合はソースコードを参照してください。

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

テストがパスし、サーバーを実行して、`/game`で勝者を宣言します。`/league`に記録されているはずです。当選者を獲得するたびに、_接続を閉じる_ ことを思い出してください。接続を再度開くには、ページを更新する必要があります。

ユーザーがゲームの勝者を記録できる簡単なWebフォームを作成しました。これを繰り返して、ユーザーが多数のプレーヤーを提供することでゲームを開始できるようにします。サーバーはクライアントにメッセージをプッシュし、時間の経過とともにブラインドの値を通知します。

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

主な変更点は、プレーヤー数を入力するセクションとブラインド値を表示するセクションを取り込むことです。ゲームのステージに応じて、ユーザーインターフェイスを表示/非表示にするための小さなロジックがあります。

`conn.onmessage`を介して受信するメッセージはすべてブラインドアラートであると想定するため、それに応じて`blindContainer.innerText`を設定します。

ブラインドアラートを送信するにはどうすればよいですか？前の章では、`Game`のアイデアを紹介しました。これにより、CLIコードが`Game`を呼び出すことができるようになり、ブラインドアラートのスケジュール設定など、他のすべての処理が行われます。これは、心配事の良い分離であることがわかりました。

```go
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}
```

ユーザーがCLIでプレーヤーの数を求められると、ゲームを「開始`Start`」してブラインドアラートを開始し、ユーザーが勝者を宣言すると「終了`Finish`」します。これは現在の要件と同じですが、入力を取得する方法が異なります。したがって、可能であれば、この概念を再利用する必要があります。

`Game`の「実際の」実装は「`TexasHoldem`」です。

```go
type TexasHoldem struct {
	alerter BlindAlerter
	store   PlayerStore
}
```

`BlindAlerter`を送信することで、`TexasHoldem`は _wherever_ に送信されるブラインドアラートをスケジュールできます

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

常にアラートを「`os.Stdout`」に送信したいので、これはCLIで機能しますが、これはWebサーバーでは機能しません。リクエストごとに新しい`http.ResponseWriter`を取得し、それを`*websocket.Conn`にアップグレードします。
したがって、アラートを送信する必要がある依存関係を構築するときに、それを知ることはできません。

そのため、 `BlindAlerter.ScheduleAlertAt`を変更してアラートの宛先を取得し、Webサーバーで再利用できるようにする必要があります。

`blind_alerter.go`を開き、パラメータ「`to io.Writer`」を追加します

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

コンパイルしてみると、宛先なしで`ScheduleAlertAt`を呼び出しているため、`TexasHoldem`で失敗します。これにより、「今すぐ」コンパイルして、`os.Stdout`にハードコードします。

テストを実行してみてください。

`SpyBlindAlerter`は`BlindAlerter`を実装していないため、テストは失敗します。
これを修正するには、`ScheduleAlertAt`のシグネチャを更新し、テストを実行します。これで緑色のままになります。

`TexasHoldem`がブラインドアラートの送信先を知っていることは意味がありません。
ここで、`Game`を更新して、ゲームを開始するときにアラートの送信先をどこに宣言するようにします。

```go
type Game interface {
	Start(numberOfPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}
```

修正する必要があることをコンパイラーに指示させます。変更はそれほど悪くありません。

* `TexasHoldem`を更新して、`Game`を適切に実装します。
* ゲームを開始するときに、`CLI`で`out`プロパティを渡します（`cli.game.Start(numberOfPlayers、cli.out)`）
* `TexasHoldem`のテストでは、`game.Start(5, io.Discard)`を使用してコンパイルの問題を修正し、アラート出力が破棄されるように設定します

すべてが正しければ、すべてがグリーンになるはずです。
これで、`Server`内で`Game`を試して使用できます。

## 最初にテストを書く

`CLI`と`Server`の要件は同じです！配信メカニズムが異なるだけです。

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

`GameSpy`を使用して、同様の結果をテストすることができるはずです

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
* 次に、ゲームのWebソケットメッセージを送信します。
* 最後に、ゲームが開始され、期待どおりに終了したと断言します。

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

最終的に！

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

次に、コンストラクタに割り当てます

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

それについて考える必要があるので、「まだ」どこにもブラインドメッセージを送信するつもりはありません。`game.Start`を呼び出すと、それに書き込まれたメッセージを破棄する`io.Discard`を送信します。

ここでは、Webサーバーを起動します。`Game`を`PlayerServer`に渡すには、`main.go`を更新する必要があります

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

ブラインドアラートをまだ取得していないという事実を無視して、アプリは機能します！ `PlayerServer`で`Game`を再利用することに成功し、すべての詳細を処理しました。ブラインドアラートを破棄するのではなく、Webソケットに送信する方法を見つけたら、すべて正常に機能するはずです。

その前に、いくつかのコードを整理しましょう。

## リファクタリング♪

WebSocketを使用する方法はかなり基本的で、エラー処理はかなり単純なので、サーバーコードからその煩雑さを取り除くためだけにタイプにカプセル化したかったのです。後でもう一度見たいと思うかもしれませんが、今のところこれは少し整理されます

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

今サーバーコードは少し簡略化されています

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

### テストを書かないようにしよう！

どうしたらよいかわからないときは、遊んで試してみるのが一番です。作業が最初にコミットされていることを確認してください。方法を見つけたら、テストを実行する必要があります。

私たちが持っている問題のあるコード行は

```go
p.game.Start(numberOfPlayers, io.Discard) //todo: Don't discard the blinds messages!
```

ゲームがブラインドアラートを書き込むには、`io.Writer`を渡す必要があります。

以前から`playerServerWS`を渡せたらいいですね。これはWebSocketのラッパーなので、 `Game`に送信してメッセージを送信できるように感じます。

試してみてね

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

当然のことのようですが、`playerServerWS`が _does_ に`io.Writer`を実装するようにすることです。そのためには、基になる`*websocket.Conn`を使用して、`WriteMessage`を使用し、メッセージをWebSocketに送信します。

```go
func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
```

これは簡単すぎるようです！アプリケーションを試して実行し、機能するかどうかを確認します。

事前に`TexasHoldem`を編集して、ブラインドインクリメント時間を短くし、実際に動作するようにしてください

```go
blindIncrement := time.Duration(5+numberOfPlayers) * time.Second // (rather than a minute)
```

あなたはそれが機能するのを見るはずです！ブラインド量は、まるで魔法のようにブラウザで増加します。

コードを元に戻して、テストする方法を考えましょう。これを _具現化_ するために、`StartGame`へのパススルーは`io.Discard`ではなく`playerServerWS`でしたので、呼び出しをスパイして動作を確認する必要があるかもしれません。

スパイは素晴らしく、実装の詳細をチェックするのに役立ちますが、可能であれば _実際_ の動作をテストすることを常に推奨する必要があります。なぜなら、リファクタリングする場合、通常、変更しようとしている実装の詳細をチェックしているため、失敗するスパイテストがよくあるためです。 。

テストでは現在、実行中のサーバーへのWebソケット接続を開き、メッセージを送信してそれを実行します。
同様に、サーバーがWebSocket接続を介して送り返すメッセージをテストできるはずです。

## 最初にテストを書く

既存のテストを編集します。

現在、`GameSpy`は、`Start`を呼び出しても、データを`out`に送信しません。これを変更して、定型メッセージを送信するように設定し、そのメッセージがWebSocketに送信されることを確認できるようにする必要があります。これにより、必要な実際の動作を実行しながら、正しく構成したことを確信できます。

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

あらかじめ用意されたメッセージを`out`に送信するように`GameSpy``Start`を更新します。

```go
func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartCalledWith = numberOfPlayers
	out.Write(g.BlindAlert)
}
```

つまり、ゲームを`Start`しようとしたときに`PlayerServer`を実行すると、正常に機能していれば、WebSocketを介してメッセージを送信することになります。

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
* メッセージがwebsocket接続で送信されることを期待して、`ws.ReadMessage()`への呼び出しを追加して、メッセージが送信されるのを待ってから、メッセージが期待どおりであることを確認します。

## テストを実行してみます

テストが永久にハングすることを確認する必要があります。これは、`ws.ReadMessage()`がメッセージを取得するまでブロックするためです。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

テストがハングすることは絶対にあってはならないことなので、タイムアウトさせたいコードを処理する方法を紹介します。

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

`within`が行うことは、関数`assert`を引数として取り、それをgoルーチンで実行することです。関数が完了すると、`完了（done）`チャネルを介して行われたことを通知します。

その間、チャネルがメッセージを送信するのを待機できるようにする`select`ステートメントを使用します。ここからは、`assert`関数と`time.After`の間の競争であり、期間が発生したときに信号を送信します。

最後に、アサーションのヘルパー関数を作成しました。

```go
func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}
```

これがテストの今の読み方です

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

最後に、サーバーコードを変更して、開始時にWebSocket接続をゲームに送信できるようにします。

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

ヘルパーの`assertGameStartedWith`と`assertFinishCalledWith`をリファクタリングして、失敗する前に短い間アサーションを再試行できます。

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

これでアプリケーションが完成しました。ポーカーゲームはWebブラウザーを介して開始でき、ユーザーはWebSocketを介して時間が経過するにつれてブラインドベットの価値を知らされます。ゲームが終了すると、数章前に書いたコードを使用して永続化された勝者を記録できます。
プレーヤーは、ウェブサイトの `/league`エンドポイントを使用して、誰が（または最も幸運な）ポーカープレーヤーであるかを見つけることができます。

旅の途中で間違いを犯しましたが、TDDフローを使用して、作業用ソフトウェアから遠く離れたことはありませんでした。私たちは自由に反復して実験を続けることができました。

最後の章では、アプローチ、私たちが到達した設計について振り返り、いくつかの緩い目的を結びつけます。

この章ではいくつかのことを説明しました

### WebSockets

* クライアントがサーバーをポーリングし続ける必要がないクライアントとサーバー間でメッセージを送信する便利な方法。私たちが持っているクライアントとサーバーのコードはどちらも非常に単純です。
* テストは簡単ですが、テストの非同期の性質に注意する必要があります

### 遅延する可能性がある、または終了しないテストでのコードの処理

* アサーションを再試行してタイムアウトを追加するヘル​​パー関数を作成します。
* ゴルーチンを使用して、アサーションが何もブロックしないようにし、チャネルを使用して、終了したかどうかを通知できます。
* `time`パッケージには、タイムアウトを設定できるように、時間内のイベントに関する信号をチャンネル経由で送信するいくつかの便利な関数があります
