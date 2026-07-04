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

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- `finalizeAutoCaptureImage()` は埋め込みmetadata書き込み失敗時も、画像保存、sidecar JSON、Discord投稿、履歴追加を可能な限り継続する。
- 現行の `Result.Error` に埋め込み失敗を入れると `AddResultsToHistory()` が履歴追加をスキップするため、警告は `Result.Warning` 追加、またはsidecarの `metadata_warnings` と診断ログに逃がす設計が必要。
- sidecarには埋め込み成功/失敗状態を記録する。例: `embedded_metadata: { attempted, written, error }`。
- Discord投稿は埋め込み後のローカル画像を添付する。ただしDiscord側が画像metadataを保持する保証はないため、Discord本文とsidecarを正本にする。

対象ファイル:

- `src/internal/appcore/autocapture.go`
- `src/internal/appcore/types.go`
- `src/internal/appcore/history.go`
- `src/internal/appcore/metadata.go`
- `src/internal/appcore/autocapture_test.go`

小タスク:

- `Result` / `HistoryEntry` に警告を持たせるか、sidecarだけに警告を持たせるかを決める。
- `WriteAutoCaptureEmbeddedMetadata()` の unsupported形式は警告扱いにし、PNG/JPEGの破損など実際の書き込み失敗もsidecar/ログへ残して継続する。
- sidecar SHA256は埋め込み後の画像に対して計算する。
- Photo方式でVRChat写真へ直接追記するため、診断ログに元ファイル改変であることを出す。

確認方法:

- 埋め込み成功時、sidecar SHA256が画像と一致する。
- 埋め込み失敗を人工的に起こしてもsidecar JSONと履歴が作られる。
- Discord投稿ONでも、埋め込み失敗が投稿全体の失敗にならない。
