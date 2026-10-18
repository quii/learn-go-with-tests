# テスト駆動開発でGo言語を学びましょう（Learn Go with Tests 日本語版）

<p align="center">
  <img src="https://raw.githubusercontent.com/quii/learn-go-with-tests/main/red-green-blue-gophers-smaller.png" />
</p>

[Art by Denise](https://github.com/deniseyu)

このリポジトリは、Chris James（[@quii](https://github.com/quii)）による [Learn Go with Tests](https://github.com/quii/learn-go-with-tests) の日本語訳です。

## 読む

- [GitBook（日本語版）](https://andmorefine.gitbook.io/learn-go-with-tests/)
- [原文（英語版）](https://quii.gitbook.io/learn-go-with-tests)
- [原文の EPUB / PDF](https://github.com/quii/learn-go-with-tests/releases)

その他の翻訳: [中文](https://studygolang.gitbook.io/learn-go-with-tests) / [Português](https://larien.gitbook.io/aprenda-go-com-testes/) / [Français](https://goosegeesejeez.gitbook.io/apprendre-go-par-les-tests) / [한국어](https://miryang.gitbook.io/learn-go-with-tests/) / [Türkçe](https://halilkocaoz.gitbook.io/go-programlama-dilini-ogren/) / [فارسی](https://go-yaad-begir.gitbook.io/go-ba-test/) / [Nederlands](https://bobkosse.gitbook.io/leer-go-met-tests/) / [Tiếng Việt](https://sons-organization-15.gitbook.io/learn-go-with-tests) / [Русский](https://eda-1.gitbook.io/lgwt/)

## この本について

* テストを書きながらGo言語を探求します
* **TDD（テスト駆動開発）の基礎を身につけます。** Goはシンプルで学びやすく、テストの仕組みが言語に組み込まれているため、TDDを学ぶのに適した言語です
* 堅牢で、よくテストされたシステムをGoで書き始められるという自信をつけます
* [なぜユニットテストとTDDが重要なのか（動画・記事）](meta/why.md)

本の背景や対象読者は[表紙のページ](gb-readme.md)を参照してください。

## 目次

### Go言語：基礎

1. [Goをインストールする](go-fundamentals/install-go.md) - 生産的に作業するための環境を整えます。
2. [Hello, World](go-fundamentals/hello-world.md) - 変数、定数、if/else、switch を使い、最初のGoプログラムと最初のテストを書きます。サブテストの書き方とクロージャも学びます。
3. [整数](go-fundamentals/integers.md) - 関数宣言の構文をさらに掘り下げ、コードのドキュメントを充実させる方法を学びます。
4. [反復、繰り返し](go-fundamentals/iteration.md) - `for` とベンチマークについて学びます。
5. [配列とスライス](go-fundamentals/arrays-and-slices.md) - 配列、スライス、`len`、可変長引数、`range`、テストカバレッジについて学びます。
6. [構造体、メソッド、インターフェース](go-fundamentals/structs-methods-and-interfaces.md) - `struct`、メソッド、`interface`、テーブル駆動テストについて学びます。
7. [ポインタとエラー](go-fundamentals/pointers-and-errors.md) - ポインタとエラーについて学びます。
8. [マップ](go-fundamentals/maps.md) - マップというデータ構造に値を保存する方法を学びます。
9. [依存性注入](go-fundamentals/dependency-injection.md) - 依存性注入とインターフェースの関係、そして io パッケージの入門です。
10. [スタブ・モック](go-fundamentals/mocking.md) - テストのない既存のコードに、依存性注入とモックを使ってテストを書きます。
11. [並行性](go-fundamentals/concurrency.md) - 並行処理のコードを書いてソフトウェアを速くする方法を学びます。
12. [Select](go-fundamentals/select.md) - 非同期処理をエレガントに同期させる方法を学びます。
13. [リフレクション](go-fundamentals/reflection.md) - リフレクションについて学びます。
14. [同期](go-fundamentals/sync.md) - `WaitGroup` や `Mutex` など、sync パッケージの機能を学びます。
15. [コンテキスト](go-fundamentals/context.md) - context パッケージを使って、長時間実行される処理を管理・キャンセルします。
16. [プロパティベースのテスト概要](go-fundamentals/roman-numerals.md) - ローマ数字の型（カタ）でTDDを練習し、プロパティベーステストを簡単に紹介します。
17. [数学](go-fundamentals/math.md) - `math` パッケージを使ってSVGの時計を描きます。
18. [ファイルの読み込み](go-fundamentals/reading-files.md) - ファイルを読み込んで処理します。
19. [テンプレート](go-fundamentals/html-templates.md) - html/template パッケージでデータからHTMLを生成し、Approval Tests についても学びます。
20. [ジェネリクス](go-fundamentals/generics.md) - ジェネリックな引数を取る関数や、独自のジェネリックなデータ構造の書き方を学びます。
21. [ジェネリクスで配列とスライスを再考する](go-fundamentals/revisiting-arrays-and-slices-with-generics.md) - 独自の `Reduce` 関数を書き、コレクションを扱うよくあるパターンを整理します。

### テストの基礎

* [受け入れテスト入門](testing-fundamentals/intro-to-acceptance-tests.md) - HTTPサーバーのグレースフルシャットダウンを題材に、受け入れテストの書き方を学びます。
* [受け入れテストのスケーリング](testing-fundamentals/scaling-acceptance-tests.md) - 大きなシステムで受け入れテストの複雑さを管理する手法を学びます。
* [モックを使わずに作業する](testing-fundamentals/working-without-mocks.md) - フェイクと契約（コントラクト）テストで、より現実的で保守しやすいテストを書く方法を学びます。
* [リファクタリングチェックリスト](testing-fundamentals/refactoring-checklist.md) - リファクタリングとは何か、どう進めるかの基本的なヒントです。

### アプリケーションを構築する

* [はじめに](build-an-application/app-intro.md)
* [HTTPサーバー](build-an-application/http-server.md) - HTTPリクエストを受け取り、応答するアプリケーションを作ります。
* [JSON、ルーティング、埋め込み](build-an-application/json.md) - エンドポイントがJSONを返すようにし、ルーティングの方法を学びます。
* [IO、並び替え](build-an-application/io.md) - データをディスクに永続化して読み込み、並び替えを行います。
* [コマンドライン、パッケージ構造](build-an-application/command-line.md) - 1つのコードベースで複数のアプリケーションをサポートし、コマンドラインから入力を読み取ります。
* [時間](build-an-application/time.md) - `time` パッケージを使って処理をスケジュールします。
* [testing/synctest で時間を再考する](build-an-application/revisiting-time-with-synctest.md) - Go 1.25 の `testing/synctest` パッケージで、タイマーを使うコードを遅くも不安定にもならずにテストします。
* [ウェブソケット](build-an-application/websockets.md) - WebSocket を使うサーバーの書き方とテスト方法を学びます。

### 質問と答え

* [OS実行](questions-and-answers/os-exec.md) - OSのコマンドを実行してデータを取得しつつ、ビジネスロジックをテスト可能に保つ例です。
* [エラーの種類](questions-and-answers/error-types.md) - 独自のエラー型を作ってテストを改善し、コードを扱いやすくする例です。
* [コンテキスト認識リーダー](questions-and-answers/context-aware-reader.md) - `io.Reader` にキャンセル機能を加える処理をTDDで作ります。
* [HTTPハンドラーの再検討](questions-and-answers/http-handlers-revisited.md) - HTTPハンドラーを正しく設計するための論点を掘り下げます。

### メタ

* [なぜユニットテストなのか、そしてユニットテストを役立てる方法](meta/why.md) - 動画を見るか、記事を読んで、ユニットテストとTDDがなぜ重要なのかを知ります。
* [アンチパターン](meta/anti-patterns.md) - TDDとユニットテストのアンチパターンを紹介する短い章です。
* [貢献](meta/contributing.md)

## 翻訳への貢献

* 誤訳・誤字の指摘や改善の提案は、このリポジトリに issue または PR をお願いします。
* 原文の内容そのものへの指摘は、[原文のリポジトリ](https://github.com/quii/learn-go-with-tests)へお願いします。
* 各章のサンプルコードは原文と同一に保っています（コードの変更は原文側で行われたものを取り込みます）。

## 原作者を応援する

* [Twitter で @quii にツイートする](https://twitter.com/quii)
* <a rel="me" href="https://mastodon.cloud/@quii">Mastodon</a>
* [コーヒーをおごる :coffee:](https://www.buymeacoffee.com/quii)
* [GitHub でスポンサーになる](https://github.com/sponsors/quii)

[MIT license](LICENSE.md)

[Logo is by egonelbre](https://github.com/egonelbre) What a star!
