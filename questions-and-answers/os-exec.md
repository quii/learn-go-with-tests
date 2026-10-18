---
description: OS実行
---

# OS実行

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/os-exec)

[keith6014](https://www.reddit.com/user/keith6014)が[reddit](https://www.reddit.com/r/golang/comments/aaz8ji/testdata_and_function_setup_help/)で質問しています。

> XMLデータを生成するコマンドを`os/exec.Command()`を使って実行しています。コマンドは、`GetData()`という関数で実行されます。
>
> `GetData()`をテストするために、自分で作成したテストデータ（`testdata`）がいくつかあります。
>
> `_test.go`には`GetData()`を呼び出す`TestGetData`がありますが、それだと`os.exec`が使われてしまいます。代わりに、自分のテストデータを使うようにしたいのです。
>
> これを実現する良い方法は何でしょうか？ `GetData`を呼び出すときに「テスト」フラグのモードを用意して、ファイルを読み込むようにすべきでしょうか？ つまり`GetData(mode string)`のように。

いくつか言っておきたいことがあります。

* 何かをテストするのが難しい場合、関心の分離がうまくできていないことが原因であることがよくあります
* コードに「テストモード」を追加しないでください。代わりに、[依存性注入（DI）](../go-fundamentals/dependency-injection.md)を使用して、依存関係をモデル化し、関心を分離できるようにしましょう。

勝手ながら、コードがどのようなものかを推測してみました。

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
* `cmd.StdoutPipe`で出力をキャプチャします。これは`io.ReadCloser`を返します（これが重要になります）
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

テスト可能なコードは疎結合で、単一の目的を持っています。
私には、このコードには主に2つの関心事があるように感じられます。

1. 未加工のXMLデータを取得する。
2. XMLデータをデコードし、ビジネスロジックを適用します（この場合は`<message>`の`strings.ToUpper`です）

最初の部分は、標準ライブラリのサンプルをコピーしただけです。

2番目の部分はビジネスロジックがある場所です。コードを見ると、ロジックの「継ぎ目（seam）」がどこから始まるかがわかります。それは`io.ReadCloser`を取得する箇所です。この既存の抽象化を使用して、関心を分離し、コードをテスト可能にすることができます。

**GetDataの問題は、ビジネスロジックがXMLを取得する手段と結合していることです。デザインを改善するには、それらを切り離す必要があります**

`TestGetData`は2つの関心事の間の統合テストとして機能できるので、引き続き動作することを確認するために残しておきます。

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

`GetData`が入力を`io.Reader`から取得するようになったので、テスト可能になり、データがどのように取得されるかに関与しなくなりました。`io.Reader`を返すもの（これは非常によくあります）であれば、何とでもこの関数を再利用できます。

たとえば、コマンドラインではなくURLからXMLを取得するように変えることもできます。

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

関心を分離し、Goの既存の抽象化を使用することで、重要なビジネスロジックのテストがとても簡単になります。
