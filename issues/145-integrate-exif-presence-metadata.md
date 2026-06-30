# 自動撮影保存処理へEXIF/PNGメタデータ書き込みを統合する

## 問題

`WriteEXIF` と `WriteUserListToEXIF` は設定値として存在するが、`finalizeAutoCaptureImage()` に接続されていない。
現在の処理順では、sidecar JSON作成後に画像へメタデータを書き込むとsidecar内SHA256が不一致になる。

## 期待する挙動

自動撮影で保存された画像に対し、設定に従って同席ユーザー情報を含む埋め込みメタデータを書き込み、その後にsidecar JSONとDiscord投稿を実行する。

## 受け入れ条件

- `finalizeAutoCaptureImage()` の処理順を「埋め込みメタデータ書き込み -> sidecar JSON作成 -> Discord投稿」にする。
- `WriteEXIF=false` の場合は現行挙動を変えない。
- `WriteEXIF=true` かつ `WriteUserListToEXIF=false` の場合は、ユーザー件数やshot IDなどの最小情報だけを埋め込む。
- `WriteUserListToEXIF=true` の場合は、schemaで許可された範囲のユーザー一覧を埋め込む。
- Photo方式でVRChat写真ファイルへ直接追記する場合の注意をログに残し、元ファイル改変を仕様として明記する。
- Stream/Spout方式の出力画像にも同じ処理を適用する。
- 埋め込み失敗時は画像保存、sidecar JSON、Discord投稿を可能な限り継続し、結果と診断ログへ理由を出す。
