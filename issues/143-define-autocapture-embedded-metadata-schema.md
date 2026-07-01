# 自動撮影の埋め込みメタデータschemaを確定する

## 問題

sidecar JSONには自動撮影メタデータが保持されているが、画像単体に埋め込むEXIF/PNGメタデータのschemaがない。
表示名やVRChatユーザーIDは個人識別性があるため、単純にsidecar全体を埋め込むとサイズとプライバシーの問題が出る。

## 期待する挙動

画像へ埋め込む最小メタデータschemaを定義し、sidecar JSONを正本として参照しつつ、画像単体でも撮影方式、構図、同席ユーザー情報を追跡できる。

## 受け入れ条件

- schema version、アプリ名/バージョン、batch ID、shot ID、撮影時刻、撮影方式、構図ID/名前、同席ユーザー件数を定義する。
- `WriteUserListToEXIF` 有効時に入れるユーザー情報の項目を、表示名、ユーザーID、status、source、confidenceごとに定義する。
- ユーザーIDは表示名とは別制御にし、既存の `IncludeUserIDsInSidecar` / Discord設定と混同しない方針を決める。
- 埋め込みpayloadの最大サイズと、上限超過時の省略ルールを決める。
- sidecar JSONのSHA256計算順に影響するため、埋め込みmetadataがsidecar作成前に確定することを明記する。
- 画像単体で外部へ共有される場合の注意文を、設定UI issueから参照できる形で整理する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- schemaはsidecar JSONを正本とし、画像埋め込みは「画像単体で最低限追跡できる情報」と「設定で許可された同席ユーザー情報」に絞る。
- `AutoCaptureEmbeddedMetadata` に以下を追加する。
  - `app_version`
  - `user_count`
  - `users_truncated`
  - `resolved_pose`
  - `metadata_warnings`
- `Users` は `WriteUserListToEXIF=true` のときだけ入れる。ユーザーIDは `WriteUserIDsToEXIF=true` のときだけ入れる。
- JPEG APP1のサイズ上限を考慮し、埋め込みJSONの最大値は60000 bytes未満に抑える。超過時はユーザー一覧を件数だけに縮退し、`users_truncated=true` と警告を入れる。
- sidecarの画像SHA256は、埋め込み後の画像に対して計算する。埋め込みmetadata内にsidecar SHAは入れない。

対象ファイル:

- `src/internal/appcore/metadata.go`
- `src/internal/appcore/autocapture.go`
- `src/internal/appcore/metadata_test.go`
- `src/SPEC.md`

小タスク:

- `BuildAutoCaptureEmbeddedMetadata()` にapp version相当の入力を追加するか、現行では空/`dev` とする方針を決める。
- `user_count` は埋め込み対象ユーザー数ではなく撮影時に取得できた同席ユーザー総数を入れる。
- payload作成時にサイズを測り、ユーザー一覧省略処理を行う。
- `resolved_pose` は [#140](140-implement-player-local-pose-transform.md) の送信world poseを使う。
- schema versionを1のまま拡張するか、2へ上げるかを実装時に決め、README/SPECへ明記する。

確認方法:

- ユーザーID埋め込みOFFでIDが入らない。
- ユーザー数が多い場合に埋め込みが失敗せず、`users_truncated=true` になる。
- sidecarの画像SHA256が埋め込み後ファイルと一致する。
