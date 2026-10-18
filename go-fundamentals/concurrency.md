---
description: 並行性
---

# 並行性

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/concurrency)**

状況はこうです。同僚がURLリストのステータスを確認する関数`CheckWebsites`を作成しました。

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = wc(url)
	}

	return results
}
```

チェックされた各URLをブール値に対応させたマップを返します。良好な応答の場合は`true`、不良応答の場合は`false`です。

単一のURLを取得してブール値を返す`WebsiteChecker`も渡す必要があります。
これは、すべてのWebサイトをチェックする機能によって使用されます。

[依存性の注入(DI)](dependency-injection.md)を使用すると、実際のHTTP呼び出しを行わずに関数をテストできるため、信頼性が高く、高速です。

ここに彼らが書いたテストがあります。

```go
package concurrency

import (
	"reflect"
	"testing"
)

func mockWebsiteChecker(url string) bool {
	return url != "waat://furhurterwe.geds"
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	want := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("wanted %v, got %v", want, got)
	}
}
```

この機能は運用中であり、数百のWebサイトをチェックするために使用されています。
しかし、あなたの同僚はそれが遅いという不満を持ち始めたので、彼らはあなたにそれを速くするのを手伝うように頼みました。

## テストを書く

ベンチマークを使用して、`CheckWebsites`の速度をテストし、変更の影響を確認してみましょう。

```go
package concurrency

import (
	"testing"
	"time"
)

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}

	for b.Loop() {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}
```

ベンチマークは、100個のURLのスライスを使用して`CheckWebsites`をテストし、`WebsiteChecker`の新しい偽の実装を使用します。`slowStubWebsiteChecker`は意図的に遅いです。
`time.Sleep`を使用して正確に20ミリ秒待機してからtrueを返します。

`go test -bench=.`を使用してベンチマークを実行する場合（またはWindows Powershellの場合`go test -bench="."`）。

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v0
BenchmarkCheckWebsites-4               1        2249228637 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v0        2.268s
```

`CheckWebsites`のベンチマークは2249228637ナノ秒で、約2.25秒です。

これをもっと速くしてみましょう。

### 成功させるのに十分なコードを書く

ようやく並行性について話すことができます。ここでの並行性とは、「複数の処理が進行中である」ことを意味します。これは私たちが毎日自然に行うことです。

たとえば、今朝私はお茶を作りました。
私はやかんを火にかけ、沸騰するのを待っている間に、冷蔵庫から牛乳を取り出し、食器棚からお茶を取り出し、お気に入りのマグカップを見つけ、ティーバッグをカップに入れました。
そして、やかんが沸騰したら、カップにお湯を注ぎました。

私が _しなかった_ のは、やかんを火にかけ、沸騰するまでその場に立ってぼんやりやかんを見つめ続け、沸騰してから残りのことを全部やる、というやり方です。

最初のやり方でお茶を入れる方が速い理由を理解できれば、`CheckWebsites`をより速くする方法を理解できます。 Webサイトの応答を待ってから次のWebサイトに要求を送信する代わりに、待機中に次の要求を行うようにコンピューターに指示します。

通常、Goで関数 `doSomething()`を呼び出すと、関数が返されるのを待ちます（返される値がない場合でも、関数が終了するのを待ちます）。この操作は _ブロッキング（blocking）_ であると言います。完了するまで待機させられるということです。

Goでブロックしない操作は、ゴルーチン _goroutine_ と呼ばれる別のプロセスで実行されます。
プロセスは、Goコードのページを上から下に読み取り、呼び出されたときに各関数の「内部」に移動して、その機能を読み取るものと考えてください。
別のプロセスが開始されると、別の読み手が関数の中を読み始め、元の読み手はそのままページを読み進めていく、というイメージです。


Goに新しいゴルーチンを開始するよう指示するには、関数呼び出しの前にキーワード`go`を置くことで、関数呼び出しを`go`ステートメントに変換します：「`go doSomething()`」。

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()
	}

	return results
}
```

ゴルーチンを開始する唯一の方法は、`go`を関数呼び出しの前に置くことなので、ゴルーチンを開始したい場合は、しばしば _無名関数_ を使用します。無名関数リテラルは、通常の関数宣言とまったく同じように見えますが、（当然ながら）名前がありません。
上記は、`for`ループの本体で確認できます。

無名関数には、それらを便利にするいくつかの機能があり、そのうち2つは上記で使用しています。まず、宣言と同時に実行できます-これは、無名関数の最後にある`()`が行っていることです。次に、定義されている字句スコープへのアクセスを維持します。無名関数を宣言した時点で使用可能なすべての変数は、関数の本体でも使用できます。

上記の無名関数の本体は、以前のループ本体とまったく同じです。唯一の違いは、ループの各反復が現在のプロセス（`WebsiteChecker`関数）と並行して新しいゴルーチンを開始することです。各ゴルーチンは、その結果を結果マップに追加します。

しかし、`go test`を実行すると

```sh
--- FAIL: TestCheckWebsites (0.00s)
        CheckWebsites_test.go:31: Wanted map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], got map[]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s

