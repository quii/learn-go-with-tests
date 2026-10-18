---
description: Time
---

# 時間

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/time)

製品の所有者は、[テキサス・ホールデム・ポーカー](https://ja.wikipedia.org/wiki/%E3%83%86%E3%82%AD%E3%82%B5%E3%82%B9%E3%83%BB%E3%83%9B%E3%83%BC%E3%83%AB%E3%83%87%E3%83%A0)をプレイする人々のグループを支援することで、コマンドラインアプリケーションの機能を拡張することを望んでいます。

## ポーカーに関する十分な情報

ポーカーについて多くのことを知る必要はありません。一定の時間間隔で、すべてのプレーヤーに着実に増加する「ブラインド（強制bet）」値を通知する必要があるということだけです。

私たちのアプリケーションは、ブラインドがいつ上がるべきか、そしてどれくらいあるべきかを追跡するのに役立ちます。

* 開始すると、何人のプレイヤーがプレイしているかを尋ねます。これは、「ブラインド」ベットが上がるまでの時間を決定します。
  * 基本時間は5分です。
  * すべてのプレーヤーについて、1分が追加されます。
  * 例：6人のプレーヤーは、ブラインドでは11分に相当します。
* ブラインドタイムが終了した後、ゲームはプレイヤーにブラインドベットの新しい金額を警告する必要があります。
* ブラインドは100チップから始まり、その後200、400、600、1000、2000となり、ゲームが終了するまで2倍になります（「ルース勝利」の以前の機能は、ゲームを終了するはずです）。

## コードの注意

前の章では、`{name} wins`のコマンドを既に受け入れているコマンドラインアプリケーションを開始しました。これが現在の`CLI`コードの外観ですが、開始する前に他のコードについてもよく理解してください。

```go
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}

func (cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
```

### `time.AfterFunc`

プレーヤーの数に応じて特定の期間にブラインドベット値を印刷するようにプログラムをスケジュールできるようにしたいと考えています。

必要なことの範囲を制限するために、ここではプレーヤーの数を忘れて5人のプレーヤーがいると仮定し、_ブラインドベットの新しい値が10分ごとに出力されるかどうかをテストします_。

いつものように、標準ライブラリは[`func AfterFunc(d Duration, f func()) *Timer`](https://golang.org/pkg/time/#AfterFunc)でカバーされています。

> `AfterFunc`は継続時間が経過するのを待ってから、独自のgoroutineでfを呼び出します。これは、Stopメソッドを使用して呼び出しをキャンセルするために使用できる`Timer`を返します。

### [`time.Duration`](https://golang.org/pkg/time/#Duration)

> 期間は、2つの瞬間間の経過時間をint64ナノ秒カウントとして表します。

タイムライブラリには、これらのナノ秒を乗算できるようにする定数がいくつかあります。これにより、これらの定数は、私たちが行うシナリオの種類に対して、より読みやすくなります。

```
5 * time.Second
```

「`PlayPoker`」を呼び出すと、すべてのブラインドアラートをスケジュールします。

これをテストするのは少し難しいかもしれません。各期間が正しいブラインド量でスケジュールされていることを確認する必要がありますが、`time.AfterFunc`のシグネチャを見ると、2番目の引数は実行される関数です。
Goでは関数を比較できないため、送信された関数をテストすることはできません。そのため、実行に時間と印刷する量を必要とする、`time.AfterFunc`の周りにある種のラッパーを書く必要があります。私たちはそれをスパイすることができます。

## 最初にテストを書く

スイートに新しいテストを追加する

```go
t.Run("it schedules printing of blind values", func(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &poker.StubPlayerStore{}
	blindAlerter := &SpyBlindAlerter{}

	cli := poker.NewCLI(playerStore, in, blindAlerter)
	cli.PlayPoker()

	if len(blindAlerter.alerts) != 1 {
		t.Fatal("expected a blind alert to be scheduled")
	}
})
```

`SpyBlindAlerter`を作成し、それを`CLI`に挿入しようとしています。次に、 `PlayPoker`を呼び出した後、アラートがスケジュールされていることを確認します。

（最初に最も単純なシナリオを実行することを思い出して、次に反復します。）

ここに`SpyBlindAlerter`の定義があります

```go
type SpyBlindAlerter struct {
	alerts []struct {
		scheduledAt time.Duration
		amount      int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct {
		scheduledAt time.Duration
		amount      int
	}{duration, amount})
}
```

## テストを実行してみます

```
./CLI_test.go:32:27: too many arguments in call to poker.NewCLI
	have (*poker.StubPlayerStore, *strings.Reader, *SpyBlindAlerter)
	want (poker.PlayerStore, io.Reader)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

新しい引数を追加しましたが、コンパイラは不満を言っています。
厳密に言えば、最小限のコードは`NewCLI`が`*SpyBlindAlerter`を受け入れるようにすることですが、少し騙して、依存関係をインターフェイスとして定義しましょう。

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}
```

そして、それをコンストラクタに追加します。

```go
func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI
```

他のテストは`NewCLI`に`BlindAlerter`が渡されていないので失敗します。

`BlindAlerter`をスパイすることは他のテストには関係ないので、テストファイルに追加します。

```go
var dummySpyAlerter = &SpyBlindAlerter{}
```

そして、コンパイルの問題を修正するために他のテストでそれを使ってください。
これを「ダミー`"dummy"`」と表記することで、テストの読者には重要ではないことが明らかになります。

[ダミーオブジェクトは渡されますが、実際に使われることはありません。通常はパラメータリストを埋めるために使われるだけです](https://martinfowler.com/articles/mocksArentStubs.html)

これでテストはコンパイルされ、新しいテストは失敗します。

```
=== RUN   TestCLI
=== RUN   TestCLI/it_schedules_printing_of_blind_values
--- FAIL: TestCLI (0.00s)
    --- FAIL: TestCLI/it_schedules_printing_of_blind_values (0.00s)
    	CLI_test.go:38: expected a blind alert to be scheduled
```

## 通過するのに十分なコードを書く

これを`PlayPoker`メソッドで参照できるようにするために`BlindAlerter`を`CLI`のフィールドとして追加する必要があります。

```go
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	alerter     BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		alerter:     alerter,
	}
}
```

テストを通過させるために、`BlindAlerter`を好きなもので呼び出すことができます。

```go
func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5*time.Second, 100)
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

次は、5人のプレイヤーのために、私たちが望むすべてのアラートのスケジュールを確認したいと思います。

## 最初にテストを書く

```go
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		cases := []struct {
			expectedScheduleTime time.Duration
			expectedAmount       int
		}{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduleTime), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				amountGot := alert.amount
				if amountGot != c.expectedAmount {
					t.Errorf("got amount %d, want %d", amountGot, c.expectedAmount)
				}

				gotScheduledTime := alert.scheduledAt
				if gotScheduledTime != c.expectedScheduleTime {
					t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, c.expectedScheduleTime)
				}
			})
		}
	})
