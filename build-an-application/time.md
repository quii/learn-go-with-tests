---
description: 時間
---

# 時間

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/time)

プロダクトオーナーは、[テキサス・ホールデム・ポーカー](https://ja.wikipedia.org/wiki/%E3%83%86%E3%82%AD%E3%82%B5%E3%82%B9%E3%83%BB%E3%83%9B%E3%83%BC%E3%83%AB%E3%83%87%E3%83%A0)をプレイする人たちのグループを手助けできるよう、コマンドラインアプリケーションの機能を拡張してほしいと考えています。

## ポーカーについて必要最低限の知識

ポーカーについて多くのことを知る必要はありません。一定の時間間隔で、すべてのプレーヤーに着実に増加する「ブラインド（強制bet）」値を通知する必要があるということだけです。

私たちのアプリケーションは、ブラインドをいつ上げるべきか、そしていくらにすべきかを把握するのに役立ちます。

* 開始すると、何人のプレーヤーがプレイしているかを尋ねます。これは、「ブラインド」ベットが上がるまでの時間を決定します。
  * 基本時間は5分です。
  * プレーヤー1人につき、1分が追加されます。
  * 例：プレーヤーが6人の場合、ブラインドの時間は11分になります。
* ブラインドの時間が過ぎたら、ゲームはプレーヤーに新しいブラインドベットの額を知らせる必要があります。
* ブラインドは100チップから始まり、その後200、400、600、1000、2000となり、その後はゲームが終了するまで2倍ずつ増えていきます（以前実装した「Ruth wins」でゲームを終了する機能は引き続き動作する必要があります）。

## これまでのコードのおさらい

前の章では、`{name} wins`というコマンドを受け付けるコマンドラインアプリケーションを作り始めました。現在の`CLI`のコードは次のとおりですが、始める前に他のコードにも目を通しておいてください。

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

プレーヤーの数に応じた一定の時間間隔で、ブラインドベットの値を出力するようプログラムをスケジュールできるようにしたいと考えています。

作業の範囲を絞るため、今のところプレーヤーの数のことは忘れて5人のプレーヤーがいると仮定し、_10分ごとにブラインドベットの新しい値が出力されること_ をテストします。

いつものように、標準ライブラリには[`func AfterFunc(d Duration, f func()) *Timer`](https://golang.org/pkg/time/#AfterFunc)が用意されています。

> `AfterFunc`は継続時間が経過するのを待ってから、自身のゴルーチンでfを呼び出します。これは、Stopメソッドを使用して呼び出しをキャンセルするために使用できる`Timer`を返します。

### [`time.Duration`](https://golang.org/pkg/time/#Duration)

> Durationは、2つの時点間の経過時間をint64のナノ秒数として表します。

timeライブラリには、ナノ秒に掛け合わせるための定数がいくつか用意されているので、これから扱うようなシナリオでも読みやすく書けます。

```
5 * time.Second
```

`PlayPoker`を呼び出したときに、すべてのブラインドアラートをスケジュールします。

これをテストするのは少し難しいかもしれません。各時刻が正しいブラインド額でスケジュールされていることを確認する必要がありますが、`time.AfterFunc`のシグネチャを見ると、2番目の引数は実行される関数です。
Goでは関数を比較できないため、どの関数が渡されたかをテストすることはできません。そこで、実行する時刻と出力する額を受け取る、`time.AfterFunc`のラッパーのようなものを書いて、それをスパイできるようにする必要があります。

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

`SpyBlindAlerter`を作成して`CLI`に注入し、`PlayPoker`を呼び出した後にアラートがスケジュールされていることを確認している点に気づくでしょう。

（まずは最も単純なシナリオから始めて、その後で繰り返し改良していくことを思い出してください。）

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
_厳密に言えば_、最小限のコードは`NewCLI`が`*SpyBlindAlerter`を受け入れるようにすることですが、少しズルをして、依存関係をインターフェースとして定義してしまいましょう。

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

`BlindAlerter`をスパイすることは他のテストには関係ないので、テストファイルに次を追加します。

```go
var dummySpyAlerter = &SpyBlindAlerter{}
```

そして、コンパイルの問題を修正するために他のテストでそれを使ってください。
これを「ダミー」と名付けることで、テストの読者にはそれが重要ではないことが明らかになります。

[ダミーオブジェクトは渡されますが、実際に使われることはありません。通常はパラメータリストを埋めるために使われるだけです](https://martinfowler.com/articles/mocksArentStubs.html)

これでテストはコンパイルされ、新しいテストは失敗します。

```
=== RUN   TestCLI
=== RUN   TestCLI/it_schedules_printing_of_blind_values
--- FAIL: TestCLI (0.00s)
    --- FAIL: TestCLI/it_schedules_printing_of_blind_values (0.00s)
    	CLI_test.go:38: expected a blind alert to be scheduled
```

## 成功させるのに十分なコードを書く

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

テストを通過させるには、`BlindAlerter`を好きな引数で呼び出せば十分です。

```go
func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5*time.Second, 100)
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

次は、5人のプレーヤーの場合に期待するアラートがすべてスケジュールされているかを確認したいと思います。

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

テーブルベースのテストはここではうまく機能し、要件が何であるかを明確に示しています。テーブルを順に処理し、`SpyBlindAlerter`をチェックして、アラートが正しい値でスケジュールされているかどうかを確認します。

## テストを実行してみます

次のような失敗がたくさん表示されるはずです。

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

## 成功させるのに十分なコードを書く

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

`PlayPoker`が少し読みやすくなるように、アラートのスケジュール処理をメソッドにカプセル化できます。

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

最後に、テストは少し不格好に見えます。同じもの（`ScheduledAlert`）を表す2つの匿名構造体があります。それを新しい型にリファクタリングして、比較用のヘルパーをいくつか作りましょう。

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

型に`String()`メソッドを追加したので、テストが失敗した場合にうまく表示されます。

新しい型を使用するようにテストを更新します。

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

`assertScheduledAlert`は自分で実装してみてください。

ここまでかなりの時間をかけてテストを書いてきましたが、アプリケーションに組み込まずにいたのは少し行儀が悪かったですね。これ以上の要件に取り掛かる前に、その問題に取り組みましょう。

アプリを実行しようとしても、`NewCLI`への引数が足りないというエラーでコンパイルできません。

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

`struct`だけでなく、どんな _型_ でもインターフェースを実装できることを思い出してください。関数が1つだけ定義されたインターフェースを公開するライブラリを作る場合、`MyInterfaceFunc`型もあわせて公開するのはよくあるイディオムです。

この型は、インターフェースを実装する`func`になります。こうすることで、インターフェースの利用者は関数だけでインターフェースを実装できます。
空の`struct`型を作成する必要はありません。

次に、関数と同じシグネチャを持つ関数`StdOutAlerter`を作成し、`time.AfterFunc`を使用して、それを`os.Stdout`に出力するようにスケジュールします。

この動作を確認するには、`NewCLI`を作成する`main`を更新します。

```go
poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter)).PlayPoker()
```

実行する前に、`CLI`の`blindTime`の増分を10分ではなく10秒に変更すると、実際の動作を確認できます。

期待どおり、10秒ごとにブラインドの値が出力されるはずです。また、CLIに`Shaun wins`と入力すれば、これまでどおり期待したとおりにプログラムが停止することにも注目してください。

ゲームは常に5人でプレイされるとは限らないため、ゲームを開始する前に、プレーヤーの人数を入力するようユーザーに促す必要があります。

## 最初にテストを書く

プレーヤーの人数の入力を促していることを確認するために、StdOutに書き込まれた内容を記録します。
これはもう何度かやってきました。`os.Stdout`は`io.Writer`なので、テストで依存性注入を使って`bytes.Buffer`を渡せば、コードが何を書き込むかを確認できます。

このテストでは、他の協調オブジェクトについてはまだ気にしないので、テストファイルにダミーをいくつか用意しました。

`CLI`の依存関係が4つになったことには少し警戒すべきでしょう。
責任を持ちすぎ始めているように感じられます。とりあえず今はこのままにして、この新しい機能を追加する中でリファクタリングの方向性が見えてくるか様子を見ましょう。

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

他のテスト用に`dummyStdOut`を追加します。

新しいテストは次のように失敗するはずです。

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

最後に、ゲームの開始時にプロンプトを書き込むことができます。

```go
func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, "Please enter the number of players: ")
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}
```

## リファクタリング♪

プロンプトの文字列が重複しているので、定数に抽出しましょう。

```go
const PlayerPrompt = "Please enter the number of players: "
```

テストコードと`CLI`の両方でこれを使用してください。

次に、数値を送り込んで、それを取り出す必要があります。
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

うわっ！変更がたくさんありますね。

* StdInのダミーを削除し、代わりに、7を入力するユーザーを表すモックバージョンを送信します。
* また、ブラインドアラーターのダミーも削除して、プレーヤー数がスケジュールに影響を与えていることを確認できるようにしました
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

これを動かすためなら、どんな罪を犯しても構わないことを思い出してください。
動くソフトウェアができたら、これから作る散らかったコードのリファクタリングに取り掛かればよいのです！

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
* そして、`scheduleBlindAlerts`がプレーヤーの人数を受け取るように変更します。次に、ブラインド額を順に処理するときに`blindTime`に加算する`blindIncrement`の時間を計算します

新しいテストは成功するようになりましたが、ゲームの開始時にユーザーが数値を入力しないとシステムが動作しなくなったため、他の多くのテストが失敗しています。
ユーザー入力を変更して、数値とそれに続く改行が追加されるようにテストを修正する必要があります。
（これは、現在のアプローチのさらなる欠陥を浮き彫りにしています）。

## リファクタリング♪

なんだかひどい感じですよね？**テストの声に耳を傾けましょう**。

* いくつかのアラートをスケジュールしていることをテストするために、4つの異なる依存関係を設定しました。システム内のある _もの_ に依存関係がたくさんある場合、それはそのものが多くのことをやりすぎていることを意味します。視覚的には、テストが雑然としていることがわかります。
* 私には、**ユーザー入力の読み取りと、実行したいビジネスロジックとの間をより明確に抽象化する必要があるように感じます**
* より良いテストは、_このユーザー入力が与えられたとき、新しい型`Game`を正しいプレーヤー数で呼び出すかどうか_ です。
* 次に、スケジューリングのテストを新しい`Game`のテストに抽出します。

まず`Game`に向けてリファクタリングを進めても、テストは引き続き成功するはずです。望む構造の変更が済んだら、新しい関心の分離を反映するようにテストをリファクタリングする方法を考えられます

リファクタリングで変更を加えるときは、変更をできるだけ小さくして、テストを再実行し続けるようにしてください。

まず自分で試してみてください。`Game`が提供するものと`CLI`が行うべきことの境界について考えてみましょう。

今のところは、`NewCLI`の外部インターフェースを**変更しないでください**。
テストコードとクライアントコードを同時に変更すると、一度に扱うことが多すぎて、何かを壊してしまいかねないからです。

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

* 何人がプレイするかを指定して`Game`を`Start`したい
* 勝者を宣言して`Game`を`Finish`したい

新しい`Game`型がこれをカプセル化します。

この変更で、`BlindAlerter`と`PlayerStore`を`Game`に渡すようにしました。
今後は`Game`がアラートと結果の保存を担当するからです。

これで`CLI`が担うのは次のことだけになりました。

* 既存の依存関係を使用して`Game`を構築します（これは次にリファクタリングします）
* ユーザー入力を`Game`のメソッド呼び出しとして解釈する

テストが長時間失敗したままになるような「大きな」リファクタリングは避けたいところです。
そうした状態はミスの可能性を高めるからです。
（大規模なチームや分散したチームで作業している場合、これは特に重要です）

まず最初に、`Game`をリファクタリングして、`CLI`に注入します。
そのためのテストへの変更は最小限にとどめ、その後、テストをユーザー入力の解析とゲーム管理というテーマに分割する方法を考えます。

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
依存関係が減り、_依存関係のリストが全体的な設計目標を反映しています_。つまり、CLIは入出力に専念し、ゲーム固有のアクションは`Game`に委任するということです。

コンパイルしてみると問題があります。これらの問題は自分で修正できるはずです。
今すぐ`Game`のモックを作成する必要はありません。すべてをコンパイルしてテストをグリーンにするために、_本物の_ `Game`を初期化するだけで構いません。

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

テストを修正して再びグリーンに戻すのに大した労力はかからないはずです（それが狙いです！）。
次の段階の前に`main.go`も修正するようにしてください。

```go
// main.go
game := poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
cli := poker.NewCLI(os.Stdin, os.Stdout, game)
cli.PlayPoker()
```

`Game`を抽出したので、ゲーム固有のアサーションをCLIとは別のテストに移動する必要があります。

これは、`CLI`のテストを依存関係を減らしてコピーするだけの作業です。

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

ゲーム終了時のテストも忘れずに移動してください。

ゲームロジックのテストを移し終えたら、
CLIのテストを簡略化して、意図した責任がより明確に反映されるようにできます。

* ユーザー入力を処理し、必要に応じて`Game`のメソッドを呼び出します
* 出力を送信
* 重要なのは、ゲームが実際にどう動くかについては何も知らないこと

これを行うには、`CLI`が具体的な`Game`型に依存せず、代わりに`Start(numberOfPlayers)`と`Finish(winner)`を持つインターフェースを受け取るようにする必要があります。そうすれば、その型のスパイを作成して、正しい呼び出しが行われたことを確認できます。

ここで、名前付けが時に厄介であることに気づきます。`Game`の名前を`TexasHoldem`に変更します（これが現在プレイしているゲームの _種類_ であるため）。新しいインターフェースを`Game`と呼ぶことにします。こうすることで、CLIは実際にプレイしているゲームや、`Start`・`Finish`したときに何が起こるかを知らない、という考え方に忠実でいられます。

```go
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}
```


`CLI`内の`*Game`への参照をすべて`Game`（新しいインターフェース）に置き換えます。
いつものように、リファクタリングしている間はテストを再実行し続けて、すべてがグリーンであることを確認します。

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

関心をきれいに分離できたので、`CLI`でのIOに関するエッジケースのチェックが簡単になるはずです。

プレーヤー数の入力を求められたときにユーザーが非数値を入力するシナリオに対処する必要があります。

その場合、コードはゲームを開始せず、ユーザーにわかりやすいエラーを出力して終了する必要があります。


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

`GameSpy`に、`Start`が呼び出された場合にのみ設定されるフィールド`StartCalled`を追加する必要があります

## テストを実行してみます

```
=== RUN   TestCLI/it_prints_an_error_when_a_non_numeric_value_is_entered_and_does_not_start_the_game
    --- FAIL: TestCLI/it_prints_an_error_when_a_non_numeric_value_is_entered_and_does_not_start_the_game (0.00s)
        CLI_test.go:62: game should not have started