```

### 余談ですが、並行性の世界へ...

この結果が得られない可能性があります。

後で説明するパニックメッセージが表示される場合があります。それが得られても心配しないでください。

上記の結果が _実際に_ 得られるまでテストを実行し続けてください。あるいは、得られたことにしておいてください。
お好きなほうで。並行処理の世界へようこそ。

正しく処理されない場合、何が起こるかを予測することは困難です。心配しないでください。
だからこそ、並行処理を予測可能な形で扱えているかを確かめるために、テストを書いているのです。

### ...戻ってきました

元からある`CheckWebsites`のテストで失敗しています。空のマップが返されるようになってしまいました。何が悪かったのでしょうか？

`for`ループが開始したゴルーチンには、結果を`results`マップに追加するのに十分な時間がありませんでした。 `CheckWebsites`関数は速すぎて、まだ空のマップを返してしまいます。

これを修正するには、すべてのゴルーチンが作業を行っている間待機してから、戻ることができます。

2秒でできるはずですよね？

```go
package concurrency

import "time"

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		go func() {
			results[url] = wc(url)
		}()
	}

	time.Sleep(2 * time.Second)

	return results
}
```

運が良ければ、次のようになります。

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        2.012s
```

しかし、運が悪い場合は次のようになります（ベンチマークで実行すると試行回数が増えるため、こちらの可能性が高くなります）。

```sh
fatal error: concurrent map writes

goroutine 8 [running]:
runtime.throw(0x12c5895, 0x15)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/panic.go:605 +0x95 fp=0xc420037700 sp=0xc4200376e0 pc=0x102d395
runtime.mapassign_faststr(0x1271d80, 0xc42007acf0, 0x12c6634, 0x17, 0x0)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:783 +0x4f5 fp=0xc420037780 sp=0xc420037700 pc=0x100eb65
github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1(0xc42007acf0, 0x12d3938, 0x12c6634, 0x17)
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x71 fp=0xc4200377c0 sp=0xc420037780 pc=0x12308f1
runtime.goexit()
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/asm_amd64.s:2337 +0x1 fp=0xc4200377c8 sp=0xc4200377c0 pc=0x105cf01
created by github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xa1

        ... many more scary lines of text ...
```

これは長くて怖いですが、私たちがする必要があるのは、一息ついてスタックトレースを読むだけです。

`fatal error: concurrent map writes`（致命的なエラー：マップへの同時書き込み）。

テストを実行すると、2つのゴルーチンがまったく同時に結果マップに書き込む場合があります。Goのマップは、一度に複数のものが書き込もうとすると気に入らないため、「致命的なエラー（`fatal error`）」が発生します。

これは _データ競合（data race）_ です。2つ以上のゴルーチンが同じメモリ位置に並行してアクセスし、そのうち少なくとも1つが書き込みである場合に発生するバグです。各ゴルーチンがいつ実行されるかを正確に制御することはできないため、複数のゴルーチンが`results`マップにまったく同時に書き込もうとする危険があります。Goのマップは並行書き込みに対して安全ではないため、ランタイムはメモリ破壊を防ぐために致命的なエラーを投げます。

