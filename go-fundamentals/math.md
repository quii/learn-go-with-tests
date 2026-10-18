---
description: 数学
---

# 数学

[**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/math)

現代のコンピューターは膨大な計算を驚異的な速さでこなす力を持っていますが、平均的な開発者が仕事で数学を使うことはほとんどありません。

でも今日は違います！

今日は、数学を使って _実際の_ 問題を解決します。退屈な数学ではありません。
三角法やベクトルなど、高校を卒業したらもう使うことはないと言い続けてきたあらゆるものを使います。

## 問題

時計のSVGを作成したい。

デジタル時計ではありません。いや、それでは簡単すぎます。針のある _アナログ_ 時計です。
凝ったものは求めていません。欲しいのは、`time`パッケージから`Time`を受け取り、時、分、秒のすべての針が正しい方向を向いている時計の _SVG_ を出力するだけの素敵な関数です。

これがどれほど難しいことなのでしょうか？

まず、いじって試すための時計の _SVG_ が必要です。
SVGは、XMLで記述された一連の形状として記述されているため、プログラムで操作するための素晴らしい画像フォーマットです。

![時計のSVG](../.gitbook/assets/example_clock.svg)

このように記述されています。

```xml
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="300"
     height="300"
     viewBox="0 0 300 300"
     version="2.0">

  <!-- bezel -->
  <circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>

  <!-- hour hand -->
  <line x1="150" y1="150" x2="114.150000" y2="132.260000"
        style="fill:none;stroke:#000;stroke-width:7px;"/>

  <!-- minute hand -->
  <line x1="150" y1="150" x2="101.290000" y2="99.730000"
        style="fill:none;stroke:#000;stroke-width:7px;"/>

  <!-- second hand -->
  <line x1="150" y1="150" x2="77.190000" y2="202.900000"
        style="fill:none;stroke:#f00;stroke-width:3px;"/>
</svg>
```

これは3本の線が描かれた円で、各線は円の中心（x=150、y=150）から始まり、ある程度離れたところで終わっています。

そこで、上記をどうにかして再構築しつつ、線を変更して、与えられた時刻に応じて適切な方向を指すようにします。

## 受け入れテスト

深入りする前に、受け入れテストについて考えてみましょう。

待って、受け入れテストが何なのかまだ知らないですよね。では、説明してみましょう。

質問させてください。勝利とはどのようなものでしょうか？作業が終わったことはどうすればわかるでしょうか？TDDは、作業が終わったことを知るための良い方法を提供してくれます。テストが成功したときです。ときには（実際にはほとんどの場合）、使用可能な機能全体を書き終えたことを教えてくれるテストを書くと便利です。特定の関数が期待どおりに動作していることを教えてくれるだけのテストではなく、達成しようとしているもの全体、つまり「機能（feature）」が完成したことを教えてくれるテストです。

このようなテストは「受け入れテスト（acceptance tests）」と呼ばれることもあれば、「機能テスト（feature tests）」と呼ばれることもあります。考え方としては、達成しようとしていることを説明する本当に高いレベルのテストを書くということです。たとえば、ユーザーがWebサイトのボタンをクリックすると、捕まえたポケモンの完全なリストが表示される、といったものです。そのテストを書いたら、受け入れテストに合格する動作するシステムに向けて、さらにテスト（ユニットテスト）を書いていくことができます。この例でいえば、それらのテストは、ボタンのあるWebページのレンダリング、Webサーバーのルートハンドラーのテスト、データベースの検索の実行などに関するものになるでしょう。これらはすべてTDDで進められ、そのすべてが元の受け入れテストを成功させることにつながっていきます。

Nat PryceとSteve Freemanによる、この _古典的な_ 図のようなものです。

![TDDにおけるアウトサイドインのフィードバックループ](../.gitbook/assets/TDD-outside-in.jpg)

とにかく、受け入れテスト、つまり完了したときにそれを教えてくれるテストを書いてみましょう。

時計の例があるので、重要なパラメータがどうなるかを考えてみましょう。

```text
<line x1="150" y1="150" x2="114.150000" y2="132.260000"
        style="fill:none;stroke:#000;stroke-width:7px;"/>
```

時計の中心（この線の属性`x1`および`y1`）は、時計の各針で同じです。
時計の針ごとに変更する必要がある数値（SVGを構築するためのパラメータ）は、`x2`および`y2`属性です。時計の針ごとに`X`と`Y`が必要です。

パラメータはもっと考えることも _できます_ 。文字盤の円の半径、SVGのサイズ、針の色や形など...。しかし、まずはシンプルで具体的な問題をシンプルで具体的な解決策で解き、それから一般化のためにパラメータを追加していくほうが良いでしょう。

そこで、次のように決めます。

* すべての時計の中心は（150、150）です
* 時針の長さは50
* 分針の長さは80
* 秒針の長さは90

SVGに関する注意点：原点（点(0, 0)）は、予想されるような _左下_ ではなく、 _左上_ にあります。線にどんな数値を当てはめるかを考える際には、これを覚えておくことが重要です。

