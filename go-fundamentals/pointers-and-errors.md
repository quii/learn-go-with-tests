---
description: ポインタとエラー
---

# ポインタとエラー

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/pointers)

前のセクションでコンセプトに関連する、いくつかの値を取得できる「構造体」について学びました。

いずれ、構造体を使って状態を管理し、メソッドを公開することで、あなたが制御できる方法でユーザーに状態を変更させたくなることもあるでしょう。

**フィンテックはGoが大好き**です。それと、えーと、ビットコインも？ というわけで、どれほど素晴らしい銀行システムを作れるかお見せしましょう。

`Bitcoin`を預金する`Wallet`構造体を作成しましょう。

## 最初にテストを書く

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()
	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

[前の例](structs-methods-and-interfaces.md)では、フィールド名を使用してフィールドに直接アクセスしましたが、 _非常に安全なウォレット_ では、内部状態を外の世界に公開したくありません。メソッドを介してアクセスを制御したいのです。

## テストを実行してみます

`./wallet_test.go:7:12: undefined: Wallet`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

コンパイラは`Wallet`が何であるかを知らないので、それを伝えましょう。

```go
type Wallet struct{}
```

これでウォレットができたので、もう一度テストを実行してみてください

```
./wallet_test.go:9:8: wallet.Deposit undefined (type Wallet has no field or method Deposit)
./wallet_test.go:11:15: wallet.Balance undefined (type Wallet has no field or method Balance)
```

これらのメソッドを定義する必要があります。

テストを実行できるようにするのに必要なことだけを行うよう心がけてください。
テストが正しく失敗し、明確なエラーメッセージが表示されることを確認する必要があります。

```go
func (w Wallet) Deposit(amount int) {

}

func (w Wallet) Balance() int {
	return 0
}
```

この構文に慣れていない場合は、戻って構造体のセクションを読んでください。

テストがコンパイルされ、実行されるはずです

`wallet_test.go:15: got 0 want 10`

## 成功させるのに十分なコードを書く

状態を保存するには、構造体に何らかの _balance_ 変数が必要です

```go
type Wallet struct {
	balance int
}
```

Goでは、シンボル（変数、型、関数など）が小文字で始まる場合、それは _定義されているパッケージの外からは_ プライベートになります。

私たちのケースでは、メソッドがこの値を操作できるようにしたいが、他の誰も操作できないようにしたい。

「レシーバー`receiver`」変数を使用して、構造体の内部の`balance`フィールドにアクセスできることを覚えておいてください。

```go
func (w Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w Wallet) Balance() int {
	return w.balance
}
```

これでフィンテックでのキャリアは安泰です。テストスイートを実行して、テストが成功する喜びに浸りましょう

`wallet_test.go:15: got 0 want 10`

### 何かがおかしい

これはわけがわかりません。コードは動くはずに見えます。
新しい金額を残高に加算し、`Balance`メソッドはその現在の状態を返すはずです。

Goでは、**関数またはメソッドを呼び出すと、引数は _コピーされます_**。

`func (w Wallet) Deposit(amount int)`を呼び出すとき、 `w`はメソッドの呼び出し元のコピーです。

あまりコンピュータサイエンス寄りの話には深入りしませんが、ウォレットのような値を作成すると、それはメモリのどこかに保存されます。そのメモリ領域の _アドレス_ は`&myVal`で確認できます。

コードにいくつかのプリントを追加して実験します

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()

	fmt.Printf("address of balance in test is %p \n", &wallet.balance)

	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

```go
func (w Wallet) Deposit(amount int) {
	fmt.Printf("address of balance in Deposit is %p \n", &w.balance)
	w.balance += amount
}
```

`%p`プレースホルダーはメモリアドレスを先頭に`0x`の付いた16進表記で出力し、`\n`エスケープ文字は改行を出力します。何かのポインタ（メモリアドレス）は、シンボルの先頭に`&`文字を置くことで取得できることに注意してください。

テストを再実行します

```
address of balance in Deposit is 0xc420012268
address of balance in test is 0xc420012260
```

