# 監査報告書の配置と説明を整理する

## 問題

監査報告書が `reports/` 直下に混在しており、日付単位のスナップショットとして見通しにくい。また、診断データの平文zip、Release workflowの権限、Windows ACLについて説明が不足している。

## 期待する挙動

報告書は作成日ごとのディレクトリへ格納され、監査時点の成果物をまとめて確認できる。平文zipの扱い、Webhook tokenの除外方針、Release workflowの権限/環境保護、Windows ACLの意味が報告書上で理解できる。

## 受け入れ条件

- 2026-06-25監査報告書を `reports/2026-06-25/` に移動する。
- 既存の2026-06-21、2026-06-24報告書も日付ディレクトリへ移動する。
- `reports/README.md` と関連issueの参照パスを更新する。
- 平文zipはユーザー確認用として残す前提にしつつ、Webhook tokenなどをzipに含めない方針を報告書へ反映する。
- Release workflowの権限/環境保護とWindows ACLについて、追加説明を報告書へ記載する。

## 対応内容

- 報告書を作成日ごとの `reports/YYYY-MM-DD/` ディレクトリへ移動した。
- `reports/README.md` と関連issue内の参照パスを移動後の配置へ更新した。
- セキュリティ監査報告書に、確認用平文zipを残す理由、Webhook tokenのダミー化方針、Release workflowの権限/環境保護、Windows ACL/DACLの説明を追記した。