最後に、SVGを _どのように_ 構築するかはまだ決めていません。
[`text/template`](https://golang.org/pkg/text/template/)パッケージのテンプレートを使うこともできますし、単にバイトを`bytes.Buffer`やライターに送り込むこともできます。しかし、これらの数値が必要になることはわかっているので、それらを作成するもののテストに焦点を合わせましょう。

### 最初にテストを書く

だから私の最初のテストは次のようになります。

```go
package clockface_test

import (
	"projectpath/clockface"
	"testing"
	"time"
)

func TestSecondHandAtMidnight(t *testing.T) {
	tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

	want := clockface.Point{X: 150, Y: 150 - 90}
	got := clockface.SecondHand(tm)

	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}
```

ここでの`projectpath`はプレースホルダーです。自分のモジュールのパスに`/clockface`を加えたものに置き換えてください（たとえば、`go.mod`に`module example.com/learning-go`と書かれていれば、このインポートは`"example.com/learning-go/clockface"`になります）。これが機能するのは、`clockface_test.go`を`clockface`という名前の専用ディレクトリに、（次に作成する）`clockface.go`と並べて置いているからです。ディレクトリ名をその中のパッケージ名と一致させることで、Goはエイリアスなしでインポートを解決できます。同じディレクトリ内の`package clockface_test`は、Goが認めている特別なケースです。同じフォルダ内に`package clockface`のファイルと並んで存在できる外部テストパッケージです。

SVGでは座標を左上の角から取ることを覚えていますか？
秒針を真夜中（0時0分0秒）の位置に置くには、X軸上では文字盤の中心から動いておらず（150のまま）、Y軸上では中心から針の長さだけ「上」、つまり150マイナス90になると期待します。

### テストを実行してみます

これにより、不足している関数と型に関する、予想どおりの失敗が引き出されます。

```text
--- FAIL: TestSecondHandAtMidnight (0.00s)
./clockface_test.go:13:10: undefined: clockface.Point
./clockface_test.go:14:9: undefined: clockface.SecondHand
```

つまり、秒針の先端が来るべき位置を表す`Point`と、それを取得する関数が必要です。

### テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

これらの型を実装して、コードをコンパイルしてみましょう

```go
package clockface

import "time"

// A Point represents a two-dimensional Cartesian coordinate
type Point struct {
	X float64
	Y float64
}

// SecondHand is the unit vector of the second hand of an analogue clock at time `t`
// represented as a Point.
func SecondHand(t time.Time) Point {
	return Point{}
}
```

すると、次のようになります。

```text
--- FAIL: TestSecondHandAtMidnight (0.00s)
    clockface_test.go:17: Got {0 0}, wanted {150 60}
FAIL
exit status 1
FAIL	learn-go-with-tests/math/clockface	0.006s
```

### 成功させるのに十分なコードを書く

予想どおりの失敗が得られたら、`SecondHand`の戻り値を埋めることができます。

```go
// SecondHand is the unit vector of the second hand of an analogue clock at time `t`
// represented as a Point.
func SecondHand(t time.Time) Point {
    return Point{150, 60}
}
```

見てください。テストが成功しました。

```text
PASS
ok  	    clockface	0.006s
```

### リファクタリング♪

まだリファクタリングする必要はありません。コードはほとんどありません！

### 新しい要件について繰り返します

どんな時刻でも真夜中を示す時計を返すだけではなく、ここで何かしら作業をする必要がありそうです...

### 最初にテストを書く

```go
func TestSecondHandAt30Seconds(t *testing.T) {
    tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)

    want := clockface.Point{X: 150, Y: 150 + 90}
    got := clockface.SecondHand(tm)

    if got != want {
        t.Errorf("Got %v, wanted %v", got, want)
    }
}
```

考え方は同じですが、今度は秒針が _下_ を向いているので、Y軸に針の長さを _足し_ ます。

これはコンパイルできます...しかし、どうすれば成功させられるでしょうか？

## 思考時間

この問題をどのように解決しますか？

秒針は毎分同じ60の状態を通過し、60の異なる方向を指しています。 0秒の場合は文字盤の上部を指し、30秒の場合は文字盤の下部を指します。簡単です。

つまり、たとえば37秒のときに秒針がどの方向を向いているかを考えたい場合、12時の位置から円周を37/60だけ回ったところまでの角度が必要です。
度数法では `(360 / 60 ) * 37 = 222`度ですが、1回転の`37/60`だと覚えておくほうが簡単です。

しかし、角度だけでは話の半分にすぎません。秒針の先端が指しているX座標とY座標を知る必要があります。どうすれば求められるでしょうか？

## 数学

原点（座標`0, 0`）を中心に描かれた半径1の円を想像してください。

![単位円の画像](../.gitbook/assets/unit_circle.png)

これは「単位円（`unit circle`）」と呼ばれます。なぜなら...そう、半径が1単位だからです！

円周は、グリッド上の点、つまり座標の集まりでできています。
これらの各座標のx成分とy成分は三角形を形成し、その斜辺は常に1（つまり円の半径）です。

![円周上にポイントが定義されている単位円の画像](../.gitbook/assets/unit_circle_coords.png)

三角法を使えば、原点でなす角度がわかっている場合に、各三角形のXとYの長さを求めることができます。X座標はcos(a)、Y座標はsin(a)になります。ここで、aは線と（正の）x軸とのなす角度です。

![それぞれcos（a）とsin（a）として定義された光線のx要素とy要素を含む単位円の画像。ここで、aはx軸と光線がなす角度](<../.gitbook/assets/unit_circle_params (1).png>)

（これを信じない場合は、[ウィキペディアを見て...](https://en.wikipedia.org/wiki/Sine#Unit_circle_definition)）

最後にひとひねり。X軸（3時）からではなく12時からの角度を測りたいので、軸を入れ替える必要があります。つまり、x = sin(a)、y = cos(a)となります。

![y軸からの角度で定義された単位円光線](../.gitbook/assets/unit_circle_12_oclock.png)

これで、秒針の角度（1秒ごとに円の1/60）とX座標とY座標を取得する方法がわかりました。`sin`と`cos`の両方の関数が必要です。

## `math`

幸い、Goの`math`パッケージには両方ありますが、1つだけ理解しておかなければならない小さな落とし穴があります。
[`math.Cos`](https://golang.org/pkg/math/#Cos)の説明を見てみましょう。

> Cosは、ラジアン引数xの余弦を返します。

角度はラジアンで与える必要があります。では、ラジアンとは何ですか？円の完全な回転を360度で構成するのではなく、完全な回転を2πラジアンとして定義します。こうするのには十分な理由がありますが、ここでは立ち入りません。

読んだり、学習したり、考えたりしたので、次のテストを書くことができます。

### 最初にテストを書く

こうした数学は難しくて混乱します。何が起きているのかを理解できている自信がありません。
それではテストを書きましょう！問題全体を一度に解決する必要はありません。特定の時間の秒針について、ラジアン単位で正しい角度を計算することから始めましょう。

これらのテストに取り組んでいる間は、取り組んでいた受け入れテストを _コメントアウト_ しておきます。こちらのテストを成功させている間に、そのテストに気を取られたくないからです。

### パッケージのおさらい

現時点では、受け入れテストは`clockface_test`パッケージにあります。テストは`clockface`パッケージの外にあっても構いません。名前が`_test.go`で終わっていれば実行できます。

これらのラジアンのテストは`clockface`パッケージ _内_ で記述します。それらがエクスポートされることはないかもしれませんし、何が起こっているのかをしっかり把握したら削除される（または移動される）かもしれません。受け入れテストのファイル名を`clockface_acceptance_test.go`に変更し、秒をラジアンでテストするための _新しい_ ファイル`clockface_test`を作成できるようにします。

```go
package clockface

import (
    "math"
    "testing"
    "time"
)

func TestSecondsInRadians(t *testing.T) {
    thirtySeconds := time.Date(312, time.October, 28, 0, 0, 30, 0, time.UTC)
    want := math.Pi
    got := secondsInRadians(thirtySeconds)

    if want != got {
        t.Fatalf("Wanted %v radians, but got %v", want, got)
    }
}
```

ここでは、毎分30秒のときに秒針が時計をちょうど半周したところにあることをテストしています。そして、これが`math`パッケージの初めての利用です！円の1回転が2πラジアンなら、半周はちょうどπラジアンになるはずです。`math.Pi`はπの値を提供します。

### テストを実行してみます

```text
./clockface_test.go:12:9: undefined: secondsInRadians
```

### テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func secondsInRadians(t time.Time) float64 {
    return 0
}
```

```text
clockface_test.go:15: Wanted 3.141592653589793 radians, but got 0
```

### 成功させるのに十分なコードを書く

```go
func secondsInRadians(t time.Time) float64 {
    return math.Pi
}
```

```text
PASS
ok  	clockface	0.011s
```

### リファクタリング♪

まだリファクタリングは必要ありません。

### 新しい要件について繰り返します

テストを拡張して、さらにいくつかのシナリオをカバーできます。
少し先に進んで、すでにリファクタリング済みのテストコードを示します。どうやってここにたどり着いたかは十分明らかなはずです。

```go
func TestSecondsInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(0, 0, 30), math.Pi},
        {simpleTime(0, 0, 0), 0},
        {simpleTime(0, 0, 45), (math.Pi / 2) * 3},
        {simpleTime(0, 0, 7), (math.Pi / 30) * 7},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := secondsInRadians(c.time)
            if got != c.angle {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

このテーブルベースのテストを書く手間を少し減らすために、いくつかのヘルパー関数を追加しました。`testName`は時刻をデジタル時計形式（HH:MM:SS）に変換し、`simpleTime`は実際に関心のある部分（繰り返しますが、時、分、秒）だけを使って`time.Time`を構築します。これがそのコードです。

```go
func simpleTime(hours, minutes, seconds int) time.Time {
    return time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
}

func testName(t time.Time) string {
    return t.Format("15:04:05")
}
```

これらの2つの関数は、これらのテスト（および将来のテスト）の記述と保守を少し簡単にするのに役立ちます。

これにより、優れたテスト出力が得られます。

```text
clockface_test.go:24: Wanted 0 radians, but got 3.141592653589793

clockface_test.go:24: Wanted 4.71238898038469 radians, but got 3.141592653589793
```

上で話していた数学をすべて実装する時間です。

```go
func secondsInRadians(t time.Time) float64 {
    return float64(t.Second()) * (math.Pi / 30)
}
```

1秒は（2π/60）ラジアンです... 2を約分すると、π/30ラジアンになります。これに秒数（ `float64`として）を掛けると、すべてのテストに合格するはずです...

```text
clockface_test.go:24: Wanted 3.141592653589793 radians, but got 3.1415926535897936
```

待って、なになに？

### 浮動小数点数（`Floats`）は恐ろしい

浮動小数点演算は[悪名高いほど不正確](https://0.30000000000000004.com/)です。
コンピュータが本当に扱えるのは整数と、ある程度までの有理数だけです。小数は、特に`secondsInRadians`関数で行っているように割ったり掛けたりすると、不正確になり始めます。
`math.Pi`を30で割ってから30を掛けたことで、 _もはや`math.Pi`と同じではない数値_ になってしまったのです。

これを回避するには2つの方法があります。

1. それを受け入れる
2. 方程式をリファクタリングして関数をリファクタリングする

（1）はそれほど魅力的ではないように思われるかもしれませんが、浮動小数点の等価性を機能させる唯一の方法であることがよくあります。時計の文字盤を描画する目的では、ごくわずかな誤差があっても率直に言って問題にならないので、角度に対して「十分に近い」等式を定義する関数を書くことができます。

しかし、精度を取り戻す簡単な方法があります。割ってから掛けるということをしないように、式を変形するのです。割り算だけですべて済ませられます。

つまり、次の式の代わりに

```text
numberOfSeconds * π / 30
```

と書くことができます。

```text
π / (30 / numberOfSeconds)
```

これは同等です。

Goでは次のようになります。

```go
func secondsInRadians(t time.Time) float64 {
    return (math.Pi / (30 / (float64(t.Second()))))
}
```

これでテストが成功します。

```text
PASS
ok      clockface     0.005s
```

すべて[このようになっている](https://github.com/quii/learn-go-with-tests/tree/main/math/v3/clockface)はずです。

### ゼロ除算についての注意

コンピューターは、無限大が少し奇妙なものであるため、ゼロで割ることを好まないことがよくあります。

Goでは、明示的にゼロで割ろうとするとコンパイルエラーになります。

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println(10.0 / 0.0) // fails to compile
}
```

もちろん、私たちの`t.Second()`のように、ゼロで割ることになるかどうかをコンパイラが常に予測できるわけではありません。

これを試してみてください。

```go
func main() {
	fmt.Println(10.0 / zero())
}

func zero() float64 {
	return 0.0
}
```

これは`+Inf`（無限大）を出力します。+Infで割るとゼロになるようです。次のコードで確認できます。

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(secondsinradians())
}

func zero() float64 {
	return 0.0
}

func secondsinradians() float64 {
	return (math.Pi / (30 / (float64(zero()))))
}
```

### 新しい要件について繰り返します

したがって、ここで最初の部分をカバーしました。
秒針がラジアンで指し示す角度はわかっています。次に、座標を計算する必要があります。

ここでも、できるだけシンプルに保ち、 _単位円_ 、つまり半径1の円だけを扱うことにします。これは針の長さがすべて1になることを意味しますが、良い面としては、数学がずっと飲み込みやすくなります。

### 最初にテストを書く

```go
func TestSecondHandPoint(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 30), Point{0, -1}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if got != c.point {
				t.Fatalf("Wanted %v Point, but got %v", c.point, got)
			}
		})
	}
}
```

### テストを実行してみます

```text
./clockface_test.go:40:11: undefined: secondHandPoint
```

### テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func secondHandPoint(t time.Time) Point {
    return Point{}
}
```

```text
clockface_test.go:42: Wanted {0 -1} Point, but got {0 0}
```

### 成功させるのに十分なコードを書く

```go
func secondHandPoint(t time.Time) Point {
    return Point{0, -1}
}
```

```text
PASS
ok  	clockface	0.007s
```

### 新しい要件について繰り返します

```go
func TestSecondHandPoint(t *testing.T) {
    cases := []struct {
        time  time.Time
        point Point
    }{
        {simpleTime(0, 0, 30), Point{0, -1}},
        {simpleTime(0, 0, 45), Point{-1, 0}},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := secondHandPoint(c.time)
            if got != c.point {
                t.Fatalf("Wanted %v Point, but got %v", c.point, got)
            }
        })
    }
}
```

### テストを実行してみます

```text
clockface_test.go:43: Wanted {-1 0} Point, but got {0 -1}
```

### 成功させるのに十分なコードを書く

単位円の画像を覚えていますか？

![それぞれcos\(a\)とsin\(a\)として定義された光線のx要素とy要素を持つ単位円の画像。ここで、aはx軸と光線がなす角度](<../.gitbook/assets/unit_circle_params (1).png>)

また、秒針と3時の間の角度を測るX軸からではなく、Y軸である12時から角度を測りたいことも思い出してください。

![y軸からの角度で定義された単位円光線](../.gitbook/assets/unit_circle_12_oclock.png)

ここで、XとYを求める式が必要です。それを秒針の関数に書いてみましょう。

```go
func secondHandPoint(t time.Time) Point {
    angle := secondsInRadians(t)
    x := math.Sin(angle)
    y := math.Cos(angle)

    return Point{x, y}
}
```

すると、次のようになります。

```text
clockface_test.go:43: Wanted {0 -1} Point, but got {1.2246467991473515e-16 -1}

