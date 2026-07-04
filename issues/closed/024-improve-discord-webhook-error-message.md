# Discord Webhook URLエラーの案内改善

## 問題

Webhook URLが空、不正、または無効なtokenの場合に、`status=401 body=...` のような技術的な表示になり、ユーザーが何をすればよいか分かりにくい。

## 期待する挙動

Webhook URLがおかしい可能性が高いことを伝え、DiscordでWebhook URLを再取得して設定し直すよう案内する。

## 受け入れ条件

- Webhook URL未設定時に再取得・再設定を促す。
- Webhook URL形式不正時に再取得・再設定を促す。
- Discord APIの401/404など、無効なWebhookを示す失敗時に再取得・再設定を促す。
- status/bodyの詳細は必要最小限に留め、ユーザー向け説明を先に表示する。
- `go test ./...` が成功する。