```

テーブルベースのテストはここではうまく機能し、要件が何であるかを明確に示しています。テーブルを実行し、`SpyBlindAlerter`をチェックして、アラートが正しい値でスケジュールされているかどうかを確認します。

## テストを実行してみる

こんな感じで失敗を重ねるといいですよ。

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/it_schedules_printing_of_blind_values
    --- FAIL: TestCLI/it_schedules_printing_of_blind_values (0.00s)
=== RUN   TestCLI/it_schedules_printing_of_blind_values/100_scheduled_for_0s
        --- FAIL: TestCLI/it_schedules_printing_of_blind_values/100_scheduled_for_0s (0.00s)
        	CLI_test.go:71: got scheduled time of 5s, want 0s
=== RUN   TestCLI/it_schedules_printing_of_blind_values/200_scheduled_for_10m0s
        --- FAIL: TestCLI/it_schedules_printing_of_blind_values/200_scheduled_for_10m0s (0.00s)
        	CLI_test.go:59: alert 1 was not scheduled [{5000000000 100}]
```

## 通過するのに十分なコードを書く

```go
func (cli *CLI) PlayPoker() {

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}

	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

それは、私たちがすでに持っていたものよりもそれほど複雑ではありません。
私たちは今、`blinds`の配列を反復処理し、増加する`blindTime`でスケジューラを呼び出しています

## リファクタリング♪

「`PlayPoker`」を少し明確にするために、スケジュールされたアラートをメソッドにカプセル化できます。

```go
func (cli *CLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}
```

最後に、テストは少し不格好に見えます。同じものを表す2つの匿名構造体、`ScheduledAlert`があります。それを新しい型にリファクタリングして、いくつかのヘルパーを比較してみましょう。

```go
type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}
```

タイプに `String()`メソッドを追加したので、テストが失敗した場合にうまく表示されます。

新しいタイプを使用するようにテストを更新します。

```go
t.Run("it schedules printing of blind values", func(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &poker.StubPlayerStore{}
	blindAlerter := &SpyBlindAlerter{}

	cli := poker.NewCLI(playerStore, in, blindAlerter)
	cli.PlayPoker()

	cases := []scheduledAlert{
		{0 * time.Second, 100},
		{10 * time.Minute, 200},
		{20 * time.Minute, 300},
		{30 * time.Minute, 400},
		{40 * time.Minute, 500},
		{50 * time.Minute, 600},
		{60 * time.Minute, 800},
		{70 * time.Minute, 1000},
		{80 * time.Minute, 2000},
		{90 * time.Minute, 4000},
		{100 * time.Minute, 8000},
	}

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
})
```

自分で「`assertScheduledAlert`」を実装します。

ここではかなりの時間をかけてテストを作成しており、アプリケーションと統合しないことはややエッチです。これ以上の要件に取り掛かる前に、その問題に取り組みましょう。

アプリを実行してみてください。アプリがコンパイルされず、「`NewCLI`」への引数が足りないという不満があります。

アプリケーションで使用できる`BlindAlerter`の実装を作成してみましょう。

`blind_alerter.go`を作成し、`BlindAlerter`インターフェースを移動して、以下の新しいものを追加します

```go
package poker

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
```

_type_ は、`struct`だけでなく、インターフェイスを実装できることに注意してください。 1つの関数が定義されたインターフェースを公開するライブラリを作成する場合、`MyInterfaceFunc`タイプも公開することは一般的な慣用法です。

このタイプはあなたのインターフェースも実装する`func`になります。こうすることで、インターフェイスのユーザーは、関数だけでインターフェイスを実装することができます。
空の`struct`タイプを作成する必要はありません。

次に、関数と同じシグネチャを持つ関数`StdOutAlerter`を作成し、`time.AfterFunc`を使用して、それを`os.Stdout`に出力するようにスケジュールします。

この動作を確認するには、`NewCLI`を作成する`main`を更新します。

```go
poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter)).PlayPoker()
```

実行する前に、「`CLI`」の「`blindTime`」の増分を10分ではなく10秒に変更すると、実際の動作を確認できます。

10秒ごとに予想されるように、ブラインド値が出力されるはずです。 CLIに「`Shaun wins`」と入力すると、プログラムが予期したとおりに停止することに注意してください。

ゲームは常に5人でプレイされるとは限らないため、ゲームを開始する前に、ユーザー数を入力するようユーザーに促す必要があります。

## 最初にテストを書く

確認するために、StdOutに書き込まれた内容を記録するプレーヤーの数を入力するように求めています。
これを数回実行しました。`os.Stdout`は`io.Writer`であることを知っているので、テストで依存性注入を使用して`bytes.Buffer`を渡し、私たちのコードが何を書くか見てください。

このテストでは、他の共同編集者についてはまだ気にしていないため、テストファイルでダミーを作成しました。

`CLI`に4つの依存関係があることに少し注意する必要があります。
これは、多すぎる責任を持ち始めているように思われます。とりあえずそれと共存して、この新しい機能を追加するときにリファクタリングが出現するかどうか見てみましょう。

```go
var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}
```

これが新しいテストです。

```go
t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
	stdout := &bytes.Buffer{}
	cli := poker.NewCLI(dummyPlayerStore, dummyStdIn, stdout, dummyBlindAlerter)
	cli.PlayPoker()

	got := stdout.String()
	want := "Please enter the number of players: "

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
})
```

`main`で`os.Stdout`になるものを渡し、何が書き込まれるかを確認します。

## テストを実行してみます

```
./CLI_test.go:38:27: too many arguments in call to poker.NewCLI
	have (*poker.StubPlayerStore, *bytes.Buffer, *bytes.Buffer, *SpyBlindAlerter)
	want (poker.PlayerStore, io.Reader, poker.BlindAlerter)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

