# Windows実機でACLを確認する

## 問題

Goの `0600` / `0700` 指定だけでは、Windows上の実際のDACLが期待通りか判断できない。

## 期待する挙動

Windows実機で `config.json`、`history.json`、`logs/`、`diagnostics/` のACLを確認し、想定外に広い権限がないことを判断できる。

## 受け入れ条件

- Windows実機で `icacls` を実行する。
- `Everyone` や想定外の広い `Users` 権限がないことを確認する。
- 確認結果を人間確認手順に反映する。

## 備考

これは人間確認が必要な作業として扱う。

## 現状

- `reports/2026-06-25/human-verification-guide.md` に確認手順と現時点の確認結果を記録済み。
- 最終判断はWindows実機での追加確認に委ねる。
