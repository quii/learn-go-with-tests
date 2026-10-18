---
description: OS Exec
---

# OS実行

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/os-exec)

[keith6014](https://www.reddit.com/user/keith6014)が[reddit](https://www.reddit.com/r/golang/comments/aaz8ji/testdata_and_function_setup_help/)で質問する

> XMLデータを生成した`os/exec.Command()`を使用してコマンドを実行しています。コマンドは、`GetData()`という関数で実行されます。
>
> `GetData()`をテストするために、作成した`testdata`をいくつか持っています。
>
> 私の`_test.go`には、`GetData()`を呼び出す`TestGetData`がありますが、`os.exec`を使用しますが、代わりに`testdata`を使用します。
>
> これを達成するための良い方法は何ですか？ `GetData`を呼び出すときは、「テスト」フラグモードを使用して、`GetData(mode string)`というファイルを読み取る必要がありますか？

いくつかのこと

* 何かをテストするのが難しい場合、懸念の分離が正しくないことが原因であることがよくあります
* コードに「テストモード」を追加しないでください。代わりに、[依存関係の注入（DI）](../go-fundamentals/dependency-injection.md)を使用して、依存関係をモデル化し、懸念事項を分離できるようにします。

私は勝手にコードがどのように見えるかを推測しました。

```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData() string {
	cmd := exec.Command("cat", "msg.xml")

	out, _ := cmd.StdoutPipe()
	var payload Payload
	decoder := xml.NewDecoder(out)

	// these 3 can return errors but I'm ignoring for brevity
	cmd.Start()
	decoder.Decode(&payload)
	cmd.Wait()

	return strings.ToUpper(payload.Message)
}
```

* プロセスへの外部コマンドを実行できる`exec.Command`を使用します
* 出力を`cmd.StdoutPipe`にキャプチャして、`io.ReadCloser`を返します（これは重要になります）
* コードの残りの部分は、[優れたドキュメント](https://golang.org/pkg/os/exec/#example_Cmd_StdoutPipe)から多かれ少なかれコピーして貼り付けたものです。
  * stdoutからの出力をキャプチャして`io.ReadCloser`に取り込み、次にコマンドを`Start`してから、`Wait`を呼び出してすべてのデータが読み取られるのを待ちます。これらの2つの呼び出しの間で、データを`Payload`構造体にデコードします。

これが `msg.xml`内に含まれているものです。

```xml
<payload>
    <message>Happy New Year!</message>
</payload>
```

実際の動作を示す簡単なテストを作成しました。

```go
func TestGetData(t *testing.T) {
	got := GetData()
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

## テスト可能なコード

テスト可能なコードは分離され、単一の目的です。
私には、このコードには2つの主な懸念があるように感じます

1. 未加工のXMLデータを取得する。
2. XMLデータをデコードし、ビジネスロジックを適用します（この場合は`<message>`の`strings.ToUpper`です）

最初の部分は、サンプルを標準`lib`からコピーすることです。

2番目の部分はビジネスロジックがある場所です。コードを見ると、ロジックの「シーム`"seam"`」がどこから始まるかがわかります。ここで、`io.ReadCloser`を取得します。この既存の抽象化を使用して、問題を分離し、コードをテスト可能にすることができます。

**GetDataの問題は、ビジネスロジックがXMLを取得する手段と結合していることです。デザインを改善するには、それらを切り離す必要があります**

私たちの`TestGetData`は2つの懸念事項の間の統合テストとして機能することができるため、それが機能し続けることを確認するためにそれを保持します。

新しく分離されたコードは次のようになります

```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := io.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

func TestGetDataIntegration(t *testing.T) {
	got := GetData(getXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

`GetData`が入力を`io.Reader`から取得するようになったので、これをテスト可能にして、データの取得方法を考慮しなくなりました。人々は`io.Reader`（これは非常に一般的です）を返すもので関数を再利用できます。

たとえば、コマンドラインではなくURLからXMLのフェッチを開始できます。

```go
func TestGetData(t *testing.T) {
	input := strings.NewReader(`
<payload>
    <message>Cats are the best animal</message>
</payload>`)

	got := GetData(input)
	want := "CATS ARE THE BEST ANIMAL"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

```

これは`GetData`の単体テストの例です。

懸念を分離し、Goのテストで既存の抽象化を使用することで、重要なビジネスロジックが簡単になります。