clockface_test.go:43: Wanted {-1 0} Point, but got {-1 -1.8369701987210272e-16}
```

待って、（また）何だって？どうやら再び浮動小数点数に呪われたようです。予期しない数値はどちらも _ごくごく小さな_ 値で、小数点以下16桁目あたりの話です。
ここでも、精度を上げようとするか、だいたい等しいということにして先へ進むかを選べます。

これらの角度の精度を上げる選択肢の1つは、`math/big`パッケージの有理数型`Rat`を使うことです。
しかし、目的は月面着陸ではなくSVGを描くことなので、多少のあいまいさは許容できると思います。

```go
func TestSecondHandPoint(t *testing.T) {
    cases := []struct {
        time  time.Time
        point Point
    }{
        {simpleTime(0, 0, 30), Point{0, -1}},
        {simpleTime(0, 0, 45), Point{-1, 0}},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := secondHandPoint(c.time)
            if !roughlyEqualPoint(got, c.point) {
                t.Fatalf("Wanted %v Point, but got %v", c.point, got)
            }
        })
    }
}

func roughlyEqualFloat64(a, b float64) bool {
    const equalityThreshold = 1e-7
    return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b Point) bool {
    return roughlyEqualFloat64(a.X, b.X) &&
        roughlyEqualFloat64(a.Y, b.Y)
}
```

2つの`Points`がおおよそ等しいことを判定する2つの関数を定義しました。X要素とY要素の差がそれぞれ0.0000001以内であれば等しいとみなします。それでも十分に正確です。

そして今、私たちはテストに成功します。

```text
PASS
ok  	clockface	0.007s
```

### リファクタリング♪

これでもかなり満足しています。

[現在はこのようになっています](https://github.com/quii/learn-go-with-tests/tree/main/math/v4/clockface)。

### 新しい要件について繰り返します

まあ、 _新しい_ というのは正確ではありません。実際には、ここであの受け入れテストを成功させることができるのです！どんなテストだったか思い出してみましょう。

```go
func TestSecondHandAt30Seconds(t *testing.T) {
    tm := time.Date(1337, time.January, 1, 0, 0, 30, 0, time.UTC)

    want := clockface.Point{X: 150, Y: 150 + 90}
    got := clockface.SecondHand(tm)

    if got != want {
        t.Errorf("Got %v, wanted %v", got, want)
    }
}
```

### テストを実行してみます

```text
clockface_acceptance_test.go:28: Got {150 60}, wanted {150 240}
```

### 成功させるのに十分なコードを書く

単位ベクトルをSVG上の点に変換するには、3つのことを行う必要があります。

1. 針の長さに合わせて拡大縮小します
2. SVGの原点が左上隅にあることを考慮して、X軸で反転します
3. 正しい位置に平行移動します（原点（150, 150）から伸びるようにします）

楽しい時間ですね！

```go
// SecondHand is the unit vector of the second hand of an analogue clock at time `t`
// represented as a Point.
func SecondHand(t time.Time) Point {
    p := secondHandPoint(t)
    p = Point{p.X * 90, p.Y * 90}   // scale
    p = Point{p.X, -p.Y}            // flip
    p = Point{p.X + 150, p.Y + 150} // translate
    return p
}
```

正確にその順序でスケーリング、反転、平行移動します。数学ばんざい！

```text
PASS
ok  	clockface	0.007s
```

### リファクタリング♪

定数として切り出すべきマジックナンバーがいくつかあるので、そうしましょう。

```go
const secondHandLength = 90
const clockCentreX = 150
const clockCentreY = 150

