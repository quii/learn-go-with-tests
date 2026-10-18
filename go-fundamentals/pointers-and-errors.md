---
description: Pointers & errors
---

# ポインタとエラー

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/pointers)

前のセクションでコンセプトに関連する、いくつかの値を取得できる「構造体」について学びました。

ある時点で、構造体を使用して状態を管理し、ユーザーが制御できる方法で状態を変更できるようにするメソッドを公開することができます。

**金融技術はGoを愛してます** 、、ビットコイン? それでは、私たちがどんな素晴らしい銀行システムを作ることができるかを示しましょう。

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

これで財布ができました。もう一度テストを実行してください

```
./wallet_test.go:9:8: wallet.Deposit undefined (type Wallet has no field or method Deposit)
./wallet_test.go:11:15: wallet.Balance undefined (type Wallet has no field or method Balance)
```

これらのメソッドを定義する必要があります。

テストを実行するのに十分なことだけを忘れないでください。
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

Goでは、シンボル（変数`var`、タイプ`type`、関数`func`）が小文字の記号で始まっている場合は、それは定義されているパッケージの外側のプライベートなものです。

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

フィンテックでのキャリアが確保されたら、テストスイートを実行し、合格したテストを浴びます

`wallet_test.go:15: got 0 want 10`

### 何かがおかしい

これは混乱を招きます。
コードは機能するように見え、新しい金額を残高に追加し、次に`balance`メソッドはその現在の状態を返す必要があります。

Goでは、**関数またはメソッドを呼び出すと、引数は _コピーされます_**。

`func (w Wallet) Deposit(amount int)`を呼び出すとき、 `w`はメソッドの呼び出し元のコピーです。

あまりにもコンピュータサイエンシーにならずに、ウォレットのように値を作成すると、メモリのどこかに保存されます。そのメモリのビットの _アドレス_ は、 `&myVal`で確認できます。

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

2つの`balances`のアドレスが異なることがわかります。
したがって、コード内の`balances`の値を変更する場合、テストから得られたもののコピーに取り組んでいます。
したがって、テストのバランスは変更されません。

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

違いは、レシーバーのタイプが`Wallet`ではなく`*Wallet`であり、**`Wallet`へのポインター**として読み取ることができることです。

テストを試行して再実行すると、テストに合格するはずです。

不思議に思うかもしれませんが、なぜ合格したのですか？
次のように、関数のポインターを逆参照しませんでした。

```go
func (w *Wallet) Balance() int {
	return (*w).balance
}
```

