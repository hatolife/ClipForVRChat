# セキュリティ監査報告書を作成する

## 問題

`reports/security_audit_prompt.md` に基づく、現行リポジトリ全体の体系的なセキュリティレビュー結果が未作成である。

## 期待する挙動

プロジェクト構成、脅威モデル、環境固有・言語固有の確認、依存関係とサプライチェーン、発見事項、リリース前チェックリストを `reports/` 配下の報告書として確認できる。

## 受け入れ条件

- `reports/security_audit_prompt.md` の指定成果物を作成する。
- 実施した確認作業、読んだ主要ファイル、実行コマンド、未確認事項を記録する。
- 発見事項は重大度、影響、推奨修正、確認方法を含めて整理する。
- `findings.json` を機械処理しやすい JSON として作成する。

## 対応内容

- `reports/project-profile.md` ほか、指定された監査報告書一式を作成した。
- `govulncheck`、`gosec`、`go test`、`npm audit`、`npm run build` の実行結果を `reports/review-log.md` に記録した。
- `findings.json` にSEC-001からSEC-009までの発見事項をJSON形式で整理した。
