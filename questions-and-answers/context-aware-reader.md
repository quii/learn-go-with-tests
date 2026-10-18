---
description: Context-aware Reader
---

# コンテキスト認識リーダー

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/context-aware-reader)

この章では、_Mat・Ryer_ と _David・Hernandez_ が [The Pace Dev Blog](https://pace.dev/blog/2020/02/03/context-aware-ioreader-for-golang-by-mat-ryer)で書いた、コンテキストを意識した`io.Reader`をテストドライブする方法を示します。

## コンテキスト認識リーダー？

まず、`io.Reader`の簡単な入門書です。

この本の他の章を読んだことがあれば、ファイルを開いたり、JSONをエンコードしたり、その他の一般的なタスクを実行したりしたときに、`io.Reader`に出くわすことになります。 _something_ からのデータの読み取りに関する単純な抽象化です。

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

`io.Reader`を使用すると、標準ライブラリから多くの再利用を得ることができます。
これは、非常に一般的に使用される抽象化です（対応する` io.Writer`とともに）

### コンテキスト認識？

[前の章で](../go-fundamentals/context.md)キャンセルを提供するために「コンテキスト`context`」を使用する方法について説明しました。これは、計算コストがかかる可能性のあるタスクを実行していて、それらを停止できるようにしたい場合に特に便利です。

`io.Reader`を使用している場合、速度について保証はありません。1ナノ秒または数百時間かかる場合があります。自分のアプリケーションでこの種のタスクをキャンセルできると便利だと思うかもしれませんが、それはMatとDavidが書い​​たものです。

彼らはこの問題を解決するために2つの単純な抽象化（`context.Context`と`io.Reader`）を組み合わせました。

`io.Reader`をラップしてキャンセルできるように、いくつかの機能をTDDで試してみましょう。

これをテストすることは興味深い挑戦を引き起こします。通常、`io.Reader`を使用するときは、通常、他の関数にそれを提供しているので、詳細に気を使う必要はありません。
`json.NewDecoder`や`io.ReadAll`など。

デモしたいのは、次のようなものです。

> "ABCDEF"を持つ`io.Reader`が与えられたとき、途中でキャンセル信号を送っても、読み続けようとすると何も出ないので、"ABC"しか出ません。

インターフェースをもう一度見てみましょう。

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

`Reader`の`Read`メソッドは、その内容を、提供する `[]byte`に読み込みます。

したがって、すべてを読むのではなく、次のことができます。

* すべてのコンテンツに適合しない固定サイズのバイト配列を提供します
* キャンセル信号を送信します
* 再試行してもう一度読み取ると、0バイトの読み取りエラーが返されます

とりあえず、キャンセルのない「ハッピーパス」テストを書いてみましょう。これは、まだ本番用のコードを記述しなくても問題に慣れることができるようにするためです。

```go
func TestContextAwareReader(t *testing.T) {
	t.Run("lets just see how a normal reader works", func(t *testing.T) {
		rdr := strings.NewReader("123456")
		got := make([]byte, 3)
		_, err := rdr.Read(got)

		if err != nil {
			t.Fatal(err)
		}

		assertBufferHas(t, got, "123")

		_, err = rdr.Read(got)

		if err != nil {
			t.Fatal(err)
		}

		assertBufferHas(t, got, "456")
	})
}

func assertBufferHas(t testing.TB, buf []byte, want string) {
	t.Helper()
	got := string(buf)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

* いくつかのデータを含む文字列から`io.Reader`を作成します
* 読み込む内容がリーダーの内容よりも小さいバイト配列
* 呼び出しを読んで、内容を確認し、繰り返します

これから、2回目の読み取りの前に動作を変更するために何らかのキャンセル信号を送信することを想像できます。

それがどのように機能するかを見てきましたので、残りの機能をTDDします。

## 最初にテストを書く

`io.Reader`を`context.Context`で作成できるようにしたいと考えています。

TDDでは、希望するAPIを想像することから始めて、そのためのテストを作成するのが最善です。

そこから、コンパイラーと失敗したテスト出力で解決策を導きましょう。

```go
t.Run("behaves like a normal reader", func(t *testing.T) {
	rdr := NewCancellableReader(strings.NewReader("123456"))
	got := make([]byte, 3)
	_, err := rdr.Read(got)

	if err != nil {
		t.Fatal(err)
	}

	assertBufferHas(t, got, "123")

	_, err = rdr.Read(got)

	if err != nil {
		t.Fatal(err)
	}

	assertBufferHas(t, got, "456")
})
```

## テストを実行してみます

```
./cancel_readers_test.go:12:10: undefined: NewCancellableReader
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

この関数を定義する必要があり、`io.Reader`を返す必要があります

```go
func NewCancellableReader(rdr io.Reader) io.Reader {
	return nil
}
```

試しに実行してみると

```
=== RUN   TestCancelReaders
=== RUN   TestCancelReaders/behaves_like_a_normal_reader
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x10f8fb5]
```

予想通り

## 成功させるのに十分なコードを書く

今のところ、渡した`io.Reader`を返すだけです

```go
func NewCancellableReader(rdr io.Reader) io.Reader {
	return rdr
}
```

これでテストに成功するはずです。

わかっています、わかっています、これは馬鹿げていて衒学的に見えますが、派手な作業に取り掛かる前に、`io.Reader`の「通常の」振る舞いを壊していないかどうかを、ある程度検証することが重要です。

## 最初にテストを書く

次に、キャンセルしてみる必要があります。

```go
t.Run("stops reading when cancelled", func(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rdr := NewCancellableReader(ctx, strings.NewReader("123456"))
	got := make([]byte, 3)
	_, err := rdr.Read(got)

	if err != nil {
		t.Fatal(err)
	}

	assertBufferHas(t, got, "123")

	cancel()

	n, err := rdr.Read(got)

	if err == nil {
		t.Error("expected an error after cancellation but didn't get one")
	}

	if n > 0 {
		t.Errorf("expected 0 bytes to be read after cancellation but %d were read", n)
	}
})
```

最初のテストは多かれ少なかれコピーできますが、今は次のようになっています。

* 最初の読み取り後に「キャンセル`cancel`」できるように、キャンセル付きの`context.Context`を作成します
* コードを機能させるには、関数に`ctx`を渡す必要があります
* そして、`cancel`後に何も読まれなかったことを主張します。

## テストを実行してみます

```
./cancel_readers_test.go:33:30: too many arguments in call to NewCancellableReader
	have (context.Context, *strings.Reader)
	want (io.Reader)
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

コンパイラは何をすべきかを指示しています。コンテキストを受け入れるように署名を更新します。

```go
func NewCancellableReader(ctx context.Context, rdr io.Reader) io.Reader {
	return rdr
}
```

（最初のテストが `context.Background` を通過するように更新する必要があります。）

これで、非常に明確な失敗したテスト出力を見ることができるはずです。

```
=== RUN   TestCancelReaders
=== RUN   TestCancelReaders/stops_reading_when_cancelled
--- FAIL: TestCancelReaders (0.00s)
    --- FAIL: TestCancelReaders/stops_reading_when_cancelled (0.00s)
        cancel_readers_test.go:48: expected an error but didn't get one
        cancel_readers_test.go:52: expected 0 bytes to be read after cancellation but 3 were read
```

## 成功させるのに十分なコードを書く

この時点では、MatとDavidによる元の投稿からコピーして貼り付けていますが、ゆっくりと繰り返し実行します。

読み込んだ`io.Reader`と`context.Context`をカプセル化するタイプが必要であることはわかっているので、それを作成して、元の`io.Reader`の代わりに関数からそれを返してみましょう

```go
func NewCancellableReader(ctx context.Context, rdr io.Reader) io.Reader {
	return &readerCtx{
		ctx:      ctx,
		delegate: rdr,
	}
}

type readerCtx struct {
	ctx      context.Context
	delegate io.Reader
}
```

このサイトで何度も強調してきたように、ゆっくりと進み、コンパイラに助けてもらいましょう。

```
./cancel_readers_test.go:60:3: cannot use &readerCtx literal (type *readerCtx) as type io.Reader in return argument:
	*readerCtx does not implement io.Reader (missing Read method)
```

抽象化はいい感じだけど、必要なインターフェースが実装されてないから、メソッドを追加しよう。

```go
func (r *readerCtx) Read(p []byte) (n int, err error) {
	panic("implement me")
}
```

テストを実行すると、それらはコンパイルするはずですが、パニックになります。これはまだ進行中です。

最初のテストを通過させるために、基礎となる`io.Reader`への呼び出しを委任してみましょう。

```go
func (r readerCtx) Read(p []byte) (n int, err error) {
	return r.delegate.Read(p)
}
```

この時点で、私たちは再び私たちの幸せなパスのテストが合格し、それは私たちのものがきれいに抽象化されているように感じています。

2回目のテストをパスするためには、`context.Context`がキャンセルされたかどうかを確認する必要があります。

```go
func (r readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.delegate.Read(p)
}
```

これですべてのテストが通過するはずです。エラーを`context.Context`から返していることに気づくでしょう。これにより、コードの呼び出し元がキャンセルが発生した様々な理由を調べることができます。

## まとめ

* 小さなインターフェイスが良く、構成が簡単
* あるもの ( `io.Reader`) を別のもので拡張しようとするとき、通常は[委任パターン（delegation pattern）](https://en.wikipedia.org/wiki/Delegation_pattern)に到達したいと思います。

> ソフトウェア工学では、デリゲーションパターンはオブジェクト指向の設計パターンであり、オブジェクトを構成して継承と同じコードの再利用を実現することができます。

* この種の作業を始める簡単な方法は、他の部分の動作を変更するために他の部分を合成し始める前に、デリゲートをラップして、デリゲートが通常どのように動作するかを保証するテストを書くことです。これは、目標に向かってコードを書く際に、物事を正しく動作させておくのに役立ちます。