一見オブジェクトに直接対処したようです。実際、 `(*w)`を使用した上記のコードは完全に有効です。
ただし、Goの作成者はこの表記を扱いにくいと判断したため、この言語では、明示的な逆参照なしで`w.balance`を記述できます。
構造体へのこれらのポインタには、独自の名前 _構造体ポインタ（struct pointers）_ があり、[自動的に逆参照されます](https://golang.org/ref/spec#Method_values)。

技術的には、バランスのコピーを取ることは問題ないので、ポインターレシーバーを使用するために`Balance`を変更する必要はありません。
ただし、慣例では、一貫性を保つために、メソッドレシーバーのタイプを同じに保つ必要があります。

## リファクタリング♪

私たちはビットコインの財布を作っていると述べましたが、これまでのところ言及していません。
`int`を使用しているのは、物事を数えるのに適したタイプだからです。

このための`struct`を作成するのは少しやり過ぎのようです。
`int`は、動作の点では問題ありませんが、説明的ではありません。

Goでは、既存のタイプから新しいタイプを作成できます。

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

これにより、新しい型を作成し、それらに対して _methods_ を宣言できます。
これは、既存のタイプの上にドメイン固有の機能を追加する場合に非常に役立ちます。

[Stringer](https://golang.org/pkg/fmt/#Stringer) をビットコインに実装しましょう

```go
type Stringer interface {
	String() string
}
```

このインターフェイスは `fmt`パッケージで定義されており、`%s`形式の文字列で出力したときにタイプがどのように印刷されるかを定義できます。

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

この動作を確認するには、意図的にテストを中断して、確認できるようにします。

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

アカウントに残っている以上に「撤回`Withdraw`」しようとするとどうなりますか？
当面の要件は、当座貸越施設がないことを前提としています。

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

`Withdraw`がエラーを返すようにしたいのですが、あなたが持っている以上のものを取り出そうとしても、`balance`は変わらないはずです。

次に、 `nil`の場合、テストに失敗してエラーが返されたことを確認します。

`nil`は他のプログラミング言語の`null`と同義です。
`Withdraw`の戻り値の型はインターフェイスである`error`になるため、エラーは `nil`になる可能性があります。
引数を取る関数、またはインターフェイスである値を返す関数がある場合、それらは nillable(潰しが利かない) になる可能性があります。

`null`のように`nil`である値にアクセスしようとすると、**ランタイムパニック**がスローされます。
これは悪いです！ `nil`をチェックすることを確認する必要があります。

## テストを実行してみます

`./wallet_test.go:31:25: wallet.Withdraw(Bitcoin(100)) used as value`

言い回しは少し不明瞭かもしれませんが、`Withdraw`を使用する以前の意図は単にそれを呼び出すことであり、決して値を返しません。
このコンパイルを行うには、戻り値の型を持つように変更する必要があります。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
	w.balance -= amount
	return nil
}
```

繰り返しになりますが、コンパイラーを満足させるのに十分なコードを書くことが非常に重要です。
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

うまくいけば、 "oh no"のエラーを返すときに、返すのはそれほど便利ではないように思われるので、私たちはそのことを繰り返していると考えていました。

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

そして、発信者を更新します。

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
これは、周りにエラーがない場合に返されるエラーについてこれ以上アサーションを作成したくないためです。
これがなければ、テストは次のステップに進み、nilポインターのためにパニックになります。

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

誰かがエラーを書き直したい場合、テストが失敗するのは本当にうっとうしいことであり、それは私たちのテストには余りにも詳細です。
正確な表現が何であるかについては実際には気にしません。
特定の条件が与えられた場合、引き出しに関する何らかの意味のあるエラーが返されるだけです。

Goでは、エラーは値なので、それを変数にリファクタリングして、単一の真のソースを持つことができます。

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

そして今、テストも従うのが簡単です。

ヘルパーをメインのテスト関数から移動したので、誰かがファイルを開いたときに、ヘルパーではなく、最初にアサーションの読み取りを開始できます。

テストのもう1つの便利な特性は、コードの _実際の_ 使用法を理解するのに役立ち、同情的なコードを作成できることです。ここで、開発者は単にコードを呼び出して、 `ErrInsufficientFunds`に対して等号チェックを行い、それに応じて行動できることがわかります。

### 未チェックのエラー

Goコンパイラーは大いに役立ちますが、まだ見逃していて、エラー処理が難しい場合もあります。

テストしていないシナリオが1つあります。これを見つけるには、ターミナルで次のコマンドを実行して、Goで使用できる多くのリンターの1つである`errcheck`をインストールします。

`go install github.com/kisielk/errcheck@latest`

次に、コードを含むディレクトリ内で `errcheck .`を実行します。

あなたは次のようなものを取得する必要があります

`wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))`

これは、そのコード行で返されているエラーをチェックしていないことを示しています。
私のコンピューターのそのコード行は、通常の _withdraw_ シナリオに対応しています。`Withdraw`が成功した場合にエラーが返されないことを確認していないためです。

これを説明する最終的なテストコードを次に示します。

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

* Goは、値を関数/メソッドに渡すときに値をコピーするので、状態を変更する必要がある関数を作成している場合は、変更したいものへのポインターを取得する必要があります。
* Goが値のコピーを取得するという事実は、多くの場合有用ですが、システムに何かのコピーを作成させたくない場合があり、その場合は参照を渡す必要があります。例としては、非常に大きなデータ構造を参照する場合や、\(データベース接続プールのように\)インスタンスが1つだけ必要なものが挙げられます。

### nil

* ポインタはnilにすることができます
* 関数が何かへのポインターを返すとき、それが`nil`であるかどうかを確認する必要があります。そうでない場合、ランタイム例外が発生する可能性があります。コンパイラーはここでは役立ちません。
* 欠落している可能性のある値を説明する場合に役立ちます

### エラー

* エラーは、関数/メソッドを呼び出すときの失敗を示す方法です。
* テストを聞いて、エラーの文字列をチェックすると不安定なテストになると結論付けました。そのため、代わりに意味のある値を使用するように実装をリファクタリングしました。これにより、コードのテストが容易になり、APIのユーザーにとっても簡単になると結論付けました。
* これはエラー処理の話の終わりではなく、より高度なことを行うことができますが、これは単なる紹介です。以降のセクションでは、より多くの戦略について説明します。
* [エラーをチェックするだけでなく、適切に処理する](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### 既存のものから新しいタイプを作成する

* 値にドメイン固有の意味を追加するのに役立ちます
* インターフェイスを実装できます

ポインタとエラーはGoを書く上で重要な部分であり、慣れる必要があります。
ありがたいことに、コンパイラは、通常、何か間違ったことをした場合にあなたを助けます。
時間をかけてエラーを読んでください。