Goは、組み込みの[_レースディテクター（race detector）_](https://blog.golang.org/race-detector)で競合状態を特定するのに役立ちます。
この機能を有効にするには、`race`フラグを指定してテストを実行します。`go test -race`。

次のような出力が表示されます。

```sh
==================
WARNING: DATA RACE
Write at 0x00c420084d20 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Previous write at 0x00c420084d20 by goroutine 7:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Goroutine 8 (running) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c

Goroutine 7 (finished) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c
==================
```

詳細は読みにくいですが、「警告：データレース（`WARNING: DATA RACE`）」は非常に明確です。エラーの本文を読むと、マップで書き込みを実行する2つの異なるゴルーチンがわかります。

'goroutine 8' で '0x00c420084d20'に書き込む（`Write at 0x00c420084d20 by goroutine 8:`）

と同じメモリブロックに書き込んでいます

'goroutine 7'による'0x00c420084d20'での以前の書き込み（`Previous write at 0x00c420084d20 by goroutine 7`）

さらに、書き込みが行われているコード行を確認できます。

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12`

そして、ゴルーチン'7'と'8'が開始されるコード行も確認できます。

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11`

知る必要のあることはすべてターミナルに表示されています。
あとは、根気よくそれを読むだけです。

### チャネル

_チャネル（channels）_ を使用してゴルーチンを調整することで、このデータ競合を解決できます。
チャネルは、値の受信と送信の両方が可能なGoデータ構造です。これらの操作とその詳細により、異なるプロセス間の通信が可能になります。

今回は、親プロセスと、URLを渡して`WebsiteChecker`関数を実行する作業のために親プロセスが作成する各ゴルーチンとの間の通信について考えます。

```go
package concurrency

type WebsiteChecker func(string) bool
type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {
		go func() {
			resultChannel <- result{url, wc(url)}
		}()
	}

	for i := 0; i < len(urls); i++ {
		r := <-resultChannel
		results[r.string] = r.bool
	}

	return results
}
```

> **ゴルーチン内の`url`についての注意** ループの各反復では、`url`を明示的に渡すことなく、`url`を参照する新しいゴルーチンを開始しています。Go 1.22以降、これは安全です。言語仕様が変更され、`url`は反復ごとに新しい変数になったため、各ゴルーチンはそれぞれ独自のコピーをキャプチャします。
>
> `go.mod`で宣言されている`go`のバージョンが`1.22`より *古い* プロジェクトでこれを実行すると、古い挙動になります。つまり、`url`はすべての反復で共有・再利用される単一の変数となるため、ゴルーチンが実行される頃には、すべてのゴルーチンが同じ（おそらく最後の）`url`の値を見ることになるでしょう。Goの *ツールチェーン* が新しくても`go.mod`の`go`ディレクティブが古い場合があるため、これはよく人をつまずかせます。ツールチェーンは、宣言されたバージョンが示すループ変数のセマンティクスに従うのです。古い`go.mod`での修正方法は、`url`をゴルーチンに明示的に渡すことです: `go func(url string) { ... }(url)`

`results`マップに加えて、同じ方法で`make`する`resultChannel`があります。`chan result`はチャネルの型で、`result`のチャネルです。新しい型`result`は、`WebsiteChecker`の戻り値をチェック対象のURLに関連付けるために作成されました。これは、`string`および`bool`の構造体です。

どちらの値にも名前を付ける必要はないので、それらはそれぞれ構造体内で匿名です。
これは、値に名前を付けるのが難しい場合に役立ちます。

URLを反復処理するとき、`map`に直接書き込む代わりに、 _send statement_ を使用して、`wc`への各呼び出しの `result`構造体を`resultChannel`に送信します。これは、`<-`演算子を使用して、左側にチャネル、右側に値を取ります。

```go
// Send statement
resultChannel <- result{url, wc(url)}
```

次の`for`ループは、各URLに対して1回反復します。
内部では、チャネルから受信した値を変数に割り当てる _receive式_ を使用しています。
これも`<-`演算子を使用していますが、2つのオペランドが逆になっています。チャネルが右側にあり、代入する変数が左側にあります。

```go
// Receive expression
r := <-resultChannel
```

次に、受け取った`result`を使用してマップを更新します。

結果をチャネルに送信することにより、結果マップへの各書き込みのタイミングを制御して、一度に1つずつ発生するようにします。

`wc`の各呼び出しと結果チャネルへの送信はそれぞれ独自のプロセス内で並行して発生しますが、receive式で結果チャネルから値を取り出すため、結果はそれぞれ一度に1つずつ処理されます。

高速化したいコードの部分には並行性を使い、同時に実行できない部分は引き続き線形的に発生するようにしました。
また、チャネルを使用することにより、関連する複数のプロセス間で通信しました。

ベンチマークを実行すると

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v2
BenchmarkCheckWebsites-8             100          23406615 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        2.377s
```

23406615ナノ秒（0.023秒）で、元の関数の約100倍の速さです。

大成功。

## まとめ

この演習は、いつもよりTDDの要素が少し控えめでした。
ある意味では、`CheckWebsites`関数の1つの長いリファクタリングに参加しています。

入力と出力が変更されることはなく、ただ速くなっただけです。しかし、実施したテストと作成したベンチマークにより、ソフトウェアがまだ機能しているという確信を維持しながら、実際に高速になったことを示す方法で`CheckWebsites`をリファクタリングすることができました。

それをより速くすることで、私たちは次のことを学びました。

* _ゴルーチン（goroutines）_ 、Goの並行処理の基本単位です。これにより、複数のウェブサイトチェックリクエストを管理できました。
* _無名関数（anonymous functions）_ 、ウェブサイトをチェックする各並行プロセスを開始するために使用しました。
* _チャネル（channels）_ 、さまざまなプロセス間の通信を整理および制御するのに役立ち、 _競合状態（race condition）_ のバグを回避できました。
* _レースディテクター（the race detector）_ 、並行コードに関する問題のデバッグに役立ちました。

### 速くする

ソフトウェアをアジャイルに構築する方法を定式化したものの1つで、しばしば誤ってケント・ベックの言葉とされているのが次のものです。

> [機能させる、正しくする、速くする](http://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast)

「機能させる」とはテストに合格すること、「正しくする」とはコードをリファクタリングすること、そして「速くする」とはコードを最適化して、たとえばすばやく実行することです。
「速くする」ことができるのは、それを機能させて正しくした後だけです。

与えられたコードは既に機能していることが実証されており、リファクタリングする必要がないことは幸運でした。他の2つのステップを実行する前に、「速くする」ことを試みるべきではありません。なぜなら、

> [早期最適化はすべての悪の根源](http://wiki.c2.com/?PrematureOptimization) -- Donald Knuth