新しい依存関係があるため、`NewCLI`を更新する必要があります

```go
func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI
```

その他のテストは、`NewCLI`に渡される`io.Writer`がないため、コンパイルに失敗します。

他のテスト用に「`dummyStdOut`」を追加します。

新しいテストはそのように失敗するはずです。

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players
    --- FAIL: TestCLI/it_prompts_the_user_to_enter_the_number_of_players (0.00s)
    	CLI_test.go:46: got '', want 'Please enter the number of players: '
FAIL
```

## 成功させるのに十分なコードを書く

`CLI`に新しい依存関係を追加して、`PlayPoker`で参照できるようにする必要があります。

```go
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
	out         io.Writer
	alerter     BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		out:         out,
		alerter:     alerter,
	}
}
```

最後に、ゲームの開始時にプロンプ​​トを書き込むことができます。

```go
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, "Please enter the number of players: ")
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

## リファクタリング♪

定数に抽出する必要があるプロンプトの重複文字列があります。

```go
const PlayerPrompt = "Please enter the number of players: "
```

テストコードと`CLI`の両方でこれを使用してください。

次に、番号を送信して抽出する必要があります。
目的の効果があったかどうかを確認する唯一の方法は、スケジュールされたブラインドアラートを確認することです。

## 最初にテストを書く

