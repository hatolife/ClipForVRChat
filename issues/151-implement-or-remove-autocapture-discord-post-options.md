# 自動撮影DiscordのpostMode/includeImagesを実装または削除する

## 問題

`AutoCaptureDiscordConfig` には `postMode` と `includeImages` があり、正規化と診断出力の対象になっている。
しかし `finalizeAutoCaptureImage()` は、Discord自動投稿がONなら常に各Shotの画像を1枚ずつ添付投稿する。
Batch投稿や画像なし投稿は実装されていない。

## 期待する挙動

設定として残す場合は、`postMode=shot/batch` と `includeImages` が実際の投稿方式に反映される。
v0.1.8で扱わない場合は、未実装の設定値として保存/診断/UIへ出さない。

## 受け入れ条件

- `postMode=shot` は現行どおりShotごとに投稿される。
- `postMode=batch` を残す場合は、Batch単位の投稿本文と添付方針を実装する。
- `includeImages=false` を残す場合は、画像添付なしの本文投稿、またはsidecar/ローカル保存のみの扱いを実装する。
- 未対応値は正規化で実装済み値へ戻し、診断ログに理由を残す。
- Discord投稿の単体テストで、payloadと添付有無を確認する。
- UI/README/SPECが実際に対応する投稿方式だけを案内する。
