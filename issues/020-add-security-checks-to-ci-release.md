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
