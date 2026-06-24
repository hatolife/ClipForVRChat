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

## 確認結果

- ブランチCIの初回runでWails CLI cache miss、`go install` 実行、cache保存を確認した。
- ブランチCIの次回runでWails CLI cache hit、`Install Wails` step skip、`wails version` とアプリビルド成功を確認した。
- feature branch上の一時タグReleaseではWails CLI cache missとなったが、`go install`、署名、zip検査、draft Release作成は成功した。
- タグrunのcache hitは、default branchへ反映後にdefault branch cacheが作成された状態で再確認する。