// SecondHand is the unit vector of the second hand of an analogue clock at time `t`
// represented as a Point.
func SecondHand(t time.Time) Point {
    p := secondHandPoint(t)
    p = Point{p.X * secondHandLength, p.Y * secondHandLength}
    p = Point{p.X, -p.Y}
    p = Point{p.X + clockCentreX, p.Y + clockCentreY} //translate
    return p
}
```

## 時計を描く

まあ...少なくとも秒針だけは...

さあ、やってみましょう。世に出て人々を驚かせるのを今か今かと待っている価値を届けないことほど、悪いことはありませんから。秒針を描いてみましょう！

メインの`clockface`パッケージディレクトリの下に、（紛らわしいことに）`clockface`という名前の新しいディレクトリを追加します。そこにSVGをビルドするバイナリを作成する`main`パッケージを置きます。

```text
|-- clockface
|       |-- main.go
|-- clockface.go
|-- clockface_acceptance_test.go
|-- clockface_test.go
```

`main.go`の中では、次のコードから始めますが、clockfaceパッケージのインポートは自分のバージョンを指すように変更してください。

```go
package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"learn-go-with-tests/math/clockface" // REPLACE THIS!
)

func main() {
	t := time.Now()
	sh := clockface.SecondHand(t)
	io.WriteString(os.Stdout, svgStart)
	io.WriteString(os.Stdout, bezel)
	io.WriteString(os.Stdout, secondHandTag(sh))
	io.WriteString(os.Stdout, svgEnd)
}

func secondHandTag(p clockface.Point) string {
	return fmt.Sprintf(`<line x1="150" y1="150" x2="%f" y2="%f" style="fill:none;stroke:#f00;stroke-width:3px;"/>`, p.X, p.Y)
}

const svgStart = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="100%"
     height="100%"
     viewBox="0 0 300 300"
     version="2.0">`

const bezel = `<circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>`

const svgEnd = `</svg>`
```

いやはや、 _この_ ごちゃごちゃしたコードで美しいコードの賞を狙っているわけではありません。
しかし、役目は果たします。これは`os.Stdout`に SVG を書き出しています。一度に1つの文字列ずつ。

これをビルドし

```text
go build
```

実行して、出力をファイルに書き出すと

```text
./clockface > clock.svg
```

次のようなものが表示されるはずです。

![秒針のみの時計](../.gitbook/assets/clock.svg)

