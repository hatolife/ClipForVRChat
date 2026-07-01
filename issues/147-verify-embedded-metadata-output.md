# 埋め込みメタデータの読み戻し/Discord投稿後/実機検証を整備する

## 問題

画像メタデータはビューアやアップロード先によって保持/削除の挙動が異なる。
書き込み実装だけでは、EXIF/PNGメタデータに同席ユーザー情報が保持されていることをユーザーが確認できない。

## 期待する挙動

PNG/JPEGそれぞれでメタデータが読み戻せること、sidecar JSONのSHA256と整合すること、Discord投稿後の保持/削除挙動が把握できることを確認できる。

## 受け入れ条件

- PNGとJPEGのサンプルに対し、埋め込みpayloadを自前readerで読み戻す自動テストを追加する。
- Windows実機で、保存画像、sidecar JSON、Discord投稿結果を使った確認手順を作る。
- `exiftool` など外部ツールがある場合の確認コマンドを手順に記載する。ただし外部ツール未導入でも自前readerで確認できるようにする。
- Discord投稿後にメタデータが保持されるか削除されるかを確認し、仕様として記録する。
- ユーザー数が多い場合、ユーザーIDを含む場合、表示名が日本語の場合の確認項目を入れる。
- 失敗時に必要な画像、sidecar、ログ、configの収集手順を記載する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- [#144](144-implement-image-metadata-writer.md) で追加する `ReadAutoCaptureEmbeddedMetadata()` を使い、PNG/JPEGの自動テストでJSONを構造体として読み戻す。
- Discord投稿後のmetadata保持はDiscord API上保証されていない。公式APIはファイル添付方法を定義しているが、EXIF/PNG metadata保持は仕様化していないため、実機確認手順では「Discord本文とsidecarを正本」「添付ファイルmetadataは観測結果として記録」と扱う。
- 外部ツール確認は任意で `exiftool` を使う。ただしCIは自前readerを正とする。

対象ファイル:

- `src/internal/appcore/metadata_test.go`
- `docs/v0.1.8-embedded-metadata-verification.md` 新設
- `README.md`
- `src/SPEC.md`

小タスク:

- PNG iTXt/eXIfの読み戻しテストを追加する。
- JPEG APP1 EXIFの読み戻しテストを追加する。
- ユーザーIDあり/なし、日本語表示名、ユーザー数過多のテストを追加する。
- sidecar JSONの `files.sha256` と実画像SHA256の一致テストを追加する。
- 実機手順にDiscord投稿後の確認、Discordからダウンロードしたファイルのreader確認、metadataが消えている場合の記録方法を書く。

確認方法:

- `go test ./...` でmetadata読み戻しが通る。
- Discord投稿あり/なしの両方で、ローカル画像とsidecarの紐づけが確認できる。
- metadataがDiscord側で削除された場合でも、投稿本文に同席ユーザー情報が残る。
