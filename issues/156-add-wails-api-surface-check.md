# Wails公開APIとフロント呼び出しの同期チェックを追加する

## 問題

フロントエンドは `window.go?.main?.App` からWails公開APIを直接呼び出している。
`src/frontend/wailsjs/` はgit管理外の生成物であり、公開メソッド追加後に生成APIや型情報の同期漏れがあっても、通常の差分確認では気づきにくい。
実機で「APIが利用できません」と表示されると、生成漏れ、ビルド漏れ、ランタイム問題の切り分けが難しい。

## 期待する挙動

CIまたはローカル検証で、フロントエンドが呼び出すWails APIがGo側に存在し、Wailsビルドで公開されることを確認できる。

## 受け入れ条件

- フロントエンドが参照する `api.*` メソッド一覧を抽出する検査を追加する。
- Go側の `App` 公開メソッド一覧と照合し、存在しないAPI参照をCIで失敗させる。
- 必要に応じて `wails generate` 相当の生成確認、またはWails build後の生成物確認をCIに追加する。
- `autoCapture` 設定型のような主要モデルがフロントで扱えることを検査する。
- API未公開時の画面エラーはユーザー向け原因と開発向け診断ログを分ける。
- 生成物をgit管理しない方針を維持する場合、その理由と検証方法をドキュメント化する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- Wails生成物をgit管理しない方針は維持し、静的検査スクリプトでfrontendの `api.*` 呼び出しとGo `App` 公開メソッドを照合する。
- Nodeスクリプト `scripts/check-wails-api-surface.mjs` を追加する。
- CIではfrontend build前後どちらでもよいが、GoメソッドとJS参照の不一致を早く検出するため、`npm run build` と `go test ./...` の間に実行する。

対象ファイル:

- `scripts/check-wails-api-surface.mjs` 新設
- `.github/workflows/ci.yml`
- `src/frontend/package.json`
- `src/app.go`
- `src/frontend/src/main.js`

小タスク:

- `src/frontend/src/main.js` から `api.<Identifier>` を正規表現で抽出する。
- `src/app.go` から `func (a *App) <ExportedName>(` を抽出する。
- frontend参照に存在してGo公開メソッドに存在しないものがあれば失敗する。
- 逆方向の未使用Goメソッドは失敗ではなく情報表示にする。
- `api?.Foo`、`api.Foo` の両方を検出する。
- `autoCapture` configの主要キーは型生成ではなく実データJSONとして扱っているため、このissueではAPI名検査を必須範囲にする。

確認方法:

- 存在しない `api.DoesNotExist` を一時的に追加するとスクリプトが失敗する。
- 現状の `main.js` と `app.go` ではスクリプトが成功する。
- Windows CIで同じ検査が走る。