そして、これが[コードの様子](https://github.com/quii/learn-go-with-tests/tree/main/math/v6/clockface)です。

### リファクタリング♪

これは臭います。まあ、 _ひどく_ 臭うというほどではありませんが、満足はしていません。

1. `SecondHand`関数全体が、SVGであることと _ものすごく_ 密接に結びついています...SVGに言及してもいないし、実際にSVGを生成してもいないのに...

2. ... 同時に、私はSVGコードをテストしていません。

ええ、しくじったようです。これは間違っている気がします。
もっとSVGを中心に据えたテストで立て直してみましょう。

どんな選択肢があるでしょうか？
そうですね、`SVGWriter`から吐き出される文字に、特定の時間に予想されるSVGタグのようなものが含まれていることをテストしてみることができます。例えば：

```go
func TestSVGWriterAtMidnight(t *testing.T) {
    tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

    var b strings.Builder
    clockface.SVGWriter(&b, tm)
    got := b.String()

    want := `<line x1="150" y1="150" x2="150" y2="60"`

    if !strings.Contains(got, want) {
        t.Errorf("Expected to find the second hand %v, in the SVG output %v", want, got)
    }
}
```

しかし、これは本当に改善されたのでしょうか？

これでは、有効なSVGを生成しなかった場合でも（出力に文字列が含まれることをテストしているだけなので）成功してしまうだけでなく、その文字列にほんのささいな、重要でない変更を加えただけでも失敗してしまいます。たとえば、属性の間にスペースを1つ追加しただけでもです。

_最大の_ 臭いは、データ構造（XML）を、一連の文字（文字列）としての表現を見てテストしていることです。これは _決して_ 、 _絶対に_ 良い考えではありません。上で挙げたような問題、つまり、壊れやすすぎるうえに感度も足りないテストを生み出してしまうからです。間違ったものをテストしているテストなのです！

したがって、唯一の解決策は、出力を _XMLとして_ テストすることです。そのためには、それを解析する必要があります。

## XMLの解析

[`encoding/xml`](https://pkg.go.dev/encoding/xml)は、シンプルなXML解析で行うすべてのことを処理できるGoパッケージです。

関数[`xml.Unmarshal`](https://pkg.go.dev/encoding/xml#Unmarshal)は、XMLデータの`[]byte`と、アンマーシャル先の構造体へのポインタを受け取ります。

したがって、XMLをアンマーシャルするための構造体が必要です。
すべてのノードと属性の正しい名前と正しい構造の書き方を検討するのに少し時間を費やすこともできますが、幸い、誰かがそのハードワークのすべてを自動化してくれるプログラム[`zek`](https://github.com/miku/zek)を書いてくれています。
さらに良いことに、[https://xml-to-go.github.io/](https://xml-to-go.github.io/)にオンラインバージョンがあります。ファイルの上部にあるSVGを片方のボックスに貼り付けるだけで、ほら、次のものが出てきます。

```go
type Svg struct {
    XMLName xml.Name `xml:"svg"`
    Text    string   `xml:",chardata"`
    Xmlns   string   `xml:"xmlns,attr"`
    Width   string   `xml:"width,attr"`
    Height  string   `xml:"height,attr"`
    ViewBox string   `xml:"viewBox,attr"`
    Version string   `xml:"version,attr"`
    Circle  struct {
        Text  string `xml:",chardata"`
        Cx    string `xml:"cx,attr"`
        Cy    string `xml:"cy,attr"`
        R     string `xml:"r,attr"`
        Style string `xml:"style,attr"`
    } `xml:"circle"`
    Line []struct {
        Text  string `xml:",chardata"`
        X1    string `xml:"x1,attr"`
        Y1    string `xml:"y1,attr"`
        X2    string `xml:"x2,attr"`
        Y2    string `xml:"y2,attr"`
        Style string `xml:"style,attr"`
    } `xml:"line"`
}
```

（構造体の名前を`SVG`に変更するなど）必要に応じてこれを調整することもできますが、手始めとしては十分です。この構造体を`clockface_acceptance_test`ファイルに貼り付けて、これを使ったテストを書いてみましょう。

```go
func TestSVGWriterAtMidnight(t *testing.T) {
    tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)

    b := bytes.Buffer{}
    clockface.SVGWriter(&b, tm)

    svg := Svg{}
    xml.Unmarshal(b.Bytes(), &svg)

    x2 := "150"
    y2 := "60"

    for _, line := range svg.Line {
        if line.X2 == x2 && line.Y2 == y2 {
            return
        }
    }

    t.Errorf("Expected to find the second hand with x2 of %+v and y2 of %+v, in the SVG output %v", x2, y2, b.String())
}
```

`clockface.SVGWriter`の出力を`bytes.Buffer`に書き込み、次にそれを`Svg`に`Unmarshal`します。次に、`Svg`内の各`Line`を見て、それらのいずれかが期待される`X2`および`Y2`値を持っているかどうかを確認します。
一致するものがあれば早期にreturnし（テストは成功）、なければ（できれば）わかりやすいメッセージとともに失敗させます。

```sh
./clockface_acceptance_test.go:41:2: undefined: clockface.SVGWriter
```

`SVGWriter.go`を作成したほうがよさそうです...

```go
package clockface

import (
    "fmt"
    "io"
    "time"
)

const (
    secondHandLength = 90
    clockCentreX     = 150
    clockCentreY     = 150
)

//SVGWriter writes an SVG representation of an analogue clock, showing the time t, to the writer w
func SVGWriter(w io.Writer, t time.Time) {
    io.WriteString(w, svgStart)
    io.WriteString(w, bezel)
    secondHand(w, t)
    io.WriteString(w, svgEnd)
}

func secondHand(w io.Writer, t time.Time) {
    p := secondHandPoint(t)
    p = Point{p.X * secondHandLength, p.Y * secondHandLength} // scale
    p = Point{p.X, -p.Y}                                      // flip
    p = Point{p.X + clockCentreX, p.Y + clockCentreY}         // translate
    fmt.Fprintf(w, `<line x1="150" y1="150" x2="%f" y2="%f" style="fill:none;stroke:#f00;stroke-width:3px;"/>`, p.X, p.Y)
}

const svgStart = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="100%"
     height="100%"
     viewBox="0 0 300 300"
     version="2.0">`

const bezel = `<circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>`

const svgEnd = `</svg>`
```

最も美しいSVGライターですか？いいえ。でも、役目は果たしてくれるはずです...

```text
clockface_acceptance_test.go:56: Expected to find the second hand with x2 of 150 and y2 of 60, in the SVG output <?xml version="1.0" encoding="UTF-8" standalone="no"?>
    <!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
    <svg xmlns="http://www.w3.org/2000/svg"
         width="100%"
         height="100%"
         viewBox="0 0 300 300"
         version="2.0"><circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/><line x1="150" y1="150" x2="150.000000" y2="60.000000" style="fill:none;stroke:#f00;stroke-width:3px;"/></svg>
```

おっと！

`%f`フォーマットディレクティブは、座標をデフォルトの精度レベル（小数点以下6桁）に出力します。座標に期待する精度のレベルを明示する必要があります。小数点第3位としましょう。

```go
	fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#f00;stroke-width:3px;"/>`, p.X, p.Y)
```

テストでの期待を更新した後

```go
    x2 := "150.000"
    y2 := "60.000"
```

結果は次のとおりです。

```text
PASS
ok  	clockface	0.006s
```

これで`main`関数を短くできます。

```go
package main

import (
	"os"
	"time"

	"learn-go-with-tests/math/clockface"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
```

[現在はこのようになっている](https://github.com/quii/learn-go-with-tests/tree/main/math/v7b/clockface)はずです。

そして、同じパターンに従って別の時間のテストを書くことができますが、その前に...

### リファクタリング♪

3つの点が気になります。

1. 存在を確認する必要がある情報のすべてを実際にテストしているわけではありません。たとえば、`x1`の値はどうでしょうか？
2. また、`x1`などの属性は実際には`strings`ではありませんよね？数値です！
3. 針の`style`を本当に気にしていますか？あるいは、`zak`によって生成された空の`Text`ノードはどうでしょうか？

私たちはもっとうまくやることができます。
すべてをシャープにするために、`Svg`構造体とテストにいくつかの調整を加えましょう。

```go
type SVG struct {
    XMLName xml.Name `xml:"svg"`
    Xmlns   string   `xml:"xmlns,attr"`
    Width   string   `xml:"width,attr"`
    Height  string   `xml:"height,attr"`
    ViewBox string   `xml:"viewBox,attr"`
    Version string   `xml:"version,attr"`
    Circle  Circle   `xml:"circle"`
    Line    []Line   `xml:"line"`
}

type Circle struct {
    Cx float64 `xml:"cx,attr"`
    Cy float64 `xml:"cy,attr"`
    R  float64 `xml:"r,attr"`
}

type Line struct {
    X1 float64 `xml:"x1,attr"`
    Y1 float64 `xml:"y1,attr"`
    X2 float64 `xml:"x2,attr"`
    Y2 float64 `xml:"y2,attr"`
}
```

ここで私は

* 構造体の重要な部分を名前付きの型にしました -- `Line` と `Circle`

* 数値属性を`string`ではなく`float64`に変更しました。
* `Style`や`Text`などの未使用の属性を削除しました。
* `Svg`を`SVG`に名前を変更しました。 _それが正しいことだからです_ 。

これにより、探している線について、より正確にアサーションできます。

```go
func TestSVGWriterAtMidnight(t *testing.T) {
    tm := time.Date(1337, time.January, 1, 0, 0, 0, 0, time.UTC)
    b := bytes.Buffer{}

    clockface.SVGWriter(&b, tm)

    svg := SVG{}

    xml.Unmarshal(b.Bytes(), &svg)

    want := Line{150, 150, 150, 60}

    for _, line := range svg.Line {
        if line == want {
            return
        }
    }

    t.Errorf("Expected to find the second hand line %+v, in the SVG lines %+v", want, svg.Line)
}
```

最後に、単体テストのテーブルを参考にして、ヘルパー関数`containsLine(line Line, lines []Line) bool`を記述して、これらのテストを本当に輝かせることができます。`simpleTime`と`testName`も再び使いますが、このファイル（`clockface_acceptance_test.go`）は`package clockface_test`にあり、それらを最初に書いた`clockface_test.go`ファイル（`package clockface`）とは別のパッケージです。これらはエクスポートされていないため、宣言されたパッケージの外からは見えません。そのため、ここに独自のコピーが必要です。

```go
func TestSVGWriterSecondHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{
			simpleTime(0, 0, 0),
			Line{150, 150, 150, 60},
		},
		{
			simpleTime(0, 0, 30),
			Line{150, 150, 150, 240},
		},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			clockface.SVGWriter(&b, c.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(c.line, svg.Line) {
				t.Errorf("Expected to find the second hand line %+v, in the SVG lines %+v", c.line, svg.Line)
			}
		})
	}
}

func containsLine(l Line, ls []Line) bool {
	for _, line := range ls {
		if line == l {
			return true
		}
	}
	return false
}

func simpleTime(hours, minutes, seconds int) time.Time {
	return time.Date(312, time.October, 28, hours, minutes, seconds, 0, time.UTC)
}

func testName(t time.Time) string {
	return t.Format("15:04:05")
}
```

[このようになります](https://github.com/quii/learn-go-with-tests/tree/main/math/v7c/clockface)。

さて、 _これこそ_ 私が受け入れテストと呼ぶものです！

### 最初にテストを書く

これで秒針は完成です。それでは分針に取りかかりましょう。

```go
func TestSVGWriterMinuteHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{
			simpleTime(0, 0, 0),
			Line{150, 150, 150, 70},
		},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			clockface.SVGWriter(&b, c.time)

			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(c.line, svg.Line) {
				t.Errorf("Expected to find the minute hand line %+v, in the SVG lines %+v", c.line, svg.Line)
			}
		})
	}
}
```

### テストを実行してみます

```text
clockface_acceptance_test.go:87: Expected to find the minute hand line {X1:150 Y1:150 X2:150 Y2:70}, in the SVG lines [{X1:150 Y1:150 X2:150 Y2:60}]
```

そろそろ他の針も作り始めたほうがよいでしょう。
秒針のテストを作成したのと同じように、反復して次の一連のテストを作成できます。ここでも、これが動くようになるまで受け入れテストはコメントアウトしておきます。

```go
func TestMinutesInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(0, 30, 0), math.Pi},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := minutesInRadians(c.time)
            if got != c.angle {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

### テストを実行してみます

```text
./clockface_test.go:59:11: undefined: minutesInRadians
```

### テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func minutesInRadians(t time.Time) float64 {
    return math.Pi
}
```

### 新しい要件について繰り返します

さて、では _本格的な_ 作業に取りかかりましょう。
分針を、1分ごとにしか動かないようにモデル化することもできます。
つまり、30分から31分へ、間を動くことなく「ジャンプ」するわけです。しかし、それでは少し見栄えが悪いでしょう。私たちが望むのは、分針が毎秒 _ほんの少しずつ_ 動くことです。

```go
func TestMinutesInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(0, 30, 0), math.Pi},
        {simpleTime(0, 0, 7), 7 * (math.Pi / (30 * 60))},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := minutesInRadians(c.time)
            if got != c.angle {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

その「ほんの少し」とはどのくらいでしょうか？そうですね...

* 1分は60秒
* 円の半回転（`math.Pi`ラジアン）は30分
* つまり半回転は`30 * 60`秒
* したがって、時刻が正時から7秒過ぎであれば...
* ...分針は12の位置から`7 * (math.Pi / (30 * 60))`ラジアン進んだところにあるはずです

### テストを実行してみます

```
clockface_test.go:62: Wanted 0.012217304763960306 radians, but got 3.141592653589793
```

### 成功させるのに十分なコードを書く

ジェニファー・アニストンの不朽の名言を借りれば、[ここからは科学の時間です](https://www.youtube.com/watch?v=29Im23SPNok)。

```go
func minutesInRadians(t time.Time) float64 {
    return (secondsInRadians(t) / 60) +
        (math.Pi / (30 / float64(t.Minute())))
}
```

1秒ごとに分針を文字盤上でどれだけ進めるかを一から計算する代わりに、ここでは`secondsInRadians`関数を活用できます。
1秒ごとに、分針が秒針の角度の1/60ずつ移動します。

```go
secondsInRadians(t) / 60
```

あとは、秒針の動きと同じように、分の動きを足すだけです。

```go
math.Pi / (30 / float64(t.Minute()))
```

そして...

```
PASS
ok  	clockface	0.007s
```

簡単です。[現在はこのようになっています](https://github.com/quii/learn-go-with-tests/tree/main/math/v8/clockface/clockface_acceptance_test.go)。

### 新しい要件について繰り返します

`minutesInRadians`テストにさらにケースを追加する必要がありますか？現時点では2つしかありません。`minuteHandPoint`関数のテストに進む前に、いくつのケースが必要ですか？

私のお気に入りのTDDの引用の1つは、しばしばケント・ベックの言葉とされている次のものです。

> 恐怖が退屈に変わるまでテストを書いてください。

そして率直に言って、私はその関数のテストに飽きてきました。どう動くかは理解できている自信があります。では次に進みましょう。

### 最初にテストを書く

```go
func TestMinuteHandPoint(t *testing.T) {
    cases := []struct {
        time  time.Time
        point Point
    }{
        {simpleTime(0, 30, 0), Point{0, -1}},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := minuteHandPoint(c.time)
            if !roughlyEqualPoint(got, c.point) {
                t.Fatalf("Wanted %v Point, but got %v", c.point, got)
            }
        })
    }
}
```

### テストを実行してみます

```text
./clockface_test.go:79:11: undefined: minuteHandPoint
```

### テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func minuteHandPoint(t time.Time) Point {
    return Point{}
}
```

```text
clockface_test.go:80: Wanted {0 -1} Point, but got {0 0}
```

### 成功させるのに十分なコードを書く

```go
func minuteHandPoint(t time.Time) Point {
    return Point{0, -1}
}
```

```text
PASS
ok  	clockface	0.007s
```

### 新しい要件について繰り返します

さて、いよいよ本格的な作業です。

```go
func TestMinuteHandPoint(t *testing.T) {
    cases := []struct {
        time  time.Time
        point Point
    }{
        {simpleTime(0, 30, 0), Point{0, -1}},
        {simpleTime(0, 45, 0), Point{-1, 0}},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := minuteHandPoint(c.time)
            if !roughlyEqualPoint(got, c.point) {
                t.Fatalf("Wanted %v Point, but got %v", c.point, got)
            }
        })
    }
}
```

```text
clockface_test.go:81: Wanted {-1 0} Point, but got {0 -1}
```

### 成功させるのに十分なコードを書く

`secondHandPoint`関数をさっとコピー＆ペーストして少し手を加えれば済むはずです...

```go
func minuteHandPoint(t time.Time) Point {
    angle := minutesInRadians(t)
    x := math.Sin(angle)
    y := math.Cos(angle)

    return Point{x, y}
}
```

```text
PASS
ok  	clockface	0.009s
```

### リファクタリング♪

`minuteHandPoint`と`secondHandPoint`には間違いなく少し重複があります。一方をコピー＆ペーストしてもう一方を作ったのだから当然です。関数を使ってDRYにしましょう。

```go
func angleToPoint(angle float64) Point {
    x := math.Sin(angle)
    y := math.Cos(angle)

    return Point{x, y}
}
```

また、`minuteHandPoint`と`secondHandPoint`をワンライナーとして書き換えることもできます。

```go
func minuteHandPoint(t time.Time) Point {
    return angleToPoint(minutesInRadians(t))
}
```

```go
func secondHandPoint(t time.Time) Point {
    return angleToPoint(secondsInRadians(t))
}
```

```text
PASS
ok  	clockface	0.007s
```

これで、受け入れテストのコメントを外して、分針を描き始めることができます。

### 成功させるのに十分なコードを書く

`minuteHand`関数は、`minuteHandLength`を宣言するなどの小さな調整を加えた`secondHand`のコピー＆ペーストです。

```go
const minuteHandLength = 80

