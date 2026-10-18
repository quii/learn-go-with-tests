---
description: スタブ・モック
---

# スタブ・モック

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/mocking)**

3からカウントダウンし、各数値を（1秒の間隔を置いて）1行ずつ表示し、ゼロに達したら「Go!」と表示して終了するプログラムを書くよう依頼されました。

```
3
2
1
Go!
```

これに取り組むには、`Countdown`という関数を作成します。この関数を`main`プログラム内に配置して、次のようにします。

```go
package main

func main() {
	Countdown()
}
```

これはかなり簡単なプログラムですが、完全にテストするには、いつものように反復的、テストドリブンのアプローチを取る必要があります。

反復的とはどういう意味でしょうか？ _役に立つソフトウェア_ を手にするために、できる限り小さなステップを踏むようにするということです。

少しハックすれば理論上は動くはずのコードに長い時間を費やしたくはありません。開発者が深みにはまる（うさぎの穴に落ちる）のは、たいていそういうときだからです。 **要件をできる限り小さくスライスして、 _動作するソフトウェア_ を手にできるようにすることは重要なスキルです。**

作業を分割して反復する方法は次のとおりです。

* 3を表示する
* 3、2、1 と Go! を表示する
* 各行の間で1秒待ちます

## 最初にテストを書く

私たちのソフトウェアはstdoutに出力する必要があり、依存性注入（DI）を使用してこれをテストする方法をDIセクションで確認しました。

```go
func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	Countdown(buffer)

	got := buffer.String()
	want := "3"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

`buffer`のようなものに不慣れな場合は、[前のセクション](dependency-injection.md)をもう一度お読みください。

`Countdown`関数でどこかにデータを書き込む必要があることはわかっています。それをGoのインターフェースとして表す事実上の標準が`io.Writer`です。

* `main`で`os.Stdout`に送信して、ターミナルに出力されたカウントダウンをユーザーに表示します。
* テストでは、`bytes.Buffer`に送信して、生成されるデータをテストでキャプチャできるようにします。

## テストを実行してみます

`./countdown_test.go:11:2: undefined: Countdown`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`Countdown`を定義します

```go
func Countdown() {}
```

再試行

```
./countdown_test.go:11:11: too many arguments in call to Countdown
    have (*bytes.Buffer)
    want ()
```

コンパイラが関数のシグネチャをどうすればよいか教えてくれているので、そのとおりに更新します。

```go
func Countdown(out *bytes.Buffer) {}
```

`countdown_test.go:17: got '' want '3'`

パーフェクト！

## 成功させるのに十分なコードを書く

```go
func Countdown(out *bytes.Buffer) {
	fmt.Fprint(out, "3")
}
```

`fmt.Fprint`を使用しています。これは、`io.Writer`（`*bytes.Buffer`など）を受け取り、それに`string`を送信します。テストは成功するはずです。

## リファクタリング♪

`*bytes.Buffer`は機能しますが、代わりに汎用インターフェースを使用する方がよいことはわかっています。

```go
func Countdown(out io.Writer) {
	fmt.Fprint(out, "3")
}
```

テストを再実行すると、テストに合格するはずです。

仕上げとして、関数を`main`に組み込みましょう。そうすれば動作するソフトウェアが手に入り、着実に前進していることを実感できます。

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func Countdown(out io.Writer) {
	fmt.Fprint(out, "3")
}

func main() {
	Countdown(os.Stdout)
}
```

プログラムを実行して、自分の仕事ぶりに感嘆してください。

ささいなことに思えるかもしれませんが、このアプローチはどんなプロジェクトにもお勧めします。
**機能を薄くスライスして切り出し、テストに裏打ちされた状態でエンドツーエンドに動作させましょう。**

次に、2、1と表示してから「Go!」を表示するようにします。

## 最初にテストを書く

配管全体を正しく機能させるために投資することで、ソリューションを安全かつ簡単に反復できます。
すべてのロジックがテストされるので、プログラムが機能していることを確認するために、プログラムを停止して再実行する必要がなくなります。

```go
func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	Countdown(buffer)

	got := buffer.String()
	want := `3
2
1
Go!`

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

バックティック構文は`string`を作成するもう1つの方法ですが、改行などを含めることができるので、テストに最適です。

## テストを実行してみます

```
countdown_test.go:21: got '3' want '3
        2
        1
        Go!'
