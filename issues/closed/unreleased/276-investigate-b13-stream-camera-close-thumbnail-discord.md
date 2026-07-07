# b13でStream自動撮影後にカメラ終了・サムネイル未表示・Discord未投稿になる

## 指示

> b13でなにかの拍子にカメラが閉じられるようなので調査してほしい
> '/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b12-windows-amd64'
> logとか設定とか見て
>
> ---
> またstreamで自動撮影された画像のサムネイルが表示されません
> ---
> またstreamで自動撮影された画像がdiscordに投稿されません

## 文脈

v0.1.8-b13作成後、Stream方式の自動撮影周辺で、VRChatカメラが閉じられる、履歴サムネイルが出ない、Discordへ投稿されないという報告がある。
ユーザー指定の確認対象フォルダは `ClipForVRChat-v0.1.8-b12-windows-amd64` で、ログ、設定、出力画像を確認できる。

## 解釈

b13で顕在化した可能性を含め、指定フォルダの実行ログと設定から、Camera OSC復元/終了送信、Stream撮影結果の履歴反映、Discord投稿条件を切り分ける。

## 問題

- Stream自動撮影後または何らかのタイミングでVRChatカメラが閉じられる可能性がある。
- Stream自動撮影された画像のサムネイルがUIに表示されない。
- Stream自動撮影された画像がDiscordに投稿されない。

## 期待する挙動

- カメラ自動終了設定がOFFなら、自動撮影後に `/usercamera/Close` や `Streaming=false` を送ってカメラを閉じない。
- Stream自動撮影で保存された画像は履歴/結果にサムネイル付きで表示される。
- Discord投稿が有効でWebhookが設定されている場合、Stream自動撮影画像もDiscordへ投稿される。

## 受け入れ条件

- 指定フォルダの設定とログから、カメラ終了OSC送信の有無と条件が説明されている。
- サムネイル未表示の原因が、保存処理、履歴記録、frontend表示、画像パス解決のいずれかに切り分けられている。
- Discord未投稿の原因が、設定、Webhook選択、投稿処理未呼び出し、投稿エラーのいずれかに切り分けられている。
- 必要な修正がある場合は実装し、関連テストを追加または更新する。

## 調査結果

- 指定フォルダの `config.json` では `autoCapture.capture.openCameraBeforeBatch=false`、`closeCameraAfterBatch=false` だが、`autoCapture.restore.enabled=true` かつfallbackが `mode=0`、`streaming=false` だった。
- `logs/2026-07-07.log` では自動撮影後に `camera close skipped` が出ており `/usercamera/Close` は送っていない。一方で直後の `camera restore` が `/usercamera/Streaming=false` と `/usercamera/Mode=0` を送っていたため、実質的にStream Cameraを閉じる方向へ動いていた。
- Spout sender復旧処理も、sender一覧が空のときに `/usercamera/Streaming=false` を送ってからtrueへ戻していた。自動起動OFF時でもこの経路が動くため、撮影中にカメラが閉じたように見える可能性があった。
- `history.json` の23:43以降のStream自動撮影履歴は `sourcePath` / `outputPath` は入っているが `thumbnail`、`url`、`discordMessageId` が空だった。
- コード上、通常の画像処理は `Processor` がサムネイルを生成するが、自動撮影の `finalizeAutoCaptureImage` は保存画像からサムネイルを生成していなかった。
- `config.json` では通常Discord投稿 `output.uploadDiscord=true` だが、自動撮影専用の `autoCapture.discord.enabled=false` だった。実装は専用設定だけを見ていたため、通常Discord投稿ONでも自動撮影画像は投稿対象になっていなかった。

## 対応

- 自動起動/終了OFF時の撮影後復元では、撮影前にVRChatから受信できたMode/Streamingだけを戻し、未受信時のfallback `mode=0` / `streaming=false` ではカメラを閉じないようにした。
- Spout sender復旧処理は、自動起動OFF時に `Streaming=false` を送らず、`Streaming=true` の再送だけで復旧を試すようにした。
- 自動撮影保存画像からサムネイルdata URLを生成し、結果/履歴に表示できるようにした。
- 通常Discord投稿ONも自動撮影Discord投稿の有効条件に含め、現在の設定のままでも自動撮影画像を投稿対象にするようにした。
- 設定画面の自動撮影Discord文言、Webhook未設定警告、画像添付設定の有効条件を実装に合わせた。
- 関連するGoテストを追加した。
