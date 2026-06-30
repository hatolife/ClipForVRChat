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
