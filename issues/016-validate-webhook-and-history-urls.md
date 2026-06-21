# Webhook URL と履歴 URL の検証強化

## 問題

Discord Webhook URL が Discord ドメインに制限されておらず、設定ファイル経由で任意サーバーへ画像を送信できる。履歴 URL も任意 URL のまま表示・到達性確認されるため、改ざんされた履歴から外部またはローカルネットワークへアクセスする可能性がある。

## 期待する挙動

Discord Webhook として妥当な HTTPS URL だけを投稿・削除に使う。履歴に保存、表示、到達性確認する画像 URL も信頼できる HTTPS URL に限定する。

## 受け入れ条件

- Discord Webhook URL は `discord.com` または `discordapp.com` の `/api/webhooks/{id}/{token}` 形式だけ許可する。
- 履歴画像 URL は HTTPS の Discord/CDN 系 URL だけ到達性確認と表示対象にする。
- 不正 URL の場合はユーザーに分かるエラーを返す。
- URL 検証の単体テストがある。