2つの`balance`のアドレスが異なることがわかります。
つまり、コード内で`balance`の値を変更しても、それはテストから渡されたもののコピーを操作しているにすぎません。
そのため、テスト側の`balance`は変更されないのです。

これは _ポインタ_ で修正できます。
[ポインタ](https://gobyexample.com/pointers)を使うと、いくつかの値を _指し示し_ 、それらを変更できます。
したがって、ウォレット全体のコピーを取得するのではなく、代わりにそのウォレットへのポインタを取得して、その中の元の値を変更できるようにします。

```go
func (w *Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w *Wallet) Balance() int {
	return w.balance
}
```

違いは、レシーバーのタイプが`Wallet`ではなく`*Wallet`であり、「ウォレットへのポインタ」と読むことができます。

テストを再実行してみてください。成功するはずです。

なぜ成功したのか、不思議に思うかもしれません。
関数内で、次のようにポインタを逆参照してはいませんでした。

```go
func (w *Wallet) Balance() int {
	return (*w).balance
}
```

そして、一見オブジェクトを直接扱っているように見えます。実際、 `(*w)`を使用した上記のコードは完全に有効です。
ただし、Goの作成者はこの表記を扱いにくいと判断したため、この言語では、明示的な逆参照なしで`w.balance`を記述できます。
構造体へのこれらのポインタには、独自の名前 _構造体ポインタ（struct pointers）_ があり、[自動的に逆参照されます](https://golang.org/ref/spec#Method_values)。

残高のコピーを取っても問題ないので、技術的には`Balance`をポインタレシーバーに変更する必要はありません。
ただし、慣例として、一貫性を保つためにメソッドのレシーバーの型は揃えておくべきです。

## リファクタリング♪

ビットコインのウォレットを作ると言いましたが、ここまでビットコインについては触れていません。
`int`を使ってきたのは、物を数えるのに適した型だからです！

このための`struct`を作成するのは少しやり過ぎのようです。
`int`は動作の点では問題ありませんが、意味が伝わりません。

Goでは、既存の型から新しい型を作成できます。

構文は、`type MyName OriginalType` です。

```go
type Bitcoin int

type Wallet struct {
	balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}
```

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(Bitcoin(10))

	got := wallet.Balance()

	want := Bitcoin(10)

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

`Bitcoin`を作成するには、`Bitcoin(999)`という構文を使用します。

これにより、新しい型を作成し、それに対して _メソッド_ を宣言できます。
これは、既存の型の上にドメイン固有の機能を追加する場合に非常に役立ちます。

[Stringer](https://golang.org/pkg/fmt/#Stringer) をビットコインに実装しましょう

```go
type Stringer interface {
	String() string
}
```

このインターフェースは`fmt`パッケージで定義されており、`%s`フォーマット文字列で出力したときに、型がどのように表示されるかを定義できます。

```go
func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

ご覧のとおり、型宣言でメソッドを作成するための構文は、構造体での構文と同じです。

ここでは規律が欠けていました。最初にテストを書かずにメソッドを追加してしまったのです。それでも構いません。私たちはいつも聖人君子というわけではありませんから。しかし、見逃してしまうべきでもありません。`go test -cover`を実行すれば`String`がカバーされていないことがわかり、後からでもテストする価値があるかを考え直す良いきっかけになります。カバレッジ100%をそれ自体のために追い求めるべきではありませんが、今回の`String`には固定しておく価値のある独自のロジック（`fmt.Sprintf`）があるので、テストを追加しましょう。

```go
t.Run("Bitcoin String", func(t *testing.T) {
	btc := Bitcoin(10)
	got := btc.String()
	want := "10 BTC"

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
})
```

次にテストフォーマット文字列を更新して、代わりに `String()`を使用する必要があります。

```go
if got != want {
	t.Errorf("got %s want %s", got, want)
}
```

この動作を確認するために、わざとテストを失敗させてみましょう。

`wallet_test.go:18: got 10 BTC want 20 BTC`

これにより、テストで何が起こっているかが明確になります。

次の要件は、 `Withdraw`関数です。

## 最初にテストを書く

`Deposit()`のほぼ逆

```go
func TestWallet(t *testing.T) {

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}

		wallet.Deposit(Bitcoin(10))

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}

		wallet.Withdraw(Bitcoin(10))

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
```

## テストを実行してみます

`./wallet_test.go:26:9: wallet.Withdraw undefined (type Wallet has no field or method Withdraw)`

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func (w *Wallet) Withdraw(amount Bitcoin) {

}
```

`wallet_test.go:33: got 20 BTC want 10 BTC`

## 成功させるのに十分なコードを書く

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
	w.balance -= amount
}
```

## リファクタリング♪

テストには重複があります。それをリファクタリングしましょう。

```go
func TestWallet(t *testing.T) {

	assertBalance := func(t testing.TB, wallet Wallet, want Bitcoin) {
		t.Helper()
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

}
```

口座に残っている以上の額を`Withdraw`しようとすると、どうなるべきでしょうか？
今のところ、要件としては当座貸越（残高を超えた引き出し）の機能はないものとします。

`Withdraw`を使用する場合、どのように問題を通知しますか？

Goでは、エラーを示したい場合、呼び出し側がチェックして対処するために関数が `err`を返すことは慣用的です。

これをテストで試してみましょう。

## 最初にテストを書く

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
	startingBalance := Bitcoin(20)
	wallet := Wallet{startingBalance}
	err := wallet.Withdraw(Bitcoin(100))

	assertBalance(t, wallet, startingBalance)

	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
})
```

持っている以上の額を引き出そうとした _場合_ に`Withdraw`がエラーを返し、残高は変わらないままにしたいです。

そして、エラーが`nil`の場合はテストを失敗させることで、エラーが返されたことを確認します。

`nil`は他のプログラミング言語の`null`と同義です。
`Withdraw`の戻り値の型はインターフェースである`error`になるため、エラーは`nil`になりえます。
引数や戻り値の型がインターフェースである関数を見かけたら、それらはnilになりうる（nillable）ということです。

`null`のように`nil`である値にアクセスしようとすると、**ランタイムパニック**がスローされます。
これはまずいです！ 必ず`nil`をチェックするようにしてください。

## テストを実行してみます

`./wallet_test.go:31:25: wallet.Withdraw(Bitcoin(100)) used as value`

エラーの言い回しは少しわかりにくいかもしれませんが、これまでの`Withdraw`は単に呼び出すだけのもので、値を返すことはありませんでした。
コンパイルを通すには、戻り値の型を持つように変更する必要があります。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
	w.balance -= amount
	return nil
}
```

繰り返しになりますが、コンパイラを満足させるのに十分なコードだけを書くことが非常に重要です。
`Withdraw`メソッドを修正して`error`を返すようにしました。
ここでは _何か_ を返す必要があるので、`nil`を返します。

## 成功させるのに十分なコードを書く

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("oh no")
	}

	w.balance -= amount
	return nil
}
```

`errors`をコードにインポートすることを忘れないでください。

`errors.New`は、選択したメッセージで新しい`error`を作成します。

## リファクタリング♪

テストの可読性を高めるために、エラーチェック用のクイックテストヘルパーを作成しましょう。

```go
assertError := func(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}
```

そして、私たちのテストでは

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
	startingBalance := Bitcoin(20)
	wallet := Wallet{startingBalance}
	err := wallet.Withdraw(Bitcoin(100))

	assertError(t, err)
	assertBalance(t, wallet, startingBalance)
})
```

"oh no"というエラーを返したとき、それを返してもあまり役に立ちそうにないので、 _いずれ_ 改善することになるだろうと思っていたのではないでしょうか。

エラーが最終的にユーザーに返されると仮定して、エラーの存在だけでなく、何らかのエラーメッセージを評価するようにテストを更新しましょう。

## 最初にテストを書く

比較する `string`のヘルパーを更新します。

```go
assertError := func(t testing.TB, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got.Error() != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

ご覧のとおり、`Error`は`.Error()`メソッドで文字列に変換できます。ここでは、欲しい文字列と比較するためにそうしています。また、`nil`に対して`.Error()`を呼び出さないように、エラーが`nil`でないことも確認しています。

そして、呼び出し側を更新します。

```go
t.Run("withdraw insufficient funds", func(t *testing.T) {
	startingBalance := Bitcoin(20)
	wallet := Wallet{startingBalance}
	err := wallet.Withdraw(Bitcoin(100))

	assertError(t, err, "cannot withdraw, insufficient funds")
	assertBalance(t, wallet, startingBalance)
})
```

呼び出された場合にテストを停止する `t.Fatal`を導入しました。
エラーが返されていないのに、そのエラーについてさらにアサーションを行いたくはないからです。
これがなければ、テストは次のステップに進み、nilポインタが原因でパニックになります。

## テストを実行してみます

`wallet_test.go:61: got err 'oh no' want 'cannot withdraw, insufficient funds'`

## 成功させるのに十分なコードを書く

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("cannot withdraw, insufficient funds")
	}

	w.balance -= amount
	return nil
}
```

## リファクタリング♪

テストコードと `Withdraw`コードの両方でエラーメッセージが重複しています。

誰かがエラーの文言を変えたいだけなのにテストが失敗するのは本当に煩わしいですし、テストにとっては細かすぎる情報です。
正確な文言は _実のところ_ どうでもよいのです。
大事なのは、特定の条件下で引き出しに関する何らかの意味のあるエラーが返されることだけです。

Goでは、エラーは値なので、それを変数にリファクタリングして、信頼できる唯一の情報源（single source of truth）にすることができます。

```go
var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return ErrInsufficientFunds
	}

	w.balance -= amount
	return nil
}
```

`var`キーワードを使用すると、パッケージにグローバルな値を定義できます。

これで、 `Withdraw`関数が非常に明確になったので、これ自体は前向きな変更です。

次に、特定の文字列の代わりにこの値を使用するようにテストコードをリファクタリングできます。

```go
func TestWallet(t *testing.T) {

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, Bitcoin(20))
	})
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
	t.Helper()
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

