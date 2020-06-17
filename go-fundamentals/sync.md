---
description: Sync
---

# 同期

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/master/sync)

安全に併用できる並行処理を作りたい。

安全でない並行処理から始めて、その動作がシングルスレッド環境で機能することを確認します。

次に、複数のゴルーチンがテストを介してそれを使用して修正することで、安全でないことを実行します。

## 最初にテストを書く

APIで、カウンターをインクリメントしてその値を取得するメソッドを提供する必要があります。

```go
func TestCounter(t *testing.T) {
    t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
        counter := Counter{}
        counter.Inc()
        counter.Inc()
        counter.Inc()

        if counter.Value() != 3 {
            t.Errorf("got %d, want %d", counter.Value(), 3)
        }
    })
}
```

## テストを実行してみます

```text
./sync_test.go:9:14: undefined: Counter
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`Counter`を定義しましょう。

```go
type Counter struct {

}
```

再試行して、次のように失敗します。

```text
./sync_test.go:14:10: counter.Inc undefined (type Counter has no field or method Inc)
./sync_test.go:18:13: counter.Value undefined (type Counter has no field or method Value)
```

最終的にテストを実行するために、これらのメソッドを定義できます。

```go
func (c *Counter) Inc() {

}

func (c *Counter) Value() int {
    return 0
}
```

実行して失敗するはずです。

```text
=== RUN   TestCounter
=== RUN   TestCounter/incrementing_the_counter_3_times_leaves_it_at_3
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/incrementing_the_counter_3_times_leaves_it_at_3 (0.00s)
        sync_test.go:27: got 0, want 3
```

## 成功させるのに十分なコードを書く

これは、私たちのようなGoの専門家にとっては簡単なことです。データ型のカウンターの状態を保持し、`Inc`を呼び出すたびにインクリメントする必要があります

```go
type Counter struct {
    value int
}

func (c *Counter) Inc() {
    c.value++
}

func (c *Counter) Value() int {
    return c.value
}
```

## リファクタリング

リファクタリングすることはそれほど多くありませんが、`Counter`を中心にさらに多くのテストを作成するので、テストが少し明確に読み取れるように、小さなアサーション関数`assertCount`を作成します。

```go
t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
    counter := Counter{}
    counter.Inc()
    counter.Inc()
    counter.Inc()

    assertCounter(t, counter, 3)
})

func assertCounter(t *testing.T, got Counter, want int)  {
    t.Helper()
    if got.Value() != want {
        t.Errorf("got %d, want %d", got.Value(), want)
    }
}
```

## 次のステップ

それは十分に簡単でしたが、現在は、並行環境で使用しても安全である必要があるという要件があります。これを実行するには、失敗するテストを作成する必要があります。

## 最初にテストを書く

```go
t.Run("it runs safely concurrently", func(t *testing.T) {
    wantedCount := 1000
    counter := Counter{}

    var wg sync.WaitGroup
    wg.Add(wantedCount)

    for i := 0; i < wantedCount; i++ {
        go func(w *sync.WaitGroup) {
            counter.Inc()
            w.Done()
        }(&wg)
    }
    wg.Wait()

    assertCounter(t, counter, wantedCount)
})
```

これは、`wantedCount`をループし、**goroutine**を起動して、`counter.Inc()`を呼び出します。

並行プロセスを同期する便利な方法である[`sync.WaitGroup`](https://golang.org/pkg/sync/#WaitGroup)を使用しています。


> `WaitGroup`は、ゴルーチンのコレクションが完了するのを待ちます。メインのゴルーチンは`Add`を呼び出して、待機するゴルーチンの数を設定します。次に、各ゴルーチンが実行され、完了したら`Done`を呼び出します。同時に、すべてのゴルーチンが完了するまで、`Wait`を使用してブロックすることができます。

アサーションを作成する前に`wg.Wait()`が完了するのを待つことで、すべてのゴルーチンが`Counter`を`Inc`しようとしたことを確認できます。

## テストを実行してみます

```text
=== RUN   TestCounter/it_runs_safely_in_a_concurrent_envionment
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/it_runs_safely_in_a_concurrent_envionment (0.00s)
        sync_test.go:26: got 939, want 1000
