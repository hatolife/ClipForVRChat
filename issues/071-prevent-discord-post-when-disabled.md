# 071 Discord投稿OFFでも自動投稿でDiscordへ投稿される

## 問題

設定の「Discord投稿」をOFFにしていても、VRChat写真自動投稿やスクリーンショット自動投稿の経路でDiscord投稿される場合がある。

## 期待する挙動

「Discord投稿」がOFFの場合は、通常処理、自動投稿、スクリーンショット自動投稿のどの経路でもDiscordへ投稿しない。

## 受け入れ条件

- Discord投稿OFFの場合、自動投稿 watcher がDiscord投稿を開始しない。
- AutoPhotoWatcher内部で `UploadDiscord` を強制的にONへ変更しない。
- 既存configで自動投稿ONが残っていても、Discord投稿OFFならDiscordへ投稿しない。
- 回帰テストでDiscord投稿OFF時にWebhookエラーにならず、ローカル保存など設定通りに処理されることを確認する。

## 対応

- `AutoPhotoWatcher.process` で `UploadDiscord` を強制的にONへ変更しないようにした。
- `restartAutoPhotoWatcher` でDiscord投稿OFFの場合は自動投稿 watcher を起動しないようにした。
- Discord投稿OFF時に自動投稿処理がWebhookエラーにならず、ローカル保存として完了する回帰テストを追加した。