```

## 成功させるのに十分なコードを書く

`Atoi`を呼び出している箇所で、エラーをチェックするだけです

```go
numberOfPlayers, err := strconv.Atoi(cli.readLine())

if err != nil {
	return
}
```

次に、ユーザーに何を間違えたのかを知らせる必要があるので、`stdout`に出力される内容をアサートします。

## 最初にテストを書く

`stdout`に出力された内容のアサートは前にもやったので、今のところはそのコードをコピーできます

```go
gotPrompt := stdout.String()

wantPrompt := poker.PlayerPrompt + "you're so silly"

if gotPrompt != wantPrompt {
	t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
}
```

`stdout`に書き込まれるものは _すべて_ 保存しているので、`poker.PlayerPrompt`も引き続き期待されます。
そのうえで、追加で何かが出力されていることを確認するだけです。
今のところ正確な文言にはあまりこだわらず、リファクタリングのときに対処します。

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

次に、メッセージを`PlayerPrompt`のような定数にリファクタリングします

```go
wantPrompt := poker.PlayerPrompt + poker.BadPlayerInputErrMsg
```

より適切なメッセージを入れます

```go
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
```

最後に、`stdout`に送信された内容に関するテストがかなり冗長なので、整理するためのアサート関数を作成しましょう。

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

さまざまな数のメッセージをアサートする必要があるので、ここでは可変長引数の構文（`...string`）が便利です。

このヘルパーは、ユーザーに送信されるメッセージに対してアサートする両方のテストで使用します。

`assertX`関数を使えば改善できるテストがいくつかあるので、テストを読みやすく整理してリファクタリングの練習をしてみてください。

少し時間をとって、これまでに書いてきたテストの価値について考えてみてください。
必要以上のテストは不要だということを思い出してください。
いくつかのテストをリファクタリング/削除しても、_すべてが動作するという確信を持ち続けられるでしょうか_？

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

これでテストはCLIの主な機能を反映するようになりました。CLIは、何人がプレイするか、誰が勝ったかについてのユーザー入力を読み取ることができ、プレーヤー数に不正な値が入力された場合にも対処できます。
こうすることで、`CLI`が何をするのか、そして何をしないのかが読者にとって明確になります。

ユーザーが`Ruth wins`の代わりに`Lloyd is a killer`と入力したらどうなるでしょうか？

このシナリオのテストを書いて成功させ、この章を締めくくりましょう。

## まとめ

### プロジェクトの簡単な要約

ここまでの5つの章で、かなりの量のコードを少しずつTDDで書いてきました

* コマンドラインアプリケーションとWebサーバーの2つのアプリケーションがあります。
* これらのアプリケーションは両方とも、勝者を記録するために`PlayerStore`に依存しています
* Webサーバーは、最も多くのゲームに勝っている人のリーグテーブルを表示することもできます
* コマンドラインアプリは、現在のブラインドの値を追跡することで、プレーヤーがポーカーをプレイするのを手助けします。

### time.Afterfunc

特定の時間が経過した後に関数呼び出しをスケジュールするための、とても便利な方法です。[`time`のドキュメントを読む](https://golang.org/pkg/time/)ことに時間を投資する価値は十分にあります。時間を節約できる便利な関数やメソッドがたくさん用意されています。

私のお気に入りをいくつか紹介します。

* 期間が終了すると、`time.After(duration)`は`chan Time`を返します。したがって、特定の時間の後に何かをしたい場合は、これが役立ちます。
* `time.NewTicker(duration)`は、チャネルを返すという点で上記と同様の `Ticker`を返しますが、こちらは1回だけではなく、期間ごとに「ティック」します。これは、`N duration`ごとにコードを実行する場合に非常に便利です。

### 関心の分離のさらなる良い例

一般的には、ユーザー入力と応答の処理の責任をドメインコードから分離することをお勧めします。
コマンドラインアプリケーションとWebサーバーで確認できます。

テストが乱雑になりました。アサーションが多すぎ（この入力をチェックし、これらのアラートをスケジュールし、など）、依存関係も多すぎました。
散らかっていることを視覚的に確認できました。 **テストに耳を傾けることはとても重要です**。

* テストが乱雑に見える場合は、リファクタリングしてみてください。
* これを行ってもまだ混乱している場合は、デザインの欠陥を指摘している可能性が高いです。
* これはテストの真の強みの1つです。

テストも本番コードも少し雑然としていましたが、テストに支えられて自由にリファクタリングできました。

これらの状況に陥ったときは、常に小さなステップを踏んで、変更のたびにテストを再実行することを忘れないでください。

テストコード _と_ 本番コードの両方を同時にリファクタリングするのは危険でした。
そのため、インターフェースを変更せずに、最初に本番コードをリファクタリングしました（現在の状態では、テストを大幅に改善することはできませんでした）。
こうすることで、変更している間もできる限りテストに頼ることができました。_その後で_、設計が改善されてからテストをリファクタリングしました。

リファクタリング後、依存関係のリストは設計目標を反映するようになりました。
これは、意図を文書化することが多いという点で、DIのもう1つの利点です。グローバル変数に依存すると、責任が非常に不明確になります。

## インターフェースを実装する関数の例

1つのメソッドでインターフェースを定義するとき、ユーザーが関数だけでインターフェースを実装できるように、それを補完するために`MyInterfaceFunc`型を定義することを検討することができます。

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

こうすることで、あなたのライブラリを利用する人は、関数だけであなたのインターフェースを実装することができます。彼らは [型変換](https://go.dev/tour/basics/13) を使って関数を `BlindAlerterFunc` に変換し、それを BlindAlerter として使うことができます (`BlindAlerterFunc` は `BlindAlerter` を実装しているので)。

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

ここでは、非常にシンプルな「ブログ」を実装するHTTPハンドラーを作成しました。このハンドラーでは、URLパスをキーとして、マップに格納された投稿を表示します。
