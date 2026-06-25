# 人間が確認する必要がある作業を手順化する

## 問題

セキュリティ監査後に、人間が確認すべき作業と自動検証で十分な作業が混在している。Windows実機、Discord実通信、GitHub Release運用など、環境依存の確認手順が一覧化されていない。

## 期待する挙動

人間が確認する必要がある作業、確認タイミング、手順、判断基準がMarkdownで確認できる。

## 受け入れ条件

- 人間確認用の手順書を `reports/2026-06-25/` に追加する。
- `GO-2026-4550` のように主作業は実装・自動検証で見る項目と、人間が見るべき項目を分ける。
- Windows ACL、CLI出力、診断zip、Webhook token保存方針について現時点の確認結果と判断を記録する。
- `reports/README.md` から参照できる。

## 対応内容

- `reports/2026-06-25/human-verification-guide.md` を追加した。
- `GO-2026-4550` は自動検証主体、暗号化/復号の実機互換性は人間確認対象として分離した。
- Windows ACL、CLI出力、診断zip確認タイミング、Webhook token保存方針を記録した。
- `reports/README.md` と `issues/README.md` へ参照を追加した。
