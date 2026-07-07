# 自動撮影用Webhook URLが想定外にグレーアウトする

## 指示

> 自動撮影用Webhook URL
> 通常投稿とは別の投稿先にしたい場合だけ入力します。空の場合は通常投稿用Webhook URLへ投稿します。
>
> これが想定しないタイミングでグレーアウトする不具合がある
> '/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b11-windows-amd64'
> の設定とlogを参照 現状グレーアウトしている状態

## 文脈

`v0.1.8-b11` の配布フォルダで、自動撮影用Webhook URL欄が現在グレーアウトしている。
この欄は通常投稿とは別の投稿先にしたい場合だけ入力する任意設定であり、空欄時は通常投稿用Webhook URLへフォールバックする。

## 解釈

自動撮影用Webhook URL欄は、Discord投稿や自動撮影の実効状態に応じて編集可能であるべきだが、設定値またはランタイム状態の組み合わせで不要にdisabledになっている可能性がある。
指定された `config.json` とログを確認し、グレーアウト条件が仕様に合っているかを切り分ける。

## 問題

自動撮影用Webhook URL欄が、ユーザーの想定しないタイミングでグレーアウトし、専用Webhook URLを入力できない。

指定配布フォルダの現状設定では、通常のDiscord投稿はONで通常投稿用Webhook URLも入力済みだが、`autoCapture.discord.enabled=false` のため、フロントエンドの `!state.config.output.uploadDiscord || !state.config.autoCapture.discord.enabled` 条件により自動撮影用Webhook URL欄がdisabledになっていた。
ログでも `output.uploadDiscord` と通常投稿用Webhook URLの変更後、自動撮影スケジュールはONになっている一方、`autoCapture.discord.enabled` の変更は記録されていない。

## 期待する挙動

自動撮影用Webhook URL欄は、設定として編集すべき状態ではグレーアウトしない。
グレーアウトする場合は、親機能OFFなど利用者が納得できる明確な条件に限定される。

## 受け入れ条件

- [x] 指定配布フォルダの `config.json` とログから、現状のグレーアウト理由を説明できる。
- [x] 自動撮影用Webhook URL欄のdisabled条件を仕様に沿って修正する。
- [x] 通常投稿用Webhook URLへのフォールバック挙動を壊さない。
- [x] フロントエンドtemplate/API surfaceチェックを実行する。

## 確認

- `node scripts/check-frontend-template-literals.mjs`
- `node scripts/check-wails-api-surface.mjs`
