# release/v0.1.8 の分岐を解消する

## 問題

`release/v0.1.8` が `origin/release/v0.1.8` に対して `ahead 3, behind 2` となり、ローカル修正とremote上のreport追加/rc13 status記録が分岐している。

## 期待する挙動

共有済みブランチの履歴を書き換えず、remote側の変更を取り込み、現在のセキュリティ修正とレポート配置を維持した状態でpushできる。

## 受け入れ条件

- [ ] `origin/release/v0.1.8` の変更を取り込み、分岐が解消される。
- [ ] Codex Security finding本文はユーザー指定の `reports/security/2026-07-01T04-48-55.763Z/` 配下で保持する。
- [ ] 重複する誤配置ファイルが残らない。
- [ ] 署名付きコミットで整理される。
- [ ] 主要なローカル検証が成功する。
