# インポート設定由来の自動処理監視フォルダを確認・制限する

## 指示

> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/03a74009d25c8191ad9d77784977701a'を日本語で解説して

## 文脈

`reports/security/2026-07-07T17-53-16.203Z/03a74009d25c8191ad9d77784977701a` は、外部JSON設定から隠れた自動処理監視フォルダとWebhookを読み込むことで、利用者が意図しない画像をDiscordへ自動投稿できると指摘している。
解説調査で、現在HEADにも主要な経路が残っていることを確認した。

## 解釈

このチケットは、findingの解説ではなく実装対応を追跡する。
外部設定由来の自動処理設定を、利用者が保存前に確認・修正でき、意図しない監視フォルダとWebhookの組み合わせがそのまま有効にならない状態にする。

## 問題

- `autoPhoto.photoDirectory` / `screenshotAutoPost.screenshotDirectory` はJSON設定として保存・読み込みされる。
- 設定画面では自動処理の監視フォルダが通常表示されず、インポート設定の内容を確認しにくい。
- 保存時確認はWebhook未設定時に限られ、攻撃者がWebhook URLも設定したJSONでは表示されない。
- 保存後に `restartAutoPhotoWatcher` が設定由来の監視フォルダとWebhookでwatcherを起動する。

## 期待する挙動

- 外部設定由来の自動処理監視フォルダとWebhookが保存前に明示表示される。
- 利用者が確認しないまま、隠れた監視フォルダからDiscord投稿が有効にならない。
- 既定のVRChat写真フォルダやスクリーンショットフォルダの通常利用は維持される。

## 受け入れ条件

- [x] 自動処理ONかつDiscord投稿ONの場合、監視フォルダと有効Webhookを保存前確認に含める。
- [x] Webhookが設定済みの場合でも、外部設定由来または非表示フィールド由来の監視フォルダを確認対象にする。
- [x] 自動処理の監視フォルダを設定画面から確認・変更できる。
- [x] 広すぎる監視フォルダや危険な既知フォルダを安全側に扱う方針を実装または明文化する。
- [x] 回帰テストを追加する。

## 対応内容

- 自動写真処理またはスクリーンショット自動処理がONで、Discord投稿もONの場合は、Webhook URLが設定済みでも保存前確認を表示するようにした。
- 保存前確認に監視フォルダ、有効な送信先Webhook、広い既知フォルダの警告を表示するようにした。
- 機能設定タブで `autoPhoto.photoDirectory` と `screenshotAutoPost.screenshotDirectory` を確認・変更できるようにした。
- Desktop、Downloads、Documents、Pictures、Videos、Music、OneDrive、ユーザーフォルダ、ドライブ直下、ルートフォルダを広すぎる監視先として警告するようにした。
- `scripts/check-auto-processing-confirmation.mjs` を追加し、保存前確認条件と監視フォルダUIが消えた場合に検出できるようにした。
