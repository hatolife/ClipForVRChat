# EXIFへ同席ユーザー情報を書き込む方式を調査する

## 問題

v0.1.8の最小要件では、撮影時点で同じインスタンスにいたユーザー情報を画像と紐づけ、Discord投稿やEXIFなどに保持する必要がある。
sidecar JSONとDiscord本文は実装済みだが、EXIFへの同席ユーザー情報書き込み方式は未確定である。

## 期待する挙動

自動撮影で保存した画像へ、撮影方式、構図、同席ユーザー情報を、画像ビューアや後続処理で壊れにくいEXIF/XMP/PNGメタデータとして保持できる方式を調査し、v0.1.8で実装可能な方針と作業チケットへ分割する。

## 受け入れ条件

- [x] 現行の画像保存/sidecar JSON/Discord/履歴処理と `WriteEXIF` 設定の扱いを確認する。
- [x] PNG/JPEGそれぞれで現実的に書き込めるメタデータ形式を比較する。
- [x] 同席ユーザー情報に表示名とユーザーIDを含める場合のプライバシー/サイズ/互換性を整理する。
- [x] v0.1.8で採用する実装方針を決める。
- [x] 実装作業を細かい日本語issueへ分割する。
- [x] この調査では実装コード変更を行わない。

## 調査結果

- `AutoCaptureOutputConfig` には `WriteEXIF` と `WriteUserListToEXIF` があるが、現状は診断表示以外で参照されず、画像への埋め込み処理は未実装である。
- sidecar JSONは `finalizeAutoCaptureImage()` から `WriteAutoCaptureSidecar()` で作成され、同席ユーザー情報、構図、撮影方式、画像SHA256を保持している。
- Discord本文は `autoCaptureDiscordContent()` で同席ユーザーを表示できるが、履歴にはsidecar相当のユーザー一覧は取り込んでいない。
- Go標準の `image/png` と `image/jpeg` は任意のEXIF/XMP/PNGテキストチャンクを書き込む用途には足りないため、専用writerか依存追加が必要である。
- sidecar JSONは画像SHA256を保存するため、埋め込みメタデータを書き込む場合は「画像メタデータ書き込み -> sidecar JSON作成 -> Discord投稿」の順にする必要がある。

## 採用方針

sidecar JSONを完全な正本として維持し、画像埋め込みメタデータは「画像単体で最低限追跡できるコピー」として実装する。
`WriteEXIF` は画像埋め込みメタデータ全体の有効化、`WriteUserListToEXIF` は同席ユーザー一覧の埋め込み有効化として扱う。

JPEGはEXIF APP1の `UserComment` などにコンパクトJSONを入れる。
PNGはW3C PNG仕様上の `eXIf` チャンクへEXIFを入れ、互換性補助としてUTF-8 JSONを `iTXt` にも入れる方針とする。
ユーザーIDは永続的な識別子なので、表示名とは別制御し、既定では埋め込まない方向でUIと設定を整理する。

埋め込みpayloadにはサイズ上限を設け、上限を超えた場合はユーザー一覧を省略してsidecar参照と件数だけを残す。
メタデータ書き込みに失敗した場合は、画像保存自体を壊さず診断ログと結果メッセージへ理由を出す。

## 分割issue

- [#143](143-define-autocapture-embedded-metadata-schema.md) 自動撮影の埋め込みメタデータschemaを確定する。
- [#144](144-implement-image-metadata-writer.md) PNG/JPEGの埋め込みメタデータwriterを追加する。
- [#145](145-integrate-exif-presence-metadata.md) 自動撮影保存処理へEXIF/PNGメタデータ書き込みを統合する。
- [#146](146-expose-exif-privacy-settings.md) EXIF/埋め込みメタデータ設定とプライバシー説明をUIへ出す。
- [#147](147-verify-embedded-metadata-output.md) 埋め込みメタデータの読み戻し/Discord投稿後/実機検証を整備する。

## 参照

- PNG Specification Third Edition eXIf: https://www.w3.org/TR/png-3/#11eXIf
- PNG Specification Third Edition iTXt: https://www.w3.org/TR/png-3/#11iTXt