```go
t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
	stdout := &bytes.Buffer{}
	in := strings.NewReader("7\n")
	blindAlerter := &SpyBlindAlerter{}

	cli := poker.NewCLI(dummyPlayerStore, in, stdout, blindAlerter)
	cli.PlayPoker()

	got := stdout.String()
	want := poker.PlayerPrompt

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	cases := []scheduledAlert{
		{0 * time.Second, 100},
		{12 * time.Minute, 200},
		{24 * time.Minute, 300},
		{36 * time.Minute, 400},
	}

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
})
```

いたたたた！たくさんの変更。。

* StdInのダミーを削除し、代わりに、7を入力するユーザーを表すモックバージョンを送信します。
* また、ブラインドアラートでダミーを削除して、プレーヤー数がスケジュールに影響を与えていることを確認できるようにしました
* スケジュールされているアラートをテストします

## テストを実行してみます

テストはコンパイルされ、失敗するはずですが、5人のプレーヤーに基づいてゲームをハードコーディングしているため、スケジュールされた時間が間違っていると報告されます。

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players
    --- FAIL: TestCLI/it_prompts_the_user_to_enter_the_number_of_players (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players/100_chips_at_0s
        --- PASS: TestCLI/it_prompts_the_user_to_enter_the_number_of_players/100_chips_at_0s (0.00s)
=== RUN   TestCLI/it_prompts_the_user_to_enter_the_number_of_players/200_chips_at_12m0s
```

## 成功させるのに十分なコードを書く

この作業を行うために必要なすべての罪を自由に犯すことを忘れないでください。
作業用ソフトウェアができたら、作成しようとしている混乱のリファクタリングに取り掛かることができます。

```go
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, _ := strconv.Atoi(cli.readLine())

	cli.scheduleBlindAlerts(numberOfPlayers)

	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}
