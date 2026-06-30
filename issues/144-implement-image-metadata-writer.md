# PNG/JPEGの埋め込みメタデータwriterを追加する

## 問題

Go標準の画像エンコード処理だけでは、自動撮影画像へEXIF/XMP/PNGテキストチャンクを非破壊に追加できない。
Photo方式ではVRChatが保存した画像、Stream方式ではSpout helper等が保存した画像に対して、後処理で安全にメタデータを挿入する必要がある。

## 期待する挙動

PNG/JPEG画像へ、定義済みschemaのメタデータを壊れない形で書き込み、書き込み後に読み戻しテストで確認できる。

## 受け入れ条件

- JPEGはEXIF APP1へコンパクトJSONまたは同等のUserCommentを挿入できる。
- PNGは `eXIf` チャンクへEXIFを挿入し、UTF-8互換性補助として `iTXt` に `ClipForVRChat:AutoCapture` JSONも挿入できる。
- 既存のPNG/JPEGファイルを再エンコードせず、チャンク/セグメント挿入で画素データを変えない。
- ファイル拡張子ではなくマジックバイトでPNG/JPEGを判定する。
- 既存EXIF/テキストチャンクがある場合の置換/追記方針を実装し、重複で壊れないようにする。
- 書き込み後の画像がGo標準デコーダで読めること、埋め込みpayloadを自前readerで読み戻せることをテストする。
- 未対応形式では画像保存を失敗させず、埋め込み未対応として呼び出し側へ返す。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。
