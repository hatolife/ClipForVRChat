# CI/Release のセキュリティチェック追加

## 問題

CI/Release の Go バージョンが `go.mod` と揃っておらず、`govulncheck` や `npm audit` がリリース前に実行されていない。Release zip の改ざん確認用チェックサムもない。

## 期待する挙動

CI/Release で依存関係の既知脆弱性を確認し、Release artifact のチェックサムを生成する。

## 受け入れ条件

- CI/Release の Go バージョンが `go.mod` と一致する。
- CI/Release で `govulncheck ./...` を実行する。
- CI/Release で `npm audit --omit=dev` を実行する。
- Release に zip と SHA256 チェックサムを添付する。

## 追記: govulncheck 実行バージョンの固定

## 問題

CI/Release workflow の `govulncheck` が未固定バージョンで実行されると、workflow 実行時に解決された未固定の Go モジュールコードが実行される。

## 期待する挙動

CI/Release workflow はレビュー済みの固定バージョンの `govulncheck` を実行し、将来公開される未確認コードを自動実行しない。

## 受け入れ条件

- CI workflow の `govulncheck` 実行が `@latest` を使わない。
- Release workflow の `govulncheck` 実行が `@latest` を使わない。
- 使用する `govulncheck` バージョンが workflow 内で明示されている。
