# 監査報告書を日付付きで専用ディレクトリへ移動する

## 問題

`SECURITY_REVIEW.md` と `PRODUCT_ISSUE_REPORT.md` は作成時点の監査結果であり、時間が経つと現行実装とずれる。リポジトリ直下に固定名で置くと、最新状態の報告書であるように見えやすい。

## 期待する挙動

監査報告書を専用ディレクトリへ配置し、作成日ごとのディレクトリに格納することで、いつ時点の報告書か分かるようにする。

## 受け入れ条件

- 監査報告書専用ディレクトリが作成されている。
- `SECURITY_REVIEW.md` と `PRODUCT_ISSUE_REPORT.md` が日付ディレクトリへ移動されている。
- 既存のチケットや一覧から新しい配置へ参照できる。

## 対応内容

- `reports/README.md`
- `reports/2026-06-21/security-review.md`
- `reports/2026-06-24/product-issue-report.md`