```

* `numberOfPlayersInput`を文字列に読み込みます
* ユーザーから入力を取得するために `cli.readLine()`を使用し、次にエラーシナリオを無視して、`Atoi`を呼び出して整数に変換します。後でそのシナリオのテストを作成する必要があります。
* ここから、`scheduleBlindAlerts`を変更して、多数のプレーヤーを受け入れます。次に、ブラインド量を反復するときに「`blindTime`」に追加するために使用する「`blindIncrement`」時間を計算します

新しいテストは修正されましたが、システムが機能するのは、ユーザーが数字を入力してゲームを開始した場合に限られるためです。
ユーザー入力を変更してテストを修正する必要があります。
これにより、番号とそれに続く改行が追加されます（これは、現在のアプローチのさらに多くの欠陥を強調しています）。

## リファクタリング♪

これは少し恐ろしいことだと思いますか？ **テストを聞いてみましょう**。

* いくつかのアラートをスケジュールしていることをテストするために、4つの異なる依存関係を設定しました。システムに _thing_ の依存関係がたくさんある場合、それは多すぎることを意味します。視覚的には、テストが雑然としていることがわかります。
* 私には、**ユーザー入力の読み取りと、実行したいビジネスロジックとの間をより明確に抽象化する必要があるように感じます**
* より適切なテストは、_このユーザー入力が与えられた場合、正しいタイプのプレーヤーで新しいタイプ`Game`を呼び出すかどうかです。
* 次に、スケジューリングのテストを新しい`Game`のテストに抽出します。

最初に「`Game`」に向けてリファクタリングでき、テストは引き続き成功します。必要な構造変更を行ったら、懸念の分離を反映するためにテストをリファクタリングする方法について考えることができます

リファクタリングに変更を加えるときは、それらをできるだけ小さくして、テストを再実行し続けるようにしてください。

まず自分で試してください。 「`Game`」が提供するものと「`CLI`」が行うべきことの境界について考えてください。

今のところは、`NewCLI`の外部インターフェースを**変更しないでください**。
テストコードとクライアントコードを同時に変更したくないからです。

これは私が思いついたものです。

```go
// game.go
type Game struct {
	alerter BlindAlerter
	store   PlayerStore
}

func (p *Game) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (p *Game) Finish(winner string) {
	p.store.RecordWin(winner)
}

// cli.go
type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game *Game
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		in:  bufio.NewScanner(in),
		out: out,
		game: &Game{
			alerter: alerter,
			store:   store,
		},
	}
}

const PlayerPrompt = "Please enter the number of players: "

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayersInput := cli.readLine()
	numberOfPlayers, _ := strconv.Atoi(strings.Trim(numberOfPlayersInput, "\n"))

	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)

	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins\n", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
```

「ドメイン」の観点から

* 何人がプレイしているかを示す「`Game`」を「開始`Start`」したい
* 勝者を宣言して「`Game`」を「終了`Finish`」したい

新しい`Game`タイプはこれをカプセル化します。

この変更により、`BlindAlerter`と`PlayerStore`が`Game`に渡されました。
これは、結果のアラートと保存を担当するためです。

私たちの`CLI`は今ちょうど関係しています

* 既存の依存関係を使用して`Game`を構築します（これは次にリファクタリングします）
* ユーザー入力を`Game`のメソッド呼び出しとして解釈する

「大きな」リファクタリングを行わないようにしたいと考えています。
これにより、長期間にわたってテストが失敗し、ミスの可能性が高まるためです。
（大規模な分散型チームで作業している場合、これは非常に重要です）

まず最初に、`Game`をリファクタリングして、`CLI`に注入します。
テストを最小限に変更してそれを容易にし、次に、テストをユーザー入力の解析とゲーム管理のテーマに分割する方法を確認します。

今必要なのは、`NewCLI`を変更することだけです。

```go
func NewCLI(in io.Reader, out io.Writer, game *Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}
```

これはすでに改善のように感じられます。
依存関係は少なく、私たちの依存関係リストは、CLIが入出力に関係し、ゲーム固有のアクションを`Game`に委任するという、全体的な設計目標を反映しています。

コンパイルしてみると問題があります。これらの問題は自分で修正できるはずです。
今すぐ`Game`のモックを作成する必要はありません。すべてをコンパイルしてテストするために「実際の」`Game`を初期化するだけです。

これを行うには、コンストラクタを作成する必要があります。

```go
func NewGame(alerter BlindAlerter, store PlayerStore) *Game {
	return &Game{
		alerter: alerter,
		store:   store,
	}
}
```

これは、修正されたテストの設定の1つの例です。

```go
stdout := &bytes.Buffer{}
in := strings.NewReader("7\n")
blindAlerter := &SpyBlindAlerter{}
game := poker.NewGame(blindAlerter, dummyPlayerStore)

