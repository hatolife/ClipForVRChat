# v0.1.8-b10 betaを作成する

## 指示

> 変更含めてbetaつくって

## 文脈

`v0.1.8-b9` 以降、`develop` には `/usercamera/Close` 送信停止とOSC送受信診断ログの永続化/JSON Lines化が入っている。既存betaタグは `v0.1.8-b9` まで作成済み。

## 解釈

現在の `develop` の変更を含めて次のbetaタグ `v0.1.8-b10` を作成し、GitHubへpushしてRelease workflowによる配布物作成を確認する。

## 問題

- `v0.1.8-b9` には直近の診断ログ改善と `/usercamera/Close` 停止が含まれていない。
- beta配布物を作成するには、タグを現在の変更込みのcommitへ付ける必要がある。

## 期待する挙動

- `v0.1.8-b10` タグが直近変更を含むcommitを指す。
- Release workflowが起動し、GitHub Releaseと規定の配布Assetが作成される。

## 受け入れ条件

- [ ] `v0.1.8-b10` に #243 と #255 の変更が含まれる。
- [ ] 関連チェックが通る。
- [ ] `develop` と `v0.1.8-b10` タグをGitHubへpushする。
- [ ] Release workflowが成功する。
- [ ] GitHub Releaseがprereleaseとして作成され、公開Assetが規定の4種類である。