```

## 成功させるのに十分なコードを書く

```go
func Countdown(out io.Writer) {
	for i := 3; i > 0; i-- {
		fmt.Fprintln(out, i)
	}
	fmt.Fprint(out, "Go!")
}
```
`i--`で逆にカウントする`for`ループを使用し、`fmt.Fprintln`を使用して、番号と改行文字を続けて`out`に出力します。最後に`fmt.Fprint`を使用して「Go!」を送信します。

## リファクタリング♪

いくつかのマジックバリューを名前付き定数にリファクタリングする以外に、リファクタリングすることは多くありません。

```go
const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
	}
	fmt.Fprint(out, finalWord)
}
```

ここでプログラムを実行すると、目的の出力が得られるはずですが、1秒ずつ間を置く劇的なカウントダウンにはなっていません。

Goでは、`time.Sleep`でこれを実現できます。コードに追加してみてください。

```go
func Countdown(out io.Writer) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		time.Sleep(1 * time.Second)
	}

	fmt.Fprint(out, finalWord)
}
```

プログラムを実行すると、期待どおりに機能します。

## モック

テストは引き続き成功し、ソフトウェアは意図したとおりに機能しますが、いくつかの問題があります。

* テストの実行には3秒かかります。
  * ソフトウェア開発に関する先進的な記事はどれも、迅速なフィードバックループの重要性を強調しています。
  * **遅いテストは開発者の生産性を台無しにします**。
  * 要件がさらに洗練され、より多くのテストが必要になると想像してください。 `Countdown`のすべての新しいテストのテスト実行に3秒が追加されることに満足できるでしょうか？
* 関数の重要なプロパティはテストしていません。

`Sleep`への依存があるので、テストで制御できるようにそれを抽出する必要があります。

`time.Sleep`を _モック_ できれば、 _依存性注入_ を使って「実際の」`time.Sleep`の代わりにそれを使い、**呼び出しをスパイ**してアサーションを行えます。

## 最初にテストを書く

依存関係をインターフェースとして定義しましょう。これにより、`main`で「実際の」Sleeperを使用し、テストで _spy sleeper_ を使用できるようになります。インターフェースを使用することにより、`Countdown`関数はこれを意識せず、呼び出し側に柔軟性を追加します。

```go
type Sleeper interface {
	Sleep()
}
```

`Countdown`関数はスリープ時間の長さには責任を持たない、という設計上の判断をしました。
これにより、少なくとも今のところはコードが少し単純になり、関数の利用者がスリープの長さを好きなように設定できるようになります。

次に、テストで使用するために**モック**を作成する必要があります。

```go
type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}
```

_Spies_ は、依存関係の使用方法を記録できるモックの一種です。送信された引数、それが呼び出された回数などを記録できます。この例では、`Sleep()`が呼び出された回数を追跡しているので、テストで確認できます。

テストを更新してスパイを依存関係として注入し、スリープが3回呼び出されたことをアサートします。

```go
func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}
	spySleeper := &SpySleeper{}

	Countdown(buffer, spySleeper)

	got := buffer.String()
	want := `3
2
1
Go!`

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}

	if spySleeper.Calls != 3 {
		t.Errorf("not enough calls to sleeper, want 3 got %d", spySleeper.Calls)
	}
}
```

## テストを実行してみます

```
too many arguments in call to Countdown
    have (*bytes.Buffer, *SpySleeper)
    want (io.Writer)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`Sleeper`を受け入れるには`Countdown`を更新する必要があります

```go
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		time.Sleep(1 * time.Second)
	}

	fmt.Fprint(out, finalWord)
}
```

もう一度試すと、同じ理由で`main`はコンパイルできなくなります

```
./main.go:26:11: not enough arguments in call to Countdown
    have (*os.File)
    want (io.Writer, Sleeper)
```

必要なインターフェースを実装する「実際の」スリーパーを作成しましょう

```go
type DefaultSleeper struct{}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}
```

それを実際のアプリケーションで使用することができます

```go
func main() {
	sleeper := &DefaultSleeper{}
	Countdown(os.Stdout, sleeper)
}
```

## 成功させるのに十分なコードを書く

テストはコンパイルされていますが、注入された依存関係ではなく、まだ`time.Sleep`を呼び出しているため、成功していません。修正しましょう。