cli := poker.NewCLI(in, stdout, game)
cli.PlayPoker()
```

テストを修正して再び緑に戻るのにそれほどの労力は必要ありませんが、それがポイントです！
次の段階の前に`main.go`も修正するようにしてください。

```go
// main.go
game := poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
cli := poker.NewCLI(os.Stdin, os.Stdout, game)
cli.PlayPoker()
```

「ゲーム`Game`」を抽出したので、ゲーム固有のアサーションをCLIとは別のテストに移動する必要があります。

これは、`CLI`テストをコピーするための単なる練習ですが、依存関係は少なくなっています。

```go
func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Minute, Amount: 200},
			{At: 20 * time.Minute, Amount: 300},
			{At: 30 * time.Minute, Amount: 400},
			{At: 40 * time.Minute, Amount: 500},
			{At: 50 * time.Minute, Amount: 600},
			{At: 60 * time.Minute, Amount: 800},
			{At: 70 * time.Minute, Amount: 1000},
			{At: 80 * time.Minute, Amount: 2000},
			{At: 90 * time.Minute, Amount: 4000},
			{At: 100 * time.Minute, Amount: 8000},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

}

func TestGame_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}
```

ポーカーのゲームが始まったときに何が起きるかの背後にある意図が今やはるかに明確になりました。

ゲームがいつ終了するかについても、テストの上を移動してください。

満足したら、ゲームロジックのテストを移行しました。
CLIテストを簡略化して、意図した責任をより明確に反映できます。

* ユーザー入力を処理し、必要に応じて`Game`のメソッドを呼び出します
* 出力を送信
* 重要なのは、ゲームがどのように機能するかについての実際の仕組みについて知らない

これを行うには、`CLI`が具体的な`Game`タイプに依存しないようにする必要がありますが、代わりに`Start(numberOfPlayers)`および`Finish(winner)`のインターフェースを受け入れます。次に、そのタイプのスパイを作成し、正しい呼び出しが行われたことを確認できます。

ここに、ネーミングが時々厄介であることがわかります。`Game`の名前を`TexasHoldem`に変更します（これが現在プレイしているゲームの _種類_ であるため）。新しいインターフェイスは`Game`と呼ばれます。これは、CLIが私たちがプレイしている実際のゲームに気づいていないという概念と、「`Start`」および「`Finish`」したときに何が起こるかを忠実に保ちます。

```go
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}
```


`CLI`内の`*Game`へのすべての参照を置き換え、それらを`Game`（our new interface）に置き換えます。
いつものように、リファクタリングしている間は、テストを再実行してすべてが緑色であることを確認します。

`TexasHoldem`から`CLI`を分離したので、スパイを使用して、正しい引数で、期待どおりに`Start`と`Finish`が呼び出されることを確認できます。

`Game`を実装するスパイを作成する

```go
type GameSpy struct {
	StartedWith  int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}
```

ゲーム固有のロジックをテストするすべての`CLI`テストを、`GameSpy`の呼び出し方法のチェックに置き換えます。これは、テストにおけるCLIの責任を明確に反映します。

これは、修正されたテストの1つの例です。
残りを自分で試して、行き詰まった場合はソースコードを確認してください。

```go
	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}
	})
