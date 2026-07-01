# govulncheckのGO-2026-4550を解消する

## 問題

2026-06-25のセキュリティ監査で、`github.com/cloudflare/circl@v1.6.2` の `GO-2026-4550` が診断パッケージ暗号化経路から到達可能として検出された。

## 期待する挙動

暗号関連依存関係を修正版へ更新し、`govulncheck ./...` が成功する。

## 受け入れ条件

- `github.com/cloudflare/circl` が `v1.6.3` 以上、可能なら `v1.6.4` 以上へ更新されている。
- `go test ./...` が成功する。
- 固定バージョンの `govulncheck` が成功する。
- 診断データ暗号化とzip引数暗号化が壊れていない。

## 対応内容

- `github.com/cloudflare/circl` を `v1.6.4` へ更新した。
- 関連する `golang.org/x/*` 依存関係を更新した。
- `go test ./...` と `govulncheck ./...` の成功を確認した。
