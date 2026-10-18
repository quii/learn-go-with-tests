---
description: Intro to property based tests
---

# プロパティベースのテスト概要

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/roman-numerals)

一部の企業は、インタビュープロセスの一環として、[ローマ数字のカタ](http://codingdojo.org/kata/RomanNumerals/)を実行するように求めます。この章では、TDDでこれに取り組む方法を示します。

[アラビア数字](https://en.wikipedia.org/wiki/Arabic_numerals)（数字0〜9）をローマ数字に変換する関数を記述します。

[ローマ数字](https://en.wikipedia.org/wiki/Roman_numerals)について聞いたことがない場合は、ローマ人が数字を書き留めた方法です。

あなたはシンボルを一緒に貼り付けることによってそれらを構築し、それらのシンボルは数字を表します

つまり`I`は「1」です。`III`は「3」です。

簡単に見えますが、いくつか興味深いルールがあります。
`V`は「5」を意味しますが、`IV`は「4」です（ `IIII`ではありません）。

`MCMLXXXIV`は「1984」です。それは複雑に見え、これを最初から理解するためのコードをどのように書くことができるか想像することは困難です。

この本で強調しているように、ソフトウェア開発者にとって重要なスキルは、「有用な」機能の「薄い垂直スライス」を特定して特定し、**繰り返す**ことです。 TDDワークフローは、反復的な開発を容易にするのに役立ちます。

したがって、「1984」ではなく、「1」から始めましょう。

## 最初にテストを書く

```go
func TestRomanNumerals(t *testing.T) {
    got := ConvertToRoman(1)
    want := "I"

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

あなたがこの本でこれまでに持っているならば、これはうまくいけばあなたにとって非常に退屈で日常的な感じです。それは良いことです。

## テストを実行してみます

```console
./numeral_test.go:6:9: undefined: ConvertToRoman
```

コンパイラーに道を案内する

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

関数を作成しますが、まだテストに合格しないでください。常に、期待どおりにテストが失敗することを確認してください。

```go
func ConvertToRoman(arabic int) string {
    return ""
}
```

今すぐ実行されます。

```console
=== RUN   TestRomanNumerals
--- FAIL: TestRomanNumerals (0.00s)
    numeral_test.go:10: got '', want 'I'
FAIL
```

## 成功させるのに十分なコードを書く

```go
func ConvertToRoman(arabic int) string {
    return "I"
}
```

## リファクタリング♪

まだリファクタリングする必要はありません。

結果をハードコーディングするだけでは変だと感じますが、TDDではできるだけ長く「赤字」を避けたいと考えています。あまり達成していないように感じるかもしれませんが、APIを定義し、ルールの1つをキャプチャするテストを取得しました。「実際の」コードがかなりばかげていても。

ここで、その不安な気持ちを使って新しいテストを記述し、少しだけ馬鹿なコードを書くように強制します。

## 最初にテストを書く

サブテストを使用してテストを適切にグループ化できます

```go
func TestRomanNumerals(t *testing.T) {
    t.Run("1 gets converted to I", func(t *testing.T) {
        got := ConvertToRoman(1)
        want := "I"

        if got != want {
            t.Errorf("got %q, want %q", got, want)
        }
    })

    t.Run("2 gets converted to II", func(t *testing.T) {
        got := ConvertToRoman(2)
        want := "II"

        if got != want {
            t.Errorf("got %q, want %q", got, want)
        }
    })
}
```

## テストを実行してみます

```console
=== RUN   TestRomanNumerals/2_gets_converted_to_II
    --- FAIL: TestRomanNumerals/2_gets_converted_to_II (0.00s)
        numeral_test.go:20: got 'I', want 'II'
```

それほど驚きはありません。

## 成功させるのに十分なコードを書く

```go
func ConvertToRoman(arabic int) string {
    if arabic == 2 {
        return "II"
    }
    return "I"
}
```

ええ、まだ問題に取り組んでいないようです。したがって、前進させるために、さらに多くのテストを作成する必要があります。

## リファクタリング♪

テストでいくつかの繰り返しがあります。 「与えられた入力X、Yを期待する」の問題のように感じる何かをテストしているときは、おそらくテーブルベースのテストを使用する必要があります。

```go
func TestRomanNumerals(t *testing.T) {
    cases := []struct {
        Description string
        Arabic      int
        Want        string
    }{
        {"1 gets converted to I", 1, "I"},
        {"2 gets converted to II", 2, "II"},
    }

    for _, test := range cases {
        t.Run(test.Description, func(t *testing.T) {
            got := ConvertToRoman(test.Arabic)
            if got != test.Want {
                t.Errorf("got %q, want %q", got, test.Want)
            }
        })
    }
}
```

これ以上テストのボイラプレートを書かなくても、簡単にケースを追加できるようになりました。

頑張って「3」に行きましょう

## 最初にテストを書く

以下をケースに追加してください

```
{"3 gets converted to III", 3, "III"},
```

## テストを実行してみます

```console
=== RUN   TestRomanNumerals/3_gets_converted_to_III
    --- FAIL: TestRomanNumerals/3_gets_converted_to_III (0.00s)
        numeral_test.go:20: got 'I', want 'III'
```

## 成功させるのに十分なコードを書く

```go
func ConvertToRoman(arabic int) string {
    if arabic == 3 {
        return "III"
    }
    if arabic == 2 {
        return "II"
    }
    return "I"
}
```

## リファクタリング♪

わかりましたので、これらのifステートメントを楽しんでいないようになり、コードを十分に見てみると、`arabic`のサイズに基づいて`I`の文字列を構築していることがわかります。

より複雑な数値については、ある種の算術および文字列連結を行うことを "知っています"。

これらの考えを念頭に置いてリファクタリングを試してみましょう。
それは最終的なソリューションには適さないかもしれませんが、それは問題ありません。
私たちはいつでもコードを捨てて、私たちをガイドする必要のあるテストからやり直すことができます。

```go
func ConvertToRoman(arabic int) string {

    var result strings.Builder

    for i:=0; i<arabic; i++ {
        result.WriteString("I")
    }

    return result.String()
}
```

[`strings.Builder`](https://golang.org/pkg/strings/#Builder)は、[ベンチマーク](iteration.md#ベンチマーク)についての説明で出てきたことを覚えているかもしれません

> Builderは、Writeメソッドを使用して文字列を効率的に構築するために使用されます。メモリのコピーを最小限に抑えます。

通常、実際のパフォーマンスの問題が発生するまで、このような最適化に悩まされることはありませんが、コードの量は文字列に追加する「手動」よりも大きくないため、より高速なアプローチを使用することもできます。

コードは私にはよく見え、ドメインを _私たちが今知っているように_ 記述しています。

### ローマ人もDRYに夢中になりました...

物事は今より複雑になり始めています。ローマ人はその知恵の中で、繰り返し登場する人物は読みにくく、数え難くなるだろうと考えていました。したがって、ローマ数字のルールでは、同じ文字を3回以上続けて繰り返すことはできません。

代わりに、次に高いシンボルを取り、その左側にシンボルを置くことで「減算」します。すべてのシンボルを減算器として使用できるわけではありません。I（1）、X（10）、C（100）のみ。

たとえば、ローマ数字の「5」は`V`です。 「4」を作成するには、 `IIII`ではなく、`IV`を実行します。

減算器は、その「系列」で次に大きい2つの記号のいずれかの前にしか置くことができません。`I`は`V`または`X`の前にのみ置くことができ（4は`IV`、9は`IX`）、`X`は`L`または`C`の前にのみ（40は`XL`、90は`XC`）、`C`は`D`または`M`の前にのみ置くことができます（400は`CD`、900は`CM`）。つまり、`I`、`X`、`C`はすべて有効な減算器ですが、自由に組み合わせることはできません。99は`IC`ではなく（`XCIX`、つまり90 + 9です）、499は`ID`ではありません（`CDXCIX`、つまり400 + 90 + 9です）。

## 最初にテストを書く

```
{"4 gets converted to IV (can't repeat more than 3 times)", 4, "IV"},
```

## テストを実行してみます

```console
=== RUN   TestRomanNumerals/4_gets_converted_to_IV_(cant_repeat_more_than_3_times)
    --- FAIL: TestRomanNumerals/4_gets_converted_to_IV_(cant_repeat_more_than_3_times) (0.00s)
        numeral_test.go:24: got 'IIII', want 'IV'
```

## 成功させるのに十分なコードを書く

```go
func ConvertToRoman(arabic int) string {

    if arabic == 4 {
        return "IV"
    }

    var result strings.Builder

    for i:=0; i<arabic; i++ {
        result.WriteString("I")
    }

    return result.String()
}
```

## リファクタリング♪

文字列の構築パターンを壊したことが「好き」ではないので、それを続けたいと思います。

```go
func ConvertToRoman(arabic int) string {

    var result strings.Builder

    for i := arabic; i > 0; i-- {
        if i == 4 {
            result.WriteString("IV")
            break
        }
        result.WriteString("I")
    }

    return result.String()
}
```

「4」を現在の考えに「合わせる」ために、アラビア数字からカウントダウンし、進行中に文字列に記号を追加します。これが長期的に機能するかどうかはわかりませんが、見てみましょう！

「５」を作成しましょう

## 最初にテストを書く

```
{"5 gets converted to V", 5, "V"},
```

## テストを実行してみます

```console
=== RUN   TestRomanNumerals/5_gets_converted_to_V
    --- FAIL: TestRomanNumerals/5_gets_converted_to_V (0.00s)
        numeral_test.go:25: got 'IIV', want 'V'
```

## 成功させるのに十分なコードを書く

「4」で行ったアプローチをコピーするだけです

```go
func ConvertToRoman(arabic int) string {

    var result strings.Builder

    for i := arabic; i > 0; i-- {
        if i == 5 {
            result.WriteString("V")
            break
        }
        if i == 4 {
            result.WriteString("IV")
            break
        }
        result.WriteString("I")
    }

    return result.String()
}
```

## リファクタリング♪

このようなループでの繰り返しは、通常、呼び出されるのを待っている抽象化の兆候です。ループを短絡することは読みやすさのための効果的なツールであるかもしれませんが、それはまたあなたに何か他のものを伝えているかもしれません。

アラビア数字をループしていて、特定の記号にぶつかった場合は`break`と呼びますが、実際に行っているのは、`i`を手間をかけて減算することです。

```go
func ConvertToRoman(arabic int) string {

    var result strings.Builder

    for arabic > 0 {
        switch {
        case arabic > 4:
            result.WriteString("V")
            arabic -= 5
        case arabic > 3:
            result.WriteString("IV")
            arabic -= 4
        default:
            result.WriteString("I")
            arabic--
        }
    }

    return result.String()
}
```

* いくつかの非常に基本的なシナリオのテストから得られたコードから読み取っている信号を考えると、ローマ数字を作成するには、シンボルを適用するときに「アラビア語」から減算する必要があることがわかります
* `for`ループはもはや`i`に依存せず、代わりに、`arabic`から十分な数のシンボルを減算するまで文字列を構築し続けます。

このアプローチが「6」（VI）、「7」（VII）、「8」（VIII）にも有効であると確信しています。
それでも、ケースをテストスイートに追加し、（簡潔にするためにコードは含めません。不明な場合はgithubのサンプルを確認してください）を確認してください。

「9」は、次の数の表現から「I」を減算する必要があるという点で、「4」と同じルールに従います。 「10」はローマ数字で「X」で表されます。したがって、「9」は「IX」になります。

## 最初にテストを書く

```
{"9 gets converted to IX", 9, "IX"},
```

## テストを実行してみます

```console
=== RUN   TestRomanNumerals/9_gets_converted_to_IX
    --- FAIL: TestRomanNumerals/9_gets_converted_to_IX (0.00s)
        numeral_test.go:29: got 'VIV', want 'IX'
```

## 成功させるのに十分なコードを書く

以前と同じアプローチを採用できるはずです

```
case arabic > 8:
    result.WriteString("IX")
    arabic -= 9
```

## リファクタリング♪

それはコードがまだどこかにリファクタリングがあることを私たちに伝えているように感じますが、私には完全に明白ではないので、続けましょう。

これもコードはスキップしますが、テストケースに「10」のテストを追加します。「10」は「X」である必要があり、先に進む前に合格にします。

「39」までのコードが機能すると確信しているので、ここに追加したいくつかのテストがあります。

```
{"10 gets converted to X", 10, "X"},
{"14 gets converted to XIV", 14, "XIV"},
{"18 gets converted to XVIII", 18, "XVIII"},
{"20 gets converted to XX", 20, "XX"},
{"39 gets converted to XXXIX", 39, "XXXIX"},
```

これまでにオブジェクト指向プログラミングを行ったことがある場合は、少し疑惑を抱いて`switch`ステートメントを表示する必要があることがわかります。通常、実際には代わりにクラス構造でキャプチャできる場合でも、一部の命令コード内でコンセプトまたはデータをキャプチャしています。

Goは厳密にはオブジェクト指向ではありませんが、オブジェクト指向が提供するレッスンを完全に無視することを意味するわけではありません（あなたが伝えたいだけのことです）。

私たちの切り替えステートメントは、動作とともにローマ数字に関するいくつかの真実を説明しています。

データを動作から分離することで、これをリファクタリングできます。

```go
type RomanNumeral struct {
    Value  int
    Symbol string
}

var allRomanNumerals = []RomanNumeral {
    {10, "X"},
    {9, "IX"},
    {5, "V"},
    {4, "IV"},
    {1, "I"},
}

func ConvertToRoman(arabic int) string {

    var result strings.Builder

    for _, numeral := range allRomanNumerals {
        for arabic >= numeral.Value {
            result.WriteString(numeral.Symbol)
            arabic -= numeral.Value
        }
    }

    return result.String()
}
```

これはずっと気分が良いです。数値に関連するいくつかのルールをアルゴリズムで非表示にするのではなくデータとして宣言しました。アラビア数字を処理して、適合する場合は結果に記号を追加する方法を確認できます。

この抽象化はより大きな数で機能しますか？ 「50」のローマ数字（`L`）で機能するようにテストスイートを拡張します。

ここにいくつかのテストケースがあります。

```
{"40 gets converted to XL", 40, "XL"},
{"47 gets converted to XLVII", 47, "XLVII"},
{"49 gets converted to XLIX", 49, "XLIX"},
{"50 gets converted to L", 50, "L"},
```

助けが必要？

追加するシンボルは[この要点](https://gist.github.com/pamelafox/6c7b948213ba55332d86efd0f0b037de)で確認できます。

## そして残りの部分！

残りの記号は次のとおりです

| Arabic | Roman |
| :--- | :---: |
| 100 | C |
| 500 | D |
| 1000 | M |

残りのシンボルについても同じ方法を使用します。テストとシンボルの配列の両方にデータを追加するだけです。

あなたのコードは`1984`:`MCMLXXXIV`で動作しますか？

これが私の最後のテストスイートです

```go
func TestRomanNumerals(t *testing.T) {
    cases := []struct {
        Arabic int
        Roman  string
    }{
        {Arabic: 1, Roman: "I"},
        {Arabic: 2, Roman: "II"},
        {Arabic: 3, Roman: "III"},
        {Arabic: 4, Roman: "IV"},
        {Arabic: 5, Roman: "V"},
        {Arabic: 6, Roman: "VI"},
        {Arabic: 7, Roman: "VII"},
        {Arabic: 8, Roman: "VIII"},
        {Arabic: 9, Roman: "IX"},
        {Arabic: 10, Roman: "X"},
        {Arabic: 14, Roman: "XIV"},
        {Arabic: 18, Roman: "XVIII"},
        {Arabic: 20, Roman: "XX"},
        {Arabic: 39, Roman: "XXXIX"},
        {Arabic: 40, Roman: "XL"},
        {Arabic: 47, Roman: "XLVII"},
        {Arabic: 49, Roman: "XLIX"},
        {Arabic: 50, Roman: "L"},
        {Arabic: 100, Roman: "C"},
        {Arabic: 90, Roman: "XC"},
        {Arabic: 400, Roman: "CD"},
        {Arabic: 500, Roman: "D"},
        {Arabic: 900, Roman: "CM"},
        {Arabic: 1000, Roman: "M"},
        {Arabic: 1984, Roman: "MCMLXXXIV"},
        {Arabic: 3999, Roman: "MMMCMXCIX"},
        {Arabic: 2014, Roman: "MMXIV"},
        {Arabic: 1006, Roman: "MVI"},
        {Arabic: 798, Roman: "DCCXCVIII"},
    }
    for _, test := range cases {
        t.Run(fmt.Sprintf("%d gets converted to %q", test.Arabic, test.Roman), func(t *testing.T) {
            got := ConvertToRoman(test.Arabic)
            if got != test.Roman {
                t.Errorf("got %q, want %q", got, test.Roman)
            }
        })
    }
}
```

* データに十分な情報が記載されていると感じたため、「コメント」を削除しました。
* もう少し自信を与えるために見つけた他のエッジケースをいくつか追加しました。テーブルベースのテストでは、これは非常に安価です。

アルゴリズムは変更しませんでした。`allRomanNumerals`配列を更新するだけで済みました。

```go
var allRomanNumerals = []RomanNumeral{
    {1000, "M"},
    {900, "CM"},
    {500, "D"},
    {400, "CD"},
    {100, "C"},
    {90, "XC"},
    {50, "L"},
    {40, "XL"},
    {10, "X"},
    {9, "IX"},
    {5, "V"},
    {4, "IV"},
    {1, "I"},
}
```

## ローマ数字の解析

まだ終わっていません。次に、ローマ数字の _from_ を`int`に変換する関数を書きます

## 最初にテストを書く

テストケースを少しリファクタリングして再利用できます

`case`変数を`var`ブロックのパッケージ変数としてテストの外に移動します。

```go
func TestConvertingToArabic(t *testing.T) {
    for _, test := range cases[:1] {
        t.Run(fmt.Sprintf("%q gets converted to %d", test.Roman, test.Arabic), func(t *testing.T) {
            got := ConvertToArabic(test.Roman)
            if got != test.Arabic {
                t.Errorf("got %d, want %d", got, test.Arabic)
            }
        })
    }
}
```

スライス機能を使用して、今のところテストの1つだけを実行していることに注意してください。（`cases [:1]`）これらのテストを一度にすべて合格にしようとすると、飛躍的に大きくなります。

## テストを実行してみます

```console
./numeral_test.go:60:11: undefined: ConvertToArabic
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

新しい関数定義を追加します

```go
func ConvertToArabic(roman string) int {
    return 0
}
```

テストが実行され、失敗するはずです

```console
--- FAIL: TestConvertingToArabic (0.00s)
    --- FAIL: TestConvertingToArabic/'I'_gets_converted_to_1 (0.00s)
        numeral_test.go:62: got 0, want 1
```

## 成功させるのに十分なコードを書く

あなたは何をするべきか知っています

```go
func ConvertToArabic(roman string) int {
    return 1
}
```

次に、テストのスライスインデックスを変更して、次のテストケースに移動します（例:`cases [:2]`）。
3つ目のケースについても、考えられる最もおかしなコードを理解し、（最高の本で間違いありませんか？）これが私のばかげたコードです。

```go
func ConvertToArabic(roman string) int {
    if roman == "III" {
        return 3
    }
    if roman == "II" {
        return 2
    }
    return 1
}
```

機能する実際のコードの愚かさを通して、以前のようなパターンを見ることができます。入力を反復処理して _something_ を構築する必要があります。この場合は合計です。

```go
func ConvertToArabic(roman string) int {
    total := 0
    for range roman {
        total++
    }
    return total
}
```

## 最初にテストを書く

次に、`cases[:4]`（`IV`）に移動します。これは、文字列の長さである2を返すため、失敗します。

## 成功させるのに十分なコードを書く

```go
// earlier..
var allRomanNumerals = RomanNumerals{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

// later..
func ConvertToArabic(roman string) int {
	var arabic = 0

	for _, numeral := range allRomanNumerals {
		for strings.HasPrefix(roman, numeral.Symbol) {
			arabic += numeral.Value
			roman = strings.TrimPrefix(roman, numeral.Symbol)
		}
	}

	return arabic
}
```

これは基本的に、`ConvertToRoman(int)`のアルゴリズムを逆向きに実装したものです。ここでは、与えられたローマ数字の文字列に対してループします。
* `allRomanNumerals`から取り出したローマ数字の記号を、大きいものから小さいものの順に、文字列の先頭で探します。
* 接頭辞が見つかったら、その値を`arabic`に加え、接頭辞を取り除きます。

最後に、合計をアラビア数字として返します。

`HasPrefix(s, prefix)`は文字列`s`が`prefix`で始まるかどうかをチェックし、`TrimPrefix(s, prefix)`は`s`から`prefix`を取り除きます。これにより、残りのローマ数字の記号の処理を続けることができます。これは`IV`や他のすべてのテストケースで機能します。

これを再帰関数として実装することもできます。その方が（私の意見では）エレガントですが、遅くなるかもしれません。これはあなたと`Benchmark...`テストに任せることにします。

アラビア数字をローマ数字に変換する関数とその逆の関数ができたので、テストをさらに一歩進めることができます。

## プロパティベースのテストの概要

この章で使用したローマ数字のドメインにはいくつかのルールがあります。

* 3つ以上の連続したシンボルは使用できません
* `I` （1）、 `X` （10） 、 `C` （100） のみが「減算器」になります
* `ConvertToRoman(N)`の結果を取得して`ConvertToArabic`に渡すと、`N`が返されます

これまでに作成したテストは、「例」ベースのテストとして説明できます。ツールが検証するための _例_ を私たちが提供しているからです。

ドメインについて知っているこれらのルールを採用し、コードに対して何らかの方法でそれらを実行できるとしたらどうでしょうか。

プロパティベースのテストは、コードにランダムデータを投げ、記述したルールが常に正しいことを確認することで、これを行うのに役立ちます。多くの人は、プロパティベースのテストは主にランダムデータについてであると考えていますが、それは間違いです。プロパティベースのテストに関する本当の課題は、ドメインをよく理解して、これらのプロパティを記述できるようにすることです。

説明はこれくらいにして、コードを見てみましょう。

> **⚠️ Linuxユーザーの方へ:** 以下のテストを **すぐには実行しないでください**。システムがフリーズする（ハードリブートが必要になる）可能性が高いです。
>
> <details>
> <summary>理由を見るにはここをクリック（技術的な説明）</summary>
>
> `testing/quick`パッケージは`int64`の最大値までのランダムな整数を生成します。現在の素朴な実装は、その長さの文字列（数千兆文字）をメモリ上に構築しようとします。
>
> macOSやWindowsではたいていうまく処理されます（UIは応答し続けます）が、Linuxカーネルでは通常「スワップスラッシング」が発生し、プロセスを強制終了する前にシステム全体がフリーズしてしまいます。
> </details>

```go
func TestPropertiesOfConversion(t *testing.T) {
    assertion := func(arabic int) bool {
        roman := ConvertToRoman(arabic)
        fromRoman := ConvertToArabic(roman)
        return fromRoman == arabic
    }

    if err := quick.Check(assertion, nil); err != nil {
        t.Error("failed checks", err)
    }
}
```

### プロパティの根拠

最初のテストでは、数値をローマ数字に変換する場合に、他の関数を使用して、最初に取得した数値に変換するときにチェックします。

* 与えられた乱数（例えば`4`）。
* 乱数で`ConvertToRoman`を呼び出します（`4`の場合は`IV`を返す必要があります）。
* 上記の結果を受け取り、それを`ConvertToArabic`に渡します。
* 上記により、元の入力（`4`）が得られます。

これは、どちらかにバグがあると壊れるので、自信をつける良いテストのように感じます。合格できる唯一の方法は、彼らに同じ種類のバグがあるかどうかです。これは不可能ではありませんが、ありそうにありません。

### 技術的な説明

標準ライブラリの[testing/quick](https://golang.org/pkg/testing/quick/)パッケージを使用しています

下から読んで、`quick.Check`関数をいくつかのランダムな入力に対して実行する関数を提供します。関数が`false`を返す場合、チェックに失敗したと見なされます。

上記の`assertion`関数は乱数を受け取り、関数を実行してプロパティをテストします。

### テストを実行する

実行してみてください。お使いのコンピュータがしばらくハングアップする可能性があるので、退屈したら強制終了してください :)

どうしたの？
以下をアサーションコードに追加してみてください。

```go
assertion := func(arabic int) bool {
    if arabic <0 || arabic > 3999 {
        log.Println(arabic)
        return true
    }
    roman := ConvertToRoman(arabic)
    fromRoman := ConvertToArabic(roman)
    return fromRoman == arabic
}
```

次のようなものが表示されます。

```console
=== RUN   TestPropertiesOfConversion
2019/07/09 14:41:27 6849766357708982977
2019/07/09 14:41:27 -7028152357875163913
2019/07/09 14:41:27 -6752532134903680693
2019/07/09 14:41:27 4051793897228170080
2019/07/09 14:41:27 -1111868396280600429
2019/07/09 14:41:27 8851967058300421387
2019/07/09 14:41:27 562755830018219185
```

この非常に単純なプロパティを実行するだけで、実装に欠陥が明らかになりました。入力として`int`を使用しましたが、

* ローマ数字では負の数を実行できません
* 最大3つの連続する記号という規則を考慮すると、「3999」を超える値を表すことはできません（[まあ、ちょっと](https://www.quora.com/Which-is-the-maximum-number-in-Roman-numerals)）および `int`の最大値は「3999」よりはるかに大きくなっています。

これは素晴らしい！
プロパティベースのテストの真の強みであるドメインについて、より深く考えることを余儀なくされました。

明らかに`int`は素晴らしいタイプではありません。もう少し適切なものを試した場合はどうなりますか？

### [`uint16`](https://golang.org/pkg/builtin/#uint16)

Goには _unsigned integers_ の型があります。つまり、負の値にはできません。
これにより、コード内の1つのクラスのバグがすぐに除外されます。
16を追加することは、最大「65535」を格納できる16ビット整数であることを意味します。
これはまだ大きすぎますが、必要なものに近づきます。

コードを更新して`int`ではなく`uint16`を使うようにしてみてください。テストの`assertion`を更新して、もう少し見やすくしました。

> コード内の変数`arabic`も`uint16`に変更する必要があることに注意してください（テストが教えてくれます）。もう少し手間がかかるかもしれないのは、`arabic += numeral.Value`の行で発生するエラーです。このエラーは、`ConvertToArabic`の中で`arabic`を`var arabic = 0`と宣言したために発生します。この宣言自体は正しいのですが、`0`は`int`の値として扱われます。`int`の値と`uint16`の値を足そうとすると失敗します。Goは型付き言語なので、これはうまくいきません。そのため、`var arabic = 0`を`var arabic uint16 = 0`に変更して、変数`arabic`を`uint16`にすることを忘れないでください。

```go
assertion := func(arabic uint16) bool {
	t.Log("testing", arabic)
	roman := ConvertToRoman(arabic)
	fromRoman := ConvertToArabic(roman)
	return fromRoman == arabic
}
```

これで、テストフレームワークの`log`メソッドを使って入力をログ出力していることに注目してください。追加の出力を表示するために、`go test`コマンドを`-v`フラグ付きで実行してください（`go test -v`）。

テストを実行すると、実際に実行され、何がテストされているかを見ることができます。これは、私たちのコードが望むように動作していることに大きな自信を与えてくれます。

しかし、まだ微妙な抜け穴があります。`uint16`は65535まで扱えますが、3999を超える値は有効なローマ数字ではありません。その範囲の値も渡されることになり、私たちの関数がそれを処理できない場合、往復変換は黙って誤った結果を生み出し、私たちはそれを検出できません。`if arabic > 3999 { return true }`のようにして防ぐこともできますが、それは問題を取り繕っているだけで、それらの反復は何もしないことになってしまいます。

より良いアプローチは、私たちが実際に関心のある値だけを生成するよう`quick.Check`に指示することです。これは`quick.Config`の`Values`フィールドを使って行えます。

```go
if err := quick.Check(assertion, &quick.Config{
	MaxCount: 1000,
	Values: func(args []reflect.Value, r *rand.Rand) {
		args[0] = reflect.ValueOf(uint16(r.Intn(4000)))
	},
}); err != nil {
	t.Error("failed checks", err)
}
```

`Values`は、`quick.Check`がアサーションに渡す引数の値のスライスと、乱数のソースを受け取る関数です。`args[0]`を`[0, 3999]`の範囲から取り出した`uint16`で埋めます。これで1000回の実行のすべてが実際に変換を検証するようになり、無駄な反復はなくなります。

### さらなる作業

* 私たちが説明した他のプロパティをチェックするプロパティテストを書くことができますか？
* 誰かが「3999」以上の番号で私たちのコードを呼び出すことができないようにする方法を考えられますか？
  * エラーを返すことができます。
  * または、「3999」を表現できない新しい型を作成します。
    * あなたは何が一番いいと思いますか？

## まとめ

### 反復開発でTDDの練習を増やす

「1984」を`MCMLXXXIV`に変換するコードを書くことは、最初は怖く感じましたか？私は長い間ソフトウェアを書いてきました。

いつものことですが、コツは、**簡単なことから始めて** **小さなステップ**を踏むことです。

このプロセスでは、大きな飛躍をしたり、大きなリファクタリングをしたり、混乱に陥ったりすることはありませんでした。

誰かが「これはただの型だ」と皮肉を言っているのが聞こえてきます。これに異論はありませんが、私は今でも自分が取り組むすべてのプロジェクトで同じアプローチをとっています。私は最初のステップで大きな分散システムを出荷することはありません。チームが出荷できる最も単純なものを見つけて、「Hello world」というウェブサイトを出荷し、その後、管理可能な小さな断片の機能を反復していくのです。

スキルとは、作業を分割する _方法_ を知っていることです。それは練習と、手助けとなる素敵なTDDによって身につきます。

### プロパティベースのテスト

* 標準ライブラリに組み込まれている
* ドメインルールをコードで記述する方法を考えることができれば、自信をつけるための優れたツールになります。
* ドメインについて深く考えさせられる
* テストスイートを補完するものになる可能性があります。

## 追記

この本はコミュニティからの貴重なフィードバックに支えられています。
[Dave](http://github.com/gypsydave5)は、ほぼすべての章で多大な助けとなっています。
しかし彼は、この章で私が「アラビア数字（Arabic numerals）」という言葉を使っていることについて
本気で熱弁をふるいました。そこで、包み隠さずお伝えするために、彼の言葉をここに載せておきます。

> `int`型の値がなぜ本当は「アラビア数字（arabic numeral）」ではないのかを書いておこうと思う。
> これは僕が細かすぎるだけかもしれないので、うるさいと言われても完全に理解するよ。
>
> _数字（digit）_ とは、数を表すのに使われる文字のことだ。ラテン語で「指」を意味する言葉に由来していて、
> 僕たちは普通、指を10本持っている。アラビア（ヒンドゥー・アラビアとも呼ばれる）記数法には10個の数字がある。
> そのアラビア数字とは次のとおりだ。
>
> ```console
>   0 1 2 3 4 5 6 7 8 9
> ```
>
> _数詞（numeral）_ とは、数字の集まりを使って数を表現したものだ。
> アラビア数詞とは、10進の位取り記数法においてアラビア数字で表された数のことだ。
> 「位取り」と言うのは、各数字が数詞の中での位置によって異なる値を持つからだ。だから
>
> ```console
>   1337
> ```
>
> この`1`は、4桁の数詞の最初の数字なので、千という値を持つ。
>
> ローマ数字は、主に数詞を作るための値として、少ない数の数字（`I`、`V`など）を使って作られる。
> 位置に関係するものも少しはあるけど、ほとんどの場合`I`は常に「1」を表す。
>
> では、これを踏まえると、`int`は「アラビア数字」なのだろうか？数という概念は、
> その表現とはまったく結びついていない。次の数の正しい表現は何かと自問してみれば、それがわかる。
>
> ```console
> 255
> 11111111
> two-hundred and fifty-five
> FF
> 377
> ```
>
> そう、これはひっかけ問題だ。どれも正しい。それぞれ10進数、2進数、英語、16進数、8進数の記数法で
> 同じ数を表現したものなんだ。
>
> 数を数詞として表現することは、数としての性質とは _独立_ している。Goの整数リテラルを見ればそれがわかる。
>
> ```go
> 	0xFF == 255 // true
> ```
>
> そして、書式文字列で整数を出力する方法を見てもわかる。
>
> ```go
> n := 255
> fmt.Printf("%b %c %d %o %q %x %X %U", n, n, n, n, n, n, n, n)
> // 11111111 ÿ 255 377 'ÿ' ff FF U+00FF
> ```
>
> 同じ整数を、16進数でもアラビア（10進）数詞でも書くことができる。
>
> だから、関数のシグネチャが`ConvertToRoman(arabic int) string`のようになっていると、
> それがどのように呼び出されるかについて少し思い込みをしていることになる。
> なぜなら、`arabic`は10進の整数リテラルとして書かれることもあれば
>
> ```go
> 	ConvertToRoman(255)
> ```
>
> 同じように次のように書かれることもあるからだ。
>
> ```go
> 	ConvertToRoman(0xFF)
> ```
>
> 実際のところ、僕たちはアラビア数詞から「変換」しているのではまったくなく、`int`をローマ数字として
> 「出力」、つまり表現しているんだ。そして`int`はアラビアであれ何であれ数詞ではなく、単なる数だ。
> `ConvertToRoman`関数は、`int`を`string`に変えるという点で、むしろ`strconv.Itoa`に近い。
>
> でも、このカタの他のどのバージョンもこの区別を気にしていないから、
> :shrug:
