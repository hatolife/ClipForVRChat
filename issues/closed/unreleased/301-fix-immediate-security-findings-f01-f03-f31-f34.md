# F01/F03/F31/F32/F33/F34 の即時修正

## 指示

> Task: fix immediate no-spec findings F01, F03, F31/F32, F33/F34 from reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md.
> - Preserve explicit RemotePlayer=false in AutoCaptureConfig.Normalize; only default nil values.
> - Count skipped auto-photo files toward MaxAutoPhotoProcessPerTick after stability check so skip paths cannot bypass the per-tick cap.
> - Suppress repeated identical auto-photo scan error/status events so persistent missing/limit errors do not grow state/UI unbounded.
> - Ensure screenshot auto-post does not fall back to AutoPhoto.WebhookURL when ScreenshotAutoPost.WebhookURL is empty; it should use the normal Discord webhook or no override.
> Add focused Go tests. Run relevant go tests from src if feasible. Edit files directly and report changed files/tests.

## 文脈

`reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md` で即時修正対象に分類された low/medium finding を、仕様変更を伴わない範囲でまとめて修正する。

## 解釈

このチケットでは、設定Normalizeの既定値補完、auto-photo watcherの処理上限、scan statusの重複抑制、スクリーンショット自動投稿のWebhook解決を修正対象とする。

## 問題

- `RemotePlayer=false` の明示設定が Normalize で既定値へ戻る可能性がある。
- auto-photo の skip パスが per-tick 上限を迂回できる。
- 同一の scan error/status が毎 tick 追加され、state/UI が増え続ける。
- screenshot auto-post が空の `ScreenshotAutoPost.WebhookURL` から `AutoPhoto.WebhookURL` へ誤ってフォールバックしうる。

## 期待する挙動

- 明示された `RemotePlayer=false` は維持される。
- auto-photo は stable 判定後の skip も含めて per-tick 上限を消費する。
- 同じ scan error/status は重複して通知されない。
- screenshot auto-post は通常の Discord webhook を使うか、明示 override がなければ上書きしない。

## 受け入れ条件

- [ ] `AutoCaptureConfig.Normalize` で `RemotePlayer=false` が保持される。
- [ ] skip された auto-photo ファイルも per-tick 上限に数えられる。
- [ ] 同一の auto-photo scan error/status が重複追加されない。
- [ ] screenshot auto-post が `AutoPhoto.WebhookURL` にフォールバックしない。
- [ ] 影響箇所に対応する Go テストが追加される。
