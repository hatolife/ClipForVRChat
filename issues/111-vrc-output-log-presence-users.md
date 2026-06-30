# VRChat output logからの同席ユーザー保持

## 問題

自動撮影時点で同じインスタンスにいたユーザー情報を取得・保持できず、画像、Discord投稿、メタデータへ紐づけられない。

## 期待する挙動

VRChat output logを監視し、join/leaveやワールド移動を検出して現在インスタンスのユーザー集合を保持する。撮影Batch/Shot時点のスナップショットをsidecar JSONへ保存し、設定で許可された場合のみDiscord投稿本文へ表示名またはユーザーIDを含める。

## 受け入れ条件

- 最新の `output_log_*.txt` を監視対象として検出できる。
- join/leaveログから表示名と、含まれる場合は `usr_...` 形式のユーザーIDを抽出できる。
- ワールド/インスタンス移動らしきログで現在ユーザー集合をリセットできる。
- ログ形式が未知でもアプリが落ちず、ユーザー情報の信頼度を `partial` または `unknown` として保存できる。
- sidecar JSONには撮影時点のユーザー一覧が画像と紐づいて保存される。
- Discord投稿へのユーザー名/ID出力は既定OFFで、設定ON時だけ行われる。

## 実装メモ

- v0.1.8 RCでは最新の `output_log_*.txt` を撮影時に読み取り、現在ユーザー集合を推定する。
- ユーザーIDをsidecar JSONへ含めるか、Discord本文へ表示名/IDを含めるかは設定で制御する。
- output log形式の揺れがあるため、実機確認で追加パターンが見つかった場合はパーサを拡張する。

## 監査メモ

- 現行実装は、Batch開始時に最新 `output_log_*.txt` を先頭から読み直すsnapshot方式であり、アプリ起動中にtailし続ける監視方式ではない。
- sidecar JSONの `world_id` / `instance_id` は型だけ存在し、現行実装では値を入れていない。
- `IncludeUserIDsInSidecar=false` の場合に、Discord本文用のユーザー一覧からもユーザーIDが消えるため、Discord用ユーザーID設定がsidecar設定から独立していない。
- v0.1.8の最小要件を満たすには、snapshot方式を正式な制約としてUI/仕様へ明記するか、継続監視へ実装を進める必要がある。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。
