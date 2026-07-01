# Codex Security findingsを現在HEADで再検証し修正する

## 問題

Codex Securityが過去commitに対して作成したfindingが、現在HEADでも成立するか未確認のまま残っている。
finding作成後にRelease workflow、設定画面、自動投稿まわりの実装が進んでいるため、古いcommit前提で判断すると過剰修正または修正漏れが起きる。

## 期待する挙動

`reports/security/2026-07-01T04-48-55.763Z/codex-security-findings-2026-07-01T04-48-55.763Z.csv` を現在HEADで再検証し、未修正または部分修正の問題だけを安全側に修正する。
重複findingは統合し、Release workflow全体、imported configの同意導線、PGP検証説明を一貫した対策にする。

## 受け入れ条件

- [x] findingごとの現在HEAD判定、理由、修正方針を `reports/security/codex-security-remediation-2026-07-01.md` に記録する。
- [x] 外部/imported configから、保存・確認前にauto-photo/screenshot auto-post watcherが開始しない。
- [x] Discord webhookや自動投稿の重要設定を、デフォルトタブから見えないまま保存・有効化しにくい警告または確認導線にする。
- [x] Release workflowで未検証のtag/ref名をshellへ直接埋め込まない。
- [x] Release workflowでキャッシュ復元された `wails.exe` を信頼して実行しない。
- [x] Release workflowで `go run ...@latest` を使わず、権限とthird-party action利用を最小化する。
- [x] PGP検証説明で、Release同梱公開鍵だけを信頼根拠にしないことを明示する。
- [x] 既存テストとCI相当チェックを実行し、結果を報告書へ記録する。