これでテストも読みやすくなりました。

ファイルを開いた人がヘルパーではなくアサーションから読み始められるように、ヘルパーをメインのテスト関数の外に移動しました。

テストのもう1つの便利な特性は、コードの _実際の_ 使われ方を理解するのに役立ち、使う側に配慮したコードを書けるようになることです。ここで、開発者は単にコードを呼び出して、 `ErrInsufficientFunds`に対して等号チェックを行い、それに応じて行動できることがわかります。

### 未チェックのエラー

Goのコンパイラは大いに助けてくれますが、それでも見落としてしまうことはありますし、エラー処理は厄介な場合もあります。

テストしていないシナリオが1つあります。これを見つけるには、ターミナルで次のコマンドを実行して、Goで使用できる多くのリンターの1つである`errcheck`をインストールします。

`go install github.com/kisielk/errcheck@latest`

次に、コードを含むディレクトリ内で `errcheck .`を実行します。

次のような出力が得られるはずです

`wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))`

これは、そのコード行で返されているエラーをチェックしていないことを示しています。
私のコンピューターでは、その行は通常の引き出しのシナリオに対応しています。`Withdraw`が成功した場合にエラーが返され _ない_ ことをチェックしていなかったのです。

これに対処した最終的なテストコードは次のとおりです。

