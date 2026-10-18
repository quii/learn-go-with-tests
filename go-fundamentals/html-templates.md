---
description: HTML Templates
---

# HTML テンプレート

**[この章のすべてのコードはここにあります](https://github.com/quii/learn-go-with-tests/tree/main/blogrenderer)**

私たちは、誰もが、何ギガバイトものトランスパイルされた JavaScript の上に構築された今月流行りのフロントエンドフレームワークを使い、複雑怪奇なビルドシステムと格闘しながら Web アプリケーションを作りたがる世界に生きています。[しかし、それは必ずしも必要ではないかもしれません](https://quii.dev/The_Web_I_Want)。

ほとんどの Go 開発者は、シンプルで安定した高速なツールチェーンを重視していると思いますが、フロントエンドの世界はこの点でしばしば期待に応えられていません。

多くの Web サイトは [SPA](https://en.wikipedia.org/wiki/Single-page_application) である必要はありません。**HTML と CSS はコンテンツを届けるための素晴らしい手段**であり、Go を使って HTML を配信する Web サイトを作ることができます。

それでも動的な要素が欲しい場合は、クライアントサイドの JavaScript を少し加えることもできますし、サーバーサイドのアプローチで動的な体験を提供できる [Hotwire](https://hotwired.dev) を試してみるのもよいでしょう。

[`fmt.Fprintf`](https://pkg.go.dev/fmt#Fprintf) を駆使して Go で HTML を生成することもできますが、この章では、Go の標準ライブラリに HTML をよりシンプルで保守しやすい方法で生成するためのツールがあることを学びます。また、これまで出会ったことがないかもしれない、この種のコードをより効果的にテストする方法についても学びます。

## これから作るもの

[ファイルの読み込み](reading-files.md)の章では、[`fs.FS`](https://pkg.go.dev/io/fs)（ファイルシステム）を受け取り、見つかった各マークダウンファイルに対応する `Post` のスライスを返すコードを書きました。

```go
posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
```

`Post` は次のように定義しました。

```go
type Post struct {
	Title, Description, Body string
	Tags                     []string
}
```

パース可能なマークダウンファイルの例を示します。

```markdown
Title: Welcome to my blog
Description: Introduction to my blog
Tags: cooking, family, live-laugh-love
---
# First recipe!
Welcome to my **amazing recipe blog**. I am going to write about my family recipes, and make sure I write a long, irrelevant and boring story about my family before you get to the actual instructions.
```

ブログソフトウェアを書く旅を続けるなら、このデータを受け取って HTML を生成し、Web サーバーが HTTP リクエストへのレスポンスとして返せるようにすることになるでしょう。

このブログでは、2 種類のページを生成したいと考えています。

1. **投稿の表示**。特定の投稿をレンダリングします。`Post` の `Body` フィールドはマークダウンを含む文字列なので、HTML に変換する必要があります。
2. **インデックス**。すべての投稿を一覧表示し、特定の投稿を表示するためのハイパーリンクを付けます。

また、サイト全体で一貫したルック＆フィールにしたいので、各ページには `<html>` や、CSS スタイルシートへのリンクなどを含む `<head>` といった、おなじみの HTML の骨組みを用意します。

ブログソフトウェアを作るとき、HTML を構築してユーザーのブラウザに送る方法にはいくつかの選択肢があります。

ここでは `io.Writer` を受け取るようにコードを設計します。これにより、コードの呼び出し側は柔軟に次のことができます。

- [os.File](https://pkg.go.dev/os#File) に書き込み、静的に配信できるようにする
- [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter) に HTML を直接書き出す
- あるいは、本当に何にでも書き込めます！`io.Writer` を実装してさえいれば、ユーザーは `Post` から HTML を生成できます

## 最初にテストを書く

いつものように、急いで飛び込む前に要件について考えることが重要です。このやや大きめの要件を、集中して取り組める小さく達成可能なステップにどう分解すればよいでしょうか？

私の考えでは、実際にコンテンツを表示することのほうがインデックスページよりも優先度が高いです。この製品をリリースして、素晴らしいコンテンツへの直接リンクを共有することもできます。実際のコンテンツにリンクできないインデックスページは役に立ちません。

それでも、先ほど説明したような投稿のレンダリングはまだ大きく感じます。HTML の骨組み、本文のマークダウンを HTML に変換すること、タグの一覧表示などなど。

この段階では、具体的なマークアップについてはあまり気にしていません。簡単な最初のステップとして、投稿のタイトルを `<h1>` としてレンダリングできることを確認するだけにしましょう。これが、少しでも前に進める最小の第一歩だと*感じられます*。

```go
package blogrenderer_test

import (
	"bytes"
	"github.com/quii/learn-go-with-tests/blogrenderer"
	"testing"
)

func TestRender(t *testing.T) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := blogrenderer.Render(&buf, aPost)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>hello world</h1>`
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}
```

`io.Writer` を受け取るという決定は、テストも簡単にしてくれます。この場合は [`bytes.Buffer`](https://pkg.go.dev/bytes#Buffer) に書き込み、後でその内容を検査できます。

## テストを実行してみます

この本のこれまでの章を読んできたなら、もう慣れたものでしょう。パッケージも `Render` 関数も定義していないので、テストは実行できません。コンパイラのメッセージに自分で従って、テストを実行でき、明確なメッセージで失敗することを確認できる状態にしてみてください。

テストが失敗することを実際に確かめるのは本当に重要です。6 か月後に誤ってテストを失敗させてしまったとき、明確なメッセージで失敗することを*今*確認する手間をかけた自分に感謝することでしょう。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

これがテストを実行するための最小限のコードです。

```go
package blogrenderer

// if you're continuing from the read files chapter, you shouldn't redefine this
type Post struct {
	Title, Description, Body string
	Tags                     []string
}

func Render(w io.Writer, p Post) error {
	return nil
}
```

テストは、空文字列が期待する値と等しくないと文句を言うはずです。

## 成功させるのに十分なコードを書く

```go
func Render(w io.Writer, p Post) error {
	_, err := fmt.Fprintf(w, "<h1>%s</h1>", p.Title)
	return err
}
```

ソフトウェア開発は主に学習活動であることを忘れないでください。作業しながら発見し学ぶためには、頻繁で質の高いフィードバックループが得られる方法で作業する必要があり、そのための最も簡単な方法は小さなステップで作業することです。

ですから、今はテンプレートライブラリを使うことについて心配していません。「普通の」文字列テンプレートだけでも十分に HTML を作れますし、テンプレートの部分を飛ばすことで、ちょっとした有用な振る舞いを検証でき、パッケージの API について少し設計作業をしたことになります。

## リファクタリング♪

まだリファクタリングすることはあまりないので、次のイテレーションに進みましょう。

## 最初にテストを書く

ごく基本的なバージョンが動くようになったので、テストを反復して機能を拡張していきます。今回は、`Post` からより多くの情報をレンダリングします。

```go
	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := blogrenderer.Render(&buf, aPost)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>hello world</h1>
<p>This is a description</p>
Tags: <ul><li>go</li><li>tdd</li></ul>`

		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
```

これを書いていて、ぎこちなく*感じる*ことに注目してください。テストの中にこれだけのマークアップがあるのは気分が悪いですし、まだ本文も、`<head>` の内容や必要なページの骨組みを含む実際の HTML も入れていません。

それでも、*今のところは*この苦痛を我慢しましょう。

## テストを実行してみます

説明とタグをレンダリングしていないので、期待する文字列がないと文句を言って失敗するはずです。

## 成功させるのに十分なコードを書く

コードをコピーするのではなく、自分でやってみてください。このテストを通すのは _少し面倒だ_ と気づくはずです！私が試したとき、最初の試みでは次のエラーが出ました。

```
=== RUN   TestRender
=== RUN   TestRender/it_converts_a_single_post_into_HTML
    renderer_test.go:32: got '<h1>hello world</h1><p>This is a description</p><ul><li>go</li><li>tdd</li></ul>' want '<h1>hello world</h1>
        <p>This is a description</p>
        Tags: <ul><li>go</li><li></li></ul>'
```

改行！誰が気にするでしょうか？そう、テストは気にします。正確な文字列の値で照合しているからです。そうあるべきでしょうか？とりあえずテストを通すために、改行を削除しました。

```go
func Render(w io.Writer, p Post) error {
	_, err := fmt.Fprintf(w, "<h1>%s</h1>\n<p>%s</p>\n", p.Title, p.Description)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "Tags: <ul>")
	if err != nil {
		return err
	}

	for _, tag := range p.Tags {
		_, err = fmt.Fprintf(w, "<li>%s</li>", tag)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "</ul>")
	if err != nil {
		return err
	}

	return nil
}
```

**うわぁ**。私が書いた中で最もきれいなコードとは言えませんし、マークアップの実装はまだごく初期の段階です。ページにはもっと多くのコンテンツや要素が必要になるので、このアプローチが適切でないことはすぐにわかります。

しかし決定的なのは、テストが通っているということです。動くソフトウェアがあるのです。

## リファクタリング♪

動くコードに対して通るテストというセーフティネットがあるので、リファクタリングの段階で実装のアプローチを変えることを考えられます。

### テンプレートの導入

Go には [text/template](https://pkg.go.dev/text/template) と [html/template](https://pkg.go.dev/html/template) という 2 つのテンプレートパッケージがあり、同じインターフェースを共有しています。どちらも、テンプレートとデータを組み合わせて文字列を生成できます。

HTML 版との違いは何でしょうか？

> Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. It provides the same interface as package text/template and should be used instead of text/template whenever the output is HTML.
>
> （template パッケージ（html/template）は、コードインジェクションに対して安全な HTML 出力を生成するためのデータ駆動型テンプレートを実装します。text/template パッケージと同じインターフェースを提供し、出力が HTML である場合は常に text/template の代わりに使用すべきです。）

このテンプレート言語は [Mustache](https://mustache.github.io) によく似ており、関心の分離をきれいに保ちながら、非常にすっきりとした方法でコンテンツを動的に生成できます。あなたが使ったことのある他のテンプレート言語と比べると、非常に制約が強い、Mustache の言葉を借りれば「ロジックレス」なものです。これは重要な、**そして意図的な**設計上の決定です。

ここでは HTML の生成に焦点を当てていますが、もしプロジェクトで複雑な文字列連結や呪文のような処理をしているなら、`text/template` を使ってコードを整理するとよいかもしれません。

### コードに戻る

これがブログ用のテンプレートです。

`<h1>{{.Title}}</h1><p>{{.Description}}</p>Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>`

この文字列をどこで定義しましょうか？いくつか選択肢がありますが、ステップを小さく保つために、まずは普通の文字列から始めましょう。

```go
package blogrenderer

import (
	"html/template"
	"io"
)

const (
	postTemplate = `<h1>{{.Title}}</h1><p>{{.Description}}</p>Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>`
)

func Render(w io.Writer, p Post) error {
	templ, err := template.New("blog").Parse(postTemplate)
	if err != nil {
		return err
	}

	return templ.Execute(w, p)
}
```

名前を付けて新しいテンプレートを作成し、テンプレート文字列をパースします。その後、`Execute` メソッドにデータ（この場合は `Post`）を渡して使うことができます。

テンプレートは `{{.Description}}` のようなものを `p.Description` の内容で置き換えます。テンプレートには、値をループする `range` や `if` といったプログラミングの基本要素もあります。詳細は [text/template のドキュメント](https://pkg.go.dev/text/template)を参照してください。

*これは純粋なリファクタリングであるべきです。* テストを変更する必要はなく、テストは引き続き通るはずです。重要なのは、コードが読みやすくなり、厄介なエラー処理がずっと少なくなったことです。

Go のエラー処理の冗長さについて不満を言う人はよくいますが、ここでのように、そもそもエラーが起きにくいコードを書くより良い方法が見つかるかもしれません。

### さらなるリファクタリング

`html/template` を使ったことは確かに改善になりましたが、コードの中に文字列定数として置いておくのはあまり良くありません。

- まだかなり読みにくいです。
- IDE やエディタに優しくありません。シンタックスハイライトも、再フォーマットやリファクタリングなどの機能も使えません。
- HTML のように見えますが、「普通の」HTML ファイルのようには扱えません。

テンプレートを別のファイルに置くことで、より整理しやすくし、HTML ファイルのように扱えるようにしたいと思います。

"templates" というフォルダを作成し、その中に `blog.gohtml` というファイルを作って、テンプレートをファイルに貼り付けます。

次に、[Go 1.16 で導入された埋め込み機能](https://pkg.go.dev/embed)を使ってファイルシステムを埋め込むようにコードを変更します。

```go
package blogrenderer

import (
	"embed"
	"html/template"
	"io"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

func Render(w io.Writer, p Post) error {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return err
	}

	return templ.Execute(w, p)
}
```

「ファイルシステム」をコードに埋め込むことで、複数のテンプレートを読み込み、自由に組み合わせることができます。これは、HTML ページの上部のヘッダーやフッターのように、異なるテンプレート間でレンダリングロジックを共有したいときに役立ちます。

### Embed とは？

Embed については[ファイルの読み込み](reading-files.md)で軽く触れました。[標準ライブラリのドキュメントでは次のように説明されています](https://pkg.go.dev/embed)。

> Package embed provides access to files embedded in the running Go program.
>
> Go source files that import "embed" can use the //go:embed directive to initialize a variable of type string, []byte, or FS with the contents of files read from the package directory or subdirectories at compile time.
>
> （embed パッケージは、実行中の Go プログラムに埋め込まれたファイルへのアクセスを提供します。
>
> "embed" をインポートする Go ソースファイルは、//go:embed ディレクティブを使って、コンパイル時にパッケージディレクトリまたはサブディレクトリから読み込んだファイルの内容で、string、[]byte、または FS 型の変数を初期化できます。）

なぜこれを使いたいのでしょうか？代わりに、「普通の」ファイルシステムからテンプレートを読み込むことも _できます_ 。しかしその場合、このソフトウェアを使う場所がどこであっても、テンプレートが正しいファイルパスにあることを保証しなければなりません。仕事では、開発、ステージング、本番といったさまざまな環境があるでしょう。これを機能させるには、テンプレートが正しい場所にコピーされていることを確認する必要があります。

embed を使えば、ビルド時にファイルが Go プログラムに含まれます。つまり、一度プログラムをビルドすれば（ビルドは一度だけ行うべきです）、ファイルは常に利用可能です。

便利なのは、個々のファイルだけでなくファイルシステムも埋め込めることです。そしてそのファイルシステムは [io/fs](https://pkg.go.dev/io/fs) を実装しているので、コードはどんな種類のファイルシステムを扱っているかを気にする必要がありません。

ただし、設定に応じて異なるテンプレートを使いたい場合は、より従来的な方法でディスクからテンプレートを読み込むほうがよいかもしれません。

## 次へ: テンプレートを「きれい」にする

テンプレートを 1 行の文字列として定義したくはありません。読みやすく扱いやすいように、次のように間隔を空けられるようにしたいのです。

```handlebars
<h1>{{.Title}}</h1>

<p>{{.Description}}</p>

Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
```

しかしこうすると、テストが失敗します。テストは非常に特定の文字列が返されることを期待しているからです。

しかし実際には、空白についてはまったく気にしていません。マークアップに少し変更を加えるたびにアサーションの文字列を苦労して更新しなければならないとしたら、このテストの保守は悪夢になるでしょう。テンプレートが大きくなるにつれて、この種の編集は管理が難しくなり、作業コストは手に負えなくなります。

## Approval Tests の導入

[Go Approval Tests](https://github.com/approvals/go-approval-tests)

> ApprovalTests allows for easy testing of larger objects, strings and anything else that can be saved to a file (images, sounds, CSV, etc...)
>
> （ApprovalTests を使うと、大きなオブジェクトや文字列、その他ファイルに保存できるあらゆるもの（画像、音声、CSV など）を簡単にテストできます）

この考え方は「ゴールデン」ファイルやスナップショットテストに似ています。テストファイルの中で文字列をぎこちなく保守する代わりに、approval ツールが、あなたが作成した「承認済み（approved）」ファイルと出力を比較してくれます。新しいバージョンを承認するなら、それをコピーするだけです。テストを再実行すれば、またグリーンに戻ります。

プロジェクトに `"github.com/approvals/go-approval-tests"` への依存関係を（`go get github.com/approvals/go-approval-tests` コマンドで）追加し、テストを次のように編集します。

```go
func TestRender(t *testing.T) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}

		if err := blogrenderer.Render(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}
```

初めて実行すると、まだ何も承認していないので失敗します。

```
=== RUN   TestRender
=== RUN   TestRender/it_converts_a_single_post_into_HTML
    renderer_test.go:29: Failed Approval: received does not match approved.
```

次のような 2 つのファイルが作成されているはずです。

- `renderer_test.TestRender.it_converts_a_single_post_into_HTML.received.txt`
- `renderer_test.TestRender.it_converts_a_single_post_into_HTML.approved.txt`

received ファイルには、まだ承認されていない新しいバージョンの出力が含まれています。それを空の approved ファイルにコピーして、テストを再実行してください。

新しいバージョンをコピーすることで変更を「承認」したことになり、テストは通るようになります。

このワークフローを実際に見るために、先ほど話したように読みやすくなるようテンプレートを編集してみましょう（ただし意味的には同じです）。

```handlebars
<h1>{{.Title}}</h1>

<p>{{.Description}}</p>

Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
```

テストを再実行します。コードの出力が承認済みのバージョンと異なるので、新しい "received" ファイルが生成されます。それらを見比べて、変更に満足したら、新しいバージョンをコピーしてテストを再実行するだけです。承認済みファイルは必ずソース管理にコミットしてください。

このアプローチにより、HTML のような大きくて見苦しいものへの変更の管理がずっと簡単になります。差分ツールを使って違いを確認・管理でき、テストコードもすっきり保てます。

![差分ツールを使って変更を管理する](https://i.imgur.com/0MoNdva.png)

これは実は approval tests のごく控えめな使い方に過ぎず、approval tests はテストの武器庫の中でも非常に有用なツールです。[Emily Bache](https://twitter.com/emilybache) は、[テストがまったくない複雑なコードベースに approval tests を使って非常に広範なテスト群を追加する興味深い動画](https://www.youtube.com/watch?v=zyM2Ep28ED8)を公開しています。「組み合わせテスト（Combinatorial Testing）」は間違いなく調べてみる価値があります。

この変更を行ったことで、コードが十分にテストされている恩恵を引き続き受けつつ、マークアップをいじっているときにテストが邪魔になりすぎることもなくなりました。

### まだ TDD をしているのか？

このアプローチの興味深い副作用は、TDD から離れてしまうことです。もちろん、承認済みファイルを望む状態に手動で編集し、テストを実行してから、定義した通りに出力されるようテンプレートを修正することも _できます_ 。

でもそれは馬鹿げています！TDD は作業、特に設計を行うための手法ですが、だからといって**あらゆること**に教条的に使わなければならないわけではありません。

重要なのは、私たちが正しいことをして、TDD を**設計ツール**として使い、パッケージの API を設計したということです。テンプレートの変更については、プロセスは次のようになります。

- テンプレートに小さな変更を加える
- approval test を実行する
- 出力を目で見て正しそうか確認する
- 承認する
- 繰り返す

小さく達成可能なステップで作業することの価値を手放すべきではありません。変更を小さくする方法を探し、テストを再実行し続けて、自分のしていることについて本物のフィードバックを得るようにしましょう。

もしテンプレートの _周りの_ コードを変更するようなことを始めるなら、もちろん TDD の作業方法に戻るのが妥当かもしれません。

## マークアップを拡張する

ほとんどの Web サイトは、今の私たちのものより豊かな HTML を持っています。まず手始めに `html` 要素、それに `head`、おそらく `nav` もあるでしょう。たいていはフッターという考え方もあります。

サイトに複数の異なるページがある場合、サイトの見た目の一貫性を保つために、これらのものを 1 か所で定義したいでしょう。Go のテンプレートは、セクションを定義して他のテンプレートにインポートすることをサポートしています。

既存のテンプレートを編集して、top と bottom のテンプレートをインポートします。

```handlebars
{{template "top" .}}
<h1>{{.Title}}</h1>

<p>{{.Description}}</p>

Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>
{{template "bottom" .}}
```

次に、以下の内容で `top.gohtml` を作成します。

```handlebars
{{define "top"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <title>My amazing blog!</title>
    <meta charset="UTF-8"/>
    <meta name="description" content="Wow, like and subscribe, it really helps the channel guys" lang="en"/>
</head>
<body>
<nav role="navigation">
    <div>
        <h1>Budding Gopher's blog</h1>
        <ul>
            <li><a href="/">home</a></li>
            <li><a href="about">about</a></li>
            <li><a href="archive">archive</a></li>
        </ul>
    </div>
</nav>
<main>
{{end}}
```

そして `bottom.gohtml` です。

```handlebars
{{define "bottom"}}
</main>
<footer>
    <ul>
        <li><a href="https://twitter.com/quii">Twitter</a></li>
        <li><a href="https://github.com/quii">GitHub</a></li>
    </ul>
</footer>
</body>
</html>
{{end}}
```

（もちろん、好きなマークアップを自由に入れてください！）

ここで、実行するテンプレートを明示的に指定する必要があります。ブログレンダラーで、`Execute` コマンドを `ExecuteTemplate` に変更します。

```go
if err := templ.ExecuteTemplate(w, "blog.gohtml", p); err != nil {
	return err
}
```

テストを再実行します。新しい "received" ファイルが作られ、テストは失敗するはずです。内容を確認して問題なければ、古いバージョンに上書きコピーして承認します。もう一度テストを再実行すれば、通るはずです。

## ベンチマークで遊んでみる口実

先に進む前に、私たちのコードが何をしているかを考えてみましょう。

```go
func Render(w io.Writer, p Post) error {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return err
	}

	return templ.ExecuteTemplate(w, "blog.gohtml", p)
}
```

- テンプレートをパースする
- テンプレートを使って投稿を `io.Writer` にレンダリングする

投稿ごとにテンプレートを再パースすることによるパフォーマンスへの影響は、ほとんどの場合ごくわずかですが、それを*しない*ようにする手間もまたごくわずかで、コードも少し整理されるはずです。

このパースを何度も繰り返さないことの影響を確認するために、ベンチマークツールを使って関数がどれくらい速いかを見てみましょう。

```go
func BenchmarkRender(b *testing.B) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	for b.Loop() {
		blogrenderer.Render(io.Discard, aPost)
	}
}
```

私のコンピューターでの結果は次の通りです。

```
BenchmarkRender-8 22124 53812 ns/op
```

テンプレートを何度も再パースしなくて済むように、パース済みのテンプレートを保持し、レンダリングを行うメソッドを持つ型を作成します。

```go
type PostRenderer struct {
	templ *template.Template
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &PostRenderer{templ: templ}, nil
}

func (r *PostRenderer) Render(w io.Writer, p Post) error {
	return r.templ.ExecuteTemplate(w, "blog.gohtml", p)
}
```

これはコードのインターフェースを変更するので、テストを更新する必要があります。

```go
func TestRender(t *testing.T) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	postRenderer, err := blogrenderer.NewPostRenderer()

	if err != nil {
		t.Fatal(err)
	}

	t.Run("it converts a single post into HTML", func(t *testing.T) {
		buf := bytes.Buffer{}

		if err := postRenderer.Render(&buf, aPost); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
}
```

そしてベンチマークも更新します。

```go
func BenchmarkRender(b *testing.B) {
	var (
		aPost = blogrenderer.Post{
			Title:       "hello world",
			Body:        "This is a post",
			Description: "This is a description",
			Tags:        []string{"go", "tdd"},
		}
	)

	postRenderer, err := blogrenderer.NewPostRenderer()

	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		postRenderer.Render(io.Discard, aPost)
	}
}
```

テストは引き続き通るはずです。ベンチマークはどうでしょうか？

`BenchmarkRender-8 362124 3131 ns/op`。以前の 1 操作あたりのナノ秒は `53812 ns/op` だったので、これはかなりの改善です！インデックスページなど、レンダリングのための他のメソッドを追加していくと、テンプレートのパースを重複させる必要がないので、コードもシンプルになるはずです。

## 本来の作業に戻る

投稿のレンダリングに関して、残っている重要な部分は、実際に `Body` をレンダリングすることです。思い出してほしいのですが、これは著者が書いたマークダウンなので、HTML に変換する必要があります。

これは読者であるあなたへの演習として残しておきます。これを行ってくれる Go のライブラリが見つかるはずです。approval test を使って、自分のしていることを検証してください。

### サードパーティライブラリのテストについて

**注意**。ユニットテストで、サードパーティライブラリがどう振る舞うかを明示的にテストすることにあまりこだわりすぎないよう気をつけてください。

自分が制御できないコードに対してテストを書くのは無駄であり、保守のオーバーヘッドを増やします。テストのために依存関係を制御してその振る舞いをモックしたい場合は、[依存性注入](./dependency-injection.md)を使うこともあるでしょう。

しかしこの場合、私はマークダウンを HTML に変換することをレンダリングの実装の詳細と見なしており、approval tests で十分な確信が得られるはずです。

### インデックスのレンダリング

次に実装する機能は、投稿を HTML の順序付きリストとして一覧表示するインデックスのレンダリングです。

API を拡張するので、再び TDD の帽子をかぶりましょう。

## 最初にテストを書く

一見するとインデックスページは単純に思えますが、テストを書くことで、やはりいくつかの設計上の選択を迫られます。

```go
t.Run("it renders an index of posts", func(t *testing.T) {
	buf := bytes.Buffer{}
	posts := []blogrenderer.Post{{Title: "Hello World"}, {Title: "Hello World 2"}}

	if err := postRenderer.RenderIndex(&buf, posts); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	want := `<ol><li><a href="/post/hello-world">Hello World</a></li><li><a href="/post/hello-world-2">Hello World 2</a></li></ol>`

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
})
```

1. `Post` のタイトルフィールドを URL のパスの一部として使っていますが、URL にスペースは入れたくないので、ハイフンに置き換えています。
2. `PostRenderer` に `RenderIndex` メソッドを追加しました。これも `io.Writer` と `Post` のスライスを受け取ります。

もしここで、テストを後から書く approval tests のアプローチを続けていたら、これらの疑問に制御された環境で答えることはできなかったでしょう。**テストは私たちに考える余地を与えてくれます**。

## テストを実行してみます

```
./renderer_test.go:41:13: undefined: blogrenderer.RenderIndex
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

```go
func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	return nil
}
```

上記で、次のようなテストの失敗が得られるはずです。

```
=== RUN   TestRender
=== RUN   TestRender/it_renders_an_index_of_posts
    renderer_test.go:49: got "" want "<ol><li><a href=\"/post/hello-world\">Hello World</a></li><li><a href=\"/post/hello-world-2\">Hello World 2</a></li></ol>"
--- FAIL: TestRender (0.00s)
```

## 成功させるのに十分なコードを書く

簡単そうに _感じられる_ にもかかわらず、少し厄介です。私は複数のステップに分けて行いました。

```go
func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	indexTemplate := `<ol>{{range .}}<li><a href="/post/{{.Title}}">{{.Title}}</a></li>{{end}}</ol>`

	templ, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return err
	}

	if err := templ.Execute(w, posts); err != nil {
		return err
	}

	return nil
}
```

最初は別のテンプレートファイルに煩わされたくなく、ただ動くようにしたかったのです。テンプレートを事前にパースしたり分離したりするのは、後でできるリファクタリングだと考えています。

これは通りませんが、惜しいところまで来ています。

```
=== RUN   TestRender
=== RUN   TestRender/it_renders_an_index_of_posts
    renderer_test.go:49: got "<ol><li><a href=\"/post/Hello%20World\">Hello World</a></li><li><a href=\"/post/Hello%20World%202\">Hello World 2</a></li></ol>" want "<ol><li><a href=\"/post/hello-world\">Hello World</a></li><li><a href=\"/post/hello-world-2\">Hello World 2</a></li></ol>"
--- FAIL: TestRender (0.00s)
    --- FAIL: TestRender/it_renders_an_index_of_posts (0.00s)
```

テンプレートのコードが `href` 属性内のスペースをエスケープしていることがわかります。スペースをハイフンに文字列置換する方法が必要です。アンカーではユーザーにスペースを表示したいので、単純に `[]Post` をループしてメモリ上で置き換えるわけにはいきません。

いくつかの選択肢があります。最初に検討するのは、テンプレートに関数を渡すことです。

### テンプレートに関数を渡す

```go
func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	indexTemplate := `<ol>{{range .}}<li><a href="/post/{{sanitiseTitle .Title}}">{{.Title}}</a></li>{{end}}</ol>`

	templ, err := template.New("index").Funcs(template.FuncMap{
		"sanitiseTitle": func(title string) string {
			return strings.ToLower(strings.Replace(title, " ", "-", -1))
		},
	}).Parse(indexTemplate)
	if err != nil {
		return err
	}

	if err := templ.Execute(w, posts); err != nil {
		return err
	}

	return nil
}
```

_テンプレートをパースする前に_ 、テンプレートに `template.FuncMap` を追加できます。これにより、テンプレート内で呼び出せる関数を定義できます。今回は `sanitiseTitle` 関数を作成し、テンプレート内で `{{sanitiseTitle .Title}}` として呼び出しています。

これは強力な機能で、テンプレートに関数を渡せると非常にクールなことができますが、そうすべきでしょうか？Mustache とロジックレステンプレートの原則に立ち返ってみると、なぜ彼らはロジックレスを提唱したのでしょうか？**テンプレートにロジックがあると何がいけないのでしょうか？**

これまで見てきたように、テンプレートをテストするために、*まったく別の種類のテストを導入しなければなりませんでした*。

いくつかの異なる振る舞いのパターンやエッジケースを持つ関数をテンプレートに導入したと想像してみてください。**それをどうテストしますか？** 現在の設計では、このロジックをテストする唯一の手段は _HTML をレンダリングして文字列を比較すること_ です。これはロジックをテストする簡単で健全な方法ではなく、 _重要な_ ビジネスロジックに対して望む方法では決してありません。

approval tests の手法によってこれらのテストを保守するコストは下がりましたが、それでもあなたが書くほとんどのユニットテストよりも保守コストは高いです。マークアップにわずかな変更を加えただけでも影響を受けやすいことに変わりはなく、ただ管理しやすくしただけです。それでも、テンプレートの周りに多くのテストを書かなくて済むようにコードを設計し、レンダリングコードの中に置く必要のないロジックはきちんと分離されるよう関心を分離するよう努めるべきです。

Mustache の影響を受けたテンプレートエンジンが与えてくれるのは有用な制約です。それをあまり頻繁に回避しようとしないでください。**流れに逆らわないでください**。代わりに、[ビューモデル](https://stackoverflow.com/a/11074506/3193)という考え方を受け入れましょう。レンダリングに必要なデータを、テンプレート言語にとって扱いやすい形で保持する専用の型を構築するのです。

こうすることで、そのデータの集まりを生成するために使う重要なビジネスロジックは何であれ、HTML やテンプレートの厄介な世界から離れて、個別にユニットテストできます。

### 関心の分離

では、代わりに何ができるでしょうか？

#### `Post` にメソッドを追加し、テンプレートからそれを呼び出す

テンプレートのコードでは、渡した型のメソッドを呼び出せるので、`Post` に `SanitisedTitle` メソッドを追加できます。これによりテンプレートはシンプルになり、必要であればこのロジックを個別に簡単にユニットテストできます。これがおそらく最も簡単な解決策ですが、必ずしも最もシンプルとは限りません。

このアプローチの欠点は、これがやはり _ビュー_ のロジックであることです。システムの他の部分にとっては関心のないものなのに、コアとなるドメインオブジェクトの API の一部になってしまいます。この種のアプローチは、時間が経つにつれて [God Object](https://en.wikipedia.org/wiki/God_object) を生み出すことにつながりかねません。

#### 必要なデータだけを正確に持つ、`PostViewModel` のような専用のビューモデル型を作る

レンダリングコードをドメインオブジェクトである `Post` に結合させるのではなく、代わりにビューモデルを受け取るようにします。

```go
type PostViewModel struct {
	Title, SanitisedTitle, Description, Body string
	Tags                                     []string
}
```

コードの呼び出し側は、`SanitizedTitle` を生成しながら `[]Post` から `[]PostView` へマッピングしなければならなくなります。これをきれいに保つ方法の 1 つは、マッピングをカプセル化する `func NewPostView(p Post) PostView` を用意することです。

これによりレンダリングコードはロジックレスに保たれ、おそらく最も厳密な関心の分離になりますが、その代償として、投稿をレンダリングするまでの過程がやや複雑になります。

どちらの選択肢も問題ありませんが、今回は前者を選びたいと思います。システムを発展させていく中で、レンダリングを円滑にするためだけにその場しのぎのメソッドをどんどん追加することには注意すべきです。ドメインオブジェクトとビューの間の変換がより複雑になるにつれて、専用のビューモデルはより有用になります。

では、`Post` にメソッドを追加しましょう。

```go
func (p Post) SanitisedTitle() string {
	return strings.ToLower(strings.Replace(p.Title, " ", "-", -1))
}
```

そして、レンダリングコードはよりシンプルな世界に戻れます。

```go
func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	indexTemplate := `<ol>{{range .}}<li><a href="/post/{{.SanitisedTitle}}">{{.Title}}</a></li>{{end}}</ol>`

	templ, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return err
	}

	if err := templ.Execute(w, posts); err != nil {
		return err
	}

	return nil
}
```

## リファクタリング♪

ついにテストが通るはずです。これでテンプレートをファイル（`templates/index.gohtml`）に移動し、レンダラーを構築するときに一度だけ読み込むようにできます。

```go
package blogrenderer

import (
	"embed"
	"html/template"
	"io"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

type PostRenderer struct {
	templ *template.Template
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &PostRenderer{templ: templ}, nil
}

func (r *PostRenderer) Render(w io.Writer, p Post) error {
	return r.templ.ExecuteTemplate(w, "blog.gohtml", p)
}

func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	return r.templ.ExecuteTemplate(w, "index.gohtml", posts)
}
```

複数のテンプレートを `templ` にパースしたので、`ExecuteTemplate` を呼び出して、 _どの_ テンプレートをレンダリングしたいかを適宜指定する必要があります。しかし、たどり着いたコードは素晴らしいものになったと同意していただけると思います。

誰かがテンプレートファイルの名前を変更するとバグが入り込むという _わずかな_ リスクはありますが、高速に実行できるユニットテストがすぐにそれを検出してくれるでしょう。

パッケージの API 設計に満足し、TDD で基本的な振る舞いを導き出せたので、テストを approvals を使うように変更しましょう。

```go
	t.Run("it renders an index of posts", func(t *testing.T) {
		buf := bytes.Buffer{}
		posts := []blogrenderer.Post{{Title: "Hello World"}, {Title: "Hello World 2"}}

		if err := postRenderer.RenderIndex(&buf, posts); err != nil {
			t.Fatal(err)
		}

		approvals.VerifyString(t, buf.String())
	})
```

テストを実行して失敗することを確認し、それから変更を承認するのを忘れないでください。

最後に、インデックスページにページの骨組みを追加できます。

```handlebars
{{template "top" .}}
<ol>{{range .}}<li><a href="/post/{{.SanitisedTitle}}">{{.Title}}</a></li>{{end}}</ol>
{{template "bottom" .}}
```

テストを再実行し、変更を承認すれば、インデックスは完成です！

## マークダウンの本文をレンダリングする

自分で試してみることをお勧めしましたが、私が最終的に採用したアプローチを紹介します。

```go
package blogrenderer

import (
	"embed"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"html/template"
	"io"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

type PostRenderer struct {
	templ *template.Template
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &PostRenderer{templ: templ}, nil
}

func (r *PostRenderer) Render(w io.Writer, p Post) error {
	return r.templ.ExecuteTemplate(w, "blog.gohtml", newPostVM(p))
}

func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	return r.templ.ExecuteTemplate(w, "index.gohtml", posts)
}

type postViewModel struct {
	Post
	HTMLBody template.HTML
}

func newPostVM(p Post) postViewModel {
	vm := postViewModel{Post: p}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)
	vm.HTMLBody = template.HTML(markdown.ToHTML([]byte(p.Body), mdParser, nil))
	return vm
}
```

私は素晴らしい [gomarkdown](https://github.com/gomarkdown/markdown) ライブラリを使いましたが、期待通りに動いてくれました。

一度だけ構築して `PostRenderer` に保持するのではなく、`newPostVM` を呼び出すたびに新しい `parser.Parser` を作成していることに注目してください。通常、何かを一度構築して再利用できるなら、そうする価値があります。繰り返し構築するコストを節約できるからです。しかし、その一般的な直感はここでは通用しません。gomarkdown のパーサーはドキュメントツリーを構築する間に内部状態を保持するため、複数の `Parse` 呼び出しにまたがって再利用するのは安全ではありません。同じパーサーを 2 回目に使うと panic が発生します（古いバージョンでは、ずっと不親切な nil ポインタ参照エラーになります）。`PostRenderer` はその生存期間中に多くの投稿を `Render` することを想定しているので、呼び出しごとに新しいパーサーを作るのが正しい使い方です。いずれにせよ、パース処理そのものに比べればパーサーの構築は取るに足らないので、コストは安く済みます。

自分でこれをやってみた場合、本文のレンダリングで HTML がエスケープされてしまったことに気づいたかもしれません。これは、悪意のあるサードパーティの HTML が出力されるのを防ぐための、Go の html/template パッケージのセキュリティ機能です。

これを回避するには、レンダリングに渡す型の中で、信頼できる HTML を [template.HTML](https://pkg.go.dev/html/template#HTML) でラップする必要があります。

> HTML encapsulates a known safe HTML document fragment. It should not be used for HTML from a third-party, or HTML with unclosed tags or comments. The outputs of a sound HTML sanitiser and a template escaped by this package are fine for use with HTML.
>
> Use of this type presents a security risk: the encapsulated content should come from a trusted source, as it will be included verbatim in the template output.
>
> （HTML は、安全であることがわかっている HTML ドキュメントの断片をカプセル化します。サードパーティからの HTML や、閉じられていないタグやコメントを含む HTML に使用すべきではありません。健全な HTML サニタイザーの出力や、このパッケージによってエスケープされたテンプレートの出力は、HTML で使用しても問題ありません。
>
> この型の使用にはセキュリティ上のリスクがあります。カプセル化されたコンテンツはテンプレートの出力にそのまま含まれるため、信頼できるソースから来たものでなければなりません。）

そこで私は**エクスポートされない**ビューモデル（`postViewModel`）を作成しました。これはやはりレンダリングの内部的な実装の詳細だと考えたからです。これを個別にテストする必要はありませんし、API を汚したくもありません。

レンダリング時にこれを構築して `Body` を `HTMLBody` にパースし、テンプレートではそのフィールドを使って HTML をレンダリングしています。

## まとめ

[ファイルの読み込み](reading-files.md)の章とこの章で学んだことを組み合わせれば、十分にテストされたシンプルな静的サイトジェネレーターを無理なく作り、自分のブログを立ち上げることができます。CSS のチュートリアルを探せば、見た目をきれいにすることもできます。

このアプローチはブログにとどまりません。データベース、API、ファイルシステムなど、あらゆるソースからデータを取得して HTML に変換し、サーバーから返すというのは、何十年にもわたって使われてきたシンプルな手法です。人々は現代の Web 開発の複雑さを嘆きがちですが、その複雑さを自ら背負い込んでいるだけではないと言い切れるでしょうか？

Go は Web 開発に最適です。特に、作ろうとしている Web サイトの本当の要件が何なのかを明確に考えているときにはそうです。サーバー上で HTML を生成することは、React のような技術で「Web アプリケーション」を作るよりも、多くの場合より良く、よりシンプルで、より高性能なアプローチです。

### 学んだこと

- HTML テンプレートを作成してレンダリングする方法。
- テンプレートを組み合わせ、関連するマークアップを [DRY](https://en.wikipedia.org/wiki/Don't_repeat_yourself) にして、一貫したルック＆フィールを保つ方法。
- テンプレートに関数を渡す方法と、それについてよく考えるべき理由。
- テンプレートレンダラーのような大きくて見苦しい出力をテストするのに役立つ「Approval Tests」の書き方。

### ロジックレステンプレートについて

いつものように、これはすべて**関心の分離**に関することです。システムのさまざまな部分の責務が何であるかを考えることが重要です。重要なビジネスロジックをテンプレートに漏れ出させ、関心を混同させて、システムを理解・保守・テストしにくくしてしまうことがあまりにも多いのです。

### HTML だけではない

Go には、テンプレートから他の種類のデータを生成するための `text/template` があることを忘れないでください。データを何らかの構造化された出力に変換する必要が出てきたときには、この章で紹介した手法が役立つでしょう。

### 参考資料とさらなる学習

- [John Calhoun の 'Learn Web Development with Go'](https://www.calhoun.io/intro-to-templates-p1-contextual-encoding/) には、テンプレートに関する優れた記事が多数あります。
- [Hotwire](https://hotwired.dev) - これらの手法を使って Hotwire の Web アプリケーションを作成できます。主に Ruby on Rails を使う Basecamp によって作られましたが、サーバーサイドの技術なので Go でも使えます。
