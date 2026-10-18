---
description: Error types
---

# エラーの種類

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/error-types)

**エラー用に独自のタイプを作成することは、コードを整頓し、コードを使いやすくテストするための洗練された方法になる場合があります。**

`Gopher Slack`の`Pedro`が尋ねる

> `fmt.Errorf("%s is foo、got%s"、bar、baz)`のようなエラーを作成している場合、文字列値を比較せずに同等性をテストする方法はありますか？

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

`Pedro`が言うように、ステータスエラーのテストをそのように書くことができました。

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

このサイトは _テストに耳を傾ける_ ことを強調しようとしていますが、このテストは良いとは感じません。

* テストのために本番コードと同じ文字列を作成しています
* 読み書きが面倒
* 正確なエラーメッセージ文字列は、実際に関係しているものですか？

これは何を教えてくれますか？
テストの人間工学は、コードを使用しようとする別のコードに反映されます。

コードのユーザーは、返される特定の種類のエラーにどのように反応しますか？
彼らができる最善のことは、非常にエラーが発生しやすく、恐ろしいエラー文字列を調べることです。

## 私たちがすべきこと

TDDを使用すると、以下の考え方に入ることができます。

> このコードをどのように使用したいですか？

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

エラーの正確な文字列をチェックする代わりに、エラーに対して[タイプアサーション（type assertion）](https://tour.golang.org/methods/15)を実行して、エラーが `BadStatusError`であるかどうかを確認しています。

これは、エラーをクリアするという種類の要望を反映しています。アサーションがパスすると仮定して、エラーのプロパティが正しいことを確認できます。

テストを実行すると、正しい種類のエラーが返されなかったことがわかります

```
--- FAIL: TestDumbGetter (0.00s)
    --- FAIL: TestDumbGetter/when_you_dont_get_a_200_you_get_a_status_error (0.00s)
    	error-types_test.go:56: was not a BadStatusError, got *errors.errorString
```

タイプを使用するようにエラー処理コードを更新して、`DumbGetter`を修正しましょう

```go
if res.StatusCode != http.StatusOK {
	return "", BadStatusError{URL: url, Status: res.StatusCode}
}
```

この変更は、いくつかの _現実的な_ プラスの効果をもたらしました。

* `DumbGetter`関数がシンプルになりました。エラー文字列の複雑さに関係することはなくなり、`BadStatusError`を作成するだけです。
* 私たちのテストは、コードのユーザーがロギングだけではなく、より高度なエラー処理を実行することを決定した場合に、およびドキュメントを反映しています。タイプアサーションを実行するだけで、エラーのプロパティに簡単にアクセスできます。
* それでも「単なるエラー（`error`）」なので、彼らが選択した場合、コールスタックに渡すか、他の「エラー`error`」と同様にログに記録できます。

## まとめ

複数のエラー条件をテストしていることに気づいたら、エラーメッセージを比較するという罠にはまらないようにしてください。

これは、不完全で読み書きの難しいテストになります。また、発生したエラーの種類に応じて異なることを始める必要がある場合に、 コードのユーザが抱える困難さを反映します。

あなたがどのようにコードを使いたいかをテストに反映させるようにしてください。この点で、エラーの種類をカプセル化するためにエラータイプを作成することを検討してください。これにより、異なる種類のエラーの処理がコードのユーザにとって容易になり、また、エラー処理のコードをよりシンプルで読みやすく書くことができるようになります。

## 補遺

Go1.13では、標準ライブラリのエラーを扱う新しい方法があります。[Goブログ](https://blog.golang.org/go1.13-errors)

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
