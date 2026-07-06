# Discord投稿ONで通常投稿用Webhook URLが空欄の場合に保存時警告を出す

## 問題

Discord投稿をONにしても通常投稿用Webhook URLが空欄のまま保存できるため、実際に処理したときに投稿できずユーザーが原因に気づきにくい。

## 期待する挙動

Discord投稿がONで通常投稿用Webhook URLが空欄の状態で設定を保存した場合、画面上部にバナー形式の警告を表示し、通常投稿用Webhook URLの設定を促す。
設定されていなくても保存や画面遷移はできる。

## 受け入れ条件

- Discord投稿がON、かつ通常投稿用Webhook URLが空欄の状態で保存すると、更新通知と同じ位置にバナー形式の警告が表示される。
- 警告からDiscord投稿タブへ移動できる。
- Discord投稿がON、かつ通常投稿用Webhook URLが空欄の場合、通常投稿用Webhook URLの入力欄がオレンジ系の注意表示になる。
- 通常投稿用Webhook URLが空欄でも、保存して設定画面から戻れる。
- Discord投稿がOFFの場合は、通常投稿用Webhook URLが空欄でも警告しない。
- 通常投稿用Webhook URLが入力されている場合は警告しない。
