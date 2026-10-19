---
description: testing/synctest で時間を再考する
---

# testing/synctest で時間を再考する

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/synctest)**

[時間の章](time.md)では、`time.AfterFunc` を使って、ポーカーの CLI に「ブラインドが上がりました」という通知をスケジュールする機能を追加しました。これをテストするのは厄介でした。`time.AfterFunc` は実際の時間が経過した後に、自分専用のゴルーチンでコールバックを実行します。また Go では関数同士を比較できないため、何がスケジュールされたのかを簡単に調べることができませんでした。

私たちはこれを、おなじみの道具である依存性注入で解決しました。`BlindAlerter` インターフェースを定義し、テストでは本物の実装を、スケジュールを依頼された内容を記録するだけのスパイに差し替えました。だいたい次のようなものです。

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type SpyBlindAlerter struct {
	Alerts []struct {
		At     time.Duration
		Amount int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.Alerts = append(s.Alerts, struct {
		At     time.Duration
		Amount int
	}{at, amount})
}
```

これは良い設計であり、これによって書けるテストは高速で信頼できます。しかし、それらのテストが実際に何をカバーしているのかをよく見てください。テストがアサートしているのは「`ScheduleAlertAt` が `10 * time.Minute` と `200` で呼び出された」といったことです。*本物の* アラーターが実行されることは一度もありません。実際に `time.AfterFunc` を呼び出して何かを出力する唯一のコード、つまり本当にバグが潜んでいる可能性のある部分は、どのテストでも実行されないのです。

これは見落としではなく、トレードオフです。`time.AfterFunc` を使う本物のアラーターをきちんとテストしようとすれば、テストが終わるまで実際に何分も待つか、時間をミリ秒単位まで縮めて、テストを実行しているマシンが忙しすぎて追いつけない、ということがないよう祈るかのどちらかになります。どちらも魅力的ではありません。そのため、これまではテストしてこなかったのです。

Go 1.25 からは、3つ目の選択肢があります。[`testing/synctest`](https://pkg.go.dev/testing/synctest) です。これを使うと、`time.Sleep` や `time.AfterFunc` などを使う、手を加えていない本物のコードを、時間を完全に制御できるテストの中で実行できます。待つ必要もなく、不安定さもなく、時間を縮めてうまくいくよう祈る必要もありません。

この章では、その本物のアラーターを一から作り、`synctest` でテストします。ただし、まずは `synctest` が何から私たちを救ってくれるのかを正確に確認して、その価値を実感してから使うことにしましょう。

## `testing/synctest` について最低限知っておくこと

`testing/synctest` は、関数を隔離された **バブル（bubble）** の中で実行します。そのバブルの中では、

* `time` パッケージは偽の時計を使います。この時計は 2000年1月1日 UTC の午前0時から始まり、前にしか進みません。
* 偽の時間は、バブル内のすべてのゴルーチンが **永続的にブロック（durably blocked）** されているときにだけ進みます。永続的なブロックとは、同じバブル内の別のゴルーチンによってしか解除できない形でブロックされている状態のことです。`time.Sleep`、バブル内で作成されたチャネルからのブロッキングな受信、`sync.Cond.Wait`、`sync.WaitGroup.Wait` はすべてこれに該当します。`sync.Mutex` のロックは該当しません。ミューテックスは通常、ごく短い時間しか保持されないからです。
* `synctest.Test(t, f)` は新しいバブルの中で `f` を実行し、その中で生成されたすべてのゴルーチンが終了するまで戻りません。バブルが永続的にブロックされ、先に進む手段がなくなった場合は、永遠にハングするのではなく、デッドロックとしてテストを失敗させます。
* `synctest.Wait()` は、バブル内の *他の* すべてのゴルーチンが永続的にブロックされるまで呼び出し元のゴルーチンをブロックし、その後で戻ります。これは、アサーションを行う前にバックグラウンドの処理を落ち着かせるための方法です。

始めるにはこれで十分です。注意すべき細かい点については、進めながらいくつか拾っていきます。

## 最初にテストを書く

`synctest` に手を伸ばす前に、まずは素直なところから始めましょう。期間と金額を与えると、その期間だけ待ってからどこかにメッセージを書き込む `BlindAlerter` を、実際の時間を使ってテストします。

```go
package poker

import (
	"bytes"
	"testing"
	"time"
)

func TestStdOutAlerter(t *testing.T) {
	out := &bytes.Buffer{}
	alerter := StdOutAlerter(out)

	alerter.ScheduleAlertAt(5*time.Second, 100)

	time.Sleep(6 * time.Second)

	want := "Blind is now 100"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}
```

スケジュールしたアラートよりも少し長く（5秒ではなく6秒）スリープしているのは、確認する前に、`time.AfterFunc` が生成したゴルーチンが実際に実行される時間を少し与えるためです。

## テストを実行してみます

まだプロダクションコードを書いていないので、これはコンパイルできません。

```
./blind_alerter_test.go:11:13: undefined: StdOutAlerter
```

## 成功させるのに十分なコードを書く

```go
package poker

import (
	"fmt"
	"io"
	"time"
)

// BlindAlerter schedules alerts for blind amounts.
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

// BlindAlerterFunc allows you to implement BlindAlerter with a function.
type BlindAlerterFunc func(duration time.Duration, amount int)

// ScheduleAlertAt is BlindAlerterFunc's implementation of BlindAlerter.
func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

// StdOutAlerter returns a BlindAlerterFunc that schedules alerts and prints them to out.
func StdOutAlerter(out io.Writer) BlindAlerterFunc {
	return func(duration time.Duration, amount int) {
		time.AfterFunc(duration, func() {
			fmt.Fprintf(out, "Blind is now %d", amount)
		})
	}
}
```

実行してみます。

```
=== RUN   TestStdOutAlerter
--- PASS: TestStdOutAlerter (6.00s)
PASS
ok  	github.com/quii/learn-go-with-tests/synctest/v1	6.191s
```

成功しました。ただし、アラート1つのために実際に6秒かかっています。実際のゲームでは11個のアラートをスケジュールし、中には1時間40分も離れているものもあります。テストを実行するたびにそれを待つ人はいないでしょう。そのため実際には、このテストはスキップされるか、本物とはなんとなく似ているだけの非現実的に短い時間で書かれることになります。`synctest` は、まさにこの問題を解決するために存在します。

## synctest の導入

同じテストをバブルで包み、「偽の時計」が何をしてくれるのかについて少し楽観的すぎた場合に何が起こるかを見てみましょう。

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

スリープを完全に取り除きました。きっと偽の時計がそれを肩代わりしてくれるはずですよね？ そうはなりません。

```
=== RUN   TestStdOutAlerter
    blind_alerter_test.go:19: got "", want "Blind is now 100"
--- FAIL: TestStdOutAlerter (0.00s)
FAIL
```

バブルの偽の時計は、独自のタイマーで勝手に進むわけではありません。バブル内の何かが時間の経過を待って永続的にブロックしたときにだけ進みます。何を待っているのかは依然として私たちが示す必要があります。ただ、その代償を実際の秒数で払う必要がなくなっただけです。スリープを元に戻しましょう。

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		time.Sleep(6 * time.Second)

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

```
=== RUN   TestStdOutAlerter
--- PASS: TestStdOutAlerter (0.00s)
PASS
ok  	github.com/quii/learn-go-with-tests/synctest/v2	0.133s
```

成功し、しかも一瞬です。バブル内の `time.Sleep(6 * time.Second)` は、実時間では何のコストもかかりません。これで完成のように見えます。習慣として、`-race` でも問題ないか念のため確認しておきましょう。

```
go test -race ./...
```

```
==================
WARNING: DATA RACE
Read at 0x00c00010e630 by goroutine 9:
  bytes.(*Buffer).String()
      /usr/local/go/src/bytes/buffer.go:77 +0x174
  github.com/quii/learn-go-with-tests/synctest/v2.TestStdOutAlerter.func1()
      blind_alerter_test.go:20 +0x15c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x164

Previous write at 0x00c00010e630 by goroutine 11:
  bytes.(*Buffer).grow()
      /usr/local/go/src/bytes/buffer.go:143 +0x354
  bytes.(*Buffer).Write()
      /usr/local/go/src/bytes/buffer.go:185 +0xb4
  fmt.Fprintf()
      /usr/local/go/src/fmt/print.go:225 +0x94
  github.com/quii/learn-go-with-tests/synctest/v2.TestStdOutAlerter.func1.BlindAlerterFunc.ScheduleAlertAt.TestStdOutAlerter.func1.StdOutAlerter.1.2()
      blind_alerter.go:26 +0x6c

Goroutine 11 (finished) created at:
  time.goFunc()
      /usr/local/go/src/time/sleep.go:215 +0x40
==================
--- FAIL: TestStdOutAlerter (0.00s)
    testing.go:1617: race detected during execution of test
FAIL
```

痛いですね。「Goroutine 11 (finished)」に注目してください。実行するたびに、書き込みは確かに読み込みの前に行われていました。機能的には、このテストがアサーションで失敗することは実際にはあり得ません。しかし、レース検出器が問うているのは「今回うまくいかなかったか？」ではありません。「それが起こり得ないことを *保証* するものは何かあるか？」と問うているのです。私たちのゴルーチンでの `time.Sleep` と、`time.AfterFunc` のコールバックを実行する専用のゴルーチンは、独立してスケジュールされる2つのゴルーチンにすぎません。「しばらくスリープした」ということは、どれだけ長くスリープしたとしても、別のゴルーチンが共有メモリへのアクセスを終えたことを何ら約束しません。スリープをどれだけ長くしても、これは直りません。`-race` を付けないときに失敗する頻度が下がるだけです。

少し立ち止まって考える価値があります。これとまったく同じデータ競合は、本物の `time.Sleep` を使った、最初の素朴な実時間版にもすでに潜んでいました。ただ、確認しようと思わなかっただけです。実行するだけで実際に6秒もかかるテストに対して、`-race` を何度も繰り返し実行する人がいるでしょうか？

## 最初にテストを書く

必要なのは本物の同期ポイントです。書き込みが *おそらく* 先に起こるようにするだけでなく、実際に先に起こったことを保証するものです。まさにそのためにあるのが `synctest.Wait()` です。

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		time.Sleep(6 * time.Second)
		synctest.Wait()

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

## リファクタリング♪

プロダクションコードには変更するところはありませんが、今度こそ `-race` のもとで繰り返し実行しても問題がないことを確かめましょう。

```
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.153s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.143s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.148s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.147s
ok  	github.com/quii/learn-go-with-tests/synctest/v3	1.143s
```

毎回クリーンです。とても良いですね。

## 最初にテストを書く

テストする価値のあることがもう1つあります。それは、まだ何も起きていないということです。実際の時間では、確信が持てる程度には長く、しかしテストがだらだらと長引かない程度のスリープ時間を当て推量で決める必要があります。synctest ならその当て推量は不要です。

```go
func TestStdOutAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		out := &bytes.Buffer{}
		alerter := StdOutAlerter(out)

		alerter.ScheduleAlertAt(5*time.Second, 100)

		synctest.Wait()
		if out.String() != "" {
			t.Errorf("did not expect anything to be printed yet, got %q", out.String())
		}

		time.Sleep(5 * time.Second)
		synctest.Wait()

		want := "Blind is now 100"
		if out.String() != want {
			t.Errorf("got %q, want %q", out.String(), want)
		}
	})
}
```

アラートをスケジュールした時点では、バブル内の他のものは何もしていないので、`Wait()` はすぐに戻るはずです。まだ *待つべき* ものが何もないので、`out` はまだ空のはずです。実行してみます。

```
=== RUN   TestStdOutAlerter
--- PASS: TestStdOutAlerter (0.00s)
PASS
```

成功しました。もう学んだとおり、`-race` でも確認しておきましょう。

```
==================
WARNING: DATA RACE
Write at 0x00c00011c630 by goroutine 11:
  bytes.(*Buffer).grow()
      /usr/local/go/src/bytes/buffer.go:143 +0x354
  bytes.(*Buffer).Write()
      /usr/local/go/src/bytes/buffer.go:185 +0xb4
  fmt.Fprintf()
      /usr/local/go/src/fmt/print.go:225 +0x94
  github.com/quii/learn-go-with-tests/synctest/v4.TestStdOutAlerter.func1.BlindAlerterFunc.ScheduleAlertAt.TestStdOutAlerter.func1.StdOutAlerter.1.2()
      blind_alerter.go:26 +0x6c

