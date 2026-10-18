# Learn Go with Tests 日本語版 — 翻訳・保守ガイド

このリポジトリは [quii/learn-go-with-tests](https://github.com/quii/learn-go-with-tests)（以下「原文」）の日本語訳フォークです。作業の中心は「原文への追従」と「訳文の品質維持」であり、本の内容そのものを書き換えることではありません。

## リポジトリ構成

- 章の Markdown は原文と違い、サブディレクトリに置いている（原文はルート直下）。

  | ディレクトリ | 章 |
  |---|---|
  | `go-fundamentals/` | install-go, hello-world, integers, iteration, arrays-and-slices, structs-methods-and-interfaces, pointers-and-errors, maps, dependency-injection, mocking, concurrency, select, reflection, sync, context, roman-numerals, math, reading-files, html-templates, generics, revisiting-arrays-and-slices-with-generics |
  | `testing-fundamentals/` | intro-to-acceptance-tests, scaling-acceptance-tests, working-without-mocks, refactoring-checklist |
  | `build-an-application/` | app-intro, http-server, json, io, command-line, time, revisiting-time-with-synctest, websockets |
  | `questions-and-answers/` | os-exec, error-types, context-aware-reader, http-handlers-revisited |
  | `meta/` | why, anti-patterns, contributing |
  | ルート | `gb-readme.md`（GitBook の表紙）, `SUMMARY.md`（GitBook の目次）, `README.md` |

- 画像は `.gitbook/assets/` に置き、章からは `../.gitbook/assets/<file>` で参照する。
- Go のコード（各章のディレクトリ、`go.mod` など）は **原文と同一に保つ**。翻訳側で独自にコードを変更しない。
- 新しい章を追加したら `SUMMARY.md`、`README.md` の目次、`build.books.sh`（PDF/EPUB の章リスト）、`build.sh`（対象ディレクトリ）を更新する。

## 原文への追従手順

1. `git remote add upstream https://github.com/quii/learn-go-with-tests.git`（未設定なら）→ `git fetch upstream`
2. 作業ブランチで `git merge upstream/main`。
   - Go コード・`go.mod`・`go.sum`・アセットは upstream 側を採用する。
   - 原文の章 `.md` はルート直下に復活してくる（modify/delete コンフリクト）。英語版の内容を参照用に退避したうえで削除し、対応する日本語版（サブディレクトリ側）に差分を反映する。
3. 既存章は、前回取り込み時点からの英語の差分（`git diff <前回のマージベース> upstream/main -- <章>.md`）を見て、変わった箇所だけを訳文に反映する。変わっていない段落の既存訳は無理に書き換えない。
4. 新しい章は全文を翻訳し、上の配置表に従って置く。
5. 検証: 各章の見出し数とコードブロック（```）数が原文と一致すること、相対リンク・画像のリンク切れがないこと、CI（`build.sh` / `build.books.sh`）が通ること。

## 翻訳スタイル

- です・ます調。原文の軽い口調（"Yuck." "Perfect!" などの一言）は、自然な日本語の一言で表す。
- 直訳調・機械翻訳調を避け、意味が通る自然な日本語にする。ただし原文にない内容は足さない。
- コードブロック内（コード、コンパイラ出力、テスト出力）は原文のまま。コード中のコメントも訳さない。
- 外部文献の引用（`>`）は訳す。必要なら原文を併記してよい。
- 章間リンクは配置表に合わせた相対パス（例: `../build-an-application/http-server.md`）。「この章のすべてのコードはここにあります」のリンクは原文と同じ `https://github.com/quii/learn-go-with-tests/tree/main/...` を使う。
- 各章の先頭には GitBook 用の front matter（`description: <日本語の章タイトル>`）を付ける。

### TDD の定型見出し

| 原文 | 訳 |
|---|---|
| Write the test first | 最初にテストを書く |
| Try to run the test | テストを実行してみます |
| Write the minimal amount of code for the test to run and check the failing test output | テストを実行するための最小限のコードを記述し、失敗したテスト出力を確認します |
| Write enough code to make it pass | 成功させるのに十分なコードを書く |
| Refactor | リファクタリング♪ |
| Wrapping up | まとめ |

### 用語

本文中の表記を統一する（コード・コメントは対象外）。

インターフェース / ゴルーチン / チャネル / プレーヤー / ポインタ / パラメータ / サーバー / ハンドラー / コンパイラ / ヘルパー / アサーション / メソッド / スライス / マップ / 構造体 / 埋め込み / 依存性注入 / モック / スパイ / スタブ / Approval Tests

## コミット・PR

- コミットメッセージと PR 本文に `Co-Authored-By: Claude ...` や「Generated with Claude Code」の行を入れない。
