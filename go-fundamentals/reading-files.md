---
description: ファイルの読み込み
---

# ファイルの読み込み

* [**この章のすべてのコードはここにあります**](https://github.com/quii/learn-go-with-tests/tree/main/reading-files)
* [Twitch ストリームでこの問題に取り組み、質問に答えている動画はこちらです](https://www.youtube.com/watch?v=nXts4dEJnkU)

この章では、いくつかのファイルを読み込み、そこからデータを取り出し、何か役に立つことをする方法を学びます。

あなたが友人と一緒にブログソフトウェアを作っていると想像してください。作者は記事をマークダウンで書き、ファイルの先頭にいくつかのメタデータを記述するという考えです。起動時に、Web サーバーはフォルダを読み込んでいくつかの `Post` を作成し、別の `NewHandler` 関数がそれらの `Post` をブログの Web サーバーのデータソースとして使用します。

私たちは、指定されたブログ記事ファイルのフォルダを `Post` のコレクションに変換するパッケージを作成するよう依頼されました。

### データの例

hello world.md

```markdown
Title: Hello, TDD world!
Description: First post on our wonderful blog
Tags: tdd, go
---
Hello world!

The body of posts starts after the `---`
```

### 期待するデータ

```go
type Post struct {
	Title, Description, Body string
	Tags                     []string
}
```

## 反復的なテスト駆動開発

私たちは、常にシンプルで安全なステップを踏んでゴールに向かう、反復的なアプローチをとります。

そのためには作業を分割する必要がありますが、[「ボトムアップ」](https://en.wikipedia.org/wiki/Top-down_and_bottom-up_design)なアプローチをとるという罠に陥らないよう注意しなければなりません。

作業を始めるときに、自分の過剰な想像力を信じてはいけません。`BlogPostFileParser` のような、すべてを組み合わせて初めて検証できる何らかの抽象化を作りたくなるかもしれません。

これは反復的では _なく_、TDD がもたらすはずの短いフィードバックループを逃しています。

Kent Beck はこう言っています。

> 楽観主義はプログラミングの職業病である。フィードバックはその治療法だ。
>
> Optimism is an occupational hazard of programming. Feedback is the treatment.

そうではなく、私たちのアプローチは、できるだけ早く _本当の_ 利用者価値を届けることに可能な限り近づくよう努めるべきです（これはしばしば「ハッピーパス」と呼ばれます）。エンドツーエンドで少量の利用者価値を届けることができれば、残りの要件についてさらに反復することは、たいてい簡単です。

## どのようなテストを見たいかを考える

始めるにあたって、私たちの考え方とゴールを思い出しましょう。

* **見たいテストを書く**。これから書くコードを、利用者の視点からどのように使いたいかを考えます。
* _何を_ と _なぜ_ に集中し、_どのように_ に気を取られないようにします。

私たちのパッケージは、フォルダを指定すると、いくつかの記事を返してくれる関数を提供する必要があります。

```go
var posts []blogposts.Post
posts = blogposts.NewPostsFromFS("some-folder")
```

これについてのテストを書くには、いくつかのサンプル記事が入ったテスト用のフォルダが必要になるでしょう。_これはそれほど悪いことではありません_ が、いくつかのトレードオフがあります。

* テストごとに、特定の振る舞いをテストするための新しいファイルを作成する必要があるかもしれません
* ファイルの読み込みに失敗する場合など、テストするのが難しい振る舞いがあります
* ファイルシステムにアクセスする必要があるため、テストの実行が少し遅くなります

また、ファイルシステムの特定の実装に不必要に結合してしまいます。

### Go 1.16 で導入されたファイルシステムの抽象化

Go 1.16 では、ファイルシステムの抽象化である [io/fs](https://golang.org/pkg/io/fs/) パッケージが導入されました。

> fs パッケージは、ファイルシステムへの基本的なインターフェースを定義します。ファイルシステムはホストのオペレーティングシステムによって提供されることもあれば、他のパッケージによって提供されることもあります。
>
> Package fs defines basic interfaces to a file system. A file system can be provided by the host operating system but also by other packages.

これにより、特定のファイルシステムへの結合を緩めることができ、必要に応じて異なる実装を注入できるようになります。

> [インターフェースの提供側では、新しい embed.FS 型が fs.FS を実装しており、zip.Reader も同様です。新しい os.DirFS 関数は、オペレーティングシステムのファイルツリーに裏付けられた fs.FS の実装を提供します。](https://golang.org/doc/go1.16#fs)
>
> On the producer side of the interface, the new embed.FS type implements fs.FS, as does zip.Reader. The new os.DirFS function provides an implementation of fs.FS backed by a tree of operating system files.

このインターフェースを使えば、私たちのパッケージの利用者は、標準ライブラリに組み込まれた多くの選択肢を使うことができます。Go の標準ライブラリで定義されたインターフェース（例: `io.fs`、[`io.Reader`](https://golang.org/pkg/io/#Reader)、[`io.Writer`](https://golang.org/pkg/io/#Writer)）を活用する方法を学ぶことは、疎結合なパッケージを書くために不可欠です。そうしたパッケージは、あなたが想像していたのとは異なるコンテキストでも、利用者にほとんど手間をかけさせることなく再利用できます。

私たちの場合、利用者は「本物の」ファイルシステム上のファイルではなく、記事を Go のバイナリに埋め込みたいと思うかもしれません。いずれにせよ、_私たちのコードはそれを気にする必要はありません_。

テストのために、[testing/fstest](https://golang.org/pkg/testing/fstest/) パッケージが [io/FS](https://golang.org/pkg/io/fs/#FS) の実装を提供してくれます。これは [net/http/httptest](https://golang.org/pkg/net/http/httptest/) でおなじみのツールと似ています。

この情報を踏まえると、次のほうがより良いアプローチだと感じられます。

```go
var posts []blogposts.Post
posts = blogposts.NewPostsFromFS(someFS)
```

## 最初にテストを書く

スコープはできるだけ小さく、かつ役に立つものに保つべきです。ディレクトリ内のすべてのファイルを読み込めることを証明できれば、良いスタートになるでしょう。これにより、私たちが書いているソフトウェアに自信が持てます。返された `[]Post` の数が、偽のファイルシステム内のファイル数と同じであることを確認できます。

この章に取り組むための新しいプロジェクトを作成しましょう。

* `mkdir blogposts`
* `cd blogposts`
* `go mod init github.com/{your-name}/blogposts`
* `touch blogposts_test.go`

```go
package blogposts_test

import (
	"testing"
	"testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte("hi")},
		"hello-world2.md": {Data: []byte("hola")},
	}

	posts := blogposts.NewPostsFromFS(fs)

	if len(posts) != len(fs) {
		t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
	}
}
```

テストのパッケージが `blogposts_test` であることに注目してください。TDD をうまく実践するときは、_利用者駆動_ のアプローチをとることを思い出してください。_利用者_ は内部の詳細を気にしないので、私たちも内部の詳細をテストしたくありません。意図したパッケージ名に `_test` を付けることで、パッケージのエクスポートされたメンバーにのみアクセスすることになります。これは、パッケージの実際のユーザーとまったく同じです。

[`testing/fstest`](https://golang.org/pkg/testing/fstest/) をインポートしたことで、[`fstest.MapFS`](https://golang.org/pkg/testing/fstest/#MapFS) 型にアクセスできるようになりました。偽のファイルシステムとして、`fstest.MapFS` をパッケージに渡します。

> MapFS は、テストで使用するためのシンプルなインメモリファイルシステムで、パス名（Open への引数）から、それが表すファイルやディレクトリの情報へのマップとして表現されます。
>
> A MapFS is a simple in-memory file system for use in tests, represented as a map from path names (arguments to Open) to information about the files or directories they represent.

これは、テストファイルのフォルダを管理するよりもシンプルに感じられますし、実行も速くなります。

最後に、利用者の視点から API の使い方をコードとして表現し、正しい数の記事が作成されるかどうかを確認しました。

## テストを実行してみます

```
./blogpost_test.go:15:12: undefined: blogposts
```

## テストを実行するための最小限のコードを記述し、_失敗したテスト出力を確認します_

パッケージが存在しません。新しいファイル `blogposts.go` を作成し、その中に `package blogposts` と書きます。次に、そのパッケージをテストにインポートする必要があります。私の場合、インポートは次のようになりました。

```go
import (
	blogposts "github.com/quii/learn-go-with-tests/reading-files"
	"testing"
	"testing/fstest"
)
```

これで、新しいパッケージに何らかのコレクションを返す `NewPostsFromFS` 関数がないため、テストはコンパイルされなくなります。

```
./blogpost_test.go:16:12: undefined: blogposts.NewPostsFromFS
```

これにより、テストを実行するために関数の骨組みを作ることを強いられます。この時点ではコードについて考えすぎないように注意してください。私たちは、テストを実行できるようにし、期待どおりに失敗することを確認しようとしているだけです。このステップを飛ばすと、前提を見落とし、役に立たないテストを書いてしまうかもしれません。

```go
package blogposts

import "testing/fstest"

type Post struct {
}

func NewPostsFromFS(fileSystem fstest.MapFS) []Post {
	return nil
}
```

これでテストは正しく失敗するはずです。

```
=== RUN   TestNewBlogPosts
    blogposts_test.go:48: got 0 posts, wanted 2 posts
```

## 成功させるのに十分なコードを書く

テストを通すために、これを[「スライム（slime）」](https://deniseyu.github.io/leveling-up-tdd/)することも _できます_。

```go
func NewPostsFromFS(fileSystem fstest.MapFS) []Post {
	return []Post{{}, {}}
}
```

しかし、Denise Yu はこう書いています。

> スライミングは、オブジェクトに「骨組み」を与えるのに役立ちます。インターフェースを設計することとロジックを実行することは 2 つの別の関心事であり、テストを戦略的にスライミングすることで、一度に 1 つのことに集中できます。
>
> Sliming is useful for giving a “skeleton” to your object. Designing an interface and executing logic are two concerns, and sliming tests strategically lets you focus on one at a time.

私たちはすでに構造を持っています。では、代わりに何をすればよいでしょうか？

スコープを絞ったので、必要なのはディレクトリを読み込み、見つかったファイルごとに記事を作成することだけです。ファイルを開いて解析することについては、まだ心配する必要はありません。

```go
func NewPostsFromFS(fileSystem fstest.MapFS) []Post {
	dir, _ := fs.ReadDir(fileSystem, ".")
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts
}
```

[`fs.ReadDir`](https://golang.org/pkg/io/fs/#ReadDir) は、指定された `fs.FS` 内のディレクトリを読み込み、[`[]DirEntry`](https://golang.org/pkg/io/fs/#DirEntry) を返します。

エラーが起こりうるため、私たちの理想化された世界観はすでに崩れています。しかし、今の焦点は設計を変えることではなく _テストを通すこと_ であることを思い出してください。なので、今のところエラーは無視します。

残りのコードは簡単です。エントリを反復処理し、それぞれに対して `Post` を作成し、スライスを返します。

## リファクタリング♪

テストは通っていますが、新しいパッケージは具体的な実装である `fstest.MapFS` に結合しているため、このコンテキストの外では使えません。しかし、そうである必要はありません。`NewPostsFromFS` 関数の引数を、標準ライブラリのインターフェースを受け取るように変更しましょう。

```go
func NewPostsFromFS(fileSystem fs.FS) []Post {
	dir, _ := fs.ReadDir(fileSystem, ".")
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts
}
```

テストを再実行してください。すべてうまく動くはずです。

### エラー処理

先ほど、ハッピーパスを動かすことに集中していたときに、エラー処理を保留にしました。機能の反復を続ける前に、ファイルを扱う際にはエラーが起こりうることを認めるべきです。ディレクトリの読み込み以外にも、個々のファイルを開くときに問題が起こる可能性があります。`error` を返せるように API を変更しましょう（もちろん、まずはテストから）。

```go
func TestNewBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte("hi")},
		"hello-world2.md": {Data: []byte("hola")},
	}

	posts, err := blogposts.NewPostsFromFS(fs)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != len(fs) {
		t.Errorf("got %d posts, wanted %d posts", len(posts), len(fs))
	}
}
```

テストを実行してください。戻り値の数が間違っていると文句を言われるはずです。コードの修正は簡単です。

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for range dir {
		posts = append(posts, Post{})
	}
	return posts, nil
}
```

これでテストは通ります。あなたの中の TDD 実践者は、`fs.ReadDir` からのエラーを伝播させるコードを書く前に、失敗するテストを見なかったことにいら立つかもしれません。これを「正しく」行うには、失敗する `fs.FS` のテストダブルを注入して、`fs.ReadDir` が `error` を返すようにする新しいテストが必要になります。

```go
type StubFailingFS struct {
}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}
```

```go
// later
_, err := blogposts.NewPostsFromFS(StubFailingFS{})
```

これで、私たちのアプローチに自信が持てるはずです。使っているインターフェースにはメソッドが 1 つしかないため、さまざまなシナリオをテストするためのテストダブルを作るのは簡単です。

エラー処理をテストすることが現実的な場合もありますが、今回のケースでは、エラーに対して何も _興味深い_ ことをしておらず、単に伝播させているだけなので、新しいテストを書く手間をかける価値はありません。

論理的に考えて、次の反復は、`Post` 型を拡張して何か役に立つデータを持たせることになります。

## 最初にテストを書く

提案されたブログ記事のスキーマの最初の行である、タイトルフィールドから始めましょう。

テストファイルの内容を指定されたものと一致するように変更し、それが正しく解析されることをアサーションします。

```go
func TestNewBlogPosts(t *testing.T) {
	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte("Title: Post 1")},
		"hello-world2.md": {Data: []byte("Title: Post 2")},
	}

	// rest of test code cut for brevity
	got := posts[0]
	want := blogposts.Post{Title: "Post 1"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
```

## テストを実行してみます

```
./blogpost_test.go:58:26: unknown field 'Title' in struct literal of type blogposts.Post
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

テストが実行できるように、`Post` 型に新しいフィールドを追加します。

```go
type Post struct {
	Title string
}
```

テストを再実行すると、明確に失敗するテストが得られるはずです。

```
=== RUN   TestNewBlogPosts
=== RUN   TestNewBlogPosts/parses_the_post
    blogpost_test.go:61: got {Title:}, want {Title:Post 1}
```

## 成功させるのに十分なコードを書く

各ファイルを開いて、タイトルを抽出する必要があります。

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(fileSystem, f)
		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
	postFile, err := fileSystem.Open(f.Name())
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()

	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

この時点での焦点は、エレガントなコードを書くことではなく、動くソフトウェアがある地点にたどり着くことだけだということを忘れないでください。

小さな一歩前進のように感じられますが、それでもかなりの量のコードを書き、エラー処理に関していくつかの前提を置く必要がありました。ここは同僚と話し合い、最善のアプローチを決めるべきポイントでしょう。

反復的なアプローチにより、要件についての私たちの理解が不完全であるという素早いフィードバックが得られました。

`fs.FS` は、`Open` メソッドによって、その中のファイルを名前で開く方法を提供してくれます。そこからファイルのデータを読み込みます。今のところ高度な解析は必要なく、文字列をスライスして `Title:` というテキストを切り取るだけです。

## リファクタリング♪

「ファイルを開くコード」と「ファイルの内容を解析するコード」を分離すると、コードが理解しやすく、扱いやすくなります。

```go
func getPost(fileSystem fs.FS, f fs.DirEntry) (Post, error) {
	postFile, err := fileSystem.Open(f.Name())
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile fs.File) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

新しい関数やメソッドをリファクタリングで切り出すときは、注意して引数について考えてください。ここで行っているのは設計であり、テストが通っているので、何が適切かを深く考えることができます。結合度と凝集度について考えましょう。今回のケースでは、自分自身に次のように問いかけるべきです。

> `newPost` は `fs.File` に結合している必要があるだろうか？ この型のメソッドやデータをすべて使っているだろうか？ _本当に_ 必要なものは何だろうか？

私たちのケースでは、`io.Reader` を必要とする `io.ReadAll` の引数としてしか使っていません。したがって、関数の結合を緩めて、`io.Reader` を要求するようにすべきです。

```go
func newPost(postFile io.Reader) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

同様の議論は `getPost` 関数にも当てはまります。この関数は `fs.DirEntry` を引数に取りますが、ファイル名を取得するために `Name()` を呼び出しているだけです。そのすべては必要ありません。その型から切り離し、ファイル名を文字列として渡すようにしましょう。完全にリファクタリングしたコードは次のとおりです。

```go
func NewPostsFromFS(fileSystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(fileSystem, f.Name())
		if err != nil {
			return nil, err //todo: needs clarification, should we totally fail if one file fails? or just ignore?
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(fileSystem fs.FS, fileName string) (Post, error) {
	postFile, err := fileSystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile io.Reader) (Post, error) {
	postData, err := io.ReadAll(postFile)
	if err != nil {
		return Post{}, err
	}

	post := Post{Title: string(postData)[7:]}
	return post, nil
}
```

これからは、作業の大部分を `newPost` の中にきれいに収めることができます。ファイルを開いて反復処理するという関心事は片付いたので、`Post` 型のデータを抽出することに集中できます。技術的には必須ではありませんが、ファイルは関連するものを論理的にまとめる良い方法なので、私は `Post` 型と `newPost` を新しい `post.go` ファイルに移動しました。

### テストヘルパー

テストにも気を配るべきです。`Posts` に対してアサーションを何度も行うことになるので、それを助けるコードを書きましょう。

```go
func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
```

```go
assertPost(t, posts[0], blogposts.Post{Title: "Post 1"})
```

## 最初にテストを書く

テストをさらに拡張して、ファイルの次の行である説明（description）を抽出しましょう。テストを通すところまでは、もう慣れ親しんだ快適なものに感じられるはずです。

```go
func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1`
		secondBody = `Title: Post 2
Description: Description 2`
	)

	fs := fstest.MapFS{
		"hello world.md":  {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	// rest of test code cut for brevity
	assertPost(t, posts[0], blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
	})

}
```

## テストを実行してみます

```
./blogpost_test.go:47:58: unknown field 'Description' in struct literal of type blogposts.Post
```

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`Post` に新しいフィールドを追加します。

```go
type Post struct {
	Title       string
	Description string
}
```

これでテストはコンパイルされ、失敗するはずです。

```
=== RUN   TestNewBlogPosts
    blogpost_test.go:47: got {Title:Post 1
        Description: Description 1 Description:}, want {Title:Post 1 Description:Description 1}
```

## 成功させるのに十分なコードを書く

標準ライブラリには、データを 1 行ずつ走査するのに便利なライブラリ [`bufio.Scanner`](https://golang.org/pkg/bufio/#Scanner) があります。

> Scanner は、改行で区切られたテキスト行からなるファイルなどのデータを読み込むための便利なインターフェースを提供します。
>
> Scanner provides a convenient interface for reading data such as a file of newline-delimited lines of text.

```go
func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	scanner.Scan()
	titleLine := scanner.Text()

	scanner.Scan()
	descriptionLine := scanner.Text()

	return Post{Title: titleLine[7:], Description: descriptionLine[13:]}, nil
}
```

都合のよいことに、これも読み込み対象として `io.Reader` を受け取るので（ここでも疎結合に感謝です）、関数の引数を変更する必要はありません。

`Scan` を呼び出して 1 行読み込み、`Text` を使ってデータを取り出します。

この関数が `error` を返すことは決してありません。この時点で戻り値の型からエラーを取り除きたくなるかもしれませんが、後で不正なファイル構造を処理しなければならないことがわかっているので、そのままにしておいてもよいでしょう。

## リファクタリング♪

行を走査してからテキストを読み込むという処理に繰り返しがあります。この操作は少なくとももう一度行うことがわかっているので、DRY にするための簡単なリファクタリングから始めましょう。

```go
func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	title := readLine()[7:]
	description := readLine()[13:]

	return Post{Title: title, Description: description}, nil
}
```

コードの行数はほとんど減っていませんが、それがリファクタリングの目的であることはめったにありません。ここで私がやろうとしているのは、行を読み込むことの _何を_ と _どのように_ を分離して、読み手にとってコードを少し宣言的にすることだけです。

7 と 13 というマジックナンバーは役目を果たしてはいますが、あまり説明的ではありません。

```go
const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLine := func() string {
		scanner.Scan()
		return scanner.Text()
	}

	title := readLine()[len(titleSeparator):]
	description := readLine()[len(descriptionSeparator):]

	return Post{Title: title, Description: description}, nil
}
```

リファクタリングの創造的な頭でコードを眺めていると、`readLine` 関数にタグの除去も任せてみたくなりました。また、`strings.TrimPrefix` 関数を使えば、文字列からプレフィックスを取り除くのをもっと読みやすく書けます。

```go
func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Post{
		Title:       readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
	}, nil
}
```

このアイデアを気に入るかどうかはあなた次第ですが、私は気に入っています。重要なのは、リファクタリングの段階では内部の詳細を自由にいじることができ、テストを実行し続けることで物事が正しく動作しているかを確認できるということです。気に入らなければ、いつでも以前の状態に戻ることができます。TDD のアプローチは、アイデアを頻繁に試す許可を与えてくれるので、素晴らしいコードを書くチャンスが増えるのです。

次の要件は、記事のタグを抽出することです。もし一緒に進めているなら、読み進める前に自分で実装してみることをおすすめします。今ではもう、反復の良いリズムがつかめていて、次の行を抽出してデータを解析することに自信を持てているはずです。

簡潔にするため TDD のステップは省略しますが、タグを追加したテストは次のとおりです。

```go
func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker`
	)

	// rest of test code cut for brevity
	assertPost(t, posts[0], blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Tags:        []string{"tdd", "go"},
	})
}
```

私が書いたものをただコピー＆ペーストするだけでは、自分自身をごまかしているだけです。認識を合わせるために、タグの抽出を含む私のコードを示します。

```go
const (
	titleSeparator       = "Title: "
	descriptionSeparator = "Description: "
	tagsSeparator        = "Tags: "
)

func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Post{
		Title:       readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
		Tags:        strings.Split(readMetaLine(tagsSeparator), ", "),
	}, nil
}
```

特に驚くようなことはないと思います。`readMetaLine` を再利用してタグ用の次の行を取得し、`strings.Split` を使ってそれらを分割することができました。

ハッピーパスにおける最後の反復は、本文を抽出することです。

提案されたファイル形式をもう一度確認しておきましょう。

```markdown
Title: Hello, TDD world!
Description: First post on our wonderful blog
Tags: tdd, go
---
Hello world!

The body of posts starts after the `---`
```

最初の 3 行はすでに読み込みました。次にもう 1 行読み込んでそれを捨てれば、ファイルの残りの部分が記事の本文になります。

## 最初にテストを書く

テストデータを変更して、区切り線と、いくつかの改行を含む本文を持たせ、すべての内容を取得できているかを確認します。

```go
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
	)
```

他のものと同様に、アサーションに追加します。

```go
	assertPost(t, posts[0], blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Tags:        []string{"tdd", "go"},
		Body: `Hello
World`,
	})
```

## テストを実行してみます

```
./blogpost_test.go:60:3: unknown field 'Body' in struct literal of type blogposts.Post
```

予想どおりです。

## テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します

`Post` に `Body` を追加すると、テストは失敗するはずです。

```
=== RUN   TestNewBlogPosts
    blogposts_test.go:38: got {Title:Post 1 Description:Description 1 Tags:[tdd go] Body:}, want {Title:Post 1 Description:Description 1 Tags:[tdd go] Body:Hello
        World}
```

## 成功させるのに十分なコードを書く

1. 次の行を走査して、`---` の区切り線を無視します。
2. 走査するものがなくなるまで走査を続けます。

```go
func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	title := readMetaLine(titleSeparator)
	description := readMetaLine(descriptionSeparator)
	tags := strings.Split(readMetaLine(tagsSeparator), ", ")

	scanner.Scan() // ignore a line

	var b strings.Builder
	for scanner.Scan() {
		fmt.Fprintln(&b, scanner.Text())
	}
	body := strings.TrimSuffix(b.String(), "\n")

	return Post{
		Title:       title,
		Description: description,
		Tags:        tags,
		Body:        body,
	}, nil
}
```

* `scanner.Scan()` は、走査するデータがまだあるかどうかを示す `bool` を返すので、それを `for` ループと組み合わせて、最後までデータを読み続けることができます。
* `Scan()` のたびに、`fmt.Fprintln` を使ってデータをバッファに書き込みます。スキャナーは各行から改行を取り除いてしまいますが、私たちは改行を維持する必要があるため、改行を追加するバージョンを使います。
* 上記の理由により、末尾に余計な改行が残らないよう、最後の改行を取り除く必要があります。

## リファクタリング♪

残りのデータを取得するという考えを関数にカプセル化すると、将来の読み手は実装の詳細を気にすることなく、`newPost` で _何が_ 起きているのかをすばやく理解できるようになります。

```go
func newPost(postBody io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postBody)

	readMetaLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}

	return Post{
		Title:       readMetaLine(titleSeparator),
		Description: readMetaLine(descriptionSeparator),
		Tags:        strings.Split(readMetaLine(tagsSeparator), ", "),
		Body:        readBody(scanner),
	}, nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan() // ignore a line
	var b strings.Builder
	for scanner.Scan() {
		fmt.Fprintln(&b, scanner.Text())
	}
	return strings.TrimSuffix(b.String(), "\n")
}
```

## さらなる反復

ハッピーパスに到達するための最短ルートをとり、機能の「スチールスレッド（steel thread）」を作りました。しかし、本番環境で使えるようになるまでには、まだ明らかに距離があります。

次のことはまだ処理していません。

* ファイルの形式が正しくない場合
* ファイルが `.md` ではない場合
* メタデータのフィールドの順序が異なる場合はどうするか？ それは許可されるべきか？ 処理できるようにすべきか？

しかし重要なのは、動くソフトウェアがあり、インターフェースを定義できたということです。上記は単なるさらなる反復であり、振る舞いを駆動するために書くべきテストが増えるだけです。上記のいずれかをサポートするために _設計_ を変える必要はなく、実装の詳細を変えるだけで済むはずです。

ゴールに集中し続けたことで、全体の設計に影響しない事柄にはまり込むのではなく、重要な決定を下し、それを望ましい振る舞いに照らして検証することができました。

## まとめ

`fs.FS` をはじめとする Go 1.16 の変更により、ファイルシステムからデータを読み込み、それを簡単にテストするためのエレガントな方法が手に入りました。

コードを「本当に」試してみたい場合は、次のようにします。

* プロジェクト内に `cmd` フォルダを作成し、`main.go` ファイルを追加します
* 次のコードを追加します

```go
import (
	blogposts "github.com/quii/fstest-spike"
	"log"
	"os"
)

func main() {
	posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(posts)
}
```

* `posts` フォルダにマークダウンファイルをいくつか追加して、プログラムを実行してみましょう！

本番コード

```go
posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
```

とテスト

```go
posts, err := blogposts.NewPostsFromFS(fs)
```

の対称性に注目してください。

これこそ、利用者駆動のトップダウンな TDD が _正しいと感じられる_ ときです。

パッケージのユーザーは、テストを見れば、それが何をするためのもので、どのように使うのかをすぐに把握できます。メンテナーとしては、_テストが利用者の視点から書かれているので、テストが役に立つものだと自信を持てます_。実装の詳細やその他の付随的な詳細をテストしているわけではないので、リファクタリングの際にテストが邪魔になるのではなく、助けになるとかなりの自信を持てます。

[**依存性注入**](dependency-injection.md)のような優れたソフトウェアエンジニアリングのプラクティスに頼ることで、コードはテストも再利用も簡単になります。

パッケージを作成するときは、たとえそれがプロジェクト内部でのみ使われるものであっても、トップダウンの利用者駆動のアプローチを優先してください。これにより、設計を過剰に想像したり、必要ないかもしれない抽象化を作ったりするのを防ぎ、書くテストが役に立つものになることを確実にするのに役立ちます。

反復的なアプローチによってすべてのステップを小さく保つことができ、継続的なフィードバックによって、場当たり的な他のアプローチよりもおそらく早く、不明確な要件を明らかにすることができました。

### 書き込みは？

これらの新機能には、ファイルを _読み込む_ ための操作しかないことに注意してください。書き込みが必要な作業であれば、他を当たる必要があります。標準ライブラリが現在何を提供しているかについて考え続けることを忘れないでください。データを書き込むのであれば、コードを疎結合かつ再利用可能に保つために、おそらく `io.Writer` のような既存のインターフェースを活用することを検討すべきでしょう。

### 参考資料

* これは `io/fs` の軽い紹介でした。[Ben Congdon が素晴らしい記事を書いており](https://benjamincongdon.me/blog/2021/01/21/A-Tour-of-Go-116s-iofs-package/)、この章を書く上で大いに役立ちました。
* [ファイルシステムインターフェースに関する議論](https://github.com/golang/go/issues/41190)
