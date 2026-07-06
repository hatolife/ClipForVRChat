# CIのWails Build application工程が長時間完了しない

## 問題

`release/v0.1.8` のGitHub Actions CIで、Windows jobの `Build application` が長時間 `in_progress` のまま完了しない。

`v0.1.8-a31` だけでなく、`v0.1.8-a30` 作成時のbranch CIでも同じ `Build application` 工程で止まっているため、a31で追加した変更だけが原因とは限らない。

`v0.1.8-a30` Release workflow自体はGo test cleanup失敗で停止しており、このWails build長時間化とは別問題として扱う。

## 期待する挙動

Windows CI/Release workflowの `wails build` が通常時間内に完了し、失敗時もログから原因を確認できる。

## 受け入れ条件

- [x] CI/Releaseの `Build application` が長時間無出力で止まる原因候補を特定する。
- [x] Wails buildが正常完了する、または失敗ログを取得できるようにtimeout/verbose/診断を追加する。
- [x] `v0.1.8-a31` 以降のRelease workflowがBuild application工程で止まらない。
- [x] a30のGo test cleanup失敗とは別問題として記録し、混同しない。

## 実装メモ

- Wails v2.12.0のWindows buildは `-webview2` でWebView2 Runtime未検出時の扱いを選べる。
- 公式ドキュメント上、`download` はアプリ実行時に公式bootstrapperのダウンロード/実行を提示し、`browser` は公式WebView2ページをブラウザで開く方式。
- CI上の `wails build` が長時間無出力で止まっているため、bootstrapper処理や外部取得の影響を避ける目的で `-webview2 browser` を明示する。
- ローカル検証で `wails build` が `go.mod` を同期する副作用を出したため、CIでは事前の `go test` / `govulncheck` を信頼し、Wails build側は `-m -nosyncgomod` でmod同期/整理を行わない。
- a33の失敗ログで、20分timeout時点の停止箇所が `Generating bindings` と判明した。フロントエンドは生成bindingsをimportせず `window.go.main.App` を直接使い、CIでは `check-wails-api-surface` も実行しているため、Wails build側は `-skipbindings` でWindows上のbindings生成を避ける。
- `wails build` に `-v 2` を付け、CI stepへ `timeout-minutes: 20` を追加して、再発時に無期限進行中にならないようにする。

## 検証観点

- GitHub Actions run:
  - a30 branch CI: `28690773566`
  - a31 branch CI: `28691252927`
  - a31 Release: `28691253567`
- [x] ローカル `wails build` はbindings/frontend生成まで進み、Linuxの `webkit2gtk-4.0` 不足で明示的に失敗するため、bindings生成の停止ではない。
- [x] CIで `wails build -v 2 -m -nosyncgomod -skipbindings -webview2 browser -tags embeddedspout ...` が完了するか。

## 検証結果

- `v0.1.8-a33` のCI/Releaseは `Generating bindings` で20分timeoutし、停止箇所を特定できた。
- `v0.1.8-a34` のbranch CIとRelease workflowは成功した。
- Release asset一覧と通常zip/separated zip/AvatarBeacon source zipの中身を確認した。