FAIL
```

テストは別の数値で _おそらく_ 失敗しますが、それでも複数のゴルーチンが同時にカウンターの値を変更しようとしている場合は機能しないことを示しています。

## 成功させるのに十分なコードを書く

簡単な解決策は、`Counter`、[`Mutex`](https://golang.org/pkg/sync/#Mutex)にロックを追加することです。

> `Mutex`は相互排他ロックです。ミューテックスのゼロ値は、ロックされていないミューテックスです。

```go
type Counter struct {
    mu sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

これが意味することは、`Inc`を呼び出す**goroutine**が最初にある場合、`Counter`のロックを取得することです。他のすべてのゴルーチンは、アクセスを取得する前に、それが`Unlock`されるのを待つ必要があります。

ここでテストを再実行すると、変更を行う前に各ゴルーチンが順番を待たなければならないため、テストに合格するはずです。

## `sync.Mutex`が構造体に埋め込まれている他の例を見てきました。

あなたはこのような例を見るかもしれません

```go
type Counter struct {
    sync.Mutex
    value int
}
```

それはコードをもう少しエレガントにすることができると主張することができます。

```go
func (c *Counter) Inc() {
    c.Lock()
    defer c.Unlock()
    c.value++
}
```

これは見栄えが良いですが、プログラミングは非常に主観的な分野ですが、これは**悪いことであり間違っています**。

場合によっては、型の埋め込みがその型のメソッドが _public_ インターフェースの一部になることを忘れることがあります。そしてあなたはしばしばそれを望まないでしょう。
公開APIには細心の注意を払う必要があることを忘れないでください。何かを公開する瞬間は、他のコードがそれに結合できる瞬間です。私たちは常に不必要な結合を避けたいと思っています。

`Lock`と`Unlock`を公開することはせいぜい混乱を招きますが、最悪の場合、同じタイプの呼び出し元がこれらのメソッドの呼び出しを開始すると、ソフトウェアに非常に有害な可能性があります。

![このAPIのユーザーがロックの状態を誤って変更する方法を示します](https://i.imgur.com/SWYNpwm.png)

_これは本当に悪い考えのようです_

## `mutexes`のコピー

テストはパスしましたが、コードはまだ少し危険です

コードで`go vet`を実行すると、次のようなエラーが表示されます

```text
sync/v2/sync_test.go:16: call of assertCounter copies lock value: v1.Counter contains sync.Mutex
sync/v2/sync_test.go:39: assertCounter passes lock by value: v1.Counter contains sync.Mutex
```

[`sync.Mutex`](https://golang.org/pkg/sync/#Mutex)のドキュメントを見ると理由がわかります

> ミューテックスは、最初の使用後にコピーしてはなりません。

`Counter`（by value）を`assertCounter`に渡すと、ミューテックスのコピーが作成されます。

これを解決するには、代わりに`Counter`へのポインターを渡す必要があるため、`assertCounter`のシグネチャを変更します

```go
func assertCounter(t *testing.T, got *Counter, want int)
```

`*Counter`ではなく`Counter`を渡そうとしているため、テストはコンパイルされなくなりました。これを解決するには、自分で型を初期化しない方がよいことをAPIのリーダーに示すコンストラクタを作成することをお勧めします。

```go
func NewCounter() *Counter {
    return &Counter{}
}
```

`Counter`を初期化するときに、この関数をテストで使用します。

## まとめ

[同期パッケージ(sync package)](https://golang.org/pkg/sync/)のいくつかをカバーしました

* `Mutex`を使用すると、データにロックを追加できます
* `Waitgroup`は、ゴルーチンがジョブを完了するのを待つ手段です

### チャネルとゴルーチンにロックを使用するのはいつですか？

[私たちは以前に最初の並行性の章でゴルーチンをカバーしました](concurrency.md)これで安全な並行コードを書くことができるので、なぜロックを使うのでしょうか？
[go wikiには、このトピック専用のページがあります。ミューテックスまたはチャネル](https://github.com/golang/go/wiki/MutexOrChannel)

> Goの初心者によくある間違いは、それが可能であったり、楽しいからといって、チャネルやゴルーチンを使いすぎてしまうことです。`sync.Mutex`が問題に最も適している場合は、恐れずに同期を使用してください。 Goは、問題を最もよく解決するツールを使用できるようにし、1つのスタイルのコードに強制するのではなく、実用的です。

言い換え

* **データの所有権を渡すときにチャネルを使用する**
* **状態の管理にミューテックスを使用する**

### go vet

ビルドスクリプトで`go vet`を使用することを忘れないでください。貧弱なユーザーに影響が及ぶ前に、コード内のいくつかの微妙なバグを警告することができます。

### 便利なので埋め込みを使わないでください

* 埋め込みがパブリックAPIに与える影響について考えてください。
* これらのメソッドを _本当に_ 公開したいですか？
* ミューテックスに関しては、これは非常に予測不能で奇妙な方法で潜在的に悲惨なものになる可能性があります。あるべきでないミューテックスをロック解除するいくつかの悪意のあるコードを想像してください。これは非常に奇妙なバグを引き起こし、追跡が困難になります。
