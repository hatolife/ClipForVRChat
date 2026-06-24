# GitHub ActionsのNode 20警告とGo cache警告を解消する

## 問題

GitHub Actionsで、利用中のactionがNode.js 20実行系を使っている警告と、Go cache restore時にリポジトリルートの `go.sum` が見つからない警告が出ている。

## 期待する挙動

- GitHub Actionsの主要actionをNode.js 24対応版へ更新する。
- Go module cacheの依存ファイルとして `src/go.sum` を明示する。
- CI/Release workflowの警告を減らし、正式リリース前の確認を安定させる。

## 受け入れ条件

- CI workflowでNode.js 20 deprecation警告とGo cache dependency警告が出ない。
- Release workflowでNode.js 20 deprecation警告とGo cache dependency警告が出ない。
- CI/Release workflowの構文検査が通る。