```go
func TestWallet(t *testing.T) {

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(10))

		assertNoError(t, err)
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, Bitcoin(20))
	})
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
	t.Helper()
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
```

## エラーのラップでコンテキストを追加する

ここでは`Withdraw`だけが`ErrInsufficientFunds`を生成できるので、`err`を`ErrInsufficientFunds`と比較するのはうまく機能します。しかし、実際のプログラムには層があります。`Wallet`は`Bank`から使われ、`Bank`はHTTPハンドラーから使われる、といった具合です。各層が受け取ったエラーをそのまま返すだけだと、数層上の呼び出し元に見えるのは`insufficient funds`だけで、実際にどのアカウントの、どの操作が失敗したのかはまったくわかりません。

名前付きのアカウントに対して引き出しを処理する関数があるとしましょう。

```go
func ProcessWithdrawal(wallet *Wallet, accountID string, amount Bitcoin) error {
	if err := wallet.Withdraw(amount); err != nil {
		return fmt.Errorf("processing withdrawal for account %s: %w", accountID, err)
	}
	return nil
}
```

`%w`動詞\(Go 1.13で導入\)は、`err`の文字列を埋め込むという点では`%v`と同じですが、`%v`にはないこともします。メッセージを新しい無関係な文字列にコピーするだけでなく、`fmt.Errorf`が返す新しいエラーの中に`err`を _ラップ_ するのです。試してみましょう。