//...

func minuteHand(w io.Writer, t time.Time) {
	p := minuteHandPoint(t)
	p = Point{p.X * minuteHandLength, p.Y * minuteHandLength}
	p = Point{p.X, -p.Y}
	p = Point{p.X + clockCentreX, p.Y + clockCentreY}
	fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#000;stroke-width:3px;"/>`, p.X, p.Y)
}
```

そして、`SVGWriter`関数の中でそれを呼び出します。

```go
func SVGWriter(w io.Writer, t time.Time) {
	io.WriteString(w, svgStart)
	io.WriteString(w, bezel)
	secondHand(w, t)
	minuteHand(w, t)
	io.WriteString(w, svgEnd)
}
```

これで`TestSVGWriterMinuteHand`が成功するはずです。

```text
PASS
ok  	clockface	0.006s
```

しかし、論より証拠です。ここで`clockface`プログラムをコンパイルして実行すると、以下のようなものが表示されるはずです。

![秒針と分針のある時計](<../.gitbook/assets/clock (1).svg>)

### リファクタリング♪

`secondHand`関数と`minuteHand`関数から重複を取り除き、拡大縮小、反転、平行移動のロジックをすべて1か所にまとめましょう。

```go
func secondHand(w io.Writer, t time.Time) {
    p := makeHand(secondHandPoint(t), secondHandLength)
    fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#f00;stroke-width:3px;"/>`, p.X, p.Y)
}

