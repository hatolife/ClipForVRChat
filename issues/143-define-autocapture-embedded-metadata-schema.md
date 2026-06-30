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
