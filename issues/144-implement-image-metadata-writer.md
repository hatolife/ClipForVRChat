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

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 進行中メモ

- 2026-07-01: PNG は `iTXt` だけでなく `eXIf` も書き込み、JPEG は APP1 を idempotent に置換する実装へ寄せる。
- 2026-07-01: 画像形式判定は拡張子ではなく magic bytes に切り替え、自前 reader `ReadAutoCaptureEmbeddedMetadata(path)` を追加する。
- 2026-07-01: `finalizeAutoCaptureImage()` 側は metadata 書き込み失敗を画像保存失敗にせず、warning として継続できるようにする。

## 実装前調査メモ

実装方針:

- ファイル形式判定は拡張子ではなくマジックバイトで行う。
- PNGは `eXIf` チャンクと `iTXt` チャンクの両方を書き込む。PNG拡張仕様では `eXIf` はJPEG APP1 marker/length/`Exif\0\0` を含まないExif profileを持ち、IDATチャンク列の途中以外、IHDRからIENDの間に1つだけ置ける。
- JPEGは現行のAPP1 EXIF挿入を維持しつつ、readerで読み戻せる構造にする。タグは現行の `ImageDescription` を継続するか、`UserComment` へ移行するかを実装時に決める。互換性を優先するなら当面は `ImageDescription` 継続でよい。
- 既存の自アプリmetadataがある場合は重複追記せず置換する。無関係なEXIF/PNGチャンクは残す。
- 未対応形式は `ErrUnsupportedEmbeddedMetadataFormat` のような識別可能エラーを返し、呼び出し側で警告扱いにできるようにする。
- 自前reader `ReadAutoCaptureEmbeddedMetadata(path)` を追加し、PNG/JPEGの埋め込みJSONを構造体として読み戻す。

対象ファイル:

- `src/internal/appcore/metadata.go`
- `src/internal/appcore/metadata_test.go`

小タスク:

- `detectImageMetadataFormat(data []byte)` を追加する。
- PNG chunk iteratorを実装し、CRC検証、`iTXt` keyword抽出、`eXIf` 抽出、同一keyword置換を行う。
- `jpegExifDescriptionSegment()` をTIFF payload生成とAPP1 segment生成に分離し、PNG `eXIf` でもTIFF payloadを再利用する。
- JPEG APP1 iteratorを実装し、既存の自アプリJSONを含むAPP1だけ置換する。
- `ReadAutoCaptureEmbeddedMetadata()` はPNGでは `iTXt` を優先し、なければ `eXIf`、JPEGではAPP1 EXIFを読む。

確認方法:

- PNG/JPEGとも、書き込み後の画像がGo標準decoderで読める。
- PNG/JPEGとも、自前readerでJSONを読み戻せる。
- 拡張子が `.bin` でもPNGマジックならPNGとして扱える。
- 2回書き込みしても自アプリmetadataが重複しない。
