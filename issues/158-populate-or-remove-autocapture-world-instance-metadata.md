# 自動撮影sidecarのworld/instance metadataを取得または削除する

## 問題

`AutoCaptureVRChatMetadata` には `world_id` と `instance_id` があるが、現行実装では常に空のままである。
sidecar JSONのschema上はVRChatワールド/インスタンス情報を保持できるように見えるが、実際には取得処理へ接続されていない。

## 期待する挙動

VRChat output logなどからworld ID / instance IDを取得できる場合はsidecar JSONへ保存する。
取得しない方針なら、空フィールドをschemaから外すか、未取得であることが分かる表現にする。

## 受け入れ条件

- output logからworld/instance移動ログを安全に解析できるか調査する。
- 取得できる場合はBatch開始時点のworld ID / instance IDをsidecar JSONへ保存する。
- 取得できない場合は、sidecar schemaとREADME/SPECで未取得として明示する。
- 空文字を「取得済みだが空」と誤解しないJSON表現にする。
- world/instance IDが個人情報や参加導線になり得る点をプライバシー説明に含める。
- ログ形式の揺れに対する単体テストを追加する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- `parseVRChatWorldMetadata()` の正規表現を、括弧を含むinstance IDを落とさないパーサへ変更する。
- `wrld_...` の開始位置を見つけ、空白、引用符、制御文字、明らかなログ区切りまでをworld location tokenとして読む方式にする。
- token末尾の句点、カンマ、閉じ括弧などログ文の句読点だけを安全にtrimする。ただしinstance ID内の `~region(us)`、`~private(usr_...)`、`~nonce(...)` の括弧は保持する。
- world IDとinstance IDは、取得できない場合は空にし、sidecarでは `omitempty` により未取得と区別する現行方針を維持する。
- world/instance IDは参加導線やプライバシー情報になり得るため、README/SPECのプライバシー説明へ追記する。

対象ファイル:

- `src/internal/appcore/autocapture.go`
- `src/internal/appcore/autocapture_test.go`
- `README.md`
- `src/SPEC.md`

小タスク:

- `extractVRChatWorldToken(line string) string` のような小関数へ切り出す。
- 複数行ログでは最後のworld tokenを採用する現行仕様を維持する。
- テストケースを追加する。
  - `wrld_x:12345~region(jp)`
  - `wrld_x:12345~private(usr_aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa)~region(us)`
  - `wrld_x:12345~hidden(usr_x)~nonce(abc)~region(eu)`
  - 行末に `.` や `)` があるログ文
- sidecar JSONの `vrchat.instance_id` に括弧込みの完全な値が入ることを確認する。

確認方法:

- 既存テストの期待値を `67890~region(us)` へ修正して通す。
- 実機ログでworld/instanceが取れない場合でも撮影自体は失敗しない。