Previous read at 0x00c00011c630 by goroutine 9:
  bytes.(*Buffer).String()
      /usr/local/go/src/bytes/buffer.go:77 +0x164
  github.com/quii/learn-go-with-tests/synctest/v4.TestStdOutAlerter.func1()
      blind_alerter_test.go:18 +0x14c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x164

Goroutine 11 (running) created at:
  time.goFunc()
      /usr/local/go/src/time/sleep.go:215 +0x40
==================
--- FAIL: TestStdOutAlerter (0.00s)
    testing.go:1617: race detected during execution of test
FAIL
```

データ競合は18行目、つまり一番最初のチェックで起きています。「まだ待つべきものは何もない」という理由で、私たちが自信を持っていたチェックです。しかも毎回確実に起こります。

問題は、`Wait()` だけでは、偽の時計をどこまで進めてよいかを制限するものが何もないことです。その時点でバブル内に他のゴルーチンはなく、*私たち* が待つべき次の期限もありません。残っているのは、ランタイムのタイマーヒープの中に置かれた、アラートの5秒のタイマーだけです。そこでランタイムは、できる唯一の有益なことをします。処理を先に進めるためにタイマーを発火させ、コールバックを実行するゴルーチン11を生成するのです。それも、`Wait()` が戻るかどうかをまだ判断している *最中に* です。ゴルーチン11は私たちのチェックの前ではなく、チェックと本当に並行して実行されています。その書き込みが先に終わることもあれば、そうでないこともあります。そして両方のゴルーチンが触っている `bytes.Buffer` にはロックがないので、どちらにしても安全にしてくれるものは何もありません。

これは、他のデータ競合を直すのと同じ方法で直すこともできます。`out` を `sync.Mutex` で包むか、もっと軽い手段として `sync/atomic` の `atomic.Pointer[T]` を使うかです。しかし、`StdOutAlerter` が実際に何を求められているのかを見てください。メッセージを決めること *と*、それを出力するという副作用を実行することです。ポーカーのブラインドのアラートをスケジュールするのに、`io.Writer` について知る必要はまったくありません。それは呼び出し側の関心事です。これは、本書が何度も立ち返ってきた種類の緊張関係です。[テストが苦痛を引き起こしているなら、そのシグナルに耳を傾け、コードの設計について考えましょう](../questions-and-answers/http-handlers-revisited.md)。アラーターは時間になったらメッセージを生成するだけにして、それをどうするかは呼び出し側に決めさせましょう。

## リファクタリング♪

ゴルーチンの境界を越えてやり取りするためにこそ、チャネルは存在します。Go のことわざにあるように、[メモリを共有することで通信するのではなく、通信することでメモリを共有せよ（don't communicate by sharing memory; share memory by communicating）](https://go.dev/blog/codelab-share)。共有された `out` に書き込む代わりに、アラーターは完成したメッセージをチャネルに送ることができます。

```go
func NewAlerter() (BlindAlerterFunc, <-chan string) {
	alerts := make(chan string)

	scheduleAlertAt := func(duration time.Duration, amount int) {
		time.AfterFunc(duration, func() {
			alerts <- fmt.Sprintf("Blind is now %d", amount)
		})
	}

	return scheduleAlertAt, alerts
}
```

`BlindAlerter` と `BlindAlerterFunc` はそのままです。`StdOutAlerter`（示唆的なことに、もはや標準出力とは何の関係もなくなっていました）だけが姿を消し、代わりに、アラーターと受信用のチャネルの両方を返す `NewAlerter` が登場しました。これらのアラートを出力したいもの（たとえば `main`）は、そのチャネルを `range` でループして出力すればよいのです。それはもはやこのパッケージの問題ではありません。

テストもシンプルになります。

```go
func TestNewAlerter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		alerter, alerts := NewAlerter()

		alerter.ScheduleAlertAt(5*time.Second, 100)

		select {
		case got := <-alerts:
			t.Fatalf("did not expect an alert yet, got %q", got)
		default:
		}

		got := <-alerts
		want := "Blind is now 100"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