func minuteHand(w io.Writer, t time.Time) {
    p := makeHand(minuteHandPoint(t), minuteHandLength)
    fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#000;stroke-width:3px;"/>`, p.X, p.Y)
}

func makeHand(p Point, length float64) Point {
    p = Point{p.X * length, p.Y * length}
    p = Point{p.X, -p.Y}
    return Point{p.X + clockCentreX, p.Y + clockCentreY}
}
```

```text
PASS
ok  	clockface	0.007s
```

[現在はここまで来ています](https://github.com/quii/learn-go-with-tests/tree/main/math/v9/clockface)。

さあ...あとは時針だけです！

### 最初にテストを書く

```go
func TestSVGWriterHourHand(t *testing.T) {
    cases := []struct {
        time time.Time
        line Line
    }{
        {
            simpleTime(6, 0, 0),
            Line{150, 150, 150, 200},
        },
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            b := bytes.Buffer{}
            clockface.SVGWriter(&b, c.time)

            svg := SVG{}
            xml.Unmarshal(b.Bytes(), &svg)

            if !containsLine(c.line, svg.Line) {
                t.Errorf("Expected to find the hour hand line %+v, in the SVG lines %+v", c.line, svg.Line)
            }
        })
    }
}
```

### テストを実行してみます

```text
clockface_acceptance_test.go:113: Expected to find the hour hand line {X1:150 Y1:150 X2:150 Y2:200}, in the SVG lines [{X1:150 Y1:150 X2:150 Y2:60} {X1:150 Y1:150 X2:150 Y2:70}]
```

ここでも、より下位レベルのテストである程度カバーできるまで、これはコメントアウトしておきましょう。

### 最初にテストを書く

```go
func TestHoursInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(6, 0, 0), math.Pi},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := hoursInRadians(c.time)
            if got != c.angle {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

### テストを実行してみます

```text
./clockface_test.go:97:11: undefined: hoursInRadians
```

### テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func hoursInRadians(t time.Time) float64 {
    return math.Pi
}
```

```text
PASS
ok  	clockface	0.007s
```

### 新しい要件について繰り返します

```go
func TestHoursInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(6, 0, 0), math.Pi},
        {simpleTime(0, 0, 0), 0},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := hoursInRadians(c.time)
            if got != c.angle {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

### テストを実行してみます

```text
clockface_test.go:100: Wanted 0 radians, but got 3.141592653589793
```

### 成功させるのに十分なコードを書く

```go
func hoursInRadians(t time.Time) float64 {
    return (math.Pi / (6 / float64(t.Hour())))
}
```

### 新しい要件について繰り返します

```go
func TestHoursInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(6, 0, 0), math.Pi},
        {simpleTime(0, 0, 0), 0},
        {simpleTime(21, 0, 0), math.Pi * 1.5},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := hoursInRadians(c.time)
            if got != c.angle {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

### テストを実行してみます

```text
clockface_test.go:101: Wanted 4.71238898038469 radians, but got 10.995574287564276
```

### 成功させるのに十分なコードを書く

```go
func hoursInRadians(t time.Time) float64 {
    return (math.Pi / (6 / (float64(t.Hour() % 12))))
}
```

これは24時間時計ではないことに注意してください。剰余演算子を使って、現在の時を12で割った余りを求める必要があります。

```text
PASS
ok  	learn-go-with-tests/math/clockface	0.008s
```

### 最初にテストを書く

次に、経過した分と秒に基づいて、時針を文字盤の周りで動かしてみましょう。

```go
func TestHoursInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(6, 0, 0), math.Pi},
        {simpleTime(0, 0, 0), 0},
        {simpleTime(21, 0, 0), math.Pi * 1.5},
        {simpleTime(0, 1, 30), math.Pi / ((6 * 60 * 60) / 90)},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := hoursInRadians(c.time)
            if got != c.angle {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

### テストを実行してみます

```text
clockface_test.go:102: Wanted 0.013089969389957472 radians, but got 0
```

### 成功させるのに十分なコードを書く

ここでもまた、少し考える必要があります。分と秒の両方に応じて、時針を少しずつ進める必要があります。幸いにも、分と秒に対応する角度はすでに手元にあります。`minutesInRadians`によって返される角度です。再利用できます！

したがって、残る問題は、その角度をどれだけの割合で小さくするかだけです。 1回転は分針の場合は1時間ですが、時針の場合は12時間です。したがって、`minutesInRadians`によって返される角度を12で割ります。

```go
func hoursInRadians(t time.Time) float64 {
    return (minutesInRadians(t) / 12) +
        (math.Pi / (6 / float64(t.Hour()%12)))
}
```

すると、どうでしょう。

```text
clockface_test.go:104: Wanted 0.013089969389957472 radians, but got 0.01308996938995747
```

浮動小数点演算がまたやってくれました。

角度の比較に`roughlyEqualFloat64`を使用するようにテストを更新してみましょう。

```go
func TestHoursInRadians(t *testing.T) {
    cases := []struct {
        time  time.Time
        angle float64
    }{
        {simpleTime(6, 0, 0), math.Pi},
        {simpleTime(0, 0, 0), 0},
        {simpleTime(21, 0, 0), math.Pi * 1.5},
        {simpleTime(0, 1, 30), math.Pi / ((6 * 60 * 60) / 90)},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := hoursInRadians(c.time)
            if !roughlyEqualFloat64(got, c.angle) {
                t.Fatalf("Wanted %v radians, but got %v", c.angle, got)
            }
        })
    }
}
```

```text
PASS
ok  	clockface	0.007s
```

### リファクタリング♪

ラジアンのテストの _1つ_ で`roughlyEqualFloat64`を使うのなら、 _すべて_ のテストで使うべきでしょう。それは素晴らしくて単純なリファクタリングで、[このような状態](https://github.com/quii/learn-go-with-tests/tree/main/math/v10/clockface)になります。

## 時針の位置

さて、今度は単位ベクトルを求めて、時針の先端がどこに来るかを計算します。

### 最初にテストを書く

```go
func TestHourHandPoint(t *testing.T) {
    cases := []struct {
        time  time.Time
        point Point
    }{
        {simpleTime(6, 0, 0), Point{0, -1}},
        {simpleTime(21, 0, 0), Point{-1, 0}},
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            got := hourHandPoint(c.time)
            if !roughlyEqualPoint(got, c.point) {
                t.Fatalf("Wanted %v Point, but got %v", c.point, got)
            }
        })
    }
}
```

待って、私は _2つ_ のテストケースを _一度に_ 書くのですか？これは _悪いTDD_ ではありませんか？

### TDD熱狂について

テスト駆動開発は宗教ではありません。一部の人々はそのように振る舞うかもしれません。
たいていは、TDDを実践していないのに、TDDは熱狂者だけがやるものだ、自分たちがテストを書かないのは「実用的」だからだ、とTwitterやDev.toで喜んで不満を言う人々です。しかしそれは宗教ではありません。それはツールです。

私は2つのテストがどうなるかを _知っています_ 。他の2つの時計の針をまったく同じ方法でテストしました。そして、実装がどうなるかもすでに知っています。分針のイテレーションで、角度を点に変換する一般的なケースの関数を書いたからです。

形式のためだけにTDDの儀式をこなすつもりはありません。TDDは、自分が書いているコード、そしてこれから書こうとしているコードをよりよく理解するのに役立つテクニックです。TDDはフィードバック、知識、洞察を与えてくれます。しかし、すでにその知識を持っているのであれば、理由もなく儀式をこなすつもりはありません。テストもTDDもそれ自体が目的ではありません。

自信がついたので、大きく前進できると感じています。いくつかのステップを「スキップ」します。自分がどこにいて、どこへ向かうのかがわかっているし、この道は以前にも通ったことがあるからです。

しかし、また注意してください。
私はテストを書くことを完全にスキップしているわけではありません。今でもテストを最初に書いています。ただ、より粒度の粗い単位で登場しているだけです。

### テストを実行してみます

```text
./clockface_test.go:119:11: undefined: hourHandPoint
```

### 成功させるのに十分なコードを書く

```go
func hourHandPoint(t time.Time) Point {
    return angleToPoint(hoursInRadians(t))
}
```

言ったとおり、自分がどこにいて、どこへ向かうのかはわかっています。
わからないふりをする必要がどこにあるでしょう？間違っていれば、テストがすぐに教えてくれます。

```text
PASS
ok  	learn-go-with-tests/math/clockface	0.009s
```

## 時針を描く

そしていよいよ、時針を描きます。受け入れテストのコメントを外して復活させましょう。

```go
func TestSVGWriterHourHand(t *testing.T) {
    cases := []struct {
        time time.Time
        line Line
    }{
        {
            simpleTime(6, 0, 0),
            Line{150, 150, 150, 200},
        },
    }

    for _, c := range cases {
        t.Run(testName(c.time), func(t *testing.T) {
            b := bytes.Buffer{}
            clockface.SVGWriter(&b, c.time)

            svg := SVG{}
            xml.Unmarshal(b.Bytes(), &svg)

            if !containsLine(c.line, svg.Line) {
                t.Errorf("Expected to find the hour hand line %+v, in the SVG lines %+v", c.line, svg.Line)
            }
        })
    }
}
```

### テストを実行してみます

```text
clockface_acceptance_test.go:113: Expected to find the hour hand line {X1:150 Y1:150 X2:150 Y2:200},
    in the SVG lines [{X1:150 Y1:150 X2:150 Y2:60} {X1:150 Y1:150 X2:150 Y2:70}]
```

### 成功させるのに十分なコードを書く

そして、SVGを書き出す定数と関数に最終調整を加えることができます。

```go
const (
    secondHandLength = 90
    minuteHandLength = 80
    hourHandLength   = 50
    clockCentreX     = 150
    clockCentreY     = 150
)

//SVGWriter writes an SVG representation of an analogue clock, showing the time t, to the writer w
func SVGWriter(w io.Writer, t time.Time) {
    io.WriteString(w, svgStart)
    io.WriteString(w, bezel)
    secondHand(w, t)
    minuteHand(w, t)
    hourHand(w, t)
    io.WriteString(w, svgEnd)
}

// ...

func hourHand(w io.Writer, t time.Time) {
    p := makeHand(hourHandPoint(t), hourHandLength)
    fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#000;stroke-width:3px;"/>`, p.X, p.Y)
}
```

すると...

```text
ok  	clockface	0.007s
```

`clockface`プログラムをコンパイルして実行して確認してみましょう。

![時計](<../.gitbook/assets/clock (2).svg>)

### リファクタリング♪

`clockface.go`を見ると、いくつかの「マジックナンバー」が散らばっています。これらはすべて、文字盤の半回転が何時間/何分/何秒にあたるかに基づいています。
リファクタリングして、意味を明確にしましょう。

```go
const (
    secondsInHalfClock = 30
    secondsInClock     = 2 * secondsInHalfClock
    minutesInHalfClock = 30
    minutesInClock     = 2 * minutesInHalfClock
    hoursInHalfClock   = 6
    hoursInClock       = 2 * hoursInHalfClock
)
```

なぜこうするのでしょうか？それは、式の中の各数値が何を _意味する_ のかを明示できるからです。このコードに戻ってくることがあれば、いや、 _いずれ必ず_ 戻ってくるときには、これらの名前が何が起きているのかを理解する助けになります。

さらに、本当に、本当に「奇妙な」時計、たとえば時針が4時間で、秒針が20秒で1周するような時計を作りたくなったとしても、これらの定数は簡単にパラメータにできます。
その扉を開けたままにしておけるわけです（たとえ実際にその扉をくぐることがないとしても）。

## まとめ

他に何かする必要がありますか？

まず、自分をほめてあげましょう。SVGの文字盤を作るプログラムを書き上げました。ちゃんと動作しますし、素晴らしいものです。

作れる文字盤は1種類だけですが、それで問題ありません！ 1種類の文字盤だけが _欲しい_ のかもしれません。特定の問題だけを解決し、それ以外は何もしないプログラムでも、何も悪いことはありません。

### プログラム...とライブラリ

しかし、私たちが書いたコードは、文字盤の描画に関するより一般的な問題を _確かに_ 解決しています。問題の小さな部分を1つずつ切り離して考えるためにテストを使い、その切り離しを関数としてコード化したので、文字盤の計算のための、かなり手頃な小さい API ができあがりました。

私たちはこのプロジェクトに取り組んで、より一般的なもの、つまり文字盤の角度やベクトルを計算するためのライブラリにすることができます。

実際、このライブラリをプログラムと一緒に提供するのは本当に良いアイデアです。何のコストもかからないし、プログラムの有用性を高め、それがどのように動作するかを文書化するのにも役立ちます。

> API はプログラムと一緒に提供されるべきであり、その逆もまた然りです。使用するために C コードを書かなければならず、コマンドラインから簡単に呼び出すことができない API は、学習して使用するのが難しくなります。そして逆に、唯一のオープンでドキュメント化された形がプログラムであり、Cプログラムから簡単に呼び出すことができないインターフェースを持つことは、非常に苦痛です。-- ヘンリー・スペンサー、_The Art of Unix Programming_の中で

[このプログラムの最終版](https://github.com/quii/learn-go-with-tests/tree/main/math/vFinal/clockface)では、`clockface`の中にあるエクスポートされていない関数をライブラリのパブリックAPIにして、時計の針の角度と単位ベクトルを計算する関数を用意しました。また、SVG生成部分を独自のパッケージである`svg`に分割し、`clockface`プログラムが直接使用するようにしました。当然のことながら、それぞれの関数とパッケージのドキュメントを作成しました。

SVGといえば

### 最も価値のあるテスト

SVG を扱うための最も洗練されたコードは、アプリケーションコードには全くなく、テストコードにあることにお気づきでしょう。これは落ち着かないことでしょうか？以下のようなことをすべきではないでしょうか？

* `text/template`のテンプレートを使う
* （テストでやっているのと同じように）XMLライブラリを使う
* SVGライブラリを使う

これらのことを行うためにコードをリファクタリングすることができます。それができるのは、SVG を _どのように_ 生成するかは重要ではなく、重要なのは _何を_ 生成するか、つまり _SVG_ だからです。そのため、システムの中で SVG について最もよく知っている必要があり、何が SVG を構成するかについて最も厳密である必要がある部分は、SVG 出力のためのテストです。SVG を出力していると確信できるように、SVG とは何かについての十分なコンテキストと知識を持っている必要があります。SVG の _何_ はテストの中にあり、 _どのように_ はコードの中にあります。

XML ライブラリのインポート、XML の解析、構造体のリファクタリングなど、SVG テストに多くの時間と労力を費やしていることに違和感を感じたかもしれませんが、このテストコードはコードベースの貴重な部分であり、現在の本番コードよりも価値があるかもしれません。生成に何を使うことにしても、出力が常に有効な SVG であることを保証するのに役立ちます。

テストは二級市民ではありません。「使い捨て」のコードでもありません。良いテストは、テストしているコードのバージョンよりもずっと長持ちします。テストを書くのに「時間がかかりすぎる」と感じるべきではありません。それは投資です。

1. 簡単に言えば、普通の度数を使うと角度としてπが何度も登場するため、円を使った微積分が簡単になります。角度をπ単位で数えれば、すべての方程式がより簡単になります。
