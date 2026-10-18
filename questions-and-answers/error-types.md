---
description: エラーの種類
---

# エラーの種類

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/error-types)

**エラー用に独自の型を作成することは、コードを整理し、使いやすくテストしやすくするための洗練された方法になり得ます。**

Gopher SlackでPedroがこんな質問をしていました。

> `fmt.Errorf("%s must be foo, got %s", bar, baz)` のようなエラーを作成している場合、文字列値を比較せずに同等性をテストする方法はありますか？

このアイデアを探索するのに役立つ関数を作りましょう。

```go
// DumbGetter will get the string body of url if it gets a 200
func DumbGetter(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("problem fetching from %s, %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("did not get 200 from %s, got %d", url, res.StatusCode)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body) // ignoring err for brevity

	return string(body), nil
}
```

さまざまな理由で失敗する可能性のある関数を作成することは珍しいことではなく、各シナリオを正しく処理できるようにしたいと考えています。

Pedroが言うように、ステータスエラーのテストを次のように書くことも _できます_ 。

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	want := fmt.Sprintf("did not get 200 from %s, got %d", svr.URL, http.StatusTeapot)
	got := err.Error()

	if got != want {
		t.Errorf(`got "%v", want "%v"`, got, want)
	}
})
```

このテストでは、常に`StatusTeapot`を返すサーバーを作成し、そのURLを`DumbGetter`の引数として使用して、`200`以外の応答を正しく処理できることを確認します。

## このテスト方法の問題

この本では _テストに耳を傾ける_ ことを強調していますが、このテストは良いものだと _感じられません_ 。

* テストのために本番コードと同じ文字列を作成しています
* 読み書きが面倒
* 正確なエラーメッセージ文字列こそが、私たちが _本当に気にしている_ ことでしょうか？

これは何を物語っているのでしょうか？
このテストの使い勝手は、私たちのコードを使おうとする他のコードの使い勝手にもそのまま表れる、ということです。

コードのユーザーは、返される特定の種類のエラーにどのように反応しますか？
せいぜいエラー文字列を調べることぐらいしかできませんが、これは非常に間違いが起きやすく、書くのもひどく面倒です。

## 私たちがすべきこと

TDDには、次のような考え方ができるという利点があります。

> _自分なら_ このコードをどのように使いたいだろうか？

`DumbGetter`にできることは、ユーザーが型システムを使用して発生したエラーの種類を理解する方法を提供することです。

もしも`DumbGetter`が次のようなものを返してくれたらどうでしょうか？

```go
type BadStatusError struct {
	URL    string
	Status int
}
```

魔法の文字列ではなく、実際に使用する _データ_ があります。

このニーズを反映するように既存のテストを変更しましょう。

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	got, isStatusErr := err.(BadStatusError)

	if !isStatusErr {
		t.Fatalf("was not a BadStatusError, got %T", err)
	}

	want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
})
```

`BadStatusError`にエラーインターフェースを実装させる必要があります。

```go
func (b BadStatusError) Error() string {
	return fmt.Sprintf("did not get 200 from %s, got %d", b.URL, b.Status)
}
```

### テストは何をしますか？

エラーの正確な文字列をチェックする代わりに、エラーに対して[型アサーション（type assertion）](https://tour.golang.org/methods/15)を実行して、エラーが `BadStatusError`であるかどうかを確認しています。

これにより、エラーの _種類_ を確かめたいという要望がより明確に表れます。アサーションがパスすると仮定して、エラーのプロパティが正しいことを確認できます。

テストを実行すると、正しい種類のエラーが返されなかったことがわかります

```
--- FAIL: TestDumbGetter (0.00s)
    --- FAIL: TestDumbGetter/when_you_dont_get_a_200_you_get_a_status_error (0.00s)
    	error-types_test.go:56: was not a BadStatusError, got *errors.errorString
```

独自の型を使用するようにエラー処理コードを更新して、`DumbGetter`を修正しましょう

```go
if res.StatusCode != http.StatusOK {
	return "", BadStatusError{URL: url, Status: res.StatusCode}
}
```

この変更は、いくつかの _本当に良い効果_ をもたらしました。

* `DumbGetter`関数がシンプルになりました。エラー文字列の複雑さに関係することはなくなり、`BadStatusError`を作成するだけです。
* テストは、コードのユーザーがログ出力だけでなく、より高度なエラー処理をしたいと考えた場合に何が _できるか_ を反映（そして文書化）するようになりました。型アサーションを行うだけで、エラーのプロパティに簡単にアクセスできます。
* それでも「単なる」`error`なので、ユーザーは望むなら、他の`error`と同様にコールスタックの上位に渡したり、ログに記録したりできます。

## まとめ

複数のエラー条件をテストしていることに気づいたら、エラーメッセージを比較するという罠にはまらないようにしてください。

そうすると、テストは不安定で読み書きの難しいものになります。そしてそれは、発生したエラーの種類に応じて処理を変える必要が出てきたときに、コードのユーザーが抱えることになる困難をそのまま表しています。

_あなた自身が_ コードをどのように使いたいかを、常にテストに反映させるようにしてください。その意味で、エラーの種類をカプセル化するためにエラー型を作成することを検討してください。これにより、異なる種類のエラーの処理がコードのユーザーにとって容易になり、また、エラー処理のコードをよりシンプルで読みやすく書くことができるようになります。

## 補遺

Go 1.13以降、標準ライブラリにはエラーを扱う新しい方法が用意されています。これについては[Goブログ](https://blog.golang.org/go1.13-errors)で解説されています。

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	var got BadStatusError
	isBadStatusError := errors.As(err, &got)
	want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

	if !isBadStatusError {
		t.Fatalf("was not a BadStatusError, got %T", err)
	}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
})
```

この場合、[`errors.As`](https://pkg.go.dev/errors#example-As)を使ってエラーをカスタム型に抽出しています。これは成功を示すために`bool`を返し、それを`got`に抽出してくれます。
