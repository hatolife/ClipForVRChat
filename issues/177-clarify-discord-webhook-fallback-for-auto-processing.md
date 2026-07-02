# 自動処理Webhookの通常投稿フォールバックを明確にする

## 問題

VRChat写真自動処理、スクリーンショット自動処理、自動撮影用のWebhook URLが空欄のとき、通常投稿用Webhook URLへフォールバックされるのか、保存時確認ダイアログ上で分かりにくい。
通常投稿用Webhook URLが設定済みでも、自動処理の保存確認ダイアログが出るため、専用Webhook未設定が危険な状態に見える。

## 期待する挙動

通常投稿用Webhook URLが設定されている場合、VRChat写真自動処理、スクリーンショット自動処理、自動撮影の専用Webhook URLが空欄でも通常投稿用Webhook URLへ送信されることがUI上で分かる。
保存時確認ダイアログは、Discord投稿ONかつ有効な自動処理の実効Webhookが未設定の場合だけ表示する。

## 受け入れ条件

- [x] VRChat写真用Webhook URL欄で、空欄時は通常投稿用Webhook URLを使うことが表示される。
- [x] スクリーンショット用Webhook URL欄で、空欄時は通常投稿用Webhook URLを使うことが表示される。
- [x] 自動撮影用Webhook URL欄で、空欄時は通常投稿用Webhook URLを使うことが表示される。
- [x] 通常投稿用Webhook URLが設定済みなら、専用Webhook空欄だけを理由に保存時確認ダイアログを出さない。
- [x] 通常投稿用Webhook URLも専用Webhook URLも空欄で、Discord投稿ONかつ対象自動処理が有効な場合は、送信先未設定として確認ダイアログを出す。
- [x] 実際のVRChat写真/スクリーンショット/自動撮影の投稿処理が通常投稿用Webhook URLへフォールバックする。

## メモ

- 2026-07-02: VRChat写真/スクリーンショットの自動処理は、専用Webhookが空欄なら `Config.Discord.WebhookURL` が通常投稿用として使われる既存挙動を確認した。
- 2026-07-02: 自動撮影は `AutoCaptureRunner.finalizeAutoCaptureImage` で自動撮影用Webhookが空欄の場合に通常投稿用Webhookへフォールバックする既存挙動を確認した。
- 2026-07-02: フロントエンドの保存前確認を、実効Webhook URLが未設定の自動処理だけに限定した。
