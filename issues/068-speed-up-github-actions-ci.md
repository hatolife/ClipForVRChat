# GitHub Actions CIを高速化する

## 問題

GitHub ActionsのCI/Releaseで、Goセットアップ、Goテスト、Wails CLIインストールに時間がかかっている。特に `go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0` は毎回実行され、cache hit時に省略できる余地がある。

## 期待する挙動

- Wails CLIをActions cacheで再利用し、cache hit時はインストールを省略する。
- CI/Release workflowの動作は変えず、ビルドとRelease成果物の検証は維持する。
- 高速化後のActionsログで効果と残課題を確認できる。

## 受け入れ条件

- CI workflowでWails CLI cacheが使われる。
- Release workflowでWails CLI cacheが使われる。
- cache miss時も従来通りWails CLIをインストールできる。
- CI/Release workflowの構文検査が通る。