```

懸念事項を明確に分離したので、`CLI`でIOの周辺のケースをチェックするのが簡単になります。

プレイヤー数の入力を求められたときにユーザーが非数値を入力するシナリオに対処する必要があります。

私たちのコードはゲームを開始してはならず、ユーザーに便利なエラーを出力して終了します。


## 最初にテストを書く

ゲームが開始しないことを確認することから始めます

```go
t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
	stdout := &bytes.Buffer{}
	in := strings.NewReader("Pies\n")
	game := &GameSpy{}

	cli := poker.NewCLI(in, stdout, game)
	cli.PlayPoker()

	if game.StartCalled {
		t.Errorf("game should not have started")
	}
})
```

「`GameSpy`」にフィールド「`StartCalled`」を追加する必要があります。これは「`Start`」が呼び出された場合にのみ設定されます

## テストを実行してみます

```
=== RUN   TestCLI/it_prints_an_error_when_a_non_numeric_value_is_entered_and_does_not_start_the_game
    --- FAIL: TestCLI/it_prints_an_error_when_a_non_numeric_value_is_entered_and_does_not_start_the_game (0.00s)
        CLI_test.go:62: game should not have started
```

## 成功させるのに十分なコードを書く

`Atoi`と呼ぶところは、エラーをチェックするだけです

```go
numberOfPlayers, err := strconv.Atoi(cli.readLine())

if err != nil {
	return
}
```

次に、ユーザーが間違ったことをユーザーに通知する必要があるため、`stdout`に出力される内容を評価します。

## 最初にテストを書く

前に「`stdout`」に出力されたものをアサートしたので、今のところそのコードをコピーできます

```go
gotPrompt := stdout.String()

wantPrompt := poker.PlayerPrompt + "you're so silly"

if gotPrompt != wantPrompt {
	t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
}
```

`stdout`に書き込まれる _常に_ を保存しているので、`poker.PlayerPrompt`がまだ必要です。
次に、追加のものが表示されることを確認します。
今のところ、正確な表現についてあまり気にせず、リファクタリングするときに対処します。

## テストを実行してみます

```
=== RUN   TestCLI/it_prints_an_error_when_a_non_numeric_value_is_entered_and_does_not_start_the_game
    --- FAIL: TestCLI/it_prints_an_error_when_a_non_numeric_value_is_entered_and_does_not_start_the_game (0.00s)
        CLI_test.go:70: got 'Please enter the number of players: ', want 'Please enter the number of players: you're so silly'