```

何がなくなったかに注目してください。`synctest.Wait()` はどこにもなく、共有された値を守るための独自の型もなく、`time.Sleep` もまったくありません。もはやどれも必要ないのです。

[select の章](../go-fundamentals/select.md)でおなじみの、`default` ケースを持つ `select` は決してブロックしません。準備のできたケースを取るか、すぐに `default` に抜けるかのどちらかです。この時点では仮想時間はまだゼロのままで、アラートまでは（偽の時間で）まるまる5秒あるので、同期すべきものは何もありません。当然、まだ何も届いていません。そして `got := <-alerts` にも後押しは必要ありません。これは普通のブロッキングな受信なので、バブルはまさに設計どおりに動作します。バブル内の誰もが待っているのはそのタイマーだけなので、偽の時間はタイマーが発火する瞬間までまっすぐ飛び、`time.AfterFunc` のゴルーチンを起こし、私たちの受信のブロックを解除します。

これを `-race` 付きで繰り返し実行してみてください。グリーンのままです。競合の対象となる共有メモリがもう残っていないからです。

これを、以前の章の `SpyBlindAlerter` のアプローチと比べてみましょう。あちらも、別の仕事のための道具としては今でも十分に優れています。あちらは、実際のタイミングをまったく気にせずに、*何が* スケジュールされるか（あるプレーヤー数に対して、正しい金額が正しいオフセットでスケジュールされているか？）を確認します（タイミングではなく計算をテストしているときに便利です）。`synctest` がその必要性を置き換えたわけではありません。純粋な「何を依頼されたか」のスパイでは必ず残ってしまうカバレッジの穴、つまり、スケジューリングの仕組みそのものが実際に動作するのか、という穴を埋めてくれるのです。

## まとめ

### 学んだこと

* `synctest.Test` と、隔離された偽の時計を持つバブルという考え方。その時計は勝手に進むのではなく、何かが永続的にブロックしたときにだけ進みます。
* 「永続的にブロック」：偽の時間を進めるための条件であり、`sync.Mutex` が意図的にそれに該当しない理由（一方でチャネルの受信は該当します）。
* 「十分に長く」スリープすることは、同期することと同じではありません。たっぷりと時間を取った、実際には常に正しく動くスリープであっても、順序を強制するものが何もなければ、それは本物のデータ競合です。
* `synctest.Wait()` は1回のチェックについてはそれを解決しますが、バブル内にそれを制限する他のものが何もない状態で呼び出すと、偽の時計が意図したよりも先まで進んでしまうことがあります。まさに「まだ何も起きていない」ケースをテストしたときに起きたことです。
* テストが単なるバグではなく設計上の問題を浮かび上がらせ、ロックに手を伸ばす代わりに設計を直したこと。

### 注意すべき落とし穴

* バブルのルート関数が戻ったときに、バブル内のゴルーチンがまだ永続的にブロックされていると、`synctest.Test` はハングするのではなくデッドロックとしてテストを失敗させます。バックグラウンドのゴルーチンが実際に終了するようにしてください。
* チャネル、タイマー、ティッカーは、それらが作成されたバブルに結び付けられています。バブルの外から使うとパニックします。
* 実際のネットワークやファイルの I/O は永続的なブロックではないため、`synctest` の偽の時計で直接駆動することはできません。インメモリの代替が必要な場合は、`net.Pipe` のようなものを使いましょう。
* `time.AfterFunc` のコールバックは専用のゴルーチンで実行され、チャネルが与えてくれるような同期の保証は一切ありません。コールバックが共有状態に直接触れる必要がある場合は、他の並行コードと同じく、その状態には独自のロックが必要です。

### 追加資料

* [Go blog: Testing concurrent code with testing/synctest](https://go.dev/blog/synctest)
* [`testing/synctest` パッケージのドキュメント](https://pkg.go.dev/testing/synctest)
* [Go 1.25 リリースノート](https://go.dev/doc/go1.25)

### この章の書き方について

これは本書で初めて、AI の支援（Claude）を受けて書かれた章です。そのことと、ここでの「支援」が実際に何を意味していたのかについて、率直にお伝えしておきたいと思います。それは「章の内容を説明したら章ができあがった」というものではなかったからです。

そのプロセスは、本書がずっと教えてきた TDD のループによく似ていました。推測するのではなく本物の `testing/synctest` のドキュメントとソースを調べ、本文を1語でも書く前に小さな使い捨てのプログラムでスパイクして主張を確認し、ドキュメントのあらゆる主張を鵜呑みにせず、実際に `go test -race` を実行して検証すべきものとして扱いました。この章の落とし穴のいくつか、とりわけ `Wait()` とレース検出器の相互作用は、何度か実際に実行するうちにテストが予想外の形で本当に失敗し、何を書くかを決める前にその理由を掘り下げる必要があったからこそ、ここに載っているのです。

設計そのものも途中で形を変えました。最初の草稿では `StdOutAlerter` が `io.Writer` に直接書き込んでいましたが、これはまさに、テストが苦痛になり始めたときに本書が常に異を唱えてきた類のものです。[テストが苦痛を引き起こしているなら、そのシグナルに耳を傾け、コードの設計について考えましょう](../questions-and-answers/http-handlers-revisited.md)。そこで私たちはそうしました。その結果が上のチャネルベースのバージョンであり、そのおかげで章はより短く、より良いものになりました。

ここにあるすべてのものはレビューと編集を経ており、何度も差し戻されました。トーンがしっくりこないとき、あるセクションが実際の教訓に見合う以上に膨らんでしまったとき、実際のテストの失敗を見せたほうがうまく伝わるところで説明が専門用語に頼ってしまったとき、などです。それでもまだ読みにくいところがあれば、それはツールではなく私の責任です。