```go
wallet := Wallet{Bitcoin(10)}
err := ProcessWithdrawal(&wallet, "acc-123", Bitcoin(100))
fmt.Println(err)
// processing withdrawal for account acc-123: cannot withdraw, insufficient funds
```

「どのアカウントの、どの操作か」という有用なコンテキストを保ちつつ、その下で実際に何が問題だったのかという詳細も失っていません。

### ラップされたエラーをチェックする

`ProcessWithdrawal`は`ErrInsufficientFunds`とは _異なる_ エラー値を返すようになったので、根本的な原因は同じであっても、`==`\(や上の`assertError`ヘルパー\)で比較すると失敗してしまいます。これこそまさに[`errors.Is`](https://pkg.go.dev/errors#Is)が解決する問題です。`errors.Is`は、一番外側のエラーだけでなく、ラップされたエラーの連鎖をチェックします。

```go
if errors.Is(err, ErrInsufficientFunds) {
	// still true, even though err's message now also mentions the account
}
```

`errors.Is`は、一致するものが見つかるか、アンラップするエラーがなくなるまで、`err`に対して[`errors.Unwrap`](https://pkg.go.dev/errors#Unwrap)を繰り返し呼び出すことで動作します\(`%w`を指定した`fmt.Errorf`は`Unwrap() error`メソッドを持つ値を生成するので、`errors.Unwrap`はラップされたエラーを取り出す方法を知っています\)。まさにこの理由から、[マップ](maps.md)の章の`assertError`ヘルパーでも`==`の代わりに`errors.Is`が使われているのを目にするでしょう。これは、一度もラップされていないエラーに対しても機能する安全なデフォルトです。

\(`ErrInsufficientFunds`のような番兵値と比較するのではなく\)ラップされたエラーの連鎖から _型付きの_ エラーを取り出す必要がある場合は、そのための同等の関数もあります。[`errors.As`](https://pkg.go.dev/errors#As)です。[エラー型](../questions-and-answers/error-types.md)の章で、これについてより詳しく説明しています。

## まとめ

### ポインタ

* Goは、値を関数/メソッドに渡すときに値をコピーするので、状態を変更する必要がある関数を作成している場合は、変更したいもののポインタを受け取るようにする必要があります。
* Goが値のコピーを取得するという事実は、多くの場合有用ですが、システムに何かのコピーを作成させたくない場合があり、その場合は参照を渡す必要があります。例としては、非常に大きなデータ構造を参照する場合や、\(データベース接続プールのように\)インスタンスが1つだけ必要なものが挙げられます。

### nil

* ポインタはnilにすることができます
* 関数が何かへのポインタを返すときは、それが`nil`かどうかを必ずチェックする必要があります。そうしないとランタイム例外が発生する可能性があり、ここではコンパイラは助けてくれません。
* 存在しない可能性のある値を表現したい場合に便利です

### エラー

* エラーは、関数/メソッドを呼び出すときの失敗を示す方法です。
* テストの声に耳を傾けた結果、エラーの文字列をチェックすると不安定なテストになると結論付けました。そのため、代わりに意味のある値を使用するように実装をリファクタリングしました。これにより、コードのテストが容易になり、APIのユーザーにとっても簡単になると結論付けました。
* これはエラー処理の話の終わりではなく、より高度なことを行うことができますが、これは単なる紹介です。以降のセクションでは、より多くの戦略について説明します。
* [エラーをチェックするだけでなく、適切に処理する](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### 既存の型から新しい型を作成する

* 値にドメイン固有の意味を追加するのに役立ちます
* インターフェースを実装できます

ポインタとエラーはGoを書く上で重要な部分であり、慣れる必要があります。
ありがたいことに、何か間違ったことをしても _たいていは_ コンパイラが助けてくれます。
時間をかけてエラーを読んでください。