```

## 成功させるのに十分なコードを書く

エラー処理コードを変更する

```go
if err != nil {
	fmt.Fprint(cli.out, "you're so silly")
	return
}
```

## リファクタリング♪

次に、メッセージを「`PlayerPrompt`」のような定数にリファクタリングします

```go
wantPrompt := poker.PlayerPrompt + poker.BadPlayerInputErrMsg
```

より適切なメッセージを入れます

```go
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
```

最後に、`stdout`に送信されたものに関するテストは非常に詳細です。クリーンアップするアサート関数を作成しましょう。

```go
func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}
```

可変長の構文（`...string`）を使用すると、さまざまな量のメッセージに対してアサートする必要があるため、ここでは便利です。

このヘルパーは、ユーザーに送信されるメッセージに対してアサートする両方のテストで使用します。

一部の「`assertX`」関数で役立つ可能性のあるテストがいくつかあるので、テストをクリーンアップして読みやすくすることでリファクタリングを練習してください。

時間をかけて、私たちが追い出したいくつかのテストの価値について考えてください。
必要以上のテストは必要ありません。
それらの一部をリファクタリング/削除しても、すべてが機能することを確信できますか？

これが私が思いついたものです。

```go
func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Chris wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		game := &GameSpy{}

		stdout := &bytes.Buffer{}
		in := userSends("pies")

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})
}
```

テストはCLIの主な機能を反映するようになり、何人のプレイヤーに不正な値が入力されたときに何人がプレイし、誰が勝って処理するかという観点からユーザー入力を読み取ることができます。
これを行うことにより、`CLI`が何をするかを読者に明らかにするだけでなく、何をしないかもわかります。

「ルースが勝つ`Ruth wins`」の代わりに「ロイドはキラー`Lloyd is a killer`」にユーザーが入れるとどうなりますか？

このシナリオのテストを記述して成功させることにより、この章を終了します。

## まとめ

### プロジェクトの簡単な要約

過去5つの章については、かなりの量のコードをゆっくりTDDしました

* コマンドラインアプリケーションとWebサーバーの2つのアプリケーションがあります。
* これらのアプリケーションは両方とも、勝者を記録するために「`PlayerStore`」に依存しています
* Webサーバーは、最も多くのゲームに勝っている人のリーグテーブルを表示することもできます
* コマンドラインアプリは、プレーヤーが現在のブラインド値を追跡することでポーカーゲームをプレイするのに役立ちます。

### time.Afterfunc

特定の期間後に関数呼び出しをスケジュールする非常に便利な方法。時間を費やす価値があります[`time`のドキュメントを見る](https://golang.org/pkg/time/)。時間を節約するための関数やメソッドがたくさんあるので、作業するのに役立ちます。

私のお気に入りのいくつかは

* 期間が終了すると、`time.After(duration)`は`chan Time`を返します。したがって、特定の時間の後に何かをしたい場合は、これが役立ちます。
* `time.NewTicker(duration)`は、チャネルを返すという点で上記と同様の `Ticker`を返しますが、これは1回だけではなく、すべての期間を「ティック`("ticks")`」します。これは、`N duration`ごとにコードを実行する場合に非常に便利です。

### 懸念事項の適切な分離のその他の例

一般的には、ユーザー入力と応答の処理の責任をドメインコードから分離することをお勧めします。
コマンドラインアプリケーションとWebサーバーで確認できます。

テストが乱雑になりました。アサーションが多すぎ、この入力を確認し、これらのアラートをスケジュールなどをして、依存関係が多すぎた。
散らかっていることを視覚的に確認できました。 **テストに耳を傾けることはとても重要です**。

* テストが乱雑に見える場合は、リファクタリングしてみてください。
* これを行ってもまだ混乱している場合は、デザインの欠陥を指摘している可能性が高いです。
* これはテストの真の強みの1つです。

テストと製品コードは少し雑然としていましたが、テストに基づいて自由にリファクタリングできました。

これらの状況に陥ったときは、常に小さなステップを踏んで、変更のたびにテストを再実行することを忘れないでください。

テストコードと本番コードの両方を同時にリファクタリングするのは危険でした。
そのため、インターフェースを変更せずに、最初に本番コードをリファクタリングしました（現在の状態では、テストを大幅に改善することはできませんでした）。
物事を変えている間、私たちができる限りテストに依存することができました。 _そして_ デザインが改善された後、テストをリファクタリングしました。

依存関係リストをリファクタリングした後、設計目標を反映しました。
これは、意図を文書化することが多いという点で、DIのもう1つの利点です。グローバル変数に依存すると、責任が非常に不明確になります。

## インターフェースを実装する関数の例

1つのメソッドでインターフェースを定義するとき、ユーザーが関数だけでインターフェースを実装できるように、それを補完するために「`MyInterfaceFunc`」型を定義することを検討することができます。

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

// BlindAlerterFunc allows you to implement BlindAlerter with a function
type BlindAlerterFunc func(duration time.Duration, amount int)

// ScheduleAlertAt is BlindAlerterFunc implementation of BlindAlerter
func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}
```

こうすることで、あなたのライブラリを利用する人は、関数だけであなたのインタフェースを実装することができます。彼らは [Type Conversion](https://go.dev/tour/basics/13) を使って関数を `BlindAlerterFunc` に変換し、それを BlindAlerter として使うことができます (`BlindAlerterFunc` は `BlindAlerter` を実装しているので)。

```go
game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
```

ここで重要なのは、Goでは構造体だけでなく、_型_ にもメソッドを追加できるということです。これは非常に強力な機能で、これを利用すればより便利な方法でインターフェースを実装することができます。

関数の型を定義するだけでなく、他の型の周りに型を定義して、その型にメソッドを追加することができることを考慮してください。

```go
type Blog map[string]string

func (b Blog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, b[r.URL.Path])
}
```

ここでは、非常にシンプルな「ブログ」を実装するHTTPハンドラを作成しました。このハンドラでは、URLパスをキーとして、マップに格納された投稿を表示します。