```go
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		sleeper.Sleep()
	}

	fmt.Fprint(out, finalWord)
}
```

テストに合格し、もう3秒はかからなくなるはずです。

### まだいくつかの問題

テストしていない重要なプロパティがまだあります。

`Countdown` は次の表示の前にそれぞれスリープする必要があります。

例：

* `Print N`
* `Sleep`
* `Print N-1`
* `Sleep`
* `Print Go!`
* etc

私たちの最新の変更は、それが3回スリープしたとアサートしているだけですが、それらのスリープは正しい順序で起こっていない可能性があります。

テストを書いていて、そのテストで十分な確信が得られているか自信がない場合は、とにかくコードを壊してみましょう！
（ただし、その前に変更をソース管理にコミットしておいてください）。

コードを次のように変更します


```go
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
	}

	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
	}

	fmt.Fprint(out, finalWord)
}
```

テストを実行すると、実装が間違っていてもテストは成功するはずです。

新しいテストで再度スパイを使用して、操作の順序が正しいことを確認してみましょう。

2つの異なる依存関係があり、それらのすべての操作を1つのリストに記録したいと考えています。
そこで、 _両方を兼ねるスパイを1つ_ 作成します。

```go
type SpyCountdownOperations struct {
	Calls []string
}

func (s *SpyCountdownOperations) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

const write = "write"
const sleep = "sleep"
```

`SpyCountdownOperations`は`io.Writer`と`Sleeper`の両方を実装し、すべての呼び出しを1つのスライスに記録します。このテストでは操作の順序のみを考慮しているため、名前付き操作のリストとして記録するだけで十分です。

テストスイートにサブテストを追加して、スリープと表示が希望する順序で動作することを確認します。

```go
t.Run("sleep before every print", func(t *testing.T) {
	spySleepPrinter := &SpyCountdownOperations{}
	Countdown(spySleepPrinter, spySleepPrinter)

	want := []string{
		write,
		sleep,
		write,
		sleep,
		write,
		sleep,
		write,
	}

	if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
		t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
	}
})
```

これでこのテストは失敗するはずです。テストを直すために、`Countdown`を元の状態に戻しましょう。

`Sleeper`をスパイする2つのテストができたので、テストをリファクタリングして、1つは表示されているものをテストし、もう1つは表示の合間に確実にスリープするようにします。
最後に、最初のスパイは使用されなくなったので削除できます。

```go
func TestCountdown(t *testing.T) {

	t.Run("prints 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		Countdown(buffer, &SpyCountdownOperations{})

		got := buffer.String()
		want := `3
2
1
Go!`

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("sleep before every print", func(t *testing.T) {
		spySleepPrinter := &SpyCountdownOperations{}
		Countdown(spySleepPrinter, spySleepPrinter)

		want := []string{
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
			t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
		}
	})
}
```

これで、関数とその2つの重要な性質が適切にテストされました。

## スリーパーを構成可能に拡張

`Sleeper`が構成可能であることは素晴らしい機能です。これは、メインプログラムでスリープ時間を調整できることを意味します。

### 最初にテストを書く

まず、設定とテストに必要なものを受け入れる `ConfigurableSleeper`の新しいタイプを作成しましょう。

```go
type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}
```

スリープ時間を設定するために`duration`を使用し、スリープ関数を渡す方法として`sleep`を使用しています。 `sleep`のシグネチャは`time.Sleep`のシグネチャと同じであり、実際の実装で`time.Sleep`を使用して、テストで次のスパイを使用できます。

```go
type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) SetDurationSlept(duration time.Duration) {
	s.durationSlept = duration
}
```

スパイを配置したら、構成可能なスリーパーの新しいテストを作成できます。

```go
func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second

	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{sleepTime, spyTime.SetDurationSlept}
	sleeper.Sleep()

	if spyTime.durationSlept != sleepTime {
		t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
	}
}
```

このテストには何も新しいものはなく、以前のモックを使ったテストと非常によく似た準備になっています。

### テストを実行してみます

```
sleeper.Sleep undefined (type ConfigurableSleeper has no field or method Sleep, but does have sleep)

```

`ConfigurableSleeper`で作成された`Sleep`メソッドがないことを示す非常に明確なエラーメッセージが表示されます。

### テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func (c *ConfigurableSleeper) Sleep() {
}
```

新しい`Sleep`関数を実装したので、今度はテストが失敗します。

```
countdown_test.go:56: should have slept for 5s but slept for 0s
```

### 成功させるのに十分なコードを書く

ここで行う必要があるのは、`ConfigurableSleeper`の`Sleep`関数を実装することだけです。

```go
func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}
```

この変更により、すべてのテストが再び成功するはずです。メインプログラムはまったく変わっていないのに、なぜこんな手間をかけたのか不思議に思うかもしれません。
次のセクションを読めば明らかになるでしょう。

### クリーンアップとリファクタリング

最後に行う必要があるのは、メイン関数で実際に `ConfigurableSleeper`を使用することです。

```go
func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
	Countdown(os.Stdout, sleeper)
}
```

テストとプログラムを手動で実行すると、すべての動作が同じであることがわかります。

`ConfigurableSleeper`を使用しているので、`DefaultSleeper`実装を削除しても安全です。
これでプログラムは完成し、任意の長さのカウントダウンに対応できる、より[汎用的な](https://stackoverflow.com/questions/19291776/whats-the-difference-between-abstraction-and-generalization)Sleeperが手に入りました。

## しかし、モックを行うのは悪ではないのか？

モックは悪だと聞いたことがあるかもしれません。ソフトウェア開発におけるあらゆるものと同じく、[DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)と同様に悪用されることもあります。

悪い状態に陥るのは、たいてい _テストの声に耳を傾けず_ 、 _リファクタリングの段階をおろそかにしている_ ときです。

モックのコードが複雑になってきたり、何かをテストするために多くのものをモックしなければならなくなったりしたら、その嫌な感じに _耳を傾け_ 、コードについて考え直すべきです。それは通常、次のような兆候です。

* テスト対象が多くのことをやりすぎている（そのため、モックすべき依存関係が多すぎる）
  * モジュールを分割して、それぞれが担うことを減らす
* 依存関係が細かすぎる
  * これらの依存関係のいくつかを、意味のある1つのモジュールに統合する方法を考える
* テストが実装の詳細に関心を持ちすぎている
  * 実装ではなく期待される振る舞いのテストを優先する

通常、モックが多いことは、コードの _抽象化が悪い_ ことを示しています。

**人々はここにTDDの弱点を見ますが、それは実際には長所です**。多くの場合、貧弱なテストコードは悪い設計の結果です。より前向きに言えば、適切に設計されたコードはテストが簡単です。

### でも、モックやテストのせいで生活が苦しくなってきました!

この状況に遭遇したことがありますか？

* リファクタリングを行いたい
* これを行うには、多くのテストを変更することになります
* TDDに疑問を抱き、Mediumに「モックは有害である」というタイトルの記事を投稿する

これは通常、 _実装の詳細_ をテストしすぎていることを示しています。
システムの動作にとって実装が本当に重要でない限り、テストが _有用な振る舞い_ をテストするようにしてください。

正確に _どのレベル_ でテストすべきかを判断するのは難しい場合もありますが、私が従うようにしている考え方とルールをいくつか紹介します。

* **リファクタリングの定義では、コードは変更されますが、動作は同じです**。リファクタリングをすると決めたなら、理論上はテストを一切変更せずにコミットできるはずです。だからテストを書くときは自問してください
  * 自分が求める振る舞いをテストしているのか、それとも実装の詳細をテストしているのか？
  * このコードをリファクタリングする場合、テストに多くの変更を加える必要がありますか？
* Goではプライベート関数をテストできますが、プライベート関数はパブリックな振る舞いを支える実装の詳細なので、私は避けたいと思います。パブリックな振る舞いをテストしましょう。Sandi Metzはプライベート関数を「安定性が低い」ものと表現しており、テストをそれらに結合させたくはありません。
* テストが**3つ以上のモックで動作している場合、それは危険信号**であるように感じます（設計を見直すべきときです）
* スパイは注意して使用してください。スパイを使用すると、作成中のアルゴリズムの内部を確認できます。これは非常に便利ですが、テストコードと実装の間の結合がより緊密になることを意味します。 **スパイするなら、その詳細が本当に気にすべきものであることを確認してください**

#### モックフレームワークを使えばいいのでは？

モックには魔法は不要で、比較的シンプルです。フレームワークを使うと、モックが実際よりも複雑に見えてしまうことがあります。この章では、次のことを得るために自動モックを使用していません。

* モックの方法をよりよく理解する
* インターフェースを実装する練習をする

共同プロジェクトでは、モックを自動生成することに価値があります。チームでは、モック生成ツールがテストダブルの一貫性を明文化します。これにより、一貫性のない書き方のテストダブルを避けることができ、それが一貫性のない書き方のテストにつながることも避けられます。

モックジェネレーターは、インターフェースに対してテストダブルを生成するものだけを使うべきです。テストの書き方を過度に指示するツールや、多くの「魔法」を使うツールは、さっさと捨ててしまいましょう。

## まとめ

### TDDアプローチの詳細

* ささいとは言えない例に直面した場合は、問題を「薄い垂直スライス」に分解してください。深みにはまったり「ビッグバン」アプローチをとったりしないように、できるだけ早く _テストに裏打ちされた動作するソフトウェア_ を手にできるようにしてください。
* 動作するソフトウェアを手にしたら、必要なソフトウェアにたどり着くまで _小さなステップで反復する_ のが簡単になるはずです。

> 「反復開発を使用する場合？反復開発は、成功させたいプロジェクトでのみ使用する必要があります。」

マーティン・ファウラー（Martin Fowler）。

### モック

* **モックがなければ、コードの重要な部分がテストされないままになります**。私たちのケースでは、各表示の間にコードが一時停止することをテストすることはできませんが、他にも数え切れないほどの例があります。失敗する _可能性のある_ サービスを呼び出していますか？特定の状態でシステムをテストしたいですか？モックなしでこれらのシナリオをテストすることは非常に困難です。
* モックがないと、単純なビジネスルールをテストするためだけに、データベースや他のサードパーティの設定が必要になる場合があります。テストが遅くなり、**フィードバックループが遅くなる**可能性があります。
* 何かをテストするためにデータベースまたはWebサービスをスピンアップする必要があるため、そのようなサービスの信頼性の低さから、**壊れやすいテスト**になりがちです。

開発者がモックについて学ぶと、システムが _何をするか_ ではなく _どのように動くか_ の観点から、システムのあらゆる側面を過剰にテストすることが非常に簡単になります。 **テストの価値**と、それらが将来のリファクタリングに与える影響について常に注意してください。

モックに関するこの投稿では、モックの一種である**スパイ**のみを取り上げました。モックは「テストダブル」の一種です。

> [テストダブル（Test Double）とは、テスト目的で本番オブジェクトを置き換えるあらゆる場合を指す総称です。](https://martinfowler.com/bliki/TestDouble.html)

テストダブルには、スタブ、スパイ、そしてもちろんモックなど、さまざまな種類があります！詳細については[Martin Fowlerの記事](https://martinfowler.com/bliki/TestDouble.html)を参照してください。

## ボーナス - Go 1.23 のイテレータの例

Go 1.23で[イテレータが導入されました](https://tip.golang.org/doc/go1.23)。イテレータはさまざまな方法で使えますが、今回はカウントダウンする数値を逆順に返す`countdownFrom`イテレータを作ってみましょう。

カスタムイテレータの書き方に入る前に、まずその使い方を見てみましょう。数値からカウントダウンするために、かなり命令的に見えるループを書く代わりに、カスタムの`countdownFrom`イテレータを`range`することで、このコードをより表現力豊かにできます。

```go
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := range countDownFrom(3) {
		fmt.Fprintln(out, i)
		sleeper.Sleep()
	}

	fmt.Fprint(out, finalWord)
}
```

`countDownFrom`のようなイテレータを書くには、特定の方法で関数を書く必要があります。ドキュメントより：

    The “range” clause in a “for-range” loop now accepts iterator functions of the following types
        func(func() bool)
        func(func(K) bool)
        func(func(K, V) bool)

（「for-range」ループの「range」句は、上記の型のイテレータ関数を受け付けるようになりました。`K`と`V`はそれぞれキーと値の型を表します。）

今回のケースでは、キーはなく値だけです。Goには便利な型`iter.Seq[T]`も用意されており、これは`func(func(T) bool)`の型エイリアスです。

```go
func countDownFrom(from int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := from; i > 0; i-- {
			if !yield(i) {
				return
			}
		}
	}
}
```

これはシンプルなイテレータで、`from`から始めて逆順に数値をyieldします。まさに私たちのユースケースにぴったりです。
